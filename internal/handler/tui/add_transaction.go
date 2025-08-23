package tui

import (
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/yourusername/expense-tracker/internal/core/domain"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
	"github.com/yourusername/expense-tracker/internal/handler/tui/components"
)

// AddTransactionModel represents the add transaction form
type AddTransactionModel struct {
	transactionUseCase *usecase.TransactionUseCase
	categoryUseCase    *usecase.CategoryUseCase
	
	width  int
	height int
	
	// Form components
	form           *components.Form
	typeToggle     *components.Toggle
	amountInput    *components.Input
	categoryDropdown *components.Dropdown
	descriptionInput *components.Input
	dateInput      *components.Input
	
	// Data
	categories []*domain.Category
	
	// State
	loading     bool
	saving      bool
	saved       bool
	error       string
	showAddAnother bool
}

// NewAddTransactionModel creates a new add transaction model
func NewAddTransactionModel(transactionUseCase *usecase.TransactionUseCase, categoryUseCase *usecase.CategoryUseCase) *AddTransactionModel {
	m := &AddTransactionModel{
		transactionUseCase: transactionUseCase,
		categoryUseCase:    categoryUseCase,
	}
	
	m.initializeForm()
	return m
}

// initializeForm sets up the form components
func (m *AddTransactionModel) initializeForm() {
	// Create form
	m.form = components.NewForm("Add New Transaction")
	
	// Transaction type toggle
	m.typeToggle = components.NewToggle("Type").
		SetHorizontal(true).
		SetOptions([]components.ToggleOption{
			{Label: "Income", Value: domain.TypeIncome, Color: "46"},
			{Label: "Expense", Value: domain.TypeExpense, Color: "196"},
		}).
		SetRequired(true)
	
	// Amount input
	m.amountInput = components.NewInput("Amount").
		SetInputType(components.InputNumber).
		SetPlaceholder("0.00").
		SetWidth(15).
		SetRequired(true).
		SetValidation(components.ValidateAmount)
	
	// Category dropdown
	m.categoryDropdown = components.NewDropdown("Category").
		SetPlaceholder("Select a category").
		SetWidth(25).
		SetRequired(true)
	
	// Description input
	m.descriptionInput = components.NewInput("Description").
		SetPlaceholder("Optional note about this transaction").
		SetWidth(40).
		SetValidation(components.ValidateDescription)
	
	// Date input (default to today)
	m.dateInput = components.NewInput("Date").
		SetPlaceholder("YYYY-MM-DD").
		SetWidth(12).
		SetRequired(true).
		SetValidation(m.validateDate)
	
	// Set today's date as default
	m.dateInput.SetValue(time.Now().Format("2006-01-02"))
	
	// Set default transaction type to expense (more common)
	m.typeToggle.SetSelected(1) // Index 1 is expense
	
	// Add fields to form
	m.form.AddToggle(m.typeToggle).
		AddInput(m.amountInput).
		AddDropdown(m.categoryDropdown).
		AddInput(m.descriptionInput).
		AddInput(m.dateInput)
	
	// Set form callbacks
	m.form.SetOnSubmit(m.handleSubmit).
		SetOnCancel(m.handleCancel).
		SetOnFieldChange(m.handleFieldChange)
}

// SetSize sets the size of the add transaction view
func (m *AddTransactionModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.form.SetSize(width, height-4) // Account for header/footer
}

// Init initializes the add transaction view
func (m *AddTransactionModel) Init() tea.Cmd {
	m.loading = true
	m.saved = false
	m.showAddAnother = false
	m.error = ""
	
	// Reset form
	m.form.Reset()
	
	// Load categories for the current transaction type
	return m.loadCategories()
}

// Update handles messages and updates the add transaction state
func (m *AddTransactionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	if m.loading {
		switch msg := msg.(type) {
		case categoriesLoadedMsg:
			m.categories = msg.categories
			m.updateCategoryDropdown()
			m.loading = false
			return m, nil
		case dataLoadErrorMsg:
			m.error = msg.error
			m.loading = false
			return m, nil
		}
		return m, nil
	}
	
	if m.saved && m.showAddAnother {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "a", "A":
				// Add another transaction
				m.saved = false
				m.showAddAnother = false
				m.form.Reset()
				m.form.ClearMessage()
				return m, m.loadCategories()
			case "esc", "q":
				// Return to dashboard
				return m, func() tea.Msg {
					return NavigateMsg{View: DashboardView}
				}
			}
		}
		return m, nil
	}
	
	if m.saving {
		switch msg := msg.(type) {
		case transactionSavedMsg:
			m.saving = false
			m.saved = true
			m.showAddAnother = true
			m.form.SetMessage("Transaction saved successfully!", components.MessageSuccess)
			return m, nil
		case saveErrorMsg:
			m.saving = false
			m.form.SetMessage("Error: "+msg.error, components.MessageError)
			return m, nil
		}
		return m, nil
	}
	
	// Check if transaction type changed before updating form
	oldType := m.typeToggle.GetSelectedValue().(domain.TransactionType)
	
	// Update form
	form, formCmd := m.form.Update(msg)
	m.form = form
	
	// Check if transaction type changed after update
	newType := m.typeToggle.GetSelectedValue().(domain.TransactionType)
	if oldType != newType {
		m.updateCategoryDropdown()
	}
	
	return m, formCmd
}

// View renders the add transaction view
func (m *AddTransactionModel) View() string {
	if m.loading {
		return m.renderLoading()
	}
	
	if m.error != "" {
		return m.renderError()
	}
	
	if m.saved && m.showAddAnother {
		return m.renderSuccess()
	}
	
	if m.saving {
		return m.renderSaving()
	}
	
	return m.form.View()
}

// renderLoading renders the loading state
func (m *AddTransactionModel) renderLoading() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(2, 0).
		Render("Loading categories...")
}

// renderError renders the error state
func (m *AddTransactionModel) renderError() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("196")).
		Margin(2, 0).
		Render("Error loading categories: " + m.error)
}

// renderSaving renders the saving state
func (m *AddTransactionModel) renderSaving() string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("39")).
		Margin(2, 0).
		Render("Saving transaction...")
}

// renderSuccess renders the success state with add another option
func (m *AddTransactionModel) renderSuccess() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")).
		Bold(true).
		Margin(1, 0)
	
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Margin(1, 0)
	
	buttonStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 2).
		Margin(1, 1, 0, 0)
	
	addAnotherButton := buttonStyle.Copy().
		BorderForeground(lipgloss.Color("46")).
		Foreground(lipgloss.Color("46")).
		Render("Add Another (A)")
	
	dashboardButton := buttonStyle.Copy().
		BorderForeground(lipgloss.Color("39")).
		Foreground(lipgloss.Color("39")).
		Render("Dashboard (ESC)")
	
	buttons := lipgloss.JoinHorizontal(lipgloss.Left, addAnotherButton, dashboardButton)
	
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("✓ Transaction Saved Successfully!"),
		messageStyle.Render("Your transaction has been added to your expense tracker."),
		buttons,
	)
}

// loadCategories loads categories for the dropdown
func (m *AddTransactionModel) loadCategories() tea.Cmd {
	return func() tea.Msg {
		categories, err := m.categoryUseCase.GetAllCategories()
		if err != nil {
			return dataLoadErrorMsg{error: err.Error()}
		}
		
		return categoriesLoadedMsg{categories: categories}
	}
}

// updateCategoryDropdown updates the category dropdown based on transaction type
func (m *AddTransactionModel) updateCategoryDropdown() {
	selectedType := m.typeToggle.GetSelectedValue().(domain.TransactionType)
	
	var options []components.DropdownOption
	for _, category := range m.categories {
		if category.Type == selectedType {
			options = append(options, components.DropdownOption{
				Label: category.Name,
				Value: category.ID,
				ID:    category.ID,
			})
		}
	}
	
	m.categoryDropdown.SetOptions(options)
	
	// Reset selection when type changes and set placeholder appropriately
	if len(options) > 0 {
		m.categoryDropdown.SetSelected(0)
		if selectedType == domain.TypeIncome {
			m.categoryDropdown.SetPlaceholder("Select income category")
		} else {
			m.categoryDropdown.SetPlaceholder("Select expense category")
		}
	} else {
		// No categories available for this type
		m.categoryDropdown.SetPlaceholder("No categories available")
	}
}


// handleSubmit is called when the form is submitted
func (m *AddTransactionModel) handleSubmit() tea.Cmd {
	if !m.form.IsValid() {
		return nil
	}
	
	m.saving = true
	
	return func() tea.Msg {
		// Parse form data
		transactionType := m.typeToggle.GetSelectedValue().(domain.TransactionType)
		
		amountStr := m.amountInput.GetValue()
		amountStr = lipgloss.NewStyle().Render(amountStr) // Clean any formatting
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			return saveErrorMsg{error: "Invalid amount format"}
		}
		
		categoryID := m.categoryDropdown.GetSelectedValue().(int)
		description := m.descriptionInput.GetValue()
		
		dateStr := m.dateInput.GetValue()
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return saveErrorMsg{error: "Invalid date format"}
		}
		
		// Create transaction
		_, err = m.transactionUseCase.AddTransaction(transactionType, amount, categoryID, description, date)
		if err != nil {
			return saveErrorMsg{error: err.Error()}
		}
		
		return transactionSavedMsg{}
	}
}

// handleCancel is called when the form is cancelled
func (m *AddTransactionModel) handleCancel() tea.Cmd {
	return func() tea.Msg {
		return NavigateMsg{View: DashboardView}
	}
}

// handleFieldChange is called when the focused field changes
func (m *AddTransactionModel) handleFieldChange(fieldIndex int) tea.Cmd {
	// Field 0 is the type toggle - update categories when we leave it
	if fieldIndex != 0 {
		return func() tea.Msg {
			// This will trigger category update in the main Update method
			return tea.KeyMsg{}
		}
	}
	return nil
}

// validateDate validates date input
func (m *AddTransactionModel) validateDate(value string) error {
	if value == "" {
		return nil // Allow empty for required validation to handle
	}
	
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}
	
	if date.After(time.Now()) {
		return domain.ErrFutureDate
	}
	
	return nil
}

// Custom messages for add transaction
type categoriesLoadedMsg struct {
	categories []*domain.Category
}

type transactionSavedMsg struct{}

type saveErrorMsg struct {
	error string
}