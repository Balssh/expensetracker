package sqlite

import (
	"context"
	"fmt"

	"expense-tracker/internal/core/domain"
)

type CategoryRepository struct {
	db *Database
}

func NewCategoryRepository(db *Database) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateExpenseCategory(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO expense_categories (name) VALUES (?)`
	
	result, err := r.db.DB().ExecContext(ctx, query, category.Name)
	if err != nil {
		return fmt.Errorf("failed to create expense category: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	category.ID = int(id)
	return nil
}

func (r *CategoryRepository) CreateIncomeCategory(ctx context.Context, category *domain.Category) error {
	query := `INSERT INTO income_categories (name) VALUES (?)`
	
	result, err := r.db.DB().ExecContext(ctx, query, category.Name)
	if err != nil {
		return fmt.Errorf("failed to create income category: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	category.ID = int(id)
	return nil
}

func (r *CategoryRepository) GetExpenseCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name FROM expense_categories ORDER BY name`

	rows, err := r.db.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense category: %w", err)
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryRepository) GetIncomeCategories(ctx context.Context) ([]*domain.Category, error) {
	query := `SELECT id, name FROM income_categories ORDER BY name`

	rows, err := r.db.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get income categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan income category: %w", err)
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryRepository) GetExpenseCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	query := `SELECT id, name FROM expense_categories WHERE id = ?`

	var category domain.Category
	err := r.db.DB().QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense category by id: %w", err)
	}

	return &category, nil
}

func (r *CategoryRepository) GetIncomeCategoryByID(ctx context.Context, id int) (*domain.Category, error) {
	query := `SELECT id, name FROM income_categories WHERE id = ?`

	var category domain.Category
	err := r.db.DB().QueryRowContext(ctx, query, id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get income category by id: %w", err)
	}

	return &category, nil
}