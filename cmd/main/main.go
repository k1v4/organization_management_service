package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/k1v4/organization_management_service/internal/config"
	v1 "github.com/k1v4/organization_management_service/internal/controller/http/v1"
	"github.com/k1v4/organization_management_service/internal/repository"
	"github.com/k1v4/organization_management_service/internal/usecase"
	"github.com/k1v4/organization_management_service/pkg/database/postgres"
	"github.com/k1v4/organization_management_service/pkg/httpserver"
	"github.com/k1v4/organization_management_service/pkg/jwtpkg"
	"github.com/k1v4/organization_management_service/pkg/logger"
	"github.com/labstack/echo/v4"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	ctx := context.Background()
	serviceLogger := logger.NewLogger()
	ctx = context.WithValue(ctx, logger.LoggerKey, serviceLogger)

	cfg, err := config.LoadConfig()
	if err != nil {
		serviceLogger.Error(ctx, err.Error())
		return
	}

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	pg, err := postgres.New(url, postgres.MaxPoolSize(cfg.PoolMax))
	if err != nil {
		serviceLogger.Error(ctx, fmt.Sprintf("app - Run - postgres.New: %s", err))
		return
	}
	defer pg.Close()

	serviceLogger.Info(ctx, "connected to database successfully")

	m, err := migrate.New(
		"file://migrations",
		url,
	)
	if err != nil {
		serviceLogger.Error(ctx, fmt.Sprintf("Migration setup failed: %v", err))
		return
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		serviceLogger.Error(ctx, fmt.Sprintf("Migration failed: %v", err))
		return
	}

	serviceLogger.Info(ctx, "Migrations applied successfully!")

	rulesRepo := repository.NewRulesRepository(pg)
	organizationsRepo := repository.NewOrganizationRepository(pg)

	serviceRules := usecase.NewRuleUseCase(rulesRepo)
	serviceOrg := usecase.NewOrganizationUseCase(organizationsRepo)

	handler := echo.New()
	// TODO через конфиг урл
	tv, err := jwtpkg.NewTokenVerifier(ctx, "")
	if err != nil {
		serviceLogger.Error(ctx, fmt.Sprintf("app - Run - jwtpkg.NewTokenVerifier: %s", err))
		return
	}

	settings := v1.FillRouterSettings(handler, serviceLogger, serviceOrg, serviceRules, cfg, tv)

	v1.NewRouter(*settings)

	httpServer := httpserver.New(handler, httpserver.Port(strconv.Itoa(cfg.RestServerPort)))

	// signal for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		serviceLogger.Info(ctx, "app-Run-signal: "+s.String())
	case err = <-httpServer.Notify():
		serviceLogger.Error(ctx, fmt.Sprintf("app-Run-httpServer.Notify: %s", err))
	}

	// shutdown
	err = httpServer.Shutdown()
	if err != nil {
		serviceLogger.Error(ctx, fmt.Sprintf("app-Run-httpServer.Shutdown: %s", err))
	}
}
