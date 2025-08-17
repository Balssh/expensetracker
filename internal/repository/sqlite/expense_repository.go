package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"expense-tracker/internal/core/domain"
)

type ExpenseRepository struct {
	db *Database
}

func NewExpenseRepository(db *Database) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	query := `
		INSERT INTO expenses (description, amount, date, category_id)
		VALUES (?, ?, ?, ?)
	`
	
	result, err := r.db.DB().ExecContext(ctx, query, 
		expense.Description, 
		expense.Amount, 
		expense.Date.Format(time.RFC3339), 
		expense.CategoryID,
	)
	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	expense.ID = int(id)
	return nil
}

func (r *ExpenseRepository) GetByID(ctx context.Context, id int) (*domain.Expense, error) {
	query := `
		SELECT e.id, e.description, e.amount, e.date, e.category_id, c.name
		FROM expenses e
		LEFT JOIN expense_categories c ON e.category_id = c.id
		WHERE e.id = ?
	`

	var expense domain.Expense
	var category domain.Category
	var dateStr string
	var categoryName sql.NullString

	err := r.db.DB().QueryRowContext(ctx, query, id).Scan(
		&expense.ID,
		&expense.Description,
		&expense.Amount,
		&dateStr,
		&expense.CategoryID,
		&categoryName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense by id: %w", err)
	}

	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	expense.Date = parsedDate

	if categoryName.Valid {
		category.ID = expense.CategoryID
		category.Name = categoryName.String
		expense.Category = &category
	}

	return &expense, nil
}

func (r *ExpenseRepository) GetAll(ctx context.Context, offset, limit int) ([]*domain.Expense, error) {
	query := `
		SELECT e.id, e.description, e.amount, e.date, e.category_id, c.name
		FROM expenses e
		LEFT JOIN expense_categories c ON e.category_id = c.id
		ORDER BY e.date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses: %w", err)
	}
	defer rows.Close()

	var expenses []*domain.Expense
	for rows.Next() {
		var expense domain.Expense
		var category domain.Category
		var dateStr string
		var categoryName sql.NullString

		err := rows.Scan(
			&expense.ID,
			&expense.Description,
			&expense.Amount,
			&dateStr,
			&expense.CategoryID,
			&categoryName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}

		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		expense.Date = parsedDate

		if categoryName.Valid {
			category.ID = expense.CategoryID
			category.Name = categoryName.String
			expense.Category = &category
		}

		expenses = append(expenses, &expense)
	}

	return expenses, nil
}

func (r *ExpenseRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Expense, error) {
	query := `
		SELECT e.id, e.description, e.amount, e.date, e.category_id, c.name
		FROM expenses e
		LEFT JOIN expense_categories c ON e.category_id = c.id
		WHERE e.date >= ? AND e.date <= ?
		ORDER BY e.date DESC
	`

	rows, err := r.db.DB().QueryContext(ctx, query, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("failed to get expenses by date range: %w", err)
	}
	defer rows.Close()

	var expenses []*domain.Expense
	for rows.Next() {
		var expense domain.Expense
		var category domain.Category
		var dateStr string
		var categoryName sql.NullString

		err := rows.Scan(
			&expense.ID,
			&expense.Description,
			&expense.Amount,
			&dateStr,
			&expense.CategoryID,
			&categoryName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}

		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		expense.Date = parsedDate

		if categoryName.Valid {
			category.ID = expense.CategoryID
			category.Name = categoryName.String
			expense.Category = &category
		}

		expenses = append(expenses, &expense)
	}

	return expenses, nil
}

func (r *ExpenseRepository) GetTotalByDateRange(ctx context.Context, start, end time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE date >= ? AND date <= ?
	`

	var total float64
	err := r.db.DB().QueryRowContext(ctx, query, start.Format(time.RFC3339), end.Format(time.RFC3339)).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total expenses by date range: %w", err)
	}

	return total, nil
}

func (r *ExpenseRepository) Update(ctx context.Context, expense *domain.Expense) error {
	query := `
		UPDATE expenses
		SET description = ?, amount = ?, date = ?, category_id = ?
		WHERE id = ?
	`

	_, err := r.db.DB().ExecContext(ctx, query,
		expense.Description,
		expense.Amount,
		expense.Date.Format(time.RFC3339),
		expense.CategoryID,
		expense.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update expense: %w", err)
	}

	return nil
}

func (r *ExpenseRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM expenses WHERE id = ?`

	_, err := r.db.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}

	return nil
}