package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/yourusername/expense-tracker/internal/core/domain"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
)

// TransactionRepository implements the TransactionRepository interface for SQLite
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository creates a new SQLite transaction repository
func NewTransactionRepository(repository *Repository) *TransactionRepository {
	return &TransactionRepository{
		db: repository.db,
	}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, type, amount, category_id, description, date, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query,
		transaction.ID,
		string(transaction.Type),
		transaction.Amount,
		transaction.CategoryID,
		transaction.Description,
		transaction.Date.Format(time.RFC3339),
		transaction.CreatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by its ID
func (r *TransactionRepository) GetByID(id string) (*domain.Transaction, error) {
	query := `
		SELECT id, type, amount, category_id, description, date, created_at
		FROM transactions
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	transaction, err := r.scanTransaction(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return transaction, nil
}

// List retrieves transactions with optional limit and offset for pagination
func (r *TransactionRepository) List(limit, offset int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, type, amount, category_id, description, date, created_at
		FROM transactions
		ORDER BY date DESC, created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// ListByType retrieves transactions of a specific type
func (r *TransactionRepository) ListByType(transactionType domain.TransactionType, limit, offset int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, type, amount, category_id, description, date, created_at
		FROM transactions
		WHERE type = ?
		ORDER BY date DESC, created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, string(transactionType), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions by type: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// ListByDateRange retrieves transactions within a date range
func (r *TransactionRepository) ListByDateRange(start, end time.Time) ([]*domain.Transaction, error) {
	query := `
		SELECT id, type, amount, category_id, description, date, created_at
		FROM transactions
		WHERE date >= ? AND date <= ?
		ORDER BY date DESC, created_at DESC
	`

	rows, err := r.db.Query(query,
		start.Format(time.RFC3339),
		end.Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions by date range: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// ListByCategory retrieves transactions for a specific category
func (r *TransactionRepository) ListByCategory(categoryID int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, type, amount, category_id, description, date, created_at
		FROM transactions
		WHERE category_id = ?
		ORDER BY date DESC, created_at DESC
	`

	rows, err := r.db.Query(query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to list transactions by category: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// Update updates an existing transaction
func (r *TransactionRepository) Update(transaction *domain.Transaction) error {
	query := `
		UPDATE transactions
		SET type = ?, amount = ?, category_id = ?, description = ?, date = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query,
		string(transaction.Type),
		transaction.Amount,
		transaction.CategoryID,
		transaction.Description,
		transaction.Date.Format(time.RFC3339),
		transaction.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

// Delete deletes a transaction by its ID
func (r *TransactionRepository) Delete(id string) error {
	query := `DELETE FROM transactions WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return usecase.ErrNotFound
	}

	return nil
}

// Count returns the total number of transactions
func (r *TransactionRepository) Count() (int, error) {
	query := `SELECT COUNT(*) FROM transactions`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	return count, nil
}

// CountByType returns the number of transactions of a specific type
func (r *TransactionRepository) CountByType(transactionType domain.TransactionType) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE type = ?`

	var count int
	err := r.db.QueryRow(query, string(transactionType)).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count transactions by type: %w", err)
	}

	return count, nil
}

// scanTransaction scans a single row into a Transaction struct
func (r *TransactionRepository) scanTransaction(row *sql.Row) (*domain.Transaction, error) {
	var transaction domain.Transaction
	var typeStr, dateStr, createdAtStr string

	err := row.Scan(
		&transaction.ID,
		&typeStr,
		&transaction.Amount,
		&transaction.CategoryID,
		&transaction.Description,
		&dateStr,
		&createdAtStr,
	)

	if err != nil {
		return nil, err
	}

	// Parse transaction type
	transaction.Type = domain.TransactionType(typeStr)

	// Parse dates
	transaction.Date, err = time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	transaction.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	return &transaction, nil
}

// scanTransactions scans multiple rows into Transaction structs
func (r *TransactionRepository) scanTransactions(rows *sql.Rows) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction

	for rows.Next() {
		var transaction domain.Transaction
		var typeStr, dateStr, createdAtStr string

		err := rows.Scan(
			&transaction.ID,
			&typeStr,
			&transaction.Amount,
			&transaction.CategoryID,
			&transaction.Description,
			&dateStr,
			&createdAtStr,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Parse transaction type
		transaction.Type = domain.TransactionType(typeStr)

		// Parse dates
		transaction.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}

		transaction.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		transactions = append(transactions, &transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}