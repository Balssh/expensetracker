package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"expense-tracker/internal/core/domain"
	"expense-tracker/internal/repository/sqlite"
)

type CategoryRepositoryIntegrationSuite struct {
	suite.Suite
	db     *sqlite.Database
	repo   *sqlite.CategoryRepository
	ctx    context.Context
	testDB string
}

func (suite *CategoryRepositoryIntegrationSuite) SetupSuite() {
	suite.ctx = context.Background()

	// Create a temporary database file for testing
	tempDir := os.TempDir()
	suite.testDB = filepath.Join(tempDir, "test_category_expense_tracker.db")

	var err error
	suite.db, err = sqlite.NewDatabase(suite.testDB)
	suite.Require().NoError(err)

	suite.repo = sqlite.NewCategoryRepository(suite.db)
}

func (suite *CategoryRepositoryIntegrationSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
	os.Remove(suite.testDB)
}

func (suite *CategoryRepositoryIntegrationSuite) SetupTest() {
	// Clean up test data before each test
	suite.cleanupTestData()
}

func (suite *CategoryRepositoryIntegrationSuite) TearDownTest() {
	// Clean up test data after each test
	suite.cleanupTestData()
}

func (suite *CategoryRepositoryIntegrationSuite) cleanupTestData() {
	// Delete test categories
	_, err := suite.db.DB().Exec("DELETE FROM categories WHERE name NOT IN ('Food & Dining', 'Transportation', 'Shopping', 'Entertainment', 'Bills & Utilities', 'Healthcare', 'Other', 'Salary', 'Freelance', 'Investment', 'Gift')")
	suite.Require().NoError(err)
}

func TestCategoryRepositoryIntegrationSuite(t *testing.T) {
	suite.Run(t, new(CategoryRepositoryIntegrationSuite))
}

func (suite *CategoryRepositoryIntegrationSuite) TestCreateAndGetExpenseCategory() {
	assert := assert.New(suite.T())

	category := &domain.Category{
		Name: "Test Food Category",
	}

	// Create expense category
	err := suite.repo.CreateCategory(suite.ctx, category, "expense")
	assert.NoError(err)
	assert.NotZero(category.ID)

	// Retrieve by ID
	retrieved, err := suite.repo.GetCategoryByID(suite.ctx, category.ID, "expense")
	assert.NoError(err)
	assert.NotNil(retrieved)
	assert.Equal(category.Name, retrieved.Name)
	assert.Equal(category.ID, retrieved.ID)
}

func (suite *CategoryRepositoryIntegrationSuite) TestCreateAndGetIncomeCategory() {
	assert := assert.New(suite.T())

	category := &domain.Category{
		Name: "Test Salary Category",
	}

	// Create income category
	err := suite.repo.CreateCategory(suite.ctx, category, "income")
	assert.NoError(err)
	assert.NotZero(category.ID)

	// Retrieve by ID
	retrieved, err := suite.repo.GetCategoryByID(suite.ctx, category.ID, "income")
	assert.NoError(err)
	assert.NotNil(retrieved)
	assert.Equal(category.Name, retrieved.Name)
	assert.Equal(category.ID, retrieved.ID)
}

func (suite *CategoryRepositoryIntegrationSuite) TestGetCategoryByID_NotFound() {
	assert := assert.New(suite.T())

	// Try to retrieve non-existent category
	retrieved, err := suite.repo.GetCategoryByID(suite.ctx, 9999, "expense")
	assert.Error(err)
	assert.Nil(retrieved)
}

func (suite *CategoryRepositoryIntegrationSuite) TestGetCategories() {
	assert := assert.New(suite.T())

	// Create expense categories
	expenseCategories := []*domain.Category{
		{Name: "Test Food"},
		{Name: "Test Transport"},
		{Name: "Test Entertainment"},
	}

	for _, cat := range expenseCategories {
		err := suite.repo.CreateCategory(suite.ctx, cat, "expense")
		assert.NoError(err)
	}

	// Create income categories
	incomeCategories := []*domain.Category{
		{Name: "Test Salary"},
		{Name: "Test Freelance"},
	}

	for _, cat := range incomeCategories {
		err := suite.repo.CreateCategory(suite.ctx, cat, "income")
		assert.NoError(err)
	}

	// Retrieve expense categories (includes default + test categories)
	retrievedExpenses, err := suite.repo.GetCategories(suite.ctx, "expense")
	assert.NoError(err)
	assert.GreaterOrEqual(len(retrievedExpenses), 3)

	// Check that our test categories are in the results
	names := make([]string, len(retrievedExpenses))
	for i, cat := range retrievedExpenses {
		names[i] = cat.Name
	}
	assert.Contains(names, "Test Food")
	assert.Contains(names, "Test Transport")
	assert.Contains(names, "Test Entertainment")

	// Retrieve income categories (includes default + test categories)
	retrievedIncomes, err := suite.repo.GetCategories(suite.ctx, "income")
	assert.NoError(err)
	assert.GreaterOrEqual(len(retrievedIncomes), 2)

	// Check that our test categories are in the results
	incomeNames := make([]string, len(retrievedIncomes))
	for i, cat := range retrievedIncomes {
		incomeNames[i] = cat.Name
	}
	assert.Contains(incomeNames, "Test Salary")
	assert.Contains(incomeNames, "Test Freelance")
}

func (suite *CategoryRepositoryIntegrationSuite) TestCreateCategory_DuplicateName() {
	assert := assert.New(suite.T())

	// Create first category
	category1 := &domain.Category{Name: "Test Unique Food"}
	err := suite.repo.CreateCategory(suite.ctx, category1, "expense")
	assert.NoError(err)

	// Try to create second category with same name and type
	category2 := &domain.Category{Name: "Test Unique Food"}
	err = suite.repo.CreateCategory(suite.ctx, category2, "expense")
	assert.Error(err) // Should fail due to unique constraint
}

func (suite *CategoryRepositoryIntegrationSuite) TestCreateCategory_SameNameDifferentType() {
	assert := assert.New(suite.T())

	// Create expense category
	expenseCategory := &domain.Category{Name: "Test Business"}
	err := suite.repo.CreateCategory(suite.ctx, expenseCategory, "expense")
	assert.NoError(err)

	// Create income category with same name (should be allowed)
	incomeCategory := &domain.Category{Name: "Test Business"}
	err = suite.repo.CreateCategory(suite.ctx, incomeCategory, "income")
	assert.NoError(err)

	// Both should exist
	assert.NotZero(expenseCategory.ID)
	assert.NotZero(incomeCategory.ID)
	assert.NotEqual(expenseCategory.ID, incomeCategory.ID)
}

func (suite *CategoryRepositoryIntegrationSuite) TestCategoryTypeSeparation() {
	assert := assert.New(suite.T())

	// Create expense category
	expenseCategory := &domain.Category{Name: "Test Separation Food"}
	err := suite.repo.CreateCategory(suite.ctx, expenseCategory, "expense")
	assert.NoError(err)

	// Create income category
	incomeCategory := &domain.Category{Name: "Test Separation Salary"}
	err = suite.repo.CreateCategory(suite.ctx, incomeCategory, "income")
	assert.NoError(err)

	// Try to get expense category as income category
	retrieved, err := suite.repo.GetCategoryByID(suite.ctx, expenseCategory.ID, "income")
	assert.Error(err)
	assert.Nil(retrieved)

	// Try to get income category as expense category
	retrieved, err = suite.repo.GetCategoryByID(suite.ctx, incomeCategory.ID, "expense")
	assert.Error(err)
	assert.Nil(retrieved)

	// But retrieving with correct types should work
	retrieved, err = suite.repo.GetCategoryByID(suite.ctx, expenseCategory.ID, "expense")
	assert.NoError(err)
	assert.NotNil(retrieved)

	retrieved, err = suite.repo.GetCategoryByID(suite.ctx, incomeCategory.ID, "income")
	assert.NoError(err)
	assert.NotNil(retrieved)
}
