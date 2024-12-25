package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	dsn := flag.String("dsn", "postgres://postgres:test-password@localhost:5432/textvault?sslmode=disable", "PostgreSQL data source name")
	up := flag.Bool("up", false, "Run migrations up")
	down := flag.Bool("down", false, "Run migrations down")

	flag.Parse()

	pool, err := pgxpool.New(context.Background(), *dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	migrationsPath := "./migrations"

	if *up {
		if err := goose.Up(db, migrationsPath); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	if *down {
		if err := goose.Down(db, migrationsPath); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	fmt.Println("Migrations applied successfully!")
}
