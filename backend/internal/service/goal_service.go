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

type GoalService struct {
	goalRepo        *repository.GoalRepository
	notificationSvc *NotificationService
	logger          *zerolog.Logger
}

func NewGoalService(
	goalRepo *repository.GoalRepository,
	notificationSvc *NotificationService,
	logger *zerolog.Logger,
) *GoalService {
	return &GoalService{
		goalRepo:        goalRepo,
		notificationSvc: notificationSvc,
		logger:          logger,
	}
}

// Create creates a new goal
func (s *GoalService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateGoalRequest) (*model.Goal, error) {
	// Verificar limite de metas ativas
	count, err := s.goalRepo.CountActiveByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count goals: %w", err)
	}
	if count >= validation.MaxActiveGoals {
		return nil, fmt.Errorf(validation.ErrMaxActiveGoalsExceeded) //nolint:ST1005 // user-facing message
	}

	goal := &model.Goal{
		UserID:       userID,
		Name:         req.Name,
		TargetAmount: req.TargetAmount,
		TargetDate:   req.TargetDate,
	}

	if req.Description != "" {
		goal.Description = &req.Description
	}
	if req.Icon != "" {
		goal.Icon = req.Icon
	} else {
		goal.Icon = "target"
	}
	if req.Color != "" {
		goal.Color = req.Color
	} else {
		goal.Color = "#F59E0B"
	}
	if req.Priority != nil {
		goal.Priority = *req.Priority
	} else {
		goal.Priority = 1
	}

	if err := s.goalRepo.Create(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to create goal: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("goal_id", goal.ID.String()).
		Str("name", goal.Name).
		Msg("goal created")

	return goal, nil
}

// GetByID retrieves a goal by ID
func (s *GoalService) GetByID(ctx context.Context, userID, goalID uuid.UUID) (*model.Goal, error) {
	goal, err := s.goalRepo.GetByIDAndUser(ctx, goalID, userID)
	if err != nil {
		return nil, err
	}

	// Load contributions
	contributions, err := s.goalRepo.GetContributionsByGoal(ctx, goalID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to load goal contributions")
	} else {
		goal.Contributions = contributions
	}

	return goal, nil
}

// GetAll retrieves all goals for a user
func (s *GoalService) GetAll(ctx context.Context, userID uuid.UUID) ([]*model.Goal, error) {
	goals, err := s.goalRepo.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goals: %w", err)
	}
	return goals, nil
}

// GetActive retrieves active (in progress) goals
func (s *GoalService) GetActive(ctx context.Context, userID uuid.UUID) ([]*model.Goal, error) {
	goals, err := s.goalRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active goals: %w", err)
	}
	return goals, nil
}

// Update updates a goal
func (s *GoalService) Update(ctx context.Context, userID, goalID uuid.UUID, req *model.UpdateGoalRequest) (*model.Goal, error) {
	goal, err := s.goalRepo.GetByIDAndUser(ctx, goalID, userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Name != nil {
		goal.Name = *req.Name
	}
	if req.Description != nil {
		goal.Description = req.Description
	}
	if req.TargetAmount != nil {
		goal.TargetAmount = *req.TargetAmount
	}
	if req.TargetDate != nil {
		goal.TargetDate = req.TargetDate
	}
	if req.Icon != nil {
		goal.Icon = *req.Icon
	}
	if req.Color != nil {
		goal.Color = *req.Color
	}
	if req.Priority != nil {
		goal.Priority = *req.Priority
	}
	if req.Status != nil {
		goal.Status = *req.Status
		if *req.Status == model.GoalStatusCompleted {
			now := time.Now()
			goal.CompletedAt = &now
		}
	}

	if err := s.goalRepo.Update(ctx, goal); err != nil {
		return nil, fmt.Errorf("failed to update goal: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("goal_id", goalID.String()).
		Msg("goal updated")

	return goal, nil
}

// Delete deletes a goal
func (s *GoalService) Delete(ctx context.Context, userID, goalID uuid.UUID) error {
	if err := s.goalRepo.Delete(ctx, goalID, userID); err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("goal_id", goalID.String()).
		Msg("goal deleted")

	return nil
}

// AddContribution adds a contribution to a goal
func (s *GoalService) AddContribution(ctx context.Context, userID, goalID uuid.UUID, req *model.CreateGoalContributionRequest) (*model.GoalContribution, error) {
	goal, err := s.goalRepo.GetByIDAndUser(ctx, goalID, userID)
	if err != nil {
		return nil, err
	}

	if goal.Status != model.GoalStatusInProgress {
		return nil, fmt.Errorf("cannot add contribution to a %s goal", goal.Status)
	}

	contribution := &model.GoalContribution{
		GoalID: goalID,
		Amount: req.Amount,
	}

	if req.Note != "" {
		contribution.Note = &req.Note
	}

	if req.ContributionDate != nil {
		contribution.ContributionDate = *req.ContributionDate
	} else {
		contribution.ContributionDate = time.Now()
	}

	if err := s.goalRepo.AddContribution(ctx, contribution); err != nil {
		return nil, fmt.Errorf("failed to add contribution: %w", err)
	}

	// Update goal current amount (trigger will handle this, but we do it here for consistency)
	newAmount := goal.CurrentAmount.Add(req.Amount)
	if err := s.goalRepo.UpdateCurrentAmount(ctx, goalID, newAmount); err != nil {
		s.logger.Warn().Err(err).Msg("failed to update goal current amount")
	}

	// Check if goal is now completed
	if newAmount.GreaterThanOrEqual(goal.TargetAmount) {
		if err := s.goalRepo.MarkCompleted(ctx, goalID); err != nil {
			s.logger.Warn().Err(err).Msg("failed to mark goal as completed")
		}

		if err := s.notificationSvc.SendGoalAchievedNotification(ctx, userID, goal.Name, goal.TargetAmount.InexactFloat64()); err != nil {
			s.logger.Warn().Err(err).Msg("failed to create goal achieved notification")
		}

		s.logger.Info().
			Str("user_id", userID.String()).
			Str("goal_id", goalID.String()).
			Msg("goal completed")
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("goal_id", goalID.String()).
		Str("amount", req.Amount.String()).
		Msg("contribution added to goal")

	return contribution, nil
}

// GetContributions retrieves contributions for a goal
func (s *GoalService) GetContributions(ctx context.Context, userID, goalID uuid.UUID) ([]model.GoalContribution, error) {
	// Verify ownership
	_, err := s.goalRepo.GetByIDAndUser(ctx, goalID, userID)
	if err != nil {
		return nil, err
	}

	contributions, err := s.goalRepo.GetContributionsByGoal(ctx, goalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributions: %w", err)
	}

	return contributions, nil
}

// DeleteContribution deletes a contribution from a goal
func (s *GoalService) DeleteContribution(ctx context.Context, userID, goalID, contributionID uuid.UUID) error {
	// Verify ownership
	_, err := s.goalRepo.GetByIDAndUser(ctx, goalID, userID)
	if err != nil {
		return err
	}

	if err := s.goalRepo.DeleteContribution(ctx, contributionID, goalID); err != nil {
		return fmt.Errorf("failed to delete contribution: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("goal_id", goalID.String()).
		Str("contribution_id", contributionID.String()).
		Msg("contribution deleted from goal")

	return nil
}

// GetSummary returns a summary of all goals for a user
func (s *GoalService) GetSummary(ctx context.Context, userID uuid.UUID) (*model.GoalSummary, error) {
	summary, err := s.goalRepo.GetGoalSummary(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get goal summary: %w", err)
	}
	return summary, nil
}

// CalculateRequiredContribution calculates the required monthly contribution to reach a goal
func (s *GoalService) CalculateRequiredContribution(targetAmount, currentAmount decimal.Decimal, targetDate *time.Time) *decimal.Decimal {
	if targetDate == nil {
		return nil
	}

	remaining := targetAmount.Sub(currentAmount)
	if remaining.LessThanOrEqual(decimal.Zero) {
		zero := decimal.Zero
		return &zero
	}

	daysRemaining := time.Until(*targetDate).Hours() / 24
	if daysRemaining <= 0 {
		return nil
	}

	monthsRemaining := daysRemaining / 30.0
	if monthsRemaining < 1 {
		return &remaining
	}

	monthly := remaining.Div(decimal.NewFromFloat(monthsRemaining))
	return &monthly
}
