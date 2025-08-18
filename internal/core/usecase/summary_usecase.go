package usecase

import (
	"context"
	"time"

	"expense-tracker/internal/core/domain"
)

type SummaryUseCase struct {
	transactionRepo TransactionRepository
}

func NewSummaryUseCase(transactionRepo TransactionRepository) *SummaryUseCase {
	return &SummaryUseCase{
		transactionRepo: transactionRepo,
	}
}

func (uc *SummaryUseCase) GetMonthlySummary(ctx context.Context, year int, month time.Month) (*domain.Summary, error) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	totalIncome, err := uc.transactionRepo.GetTotalByDateRange(ctx, start, end, "income")
	if err != nil {
		return nil, err
	}

	totalExpense, err := uc.transactionRepo.GetTotalByDateRange(ctx, start, end, "expense")
	if err != nil {
		return nil, err
	}

	return &domain.Summary{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetBalance:   totalIncome - totalExpense,
	}, nil
}

func (uc *SummaryUseCase) GetRecentTransactions(ctx context.Context, limit int) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetRecentTransactions(ctx, limit)
}

func (uc *SummaryUseCase) GetAllTransactions(ctx context.Context, offset, limit int) ([]*domain.Transaction, error) {
	return uc.transactionRepo.GetAll(ctx, offset, limit)
}

func (uc *SummaryUseCase) SearchTransactions(ctx context.Context, query string, offset, limit int) ([]*domain.Transaction, error) {
	return uc.transactionRepo.SearchTransactions(ctx, query, offset, limit)
}