package usecase

import (
	"fmt"
	"time"

	"github.com/yourusername/expense-tracker/internal/core/domain"
)

// TransactionUseCase handles business logic for transactions
type TransactionUseCase struct {
	transactionRepo TransactionRepository
	categoryRepo    CategoryRepository
}

// NewTransactionUseCase creates a new TransactionUseCase
func NewTransactionUseCase(transactionRepo TransactionRepository, categoryRepo CategoryRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
	}
}

// AddTransaction adds a new transaction with validation
func (uc *TransactionUseCase) AddTransaction(transactionType domain.TransactionType, amount float64, categoryID int, description string, date time.Time) (*domain.Transaction, error) {
	// Validate that the category exists and matches the transaction type
	category, err := uc.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if category.Type != transactionType {
		return nil, fmt.Errorf("category type %s does not match transaction type %s", category.Type, transactionType)
	}

	// Create the transaction
	transaction := domain.NewTransaction(transactionType, amount, categoryID, description, date)

	// Validate the transaction
	if err := transaction.Validate(); err != nil {
		return nil, fmt.Errorf("invalid transaction: %w", err)
	}

	// Save to repository
	if err := uc.transactionRepo.Create(transaction); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return transaction, nil
}

// GetTransaction retrieves a transaction by ID
func (uc *TransactionUseCase) GetTransaction(id string) (*domain.Transaction, error) {
	transaction, err := uc.transactionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	return transaction, nil
}

// ListTransactions retrieves transactions with pagination
func (uc *TransactionUseCase) ListTransactions(limit, offset int) ([]*domain.Transaction, error) {
	transactions, err := uc.transactionRepo.List(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	return transactions, nil
}

// ListTransactionsByType retrieves transactions of a specific type
func (uc *TransactionUseCase) ListTransactionsByType(transactionType domain.TransactionType, limit, offset int) ([]*domain.Transaction, error) {
	transactions, err := uc.transactionRepo.ListByType(transactionType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions by type: %w", err)
	}
	return transactions, nil
}

// GetRecentTransactions retrieves transactions from the last N days
func (uc *TransactionUseCase) GetRecentTransactions(days int) ([]*domain.Transaction, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	
	transactions, err := uc.transactionRepo.ListByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	return transactions, nil
}

// UpdateTransaction updates an existing transaction
func (uc *TransactionUseCase) UpdateTransaction(transaction *domain.Transaction) error {
	// Validate the updated transaction
	if err := transaction.Validate(); err != nil {
		return fmt.Errorf("invalid transaction: %w", err)
	}

	// Validate that the category exists and matches the transaction type
	category, err := uc.categoryRepo.GetByID(transaction.CategoryID)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}

	if category.Type != transaction.Type {
		return fmt.Errorf("category type %s does not match transaction type %s", category.Type, transaction.Type)
	}

	// Update in repository
	if err := uc.transactionRepo.Update(transaction); err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

// DeleteTransaction deletes a transaction by ID
func (uc *TransactionUseCase) DeleteTransaction(id string) error {
	// Check if transaction exists
	_, err := uc.transactionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	// Delete from repository
	if err := uc.transactionRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	return nil
}

// GetMonthlyBalance calculates the balance for a specific month
func (uc *TransactionUseCase) GetMonthlyBalance(year, month int) (income, expense, balance float64, err error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	transactions, err := uc.transactionRepo.ListByDateRange(startDate, endDate)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get monthly transactions: %w", err)
	}

	for _, transaction := range transactions {
		switch transaction.Type {
		case domain.TypeIncome:
			income += transaction.Amount
		case domain.TypeExpense:
			expense += transaction.Amount
		}
	}

	balance = income - expense
	return income, expense, balance, nil
}

// GetCategoryBreakdown returns spending/income breakdown by category for a given month
func (uc *TransactionUseCase) GetCategoryBreakdown(transactionType domain.TransactionType, year, month int) (map[string]float64, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	transactions, err := uc.transactionRepo.ListByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	breakdown := make(map[string]float64)
	categoryNames := make(map[int]string)

	// Get category names
	categories, err := uc.categoryRepo.ListByType(transactionType)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	for _, category := range categories {
		categoryNames[category.ID] = category.Name
	}

	// Calculate breakdown
	for _, transaction := range transactions {
		if transaction.Type == transactionType {
			categoryName := categoryNames[transaction.CategoryID]
			if categoryName == "" {
				categoryName = "Unknown"
			}
			breakdown[categoryName] += transaction.Amount
		}
	}

	return breakdown, nil
}