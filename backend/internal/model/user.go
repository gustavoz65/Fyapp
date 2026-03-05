// Pacote model define os modelos de domínio e objetos de transferência de dados da aplicação.
package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleUser    UserRole = "user"
	UserRoleAdmin   UserRole = "admin"
	UserRolePremium UserRole = "premium"
)

type User struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	Email             string     `json:"email" db:"email"`
	PasswordHash      string     `json:"-" db:"password_hash"`
	FirstName         string     `json:"first_name" db:"first_name"`
	LastName          string     `json:"last_name" db:"last_name"`
	Phone             *string    `json:"phone,omitempty" db:"phone"`
	AvatarURL         *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	PreferredCurrency string     `json:"preferred_currency" db:"preferred_currency"`
	PreferredLanguage string     `json:"preferred_language" db:"preferred_language"`
	Timezone          string     `json:"timezone" db:"timezone"`
	Role              UserRole   `json:"role" db:"role"`
	EmailVerified     bool       `json:"email_verified" db:"email_verified"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at,omitempty" db:"email_verified_at"`
	LastLoginAt          *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
	IsActive             bool       `json:"is_active" db:"is_active"`
	OnboardingCompleted  bool       `json:"onboarding_completed" db:"onboarding_completed"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u *User) IsPremium() bool {
	return u.Role == UserRolePremium || u.Role == UserRoleAdmin
}

type UserSession struct {
	ID               uuid.UUID `json:"id" db:"id"`
	UserID           uuid.UUID `json:"user_id" db:"user_id"`
	RefreshTokenHash string    `json:"-" db:"refresh_token_hash"`
	DeviceInfo       *string   `json:"device_info,omitempty" db:"device_info"`
	IPAddress        *string   `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent        *string   `json:"user_agent,omitempty" db:"user_agent"`
	ExpiresAt        time.Time `json:"expires_at" db:"expires_at"`
	IsRevoked        bool      `json:"is_revoked" db:"is_revoked"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *UserSession) IsValid() bool {
	return !s.IsRevoked && !s.IsExpired()
}

type UserSettings struct {
	ID                      uuid.UUID `json:"id" db:"id"`
	UserID                  uuid.UUID `json:"user_id" db:"user_id"`
	NotificationEmail       bool      `json:"notification_email" db:"notification_email"`
	NotificationPush        bool      `json:"notification_push" db:"notification_push"`
	NotificationSMS         bool      `json:"notification_sms" db:"notification_sms"`
	BudgetAlerts            bool      `json:"budget_alerts" db:"budget_alerts"`
	BillReminders           bool      `json:"bill_reminders" db:"bill_reminders"`
	BillReminderDays        int       `json:"bill_reminder_days" db:"bill_reminder_days"`
	WeeklySummary           bool      `json:"weekly_summary" db:"weekly_summary"`
	MonthlyReport           bool      `json:"monthly_report" db:"monthly_report"`
	LowBalanceAlert         bool      `json:"low_balance_alert" db:"low_balance_alert"`
	LowBalanceThreshold     float64   `json:"low_balance_threshold" db:"low_balance_threshold"`
	AllowManualTransactions bool      `json:"allow_manual_transactions" db:"allow_manual_transactions"`
	Theme                   string    `json:"theme" db:"theme"`
	DashboardLayout         *string   `json:"dashboard_layout,omitempty" db:"dashboard_layout"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time `json:"updated_at" db:"updated_at"`
}
