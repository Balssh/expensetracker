package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"expense-tracker/internal/core/domain"
	"expense-tracker/internal/repository/sqlite"
)

type TransactionRepositoryIntegrationSuite struct {
	suite.Suite
	db           *sqlite.Database
	repo         *sqlite.TransactionRepository
	categoryRepo *sqlite.CategoryRepository
	ctx          context.Context
	testDB       string
}

func (suite *TransactionRepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Create a temporary database file for testing
	tempDir := os.TempDir()
	suite.testDB = filepath.Join(tempDir, "test_expense_tracker.db")

	var err error
	suite.db, err = sqlite.NewDatabase(suite.testDB)
	suite.Require().NoError(err)

	suite.repo = sqlite.NewTransactionRepository(suite.db)
	suite.categoryRepo = sqlite.NewCategoryRepository(suite.db)
}

func (suite *TransactionRepositoryIntegrationSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
	os.Remove(suite.testDB)
}

func (suite *TransactionRepositoryIntegrationSuite) SetupTest() {
	// Clean up test data before each test
	suite.cleanupTestData()
}

func (suite *TransactionRepositoryIntegrationSuite) TearDownTest() {
	// Clean up test data after each test
	suite.cleanupTestData()
}

func (suite *TransactionRepositoryIntegrationSuite) cleanupTestData() {
	// Delete test transactions and categories
	_, err := suite.db.DB().Exec("DELETE FROM transactions")
	suite.Require().NoError(err)
	_, err = suite.db.DB().Exec("DELETE FROM categories WHERE name NOT IN ('Food & Dining', 'Transportation', 'Shopping', 'Entertainment', 'Bills & Utilities', 'Healthcare', 'Other', 'Salary', 'Freelance', 'Investment', 'Gift')")
	suite.Require().NoError(err)
}

func TestTransactionRepositoryIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryIntegrationSuite))
}

func (suite *TransactionRepositoryIntegrationSuite) TestCreateAndRetrieve() {
	assert := assert.New(suite.T())

	// Create a test category first
	category := &domain.Category{Name: "Test Category"}
	err := suite.categoryRepo.CreateCategory(suite.ctx, category, "expense")
	assert.NoError(err)
	assert.NotZero(category.ID)

	// Create a transaction
	transaction := &domain.Transaction{
		Description: "Test Transaction",
		Amount:      123.45,
		Type:        "expense",
		Date:        time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
		Category:    category,
	}

	err = suite.repo.Create(suite.ctx, transaction)
	assert.NoError(err)
	assert.NotZero(transaction.ID)

	// Retrieve the transaction
	retrieved, err := suite.repo.GetByID(suite.ctx, transaction.ID)
	assert.NoError(err)
	assert.NotNil(retrieved)
	assert.Equal(transaction.Description, retrieved.Description)
	assert.Equal(transaction.Amount, retrieved.Amount)
	assert.Equal(transaction.Type, retrieved.Type)
	assert.Equal(transaction.Date.Truncate(time.Second), retrieved.Date.Truncate(time.Second))
	assert.NotNil(retrieved.Category)
	assert.Equal(category.ID, retrieved.Category.ID)
	assert.Equal(category.Name, retrieved.Category.Name)
}

func (suite *TransactionRepositoryIntegrationSuite) TestCreateWithoutCategory() {
	assert := assert.New(suite.T())

	transaction := &domain.Transaction{
		Description: "Test Transaction Without Category",
		Amount:      50.0,
		Type:        "income",
		Date:        time.Now(),
	}

	err := suite.repo.Create(suite.ctx, transaction)
	assert.NoError(err)
	assert.NotZero(transaction.ID)

	// Retrieve and verify
	retrieved, err := suite.repo.GetByID(suite.ctx, transaction.ID)
	assert.NoError(err)
	assert.NotNil(retrieved)
	assert.Nil(retrieved.Category)
}

func (suite *TransactionRepositoryIntegrationSuite) TestGetByID_NotFound() {
	assert := assert.New(suite.T())

	retrieved, err := suite.repo.GetByID(suite.ctx, 9999)
	assert.Error(err)
	assert.Nil(retrieved)
}

func (suite *TransactionRepositoryIntegrationSuite) TestGetAll() {
	assert := assert.New(suite.T())

	// Create test transactions
	transactions := []*domain.Transaction{
		{Description: "Transaction 1", Amount: 100.0, Type: "expense", Date: time.Now()},
		{Description: "Transaction 2", Amount: 200.0, Type: "income", Date: time.Now().Add(-time.Hour)},
		{Description: "Transaction 3", Amount: 300.0, Type: "expense", Date: time.Now().Add(-2 * time.Hour)},
	}

	for _, tx := range transactions {
		err := suite.repo.Create(suite.ctx, tx)
		assert.NoError(err)
	}

	// Retrieve all transactions
	retrieved, err := suite.repo.GetAll(suite.ctx, 0, 10)
	assert.NoError(err)
	assert.Len(retrieved, 3)

	// Test pagination
	retrieved, err = suite.repo.GetAll(suite.ctx, 0, 2)
	assert.NoError(err)
	assert.Len(retrieved, 2)

	retrieved, err = suite.repo.GetAll(suite.ctx, 2, 2)
	assert.NoError(err)
	assert.Len(retrieved, 1)
}

func (suite *TransactionRepositoryIntegrationSuite) TestGetByType() {
	assert := assert.New(suite.T())

	// Create mixed transactions
	expenseTransaction := &domain.Transaction{
		Description: "Expense Transaction",
		Amount:      100.0,
		Type:        "expense",
		Date:        time.Now(),
	}
	incomeTransaction := &domain.Transaction{
		Description: "Income Transaction",
		Amount:      200.0,
		Type:        "income",
		Date:        time.Now(),
	}

	err := suite.repo.Create(suite.ctx, expenseTransaction)
	assert.NoError(err)
	err = suite.repo.Create(suite.ctx, incomeTransaction)
	assert.NoError(err)

	// Test expense filtering
	expenses, err := suite.repo.GetByType(suite.ctx, "expense", 0, 10)
	assert.NoError(err)
	assert.Len(expenses, 1)
	assert.Equal("expense", expenses[0].Type)

	// Test income filtering
	incomes, err := suite.repo.GetByType(suite.ctx, "income", 0, 10)
	assert.NoError(err)
	assert.Len(incomes, 1)
	assert.Equal("income", incomes[0].Type)
}

func (suite *TransactionRepositoryIntegrationSuite) TestGetByDateRange() {
	assert := assert.New(suite.T())

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	tomorrow := now.AddDate(0, 0, 1)

	// Create transactions with different dates
	transactions := []*domain.Transaction{
		{Description: "Yesterday", Amount: 100.0, Type: "expense", Date: yesterday},
		{Description: "Today", Amount: 200.0, Type: "expense", Date: now},
		{Description: "Tomorrow", Amount: 300.0, Type: "expense", Date: tomorrow},
	}

	for _, tx := range transactions {
		err := suite.repo.Create(suite.ctx, tx)
		assert.NoError(err)
	}

	// Test date range filtering
	retrieved, err := suite.repo.GetByDateRange(suite.ctx, yesterday.Add(-time.Hour), now.Add(time.Hour))
	assert.NoError(err)
	assert.Len(retrieved, 2) // Should include yesterday and today, not tomorrow
}

func (suite *TransactionRepositoryIntegrationSuite) TestUpdate() {
	assert := assert.New(suite.T())

	// Create initial transaction
	transaction := &domain.Transaction{
		Description: "Original Description",
		Amount:      100.0,
		Type:        "expense",
		Date:        time.Now(),
	}

	err := suite.repo.Create(suite.ctx, transaction)
	assert.NoError(err)
	originalID := transaction.ID

	// Update the transaction
	transaction.Description = "Updated Description"
	transaction.Amount = 150.0

	err = suite.repo.Update(suite.ctx, transaction)
	assert.NoError(err)

	// Retrieve and verify update
	retrieved, err := suite.repo.GetByID(suite.ctx, originalID)
	assert.NoError(err)
	assert.Equal("Updated Description", retrieved.Description)
	assert.Equal(150.0, retrieved.Amount)
	assert.Equal(originalID, retrieved.ID)
}

func (suite *TransactionRepositoryIntegrationSuite) TestDelete() {
	assert := assert.New(suite.T())

	// Create transaction
	transaction := &domain.Transaction{
		Description: "To Be Deleted",
		Amount:      100.0,
		Type:        "expense",
		Date:        time.Now(),
	}

	err := suite.repo.Create(suite.ctx, transaction)
	assert.NoError(err)
	transactionID := transaction.ID

	// Delete the transaction
	err = suite.repo.Delete(suite.ctx, transactionID)
	assert.NoError(err)

	// Verify deletion
	retrieved, err := suite.repo.GetByID(suite.ctx, transactionID)
	assert.Error(err)
	assert.Nil(retrieved)
}

func (suite *TransactionRepositoryIntegrationSuite) TestGetTotalByDateRange() {
	assert := assert.New(suite.T())

	startDate := time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	// Create test transactions
	transactions := []*domain.Transaction{
		{Description: "Income 1", Amount: 1000.0, Type: "income", Date: time.Date(2023, 12, 15, 10, 0, 0, 0, time.UTC)},
		{Description: "Income 2", Amount: 500.0, Type: "income", Date: time.Date(2023, 12, 20, 10, 0, 0, 0, time.UTC)},
		{Description: "Expense 1", Amount: 300.0, Type: "expense", Date: time.Date(2023, 12, 10, 10, 0, 0, 0, time.UTC)},
		{Description: "Expense 2", Amount: 200.0, Type: "expense", Date: time.Date(2023, 12, 25, 10, 0, 0, 0, time.UTC)},
		{Description: "Outside Range", Amount: 1000.0, Type: "income", Date: time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC)},
	}

	for _, tx := range transactions {
		err := suite.repo.Create(suite.ctx, tx)
		assert.NoError(err)
	}

	// Test income total
	totalIncome, err := suite.repo.GetTotalByDateRange(suite.ctx, startDate, endDate, "income")
	assert.NoError(err)
	assert.Equal(1500.0, totalIncome) // 1000 + 500

	// Test expense total
	totalExpense, err := suite.repo.GetTotalByDateRange(suite.ctx, startDate, endDate, "expense")
	assert.NoError(err)
	assert.Equal(500.0, totalExpense) // 300 + 200
}

func (suite *TransactionRepositoryIntegrationSuite) TestSearchTransactions() {
	assert := assert.New(suite.T())

	// Create test transactions
	transactions := []*domain.Transaction{
		{Description: "Grocery shopping at supermarket", Amount: 100.0, Type: "expense", Date: time.Now()},
		{Description: "Salary payment", Amount: 2000.0, Type: "income", Date: time.Now()},
		{Description: "Coffee shop expense", Amount: 5.0, Type: "expense", Date: time.Now()},
		{Description: "Freelance income", Amount: 500.0, Type: "income", Date: time.Now()},
	}

	for _, tx := range transactions {
		err := suite.repo.Create(suite.ctx, tx)
		assert.NoError(err)
	}

	// Test search functionality
	results, err := suite.repo.SearchTransactions(suite.ctx, "shop", 0, 10)
	assert.NoError(err)
	assert.GreaterOrEqual(len(results), 1) // Should find at least "Grocery shopping"

	results, err = suite.repo.SearchTransactions(suite.ctx, "income", 0, 10)
	assert.NoError(err)
	assert.GreaterOrEqual(len(results), 1) // Should find at least one income transaction

	results, err = suite.repo.SearchTransactions(suite.ctx, "nonexistent", 0, 10)
	assert.NoError(err)
	assert.Len(results, 0)
}
