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
	"github.com/charmbracelet/lipgloss"
)

type transactionCategoriesMsg struct {
	categories []*domain.Category
	err        error
}

type transactionSubmissionMsg struct {
	success bool
	err     error
}

type TransactionType string

const (
	TransactionTypeExpense TransactionType = "expense"
	TransactionTypeIncome  TransactionType = "income"
)

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

type AddTransactionModel struct {
	transactionUseCase *usecase.TransactionUseCase
	transactionType    TransactionType
	categories         []*domain.Category
	inputs             []textinput.Model
	currentField       formField
	selectedCategory   int
	currentMode        editMode
	loading            bool
	err                error
	successMsg         string
	shouldReturn       bool
	width              int
	height             int
}

func NewAddTransactionModel(transactionUseCase *usecase.TransactionUseCase, transactionType TransactionType) *AddTransactionModel {
	m := &AddTransactionModel{
		transactionUseCase: transactionUseCase,
		transactionType:    transactionType,
		inputs:             make([]textinput.Model, 3),
		currentField:       fieldDescription,
		currentMode:        modeNavigate,
	}

	// Description input
	m.inputs[0] = textinput.New()
	m.inputs[0].Placeholder = "Enter " + string(transactionType) + " description"
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

func (m *AddTransactionModel) Init() tea.Cmd {
	return m.fetchCategories()
}

func (m *AddTransactionModel) Reset() {
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

func (m *AddTransactionModel) SetDimensions(width, height int) {
	m.width = width
	m.height = height
}

func (m *AddTransactionModel) fetchCategories() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		categories, err := m.transactionUseCase.GetCategories(ctx, string(m.transactionType))
		return transactionCategoriesMsg{categories: categories, err: err}
	})
}

func (m *AddTransactionModel) submitTransaction() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()

		amount, err := strconv.ParseFloat(m.inputs[1].Value(), 64)
		if err != nil || amount <= 0 {
			return transactionSubmissionMsg{err: fmt.Errorf("invalid amount")}
		}

		dateStr := strings.TrimSpace(m.inputs[2].Value())
		var date time.Time
		if dateStr == "" {
			date = time.Now()
		} else {
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return transactionSubmissionMsg{err: fmt.Errorf("invalid date format (use YYYY-MM-DD)")}
			}
			date = parsedDate
		}

		if len(m.categories) == 0 {
			return transactionSubmissionMsg{err: fmt.Errorf("no categories available")}
		}

		transaction := &domain.Transaction{
			Description: m.inputs[0].Value(),
			Amount:      amount,
			Date:        date,
			Type:        string(m.transactionType),
			Category:    m.categories[m.selectedCategory],
		}

		err = m.transactionUseCase.AddTransaction(ctx, transaction)
		if err != nil {
			return transactionSubmissionMsg{err: err}
		}

		return transactionSubmissionMsg{success: true}
	})
}

func (m *AddTransactionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case transactionCategoriesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.categories = msg.categories
		}
		return m, nil

	case transactionSubmissionMsg:
		if msg.err != nil {
			m.err = msg.err
		} else if msg.success {
			m.successMsg = strings.Title(string(m.transactionType)) + " added successfully!"
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

// Navigation and input handling methods (similar to original but unified)
func (m *AddTransactionModel) handleNavigateMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		return m.attemptSubmit()
	}
	return m, nil
}

func (m *AddTransactionModel) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.currentMode = modeNavigate
		m.inputs[m.getInputIndex()].Blur()
		return m, nil
	case "enter":
		m.currentMode = modeNavigate
		m.inputs[m.getInputIndex()].Blur()
		m.navigateDown()
		return m, nil
	case "ctrl+u":
		m.inputs[m.getInputIndex()].SetValue("")
		return m, nil
	}

	if idx := m.getInputIndex(); idx >= 0 {
		var cmd tea.Cmd
		m.inputs[idx], cmd = m.inputs[idx].Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *AddTransactionModel) handleCategorySelectMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
	}
	return m, nil
}

// Helper methods
func (m *AddTransactionModel) navigateUp() {
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

func (m *AddTransactionModel) navigateDown() {
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

func (m *AddTransactionModel) enterCurrentField() (tea.Model, tea.Cmd) {
	switch m.currentField {
	case fieldDescription, fieldAmount, fieldDate:
		m.currentMode = modeEdit
		idx := m.getInputIndex()
		if idx >= 0 {
			m.inputs[idx].Focus()
		}
		return m, nil
	case fieldCategory:
		if len(m.categories) > 0 {
			m.currentMode = modeCategorySelect
		}
		return m, nil
	case fieldSubmit:
		return m.attemptSubmit()
	}
	return m, nil
}

func (m *AddTransactionModel) getInputIndex() int {
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

func (m *AddTransactionModel) attemptSubmit() (tea.Model, tea.Cmd) {
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
	return m, m.submitTransaction()
}

func (m *AddTransactionModel) View() string {
	config := NewCenterConfig(m.width, m.height)
	
	if m.loading {
		content := loadingStyle.Render("Saving " + string(m.transactionType) + "...")
		return lipgloss.Place(config.Width, config.Height, 
			lipgloss.Center, lipgloss.Center, content)
	}

	// Create form content
	formContent := m.createFormContent()
	helpContent := m.createFormHelpText()
	
	// Dynamic title based on transaction type
	var icon string
	if m.transactionType == TransactionTypeExpense {
		icon = "💳"
	} else {
		icon = "💰"
	}
	title := titleStyle.Render(icon + " Add " + strings.Title(string(m.transactionType)))
	
	// Create base layout
	baseContent := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		formContent,
		"",
		helpContent,
	)
	
	// If category selection is active, overlay the popup
	if m.currentMode == modeCategorySelect {
		popup := m.createCategoryPopup()
		// Overlay popup on top of the base content
		return m.overlayPopup(baseContent, popup, config)
	}
	
	return lipgloss.Place(config.Width, config.Height, 
		lipgloss.Center, lipgloss.Center, baseContent)
}

func (m *AddTransactionModel) createFormContent() string {
	var b strings.Builder
	
	// Show status messages
	if m.err != nil {
		b.WriteString(errorStyle.Render("❌ " + m.err.Error()) + "\n\n")
	}
	if m.successMsg != "" {
		b.WriteString(successStyle.Render("✅ " + m.successMsg) + "\n\n")
	}
	
	// Create form fields
	fields := []struct {
		label    string
		field    formField
		content  string
		required bool
	}{
		{"Description", fieldDescription, m.renderFormField(fieldDescription), true},
		{"Amount", fieldAmount, m.renderFormField(fieldAmount), true},
		{"Date", fieldDate, m.renderFormField(fieldDate), false},
		{"Category", fieldCategory, m.renderFormField(fieldCategory), true},
	}
	
	for _, field := range fields {
		// Navigation indicator
		indicator := "  "
		if m.currentField == field.field && m.currentMode == modeNavigate {
			indicator = "▶ "
		}
		
		// Required indicator
		required := ""
		if field.required {
			required = " " + errorStyle.Render("*")
		}
		
		label := formFieldLabelStyle.Render(field.label + required + ":")
		fieldLine := fmt.Sprintf("%s%s\n   %s", indicator, label, field.content)
		b.WriteString(fieldLine + "\n\n")
	}
	
	// Submit button - no border as requested
	submitIndicator := "  "
	if m.currentField == fieldSubmit && m.currentMode == modeNavigate {
		submitIndicator = "▶ "
	}
	
	buttonText := "💾 Save " + strings.Title(string(m.transactionType))
	submitButton := m.renderSubmitButton(buttonText)
	b.WriteString(fmt.Sprintf("%s%s", submitIndicator, submitButton))
	
	// Create bordered form panel
	panelWidth := NewCenterConfig(m.width, m.height).CalculateContentWidth() - 8
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Width(panelWidth).
		Padding(1, 2)
	
	return style.Render(b.String())
}

func (m *AddTransactionModel) renderFormField(field formField) string {
	switch field {
	case fieldDescription, fieldAmount, fieldDate:
		return m.renderTextInput(int(field))
	case fieldCategory:
		return m.renderCategoryField()
	default:
		return ""
	}
}

func (m *AddTransactionModel) renderTextInput(index int) string {
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

func (m *AddTransactionModel) renderCategoryField() string {
	if len(m.categories) == 0 {
		return inputStyle.Render(inputPlaceholderStyle.Render("Loading categories..."))
	}
	
	selectedCategory := m.categories[m.selectedCategory].Name
	
	if m.currentMode == modeCategorySelect {
		return inputFocusedStyle.Render(selectedCategory + " ▼")
	}
	
	return inputStyle.Render(selectedCategory)
}

func (m *AddTransactionModel) renderSubmitButton(text string) string {
	// Remove border styling from button as requested
	if m.currentField == fieldSubmit && m.currentMode == modeNavigate {
		return successStyle.Render(text) // Highlighted but no border
	}
	return helpStyle.Render(text)
}


func (m *AddTransactionModel) createFormHelpText() string {
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
	
	return strings.Join(helpTexts, " • ")
}

// createCategoryPopup creates a centered popup for category selection
func (m *AddTransactionModel) createCategoryPopup() string {
	var b strings.Builder
	
	// Header
	b.WriteString("Select Category\n\n")
	
	// Category list (limit height to stay within borders)
	maxCategories := 8 // Limit to keep popup reasonable size
	startIdx := 0
	endIdx := len(m.categories)
	
	if len(m.categories) > maxCategories {
		// Show categories around the selected one
		startIdx = m.selectedCategory - maxCategories/2
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + maxCategories
		if endIdx > len(m.categories) {
			endIdx = len(m.categories)
			startIdx = endIdx - maxCategories
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}
	
	for i := startIdx; i < endIdx; i++ {
		category := m.categories[i]
		if i == m.selectedCategory {
			b.WriteString(dropdownItemSelectedStyle.Render("▶ " + category.Name))
		} else {
			b.WriteString("  " + category.Name)
		}
		b.WriteString("\n")
	}
	
	// Show scroll indicators if needed
	if len(m.categories) > maxCategories {
		if startIdx > 0 {
			b.WriteString(helpStyle.Render("  ↑ More above") + "\n")
		}
		if endIdx < len(m.categories) {
			b.WriteString(helpStyle.Render("  ↓ More below") + "\n")
		}
	}
	
	// Create popup with border
	popupStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Background(lipgloss.Color("#1a1a1a")).
		Padding(1, 2).
		Width(30) // Fixed reasonable width
	
	return popupStyle.Render(b.String())
}

// overlayPopup overlays a popup on top of base content using proper layering
func (m *AddTransactionModel) overlayPopup(baseContent, popup string, config CenterConfig) string {
	// Simple approach: use lipgloss PlaceHorizontal and PlaceVertical for proper layering
	
	// First, place the base content normally
	base := lipgloss.Place(config.Width, config.Height, 
		lipgloss.Center, lipgloss.Center, baseContent)
	
	// Create a semi-transparent overlay effect by dimming colors in the base
	dimmedBase := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render(base)
	
	
	// Position the popup exactly in center using proper coordinates
	popupLines := strings.Split(popup, "\n")
	baseLines := strings.Split(dimmedBase, "\n")
	
	// Ensure we have enough lines
	for len(baseLines) < config.Height {
		baseLines = append(baseLines, "")
	}
	
	// Calculate popup position
	popupHeight := len(popupLines)
	startY := (config.Height - popupHeight) / 2
	if startY < 0 { startY = 0 }
	
	// Overlay popup lines onto base, maintaining proper spacing
	result := make([]string, len(baseLines))
	copy(result, baseLines)
	
	popupWidth := 32 // Fixed width from createCategoryPopup
	startX := (config.Width - popupWidth) / 2
	if startX < 0 { startX = 0 }
	
	for i, popupLine := range popupLines {
		if startY+i < len(result) && startY+i >= 0 {
			// Replace the middle section of the base line with the popup line
			baseLine := result[startY+i]
			
			// Pad or trim base line to config width
			if len(baseLine) < config.Width {
				baseLine += strings.Repeat(" ", config.Width-len(baseLine))
			}
			if len(baseLine) > config.Width {
				baseLine = baseLine[:config.Width]
			}
			
			// Insert popup line at center position
			if len(popupLine) > 0 && startX+len(popupLine) <= len(baseLine) {
				result[startY+i] = baseLine[:startX] + popupLine + baseLine[startX+len(popupLine):]
			}
		}
	}
	
	return strings.Join(result, "\n")
}