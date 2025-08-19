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
	addExpenseModel    *AddExpenseModel
	addIncomeModel     *AddIncomeModel
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
	m.addExpenseModel = NewAddExpenseModel(transactionUseCase)
	m.addIncomeModel = NewAddIncomeModel(transactionUseCase)
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
		// TODO: Update form models when they get SetDimensions methods
		
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
				m.state = addExpenseView
				m.addExpenseModel.Reset()
				return m, m.addExpenseModel.Init()

			case "i":
				m.state = addIncomeView
				m.addIncomeModel.Reset()
				return m, m.addIncomeModel.Init()

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

	case addExpenseView:
		addExpenseModel, cmd := m.addExpenseModel.Update(msg)
		m.addExpenseModel = addExpenseModel.(*AddExpenseModel)
		if m.addExpenseModel.shouldReturn {
			m.state = dashboardView
			return m, tea.Batch(cmd, m.dashboardModel.Refresh())
		}
		return m, cmd

	case addIncomeView:
		addIncomeModel, cmd := m.addIncomeModel.Update(msg)
		m.addIncomeModel = addIncomeModel.(*AddIncomeModel)
		if m.addIncomeModel.shouldReturn {
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
	case addExpenseView:
		return m.addExpenseModel.View()
	case addIncomeView:
		return m.addIncomeModel.View()
	case listTransactionsView:
		return m.transactionsModel.View()
	default:
		return "Unknown view"
	}
}
