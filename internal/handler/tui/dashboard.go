package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"expense-tracker/internal/core/domain"
	"expense-tracker/internal/core/usecase"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type summaryMsg struct {
	summary      *domain.Summary
	transactions []*domain.Transaction
	err          error
}

type DashboardModel struct {
	summaryUseCase *usecase.SummaryUseCase
	summary        *domain.Summary
	transactions   []*domain.Transaction
	loading        bool
	err            error
	width          int
	height         int
}

func NewDashboardModel(summaryUseCase *usecase.SummaryUseCase) *DashboardModel {
	return &DashboardModel{
		summaryUseCase: summaryUseCase,
		loading:        true,
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	return m.fetchData()
}

func (m *DashboardModel) Refresh() tea.Cmd {
	m.loading = true
	return m.fetchData()
}

func (m *DashboardModel) fetchData() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()

		// Use the enhanced summary system with intelligent defaults (current month)
		summary, err := m.summaryUseCase.GetSummary(ctx)
		if err != nil {
			return summaryMsg{err: err}
		}

		transactions, err := m.summaryUseCase.GetRecentTransactions(ctx, 10)
		if err != nil {
			return summaryMsg{err: err}
		}

		return summaryMsg{
			summary:      summary,
			transactions: transactions,
		}
	})
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case summaryMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.summary = msg.summary
			m.transactions = msg.transactions
			m.err = nil
		}
		return m, nil
	}
	return m, nil
}

func (m *DashboardModel) View() string {
	if m.loading {
		centered := CenterHorizontally(loadingStyle.Render("Loading expense tracker..."), NewCenterConfig(m.width, m.height))
		return CreateMainApplicationBorder(centered, NewCenterConfig(m.width, m.height))
	}

	if m.err != nil {
		errorContent := titleStyle.Render("Expense Tracker") + "\n\n" +
			errorStyle.Render("Error: "+m.err.Error())
		centered := CenterHorizontally(errorContent, NewCenterConfig(m.width, m.height))
		return CreateMainApplicationBorder(centered, NewCenterConfig(m.width, m.height))
	}

	// Create the three panels
	topPanel := m.createSummaryPanel()
	middlePanel := m.createTransactionsPanel()
	bottomPanel := m.createHelpPanel()

	// Combine panels using the layout system
	config := NewCenterConfig(m.width, m.height)
	panelLayout := CreateThreePanelLayout(topPanel, middlePanel, bottomPanel, config)
	
	// Add the main title at the top
	title := titleStyle.Render("💰 Expense Tracker")
	
	// Center everything horizontally
	fullContent := title + "\n\n" + panelLayout
	centered := CenterHorizontally(fullContent, config)
	
	return CreateMainApplicationBorder(centered, config)
}

// createSummaryPanel creates the top panel with monthly summary and expense breakdown
func (m *DashboardModel) createSummaryPanel() string {
	var b strings.Builder
	
	// Header with current month/year
	now := time.Now()
	monthYear := fmt.Sprintf("%s %d", now.Month().String(), now.Year())
	header := summaryHeaderStyle.Render("📊 " + monthYear + " Summary")
	b.WriteString(header + "\n\n")
	
	// Financial summary with better formatting
	summaryLine1 := fmt.Sprintf("%-12s %s", "Income:", incomeStyle.Render(fmt.Sprintf("$%.2f", m.summary.TotalIncome)))
	summaryLine2 := fmt.Sprintf("%-12s %s", "Expenses:", expenseStyle.Render(fmt.Sprintf("$%.2f", m.summary.TotalExpense)))
	summaryLine3 := fmt.Sprintf("%-12s %s", "Balance:", m.formatBalance(m.summary.NetBalance))
	
	b.WriteString(summaryLine1 + "\n")
	b.WriteString(summaryLine2 + "\n")
	b.WriteString(summaryLine3 + "\n\n")
	
	// Add a simple expense breakdown bar (dummy data for now)
	b.WriteString(m.createExpenseBreakdownBar())
	
	return b.String()
}

// createExpenseBreakdownBar creates a visual representation of expense categories
func (m *DashboardModel) createExpenseBreakdownBar() string {
	// For now, create a dummy breakdown bar
	// TODO: Replace with actual category data
	categories := []struct {
		name       string
		percentage int
		color      lipgloss.Style
	}{
		{"Food", 35, expenseStyle},
		{"Transport", 20, warningStyle},
		{"Bills", 25, errorStyle},
		{"Other", 20, infoStyle},
	}
	
	var b strings.Builder
	b.WriteString("Expense Breakdown:\n")
	
	// Create a horizontal bar
	totalWidth := 50
	bar := ""
	
	for _, cat := range categories {
		segmentWidth := (cat.percentage * totalWidth) / 100
		if segmentWidth < 1 {
			segmentWidth = 1
		}
		segment := strings.Repeat("█", segmentWidth)
		bar += cat.color.Render(segment)
	}
	
	b.WriteString(bar + "\n")
	
	// Add legend
	for _, cat := range categories {
		legend := fmt.Sprintf("%s %d%%", cat.name, cat.percentage)
		b.WriteString(cat.color.Render("■") + " " + legend + "  ")
	}
	
	return b.String()
}

// createTransactionsPanel creates the middle panel with recent transactions
func (m *DashboardModel) createTransactionsPanel() string {
	var b strings.Builder
	
	// Header
	header := summaryHeaderStyle.Render("📋 Recent Transactions")
	b.WriteString(header + "\n\n")
	
	if len(m.transactions) == 0 {
		b.WriteString(helpStyle.Render("No transactions found. Press 'a' to add an expense or 'i' to add income."))
		return b.String()
	}
	
	// Create table with proper alignment
	columns := []TableColumn{
		{Header: "Date", Width: 12, Alignment: lipgloss.Left},
		{Header: "Category", Width: 15, Alignment: lipgloss.Left},
		{Header: "Description", Width: 20, Alignment: lipgloss.Left},
		{Header: "Type", Width: 8, Alignment: lipgloss.Left},
		{Header: "Amount", Width: 10, Alignment: lipgloss.Right},
	}
	
	// Table header
	tableHeader := CreateTableHeader(columns)
	b.WriteString(tableHeader + "\n")
	
	// Separator
	b.WriteString(CreateTableSeparator(65) + "\n")
	
	// Transaction rows
	for _, transaction := range m.transactions {
		categoryName := "Uncategorized"
		if transaction.Category != nil {
			categoryName = transaction.Category.Name
		}
		
		values := []string{
			transaction.Date.Format("Jan 02"),
			TruncateWithEllipsis(categoryName, 15),
			TruncateWithEllipsis(transaction.Description, 20),
			m.formatTransactionType(transaction.Type),
			fmt.Sprintf("$%.2f", transaction.Amount),
		}
		
		row := FormatTableRow(columns, values)
		b.WriteString(row + "\n")
	}
	
	return b.String()
}

// createHelpPanel creates the bottom panel with available keybindings
func (m *DashboardModel) createHelpPanel() string {
	helpText := m.createHelpText()
	return navigationStyle.Render(helpText)
}

// createHelpText creates context-aware help text
func (m *DashboardModel) createHelpText() string {
	keys := []struct {
		key  string
		desc string
	}{
		{"a", "Add Expense"},
		{"i", "Add Income"},
		{"l", "List All"},
		{"r", "Refresh"},
		{"?", "Help"},
		{"q", "Quit"},
	}
	
	var parts []string
	for _, k := range keys {
		part := helpKeyStyle.Render("("+k.key+")") + " " + k.desc
		parts = append(parts, part)
	}
	
	return strings.Join(parts, " • ")
}

// SetDimensions updates the model's width and height for responsive layout
func (m *DashboardModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

func (m *DashboardModel) formatBalance(balance float64) string {
	if balance >= 0 {
		return balancePositiveStyle.Render(fmt.Sprintf("$%.2f", balance))
	}
	return balanceNegativeStyle.Render(fmt.Sprintf("$%.2f", balance))
}

func (m *DashboardModel) formatTransactionType(transactionType string) string {
	if transactionType == "income" {
		return incomeStyle.Render("Income")
	}
	return expenseStyle.Render("Expense")
}
