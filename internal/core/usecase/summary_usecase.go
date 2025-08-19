package usecase

import (
	"context"
	"fmt"
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

// GetSummary provides intelligent defaults for period-based summaries
func (uc *SummaryUseCase) GetSummary(ctx context.Context, periodType ...domain.PeriodType) (*domain.Summary, error) {
	period := domain.PeriodTypeMonth // Default to current month
	if len(periodType) > 0 {
		period = periodType[0]
	}
	
	var dateRange *domain.DateRange
	switch period {
	case domain.PeriodTypeWeek:
		dateRange = domain.GetCurrentWeekRange()
	case domain.PeriodTypeMonth:
		dateRange = domain.GetCurrentMonthRange()
	case domain.PeriodTypeQuarter:
		dateRange = domain.GetCurrentQuarterRange()
	case domain.PeriodTypeYear:
		dateRange = domain.GetCurrentYearRange()
	default:
		return nil, fmt.Errorf("unsupported period type: %s", period)
	}
	
	return uc.GetSummaryByDateRange(ctx, dateRange.Start, dateRange.End, period)
}

// GetWeeklySummary gets summary for a specific week
func (uc *SummaryUseCase) GetWeeklySummary(ctx context.Context, year, week int) (*domain.Summary, error) {
	dateRange := domain.GetWeekRange(year, week)
	if dateRange == nil {
		return nil, fmt.Errorf("invalid week: %d/%d", year, week)
	}
	return uc.GetSummaryByDateRange(ctx, dateRange.Start, dateRange.End, domain.PeriodTypeWeek)
}

// GetQuarterlySummary gets summary for a specific quarter
func (uc *SummaryUseCase) GetQuarterlySummary(ctx context.Context, year, quarter int) (*domain.Summary, error) {
	dateRange := domain.GetQuarterRange(year, quarter)
	if dateRange == nil {
		return nil, fmt.Errorf("invalid quarter: %d/Q%d", year, quarter)
	}
	return uc.GetSummaryByDateRange(ctx, dateRange.Start, dateRange.End, domain.PeriodTypeQuarter)
}

// GetYearlySummary gets summary for a specific year
func (uc *SummaryUseCase) GetYearlySummary(ctx context.Context, year int) (*domain.Summary, error) {
	dateRange := domain.GetYearRange(year)
	return uc.GetSummaryByDateRange(ctx, dateRange.Start, dateRange.End, domain.PeriodTypeYear)
}

// GetCustomSummary gets summary for a custom date range
func (uc *SummaryUseCase) GetCustomSummary(ctx context.Context, start, end time.Time) (*domain.Summary, error) {
	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}
	return uc.GetSummaryByDateRange(ctx, start, end, domain.PeriodTypeCustom)
}

// GetSummaryByDateRange is the core method that builds enhanced summaries
func (uc *SummaryUseCase) GetSummaryByDateRange(ctx context.Context, start, end time.Time, periodType domain.PeriodType) (*domain.Summary, error) {
	// Get basic totals
	totalIncome, err := uc.transactionRepo.GetTotalByDateRange(ctx, start, end, "income")
	if err != nil {
		return nil, fmt.Errorf("failed to get income total: %w", err)
	}

	totalExpense, err := uc.transactionRepo.GetTotalByDateRange(ctx, start, end, "expense")
	if err != nil {
		return nil, fmt.Errorf("failed to get expense total: %w", err)
	}

	// Create enhanced summary
	dateRange := domain.NewDateRange(start, end)
	summary := domain.NewEnhancedSummary(totalIncome, totalExpense, periodType, dateRange)

	// Get category breakdowns
	incomeBreakdown, err := uc.transactionRepo.GetCategoryTotalsByDateRange(ctx, start, end, "income")
	if err != nil {
		return nil, fmt.Errorf("failed to get income breakdown: %w", err)
	}

	expenseBreakdown, err := uc.transactionRepo.GetCategoryTotalsByDateRange(ctx, start, end, "expense")
	if err != nil {
		return nil, fmt.Errorf("failed to get expense breakdown: %w", err)
	}

	summary.SetCategoryBreakdowns(incomeBreakdown, expenseBreakdown)

	// Get transaction counts
	totalCount, err := uc.transactionRepo.GetTransactionCountByDateRange(ctx, start, end, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get total transaction count: %w", err)
	}

	incomeCount, err := uc.transactionRepo.GetTransactionCountByDateRange(ctx, start, end, "income")
	if err != nil {
		return nil, fmt.Errorf("failed to get income transaction count: %w", err)
	}

	expenseCount, err := uc.transactionRepo.GetTransactionCountByDateRange(ctx, start, end, "expense")
	if err != nil {
		return nil, fmt.Errorf("failed to get expense transaction count: %w", err)
	}

	summary.SetTransactionCounts(totalCount, incomeCount, expenseCount)

	return summary, nil
}

// GetSummaryWithComparison gets summary with comparison to previous period
func (uc *SummaryUseCase) GetSummaryWithComparison(ctx context.Context, periodType domain.PeriodType, reference ...time.Time) (*domain.Summary, error) {
	refTime := time.Now()
	if len(reference) > 0 {
		refTime = reference[0]
	}

	// Get current period summary
	currentRange := domain.GetPeriodRange(periodType, refTime)
	if currentRange == nil {
		return nil, fmt.Errorf("failed to get period range for type: %s", periodType)
	}

	summary, err := uc.GetSummaryByDateRange(ctx, currentRange.Start, currentRange.End, periodType)
	if err != nil {
		return nil, err
	}

	// Calculate previous period
	var prevStart, prevEnd time.Time
	switch periodType {
	case domain.PeriodTypeWeek:
		prevStart = currentRange.Start.AddDate(0, 0, -7)
		prevEnd = currentRange.Start.Add(-time.Nanosecond)
	case domain.PeriodTypeMonth:
		prevStart = currentRange.Start.AddDate(0, -1, 0)
		prevEnd = currentRange.Start.Add(-time.Nanosecond)
	case domain.PeriodTypeQuarter:
		prevStart = currentRange.Start.AddDate(0, -3, 0)
		prevEnd = currentRange.Start.Add(-time.Nanosecond)
	case domain.PeriodTypeYear:
		prevStart = currentRange.Start.AddDate(-1, 0, 0)
		prevEnd = currentRange.Start.Add(-time.Nanosecond)
	default:
		return summary, nil // No comparison for custom ranges
	}

	// Get previous period totals
	prevIncome, err := uc.transactionRepo.GetTotalByDateRange(ctx, prevStart, prevEnd, "income")
	if err != nil {
		return summary, nil // Return summary without comparison if previous data fails
	}

	prevExpense, err := uc.transactionRepo.GetTotalByDateRange(ctx, prevStart, prevEnd, "expense")
	if err != nil {
		return summary, nil
	}

	// Calculate comparison
	comparison := &domain.PeriodComparison{
		PreviousPeriodIncome:  prevIncome,
		PreviousPeriodExpense: prevExpense,
		IncomeChange:          summary.TotalIncome - prevIncome,
		ExpenseChange:         summary.TotalExpense - prevExpense,
	}

	// Calculate percentage changes
	if prevIncome > 0 {
		comparison.IncomeChangePercent = (comparison.IncomeChange / prevIncome) * 100
	}
	if prevExpense > 0 {
		comparison.ExpenseChangePercent = (comparison.ExpenseChange / prevExpense) * 100
	}

	summary.SetComparison(comparison)
	return summary, nil
}
