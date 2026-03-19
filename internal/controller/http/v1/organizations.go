package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/k1v4/organization_management_service/internal/entity"
	"github.com/k1v4/organization_management_service/pkg/jwtpkg"
	"github.com/k1v4/organization_management_service/pkg/logger"
	"github.com/labstack/echo/v4"
)

type IOrganizationService interface {
	CreateOrganization(ctx context.Context, org *entity.Organization) (*entity.Organization, error)
	GetOrganizationByID(ctx context.Context, organizationID, userID string) (*entity.Organization, error)
	UpdateOrganization(ctx context.Context, org *entity.Organization) (*entity.Organization, error)
	ArchiveOrganization(ctx context.Context, id uuid.UUID) error
	UpdateOrganizationOwner(ctx context.Context, id uuid.UUID, ownerIdentityID string) error
}

type organizationRoutes struct {
	t  IOrganizationService
	l  logger.Logger
	tv jwtpkg.TokenVerifier
}

func newOrganizationRoutes(handler *echo.Group, t IOrganizationService, l logger.Logger, tv jwtpkg.TokenVerifier) {
	r := &organizationRoutes{t, l, tv}

	// GET /api/v1/articles/{id}
	// handler.GET("/articles/:id", r.GetArticle)

	// POST /api/v1/organizations

	// GET /api/v1/organizations/{orgId}
	handler.GET("/organizations/:orgId", r.GetOrganization)

	// PUT /api/v1/organizations/{orgId}

	// DELETE /api/v1/organizations/{orgId}

	// PUT /api/v1/organizations/{orgId}/owner

	// GET /api/v1/organizations/{orgId}/policy
	// PUT /api/v1/organizations/{orgId}/policy

	_ = r
}

func (o *organizationRoutes) GetOrganization(c echo.Context) error {
	const op = "Controller.GetOrganization"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		err := errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	// получаем user id
	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	organizationID := c.Param("orgId")

	if len(strings.TrimSpace(organizationID)) == 0 {
		err = errorResponse(c, http.StatusBadRequest, "bad request")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}

		return fmt.Errorf("%s: %s", op, "item name is required")
	}

	organization, err := o.t.GetOrganizationByID(ctx, organizationID, userID)
	if err != nil {
		err = errorResponse(c, http.StatusInternalServerError, "internal error")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, organization)
}
