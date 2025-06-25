// Package main is the entry point for the MyExpenses API application
// This file orchestrates the startup of the entire application
// It wires together all the layers following Clean Architecture principles
package main

import (
	"log" // For logging application startup and errors
	"os"  // For reading environment variables and getting port

	"myexpenses/internal/db"                   // Database configuration
	"myexpenses/internal/expenses/application" // Business logic layer

	// Domain layer (for error types)
	"myexpenses/internal/expenses/infrastructure/http"     // HTTP handlers and routes
	"myexpenses/internal/expenses/infrastructure/postgres" // Database implementation

	"github.com/gin-gonic/gin" // HTTP web framework
	"github.com/joho/godotenv" // For loading .env files
)

// main is the entry point function that gets called when the application starts
// This function follows the dependency injection pattern to wire up all components
func main() {
	// Step 1: Load environment variables from .env file
	// godotenv.Load() reads a .env file and sets environment variables
	// This is useful for development - in production, environment variables are set differently
	if err := godotenv.Load(); err != nil {
		// If no .env file is found, that's okay - we'll use system environment variables
		log.Println("No .env file found, using system environment variables")
	}

	// Step 2: Initialize database configuration
	// NewConfig() reads database settings from environment variables
	// It provides sensible defaults if environment variables are not set
	dbConfig := db.NewConfig()

	// Step 3: Connect to the database
	// Connect() establishes a connection to PostgreSQL using the configuration
	database, err := db.Connect(dbConfig)
	if err != nil {
		// If database connection fails, log the error and exit
		// log.Fatalf() prints the error and calls os.Exit(1)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Step 4: Initialize the repository layer
	// NewRepository() creates a PostgreSQL implementation of the repository interface
	// This is where we choose which database implementation to use
	repo := postgres.NewRepository(database)

	// Step 5: Run database migrations
	// AutoMigrate() creates database tables based on our struct definitions
	// It ensures the database schema matches our domain models
	if err := repo.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Step 6: Initialize the application service layer
	// NewService() creates the business logic layer with the repository dependency
	// This follows dependency injection - the service gets its dependencies from outside
	service := application.NewService(repo)

	// Step 7: Initialize the HTTP server
	// gin.Default() creates a new Gin router with default middleware
	// Default middleware includes logging and panic recovery
	router := gin.Default()

	// Step 8: Add additional middleware
	// Middleware functions process requests before they reach handlers
	// They can add logging, authentication, CORS, etc.
	router.Use(gin.Logger())   // Logs HTTP requests (method, path, status, duration)
	router.Use(gin.Recovery()) // Recovers from panics and returns 500 errors

	// Step 9: Setup API routes
	// SetupRoutes() configures all the expense endpoints
	// It maps HTTP requests to the appropriate handler methods
	http.SetupRoutes(router, service)

	// Step 10: Add a health check endpoint
	// This endpoint is useful for load balancers and monitoring systems
	// It allows external systems to check if the API is running
	router.GET("/health", func(c *gin.Context) {
		// Return a simple JSON response indicating the service is healthy
		c.JSON(200, gin.H{
			"status":  "ok",             // Health status
			"service": "MyExpenses API", // Service name
		})
	})

	// Step 11: Get the port from environment or use default
	// os.Getenv("PORT") reads the PORT environment variable
	port := os.Getenv("PORT")
	if port == "" {
		// If no PORT is set, use the default port 8080
		port = "8080"
	}

	// Step 12: Start the HTTP server
	// Log that we're starting the server
	log.Printf("Starting server on port %s", port)

	// router.Run() starts the HTTP server and blocks until the server stops
	// It listens for incoming HTTP requests on the specified port
	if err := router.Run(":" + port); err != nil {
		// If the server fails to start, log the error and exit
		log.Fatalf("Failed to start server: %v", err)
	}

	// Note: The application will run indefinitely until interrupted
	// To stop the server, send a SIGINT signal (Ctrl+C) or SIGTERM
}
