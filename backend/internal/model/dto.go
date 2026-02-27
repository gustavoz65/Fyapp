package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ========================================
// Auth DTOs
// ========================================

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email,max=255"`
	Password  string `json:"password" validate:"required,min=8,max=72"`
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Phone     string `json:"phone,omitempty" validate:"omitempty,max=20"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
}

type SetPasswordRequest struct {
	NewPassword     string `json:"new_password" validate:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

// Social Login/Register DTOs
type SocialLoginRequest struct {
	Provider   string `json:"provider" validate:"required,oneof=google facebook github"`
	IDToken    string `json:"id_token" validate:"required"`
	DeviceInfo string `json:"device_info,omitempty" validate:"omitempty,max=500"`
}

type LinkProviderRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google facebook github"`
	IDToken  string `json:"id_token" validate:"required"`
}

type UnlinkProviderRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google facebook github"`
}

type LinkedProviderResponse struct {
	Provider  string    `json:"provider"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	IsPrimary bool      `json:"is_primary"`
	LinkedAt  time.Time `json:"linked_at"`
}

type ListProvidersResponse struct {
	Providers   []LinkedProviderResponse `json:"providers"`
	HasPassword bool                     `json:"has_password"`
}

// ========================================
// User DTOs
// ========================================

type UpdateUserRequest struct {
	FirstName         *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	LastName          *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone             *string `json:"phone,omitempty" validate:"omitempty,max=20"`
	PreferredCurrency *string `json:"preferred_currency,omitempty" validate:"omitempty,len=3"`
	PreferredLanguage *string `json:"preferred_language,omitempty" validate:"omitempty,max=5"`
	Timezone          *string `json:"timezone,omitempty" validate:"omitempty,max=50"`
}

type UpdateUserSettingsRequest struct {
	NotificationEmail       *bool    `json:"notification_email,omitempty"`
	NotificationPush        *bool    `json:"notification_push,omitempty"`
	NotificationSMS         *bool    `json:"notification_sms,omitempty"`
	BudgetAlerts            *bool    `json:"budget_alerts,omitempty"`
	BillReminders           *bool    `json:"bill_reminders,omitempty"`
	BillReminderDays        *int     `json:"bill_reminder_days,omitempty" validate:"omitempty,min=1,max=30"`
	WeeklySummary           *bool    `json:"weekly_summary,omitempty"`
	MonthlyReport           *bool    `json:"monthly_report,omitempty"`
	LowBalanceAlert         *bool    `json:"low_balance_alert,omitempty"`
	LowBalanceThreshold     *float64 `json:"low_balance_threshold,omitempty" validate:"omitempty,gte=0"`
	AllowManualTransactions *bool    `json:"allow_manual_transactions,omitempty"`
	Theme                   *string  `json:"theme,omitempty" validate:"omitempty,oneof=light dark system"`
}

// ========================================
// Category DTOs
// ========================================

type CreateCategoryRequest struct {
	Name        string       `json:"name" validate:"required,min=2,max=100"`
	Description string       `json:"description,omitempty" validate:"omitempty,max=255"`
	Type        CategoryType `json:"type" validate:"required,oneof=income expense"`
	Color       string       `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon        string       `json:"icon,omitempty" validate:"omitempty,max=50"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=255"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon        *string `json:"icon,omitempty" validate:"omitempty,max=50"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// ========================================
// Bank Account DTOs
// ========================================

type CreateBankAccountRequest struct {
	Name           string           `json:"name" validate:"required,min=2,max=100"`
	BankName       string           `json:"bank_name,omitempty" validate:"omitempty,max=100"`
	BankCode       string           `json:"bank_code,omitempty" validate:"omitempty,max=10"`
	AccountType    AccountType      `json:"account_type" validate:"required,oneof=checking savings credit_card investment cash other"`
	AccountNumber  string           `json:"account_number,omitempty" validate:"omitempty,max=50"`
	Agency         string           `json:"agency,omitempty" validate:"omitempty,max=20"`
	InitialBalance decimal.Decimal  `json:"initial_balance" validate:"required"`
	CreditLimit    *decimal.Decimal `json:"credit_limit,omitempty"`
	ClosingDay     *int             `json:"closing_day,omitempty" validate:"omitempty,min=1,max=31"`
	DueDay         *int             `json:"due_day,omitempty" validate:"omitempty,min=1,max=31"`
	Currency       string           `json:"currency,omitempty" validate:"omitempty,len=3"`
	Color          string           `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon           string           `json:"icon,omitempty" validate:"omitempty,max=50"`
	IncludeInTotal *bool            `json:"include_in_total,omitempty"`
}

type UpdateBankAccountRequest struct {
	Name           *string          `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	BankName       *string          `json:"bank_name,omitempty" validate:"omitempty,max=100"`
	BankCode       *string          `json:"bank_code,omitempty" validate:"omitempty,max=10"`
	AccountNumber  *string          `json:"account_number,omitempty" validate:"omitempty,max=50"`
	Agency         *string          `json:"agency,omitempty" validate:"omitempty,max=20"`
	CreditLimit    *decimal.Decimal `json:"credit_limit,omitempty"`
	ClosingDay     *int             `json:"closing_day,omitempty" validate:"omitempty,min=1,max=31"`
	DueDay         *int             `json:"due_day,omitempty" validate:"omitempty,min=1,max=31"`
	Color          *string          `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Icon           *string          `json:"icon,omitempty" validate:"omitempty,max=50"`
	IsActive       *bool            `json:"is_active,omitempty"`
	IncludeInTotal *bool            `json:"include_in_total,omitempty"`
}

// ========================================
// Transaction DTOs
// ========================================

type CreateTransactionRequest struct {
	BankAccountID     uuid.UUID       `json:"bank_account_id" validate:"required"`
	CategoryID        *uuid.UUID      `json:"category_id,omitempty"`
	Type              TransactionType `json:"type" validate:"required,oneof=income expense"`
	Amount            decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Description       string          `json:"description" validate:"required,min=2,max=255"`
	Notes             string          `json:"notes,omitempty" validate:"omitempty,max=1000"`
	TransactionDate   time.Time       `json:"transaction_date" validate:"required"`
	DueDate           *time.Time      `json:"due_date,omitempty"`
	IsPaid            bool            `json:"is_paid"`
	AutoPay           bool            `json:"auto_pay"`
	Tags              []string        `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
	TotalInstallments *int            `json:"total_installments,omitempty" validate:"omitempty,min=2,max=120"`
}

type UpdateTransactionRequest struct {
	CategoryID      *uuid.UUID       `json:"category_id,omitempty"`
	Amount          *decimal.Decimal `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Description     *string          `json:"description,omitempty" validate:"omitempty,min=2,max=255"`
	Notes           *string          `json:"notes,omitempty" validate:"omitempty,max=1000"`
	TransactionDate *time.Time       `json:"transaction_date,omitempty"`
	DueDate         *time.Time       `json:"due_date,omitempty"`
	IsPaid          *bool            `json:"is_paid,omitempty"`
	Tags            []string         `json:"tags,omitempty" validate:"omitempty,dive,max=50"`
}

type TransactionFilter struct {
	UserID        uuid.UUID          `json:"-"`
	AccountID     *uuid.UUID         `json:"account_id,omitempty"`
	CategoryID    *uuid.UUID         `json:"category_id,omitempty"`
	Type          *TransactionType   `json:"type,omitempty"`
	Source        *TransactionSource `json:"source,omitempty"`
	StartDate     *time.Time         `json:"start_date,omitempty"`
	EndDate       *time.Time         `json:"end_date,omitempty"`
	IsPaid        *bool              `json:"is_paid,omitempty"`
	IsRecurring   *bool              `json:"is_recurring,omitempty"`
	MinAmount     *decimal.Decimal   `json:"min_amount,omitempty"`
	MaxAmount     *decimal.Decimal   `json:"max_amount,omitempty"`
	SearchTerm    string             `json:"search_term,omitempty"`
	Tags          []string           `json:"tags,omitempty"`
	Page          int                `json:"page,omitempty"`
	PageSize      int                `json:"page_size,omitempty"`
	SortBy        string             `json:"sort_by,omitempty"`
	SortDirection string             `json:"sort_direction,omitempty"`
}

// ========================================
// Transfer DTOs
// ========================================

type CreateTransferRequest struct {
	FromAccountID uuid.UUID       `json:"from_account_id" validate:"required"`
	ToAccountID   uuid.UUID       `json:"to_account_id" validate:"required,nefield=FromAccountID"`
	Amount        decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Description   string          `json:"description,omitempty" validate:"omitempty,max=255"`
	TransferDate  time.Time       `json:"transfer_date" validate:"required"`
}

// ========================================
// Recurring Transaction DTOs
// ========================================

type CreateRecurringTransactionRequest struct {
	BankAccountID uuid.UUID          `json:"bank_account_id" validate:"required"`
	CategoryID    *uuid.UUID         `json:"category_id,omitempty"`
	Type          TransactionType    `json:"type" validate:"required,oneof=income expense"`
	Amount        decimal.Decimal    `json:"amount" validate:"required,gt=0"`
	Description   string             `json:"description" validate:"required,min=2,max=255"`
	Frequency     RecurringFrequency `json:"frequency" validate:"required,oneof=daily weekly biweekly monthly quarterly yearly"`
	DayOfMonth    *int               `json:"day_of_month,omitempty" validate:"omitempty,min=1,max=31"`
	DayOfWeek     *int               `json:"day_of_week,omitempty" validate:"omitempty,min=0,max=6"`
	StartDate     time.Time          `json:"start_date" validate:"required"`
	EndDate       *time.Time         `json:"end_date,omitempty"`
	AutoConfirm   bool               `json:"auto_confirm"`
}

type UpdateRecurringTransactionRequest struct {
	CategoryID  *uuid.UUID          `json:"category_id,omitempty"`
	Amount      *decimal.Decimal    `json:"amount,omitempty" validate:"omitempty,gt=0"`
	Description *string             `json:"description,omitempty" validate:"omitempty,min=2,max=255"`
	Frequency   *RecurringFrequency `json:"frequency,omitempty" validate:"omitempty,oneof=daily weekly biweekly monthly quarterly yearly"`
	DayOfMonth  *int                `json:"day_of_month,omitempty" validate:"omitempty,min=1,max=31"`
	DayOfWeek   *int                `json:"day_of_week,omitempty" validate:"omitempty,min=0,max=6"`
	EndDate     *time.Time          `json:"end_date,omitempty"`
	IsActive    *bool               `json:"is_active,omitempty"`
	AutoConfirm *bool               `json:"auto_confirm,omitempty"`
}

// ========================================
// Budget DTOs
// ========================================

type CreateBudgetRequest struct {
	CategoryID     *uuid.UUID       `json:"category_id,omitempty"`
	Name           string           `json:"name" validate:"required,min=2,max=100"`
	Amount         decimal.Decimal  `json:"amount" validate:"required,gt=0"`
	PeriodType     BudgetPeriodType `json:"period_type" validate:"required,oneof=monthly quarterly yearly custom"`
	StartDate      time.Time        `json:"start_date" validate:"required"`
	EndDate        time.Time        `json:"end_date" validate:"required,gtfield=StartDate"`
	AlertThreshold *decimal.Decimal `json:"alert_threshold,omitempty" validate:"omitempty,gt=0,lte=100"`
}

type UpdateBudgetRequest struct {
	Name           *string          `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Amount         *decimal.Decimal `json:"amount,omitempty" validate:"omitempty,gt=0"`
	AlertThreshold *decimal.Decimal `json:"alert_threshold,omitempty" validate:"omitempty,gt=0,lte=100"`
	IsActive       *bool            `json:"is_active,omitempty"`
}

// ========================================
// Goal DTOs
// ========================================

type CreateGoalRequest struct {
	Name         string          `json:"name" validate:"required,min=2,max=100"`
	Description  string          `json:"description,omitempty" validate:"omitempty,max=500"`
	TargetAmount decimal.Decimal `json:"target_amount" validate:"required,gt=0"`
	TargetDate   *time.Time      `json:"target_date,omitempty"`
	Icon         string          `json:"icon,omitempty" validate:"omitempty,max=50"`
	Color        string          `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Priority     *int            `json:"priority,omitempty" validate:"omitempty,min=1,max=5"`
}

type UpdateGoalRequest struct {
	Name         *string          `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description  *string          `json:"description,omitempty" validate:"omitempty,max=500"`
	TargetAmount *decimal.Decimal `json:"target_amount,omitempty" validate:"omitempty,gt=0"`
	TargetDate   *time.Time       `json:"target_date,omitempty"`
	Icon         *string          `json:"icon,omitempty" validate:"omitempty,max=50"`
	Color        *string          `json:"color,omitempty" validate:"omitempty,hexcolor"`
	Priority     *int             `json:"priority,omitempty" validate:"omitempty,min=1,max=5"`
	Status       *GoalStatus      `json:"status,omitempty" validate:"omitempty,oneof=in_progress completed cancelled"`
}

type CreateGoalContributionRequest struct {
	Amount           decimal.Decimal `json:"amount" validate:"required,gt=0"`
	Note             string          `json:"note,omitempty" validate:"omitempty,max=255"`
	ContributionDate *time.Time      `json:"contribution_date,omitempty"`
}

// ========================================
// Report DTOs
// ========================================

type GenerateReportRequest struct {
	Type          ReportType  `json:"type" validate:"required,oneof=cash_flow expense_by_category income_vs_expense budget_analysis custom"`
	StartDate     time.Time   `json:"start_date" validate:"required"`
	EndDate       time.Time   `json:"end_date" validate:"required,gtfield=StartDate"`
	AccountIDs    []uuid.UUID `json:"account_ids,omitempty"`
	CategoryIDs   []uuid.UUID `json:"category_ids,omitempty"`
	Format        string      `json:"format,omitempty" validate:"omitempty,oneof=json pdf excel"`
	IncludeCharts bool        `json:"include_charts"`
}

type CreateScheduledReportRequest struct {
	Name       string         `json:"name" validate:"required,min=2,max=100"`
	Type       ReportType     `json:"type" validate:"required,oneof=cash_flow expense_by_category income_vs_expense budget_analysis"`
	Schedule   ReportSchedule `json:"schedule" validate:"required,oneof=daily weekly monthly"`
	Parameters string         `json:"parameters" validate:"required"`
}

// ========================================
// Category Suggestion DTOs
// ========================================

type SuggestCategoryRequest struct {
	Description string `json:"description" validate:"required,min=1,max=255"`
}

type SuggestCategoryResponse struct {
	Suggestions []CategorySuggestion `json:"suggestions"`
}

// ========================================
// Pagination & Response DTOs
// ========================================

type PaginationParams struct {
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"page_size" validate:"min=1,max=100"`
	SortBy   string `json:"sort_by,omitempty"`
	SortDir  string `json:"sort_dir,omitempty" validate:"omitempty,oneof=asc desc"`
}

type PaginatedResponse[T any] struct {
	Data       []T   `json:"data"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
	HasMore    bool  `json:"has_more"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}
