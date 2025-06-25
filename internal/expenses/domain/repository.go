// Package domain contains the core business logic and entities
// This file defines the repository interface - a contract for data access operations
package domain

import (
	"context" // Go's package for handling request context (cancellation, timeouts, etc.)
)

// Repository defines the interface for expense data operations
// This is an "interface" - it defines what methods must be implemented, but not how
// This follows the "Interface Segregation Principle" - keep interfaces small and focused
// The repository pattern abstracts data access from business logic
type Repository interface {
	// Create adds a new expense to the repository (database)
	// ctx is the context for this operation (allows cancellation, timeouts)
	// expense is a pointer to the expense we want to save
	// Returns an error if the operation fails
	Create(ctx context.Context, expense *Expense) error

	// GetByID retrieves an expense by its unique identifier
	// ctx is the context for this operation
	// id is the string representation of the expense's UUID
	// Returns a pointer to the expense if found, or an error if not found/failed
	GetByID(ctx context.Context, id string) (*Expense, error)

	// GetAll retrieves all expenses with optional filtering
	// ctx is the context for this operation
	// filters is a map of filter criteria (e.g., {"category": "Food", "min_amount": 10.0})
	// Returns a slice of expense pointers and an error if the operation fails
	// A slice is Go's dynamic array type (like ArrayList in Java)
	GetAll(ctx context.Context, filters map[string]interface{}) ([]*Expense, error)

	// Update modifies an existing expense in the repository
	// ctx is the context for this operation
	// expense is a pointer to the expense with updated values
	// Returns an error if the operation fails
	Update(ctx context.Context, expense *Expense) error

	// Delete removes an expense from the repository by its ID
	// ctx is the context for this operation
	// id is the string representation of the expense's UUID
	// Returns an error if the operation fails
	Delete(ctx context.Context, id string) error

	// Exists checks if an expense with the given ID exists in the repository
	// ctx is the context for this operation
	// id is the string representation of the expense's UUID
	// Returns true if the expense exists, false if not, and an error if the operation fails
	// This is useful for validation before performing operations
	Exists(ctx context.Context, id string) (bool, error)
}
