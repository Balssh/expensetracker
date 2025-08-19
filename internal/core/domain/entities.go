package domain

import (
	"fmt"
	"strings"
	"time"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *Category) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	if len(c.Name) > 50 {
		return fmt.Errorf("category name cannot exceed 50 characters")
	}
	return nil
}

type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	NetBalance   float64 `json:"net_balance"`
}

func NewSummary(totalIncome, totalExpense float64) *Summary {
	return &Summary{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetBalance:   totalIncome - totalExpense,
	}
}

type Transaction struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"` // "income" or "expense"
	Category    *Category `json:"category,omitempty"`
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return fmt.Errorf("transaction amount must be positive")
	}

	if strings.TrimSpace(t.Description) == "" {
		return fmt.Errorf("transaction description cannot be empty")
	}

	if len(t.Description) > 200 {
		return fmt.Errorf("transaction description cannot exceed 200 characters")
	}

	if t.Type != "income" && t.Type != "expense" {
		return fmt.Errorf("transaction type must be 'income' or 'expense'")
	}

	if t.Date.IsZero() {
		return fmt.Errorf("transaction date cannot be zero")
	}

	if t.Category != nil {
		if err := t.Category.Validate(); err != nil {
			return fmt.Errorf("invalid category: %w", err)
		}
	}

	return nil
}

func (t *Transaction) IsIncome() bool {
	return t.Type == "income"
}

func (t *Transaction) IsExpense() bool {
	return t.Type == "expense"
}
