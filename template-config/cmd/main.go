package main

import (
	"fmt"
	"log"
	"template-config/config"
	"template-config/db"
	"template-config/routes"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func buildPostgresDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
}

func runMigrations(databaseURL string, migrationsPath string) {
	m, err := migrate.New(
		"file://"+migrationsPath,
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migrated successfully")
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Build DSN
	dsn := buildPostgresDSN(cfg)

	// Run migrations before anything else
	runMigrations(dsn, "./migrations")

	// Setup database
	dbConn, err := db.ConnectDSN(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup routes
	router := routes.SetupRoutes(dbConn, cfg)

	// Start server
	log.Printf("Starting server on :%s", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
