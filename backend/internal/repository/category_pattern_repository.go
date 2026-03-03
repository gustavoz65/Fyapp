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
	ErrPatternNotFound = errors.New("category pattern not found")
)

type CategoryPatternRepository struct {
	*BaseRepository
}

func NewCategoryPatternRepository(db *database.Database, logger *zerolog.Logger) *CategoryPatternRepository {
	return &CategoryPatternRepository{
		BaseRepository: NewBaseRepository(db, logger),
	}
}

// Upsert cria ou incrementa a confianca de um padrao existente
func (r *CategoryPatternRepository) Upsert(ctx context.Context, pattern *model.CategoryPattern) error {
	query := `
		INSERT INTO category_patterns (id, user_id, keyword, category_id, confidence, source, created_at, updated_at)
		VALUES (?, ?, ?, ?, 1, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			confidence = confidence + 1,
			updated_at = ?
	`

	now := time.Now()
	pattern.ID = uuid.New()
	pattern.CreatedAt = now
	pattern.UpdatedAt = now

	_, err := r.ExecContext(ctx, query,
		pattern.ID.String(),
		pattern.UserID.String(),
		pattern.Keyword,
		pattern.CategoryID.String(),
		pattern.Source,
		now,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("falha ao upsert padrao de categoria: %w", err)
	}

	return nil
}

// FindByKeywords busca padroes que batem com qualquer uma das keywords para um usuario
func (r *CategoryPatternRepository) FindByKeywords(ctx context.Context, userID uuid.UUID, keywords []string) ([]*model.CategoryPattern, error) {
	if len(keywords) == 0 {
		return nil, nil
	}

	placeholders := ""
	args := []interface{}{userID.String()}
	for i, kw := range keywords {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
		args = append(args, kw)
	}

	query := fmt.Sprintf(`
		SELECT cp.id, cp.user_id, cp.keyword, cp.category_id, cp.confidence, cp.source,
			cp.created_at, cp.updated_at,
			c.name as category_name
		FROM category_patterns cp
		JOIN categories c ON cp.category_id = c.id AND c.is_active = TRUE
		WHERE cp.user_id = ? AND cp.keyword IN (%s)
		ORDER BY cp.confidence DESC
	`, placeholders)

	rows, err := r.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar padroes por keywords: %w", err)
	}
	defer rows.Close()

	return r.scanPatternsWithCategory(rows)
}

// FindByPartialKeyword busca padroes com keyword parcial (LIKE) para um usuario
func (r *CategoryPatternRepository) FindByPartialKeyword(ctx context.Context, userID uuid.UUID, partial string) ([]*model.CategoryPattern, error) {
	query := `
		SELECT cp.id, cp.user_id, cp.keyword, cp.category_id, cp.confidence, cp.source,
			cp.created_at, cp.updated_at,
			c.name as category_name
		FROM category_patterns cp
		JOIN categories c ON cp.category_id = c.id AND c.is_active = TRUE
		WHERE cp.user_id = ? AND cp.keyword LIKE ?
		ORDER BY cp.confidence DESC
		LIMIT 20
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), "%"+partial+"%")
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar padroes parciais: %w", err)
	}
	defer rows.Close()

	return r.scanPatternsWithCategory(rows)
}

// GetTopPatternsForUser retorna os padroes mais usados de um usuario
func (r *CategoryPatternRepository) GetTopPatternsForUser(ctx context.Context, userID uuid.UUID, limit int) ([]*model.CategoryPattern, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT cp.id, cp.user_id, cp.keyword, cp.category_id, cp.confidence, cp.source,
			cp.created_at, cp.updated_at,
			c.name as category_name
		FROM category_patterns cp
		JOIN categories c ON cp.category_id = c.id AND c.is_active = TRUE
		WHERE cp.user_id = ?
		ORDER BY cp.confidence DESC
		LIMIT ?
	`

	rows, err := r.QueryContext(ctx, query, userID.String(), limit)
	if err != nil {
		return nil, fmt.Errorf("falha ao buscar top padroes: %w", err)
	}
	defer rows.Close()

	return r.scanPatternsWithCategory(rows)
}

// DeleteByUser remove todos os padroes de um usuario
func (r *CategoryPatternRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM category_patterns WHERE user_id = ?`
	_, err := r.ExecContext(ctx, query, userID.String())
	if err != nil {
		return fmt.Errorf("falha ao deletar padroes do usuario: %w", err)
	}
	return nil
}

// scanPatternsWithCategory escaneia padroes com nome da categoria
func (r *CategoryPatternRepository) scanPatternsWithCategory(rows *sql.Rows) ([]*model.CategoryPattern, error) {
	var patterns []*model.CategoryPattern

	for rows.Next() {
		p := &model.CategoryPattern{}
		var categoryName string

		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Keyword,
			&p.CategoryID,
			&p.Confidence,
			&p.Source,
			&p.CreatedAt,
			&p.UpdatedAt,
			&categoryName,
		)

		if err != nil {
			return nil, fmt.Errorf("falha ao escanear padrao: %w", err)
		}

		p.Category = &model.Category{
			ID:   p.CategoryID,
			Name: categoryName,
		}

		patterns = append(patterns, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar padroes: %w", err)
	}

	return patterns, nil
}
