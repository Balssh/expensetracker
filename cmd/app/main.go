package main

import (
	"fmt"
	"log"
	"os"

	"expense-tracker/internal/core/usecase"
	"expense-tracker/internal/handler/tui"
	"expense-tracker/internal/repository/sqlite"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	// Ensure the database directory exists
	dbDir := homeDir + "/.local/share/expensetracker"
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	db, err := sqlite.NewDatabase(dbDir + "/expense_tracker.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	transactionRepo := sqlite.NewTransactionRepository(db)
	categoryRepo := sqlite.NewCategoryRepository(db)

	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, categoryRepo)
	summaryUseCase := usecase.NewSummaryUseCase(transactionRepo)

	model := tui.NewModel(transactionUseCase, summaryUseCase)

	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
