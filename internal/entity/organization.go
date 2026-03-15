package entity

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID              uuid.UUID
	Name            string
	Description     *string
	Status          string
	OwnerIdentityID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
