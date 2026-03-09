package service

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/gustavoz65/Fyapp/internal/lib/categorization"
	"github.com/gustavoz65/Fyapp/internal/lib/similarity"
	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
)

type CategorizationService struct {
	patternRepo  *repository.CategoryPatternRepository
	categoryRepo *repository.CategoryRepository
	rules        *categorization.Rules
	logger       *zerolog.Logger
}

func NewCategorizationService(
	patternRepo *repository.CategoryPatternRepository,
	categoryRepo *repository.CategoryRepository,
	logger *zerolog.Logger,
) *CategorizationService {
	rules, err := categorization.LoadRules("rules.json")
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to load categorization rules, continuing without them")
	}

	return &CategorizationService{
		patternRepo:  patternRepo,
		categoryRepo: categoryRepo,
		rules:        rules,
		logger:       logger,
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

// LearnFromCorrection aprende quando usuario corrige uma categorizacao
func (s *CategorizationService) LearnFromCorrection(
	ctx context.Context,
	userID uuid.UUID,
	description string,
	wrongCategoryID uuid.UUID,
	correctCategoryID uuid.UUID,
) error {
	keywords := extractKeywords(description)
	if len(keywords) == 0 {
		return nil
	}

	// Decrementa confianca da categoria errada
	for _, keyword := range keywords {
		if err := s.patternRepo.DecrementConfidence(ctx, userID, keyword, wrongCategoryID); err != nil {
			s.logger.Error().Err(err).Msg("Failed to decrement confidence")
		}
	}

	// Incrementa confianca da categoria correta
	for _, keyword := range keywords {
		pattern := &model.CategoryPattern{
			UserID:     userID,
			Keyword:    keyword,
			CategoryID: correctCategoryID,
			Source:     "user_correction",
		}
		if err := s.patternRepo.Upsert(ctx, pattern); err != nil {
			s.logger.Error().Err(err).Msg("Failed to increment confidence")
		}
	}

	s.logger.Debug().
		Str("user_id", userID.String()).
		Str("description", description).
		Str("correct_category", correctCategoryID.String()).
		Msg("Learned from user correction")

	return nil
}

// SuggestCategory sugere categorias baseado na descricao da transacao.
// Usa sistema de votacao: cada keyword contribui com seus padroes,
// e a categoria com mais votos ponderados pela confianca vence.
func (s *CategorizationService) SuggestCategory(ctx context.Context, userID uuid.UUID, description string) (*model.SuggestCategoryResponse, error) {
	descNormalized := normalizeText(description)

	// 1. Try regex patterns first (highest priority)
	if s.rules != nil {
		for _, pattern := range s.rules.RegexPatterns {
			if matched, _ := regexp.MatchString(pattern.Pattern, description); matched {
				category, err := s.categoryRepo.GetByNameAndType(ctx, userID, pattern.Category, model.CategoryTypeExpense)
				if err == nil && category != nil {
					return &model.SuggestCategoryResponse{
						Suggestions: []model.CategorySuggestion{{
							CategoryID:   category.ID,
							CategoryName: category.Name,
							Confidence:   100,
							MatchCount:   1,
						}},
					}, nil
				}
			}
		}
	}

	// 2. Try keyword rules (high priority)
	if s.rules != nil {
		// Try expense patterns
		for _, expPattern := range s.rules.ExpensePatterns {
			for _, keyword := range expPattern.Keywords {
				if strings.Contains(descNormalized, keyword) {
					category, err := s.categoryRepo.GetByNameAndType(ctx, userID, expPattern.Category, model.CategoryTypeExpense)
					if err == nil && category != nil {
						return &model.SuggestCategoryResponse{
							Suggestions: []model.CategorySuggestion{{
								CategoryID:   category.ID,
								CategoryName: category.Name,
								Confidence:   expPattern.Confidence * 20, // scale to 0-100
								MatchCount:   1,
							}},
						}, nil
					}
				}
			}
		}

		// Try income patterns
		for _, incPattern := range s.rules.IncomePatterns {
			for _, keyword := range incPattern.Keywords {
				if strings.Contains(descNormalized, keyword) {
					category, err := s.categoryRepo.GetByNameAndType(ctx, userID, incPattern.Category, model.CategoryTypeIncome)
					if err == nil && category != nil {
						return &model.SuggestCategoryResponse{
							Suggestions: []model.CategorySuggestion{{
								CategoryID:   category.ID,
								CategoryName: category.Name,
								Confidence:   incPattern.Confidence * 20, // scale to 0-100
								MatchCount:   1,
							}},
						}, nil
					}
				}
			}
		}
	}

	// 3. Try exact keyword match from user patterns (medium priority)
	keywords := extractKeywords(description)
	if len(keywords) > 0 {
		patterns, err := s.patternRepo.FindByKeywords(ctx, userID, keywords)
		if err == nil && len(patterns) > 0 {
			return s.buildSuggestionsFromPatterns(patterns), nil
		}
	}

	// 4. Try similarity-based matching (fallback)
	allPatterns, err := s.patternRepo.GetTopPatternsForUser(ctx, userID, 50)
	if err != nil || len(allPatterns) == 0 {
		return &model.SuggestCategoryResponse{Suggestions: []model.CategorySuggestion{}}, nil
	}

	type scoredPattern struct {
		pattern    *model.CategoryPattern
		similarity float64
	}

	var scored []scoredPattern
	for _, p := range allPatterns {
		sim := similarity.CalculateSimilarity(descNormalized, normalizeText(p.Keyword))
		if sim >= 0.75 { // 75% similarity threshold
			scored = append(scored, scoredPattern{pattern: p, similarity: sim})
		}
	}

	if len(scored) == 0 {
		return &model.SuggestCategoryResponse{Suggestions: []model.CategorySuggestion{}}, nil
	}

	// Sort by similarity DESC
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].similarity > scored[j].similarity
	})

	// Group by category
	categoryScores := make(map[uuid.UUID]*model.CategorySuggestion)
	for _, s := range scored {
		if existing, ok := categoryScores[s.pattern.CategoryID]; ok {
			existing.Confidence += int(s.similarity * 100)
			existing.MatchCount++
		} else {
			categoryName := ""
			if s.pattern.Category != nil {
				categoryName = s.pattern.Category.Name
			}
			categoryScores[s.pattern.CategoryID] = &model.CategorySuggestion{
				CategoryID:   s.pattern.CategoryID,
				CategoryName: categoryName,
				Confidence:   int(s.similarity * 100),
				MatchCount:   1,
			}
		}
	}

	suggestions := make([]model.CategorySuggestion, 0, len(categoryScores))
	for _, sugg := range categoryScores {
		suggestions = append(suggestions, *sugg)
	}

	// Sort by confidence
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Confidence > suggestions[j].Confidence
	})

	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return &model.SuggestCategoryResponse{Suggestions: suggestions}, nil
}

func (s *CategorizationService) buildSuggestionsFromPatterns(patterns []*model.CategoryPattern) *model.SuggestCategoryResponse {
	votes := make(map[uuid.UUID]*struct {
		categoryName string
		totalScore   int
		matchCount   int
	})

	for _, p := range patterns {
		v, exists := votes[p.CategoryID]
		if !exists {
			categoryName := ""
			if p.Category != nil {
				categoryName = p.Category.Name
			}
			v = &struct {
				categoryName string
				totalScore   int
				matchCount   int
			}{categoryName: categoryName}
			votes[p.CategoryID] = v
		}
		v.totalScore += p.Confidence
		v.matchCount++
	}

	suggestions := make([]model.CategorySuggestion, 0, len(votes))
	for catID, v := range votes {
		suggestions = append(suggestions, model.CategorySuggestion{
			CategoryID:   catID,
			CategoryName: v.categoryName,
			Confidence:   v.totalScore,
			MatchCount:   v.matchCount,
		})
	}

	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Confidence > suggestions[j].Confidence
	})

	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return &model.SuggestCategoryResponse{Suggestions: suggestions}
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
