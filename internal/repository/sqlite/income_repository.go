package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"expense-tracker/internal/core/domain"
)

type IncomeRepository struct {
	db *Database
}

func NewIncomeRepository(db *Database) *IncomeRepository {
	return &IncomeRepository{db: db}
}

func (r *IncomeRepository) Create(ctx context.Context, income *domain.Income) error {
	query := `
		INSERT INTO income (description, amount, date, category_id)
		VALUES (?, ?, ?, ?)
	`
	
	result, err := r.db.DB().ExecContext(ctx, query, 
		income.Description, 
		income.Amount, 
		income.Date.Format(time.RFC3339), 
		income.CategoryID,
	)
	if err != nil {
		return fmt.Errorf("failed to create income: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	income.ID = int(id)
	return nil
}

func (r *IncomeRepository) GetByID(ctx context.Context, id int) (*domain.Income, error) {
	query := `
		SELECT i.id, i.description, i.amount, i.date, i.category_id, c.name
		FROM income i
		LEFT JOIN income_categories c ON i.category_id = c.id
		WHERE i.id = ?
	`

	var income domain.Income
	var category domain.Category
	var dateStr string
	var categoryName sql.NullString

	err := r.db.DB().QueryRowContext(ctx, query, id).Scan(
		&income.ID,
		&income.Description,
		&income.Amount,
		&dateStr,
		&income.CategoryID,
		&categoryName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get income by id: %w", err)
	}

	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}
	income.Date = parsedDate

	if categoryName.Valid {
		category.ID = income.CategoryID
		category.Name = categoryName.String
		income.Category = &category
	}

	return &income, nil
}

func (r *IncomeRepository) GetAll(ctx context.Context, offset, limit int) ([]*domain.Income, error) {
	query := `
		SELECT i.id, i.description, i.amount, i.date, i.category_id, c.name
		FROM income i
		LEFT JOIN income_categories c ON i.category_id = c.id
		ORDER BY i.date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB().QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get incomes: %w", err)
	}
	defer rows.Close()

	var incomes []*domain.Income
	for rows.Next() {
		var income domain.Income
		var category domain.Category
		var dateStr string
		var categoryName sql.NullString

		err := rows.Scan(
			&income.ID,
			&income.Description,
			&income.Amount,
			&dateStr,
			&income.CategoryID,
			&categoryName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan income: %w", err)
		}

		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		income.Date = parsedDate

		if categoryName.Valid {
			category.ID = income.CategoryID
			category.Name = categoryName.String
			income.Category = &category
		}

		incomes = append(incomes, &income)
	}

	return incomes, nil
}

func (r *IncomeRepository) GetByDateRange(ctx context.Context, start, end time.Time) ([]*domain.Income, error) {
	query := `
		SELECT i.id, i.description, i.amount, i.date, i.category_id, c.name
		FROM income i
		LEFT JOIN income_categories c ON i.category_id = c.id
		WHERE i.date >= ? AND i.date <= ?
		ORDER BY i.date DESC
	`

	rows, err := r.db.DB().QueryContext(ctx, query, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, fmt.Errorf("failed to get incomes by date range: %w", err)
	}
	defer rows.Close()

	var incomes []*domain.Income
	for rows.Next() {
		var income domain.Income
		var category domain.Category
		var dateStr string
		var categoryName sql.NullString

		err := rows.Scan(
			&income.ID,
			&income.Description,
			&income.Amount,
			&dateStr,
			&income.CategoryID,
			&categoryName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan income: %w", err)
		}

		parsedDate, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %w", err)
		}
		income.Date = parsedDate

		if categoryName.Valid {
			category.ID = income.CategoryID
			category.Name = categoryName.String
			income.Category = &category
		}

		incomes = append(incomes, &income)
	}

	return incomes, nil
}

func (r *IncomeRepository) GetTotalByDateRange(ctx context.Context, start, end time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM income
		WHERE date >= ? AND date <= ?
	`

	var total float64
	err := r.db.DB().QueryRowContext(ctx, query, start.Format(time.RFC3339), end.Format(time.RFC3339)).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total incomes by date range: %w", err)
	}

	return total, nil
}

func (r *IncomeRepository) Update(ctx context.Context, income *domain.Income) error {
	query := `
		UPDATE income
		SET description = ?, amount = ?, date = ?, category_id = ?
		WHERE id = ?
	`

	_, err := r.db.DB().ExecContext(ctx, query,
		income.Description,
		income.Amount,
		income.Date.Format(time.RFC3339),
		income.CategoryID,
		income.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update income: %w", err)
	}

	return nil
}

func (r *IncomeRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM income WHERE id = ?`

	_, err := r.db.DB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete income: %w", err)
	}

	return nil
}