package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/expense-tracker/internal/core/domain"
	"github.com/yourusername/expense-tracker/internal/core/usecase"
)

// CategoryRepository implements the CategoryRepository interface for SQLite
type CategoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository creates a new SQLite category repository
func NewCategoryRepository(repository *Repository) *CategoryRepository {
	return &CategoryRepository{
		db: repository.db,
	}
}

// Create creates a new category
func (r *CategoryRepository) Create(category *domain.Category) error {
	query := `
		INSERT INTO categories (name, type)
		VALUES (?, ?)
	`

	result, err := r.db.Exec(query, category.Name, string(category.Type))
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	// Get the generated ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get inserted ID: %w", err)
	}

	category.ID = int(id)
	return nil
}

// GetByID retrieves a category by its ID
func (r *CategoryRepository) GetByID(id int) (*domain.Category, error) {
	query := `
		SELECT id, name, type
		FROM categories
		WHERE id = ?
	`

	row := r.db.QueryRow(query, id)

	category, err := r.scanCategory(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// GetByName retrieves a category by its name and type
func (r *CategoryRepository) GetByName(name string, categoryType domain.TransactionType) (*domain.Category, error) {
	query := `
		SELECT id, name, type
		FROM categories
		WHERE name = ? AND type = ?
	`

	row := r.db.QueryRow(query, name, string(categoryType))

	category, err := r.scanCategory(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, usecase.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get category by name: %w", err)
	}

	return category, nil
}

// ListByType retrieves all categories of a specific type
func (r *CategoryRepository) ListByType(categoryType domain.TransactionType) ([]*domain.Category, error) {
	query := `
		SELECT id, name, type
		FROM categories
		WHERE type = ?
		ORDER BY name
	`

	rows, err := r.db.Query(query, string(categoryType))
	if err != nil {
		return nil, fmt.Errorf("failed to list categories by type: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// List retrieves all categories
func (r *CategoryRepository) List() ([]*domain.Category, error) {
	query := `
		SELECT id, name, type
		FROM categories
		ORDER BY type, name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// Update updates an existing category
func (r *CategoryRepository) Update(category *domain.Category) error {
	query := `
		UPDATE categories
		SET name = ?, type = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, category.Name, string(category.Type), category.ID)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
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

// Delete deletes a category by its ID
func (r *CategoryRepository) Delete(id int) error {
	query := `DELETE FROM categories WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
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

// GetUsageCount returns the number of transactions using this category
func (r *CategoryRepository) GetUsageCount(id int) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE category_id = ?`

	var count int
	err := r.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get category usage count: %w", err)
	}

	return count, nil
}

// Exists checks if a category with the given name and type already exists
func (r *CategoryRepository) Exists(name string, categoryType domain.TransactionType) (bool, error) {
	query := `SELECT COUNT(*) FROM categories WHERE name = ? AND type = ?`

	var count int
	err := r.db.QueryRow(query, name, string(categoryType)).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if category exists: %w", err)
	}

	return count > 0, nil
}

// scanCategory scans a single row into a Category struct
func (r *CategoryRepository) scanCategory(row *sql.Row) (*domain.Category, error) {
	var category domain.Category
	var typeStr string

	err := row.Scan(
		&category.ID,
		&category.Name,
		&typeStr,
	)

	if err != nil {
		return nil, err
	}

	category.Type = domain.TransactionType(typeStr)
	return &category, nil
}

// scanCategories scans multiple rows into Category structs
func (r *CategoryRepository) scanCategories(rows *sql.Rows) ([]*domain.Category, error) {
	var categories []*domain.Category

	for rows.Next() {
		var category domain.Category
		var typeStr string

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&typeStr,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		category.Type = domain.TransactionType(typeStr)
		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return categories, nil
}