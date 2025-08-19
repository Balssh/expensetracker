package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"expense-tracker/internal/core/domain"
	"expense-tracker/test/mocks"
)

type SummaryUseCaseTestSuite struct {
	suite.Suite
	useCase         *SummaryUseCase
	transactionRepo *mocks.MockTransactionRepository
	ctx             context.Context
}

func (suite *SummaryUseCaseTestSuite) SetupTest() {
	suite.transactionRepo = mocks.NewMockTransactionRepository(suite.T())
	suite.useCase = NewSummaryUseCase(suite.transactionRepo)
	suite.ctx = context.Background()
}

func TestSummaryUseCaseSuite(t *testing.T) {
	suite.Run(t, new(SummaryUseCaseTestSuite))
}

func (suite *SummaryUseCaseTestSuite) TestGetMonthlySummary_Success() {
	assert := assert.New(suite.T())

	year := 2023
	month := time.December

	// Calculate expected date range
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "income").Return(2500.0, nil)
	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "expense").Return(1800.0, nil)

	summary, err := suite.useCase.GetMonthlySummary(suite.ctx, year, month)

	assert.NoError(err)
	assert.NotNil(summary)
	assert.Equal(2500.0, summary.TotalIncome)
	assert.Equal(1800.0, summary.TotalExpense)
	assert.Equal(700.0, summary.NetBalance) // 2500 - 1800

	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *SummaryUseCaseTestSuite) TestGetMonthlySummary_IncomeError() {
	assert := assert.New(suite.T())

	year := 2023
	month := time.December

	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "income").Return(0.0, errors.New("database error"))

	summary, err := suite.useCase.GetMonthlySummary(suite.ctx, year, month)

	assert.Error(err)
	assert.Nil(summary)
	assert.Contains(err.Error(), "database error")

	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *SummaryUseCaseTestSuite) TestGetMonthlySummary_ExpenseError() {
	assert := assert.New(suite.T())

	year := 2023
	month := time.December

	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "income").Return(2500.0, nil)
	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "expense").Return(0.0, errors.New("database error"))

	summary, err := suite.useCase.GetMonthlySummary(suite.ctx, year, month)

	assert.Error(err)
	assert.Nil(summary)
	assert.Contains(err.Error(), "database error")

	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *SummaryUseCaseTestSuite) TestGetMonthlySummary_NegativeBalance() {
	assert := assert.New(suite.T())

	year := 2023
	month := time.December

	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "income").Return(1000.0, nil)
	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "expense").Return(1500.0, nil)

	summary, err := suite.useCase.GetMonthlySummary(suite.ctx, year, month)

	assert.NoError(err)
	assert.NotNil(summary)
	assert.Equal(1000.0, summary.TotalIncome)
	assert.Equal(1500.0, summary.TotalExpense)
	assert.Equal(-500.0, summary.NetBalance) // 1000 - 1500

	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *SummaryUseCaseTestSuite) TestGetMonthlySummary_ZeroValues() {
	assert := assert.New(suite.T())

	year := 2023
	month := time.December

	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "income").Return(0.0, nil)
	suite.transactionRepo.On("GetTotalByDateRange", suite.ctx, start, end, "expense").Return(0.0, nil)

	summary, err := suite.useCase.GetMonthlySummary(suite.ctx, year, month)

	assert.NoError(err)
	assert.NotNil(summary)
	assert.Equal(0.0, summary.TotalIncome)
	assert.Equal(0.0, summary.TotalExpense)
	assert.Equal(0.0, summary.NetBalance)

	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *SummaryUseCaseTestSuite) TestGetRecentTransactions_Success() {
	assert := assert.New(suite.T())

	expectedTransactions := []*domain.Transaction{
		{ID: 1, Description: "Recent expense", Amount: 50.0, Type: "expense"},
		{ID: 2, Description: "Recent income", Amount: 100.0, Type: "income"},
	}

	suite.transactionRepo.On("GetRecentTransactions", suite.ctx, 5).Return(expectedTransactions, nil)

	transactions, err := suite.useCase.GetRecentTransactions(suite.ctx, 5)

	assert.NoError(err)
	assert.Equal(expectedTransactions, transactions)

	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *SummaryUseCaseTestSuite) TestGetRecentTransactions_Error() {
	assert := assert.New(suite.T())

	suite.transactionRepo.On("GetRecentTransactions", suite.ctx, 5).Return(nil, errors.New("database error"))

	transactions, err := suite.useCase.GetRecentTransactions(suite.ctx, 5)

	assert.Error(err)
	assert.Nil(transactions)
	assert.Contains(err.Error(), "database error")

	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *SummaryUseCaseTestSuite) TestGetAllTransactions_Success() {
	assert := assert.New(suite.T())

	expectedTransactions := []*domain.Transaction{
		{ID: 1, Description: "Transaction 1", Amount: 100.0},
		{ID: 2, Description: "Transaction 2", Amount: 200.0},
	}

	suite.transactionRepo.On("GetAll", suite.ctx, 0, 10).Return(expectedTransactions, nil)

	transactions, err := suite.useCase.GetAllTransactions(suite.ctx, 0, 10)

	assert.NoError(err)
	assert.Equal(expectedTransactions, transactions)

	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *SummaryUseCaseTestSuite) TestSearchTransactions_Success() {
	assert := assert.New(suite.T())

	expectedTransactions := []*domain.Transaction{
		{ID: 1, Description: "Grocery shopping", Amount: 50.0},
	}

	suite.transactionRepo.On("SearchTransactions", suite.ctx, "grocery", 0, 10).Return(expectedTransactions, nil)

	transactions, err := suite.useCase.SearchTransactions(suite.ctx, "grocery", 0, 10)

	assert.NoError(err)
	assert.Equal(expectedTransactions, transactions)

	suite.transactionRepo.AssertExpectations(suite.T())
}
