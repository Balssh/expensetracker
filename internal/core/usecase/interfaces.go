package usecase

import (
	"fmt"
	"time"

	"github.com/yourusername/expense-tracker/internal/core/domain"
)

// TransactionRepository defines the interface for transaction persistence operations
type TransactionRepository interface {
	// Create creates a new transaction
	Create(transaction *domain.Transaction) error
	
	// GetByID retrieves a transaction by its ID
	GetByID(id string) (*domain.Transaction, error)
	
	// List retrieves transactions with optional limit and offset for pagination
	List(limit, offset int) ([]*domain.Transaction, error)
	
	// ListByType retrieves transactions of a specific type
	ListByType(transactionType domain.TransactionType, limit, offset int) ([]*domain.Transaction, error)
	
	// ListByDateRange retrieves transactions within a date range
	ListByDateRange(start, end time.Time) ([]*domain.Transaction, error)
	
	// ListByCategory retrieves transactions for a specific category
	ListByCategory(categoryID int) ([]*domain.Transaction, error)
	
	// Update updates an existing transaction
	Update(transaction *domain.Transaction) error
	
	// Delete deletes a transaction by its ID
	Delete(id string) error
	
	// Count returns the total number of transactions
	Count() (int, error)
	
	// CountByType returns the number of transactions of a specific type
	CountByType(transactionType domain.TransactionType) (int, error)
}

// CategoryRepository defines the interface for category persistence operations
type CategoryRepository interface {
	// Create creates a new category
	Create(category *domain.Category) error
	
	// GetByID retrieves a category by its ID
	GetByID(id int) (*domain.Category, error)
	
	// GetByName retrieves a category by its name and type
	GetByName(name string, categoryType domain.TransactionType) (*domain.Category, error)
	
	// ListByType retrieves all categories of a specific type
	ListByType(categoryType domain.TransactionType) ([]*domain.Category, error)
	
	// List retrieves all categories
	List() ([]*domain.Category, error)
	
	// Update updates an existing category
	Update(category *domain.Category) error
	
	// Delete deletes a category by its ID
	Delete(id int) error
	
	// GetUsageCount returns the number of transactions using this category
	GetUsageCount(id int) (int, error)
	
	// Exists checks if a category with the given name and type already exists
	Exists(name string, categoryType domain.TransactionType) (bool, error)
}

// RepositoryError represents errors that can occur in repository operations
type RepositoryError struct {
	Op  string // Operation that failed
	Err error  // Underlying error
}

func (e *RepositoryError) Error() string {
	return e.Op + ": " + e.Err.Error()
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// Common repository errors
var (
	ErrNotFound      = &RepositoryError{Op: "repository", Err: fmt.Errorf("record not found")}
	ErrDuplicate     = &RepositoryError{Op: "repository", Err: fmt.Errorf("duplicate record")}
	ErrInvalidInput  = &RepositoryError{Op: "repository", Err: fmt.Errorf("invalid input")}
	ErrDatabase      = &RepositoryError{Op: "repository", Err: fmt.Errorf("database error")}
)