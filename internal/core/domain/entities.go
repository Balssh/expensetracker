package domain

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type PeriodType string

const (
	PeriodTypeWeek    PeriodType = "week"
	PeriodTypeMonth   PeriodType = "month"
	PeriodTypeQuarter PeriodType = "quarter"
	PeriodTypeYear    PeriodType = "year"
	PeriodTypeCustom  PeriodType = "custom"
)

func (p PeriodType) IsValid() bool {
	switch p {
	case PeriodTypeWeek, PeriodTypeMonth, PeriodTypeQuarter, PeriodTypeYear, PeriodTypeCustom:
		return true
	default:
		return false
	}
}

type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func NewDateRange(start, end time.Time) *DateRange {
	return &DateRange{
		Start: start,
		End:   end,
	}
}

func (dr *DateRange) Validate() error {
	if dr.Start.IsZero() {
		return fmt.Errorf("start date cannot be zero")
	}
	if dr.End.IsZero() {
		return fmt.Errorf("end date cannot be zero")
	}
	if dr.Start.After(dr.End) {
		return fmt.Errorf("start date cannot be after end date")
	}
	return nil
}

func (dr *DateRange) Contains(t time.Time) bool {
	return (t.Equal(dr.Start) || t.After(dr.Start)) && (t.Equal(dr.End) || t.Before(dr.End))
}

func (dr *DateRange) Duration() time.Duration {
	return dr.End.Sub(dr.Start)
}

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

type CategoryBreakdown struct {
	Category         *Category `json:"category"`
	TotalAmount      float64   `json:"total_amount"`
	TransactionCount int       `json:"transaction_count"`
	Percentage       float64   `json:"percentage"`
}

type PeriodComparison struct {
	PreviousPeriodIncome  float64 `json:"previous_period_income"`
	PreviousPeriodExpense float64 `json:"previous_period_expense"`
	IncomeChange          float64 `json:"income_change"`
	ExpenseChange         float64 `json:"expense_change"`
	IncomeChangePercent   float64 `json:"income_change_percent"`
	ExpenseChangePercent  float64 `json:"expense_change_percent"`
}

type Summary struct {
	TotalIncome             float64              `json:"total_income"`
	TotalExpense            float64              `json:"total_expense"`
	NetBalance              float64              `json:"net_balance"`
	Period                  PeriodType           `json:"period_type"`
	DateRange               *DateRange           `json:"date_range"`
	IncomeBreakdown         []*CategoryBreakdown `json:"income_breakdown,omitempty"`
	ExpenseBreakdown        []*CategoryBreakdown `json:"expense_breakdown,omitempty"`
	TransactionCount        int                  `json:"transaction_count"`
	IncomeTransactionCount  int                  `json:"income_transaction_count"`
	ExpenseTransactionCount int                  `json:"expense_transaction_count"`
	Comparison              *PeriodComparison    `json:"comparison,omitempty"`
}

func NewSummary(totalIncome, totalExpense float64) *Summary {
	return &Summary{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetBalance:   totalIncome - totalExpense,
	}
}

func NewEnhancedSummary(totalIncome, totalExpense float64, period PeriodType, dateRange *DateRange) *Summary {
	return &Summary{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		NetBalance:   totalIncome - totalExpense,
		Period:       period,
		DateRange:    dateRange,
	}
}

func (s *Summary) SetCategoryBreakdowns(incomeBreakdown, expenseBreakdown []*CategoryBreakdown) {
	s.IncomeBreakdown = incomeBreakdown
	s.ExpenseBreakdown = expenseBreakdown

	s.calculatePercentages()
}

func (s *Summary) SetTransactionCounts(total, income, expense int) {
	s.TransactionCount = total
	s.IncomeTransactionCount = income
	s.ExpenseTransactionCount = expense
}

func (s *Summary) SetComparison(comparison *PeriodComparison) {
	s.Comparison = comparison
}

func (s *Summary) calculatePercentages() {
	for _, breakdown := range s.IncomeBreakdown {
		if s.TotalIncome > 0 {
			breakdown.Percentage = (breakdown.TotalAmount / s.TotalIncome) * 100
			breakdown.Percentage = math.Round(breakdown.Percentage*100) / 100
		}
	}

	for _, breakdown := range s.ExpenseBreakdown {
		if s.TotalExpense > 0 {
			breakdown.Percentage = (breakdown.TotalAmount / s.TotalExpense) * 100
			breakdown.Percentage = math.Round(breakdown.Percentage*100) / 100
		}
	}
}

func (s *Summary) GetTopExpenseCategories(limit int) []*CategoryBreakdown {
	if len(s.ExpenseBreakdown) <= limit {
		return s.ExpenseBreakdown
	}
	return s.ExpenseBreakdown[:limit]
}

func (s *Summary) GetTopIncomeCategories(limit int) []*CategoryBreakdown {
	if len(s.IncomeBreakdown) <= limit {
		return s.IncomeBreakdown
	}
	return s.IncomeBreakdown[:limit]
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

func GetCurrentWeekRange() *DateRange {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := now.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return NewDateRange(start, end)
}

func GetCurrentMonthRange() *DateRange {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return NewDateRange(start, end)
}

func GetCurrentQuarterRange() *DateRange {
	now := time.Now()
	quarter := (int(now.Month()) - 1) / 3
	startMonth := time.Month(quarter*3 + 1)
	start := time.Date(now.Year(), startMonth, 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 3, 0).Add(-time.Nanosecond)
	return NewDateRange(start, end)
}

func GetCurrentYearRange() *DateRange {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	return NewDateRange(start, end)
}

func GetWeekRange(year, week int) *DateRange {
	jan1 := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	_, w := jan1.ISOWeek()

	var start time.Time
	if w == 1 {
		start = jan1.AddDate(0, 0, (week-1)*7)
	} else {
		mondayOfWeek1 := jan1.AddDate(0, 0, -int(jan1.Weekday())+1)
		if jan1.Weekday() == time.Sunday {
			mondayOfWeek1 = mondayOfWeek1.AddDate(0, 0, 1)
		}
		start = mondayOfWeek1.AddDate(0, 0, (week-1)*7)
	}

	end := start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return NewDateRange(start, end)
}

func GetMonthRange(year int, month time.Month) *DateRange {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return NewDateRange(start, end)
}

func GetQuarterRange(year, quarter int) *DateRange {
	if quarter < 1 || quarter > 4 {
		return nil
	}
	startMonth := time.Month((quarter-1)*3 + 1)
	start := time.Date(year, startMonth, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 3, 0).Add(-time.Nanosecond)
	return NewDateRange(start, end)
}

func GetYearRange(year int) *DateRange {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	return NewDateRange(start, end)
}

func GetPeriodRange(periodType PeriodType, reference time.Time) *DateRange {
	switch periodType {
	case PeriodTypeWeek:
		weekday := int(reference.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		start := reference.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
		end := start.AddDate(0, 0, 7).Add(-time.Nanosecond)
		return NewDateRange(start, end)
	case PeriodTypeMonth:
		start := time.Date(reference.Year(), reference.Month(), 1, 0, 0, 0, 0, reference.Location())
		end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		return NewDateRange(start, end)
	case PeriodTypeQuarter:
		quarter := (int(reference.Month()) - 1) / 3
		startMonth := time.Month(quarter*3 + 1)
		start := time.Date(reference.Year(), startMonth, 1, 0, 0, 0, 0, reference.Location())
		end := start.AddDate(0, 3, 0).Add(-time.Nanosecond)
		return NewDateRange(start, end)
	case PeriodTypeYear:
		start := time.Date(reference.Year(), 1, 1, 0, 0, 0, 0, reference.Location())
		end := start.AddDate(1, 0, 0).Add(-time.Nanosecond)
		return NewDateRange(start, end)
	default:
		return nil
	}
}
