package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeCreditCard AccountType = "credit_card"
	AccountTypeInvestment AccountType = "investment"
	AccountTypeCash       AccountType = "cash"
	AccountTypeOther      AccountType = "other"
)

type BankAccount struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	UserID         uuid.UUID        `json:"user_id" db:"user_id"`
	Name           string           `json:"name" db:"name"`
	BankName       *string          `json:"bank_name,omitempty" db:"bank_name"`
	BankCode       *string          `json:"bank_code,omitempty" db:"bank_code"`
	AccountType    AccountType      `json:"account_type" db:"account_type"`
	AccountNumber  *string          `json:"account_number,omitempty" db:"account_number"`
	Agency         *string          `json:"agency,omitempty" db:"agency"`
	InitialBalance decimal.Decimal  `json:"initial_balance" db:"initial_balance"`
	CurrentBalance decimal.Decimal  `json:"current_balance" db:"current_balance"`
	CreditLimit    *decimal.Decimal `json:"credit_limit,omitempty" db:"credit_limit"`
	ClosingDay     *int             `json:"closing_day,omitempty" db:"closing_day"`
	DueDay         *int             `json:"due_day,omitempty" db:"due_day"`
	Currency       string           `json:"currency" db:"currency"`
	Color          string           `json:"color" db:"color"`
	Icon           string           `json:"icon" db:"icon"`
	IsActive       bool             `json:"is_active" db:"is_active"`
	IncludeInTotal bool             `json:"include_in_total" db:"include_in_total"`
	LastSyncAt     *time.Time       `json:"last_sync_at,omitempty" db:"last_sync_at"`
	ExternalID     *string          `json:"external_id,omitempty" db:"external_id"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}

func (a *BankAccount) IsCreditCard() bool {
	return a.AccountType == AccountTypeCreditCard
}

func (a *BankAccount) AvailableCredit() decimal.Decimal {
	if a.CreditLimit == nil {
		return decimal.Zero
	}
	return a.CreditLimit.Sub(a.CurrentBalance.Abs())
}

func (a *BankAccount) UsedCreditPercentage() decimal.Decimal {
	if a.CreditLimit == nil || a.CreditLimit.IsZero() {
		return decimal.Zero
	}
	return a.CurrentBalance.Abs().Div(*a.CreditLimit).Mul(decimal.NewFromInt(100))
}

type BankIntegration struct {
	ID                    uuid.UUID  `json:"id" db:"id"`
	UserID                uuid.UUID  `json:"user_id" db:"user_id"`
	BankAccountID         *uuid.UUID `json:"bank_account_id,omitempty" db:"bank_account_id"`
	Provider              string     `json:"provider" db:"provider"`
	ProviderAccountID     *string    `json:"provider_account_id,omitempty" db:"provider_account_id"`
	AccessTokenEncrypted  *string    `json:"-" db:"access_token_encrypted"`
	RefreshTokenEncrypted *string    `json:"-" db:"refresh_token_encrypted"`
	TokenExpiresAt        *time.Time `json:"token_expires_at,omitempty" db:"token_expires_at"`
	ConsentExpiresAt      *time.Time `json:"consent_expires_at,omitempty" db:"consent_expires_at"`
	Status                string     `json:"status" db:"status"`
	LastSyncAt            *time.Time `json:"last_sync_at,omitempty" db:"last_sync_at"`
	LastError             *string    `json:"last_error,omitempty" db:"last_error"`
	Metadata              *string    `json:"metadata,omitempty" db:"metadata"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

func (bi *BankIntegration) IsActive() bool {
	return bi.Status == "active"
}

func (bi *BankIntegration) NeedsRefresh() bool {
	if bi.TokenExpiresAt == nil {
		return true
	}
	return time.Now().After(bi.TokenExpiresAt.Add(-5 * time.Minute))
}
