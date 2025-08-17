package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"expense-tracker/internal/core/domain"
	"expense-tracker/internal/core/usecase"
)

type expenseCategoriesMsg struct {
	categories []*domain.Category
	err        error
}

type expenseSubmissionMsg struct {
	success bool
	err     error
}

type AddExpenseModel struct {
	expenseUseCase   *usecase.ExpenseUseCase
	categories       []*domain.Category
	inputs           []textinput.Model
	focused          int
	selectedCategory int
	showCategories   bool
	loading          bool
	err              error
	successMsg       string
	shouldReturn     bool
}

const (
	expenseDescriptionInput = iota
	expenseAmountInput
	expenseDateInput
)

func NewAddExpenseModel(expenseUseCase *usecase.ExpenseUseCase) *AddExpenseModel {
	m := &AddExpenseModel{
		expenseUseCase: expenseUseCase,
		inputs:         make([]textinput.Model, 3),
	}

	m.inputs[expenseDescriptionInput] = textinput.New()
	m.inputs[expenseDescriptionInput].Placeholder = "Enter description"
	m.inputs[expenseDescriptionInput].Focus()

	m.inputs[expenseAmountInput] = textinput.New()
	m.inputs[expenseAmountInput].Placeholder = "0.00"

	m.inputs[expenseDateInput] = textinput.New()
	m.inputs[expenseDateInput].Placeholder = "YYYY-MM-DD (or leave empty for today)"

	return m
}

func (m *AddExpenseModel) Init() tea.Cmd {
	return m.fetchCategories()
}

func (m *AddExpenseModel) Reset() {
	m.inputs[expenseDescriptionInput].SetValue("")
	m.inputs[expenseAmountInput].SetValue("")
	m.inputs[expenseDateInput].SetValue("")
	m.focused = 0
	m.selectedCategory = 0
	m.showCategories = false
	m.loading = false
	m.err = nil
	m.successMsg = ""
	m.shouldReturn = false
	m.inputs[expenseDescriptionInput].Focus()
	m.inputs[expenseAmountInput].Blur()
	m.inputs[expenseDateInput].Blur()
}

func (m *AddExpenseModel) fetchCategories() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		categories, err := m.expenseUseCase.GetExpenseCategories(ctx)
		return expenseCategoriesMsg{categories: categories, err: err}
	})
}

func (m *AddExpenseModel) submitExpense() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()

		amount, err := strconv.ParseFloat(m.inputs[expenseAmountInput].Value(), 64)
		if err != nil || amount <= 0 {
			return expenseSubmissionMsg{err: fmt.Errorf("invalid amount")}
		}

		dateStr := strings.TrimSpace(m.inputs[expenseDateInput].Value())
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

		categoryID := m.categories[m.selectedCategory].ID

		err = m.expenseUseCase.AddExpense(ctx, m.inputs[expenseDescriptionInput].Value(), amount, date, categoryID)
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
		switch msg.String() {
		case "ctrl+c", "esc":
			m.shouldReturn = true
			return m, nil

		case "enter":
			if m.showCategories {
				m.showCategories = false
				return m, nil
			}

			if m.focused < len(m.inputs)-1 {
				m.inputs[m.focused].Blur()
				m.focused++
				m.inputs[m.focused].Focus()
				return m, nil
			} else {
				if m.inputs[expenseDescriptionInput].Value() == "" {
					m.err = fmt.Errorf("description is required")
					return m, nil
				}
				m.loading = true
				return m, m.submitExpense()
			}

		case "tab":
			if !m.showCategories {
				m.showCategories = true
				return m, nil
			}

		case "up":
			if m.showCategories && m.selectedCategory > 0 {
				m.selectedCategory--
			}
			return m, nil

		case "down":
			if m.showCategories && m.selectedCategory < len(m.categories)-1 {
				m.selectedCategory++
			}
			return m, nil
		}
	}

	if !m.showCategories {
		var cmd tea.Cmd
		m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *AddExpenseModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Add Expense"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString("Loading...")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
	}

	if m.successMsg != "" {
		b.WriteString(successStyle.Render(m.successMsg))
		b.WriteString("\n\n")
	}

	formContent := fmt.Sprintf(
		"Description: %s\n\nAmount: %s\n\nDate: %s\n\nCategory: %s",
		m.renderInput(expenseDescriptionInput),
		m.renderInput(expenseAmountInput),
		m.renderInput(expenseDateInput),
		m.renderCategorySelection(),
	)

	b.WriteString(formStyle.Render(formContent))
	b.WriteString("\n\n")

	if m.showCategories {
		b.WriteString("Select Category (↑/↓ to navigate, Enter to select):\n")
		for i, category := range m.categories {
			if i == m.selectedCategory {
				b.WriteString(selectedItemStyle.Render("→ " + category.Name))
			} else {
				b.WriteString("  " + category.Name)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("Tab: Select Category • Enter: Next/Submit • Esc: Cancel"))

	return b.String()
}

func (m *AddExpenseModel) renderInput(index int) string {
	if m.focused == index && !m.showCategories {
		return focusedInputStyle.Render(m.inputs[index].View())
	}
	return inputStyle.Render(m.inputs[index].View())
}

func (m *AddExpenseModel) renderCategorySelection() string {
	if len(m.categories) == 0 {
		return "Loading categories..."
	}
	
	selectedCategoryName := m.categories[m.selectedCategory].Name
	if m.showCategories {
		return focusedInputStyle.Render(selectedCategoryName)
	}
	return inputStyle.Render(selectedCategoryName)
}