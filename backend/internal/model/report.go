package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ReportType string

const (
	ReportTypeCashFlow          ReportType = "cash_flow"
	ReportTypeExpenseByCategory ReportType = "expense_by_category"
	ReportTypeIncomeVsExpense   ReportType = "income_vs_expense"
	ReportTypeBudgetAnalysis    ReportType = "budget_analysis"
	ReportTypeCustom            ReportType = "custom"
)

type ReportSchedule string

const (
	ReportScheduleDaily   ReportSchedule = "daily"
	ReportScheduleWeekly  ReportSchedule = "weekly"
	ReportScheduleMonthly ReportSchedule = "monthly"
)

type Report struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          uuid.UUID       `json:"user_id" db:"user_id"`
	Name            string          `json:"name" db:"name"`
	Type            ReportType      `json:"type" db:"type"`
	Parameters      string          `json:"parameters" db:"parameters"`
	Schedule        *ReportSchedule `json:"schedule,omitempty" db:"schedule"`
	LastGeneratedAt *time.Time      `json:"last_generated_at,omitempty" db:"last_generated_at"`
	FileURL         *string         `json:"file_url,omitempty" db:"file_url"`
	IsActive        bool            `json:"is_active" db:"is_active"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

type ReportParameters struct {
	StartDate    time.Time   `json:"start_date"`
	EndDate      time.Time   `json:"end_date"`
	AccountIDs   []uuid.UUID `json:"account_ids,omitempty"`
	CategoryIDs  []uuid.UUID `json:"category_ids,omitempty"`
	GroupBy      string      `json:"group_by,omitempty"`
	IncludeCharts bool       `json:"include_charts"`
}

type CashFlowReport struct {
	Period          string          `json:"period"`
	StartDate       time.Time       `json:"start_date"`
	EndDate         time.Time       `json:"end_date"`
	TotalIncome     decimal.Decimal `json:"total_income"`
	TotalExpense    decimal.Decimal `json:"total_expense"`
	NetCashFlow     decimal.Decimal `json:"net_cash_flow"`
	OpeningBalance  decimal.Decimal `json:"opening_balance"`
	ClosingBalance  decimal.Decimal `json:"closing_balance"`
	DailyBreakdown  []DailyCashFlow `json:"daily_breakdown,omitempty"`
	IncomeByCategory  []CategoryAmount `json:"income_by_category"`
	ExpenseByCategory []CategoryAmount `json:"expense_by_category"`
}

type DailyCashFlow struct {
	Date     time.Time       `json:"date"`
	Income   decimal.Decimal `json:"income"`
	Expense  decimal.Decimal `json:"expense"`
	NetFlow  decimal.Decimal `json:"net_flow"`
	Balance  decimal.Decimal `json:"balance"`
}

type CategoryAmount struct {
	CategoryID   *uuid.UUID      `json:"category_id,omitempty"`
	CategoryName string          `json:"category_name"`
	Amount       decimal.Decimal `json:"amount"`
	Percentage   decimal.Decimal `json:"percentage"`
	Count        int             `json:"count"`
}

type ExpenseByCategoryReport struct {
	StartDate        time.Time        `json:"start_date"`
	EndDate          time.Time        `json:"end_date"`
	TotalExpense     decimal.Decimal  `json:"total_expense"`
	Categories       []CategoryAmount `json:"categories"`
	AveragePerDay    decimal.Decimal  `json:"average_per_day"`
	LargestExpense   *Transaction     `json:"largest_expense,omitempty"`
	MostFrequentCategory string       `json:"most_frequent_category"`
}

type IncomeVsExpenseReport struct {
	StartDate      time.Time              `json:"start_date"`
	EndDate        time.Time              `json:"end_date"`
	TotalIncome    decimal.Decimal        `json:"total_income"`
	TotalExpense   decimal.Decimal        `json:"total_expense"`
	NetIncome      decimal.Decimal        `json:"net_income"`
	SavingsRate    decimal.Decimal        `json:"savings_rate"`
	MonthlyData    []MonthlyIncomeExpense `json:"monthly_data,omitempty"`
	IncomeGrowth   decimal.Decimal        `json:"income_growth"`
	ExpenseGrowth  decimal.Decimal        `json:"expense_growth"`
}

type MonthlyIncomeExpense struct {
	Month      string          `json:"month"`
	Year       int             `json:"year"`
	Income     decimal.Decimal `json:"income"`
	Expense    decimal.Decimal `json:"expense"`
	NetIncome  decimal.Decimal `json:"net_income"`
	SavingsRate decimal.Decimal `json:"savings_rate"`
}

type BudgetAnalysisReport struct {
	StartDate       time.Time            `json:"start_date"`
	EndDate         time.Time            `json:"end_date"`
	TotalBudgeted   decimal.Decimal      `json:"total_budgeted"`
	TotalSpent      decimal.Decimal      `json:"total_spent"`
	TotalRemaining  decimal.Decimal      `json:"total_remaining"`
	OverallUsage    decimal.Decimal      `json:"overall_usage"`
	BudgetDetails   []BudgetDetail       `json:"budget_details"`
	OverBudgetCount int                  `json:"over_budget_count"`
	UnderBudgetCount int                 `json:"under_budget_count"`
}

type BudgetDetail struct {
	BudgetID       uuid.UUID       `json:"budget_id"`
	BudgetName     string          `json:"budget_name"`
	CategoryName   *string         `json:"category_name,omitempty"`
	Amount         decimal.Decimal `json:"amount"`
	SpentAmount    decimal.Decimal `json:"spent_amount"`
	RemainingAmount decimal.Decimal `json:"remaining_amount"`
	UsagePercentage decimal.Decimal `json:"usage_percentage"`
	Status         string          `json:"status"`
}

type DashboardSummary struct {
	// Balance info
	TotalBalance     decimal.Decimal `json:"total_balance"`
	TotalInAccounts  int             `json:"total_accounts"`

	// Period summary (current month)
	MonthIncome      decimal.Decimal `json:"month_income"`
	MonthExpense     decimal.Decimal `json:"month_expense"`
	MonthBalance     decimal.Decimal `json:"month_balance"`

	// Comparison with last month
	IncomeChange     decimal.Decimal `json:"income_change"`
	ExpenseChange    decimal.Decimal `json:"expense_change"`

	// Budget status
	ActiveBudgets    int             `json:"active_budgets"`
	OverBudgetCount  int             `json:"over_budget_count"`
	BudgetUsage      decimal.Decimal `json:"budget_usage"`

	// Goals
	ActiveGoals      int             `json:"active_goals"`
	GoalsProgress    decimal.Decimal `json:"goals_progress"`

	// Upcoming bills
	UpcomingBills    int             `json:"upcoming_bills"`
	UpcomingTotal    decimal.Decimal `json:"upcoming_total"`

	// Recent transactions
	RecentTransactions []Transaction `json:"recent_transactions,omitempty"`

	// Top expense categories
	TopCategories    []CategoryAmount `json:"top_categories,omitempty"`
}

type AuditLog struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Action     string     `json:"action" db:"action"`
	EntityType string     `json:"entity_type" db:"entity_type"`
	EntityID   *uuid.UUID `json:"entity_id,omitempty" db:"entity_id"`
	OldValues  *string    `json:"old_values,omitempty" db:"old_values"`
	NewValues  *string    `json:"new_values,omitempty" db:"new_values"`
	IPAddress  *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent  *string    `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}
