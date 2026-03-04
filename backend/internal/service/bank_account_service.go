package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
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

// Create creates a new bank account
func (s *BankAccountService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateBankAccountRequest) (*model.BankAccount, error) {
	// Verificar limite de contas bancárias
	count, err := s.accountRepo.CountAccounts(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count bank accounts: %w", err)
	}
	if count >= 20 { // MaxBankAccountsPerUser
		return nil, fmt.Errorf("você atingiu o limite máximo de 20 contas bancárias")
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
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("account_id", account.ID.String()).
		Str("name", account.Name).
		Msg("bank account created")

	return account, nil
}

// GetByID retrieves a bank account by ID
func (s *BankAccountService) GetByID(ctx context.Context, userID, accountID uuid.UUID) (*model.BankAccount, error) {
	account, err := s.accountRepo.GetByIDAndUser(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// GetAll retrieves all bank accounts for a user
func (s *BankAccountService) GetAll(ctx context.Context, userID uuid.UUID) ([]*model.BankAccount, error) {
	accounts, err := s.accountRepo.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank accounts: %w", err)
	}
	return accounts, nil
}

// GetByType retrieves bank accounts by type
func (s *BankAccountService) GetByType(ctx context.Context, userID uuid.UUID, accountType model.AccountType) ([]*model.BankAccount, error) {
	accounts, err := s.accountRepo.GetByType(ctx, userID, accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank accounts by type: %w", err)
	}
	return accounts, nil
}

// Update updates a bank account
func (s *BankAccountService) Update(ctx context.Context, userID, accountID uuid.UUID, req *model.UpdateBankAccountRequest) (*model.BankAccount, error) {
	account, err := s.accountRepo.GetByIDAndUser(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
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
		return nil, fmt.Errorf("failed to update bank account: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("account_id", accountID.String()).
		Msg("bank account updated")

	return account, nil
}

// Delete soft deletes a bank account
func (s *BankAccountService) Delete(ctx context.Context, userID, accountID uuid.UUID) error {
	if err := s.accountRepo.Delete(ctx, accountID, userID); err != nil {
		return fmt.Errorf("failed to delete bank account: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("account_id", accountID.String()).
		Msg("bank account deleted")

	return nil
}

// GetTotalBalance retrieves the total balance across all accounts
func (s *BankAccountService) GetTotalBalance(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error) {
	total, err := s.accountRepo.GetTotalBalance(ctx, userID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get total balance: %w", err)
	}
	return total, nil
}

// CountAccounts returns the number of active accounts
func (s *BankAccountService) CountAccounts(ctx context.Context, userID uuid.UUID) (int64, error) {
	count, err := s.accountRepo.CountAccounts(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count accounts: %w", err)
	}
	return count, nil
}

// AdjustBalance adjusts the balance of an account
func (s *BankAccountService) AdjustBalance(ctx context.Context, accountID uuid.UUID, delta decimal.Decimal) error {
	if err := s.accountRepo.AdjustBalance(ctx, accountID, delta); err != nil {
		return fmt.Errorf("failed to adjust balance: %w", err)
	}
	return nil
}

// RecalculateBalance recomputes current_balance from initial_balance + all paid transactions
func (s *BankAccountService) RecalculateBalance(ctx context.Context, userID, accountID uuid.UUID) (*model.BankAccount, error) {
	account, err := s.accountRepo.GetByIDAndUser(ctx, accountID, userID)
	if err != nil {
		return nil, err
	}

	netFromTransactions, err := s.txRepo.GetNetBalanceForAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to compute net balance: %w", err)
	}

	newBalance := account.InitialBalance.Add(netFromTransactions)

	if err := s.accountRepo.UpdateBalance(ctx, accountID, newBalance); err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	account.CurrentBalance = newBalance

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("account_id", accountID.String()).
		Str("new_balance", newBalance.String()).
		Msg("account balance recalculated")

	return account, nil
}
