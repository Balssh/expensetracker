package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"expense-tracker/internal/core/usecase"
	"expense-tracker/internal/handler/tui"
	"expense-tracker/internal/repository/sqlite"
)

func main() {
	dbPath := "expense_tracker.db"
	
	db, err := sqlite.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	expenseRepo := sqlite.NewExpenseRepository(db)
	incomeRepo := sqlite.NewIncomeRepository(db)
	categoryRepo := sqlite.NewCategoryRepository(db)
	transactionRepo := sqlite.NewTransactionRepository(db)

	expenseUseCase := usecase.NewExpenseUseCase(expenseRepo, categoryRepo)
	incomeUseCase := usecase.NewIncomeUseCase(incomeRepo, categoryRepo)
	summaryUseCase := usecase.NewSummaryUseCase(expenseRepo, incomeRepo, transactionRepo)

	model := tui.NewModel(expenseUseCase, incomeUseCase, summaryUseCase)

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}