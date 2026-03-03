package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoz65/Fyapp/internal/database"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/rs/zerolog"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCannotModifySystem    = errors.New("cannot modify system category")
)

type CategoryRepository struct {
	*BaseRepository
}

func NewCategoryRepository(db *database.Database, logger *zerolog.Logger) *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new category
func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	category.ID = uuid.New()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	query := `
		INSERT INTO categories (
			id, user_id, name, description, type, color, icon,
			is_system, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var userID sql.NullString
	if category.UserID != nil {
		userID = sql.NullString{String: category.UserID.String(), Valid: true}
	}

	_, err := r.ExecContext(ctx, query,
		category.ID.String(),
		userID,
		category.Name,
		NullString(category.Description),
		category.Type,
		category.Color,
		category.Icon,
		category.IsSystem,
		category.IsActive,
		category.CreatedAt,
		category.UpdatedAt,
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrCategoryAlreadyExists
		}
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Category, error) {
	query := `
		SELECT id, user_id, name, description, type, color, icon,
			is_system, is_active, created_at, updated_at
		FROM categories
		WHERE id = ?
	`

	return r.scanCategory(r.QueryRowContext(ctx, query, id.String()))
}

// GetByIDAndUser retrieves a category by ID ensuring it belongs to user or is system
func (r *CategoryRepository) GetByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*model.Category, error) {
	query := `
		SELECT id, user_id, name, description, type, color, icon,
			is_system, is_active, created_at, updated_at
		FROM categories
		WHERE id = ? AND (user_id = ? OR user_id IS NULL)
	`

	return r.scanCategory(r.QueryRowContext(ctx, query, id.String(), userID.String()))
}

// GetAllForUser retrieves all categories for a user (including system categories)
func (r *CategoryRepository) GetAllForUser(ctx context.Context, userID uuid.UUID) ([]*model.Category, error) {
	query := `
		SELECT id, user_id, name, description, type, color, icon,
			is_system, is_active, created_at, updated_at
		FROM categories
		WHERE (user_id = ? OR user_id IS NULL) AND is_active = TRUE
		ORDER BY is_system DESC, name ASC
	`

	rows, err := r.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// GetByType retrieves categories by type for a user
func (r *CategoryRepository) GetByType(ctx context.Context, userID uuid.UUID, categoryType model.CategoryType) ([]*model.Category, error) {
	query := `
		SELECT id, user_id, name, description, type, color, icon,
			is_system, is_active, created_at, updated_at
		FROM categories
		WHERE (user_id = ? OR user_id IS NULL) AND type = ? AND is_active = TRUE
		ORDER BY is_system DESC, name ASC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by type: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// GetSystemCategories retrieves all system categories
func (r *CategoryRepository) GetSystemCategories(ctx context.Context) ([]*model.Category, error) {
	query := `
		SELECT id, user_id, name, description, type, color, icon,
			is_system, is_active, created_at, updated_at
		FROM categories
		WHERE is_system = TRUE AND is_active = TRUE
		ORDER BY type, name ASC
	`

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get system categories: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// Update updates a category
func (r *CategoryRepository) Update(ctx context.Context, category *model.Category) error {
	if category.IsSystem {
		return ErrCannotModifySystem
	}

	category.UpdatedAt = time.Now()

	query := `
		UPDATE categories SET
			name = ?, description = ?, color = ?, icon = ?,
			is_active = ?, updated_at = ?
		WHERE id = ? AND is_system = FALSE
	`

	result, err := r.ExecContext(ctx, query,
		category.Name,
		NullString(category.Description),
		category.Color,
		category.Icon,
		category.IsActive,
		category.UpdatedAt,
		category.ID.String(),
	)

	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrCategoryAlreadyExists
		}
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

// Delete soft deletes a category (marks as inactive)
func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE categories
		SET is_active = FALSE, updated_at = ?
		WHERE id = ? AND is_system = FALSE
	`

	result, err := r.ExecContext(ctx, query, time.Now(), id.String())
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

// CountTransactionsByCategory counts transactions using a category
func (r *CategoryRepository) CountTransactionsByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE category_id = ?`

	var count int64
	err := r.QueryRowContext(ctx, query, categoryID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	return count, nil
}

// Helper functions

func (r *CategoryRepository) scanCategory(row *sql.Row) (*model.Category, error) {
	category := &model.Category{}
	var userID, description sql.NullString

	err := row.Scan(
		&category.ID,
		&userID,
		&category.Name,
		&description,
		&category.Type,
		&category.Color,
		&category.Icon,
		&category.IsSystem,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("failed to scan category: %w", err)
	}

	if userID.Valid {
		id, _ := uuid.Parse(userID.String)
		category.UserID = &id
	}
	category.Description = StringPtr(description)

	return category, nil
}

func (r *CategoryRepository) scanCategories(rows *sql.Rows) ([]*model.Category, error) {
	var categories []*model.Category

	for rows.Next() {
		category := &model.Category{}
		var userID, description sql.NullString

		err := rows.Scan(
			&category.ID,
			&userID,
			&category.Name,
			&description,
			&category.Type,
			&category.Color,
			&category.Icon,
			&category.IsSystem,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		if userID.Valid {
			id, _ := uuid.Parse(userID.String)
			category.UserID = &id
		}
		category.Description = StringPtr(description)

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}
