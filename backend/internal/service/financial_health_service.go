package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
)

type FinancialHealthService struct {
	txRepo      *repository.TransactionRepository
	budgetRepo  *repository.BudgetRepository
	goalRepo    *repository.GoalRepository
	accountRepo *repository.BankAccountRepository
	healthRepo  *repository.FinancialHealthRepository
	cache       *redis.Client
	logger      *zerolog.Logger
}

func NewFinancialHealthService(
	txRepo *repository.TransactionRepository,
	budgetRepo *repository.BudgetRepository,
	goalRepo *repository.GoalRepository,
	accountRepo *repository.BankAccountRepository,
	healthRepo *repository.FinancialHealthRepository,
	cache *redis.Client,
	logger *zerolog.Logger,
) *FinancialHealthService {
	return &FinancialHealthService{
		txRepo:      txRepo,
		budgetRepo:  budgetRepo,
		goalRepo:    goalRepo,
		accountRepo: accountRepo,
		healthRepo:  healthRepo,
		cache:       cache,
		logger:      logger,
	}
}

// CalculateScore calcula o score de saúde financeira
func (s *FinancialHealthService) CalculateScore(ctx context.Context, userID uuid.UUID, period time.Time) (*model.FinancialHealthSnapshot, error) {
	// Normalizar para primeiro dia do mês
	period = time.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0, time.UTC)

	startDate := period
	endDate := period.AddDate(0, 1, -1) // Último dia do mês

	s.logger.Info().
		Str("user_id", userID.String()).
		Time("period", period).
		Msg("Calculating financial health score")

	// 1. Buscar income/expense do período
	income, err := s.txRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get income: %w", err)
	}

	expense, err := s.txRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense: %w", err)
	}

	balance := income.Sub(expense)

	// 2. Calcular taxa de economia (40%)
	economyRateScore := calculateEconomyRate(income, expense)
	economyPercentage := float64(0)
	if income.GreaterThan(decimal.Zero) {
		economyPercentage = balance.Div(income).InexactFloat64() * 100
	}

	// 3. Buscar budgets e calcular cumprimento (20%)
	budgetScore, budgetDetails, err := s.calculateBudgetCompliance(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate budget compliance: %w", err)
	}

	// 4. Buscar goals e calcular progresso (20%)
	goalsScore, goalDetails, err := s.calculateGoalsProgress(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate goals progress: %w", err)
	}

	// 5. Calcular redução de gastos vs mês anterior (10%)
	spendingReductionScore, categoryDetails, err := s.calculateSpendingReduction(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate spending reduction: %w", err)
	}

	// 6. Calcular consistência (10%)
	consistencyScore, consistencyDetails, err := s.calculateConsistency(ctx, userID, period)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate consistency: %w", err)
	}

	// 7. Score final = (E*0.4) + (B*0.2) + (M*0.2) + (R*0.1) + (C*0.1)
	finalScore := (economyRateScore * 0.40) +
		(budgetScore * 0.20) +
		(goalsScore * 0.20) +
		(spendingReductionScore * 0.10) +
		(consistencyScore * 0.10)

	// 8. Montar breakdown JSON
	breakdown := model.HealthScoreBreakdown{
		Income:             income,
		Expense:            expense,
		Balance:            balance,
		EconomyPercentage:  economyPercentage,
		Budgets:            budgetDetails,
		Goals:              goalDetails,
		TopCategories:      categoryDetails,
		ConsistencyDetails: consistencyDetails,
	}

	breakdownJSON, err := json.Marshal(breakdown)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal breakdown: %w", err)
	}
	breakdownStr := string(breakdownJSON)

	// 9. Retornar snapshot
	snapshot := &model.FinancialHealthSnapshot{
		UserID:            userID,
		Period:            period,
		Score:             decimal.NewFromFloat(finalScore),
		ScorePercentage:   decimal.NewFromFloat(finalScore * 20), // 0-5 -> 0-100
		EconomyRate:       decimal.NewFromFloat(economyRateScore),
		BudgetCompliance:  decimal.NewFromFloat(budgetScore),
		GoalsProgress:     decimal.NewFromFloat(goalsScore),
		SpendingReduction: decimal.NewFromFloat(spendingReductionScore),
		Consistency:       decimal.NewFromFloat(consistencyScore),
		BreakdownJSON:     &breakdownStr,
	}

	// Salvar snapshot no banco
	if err := s.healthRepo.SaveSnapshot(ctx, snapshot); err != nil {
		return nil, fmt.Errorf("failed to save snapshot: %w", err)
	}

	s.logger.Info().
		Str("user_id", userID.String()).
		Float64("score", finalScore).
		Msg("Financial health score calculated successfully")

	return snapshot, nil
}

// calculateEconomyRate calcula a taxa de economia (componente E)
// Retorna score de 0-5 baseado na % de economia
func calculateEconomyRate(income, expense decimal.Decimal) float64 {
	if income.IsZero() || income.LessThanOrEqual(decimal.Zero) {
		return 0.0
	}

	balance := income.Sub(expense)
	economyRate := balance.Div(income).InexactFloat64() * 100

	switch {
	case economyRate >= 30:
		return 5.0
	case economyRate >= 20:
		return 4.0
	case economyRate >= 10:
		return 3.0
	case economyRate >= 0:
		return 2.0
	case economyRate >= -10:
		return 1.0
	default:
		return 0.0
	}
}

// calculateBudgetCompliance calcula o cumprimento de budgets (componente B)
// Retorna score de 0-5 e detalhes dos budgets
func (s *FinancialHealthService) calculateBudgetCompliance(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (float64, []model.BudgetHealthDetail, error) {
	// Buscar budgets ativos no período
	budgets, err := s.budgetRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get budgets: %w", err)
	}

	// Filtrar budgets que estão no período
	var activeBudgets []*model.Budget
	for _, budget := range budgets {
		if budget.StartDate.Before(endDate.AddDate(0, 0, 1)) && budget.EndDate.After(startDate.AddDate(0, 0, -1)) {
			activeBudgets = append(activeBudgets, budget)
		}
	}

	// Se não tem budgets, retorna score neutro
	if len(activeBudgets) == 0 {
		return 3.0, []model.BudgetHealthDetail{}, nil
	}

	var totalScore float64
	var details []model.BudgetHealthDetail

	for _, budget := range activeBudgets {
		// Recalcular spent amount para garantir precisão
		if err := s.budgetRepo.RecalculateSpentAmount(ctx, budget); err != nil {
			s.logger.Warn().Err(err).Str("budget_id", budget.ID.String()).Msg("Failed to recalculate budget")
			continue
		}

		// Buscar budget atualizado
		budget, err = s.budgetRepo.GetByID(ctx, budget.ID)
		if err != nil {
			continue
		}

		usage := float64(0)
		if budget.Amount.GreaterThan(decimal.Zero) {
			usage = budget.SpentAmount.Div(budget.Amount).InexactFloat64() * 100
		}

		// Calcular score individual do budget
		var score float64
		switch {
		case usage <= 70:
			score = 5.0 // Excelente
		case usage <= 85:
			score = 4.0 // Bom
		case usage <= 100:
			score = 3.0 // OK
		case usage <= 115:
			score = 2.0 // Atenção
		case usage <= 130:
			score = 1.0 // Ruim
		default:
			score = 0.0 // Crítico
		}

		totalScore += score

		details = append(details, model.BudgetHealthDetail{
			ID:    budget.ID,
			Name:  budget.Name,
			Limit: budget.Amount,
			Spent: budget.SpentAmount,
			Usage: usage,
			Score: score,
		})
	}

	// Retornar média dos scores
	avgScore := totalScore / float64(len(activeBudgets))
	return avgScore, details, nil
}

// calculateGoalsProgress calcula o progresso de metas (componente M)
// Retorna score de 0-5 e detalhes das metas
func (s *FinancialHealthService) calculateGoalsProgress(ctx context.Context, userID uuid.UUID) (float64, []model.GoalHealthDetail, error) {
	// Buscar metas ativas
	goals, err := s.goalRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get goals: %w", err)
	}

	// Se não tem metas, retorna score neutro
	if len(goals) == 0 {
		return 3.0, []model.GoalHealthDetail{}, nil
	}

	var totalScore float64
	var details []model.GoalHealthDetail

	for _, goal := range goals {
		progress := float64(0)
		if goal.TargetAmount.GreaterThan(decimal.Zero) {
			progress = goal.CurrentAmount.Div(goal.TargetAmount).InexactFloat64() * 100
		}

		// Calcular score individual da meta
		var score float64
		switch {
		case progress >= 80:
			score = 5.0 // Excelente
		case progress >= 60:
			score = 4.0 // Bom
		case progress >= 40:
			score = 3.0 // Regular
		case progress >= 20:
			score = 2.0 // Baixo
		case progress >= 10:
			score = 1.0 // Muito baixo
		default:
			score = 0.0 // Sem progresso
		}

		totalScore += score

		details = append(details, model.GoalHealthDetail{
			ID:            goal.ID,
			Name:          goal.Name,
			TargetAmount:  goal.TargetAmount,
			CurrentAmount: goal.CurrentAmount,
			Progress:      progress,
			Score:         score,
		})
	}

	// Retornar média dos scores
	avgScore := totalScore / float64(len(goals))
	return avgScore, details, nil
}

// calculateSpendingReduction calcula a redução de gastos (componente R)
// Compara top 3 categorias de despesa do mês atual vs anterior
func (s *FinancialHealthService) calculateSpendingReduction(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (float64, []model.CategorySpendingDetail, error) {
	// Buscar top categorias do mês atual
	currentCategories, err := s.txRepo.GetSumByCategory(ctx, userID, startDate, endDate)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get current categories: %w", err)
	}

	// Se não tem dados, retorna neutro
	if len(currentCategories) == 0 {
		return 3.0, []model.CategorySpendingDetail{}, nil
	}

	// Buscar top categorias do mês anterior
	prevStartDate := startDate.AddDate(0, -1, 0)
	prevEndDate := endDate.AddDate(0, -1, 0)
	prevCategories, err := s.txRepo.GetSumByCategory(ctx, userID, prevStartDate, prevEndDate)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get previous categories: %w", err)
	}

	// Criar map de categorias anteriores para fácil lookup
	prevMap := make(map[string]decimal.Decimal)
	for _, cat := range prevCategories {
		prevMap[cat.CategoryName] = cat.Amount
	}

	// Pegar top 3 categorias do mês atual
	topCount := 3
	if len(currentCategories) < topCount {
		topCount = len(currentCategories)
	}

	var details []model.CategorySpendingDetail
	var totalScore float64
	var scoredCount int

	for i := 0; i < topCount; i++ {
		current := currentCategories[i]
		previous := prevMap[current.CategoryName]

		changePercentage := float64(0)
		if previous.GreaterThan(decimal.Zero) {
			change := current.Amount.Sub(previous)
			changePercentage = change.Div(previous).InexactFloat64() * 100
		} else if current.Amount.GreaterThan(decimal.Zero) {
			// Categoria nova no mês atual = aumento de 100%
			changePercentage = 100
		}

		// Calcular score (redução é positiva, aumento é negativa)
		var score float64
		switch {
		case changePercentage <= -20: // Reduziu 20%+
			score = 5.0
		case changePercentage <= -10: // Reduziu 10-20%
			score = 4.0
		case changePercentage <= 0: // Reduziu até 10%
			score = 3.0
		case changePercentage <= 10: // Aumentou até 10%
			score = 2.0
		case changePercentage <= 20: // Aumentou 10-20%
			score = 1.0
		default: // Aumentou 20%+
			score = 0.0
		}

		totalScore += score
		scoredCount++

		categoryID := uuid.Nil
		if current.CategoryID != nil {
			categoryID = *current.CategoryID
		}

		details = append(details, model.CategorySpendingDetail{
			CategoryID:       categoryID,
			CategoryName:     current.CategoryName,
			CurrentAmount:    current.Amount,
			PreviousAmount:   previous,
			ChangePercentage: changePercentage,
			Score:            score,
		})
	}

	avgScore := float64(3.0) // Neutro se não tem dados
	if scoredCount > 0 {
		avgScore = totalScore / float64(scoredCount)
	}

	return avgScore, details, nil
}

// calculateConsistency calcula a consistência financeira (componente C)
// Conta meses negativos (últimos 6) e dívidas ativas
func (s *FinancialHealthService) calculateConsistency(ctx context.Context, userID uuid.UUID, currentPeriod time.Time) (float64, model.ConsistencyDetail, error) {
	// Contar meses negativos nos últimos 6 meses
	negativeMonths, err := s.countNegativeMonths(ctx, userID, currentPeriod, 6)
	if err != nil {
		return 0, model.ConsistencyDetail{}, fmt.Errorf("failed to count negative months: %w", err)
	}

	// Contar dívidas ativas (credit cards com saldo negativo)
	activeDebts, err := s.countActiveDebts(ctx, userID)
	if err != nil {
		return 0, model.ConsistencyDetail{}, fmt.Errorf("failed to count active debts: %w", err)
	}

	// Calcular score baseado em meses negativos e dívidas
	var score float64

	// Penalidade por meses negativos (0-3 pontos de penalidade)
	monthsPenalty := float64(negativeMonths) * 0.5
	if monthsPenalty > 3.0 {
		monthsPenalty = 3.0
	}

	// Penalidade por dívidas ativas (0-2 pontos de penalidade)
	debtsPenalty := float64(activeDebts) * 0.5
	if debtsPenalty > 2.0 {
		debtsPenalty = 2.0
	}

	// Score começa em 5.0 e subtrai penalidades
	score = 5.0 - monthsPenalty - debtsPenalty
	if score < 0 {
		score = 0
	}

	detail := model.ConsistencyDetail{
		NegativeMonthsCount: negativeMonths,
		ActiveDebtsCount:    activeDebts,
		Score:               score,
	}

	return score, detail, nil
}

// countNegativeMonths conta quantos meses ficaram negativos nos últimos N meses
func (s *FinancialHealthService) countNegativeMonths(ctx context.Context, userID uuid.UUID, currentPeriod time.Time, months int) (int, error) {
	count := 0

	for i := 1; i <= months; i++ {
		// Calcular período do mês a verificar
		checkPeriod := currentPeriod.AddDate(0, -i, 0)
		startDate := time.Date(checkPeriod.Year(), checkPeriod.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, -1)

		// Buscar income e expense do mês
		income, err := s.txRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startDate, endDate)
		if err != nil {
			return 0, err
		}

		expense, err := s.txRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startDate, endDate)
		if err != nil {
			return 0, err
		}

		balance := income.Sub(expense)
		if balance.LessThan(decimal.Zero) {
			count++
		}
	}

	return count, nil
}

// countActiveDebts conta credit cards com saldo negativo
func (s *FinancialHealthService) countActiveDebts(ctx context.Context, userID uuid.UUID) (int, error) {
	// Buscar todas as contas de crédito
	creditAccounts, err := s.accountRepo.GetByType(ctx, userID, model.AccountTypeCreditCard)
	if err != nil {
		return 0, fmt.Errorf("failed to get credit accounts: %w", err)
	}

	count := 0
	for _, account := range creditAccounts {
		// CurrentBalance negativo em credit card = dívida ativa
		if account.CurrentBalance.LessThan(decimal.Zero) {
			count++
		}
	}

	return count, nil
}

// GetSnapshot busca um snapshot salvo
func (s *FinancialHealthService) GetSnapshot(ctx context.Context, userID uuid.UUID, period time.Time) (*model.FinancialHealthSnapshot, error) {
	period = time.Date(period.Year(), period.Month(), 1, 0, 0, 0, 0, time.UTC)
	return s.healthRepo.GetSnapshot(ctx, userID, period)
}

// GetHistoricalSnapshots busca histórico de snapshots
func (s *FinancialHealthService) GetHistoricalSnapshots(ctx context.Context, userID uuid.UUID, months int) ([]*model.FinancialHealthSnapshot, error) {
	return s.healthRepo.GetHistoricalSnapshots(ctx, userID, months)
}
