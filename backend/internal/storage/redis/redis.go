package redis

import (
	"TextVault/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	rdb *redis.Client
}

func New(cfg *config.RedisConfig) *Storage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Storage{rdb: rdb}
}

func (s *Storage) Close() error {
	return s.rdb.Close()
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.rdb.Ping(ctx).Err()
}

func (s *Storage) Set(ctx context.Context, key string, value string) error {
	return s.rdb.Set(ctx, key, value, time.Hour*12).Err()
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	return s.rdb.Get(ctx, key).Result()
}

func (s *Storage) Delete(ctx context.Context, key string) error {
	return s.rdb.Del(ctx, key).Err()
}

func (s *Storage) Exists(ctx context.Context, key string) error {
	return s.rdb.Exists(ctx, key).Err()
}
