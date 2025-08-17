package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"expense-tracker/internal/core/domain"
	"expense-tracker/internal/core/usecase"
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
	var b strings.Builder

	b.WriteString(titleStyle.Render("All Transactions"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString("Loading...")
		return b.String()
	}

	if m.err != nil {
		b.WriteString(errorStyle.Render("Error: " + m.err.Error()))
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("q: Back to Dashboard"))
		return b.String()
	}

	searchSection := "Search: "
	if m.isSearching {
		searchSection += focusedInputStyle.Render(m.searchInput.View())
	} else {
		searchSection += inputStyle.Render(m.searchInput.View())
	}
	b.WriteString(searchSection)
	b.WriteString("\n\n")

	if len(m.transactions) == 0 {
		if m.searchInput.Value() != "" {
			b.WriteString("No transactions found matching your search.\n")
		} else {
			b.WriteString("No transactions found.\n")
		}
	} else {
		b.WriteString(m.renderTransactionTable())
	}

	b.WriteString("\n")
	
	helpText := "Controls: (/) Search • (c) Clear Search"
	if !m.isSearching {
		if m.currentPage > 0 {
			helpText += " • (p) Previous Page"
		}
		if len(m.transactions) == m.itemsPerPage {
			helpText += " • (n) Next Page"
		}
	}
	helpText += " • (q) Back"
	
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

func (m *TransactionsModel) renderTransactionTable() string {
	var b strings.Builder

	header := fmt.Sprintf("%-12s %-20s %-15s %-10s %10s",
		"Date", "Description", "Category", "Type", "Amount")
	b.WriteString(tableHeaderStyle.Render(header))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 70) + "\n")

	for _, transaction := range m.transactions {
		line := fmt.Sprintf("%-12s %-20s %-15s %-10s %10s",
			transaction.Date.Format("Jan 02, 2006"),
			m.truncateString(transaction.Description, 20),
			m.truncateString(transaction.Category, 15),
			m.formatTransactionType(transaction.Type),
			fmt.Sprintf("$%.2f", transaction.Amount),
		)
		b.WriteString(line + "\n")
	}

	if m.currentPage > 0 || len(m.transactions) == m.itemsPerPage {
		b.WriteString(strings.Repeat("─", 70) + "\n")
		pageInfo := fmt.Sprintf("Page %d", m.currentPage+1)
		if len(m.transactions) == m.itemsPerPage {
			pageInfo += " (more available)"
		}
		b.WriteString(pageInfo + "\n")
	}

	return b.String()
}

func (m *TransactionsModel) truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func (m *TransactionsModel) formatTransactionType(transactionType string) string {
	if transactionType == "income" {
		return incomeStyle.Render("Income")
	}
	return expenseStyle.Render("Expense")
}