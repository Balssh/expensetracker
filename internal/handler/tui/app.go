package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
	"github.com/yourusername/expense-tracker/internal/handler/tui/components"
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
	
	// Navigation
	navbar          *components.NavBar
	
	// View models
	dashboard       *DashboardModel
	addTransaction  *AddTransactionModel
	listTransactions *ListTransactionsModel
	categories      *CategoriesModel
	help           *HelpModel
	
	// Unified error and message handling
	errorHandler    *ErrorHandler
	toastManager    *ToastManager
	currentMessage  *UnifiedMessageMsg
}

// NewApp creates a new application instance
func NewApp(transactionUseCase *usecase.TransactionUseCase, categoryUseCase *usecase.CategoryUseCase) *App {
	app := &App{
		transactionUseCase: transactionUseCase,
		categoryUseCase:    categoryUseCase,
		state:              DashboardView,
		errorHandler:       NewErrorHandler(),
		toastManager:       NewToastManager(),
	}
	
	// Initialize navigation bar
	app.navbar = components.NewNavBar().SetTabs([]components.NavTab{
		{Key: "1", Label: "Dashboard", ViewID: int(DashboardView), Icon: "📊", Shortcut: "1"},
		{Key: "2", Label: "Add Transaction", ViewID: int(AddTransactionView), Icon: "➕", Shortcut: "2"},
		{Key: "3", Label: "Transactions", ViewID: int(ListTransactionsView), Icon: "📋", Shortcut: "3"},
		{Key: "4", Label: "Categories", ViewID: int(CategoriesView), Icon: "📁", Shortcut: "4"},
		{Key: "?", Label: "Help", ViewID: int(HelpView), Icon: "❓", Shortcut: "?"},
	})
	
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
		a.tickToasts(), // Start the toast update ticker
	)
}

// tickToasts creates a command that periodically updates toast notifications
func (a *App) tickToasts() tea.Cmd {
	return tea.Every(time.Millisecond*100, func(t time.Time) tea.Msg {
		return toastTickMsg{}
	})
}

// Update handles messages and updates the application state
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		
		// Update navbar size
		a.navbar.SetWidth(msg.Width)
		
		// Update all view models with new size (accounting for navbar space)
		contentHeight := msg.Height - 3 // Account for navbar space
		a.dashboard.SetSize(msg.Width, contentHeight)
		a.addTransaction.SetSize(msg.Width, contentHeight)
		a.listTransactions.SetSize(msg.Width, contentHeight)
		a.categories.SetSize(msg.Width, contentHeight)
		a.help.SetSize(msg.Width, contentHeight)
		
	case tea.KeyMsg:
		// Handle global key bindings
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "?", "h":
			if a.state != HelpView {
				a.state = HelpView
				a.navbar.SetActiveByViewID(int(HelpView))
			} else {
				a.state = DashboardView
				a.navbar.SetActiveByViewID(int(DashboardView))
			}
			return a, nil
		case "esc":
			// Return to dashboard from any view
			a.state = DashboardView
			a.navbar.SetActiveByViewID(int(DashboardView))
			a.clearMessages()
			return a, a.dashboard.Init()
		case "1":
			// Allow global navigation to Dashboard
			a.state = DashboardView
			a.navbar.SetActiveByViewID(int(DashboardView))
			a.clearMessages()
			return a, a.dashboard.Init()
		case "2":
			// Allow global navigation to Add Transaction
			a.state = AddTransactionView
			a.navbar.SetActiveByViewID(int(AddTransactionView))
			a.clearMessages()
			return a, a.addTransaction.Init()
		case "3":
			// Allow global navigation to List Transactions
			a.state = ListTransactionsView
			a.navbar.SetActiveByViewID(int(ListTransactionsView))
			a.clearMessages()
			return a, a.listTransactions.Init()
		case "4":
			// Allow global navigation to Categories
			a.state = CategoriesView
			a.navbar.SetActiveByViewID(int(CategoriesView))
			a.clearMessages()
			return a, a.categories.Init()
		}
		
	case NavigateMsg:
		a.state = msg.View
		a.navbar.SetActiveByViewID(int(msg.View))
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
		
	case components.NavigationRequestMsg:
		// Handle navigation requests from the navbar
		newState := ViewState(msg.ViewID)
		a.state = newState
		a.navbar.SetActiveByViewID(msg.ViewID)
		a.clearMessages()
		switch newState {
		case DashboardView:
			return a, a.dashboard.Init()
		case AddTransactionView:
			return a, a.addTransaction.Init()
		case ListTransactionsView:
			return a, a.listTransactions.Init()
		case CategoriesView:
			return a, a.categories.Init()
		case HelpView:
			return a, nil // Help view doesn't need initialization
		}
		
	case UnifiedMessageMsg:
		a.currentMessage = &msg
		
	case UnifiedErrorMsg:
		a.errorHandler.AddError(msg.Error)
		
	// Legacy message support (for backward compatibility during transition)
	case MessageMsg:
		a.currentMessage = &UnifiedMessageMsg{
			Title:       "Info",
			Message:     msg.Message,
			Severity:    ErrorSeverityInfo,
			Dismissible: true,
			AutoDismiss: 3 * time.Second,
		}
		
	case ErrorMsg:
		userErr := NewUserError("LEGACY_ERROR", "Error", msg.Message, CategoryGeneral)
		a.errorHandler.AddError(userErr)
		
	case toastTickMsg:
		// Update toast notifications
		a.toastManager.UpdateToasts()
		return a, a.tickToasts() // Continue ticking
		
	case toastMsg:
		// Add new toast notification
		a.toastManager.AddToast(msg.Toast)
	}
	
	// Update navbar
	navbar, navCmd := a.navbar.Update(msg)
	a.navbar = navbar
	if navCmd != nil {
		cmd = navCmd
	}
	
	// Update the current view
	switch a.state {
	case DashboardView:
		dashboard, dashCmd := a.dashboard.Update(msg)
		a.dashboard = dashboard.(*DashboardModel)
		if cmd == nil {
			cmd = dashCmd
		}
	case AddTransactionView:
		addTransaction, addCmd := a.addTransaction.Update(msg)
		a.addTransaction = addTransaction.(*AddTransactionModel)
		if cmd == nil {
			cmd = addCmd
		}
	case ListTransactionsView:
		listTransactions, listCmd := a.listTransactions.Update(msg)
		a.listTransactions = listTransactions.(*ListTransactionsModel)
		if cmd == nil {
			cmd = listCmd
		}
	case CategoriesView:
		categories, catCmd := a.categories.Update(msg)
		a.categories = categories.(*CategoriesModel)
		if cmd == nil {
			cmd = catCmd
		}
	case HelpView:
		help, helpCmd := a.help.Update(msg)
		a.help = help.(*HelpModel)
		if cmd == nil {
			cmd = helpCmd
		}
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
	
	// Add navigation bar
	navbar := a.navbar.View()
	
	// Add header
	header := a.renderHeader()
	
	// Add footer with navigation hints
	footer := a.renderFooter()
	
	// Add message/error display
	messageArea := a.renderMessages()
	
	// Add toast notifications
	toastArea := a.renderToasts()
	
	// Calculate available height for content
	contentHeight := a.height - lipgloss.Height(navbar) - lipgloss.Height(header) - lipgloss.Height(footer) - lipgloss.Height(messageArea) - lipgloss.Height(toastArea)
	
	// Style the content area
	contentStyle := lipgloss.NewStyle().
		Height(contentHeight).
		Width(a.width)
	
	content = contentStyle.Render(content)
	
	return lipgloss.JoinVertical(lipgloss.Left,
		navbar,
		header,
		messageArea,
		toastArea,
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
	// Since we now have a navbar, we can simplify the footer hints
	hints := "ESC: Dashboard • Ctrl+C/q: Quit"
	
	footerStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("240")).
		Foreground(lipgloss.Color("15")).
		Padding(0, 2).
		Width(a.width)
	
	return footerStyle.Render(hints)
}

// renderMessages renders success/error messages using unified error handling
func (a *App) renderMessages() string {
	var content string
	
	// Show current error if any
	if latestError := a.errorHandler.GetLatestError(); latestError != nil {
		errorStyle := lipgloss.NewStyle().
			Background(lipgloss.Color(latestError.GetSeverityColor())).
			Foreground(lipgloss.Color("15")).
			Padding(0, 2).
			Width(a.width)
		
		errorText := fmt.Sprintf("%s %s: %s", latestError.GetSeveritySymbol(), latestError.Title, latestError.Message)
		if latestError.Recoverable && latestError.RetryAction != "" {
			errorText += " • " + latestError.RetryAction
		}
		
		content = errorStyle.Render(errorText)
	} else if a.currentMessage != nil {
		// Show current message
		var messageStyle lipgloss.Style
		var symbol string
		
		switch a.currentMessage.Severity {
		case ErrorSeverityInfo:
			messageStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("39")).
				Foreground(lipgloss.Color("15"))
			symbol = "ℹ️"
		case ErrorSeverityWarning:
			messageStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("226")).
				Foreground(lipgloss.Color("0"))
			symbol = "⚠️"
		default: // Success
			messageStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("46")).
				Foreground(lipgloss.Color("0"))
			symbol = "✓"
		}
		
		messageStyle = messageStyle.Padding(0, 2).Width(a.width)
		messageText := fmt.Sprintf("%s %s", symbol, a.currentMessage.Message)
		if a.currentMessage.Title != "" && a.currentMessage.Title != "Info" {
			messageText = fmt.Sprintf("%s %s: %s", symbol, a.currentMessage.Title, a.currentMessage.Message)
		}
		
		content = messageStyle.Render(messageText)
	}
	
	return content
}

// renderToasts renders active toast notifications
func (a *App) renderToasts() string {
	activeToasts := a.toastManager.GetActiveToasts()
	if len(activeToasts) == 0 {
		return ""
	}
	
	var toastViews []string
	
	for _, toast := range activeToasts {
		// Create toast content
		toastContent := fmt.Sprintf("%s %s", toast.Icon, toast.Message)
		if toast.Title != "" {
			toastContent = fmt.Sprintf("%s %s: %s", toast.Icon, toast.Title, toast.Message)
		}
		
		// Add progress indicator for auto-dismissing toasts
		if toast.AutoDismiss > 0 && !toast.Dismissed {
			progressBar := a.renderProgressBar(toast.Progress, 20)
			toastContent += "\n" + progressBar
		}
		
		// Style the toast
		toastStyle := lipgloss.NewStyle().
			Background(lipgloss.Color(toast.GetSeverityColor())).
			Foreground(lipgloss.Color("15")).
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(toast.GetSeverityColor())).
			Width(a.width - 4). // Account for margins
			Margin(0, 1)
		
		toastViews = append(toastViews, toastStyle.Render(toastContent))
	}
	
	return lipgloss.JoinVertical(lipgloss.Left, toastViews...)
}

// renderProgressBar renders a simple progress bar for toast auto-dismiss
func (a *App) renderProgressBar(progress float64, width int) string {
	if progress <= 0 {
		return ""
	}
	
	filled := int(progress * float64(width))
	if filled > width {
		filled = width
	}
	
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Width(width)
	
	return progressStyle.Render(bar)
}

// clearMessages clears any displayed messages
func (a *App) clearMessages() {
	a.errorHandler.ClearErrors()
	a.currentMessage = nil
}

// Helper methods for creating toast notifications

// ShowSuccessToast shows a success toast notification
func (a *App) ShowSuccessToast(title, message string) tea.Cmd {
	a.toastManager.AddSuccessToast(title, message)
	return nil
}

// ShowWarningToast shows a warning toast notification
func (a *App) ShowWarningToast(title, message string) tea.Cmd {
	a.toastManager.AddWarningToast(title, message)
	return nil
}

// ShowErrorToast shows an error toast notification
func (a *App) ShowErrorToast(title, message string) tea.Cmd {
	a.toastManager.AddErrorToast(title, message)
	return nil
}

// ShowInfoToast shows an informational toast notification
func (a *App) ShowInfoToast(title, message string) tea.Cmd {
	toast := NewToastNotification(title, message, ErrorSeverityInfo, 4*time.Second)
	a.toastManager.AddToast(toast)
	return nil
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

// Toast message types
type toastTickMsg struct{}

type toastMsg struct {
	Toast *ToastNotification
}