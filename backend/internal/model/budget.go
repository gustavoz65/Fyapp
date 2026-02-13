package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BudgetPeriodType string

const (
	BudgetPeriodMonthly   BudgetPeriodType = "monthly"
	BudgetPeriodQuarterly BudgetPeriodType = "quarterly"
	BudgetPeriodYearly    BudgetPeriodType = "yearly"
	BudgetPeriodCustom    BudgetPeriodType = "custom"
)

type Budget struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	UserID         uuid.UUID        `json:"user_id" db:"user_id"`
	CategoryID     *uuid.UUID       `json:"category_id,omitempty" db:"category_id"`
	Name           string           `json:"name" db:"name"`
	Amount         decimal.Decimal  `json:"amount" db:"amount"`
	SpentAmount    decimal.Decimal  `json:"spent_amount" db:"spent_amount"`
	PeriodType     BudgetPeriodType `json:"period_type" db:"period_type"`
	StartDate      time.Time        `json:"start_date" db:"start_date"`
	EndDate        time.Time        `json:"end_date" db:"end_date"`
	AlertThreshold decimal.Decimal  `json:"alert_threshold" db:"alert_threshold"`
	AlertSent      bool             `json:"alert_sent" db:"alert_sent"`
	IsActive       bool             `json:"is_active" db:"is_active"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`

	Category *Category `json:"category,omitempty" db:"-"`
}

func (b *Budget) RemainingAmount() decimal.Decimal {
	return b.Amount.Sub(b.SpentAmount)
}

func (b *Budget) UsedPercentage() decimal.Decimal {
	if b.Amount.IsZero() {
		return decimal.Zero
	}
	return b.SpentAmount.Div(b.Amount).Mul(decimal.NewFromInt(100))
}

func (b *Budget) IsOverBudget() bool {
	return b.SpentAmount.GreaterThan(b.Amount)
}

func (b *Budget) ShouldAlert() bool {
	if b.AlertSent {
		return false
	}
	threshold := b.Amount.Mul(b.AlertThreshold).Div(decimal.NewFromInt(100))
	return b.SpentAmount.GreaterThanOrEqual(threshold)
}

func (b *Budget) DaysRemaining() int {
	if time.Now().After(b.EndDate) {
		return 0
	}
	return int(time.Until(b.EndDate).Hours() / 24)
}

func (b *Budget) DailyBudget() decimal.Decimal {
	days := b.DaysRemaining()
	if days <= 0 {
		return decimal.Zero
	}
	return b.RemainingAmount().Div(decimal.NewFromInt(int64(days)))
}

func (b *Budget) IsCurrentPeriod() bool {
	now := time.Now()
	return (now.Equal(b.StartDate) || now.After(b.StartDate)) &&
		(now.Equal(b.EndDate) || now.Before(b.EndDate))
}

type BudgetSummary struct {
	Budget          *Budget         `json:"budget"`
	TotalBudgeted   decimal.Decimal `json:"total_budgeted"`
	TotalSpent      decimal.Decimal `json:"total_spent"`
	TotalRemaining  decimal.Decimal `json:"total_remaining"`
	UsedPercentage  decimal.Decimal `json:"used_percentage"`
	TransactionCount int            `json:"transaction_count"`
}
