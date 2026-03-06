package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/database"
	"github.com/gustavoz65/Fyapp/internal/model"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
)

type TransactionRepository struct {
	*BaseRepository
}

func NewTransactionRepository(db *database.Database, logger *zerolog.Logger) *TransactionRepository {
	return &TransactionRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	tx.ID = uuid.New()
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()

	query := `
		INSERT INTO transactions (
			id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var categoryID sql.NullString
	if tx.CategoryID != nil {
		categoryID = sql.NullString{String: tx.CategoryID.String(), Valid: true}
	}

	var recurringID, installmentGroupID sql.NullString
	if tx.RecurringID != nil {
		recurringID = sql.NullString{String: tx.RecurringID.String(), Valid: true}
	}
	if tx.InstallmentGroupID != nil {
		installmentGroupID = sql.NullString{String: tx.InstallmentGroupID.String(), Valid: true}
	}

	var dueDate, paymentDate sql.NullTime
	if tx.DueDate != nil {
		dueDate = sql.NullTime{Time: *tx.DueDate, Valid: true}
	}
	if tx.PaymentDate != nil {
		paymentDate = sql.NullTime{Time: *tx.PaymentDate, Valid: true}
	}

	var tagsJSON []byte
	if len(tx.Tags) > 0 {
		tagsJSON, _ = json.Marshal(tx.Tags)
	}

	_, err := r.ExecContext(ctx, query,
		tx.ID.String(),
		tx.UserID.String(),
		tx.BankAccountID.String(),
		categoryID,
		tx.Type,
		tx.Amount.String(),
		tx.Description,
		NullString(tx.Notes),
		tx.Source,
		tx.TransactionDate,
		dueDate,
		paymentDate,
		tx.IsPaid,
		tx.AutoPay,
		tx.IsRecurring,
		recurringID,
		NullInt32(tx.InstallmentNumber),
		NullInt32(tx.TotalInstallments),
		installmentGroupID,
		tagsJSON,
		NullString(tx.AttachmentURL),
		NullString(tx.ExternalID),
		tx.CreatedAt,
		tx.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	query := `
		SELECT id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		FROM transactions
		WHERE id = ?
	`

	return r.scanTransaction(r.QueryRowContext(ctx, query, id.String()))
}

// GetByIDAndUser retrieves a transaction by ID ensuring it belongs to the user
func (r *TransactionRepository) GetByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*model.Transaction, error) {
	query := `
		SELECT id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		FROM transactions
		WHERE id = ? AND user_id = ?
	`

	return r.scanTransaction(r.QueryRowContext(ctx, query, id.String(), userID.String()))
}

// GetByFilter retrieves transactions based on filter criteria
func (r *TransactionRepository) GetByFilter(ctx context.Context, filter *model.TransactionFilter) ([]*model.Transaction, int64, error) {
	// Build WHERE clause (prefixed with t. for JOIN compatibility)
	conditions := []string{"t.user_id = ?"}
	args := []interface{}{filter.UserID.String()}

	if filter.AccountID != nil {
		conditions = append(conditions, "t.bank_account_id = ?")
		args = append(args, filter.AccountID.String())
	}

	if filter.CategoryID != nil {
		conditions = append(conditions, "t.category_id = ?")
		args = append(args, filter.CategoryID.String())
	}

	if filter.Type != nil {
		conditions = append(conditions, "t.type = ?")
		args = append(args, *filter.Type)
	}

	if filter.Source != nil {
		conditions = append(conditions, "t.source = ?")
		args = append(args, *filter.Source)
	}

	if filter.StartDate != nil {
		conditions = append(conditions, "t.transaction_date >= ?")
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		conditions = append(conditions, "t.transaction_date <= ?")
		args = append(args, *filter.EndDate)
	}

	if filter.IsPaid != nil {
		conditions = append(conditions, "t.is_paid = ?")
		args = append(args, *filter.IsPaid)
	}

	if filter.IsRecurring != nil {
		conditions = append(conditions, "t.is_recurring = ?")
		args = append(args, *filter.IsRecurring)
	}

	if filter.MinAmount != nil {
		conditions = append(conditions, "t.amount >= ?")
		args = append(args, filter.MinAmount.String())
	}

	if filter.MaxAmount != nil {
		conditions = append(conditions, "t.amount <= ?")
		args = append(args, filter.MaxAmount.String())
	}

	if filter.SearchTerm != "" {
		conditions = append(conditions, "(t.description LIKE ? OR t.notes LIKE ?)")
		searchTerm := "%" + filter.SearchTerm + "%"
		args = append(args, searchTerm, searchTerm)
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total (only active accounts)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions t INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE WHERE %s", whereClause)
	var total int64
	err := r.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Pagination (sort by t. prefix for ambiguous fields)
	pagination := PaginationParams{
		Page:     filter.Page,
		PageSize: filter.PageSize,
		SortBy:   filter.SortBy,
		SortDir:  filter.SortDirection,
	}
	pagination.Validate([]string{"transaction_date", "amount", "created_at", "description"})

	// Get data with LEFT JOINs to populate category and filter by active accounts
	query := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.bank_account_id, t.category_id, t.type, t.amount,
			t.description, t.notes, t.source, t.transaction_date, t.due_date, t.payment_date,
			t.is_paid, t.auto_pay, t.is_recurring, t.recurring_id, t.installment_number,
			t.total_installments, t.installment_group_id, t.tags, t.attachment_url,
			t.external_id, t.created_at, t.updated_at,
			c.id, c.name, c.type, c.color, c.icon, c.is_system, c.is_active
		FROM transactions t
		INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE %s
		%s
		%s
	`, whereClause, pagination.BuildOrderClause(), pagination.BuildLimitClause())

	rows, err := r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions, err := r.scanTransactionsWithCategory(rows)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetByDateRange retrieves transactions within a date range (only from active accounts)
func (r *TransactionRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*model.Transaction, error) {
	query := `
		SELECT t.id, t.user_id, t.bank_account_id, t.category_id, t.type, t.amount,
			t.description, t.notes, t.source, t.transaction_date, t.due_date, t.payment_date,
			t.is_paid, t.auto_pay, t.is_recurring, t.recurring_id, t.installment_number,
			t.total_installments, t.installment_group_id, t.tags, t.attachment_url,
			t.external_id, t.created_at, t.updated_at
		FROM transactions t
		INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE
		WHERE t.user_id = ? AND t.transaction_date BETWEEN ? AND ?
		ORDER BY t.transaction_date DESC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetUpcomingBills retrieves unpaid transactions with due dates coming up (only from active accounts)
func (r *TransactionRepository) GetUpcomingBills(ctx context.Context, userID uuid.UUID, days int) ([]*model.Transaction, error) {
	query := `
		SELECT t.id, t.user_id, t.bank_account_id, t.category_id, t.type, t.amount,
			t.description, t.notes, t.source, t.transaction_date, t.due_date, t.payment_date,
			t.is_paid, t.auto_pay, t.is_recurring, t.recurring_id, t.installment_number,
			t.total_installments, t.installment_group_id, t.tags, t.attachment_url,
			t.external_id, t.created_at, t.updated_at
		FROM transactions t
		INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE
		WHERE t.user_id = ? AND t.is_paid = FALSE AND t.type = 'expense'
			AND t.due_date IS NOT NULL
			AND t.due_date >= CURRENT_DATE
			AND t.due_date <= DATE_ADD(CURRENT_DATE, INTERVAL ? DAY)
		ORDER BY t.due_date ASC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), days)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming bills: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetRecentTransactions retrieves the most recent transactions (only from active accounts)
func (r *TransactionRepository) GetRecentTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]*model.Transaction, error) {
	query := `
		SELECT t.id, t.user_id, t.bank_account_id, t.category_id, t.type, t.amount,
			t.description, t.notes, t.source, t.transaction_date, t.due_date, t.payment_date,
			t.is_paid, t.auto_pay, t.is_recurring, t.recurring_id, t.installment_number,
			t.total_installments, t.installment_group_id, t.tags, t.attachment_url,
			t.external_id, t.created_at, t.updated_at
		FROM transactions t
		INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE
		WHERE t.user_id = ?
		ORDER BY t.transaction_date DESC, t.created_at DESC
		LIMIT ?
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetRecentByPeriod retrieves the most recent transactions within a date range, with category info
func (r *TransactionRepository) GetRecentByPeriod(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time, limit int) ([]*model.Transaction, error) {
	query := `
		SELECT t.id, t.user_id, t.bank_account_id, t.category_id, t.type, t.amount,
			t.description, t.notes, t.source, t.transaction_date, t.due_date, t.payment_date,
			t.is_paid, t.auto_pay, t.is_recurring, t.recurring_id, t.installment_number,
			t.total_installments, t.installment_group_id, t.tags, t.attachment_url,
			t.external_id, t.created_at, t.updated_at,
			c.id, c.name, c.type, c.color, c.icon, c.is_system, c.is_active
		FROM transactions t
		INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? AND t.transaction_date BETWEEN ? AND ?
		ORDER BY t.transaction_date DESC, t.created_at DESC
		LIMIT ?
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions by period: %w", err)
	}
	defer rows.Close()

	return r.scanTransactionsWithCategory(rows)
}

// Update updates a transaction
func (r *TransactionRepository) Update(ctx context.Context, tx *model.Transaction) error {
	tx.UpdatedAt = time.Now()

	query := `
		UPDATE transactions SET
			category_id = ?, amount = ?, description = ?, notes = ?,
			transaction_date = ?, due_date = ?, payment_date = ?,
			is_paid = ?, tags = ?, attachment_url = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	var categoryID sql.NullString
	if tx.CategoryID != nil {
		categoryID = sql.NullString{String: tx.CategoryID.String(), Valid: true}
	}

	var dueDate, paymentDate sql.NullTime
	if tx.DueDate != nil {
		dueDate = sql.NullTime{Time: *tx.DueDate, Valid: true}
	}
	if tx.PaymentDate != nil {
		paymentDate = sql.NullTime{Time: *tx.PaymentDate, Valid: true}
	}

	var tagsJSON []byte
	if len(tx.Tags) > 0 {
		tagsJSON, _ = json.Marshal(tx.Tags)
	}

	result, err := r.ExecContext(ctx, query,
		categoryID,
		tx.Amount.String(),
		tx.Description,
		NullString(tx.Notes),
		tx.TransactionDate,
		dueDate,
		paymentDate,
		tx.IsPaid,
		tagsJSON,
		NullString(tx.AttachmentURL),
		tx.UpdatedAt,
		tx.ID.String(),
		tx.UserID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

// MarkAsPaid marks a transaction as paid
func (r *TransactionRepository) MarkAsPaid(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		UPDATE transactions
		SET is_paid = TRUE, payment_date = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	now := time.Now()
	result, err := r.ExecContext(ctx, query, now, now, id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("failed to mark transaction as paid: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

// Delete deletes a transaction
func (r *TransactionRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM transactions WHERE id = ? AND user_id = ?`

	result, err := r.ExecContext(ctx, query, id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

// DeleteTx deleta uma transação dentro de uma transação de banco de dados existente
func (r *TransactionRepository) DeleteTx(ctx context.Context, dbTx *sql.Tx, id, userID uuid.UUID) error {
	query := `DELETE FROM transactions WHERE id = ? AND user_id = ?`

	result, err := dbTx.ExecContext(ctx, query, id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("failed to delete transaction in tx: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

// GetNetBalanceForAccount returns the net balance from all paid transactions for an account
// net = SUM(income) - SUM(expense)
func (r *TransactionRepository) GetNetBalanceForAccount(ctx context.Context, accountID uuid.UUID) (decimal.Decimal, error) {
	query := `
		SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE -amount END), 0)
		FROM transactions
		WHERE bank_account_id = ? AND is_paid = TRUE
	`

	var net string
	err := r.QueryRowContext(ctx, query, accountID.String()).Scan(&net)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get net balance for account: %w", err)
	}

	total, _ := decimal.NewFromString(net)
	return total, nil
}

// GetSumByType returns the sum of transactions by type for a date range (only from active accounts)
func (r *TransactionRepository) GetSumByType(ctx context.Context, userID uuid.UUID, txType model.TransactionType, startDate, endDate time.Time) (decimal.Decimal, error) {
	query := `
		SELECT COALESCE(SUM(t.amount), 0)
		FROM transactions t
		INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE AND ba.include_in_total = TRUE
		WHERE t.user_id = ? AND t.type = ? AND t.is_paid = TRUE
			AND t.transaction_date BETWEEN ? AND ?
	`

	var sum string
	err := r.QueryRowContext(ctx, query, userID.String(), txType, startDate, endDate).Scan(&sum)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get sum by type: %w", err)
	}

	total, _ := decimal.NewFromString(sum)
	return total, nil
}

// GetSumByCategory returns the sum of transactions by category (only from active accounts)
func (r *TransactionRepository) GetSumByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]model.CategoryAmount, error) {
	query := `
		SELECT
			t.category_id,
			COALESCE(c.name, 'Sem categoria') as category_name,
			SUM(t.amount) as amount,
			COUNT(*) as count
		FROM transactions t
		INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE
		LEFT JOIN categories c ON t.category_id = c.id
		WHERE t.user_id = ? AND t.type = 'expense' AND t.is_paid = TRUE
			AND t.transaction_date BETWEEN ? AND ?
		GROUP BY t.category_id, c.name
		ORDER BY amount DESC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get sum by category: %w", err)
	}
	defer rows.Close()

	var results []model.CategoryAmount
	var totalAmount decimal.Decimal

	for rows.Next() {
		var ca model.CategoryAmount
		var categoryID sql.NullString
		var amountStr string

		err := rows.Scan(&categoryID, &ca.CategoryName, &amountStr, &ca.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category amount: %w", err)
		}

		if categoryID.Valid {
			id, _ := uuid.Parse(categoryID.String)
			ca.CategoryID = &id
		}

		ca.Amount, _ = decimal.NewFromString(amountStr)
		totalAmount = totalAmount.Add(ca.Amount)
		results = append(results, ca)
	}

	// Calculate percentages
	for i := range results {
		if !totalAmount.IsZero() {
			results[i].Percentage = results[i].Amount.Div(totalAmount).Mul(decimal.NewFromInt(100))
		}
	}

	return results, nil
}

// Helper functions

// populateTxFields assigns nullable/converted fields onto tx after a successful Scan call.
func populateTxFields(
	tx *model.Transaction,
	amount string,
	categoryID, recurringID, installmentGroupID, notes, attachmentURL, externalID sql.NullString,
	dueDate, paymentDate sql.NullTime,
	installmentNumber, totalInstallments sql.NullInt32,
	tags []byte,
) {
	tx.Amount, _ = decimal.NewFromString(amount)

	if categoryID.Valid {
		id, _ := uuid.Parse(categoryID.String)
		tx.CategoryID = &id
	}
	if recurringID.Valid {
		id, _ := uuid.Parse(recurringID.String)
		tx.RecurringID = &id
	}
	if installmentGroupID.Valid {
		id, _ := uuid.Parse(installmentGroupID.String)
		tx.InstallmentGroupID = &id
	}

	tx.Notes = StringPtr(notes)
	tx.AttachmentURL = StringPtr(attachmentURL)
	tx.ExternalID = StringPtr(externalID)
	tx.InstallmentNumber = IntPtr(installmentNumber)
	tx.TotalInstallments = IntPtr(totalInstallments)

	if dueDate.Valid {
		tx.DueDate = &dueDate.Time
	}
	if paymentDate.Valid {
		tx.PaymentDate = &paymentDate.Time
	}

	if len(tags) > 0 {
		_ = json.Unmarshal(tags, &tx.Tags)
	}
}

func (r *TransactionRepository) scanTransaction(row *sql.Row) (*model.Transaction, error) {
	tx := &model.Transaction{}
	var categoryID, recurringID, installmentGroupID, notes, attachmentURL, externalID sql.NullString
	var dueDate, paymentDate sql.NullTime
	var installmentNumber, totalInstallments sql.NullInt32
	var amount string
	var tags []byte

	err := row.Scan(
		&tx.ID,
		&tx.UserID,
		&tx.BankAccountID,
		&categoryID,
		&tx.Type,
		&amount,
		&tx.Description,
		&notes,
		&tx.Source,
		&tx.TransactionDate,
		&dueDate,
		&paymentDate,
		&tx.IsPaid,
		&tx.AutoPay,
		&tx.IsRecurring,
		&recurringID,
		&installmentNumber,
		&totalInstallments,
		&installmentGroupID,
		&tags,
		&attachmentURL,
		&externalID,
		&tx.CreatedAt,
		&tx.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("failed to scan transaction: %w", err)
	}

	populateTxFields(tx, amount, categoryID, recurringID, installmentGroupID, notes, attachmentURL, externalID, dueDate, paymentDate, installmentNumber, totalInstallments, tags)
	return tx, nil
}

func (r *TransactionRepository) scanTransactions(rows *sql.Rows) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	for rows.Next() {
		tx := &model.Transaction{}
		var categoryID, recurringID, installmentGroupID, notes, attachmentURL, externalID sql.NullString
		var dueDate, paymentDate sql.NullTime
		var installmentNumber, totalInstallments sql.NullInt32
		var amount string
		var tags []byte

		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.BankAccountID,
			&categoryID,
			&tx.Type,
			&amount,
			&tx.Description,
			&notes,
			&tx.Source,
			&tx.TransactionDate,
			&dueDate,
			&paymentDate,
			&tx.IsPaid,
			&tx.AutoPay,
			&tx.IsRecurring,
			&recurringID,
			&installmentNumber,
			&totalInstallments,
			&installmentGroupID,
			&tags,
			&attachmentURL,
			&externalID,
			&tx.CreatedAt,
			&tx.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		populateTxFields(tx, amount, categoryID, recurringID, installmentGroupID, notes, attachmentURL, externalID, dueDate, paymentDate, installmentNumber, totalInstallments, tags)
		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// scanTransactionsWithCategory scans transaction rows that include a LEFT JOIN with categories
func (r *TransactionRepository) scanTransactionsWithCategory(rows *sql.Rows) ([]*model.Transaction, error) {
	var transactions []*model.Transaction

	for rows.Next() {
		tx := &model.Transaction{}
		var categoryID, recurringID, installmentGroupID, notes, attachmentURL, externalID sql.NullString
		var dueDate, paymentDate sql.NullTime
		var installmentNumber, totalInstallments sql.NullInt32
		var amount string
		var tags []byte

		// Category columns (nullable from LEFT JOIN)
		var catID, catName, catType, catColor, catIcon sql.NullString
		var catIsSystem, catIsActive sql.NullBool

		err := rows.Scan(
			&tx.ID,
			&tx.UserID,
			&tx.BankAccountID,
			&categoryID,
			&tx.Type,
			&amount,
			&tx.Description,
			&notes,
			&tx.Source,
			&tx.TransactionDate,
			&dueDate,
			&paymentDate,
			&tx.IsPaid,
			&tx.AutoPay,
			&tx.IsRecurring,
			&recurringID,
			&installmentNumber,
			&totalInstallments,
			&installmentGroupID,
			&tags,
			&attachmentURL,
			&externalID,
			&tx.CreatedAt,
			&tx.UpdatedAt,
			// Category fields
			&catID,
			&catName,
			&catType,
			&catColor,
			&catIcon,
			&catIsSystem,
			&catIsActive,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		populateTxFields(tx, amount, categoryID, recurringID, installmentGroupID, notes, attachmentURL, externalID, dueDate, paymentDate, installmentNumber, totalInstallments, tags)

		// Populate Category if JOIN returned data
		if catID.Valid && catName.Valid {
			catUUID, _ := uuid.Parse(catID.String)
			tx.Category = &model.Category{
				ID:       catUUID,
				Name:     catName.String,
				Type:     model.CategoryType(catType.String),
				Color:    catColor.String,
				Icon:     catIcon.String,
				IsSystem: catIsSystem.Bool,
				IsActive: catIsActive.Bool,
			}
		}

		transactions = append(transactions, tx)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// GetDueForAutoReconcile busca transações com auto_pay=true que estão vencidas e não pagas
func (r *TransactionRepository) GetDueForAutoReconcile(ctx context.Context, today time.Time) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		FROM transactions
		WHERE is_paid = FALSE AND auto_pay = TRUE AND due_date IS NOT NULL AND due_date <= ?
		ORDER BY due_date ASC
	`

	rows, err := r.QueryContext(ctx, query, today)
	if err != nil {
		return nil, fmt.Errorf("failed to get due transactions for auto reconcile: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// CountTransactionsByUserToday conta o número de transações criadas pelo usuário hoje
func (r *TransactionRepository) CountTransactionsByUserToday(ctx context.Context, userID uuid.UUID) (int64, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	query := `
		SELECT COUNT(*)
		FROM transactions
		WHERE user_id = ? AND created_at >= ? AND created_at < ?
	`

	var count int64
	err := r.QueryRowContext(ctx, query, userID.String(), today, tomorrow).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	return count, nil
}

// MonthlyRawData contém dados brutos de receita/despesa agregados por mês vindo do banco
type MonthlyRawData struct {
	Year    int
	Month   int
	Income  decimal.Decimal
	Expense decimal.Decimal
}

// GetMonthlyBreakdown retorna receitas e despesas agrupadas por mês em um intervalo de datas.
// Retorna apenas meses com transações; o chamador deve preencher meses vazios se necessário.
func (r *TransactionRepository) GetMonthlyBreakdown(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]MonthlyRawData, error) {
	query := `
		SELECT
			YEAR(t.transaction_date) AS year,
			MONTH(t.transaction_date) AS month,
			COALESCE(SUM(CASE WHEN t.type = 'income' THEN t.amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN t.type = 'expense' THEN t.amount ELSE 0 END), 0) AS expense
		FROM transactions t
		INNER JOIN bank_accounts ba ON t.bank_account_id = ba.id AND ba.is_active = TRUE AND ba.include_in_total = TRUE
		WHERE t.user_id = ? AND t.is_paid = TRUE
			AND t.transaction_date BETWEEN ? AND ?
		GROUP BY YEAR(t.transaction_date), MONTH(t.transaction_date)
		ORDER BY year ASC, month ASC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly breakdown: %w", err)
	}
	defer rows.Close()

	var results []MonthlyRawData
	for rows.Next() {
		var d MonthlyRawData
		var incomeStr, expenseStr string
		if err := rows.Scan(&d.Year, &d.Month, &incomeStr, &expenseStr); err != nil {
			return nil, fmt.Errorf("failed to scan monthly breakdown: %w", err)
		}
		d.Income, _ = decimal.NewFromString(incomeStr)
		d.Expense, _ = decimal.NewFromString(expenseStr)
		results = append(results, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating monthly breakdown: %w", err)
	}

	return results, nil
}

// CreateTx cria uma transação dentro de uma transação de banco de dados existente
func (r *TransactionRepository) CreateTx(ctx context.Context, dbTx *sql.Tx, tx *model.Transaction) error {
	tx.ID = uuid.New()
	tx.CreatedAt = time.Now()
	tx.UpdatedAt = time.Now()

	query := `
		INSERT INTO transactions (
			id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var categoryID sql.NullString
	if tx.CategoryID != nil {
		categoryID = sql.NullString{String: tx.CategoryID.String(), Valid: true}
	}

	var recurringID, installmentGroupID sql.NullString
	if tx.RecurringID != nil {
		recurringID = sql.NullString{String: tx.RecurringID.String(), Valid: true}
	}
	if tx.InstallmentGroupID != nil {
		installmentGroupID = sql.NullString{String: tx.InstallmentGroupID.String(), Valid: true}
	}

	var dueDate, paymentDate sql.NullTime
	if tx.DueDate != nil {
		dueDate = sql.NullTime{Time: *tx.DueDate, Valid: true}
	}
	if tx.PaymentDate != nil {
		paymentDate = sql.NullTime{Time: *tx.PaymentDate, Valid: true}
	}

	var tagsJSON []byte
	if len(tx.Tags) > 0 {
		tagsJSON, _ = json.Marshal(tx.Tags)
	}

	_, err := dbTx.ExecContext(ctx, query,
		tx.ID.String(),
		tx.UserID.String(),
		tx.BankAccountID.String(),
		categoryID,
		tx.Type,
		tx.Amount.String(),
		tx.Description,
		NullString(tx.Notes),
		tx.Source,
		tx.TransactionDate,
		dueDate,
		paymentDate,
		tx.IsPaid,
		tx.AutoPay,
		tx.IsRecurring,
		recurringID,
		NullInt32(tx.InstallmentNumber),
		NullInt32(tx.TotalInstallments),
		installmentGroupID,
		tagsJSON,
		NullString(tx.AttachmentURL),
		NullString(tx.ExternalID),
		tx.CreatedAt,
		tx.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create transaction in tx: %w", err)
	}

	return nil
}

// GetPaidBalanceDeltaByAccountTx calcula o delta de saldo necessário para reverter todas as
// transações pagas de uma conta, dentro de uma transação de banco de dados existente.
// O valor retornado, quando aplicado ao saldo da conta, desfaz o impacto de todas as transações pagas.
func (r *TransactionRepository) GetPaidBalanceDeltaByAccountTx(ctx context.Context, dbTx *sql.Tx, accountID, userID uuid.UUID) (decimal.Decimal, error) {
	query := `
		SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN -amount ELSE amount END), 0)
		FROM transactions
		WHERE bank_account_id = ? AND user_id = ? AND is_paid = TRUE
	`

	var deltaStr string
	if err := dbTx.QueryRowContext(ctx, query, accountID.String(), userID.String()).Scan(&deltaStr); err != nil {
		return decimal.Zero, fmt.Errorf("failed to get paid balance delta: %w", err)
	}

	delta, _ := decimal.NewFromString(deltaStr)
	return delta, nil
}

// DeleteAllByAccountTx deleta todas as transações de uma conta dentro de uma transação de banco de dados existente
func (r *TransactionRepository) DeleteAllByAccountTx(ctx context.Context, dbTx *sql.Tx, accountID, userID uuid.UUID) (int64, error) {
	query := `DELETE FROM transactions WHERE bank_account_id = ? AND user_id = ?`

	result, err := dbTx.ExecContext(ctx, query, accountID.String(), userID.String())
	if err != nil {
		return 0, fmt.Errorf("failed to delete all transactions for account: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

// GetPaidBalanceDeltasByDateRangeTx retorna o delta de saldo por conta para reverter transações pagas
// de source='manual' em um intervalo de datas, dentro de uma transação de banco de dados existente.
func (r *TransactionRepository) GetPaidBalanceDeltasByDateRangeTx(ctx context.Context, dbTx *sql.Tx, userID uuid.UUID, startDate, endDate time.Time, accountID *uuid.UUID) (map[uuid.UUID]decimal.Decimal, error) {
	query := `
		SELECT bank_account_id,
			COALESCE(SUM(CASE WHEN type = 'income' THEN -amount ELSE amount END), 0) AS delta
		FROM transactions
		WHERE user_id = ? AND is_paid = TRUE AND source = 'manual'
			AND transaction_date BETWEEN ? AND ?
	`
	args := []interface{}{userID.String(), startDate, endDate}

	if accountID != nil {
		query += ` AND bank_account_id = ?`
		args = append(args, accountID.String())
	}

	query += ` GROUP BY bank_account_id`

	rows, err := dbTx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get paid balance deltas: %w", err)
	}
	defer rows.Close()

	result := make(map[uuid.UUID]decimal.Decimal)
	for rows.Next() {
		var accountIDStr, deltaStr string
		if err := rows.Scan(&accountIDStr, &deltaStr); err != nil {
			return nil, fmt.Errorf("failed to scan balance delta: %w", err)
		}
		accID, _ := uuid.Parse(accountIDStr)
		delta, _ := decimal.NewFromString(deltaStr)
		result[accID] = delta
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating balance deltas: %w", err)
	}

	return result, nil
}

// DeleteByDateRangeTx deleta transações de source='manual' em um intervalo de datas dentro de uma tx.
// Apenas transações manuais são removidas para manter consistência com a regra CanBeDeleted().
func (r *TransactionRepository) DeleteByDateRangeTx(ctx context.Context, dbTx *sql.Tx, userID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	query := `DELETE FROM transactions WHERE user_id = ? AND source = 'manual' AND transaction_date BETWEEN ? AND ?`

	result, err := dbTx.ExecContext(ctx, query, userID.String(), startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("failed to delete transactions by date range in tx: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

// DeleteByAccountAndDateRangeTx deleta transações de source='manual' de uma conta em um intervalo de datas dentro de uma tx.
// Apenas transações manuais são removidas para manter consistência com a regra CanBeDeleted().
func (r *TransactionRepository) DeleteByAccountAndDateRangeTx(ctx context.Context, dbTx *sql.Tx, userID, accountID uuid.UUID, startDate, endDate time.Time) (int64, error) {
	query := `DELETE FROM transactions WHERE user_id = ? AND bank_account_id = ? AND source = 'manual' AND transaction_date BETWEEN ? AND ?`

	result, err := dbTx.ExecContext(ctx, query, userID.String(), accountID.String(), startDate, endDate)
	if err != nil {
		return 0, fmt.Errorf("failed to delete transactions by account and date range in tx: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

// GetExistingExternalIDs returns a map of existing external IDs for the user
func (r *TransactionRepository) GetExistingExternalIDs(ctx context.Context, userID uuid.UUID, externalIDs []string) (map[string]bool, error) {
	if len(externalIDs) == 0 {
		return make(map[string]bool), nil
	}

	// Create placeholders for IN clause
	placeholders := make([]string, len(externalIDs))
	args := make([]interface{}, 0, len(externalIDs)+1)
	args = append(args, userID.String())

	for i, id := range externalIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query := fmt.Sprintf(`
		SELECT external_id
		FROM transactions
		WHERE user_id = ? AND external_id IN (%s)
	`, strings.Join(placeholders, ","))

	rows, err := r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing external IDs: %w", err)
	}
	defer rows.Close()

	existingIDs := make(map[string]bool)
	for rows.Next() {
		var externalID sql.NullString
		if err := rows.Scan(&externalID); err != nil {
			return nil, fmt.Errorf("failed to scan external ID: %w", err)
		}
		if externalID.Valid {
			existingIDs[externalID.String] = true
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating external IDs: %w", err)
	}

	return existingIDs, nil
}
