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

type IRulesService interface {
	UpdateRule(ctx context.Context, policy *entity.BookingPolicy, userID string) (*entity.BookingPolicy, error)
	GetByOrganizationID(ctx context.Context, orgID uuid.UUID, userID string) (*entity.BookingPolicy, error)
}

type rulesRoutes struct {
	t  IRulesService
	l  logger.Logger
	tv *jwtpkg.TokenVerifier
}

func newRulesRoutes(handler *echo.Group, t IRulesService, l logger.Logger, tv *jwtpkg.TokenVerifier) {
	r := &rulesRoutes{t, l, tv}

	// GET /api/v1/organizations/{orgId}/policy
	handler.GET("/organizations/:orgId/policy", r.GetRules)

	// PUT /api/v1/organizations/{orgId}/policy
	handler.PUT("/organizations/:orgId/policy", r.UpdateRules)
}

func (rs *rulesRoutes) GetRules(c echo.Context) error {
	const op = "Controller.GetRules"

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
	userID, err := rs.tv.GetIdentityID(ctx, token)
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

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		err = errorResponse(c, http.StatusBadRequest, "invalid organization id")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "invalid organization id")
	}

	id, err := rs.t.GetByOrganizationID(ctx, orgUUID, userID)
	if err != nil {
		err = errorResponse(c, http.StatusBadRequest, "failed to get rules")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, id)
}

func (rs *rulesRoutes) UpdateRules(c echo.Context) error {
	const op = "Controller.GetRules"

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
	userID, err := rs.tv.GetIdentityID(ctx, token)
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

	u := new(entity.UpdatePolicy)
	if err = c.Bind(u); err != nil {
		err = errorResponse(c, http.StatusBadRequest, "bad request")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	orgUUID, err := uuid.Parse(organizationID)
	if err != nil {
		err = errorResponse(c, http.StatusBadRequest, "invalid organization id")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %s", op, "invalid organization id")
	}

	bp := &entity.BookingPolicy{
		OrganizationID:           orgUUID,
		MaxBookingDurationMin:    pointer.Get(u.MaxBookingDurationMin),
		BookingWindowDays:        pointer.Get(u.BookingWindowDays),
		MaxActiveBookingsPerUser: pointer.Get(u.MaxActiveBookingsPerUser),
	}

	rule, err := rs.t.UpdateRule(ctx, bp, userID)
	if err != nil {
		err = errorResponse(c, http.StatusBadRequest, "failed to get rules")
		if err != nil {
			return fmt.Errorf("%s-%s: %w", op, "failed to sent response", err)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return c.JSON(http.StatusOK, rule)
}
