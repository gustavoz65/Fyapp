package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gustavoz65/Cashing-go/backend/internal/database"
	"github.com/gustavoz65/Cashing-go/backend/internal/model"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

var (
	ErrBudgetNotFound = errors.New("budget not found")
)

type BudgetRepository struct {
	*BaseRepository
}

func NewBudgetRepository(db *database.Database, logger *zerolog.Logger) *BudgetRepository {
	return &BudgetRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Create creates a new budget
func (r *BudgetRepository) Create(ctx context.Context, budget *model.Budget) error {
	budget.ID = uuid.New()
	budget.SpentAmount = decimal.Zero
	budget.CreatedAt = time.Now()
	budget.UpdatedAt = time.Now()

	query := `
		INSERT INTO budgets (
			id, user_id, category_id, name, amount, spent_amount,
			period_type, start_date, end_date, alert_threshold,
			alert_sent, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var categoryID sql.NullString
	if budget.CategoryID != nil {
		categoryID = sql.NullString{String: budget.CategoryID.String(), Valid: true}
	}

	_, err := r.ExecContext(ctx, query,
		budget.ID.String(),
		budget.UserID.String(),
		categoryID,
		budget.Name,
		budget.Amount.String(),
		budget.SpentAmount.String(),
		budget.PeriodType,
		budget.StartDate,
		budget.EndDate,
		budget.AlertThreshold.String(),
		budget.AlertSent,
		budget.IsActive,
		budget.CreatedAt,
		budget.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create budget: %w", err)
	}

	return nil
}

// GetByID retrieves a budget by ID
func (r *BudgetRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Budget, error) {
	query := `
		SELECT id, user_id, category_id, name, amount, spent_amount,
			period_type, start_date, end_date, alert_threshold,
			alert_sent, is_active, created_at, updated_at
		FROM budgets
		WHERE id = ?
	`

	return r.scanBudget(r.QueryRowContext(ctx, query, id.String()))
}

// GetByIDAndUser retrieves a budget by ID ensuring it belongs to the user
func (r *BudgetRepository) GetByIDAndUser(ctx context.Context, id, userID uuid.UUID) (*model.Budget, error) {
	query := `
		SELECT id, user_id, category_id, name, amount, spent_amount,
			period_type, start_date, end_date, alert_threshold,
			alert_sent, is_active, created_at, updated_at
		FROM budgets
		WHERE id = ? AND user_id = ?
	`

	return r.scanBudget(r.QueryRowContext(ctx, query, id.String(), userID.String()))
}

// GetAllByUser retrieves all budgets for a user
func (r *BudgetRepository) GetAllByUser(ctx context.Context, userID uuid.UUID) ([]*model.Budget, error) {
	query := `
		SELECT id, user_id, category_id, name, amount, spent_amount,
			period_type, start_date, end_date, alert_threshold,
			alert_sent, is_active, created_at, updated_at
		FROM budgets
		WHERE user_id = ? AND is_active = TRUE
		ORDER BY start_date DESC
	`

	rows, err := r.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets: %w", err)
	}
	defer rows.Close()

	return r.scanBudgets(rows)
}

// GetActiveByUser retrieves active budgets for the current period
func (r *BudgetRepository) GetActiveByUser(ctx context.Context, userID uuid.UUID) ([]*model.Budget, error) {
	query := `
		SELECT id, user_id, category_id, name, amount, spent_amount,
			period_type, start_date, end_date, alert_threshold,
			alert_sent, is_active, created_at, updated_at
		FROM budgets
		WHERE user_id = ? AND is_active = TRUE
			AND CURRENT_DATE BETWEEN start_date AND end_date
		ORDER BY name ASC
	`

	rows, err := r.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get active budgets: %w", err)
	}
	defer rows.Close()

	return r.scanBudgets(rows)
}

// GetByCategory retrieves budgets by category
func (r *BudgetRepository) GetByCategory(ctx context.Context, userID, categoryID uuid.UUID) ([]*model.Budget, error) {
	query := `
		SELECT id, user_id, category_id, name, amount, spent_amount,
			period_type, start_date, end_date, alert_threshold,
			alert_sent, is_active, created_at, updated_at
		FROM budgets
		WHERE user_id = ? AND category_id = ? AND is_active = TRUE
		ORDER BY start_date DESC
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), categoryID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets by category: %w", err)
	}
	defer rows.Close()

	return r.scanBudgets(rows)
}

// Update updates a budget
func (r *BudgetRepository) Update(ctx context.Context, budget *model.Budget) error {
	budget.UpdatedAt = time.Now()

	query := `
		UPDATE budgets SET
			name = ?, amount = ?, alert_threshold = ?,
			is_active = ?, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.ExecContext(ctx, query,
		budget.Name,
		budget.Amount.String(),
		budget.AlertThreshold.String(),
		budget.IsActive,
		budget.UpdatedAt,
		budget.ID.String(),
		budget.UserID.String(),
	)

	if err != nil {
		return fmt.Errorf("failed to update budget: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrBudgetNotFound
	}

	return nil
}

// UpdateSpentAmount updates the spent amount of a budget
func (r *BudgetRepository) UpdateSpentAmount(ctx context.Context, id uuid.UUID, spentAmount decimal.Decimal) error {
	query := `
		UPDATE budgets
		SET spent_amount = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.ExecContext(ctx, query, spentAmount.String(), time.Now(), id.String())
	if err != nil {
		return fmt.Errorf("failed to update spent amount: %w", err)
	}

	return nil
}

// MarkAlertSent marks that the alert has been sent for a budget
func (r *BudgetRepository) MarkAlertSent(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE budgets SET alert_sent = TRUE, updated_at = ? WHERE id = ?`

	_, err := r.ExecContext(ctx, query, time.Now(), id.String())
	if err != nil {
		return fmt.Errorf("failed to mark alert sent: %w", err)
	}

	return nil
}

// GetBudgetsNeedingAlert retrieves budgets that need alerts
func (r *BudgetRepository) GetBudgetsNeedingAlert(ctx context.Context) ([]*model.Budget, error) {
	query := `
		SELECT id, user_id, category_id, name, amount, spent_amount,
			period_type, start_date, end_date, alert_threshold,
			alert_sent, is_active, created_at, updated_at
		FROM budgets
		WHERE is_active = TRUE
			AND alert_sent = FALSE
			AND CURRENT_DATE BETWEEN start_date AND end_date
			AND (spent_amount / amount * 100) >= alert_threshold
	`

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get budgets needing alert: %w", err)
	}
	defer rows.Close()

	return r.scanBudgets(rows)
}

// Delete soft deletes a budget
func (r *BudgetRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		UPDATE budgets
		SET is_active = FALSE, updated_at = ?
		WHERE id = ? AND user_id = ?
	`

	result, err := r.ExecContext(ctx, query, time.Now(), id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrBudgetNotFound
	}

	return nil
}

// RecalculateSpentAmount recalculates the spent amount for a budget based on transactions
func (r *BudgetRepository) RecalculateSpentAmount(ctx context.Context, budget *model.Budget) error {
	var query string
	var args []interface{}

	if budget.CategoryID != nil {
		query = `
			SELECT COALESCE(SUM(amount), 0)
			FROM transactions
			WHERE user_id = ? AND type = 'expense' AND is_paid = TRUE
				AND category_id = ?
				AND transaction_date BETWEEN ? AND ?
		`
		args = []interface{}{budget.UserID.String(), budget.CategoryID.String(), budget.StartDate, budget.EndDate}
	} else {
		query = `
			SELECT COALESCE(SUM(amount), 0)
			FROM transactions
			WHERE user_id = ? AND type = 'expense' AND is_paid = TRUE
				AND transaction_date BETWEEN ? AND ?
		`
		args = []interface{}{budget.UserID.String(), budget.StartDate, budget.EndDate}
	}

	var spent string
	err := r.QueryRowContext(ctx, query, args...).Scan(&spent)
	if err != nil {
		return fmt.Errorf("failed to calculate spent amount: %w", err)
	}

	spentAmount, _ := decimal.NewFromString(spent)
	return r.UpdateSpentAmount(ctx, budget.ID, spentAmount)
}

// Helper functions

func (r *BudgetRepository) scanBudget(row *sql.Row) (*model.Budget, error) {
	budget := &model.Budget{}
	var categoryID sql.NullString
	var amount, spentAmount, alertThreshold string

	err := row.Scan(
		&budget.ID,
		&budget.UserID,
		&categoryID,
		&budget.Name,
		&amount,
		&spentAmount,
		&budget.PeriodType,
		&budget.StartDate,
		&budget.EndDate,
		&alertThreshold,
		&budget.AlertSent,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, fmt.Errorf("failed to scan budget: %w", err)
	}

	if categoryID.Valid {
		id, _ := uuid.Parse(categoryID.String)
		budget.CategoryID = &id
	}

	budget.Amount, _ = decimal.NewFromString(amount)
	budget.SpentAmount, _ = decimal.NewFromString(spentAmount)
	budget.AlertThreshold, _ = decimal.NewFromString(alertThreshold)

	return budget, nil
}

func (r *BudgetRepository) scanBudgets(rows *sql.Rows) ([]*model.Budget, error) {
	var budgets []*model.Budget

	for rows.Next() {
		budget := &model.Budget{}
		var categoryID sql.NullString
		var amount, spentAmount, alertThreshold string

		err := rows.Scan(
			&budget.ID,
			&budget.UserID,
			&categoryID,
			&budget.Name,
			&amount,
			&spentAmount,
			&budget.PeriodType,
			&budget.StartDate,
			&budget.EndDate,
			&alertThreshold,
			&budget.AlertSent,
			&budget.IsActive,
			&budget.CreatedAt,
			&budget.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}

		if categoryID.Valid {
			id, _ := uuid.Parse(categoryID.String)
			budget.CategoryID = &id
		}

		budget.Amount, _ = decimal.NewFromString(amount)
		budget.SpentAmount, _ = decimal.NewFromString(spentAmount)
		budget.AlertThreshold, _ = decimal.NewFromString(alertThreshold)

		budgets = append(budgets, budget)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating budgets: %w", err)
	}

	return budgets, nil
}
