package service

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/lib/parser"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
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
	if count >= 100 { // MaxTransactionsPerDay
		return nil, fmt.Errorf("você atingiu o limite de 100 transações por dia")
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

// createInstallments creates multiple transactions for installment payments
func (s *TransactionService) createInstallments(ctx context.Context, userID uuid.UUID, baseTx *model.Transaction, totalInstallments int) (*model.Transaction, error) {
	groupID := uuid.New()
	installmentAmount := baseTx.Amount.Div(decimal.NewFromInt(int64(totalInstallments)))

	var firstTx *model.Transaction
	currentDate := baseTx.TransactionDate

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

		if err := s.txRepo.Create(ctx, tx); err != nil {
			return nil, fmt.Errorf("failed to create installment %d: %w", i, err)
		}

		if i == 1 {
			firstTx = tx

			// Update balance for first installment if paid
			if tx.IsPaid {
				delta := tx.Amount
				if tx.Type == model.TransactionTypeExpense {
					delta = delta.Neg()
				}
				_ = s.accountRepo.AdjustBalance(ctx, baseTx.BankAccountID, delta)
			}
		}

		// Move to next month
		currentDate = currentDate.AddDate(0, 1, 0)
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

	// Validate if transaction can be deleted
	if !tx.CanBeDeleted() {
		return fmt.Errorf("cannot delete transaction from source: %s (only manual transactions can be deleted)", tx.Source)
	}

	if err := s.txRepo.Delete(ctx, txID, userID); err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	// Reverse balance impact if it was paid
	if tx.IsPaid {
		delta := tx.Amount
		if tx.Type == model.TransactionTypeIncome {
			delta = delta.Neg()
		}
		_ = s.accountRepo.AdjustBalance(ctx, tx.BankAccountID, delta)
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
	TotalImported int                     `json:"total_imported"`
	Duplicates    int                     `json:"duplicates"`
	Errors        []string                `json:"errors"`
	Transactions  []*model.Transaction    `json:"transactions"`
}

// ImportTransactions imports transactions from CSV file
func (s *TransactionService) ImportTransactions(
	ctx context.Context,
	userID uuid.UUID,
	bankAccountID uuid.UUID,
	reader io.Reader,
	bankType string,
) (*ImportResult, error) {
	// Validate account belongs to user
	account, err := s.accountRepo.GetByIDAndUser(ctx, bankAccountID, userID)
	if err != nil {
		return nil, fmt.Errorf("invalid bank account: %w", err)
	}

	// Parse CSV file
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

	// Get existing external IDs for deduplication
	externalIDs := make([]string, len(importedTxs))
	for i, tx := range importedTxs {
		externalIDs[i] = tx.ExternalID
	}

	existingIDs, err := s.txRepo.GetExistingExternalIDs(ctx, userID, externalIDs)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to check for duplicates, proceeding without deduplication")
		existingIDs = make(map[string]bool)
	}

	// Process transactions
	result := &ImportResult{
		Transactions: make([]*model.Transaction, 0),
		Errors:       make([]string, 0),
	}

	for _, importedTx := range importedTxs {
		// Check for duplicates
		if existingIDs[importedTx.ExternalID] {
			result.Duplicates++
			s.logger.Debug().
				Str("external_id", importedTx.ExternalID).
				Str("description", importedTx.Description).
				Msg("skipping duplicate transaction")
			continue
		}

		// Convert TransactionImport to Transaction
		tx := &model.Transaction{
			UserID:          userID,
			BankAccountID:   bankAccountID,
			Type:            model.TransactionType(importedTx.Type),
			Amount:          importedTx.Amount,
			Description:     importedTx.Description,
			Source:          model.TransactionSourceBankSync,
			TransactionDate: importedTx.Date,
			IsPaid:          true, // Imported transactions are already executed
			ExternalID:      &importedTx.ExternalID,
		}

		// Set payment date same as transaction date for imported
		tx.PaymentDate = &importedTx.Date

		// Auto-categorize if available
		if s.categorizationSvc != nil {
			suggestion, err := s.categorizationSvc.SuggestCategory(ctx, userID, tx.Description)
			if err == nil && len(suggestion.Suggestions) > 0 {
				bestSuggestion := suggestion.Suggestions[0]
				if bestSuggestion.Confidence >= 1 { // Only use if confidence is at least 1
					tx.CategoryID = &bestSuggestion.CategoryID
					s.logger.Debug().
						Str("description", tx.Description).
						Str("category", bestSuggestion.CategoryName).
						Int("confidence", bestSuggestion.Confidence).
						Msg("auto-categorized imported transaction")
				}
			}
		}

		// Create transaction
		if err := s.txRepo.Create(ctx, tx); err != nil {
			errMsg := fmt.Sprintf("Failed to import transaction '%s': %v", tx.Description, err)
			result.Errors = append(result.Errors, errMsg)
			s.logger.Error().Err(err).
				Str("description", tx.Description).
				Msg("failed to create imported transaction")
			continue
		}

		// Update account balance (since imported transactions are already paid)
		delta := tx.Amount
		if tx.Type == model.TransactionTypeExpense {
			delta = delta.Neg()
		}
		if err := s.accountRepo.AdjustBalance(ctx, account.ID, delta); err != nil {
			s.logger.Error().Err(err).
				Str("transaction_id", tx.ID.String()).
				Msg("failed to update account balance for imported transaction")
		}

		result.TotalImported++
		result.Transactions = append(result.Transactions, tx)

		// Learn categorization pattern if categorized
		if tx.CategoryID != nil && s.categorizationSvc != nil {
			s.categorizationSvc.LearnFromTransaction(ctx, userID, tx.Description, *tx.CategoryID, "bank_sync")
		}
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Str("bank_account_id", bankAccountID.String()).
		Int("total_imported", result.TotalImported).
		Int("duplicates", result.Duplicates).
		Int("errors", len(result.Errors)).
		Msg("completed transaction import")

	return result, nil
}

// ImportTransactionsWithProgress imports transactions with progress callback
func (s *TransactionService) ImportTransactionsWithProgress(
	ctx context.Context,
	userID uuid.UUID,
	bankAccountID uuid.UUID,
	reader io.Reader,
	bankType string,
	progressCallback func(current, total int),
) (*ImportResult, error) {
	// Validate account belongs to user
	account, err := s.accountRepo.GetByIDAndUser(ctx, bankAccountID, userID)
	if err != nil {
		return nil, fmt.Errorf("invalid bank account: %w", err)
	}

	// Parse CSV file
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

	// Get existing external IDs for deduplication
	externalIDs := make([]string, len(importedTxs))
	for i, tx := range importedTxs {
		externalIDs[i] = tx.ExternalID
	}

	existingIDs, err := s.txRepo.GetExistingExternalIDs(ctx, userID, externalIDs)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to check for duplicates, proceeding without deduplication")
		existingIDs = make(map[string]bool)
	}

	// Process transactions
	result := &ImportResult{
		Transactions: make([]*model.Transaction, 0),
		Errors:       make([]string, 0),
	}

	for idx, importedTx := range importedTxs {
		// Update progress
		progressCallback(idx+1, total)

		// Check for duplicates
		if existingIDs[importedTx.ExternalID] {
			result.Duplicates++
			s.logger.Debug().
				Str("external_id", importedTx.ExternalID).
				Str("description", importedTx.Description).
				Msg("skipping duplicate transaction")
			continue
		}

		// Convert TransactionImport to Transaction
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

		// Auto-categorize if available
		if s.categorizationSvc != nil {
			suggestion, err := s.categorizationSvc.SuggestCategory(ctx, userID, tx.Description)
			if err == nil && len(suggestion.Suggestions) > 0 {
				bestSuggestion := suggestion.Suggestions[0]
				if bestSuggestion.Confidence >= 1 {
					tx.CategoryID = &bestSuggestion.CategoryID
					s.logger.Debug().
						Str("description", tx.Description).
						Str("category", bestSuggestion.CategoryName).
						Int("confidence", bestSuggestion.Confidence).
						Msg("auto-categorized imported transaction")
				}
			}
		}

		// Create transaction
		if err := s.txRepo.Create(ctx, tx); err != nil {
			errMsg := fmt.Sprintf("Failed to import transaction '%s': %v", tx.Description, err)
			result.Errors = append(result.Errors, errMsg)
			s.logger.Error().Err(err).
				Str("description", tx.Description).
				Msg("failed to create imported transaction")
			continue
		}

		// Update account balance
		delta := tx.Amount
		if tx.Type == model.TransactionTypeExpense {
			delta = delta.Neg()
		}
		if err := s.accountRepo.AdjustBalance(ctx, account.ID, delta); err != nil {
			s.logger.Error().Err(err).
				Str("transaction_id", tx.ID.String()).
				Msg("failed to update account balance for imported transaction")
		}

		result.TotalImported++
		result.Transactions = append(result.Transactions, tx)

		// Learn categorization pattern if categorized
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
		Int("errors", len(result.Errors)).
		Msg("completed transaction import with progress")

	return result, nil
}
