package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
)

// ViewState represents the current view state of the application
type ViewState int

const (
	DashboardView ViewState = iota
	AddTransactionView
	ListTransactionsView
	CategoriesView
	HelpView
)

// App represents the main application model
type App struct {
	transactionUseCase *usecase.TransactionUseCase
	categoryUseCase    *usecase.CategoryUseCase
	
	state       ViewState
	width       int
	height      int
	
	// View models
	dashboard       *DashboardModel
	addTransaction  *AddTransactionModel
	listTransactions *ListTransactionsModel
	categories      *CategoriesModel
	help           *HelpModel
	
	// Error and message handling
	message string
	error   string
}

// NewApp creates a new application instance
func NewApp(transactionUseCase *usecase.TransactionUseCase, categoryUseCase *usecase.CategoryUseCase) *App {
	app := &App{
		transactionUseCase: transactionUseCase,
		categoryUseCase:    categoryUseCase,
		state:              DashboardView,
	}
	
	// Initialize view models
	app.dashboard = NewDashboardModel(transactionUseCase, categoryUseCase)
	app.addTransaction = NewAddTransactionModel(transactionUseCase, categoryUseCase)
	app.listTransactions = NewListTransactionsModel(transactionUseCase, categoryUseCase)
	app.categories = NewCategoriesModel(categoryUseCase)
	app.help = NewHelpModel()
	
	return app
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.dashboard.Init(),
		tea.EnterAltScreen,
	)
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		
		// Update all view models with new size
		a.dashboard.SetSize(msg.Width, msg.Height)
		a.addTransaction.SetSize(msg.Width, msg.Height)
		a.listTransactions.SetSize(msg.Width, msg.Height)
		a.categories.SetSize(msg.Width, msg.Height)
		a.help.SetSize(msg.Width, msg.Height)
		
	case tea.KeyMsg:
		// Handle global key bindings
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "?", "h":
			if a.state != HelpView {
				a.state = HelpView
			} else {
				a.state = DashboardView
			}
			return a, nil
		case "esc":
			// Return to dashboard from any view
			a.state = DashboardView
			a.clearMessages()
			return a, a.dashboard.Init()
		case "1":
			if a.state == DashboardView {
				a.state = DashboardView
				a.clearMessages()
				return a, a.dashboard.Init()
			}
		case "2":
			if a.state == DashboardView {
				a.state = AddTransactionView
				a.clearMessages()
				return a, a.addTransaction.Init()
			}
		case "3":
			if a.state == DashboardView {
				a.state = ListTransactionsView
				a.clearMessages()
				return a, a.listTransactions.Init()
			}
		case "4":
			if a.state == DashboardView {
				a.state = CategoriesView
				a.clearMessages()
				return a, a.categories.Init()
			}
		}
		
	case NavigateMsg:
		a.state = msg.View
		a.clearMessages()
		switch msg.View {
		case DashboardView:
			return a, a.dashboard.Init()
		case AddTransactionView:
			return a, a.addTransaction.Init()
		case ListTransactionsView:
			return a, a.listTransactions.Init()
		case CategoriesView:
			return a, a.categories.Init()
		}
		
	case MessageMsg:
		a.message = msg.Message
		
	case ErrorMsg:
		a.error = msg.Message
	}
	
	// Update the current view
	switch a.state {
	case DashboardView:
		dashboard, dashCmd := a.dashboard.Update(msg)
		a.dashboard = dashboard.(*DashboardModel)
		cmd = dashCmd
	case AddTransactionView:
		addTransaction, addCmd := a.addTransaction.Update(msg)
		a.addTransaction = addTransaction.(*AddTransactionModel)
		cmd = addCmd
	case ListTransactionsView:
		listTransactions, listCmd := a.listTransactions.Update(msg)
		a.listTransactions = listTransactions.(*ListTransactionsModel)
		cmd = listCmd
	case CategoriesView:
		categories, catCmd := a.categories.Update(msg)
		a.categories = categories.(*CategoriesModel)
		cmd = catCmd
	case HelpView:
		help, helpCmd := a.help.Update(msg)
		a.help = help.(*HelpModel)
		cmd = helpCmd
	}
	
	return a, cmd
}

// View renders the application view
func (a *App) View() string {
	if a.width == 0 || a.height == 0 {
		return "Loading..."
	}
	
	var content string
	
	// Render the current view
	switch a.state {
	case DashboardView:
		content = a.dashboard.View()
	case AddTransactionView:
		content = a.addTransaction.View()
	case ListTransactionsView:
		content = a.listTransactions.View()
	case CategoriesView:
		content = a.categories.View()
	case HelpView:
		content = a.help.View()
	}
	
	// Add header
	header := a.renderHeader()
	
	// Add footer with navigation hints
	footer := a.renderFooter()
	
	// Add message/error display
	messageArea := a.renderMessages()
	
	// Calculate available height for content
	contentHeight := a.height - lipgloss.Height(header) - lipgloss.Height(footer) - lipgloss.Height(messageArea)
	
	// Style the content area
	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		Width(a.width)
	
	content = contentStyle.Render(content)
	
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		messageArea,
		content,
		footer,
	)
}

// renderHeader renders the application header
func (a *App) renderHeader() string {
	title := "💰 Expense Tracker"
	
	var viewName string
	switch a.state {
	case DashboardView:
		viewName = "Dashboard"
	case AddTransactionView:
		viewName = "Add Transaction"
	case ListTransactionsView:
		viewName = "Transactions"
	case CategoriesView:
		viewName = "Categories"
	case HelpView:
		viewName = "Help"
	}
	
	headerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("15")).
		Padding(0, 2).
		Width(a.width).
		Bold(true)
	
	return headerStyle.Render(title + " - " + viewName)
}

// renderFooter renders the application footer
func (a *App) renderFooter() string {
	var hints string
	
	if a.state == DashboardView {
		hints = "1: Dashboard • 2: Add Transaction • 3: View Transactions • 4: Categories • ?: Help • q: Quit"
	} else {
		hints = "ESC: Back to Dashboard • ?: Help • q: Quit"
	}
	
	footerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("15")).
		Padding(0, 2).
		Width(a.width)
	
	return footerStyle.Render(hints)
}

// renderMessages renders success/error messages
func (a *App) renderMessages() string {
	if a.message == "" && a.error == "" {
		return ""
	}
	
	var content string
	if a.error != "" {
		errorStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("196")).
			Foreground(lipgloss.Color("15")).
			Padding(0, 2).
			Width(a.width)
		content = errorStyle.Render("Error: " + a.error)
	} else if a.message != "" {
		messageStyle := lipgloss.NewStyle().
			Background(lipgloss.Color("46")).
			Foreground(lipgloss.Color("0")).
			Padding(0, 2).
			Width(a.width)
		content = messageStyle.Render("✓ " + a.message)
	}
	
	return content
}

// clearMessages clears any displayed messages
func (a *App) clearMessages() {
	a.message = ""
	a.error = ""
}

// Custom messages for navigation and feedback
type NavigateMsg struct {
	View ViewState
}

type MessageMsg struct {
	Message string
}

type ErrorMsg struct {
	Message string
}