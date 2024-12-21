package app

import (
	"TextVault/internal/config"
	"TextVault/internal/router"
	"TextVault/internal/storage/cloud"
	"TextVault/internal/storage/postgres"
	"context"
	"log/slog"
)

type App struct {
	Router *router.Router
	log    *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	ctx := context.Background()

	cloudStorage, err := cloud.New(log, cfg.S3)
	if err != nil {
		return nil, err
	}

	log.Info("Connected to s3 storage")

	if err := cloudStorage.BucketExists(ctx); err != nil {
		return nil, err
	}

	storage, err := postgres.New(ctx, log, &cfg.Postgres)
	if err != nil {
		return nil, err
	}

	defer storage.Close()

	if err := storage.Ping(ctx); err != nil {
		return nil, err
	}

	log.Info("Connected to database")

	router := router.New(storage, log)
	return &App{
		Router: router,
		log:    log,
	}, nil
}
