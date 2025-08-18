package usecase

import (
	"context"
	"time"

	"expense-tracker/internal/core/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) error
	GetByID(ctx context.Context, id int) (*domain.Transaction, error)
	GetAll(ctx context.Context, offset, limit int) ([]*domain.Transaction, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Transaction, error)
	GetByType(ctx context.Context, transactionType string, offset, limit int) ([]*domain.Transaction, error)
	GetTotalByDateRange(ctx context.Context, start, end time.Time, transactionType string) (float64, error)
	GetRecentTransactions(ctx context.Context, limit int) ([]*domain.Transaction, error)
	SearchTransactions(ctx context.Context, query string, offset, limit int) ([]*domain.Transaction, error)
	Update(ctx context.Context, transaction *domain.Transaction) error
	Delete(ctx context.Context, id int) error
}

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *domain.Category, categoryType string) error
	GetCategories(ctx context.Context, categoryType string) ([]*domain.Category, error)
	GetCategoryByID(ctx context.Context, id int, categoryType string) (*domain.Category, error)
}