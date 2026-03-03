package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type NotificationService struct {
	notificationRepo *repository.NotificationRepository
	userRepo         *repository.UserRepository
	logger           *zerolog.Logger
	jobClient        *asynq.Client
}

func NewNotificationService(
	notificationRepo *repository.NotificationRepository,
	userRepo *repository.UserRepository,
	logger *zerolog.Logger,
	jobClient *asynq.Client,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		jobClient:        jobClient,
		userRepo:         userRepo,
		logger:           logger,
	}
}

// cria uma nova notificação para um usuário
func (s *NotificationService) Create(ctx context.Context, notification *model.Notification) error {
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	s.logger.Info().
		Str("user_id", notification.UserID.String()).
		Str("notification_id", notification.ID.String()).
		Str("type", string(notification.Type)).
		Msg("notification created")

	return nil
}

// cria notificações em massa para uma lista de usuários
func (s *NotificationService) CreateBulk(ctx context.Context, userIDs []uuid.UUID, notificationType model.NotificationType, title, message string, data interface{}) error {
	var dataJSON *string
	if data != nil {
		bytes, err := json.Marshal(data)
		if err == nil {
			str := string(bytes)
			dataJSON = &str
		}
	}

	for _, userID := range userIDs {
		notification := &model.Notification{
			UserID:  userID,
			Type:    notificationType,
			Title:   title,
			Message: message,
			Data:    dataJSON,
		}

		if err := s.notificationRepo.Create(ctx, notification); err != nil {
			s.logger.Error().Err(err).
				Str("user_id", userID.String()).
				Msg("failed to create notification for user")
		}
	}

	return nil
}

// GetByID retrieves a notification by ID
func (s *NotificationService) GetByID(ctx context.Context, userID, notificationID uuid.UUID) (*model.Notification, error) {
	notification, err := s.notificationRepo.GetByIDAndUser(ctx, notificationID, userID)
	if err != nil {
		return nil, err
	}
	return notification, nil
}

// GetAll retrieves all notifications for a user with pagination
func (s *NotificationService) GetAll(ctx context.Context, userID uuid.UUID, page, pageSize int) (*model.PaginatedResponse[*model.Notification], error) {
	pagination := repository.PaginationParams{
		Page:     page,
		PageSize: pageSize,
		SortBy:   "created_at",
		SortDir:  "DESC",
	}
	pagination.Validate([]string{"created_at", "is_read"})

	notifications, total, err := s.notificationRepo.GetAllByUser(ctx, userID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &model.PaginatedResponse[*model.Notification]{
		Data:       notifications,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}, nil
}

// GetUnread retrieves unread notifications for a user
func (s *NotificationService) GetUnread(ctx context.Context, userID uuid.UUID) ([]*model.Notification, error) {
	notifications, err := s.notificationRepo.GetUnreadByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %w", err)
	}
	return notifications, nil
}

// GetUnreadCount returns the count of unread notifications
func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, userID, notificationID uuid.UUID) error {
	if err := s.notificationRepo.MarkAsRead(ctx, notificationID, userID); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	if err := s.notificationRepo.MarkAllAsRead(ctx, userID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Msg("all notifications marked as read")

	return nil
}

// Delete deletes a notification
func (s *NotificationService) Delete(ctx context.Context, userID, notificationID uuid.UUID) error {
	if err := s.notificationRepo.Delete(ctx, notificationID, userID); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("notification_id", notificationID.String()).
		Msg("notification deleted")

	return nil
}

// SendBudgetAlert é uma função de conveniência para criar uma notificação de alerta de orçamento
func (s *NotificationService) SendBudgetAlert(ctx context.Context, userID uuid.UUID, budgetName string, percentage, remaining float64) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    model.NotificationTypeBudgetAlert,
		Title:   "Alerta de Orcamento",
		Message: fmt.Sprintf("Voce atingiu %.0f%% do seu orcamento '%s'. Restam R$ %.2f.", percentage, budgetName, remaining),
	}

	return s.Create(ctx, notification)
}

// SendBillReminder sends a bill reminder notification
func (s *NotificationService) SendBillReminder(ctx context.Context, userID uuid.UUID, description string, amount float64, daysUntilDue int) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    model.NotificationTypeBillReminder,
		Title:   "Lembrete de Pagamento",
		Message: fmt.Sprintf("A conta '%s' no valor de R$ %.2f vence em %d dias.", description, amount, daysUntilDue),
	}

	return s.Create(ctx, notification)
}

// SendGoalAchievedNotification sends a goal achieved notification
func (s *NotificationService) SendGoalAchievedNotification(ctx context.Context, userID uuid.UUID, goalName string, targetAmount float64) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    model.NotificationTypeGoalAchieved,
		Title:   "Meta Alcancada!",
		Message: fmt.Sprintf("Parabens! Voce alcancou sua meta '%s' de R$ %.2f.", goalName, targetAmount),
	}

	return s.Create(ctx, notification)
}

// SendLowBalanceAlert sends a low balance alert notification
func (s *NotificationService) SendLowBalanceAlert(ctx context.Context, userID uuid.UUID, accountName string, threshold, currentBalance float64) error {
	notification := &model.Notification{
		UserID:  userID,
		Type:    model.NotificationTypeLowBalance,
		Title:   "Saldo Baixo",
		Message: fmt.Sprintf("O saldo da conta '%s' esta abaixo de R$ %.2f. Saldo atual: R$ %.2f.", accountName, threshold, currentBalance),
	}

	return s.Create(ctx, notification)
}

// ProcessScheduledNotifications processes and sends scheduled notifications
func (s *NotificationService) ProcessScheduledNotifications(ctx context.Context) error {
	notifications, err := s.notificationRepo.GetPendingScheduled(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending scheduled notifications: %w", err)
	}

	for _, notification := range notifications {
		// Check user settings
		settings, err := s.userRepo.GetSettings(ctx, notification.UserID)
		if err != nil {
			s.logger.Warn().Err(err).
				Str("user_id", notification.UserID.String()).
				Msg("failed to get user settings for notification")
			continue
		}

		// Determine if we should send based on settings and notification type
		shouldSend := s.shouldSendNotification(notification.Type, settings)

		if shouldSend {
			// Mark as sent
			if err := s.notificationRepo.MarkAsSent(ctx, notification.ID); err != nil {
				s.logger.Error().Err(err).
					Str("notification_id", notification.ID.String()).
					Msg("failed to mark notification as sent")
			}

			s.logger.Info().
				Str("notification_id", notification.ID.String()).
				Str("user_id", notification.UserID.String()).
				Str("type", string(notification.Type)).
				Msg("scheduled notification processed")
		}
	}

	return nil
}

// shouldSendNotification determines if a notification should be sent based on user settings
func (s *NotificationService) shouldSendNotification(notificationType model.NotificationType, settings *model.UserSettings) bool {
	switch notificationType {
	case model.NotificationTypeBudgetAlert:
		return settings.BudgetAlerts
	case model.NotificationTypeBillReminder:
		return settings.BillReminders
	case model.NotificationTypeGoalAchieved:
		return true // Always send goal achieved notifications
	case model.NotificationTypeLowBalance:
		return settings.LowBalanceAlert
	default:
		return true
	}
}

// CleanupOldNotifications removes notifications older than the specified number of days
func (s *NotificationService) CleanupOldNotifications(ctx context.Context, days int) (int64, error) {
	deleted, err := s.notificationRepo.DeleteOldNotifications(ctx, days)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old notifications: %w", err)
	}

	if deleted > 0 {
		s.logger.Info().
			Int64("deleted_count", deleted).
			Int("days_old", days).
			Msg("cleaned up old notifications")
	}

	return deleted, nil
}
