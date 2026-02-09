"use client";

import { DollarSign, TrendingUp, TrendingDown, Wallet, PiggyBank, Target, Receipt } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { DashboardSummary } from "@/types";
import { formatCurrency, formatPercentage } from "@/lib/format";

interface SummaryCardsProps {
  summary: DashboardSummary;
}

export function SummaryCards({ summary }: SummaryCardsProps) {
  const cards = [
    {
      title: "Saldo Total",
      value: formatCurrency(summary.total_balance),
      description: `${summary.total_accounts} contas`,
      icon: Wallet,
      trend: null,
    },
    {
      title: "Receita do Mes",
      value: formatCurrency(summary.month_income),
      description: parseFloat(summary.income_change) !== 0
        ? `${parseFloat(summary.income_change) > 0 ? "+" : ""}${formatPercentage(summary.income_change)} vs mes anterior`
        : "Sem comparacao",
      icon: TrendingUp,
      trend: parseFloat(summary.income_change),
    },
    {
      title: "Despesa do Mes",
      value: formatCurrency(summary.month_expense),
      description: parseFloat(summary.expense_change) !== 0
        ? `${parseFloat(summary.expense_change) > 0 ? "+" : ""}${formatPercentage(summary.expense_change)} vs mes anterior`
        : "Sem comparacao",
      icon: TrendingDown,
      trend: -parseFloat(summary.expense_change),
    },
    {
      title: "Balanco do Mes",
      value: formatCurrency(summary.month_balance),
      description: parseFloat(summary.month_balance) >= 0 ? "Positivo" : "Negativo",
      icon: DollarSign,
      trend: parseFloat(summary.month_balance),
    },
    {
      title: "Orcamentos",
      value: `${summary.active_budgets} ativos`,
      description: summary.over_budget_count > 0
        ? `${summary.over_budget_count} acima do limite`
        : "Todos dentro do limite",
      icon: PiggyBank,
      trend: null,
    },
    {
      title: "Metas",
      value: `${summary.active_goals} em andamento`,
      description: `${formatPercentage(summary.goals_progress)} de progresso geral`,
      icon: Target,
      trend: null,
    },
    {
      title: "Contas a Pagar",
      value: `${summary.upcoming_bills} proximas`,
      description: `Total: ${formatCurrency(summary.upcoming_total)}`,
      icon: Receipt,
      trend: null,
    },
  ];

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      {cards.map((card) => (
        <Card key={card.title} className="hover:shadow-md transition-shadow">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-3">
            <CardTitle className="text-sm font-medium text-muted-foreground">{card.title}</CardTitle>
            <div className="p-2 bg-primary/10 rounded-lg">
              <card.icon className="h-4 w-4 text-primary" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight">{card.value}</div>
            <p className={`text-sm mt-2 ${
              card.trend !== null
                ? card.trend >= 0
                  ? "text-green-600 dark:text-green-500"
                  : "text-red-600 dark:text-red-500"
                : "text-muted-foreground"
            }`}>
              {card.description}
            </p>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
