package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OrganizationStatus string

var OrganizationStatusActive OrganizationStatus = "active"
var OrganizationStatusArchive OrganizationStatus = "archived"

type Organization struct {
	ID              uuid.UUID          `json:"id"`
	Name            string             `json:"name"`
	Description     *string            `json:"description"`
	Status          OrganizationStatus `json:"status"`
	OwnerIdentityID string             `json:"owner_identity_id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type PostOrganization struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

func (o *Organization) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Organization) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}
