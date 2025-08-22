package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Repository manages the SQLite database connection and operations
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new SQLite repository
func NewRepository(dbPath string) (*Repository, error) {
	// Ensure the directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repository := &Repository{db: db}

	// Initialize database schema
	if err := repository.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repository, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}

// initSchema creates the database tables if they don't exist
func (r *Repository) initSchema() error {
	// Create categories table
	categoryTableSQL := `
	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
		UNIQUE(name, type)
	);`

	if _, err := r.db.Exec(categoryTableSQL); err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}

	// Create transactions table
	transactionTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL CHECK (type IN ('income', 'expense')),
		amount REAL NOT NULL,
		category_id INTEGER NOT NULL,
		description TEXT NOT NULL,
		date TEXT NOT NULL,
		created_at TEXT NOT NULL,
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);`

	if _, err := r.db.Exec(transactionTableSQL); err != nil {
		return fmt.Errorf("failed to create transactions table: %w", err)
	}

	// Create indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions(date);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);",
		"CREATE INDEX IF NOT EXISTS idx_transactions_category ON transactions(category_id);",
		"CREATE INDEX IF NOT EXISTS idx_categories_type ON categories(type);",
	}

	for _, indexSQL := range indexes {
		if _, err := r.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// BeginTx starts a new transaction
func (r *Repository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

// GetDB returns the database connection for testing purposes
func (r *Repository) GetDB() *sql.DB {
	return r.db
}