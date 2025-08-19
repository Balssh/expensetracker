package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EntityTestSuite struct {
	suite.Suite
}

func (suite *EntityTestSuite) SetupTest() {
	// SetupTest runs before each test method
	// Currently no setup is needed for entity tests
}

func TestEntitySuite(t *testing.T) {
	suite.Run(t, new(EntityTestSuite))
}

func (suite *EntityTestSuite) TestTransactionValidation() {
	assert := assert.New(suite.T())

	tests := []struct {
		name        string
		tx          Transaction
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid income",
			tx: Transaction{
				Description: "Salary",
				Amount:      1234,
				Type:        "income",
				Date:        time.Now(),
			},
			expectError: false,
		},
		{
			name: "valid expense",
			tx: Transaction{
				Description: "Groceries",
				Amount:      567,
				Type:        "expense",
				Date:        time.Now(),
			},
			expectError: false,
		},
		{
			name: "negative amount",
			tx: Transaction{
				Description: "Test",
				Amount:      -100,
				Type:        "expense",
				Date:        time.Now(),
			},
			expectError: true,
			errorMsg:    "transaction amount must be positive",
		},
		{
			name: "zero amount",
			tx: Transaction{
				Description: "Test",
				Amount:      0,
				Type:        "expense",
				Date:        time.Now(),
			},
			expectError: true,
			errorMsg:    "transaction amount must be positive",
		},
		{
			name: "empty description",
			tx: Transaction{
				Description: "",
				Amount:      100,
				Type:        "expense",
				Date:        time.Now(),
			},
			expectError: true,
			errorMsg:    "transaction description cannot be empty",
		},
		{
			name: "whitespace only description",
			tx: Transaction{
				Description: "   ",
				Amount:      100,
				Type:        "expense",
				Date:        time.Now(),
			},
			expectError: true,
			errorMsg:    "transaction description cannot be empty",
		},
		{
			name: "description too long",
			tx: Transaction{
				Description: "This is an extremely long description that exceeds the 200 character limit for transaction descriptions. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
				Amount:      100,
				Type:        "expense",
				Date:        time.Now(),
			},
			expectError: true,
			errorMsg:    "transaction description cannot exceed 200 characters",
		},
		{
			name: "invalid type",
			tx: Transaction{
				Description: "Test",
				Amount:      100,
				Type:        "invalid",
				Date:        time.Now(),
			},
			expectError: true,
			errorMsg:    "transaction type must be 'income' or 'expense'",
		},
		{
			name: "zero date",
			tx: Transaction{
				Description: "Test",
				Amount:      100,
				Type:        "expense",
				Date:        time.Time{},
			},
			expectError: true,
			errorMsg:    "transaction date cannot be zero",
		},
		{
			name: "valid with category",
			tx: Transaction{
				Description: "Groceries",
				Amount:      50,
				Type:        "expense",
				Date:        time.Now(),
				Category:    &Category{Name: "Food"},
			},
			expectError: false,
		},
		{
			name: "invalid category",
			tx: Transaction{
				Description: "Test",
				Amount:      100,
				Type:        "expense",
				Date:        time.Now(),
				Category:    &Category{Name: ""},
			},
			expectError: true,
			errorMsg:    "invalid category",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := tt.tx.Validate()
			if tt.expectError {
				assert.Error(err)
				assert.Contains(err.Error(), tt.errorMsg)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func (suite *EntityTestSuite) TestCategoryValidation() {
	assert := assert.New(suite.T())

	tests := []struct {
		name        string
		category    Category
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid category",
			category:    Category{Name: "Food"},
			expectError: false,
		},
		{
			name:        "empty name",
			category:    Category{Name: ""},
			expectError: true,
			errorMsg:    "category name cannot be empty",
		},
		{
			name:        "whitespace only name",
			category:    Category{Name: "   "},
			expectError: true,
			errorMsg:    "category name cannot be empty",
		},
		{
			name:        "name too long",
			category:    Category{Name: "this is a very long category name that exceeds fifty characters limit"},
			expectError: true,
			errorMsg:    "category name cannot exceed 50 characters",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			err := tt.category.Validate()
			if tt.expectError {
				assert.Error(err)
				assert.Contains(err.Error(), tt.errorMsg)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func (suite *EntityTestSuite) TestSummaryCreation() {
	assert := assert.New(suite.T())

	summary := NewSummary(1000.0, 750.0)

	assert.Equal(1000.0, summary.TotalIncome)
	assert.Equal(750.0, summary.TotalExpense)
	assert.Equal(250.0, summary.NetBalance)
}

func (suite *EntityTestSuite) TestTransactionHelperMethods() {
	assert := assert.New(suite.T())

	incomeTransaction := &Transaction{Type: "income"}
	expenseTransaction := &Transaction{Type: "expense"}
	invalidTransaction := &Transaction{Type: "invalid"}

	assert.True(incomeTransaction.IsIncome())
	assert.False(incomeTransaction.IsExpense())

	assert.False(expenseTransaction.IsIncome())
	assert.True(expenseTransaction.IsExpense())

	assert.False(invalidTransaction.IsIncome())
	assert.False(invalidTransaction.IsExpense())
}

func (suite *EntityTestSuite) TestSummaryCalculations() {
	assert := assert.New(suite.T())

	testCases := []struct {
		name        string
		income      float64
		expense     float64
		expectedNet float64
	}{
		{"positive balance", 1000.0, 750.0, 250.0},
		{"negative balance", 500.0, 800.0, -300.0},
		{"zero balance", 1000.0, 1000.0, 0.0},
		{"zero income", 0.0, 500.0, -500.0},
		{"zero expense", 1000.0, 0.0, 1000.0},
		{"both zero", 0.0, 0.0, 0.0},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			summary := NewSummary(tc.income, tc.expense)

			assert.Equal(tc.income, summary.TotalIncome)
			assert.Equal(tc.expense, summary.TotalExpense)
			assert.Equal(tc.expectedNet, summary.NetBalance)
		})
	}
}

func (suite *EntityTestSuite) TestCategoryEdgeCases() {
	assert := assert.New(suite.T())

	testCases := []struct {
		name        string
		category    Category
		expectError bool
		errorMsg    string
	}{
		{
			name:        "exactly 50 characters",
			category:    Category{Name: "12345678901234567890123456789012345678901234567890"},
			expectError: false,
		},
		{
			name:        "51 characters",
			category:    Category{Name: "123456789012345678901234567890123456789012345678901"},
			expectError: true,
			errorMsg:    "category name cannot exceed 50 characters",
		},
		{
			name:        "single character",
			category:    Category{Name: "A"},
			expectError: false,
		},
		{
			name:        "mixed whitespace",
			category:    Category{Name: " \t\n "},
			expectError: true,
			errorMsg:    "category name cannot be empty",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			err := tc.category.Validate()
			if tc.expectError {
				assert.Error(err)
				assert.Contains(err.Error(), tc.errorMsg)
			} else {
				assert.NoError(err)
			}
		})
	}
}
