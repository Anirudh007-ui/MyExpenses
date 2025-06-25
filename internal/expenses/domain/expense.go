// Package domain contains the core business logic and entities
// This is the innermost layer of Clean Architecture - it has no dependencies on other layers
package domain

import (
	"time" // Package for handling dates and times

	"github.com/google/uuid" // Package for generating unique identifiers (UUIDs)
)

// Expense represents the core business entity for an expense
// This is the main data structure that represents an expense in our system
// It contains all the fields that define what an expense is in our business domain
type Expense struct {
	// ID is a unique identifier for each expense
	// uuid.UUID is a type that represents a universally unique identifier
	// The tags below provide metadata for JSON serialization and database mapping
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`

	// Description is what the expense was for (e.g., "Coffee", "Gas", "Groceries")
	// string is Go's built-in type for text
	// gorm:"not null" means this field cannot be empty in the database
	Description string `json:"description" gorm:"not null"`

	// Amount is how much the expense cost
	// float64 is Go's type for decimal numbers (64-bit precision)
	// This allows us to store amounts like 12.99, 100.50, etc.
	Amount float64 `json:"amount" gorm:"not null"`

	// Category helps organize expenses (e.g., "Food", "Transportation", "Entertainment")
	Category string `json:"category" gorm:"not null"`

	// Date is when the expense occurred
	// time.Time is Go's type for representing dates and times
	Date time.Time `json:"date" gorm:"not null"`

	// CreatedAt is automatically set when the expense is first saved to the database
	// gorm:"autoCreateTime" tells GORM to automatically set this field
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`

	// UpdatedAt is automatically updated whenever the expense is modified
	// gorm:"autoUpdateTime" tells GORM to automatically update this field
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// NewExpense creates a new expense with validation
// This is a "factory function" - it ensures that all expenses are created with valid data
// It returns a pointer to Expense (*Expense) and an error
// The * means it's a pointer - a reference to the actual data in memory
func NewExpense(description string, amount float64, category string, date time.Time) (*Expense, error) {
	// Validation: Check if description is empty
	// In Go, "" represents an empty string
	if description == "" {
		// Return nil (no expense) and an error indicating the problem
		// nil means "no value" or "empty"
		return nil, ErrInvalidDescription
	}

	// Validation: Check if amount is less than or equal to 0
	// We don't want negative or zero amounts for expenses
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Validation: Check if category is empty
	if category == "" {
		return nil, ErrInvalidCategory
	}

	// Validation: Check if date is "zero" (uninitialized)
	// time.Time{} creates a zero time value
	if date.IsZero() {
		return nil, ErrInvalidDate
	}

	// If all validations pass, create and return a new expense
	// &Expense{...} creates a new Expense struct and returns a pointer to it
	return &Expense{
		ID:          uuid.New(),  // Generate a new unique ID
		Description: description, // Set the description
		Amount:      amount,      // Set the amount
		Category:    category,    // Set the category
		Date:        date,        // Set the date
		// Note: CreatedAt and UpdatedAt will be set automatically by GORM
	}, nil
}

// Validate checks if the expense is valid
// This method can be called on an existing expense to verify it's still valid
// It returns an error if the expense is invalid, or nil if it's valid
func (e *Expense) Validate() error {
	// The (e *Expense) part means this method belongs to the Expense struct
	// 'e' is a reference to the current expense instance (like 'this' in other languages)

	// Check if description is empty
	if e.Description == "" {
		return ErrInvalidDescription
	}

	// Check if amount is invalid
	if e.Amount <= 0 {
		return ErrInvalidAmount
	}

	// Check if category is empty
	if e.Category == "" {
		return ErrInvalidCategory
	}

	// Check if date is zero (uninitialized)
	if e.Date.IsZero() {
		return ErrInvalidDate
	}

	// If we get here, all validations passed
	// Return nil to indicate no error
	return nil
}

// Update updates the expense fields
// This method allows partial updates - only the provided fields will be changed
// It takes the new values as parameters and only updates non-empty/non-zero values
func (e *Expense) Update(description string, amount float64, category string, date time.Time) error {
	// Update description only if a new one is provided (not empty)
	if description != "" {
		e.Description = description
	}

	// Update amount only if a valid new amount is provided (greater than 0)
	if amount > 0 {
		e.Amount = amount
	}

	// Update category only if a new one is provided (not empty)
	if category != "" {
		e.Category = category
	}

	// Update date only if a valid new date is provided (not zero)
	if !date.IsZero() {
		e.Date = date
	}

	// After updating, validate the expense to ensure it's still valid
	return e.Validate()
}
