package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/k1v4/organization_management_service/internal/entity"
)

type IRuleRepository interface {
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID) (*entity.BookingPolicy, error)
	Create(ctx context.Context, orgID uuid.UUID) (*entity.BookingPolicy, error)
	Update(ctx context.Context, policy *entity.BookingPolicy) (*entity.BookingPolicy, error)
}

type RuleUseCase struct {
	repo IRuleRepository
}

func NewRuleUseCase(repo IRuleRepository) *RuleUseCase {
	return &RuleUseCase{repo: repo}
}

func (r *RuleUseCase) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) (*entity.BookingPolicy, error) {
	id, err := r.repo.GetByOrganizationID(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("GetByOrganizationID: %w", err)
	}

	return id, nil
}

func (r *RuleUseCase) CreateRule(ctx context.Context, orgID uuid.UUID) (*entity.BookingPolicy, error) {
	create, err := r.repo.Create(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("CreateRule: %w", err)
	}

	return create, nil
}

func (r *RuleUseCase) UpdateRule(ctx context.Context, policy *entity.BookingPolicy) (*entity.BookingPolicy, error) {
	// TODO можно добавить проверку, что такая компания и вправду есть

	update, err := r.repo.Update(ctx, policy)
	if err != nil {
		return nil, fmt.Errorf("UpdateRule: %w", err)
	}

	return update, nil
}
