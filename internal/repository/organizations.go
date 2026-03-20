package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/k1v4/organization_management_service/internal/entity"
	"github.com/k1v4/organization_management_service/pkg/database"
	"github.com/k1v4/organization_management_service/pkg/database/postgres"
)

type OrganizationRepository struct {
	*postgres.Postgres
}

func NewOrganizationRepository(pg *postgres.Postgres) *OrganizationRepository {
	return &OrganizationRepository{pg}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *entity.Organization) (*entity.Organization, error) {
	sql, args, err := r.Builder.
		Insert("organizations").
		Columns("name", "description", "owner_identity_id").
		Values(org.Name, org.Description, org.OwnerIdentityID).
		Suffix("RETURNING id, name, description, status, owner_identity_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("OrganizationRepository.Create - ToSql: %w", err)
	}

	created := &entity.Organization{}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.Status,
		&created.OwnerIdentityID,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("OrganizationRepository.Create - QueryRow: %w", err)
	}

	return created, nil
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error) {
	sql, args, err := r.Builder.
		Select("id", "name", "description", "status", "owner_identity_id", "created_at", "updated_at").
		From("organizations").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("OrganizationRepository.GetByID - ToSql: %w", err)
	}

	org := &entity.Organization{}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&org.ID,
		&org.Name,
		&org.Description,
		&org.Status,
		&org.OwnerIdentityID,
		&org.CreatedAt,
		&org.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrOrganizationNotFound
		}

		return nil, fmt.Errorf("OrganizationRepository.GetByID - QueryRow: %w", err)
	}

	return org, nil
}

func (r *OrganizationRepository) Update(ctx context.Context, org *entity.Organization) (*entity.Organization, error) {
	sql, args, err := r.Builder.
		Update("organizations").
		Set("name", org.Name).
		Set("description", org.Description).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": org.ID, "status": "active"}).
		Suffix("RETURNING id, name, description, status, owner_identity_id, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("OrganizationRepository.Update - ToSql: %w", err)
	}

	updated := &entity.Organization{}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.Status,
		&updated.OwnerIdentityID,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrOrganizationNotFound
		}

		return nil, fmt.Errorf("OrganizationRepository.Update - QueryRow: %w", err)
	}

	return updated, nil
}

func (r *OrganizationRepository) Archive(ctx context.Context, id uuid.UUID) error {
	sql, args, err := r.Builder.
		Update("organizations").
		Set("status", "archived").
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": id, "status": "active"}).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrganizationRepository.Archive - ToSql: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.ErrOrganizationNotFound
		}

		return fmt.Errorf("OrganizationRepository.Archive - Exec: %w", err)
	}

	return nil
}

func (r *OrganizationRepository) UpdateOwner(ctx context.Context, id uuid.UUID, ownerIdentityID, newOwnerIdentityID string) error {
	sql, args, err := r.Builder.
		Update("organizations").
		Set("owner_identity_id", ownerIdentityID).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"id": id, "status": "active"}).
		ToSql()
	if err != nil {
		return fmt.Errorf("OrganizationRepository.UpdateOwner - ToSql: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return database.ErrOrganizationNotFound
		}

		return fmt.Errorf("OrganizationRepository.UpdateOwner - Exec: %w", err)
	}

	return nil
}
