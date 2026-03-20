package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/k1v4/organization_management_service/internal/entity"
	"github.com/k1v4/organization_management_service/pkg/adapter"
)

type IRuleRepository interface {
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) (*entity.BookingPolicy, error)
	Create(ctx context.Context, orgID uuid.UUID) (*entity.BookingPolicy, error)
	Update(ctx context.Context, policy *entity.BookingPolicy) (*entity.BookingPolicy, error)
}

type RuleUseCase struct {
	repo    IRuleRepository
	adapter adapter.Client
}

func NewRuleUseCase(repo IRuleRepository) *RuleUseCase {
	return &RuleUseCase{repo: repo}
}

func (r *RuleUseCase) GetByOrganizationID(ctx context.Context, orgID uuid.UUID, userID string) (*entity.BookingPolicy, error) {
	permission, err := r.adapter.CheckPermission(ctx, userID, orgID.String(), "POLICIES_LIST")
	if err != nil {
		return nil, fmt.Errorf("CheckPermission: %w", err)
	}
	if !permission {
		return nil, fmt.Errorf("user doesnt have permission")
	}

	rules, err := r.repo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("GetByOrganizationID: %w", err)
	}

	return rules, nil
}

func (r *RuleUseCase) CreateRule(ctx context.Context, orgID uuid.UUID, userID string) (*entity.BookingPolicy, error) {
	create, err := r.repo.Create(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("CreateRule: %w", err)
	}

	return create, nil
}

func (r *RuleUseCase) UpdateRule(ctx context.Context, policy *entity.BookingPolicy, userID string) (*entity.BookingPolicy, error) {
	permission, err := r.adapter.CheckPermission(ctx, userID, policy.OrganizationID.String(), "POLICIES_MANAGE")
	if err != nil {
		return nil, fmt.Errorf("CheckPermission: %w", err)
	}
	if !permission {
		return nil, fmt.Errorf("user doesnt have permission")
	}

	update, err := r.repo.Update(ctx, policy)
	if err != nil {
		return nil, fmt.Errorf("UpdateRule: %w", err)
	}

	return update, nil
}
