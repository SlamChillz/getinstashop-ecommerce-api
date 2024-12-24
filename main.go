package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/slamchillz/getinstashop-ecommerce-api/cmd/server"
	"github.com/slamchillz/getinstashop-ecommerce-api/config"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"log"
)

func main() {
	serverConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading application config: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connPool, err := pgxpool.New(ctx, serverConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	runDBMigration(serverConfig.MigrationURL, serverConfig.DatabaseURL)
	store := db.NewStore(connPool)
	runHTTPServer(serverConfig, store)
}

func runHTTPServer(config config.Config, store db.Store) {
	apiServer, err := server.NewServer(config, store)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	err = apiServer.Start()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func runDBMigration(migrationUrl, databaseURL string) {
	fmt.Printf("Running database migration on %s\n", migrationUrl)
	migration, err := migrate.New(migrationUrl, databaseURL)
	if err != nil {
		log.Fatalf("Error creating migration: %v", err)
	}
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Error running migration: %v", err)
	}
	log.Printf("Migration completed successfully")
}
