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

type RulesRepository struct {
	*postgres.Postgres
}

func NewRulesRepository(pg *postgres.Postgres) *RulesRepository {
	return &RulesRepository{pg}
}

func (r *RulesRepository) GetByOrganizationID(ctx context.Context, orgID uuid.UUID) (*entity.BookingPolicy, error) {
	sql, args, err := r.Builder.
		Select("id", "organization_id", "max_booking_duration_min", "booking_window_days", "max_active_bookings_per_user", "created_at", "updated_at").
		From("booking_policies").
		Where(squirrel.Eq{"organization_id": orgID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("RulesRepository.GetByOrganizationID - ToSql: %w", err)
	}

	policy := &entity.BookingPolicy{}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&policy.ID,
		&policy.OrganizationID,
		&policy.MaxBookingDurationMin,
		&policy.BookingWindowDays,
		&policy.MaxActiveBookingsPerUser,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrRulesNotFound
		}

		return nil, fmt.Errorf("RulesRepository.GetByOrganizationID - QueryRow: %w", err)
	}

	return policy, nil
}

func (r *RulesRepository) Update(ctx context.Context, policy *entity.BookingPolicy) (*entity.BookingPolicy, error) {
	sql, args, err := r.Builder.
		Update("booking_policies").
		Set("max_booking_duration_min", policy.MaxBookingDurationMin).
		Set("booking_window_days", policy.BookingWindowDays).
		Set("max_active_bookings_per_user", policy.MaxActiveBookingsPerUser).
		Set("updated_at", squirrel.Expr("now()")).
		Where(squirrel.Eq{"organization_id": policy.OrganizationID}).
		Suffix("RETURNING id, organization_id, max_booking_duration_min, booking_window_days, max_active_bookings_per_user, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("RulesRepository.Update - ToSql: %w", err)
	}

	updated := &entity.BookingPolicy{}
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(
		&updated.ID,
		&updated.OrganizationID,
		&updated.MaxBookingDurationMin,
		&updated.BookingWindowDays,
		&updated.MaxActiveBookingsPerUser,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrRulesNotFound
		}

		return nil, fmt.Errorf("RulesRepository.Update - QueryRow: %w", err)
	}

	return updated, nil
}
