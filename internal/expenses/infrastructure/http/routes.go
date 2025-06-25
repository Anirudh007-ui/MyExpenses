// Package http contains the HTTP handlers for the expense API
// This file configures the routing for all expense-related endpoints
// It maps HTTP requests to the appropriate handler methods
package http

import (
	"myexpenses/internal/expenses/application" // Import our application layer

	"github.com/gin-gonic/gin" // Gin is a high-performance HTTP web framework for Go
)

// SetupRoutes configures the expense routes
// This function takes a Gin router and application service, then sets up all the routes
// It's called from main.go to wire up the HTTP layer
func SetupRoutes(router *gin.Engine, service *application.Service) {
	// Create a new handler instance with the service dependency
	// This follows dependency injection - the handler gets its dependencies from outside
	handler := NewHandler(service)

	// Create a route group for all expense-related endpoints
	// Route groups help organize related endpoints and can share middleware
	// The "/expenses" prefix will be added to all routes in this group
	expenses := router.Group("/expenses")
	{
		// POST /expenses - Create a new expense
		// This route accepts JSON data in the request body and creates a new expense
		// The empty string "" means no additional path beyond the group prefix
		expenses.POST("", handler.CreateExpense)

		// GET /expenses - Get all expenses (with optional filtering)
		// This route can accept query parameters for filtering (e.g., ?category=Food)
		expenses.GET("", handler.GetAllExpenses)

		// GET /expenses/{id} - Get a specific expense by ID
		// The {id} is a URL parameter that gets passed to the handler
		// For example, GET /expenses/123e4567-e89b-12d3-a456-426614174000
		expenses.GET("/:id", handler.GetExpense)

		// PUT /expenses/{id} - Update an existing expense
		// This route accepts JSON data in the request body and updates the specified expense
		expenses.PUT("/:id", handler.UpdateExpense)

		// DELETE /expenses/{id} - Delete an expense
		// This route removes the specified expense from the system
		expenses.DELETE("/:id", handler.DeleteExpense)
	}

	// Note: This follows RESTful conventions:
	// - POST for creating new resources
	// - GET for retrieving resources
	// - PUT for updating existing resources
	// - DELETE for removing resources
	// - URLs use nouns (expenses) not verbs
	// - HTTP status codes indicate the result (200 OK, 201 Created, 404 Not Found, etc.)
}
