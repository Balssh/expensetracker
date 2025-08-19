package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"expense-tracker/internal/core/domain"
	"expense-tracker/test/mocks"
)

type TransactionUseCaseTestSuite struct {
	suite.Suite
	useCase         *TransactionUseCase
	transactionRepo *mocks.MockTransactionRepository
	categoryRepo    *mocks.MockCategoryRepository
	ctx             context.Context
}

func (suite *TransactionUseCaseTestSuite) SetupTest() {
	suite.transactionRepo = mocks.NewMockTransactionRepository(suite.T())
	suite.categoryRepo = mocks.NewMockCategoryRepository(suite.T())
	suite.useCase = NewTransactionUseCase(suite.transactionRepo, suite.categoryRepo)
	suite.ctx = context.Background()
}

func TestTransactionUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TransactionUseCaseTestSuite))
}

func (suite *TransactionUseCaseTestSuite) TestAddTransaction_Success() {
	assert := assert.New(suite.T())

	transaction := &domain.Transaction{
		Description: "Test valid transaction",
		Amount:      100,
		Type:        "expense",
		Date:        time.Now(),
		Category:    &domain.Category{ID: 1},
	}

	// Mock category lookup
	expectedCategory := &domain.Category{ID: 1, Name: "Food"}
	suite.categoryRepo.On("GetCategoryByID", suite.ctx, 1, "expense").Return(expectedCategory, nil)

	// Mock transaction creation
	suite.transactionRepo.On("Create", suite.ctx, transaction).Return(nil)

	err := suite.useCase.AddTransaction(suite.ctx, transaction)

	assert.NoError(err)
	suite.transactionRepo.AssertExpectations(suite.T())
	suite.categoryRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestAddTransaction_InvalidAmount() {
	assert := assert.New(suite.T())

	transaction := &domain.Transaction{
		Description: "Test",
		Amount:      -100, // Invalid negative amount
		Type:        "expense",
		Date:        time.Now(),
	}

	err := suite.useCase.AddTransaction(suite.ctx, transaction)

	assert.Error(err)
	assert.Contains(err.Error(), "transaction amount must be positive")

	// Verify no repository calls were made
	suite.transactionRepo.AssertNotCalled(suite.T(), "Create")
	suite.categoryRepo.AssertNotCalled(suite.T(), "GetCategoryByID")
}

func (suite *TransactionUseCaseTestSuite) TestAddTransaction_EmptyDescription() {
	assert := assert.New(suite.T())

	transaction := &domain.Transaction{
		Description: "", // Invalid empty description
		Amount:      100,
		Type:        "expense",
		Date:        time.Now(),
	}

	err := suite.useCase.AddTransaction(suite.ctx, transaction)

	assert.Error(err)
	assert.Contains(err.Error(), "transaction description is required")

	suite.transactionRepo.AssertNotCalled(suite.T(), "Create")
	suite.categoryRepo.AssertNotCalled(suite.T(), "GetCategoryByID")
}

func (suite *TransactionUseCaseTestSuite) TestAddTransaction_InvalidType() {
	assert := assert.New(suite.T())

	transaction := &domain.Transaction{
		Description: "Test",
		Amount:      100,
		Type:        "invalid", // Invalid type
		Date:        time.Now(),
	}

	err := suite.useCase.AddTransaction(suite.ctx, transaction)

	assert.Error(err)
	assert.Contains(err.Error(), "transaction type must be 'income' or 'expense'")

	suite.transactionRepo.AssertNotCalled(suite.T(), "Create")
	suite.categoryRepo.AssertNotCalled(suite.T(), "GetCategoryByID")
}

func (suite *TransactionUseCaseTestSuite) TestAddTransaction_AutoSetDate() {
	assert := assert.New(suite.T())

	transaction := &domain.Transaction{
		Description: "Test",
		Amount:      100,
		Type:        "expense",
		Date:        time.Time{}, // Zero date should be auto-set
	}

	suite.transactionRepo.On("Create", suite.ctx, mock.AnythingOfType("*domain.Transaction")).Return(nil)

	err := suite.useCase.AddTransaction(suite.ctx, transaction)

	assert.NoError(err)
	assert.False(transaction.Date.IsZero()) // Date should be set
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestAddTransaction_CategoryNotFound() {
	assert := assert.New(suite.T())

	transaction := &domain.Transaction{
		Description: "Test",
		Amount:      100,
		Type:        "expense",
		Date:        time.Now(),
		Category:    &domain.Category{ID: 999}, // Category with ID that doesn't exist
	}

	suite.categoryRepo.On("GetCategoryByID", suite.ctx, 999, "expense").Return(nil, errors.New("category not found"))

	err := suite.useCase.AddTransaction(suite.ctx, transaction)

	assert.Error(err)
	assert.Contains(err.Error(), "invalid category")

	suite.categoryRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertNotCalled(suite.T(), "Create")
}

func (suite *TransactionUseCaseTestSuite) TestGetTransactionsByType_InvalidType() {
	assert := assert.New(suite.T())

	transactions, err := suite.useCase.GetTransactionsByType(suite.ctx, "invalid", 0, 10)

	assert.Nil(transactions)
	assert.Error(err)
	assert.Contains(err.Error(), "transaction type must be 'income' or 'expense'")

	suite.transactionRepo.AssertNotCalled(suite.T(), "GetByType")
}

func (suite *TransactionUseCaseTestSuite) TestGetTransactionsByType_Success() {
	assert := assert.New(suite.T())

	expectedTransactions := []*domain.Transaction{
		{ID: 1, Type: "income", Amount: 1000},
		{ID: 2, Type: "income", Amount: 500},
	}

	suite.transactionRepo.On("GetByType", suite.ctx, "income", 0, 10).Return(expectedTransactions, nil)

	transactions, err := suite.useCase.GetTransactionsByType(suite.ctx, "income", 0, 10)

	assert.NoError(err)
	assert.Equal(expectedTransactions, transactions)
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestUpdateTransaction_InvalidID() {
	assert := assert.New(suite.T())

	transaction := &domain.Transaction{
		ID:          0, // Invalid ID
		Description: "Test",
		Amount:      100,
		Type:        "expense",
	}

	err := suite.useCase.UpdateTransaction(suite.ctx, transaction)

	assert.Error(err)
	assert.Contains(err.Error(), "transaction ID is required for update")

	suite.transactionRepo.AssertNotCalled(suite.T(), "Update")
}

func (suite *TransactionUseCaseTestSuite) TestDeleteTransaction_InvalidID() {
	assert := assert.New(suite.T())

	err := suite.useCase.DeleteTransaction(suite.ctx, 0)

	assert.Error(err)
	assert.Contains(err.Error(), "transaction ID is required for delete")

	suite.transactionRepo.AssertNotCalled(suite.T(), "Delete")
}

func (suite *TransactionUseCaseTestSuite) TestDeleteTransaction_Success() {
	assert := assert.New(suite.T())

	suite.transactionRepo.On("Delete", suite.ctx, 1).Return(nil)

	err := suite.useCase.DeleteTransaction(suite.ctx, 1)

	assert.NoError(err)
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionUseCaseTestSuite) TestGetCategories_InvalidType() {
	assert := assert.New(suite.T())

	categories, err := suite.useCase.GetCategories(suite.ctx, "invalid")

	assert.Nil(categories)
	assert.Error(err)
	assert.Contains(err.Error(), "transaction type must be 'income' or 'expense'")

	suite.categoryRepo.AssertNotCalled(suite.T(), "GetCategories")
}

func (suite *TransactionUseCaseTestSuite) TestGetCategories_Success() {
	assert := assert.New(suite.T())

	expectedCategories := []*domain.Category{
		{ID: 1, Name: "Food"},
		{ID: 2, Name: "Transport"},
	}

	suite.categoryRepo.On("GetCategories", suite.ctx, "expense").Return(expectedCategories, nil)

	categories, err := suite.useCase.GetCategories(suite.ctx, "expense")

	assert.NoError(err)
	assert.Equal(expectedCategories, categories)
	suite.categoryRepo.AssertExpectations(suite.T())
}
