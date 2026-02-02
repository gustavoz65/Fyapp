package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoz65/Cashing-go/internal/database"
	"github.com/gustavoz65/Cashing-go/internal/model"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

var (
	ErrBankAccountNotFound = errors.New("bank account not found")
)

type BankAccountRepository struct {
	*BaseRepository
}

func NewBankAccountRepository(db *database.Database, logger *zerolog.Logger) *BankAccountRepository {
	return &BankAccountRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new bank account
func (r *BankAccountRepository) Create(ctx context.Context, account *model.BankAccount) error {
	account.ID = uuid.New()
	account.CurrentBalance = account.InitialBalance
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	query := `
		INSERT INTO bank_accounts (
			id, user_id, name, bank_name, bank_code, account_type,
			account_number, agency, initial_balance, current_balance,
			credit_limit, closing_day, due_day, currency, color, icon,
			is_active, include_in_total, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var creditLimit sql.NullString
	if account.CreditLimit != nil {
		creditLimit = sql.NullString{String: account.CreditLimit.String(), Valid: true}
	}

	_, err := r.ExecContext(ctx, query,
		account.ID.String(),
		account.UserID.String(),
		account.Name,
		NullString(account.BankName),
		NullString(account.BankCode),
		account.AccountType,
		NullString(account.AccountNumber),
		NullString(account.Agency),
		account.InitialBalance.String(),
		account.CurrentBalance.String(),
		creditLimit,
		NullInt32(account.ClosingDay),
		NullInt32(account.DueDay),
		account.Currency,
		account.Color,
		account.Icon,
		account.IsActive,
		account.IncludeInTotal,
		account.CreatedAt,
		account.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create bank account: %w", err)
	}

	return nil
}

// GetByID retrieves a bank account by ID
func (r *BankAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.BankAccount, error) {
	query := `
		SELECT id, user_id, name, bank_name, bank_code, account_type,
			account_number, agency, initial_balance, current_balance,
			credit_limit, closing_day, due_day, currency, color, icon,
			is_active, include_in_total, last_sync_at, external_id,
			created_at, updated_at
		FROM bank_accounts
		WHERE id = ?
	`

	return r.scanBankAccount(r.QueryRowContext(ctx, query, id.String()))
}

// GetByIDAndUser retrieves a bank account by ID ensuring it belongs to the user
func (r *BankAccountRepository) GetByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*model.BankAccount, error) {
	query := `
		SELECT id, user_id, name, bank_name, bank_code, account_type,
			account_number, agency, initial_balance, current_balance,
			credit_limit, closing_day, due_day, currency, color, icon,
			is_active, include_in_total, last_sync_at, external_id,
			created_at, updated_at
		FROM bank_accounts
		WHERE id = ? AND user_id = ?
	`

	return r.scanBankAccount(r.QueryRowContext(ctx, query, id.String(), userID.String()))
}

// GetAllByUser retrieves all bank accounts for a user
func (r *BankAccountRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*model.BankAccount, error) {
	query := `
		SELECT id, user_id, name, bank_name, bank_code, account_type,
			account_number, agency, initial_balance, current_balance,
			credit_limit, closing_day, due_day, currency, color, icon,
			is_active, include_in_total, last_sync_at, external_id,
			created_at, updated_at
		FROM bank_accounts
		WHERE user_id = ? AND is_active = TRUE
		ORDER BY name ASC
	`

	rows, err := r.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get bank accounts: %w", err)
	}
	defer rows.Close()

	return r.scanBankAccounts(rows)
}

// GetActiveByUser retrieves only active bank accounts for a user
func (r *BankAccountRepository) GetActiveByUser(ctx context.Context, userID uuid.UUID) ([]*model.BankAccount, error) {
	return r.GetAllByUser(ctx, userID)
}

// GetByType retrieves bank accounts by type for a user
func (r *BankAccountRepository) GetByType(ctx context.Context, userID uuid.UUID, accountType model.AccountType) ([]*model.BankAccount, error) {
	query := `
		SELECT id, user_id, name, bank_name, bank_code, account_type,
			account_number, agency, initial_balance, current_balance,
			credit_limit, closing_day, due_day, currency, color, icon,
			is_active, include_in_total, last_sync_at, external_id,
			created_at, updated_at
		FROM bank_accounts
		WHERE user_id = ? AND account_type = ? AND is_active = TRUE
		ORDER BY name ASC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), accountType)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank accounts by type: %w", err)
	}
	defer rows.Close()

	return r.scanBankAccounts(rows)
}

// Update updates a bank account
func (r *BankAccountRepository) Update(ctx context.Context, account *model.BankAccount) error {
	account.UpdatedAt = time.Now()

	query := `
		UPDATE bank_accounts SET
			name = ?, bank_name = ?, bank_code = ?, account_number = ?,
			agency = ?, credit_limit = ?, closing_day = ?, due_day = ?,
			color = ?, icon = ?, is_active = ?, include_in_total = ?,
			updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	var creditLimit sql.NullString
	if account.CreditLimit != nil {
		creditLimit = sql.NullString{String: account.CreditLimit.String(), Valid: true}
	}

	result, err := r.ExecContext(ctx, query,
		account.Name,
		NullString(account.BankName),
		NullString(account.BankCode),
		NullString(account.AccountNumber),
		NullString(account.Agency),
		creditLimit,
		NullInt32(account.ClosingDay),
		NullInt32(account.DueDay),
		account.Color,
		account.Icon,
		account.IsActive,
		account.IncludeInTotal,
		account.UpdatedAt,
		account.ID.String(),
		account.UserID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update bank account: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrBankAccountNotFound
	}

	return nil
}

// UpdateBalance updates the current balance of a bank account
func (r *BankAccountRepository) UpdateBalance(ctx context.Context, id uuid.UUID, balance decimal.Decimal) error {
	query := `UPDATE bank_accounts SET current_balance = ?, updated_at = ? WHERE id = ?`

	_, err := r.ExecContext(ctx, query, balance.String(), time.Now(), id.String())
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

// AdjustBalance adjusts the current balance by a delta amount
func (r *BankAccountRepository) AdjustBalance(ctx context.Context, id uuid.UUID, delta decimal.Decimal) error {
	query := `
		UPDATE bank_accounts
		SET current_balance = current_balance + ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.ExecContext(ctx, query, delta.String(), time.Now(), id.String())
	if err != nil {
		return fmt.Errorf("failed to adjust balance: %w", err)
	}

	return nil
}

// Delete soft deletes a bank account
func (r *BankAccountRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		UPDATE bank_accounts
		SET is_active = FALSE, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.ExecContext(ctx, query, time.Now(), id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("failed to delete bank account: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrBankAccountNotFound
	}

	return nil
}

// GetTotalBalance retrieves the total balance of all accounts for a user
func (r *BankAccountRepository) GetTotalBalance(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error) {
	query := `
		SELECT COALESCE(SUM(current_balance), 0)
		FROM bank_accounts
		WHERE user_id = ? AND is_active = TRUE AND include_in_total = TRUE
	`

	var total string
	err := r.QueryRowContext(ctx, query, userID.String()).Scan(&total)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get total balance: %w", err)
	}

	balance, err := decimal.NewFromString(total)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to parse total balance: %w", err)
	}

	return balance, nil
}

// CountAccounts counts the number of active accounts for a user
func (r *BankAccountRepository) CountAccounts(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM bank_accounts WHERE user_id = ? AND is_active = TRUE`

	var count int64
	err := r.QueryRowContext(ctx, query, userID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count accounts: %w", err)
	}

	return count, nil
}

// Helper functions

func (r *BankAccountRepository) scanBankAccount(row *sql.Row) (*model.BankAccount, error) {
	account := &model.BankAccount{}
	var bankName, bankCode, accountNumber, agency, externalID sql.NullString
	var creditLimit sql.NullString
	var closingDay, dueDay sql.NullInt32
	var lastSyncAt sql.NullTime
	var initialBalance, currentBalance string

	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&bankName,
		&bankCode,
		&account.AccountType,
		&accountNumber,
		&agency,
		&initialBalance,
		&currentBalance,
		&creditLimit,
		&closingDay,
		&dueDay,
		&account.Currency,
		&account.Color,
		&account.Icon,
		&account.IsActive,
		&account.IncludeInTotal,
		&lastSyncAt,
		&externalID,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBankAccountNotFound
		}
		return nil, fmt.Errorf("failed to scan bank account: %w", err)
	}

	account.BankName = StringPtr(bankName)
	account.BankCode = StringPtr(bankCode)
	account.AccountNumber = StringPtr(accountNumber)
	account.Agency = StringPtr(agency)
	account.ExternalID = StringPtr(externalID)
	account.ClosingDay = IntPtr(closingDay)
	account.DueDay = IntPtr(dueDay)

	if lastSyncAt.Valid {
		account.LastSyncAt = &lastSyncAt.Time
	}

	account.InitialBalance, _ = decimal.NewFromString(initialBalance)
	account.CurrentBalance, _ = decimal.NewFromString(currentBalance)

	if creditLimit.Valid {
		cl, _ := decimal.NewFromString(creditLimit.String)
		account.CreditLimit = &cl
	}

	return account, nil
}

func (r *BankAccountRepository) scanBankAccounts(rows *sql.Rows) ([]*model.BankAccount, error) {
	var accounts []*model.BankAccount

	for rows.Next() {
		account := &model.BankAccount{}
		var bankName, bankCode, accountNumber, agency, externalID sql.NullString
		var creditLimit sql.NullString
		var closingDay, dueDay sql.NullInt32
		var lastSyncAt sql.NullTime
		var initialBalance, currentBalance string

		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&bankName,
			&bankCode,
			&account.AccountType,
			&accountNumber,
			&agency,
			&initialBalance,
			&currentBalance,
			&creditLimit,
			&closingDay,
			&dueDay,
			&account.Currency,
			&account.Color,
			&account.Icon,
			&account.IsActive,
			&account.IncludeInTotal,
			&lastSyncAt,
			&externalID,
			&account.CreatedAt,
			&account.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan bank account: %w", err)
		}

		account.BankName = StringPtr(bankName)
		account.BankCode = StringPtr(bankCode)
		account.AccountNumber = StringPtr(accountNumber)
		account.Agency = StringPtr(agency)
		account.ExternalID = StringPtr(externalID)
		account.ClosingDay = IntPtr(closingDay)
		account.DueDay = IntPtr(dueDay)

		if lastSyncAt.Valid {
			account.LastSyncAt = &lastSyncAt.Time
		}

		account.InitialBalance, _ = decimal.NewFromString(initialBalance)
		account.CurrentBalance, _ = decimal.NewFromString(currentBalance)

		if creditLimit.Valid {
			cl, _ := decimal.NewFromString(creditLimit.String)
			account.CreditLimit = &cl
		}

		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating bank accounts: %w", err)
	}

	return accounts, nil
}
