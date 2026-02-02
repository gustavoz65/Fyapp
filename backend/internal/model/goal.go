package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GoalStatus string

const (
	GoalStatusInProgress GoalStatus = "in_progress"
	GoalStatusCompleted  GoalStatus = "completed"
	GoalStatusCancelled  GoalStatus = "cancelled"
)

type Goal struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	UserID        uuid.UUID       `json:"user_id" db:"user_id"`
	Name          string          `json:"name" db:"name"`
	Description   *string         `json:"description,omitempty" db:"description"`
	TargetAmount  decimal.Decimal `json:"target_amount" db:"target_amount"`
	CurrentAmount decimal.Decimal `json:"current_amount" db:"current_amount"`
	TargetDate    *time.Time      `json:"target_date,omitempty" db:"target_date"`
	Icon          string          `json:"icon" db:"icon"`
	Color         string          `json:"color" db:"color"`
	Priority      int             `json:"priority" db:"priority"`
	Status        GoalStatus      `json:"status" db:"status"`
	CompletedAt   *time.Time      `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`

	// Computed fields
	Contributions []GoalContribution `json:"contributions,omitempty" db:"-"`
}

func (g *Goal) Progress() decimal.Decimal {
	if g.TargetAmount.IsZero() {
		return decimal.Zero
	}
	return g.CurrentAmount.Div(g.TargetAmount).Mul(decimal.NewFromInt(100))
}

func (g *Goal) RemainingAmount() decimal.Decimal {
	remaining := g.TargetAmount.Sub(g.CurrentAmount)
	if remaining.IsNegative() {
		return decimal.Zero
	}
	return remaining
}

func (g *Goal) IsCompleted() bool {
	return g.Status == GoalStatusCompleted || g.CurrentAmount.GreaterThanOrEqual(g.TargetAmount)
}

func (g *Goal) DaysRemaining() *int {
	if g.TargetDate == nil {
		return nil
	}
	if time.Now().After(*g.TargetDate) {
		zero := 0
		return &zero
	}
	days := int(time.Until(*g.TargetDate).Hours() / 24)
	return &days
}

func (g *Goal) RequiredMonthlyContribution() *decimal.Decimal {
	if g.TargetDate == nil || g.IsCompleted() {
		return nil
	}
	daysRemaining := g.DaysRemaining()
	if daysRemaining == nil || *daysRemaining <= 0 {
		return nil
	}
	monthsRemaining := decimal.NewFromFloat(float64(*daysRemaining) / 30.0)
	if monthsRemaining.IsZero() {
		return nil
	}
	monthly := g.RemainingAmount().Div(monthsRemaining)
	return &monthly
}

func (g *Goal) IsOnTrack() bool {
	if g.TargetDate == nil || g.IsCompleted() {
		return true
	}

	// Calculate expected progress based on time elapsed
	totalDuration := g.TargetDate.Sub(g.CreatedAt)
	elapsedDuration := time.Since(g.CreatedAt)

	if totalDuration <= 0 {
		return true
	}

	expectedProgress := decimal.NewFromFloat(float64(elapsedDuration) / float64(totalDuration) * 100)
	actualProgress := g.Progress()

	return actualProgress.GreaterThanOrEqual(expectedProgress.Sub(decimal.NewFromInt(10)))
}

type GoalContribution struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	GoalID           uuid.UUID       `json:"goal_id" db:"goal_id"`
	Amount           decimal.Decimal `json:"amount" db:"amount"`
	Note             *string         `json:"note,omitempty" db:"note"`
	ContributionDate time.Time       `json:"contribution_date" db:"contribution_date"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

type GoalSummary struct {
	TotalGoals        int             `json:"total_goals"`
	CompletedGoals    int             `json:"completed_goals"`
	InProgressGoals   int             `json:"in_progress_goals"`
	TotalTargetAmount decimal.Decimal `json:"total_target_amount"`
	TotalSavedAmount  decimal.Decimal `json:"total_saved_amount"`
	OverallProgress   decimal.Decimal `json:"overall_progress"`
}
