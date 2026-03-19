package v1

import (
	"github.com/k1v4/organization_management_service/internal/entity"
	"github.com/labstack/echo/v4"
)

func errorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, entity.ErrorResponse{Error: msg})
}
