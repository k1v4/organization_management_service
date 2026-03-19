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
	ArchiveOrganization(ctx context.Context, id uuid.UUID, userID string) error
	UpdateOrganizationOwner(ctx context.Context, id uuid.UUID, ownerIdentityID string) error
}

type organizationRoutes struct {
	t  IOrganizationService
	l  logger.Logger
	tv jwtpkg.TokenVerifier
}

func newOrganizationRoutes(handler *echo.Group, t IOrganizationService, l logger.Logger, tv jwtpkg.TokenVerifier) {
	r := &organizationRoutes{t, l, tv}

	// POST /api/v1/organizations
	handler.GET("/organizations", r.CreateOrganization)

	// GET /api/v1/organizations/{orgId}
	handler.GET("/organizations/:orgId", r.GetOrganization)

	// PUT /api/v1/organizations/{orgId}
	handler.PUT("/organizations/:orgId", r.UpdateOrganization)

	// DELETE /api/v1/organizations/{orgId}
	handler.DELETE("/organizations/:orgId", r.ArchiveOrganization)

	// PUT /api/v1/organizations/{orgId}/owner
	handler.PUT("/organizations/:orgId/owner", r.UpdateOrganizationOwner)
}

func (o *organizationRoutes) CreateOrganization(c echo.Context) error {
	const op = "Controller.CreateOrganization"

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
		err = errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}

		return fmt.Errorf("%s: %s", op, err)
	}

	u := new(entity.PostOrganization)
	if err = c.Bind(u); err != nil {
		err = errorResponse(c, http.StatusBadRequest, "bad request")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if u.Name == "" {
		err = errorResponse(c, http.StatusBadRequest, "name is required")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}

		return fmt.Errorf("%s: %s", op, "name is required")
	}

	org, err := o.t.CreateOrganization(ctx, &entity.Organization{
		Name:            u.Name,
		Description:     u.Description,
		OwnerIdentityID: userID,
	})
	if err != nil {
		err = errorResponse(c, http.StatusInternalServerError, "internal error")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusCreated, org)
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
		err = errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}

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

func (o *organizationRoutes) UpdateOrganization(c echo.Context) error {
	const op = "Controller.UpdateOrganization"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		err := errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		err = errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, err)
	}

	organizationID := c.Param("orgId")
	if len(strings.TrimSpace(organizationID)) == 0 {
		err = errorResponse(c, http.StatusBadRequest, "bad request")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "organizationID is required")
	}

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		err = errorResponse(c, http.StatusBadRequest, "invalid organization id")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "invalid organization id")
	}

	u := new(entity.PostOrganization)
	if err = c.Bind(u); err != nil {
		err = errorResponse(c, http.StatusBadRequest, "bad request")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	org, err := o.t.UpdateOrganization(ctx, &entity.Organization{
		ID:              orgUUID,
		Name:            u.Name,
		Description:     u.Description,
		OwnerIdentityID: userID,
	})
	if err != nil {
		err = errorResponse(c, http.StatusInternalServerError, "internal error")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, org)
}

func (o *organizationRoutes) ArchiveOrganization(c echo.Context) error {
	const op = "Controller.ArchiveOrganization"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		err := errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		err = errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, err)
	}

	organizationID := c.Param("orgId")
	if len(strings.TrimSpace(organizationID)) == 0 {
		err = errorResponse(c, http.StatusBadRequest, "bad request")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "organizationID is required")
	}

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		err = errorResponse(c, http.StatusBadRequest, "invalid organization id")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "invalid organization id")
	}

	if err = o.t.ArchiveOrganization(ctx, orgUUID, userID); err != nil {
		err = errorResponse(c, http.StatusInternalServerError, "internal error")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (o *organizationRoutes) UpdateOrganizationOwner(c echo.Context) error {
	const op = "Controller.UpdateOrganizationOwner"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		err := errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		err = errorResponse(c, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, err)
	}

	organizationID := c.Param("orgId")
	if len(strings.TrimSpace(organizationID)) == 0 {
		err = errorResponse(c, http.StatusBadRequest, "bad request")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "organizationID is required")
	}

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		err = errorResponse(c, http.StatusBadRequest, "invalid organization id")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "invalid organization id")
	}

	var req struct {
		NewOwnerIdentityID string `json:"new_owner_identity_id"`
	}
	if err = c.Bind(&req); err != nil {
		err = errorResponse(c, http.StatusBadRequest, "bad request")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(strings.TrimSpace(req.NewOwnerIdentityID)) == 0 {
		err = errorResponse(c, http.StatusBadRequest, "new_owner_identity_id is required")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "new_owner_identity_id is required")
	}

	err = o.t.UpdateOrganizationOwner(ctx, orgUUID, req.NewOwnerIdentityID)
	if err != nil {
		err = errorResponse(c, http.StatusInternalServerError, "internal error")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, "")
}
