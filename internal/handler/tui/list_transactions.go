package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/expense-tracker/internal/core/domain"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
	"github.com/yourusername/expense-tracker/internal/handler/tui/components"
)

// FilterType represents different transaction filter types
type FilterType int

const (
	FilterAll FilterType = iota
	FilterIncome
	FilterExpense
)

// DateRange represents a date range filter
type DateRange int

const (
	DateRangeAll DateRange = iota
	DateRangeLast30Days
	DateRangeLast90Days
	DateRangeThisMonth
	DateRangeLastMonth
	DateRangeCustom
)

// ListTransactionsModel represents the list transactions view
type ListTransactionsModel struct {
	transactionUseCase *usecase.TransactionUseCase
	categoryUseCase    *usecase.CategoryUseCase

	width  int
	height int

	// Responsive configuration
	responsiveConfig *components.ResponsiveConfig
	currentLayout    components.LayoutConfig

	// Components
	table             *components.Table
	filterToggle      *components.Toggle
	dateRangeDropdown *components.Dropdown
	searchInput       *components.Input

	// Data
	transactions []*domain.Transaction
	categories   map[int]*domain.Category // ID -> Category mapping

	// Filter state
	currentFilter    FilterType
	currentDateRange DateRange
	searchTerm       string
	customStartDate  time.Time
	customEndDate    time.Time

	// UI state
	focusedComponent    int // 0=table, 1=filter, 2=date, 3=search
	loading             bool
	error               string
	showingDetails      bool
	selectedTransaction *domain.Transaction
}

// NewListTransactionsModel creates a new list transactions model
func NewListTransactionsModel(transactionUseCase *usecase.TransactionUseCase, categoryUseCase *usecase.CategoryUseCase) *ListTransactionsModel {
	responsiveConfig := components.NewResponsiveConfig()
	
	m := &ListTransactionsModel{
		transactionUseCase: transactionUseCase,
		categoryUseCase:    categoryUseCase,
		categories:         make(map[int]*domain.Category),
		currentFilter:      FilterAll,
		currentDateRange:   DateRangeLast30Days,
		responsiveConfig:   responsiveConfig,
		currentLayout:      responsiveConfig.GetLayoutForWidth(80), // Default to standard layout
	}

	m.initializeComponents()
	return m
}

// initializeComponents sets up the UI components
func (m *ListTransactionsModel) initializeComponents() {
	// Create table with transaction columns (optimized for smaller screens)
	m.table = components.NewTable().
		SetSelectable(true).
		SetPageSize(10). // Reduced from 15 to save vertical space
		SetColumns([]components.TableColumn{
			{Key: "date", Header: "Date", Width: 8, MinWidth: 6, MaxWidth: 12, Sortable: true, Align: lipgloss.Left},
			{Key: "type", Header: "Type", Width: 6, MinWidth: 4, MaxWidth: 8, Sortable: true, Align: lipgloss.Center},
			{Key: "category", Header: "Category", Width: 10, MinWidth: 6, MaxWidth: 15, Sortable: true, Align: lipgloss.Left},
			{Key: "description", Header: "Description", Width: 20, MinWidth: 12, MaxWidth: 35, Sortable: false, Align: lipgloss.Left},
			{Key: "amount", Header: "Amount", Width: 9, MinWidth: 7, MaxWidth: 12, Sortable: true, Align: lipgloss.Right},
		})

	// Filter toggle
	m.filterToggle = components.NewToggle("Filter").
		SetHorizontal(true).
		SetOptions([]components.ToggleOption{
			{Label: "All", Value: FilterAll, Color: "15"},
			{Label: "Income", Value: FilterIncome, Color: "46"},
			{Label: "Expense", Value: FilterExpense, Color: "196"},
		})

	// Date range dropdown
	dropdown := components.NewDropdown("Date Range").
		SetWidth(20).
		SetOptions([]components.DropdownOption{
			{Label: "All Time", Value: DateRangeAll},
			{Label: "Last 30 Days", Value: DateRangeLast30Days},
			{Label: "Last 90 Days", Value: DateRangeLast90Days},
			{Label: "This Month", Value: DateRangeThisMonth},
			{Label: "Last Month", Value: DateRangeLastMonth},
		})
	dropdown.SetSelected(1) // Default to Last 30 Days
	m.dateRangeDropdown = dropdown

	// Search input
	m.searchInput = components.NewInput("Search").
		SetPlaceholder("Search descriptions...").
		SetWidth(25)
}

// adjustTableColumns dynamically adjusts column widths using responsive configuration
func (m *ListTransactionsModel) adjustTableColumns(screenWidth int) {
	// Update current layout based on screen width
	m.currentLayout = m.responsiveConfig.GetLayoutForWidth(screenWidth)
	
	// Get current table columns
	columns := m.table.Columns
	if len(columns) == 0 {
		return
	}
	
	// Use responsive config to adjust columns
	adjustedColumns := m.responsiveConfig.AdjustTableColumns(columns, screenWidth, m.currentLayout)
	
	// Update table with adjusted columns
	m.table.SetColumns(adjustedColumns)
}

// SetSize sets the size of the list transactions view
func (m *ListTransactionsModel) SetSize(width, height int) {
	m.width = width
	m.height = height

	// The height we receive is the TOTAL terminal height, but the App component
	// will constrain our content to fit within the available content area.
	// We need to plan our layout to fit within a reasonable content height.
	
	// Estimate available content height after App header/footer (typically 4-6 lines total)
	appOverhead := 6 // Conservative estimate for header + footer + messages
	availableHeight := height - appOverhead
	if availableHeight < 10 {
		availableHeight = 10 // Minimum usable height
	}

	// Calculate space for our internal components:
	// - Title: 1 line 
	// - Filter controls: 1-3 lines (depending on screen width and layout) 
	// - Help text: 1 line (when shown)
	// - Table: remaining space
	filterLines := 1
	if width < 80 {
		filterLines = 3 // Vertical layout takes more space
	}
	tableHeight := availableHeight - 1 - filterLines - 1 // title - filters - help
	if tableHeight < 3 {
		tableHeight = 3 // Absolute minimum table height
	}

	// Adjust table size with responsive column widths
	m.adjustTableColumns(width)
	m.table.SetSize(width-2, tableHeight)
}

// Init initializes the list transactions view
func (m *ListTransactionsModel) Init() tea.Cmd {
	m.loading = true
	m.error = ""
	m.showingDetails = false
	m.focusedComponent = 0 // Focus on table by default

	// Focus the table
	m.table.Focus()

	// Apply initial sizing if width is available
	if m.width > 0 {
		m.adjustTableColumns(m.width)
		m.adjustFilterComponentWidths()
	}

	return tea.Batch(
		m.loadCategories(),
		m.loadTransactions(),
	)
}

// Update handles messages and updates the list transactions state
func (m *ListTransactionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle loading states
	if m.loading {
		switch msg := msg.(type) {
		case categoriesLoadedMsg:
			for _, category := range msg.categories {
				m.categories[category.ID] = category
			}
			return m, nil
		case transactionsLoadedMsg:
			m.transactions = msg.transactions
			m.updateTableData()
			m.loading = false
			return m, nil
		case dataLoadErrorMsg:
			m.error = msg.error
			m.loading = false
			return m, nil
		}
		return m, nil
	}

	// Handle transaction details view
	if m.showingDetails {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				m.showingDetails = false
				return m, nil
			case "e":
				// TODO: Implement edit transaction
				return m, nil
			case "d":
				// TODO: Implement delete transaction with confirmation
				return m, nil
			}
		}
		return m, nil
	}

	// Handle global keys
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			return m, m.nextComponent()
		case "shift+tab":
			return m, m.prevComponent()
		case "enter":
			if m.focusedComponent == 0 { // Table focused
				// Show transaction details
				if selectedRow := m.table.GetSelectedRow(); selectedRow != nil {
					m.showTransactionDetails(selectedRow.ID)
				}
			}
			return m, nil
		case "r":
			// Refresh data
			m.loading = true
			return m, tea.Batch(m.loadCategories(), m.loadTransactions())
		case "f":
			// Focus filter controls
			m.focusedComponent = 1
			return m, m.updateFocus()
		case "/":
			// Focus search
			m.focusedComponent = 3
			return m, m.updateFocus()
		}
	}

	// Update focused component
	switch m.focusedComponent {
	case 0: // Table
		table, tableCmd := m.table.Update(msg)
		m.table = table
		cmd = tableCmd

	case 1: // Filter toggle
		filterToggle, filterCmd := m.filterToggle.Update(msg)
		m.filterToggle = filterToggle
		cmd = filterCmd

		// Check if filter changed
		if newFilter := m.filterToggle.GetSelectedValue().(FilterType); newFilter != m.currentFilter {
			m.currentFilter = newFilter
			m.updateTableData()
		}

	case 2: // Date range dropdown
		dateDropdown, dateCmd := m.dateRangeDropdown.Update(msg)
		m.dateRangeDropdown = dateDropdown
		cmd = dateCmd

		// Check if date range changed
		if newRange := m.dateRangeDropdown.GetSelectedValue().(DateRange); newRange != m.currentDateRange {
			m.currentDateRange = newRange
			m.updateTableData()
		}

	case 3: // Search input
		searchInput, searchCmd := m.searchInput.Update(msg)
		m.searchInput = searchInput
		cmd = searchCmd

		// Check if search term changed
		if newTerm := m.searchInput.GetValue(); newTerm != m.searchTerm {
			m.searchTerm = newTerm
			m.updateTableData()
		}
	}

	return m, cmd
}

// View renders the list transactions view
func (m *ListTransactionsModel) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.error != "" {
		return m.renderError()
	}

	if m.showingDetails {
		return m.renderTransactionDetails()
	}

	return m.renderTransactionList()
}

// renderLoading renders the loading state
func (m *ListTransactionsModel) renderLoading() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(2, 0).
		Render("Loading transactions...")
}

// renderError renders the error state
func (m *ListTransactionsModel) renderError() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Margin(2, 0).
		Render("Error loading transactions: " + m.error)
}

// renderTransactionList renders the main transaction list view
func (m *ListTransactionsModel) renderTransactionList() string {
	var sections []string

	// Title (compact to save space)
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)
	sections = append(sections, titleStyle.Render("📊 Transaction History"))

	// Filter controls
	filterSection := m.renderFilterControls()
	sections = append(sections, filterSection)

	// Table (this will be sized appropriately by SetSize)
	tableView := m.table.View()
	sections = append(sections, tableView)

	// Help text based on responsive configuration
	if m.currentLayout.ShowHelp {
		// Calculate the current total height and be conservative
		currentHeight := 1 + lipgloss.Height(filterSection) + lipgloss.Height(tableView)
		appOverhead := 10 // Conservative estimate to ensure no overflow
		
		// Only show help text if we have sufficient space
		if m.height > 20 && currentHeight < m.height-appOverhead {
			helpStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

			var helpText string
			breakpointName := m.responsiveConfig.GetBreakpointName(m.width)
			switch breakpointName {
			case "mobile":
				helpText = "↵:Details Tab:Nav r:Refresh"
			case "narrow":
				helpText = "↵: Details • Tab: Nav • r: Refresh • f: Filter"
			default: // standard, wide
				helpText = "Enter: Details • Tab: Navigate • r: Refresh • f: Filter • /: Search"
			}
			sections = append(sections, helpStyle.Render(helpText))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// renderFilterControls renders the filter controls with responsive layout
func (m *ListTransactionsModel) renderFilterControls() string {
	// Style for focused components (very compact to prevent overflow)
	focusStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true) // Use bold instead of border to save space

	// Adjust component widths based on screen size
	m.adjustFilterComponentWidths()

	// Filter toggle
	filterView := m.filterToggle.View()
	if m.focusedComponent == 1 {
		filterView = focusStyle.Render(filterView)
	}

	// Date range dropdown
	dateView := m.dateRangeDropdown.View()
	if m.focusedComponent == 2 {
		dateView = focusStyle.Render(dateView)
	}

	// Search input
	searchView := m.searchInput.View()
	if m.focusedComponent == 3 {
		searchView = focusStyle.Render(searchView)
	}

	// Layout based on responsive configuration
	var filterRow string
	switch m.currentLayout.FilterLayout {
	case components.FilterLayoutVertical:
		// Stack vertically on narrow screens
		filterRow = lipgloss.JoinVertical(lipgloss.Left,
			filterView,
			dateView,
			searchView,
		)
	case components.FilterLayoutCollapsed:
		// Show only essential filters for very narrow screens
		filterRow = lipgloss.JoinVertical(lipgloss.Left,
			filterView,
			searchView,
		)
	default: // FilterLayoutHorizontal
		// Layout horizontally on wider screens
		filterRow = lipgloss.JoinHorizontal(lipgloss.Top,
			filterView, " ", // Single space instead of double
			dateView, " ",
			searchView,
		)
	}
	
	// More compact margin to save vertical space
	return lipgloss.NewStyle().Margin(0, 0, 0, 0).Render(filterRow)
}

// adjustFilterComponentWidths adjusts the widths of filter components using responsive configuration
func (m *ListTransactionsModel) adjustFilterComponentWidths() {
	// Use current layout's component width configuration
	if dropdownWidth, exists := m.currentLayout.ComponentWidths["dropdown"]; exists {
		m.dateRangeDropdown.SetWidth(dropdownWidth)
	}
	
	if searchWidth, exists := m.currentLayout.ComponentWidths["search"]; exists {
		m.searchInput.SetWidth(searchWidth)
	}
}

// renderTransactionDetails renders the transaction details view
func (m *ListTransactionsModel) renderTransactionDetails() string {
	if m.selectedTransaction == nil {
		return "No transaction selected"
	}

	tx := m.selectedTransaction

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true).
		Margin(0, 0, 1, 0)

	// Responsive label width based on screen size
	labelWidth := 12
	if m.width < 60 {
		labelWidth = 8
	} else if m.width < 80 {
		labelWidth = 10
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Bold(true).
		Width(labelWidth)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	// Get category name
	categoryName := "Unknown"
	if category, exists := m.categories[tx.CategoryID]; exists {
		categoryName = category.Name
	}

	// Amount color
	amountColor := "46" // Green for income
	if tx.Type == domain.TypeExpense {
		amountColor = "196" // Red for expense
	}

	// Truncate long description for narrow screens
	description := tx.Description
	maxDescLength := m.width - labelWidth - 10 // Account for label and padding
	if len(description) > maxDescLength && maxDescLength > 20 {
		description = description[:maxDescLength-3] + "..."
	}

	details := []string{
		titleStyle.Render("📋 Transaction Details"),
		"",
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("ID:"), valueStyle.Render(tx.ID)),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Type:"), valueStyle.Render(string(tx.Type))),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Amount:"),
			lipgloss.NewStyle().Foreground(lipgloss.Color(amountColor)).Render(tx.FormatAmount())),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Category:"), valueStyle.Render(categoryName)),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Date:"), valueStyle.Render(tx.FormatDate())),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Description:"), valueStyle.Render(description)),
		lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Created:"),
			valueStyle.Render(tx.CreatedAt.Format("2006-01-02 15:04:05"))),
		"",
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("ESC: Back • e: Edit • d: Delete"),
	}

	return lipgloss.JoinVertical(lipgloss.Left, details...)
}

// buildTransactionRows converts transactions to table rows
func (m *ListTransactionsModel) buildTransactionRows() []components.TableRow {
	var rows []components.TableRow

	for _, tx := range m.transactions {
		// Apply filters
		if !m.shouldIncludeTransaction(tx) {
			continue
		}

		// Get category name
		categoryName := "Unknown"
		if category, exists := m.categories[tx.CategoryID]; exists {
			categoryName = category.Name
		}

		// Format amount with +/- prefix and symbols for accessibility
		amountText := tx.FormatAmount()
		typeText := string(tx.Type)
		if tx.Type == domain.TypeIncome {
			amountText = "+" + amountText
			typeText = "↗️ " + typeText
		} else {
			amountText = "-" + amountText
			typeText = "↘️ " + typeText
		}

		// Determine colors
		colors := make(map[string]string)
		if tx.Type == domain.TypeIncome {
			colors["type"] = "46"   // Green for income type
			colors["amount"] = "46" // Green for income amount
		} else {
			colors["type"] = "196"   // Red for expense type
			colors["amount"] = "196" // Red for expense amount
		}

		// Create table row
		row := components.TableRow{
			ID: tx.ID,
			Data: map[string]interface{}{
				"date":        tx.Date.Format("Jan 02"),
				"type":        strings.Title(typeText),
				"category":    categoryName,
				"description": tx.Description,
				"amount":      amountText,
			},
			Colors: colors,
		}

		rows = append(rows, row)
	}

	return rows
}

// Rest of implementation...
func (m *ListTransactionsModel) shouldIncludeTransaction(tx *domain.Transaction) bool {
	// Filter by type
	switch m.currentFilter {
	case FilterIncome:
		if tx.Type != domain.TypeIncome {
			return false
		}
	case FilterExpense:
		if tx.Type != domain.TypeExpense {
			return false
		}
	}

	// Filter by date range
	now := time.Now()
	switch m.currentDateRange {
	case DateRangeLast30Days:
		if tx.Date.Before(now.AddDate(0, 0, -30)) {
			return false
		}
	case DateRangeLast90Days:
		if tx.Date.Before(now.AddDate(0, 0, -90)) {
			return false
		}
	case DateRangeThisMonth:
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		if tx.Date.Before(startOfMonth) {
			return false
		}
	case DateRangeLastMonth:
		startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
		if tx.Date.Before(startOfLastMonth) || tx.Date.After(startOfThisMonth) {
			return false
		}
	}

	// Filter by search term
	if m.searchTerm != "" {
		searchLower := strings.ToLower(m.searchTerm)
		if !strings.Contains(strings.ToLower(tx.Description), searchLower) {
			return false
		}
	}

	return true
}

// updateTableData refreshes the table with current filter settings
func (m *ListTransactionsModel) updateTableData() {
	rows := m.buildTransactionRows()
	m.table.SetRows(rows)
}

// showTransactionDetails shows details for the specified transaction
func (m *ListTransactionsModel) showTransactionDetails(transactionID string) {
	for _, tx := range m.transactions {
		if tx.ID == transactionID {
			m.selectedTransaction = tx
			m.showingDetails = true
			break
		}
	}
}

// Component navigation methods
func (m *ListTransactionsModel) nextComponent() tea.Cmd {
	m.blurAllComponents()
	m.focusedComponent = (m.focusedComponent + 1) % 4
	return m.updateFocus()
}

func (m *ListTransactionsModel) prevComponent() tea.Cmd {
	m.blurAllComponents()
	m.focusedComponent--
	if m.focusedComponent < 0 {
		m.focusedComponent = 3
	}
	return m.updateFocus()
}

func (m *ListTransactionsModel) updateFocus() tea.Cmd {
	switch m.focusedComponent {
	case 0:
		m.table.Focus()
	case 1:
		m.filterToggle.Focus()
	case 2:
		m.dateRangeDropdown.Focus()
	case 3:
		m.searchInput.Focus()
	}
	return nil
}

func (m *ListTransactionsModel) blurAllComponents() {
	m.table.Blur()
	m.filterToggle.Blur()
	m.dateRangeDropdown.Blur()
	m.searchInput.Blur()
}

// Data loading methods
func (m *ListTransactionsModel) loadCategories() tea.Cmd {
	return func() tea.Msg {
		categories, err := m.categoryUseCase.GetAllCategories()
		if err != nil {
			return dataLoadErrorMsg{error: err.Error()}
		}
		return categoriesLoadedMsg{categories: categories}
	}
}

func (m *ListTransactionsModel) loadTransactions() tea.Cmd {
	return func() tea.Msg {
		transactions, err := m.transactionUseCase.ListTransactions(1000, 0) // Load up to 1000 transactions
		if err != nil {
			return dataLoadErrorMsg{error: err.Error()}
		}
		return transactionsLoadedMsg{transactions: transactions}
	}
}

// Custom messages for list transactions
type transactionsLoadedMsg struct {
	transactions []*domain.Transaction
}
