package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoz65/finext/internal/database"
	"github.com/gustavoz65/finext/internal/model"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

var (
	ErrGoalNotFound         = errors.New("goal not found")
	ErrContributionNotFound = errors.New("contribution not found")
)

type GoalRepository struct {
	*BaseRepository
}

func NewGoalRepository(db *database.Database, logger *zerolog.Logger) *GoalRepository {
	return &GoalRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new goal
func (r *GoalRepository) Create(ctx context.Context, goal *model.Goal) error {
	goal.ID = uuid.New()
	goal.CurrentAmount = decimal.Zero
	goal.Status = model.GoalStatusInProgress
	goal.CreatedAt = time.Now()
	goal.UpdatedAt = time.Now()

	query := `
		INSERT INTO goals (
			id, user_id, name, description, target_amount, current_amount,
			target_date, icon, color, priority, status, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var targetDate sql.NullTime
	if goal.TargetDate != nil {
		targetDate = sql.NullTime{Time: *goal.TargetDate, Valid: true}
	}

	_, err := r.ExecContext(ctx, query,
		goal.ID.String(),
		goal.UserID.String(),
		goal.Name,
		NullString(goal.Description),
		goal.TargetAmount.String(),
		goal.CurrentAmount.String(),
		targetDate,
		goal.Icon,
		goal.Color,
		goal.Priority,
		goal.Status,
		goal.CreatedAt,
		goal.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create goal: %w", err)
	}

	return nil
}

// GetByID retrieves a goal by ID
func (r *GoalRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Goal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			target_date, icon, color, priority, status, completed_at,
			created_at, updated_at
		FROM goals
		WHERE id = ?
	`

	return r.scanGoal(r.QueryRowContext(ctx, query, id.String()))
}

// GetByIDAndUser retrieves a goal by ID ensuring it belongs to the user
func (r *GoalRepository) GetByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*model.Goal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			target_date, icon, color, priority, status, completed_at,
			created_at, updated_at
		FROM goals
		WHERE id = ? AND user_id = ?
	`

	return r.scanGoal(r.QueryRowContext(ctx, query, id.String(), userID.String()))
}

// GetAllByUser retrieves all goals for a user
func (r *GoalRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*model.Goal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			target_date, icon, color, priority, status, completed_at,
			created_at, updated_at
		FROM goals
		WHERE user_id = ?
		ORDER BY priority ASC, created_at DESC
	`

	rows, err := r.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get goals: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// GetActiveByUser retrieves active (in progress) goals for a user
func (r *GoalRepository) GetActiveByUser(ctx context.Context, userID uuid.UUID) ([]*model.Goal, error) {
	query := `
		SELECT id, user_id, name, description, target_amount, current_amount,
			target_date, icon, color, priority, status, completed_at,
			created_at, updated_at
		FROM goals
		WHERE user_id = ? AND status = 'in_progress'
		ORDER BY priority ASC, created_at DESC
	`

	rows, err := r.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get active goals: %w", err)
	}
	defer rows.Close()

	return r.scanGoals(rows)
}

// Update updates a goal
func (r *GoalRepository) Update(ctx context.Context, goal *model.Goal) error {
	goal.UpdatedAt = time.Now()

	query := `
		UPDATE goals SET
			name = ?, description = ?, target_amount = ?, target_date = ?,
			icon = ?, color = ?, priority = ?, status = ?, completed_at = ?,
			updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	var targetDate, completedAt sql.NullTime
	if goal.TargetDate != nil {
		targetDate = sql.NullTime{Time: *goal.TargetDate, Valid: true}
	}
	if goal.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *goal.CompletedAt, Valid: true}
	}

	result, err := r.ExecContext(ctx, query,
		goal.Name,
		NullString(goal.Description),
		goal.TargetAmount.String(),
		targetDate,
		goal.Icon,
		goal.Color,
		goal.Priority,
		goal.Status,
		completedAt,
		goal.UpdatedAt,
		goal.ID.String(),
		goal.UserID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update goal: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrGoalNotFound
	}

	return nil
}

// UpdateCurrentAmount updates the current amount of a goal
func (r *GoalRepository) UpdateCurrentAmount(ctx context.Context, id uuid.UUID, amount decimal.Decimal) error {
	query := `UPDATE goals SET current_amount = ?, updated_at = ? WHERE id = ?`

	_, err := r.ExecContext(ctx, query, amount.String(), time.Now(), id.String())
	if err != nil {
		return fmt.Errorf("failed to update current amount: %w", err)
	}

	return nil
}

// MarkCompleted marks a goal as completed
func (r *GoalRepository) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE goals
		SET status = 'completed', completed_at = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now()
	_, err := r.ExecContext(ctx, query, now, now, id.String())
	if err != nil {
		return fmt.Errorf("failed to mark goal completed: %w", err)
	}

	return nil
}

// Delete deletes a goal (hard delete)
func (r *GoalRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM goals WHERE id = ? AND user_id = ?`

	result, err := r.ExecContext(ctx, query, id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("failed to delete goal: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrGoalNotFound
	}

	return nil
}

// ========================================
// Goal Contribution Methods
// ========================================

// AddContribution adds a contribution to a goal
func (r *GoalRepository) AddContribution(ctx context.Context, contribution *model.GoalContribution) error {
	contribution.ID = uuid.New()
	contribution.CreatedAt = time.Now()

	query := `
		INSERT INTO goal_contributions (id, goal_id, amount, note, contribution_date, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.ExecContext(ctx, query,
		contribution.ID.String(),
		contribution.GoalID.String(),
		contribution.Amount.String(),
		NullString(contribution.Note),
		contribution.ContributionDate,
		contribution.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to add contribution: %w", err)
	}

	return nil
}

// GetContributionsByGoal retrieves all contributions for a goal
func (r *GoalRepository) GetContributionsByGoal(ctx context.Context, goalID uuid.UUID) ([]model.GoalContribution, error) {
	query := `
		SELECT id, goal_id, amount, note, contribution_date, created_at
		FROM goal_contributions
		WHERE goal_id = ?
		ORDER BY contribution_date DESC
	`

	rows, err := r.QueryContext(ctx, query, goalID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get contributions: %w", err)
	}
	defer rows.Close()

	var contributions []model.GoalContribution

	for rows.Next() {
		var c model.GoalContribution
		var note sql.NullString
		var amount string

		err := rows.Scan(
			&c.ID,
			&c.GoalID,
			&amount,
			&note,
			&c.ContributionDate,
			&c.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan contribution: %w", err)
		}

		c.Amount, _ = decimal.NewFromString(amount)
		c.Note = StringPtr(note)
		contributions = append(contributions, c)
	}

	return contributions, nil
}

// DeleteContribution deletes a contribution
func (r *GoalRepository) DeleteContribution(ctx context.Context, contributionID, goalID uuid.UUID) error {
	query := `DELETE FROM goal_contributions WHERE id = ? AND goal_id = ?`

	result, err := r.ExecContext(ctx, query, contributionID.String(), goalID.String())
	if err != nil {
		return fmt.Errorf("failed to delete contribution: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrContributionNotFound
	}

	return nil
}

// GetGoalSummary returns a summary of all goals for a user
func (r *GoalRepository) GetGoalSummary(ctx context.Context, userID uuid.UUID) (*model.GoalSummary, error) {
	query := `
		SELECT
			COUNT(*) as total_goals,
			COALESCE(SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END), 0) as completed_goals,
			COALESCE(SUM(CASE WHEN status = 'in_progress' THEN 1 ELSE 0 END), 0) as in_progress_goals,
			COALESCE(SUM(target_amount), 0) as total_target,
			COALESCE(SUM(current_amount), 0) as total_saved
		FROM goals
		WHERE user_id = ?
	`

	var summary model.GoalSummary
	var totalTarget, totalSaved string

	err := r.QueryRowContext(ctx, query, userID.String()).Scan(
		&summary.TotalGoals,
		&summary.CompletedGoals,
		&summary.InProgressGoals,
		&totalTarget,
		&totalSaved,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get goal summary: %w", err)
	}

	summary.TotalTargetAmount, _ = decimal.NewFromString(totalTarget)
	summary.TotalSavedAmount, _ = decimal.NewFromString(totalSaved)

	if !summary.TotalTargetAmount.IsZero() {
		summary.OverallProgress = summary.TotalSavedAmount.Div(summary.TotalTargetAmount).Mul(decimal.NewFromInt(100))
	}

	return &summary, nil
}

// Helper functions

func (r *GoalRepository) scanGoal(row *sql.Row) (*model.Goal, error) {
	goal := &model.Goal{}
	var description sql.NullString
	var targetDate, completedAt sql.NullTime
	var targetAmount, currentAmount string

	err := row.Scan(
		&goal.ID,
		&goal.UserID,
		&goal.Name,
		&description,
		&targetAmount,
		&currentAmount,
		&targetDate,
		&goal.Icon,
		&goal.Color,
		&goal.Priority,
		&goal.Status,
		&completedAt,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGoalNotFound
		}
		return nil, fmt.Errorf("failed to scan goal: %w", err)
	}

	goal.Description = StringPtr(description)
	goal.TargetAmount, _ = decimal.NewFromString(targetAmount)
	goal.CurrentAmount, _ = decimal.NewFromString(currentAmount)

	if targetDate.Valid {
		goal.TargetDate = &targetDate.Time
	}
	if completedAt.Valid {
		goal.CompletedAt = &completedAt.Time
	}

	return goal, nil
}

func (r *GoalRepository) scanGoals(rows *sql.Rows) ([]*model.Goal, error) {
	var goals []*model.Goal

	for rows.Next() {
		goal := &model.Goal{}
		var description sql.NullString
		var targetDate, completedAt sql.NullTime
		var targetAmount, currentAmount string

		err := rows.Scan(
			&goal.ID,
			&goal.UserID,
			&goal.Name,
			&description,
			&targetAmount,
			&currentAmount,
			&targetDate,
			&goal.Icon,
			&goal.Color,
			&goal.Priority,
			&goal.Status,
			&completedAt,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan goal: %w", err)
		}

		goal.Description = StringPtr(description)
		goal.TargetAmount, _ = decimal.NewFromString(targetAmount)
		goal.CurrentAmount, _ = decimal.NewFromString(currentAmount)

		if targetDate.Valid {
			goal.TargetDate = &targetDate.Time
		}
		if completedAt.Valid {
			goal.CompletedAt = &completedAt.Time
		}

		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating goals: %w", err)
	}

	return goals, nil
}
