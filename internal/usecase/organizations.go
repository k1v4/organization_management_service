package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/k1v4/organization_management_service/internal/entity"
)

type IOrganizationRepository interface {
	Create(ctx context.Context, org *entity.Organization) (*entity.Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error)
	Update(ctx context.Context, org *entity.Organization) (*entity.Organization, error)
	Archive(ctx context.Context, id uuid.UUID) error
	UpdateOwner(ctx context.Context, id uuid.UUID, ownerIdentityID string) error
}

type OrganizationUseCase struct {
	repo IOrganizationRepository
}

func NewOrganizationUseCase(repo IOrganizationRepository) *OrganizationUseCase {
	return &OrganizationUseCase{repo: repo}
}

func (uc *OrganizationUseCase) CreateOrganization(ctx context.Context, org *entity.Organization) (*entity.Organization, error) {
	return uc.repo.Create(ctx, org)
}

func (uc *OrganizationUseCase) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *OrganizationUseCase) UpdateOrganization(ctx context.Context, org *entity.Organization) (*entity.Organization, error) {
	return uc.repo.Update(ctx, org)
}

func (uc *OrganizationUseCase) ArchiveOrganization(ctx context.Context, id uuid.UUID) error {
	return uc.repo.Archive(ctx, id)
}

func (uc *OrganizationUseCase) UpdateOrganizationOwner(ctx context.Context, id uuid.UUID, ownerIdentityID string) error {
	return uc.repo.UpdateOwner(ctx, id, ownerIdentityID)
}
