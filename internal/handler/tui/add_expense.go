package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"expense-tracker/internal/core/domain"
	"expense-tracker/internal/core/usecase"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type expenseCategoriesMsg struct {
	categories []*domain.Category
	err        error
}

type expenseSubmissionMsg struct {
	success bool
	err     error
}

type formField int

const (
	fieldDescription formField = iota
	fieldAmount
	fieldDate
	fieldCategory
	fieldSubmit
)

type editMode int

const (
	modeNavigate editMode = iota
	modeEdit
	modeCategorySelect
)

type AddExpenseModel struct {
	expenseUseCase   *usecase.TransactionUseCase
	categories       []*domain.Category
	inputs           []textinput.Model
	currentField     formField
	selectedCategory int
	currentMode      editMode
	loading          bool
	err              error
	successMsg       string
	shouldReturn     bool
	width            int
	height           int
}

func NewAddExpenseModel(transactionUseCase *usecase.TransactionUseCase) *AddExpenseModel {
	m := &AddExpenseModel{
		expenseUseCase: transactionUseCase,
		inputs:         make([]textinput.Model, 3),
		currentField:   fieldDescription,
		currentMode:    modeNavigate,
	}

	// Description input
	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "Enter expense description"
	m.inputs[0].CharLimit = 200

	// Amount input
	m.inputs[1] = textinput.New()
	m.inputs[1].Placeholder = "0.00"
	m.inputs[1].CharLimit = 20

	// Date input
	m.inputs[2] = textinput.New()
	m.inputs[2].Placeholder = "YYYY-MM-DD (leave empty for today)"
	m.inputs[2].CharLimit = 10

	return m
}

func (m *AddExpenseModel) Init() tea.Cmd {
	return m.fetchCategories()
}

func (m *AddExpenseModel) Reset() {
	// Clear all input values
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}
	
	// Reset state
	m.currentField = fieldDescription
	m.selectedCategory = 0
	m.currentMode = modeNavigate
	m.loading = false
	m.err = nil
	m.successMsg = ""
	m.shouldReturn = false
}

// SetDimensions updates the model's width and height for responsive layout
func (m *AddExpenseModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

func (m *AddExpenseModel) fetchCategories() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		categories, err := m.expenseUseCase.GetCategories(ctx, "expense")
		return expenseCategoriesMsg{categories: categories, err: err}
	})
}

func (m *AddExpenseModel) submitExpense() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()

		amount, err := strconv.ParseFloat(m.inputs[1].Value(), 64)
		if err != nil || amount <= 0 {
			return expenseSubmissionMsg{err: fmt.Errorf("invalid amount")}
		}

		dateStr := strings.TrimSpace(m.inputs[2].Value())
		var date time.Time
		if dateStr == "" {
			date = time.Now()
		} else {
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return expenseSubmissionMsg{err: fmt.Errorf("invalid date format (use YYYY-MM-DD)")}
			}
			date = parsedDate
		}

		if len(m.categories) == 0 {
			return expenseSubmissionMsg{err: fmt.Errorf("no categories available")}
		}

		transaction := &domain.Transaction{
			Description: m.inputs[0].Value(),
			Amount:      amount,
			Date:        date,
			Type:        "expense",
			Category:    m.categories[m.selectedCategory],
		}

		err = m.expenseUseCase.AddTransaction(ctx, transaction)
		if err != nil {
			return expenseSubmissionMsg{err: err}
		}

		return expenseSubmissionMsg{success: true}
	})
}

func (m *AddExpenseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case expenseCategoriesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.categories = msg.categories
		}
		return m, nil

	case expenseSubmissionMsg:
		if msg.err != nil {
			m.err = msg.err
		} else if msg.success {
			m.successMsg = "Expense added successfully!"
			m.shouldReturn = true
		}
		return m, nil

	case tea.KeyMsg:
		// Handle different modes
		switch m.currentMode {
		case modeNavigate:
			return m.handleNavigateMode(msg)
		case modeEdit:
			return m.handleEditMode(msg)
		case modeCategorySelect:
			return m.handleCategorySelectMode(msg)
		}
	}

	return m, nil
}

// handleNavigateMode handles navigation between form fields
func (m *AddExpenseModel) handleNavigateMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.shouldReturn = true
		return m, nil

	case "up", "k":
		m.navigateUp()
		return m, nil

	case "down", "j":
		m.navigateDown()
		return m, nil

	case "enter":
		return m.enterCurrentField()

	case "ctrl+s":
		// Quick save shortcut
		return m.attemptSubmit()
	}

	return m, nil
}

// handleEditMode handles text input editing
func (m *AddExpenseModel) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Exit edit mode
		m.currentMode = modeNavigate
		m.inputs[m.getInputIndex()].Blur()
		return m, nil

	case "enter":
		// Confirm edit and move to next field
		m.currentMode = modeNavigate
		m.inputs[m.getInputIndex()].Blur()
		m.navigateDown()
		return m, nil

	case "ctrl+u":
		// Clear current field
		m.inputs[m.getInputIndex()].SetValue("")
		return m, nil
	}

	// Update the current input
	if idx := m.getInputIndex(); idx >= 0 {
		var cmd tea.Cmd
		m.inputs[idx], cmd = m.inputs[idx].Update(msg)
		return m, cmd
	}

	return m, nil
}

// handleCategorySelectMode handles category selection dropdown
func (m *AddExpenseModel) handleCategorySelectMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.currentMode = modeNavigate
		return m, nil

	case "up", "k":
		if m.selectedCategory > 0 {
			m.selectedCategory--
		}
		return m, nil

	case "down", "j":
		if m.selectedCategory < len(m.categories)-1 {
			m.selectedCategory++
		}
		return m, nil

	case "enter":
		m.currentMode = modeNavigate
		return m, nil

	case "/":
		// TODO: Implement category search in future iteration
		return m, nil
	}

	return m, nil
}

// Navigation helper methods
func (m *AddExpenseModel) navigateUp() {
	switch m.currentField {
	case fieldDescription:
		// Stay at description
	case fieldAmount:
		m.currentField = fieldDescription
	case fieldDate:
		m.currentField = fieldAmount
	case fieldCategory:
		m.currentField = fieldDate
	case fieldSubmit:
		m.currentField = fieldCategory
	}
}

func (m *AddExpenseModel) navigateDown() {
	switch m.currentField {
	case fieldDescription:
		m.currentField = fieldAmount
	case fieldAmount:
		m.currentField = fieldDate
	case fieldDate:
		m.currentField = fieldCategory
	case fieldCategory:
		m.currentField = fieldSubmit
	case fieldSubmit:
		// Stay at submit
	}
}

func (m *AddExpenseModel) enterCurrentField() (tea.Model, tea.Cmd) {
	switch m.currentField {
	case fieldDescription, fieldAmount, fieldDate:
		// Enter edit mode for text inputs
		m.currentMode = modeEdit
		idx := m.getInputIndex()
		if idx >= 0 {
			m.inputs[idx].Focus()
		}
		return m, nil

	case fieldCategory:
		// Open category selection
		if len(m.categories) > 0 {
			m.currentMode = modeCategorySelect
		}
		return m, nil

	case fieldSubmit:
		// Submit the form
		return m.attemptSubmit()
	}

	return m, nil
}

func (m *AddExpenseModel) getInputIndex() int {
	switch m.currentField {
	case fieldDescription:
		return 0
	case fieldAmount:
		return 1
	case fieldDate:
		return 2
	default:
		return -1
	}
}

func (m *AddExpenseModel) attemptSubmit() (tea.Model, tea.Cmd) {
	// Validate required fields
	if m.inputs[0].Value() == "" {
		m.err = fmt.Errorf("description is required")
		return m, nil
	}

	if m.inputs[1].Value() == "" {
		m.err = fmt.Errorf("amount is required")
		return m, nil
	}

	if len(m.categories) == 0 {
		m.err = fmt.Errorf("no categories available")
		return m, nil
	}

	m.loading = true
	m.err = nil
	return m, m.submitExpense()
}

func (m *AddExpenseModel) View() string {
	if m.loading {
		centered := CenterHorizontally(loadingStyle.Render("Loading expense form..."), NewCenterConfig(m.width, m.height))
		return CreateMainApplicationBorder(centered, NewCenterConfig(m.width, m.height))
	}

	// Create form content
	formContent := m.createFormContent()
	
	// Create help panel
	helpContent := m.createFormHelpText()
	
	// Main title
	title := titleStyle.Render("💳 Add Expense")
	
	// Combine all content
	fullContent := title + "\n\n" + formContent + "\n\n" + helpContent
	
	// Center and apply main border
	config := NewCenterConfig(m.width, m.height)
	centered := CenterHorizontally(fullContent, config)
	
	return CreateMainApplicationBorder(centered, config)
}

func (m *AddExpenseModel) createFormContent() string {
	var b strings.Builder
	
	// Show status messages
	if m.err != nil {
		b.WriteString(errorStyle.Render("❌ " + m.err.Error()) + "\n\n")
	}
	
	if m.successMsg != "" {
		b.WriteString(successStyle.Render("✅ " + m.successMsg) + "\n\n")
	}
	
	// Create form fields with navigation indicators
	fields := []struct {
		label       string
		field       formField
		content     string
		required    bool
	}{
		{"Description", fieldDescription, m.renderFormField(fieldDescription), true},
		{"Amount", fieldAmount, m.renderFormField(fieldAmount), true},
		{"Date", fieldDate, m.renderFormField(fieldDate), false},
		{"Category", fieldCategory, m.renderFormField(fieldCategory), true},
	}
	
	for _, field := range fields {
		// Add navigation indicator
		indicator := "  "
		if m.currentField == field.field && m.currentMode == modeNavigate {
			indicator = "▶ "
		}
		
		// Add required indicator
		required := ""
		if field.required {
			required = " " + errorStyle.Render("*")
		}
		
		// Format field
		label := formFieldLabelStyle.Render(field.label + required + ":")
		fieldLine := fmt.Sprintf("%s%s\n   %s", indicator, label, field.content)
		b.WriteString(fieldLine + "\n\n")
	}
	
	// Submit button
	submitIndicator := "  "
	if m.currentField == fieldSubmit && m.currentMode == modeNavigate {
		submitIndicator = "▶ "
	}
	
	submitButton := m.renderSubmitButton()
	b.WriteString(fmt.Sprintf("%s%s", submitIndicator, submitButton))
	
	// Category dropdown if in selection mode
	if m.currentMode == modeCategorySelect {
		b.WriteString("\n\n" + m.renderCategoryDropdown())
	}
	
	return formStyle.Render(b.String())
}

func (m *AddExpenseModel) renderFormField(field formField) string {
	switch field {
	case fieldDescription:
		return m.renderTextInput(0)
	case fieldAmount:
		return m.renderTextInput(1)
	case fieldDate:
		return m.renderTextInput(2)
	case fieldCategory:
		return m.renderCategoryField()
	default:
		return ""
	}
}

func (m *AddExpenseModel) renderTextInput(index int) string {
	input := m.inputs[index]
	
	if m.currentField == formField(index) && m.currentMode == modeEdit {
		return inputFocusedStyle.Render(input.View())
	}
	
	value := input.Value()
	if value == "" {
		value = inputPlaceholderStyle.Render(input.Placeholder)
	}
	
	return inputStyle.Render(value)
}

func (m *AddExpenseModel) renderCategoryField() string {
	if len(m.categories) == 0 {
		return inputStyle.Render(inputPlaceholderStyle.Render("Loading categories..."))
	}
	
	selectedCategory := m.categories[m.selectedCategory].Name
	
	if m.currentMode == modeCategorySelect {
		return inputFocusedStyle.Render(selectedCategory + " ▼")
	}
	
	return inputStyle.Render(selectedCategory)
}

func (m *AddExpenseModel) renderSubmitButton() string {
	buttonText := "💾 Save Expense"
	
	if m.currentField == fieldSubmit && m.currentMode == modeNavigate {
		return modalButtonFocusedStyle.Render(buttonText)
	}
	
	return modalButtonStyle.Render(buttonText)
}

func (m *AddExpenseModel) renderCategoryDropdown() string {
	var b strings.Builder
	
	b.WriteString(dropdownStyle.Render("Select Category (↑/↓ to navigate, Enter to select):") + "\n")
	
	for i, category := range m.categories {
		if i == m.selectedCategory {
			b.WriteString(dropdownItemSelectedStyle.Render("→ " + category.Name))
		} else {
			b.WriteString(dropdownItemStyle.Render("  " + category.Name))
		}
		b.WriteString("\n")
	}
	
	return b.String()
}

func (m *AddExpenseModel) createFormHelpText() string {
	var helpTexts []string
	
	switch m.currentMode {
	case modeNavigate:
		helpTexts = []string{
			helpKeyStyle.Render("↑/↓") + " Navigate",
			helpKeyStyle.Render("Enter") + " Edit/Select",
			helpKeyStyle.Render("Ctrl+S") + " Save",
			helpKeyStyle.Render("Esc") + " Cancel",
		}
	case modeEdit:
		helpTexts = []string{
			helpKeyStyle.Render("Type") + " to edit",
			helpKeyStyle.Render("Enter") + " Confirm",
			helpKeyStyle.Render("Ctrl+U") + " Clear",
			helpKeyStyle.Render("Esc") + " Cancel",
		}
	case modeCategorySelect:
		helpTexts = []string{
			helpKeyStyle.Render("↑/↓") + " Navigate",
			helpKeyStyle.Render("Enter") + " Select",
			helpKeyStyle.Render("Esc") + " Cancel",
		}
	}
	
	return navigationStyle.Render(strings.Join(helpTexts, " • "))
}
