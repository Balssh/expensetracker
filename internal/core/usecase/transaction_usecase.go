package usecase

import (
	"context"
	"fmt"
	"time"

	"expense-tracker/internal/core/domain"
)

type TransactionUseCase struct {
	transactionRepo TransactionRepository
	categoryRepo    CategoryRepository
}

func NewTransactionUseCase(transactionRepo TransactionRepository, categoryRepo CategoryRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
	}
}

func (uc *TransactionUseCase) AddTransaction(ctx context.Context, transaction *domain.Transaction) error {
	if transaction.Amount <= 0 {
		return fmt.Errorf("transaction amount must be positive")
	}

	if transaction.Description == "" {
		return fmt.Errorf("transaction description is required")
	}

	if transaction.Type != "income" && transaction.Type != "expense" {
		return fmt.Errorf("transaction type must be 'income' or 'expense'")
	}

	if transaction.Date.IsZero() {
		transaction.Date = time.Now()
	}

	if transaction.Category != nil && transaction.Category.ID > 0 {
		category, err := uc.categoryRepo.GetCategoryByID(ctx, transaction.Category.ID, transaction.Type)
		if err != nil {
			return fmt.Errorf("invalid category: %w", err)
		}
		transaction.Category = category
	}

	return uc.transactionRepo.Create(ctx, transaction)
}

func (uc *TransactionUseCase) GetTransactionByID(ctx context.Context, id int) (*domain.Transaction, error) {
	return uc.transactionRepo.GetByID(ctx, id)
}

func (uc *TransactionUseCase) GetAllTransactions(ctx context.Context, offset, limit int) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetAll(ctx, offset, limit)
}

func (uc *TransactionUseCase) GetTransactionsByType(ctx context.Context, transactionType string, offset, limit int) ([]*domain.Transaction, error) {
	if transactionType != "income" && transactionType != "expense" {
		return nil, fmt.Errorf("transaction type must be 'income' or 'expense'")
	}
	return uc.transactionRepo.GetByType(ctx, transactionType, offset, limit)
}

func (uc *TransactionUseCase) GetTransactionsByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetByDateRange(ctx, start, end)
}

func (uc *TransactionUseCase) GetRecentTransactions(ctx context.Context, limit int) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetRecentTransactions(ctx, limit)
}

func (uc *TransactionUseCase) SearchTransactions(ctx context.Context, query string, offset, limit int) ([]*domain.Transaction, error) {
	return uc.transactionRepo.SearchTransactions(ctx, query, offset, limit)
}

func (uc *TransactionUseCase) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	if transaction.ID <= 0 {
		return fmt.Errorf("transaction ID is required for update")
	}

	if transaction.Amount <= 0 {
		return fmt.Errorf("transaction amount must be positive")
	}

	if transaction.Description == "" {
		return fmt.Errorf("transaction description is required")
	}

	if transaction.Type != "income" && transaction.Type != "expense" {
		return fmt.Errorf("transaction type must be 'income' or 'expense'")
	}

	if transaction.Category != nil && transaction.Category.ID > 0 {
		category, err := uc.categoryRepo.GetCategoryByID(ctx, transaction.Category.ID, transaction.Type)
		if err != nil {
			return fmt.Errorf("invalid category: %w", err)
		}
		transaction.Category = category
	}

	return uc.transactionRepo.Update(ctx, transaction)
}

func (uc *TransactionUseCase) DeleteTransaction(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("transaction ID is required for delete")
	}
	return uc.transactionRepo.Delete(ctx, id)
}

func (uc *TransactionUseCase) GetCategories(ctx context.Context, transactionType string) ([]*domain.Category, error) {
	if transactionType != "income" && transactionType != "expense" {
		return nil, fmt.Errorf("transaction type must be 'income' or 'expense'")
	}
	return uc.categoryRepo.GetCategories(ctx, transactionType)
}