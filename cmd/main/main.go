package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/k1v4/organization_management_service/internal/config"
	"github.com/k1v4/organization_management_service/pkg/database/postgres"
	"github.com/k1v4/organization_management_service/pkg/logger"

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
}
