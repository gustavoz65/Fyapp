package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gustavoz65/finext/internal/model"
	"github.com/gustavoz65/finext/internal/repository"
	"github.com/rs/zerolog"
)

var (
	ErrCategoryInUse = errors.New("category is in use by transactions")
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	logger       *zerolog.Logger
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, logger *zerolog.Logger) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

// Create creates a new category
func (s *CategoryService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{
		UserID:   &userID,
		Name:     req.Name,
		Type:     req.Type,
		IsSystem: false,
		IsActive: true,
	}

	if req.Description != "" {
		category.Description = &req.Description
	}
	if req.Color != "" {
		category.Color = req.Color
	} else {
		category.Color = "#6366F1"
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	} else {
		category.Icon = "wallet"
	}

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		if errors.Is(err, repository.ErrCategoryAlreadyExists) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("category_id", category.ID.String()).
		Str("name", category.Name).
		Msg("category created")

	return category, nil
}

// GetByID retrieves a category by ID
func (s *CategoryService) GetByID(ctx context.Context, userID, categoryID uuid.UUID) (*model.Category, error) {
	category, err := s.categoryRepo.GetByIDAndUser(ctx, categoryID, userID)
	if err != nil {
		return nil, err
	}
	return category, nil
}

// GetAll retrieves all categories for a user (including system categories)
func (s *CategoryService) GetAll(ctx context.Context, userID uuid.UUID) ([]*model.Category, error) {
	categories, err := s.categoryRepo.GetAllForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

// GetByType retrieves categories by type
func (s *CategoryService) GetByType(ctx context.Context, userID uuid.UUID, categoryType model.CategoryType) ([]*model.Category, error) {
	categories, err := s.categoryRepo.GetByType(ctx, userID, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by type: %w", err)
	}
	return categories, nil
}

// GetIncomeCategories retrieves income categories
func (s *CategoryService) GetIncomeCategories(ctx context.Context, userID uuid.UUID) ([]*model.Category, error) {
	return s.GetByType(ctx, userID, model.CategoryTypeIncome)
}

// GetExpenseCategories retrieves expense categories
func (s *CategoryService) GetExpenseCategories(ctx context.Context, userID uuid.UUID) ([]*model.Category, error) {
	return s.GetByType(ctx, userID, model.CategoryTypeExpense)
}

// Update updates a category
func (s *CategoryService) Update(ctx context.Context, userID, categoryID uuid.UUID, req *model.UpdateCategoryRequest) (*model.Category, error) {
	category, err := s.categoryRepo.GetByIDAndUser(ctx, categoryID, userID)
	if err != nil {
		return nil, err
	}

	if category.IsSystem {
		return nil, repository.ErrCannotModifySystem
	}

	// Apply updates
	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.Description != nil {
		category.Description = req.Description
	}
	if req.Color != nil {
		category.Color = *req.Color
	}
	if req.Icon != nil {
		category.Icon = *req.Icon
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("category_id", categoryID.String()).
		Msg("category updated")

	return category, nil
}

// Delete deletes a category
func (s *CategoryService) Delete(ctx context.Context, userID, categoryID uuid.UUID) error {
	category, err := s.categoryRepo.GetByIDAndUser(ctx, categoryID, userID)
	if err != nil {
		return err
	}

	if category.IsSystem {
		return repository.ErrCannotModifySystem
	}

	// Check if category is in use
	count, err := s.categoryRepo.CountTransactionsByCategory(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to check category usage: %w", err)
	}

	if count > 0 {
		return ErrCategoryInUse
	}

	if err := s.categoryRepo.Delete(ctx, categoryID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("category_id", categoryID.String()).
		Msg("category deleted")

	return nil
}
