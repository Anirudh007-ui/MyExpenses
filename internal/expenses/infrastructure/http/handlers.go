// Package http contains the HTTP handlers for the expense API
// This is part of the infrastructure layer - it handles HTTP-specific concerns
// It translates HTTP requests into business operations and formats responses
package http

import (
	"net/http" // Go's built-in HTTP package for status codes and request/response handling
	"strconv"  // For converting strings to numbers (used for query parameters)

	// For handling dates and times
	"myexpenses/internal/expenses/application" // Import our application layer

	"github.com/gin-gonic/gin" // Gin is a high-performance HTTP web framework for Go
)

// Handler handles HTTP requests for expenses
// This struct holds a reference to the application service
// It acts as a bridge between HTTP concerns and business logic
type Handler struct {
	// service is a dependency on the application service
	// This follows dependency injection - the handler doesn't create the service, it receives it
	service *application.Service
}

// NewHandler creates a new expense handler
// This is a constructor function that implements dependency injection
func NewHandler(service *application.Service) *Handler {
	return &Handler{
		service: service, // Store the service dependency
	}
}

// CreateExpense handles POST /expenses
// This method processes HTTP POST requests to create new expenses
// It follows the REST convention where POST creates new resources
func (h *Handler) CreateExpense(c *gin.Context) {
	// Step 1: Declare a variable to hold the parsed request data
	// This struct will be populated with the JSON data from the request body
	var req application.CreateExpenseRequest

	// Step 2: Parse and validate the JSON request body
	// ShouldBindJSON automatically validates the request based on struct tags
	// If validation fails, it returns an error
	if err := c.ShouldBindJSON(&req); err != nil {
		// Step 3: Return a 400 Bad Request response if validation fails
		// gin.H is a helper for creating map literals (like map[string]interface{})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body", // User-friendly error message
			"details": err.Error(),            // Technical details for debugging
		})
		return // Exit the function early
	}

	// Step 4: Call the business logic to create the expense
	// c.Request.Context() provides the HTTP request context for cancellation/timeout
	expense, err := h.service.CreateExpense(c.Request.Context(), &req)
	if err != nil {
		// Step 5: Return a 500 Internal Server Error if business logic fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create expense",
		})
		return
	}

	// Step 6: Return a 201 Created response with the created expense
	// 201 is the standard HTTP status code for successful resource creation
	c.JSON(http.StatusCreated, gin.H{
		"message": "Expense created successfully", // Success message
		"data":    expense,                        // The created expense data
	})
}

// GetExpense handles GET /expenses/{id}
// This method processes HTTP GET requests to retrieve a specific expense
// The {id} part is a URL parameter that gets passed to this handler
func (h *Handler) GetExpense(c *gin.Context) {
	// Step 1: Extract the ID from the URL parameters
	// c.Param("id") gets the value of the "id" parameter from the URL
	id := c.Param("id")
	if id == "" {
		// Step 2: Return 400 Bad Request if no ID is provided
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Expense ID is required",
		})
		return
	}

	// Step 3: Call the business logic to get the expense
	expense, err := h.service.GetExpense(c.Request.Context(), id)
	if err != nil {
		// Step 4: Handle different types of errors
		if err.Error() == "expense not found" {
			// Return 404 Not Found if the expense doesn't exist
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Expense not found",
			})
			return
		}
		// Return 500 Internal Server Error for other errors
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get expense",
		})
		return
	}

	// Step 5: Return 200 OK with the expense data
	c.JSON(http.StatusOK, gin.H{
		"data": expense,
	})
}

// GetAllExpenses handles GET /expenses
// This method processes HTTP GET requests to retrieve all expenses with optional filtering
// It supports query parameters for filtering the results
func (h *Handler) GetAllExpenses(c *gin.Context) {
	// Step 1: Create a map to hold filter criteria
	// map[string]interface{} is a map where keys are strings and values can be any type
	filters := make(map[string]interface{})

	// Step 2: Parse query parameters and add them to filters
	// Query parameters are the part of the URL after the ? (e.g., ?category=Food&min_amount=10)

	// Check for category filter
	if category := c.Query("category"); category != "" {
		// c.Query("category") gets the value of the "category" query parameter
		// If it's not empty, add it to the filters map
		filters["category"] = category
	}

	// Check for date range filters
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		filters["date_from"] = dateFrom
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		filters["date_to"] = dateTo
	}

	// Check for amount range filters
	// strconv.ParseFloat converts a string to a float64
	if minAmountStr := c.Query("min_amount"); minAmountStr != "" {
		// Parse the string to a float64, ignore the error if parsing fails
		if minAmount, err := strconv.ParseFloat(minAmountStr, 64); err == nil {
			filters["min_amount"] = minAmount
		}
	}

	if maxAmountStr := c.Query("max_amount"); maxAmountStr != "" {
		if maxAmount, err := strconv.ParseFloat(maxAmountStr, 64); err == nil {
			filters["max_amount"] = maxAmount
		}
	}

	// Check for description filter
	if description := c.Query("description"); description != "" {
		filters["description"] = description
	}

	// Step 3: Call the business logic to get filtered expenses
	expenses, err := h.service.GetAllExpenses(c.Request.Context(), filters)
	if err != nil {
		// Step 4: Return 500 Internal Server Error if business logic fails
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get expenses",
		})
		return
	}

	// Step 5: Return 200 OK with the expenses and count
	c.JSON(http.StatusOK, gin.H{
		"data":  expenses,      // The list of expenses
		"count": len(expenses), // The number of expenses returned
	})
}

// UpdateExpense handles PUT /expenses/{id}
// This method processes HTTP PUT requests to update existing expenses
// PUT is used for complete updates (though we allow partial updates in our implementation)
func (h *Handler) UpdateExpense(c *gin.Context) {
	// Step 1: Extract the ID from the URL parameters
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Expense ID is required",
		})
		return
	}

	// Step 2: Parse and validate the JSON request body
	var req application.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Step 3: Call the business logic to update the expense
	expense, err := h.service.UpdateExpense(c.Request.Context(), id, &req)
	if err != nil {
		// Step 4: Handle different types of errors
		if err.Error() == "expense not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Expense not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update expense",
		})
		return
	}

	// Step 5: Return 200 OK with the updated expense
	c.JSON(http.StatusOK, gin.H{
		"message": "Expense updated successfully",
		"data":    expense,
	})
}

// DeleteExpense handles DELETE /expenses/{id}
// This method processes HTTP DELETE requests to remove expenses
func (h *Handler) DeleteExpense(c *gin.Context) {
	// Step 1: Extract the ID from the URL parameters
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Expense ID is required",
		})
		return
	}

	// Step 2: Call the business logic to delete the expense
	err := h.service.DeleteExpense(c.Request.Context(), id)
	if err != nil {
		// Step 3: Handle different types of errors
		if err.Error() == "expense not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Expense not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete expense",
		})
		return
	}

	// Step 4: Return 200 OK with success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Expense deleted successfully",
	})
}
