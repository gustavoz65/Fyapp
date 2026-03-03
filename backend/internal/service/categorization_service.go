package service

import (
	"context"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
	"github.com/rs/zerolog"
)

type CategorizationService struct {
	patternRepo *repository.CategoryPatternRepository
	logger      *zerolog.Logger
}

func NewCategorizationService(
	patternRepo *repository.CategoryPatternRepository,
	logger *zerolog.Logger,
) *CategorizationService {
	return &CategorizationService{
		patternRepo: patternRepo,
		logger:      logger,
	}
}

// LearnFromTransaction aprende padroes a partir de uma transacao categorizada.
// Extrai keywords da descricao e associa com a categoria escolhida.
func (s *CategorizationService) LearnFromTransaction(ctx context.Context, userID uuid.UUID, description string, categoryID uuid.UUID, source string) {
	keywords := extractKeywords(description)
	if len(keywords) == 0 {
		return
	}

	for _, keyword := range keywords {
		pattern := &model.CategoryPattern{
			UserID:     userID,
			Keyword:    keyword,
			CategoryID: categoryID,
			Source:     source,
		}

		if err := s.patternRepo.Upsert(ctx, pattern); err != nil {
			s.logger.Error().
				Err(err).
				Str("user_id", userID.String()).
				Str("keyword", keyword).
				Str("category_id", categoryID.String()).
				Msg("falha ao salvar padrao de categoria")
		}
	}

	s.logger.Debug().
		Str("user_id", userID.String()).
		Strs("keywords", keywords).
		Str("category_id", categoryID.String()).
		Msg("padroes de categorizacao aprendidos")
}

// SuggestCategory sugere categorias baseado na descricao da transacao.
// Usa sistema de votacao: cada keyword contribui com seus padroes,
// e a categoria com mais votos ponderados pela confianca vence.
func (s *CategorizationService) SuggestCategory(ctx context.Context, userID uuid.UUID, description string) (*model.SuggestCategoryResponse, error) {
	keywords := extractKeywords(description)
	if len(keywords) == 0 {
		return &model.SuggestCategoryResponse{Suggestions: []model.CategorySuggestion{}}, nil
	}

	// Busca padroes por keywords exatas
	patterns, err := s.patternRepo.FindByKeywords(ctx, userID, keywords)
	if err != nil {
		return nil, err
	}

	// Se nao encontrou por keyword exata, tenta busca parcial com a descricao completa normalizada
	if len(patterns) == 0 {
		normalizedDesc := normalizeText(description)
		if normalizedDesc != "" {
			patterns, err = s.patternRepo.FindByPartialKeyword(ctx, userID, normalizedDesc)
			if err != nil {
				return nil, err
			}
		}
	}

	if len(patterns) == 0 {
		return &model.SuggestCategoryResponse{Suggestions: []model.CategorySuggestion{}}, nil
	}

	// Sistema de votacao: agrupa por categoria e soma confianca
	type vote struct {
		categoryName string
		totalScore   int
		matchCount   int
	}
	votes := make(map[uuid.UUID]*vote)

	for _, p := range patterns {
		v, exists := votes[p.CategoryID]
		if !exists {
			categoryName := ""
			if p.Category != nil {
				categoryName = p.Category.Name
			}
			v = &vote{categoryName: categoryName}
			votes[p.CategoryID] = v
		}
		v.totalScore += p.Confidence
		v.matchCount++
	}

	// Converte para resposta ordenada por score
	suggestions := make([]model.CategorySuggestion, 0, len(votes))
	for catID, v := range votes {
		suggestions = append(suggestions, model.CategorySuggestion{
			CategoryID:   catID,
			CategoryName: v.categoryName,
			Confidence:   v.totalScore,
			MatchCount:   v.matchCount,
		})
	}

	// Ordena por confianca decrescente
	for i := 0; i < len(suggestions); i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].Confidence > suggestions[i].Confidence {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}

	// Limita a 5 sugestoes
	if len(suggestions) > 5 {
		suggestions = suggestions[:5]
	}

	return &model.SuggestCategoryResponse{Suggestions: suggestions}, nil
}

// extractKeywords extrai keywords significativas de uma descricao.
// Remove stopwords e normaliza o texto.
func extractKeywords(description string) []string {
	normalized := normalizeText(description)
	words := strings.Fields(normalized)

	var keywords []string
	seen := make(map[string]bool)

	for _, word := range words {
		if len(word) < 3 {
			continue
		}
		if isStopword(word) {
			continue
		}
		if seen[word] {
			continue
		}
		seen[word] = true
		keywords = append(keywords, word)
	}

	// Tambem adiciona a descricao completa normalizada como keyword
	// (util para nomes proprios como "netflix", "spotify", etc.)
	fullNormalized := strings.Join(keywords, " ")
	if fullNormalized != "" && len(keywords) > 1 && !seen[fullNormalized] {
		keywords = append(keywords, fullNormalized)
	}

	return keywords
}

// normalizeText converte para minusculas e remove acentos/caracteres especiais
func normalizeText(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))

	var result strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// isStopword verifica se a palavra e uma stopword (palavras comuns sem significado)
func isStopword(word string) bool {
	stopwords := map[string]bool{
		// Portugues
		"de": true, "da": true, "do": true, "das": true, "dos": true,
		"em": true, "no": true, "na": true, "nos": true, "nas": true,
		"com": true, "por": true, "para": true, "que": true,
		"uma": true, "uns": true, "umas": true,
		"the": true, "and": true, "for": true, "from": true,
		"with": true, "this": true, "that": true,
		"pag": true, "pgto": true, "ref": true,
		"pix": true, "ted": true, "doc": true, "tef": true,
	}
	return stopwords[word]
}
