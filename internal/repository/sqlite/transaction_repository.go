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

func (r *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	var categoryID interface{}
	if transaction.Category != nil {
		categoryID = transaction.Category.ID
	}

	query := `
		INSERT INTO transactions (description, amount, date, type, category_id)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.DB().ExecContext(ctx, query,
		transaction.Description,
		transaction.Amount,
		transaction.Date.Format(time.RFC3339),
		transaction.Type,
		categoryID,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	transaction.ID = int(id)
	return nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int) (*domain.Transaction, error) {
	query := `
		SELECT t.id, t.description, t.amount, t.date, t.type, c.id, c.name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.id = ?
	`

	row := r.db.DB().QueryRowContext(ctx, query, id)
	return r.scanTransaction(row)
}

func (r *TransactionRepository) GetAll(ctx context.Context, offset, limit int) ([]*domain.Transaction, error) {
	query := `
		SELECT t.id, t.description, t.amount, t.date, t.type, c.id, c.name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		ORDER BY t.date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Transaction, error) {
	query := `
		SELECT t.id, t.description, t.amount, t.date, t.type, c.id, c.name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.date BETWEEN ? AND ?
		ORDER BY t.date DESC
	`

	rows, err := r.db.DB().QueryContext(ctx, query, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) GetByType(ctx context.Context, transactionType string, offset, limit int) ([]*domain.Transaction, error) {
	query := `
		SELECT t.id, t.description, t.amount, t.date, t.type, c.id, c.name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.type = ?
		ORDER BY t.date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB().QueryContext(ctx, query, transactionType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by type: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) GetTotalByDateRange(ctx context.Context, start, end time.Time, transactionType string) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE date BETWEEN ? AND ? AND type = ?
	`

	var total float64
	err := r.db.DB().QueryRowContext(ctx, query, start.Format(time.RFC3339), end.Format(time.RFC3339), transactionType).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total by date range: %w", err)
	}

	return total, nil
}

func (r *TransactionRepository) GetRecentTransactions(ctx context.Context, limit int) ([]*domain.Transaction, error) {
	query := `
		SELECT t.id, t.description, t.amount, t.date, t.type, c.id, c.name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		ORDER BY t.date DESC
		LIMIT ?
	`

	rows, err := r.db.DB().QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) SearchTransactions(ctx context.Context, searchQuery string, offset, limit int) ([]*domain.Transaction, error) {
	query := `
		SELECT t.id, t.description, t.amount, t.date, t.type, c.id, c.name
		FROM transactions t
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.description LIKE ? OR c.name LIKE ?
		ORDER BY t.date DESC
		LIMIT ? OFFSET ?
	`

	searchTerm := "%" + searchQuery + "%"
	rows, err := r.db.DB().QueryContext(ctx, query, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

func (r *TransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	var categoryID interface{}
	if transaction.Category != nil {
		categoryID = transaction.Category.ID
	}

	query := `
		UPDATE transactions 
		SET description = ?, amount = ?, date = ?, type = ?, category_id = ?
		WHERE id = ?
	`

	_, err := r.db.DB().ExecContext(ctx, query,
		transaction.Description,
		transaction.Amount,
		transaction.Date.Format(time.RFC3339),
		transaction.Type,
		categoryID,
		transaction.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM transactions WHERE id = ?`
	_, err := r.db.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

func (r *TransactionRepository) scanTransaction(row *sql.Row) (*domain.Transaction, error) {
	var transaction domain.Transaction
	var dateStr string
	var categoryID sql.NullInt64
	var categoryName sql.NullString

	err := row.Scan(
		&transaction.ID,
		&transaction.Description,
		&transaction.Amount,
		&dateStr,
		&transaction.Type,
		&categoryID,
		&categoryName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan transaction: %w", err)
	}

	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	transaction.Date = parsedDate

	if categoryID.Valid && categoryName.Valid {
		transaction.Category = &domain.Category{
			ID:   int(categoryID.Int64),
			Name: categoryName.String,
		}
	}

	return &transaction, nil
}

func (r *TransactionRepository) scanTransactions(rows *sql.Rows) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction

	for rows.Next() {
		var transaction domain.Transaction
		var dateStr string
		var categoryID sql.NullInt64
		var categoryName sql.NullString

		err := rows.Scan(
			&transaction.ID,
			&transaction.Description,
			&transaction.Amount,
			&dateStr,
			&transaction.Type,
			&categoryID,
			&categoryName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		transaction.Date = parsedDate

		if categoryID.Valid && categoryName.Valid {
			transaction.Category = &domain.Category{
				ID:   int(categoryID.Int64),
				Name: categoryName.String,
			}
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}