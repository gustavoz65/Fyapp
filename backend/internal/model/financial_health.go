package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type FinancialHealthSnapshot struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	UserID            uuid.UUID       `json:"user_id" db:"user_id"`
	Period            time.Time       `json:"period" db:"period"` // YYYY-MM-01
	Score             decimal.Decimal `json:"score" db:"score"`   // 0.00-5.00
	ScorePercentage   decimal.Decimal `json:"score_percentage" db:"score_percentage"` // 0.00-100.00
	EconomyRate       decimal.Decimal `json:"economy_rate" db:"economy_rate"`
	BudgetCompliance  decimal.Decimal `json:"budget_compliance" db:"budget_compliance"`
	GoalsProgress     decimal.Decimal `json:"goals_progress" db:"goals_progress"`
	SpendingReduction decimal.Decimal `json:"spending_reduction" db:"spending_reduction"`
	Consistency       decimal.Decimal `json:"consistency" db:"consistency"`
	BreakdownJSON     *string         `json:"breakdown_json,omitempty" db:"breakdown_json"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
}

type HealthScoreBreakdown struct {
	Income                decimal.Decimal            `json:"income"`
	Expense               decimal.Decimal            `json:"expense"`
	Balance               decimal.Decimal            `json:"balance"`
	EconomyPercentage     float64                    `json:"economy_percentage"`
	Budgets               []BudgetHealthDetail       `json:"budgets"`
	Goals                 []GoalHealthDetail         `json:"goals"`
	TopCategories         []CategorySpendingDetail   `json:"top_categories"`
	ConsistencyDetails    ConsistencyDetail          `json:"consistency_details"`
}

type BudgetHealthDetail struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Limit      decimal.Decimal `json:"limit"`
	Spent      decimal.Decimal `json:"spent"`
	Usage      float64         `json:"usage"` // percentage
	Score      float64         `json:"score"` // 0-5
}

type GoalHealthDetail struct {
	ID              uuid.UUID       `json:"id"`
	Name            string          `json:"name"`
	TargetAmount    decimal.Decimal `json:"target_amount"`
	CurrentAmount   decimal.Decimal `json:"current_amount"`
	Progress        float64         `json:"progress"` // percentage
	Score           float64         `json:"score"`    // 0-5
}

type CategorySpendingDetail struct {
	CategoryID       uuid.UUID       `json:"category_id"`
	CategoryName     string          `json:"category_name"`
	CurrentAmount    decimal.Decimal `json:"current_amount"`
	PreviousAmount   decimal.Decimal `json:"previous_amount"`
	ChangePercentage float64         `json:"change_percentage"`
	Score            float64         `json:"score"` // 0-5
}

type ConsistencyDetail struct {
	NegativeMonthsCount int     `json:"negative_months_count"`
	ActiveDebtsCount    int     `json:"active_debts_count"`
	Score               float64 `json:"score"` // 0-5
}
