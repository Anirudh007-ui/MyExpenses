// Package application contains the business logic and use cases
// This layer orchestrates the domain logic and coordinates between different parts of the system
// It sits between the domain layer (business rules) and infrastructure layer (external concerns)
package application

import (
	"context" // For request context (cancellation, timeouts)
	"fmt"     // For formatted string operations and error wrapping
	"time"    // For handling dates and times

	"myexpenses/internal/expenses/domain" // Import our domain layer
)

// Service handles business logic for expenses
// This is the main business logic layer that coordinates between domain and infrastructure
// It implements the "Application Service" pattern from Domain-Driven Design
type Service struct {
	// repo is a dependency on the repository interface
	// This follows the Dependency Inversion Principle - depend on abstractions, not concretions
	// The actual implementation (PostgreSQL, in-memory, etc.) is injected later
	repo domain.Repository
}

// NewService creates a new expense service
// This is a constructor function that implements dependency injection
// It takes a repository implementation and returns a configured service
func NewService(repo domain.Repository) *Service {
	return &Service{
		repo: repo, // Store the repository dependency
	}
}

// CreateExpenseRequest represents the request to create an expense
// This is a DTO (Data Transfer Object) - it defines the contract for creating expenses
// It's separate from the domain model to allow for API-specific validation and flexibility
type CreateExpenseRequest struct {
	// Description is what the expense was for
	// binding:"required" is a Gin validation tag that ensures this field is provided
	Description string `json:"description" binding:"required"`

	// Amount is how much the expense cost
	// binding:"required,gt=0" ensures the amount is provided and greater than 0
	Amount float64 `json:"amount" binding:"required,gt=0"`

	// Category helps organize the expense
	Category string `json:"category" binding:"required"`

	// Date is when the expense occurred
	Date time.Time `json:"date" binding:"required"`
}

// UpdateExpenseRequest represents the request to update an expense
// Note that these fields are not marked as "required" because updates can be partial
// A client can update just the amount without changing other fields
type UpdateExpenseRequest struct {
	// All fields are optional for updates
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
}

// CreateExpense creates a new expense
// This is a use case - it represents a specific business operation
// It orchestrates the creation process: validation -> domain object creation -> persistence
func (s *Service) CreateExpense(ctx context.Context, req *CreateExpenseRequest) (*domain.Expense, error) {
	// Step 1: Create a domain object using the factory function
	// This ensures all business rules are enforced
	expense, err := domain.NewExpense(req.Description, req.Amount, req.Category, req.Date)
	if err != nil {
		// If domain validation fails, wrap the error with context
		// %w is the error wrapping verb - it preserves the original error
		return nil, fmt.Errorf("failed to create expense: %w", err)
	}

	// Step 2: Save the expense to the repository (database)
	if err := s.repo.Create(ctx, expense); err != nil {
		// If persistence fails, wrap the error with context
		return nil, fmt.Errorf("failed to save expense: %w", err)
	}

	// Step 3: Return the created expense
	return expense, nil
}

// GetExpense retrieves an expense by ID
// This is a simple query use case
func (s *Service) GetExpense(ctx context.Context, id string) (*domain.Expense, error) {
	// Delegate to the repository to fetch the expense
	expense, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Wrap any errors with context
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}
	return expense, nil
}

// GetAllExpenses retrieves all expenses with optional filtering
// This is a query use case that supports filtering
func (s *Service) GetAllExpenses(ctx context.Context, filters map[string]interface{}) ([]*domain.Expense, error) {
	// Delegate to the repository to fetch expenses with filters
	expenses, err := s.repo.GetAll(ctx, filters)
	if err != nil {
		// Wrap any errors with context
		return nil, fmt.Errorf("failed to get expenses: %w", err)
	}
	return expenses, nil
}

// UpdateExpense updates an existing expense
// This is a complex use case that involves validation and coordination
func (s *Service) UpdateExpense(ctx context.Context, id string, req *UpdateExpenseRequest) (*domain.Expense, error) {
	// Step 1: Check if the expense exists before trying to update it
	// This prevents errors and provides better user feedback
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check expense existence: %w", err)
	}
	if !exists {
		// Return a domain-specific error if the expense doesn't exist
		return nil, domain.ErrExpenseNotFound
	}

	// Step 2: Get the current expense from the repository
	expense, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	// Step 3: Update the expense fields using the domain method
	// This ensures business rules are still enforced during updates
	if err := expense.Update(req.Description, req.Amount, req.Category, req.Date); err != nil {
		return nil, fmt.Errorf("failed to update expense: %w", err)
	}

	// Step 4: Save the updated expense back to the repository
	if err := s.repo.Update(ctx, expense); err != nil {
		return nil, fmt.Errorf("failed to save updated expense: %w", err)
	}

	// Step 5: Return the updated expense
	return expense, nil
}

// DeleteExpense removes an expense
// This is a simple command use case
func (s *Service) DeleteExpense(ctx context.Context, id string) error {
	// Step 1: Check if the expense exists before trying to delete it
	exists, err := s.repo.Exists(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check expense existence: %w", err)
	}
	if !exists {
		// Return a domain-specific error if the expense doesn't exist
		return domain.ErrExpenseNotFound
	}

	// Step 2: Delete the expense from the repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	// Step 3: Return nil to indicate success
	return nil
}
