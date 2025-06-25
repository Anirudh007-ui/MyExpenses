// Package db contains database configuration and connection logic
// This package handles all database-related concerns like connection setup and configuration
// It's part of the infrastructure layer - external concerns like database connections
package db

import (
	"fmt" // For formatted string operations (building connection strings)
	"log" // For logging database connection status
	"os"  // For reading environment variables

	"gorm.io/driver/postgres" // GORM's PostgreSQL driver
	"gorm.io/gorm"            // GORM ORM library
	"gorm.io/gorm/logger"     // GORM's logging configuration
)

// Config holds database configuration
// This struct centralizes all database connection parameters
// It makes it easy to manage database settings in one place
type Config struct {
	// Host is the database server address (e.g., "localhost", "192.168.1.100")
	Host string

	// Port is the database server port (e.g., "5432" for PostgreSQL)
	Port string

	// User is the database username for authentication
	User string

	// Password is the database password for authentication
	Password string

	// DBName is the name of the database to connect to
	DBName string

	// SSLMode determines the SSL connection mode
	// Common values: "disable", "require", "verify-ca", "verify-full"
	SSLMode string
}

// NewConfig creates a new database configuration from environment variables
// This function reads database settings from environment variables with sensible defaults
// Environment variables allow for different configurations in different environments (dev, staging, prod)
func NewConfig() *Config {
	return &Config{
		// Read from environment variable DB_HOST, default to "localhost" if not set
		Host: getEnv("DB_HOST", "localhost"),

		// Read from environment variable DB_PORT, default to "5432" if not set
		Port: getEnv("DB_PORT", "5432"),

		// Read from environment variable DB_USER, default to "postgres" if not set
		User: getEnv("DB_USER", "postgres"),

		// Read from environment variable DB_PASSWORD, default to "password" if not set
		Password: getEnv("DB_PASSWORD", "password"),

		// Read from environment variable DB_NAME, default to "myexpenses" if not set
		DBName: getEnv("DB_NAME", "myexpenses"),

		// Read from environment variable DB_SSLMODE, default to "disable" if not set
		SSLMode: getEnv("DB_SSLMODE", "disable"),
	}
}

// Connect establishes a connection to PostgreSQL
// This function takes a config and returns a GORM database connection
// It handles the connection string building and connection testing
func Connect(config *Config) (*gorm.DB, error) {
	// Build the PostgreSQL connection string (DSN - Data Source Name)
	// fmt.Sprintf formats a string with the provided values
	// The format follows PostgreSQL's connection string specification
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host,     // Database host
		config.Port,     // Database port
		config.User,     // Database user
		config.Password, // Database password
		config.DBName,   // Database name
		config.SSLMode,  // SSL mode
	)

	// Open a database connection using GORM
	// postgres.Open(dsn) creates a PostgreSQL driver with our connection string
	// &gorm.Config{} provides configuration options for GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Configure GORM's logging
		// logger.Default.LogMode(logger.Info) enables SQL query logging
		// This is useful for debugging but can be verbose in production
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		// If connection fails, return an error with context
		// %w preserves the original error while adding our message
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection by getting the underlying *sql.DB and pinging it
	// This ensures the database is actually reachable and responsive
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Ping the database to verify connectivity
	// This sends a simple query to the database to ensure it's working
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Log successful connection
	log.Println("Successfully connected to PostgreSQL database")

	// Return the GORM database connection
	return db, nil
}

// getEnv gets an environment variable with a fallback default value
// This helper function simplifies reading environment variables
// It returns the environment variable value if set, otherwise returns the fallback
func getEnv(key, fallback string) string {
	// os.Getenv(key) reads the environment variable with the given key
	// It returns an empty string if the variable is not set
	if value := os.Getenv(key); value != "" {
		// If the environment variable is set and not empty, return its value
		return value
	}
	// If the environment variable is not set or is empty, return the fallback value
	return fallback
}
