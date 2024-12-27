package main

import (
	"context"
	"errors"
	"flag"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/slamchillz/getinstashop-ecommerce-api/cmd/server"
	"github.com/slamchillz/getinstashop-ecommerce-api/config"
	_ "github.com/slamchillz/getinstashop-ecommerce-api/docs"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/utils"
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"log"
	"os"
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
	// Define command line flags for email and password for admin user generation
	email := flag.String("email", "", "User email address")
	password := flag.String("password", "", "User password")

	// Parse the flags
	flag.Parse()

	// Check if both email and password are provided when one of them is
	if (*email != "" && *password == "") || (*email == "" && *password != "") {
		log.Fatal("Error: If one of 'email' or 'password' is provided, both must be provided.")
	}

	serverConfig, err := config.LoadConfig("")
	if err != nil {
		log.Fatalf("Error loading application config: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	connPool, err := pgxpool.New(ctx, serverConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	runDBMigration(serverConfig.DatabaseURL)
	store := db.NewStore(connPool)
	// If both are provided we create an admin user and exit
	if *email != "" && *password != "" {
		user, err := store.CreateUser(ctx, db.CreateUserParams{
			ID:       uuid.New(),
			Email:    *email,
			Password: *password,
		})
		if err != nil {
			log.Fatalf("Error creating admin user: %v", err)
		}
		log.Printf("Admin user created. Email: %v, Password: %v", user.Email, user.Password)
		os.Exit(0)
	}
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

func runDBMigration(databaseURL string) {
	migration, err := migrate.New(utils.GenerateMigrationFilePath(), databaseURL)
	if err != nil {
		log.Fatalf("Error creating migration: %v", err)
	}
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Error running migration: %v", err)
	}
	log.Printf("Migration completed successfully")
}
