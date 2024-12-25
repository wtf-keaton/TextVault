package app

import (
	"TextVault/internal/config"
	"TextVault/internal/router"
	"TextVault/internal/storage/postgres"
	"TextVault/internal/storage/redis"
	"TextVault/internal/storage/s3"
	"context"
	"log/slog"
)

type App struct {
	Router *router.Router
	log    *slog.Logger
}

func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	ctx := context.Background()

	s3Storage, err := s3.New(log, cfg.S3)
	if err != nil {
		return nil, err
	}

	log.Info("Connected to s3 storage")

	if err := s3Storage.BucketExists(ctx); err != nil {
		return nil, err
	}

	storage, err := postgres.New(ctx, log, &cfg.Postgres)
	if err != nil {
		return nil, err
	}

	if err := storage.Ping(ctx); err != nil {
		return nil, err
	}

	log.Info("Connected to database")

	redisStorage := redis.New(&cfg.Redis)

	if err := redisStorage.Ping(ctx); err != nil {
		return nil, err
	}

	log.Info("Connected to redis")

	router := router.New(storage, redisStorage, s3Storage, log)
	return &App{
		Router: router,
		log:    log,
	}, nil
}
