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
	_ "github.com/slamchillz/getinstashop-ecommerce-api/docs"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
)

// @title           Swagger Example API
// @version         3.0
// @description     This InstaShop e-commerce technical assessment REST API server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in							header
// @name						Authorization
// @security BearerAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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
