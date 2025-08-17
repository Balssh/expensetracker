package usecase

import (
	"context"
	"fmt"
	"time"

	"expense-tracker/internal/core/domain"
)

type ExpenseUseCase struct {
	expenseRepo  ExpenseRepository
	categoryRepo CategoryRepository
}

func NewExpenseUseCase(expenseRepo ExpenseRepository, categoryRepo CategoryRepository) *ExpenseUseCase {
	return &ExpenseUseCase{
		expenseRepo:  expenseRepo,
		categoryRepo: categoryRepo,
	}
}

func (uc *ExpenseUseCase) AddExpense(ctx context.Context, description string, amount float64, date time.Time, categoryID int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if description == "" {
		return fmt.Errorf("description cannot be empty")
	}

	category, err := uc.categoryRepo.GetExpenseCategoryByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("invalid category: %w", err)
	}

	expense := &domain.Expense{
		Description: description,
		Amount:      amount,
		Date:        date,
		CategoryID:  categoryID,
		Category:    category,
	}

	return uc.expenseRepo.Create(ctx, expense)
}

func (uc *ExpenseUseCase) GetExpenses(ctx context.Context, offset, limit int) ([]*domain.Expense, error) {
	return uc.expenseRepo.GetAll(ctx, offset, limit)
}

func (uc *ExpenseUseCase) GetExpensesByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Expense, error) {
	return uc.expenseRepo.GetByDateRange(ctx, start, end)
}

func (uc *ExpenseUseCase) GetTotalExpensesByDateRange(ctx context.Context, start, end time.Time) (float64, error) {
	return uc.expenseRepo.GetTotalByDateRange(ctx, start, end)
}

func (uc *ExpenseUseCase) GetExpenseCategories(ctx context.Context) ([]*domain.Category, error) {
	return uc.categoryRepo.GetExpenseCategories(ctx)
}

func (uc *ExpenseUseCase) AddExpenseCategory(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("category name cannot be empty")
	}

	category := &domain.Category{
		Name: name,
	}

	return uc.categoryRepo.CreateExpenseCategory(ctx, category)
}