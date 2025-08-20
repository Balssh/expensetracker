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
	// Use full terminal dimensions and auto-scale content
	config := NewCenterConfig(m.width, m.height)
	
	if m.loading {
		content := loadingStyle.Render("Loading expense tracker...")
		return m.renderWithLayout(content, config)
	}

	if m.err != nil {
		content := titleStyle.Render("💰 Expense Tracker") + "\n\n" +
			errorStyle.Render("Error: "+m.err.Error())
		return m.renderWithLayout(content, config)
	}

	// Create responsive panels that auto-adjust to terminal size
	summaryPanel := m.createSummaryPanel()
	transactionsPanel := m.createTransactionsPanel()
	helpPanel := m.createHelpPanel()

	// Main title
	title := titleStyle.Render("💰 Expense Tracker")
	
	// Create full layout using Lipgloss vertical join for better alignment
	fullContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		m.createPanelWithBorder(summaryPanel, "📊 Summary", config),
		"",
		m.createPanelWithBorder(transactionsPanel, "📋 Recent Transactions", config),
		"",
		helpPanel, // No border for help panel as requested
	)
	
	// Center the entire layout
	return lipgloss.Place(config.Width, config.Height, 
		lipgloss.Center, lipgloss.Center, fullContent)
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
	
	// Calculate table width to fill the entire panel container
	config := NewCenterConfig(m.width, m.height)
	panelWidth := config.CalculateContentWidth() - 8 // Account for panel padding
	
	columns := []TableColumn{
		{Header: "Date", Width: 12, Alignment: lipgloss.Left},
		{Header: "Category", Width: panelWidth * 25 / 100, Alignment: lipgloss.Left},  
		{Header: "Description", Width: panelWidth * 50 / 100, Alignment: lipgloss.Left},
		{Header: "Amount", Width: panelWidth * 25 / 100, Alignment: lipgloss.Right},
	}
	
	// Ensure minimum widths and adjust if needed
	if columns[1].Width < 12 { columns[1].Width = 12 }
	if columns[2].Width < 15 { columns[2].Width = 15 }
	if columns[3].Width < 10 { columns[3].Width = 10 }
	
	// Table header
	tableHeader := CreateTableHeader(columns)
	b.WriteString(tableHeader + "\n")
	
	// Calculate separator width based on actual column widths
	totalWidth := 0
	for _, col := range columns {
		totalWidth += col.Width + 1 // +1 for space between columns
	}
	
	// Separator
	b.WriteString(CreateTableSeparator(totalWidth-1) + "\n")
	
	// Transaction rows
	for _, transaction := range m.transactions {
		categoryName := "Uncategorized"
		if transaction.Category != nil {
			categoryName = transaction.Category.Name
		}
		
		// Format amount with color coding based on transaction type
		var formattedAmount string
		if transaction.Type == "income" {
			formattedAmount = incomeStyle.Render(fmt.Sprintf("+$%.2f", transaction.Amount))
		} else {
			formattedAmount = expenseStyle.Render(fmt.Sprintf("-$%.2f", transaction.Amount))
		}
		
		values := []string{
			transaction.Date.Format("Jan 02"),
			TruncateWithEllipsis(categoryName, columns[1].Width),
			TruncateWithEllipsis(transaction.Description, columns[2].Width),
			formattedAmount,
		}
		
		row := FormatTableRow(columns, values)
		b.WriteString(row + "\n")
	}
	
	return b.String()
}

// createHelpPanel creates the bottom panel with available keybindings
func (m *DashboardModel) createHelpPanel() string {
	helpText := m.createHelpText()
	// Remove border styling as requested
	return helpText
}

// renderWithLayout handles simple content with responsive centering
func (m *DashboardModel) renderWithLayout(content string, config CenterConfig) string {
	return lipgloss.Place(config.Width, config.Height, 
		lipgloss.Center, lipgloss.Center, content)
}

// createPanelWithBorder creates a bordered panel that scales with terminal size
func (m *DashboardModel) createPanelWithBorder(content, title string, config CenterConfig) string {
	// Calculate consistent panel dimensions 
	panelWidth := config.CalculateContentWidth() - 4 // Account for border padding
	panelHeight := (config.Height - 10) / 3 - 2 // Divide available height into 3 equal sections
	
	// Ensure minimum height
	if panelHeight < 6 {
		panelHeight = 6
	}
	
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorNeutral).
		Width(panelWidth).
		Height(panelHeight).
		Padding(1, 2).
		Align(lipgloss.Left)
	
	return style.Render(content)
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
