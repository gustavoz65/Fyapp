package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
)

type RecurringTransactionService struct {
	repo            *repository.RecurringTransactionRepository
	transactionRepo *repository.TransactionRepository
	accountRepo     *repository.BankAccountRepository
	logger          *zerolog.Logger
}

func NewRecurringTransactionService(
	repo *repository.RecurringTransactionRepository,
	transactionRepo *repository.TransactionRepository,
	accountRepo *repository.BankAccountRepository,
	logger *zerolog.Logger,
) *RecurringTransactionService {
	return &RecurringTransactionService{
		repo:            repo,
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
		logger:          logger,
	}
}

func (s *RecurringTransactionService) Create(ctx context.Context, rt *model.RecurringTransaction) error {
	if rt.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}

	if rt.BankAccountID == uuid.Nil {
		return errors.New("bank_account_id is required")
	}

	if rt.Amount.IsZero() || rt.Amount.IsNegative() {
		return errors.New("amount must be greater than zero")
	}

	if rt.Description == "" {
		return errors.New("description is required")
	}

	// Verificar limite de transações recorrentes ativas
	count, err := s.repo.CountActiveByUser(ctx, rt.UserID)
	if err != nil {
		return fmt.Errorf("failed to count recurring transactions: %w", err)
	}
	if count >= 50 { // MaxRecurringTransactions
		return fmt.Errorf("você atingiu o limite de 50 transações recorrentes ativas")
	}

	_, err = s.accountRepo.GetByIDAndUser(ctx, rt.BankAccountID, rt.UserID)
	if err != nil {
		return fmt.Errorf("bank account not found")
	}

	if rt.NextOccurrence.IsZero() {
		rt.NextOccurrence = rt.StartDate
	}

	if rt.StartDate.IsZero() {
		rt.StartDate = time.Now()
	}

	return s.repo.Create(ctx, rt)
}

func (s *RecurringTransactionService) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.RecurringTransaction, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *RecurringTransactionService) GetAll(ctx context.Context, userID uuid.UUID, isActive *bool) ([]*model.RecurringTransaction, error) {
	return s.repo.GetAll(ctx, userID, isActive)
}

func (s *RecurringTransactionService) Update(ctx context.Context, rt *model.RecurringTransaction) error {
	existing, err := s.repo.GetByID(ctx, rt.ID, rt.UserID)
	if err != nil {
		return err
	}

	if rt.BankAccountID != existing.BankAccountID {
		_, err := s.accountRepo.GetByIDAndUser(ctx, rt.BankAccountID, rt.UserID)
		if err != nil {
			return fmt.Errorf("bank account not found")
		}
	}

	if rt.Amount.IsZero() || rt.Amount.IsNegative() {
		return errors.New("amount must be greater than zero")
	}

	return s.repo.Update(ctx, rt)
}

func (s *RecurringTransactionService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *RecurringTransactionService) ToggleActive(ctx context.Context, id, userID uuid.UUID, isActive bool) error {
	return s.repo.ToggleActive(ctx, id, userID, isActive)
}

func (s *RecurringTransactionService) ProcessDueRecurrings(ctx context.Context) error {
	recurrings, err := s.repo.GetDueRecurrings(ctx)
	if err != nil {
		return err
	}

	s.logger.Info().Msgf("Processing %d due recurring transactions", len(recurrings))

	for _, rt := range recurrings {
		if err := s.generateTransaction(ctx, rt); err != nil {
			s.logger.Error().Err(err).Str("recurring_id", rt.ID.String()).Msg("Failed to generate transaction from recurring")
			continue
		}

		nextOccurrence := rt.CalculateNextOccurrence()
		if err := s.repo.UpdateNextOccurrence(ctx, rt.ID, nextOccurrence); err != nil {
			s.logger.Error().Err(err).Str("recurring_id", rt.ID.String()).Msg("Failed to update next occurrence")
		}
	}

	return nil
}

func (s *RecurringTransactionService) generateTransaction(ctx context.Context, rt *model.RecurringTransaction) error {
	isPaid := rt.AutoConfirm
	var paymentDate *time.Time
	if isPaid {
		now := time.Now()
		paymentDate = &now
	}

	transaction := &model.Transaction{
		UserID:          rt.UserID,
		BankAccountID:   rt.BankAccountID,
		CategoryID:      rt.CategoryID,
		Type:            rt.Type,
		Amount:          rt.Amount,
		Description:     rt.Description,
		Source:          model.TransactionSourceRecurring,
		TransactionDate: rt.NextOccurrence,
		DueDate:         &rt.NextOccurrence,
		PaymentDate:     paymentDate,
		IsPaid:          isPaid,
		IsRecurring:     true,
		RecurringID:     &rt.ID,
	}

	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return err
	}

	if isPaid {
		var amount decimal.Decimal
		if rt.Type == model.TransactionTypeExpense {
			amount = rt.Amount.Neg()
		} else {
			amount = rt.Amount
		}

		if err := s.accountRepo.UpdateBalance(ctx, rt.BankAccountID, amount); err != nil {
			s.logger.Error().Err(err).Msg("Failed to update account balance")
		}
	}

	s.logger.Info().
		Str("transaction_id", transaction.ID.String()).
		Str("recurring_id", rt.ID.String()).
		Msg("Generated transaction from recurring")

	return nil
}

func (s *RecurringTransactionService) GetStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	return s.repo.GetStats(ctx, userID)
}

func (s *RecurringTransactionService) GetUpcoming(ctx context.Context, userID uuid.UUID, days int) ([]*model.RecurringTransaction, error) {
	all, err := s.repo.GetAll(ctx, userID, nil)
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().AddDate(0, 0, days)
	var upcoming []*model.RecurringTransaction

	for _, rt := range all {
		if rt.IsActive && rt.NextOccurrence.Before(cutoff) {
			upcoming = append(upcoming, rt)
		}
	}

	return upcoming, nil
}
