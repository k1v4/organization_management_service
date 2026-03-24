package v1

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/AlekSi/pointer"
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
	UpdateOrganizationOwner(ctx context.Context, id uuid.UUID, initiatorIdentityID, newOwnerIdentityID string) error
	SelectByStatus(ctx context.Context, ownerIdentityID string, status *entity.OrganizationStatus) ([]*entity.Organization, error)
}

type organizationRoutes struct {
	t  IOrganizationService
	l  logger.Logger
	tv *jwtpkg.TokenVerifier
}

func newOrganizationRoutes(handler *echo.Group, t IOrganizationService, l logger.Logger, tv *jwtpkg.TokenVerifier) {
	r := &organizationRoutes{t, l, tv}

	// POST /api/v1/organizations
	handler.POST("/organizations", r.CreateOrganization)

	// GET /api/v1/organizations/{orgId}
	handler.GET("/organizations/:orgId", r.GetOrganization)

	// PUT /api/v1/organizations/{orgId}
	handler.PUT("/organizations/:orgId", r.UpdateOrganization)

	// DELETE /api/v1/organizations/{orgId}
	handler.DELETE("/organizations/:orgId", r.ArchiveOrganization)

	// PUT /api/v1/organizations/{orgId}/owner
	handler.PUT("/organizations/:orgId/owner", r.UpdateOrganizationOwner)

	// PUT /api/v1/organizations/active
	handler.GET("/organizations/active", r.GetActiveOrganizations)

	// PUT /api/v1/organizations/deactivated
	handler.GET("/organizations/deactivated", r.GetDeletedOrganizations)

	// PUT /api/v1/organizations
	handler.GET("/organizations", r.GetAllOrganizations)
}

func (o *organizationRoutes) CreateOrganization(c echo.Context) error {
	const op = "Controller.CreateOrganization"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	// получаем user id
	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	u := new(entity.PostOrganization)
	if err = c.Bind(u); err != nil {
		fmt.Println(err)
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	if u.Name == "" {
		errorResponse(c, http.StatusBadRequest, "name is required")

		return fmt.Errorf("%s: %s", op, "name is required")
	}

	org, err := o.t.CreateOrganization(ctx, &entity.Organization{
		Name:            u.Name,
		Description:     u.Description,
		OwnerIdentityID: userID,
	})
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusCreated, org)
}

func (o *organizationRoutes) GetOrganization(c echo.Context) error {
	const op = "Controller.GetOrganization"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

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
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "item name is required")
	}

	organization, err := o.t.GetOrganizationByID(ctx, organizationID, userID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, organization)
}

func (o *organizationRoutes) UpdateOrganization(c echo.Context) error {
	const op = "Controller.UpdateOrganization"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	organizationID := c.Param("orgId")
	if len(strings.TrimSpace(organizationID)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "organizationID is required")
	}

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid organization id")

		return fmt.Errorf("%s: %s", op, "invalid organization id")
	}

	u := new(entity.PostOrganization)
	if err = c.Bind(u); err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	org, err := o.t.UpdateOrganization(ctx, &entity.Organization{
		ID:              orgUUID,
		Name:            u.Name,
		Description:     u.Description,
		OwnerIdentityID: userID,
	})
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, org)
}

func (o *organizationRoutes) ArchiveOrganization(c echo.Context) error {
	const op = "Controller.ArchiveOrganization"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	organizationID := c.Param("orgId")
	if len(strings.TrimSpace(organizationID)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "organizationID is required")
	}

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid organization id")

		return fmt.Errorf("%s: %s", op, "invalid organization id")
	}

	if err = o.t.ArchiveOrganization(ctx, orgUUID, userID); err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (o *organizationRoutes) UpdateOrganizationOwner(c echo.Context) error {
	const op = "Controller.UpdateOrganizationOwner"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	organizationID := c.Param("orgId")
	if len(strings.TrimSpace(organizationID)) == 0 {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %s", op, "organizationID is required")
	}

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid organization id")

		return fmt.Errorf("%s: %s", op, "invalid organization id")
	}

	var req struct {
		NewOwnerIdentityID string `json:"new_owner_identity_id"`
	}
	if err = c.Bind(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "bad request")

		return fmt.Errorf("%s: %w", op, err)
	}

	if strings.EqualFold(req.NewOwnerIdentityID, userID) {
		errorResponse(c, http.StatusBadRequest, "new_owner_identity_id is equal to current owner_id")

		return fmt.Errorf("%s: %s", op, "new_owner_identity_id is equal to current owner_id")
	}

	if len(strings.TrimSpace(req.NewOwnerIdentityID)) == 0 {
		errorResponse(c, http.StatusBadRequest, "new_owner_identity_id is required")

		return fmt.Errorf("%s: %s", op, "new_owner_identity_id is required")
	}

	err = o.t.UpdateOrganizationOwner(ctx, orgUUID, userID, req.NewOwnerIdentityID)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"status": "ok",
	})
}

func (o *organizationRoutes) GetActiveOrganizations(c echo.Context) error {
	const op = "Controller.GetActiveOrganizations"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	orgs, err := o.t.SelectByStatus(ctx, userID, pointer.To(entity.OrganizationStatusActive))
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, orgs)
}

func (o *organizationRoutes) GetDeletedOrganizations(c echo.Context) error {
	const op = "Controller.GetDeletedOrganizations"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	orgs, err := o.t.SelectByStatus(ctx, userID, pointer.To(entity.OrganizationStatusArchive))
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, orgs)
}

func (o *organizationRoutes) GetAllOrganizations(c echo.Context) error {
	const op = "Controller.GetDeletedOrganizations"

	ctx := c.Request().Context()

	token := jwtpkg.ExtractToken(c)
	if token == "" {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, "token is required")
	}

	userID, err := o.tv.GetIdentityID(ctx, token)
	if err != nil {
		errorResponse(c, http.StatusUnauthorized, "Unauthorized")

		return fmt.Errorf("%s: %s", op, err)
	}

	orgs, err := o.t.SelectByStatus(ctx, userID, nil)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "internal error")

		return fmt.Errorf("%s: %s", op, err)
	}

	return c.JSON(http.StatusOK, orgs)
}
