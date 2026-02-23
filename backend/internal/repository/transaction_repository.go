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
	"github.com/gustavoz65/Cashing-go/internal/database"
	"github.com/gustavoz65/Cashing-go/internal/model"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
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
	// Build WHERE clause
	conditions := []string{"user_id = ?"}
	args := []interface{}{filter.UserID.String()}

	if filter.AccountID != nil {
		conditions = append(conditions, "bank_account_id = ?")
		args = append(args, filter.AccountID.String())
	}

	if filter.CategoryID != nil {
		conditions = append(conditions, "category_id = ?")
		args = append(args, filter.CategoryID.String())
	}

	if filter.Type != nil {
		conditions = append(conditions, "type = ?")
		args = append(args, *filter.Type)
	}

	if filter.Source != nil {
		conditions = append(conditions, "source = ?")
		args = append(args, *filter.Source)
	}

	if filter.StartDate != nil {
		conditions = append(conditions, "transaction_date >= ?")
		args = append(args, *filter.StartDate)
	}

	if filter.EndDate != nil {
		conditions = append(conditions, "transaction_date <= ?")
		args = append(args, *filter.EndDate)
	}

	if filter.IsPaid != nil {
		conditions = append(conditions, "is_paid = ?")
		args = append(args, *filter.IsPaid)
	}

	if filter.IsRecurring != nil {
		conditions = append(conditions, "is_recurring = ?")
		args = append(args, *filter.IsRecurring)
	}

	if filter.MinAmount != nil {
		conditions = append(conditions, "amount >= ?")
		args = append(args, filter.MinAmount.String())
	}

	if filter.MaxAmount != nil {
		conditions = append(conditions, "amount <= ?")
		args = append(args, filter.MaxAmount.String())
	}

	if filter.SearchTerm != "" {
		conditions = append(conditions, "(description LIKE ? OR notes LIKE ?)")
		searchTerm := "%" + filter.SearchTerm + "%"
		args = append(args, searchTerm, searchTerm)
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions WHERE %s", whereClause)
	var total int64
	err := r.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count transactions: %w", err)
	}

	// Pagination
	pagination := PaginationParams{
		Page:     filter.Page,
		PageSize: filter.PageSize,
		SortBy:   filter.SortBy,
		SortDir:  filter.SortDirection,
	}
	pagination.Validate([]string{"transaction_date", "amount", "created_at", "description"})

	// Get data
	query := fmt.Sprintf(`
		SELECT id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		FROM transactions
		WHERE %s
		%s
		%s
	`, whereClause, pagination.BuildOrderClause(), pagination.BuildLimitClause())

	rows, err := r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()

	transactions, err := r.scanTransactions(rows)
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

// GetByDateRange retrieves transactions within a date range
func (r *TransactionRepository) GetByDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		FROM transactions
		WHERE user_id = ? AND transaction_date BETWEEN ? AND ?
		ORDER BY transaction_date DESC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by date range: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetUpcomingBills retrieves unpaid transactions with due dates coming up
func (r *TransactionRepository) GetUpcomingBills(ctx context.Context, userID uuid.UUID, days int) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		FROM transactions
		WHERE user_id = ? AND is_paid = FALSE AND type = 'expense'
			AND due_date IS NOT NULL AND due_date <= DATE_ADD(CURRENT_DATE, INTERVAL ? DAY)
		ORDER BY due_date ASC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), days)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming bills: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
}

// GetRecentTransactions retrieves the most recent transactions
func (r *TransactionRepository) GetRecentTransactions(ctx context.Context, userID uuid.UUID, limit int) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, bank_account_id, category_id, type, amount,
			description, notes, source, transaction_date, due_date, payment_date,
			is_paid, auto_pay, is_recurring, recurring_id, installment_number,
			total_installments, installment_group_id, tags, attachment_url,
			external_id, created_at, updated_at
		FROM transactions
		WHERE user_id = ?
		ORDER BY transaction_date DESC, created_at DESC
		LIMIT ?
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}
	defer rows.Close()

	return r.scanTransactions(rows)
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

// GetSumByType returns the sum of transactions by type for a date range
func (r *TransactionRepository) GetSumByType(ctx context.Context, userID uuid.UUID, txType model.TransactionType, startDate, endDate time.Time) (decimal.Decimal, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE user_id = ? AND type = ? AND is_paid = TRUE
			AND transaction_date BETWEEN ? AND ?
	`

	var sum string
	err := r.QueryRowContext(ctx, query, userID.String(), txType, startDate, endDate).Scan(&sum)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get sum by type: %w", err)
	}

	total, _ := decimal.NewFromString(sum)
	return total, nil
}

// GetSumByCategory returns the sum of transactions by category
func (r *TransactionRepository) GetSumByCategory(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]model.CategoryAmount, error) {
	query := `
		SELECT
			t.category_id,
			COALESCE(c.name, 'Sem categoria') as category_name,
			SUM(t.amount) as amount,
			COUNT(*) as count
		FROM transactions t
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
