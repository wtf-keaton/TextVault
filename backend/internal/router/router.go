package router

import (
	"TextVault/internal/router/services/account"
	"TextVault/internal/router/services/pastes"
	"TextVault/internal/storage/cloud"
	"TextVault/internal/storage/postgres"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app *fiber.App
	log *slog.Logger

	accountService *account.Service
	pasteService   *pastes.Service
}

func New(storage *postgres.Storage, S3 *cloud.Storage, log *slog.Logger) *Router {
	app := fiber.New(fiber.Config{
		AppName:               "TextVault API",
		DisableStartupMessage: true,
	})

	accountService := account.New(log, storage, storage)
	pasteService := pastes.New(log, storage, storage, S3)

	return &Router{
		app:            app,
		log:            log,
		accountService: accountService,
		pasteService:   pasteService,
	}
}

func (r *Router) setupAccountRoutes(app *fiber.App) {
	accountApi := app.Group("/account")
	accountApi.Post("/register", r.accountService.Register)
	accountApi.Post("/login", r.accountService.Login)
	accountApi.Get("/pastes", r.accountService.GetUserPastes)
}

func (r *Router) setupPastesRoutes(app *fiber.App) {
	pasteApi := app.Group("/pastes")
	pasteApi.Post("/", r.pasteService.SavePaste)
	pasteApi.Get("/:hash", r.pasteService.GetPaste)
	pasteApi.Delete("/:hash", r.pasteService.DeletePaste)
}

func (r *Router) setupRoutes() {
	r.setupAccountRoutes(r.app)
	r.setupPastesRoutes(r.app)
}

func (r *Router) MustRun() {
	const prefix = "internal.router.MustRun"
	log := r.log.With(
		slog.String("op", prefix),
	)

	log.Info("Setupping routes")

	r.setupRoutes()

	log.Info("Starting router")
	if err := r.run(); err != nil {
		panic(err)
	}
}

func (r *Router) run() error {
	fmt.Println("Server started on port 8080")

	return r.app.Listen(":8080")
}
