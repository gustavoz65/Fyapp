package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeBudgetAlert      NotificationType = "budget_alert"
	NotificationTypeBillReminder     NotificationType = "bill_reminder"
	NotificationTypeGoalAchieved     NotificationType = "goal_achieved"
	NotificationTypeLowBalance       NotificationType = "low_balance"
	NotificationTypeTransactionAlert NotificationType = "transaction_alert"
	NotificationTypeSystem           NotificationType = "system"
)

type NotificationChannel string

const (
	NotificationChannelInApp NotificationChannel = "in_app"
	NotificationChannelEmail NotificationChannel = "email"
	NotificationChannelPush  NotificationChannel = "push"
	NotificationChannelSMS   NotificationChannel = "sms"
)

type Notification struct {
	ID           uuid.UUID        `json:"id" db:"id"`
	UserID       uuid.UUID        `json:"user_id" db:"user_id"`
	Type         NotificationType `json:"type" db:"type"`
	Title        string           `json:"title" db:"title"`
	Message      string           `json:"message" db:"message"`
	Data         *string          `json:"data,omitempty" db:"data"`
	IsRead       bool             `json:"is_read" db:"is_read"`
	ReadAt       *time.Time       `json:"read_at,omitempty" db:"read_at"`
	SentVia      []string         `json:"sent_via" db:"sent_via"`
	ScheduledFor *time.Time       `json:"scheduled_for,omitempty" db:"scheduled_for"`
	SentAt       *time.Time       `json:"sent_at,omitempty" db:"sent_at"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
}

func (n *Notification) MarkAsRead() {
	n.IsRead = true
	now := time.Now()
	n.ReadAt = &now
}

func (n *Notification) IsPending() bool {
	return n.SentAt == nil
}

func (n *Notification) IsScheduled() bool {
	return n.ScheduledFor != nil && n.SentAt == nil
}

func (n *Notification) ShouldSendNow() bool {
	if n.SentAt != nil {
		return false
	}
	if n.ScheduledFor == nil {
		return true
	}
	return time.Now().After(*n.ScheduledFor)
}

type NotificationPreferences struct {
	UserID          uuid.UUID             `json:"user_id"`
	EnabledChannels []NotificationChannel `json:"enabled_channels"`
	EnabledTypes    []NotificationType    `json:"enabled_types"`
	QuietHoursStart *string               `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd   *string               `json:"quiet_hours_end,omitempty"`
}

type NotificationTemplate struct {
	Type    NotificationType
	Title   string
	Message string
}

var NotificationTemplates = map[NotificationType]NotificationTemplate{
	NotificationTypeBudgetAlert: {
		Type:    NotificationTypeBudgetAlert,
		Title:   "Alerta de Orcamento",
		Message: "Voce atingiu %s%% do seu orcamento '%s'. Restam %s.",
	},
	NotificationTypeBillReminder: {
		Type:    NotificationTypeBillReminder,
		Title:   "Lembrete de Pagamento",
		Message: "A conta '%s' no valor de %s vence em %d dias.",
	},
	NotificationTypeGoalAchieved: {
		Type:    NotificationTypeGoalAchieved,
		Title:   "Meta Alcancada!",
		Message: "Parabens! Voce alcancou sua meta '%s' de %s.",
	},
	NotificationTypeLowBalance: {
		Type:    NotificationTypeLowBalance,
		Title:   "Saldo Baixo",
		Message: "O saldo da conta '%s' esta abaixo de %s. Saldo atual: %s.",
	},
	NotificationTypeTransactionAlert: {
		Type:    NotificationTypeTransactionAlert,
		Title:   "Nova Transacao",
		Message: "Uma nova transacao de %s foi registrada: %s.",
	},
}
