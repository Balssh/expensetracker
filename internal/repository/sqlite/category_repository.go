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

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category, categoryType string) error {
	query := `INSERT INTO categories (name, type) VALUES (?, ?)`
	result, err := r.db.DB().ExecContext(ctx, query, category.Name, categoryType)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	category.ID = int(id)
	return nil
}

func (r *CategoryRepository) GetCategories(ctx context.Context, categoryType string) ([]*domain.Category, error) {
	query := `SELECT id, name FROM categories WHERE type = ? ORDER BY name`
	rows, err := r.db.DB().QueryContext(ctx, query, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &category)
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int, categoryType string) (*domain.Category, error) {
	query := `SELECT id, name FROM categories WHERE id = ? AND type = ?`
	var category domain.Category
	err := r.db.DB().QueryRowContext(ctx, query, id, categoryType).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by id: %w", err)
	}
	return &category, nil
}