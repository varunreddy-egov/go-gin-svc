package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server configuration
	HTTPPort          string
	ServerContextPath string

	// Database configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Migration script configuration
	MigrationScriptPath string
	MigrationEnabled    bool
}

func Load() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it, relying on system env vars.")
	}

	return &Config{
		// Server configuration
		HTTPPort:          getEnv("HTTP_PORT", "8080"),
		ServerContextPath: getEnv("SERVER_CONTEXT_PATH", "template-config/v1"),

		// Database configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "template_config"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),

		//Migration script configuration
		MigrationScriptPath: getEnv("MIGRATION_SCRIPT_PATH", "./migrations"),
		MigrationEnabled:    getEnvAsBool("MIGRATION_ENABLED", false),
	}
}

func getEnv(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}
