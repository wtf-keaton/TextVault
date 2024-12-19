package main

import (
	"context"
	"log"

	"TextVault/internal/router"
	"TextVault/internal/storage/cloud"
	"TextVault/internal/storage/postgres"
)

func main() {
	ctx := context.Background()

	cloudStorage, err := cloud.New()
	if err != nil {
		log.Fatalf("Unable to connect to cloud: %v\n", err)
	}

	log.Println("Connected to s3 storage")

	if err := cloudStorage.BucketExists(ctx); err != nil {
		log.Fatalf("Unable to check if bucket exists: %v\n", err)
	}

	log.Println("Bucket 'textvault' exists and you already own it")

	storage, err := postgres.New(cloudStorage, ctx)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer storage.Close()

	if err := storage.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	log.Println("Connected to database")

	route := router.New(storage)

	route.MustRun()
}
