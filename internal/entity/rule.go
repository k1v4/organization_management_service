package entity

import (
	"time"

	"github.com/google/uuid"
)

type BookingPolicy struct {
	ID                       int
	OrganizationID           uuid.UUID
	MaxBookingDurationMin    int
	BookingWindowDays        int
	MaxActiveBookingsPerUser int
	CreatedAt                time.Time
	UpdatedAt                time.Time
}
