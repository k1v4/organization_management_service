package v1

import (
	"context"

	"github.com/google/uuid"
	"github.com/k1v4/organization_management_service/internal/entity"
	"github.com/k1v4/organization_management_service/pkg/logger"
	"github.com/labstack/echo/v4"
)

type IOrganizationService interface {
	CreateOrganization(ctx context.Context, org *entity.Organization) (*entity.Organization, error)
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error)
	UpdateOrganization(ctx context.Context, org *entity.Organization) (*entity.Organization, error)
	ArchiveOrganization(ctx context.Context, id uuid.UUID) error
	UpdateOrganizationOwner(ctx context.Context, id uuid.UUID, ownerIdentityID string) error
}

type organizationRoutes struct {
	t IOrganizationService
	l logger.Logger
}

func newOrganizationRoutes(handler *echo.Group, t IOrganizationService, l logger.Logger) {
	r := &organizationRoutes{t, l}

	// GET /api/v1/articles/{id}
	// handler.GET("/articles/:id", r.GetArticle)

	_ = r
}
