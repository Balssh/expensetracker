package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"expense-tracker/internal/core/domain"
)

type TransactionRepository struct {
	db *Database
}

func NewTransactionRepository(db *Database) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) GetRecentTransactions(ctx context.Context, limit int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, description, amount, date, category, type FROM (
			SELECT e.id, e.description, e.amount, e.date, c.name as category, 'expense' as type
			FROM expenses e
			LEFT JOIN expense_categories c ON e.category_id = c.id
			UNION ALL
			SELECT i.id, i.description, i.amount, i.date, c.name as category, 'income' as type
			FROM income i
			LEFT JOIN income_categories c ON i.category_id = c.id
		) transactions
		ORDER BY date DESC
		LIMIT ?
	`

	rows, err := r.db.DB().QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) GetAllTransactions(ctx context.Context, offset, limit int) ([]*domain.Transaction, error) {
	query := `
		SELECT id, description, amount, date, category, type FROM (
			SELECT e.id, e.description, e.amount, e.date, c.name as category, 'expense' as type
			FROM expenses e
			LEFT JOIN expense_categories c ON e.category_id = c.id
			UNION ALL
			SELECT i.id, i.description, i.amount, i.date, c.name as category, 'income' as type
			FROM income i
			LEFT JOIN income_categories c ON i.category_id = c.id
		) transactions
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) SearchTransactions(ctx context.Context, query string, offset, limit int) ([]*domain.Transaction, error) {
	searchQuery := `
		SELECT id, description, amount, date, category, type FROM (
			SELECT e.id, e.description, e.amount, e.date, c.name as category, 'expense' as type
			FROM expenses e
			LEFT JOIN expense_categories c ON e.category_id = c.id
			WHERE e.description LIKE ? OR c.name LIKE ?
			UNION ALL
			SELECT i.id, i.description, i.amount, i.date, c.name as category, 'income' as type
			FROM income i
			LEFT JOIN income_categories c ON i.category_id = c.id
			WHERE i.description LIKE ? OR c.name LIKE ?
		) transactions
		ORDER BY date DESC
		LIMIT ? OFFSET ?
	`

	searchTerm := "%" + query + "%"
	rows, err := r.db.DB().QueryContext(ctx, searchQuery, searchTerm, searchTerm, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) scanTransactions(rows *sql.Rows) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction

	for rows.Next() {
		var transaction domain.Transaction
		var dateStr string
		var category sql.NullString

		err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&dateStr,
			&category,
			&transaction.Type,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		transaction.Date = parsedDate

		if category.Valid {
			transaction.Category = category.String
		} else {
			transaction.Category = "Unknown"
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}