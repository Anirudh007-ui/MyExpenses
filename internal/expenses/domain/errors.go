// Package domain contains the core business logic and entities
// This file defines all the domain-specific errors that can occur in our business logic
package domain

import "errors" // Go's built-in package for creating and handling errors

// Domain errors are defined as package-level variables
// These errors represent business rule violations and domain-specific problems
// Using errors.New() creates a new error with a descriptive message
var (
	// ErrInvalidDescription occurs when trying to create an expense with an empty description
	// This enforces the business rule that every expense must have a description
	ErrInvalidDescription = errors.New("invalid description: cannot be empty")

	// ErrInvalidAmount occurs when trying to create an expense with an invalid amount
	// This enforces the business rule that expenses must have positive amounts
	ErrInvalidAmount = errors.New("invalid amount: must be greater than 0")

	// ErrInvalidCategory occurs when trying to create an expense with an empty category
	// This enforces the business rule that every expense must be categorized
	ErrInvalidCategory = errors.New("invalid category: cannot be empty")

	// ErrInvalidDate occurs when trying to create an expense with an invalid date
	// This enforces the business rule that every expense must have a valid date
	ErrInvalidDate = errors.New("invalid date: cannot be zero")

	// ErrExpenseNotFound occurs when trying to access an expense that doesn't exist
	// This is used when the database cannot find an expense with the given ID
	ErrExpenseNotFound = errors.New("expense not found")

	// ErrExpenseExists occurs when trying to create an expense that already exists
	// This prevents duplicate expenses (though not currently used in this implementation)
	ErrExpenseExists = errors.New("expense already exists")
)
