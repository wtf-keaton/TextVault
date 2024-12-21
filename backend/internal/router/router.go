package router

import (
	"TextVault/internal/router/services/paste"
	"TextVault/internal/router/services/user"
	"TextVault/internal/storage/cloud"
	"TextVault/internal/storage/postgres"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app *fiber.App
	log *slog.Logger

	UserService  *user.Service
	PasteService *paste.Service
}

func New(storage *postgres.Storage, S3 *cloud.Storage, log *slog.Logger) *Router {
	app := fiber.New(fiber.Config{
		AppName:               "TextVault API",
		DisableStartupMessage: true,
	})

	userService := user.New(log, storage, storage)
	pasteService := paste.New(log, storage, storage, S3)

	return &Router{
		app:          app,
		log:          log,
		UserService:  userService,
		PasteService: pasteService,
	}
}

func (r *Router) setupRoutes() {
	api := r.app.Group("/api/v1")

	userApi := api.Group("/user")
	userApi.Post("/register", r.UserService.Register)
	userApi.Post("/login", r.UserService.Login)
	userApi.Get("/validate", r.UserService.ValidateToken)

	pasteApi := api.Group("/paste")
	pasteApi.Post("/save", r.PasteService.SavePaste)
	pasteApi.Get("/get/:hash", r.PasteService.GetPaste)
	pasteApi.Delete("/delete/:hash", r.PasteService.DeletePaste)

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
