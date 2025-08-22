package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
)

// AddTransactionModel - stub implementation
type AddTransactionModel struct {
	transactionUseCase *usecase.TransactionUseCase
	categoryUseCase    *usecase.CategoryUseCase
	width              int
	height             int
}

func NewAddTransactionModel(transactionUseCase *usecase.TransactionUseCase, categoryUseCase *usecase.CategoryUseCase) *AddTransactionModel {
	return &AddTransactionModel{
		transactionUseCase: transactionUseCase,
		categoryUseCase:    categoryUseCase,
	}
}

func (m *AddTransactionModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *AddTransactionModel) Init() tea.Cmd {
	return nil
}

func (m *AddTransactionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *AddTransactionModel) View() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(2, 0)
	
	return style.Render("Add Transaction view - Coming soon!\nPress ESC to return to dashboard.")
}

// ListTransactionsModel - stub implementation  
type ListTransactionsModel struct {
	transactionUseCase *usecase.TransactionUseCase
	categoryUseCase    *usecase.CategoryUseCase
	width              int
	height             int
}

func NewListTransactionsModel(transactionUseCase *usecase.TransactionUseCase, categoryUseCase *usecase.CategoryUseCase) *ListTransactionsModel {
	return &ListTransactionsModel{
		transactionUseCase: transactionUseCase,
		categoryUseCase:    categoryUseCase,
	}
}

func (m *ListTransactionsModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *ListTransactionsModel) Init() tea.Cmd {
	return nil
}

func (m *ListTransactionsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *ListTransactionsModel) View() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(2, 0)
	
	return style.Render("List Transactions view - Coming soon!\nPress ESC to return to dashboard.")
}

// CategoriesModel - stub implementation
type CategoriesModel struct {
	categoryUseCase *usecase.CategoryUseCase
	width           int
	height          int
}

func NewCategoriesModel(categoryUseCase *usecase.CategoryUseCase) *CategoriesModel {
	return &CategoriesModel{
		categoryUseCase: categoryUseCase,
	}
}

func (m *CategoriesModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *CategoriesModel) Init() tea.Cmd {
	return nil
}

func (m *CategoriesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *CategoriesModel) View() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(2, 0)
	
	return style.Render("Categories view - Coming soon!\nPress ESC to return to dashboard.")
}

// HelpModel - stub implementation
type HelpModel struct {
	width  int
	height int
}

func NewHelpModel() *HelpModel {
	return &HelpModel{}
}

func (m *HelpModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *HelpModel) Init() tea.Cmd {
	return nil
}

func (m *HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *HelpModel) View() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true).
		Margin(1, 0)
	
	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(0, 2)
	
	help := []string{
		"Global Navigation:",
		"  1 - Dashboard",
		"  2 - Add Transaction", 
		"  3 - View Transactions",
		"  4 - Manage Categories",
		"  ? - Toggle Help",
		"  ESC - Return to Dashboard",
		"  q - Quit Application",
		"",
		"Tips:",
		"  - Use arrow keys to navigate within views",
		"  - Tab to move between form fields",
		"  - Enter to confirm actions",
	}
	
	content := lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("💡 Help & Keyboard Shortcuts"),
		contentStyle.Render(lipgloss.JoinVertical(lipgloss.Left, help...)),
	)
	
	return content
}