package usecase

import (
	"context"
	"fmt"
	"time"

	"expense-tracker/internal/core/domain"
)

type IncomeUseCase struct {
	incomeRepo   IncomeRepository
	categoryRepo CategoryRepository
}

func NewIncomeUseCase(incomeRepo IncomeRepository, categoryRepo CategoryRepository) *IncomeUseCase {
	return &IncomeUseCase{
		incomeRepo:   incomeRepo,
		categoryRepo: categoryRepo,
	}
}

func (uc *IncomeUseCase) AddIncome(ctx context.Context, description string, amount float64, date time.Time, categoryID int) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if description == "" {
		return fmt.Errorf("description cannot be empty")
	}

	category, err := uc.categoryRepo.GetIncomeCategoryByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("invalid category: %w", err)
	}

	income := &domain.Income{
		Description: description,
		Amount:      amount,
		Date:        date,
		CategoryID:  categoryID,
		Category:    category,
	}

	return uc.incomeRepo.Create(ctx, income)
}

func (uc *IncomeUseCase) GetIncomes(ctx context.Context, offset, limit int) ([]*domain.Income, error) {
	return uc.incomeRepo.GetAll(ctx, offset, limit)
}

func (uc *IncomeUseCase) GetIncomesByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Income, error) {
	return uc.incomeRepo.GetByDateRange(ctx, start, end)
}

func (uc *IncomeUseCase) GetTotalIncomesByDateRange(ctx context.Context, start, end time.Time) (float64, error) {
	return uc.incomeRepo.GetTotalByDateRange(ctx, start, end)
}

func (uc *IncomeUseCase) GetIncomeCategories(ctx context.Context) ([]*domain.Category, error) {
	return uc.categoryRepo.GetIncomeCategories(ctx)
}

func (uc *IncomeUseCase) AddIncomeCategory(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("category name cannot be empty")
	}

	category := &domain.Category{
		Name: name,
	}

	return uc.categoryRepo.CreateIncomeCategory(ctx, category)
}