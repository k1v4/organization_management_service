package v1

import (
	"log"
	"net/http"

	"github.com/k1v4/organization_management_service/internal/config"
	"github.com/k1v4/organization_management_service/internal/usecase"
	"github.com/k1v4/organization_management_service/pkg/jwtpkg"
	"github.com/k1v4/organization_management_service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RouterSettings struct {
	handler *echo.Echo
	l       logger.Logger
	o       IOrganizationService
	r       usecase.RuleUseCase
	cfg     config.Config
	tv      jwtpkg.TokenVerifier
}

func fillRouterSettings(
	handler *echo.Echo,
	l logger.Logger,
	o IOrganizationService,
	r usecase.RuleUseCase,
	cfg config.Config,
	tv jwtpkg.TokenVerifier,
) *RouterSettings {
	return &RouterSettings{
		handler: handler,
		l:       l,
		o:       o,
		r:       r,
		cfg:     cfg,
		tv:      tv,
	}
}

// t usecase.IArticleService
func NewRouter(rs RouterSettings) {
	// Middleware
	rs.handler.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("%s %s %d %s\n", v.Method, v.URI, v.Status, v.Latency)
			return nil
		},
	}))
	rs.handler.Use(middleware.Recover())
	rs.handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO добавить хосты пацанов через энвы
		AllowOrigins:     []string{"http://localhost:3000"},                                                                // Разрешить запросы с этого origin
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},                               // Разрешенные методы
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}, // Разрешенные заголовки
		AllowCredentials: true,                                                                                             // Разрешить передачу кук и заголовков авторизации
	}))

	rs.handler.GET("/api/article/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	h := rs.handler.Group("/api/v1")
	{
		newOrganizationRoutes(h, rs.o, rs.l, rs.tv)
	}
}
