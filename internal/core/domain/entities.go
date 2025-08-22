package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of transaction (income or expense)
type TransactionType string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
)

// Transaction represents a financial transaction in the system
type Transaction struct {
	ID          string          `json:"id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	CategoryID  int             `json:"category_id"`
	Description string          `json:"description"`
	Date        time.Time       `json:"date"`
	CreatedAt   time.Time       `json:"created_at"`
}

// Category represents a transaction category
type Category struct {
	ID   int             `json:"id"`
	Name string          `json:"name"`
	Type TransactionType `json:"type"`
}

// NewTransaction creates a new transaction with generated ID and current timestamp
func NewTransaction(transactionType TransactionType, amount float64, categoryID int, description string, date time.Time) *Transaction {
	return &Transaction{
		ID:          uuid.New().String(),
		Type:        transactionType,
		Amount:      amount,
		CategoryID:  categoryID,
		Description: strings.TrimSpace(description),
		Date:        date,
		CreatedAt:   time.Now(),
	}
}

// NewCategory creates a new category
func NewCategory(name string, categoryType TransactionType) *Category {
	return &Category{
		Name: strings.TrimSpace(name),
		Type: categoryType,
	}
}

// Validate validates the transaction fields
func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if t.Type != TypeIncome && t.Type != TypeExpense {
		return errors.New("transaction type must be either 'income' or 'expense'")
	}

	if t.CategoryID <= 0 {
		return errors.New("category ID must be positive")
	}

	if len(t.Description) > 200 {
		return errors.New("description must be 200 characters or less")
	}

	if t.Date.IsZero() {
		return errors.New("date cannot be zero")
	}

	if t.Date.After(time.Now()) {
		return errors.New("date cannot be in the future")
	}

	return nil
}

// Validate validates the category fields
func (c *Category) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("category name cannot be empty")
	}

	if len(c.Name) > 50 {
		return errors.New("category name must be 50 characters or less")
	}

	if c.Type != TypeIncome && c.Type != TypeExpense {
		return errors.New("category type must be either 'income' or 'expense'")
	}

	return nil
}

// IsIncome returns true if the transaction is an income
func (t *Transaction) IsIncome() bool {
	return t.Type == TypeIncome
}

// IsExpense returns true if the transaction is an expense
func (t *Transaction) IsExpense() bool {
	return t.Type == TypeExpense
}

// FormatAmount returns the amount as a formatted string with currency
func (t *Transaction) FormatAmount() string {
	return formatCurrency(t.Amount)
}

// FormatDate returns the date as a formatted string
func (t *Transaction) FormatDate() string {
	return t.Date.Format("2006-01-02")
}

// formatCurrency formats a float64 as currency
func formatCurrency(amount float64) string {
	return fmt.Sprintf("$%.2f", amount)
}

// DefaultExpenseCategories returns the default expense categories
func DefaultExpenseCategories() []Category {
	return []Category{
		{Name: "Food & Dining", Type: TypeExpense},
		{Name: "Transportation", Type: TypeExpense},
		{Name: "Shopping", Type: TypeExpense},
		{Name: "Entertainment", Type: TypeExpense},
		{Name: "Bills & Utilities", Type: TypeExpense},
		{Name: "Healthcare", Type: TypeExpense},
		{Name: "Other", Type: TypeExpense},
	}
}

// DefaultIncomeCategories returns the default income categories
func DefaultIncomeCategories() []Category {
	return []Category{
		{Name: "Salary", Type: TypeIncome},
		{Name: "Gifts", Type: TypeIncome},
		{Name: "Investments", Type: TypeIncome},
		{Name: "Other", Type: TypeIncome},
	}
}