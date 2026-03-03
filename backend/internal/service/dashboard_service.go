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

// GetDashboardSummary returns a comprehensive dashboard summary
func (s *DashboardService) GetDashboardSummary(ctx context.Context, userID uuid.UUID) (*model.DashboardSummary, error) {
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

	// Get current month income/expense
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	monthIncome, err := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startOfMonth, endOfMonth)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get month income")
	} else {
		summary.MonthIncome = monthIncome
	}

	monthExpense, err := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startOfMonth, endOfMonth)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get month expense")
	} else {
		summary.MonthExpense = monthExpense
	}

	summary.MonthBalance = summary.MonthIncome.Sub(summary.MonthExpense)

	// Compare with last month
	startOfLastMonth := startOfMonth.AddDate(0, -1, 0)
	endOfLastMonth := startOfMonth.AddDate(0, 0, -1)

	lastMonthIncome, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startOfLastMonth, endOfLastMonth)
	lastMonthExpense, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startOfLastMonth, endOfLastMonth)

	if !lastMonthIncome.IsZero() {
		summary.IncomeChange = summary.MonthIncome.Sub(lastMonthIncome).Div(lastMonthIncome).Mul(decimal.NewFromInt(100))
	}
	if !lastMonthExpense.IsZero() {
		summary.ExpenseChange = summary.MonthExpense.Sub(lastMonthExpense).Div(lastMonthExpense).Mul(decimal.NewFromInt(100))
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

	// Get recent transactions
	recentTransactions, err := s.transactionRepo.GetRecentTransactions(ctx, userID, 5)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to get recent transactions")
	} else {
		summary.RecentTransactions = make([]model.Transaction, len(recentTransactions))
		for i, tx := range recentTransactions {
			summary.RecentTransactions[i] = *tx
		}
	}

	// Get top expense categories
	topCategories, err := s.transactionRepo.GetSumByCategory(ctx, userID, startOfMonth, endOfMonth)
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

// GetIncomeVsExpenseReport generates an income vs expense comparison report
func (s *DashboardService) GetIncomeVsExpenseReport(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) (*model.IncomeVsExpenseReport, error) {
	report := &model.IncomeVsExpenseReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Get totals
	totalIncome, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startDate, endDate)
	totalExpense, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startDate, endDate)

	report.TotalIncome = totalIncome
	report.TotalExpense = totalExpense
	report.NetIncome = totalIncome.Sub(totalExpense)

	// Calculate savings rate
	if !totalIncome.IsZero() {
		report.SavingsRate = report.NetIncome.Div(totalIncome).Mul(decimal.NewFromInt(100))
	}

	// Generate monthly breakdown
	current := startDate
	for current.Before(endDate) || current.Equal(endDate) {
		monthStart := time.Date(current.Year(), current.Month(), 1, 0, 0, 0, 0, time.Local)
		monthEnd := monthStart.AddDate(0, 1, -1)

		if monthEnd.After(endDate) {
			monthEnd = endDate
		}

		monthIncome, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, monthStart, monthEnd)
		monthExpense, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, monthStart, monthEnd)
		netIncome := monthIncome.Sub(monthExpense)

		var savingsRate decimal.Decimal
		if !monthIncome.IsZero() {
			savingsRate = netIncome.Div(monthIncome).Mul(decimal.NewFromInt(100))
		}

		report.MonthlyData = append(report.MonthlyData, model.MonthlyIncomeExpense{
			Month:       current.Format("January"),
			Year:        current.Year(),
			Income:      monthIncome,
			Expense:     monthExpense,
			NetIncome:   netIncome,
			SavingsRate: savingsRate,
		})

		current = current.AddDate(0, 1, 0)
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

// GetMonthlyComparison returns a comparison of the current month with previous months
func (s *DashboardService) GetMonthlyComparison(ctx context.Context, userID uuid.UUID, months int) ([]model.MonthlyIncomeExpense, error) {
	var results []model.MonthlyIncomeExpense

	now := time.Now()
	for i := months - 1; i >= 0; i-- {
		monthDate := now.AddDate(0, -i, 0)
		startDate := time.Date(monthDate.Year(), monthDate.Month(), 1, 0, 0, 0, 0, time.Local)
		endDate := startDate.AddDate(0, 1, -1)

		income, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeIncome, startDate, endDate)
		expense, _ := s.transactionRepo.GetSumByType(ctx, userID, model.TransactionTypeExpense, startDate, endDate)
		netIncome := income.Sub(expense)

		var savingsRate decimal.Decimal
		if !income.IsZero() {
			savingsRate = netIncome.Div(income).Mul(decimal.NewFromInt(100))
		}

		results = append(results, model.MonthlyIncomeExpense{
			Month:       startDate.Format("January"),
			Year:        startDate.Year(),
			Income:      income,
			Expense:     expense,
			NetIncome:   netIncome,
			SavingsRate: savingsRate,
		})
	}

	return results, nil
}
