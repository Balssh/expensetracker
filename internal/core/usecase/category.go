package usecase

import (
	"fmt"

	"github.com/yourusername/expense-tracker/internal/core/domain"
)

// CategoryUseCase handles business logic for categories
type CategoryUseCase struct {
	categoryRepo    CategoryRepository
	transactionRepo TransactionRepository
}

// NewCategoryUseCase creates a new CategoryUseCase
func NewCategoryUseCase(categoryRepo CategoryRepository, transactionRepo TransactionRepository) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo:    categoryRepo,
		transactionRepo: transactionRepo,
	}
}

// GetCategories retrieves all categories of a specific type
func (uc *CategoryUseCase) GetCategories(categoryType domain.TransactionType) ([]*domain.Category, error) {
	categories, err := uc.categoryRepo.ListByType(categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

// GetAllCategories retrieves all categories
func (uc *CategoryUseCase) GetAllCategories() ([]*domain.Category, error) {
	categories, err := uc.categoryRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	return categories, nil
}

// AddCategory creates a new category
func (uc *CategoryUseCase) AddCategory(name string, categoryType domain.TransactionType) (*domain.Category, error) {
	// Check if category already exists
	exists, err := uc.categoryRepo.Exists(name, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to check if category exists: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("category '%s' already exists for type '%s'", name, categoryType)
	}

	// Create the category
	category := domain.NewCategory(name, categoryType)

	// Validate the category
	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// Save to repository
	if err := uc.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// UpdateCategory updates an existing category
func (uc *CategoryUseCase) UpdateCategory(id int, name string) error {
	// Get existing category
	category, err := uc.categoryRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	// Check if new name conflicts with existing categories
	if category.Name != name {
		exists, err := uc.categoryRepo.Exists(name, category.Type)
		if err != nil {
			return fmt.Errorf("failed to check if category exists: %w", err)
		}

		if exists {
			return fmt.Errorf("category '%s' already exists for type '%s'", name, category.Type)
		}
	}

	// Update the category
	category.Name = name

	// Validate the updated category
	if err := category.Validate(); err != nil {
		return fmt.Errorf("invalid category: %w", err)
	}

	// Save to repository
	if err := uc.categoryRepo.Update(category); err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

// DeleteCategory deletes a category if it's not in use
func (uc *CategoryUseCase) DeleteCategory(id int) error {
	// Check if category exists
	_, err := uc.categoryRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("category not found: %w", err)
	}

	// Check if category is in use
	usageCount, err := uc.categoryRepo.GetUsageCount(id)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}

	if usageCount > 0 {
		return fmt.Errorf("cannot delete category: it is used by %d transaction(s)", usageCount)
	}

	// Delete from repository
	if err := uc.categoryRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// CanDeleteCategory checks if a category can be deleted
func (uc *CategoryUseCase) CanDeleteCategory(id int) (bool, string, error) {
	// Check if category exists
	_, err := uc.categoryRepo.GetByID(id)
	if err != nil {
		return false, "Category not found", err
	}

	// Check usage count
	usageCount, err := uc.categoryRepo.GetUsageCount(id)
	if err != nil {
		return false, "Failed to check category usage", err
	}

	if usageCount > 0 {
		return false, fmt.Sprintf("Category is used by %d transaction(s)", usageCount), nil
	}

	return true, "", nil
}

// GetCategoryUsage returns usage statistics for a category
func (uc *CategoryUseCase) GetCategoryUsage(id int) (int, error) {
	usageCount, err := uc.categoryRepo.GetUsageCount(id)
	if err != nil {
		return 0, fmt.Errorf("failed to get category usage: %w", err)
	}
	return usageCount, nil
}

// InitializeDefaultCategories creates default categories if they don't exist
func (uc *CategoryUseCase) InitializeDefaultCategories() error {
	// Initialize default expense categories
	for _, category := range domain.DefaultExpenseCategories() {
		exists, err := uc.categoryRepo.Exists(category.Name, category.Type)
		if err != nil {
			return fmt.Errorf("failed to check if category exists: %w", err)
		}

		if !exists {
			if err := uc.categoryRepo.Create(&category); err != nil {
				return fmt.Errorf("failed to create default expense category '%s': %w", category.Name, err)
			}
		}
	}

	// Initialize default income categories
	for _, category := range domain.DefaultIncomeCategories() {
		exists, err := uc.categoryRepo.Exists(category.Name, category.Type)
		if err != nil {
			return fmt.Errorf("failed to check if category exists: %w", err)
		}

		if !exists {
			if err := uc.categoryRepo.Create(&category); err != nil {
				return fmt.Errorf("failed to create default income category '%s': %w", category.Name, err)
			}
		}
	}

	return nil
}
