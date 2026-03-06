package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
	"github.com/gustavoz65/Fyapp/internal/validation"
)

type BankAccountService struct {
	accountRepo *repository.BankAccountRepository
	txRepo      *repository.TransactionRepository
	logger      *zerolog.Logger
}

func NewBankAccountService(accountRepo *repository.BankAccountRepository, txRepo *repository.TransactionRepository, logger *zerolog.Logger) *BankAccountService {
	return &BankAccountService{
		accountRepo: accountRepo,
		txRepo:      txRepo,
		logger:      logger,
	}
}

// Create cria uma nova conta bancária
func (s *BankAccountService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateBankAccountRequest) (*model.BankAccount, error) {
	// Verificar limite de contas bancárias
	count, err := s.accountRepo.CountAccounts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao contar contas bancárias: %w", err)
	}
	if count >= validation.MaxBankAccountsPerUser {
		return nil, fmt.Errorf(validation.ErrMaxBankAccountsExceeded) //nolint:staticcheck // user-facing message
	}

	account := &model.BankAccount{
		UserID:         userID,
		Name:           req.Name,
		AccountType:    req.AccountType,
		InitialBalance: req.InitialBalance,
		IsActive:       true,
		IncludeInTotal: true,
	}

	if req.BankName != "" {
		account.BankName = &req.BankName
	}
	if req.BankCode != "" {
		account.BankCode = &req.BankCode
	}
	if req.AccountNumber != "" {
		account.AccountNumber = &req.AccountNumber
	}
	if req.Agency != "" {
		account.Agency = &req.Agency
	}
	if req.CreditLimit != nil {
		account.CreditLimit = req.CreditLimit
	}
	if req.ClosingDay != nil {
		account.ClosingDay = req.ClosingDay
	}
	if req.DueDay != nil {
		account.DueDay = req.DueDay
	}
	if req.Currency != "" {
		account.Currency = req.Currency
	} else {
		account.Currency = "BRL"
	}
	if req.Color != "" {
		account.Color = req.Color
	} else {
		account.Color = "#10B981"
	}
	if req.Icon != "" {
		account.Icon = req.Icon
	} else {
		account.Icon = "bank"
	}
	if req.IncludeInTotal != nil {
		account.IncludeInTotal = *req.IncludeInTotal
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, fmt.Errorf("falha ao criar conta bancária: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("account_id", account.ID.String()).
		Str("name", account.Name).
		Msg("bank account created")

	return account, nil
}

// GetByID busca uma conta bancária pelo ID
func (s *BankAccountService) GetByID(ctx context.Context, userID, accountID uuid.UUID) (*model.BankAccount, error) {
	account, err := s.accountRepo.GetByIDAndUser(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetAll retorna todas as contas bancárias do usuário
func (s *BankAccountService) GetAll(ctx context.Context, userID uuid.UUID) ([]*model.BankAccount, error) {
	accounts, err := s.accountRepo.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar contas bancárias: %w", err)
	}
	return accounts, nil
}

// GetByType busca contas bancárias por tipo
func (s *BankAccountService) GetByType(ctx context.Context, userID uuid.UUID, accountType model.AccountType) ([]*model.BankAccount, error) {
	accounts, err := s.accountRepo.GetByType(ctx, userID, accountType)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar contas por tipo: %w", err)
	}
	return accounts, nil
}

// Update atualiza uma conta bancária
func (s *BankAccountService) Update(ctx context.Context, userID, accountID uuid.UUID, req *model.UpdateBankAccountRequest) (*model.BankAccount, error) {
	account, err := s.accountRepo.GetByIDAndUser(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	// Aplicar atualizações
	if req.Name != nil {
		account.Name = *req.Name
	}
	if req.BankName != nil {
		account.BankName = req.BankName
	}
	if req.BankCode != nil {
		account.BankCode = req.BankCode
	}
	if req.AccountNumber != nil {
		account.AccountNumber = req.AccountNumber
	}
	if req.Agency != nil {
		account.Agency = req.Agency
	}
	if req.CreditLimit != nil {
		account.CreditLimit = req.CreditLimit
	}
	if req.ClosingDay != nil {
		account.ClosingDay = req.ClosingDay
	}
	if req.DueDay != nil {
		account.DueDay = req.DueDay
	}
	if req.Color != nil {
		account.Color = *req.Color
	}
	if req.Icon != nil {
		account.Icon = *req.Icon
	}
	if req.IsActive != nil {
		account.IsActive = *req.IsActive
	}
	if req.IncludeInTotal != nil {
		account.IncludeInTotal = *req.IncludeInTotal
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return nil, fmt.Errorf("falha ao atualizar conta bancária: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("account_id", accountID.String()).
		Msg("bank account updated")

	return account, nil
}

// Delete desativa (soft delete) uma conta bancária
func (s *BankAccountService) Delete(ctx context.Context, userID, accountID uuid.UUID) error {
	if err := s.accountRepo.Delete(ctx, accountID, userID); err != nil {
		return fmt.Errorf("falha ao remover conta bancária: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("account_id", accountID.String()).
		Msg("bank account deleted")

	return nil
}

// GetTotalBalance retorna o saldo total de todas as contas
func (s *BankAccountService) GetTotalBalance(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error) {
	total, err := s.accountRepo.GetTotalBalance(ctx, userID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("falha ao obter saldo total: %w", err)
	}
	return total, nil
}

// CountAccounts retorna a quantidade de contas ativas
func (s *BankAccountService) CountAccounts(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := s.accountRepo.CountAccounts(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("falha ao contar contas: %w", err)
	}
	return count, nil
}

// RecalculateBalance recalcula o saldo atual com base em todas as transações pagas
func (s *BankAccountService) RecalculateBalance(ctx context.Context, userID, accountID uuid.UUID) (*model.BankAccount, error) {
	account, err := s.accountRepo.RecalculateBalance(ctx, accountID, userID)
	if err != nil {
		return nil, fmt.Errorf("falha ao recalcular saldo: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("account_id", accountID.String()).
		Msg("bank account balance recalculated")

	return account, nil
}

// AdjustBalance ajusta o saldo da conta por um valor delta
func (s *BankAccountService) AdjustBalance(ctx context.Context, accountID uuid.UUID, delta decimal.Decimal) error {
	if err := s.accountRepo.AdjustBalance(ctx, accountID, delta); err != nil {
		return fmt.Errorf("falha ao ajustar saldo: %w", err)
	}
	return nil
}
