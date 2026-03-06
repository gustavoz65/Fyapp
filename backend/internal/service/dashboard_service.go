package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/gustavoz65/Fyapp/internal/model"
	"github.com/gustavoz65/Fyapp/internal/repository"
)

type DashboardService struct {
	accountRepo     *repository.BankAccountRepository
	transactionRepo *repository.TransactionRepository
	budgetRepo      *repository.BudgetRepository
	goalRepo        *repository.GoalRepository
	logger          *zerolog.Logger
}

func NewDashboardService(
	accountRepo *repository.BankAccountRepository,
	transactionRepo *repository.TransactionRepository,
	budgetRepo *repository.BudgetRepository,
	goalRepo *repository.GoalRepository,
	logger *zerolog.Logger,
) *DashboardService {
	return &DashboardService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		budgetRepo:      budgetRepo,
		goalRepo:        goalRepo,
		logger:          logger,
	}
}

// GetDashboardSummary returns a comprehensive dashboard summary for the given date range.
// The comparison period is automatically computed as the same-length interval before startDate.
func (s *DashboardService) GetDashboardSummary(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*model.DashboardSummary, error) {
	summary := &model.DashboardSummary{}

	// Get total balance
	totalBalance, err := s.accountRepo.GetTotalBalance(ctx, userID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get total balance")
	} else {
		summary.TotalBalance = totalBalance
	}

	// Get account count
	accountCount, err := s.accountRepo.CountAccounts(ctx, userID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to count accounts")
	} else {
		summary.TotalInAccounts = int(accountCount)
	}

	// Income/expense for the requested period
	monthIncome, err := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startDate, endDate)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get period income")
	} else {
		summary.MonthIncome = monthIncome
	}

	monthExpense, err := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startDate, endDate)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get period expense")
	} else {
		summary.MonthExpense = monthExpense
	}

	summary.MonthBalance = summary.MonthIncome.Sub(summary.MonthExpense)

	// Compute previous period (same duration, shifted back)
	periodDuration := endDate.Sub(startDate)
	prevEnd := startDate.Add(-time.Second)
	prevStart := prevEnd.Add(-periodDuration)

	prevIncome, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, prevStart, prevEnd)
	prevExpense, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, prevStart, prevEnd)

	if !prevIncome.IsZero() {
		summary.IncomeChange = summary.MonthIncome.Sub(prevIncome).Div(prevIncome).Mul(decimal.NewFromInt(100))
	}
	if !prevExpense.IsZero() {
		summary.ExpenseChange = summary.MonthExpense.Sub(prevExpense).Div(prevExpense).Mul(decimal.NewFromInt(100))
	}

	// Get budget status
	activeBudgets, err := s.budgetRepo.GetActiveByUser(ctx, userID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get active budgets")
	} else {
		summary.ActiveBudgets = len(activeBudgets)
		var totalBudgeted, totalSpent decimal.Decimal
		for _, budget := range activeBudgets {
			totalBudgeted = totalBudgeted.Add(budget.Amount)
			totalSpent = totalSpent.Add(budget.SpentAmount)
			if budget.IsOverBudget() {
				summary.OverBudgetCount++
			}
		}
		if !totalBudgeted.IsZero() {
			summary.BudgetUsage = totalSpent.Div(totalBudgeted).Mul(decimal.NewFromInt(100))
		}
	}

	// Get goal progress
	goalSummary, err := s.goalRepo.GetGoalSummary(ctx, userID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get goal summary")
	} else {
		summary.ActiveGoals = goalSummary.InProgressGoals
		summary.GoalsProgress = goalSummary.OverallProgress
	}

	// Get upcoming bills
	upcomingBills, err := s.transactionRepo.GetUpcomingBills(ctx, userID, 7)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get upcoming bills")
	} else {
		summary.UpcomingBills = len(upcomingBills)
		for _, bill := range upcomingBills {
			summary.UpcomingTotal = summary.UpcomingTotal.Add(bill.Amount)
		}
	}

	// Get recent transactions within the period (with category info)
	recentTransactions, err := s.transactionRepo.GetRecentByPeriod(ctx, userID, startDate, endDate, 10)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get recent transactions by period")
	} else {
		summary.RecentTransactions = make([]model.Transaction, len(recentTransactions))
		for i, tx := range recentTransactions {
			summary.RecentTransactions[i] = *tx
		}
	}

	// Get top expense categories for the period
	topCategories, err := s.transactionRepo.GetSumByCategory(ctx, userID, startDate, endDate)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get top categories")
	} else {
		if len(topCategories) > 5 {
			topCategories = topCategories[:5]
		}
		summary.TopCategories = topCategories
	}

	return summary, nil
}

// GetCashFlowReport generates a cash flow report
func (s *DashboardService) GetCashFlowReport(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*model.CashFlowReport, error) {
	report := &model.CashFlowReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Get total income
	totalIncome, err := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get total income: %w", err)
	}
	report.TotalIncome = totalIncome

	// Get total expense
	totalExpense, err := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get total expense: %w", err)
	}
	report.TotalExpense = totalExpense

	// Calculate net cash flow
	report.NetCashFlow = totalIncome.Sub(totalExpense)

	// Get opening balance (total balance at start of period)
	// This would require historical balance tracking - for now, we'll calculate from current balance
	currentBalance, _ := s.accountRepo.GetTotalBalance(ctx, userID)
	report.ClosingBalance = currentBalance
	report.OpeningBalance = currentBalance.Sub(report.NetCashFlow)

	// Get income by category
	report.ExpenseByCategory, _ = s.transactionRepo.GetSumByCategory(ctx, userID, startDate, endDate)

	return report, nil
}

// GetIncomeVsExpenseReport generates an income vs expense comparison report.
// Usa uma única query GROUP BY para substituir o loop anterior de N*2 queries.
func (s *DashboardService) GetIncomeVsExpenseReport(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*model.IncomeVsExpenseReport, error) {
	report := &model.IncomeVsExpenseReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	rawData, err := s.transactionRepo.GetMonthlyBreakdown(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly breakdown: %w", err)
	}

	// Indexa por (ano, mês) para lookup O(1)
	type monthKey = [2]int
	byMonth := make(map[monthKey]repository.MonthlyRawData, len(rawData))
	for _, d := range rawData {
		byMonth[monthKey{d.Year, d.Month}] = d
	}

	// Percorre todos os meses do intervalo, preenchendo zeros onde não há transações
	current := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.Local)
	endMonth := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.Local)

	for !current.After(endMonth) {
		d := byMonth[monthKey{current.Year(), int(current.Month())}]
		netIncome := d.Income.Sub(d.Expense)

		var savingsRate decimal.Decimal
		if !d.Income.IsZero() {
			savingsRate = netIncome.Div(d.Income).Mul(decimal.NewFromInt(100))
		}

		report.MonthlyData = append(report.MonthlyData, model.MonthlyIncomeExpense{
			Month:       current.Format("January"),
			Year:        current.Year(),
			Income:      d.Income,
			Expense:     d.Expense,
			NetIncome:   netIncome,
			SavingsRate: savingsRate,
		})

		report.TotalIncome = report.TotalIncome.Add(d.Income)
		report.TotalExpense = report.TotalExpense.Add(d.Expense)

		current = current.AddDate(0, 1, 0)
	}

	report.NetIncome = report.TotalIncome.Sub(report.TotalExpense)
	if !report.TotalIncome.IsZero() {
		report.SavingsRate = report.NetIncome.Div(report.TotalIncome).Mul(decimal.NewFromInt(100))
	}

	// Calculate IncomeGrowth and ExpenseGrowth comparing first and last months
	if len(report.MonthlyData) >= 2 {
		first := report.MonthlyData[0]
		last := report.MonthlyData[len(report.MonthlyData)-1]
		if !first.Income.IsZero() {
			report.IncomeGrowth = last.Income.Sub(first.Income).Div(first.Income).Mul(decimal.NewFromInt(100))
		}
		if !first.Expense.IsZero() {
			report.ExpenseGrowth = last.Expense.Sub(first.Expense).Div(first.Expense).Mul(decimal.NewFromInt(100))
		}
	}

	return report, nil
}

// GetAccountBalances returns balances for all accounts
func (s *DashboardService) GetAccountBalances(ctx context.Context, userID uuid.UUID) ([]*model.BankAccount, error) {
	accounts, err := s.accountRepo.GetAllByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balances: %w", err)
	}
	return accounts, nil
}

// GetMonthlyComparison returns a comparison of the current month with previous months.
// Usa uma única query GROUP BY para substituir o loop anterior de N*2 queries.
func (s *DashboardService) GetMonthlyComparison(ctx context.Context, userID uuid.UUID, months int) ([]model.MonthlyIncomeExpense, error) {
	now := time.Now()
	startDate := time.Date(now.AddDate(0, -(months-1), 0).Year(), now.AddDate(0, -(months-1), 0).Month(), 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	rawData, err := s.transactionRepo.GetMonthlyBreakdown(ctx, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly breakdown: %w", err)
	}

	type monthKey = [2]int
	byMonth := make(map[monthKey]repository.MonthlyRawData, len(rawData))
	for _, d := range rawData {
		byMonth[monthKey{d.Year, d.Month}] = d
	}

	results := make([]model.MonthlyIncomeExpense, 0, months)
	for i := months - 1; i >= 0; i-- {
		monthDate := now.AddDate(0, -i, 0)
		d := byMonth[monthKey{monthDate.Year(), int(monthDate.Month())}]
		netIncome := d.Income.Sub(d.Expense)

		var savingsRate decimal.Decimal
		if !d.Income.IsZero() {
			savingsRate = netIncome.Div(d.Income).Mul(decimal.NewFromInt(100))
		}

		results = append(results, model.MonthlyIncomeExpense{
			Month:       monthDate.Format("January"),
			Year:        monthDate.Year(),
			Income:      d.Income,
			Expense:     d.Expense,
			NetIncome:   netIncome,
			SavingsRate: savingsRate,
		})
	}

	return results, nil
}
