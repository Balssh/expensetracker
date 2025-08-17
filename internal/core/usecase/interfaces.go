package usecase

import (
	"context"
	"time"

	"expense-tracker/internal/core/domain"
)

type ExpenseRepository interface {
	Create(ctx context.Context, expense *domain.Expense) error
	GetByID(ctx context.Context, id int) (*domain.Expense, error)
	GetAll(ctx context.Context, offset, limit int) ([]*domain.Expense, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Expense, error)
	GetTotalByDateRange(ctx context.Context, start, end time.Time) (float64, error)
	Update(ctx context.Context, expense *domain.Expense) error
	Delete(ctx context.Context, id int) error
}

type IncomeRepository interface {
	Create(ctx context.Context, income *domain.Income) error
	GetByID(ctx context.Context, id int) (*domain.Income, error)
	GetAll(ctx context.Context, offset, limit int) ([]*domain.Income, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Income, error)
	GetTotalByDateRange(ctx context.Context, start, end time.Time) (float64, error)
	Update(ctx context.Context, income *domain.Income) error
	Delete(ctx context.Context, id int) error
}

type CategoryRepository interface {
	CreateExpenseCategory(ctx context.Context, category *domain.Category) error
	CreateIncomeCategory(ctx context.Context, category *domain.Category) error
	GetExpenseCategories(ctx context.Context) ([]*domain.Category, error)
	GetIncomeCategories(ctx context.Context) ([]*domain.Category, error)
	GetExpenseCategoryByID(ctx context.Context, id int) (*domain.Category, error)
	GetIncomeCategoryByID(ctx context.Context, id int) (*domain.Category, error)
}

type TransactionRepository interface {
	GetRecentTransactions(ctx context.Context, limit int) ([]*domain.Transaction, error)
	GetAllTransactions(ctx context.Context, offset, limit int) ([]*domain.Transaction, error)
	SearchTransactions(ctx context.Context, query string, offset, limit int) ([]*domain.Transaction, error)
}