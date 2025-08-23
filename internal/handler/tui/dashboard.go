package tui

import (
	"fmt"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/expense-tracker/internal/core/domain"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
)

// DashboardModel represents the dashboard view model
type DashboardModel struct {
	transactionUseCase *usecase.TransactionUseCase
	categoryUseCase    *usecase.CategoryUseCase
	
	width  int
	height int
	
	// Data
	monthlyIncome  float64
	monthlyExpense float64
	monthlyBalance float64
	
	recentTransactions []*domain.Transaction
	
	// Loading state
	loading bool
	error   string
}

// NewDashboardModel creates a new dashboard model
func NewDashboardModel(transactionUseCase *usecase.TransactionUseCase, categoryUseCase *usecase.CategoryUseCase) *DashboardModel {
	return &DashboardModel{
		transactionUseCase: transactionUseCase,
		categoryUseCase:    categoryUseCase,
		loading:            false,
	}
}

// SetSize sets the size of the dashboard
func (m *DashboardModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Init initializes the dashboard
func (m *DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadMonthlyData(),
		m.loadRecentTransactions(),
	)
}

// Update handles messages and updates the dashboard state
func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case monthlyDataLoadedMsg:
		m.monthlyIncome = msg.income
		m.monthlyExpense = msg.expense
		m.monthlyBalance = msg.balance
		m.loading = false
		
	case recentTransactionsLoadedMsg:
		m.recentTransactions = msg.transactions
		
	case dataLoadErrorMsg:
		m.error = msg.error
		m.loading = false
	}
	
	return m, nil
}

// View renders the dashboard view
func (m *DashboardModel) View() string {
	if m.loading {
		return m.renderLoading()
	}
	
	if m.error != "" {
		return m.renderError()
	}
	
	// Create the main dashboard layout
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderSummary(),
		"",
		m.renderRecentTransactions(),
		"",
		m.renderNavigation(),
	)
}

// renderSummary renders the monthly summary section
func (m *DashboardModel) renderSummary() string {
	now := time.Now()
	monthName := now.Format("January 2006")
	
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true).
		Margin(1, 0, 0, 0)
	
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Margin(0, 1, 0, 0).
		Width(20)
	
	incomeCard := cardStyle.Copy().
		BorderForeground(lipgloss.Color("46")).
		Render(lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true).Render("↗️ Income"),
			fmt.Sprintf("+$%.2f", m.monthlyIncome),
		))
	
	expenseCard := cardStyle.Copy().
		BorderForeground(lipgloss.Color("196")).
		Render(lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("↘️ Expenses"),
			fmt.Sprintf("-$%.2f", m.monthlyExpense),
		))
	
	balanceColor := "46" // Green for positive
	balanceSymbol := "📈"
	balancePrefix := "+"
	if m.monthlyBalance < 0 {
		balanceColor = "196" // Red for negative
		balanceSymbol = "📉"
		balancePrefix = "-"
	} else if m.monthlyBalance == 0 {
		balanceColor = "240" // Gray for zero
		balanceSymbol = "⚖️"
		balancePrefix = ""
	}
	
	balanceCard := cardStyle.Copy().
		BorderForeground(lipgloss.Color(balanceColor)).
		Render(lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color(balanceColor)).Bold(true).Render(balanceSymbol+" Balance"),
			fmt.Sprintf("%s$%.2f", balancePrefix, abs(m.monthlyBalance)),
		))
	
	cards := lipgloss.JoinHorizontal(lipgloss.Top, incomeCard, expenseCard, balanceCard)
	
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("📊 Monthly Summary - "+monthName),
		cards,
	)
}

// renderRecentTransactions renders the recent transactions section
func (m *DashboardModel) renderRecentTransactions() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true).
		Margin(1, 0, 0, 0)
	
	if len(m.recentTransactions) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Italic(true).
			Margin(1, 0)
		
		return lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("📝 Recent Transactions"),
			emptyStyle.Render("No transactions yet. Press '2' to add your first transaction."),
		)
	}
	
	// Create transaction rows
	var rows []string
	
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true).
		Underline(true)
	
	header := fmt.Sprintf("%-12s %-12s %-20s %s",
		"Date", "Type", "Description", "Amount")
	rows = append(rows, headerStyle.Render(header))
	
	for i, tx := range m.recentTransactions {
		if i >= 5 { // Show only last 5 transactions
			break
		}
		
		typeColor := "46" // Green for income
		amountPrefix := "+"
		typeSymbol := "↗️"
		if tx.Type == domain.TypeExpense {
			typeColor = "196" // Red for expense
			amountPrefix = "-"
			typeSymbol = "↘️"
		}
		
		dateStr := tx.Date.Format("Jan 02")
		typeStr := typeSymbol + " " + string(tx.Type)
		description := tx.Description
		if len(description) > 20 {
			description = description[:17] + "..."
		}
		amountStr := fmt.Sprintf("%s$%.2f", amountPrefix, tx.Amount)
		
		amountStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(typeColor))
		
		row := fmt.Sprintf("%-12s %-12s %-20s %s",
			dateStr, typeStr, description, amountStyle.Render(amountStr))
		
		rows = append(rows, row)
	}
	
	transactionList := lipgloss.JoinVertical(lipgloss.Left, rows...)
	
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("📝 Recent Transactions"),
		transactionList,
	)
}

// renderNavigation renders the navigation options
func (m *DashboardModel) renderNavigation() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true).
		Margin(1, 0, 0, 0)
	
	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(0, 0, 0, 2)
	
	options := []string{
		"2 - Add New Transaction",
		"3 - View All Transactions",
		"4 - Manage Categories",
		"? - Show Help",
		"q - Quit Application",
	}
	
	var optionStrings []string
	for _, option := range options {
		optionStrings = append(optionStrings, optionStyle.Render(option))
	}
	
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("🎯 Quick Actions"),
		lipgloss.JoinVertical(lipgloss.Left, optionStrings...),
	)
}

// renderLoading renders the loading state
func (m *DashboardModel) renderLoading() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(2, 0).
		Render("Loading dashboard data...")
}

// renderError renders the error state
func (m *DashboardModel) renderError() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Margin(2, 0).
		Render("Error loading dashboard: " + m.error)
}

// loadMonthlyData loads monthly financial data
func (m *DashboardModel) loadMonthlyData() tea.Cmd {
	return func() tea.Msg {
		now := time.Now()
		income, expense, balance, err := m.transactionUseCase.GetMonthlyBalance(now.Year(), int(now.Month()))
		if err != nil {
			return dataLoadErrorMsg{error: err.Error()}
		}
		
		return monthlyDataLoadedMsg{
			income:  income,
			expense: expense,
			balance: balance,
		}
	}
}

// loadRecentTransactions loads recent transactions
func (m *DashboardModel) loadRecentTransactions() tea.Cmd {
	return func() tea.Msg {
		transactions, err := m.transactionUseCase.GetRecentTransactions(30) // Last 30 days
		if err != nil {
			return dataLoadErrorMsg{error: err.Error()}
		}
		
		return recentTransactionsLoadedMsg{transactions: transactions}
	}
}

// Custom messages for data loading
type monthlyDataLoadedMsg struct {
	income  float64
	expense float64
	balance float64
}

type recentTransactionsLoadedMsg struct {
	transactions []*domain.Transaction
}

type dataLoadErrorMsg struct {
	error string
}

// abs returns the absolute value of a float64
func abs(x float64) float64 {
	return math.Abs(x)
}