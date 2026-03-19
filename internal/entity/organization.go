package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Description     *string   `json:"description"`
	Status          string    `json:"status"`
	OwnerIdentityID string    `json:"owner_identity_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (o *Organization) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Organization) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
