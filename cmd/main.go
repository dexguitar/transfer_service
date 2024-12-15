package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dexguitar/p2p_service/config"
	"github.com/dexguitar/p2p_service/internal/app"
	"github.com/joho/godotenv"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	err = runMigrations(cfg)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	log.Println("Starting server on :8080...")
	if err := application.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}

func runMigrations(c *config.Config) error {
	op := "db.Migrate"

	dbConn, err := newDatabase(c)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = dbConn.Ping()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer dbConn.Close()

	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func newDatabase(c *config.Config) (*sql.DB, error) {
	op := "db.NewDatabase"

	db, err := sql.Open("postgres", "postgres://postgres:qwerty@localhost:5434/postgres?sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
