package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
	"github.com/gustavoz65/Fyapp/internal/validation"
)

type BudgetService struct {
	budgetRepo      *repository.BudgetRepository
	notificationSvc *NotificationService
	logger          *zerolog.Logger
}

func NewBudgetService(
	budgetRepo *repository.BudgetRepository,
	notificationSvc *NotificationService,
	logger *zerolog.Logger,
) *BudgetService {
	return &BudgetService{
		budgetRepo:      budgetRepo,
		notificationSvc: notificationSvc,
		logger:          logger,
	}
}

// Create creates a new budget
func (s *BudgetService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateBudgetRequest) (*model.Budget, error) {
	// Verificar limite de orçamentos ativos
	count, err := s.budgetRepo.CountActiveByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count budgets: %w", err)
	}
	if count >= validation.MaxActiveBudgets {
		return nil, fmt.Errorf(validation.ErrMaxActiveBudgetsExceeded)
	}

	budget := &model.Budget{
		UserID:     userID,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Amount:     req.Amount,
		PeriodType: req.PeriodType,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		IsActive:   true,
	}

	if req.AlertThreshold != nil {
		budget.AlertThreshold = *req.AlertThreshold
	} else {
		budget.AlertThreshold = decimal.NewFromInt(80)
	}

	if err := s.budgetRepo.Create(ctx, budget); err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	// Calculate initial spent amount
	if err := s.budgetRepo.RecalculateSpentAmount(ctx, budget); err != nil {
		s.logger.Warn().Err(err).Msg("failed to calculate initial spent amount")
	}

	// Refresh to get the calculated spent amount
	budget, _ = s.budgetRepo.GetByID(ctx, budget.ID)

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("budget_id", budget.ID.String()).
		Str("name", budget.Name).
		Msg("budget created")

	return budget, nil
}

// GetByID retrieves a budget by ID
func (s *BudgetService) GetByID(ctx context.Context, userID, budgetID uuid.UUID) (*model.Budget, error) {
	budget, err := s.budgetRepo.GetByIDAndUser(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}
	return budget, nil
}

// GetAll retrieves all budgets for a user
func (s *BudgetService) GetAll(ctx context.Context, userID uuid.UUID) ([]*model.Budget, error) {
	budgets, err := s.budgetRepo.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}
	return budgets, nil
}

// GetActive retrieves active budgets for the current period
func (s *BudgetService) GetActive(ctx context.Context, userID uuid.UUID) ([]*model.Budget, error) {
	budgets, err := s.budgetRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active budgets: %w", err)
	}
	return budgets, nil
}

// Update updates a budget
func (s *BudgetService) Update(ctx context.Context, userID, budgetID uuid.UUID, req *model.UpdateBudgetRequest) (*model.Budget, error) {
	budget, err := s.budgetRepo.GetByIDAndUser(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Name != nil {
		budget.Name = *req.Name
	}
	if req.Amount != nil {
		budget.Amount = *req.Amount
		// Reset alert sent status if amount changes
		budget.AlertSent = false
	}
	if req.AlertThreshold != nil {
		budget.AlertThreshold = *req.AlertThreshold
		budget.AlertSent = false
	}
	if req.IsActive != nil {
		budget.IsActive = *req.IsActive
	}

	if err := s.budgetRepo.Update(ctx, budget); err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("budget_id", budgetID.String()).
		Msg("budget updated")

	return budget, nil
}

// Delete soft deletes a budget
func (s *BudgetService) Delete(ctx context.Context, userID, budgetID uuid.UUID) error {
	if err := s.budgetRepo.Delete(ctx, budgetID, userID); err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("budget_id", budgetID.String()).
		Msg("budget deleted")

	return nil
}

// RecalculateSpentAmount recalculates the spent amount for a budget
func (s *BudgetService) RecalculateSpentAmount(ctx context.Context, budgetID uuid.UUID) error {
	budget, err := s.budgetRepo.GetByID(ctx, budgetID)
	if err != nil {
		return err
	}

	return s.budgetRepo.RecalculateSpentAmount(ctx, budget)
}

// CheckBudgetAlerts checks all budgets and creates alerts for those exceeding thresholds
func (s *BudgetService) CheckBudgetAlerts(ctx context.Context) error {
	budgets, err := s.budgetRepo.GetBudgetsNeedingAlert(ctx)
	if err != nil {
		return fmt.Errorf("failed to get budgets needing alert: %w", err)
	}

	for _, budget := range budgets {
		if err := s.notificationSvc.SendBudgetAlert(
			ctx,
			budget.UserID,
			budget.Name,
			budget.UsedPercentage().InexactFloat64(),
			budget.RemainingAmount().InexactFloat64(),
		); err != nil {
			s.logger.Error().Err(err).
				Str("budget_id", budget.ID.String()).
				Msg("failed to create budget alert notification")
			continue
		}

		// Mark alert as sent
		if err := s.budgetRepo.MarkAlertSent(ctx, budget.ID); err != nil {
			s.logger.Error().Err(err).
				Str("budget_id", budget.ID.String()).
				Msg("failed to mark budget alert as sent")
		}

		s.logger.Info().
			Str("budget_id", budget.ID.String()).
			Str("user_id", budget.UserID.String()).
			Msg("budget alert created")
	}

	return nil
}

// GetBudgetSummary returns a summary of budget usage for a user
func (s *BudgetService) GetBudgetSummary(ctx context.Context, userID uuid.UUID) (*model.BudgetAnalysisReport, error) {
	budgets, err := s.budgetRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active budgets: %w", err)
	}

	summary := &model.BudgetAnalysisReport{
		StartDate:     time.Now(),
		EndDate:       time.Now(),
		TotalBudgeted: decimal.Zero,
		TotalSpent:    decimal.Zero,
	}

	for _, budget := range budgets {
		summary.TotalBudgeted = summary.TotalBudgeted.Add(budget.Amount)
		summary.TotalSpent = summary.TotalSpent.Add(budget.SpentAmount)

		status := "ok"
		if budget.IsOverBudget() {
			status = "over"
			summary.OverBudgetCount++
		} else if budget.ShouldAlert() {
			status = "warning"
		} else {
			summary.UnderBudgetCount++
		}

		var categoryName *string
		if budget.Category != nil {
			categoryName = &budget.Category.Name
		}

		summary.BudgetDetails = append(summary.BudgetDetails, model.BudgetDetail{
			BudgetID:        budget.ID,
			BudgetName:      budget.Name,
			CategoryName:    categoryName,
			Amount:          budget.Amount,
			SpentAmount:     budget.SpentAmount,
			RemainingAmount: budget.RemainingAmount(),
			UsagePercentage: budget.UsedPercentage(),
			Status:          status,
		})
	}

	summary.TotalRemaining = summary.TotalBudgeted.Sub(summary.TotalSpent)
	if !summary.TotalBudgeted.IsZero() {
		summary.OverallUsage = summary.TotalSpent.Div(summary.TotalBudgeted).Mul(decimal.NewFromInt(100))
	}

	return summary, nil
}

// CreateMonthlyBudget creates a monthly budget based on predefined settings
func (s *BudgetService) CreateMonthlyBudget(ctx context.Context, userID uuid.UUID, categoryID *uuid.UUID, name string, amount decimal.Decimal) (*model.Budget, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	req := &model.CreateBudgetRequest{
		CategoryID: categoryID,
		Name:       name,
		Amount:     amount,
		PeriodType: model.BudgetPeriodMonthly,
		StartDate:  startDate,
		EndDate:    endDate,
	}

	return s.Create(ctx, userID, req)
}
