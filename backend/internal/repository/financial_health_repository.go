package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/database"
	"github.com/gustavoz65/Fyapp/internal/model"
)

var ErrSnapshotNotFound = errors.New("snapshot not found")

type FinancialHealthRepository struct {
	*BaseRepository
}

func NewFinancialHealthRepository(db *database.Database, logger *zerolog.Logger) *FinancialHealthRepository {
	return &FinancialHealthRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// SaveSnapshot salva ou atualiza um snapshot (upsert)
func (r *FinancialHealthRepository) SaveSnapshot(ctx context.Context, snapshot *model.FinancialHealthSnapshot) error {
	query := `
		INSERT INTO financial_health_snapshots
		(id, user_id, period, score, score_percentage, economy_rate, budget_compliance,
		 goals_progress, spending_reduction, consistency, breakdown_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			score = VALUES(score),
			score_percentage = VALUES(score_percentage),
			economy_rate = VALUES(economy_rate),
			budget_compliance = VALUES(budget_compliance),
			goals_progress = VALUES(goals_progress),
			spending_reduction = VALUES(spending_reduction),
			consistency = VALUES(consistency),
			breakdown_json = VALUES(breakdown_json),
			created_at = VALUES(created_at)
	`

	snapshot.ID = uuid.New()
	snapshot.CreatedAt = time.Now()

	_, err := r.ExecContext(ctx, query,
		snapshot.ID.String(),
		snapshot.UserID.String(),
		snapshot.Period,
		snapshot.Score.String(),
		snapshot.ScorePercentage.String(),
		snapshot.EconomyRate.String(),
		snapshot.BudgetCompliance.String(),
		snapshot.GoalsProgress.String(),
		snapshot.SpendingReduction.String(),
		snapshot.Consistency.String(),
		NullString(snapshot.BreakdownJSON),
		snapshot.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save snapshot: %w", err)
	}

	return nil
}

// GetSnapshot busca snapshot por usuário e período
func (r *FinancialHealthRepository) GetSnapshot(ctx context.Context, userID uuid.UUID, period time.Time) (*model.FinancialHealthSnapshot, error) {
	query := `
		SELECT id, user_id, period, score, score_percentage, economy_rate, budget_compliance,
		       goals_progress, spending_reduction, consistency, breakdown_json, created_at
		FROM financial_health_snapshots
		WHERE user_id = ? AND period = ?
	`

	return r.scanSnapshot(r.QueryRowContext(ctx, query, userID.String(), period))
}

// GetHistoricalSnapshots busca últimos N snapshots do usuário
func (r *FinancialHealthRepository) GetHistoricalSnapshots(ctx context.Context, userID uuid.UUID, months int) ([]*model.FinancialHealthSnapshot, error) {
	query := `
		SELECT id, user_id, period, score, score_percentage, economy_rate, budget_compliance,
		       goals_progress, spending_reduction, consistency, breakdown_json, created_at
		FROM financial_health_snapshots
		WHERE user_id = ?
		ORDER BY period DESC
		LIMIT ?
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), months)
	if err != nil {
		return nil, fmt.Errorf("failed to get historical snapshots: %w", err)
	}
	defer rows.Close()

	var snapshots []*model.FinancialHealthSnapshot
	for rows.Next() {
		snapshot, err := r.scanSnapshotFromRow(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan snapshot: %w", err)
		}
		snapshots = append(snapshots, snapshot)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating snapshots: %w", err)
	}

	return snapshots, nil
}

// DeleteSnapshotsByUser deleta todos os snapshots de um usuário
func (r *FinancialHealthRepository) DeleteSnapshotsByUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM financial_health_snapshots WHERE user_id = ?`
	_, err := r.ExecContext(ctx, query, userID.String())
	return err
}

// scanSnapshot scans a single snapshot from a row
func (r *FinancialHealthRepository) scanSnapshot(row *sql.Row) (*model.FinancialHealthSnapshot, error) {
	var snapshot model.FinancialHealthSnapshot
	var breakdownJSON sql.NullString
	var score, scorePercentage, economyRate, budgetCompliance string
	var goalsProgress, spendingReduction, consistency string

	err := row.Scan(
		&snapshot.ID,
		&snapshot.UserID,
		&snapshot.Period,
		&score,
		&scorePercentage,
		&economyRate,
		&budgetCompliance,
		&goalsProgress,
		&spendingReduction,
		&consistency,
		&breakdownJSON,
		&snapshot.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSnapshotNotFound
		}
		return nil, fmt.Errorf("failed to scan snapshot: %w", err)
	}

	// Convert string values to decimal
	snapshot.Score, _ = decimal.NewFromString(score)
	snapshot.ScorePercentage, _ = decimal.NewFromString(scorePercentage)
	snapshot.EconomyRate, _ = decimal.NewFromString(economyRate)
	snapshot.BudgetCompliance, _ = decimal.NewFromString(budgetCompliance)
	snapshot.GoalsProgress, _ = decimal.NewFromString(goalsProgress)
	snapshot.SpendingReduction, _ = decimal.NewFromString(spendingReduction)
	snapshot.Consistency, _ = decimal.NewFromString(consistency)

	snapshot.BreakdownJSON = StringPtr(breakdownJSON)

	return &snapshot, nil
}

// scanSnapshotFromRow scans a snapshot from a Rows object
func (r *FinancialHealthRepository) scanSnapshotFromRow(rows *sql.Rows) (*model.FinancialHealthSnapshot, error) {
	var snapshot model.FinancialHealthSnapshot
	var breakdownJSON sql.NullString
	var score, scorePercentage, economyRate, budgetCompliance string
	var goalsProgress, spendingReduction, consistency string

	err := rows.Scan(
		&snapshot.ID,
		&snapshot.UserID,
		&snapshot.Period,
		&score,
		&scorePercentage,
		&economyRate,
		&budgetCompliance,
		&goalsProgress,
		&spendingReduction,
		&consistency,
		&breakdownJSON,
		&snapshot.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan snapshot: %w", err)
	}

	// Convert string values to decimal
	snapshot.Score, _ = decimal.NewFromString(score)
	snapshot.ScorePercentage, _ = decimal.NewFromString(scorePercentage)
	snapshot.EconomyRate, _ = decimal.NewFromString(economyRate)
	snapshot.BudgetCompliance, _ = decimal.NewFromString(budgetCompliance)
	snapshot.GoalsProgress, _ = decimal.NewFromString(goalsProgress)
	snapshot.SpendingReduction, _ = decimal.NewFromString(spendingReduction)
	snapshot.Consistency, _ = decimal.NewFromString(consistency)

	snapshot.BreakdownJSON = StringPtr(breakdownJSON)

	return &snapshot, nil
}
