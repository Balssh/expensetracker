package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
	"github.com/yourusername/expense-tracker/internal/handler/tui"
	"github.com/yourusername/expense-tracker/internal/repository/sqlite"
)

func main() {
	// Get database path - use user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	
	dbPath := filepath.Join(homeDir, ".expense-tracker", "expenses.db")
	
	// Initialize repository
	repository, err := sqlite.NewRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}
	defer repository.Close()
	
	// Initialize repositories
	transactionRepo := sqlite.NewTransactionRepository(repository)
	categoryRepo := sqlite.NewCategoryRepository(repository)
	
	// Initialize use cases
	transactionUseCase := usecase.NewTransactionUseCase(transactionRepo, categoryRepo)
	categoryUseCase := usecase.NewCategoryUseCase(categoryRepo, transactionRepo)
	
	// Initialize default categories
	if err := categoryUseCase.InitializeDefaultCategories(); err != nil {
		log.Printf("Warning: Failed to initialize default categories: %v", err)
	}
	
	// Create TUI application
	app := tui.NewApp(transactionUseCase, categoryUseCase)
	
	// Start the application
	program := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	
	fmt.Println("Starting Expense Tracker...")
	if _, err := program.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}