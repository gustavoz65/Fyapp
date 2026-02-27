package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

type TransactionSource string

const (
	TransactionSourceManual    TransactionSource = "manual"
	TransactionSourceBankSync  TransactionSource = "bank_sync"
	TransactionSourceRecurring TransactionSource = "recurring"
)

type Transaction struct {
	ID                 uuid.UUID         `json:"id" db:"id"`
	UserID             uuid.UUID         `json:"user_id" db:"user_id"`
	BankAccountID      uuid.UUID         `json:"bank_account_id" db:"bank_account_id"`
	CategoryID         *uuid.UUID        `json:"category_id,omitempty" db:"category_id"`
	Type               TransactionType   `json:"type" db:"type"`
	Amount             decimal.Decimal   `json:"amount" db:"amount"`
	Description        string            `json:"description" db:"description"`
	Notes              *string           `json:"notes,omitempty" db:"notes"`
	Source             TransactionSource `json:"source" db:"source"`
	TransactionDate    time.Time         `json:"transaction_date" db:"transaction_date"`
	DueDate            *time.Time        `json:"due_date,omitempty" db:"due_date"`
	PaymentDate        *time.Time        `json:"payment_date,omitempty" db:"payment_date"`
	IsPaid             bool              `json:"is_paid" db:"is_paid"`
	AutoPay            bool              `json:"auto_pay" db:"auto_pay"`
	IsRecurring        bool              `json:"is_recurring" db:"is_recurring"`
	RecurringID        *uuid.UUID        `json:"recurring_id,omitempty" db:"recurring_id"`
	InstallmentNumber  *int              `json:"installment_number,omitempty" db:"installment_number"`
	TotalInstallments  *int              `json:"total_installments,omitempty" db:"total_installments"`
	InstallmentGroupID *uuid.UUID        `json:"installment_group_id,omitempty" db:"installment_group_id"`
	Tags               []string          `json:"tags,omitempty" db:"tags"`
	AttachmentURL      *string           `json:"attachment_url,omitempty" db:"attachment_url"`
	ExternalID         *string           `json:"external_id,omitempty" db:"external_id"`
	CreatedAt          time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at" db:"updated_at"`

	Category    *Category    `json:"category,omitempty" db:"-"`
	BankAccount *BankAccount `json:"bank_account,omitempty" db:"-"`
}

func (t *Transaction) IsInstallment() bool {
	return t.InstallmentNumber != nil && t.TotalInstallments != nil
}

func (t *Transaction) IsOverdue() bool {
	if t.IsPaid || t.DueDate == nil {
		return false
	}
	return time.Now().After(*t.DueDate)
}

func (t *Transaction) DaysUntilDue() int {
	if t.DueDate == nil {
		return 0
	}
	duration := time.Until(*t.DueDate)
	return int(duration.Hours() / 24)
}

func (t *Transaction) CanBeDeleted() bool {
	return t.Source == TransactionSourceManual
}

func (t *Transaction) IsFromBank() bool {
	return t.Source == TransactionSourceBankSync
}

type RecurringFrequency string

const (
	FrequencyDaily     RecurringFrequency = "daily"
	FrequencyWeekly    RecurringFrequency = "weekly"
	FrequencyBiweekly  RecurringFrequency = "biweekly"
	FrequencyMonthly   RecurringFrequency = "monthly"
	FrequencyQuarterly RecurringFrequency = "quarterly"
	FrequencyYearly    RecurringFrequency = "yearly"
)

type RecurringTransaction struct {
	ID              uuid.UUID          `json:"id" db:"id"`
	UserID          uuid.UUID          `json:"user_id" db:"user_id"`
	BankAccountID   uuid.UUID          `json:"bank_account_id" db:"bank_account_id"`
	CategoryID      *uuid.UUID         `json:"category_id,omitempty" db:"category_id"`
	Type            TransactionType    `json:"type" db:"type"`
	Amount          decimal.Decimal    `json:"amount" db:"amount"`
	Description     string             `json:"description" db:"description"`
	Frequency       RecurringFrequency `json:"frequency" db:"frequency"`
	DayOfMonth      *int               `json:"day_of_month,omitempty" db:"day_of_month"`
	DayOfWeek       *int               `json:"day_of_week,omitempty" db:"day_of_week"`
	StartDate       time.Time          `json:"start_date" db:"start_date"`
	EndDate         *time.Time         `json:"end_date,omitempty" db:"end_date"`
	NextOccurrence  time.Time          `json:"next_occurrence" db:"next_occurrence"`
	LastGeneratedAt *time.Time         `json:"last_generated_at,omitempty" db:"last_generated_at"`
	IsActive        bool               `json:"is_active" db:"is_active"`
	AutoConfirm     bool               `json:"auto_confirm" db:"auto_confirm"`
	CreatedAt       time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" db:"updated_at"`

	Category    *Category    `json:"category,omitempty" db:"-"`
	BankAccount *BankAccount `json:"bank_account,omitempty" db:"-"`
}

func (r *RecurringTransaction) ShouldGenerate() bool {
	if !r.IsActive {
		return false
	}
	if r.EndDate != nil && time.Now().After(*r.EndDate) {
		return false
	}
	return !time.Now().Before(r.NextOccurrence)
}

func (r *RecurringTransaction) CalculateNextOccurrence() time.Time {
	current := r.NextOccurrence
	switch r.Frequency {
	case FrequencyDaily:
		return current.AddDate(0, 0, 1)
	case FrequencyWeekly:
		return current.AddDate(0, 0, 7)
	case FrequencyBiweekly:
		return current.AddDate(0, 0, 14)
	case FrequencyMonthly:
		return current.AddDate(0, 1, 0)
	case FrequencyQuarterly:
		return current.AddDate(0, 3, 0)
	case FrequencyYearly:
		return current.AddDate(1, 0, 0)
	default:
		return current.AddDate(0, 1, 0)
	}
}

type Transfer struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	UserID            uuid.UUID       `json:"user_id" db:"user_id"`
	FromAccountID     uuid.UUID       `json:"from_account_id" db:"from_account_id"`
	ToAccountID       uuid.UUID       `json:"to_account_id" db:"to_account_id"`
	Amount            decimal.Decimal `json:"amount" db:"amount"`
	Description       *string         `json:"description,omitempty" db:"description"`
	TransferDate      time.Time       `json:"transfer_date" db:"transfer_date"`
	FromTransactionID *uuid.UUID      `json:"from_transaction_id,omitempty" db:"from_transaction_id"`
	ToTransactionID   *uuid.UUID      `json:"to_transaction_id,omitempty" db:"to_transaction_id"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`

	FromAccount *BankAccount `json:"from_account,omitempty" db:"-"`
	ToAccount   *BankAccount `json:"to_account,omitempty" db:"-"`
}

type Installment struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	UserID            uuid.UUID       `json:"user_id" db:"user_id"`
	BankAccountID     uuid.UUID       `json:"bank_account_id" db:"bank_account_id"`
	CategoryID        *uuid.UUID      `json:"category_id,omitempty" db:"category_id"`
	Description       string          `json:"description" db:"description"`
	TotalAmount       decimal.Decimal `json:"total_amount" db:"total_amount"`
	InstallmentAmount decimal.Decimal `json:"installment_amount" db:"installment_amount"`
	TotalInstallments int             `json:"total_installments" db:"total_installments"`
	PaidInstallments  int             `json:"paid_installments" db:"paid_installments"`
	FirstDueDate      time.Time       `json:"first_due_date" db:"first_due_date"`
	IsActive          bool            `json:"is_active" db:"is_active"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`

	Category    *Category    `json:"category,omitempty" db:"-"`
	BankAccount *BankAccount `json:"bank_account,omitempty" db:"-"`
}

func (i *Installment) RemainingInstallments() int {
	return i.TotalInstallments - i.PaidInstallments
}

func (i *Installment) RemainingAmount() decimal.Decimal {
	return i.InstallmentAmount.Mul(decimal.NewFromInt(int64(i.RemainingInstallments())))
}

func (i *Installment) Progress() decimal.Decimal {
	if i.TotalInstallments == 0 {
		return decimal.Zero
	}
	return decimal.NewFromInt(int64(i.PaidInstallments)).
		Div(decimal.NewFromInt(int64(i.TotalInstallments))).
		Mul(decimal.NewFromInt(100))
}
