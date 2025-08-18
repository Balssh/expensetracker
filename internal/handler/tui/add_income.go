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

type incomeCategoriesMsg struct {
	categories []*domain.Category
	err        error
}

type incomeSubmissionMsg struct {
	success bool
	err     error
}

type AddIncomeModel struct {
	transactionUseCase *usecase.TransactionUseCase
	categories         []*domain.Category
	inputs             []textinput.Model
	focused            int
	selectedCategory   int
	showCategories     bool
	loading            bool
	err                error
	successMsg         string
	shouldReturn       bool
}

const (
	incomeDescriptionInput = iota
	incomeAmountInput
	incomeDateInput
)

func NewAddIncomeModel(transactionUseCase *usecase.TransactionUseCase) *AddIncomeModel {
	m := &AddIncomeModel{
		transactionUseCase: transactionUseCase,
		inputs:             make([]textinput.Model, 3),
	}

	m.inputs[incomeDescriptionInput] = textinput.New()
	m.inputs[incomeDescriptionInput].Placeholder = "Enter description"
	m.inputs[incomeDescriptionInput].Focus()

	m.inputs[incomeAmountInput] = textinput.New()
	m.inputs[incomeAmountInput].Placeholder = "0.00"

	m.inputs[incomeDateInput] = textinput.New()
	m.inputs[incomeDateInput].Placeholder = "YYYY-MM-DD (or leave empty for today)"

	return m
}

func (m *AddIncomeModel) Init() tea.Cmd {
	return m.fetchCategories()
}

func (m *AddIncomeModel) Reset() {
	m.inputs[incomeDescriptionInput].SetValue("")
	m.inputs[incomeAmountInput].SetValue("")
	m.inputs[incomeDateInput].SetValue("")
	m.focused = 0
	m.selectedCategory = 0
	m.showCategories = false
	m.loading = false
	m.err = nil
	m.successMsg = ""
	m.shouldReturn = false
	m.inputs[incomeDescriptionInput].Focus()
	m.inputs[incomeAmountInput].Blur()
	m.inputs[incomeDateInput].Blur()
}

func (m *AddIncomeModel) fetchCategories() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		categories, err := m.transactionUseCase.GetCategories(ctx, "income")
		return incomeCategoriesMsg{categories: categories, err: err}
	})
}

func (m *AddIncomeModel) submitIncome() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()

		amount, err := strconv.ParseFloat(m.inputs[incomeAmountInput].Value(), 64)
		if err != nil || amount <= 0 {
			return incomeSubmissionMsg{err: fmt.Errorf("invalid amount")}
		}

		dateStr := strings.TrimSpace(m.inputs[incomeDateInput].Value())
		var date time.Time
		if dateStr == "" {
			date = time.Now()
		} else {
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return incomeSubmissionMsg{err: fmt.Errorf("invalid date format (use YYYY-MM-DD)")}
			}
			date = parsedDate
		}

		if len(m.categories) == 0 {
			return incomeSubmissionMsg{err: fmt.Errorf("no categories available")}
		}

		transaction := &domain.Transaction{
			Description: m.inputs[incomeDescriptionInput].Value(),
			Amount:      amount,
			Date:        date,
			Type:        "income",
			Category:    m.categories[m.selectedCategory],
		}

		err = m.transactionUseCase.AddTransaction(ctx, transaction)
		if err != nil {
			return incomeSubmissionMsg{err: err}
		}

		return incomeSubmissionMsg{success: true}
	})
}

func (m *AddIncomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case incomeCategoriesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.categories = msg.categories
		}
		return m, nil

	case incomeSubmissionMsg:
		if msg.err != nil {
			m.err = msg.err
		} else if msg.success {
			m.successMsg = "Income added successfully!"
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
				if m.inputs[incomeDescriptionInput].Value() == "" {
					m.err = fmt.Errorf("description is required")
					return m, nil
				}
				m.loading = true
				return m, m.submitIncome()
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

func (m *AddIncomeModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Add Income"))
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
		m.renderInput(incomeDescriptionInput),
		m.renderInput(incomeAmountInput),
		m.renderInput(incomeDateInput),
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

func (m *AddIncomeModel) renderInput(index int) string {
	if m.focused == index && !m.showCategories {
		return focusedInputStyle.Render(m.inputs[index].View())
	}
	return inputStyle.Render(m.inputs[index].View())
}

func (m *AddIncomeModel) renderCategorySelection() string {
	if len(m.categories) == 0 {
		return "Loading categories..."
	}
	
	selectedCategoryName := m.categories[m.selectedCategory].Name
	if m.showCategories {
		return focusedInputStyle.Render(selectedCategoryName)
	}
	return inputStyle.Render(selectedCategoryName)
}