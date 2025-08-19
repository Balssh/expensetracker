package tui

import (
	"context"
	"fmt"
	"strings"

	"expense-tracker/internal/core/domain"
	"expense-tracker/internal/core/usecase"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type transactionsMsg struct {
	transactions []*domain.Transaction
	err          error
}

type TransactionsModel struct {
	summaryUseCase *usecase.SummaryUseCase
	transactions   []*domain.Transaction
	searchInput    textinput.Model
	isSearching    bool
	loading        bool
	err            error
	currentPage    int
	itemsPerPage   int
	width          int
	height         int
}

func NewTransactionsModel(summaryUseCase *usecase.SummaryUseCase) *TransactionsModel {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search transactions..."

	return &TransactionsModel{
		summaryUseCase: summaryUseCase,
		searchInput:    searchInput,
		itemsPerPage:   20,
		currentPage:    0,
	}
}

func (m *TransactionsModel) Init() tea.Cmd {
	m.loading = true
	return m.fetchTransactions()
}

// SetDimensions updates the model's width and height for responsive layout
func (m *TransactionsModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

func (m *TransactionsModel) fetchTransactions() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()

		var transactions []*domain.Transaction
		var err error

		if m.searchInput.Value() != "" {
			transactions, err = m.summaryUseCase.SearchTransactions(
				ctx,
				m.searchInput.Value(),
				m.currentPage*m.itemsPerPage,
				m.itemsPerPage,
			)
		} else {
			transactions, err = m.summaryUseCase.GetAllTransactions(
				ctx,
				m.currentPage*m.itemsPerPage,
				m.itemsPerPage,
			)
		}

		return transactionsMsg{transactions: transactions, err: err}
	})
}

func (m *TransactionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case transactionsMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.transactions = msg.transactions
			m.err = nil
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, nil

		case "/":
			if !m.isSearching {
				m.isSearching = true
				m.searchInput.Focus()
				return m, nil
			}

		case "enter":
			if m.isSearching {
				m.isSearching = false
				m.searchInput.Blur()
				m.currentPage = 0
				m.loading = true
				return m, m.fetchTransactions()
			}

		case "c":
			if !m.isSearching {
				m.searchInput.SetValue("")
				m.currentPage = 0
				m.loading = true
				return m, m.fetchTransactions()
			}

		case "n":
			if !m.isSearching && len(m.transactions) == m.itemsPerPage {
				m.currentPage++
				m.loading = true
				return m, m.fetchTransactions()
			}

		case "p":
			if !m.isSearching && m.currentPage > 0 {
				m.currentPage--
				m.loading = true
				return m, m.fetchTransactions()
			}
		}

		if m.isSearching {
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m *TransactionsModel) View() string {
	if m.loading {
		centered := CenterHorizontally(loadingStyle.Render("Loading transactions..."), NewCenterConfig(m.width, m.height))
		return CreateMainApplicationBorder(centered, NewCenterConfig(m.width, m.height))
	}

	// Create the three sections
	searchPanel := m.createSearchPanel()
	transactionsPanel := m.createTransactionsPanel()
	helpPanel := m.createTransactionsHelpPanel()

	// Combine sections using the layout system
	config := NewCenterConfig(m.width, m.height)
	panelLayout := CreateThreePanelLayout(searchPanel, transactionsPanel, helpPanel, config)
	
	// Add the main title at the top
	title := titleStyle.Render("📋 All Transactions")
	
	// Center everything horizontally
	fullContent := title + "\n\n" + panelLayout
	centered := CenterHorizontally(fullContent, config)
	
	return CreateMainApplicationBorder(centered, config)
}

// createSearchPanel creates the top panel with search functionality
func (m *TransactionsModel) createSearchPanel() string {
	var b strings.Builder
	
	// Header
	header := summaryHeaderStyle.Render("🔍 Search & Filter")
	b.WriteString(header + "\n\n")
	
	if m.err != nil {
		b.WriteString(errorStyle.Render("❌ " + m.err.Error()) + "\n\n")
	}
	
	// Search input with better styling
	searchLabel := formFieldLabelStyle.Render("Search:")
	
	var searchField string
	if m.isSearching {
		searchField = searchBoxFocusedStyle.Render(m.searchInput.View())
	} else {
		searchValue := m.searchInput.Value()
		if searchValue == "" {
			searchValue = inputPlaceholderStyle.Render("Press '/' to search transactions...")
		}
		searchField = searchBoxStyle.Render(searchValue)
	}
	
	b.WriteString(searchLabel + "\n" + searchField)
	
	// Show active filters
	if m.searchInput.Value() != "" {
		b.WriteString("\n\n" + activeFilterTagStyle.Render("🔍 " + m.searchInput.Value()))
	}
	
	return b.String()
}

// createTransactionsPanel creates the middle panel with the transaction table
func (m *TransactionsModel) createTransactionsPanel() string {
	var b strings.Builder
	
	// Header with count
	transactionCount := len(m.transactions)
	header := summaryHeaderStyle.Render(fmt.Sprintf("📊 Transactions (%d)", transactionCount))
	b.WriteString(header + "\n\n")
	
	if len(m.transactions) == 0 {
		var emptyMessage string
		if m.searchInput.Value() != "" {
			emptyMessage = "No transactions found matching your search."
		} else {
			emptyMessage = "No transactions found. Press 'q' to go back and add some!"
		}
		b.WriteString(helpStyle.Render(emptyMessage))
		return b.String()
	}
	
	// Create properly aligned table
	b.WriteString(m.renderEnhancedTransactionTable())
	
	// Pagination info
	if m.currentPage > 0 || len(m.transactions) == m.itemsPerPage {
		b.WriteString("\n" + m.renderPaginationInfo())
	}
	
	return b.String()
}

// createTransactionsHelpPanel creates the bottom panel with keybindings
func (m *TransactionsModel) createTransactionsHelpPanel() string {
	var helpTexts []string
	
	if m.isSearching {
		helpTexts = []string{
			helpKeyStyle.Render("Type") + " to search",
			helpKeyStyle.Render("Enter") + " Apply",
			helpKeyStyle.Render("Esc") + " Cancel",
		}
	} else {
		helpTexts = []string{
			helpKeyStyle.Render("/") + " Search",
			helpKeyStyle.Render("c") + " Clear",
		}
		
		if m.currentPage > 0 {
			helpTexts = append(helpTexts, helpKeyStyle.Render("p") + " Previous")
		}
		if len(m.transactions) == m.itemsPerPage {
			helpTexts = append(helpTexts, helpKeyStyle.Render("n") + " Next")
		}
		
		helpTexts = append(helpTexts, helpKeyStyle.Render("q") + " Back")
	}
	
	return navigationStyle.Render(strings.Join(helpTexts, " • "))
}

// renderEnhancedTransactionTable creates a properly aligned transaction table
func (m *TransactionsModel) renderEnhancedTransactionTable() string {
	// Define responsive column widths based on terminal size
	breakpoint := GetBreakpoint(m.width)
	
	var columns []TableColumn
	switch breakpoint {
	case BreakpointNarrow:
		columns = []TableColumn{
			{Header: "Date", Width: 8, Alignment: lipgloss.Left},
			{Header: "Description", Width: 18, Alignment: lipgloss.Left},
			{Header: "Amount", Width: 10, Alignment: lipgloss.Right},
		}
	case BreakpointStandard:
		columns = []TableColumn{
			{Header: "Date", Width: 12, Alignment: lipgloss.Left},
			{Header: "Category", Width: 15, Alignment: lipgloss.Left},
			{Header: "Description", Width: 20, Alignment: lipgloss.Left},
			{Header: "Type", Width: 8, Alignment: lipgloss.Left},
			{Header: "Amount", Width: 12, Alignment: lipgloss.Right},
		}
	default: // BreakpointWide
		columns = []TableColumn{
			{Header: "Date", Width: 14, Alignment: lipgloss.Left},
			{Header: "Category", Width: 18, Alignment: lipgloss.Left},
			{Header: "Description", Width: 25, Alignment: lipgloss.Left},
			{Header: "Type", Width: 10, Alignment: lipgloss.Left},
			{Header: "Amount", Width: 15, Alignment: lipgloss.Right},
		}
	}
	
	var b strings.Builder
	
	// Table header
	tableHeader := CreateTableHeader(columns)
	b.WriteString(tableHeader + "\n")
	
	// Calculate total width for separator
	totalWidth := 0
	for _, col := range columns {
		totalWidth += col.Width + 1 // +1 for space between columns
	}
	
	// Separator
	b.WriteString(CreateTableSeparator(totalWidth-1) + "\n")
	
	// Transaction rows with alternating styles
	for i, transaction := range m.transactions {
		values := m.getTransactionRowValues(transaction, columns)
		row := FormatTableRow(columns, values)
		
		// Apply alternating row styles for better readability
		if i%2 == 0 {
			b.WriteString(tableRowStyle.Render(row))
		} else {
			b.WriteString(tableRowAltStyle.Render(row))
		}
		b.WriteString("\n")
	}
	
	return b.String()
}

// getTransactionRowValues formats transaction data for table display
func (m *TransactionsModel) getTransactionRowValues(transaction *domain.Transaction, columns []TableColumn) []string {
	categoryName := "Uncategorized"
	if transaction.Category != nil {
		categoryName = transaction.Category.Name
	}
	
	// Format values based on available columns
	values := make([]string, len(columns))
	
	for i, col := range columns {
		switch col.Header {
		case "Date":
			if col.Width <= 8 {
				values[i] = transaction.Date.Format("Jan 02")
			} else {
				values[i] = transaction.Date.Format("Jan 02, 2006")
			}
		case "Category":
			values[i] = TruncateWithEllipsis(categoryName, col.Width)
		case "Description":
			values[i] = TruncateWithEllipsis(transaction.Description, col.Width)
		case "Type":
			values[i] = m.formatTransactionTypeColored(transaction.Type)
		case "Amount":
			values[i] = m.formatAmountColored(transaction.Amount, transaction.Type)
		}
	}
	
	return values
}

// formatTransactionTypeColored returns a colored transaction type
func (m *TransactionsModel) formatTransactionTypeColored(transactionType string) string {
	if transactionType == "income" {
		return incomeStyle.Render("Income")
	}
	return expenseStyle.Render("Expense")
}

// formatAmountColored returns a colored amount based on transaction type
func (m *TransactionsModel) formatAmountColored(amount float64, transactionType string) string {
	formatted := fmt.Sprintf("$%.2f", amount)
	if transactionType == "income" {
		return incomeStyle.Render("+" + formatted)
	}
	return expenseStyle.Render("-" + formatted)
}

// renderPaginationInfo creates pagination information display
func (m *TransactionsModel) renderPaginationInfo() string {
	pageInfo := fmt.Sprintf("Page %d", m.currentPage+1)
	if len(m.transactions) == m.itemsPerPage {
		pageInfo += " (more available)"
	}
	return infoStyle.Render(pageInfo)
}

