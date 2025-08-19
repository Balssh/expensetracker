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

// GetCategoryTotalsByDateRange returns category breakdowns for the given date range and transaction type
func (r *TransactionRepository) GetCategoryTotalsByDateRange(ctx context.Context, start, end time.Time, transactionType string) ([]*domain.CategoryBreakdown, error) {
	query := `
		SELECT 
			c.id, 
			c.name, 
			COALESCE(SUM(t.amount), 0) as total_amount,
			COUNT(t.id) as transaction_count
		FROM categories c
		LEFT JOIN transactions t ON c.id = t.category_id 
			AND t.date BETWEEN ? AND ? 
			AND t.type = ?
		WHERE c.type = ?
		GROUP BY c.id, c.name
		HAVING total_amount > 0
		ORDER BY total_amount DESC
	`

	rows, err := r.db.DB().QueryContext(ctx, query, 
		start.Format(time.RFC3339), 
		end.Format(time.RFC3339), 
		transactionType,
		transactionType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get category totals by date range: %w", err)
	}
	defer rows.Close()

	var breakdowns []*domain.CategoryBreakdown
	for rows.Next() {
		var categoryID int
		var categoryName string
		var totalAmount float64
		var transactionCount int

		err := rows.Scan(&categoryID, &categoryName, &totalAmount, &transactionCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category breakdown: %w", err)
		}

		category := &domain.Category{
			ID:   categoryID,
			Name: categoryName,
		}

		breakdown := &domain.CategoryBreakdown{
			Category:         category,
			TotalAmount:      totalAmount,
			TransactionCount: transactionCount,
		}

		breakdowns = append(breakdowns, breakdown)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate category breakdown rows: %w", err)
	}

	return breakdowns, nil
}

// GetTransactionCountByDateRange returns the count of transactions for the given date range and type
func (r *TransactionRepository) GetTransactionCountByDateRange(ctx context.Context, start, end time.Time, transactionType string) (int, error) {
	var query string
	var args []interface{}

	if transactionType == "" {
		// Count all transactions regardless of type
		query = `
			SELECT COUNT(*)
			FROM transactions
			WHERE date BETWEEN ? AND ?
		`
		args = []interface{}{start.Format(time.RFC3339), end.Format(time.RFC3339)}
	} else {
		// Count transactions of specific type
		query = `
			SELECT COUNT(*)
			FROM transactions
			WHERE date BETWEEN ? AND ? AND type = ?
		`
		args = []interface{}{start.Format(time.RFC3339), end.Format(time.RFC3339), transactionType}
	}

	var count int
	err := r.db.DB().QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction count by date range: %w", err)
	}

	return count, nil
}

// GetCategoryTransactionCount returns the count of transactions for a specific category, date range, and type
func (r *TransactionRepository) GetCategoryTransactionCount(ctx context.Context, start, end time.Time, categoryID int, transactionType string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM transactions
		WHERE date BETWEEN ? AND ? AND category_id = ? AND type = ?
	`

	var count int
	err := r.db.DB().QueryRowContext(ctx, query, 
		start.Format(time.RFC3339), 
		end.Format(time.RFC3339), 
		categoryID, 
		transactionType,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get category transaction count: %w", err)
	}

	return count, nil
}
