// Package postgres contains the PostgreSQL implementation of the repository interface
// This is part of the infrastructure layer - it handles external concerns like database operations
// It implements the domain.Repository interface, making it a concrete implementation
package postgres

import (
	"context" // For request context (cancellation, timeouts)
	"fmt"     // For formatted string operations and error wrapping

	// For string manipulation (though not used in this implementation)
	"myexpenses/internal/expenses/domain" // Import our domain layer

	"github.com/google/uuid" // For UUID parsing and validation
	"gorm.io/gorm"           // GORM is an ORM (Object-Relational Mapping) library for Go
)

// Repository implements the domain.Repository interface using PostgreSQL
// This struct holds a reference to the GORM database connection
// It provides the concrete implementation of all repository methods
type Repository struct {
	// db is the GORM database connection
	// GORM provides a convenient way to interact with databases using Go structs
	db *gorm.DB
}

// NewRepository creates a new PostgreSQL repository
// This is a constructor function that takes a GORM database connection
// It returns a configured repository instance
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db, // Store the database connection
	}
}

// Create adds a new expense to the database
// This method implements the domain.Repository.Create interface
func (r *Repository) Create(ctx context.Context, expense *domain.Expense) error {
	// Use GORM's Create method to insert the expense into the database
	// WithContext(ctx) propagates the context for cancellation/timeout handling
	// Create() automatically handles the SQL INSERT statement
	return r.db.WithContext(ctx).Create(expense).Error
}

// GetByID retrieves an expense by its ID
// This method implements the domain.Repository.GetByID interface
func (r *Repository) GetByID(ctx context.Context, id string) (*domain.Expense, error) {
	// Step 1: Parse the string ID into a UUID
	// This validates that the ID is a proper UUID format
	uuid, err := uuid.Parse(id)
	if err != nil {
		// If the ID is not a valid UUID, return an error
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	// Step 2: Declare a variable to hold the result
	// This will be populated by GORM when the query executes
	var expense domain.Expense

	// Step 3: Execute the database query
	// WithContext(ctx) - propagates context for cancellation/timeout
	// Where("id = ?", uuid) - adds a WHERE clause to filter by ID
	// First(&expense) - gets the first matching record and stores it in expense
	// .Error - gets any error that occurred during the query
	if err := r.db.WithContext(ctx).Where("id = ?", uuid).First(&expense).Error; err != nil {
		// Step 4: Handle specific error cases
		if err == gorm.ErrRecordNotFound {
			// If no record was found, return our domain-specific error
			return nil, domain.ErrExpenseNotFound
		}
		// For any other database error, wrap it with context
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	// Step 5: Return the found expense
	// &expense returns a pointer to the expense
	return &expense, nil
}

// GetAll retrieves all expenses with optional filtering
// This method implements the domain.Repository.GetAll interface
func (r *Repository) GetAll(ctx context.Context, filters map[string]interface{}) ([]*domain.Expense, error) {
	// Step 1: Declare a slice to hold the results
	// []*domain.Expense is a slice of pointers to Expense structs
	var expenses []*domain.Expense

	// Step 2: Start building the query
	// WithContext(ctx) propagates context for cancellation/timeout
	query := r.db.WithContext(ctx)

	// Step 3: Apply filters to the query
	// This loop iterates through each filter and adds WHERE clauses
	for key, value := range filters {
		switch key {
		case "category":
			// Filter by category with partial matching (case-insensitive)
			if category, ok := value.(string); ok && category != "" {
				// ILIKE is PostgreSQL's case-insensitive LIKE operator
				// %category% means "contains the category text anywhere"
				query = query.Where("category ILIKE ?", "%"+category+"%")
			}
		case "date_from":
			// Filter expenses from a specific date onwards
			if dateFrom, ok := value.(string); ok && dateFrom != "" {
				query = query.Where("date >= ?", dateFrom)
			}
		case "date_to":
			// Filter expenses up to a specific date
			if dateTo, ok := value.(string); ok && dateTo != "" {
				query = query.Where("date <= ?", dateTo)
			}
		case "min_amount":
			// Filter expenses with amount greater than or equal to min_amount
			if minAmount, ok := value.(float64); ok && minAmount > 0 {
				query = query.Where("amount >= ?", minAmount)
			}
		case "max_amount":
			// Filter expenses with amount less than or equal to max_amount
			if maxAmount, ok := value.(float64); ok && maxAmount > 0 {
				query = query.Where("amount <= ?", maxAmount)
			}
		case "description":
			// Filter by description with partial matching (case-insensitive)
			if description, ok := value.(string); ok && description != "" {
				query = query.Where("description ILIKE ?", "%"+description+"%")
			}
		}
	}

	// Step 4: Add ordering to the query
	// Order by date descending (newest expenses first)
	query = query.Order("date DESC")

	// Step 5: Execute the query and populate the expenses slice
	if err := query.Find(&expenses).Error; err != nil {
		// If the query fails, wrap the error with context
		return nil, fmt.Errorf("failed to get expenses: %w", err)
	}

	// Step 6: Return the results
	return expenses, nil
}

// Update modifies an existing expense
// This method implements the domain.Repository.Update interface
func (r *Repository) Update(ctx context.Context, expense *domain.Expense) error {
	// Use GORM's Save method to update the expense in the database
	// Save() automatically handles the SQL UPDATE statement
	// It updates all fields of the expense
	return r.db.WithContext(ctx).Save(expense).Error
}

// Delete removes an expense by its ID
// This method implements the domain.Repository.Delete interface
func (r *Repository) Delete(ctx context.Context, id string) error {
	// Step 1: Parse the string ID into a UUID
	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}

	// Step 2: Execute the delete operation
	// Where("id = ?", uuid) - filters to delete only the specific expense
	// Delete(&domain.Expense{}) - deletes records matching the WHERE clause
	// The empty struct is just a placeholder to tell GORM which table to delete from
	result := r.db.WithContext(ctx).Where("id = ?", uuid).Delete(&domain.Expense{})

	// Step 3: Check for database errors
	if result.Error != nil {
		return fmt.Errorf("failed to delete expense: %w", result.Error)
	}

	// Step 4: Check if any records were actually deleted
	// RowsAffected tells us how many rows were deleted
	if result.RowsAffected == 0 {
		// If no rows were deleted, the expense didn't exist
		return domain.ErrExpenseNotFound
	}

	// Step 5: Return nil to indicate success
	return nil
}

// Exists checks if an expense with the given ID exists
// This method implements the domain.Repository.Exists interface
func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	// Step 1: Parse the string ID into a UUID
	uuid, err := uuid.Parse(id)
	if err != nil {
		return false, fmt.Errorf("invalid UUID format: %w", err)
	}

	// Step 2: Count records with the given ID
	// Model(&domain.Expense{}) - tells GORM which table to query
	// Where("id = ?", uuid) - filters by the specific ID
	// Count(&count) - counts matching records and stores result in count
	var count int64
	if err := r.db.WithContext(ctx).Model(&domain.Expense{}).Where("id = ?", uuid).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check expense existence: %w", err)
	}

	// Step 3: Return true if count > 0, false otherwise
	return count > 0, nil
}

// AutoMigrate runs database migrations
// This method creates the database table if it doesn't exist
// It's not part of the repository interface, but a utility method
func (r *Repository) AutoMigrate() error {
	// GORM's AutoMigrate automatically creates tables based on struct definitions
	// It also adds missing columns and indexes
	return r.db.AutoMigrate(&domain.Expense{})
}
