package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoz65/Fyapp/internal/database"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/rs/zerolog"
)

var (
	ErrRecurringTransactionNotFound = errors.New("recurring transaction not found")
)

type RecurringTransactionRepository struct {
	*BaseRepository
}

func NewRecurringTransactionRepository(db *database.Database, logger *zerolog.Logger) *RecurringTransactionRepository {
	return &RecurringTransactionRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

func (r *RecurringTransactionRepository) Create(ctx context.Context, rt *model.RecurringTransaction) error {
	rt.ID = uuid.New()
	rt.CreatedAt = time.Now()
	rt.UpdatedAt = time.Now()

	query := `
		INSERT INTO recurring_transactions (
			id, user_id, bank_account_id, category_id, type, amount,
			description, frequency, day_of_month, day_of_week,
			start_date, end_date, next_occurrence, last_generated_at,
			is_active, auto_confirm, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var categoryID sql.NullString
	if rt.CategoryID != nil {
		categoryID = sql.NullString{String: rt.CategoryID.String(), Valid: true}
	}

	var dayOfMonth, dayOfWeek sql.NullInt32
	if rt.DayOfMonth != nil {
		dayOfMonth = sql.NullInt32{Int32: int32(*rt.DayOfMonth), Valid: true}
	}
	if rt.DayOfWeek != nil {
		dayOfWeek = sql.NullInt32{Int32: int32(*rt.DayOfWeek), Valid: true}
	}

	var endDate, lastGeneratedAt sql.NullTime
	if rt.EndDate != nil {
		endDate = sql.NullTime{Time: *rt.EndDate, Valid: true}
	}
	if rt.LastGeneratedAt != nil {
		lastGeneratedAt = sql.NullTime{Time: *rt.LastGeneratedAt, Valid: true}
	}

	_, err := r.ExecContext(ctx, query,
		rt.ID.String(),
		rt.UserID.String(),
		rt.BankAccountID.String(),
		categoryID,
		rt.Type,
		rt.Amount,
		rt.Description,
		rt.Frequency,
		dayOfMonth,
		dayOfWeek,
		rt.StartDate,
		endDate,
		rt.NextOccurrence,
		lastGeneratedAt,
		rt.IsActive,
		rt.AutoConfirm,
		rt.CreatedAt,
		rt.UpdatedAt,
	)

	return err
}

func (r *RecurringTransactionRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.RecurringTransaction, error) {
	query := `
		SELECT rt.*, c.name as category_name, c.icon as category_icon, c.color as category_color,
			ba.name as bank_account_name, ba.account_type as bank_account_type
		FROM recurring_transactions rt
		LEFT JOIN categories c ON rt.category_id = c.id
		LEFT JOIN bank_accounts ba ON rt.bank_account_id = ba.id
		WHERE rt.id = ? AND rt.user_id = ?
	`

	var rt model.RecurringTransaction
	var categoryID, categoryName, categoryIcon, categoryColor sql.NullString
	var dayOfMonth, dayOfWeek sql.NullInt32
	var endDate, lastGeneratedAt sql.NullTime
	var bankAccountName, bankAccountType sql.NullString

	err := r.QueryRowContext(ctx, query, id.String(), userID.String()).Scan(
		&rt.ID, &rt.UserID, &rt.BankAccountID, &categoryID, &rt.Type,
		&rt.Amount, &rt.Description, &rt.Frequency, &dayOfMonth, &dayOfWeek,
		&rt.StartDate, &endDate, &rt.NextOccurrence, &lastGeneratedAt,
		&rt.IsActive, &rt.AutoConfirm, &rt.CreatedAt, &rt.UpdatedAt,
		&categoryName, &categoryIcon, &categoryColor,
		&bankAccountName, &bankAccountType,
	)

	if err == sql.ErrNoRows {
		return nil, ErrRecurringTransactionNotFound
	}
	if err != nil {
		return nil, err
	}

	if categoryID.Valid {
		catID, _ := uuid.Parse(categoryID.String)
		rt.CategoryID = &catID
		if categoryName.Valid {
			rt.Category = &model.Category{
				ID:    catID,
				Name:  categoryName.String,
				Icon:  categoryIcon.String,
				Color: categoryColor.String,
			}
		}
	}

	if dayOfMonth.Valid {
		day := int(dayOfMonth.Int32)
		rt.DayOfMonth = &day
	}
	if dayOfWeek.Valid {
		day := int(dayOfWeek.Int32)
		rt.DayOfWeek = &day
	}
	if endDate.Valid {
		rt.EndDate = &endDate.Time
	}
	if lastGeneratedAt.Valid {
		rt.LastGeneratedAt = &lastGeneratedAt.Time
	}

	if bankAccountName.Valid {
		rt.BankAccount = &model.BankAccount{
			ID:          rt.BankAccountID,
			Name:        bankAccountName.String,
			AccountType: model.AccountType(bankAccountType.String),
		}
	}

	return &rt, nil
}

func (r *RecurringTransactionRepository) GetAll(ctx context.Context, userID uuid.UUID, isActive *bool) ([]*model.RecurringTransaction, error) {
	query := `
		SELECT rt.*, c.name as category_name, c.icon as category_icon, c.color as category_color,
			ba.name as bank_account_name, ba.account_type as bank_account_type
		FROM recurring_transactions rt
		LEFT JOIN categories c ON rt.category_id = c.id
		LEFT JOIN bank_accounts ba ON rt.bank_account_id = ba.id
		WHERE rt.user_id = ?
	`

	args := []interface{}{userID.String()}

	if isActive != nil {
		query += " AND rt.is_active = ?"
		args = append(args, *isActive)
	}

	query += " ORDER BY rt.next_occurrence ASC"

	rows, err := r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurrings []*model.RecurringTransaction

	for rows.Next() {
		var rt model.RecurringTransaction
		var categoryID, categoryName, categoryIcon, categoryColor sql.NullString
		var dayOfMonth, dayOfWeek sql.NullInt32
		var endDate, lastGeneratedAt sql.NullTime
		var bankAccountName, bankAccountType sql.NullString

		err := rows.Scan(
			&rt.ID, &rt.UserID, &rt.BankAccountID, &categoryID, &rt.Type,
			&rt.Amount, &rt.Description, &rt.Frequency, &dayOfMonth, &dayOfWeek,
			&rt.StartDate, &endDate, &rt.NextOccurrence, &lastGeneratedAt,
			&rt.IsActive, &rt.AutoConfirm, &rt.CreatedAt, &rt.UpdatedAt,
			&categoryName, &categoryIcon, &categoryColor,
			&bankAccountName, &bankAccountType,
		)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			catID, _ := uuid.Parse(categoryID.String)
			rt.CategoryID = &catID
			if categoryName.Valid {
				rt.Category = &model.Category{
					ID:    catID,
					Name:  categoryName.String,
					Icon:  categoryIcon.String,
					Color: categoryColor.String,
				}
			}
		}

		if dayOfMonth.Valid {
			day := int(dayOfMonth.Int32)
			rt.DayOfMonth = &day
		}
		if dayOfWeek.Valid {
			day := int(dayOfWeek.Int32)
			rt.DayOfWeek = &day
		}
		if endDate.Valid {
			rt.EndDate = &endDate.Time
		}
		if lastGeneratedAt.Valid {
			rt.LastGeneratedAt = &lastGeneratedAt.Time
		}

		if bankAccountName.Valid {
			rt.BankAccount = &model.BankAccount{
				ID:          rt.BankAccountID,
				Name:        bankAccountName.String,
				AccountType: model.AccountType(bankAccountType.String),
			}
		}

		recurrings = append(recurrings, &rt)
	}

	return recurrings, rows.Err()
}

func (r *RecurringTransactionRepository) Update(ctx context.Context, rt *model.RecurringTransaction) error {
	rt.UpdatedAt = time.Now()

	query := `
		UPDATE recurring_transactions
		SET bank_account_id = ?, category_id = ?, type = ?, amount = ?,
			description = ?, frequency = ?, day_of_month = ?, day_of_week = ?,
			start_date = ?, end_date = ?, next_occurrence = ?, last_generated_at = ?,
			is_active = ?, auto_confirm = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	var categoryID sql.NullString
	if rt.CategoryID != nil {
		categoryID = sql.NullString{String: rt.CategoryID.String(), Valid: true}
	}

	var dayOfMonth, dayOfWeek sql.NullInt32
	if rt.DayOfMonth != nil {
		dayOfMonth = sql.NullInt32{Int32: int32(*rt.DayOfMonth), Valid: true}
	}
	if rt.DayOfWeek != nil {
		dayOfWeek = sql.NullInt32{Int32: int32(*rt.DayOfWeek), Valid: true}
	}

	var endDate, lastGeneratedAt sql.NullTime
	if rt.EndDate != nil {
		endDate = sql.NullTime{Time: *rt.EndDate, Valid: true}
	}
	if rt.LastGeneratedAt != nil {
		lastGeneratedAt = sql.NullTime{Time: *rt.LastGeneratedAt, Valid: true}
	}

	result, err := r.ExecContext(ctx, query,
		rt.BankAccountID.String(),
		categoryID,
		rt.Type,
		rt.Amount,
		rt.Description,
		rt.Frequency,
		dayOfMonth,
		dayOfWeek,
		rt.StartDate,
		endDate,
		rt.NextOccurrence,
		lastGeneratedAt,
		rt.IsActive,
		rt.AutoConfirm,
		rt.UpdatedAt,
		rt.ID.String(),
		rt.UserID.String(),
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecurringTransactionNotFound
	}

	return nil
}

func (r *RecurringTransactionRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM recurring_transactions WHERE id = ? AND user_id = ?`

	result, err := r.ExecContext(ctx, query, id.String(), userID.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecurringTransactionNotFound
	}

	return nil
}

func (r *RecurringTransactionRepository) GetDueRecurrings(ctx context.Context) ([]*model.RecurringTransaction, error) {
	query := `
		SELECT id, user_id, bank_account_id, category_id, type, amount,
			description, frequency, day_of_month, day_of_week,
			start_date, end_date, next_occurrence, last_generated_at,
			is_active, auto_confirm, created_at, updated_at
		FROM recurring_transactions
		WHERE is_active = TRUE
			AND next_occurrence <= NOW()
			AND (end_date IS NULL OR end_date >= NOW())
		ORDER BY next_occurrence ASC
	`

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recurrings []*model.RecurringTransaction

	for rows.Next() {
		var rt model.RecurringTransaction
		var categoryID sql.NullString
		var dayOfMonth, dayOfWeek sql.NullInt32
		var endDate, lastGeneratedAt sql.NullTime

		err := rows.Scan(
			&rt.ID, &rt.UserID, &rt.BankAccountID, &categoryID, &rt.Type,
			&rt.Amount, &rt.Description, &rt.Frequency, &dayOfMonth, &dayOfWeek,
			&rt.StartDate, &endDate, &rt.NextOccurrence, &lastGeneratedAt,
			&rt.IsActive, &rt.AutoConfirm, &rt.CreatedAt, &rt.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			catID, _ := uuid.Parse(categoryID.String)
			rt.CategoryID = &catID
		}
		if dayOfMonth.Valid {
			day := int(dayOfMonth.Int32)
			rt.DayOfMonth = &day
		}
		if dayOfWeek.Valid {
			day := int(dayOfWeek.Int32)
			rt.DayOfWeek = &day
		}
		if endDate.Valid {
			rt.EndDate = &endDate.Time
		}
		if lastGeneratedAt.Valid {
			rt.LastGeneratedAt = &lastGeneratedAt.Time
		}

		recurrings = append(recurrings, &rt)
	}

	return recurrings, rows.Err()
}

func (r *RecurringTransactionRepository) ToggleActive(ctx context.Context, id, userID uuid.UUID, isActive bool) error {
	query := `UPDATE recurring_transactions SET is_active = ?, updated_at = ? WHERE id = ? AND user_id = ?`

	result, err := r.ExecContext(ctx, query, isActive, time.Now(), id.String(), userID.String())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecurringTransactionNotFound
	}

	return nil
}

func (r *RecurringTransactionRepository) UpdateNextOccurrence(ctx context.Context, id uuid.UUID, nextOccurrence time.Time) error {
	now := time.Now()
	query := `UPDATE recurring_transactions SET next_occurrence = ?, last_generated_at = ?, updated_at = ? WHERE id = ?`

	_, err := r.ExecContext(ctx, query, nextOccurrence, now, now, id.String())
	return err
}

func (r *RecurringTransactionRepository) GetStats(ctx context.Context, userID uuid.UUID) (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total,
			SUM(CASE WHEN is_active = TRUE THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END) as total_income,
			SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END) as total_expense
		FROM recurring_transactions
		WHERE user_id = ?
	`

	var total, active int
	var totalIncome, totalExpense float64

	err := r.QueryRowContext(ctx, query, userID.String()).Scan(&total, &active, &totalIncome, &totalExpense)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":         total,
		"active":        active,
		"inactive":      total - active,
		"total_income":  totalIncome,
		"total_expense": totalExpense,
	}, nil
}

// CountActiveByUser conta o número de transações recorrentes ativas de um usuário
func (r *RecurringTransactionRepository) CountActiveByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM recurring_transactions
		WHERE user_id = ? AND is_active = TRUE
	`

	var count int64
	err := r.QueryRowContext(ctx, query, userID.String()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active recurring transactions: %w", err)
	}

	return count, nil
}
