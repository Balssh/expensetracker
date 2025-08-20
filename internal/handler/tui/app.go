package tui

import (
	"expense-tracker/internal/core/usecase"

	tea "github.com/charmbracelet/bubbletea"
)

type sessionState int

const (
	dashboardView sessionState = iota
	addExpenseView
	addIncomeView
	listTransactionsView
)

type Model struct {
	state              sessionState
	width              int
	height             int
	transactionUseCase *usecase.TransactionUseCase
	summaryUseCase     *usecase.SummaryUseCase
	dashboardModel     *DashboardModel
	addTransactionModel *AddTransactionModel
	transactionsModel  *TransactionsModel
}

func NewModel(
	transactionUseCase *usecase.TransactionUseCase,
	summaryUseCase *usecase.SummaryUseCase,
) *Model {
	m := &Model{
		state:              dashboardView,
		transactionUseCase: transactionUseCase,
		summaryUseCase:     summaryUseCase,
	}

	m.dashboardModel = NewDashboardModel(summaryUseCase)
	// Start with expense as default - will be reconfigured when needed
	m.addTransactionModel = NewAddTransactionModel(transactionUseCase, TransactionTypeExpense)
	m.transactionsModel = NewTransactionsModel(summaryUseCase)

	return m
}

func (m Model) Init() tea.Cmd {
	return m.dashboardModel.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update dimensions for all models that need responsive layout
		m.dashboardModel.SetDimensions(msg.Width, msg.Height)
		m.transactionsModel.SetDimensions(msg.Width, msg.Height)
		m.addTransactionModel.SetDimensions(msg.Width, msg.Height)
		
		return m, nil

	case tea.KeyMsg:
		// Global navigation - works in all views
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
			
		case "q", "esc":
			// Context-sensitive quit/back behavior
			if m.state == dashboardView {
				return m, tea.Quit
			}
			// Go back to dashboard from other views
			m.state = dashboardView
			return m, m.dashboardModel.Refresh()
		}

		// Dashboard-specific navigation
		if m.state == dashboardView {
			switch msg.String() {
			case "a":
				// Configure for expense and switch to form
				m.addTransactionModel = NewAddTransactionModel(m.transactionUseCase, TransactionTypeExpense)
				m.addTransactionModel.SetDimensions(m.width, m.height)
				m.state = addExpenseView
				m.addTransactionModel.Reset()
				return m, m.addTransactionModel.Init()

			case "i":
				// Configure for income and switch to form
				m.addTransactionModel = NewAddTransactionModel(m.transactionUseCase, TransactionTypeIncome)
				m.addTransactionModel.SetDimensions(m.width, m.height)
				m.state = addIncomeView
				m.addTransactionModel.Reset()
				return m, m.addTransactionModel.Init()

			case "l":
				m.state = listTransactionsView
				return m, m.transactionsModel.Init()
				
			case "r":
				// Refresh data
				return m, m.dashboardModel.Refresh()
				
			case "?", "h":
				// TODO: Show help modal in future iteration
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	switch m.state {
	case dashboardView:
		dashboardModel, cmd := m.dashboardModel.Update(msg)
		m.dashboardModel = dashboardModel.(*DashboardModel)
		return m, cmd

	case addExpenseView, addIncomeView:
		// Both expense and income use the same unified model
		addTransactionModel, cmd := m.addTransactionModel.Update(msg)
		m.addTransactionModel = addTransactionModel.(*AddTransactionModel)
		if m.addTransactionModel.shouldReturn {
			m.state = dashboardView
			return m, tea.Batch(cmd, m.dashboardModel.Refresh())
		}
		return m, cmd

	case listTransactionsView:
		transactionsModel, cmd := m.transactionsModel.Update(msg)
		m.transactionsModel = transactionsModel.(*TransactionsModel)
		return m, cmd
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.state {
	case dashboardView:
		return m.dashboardModel.View()
	case addExpenseView, addIncomeView:
		return m.addTransactionModel.View()
	case listTransactionsView:
		return m.transactionsModel.View()
	default:
		return "Unknown view"
	}
}
