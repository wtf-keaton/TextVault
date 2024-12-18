package router

import (
	"TextVault/internal/router/services/paste"
	"TextVault/internal/router/services/user"
	"TextVault/internal/storage/postgres"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	app *fiber.App

	UserService  *user.Service
	PasteService *paste.Service
}

func New(storage *postgres.Storage) *Router {
	app := fiber.New(fiber.Config{
		AppName: "TextVault API",
	})

	userService := user.New(storage, storage)
	pasteService := paste.New(storage, storage)

	return &Router{
		app:          app,
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
}

func (r *Router) MustRun() {
	const prefix = "internal.router.MustRun"

	log.Println(prefix, ": Setupping routes")

	r.setupRoutes()

	log.Println(prefix, ": Starting router")
	if err := r.run(); err != nil {
		panic(err)
	}
}

func (r *Router) run() error {
	return r.app.Listen(":8080")
}
