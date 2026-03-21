package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/k1v4/organization_management_service/pkg/adapter"
	"golang.org/x/sync/errgroup"

	"github.com/k1v4/organization_management_service/internal/entity"
)

type IOrganizationRepository interface {
	Create(ctx context.Context, org *entity.Organization) (*entity.Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error)
	Update(ctx context.Context, org *entity.Organization) (*entity.Organization, error)
	Archive(ctx context.Context, id uuid.UUID) error
	UpdateOwner(ctx context.Context, id uuid.UUID, ownerIdentityID, newOwnerIdentityID string) error
}

type OrganizationUseCase struct {
	repo    IOrganizationRepository
	adapter *adapter.Client
}

func NewOrganizationUseCase(repo IOrganizationRepository, adapter *adapter.Client) *OrganizationUseCase {
	return &OrganizationUseCase{
		repo:    repo,
		adapter: adapter,
	}
}

func (uc *OrganizationUseCase) CreateOrganization(ctx context.Context, org *entity.Organization) (*entity.Organization, error) {
	id, err := uc.adapter.GetUserByIdentityID(ctx, org.OwnerIdentityID)
	if err != nil {
		return nil, fmt.Errorf("create organization: %w", err)
	}
	if id == nil {
		return nil, fmt.Errorf("failed to create organization: owner not exist")
	}

	organization, err := uc.repo.Create(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	err = uc.adapter.SetOrganizationOwner(ctx, organization.ID.String(), organization.OwnerIdentityID)
	if err != nil {
		return nil, fmt.Errorf("failed to set organization owner: %w", err)
	}

	return organization, nil
}

func (uc *OrganizationUseCase) GetOrganizationByID(ctx context.Context, organizationID, userID string) (*entity.Organization, error) {
	permission, err := uc.adapter.CheckPermission(ctx, userID, organizationID, "ORG_READ")
	if err != nil {
		return nil, fmt.Errorf("UseCase-GetOrganizationByID: permission denied: %v", err)
	}
	if !permission {
		return nil, fmt.Errorf("UseCase-GetOrganizationByID: no access to organization")
	}

	organizationUUID, err := uuid.Parse(organizationID)
	if err != nil {
		return nil, fmt.Errorf("UseCase-GetOrganizationByID: %s - %s", "failed to parse organizationID into uuid", organizationID)
	}

	organization, err := uc.repo.GetByID(ctx, organizationUUID)
	if err != nil {
		return nil, fmt.Errorf("UseCase-GetOrganizationByID: %s - %s", "failed to get organization", organizationID)
	}

	return organization, nil
}

func (uc *OrganizationUseCase) UpdateOrganization(ctx context.Context, org *entity.Organization) (*entity.Organization, error) {
	permission, err := uc.adapter.CheckPermission(ctx, org.OwnerIdentityID, org.ID.String(), "ORG_UPDATE")
	if err != nil {
		return nil, fmt.Errorf("UseCase-GetOrganizationByID: permission denied: %v", err)
	}
	if !permission {
		return nil, fmt.Errorf("UseCase-GetOrganizationByID: no access to organization")
	}

	organisation, err := uc.repo.Update(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return organisation, nil
}

func (uc *OrganizationUseCase) ArchiveOrganization(ctx context.Context, orgID uuid.UUID, userID string) error {
	permission, err := uc.adapter.CheckPermission(ctx, userID, orgID.String(), "ORG_DELETE")
	if err != nil {
		return fmt.Errorf("UseCase-ArchiveOrganization: permission denied: %v", err)
	}
	if !permission {
		return fmt.Errorf("UseCase-ArchiveOrganization: no access to organization")
	}

	err = uc.repo.Archive(ctx, orgID)
	if err != nil {
		return fmt.Errorf("UseCase-ArchiveOrganization: %w", err)
	}

	return nil
}

func (uc *OrganizationUseCase) UpdateOrganizationOwner(ctx context.Context, id uuid.UUID, initiatorIdentityID, newOwnerIdentityID, token string) error {
	var (
		permission bool
		identityID *adapter.UserProfile
	)

	// параллельно проверяем права инициатора и существование нового владельца
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		permission, err = uc.adapter.CheckPermission(gCtx, initiatorIdentityID, id.String(), "ORG_OWNER_CHANGE")
		if err != nil {
			return fmt.Errorf("UseCase-UpdateOrganizationOwner: failed to check permission: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		var err error
		identityID, err = uc.adapter.GetUserByIdentityID(gCtx, newOwnerIdentityID)
		if err != nil {
			return fmt.Errorf("UseCase-UpdateOrganizationOwner: failed to get new owner: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if !permission {
		return fmt.Errorf("UseCase-UpdateOrganizationOwner: no access to organization")
	}
	if identityID == nil {
		return fmt.Errorf("UseCase-UpdateOrganizationOwner: new owner not found")
	}

	// получаем текущего владельца из БД
	org, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("UseCase-UpdateOrganizationOwner: failed to get organization: %w", err)
	}

	// получаем membershipId текущего владельца
	membership, err := uc.adapter.GetMembership(ctx, id.String(), org.OwnerIdentityID)
	if err != nil {
		return fmt.Errorf("UseCase-UpdateOrganizationOwner: failed to get old owner membership: %w", err)
	}

	// снимаем роль ORG_OWNER у старого владельца
	if err = uc.adapter.RevokeRole(ctx, id.String(), membership.MembershipID, "ORG_OWNER", token); err != nil {
		return fmt.Errorf("UseCase-UpdateOrganizationOwner: failed to revoke old owner role: %w", err)
	}

	// назначаем нового владельца в OrgMembershipService
	if err = uc.adapter.SetOrganizationOwner(ctx, id.String(), newOwnerIdentityID); err != nil {
		return fmt.Errorf("UseCase-UpdateOrganizationOwner: failed to set new owner: %w", err)
	}

	// обновляем owner_identity_id в БД
	if err = uc.repo.UpdateOwner(ctx, id, org.OwnerIdentityID, newOwnerIdentityID); err != nil {
		return fmt.Errorf("UseCase-UpdateOrganizationOwner: failed to update owner in db: %w", err)
	}

	return nil
}
