package service

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/lib/parser"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
	"github.com/gustavoz65/Fyapp/internal/validation"
)

type TransactionService struct {
	txRepo            *repository.TransactionRepository
	accountRepo       *repository.BankAccountRepository
	budgetRepo        *repository.BudgetRepository
	userRepo          *repository.UserRepository
	categorizationSvc *CategorizationService
	logger            *zerolog.Logger
}

func NewTransactionService(
	txRepo *repository.TransactionRepository,
	accountRepo *repository.BankAccountRepository,
	budgetRepo *repository.BudgetRepository,
	userRepo *repository.UserRepository,
	categorizationSvc *CategorizationService,
	logger *zerolog.Logger,
) *TransactionService {
	return &TransactionService{
		txRepo:            txRepo,
		accountRepo:       accountRepo,
		budgetRepo:        budgetRepo,
		userRepo:          userRepo,
		categorizationSvc: categorizationSvc,
		logger:            logger,
	}
}

// Create creates a new transaction
func (s *TransactionService) Create(ctx context.Context, userID uuid.UUID, req *model.CreateTransactionRequest) (*model.Transaction, error) {
	// Check if user allows manual transactions
	settings, err := s.userRepo.GetSettings(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user settings: %w", err)
	}
	if !settings.AllowManualTransactions {
		return nil, fmt.Errorf("manual transactions are disabled in your settings")
	}

	// Verificar limite de transações por dia
	count, err := s.txRepo.CountTransactionsByUserToday(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count transactions: %w", err)
	}
	if count >= validation.MaxTransactionsPerDay {
		return nil, fmt.Errorf(validation.ErrMaxTransactionsPerDay)
	}

	// Validate account belongs to user
	account, err := s.accountRepo.GetByIDAndUser(ctx, req.BankAccountID, userID)
	if err != nil {
		return nil, fmt.Errorf("invalid bank account: %w", err)
	}

	tx := &model.Transaction{
		UserID:          userID,
		BankAccountID:   req.BankAccountID,
		CategoryID:      req.CategoryID,
		Type:            req.Type,
		Amount:          req.Amount,
		Description:     req.Description,
		Source:          model.TransactionSourceManual,
		TransactionDate: req.TransactionDate,
		DueDate:         req.DueDate,
		IsPaid:          req.IsPaid,
		AutoPay:         req.AutoPay,
		Tags:            req.Tags,
	}

	if req.Notes != "" {
		tx.Notes = &req.Notes
	}

	// Handle installments
	if req.TotalInstallments != nil && *req.TotalInstallments > 1 {
		firstTx, err := s.createInstallments(ctx, userID, tx, *req.TotalInstallments)
		if err != nil {
			return nil, err
		}
		// Aprende padrao de categorizacao para parcelas tambem
		if tx.CategoryID != nil && s.categorizationSvc != nil {
			s.categorizationSvc.LearnFromTransaction(ctx, userID, req.Description, *tx.CategoryID, "manual")
		}
		return firstTx, nil
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Aprende padrao de categorizacao
	if tx.CategoryID != nil && s.categorizationSvc != nil {
		s.categorizationSvc.LearnFromTransaction(ctx, userID, tx.Description, *tx.CategoryID, string(tx.Source))
	}

	// Update account balance if paid
	if tx.IsPaid {
		delta := tx.Amount
		if tx.Type == model.TransactionTypeExpense {
			delta = delta.Neg()
		}
		if err := s.accountRepo.AdjustBalance(ctx, account.ID, delta); err != nil {
			s.logger.Error().Err(err).Msg("failed to update account balance")
		}
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("transaction_id", tx.ID.String()).
		Str("type", string(tx.Type)).
		Str("amount", tx.Amount.String()).
		Msg("transaction created")

	return tx, nil
}

// createInstallments creates multiple transactions for installment payments within a single DB transaction
func (s *TransactionService) createInstallments(ctx context.Context, userID uuid.UUID, baseTx *model.Transaction, totalInstallments int) (*model.Transaction, error) {
	groupID := uuid.New()
	installmentAmount := baseTx.Amount.Div(decimal.NewFromInt(int64(totalInstallments)))

	var firstTx *model.Transaction
	currentDate := baseTx.TransactionDate

	err := s.txRepo.WithTx(ctx, func(dbTx *sql.Tx) error {
		for i := 1; i <= totalInstallments; i++ {
			installmentNum := i
			tx := &model.Transaction{
				UserID:             userID,
				BankAccountID:      baseTx.BankAccountID,
				CategoryID:         baseTx.CategoryID,
				Type:               baseTx.Type,
				Amount:             installmentAmount,
				Description:        fmt.Sprintf("%s (%d/%d)", baseTx.Description, i, totalInstallments),
				Source:             model.TransactionSourceManual,
				TransactionDate:    currentDate,
				IsPaid:             i == 1 && baseTx.IsPaid,
				InstallmentNumber:  &installmentNum,
				TotalInstallments:  &totalInstallments,
				InstallmentGroupID: &groupID,
				Tags:               baseTx.Tags,
			}

			if baseTx.DueDate != nil {
				dueDate := baseTx.DueDate.AddDate(0, i-1, 0)
				tx.DueDate = &dueDate
			}

			if err := s.txRepo.CreateTx(ctx, dbTx, tx); err != nil {
				return fmt.Errorf("failed to create installment %d: %w", i, err)
			}

			if i == 1 {
				firstTx = tx

				if tx.IsPaid {
					delta := tx.Amount
					if tx.Type == model.TransactionTypeExpense {
						delta = delta.Neg()
					}
					if err := s.accountRepo.AdjustBalanceTx(ctx, dbTx, baseTx.BankAccountID, delta); err != nil {
						return fmt.Errorf("failed to adjust balance for first installment: %w", err)
					}
				}
			}

			currentDate = currentDate.AddDate(0, 1, 0)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("group_id", groupID.String()).
		Int("installments", totalInstallments).
		Msg("installment transactions created")

	return firstTx, nil
}

// GetByID retrieves a transaction by ID
func (s *TransactionService) GetByID(ctx context.Context, userID, txID uuid.UUID) (*model.Transaction, error) {
	tx, err := s.txRepo.GetByIDAndUser(ctx, txID, userID)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetByFilter retrieves transactions based on filter criteria
func (s *TransactionService) GetByFilter(ctx context.Context, filter *model.TransactionFilter) (*model.PaginatedResponse[*model.Transaction], error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	transactions, total, err := s.txRepo.GetByFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &model.PaginatedResponse[*model.Transaction]{
		Data:       transactions,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalItems: total,
		TotalPages: totalPages,
		HasMore:    filter.Page < totalPages,
	}, nil
}

// GetByDateRange retrieves transactions within a date range
func (s *TransactionService) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*model.Transaction, error) {
	transactions, err := s.txRepo.GetByDateRange(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}
	return transactions, nil
}

// GetUpcomingBills retrieves upcoming unpaid bills
func (s *TransactionService) GetUpcomingBills(ctx context.Context, userID uuid.UUID, days int) ([]*model.Transaction, error) {
	if days <= 0 {
		days = 7
	}
	transactions, err := s.txRepo.GetUpcomingBills(ctx, userID, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming bills: %w", err)
	}
	return transactions, nil
}

// GetRecentTransactions retrieves recent transactions
func (s *TransactionService) GetRecentTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]*model.Transaction, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	transactions, err := s.txRepo.GetRecentTransactions(ctx, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	return transactions, nil
}

// Update updates a transaction
func (s *TransactionService) Update(ctx context.Context, userID, txID uuid.UUID, req *model.UpdateTransactionRequest) (*model.Transaction, error) {
	tx, err := s.txRepo.GetByIDAndUser(ctx, txID, userID)
	if err != nil {
		return nil, err
	}

	oldAmount := tx.Amount
	oldIsPaid := tx.IsPaid
	oldType := tx.Type

	// Apply updates
	if req.CategoryID != nil {
		tx.CategoryID = req.CategoryID
	}
	if req.Amount != nil {
		tx.Amount = *req.Amount
	}
	if req.Description != nil {
		tx.Description = *req.Description
	}
	if req.Notes != nil {
		tx.Notes = req.Notes
	}
	if req.TransactionDate != nil {
		tx.TransactionDate = *req.TransactionDate
	}
	if req.DueDate != nil {
		tx.DueDate = req.DueDate
	}
	if req.IsPaid != nil {
		tx.IsPaid = *req.IsPaid
	}
	if req.Tags != nil {
		tx.Tags = req.Tags
	}

	if err := s.txRepo.Update(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	// Aprende padrao de categorizacao quando usuario altera/define a categoria
	if req.CategoryID != nil && tx.CategoryID != nil && s.categorizationSvc != nil {
		s.categorizationSvc.LearnFromTransaction(ctx, userID, tx.Description, *tx.CategoryID, string(tx.Source))
	}

	// Handle balance adjustments
	s.handleBalanceAdjustment(ctx, tx.BankAccountID, oldAmount, oldIsPaid, oldType, tx.Amount, tx.IsPaid, tx.Type)

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("transaction_id", txID.String()).
		Msg("transaction updated")

	return tx, nil
}

// handleBalanceAdjustment handles account balance adjustments when a transaction changes
func (s *TransactionService) handleBalanceAdjustment(
	ctx context.Context,
	accountID uuid.UUID,
	oldAmount decimal.Decimal, oldIsPaid bool, oldType model.TransactionType,
	newAmount decimal.Decimal, newIsPaid bool, newType model.TransactionType,
) {
	// Calculate the delta
	var delta decimal.Decimal

	// Remove old impact
	if oldIsPaid {
		if oldType == model.TransactionTypeIncome {
			delta = delta.Sub(oldAmount)
		} else {
			delta = delta.Add(oldAmount)
		}
	}

	// Add new impact
	if newIsPaid {
		if newType == model.TransactionTypeIncome {
			delta = delta.Add(newAmount)
		} else {
			delta = delta.Sub(newAmount)
		}
	}

	if !delta.IsZero() {
		_ = s.accountRepo.AdjustBalance(ctx, accountID, delta)
	}
}

// MarkAsPaid marks a transaction as paid
func (s *TransactionService) MarkAsPaid(ctx context.Context, userID, txID uuid.UUID) (*model.Transaction, error) {
	tx, err := s.txRepo.GetByIDAndUser(ctx, txID, userID)
	if err != nil {
		return nil, err
	}

	if tx.IsPaid {
		return tx, nil // Already paid
	}

	if err := s.txRepo.MarkAsPaid(ctx, txID, userID); err != nil {
		return nil, fmt.Errorf("failed to mark transaction as paid: %w", err)
	}

	// Update balance
	delta := tx.Amount
	if tx.Type == model.TransactionTypeExpense {
		delta = delta.Neg()
	}
	_ = s.accountRepo.AdjustBalance(ctx, tx.BankAccountID, delta)

	tx.IsPaid = true
	now := time.Now()
	tx.PaymentDate = &now

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("transaction_id", txID.String()).
		Msg("transaction marked as paid")

	return tx, nil
}

// Delete deletes a transaction
func (s *TransactionService) Delete(ctx context.Context, userID, txID uuid.UUID) error {
	tx, err := s.txRepo.GetByIDAndUser(ctx, txID, userID)
	if err != nil {
		return err
	}

	if !tx.CanBeDeleted() {
		return fmt.Errorf("cannot delete transaction from source: %s (only manual transactions can be deleted)", tx.Source)
	}

	// DELETE + ajuste de saldo são atômicos: se o ajuste falhar, o DELETE é revertido
	err = s.txRepo.WithTx(ctx, func(dbTx *sql.Tx) error {
		if err := s.txRepo.DeleteTx(ctx, dbTx, txID, userID); err != nil {
			return fmt.Errorf("failed to delete transaction: %w", err)
		}

		if tx.IsPaid {
			delta := tx.Amount
			if tx.Type == model.TransactionTypeIncome {
				delta = delta.Neg()
			}
			if err := s.accountRepo.AdjustBalanceTx(ctx, dbTx, tx.BankAccountID, delta); err != nil {
				return fmt.Errorf("failed to adjust balance: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("transaction_id", txID.String()).
		Str("source", string(tx.Source)).
		Msg("transaction deleted")

	return nil
}

// GetMonthlyTotals returns income and expense totals for a month
func (s *TransactionService) GetMonthlyTotals(ctx context.Context, userID uuid.UUID, year, month int) (income, expense decimal.Decimal, err error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)

	income, err = s.txRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startDate, endDate)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	expense, err = s.txRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startDate, endDate)
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}

	return income, expense, nil
}

// GetExpensesByCategory returns expenses grouped by category
func (s *TransactionService) GetExpensesByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]model.CategoryAmount, error) {
	return s.txRepo.GetSumByCategory(ctx, userID, startDate, endDate)
}

// AutoReconcile processa transações vencidas com auto_pay=true, aplica categorização automática
// se a transação não tiver categoria, e marca como pagas atualizando o saldo da conta.
func (s *TransactionService) AutoReconcile(ctx context.Context) error {
	today := time.Now()
	transactions, err := s.txRepo.GetDueForAutoReconcile(ctx, today)
	if err != nil {
		return fmt.Errorf("failed to fetch due transactions for auto reconcile: %w", err)
	}

	s.logger.Info().Msgf("Auto-reconciling %d due transactions", len(transactions))

	for _, tx := range transactions {
		// Aplica categorização automática se a transação não tiver categoria
		if tx.CategoryID == nil && s.categorizationSvc != nil {
			suggestion, err := s.categorizationSvc.SuggestCategory(ctx, tx.UserID, tx.Description)
			if err == nil && len(suggestion.Suggestions) > 0 && suggestion.Suggestions[0].Confidence >= 1 {
				catID := suggestion.Suggestions[0].CategoryID
				tx.CategoryID = &catID

				// Persiste a categoria sugerida na transação
				if updateErr := s.txRepo.Update(ctx, tx); updateErr != nil {
					s.logger.Error().Err(updateErr).
						Str("transaction_id", tx.ID.String()).
						Msg("Failed to update category during auto reconcile")
				} else {
					s.logger.Debug().
						Str("transaction_id", tx.ID.String()).
						Str("category_id", catID.String()).
						Str("category_name", suggestion.Suggestions[0].CategoryName).
						Msg("Auto-categorized transaction during reconciliation")
				}
			}
		}

		// Marca como paga
		if err := s.txRepo.MarkAsPaid(ctx, tx.ID, tx.UserID); err != nil {
			s.logger.Error().Err(err).
				Str("transaction_id", tx.ID.String()).
				Msg("Failed to mark transaction as paid during auto reconcile")
			continue
		}

		// Atualiza saldo da conta
		delta := tx.Amount
		if tx.Type == model.TransactionTypeExpense {
			delta = delta.Neg()
		}
		if err := s.accountRepo.AdjustBalance(ctx, tx.BankAccountID, delta); err != nil {
			s.logger.Error().Err(err).
				Str("transaction_id", tx.ID.String()).
				Msg("Failed to adjust balance during auto reconcile")
		}

		s.logger.Info().
			Str("transaction_id", tx.ID.String()).
			Str("description", tx.Description).
			Str("amount", tx.Amount.String()).
			Msg("Auto-reconciled transaction")
	}

	return nil
}

// ImportResult contains the result of a transaction import
type ImportResult struct {
	TotalImported int                  `json:"total_imported"`
	Duplicates    int                  `json:"duplicates"`
	Errors        []string             `json:"errors"`
	Transactions  []*model.Transaction `json:"transactions"`
}

// ImportTransactionsWithProgress imports transactions with progress callback.
// Todo o processo de criação e ajuste de saldo é atômico: ou todas as transações
// são importadas com sucesso, ou nenhuma é persistida (rollback automático).
func (s *TransactionService) ImportTransactionsWithProgress(
	ctx context.Context,
	userID uuid.UUID,
	bankAccountID uuid.UUID,
	reader io.Reader,
	bankType string,
	forceReimport bool,
	progressCallback func(current, total int),
) (*ImportResult, error) {
	account, err := s.accountRepo.GetByIDAndUser(ctx, bankAccountID, userID)
	if err != nil {
		return nil, fmt.Errorf("invalid bank account: %w", err)
	}

	txParser := parser.NewTransactionParser()
	importedTxs, err := txParser.ParseCSV(reader, bankType)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	if len(importedTxs) == 0 {
		return &ImportResult{
			TotalImported: 0,
			Duplicates:    0,
			Errors:        []string{"No transactions found in file"},
			Transactions:  []*model.Transaction{},
		}, nil
	}

	total := len(importedTxs)
	progressCallback(0, total)

	existingIDs := make(map[string]bool)
	if !forceReimport {
		externalIDs := make([]string, len(importedTxs))
		for i, tx := range importedTxs {
			externalIDs[i] = tx.ExternalID
		}
		foundIDs, err := s.txRepo.GetExistingExternalIDs(ctx, userID, externalIDs)
		if err != nil {
			s.logger.Warn().Err(err).Msg("failed to check for duplicates, proceeding without deduplication")
		} else {
			existingIDs = foundIDs
		}
	}

	result := &ImportResult{
		Transactions: make([]*model.Transaction, 0),
		Errors:       make([]string, 0),
	}

	// Toda a inserção e ajuste de saldo ocorre dentro de uma única transação de banco de dados.
	// Se qualquer operação falhar, tudo é revertido automaticamente.
	importErr := s.txRepo.WithTx(ctx, func(dbTx *sql.Tx) error {
		for idx, importedTx := range importedTxs {
			progressCallback(idx+1, total)

			if !forceReimport && existingIDs[importedTx.ExternalID] {
				result.Duplicates++
				s.logger.Debug().
					Str("external_id", importedTx.ExternalID).
					Str("description", importedTx.Description).
					Msg("skipping duplicate transaction")
				continue
			}

			tx := &model.Transaction{
				UserID:          userID,
				BankAccountID:   bankAccountID,
				Type:            model.TransactionType(importedTx.Type),
				Amount:          importedTx.Amount,
				Description:     importedTx.Description,
				Source:          model.TransactionSourceBankSync,
				TransactionDate: importedTx.Date,
				IsPaid:          true,
				ExternalID:      &importedTx.ExternalID,
			}
			tx.PaymentDate = &importedTx.Date

			// Auto-categorize (operação de leitura, segura dentro da tx)
			if s.categorizationSvc != nil {
				suggestion, err := s.categorizationSvc.SuggestCategory(ctx, userID, tx.Description)
				if err == nil && len(suggestion.Suggestions) > 0 && suggestion.Suggestions[0].Confidence >= 1 {
					tx.CategoryID = &suggestion.Suggestions[0].CategoryID
					s.logger.Debug().
						Str("description", tx.Description).
						Str("category", suggestion.Suggestions[0].CategoryName).
						Msg("auto-categorized imported transaction")
				}
			}

			if err := s.txRepo.CreateTx(ctx, dbTx, tx); err != nil {
				return fmt.Errorf("failed to import transaction '%s': %w", tx.Description, err)
			}

			delta := tx.Amount
			if tx.Type == model.TransactionTypeExpense {
				delta = delta.Neg()
			}
			if err := s.accountRepo.AdjustBalanceTx(ctx, dbTx, account.ID, delta); err != nil {
				return fmt.Errorf("failed to adjust balance for transaction '%s': %w", tx.Description, err)
			}

			result.TotalImported++
			result.Transactions = append(result.Transactions, tx)
		}
		return nil
	})

	if importErr != nil {
		return nil, fmt.Errorf("import failed and was rolled back: %w", importErr)
	}

	// Aprende padrões de categorização após o commit (best-effort)
	for _, tx := range result.Transactions {
		if tx.CategoryID != nil && s.categorizationSvc != nil {
			s.categorizationSvc.LearnFromTransaction(ctx, userID, tx.Description, *tx.CategoryID, "bank_sync")
		}
	}

	progressCallback(total, total)

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("bank_account_id", bankAccountID.String()).
		Int("total_imported", result.TotalImported).
		Int("duplicates", result.Duplicates).
		Msg("completed transaction import with progress")

	return result, nil
}

// DeleteAllByAccount deleta todas as transações de uma conta bancária.
// Operação atômica: calcula o ajuste de saldo e deleta tudo em uma única transação de banco de dados.
func (s *TransactionService) DeleteAllByAccount(ctx context.Context, userID, accountID uuid.UUID) error {
	if _, err := s.accountRepo.GetByIDAndUser(ctx, accountID, userID); err != nil {
		return fmt.Errorf("conta bancária inválida: %w", err)
	}

	return s.txRepo.WithTx(ctx, func(dbTx *sql.Tx) error {
		// 1. Calcula o delta para reverter o impacto de todas as transações pagas
		delta, err := s.txRepo.GetPaidBalanceDeltaByAccountTx(ctx, dbTx, accountID, userID)
		if err != nil {
			return fmt.Errorf("falha ao calcular ajuste de saldo: %w", err)
		}

		// 2. Deleta todas as transações da conta em uma única query
		deleted, err := s.txRepo.DeleteAllByAccountTx(ctx, dbTx, accountID, userID)
		if err != nil {
			return fmt.Errorf("falha ao deletar transações: %w", err)
		}

		// 3. Ajusta o saldo da conta para reverter o impacto das transações deletadas
		if !delta.IsZero() {
			if err := s.accountRepo.AdjustBalanceTx(ctx, dbTx, accountID, delta); err != nil {
				return fmt.Errorf("falha ao ajustar saldo: %w", err)
			}
		}

		s.logger.Info().
			Str("user_id", userID.String()).
			Str("account_id", accountID.String()).
			Int64("deleted", deleted).
			Msg("deletadas todas as transações da conta")

		return nil
	})
}

// BulkDeleteRequest holds parameters for bulk transaction deletion
type BulkDeleteRequest struct {
	AccountID *uuid.UUID `json:"account_id,omitempty"`
	StartDate time.Time  `json:"start_date"`
	EndDate   time.Time  `json:"end_date"`
}

// BulkDeleteResult contains the result of bulk deletion
type BulkDeleteResult struct {
	DeletedCount int64 `json:"deleted_count"`
}

// BulkDeleteByDateRange deletes multiple transactions within a date range.
// Operação atômica: calcula ajustes de saldo por conta e deleta tudo em uma única transação de banco de dados.
func (s *TransactionService) BulkDeleteByDateRange(ctx context.Context, userID uuid.UUID, req *BulkDeleteRequest) (*BulkDeleteResult, error) {
	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("data final deve ser posterior à data inicial")
	}

	if req.AccountID != nil {
		if _, err := s.accountRepo.GetByIDAndUser(ctx, *req.AccountID, userID); err != nil {
			return nil, fmt.Errorf("conta bancária inválida: %w", err)
		}
	}

	var result BulkDeleteResult

	err := s.txRepo.WithTx(ctx, func(dbTx *sql.Tx) error {
		// 1. Calcula o delta por conta para reverter o impacto das transações pagas
		deltas, err := s.txRepo.GetPaidBalanceDeltasByDateRangeTx(ctx, dbTx, userID, req.StartDate, req.EndDate, req.AccountID)
		if err != nil {
			return fmt.Errorf("falha ao calcular ajustes de saldo: %w", err)
		}

		// 2. Deleta as transações em uma única query
		var deleted int64
		if req.AccountID != nil {
			deleted, err = s.txRepo.DeleteByAccountAndDateRangeTx(ctx, dbTx, userID, *req.AccountID, req.StartDate, req.EndDate)
		} else {
			deleted, err = s.txRepo.DeleteByDateRangeTx(ctx, dbTx, userID, req.StartDate, req.EndDate)
		}
		if err != nil {
			return fmt.Errorf("falha ao deletar transações: %w", err)
		}
		result.DeletedCount = deleted

		// 3. Ajusta o saldo de cada conta afetada
		for accountID, delta := range deltas {
			if err := s.accountRepo.AdjustBalanceTx(ctx, dbTx, accountID, delta); err != nil {
				return fmt.Errorf("falha ao ajustar saldo da conta %s: %w", accountID, err)
			}
		}

		s.logger.Info().
			Str("user_id", userID.String()).
			Str("start_date", req.StartDate.Format("2006-01-02")).
			Str("end_date", req.EndDate.Format("2006-01-02")).
			Int64("deleted", deleted).
			Msg("transações deletadas em lote")

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}
