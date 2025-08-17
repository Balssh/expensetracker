package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"expense-tracker/internal/core/domain"
	"expense-tracker/internal/core/usecase"
)

type summaryMsg struct {
	summary      *domain.Summary
	transactions []*domain.Transaction
	err          error
}

type DashboardModel struct {
	summaryUseCase *usecase.SummaryUseCase
	summary        *domain.Summary
	transactions   []*domain.Transaction
	loading        bool
	err            error
}

func NewDashboardModel(summaryUseCase *usecase.SummaryUseCase) *DashboardModel {
	return &DashboardModel{
		summaryUseCase: summaryUseCase,
		loading:        true,
	}
}

func (m *DashboardModel) Init() tea.Cmd {
	return m.fetchData()
}

func (m *DashboardModel) Refresh() tea.Cmd {
	m.loading = true
	return m.fetchData()
}

func (m *DashboardModel) fetchData() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		now := time.Now()
		
		summary, err := m.summaryUseCase.GetMonthlySummary(ctx, now.Year(), now.Month())
		if err != nil {
			return summaryMsg{err: err}
		}

		transactions, err := m.summaryUseCase.GetRecentTransactions(ctx, 10)
		if err != nil {
			return summaryMsg{err: err}
		}

		return summaryMsg{
			summary:      summary,
			transactions: transactions,
		}
	})
}

func (m *DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case summaryMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.summary = msg.summary
			m.transactions = msg.transactions
			m.err = nil
		}
		return m, nil
	}
	return m, nil
}

func (m *DashboardModel) View() string {
	if m.loading {
		return titleStyle.Render("Expense Tracker") + "\n\nLoading..."
	}

	if m.err != nil {
		return titleStyle.Render("Expense Tracker") + "\n\n" + 
			errorStyle.Render("Error: "+m.err.Error())
	}

	var b strings.Builder

	b.WriteString(titleStyle.Render("Expense Tracker"))
	b.WriteString("\n\n")

	now := time.Now()
	monthYear := fmt.Sprintf("%s %d", now.Month().String(), now.Year())
	b.WriteString(fmt.Sprintf("Summary for %s\n", monthYear))

	summaryContent := fmt.Sprintf(
		"%s  $%.2f\n%s $%.2f\n%s   %s",
		incomeStyle.Render("Income: "),
		m.summary.TotalIncome,
		expenseStyle.Render("Expenses:"),
		m.summary.TotalExpense,
		"Balance:",
		m.formatBalance(m.summary.NetBalance),
	)

	b.WriteString(summaryBoxStyle.Render(summaryContent))
	b.WriteString("\n")

	b.WriteString("Recent Transactions\n")
	b.WriteString(strings.Repeat("─", 60) + "\n")

	if len(m.transactions) == 0 {
		b.WriteString("No transactions found.\n")
	} else {
		for _, transaction := range m.transactions {
			line := fmt.Sprintf("%-20s %-12s %8s $%.2f",
				transaction.Date.Format("Jan 02"),
				transaction.Category,
				m.formatTransactionType(transaction.Type),
				transaction.Amount,
			)
			b.WriteString(line + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Controls: (a) Add Expense • (i) Add Income • (l) List All • (q) Quit"))

	return b.String()
}

func (m *DashboardModel) formatBalance(balance float64) string {
	if balance >= 0 {
		return balancePositiveStyle.Render(fmt.Sprintf("$%.2f", balance))
	}
	return balanceNegativeStyle.Render(fmt.Sprintf("$%.2f", balance))
}

func (m *DashboardModel) formatTransactionType(transactionType string) string {
	if transactionType == "income" {
		return incomeStyle.Render("Income")
	}
	return expenseStyle.Render("Expense")
}