package main

import (
	"log/slog"
	"os"

	"TextVault/internal/app"
	"TextVault/internal/config"
	"TextVault/internal/lib/log/sl"
)

const (
	envLocal = "local"
	envProd  = "production"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	app, err := app.New(log, cfg)
	if err != nil {
		log.Error("failed to create app", sl.Err(err))
		os.Exit(1)
	}

	app.Router.MustRun()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
