package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"expense-tracker/internal/core/usecase"
)

type sessionState int

const (
	dashboardView sessionState = iota
	addExpenseView
	addIncomeView
	listTransactionsView
)

type Model struct {
	state             sessionState
	width             int
	height            int
	expenseUseCase    *usecase.ExpenseUseCase
	incomeUseCase     *usecase.IncomeUseCase
	summaryUseCase    *usecase.SummaryUseCase
	dashboardModel    *DashboardModel
	addExpenseModel   *AddExpenseModel
	addIncomeModel    *AddIncomeModel
	transactionsModel *TransactionsModel
}

func NewModel(
	expenseUseCase *usecase.ExpenseUseCase,
	incomeUseCase *usecase.IncomeUseCase,
	summaryUseCase *usecase.SummaryUseCase,
) *Model {
	m := &Model{
		state:          dashboardView,
		expenseUseCase: expenseUseCase,
		incomeUseCase:  incomeUseCase,
		summaryUseCase: summaryUseCase,
	}

	m.dashboardModel = NewDashboardModel(summaryUseCase)
	m.addExpenseModel = NewAddExpenseModel(expenseUseCase)
	m.addIncomeModel = NewAddIncomeModel(incomeUseCase)
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
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == dashboardView {
				return m, tea.Quit
			}
			m.state = dashboardView
			return m, m.dashboardModel.Refresh()

		case "a":
			if m.state == dashboardView {
				m.state = addExpenseView
				m.addExpenseModel.Reset()
				return m, m.addExpenseModel.Init()
			}

		case "i":
			if m.state == dashboardView {
				m.state = addIncomeView
				m.addIncomeModel.Reset()
				return m, m.addIncomeModel.Init()
			}

		case "l":
			if m.state == dashboardView {
				m.state = listTransactionsView
				return m, m.transactionsModel.Init()
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