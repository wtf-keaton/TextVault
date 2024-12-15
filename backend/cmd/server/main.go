package main

import (
	"context"
	"log"

	"TextVault/internal/router"
	"TextVault/internal/storage/postgres"
)

func main() {
	storage, err := postgres.New(context.Background())
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer storage.Close()

	if err := storage.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Connected to database")

	route := router.New(storage)

	route.MustRun()
}
