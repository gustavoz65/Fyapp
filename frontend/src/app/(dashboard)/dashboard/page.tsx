"use client";

import { CashFlowChart } from "@/components/dashboard/cash-flow-chart";
import { RecentTransactions } from "@/components/dashboard/recent-transactions";
import { SummaryCards } from "@/components/dashboard/summary-cards";
import { TopCategories } from "@/components/dashboard/top-categories";
import { UpcomingBills } from "@/components/dashboard/upcoming-bills";
import { Skeleton } from "@/components/ui/skeleton";
import { api } from "@/lib/api";
import type {
  DashboardSummary,
  MonthlyIncomeExpense,
  Transaction,
} from "@/types";
import { useCallback, useEffect, useState } from "react";

export default function DashboardPage() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [monthlyData, setMonthlyData] = useState<MonthlyIncomeExpense[]>([]);
  const [upcomingBills, setUpcomingBills] = useState<Transaction[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchData = useCallback(async () => {
    try {
      const [summaryData, monthly, upcoming] = await Promise.all([
        api.get<DashboardSummary>("/dashboard"),
        api.get<MonthlyIncomeExpense[]>(
          "/dashboard/monthly-comparison?months=6",
        ),
        api.get<Transaction[]>("/transactions/upcoming"),
      ]);
      setSummary(summaryData);
      setMonthlyData(monthly ?? []);
      setUpcomingBills(upcoming ?? []);
    } catch {
      // Will show empty state
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-32" />
          ))}
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <Skeleton className="h-[350px]" />
          <Skeleton className="h-[350px]" />
        </div>
      </div>
    );
  }

  const defaultSummary: DashboardSummary = {
    total_balance: "0",
    total_accounts: 0,
    month_income: "0",
    month_expense: "0",
    month_balance: "0",
    income_change: "0",
    expense_change: "0",
    active_budgets: 0,
    over_budget_count: 0,
    budget_usage: "0",
    active_goals: 0,
    goals_progress: "0",
    upcoming_bills: 0,
    upcoming_total: "0",
    recent_transactions: [],
    top_categories: [],
  };

  const data = summary || defaultSummary;

  return (
    <div className="space-y-8 pb-8">
      <div>
        <h1 className="text-4xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Visao geral das suas financas
        </p>
      </div>

      <SummaryCards summary={data} />

      <div className="grid gap-6 md:grid-cols-2">
        <CashFlowChart data={monthlyData} />
        <TopCategories categories={data.top_categories || []} />
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <RecentTransactions transactions={data.recent_transactions || []} />
        <UpcomingBills transactions={upcomingBills} />
      </div>
    </div>
  );
}
