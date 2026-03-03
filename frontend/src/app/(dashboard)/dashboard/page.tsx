"use client";

import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { api } from "@/lib/api";
import { formatCurrency } from "@/lib/format";
import type { DashboardSummary, MonthlyIncomeExpense } from "@/types";
import {
  CreditCard,
  DollarSign,
  TrendingDown,
  TrendingUp,
  TrendingUpIcon,
  Wallet,
} from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

export default function DashboardPage() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [monthlyData, setMonthlyData] = useState<MonthlyIncomeExpense[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchData = useCallback(async () => {
    try {
      const [summaryData, monthly] = await Promise.all([
        api.get<DashboardSummary>("/dashboard"),
        api.get<MonthlyIncomeExpense[]>(
          "/dashboard/monthly-comparison?months=6",
        ),
      ]);
      setSummary(summaryData);
      setMonthlyData(monthly ?? []);
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
      <div className="space-y-8 pb-8">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground mt-2">
            Visão geral das suas finanças
          </p>
        </div>
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-32" />
          ))}
        </div>
        <Skeleton className="h-[400px]" />
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

  // Calcular variação do saldo comparando com o mês anterior
  const currentMonthBalance = parseFloat(data.month_income || "0") - parseFloat(data.month_expense || "0");
  const previousMonthData = monthlyData.length >= 2 ? monthlyData[monthlyData.length - 2] : null;
  const previousMonthBalance = previousMonthData
    ? parseFloat(previousMonthData.income || "0") - parseFloat(previousMonthData.expense || "0")
    : 0;

  const balanceChangePercent = previousMonthBalance !== 0
    ? ((currentMonthBalance - previousMonthBalance) / Math.abs(previousMonthBalance)) * 100
    : currentMonthBalance !== 0 ? 100 : 0;

  const totalBalance = parseFloat(data.total_balance || "0");
  const isNegativeBalance = totalBalance < 0;

  type HealthStatus = "good" | "warning" | "critical";

  // Saldo Total
  // critical → saldo negativo | warning → positivo mas caindo | good → positivo e crescendo
  const balanceStatus: HealthStatus =
    isNegativeBalance ? "critical" : balanceChangePercent < 0 ? "warning" : "good";

  // Receitas do Mês
  // good → crescendo | warning → caindo até -30% | critical → caindo mais de -30%
  const incomeChange = parseFloat(data.income_change || "0");
  const incomeStatus: HealthStatus =
    incomeChange > 0 ? "good" : incomeChange > -30 ? "warning" : "critical";

  // Despesas do Mês (lógica invertida: gastar menos é bom)
  // good → caindo | warning → subindo até +30% | critical → subindo mais de +30%
  const expenseChange = parseFloat(data.expense_change || "0");
  const expenseStatus: HealthStatus =
    expenseChange < 0 ? "good" : expenseChange < 30 ? "warning" : "critical";

  // Metas Ativas
  // good → >50% concluído | warning → 20-50% | critical → <20%
  const goalsProgress = parseFloat(data.goals_progress || "0");
  const goalsStatus: HealthStatus =
    goalsProgress > 50 ? "good" : goalsProgress > 20 ? "warning" : "critical";

  const metrics = [
    {
      title: "Saldo Total",
      value: formatCurrency(data.total_balance),
      change: `${balanceChangePercent > 0 ? "+" : ""}${balanceChangePercent.toFixed(1)}%`,
      trend: balanceChangePercent >= 0 ? ("up" as const) : ("down" as const),
      icon: Wallet,
      isNegative: isNegativeBalance,
      status: balanceStatus,
    },
    {
      title: "Receitas do Mês",
      value: formatCurrency(data.month_income),
      change: data.income_change
        ? `${incomeChange > 0 ? "+" : ""}${incomeChange.toFixed(1)}%`
        : "+0%",
      trend: incomeChange >= 0 ? ("up" as const) : ("down" as const),
      icon: TrendingUpIcon,
      status: incomeStatus,
    },
    {
      title: "Despesas do Mês",
      value: formatCurrency(data.month_expense),
      change: data.expense_change
        ? `${expenseChange > 0 ? "+" : ""}${expenseChange.toFixed(1)}%`
        : "+0%",
      trend: expenseChange > 0 ? ("up" as const) : ("down" as const),
      icon: CreditCard,
      status: expenseStatus,
    },
    {
      title: "Metas Ativas",
      value: data.active_goals.toString(),
      change: `${data.goals_progress}% concluído`,
      trend: goalsProgress > 50 ? ("up" as const) : ("down" as const),
      icon: DollarSign,
      status: goalsStatus,
    },
  ];

  // Transform monthly data for the chart
  const chartData = monthlyData
    .map((item) => ({
      month: new Date(item.month).toLocaleDateString("pt-BR", {
        month: "short",
      }),
      income: parseFloat(item.income || "0"),
      expense: parseFloat(item.expense || "0"),
    }))
    .reverse();

  const recentTransactions = data.recent_transactions || [];

  return (
    <div className="space-y-8 pb-8">
      <div>
        <h1 className="text-4xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Visão geral das suas finanças
        </p>
      </div>

      {/* Metric Cards */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        {metrics.map((metric, index) => {
          const Icon = metric.icon;
          const isNegative = "isNegative" in metric && metric.isNegative;
          return (
            <Card key={index} className={`hover:shadow-md transition-shadow ${isNegative ? "border-red-500 bg-red-50 dark:bg-red-950/20" : ""}`}>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {metric.title}
                </CardTitle>
                <Icon className={`h-4 w-4 ${isNegative ? "text-red-500" : "text-muted-foreground"}`} />
              </CardHeader>
              <CardContent>
                <div className={`text-2xl font-bold ${isNegative ? "text-red-600 dark:text-red-400" : ""}`}>
                  {metric.value}
                </div>
                <div className="flex items-center text-xs text-muted-foreground mt-1">
                  {(() => {
                    const colorMap = {
                      good: "text-green-500",
                      warning: "text-yellow-500",
                      critical: "text-red-500",
                    };
                    const color = colorMap[metric.status];
                    return metric.trend === "up" ? (
                      <TrendingUp className={`h-3 w-3 mr-1 ${color}`} />
                    ) : (
                      <TrendingDown className={`h-3 w-3 mr-1 ${color}`} />
                    );
                  })()}
                  <span
                    className={{
                      good: "text-green-500",
                      warning: "text-yellow-500",
                      critical: "text-red-500",
                    }[metric.status]}
                  >
                    {metric.change}
                  </span>
                  <span className="ml-1">do mês anterior</span>
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      {/* Large Area Chart */}
      <Card>
        <CardHeader>
          <CardTitle>Fluxo Financeiro</CardTitle>
          <CardDescription>
            Receitas e despesas dos últimos 6 meses
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={350}>
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="colorIncome" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="5%"
                    stopColor="#16a34a"
                    stopOpacity={0.3}
                  />
                  <stop
                    offset="95%"
                    stopColor="#16a34a"
                    stopOpacity={0}
                  />
                </linearGradient>
                <linearGradient id="colorExpense" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="5%"
                    stopColor="#ef4444"
                    stopOpacity={0.3}
                  />
                  <stop
                    offset="95%"
                    stopColor="#ef4444"
                    stopOpacity={0}
                  />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
              <XAxis
                dataKey="month"
                className="text-xs"
                tick={{ fill: "hsl(var(--muted-foreground))" }}
              />
              <YAxis
                className="text-xs"
                tick={{ fill: "hsl(var(--muted-foreground))" }}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: "hsl(var(--popover))",
                  border: "1px solid hsl(var(--border))",
                  borderRadius: "0.5rem",
                }}
                formatter={(value: number) => formatCurrency(value)}
              />
              <Area
                type="monotone"
                dataKey="income"
                stroke="#16a34a"
                fillOpacity={1}
                fill="url(#colorIncome)"
                strokeWidth={2}
                name="Receitas"
              />
              <Area
                type="monotone"
                dataKey="expense"
                stroke="#ef4444"
                fillOpacity={1}
                fill="url(#colorExpense)"
                strokeWidth={2}
                name="Despesas"
              />
            </AreaChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Tabs with Table */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Visão Geral</TabsTrigger>
          <TabsTrigger value="performance">Desempenho</TabsTrigger>
          <TabsTrigger value="categories">Categorias</TabsTrigger>
          <TabsTrigger value="reports">Relatórios</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Transações Recentes</CardTitle>
              <CardDescription>
                Últimas movimentações financeiras
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-12">
                      <Checkbox />
                    </TableHead>
                    <TableHead>Descrição</TableHead>
                    <TableHead>Tipo</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Valor</TableHead>
                    <TableHead>Categoria</TableHead>
                    <TableHead>Conta</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {recentTransactions.length === 0 ? (
                    <TableRow>
                      <TableCell
                        colSpan={7}
                        className="text-center text-muted-foreground py-8"
                      >
                        Nenhuma transação recente
                      </TableCell>
                    </TableRow>
                  ) : (
                    recentTransactions.map((transaction) => (
                      <TableRow
                        key={transaction.id}
                        className="hover:bg-muted/50"
                      >
                        <TableCell>
                          <Checkbox />
                        </TableCell>
                        <TableCell className="font-medium">
                          {transaction.description}
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={
                              transaction.type === "income"
                                ? "default"
                                : "secondary"
                            }
                          >
                            {transaction.type === "income"
                              ? "Receita"
                              : "Despesa"}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={
                              transaction.is_paid ? "outline" : "secondary"
                            }
                          >
                            {transaction.is_paid ? "Pago" : "Pendente"}
                          </Badge>
                        </TableCell>
                        <TableCell
                          className={
                            transaction.type === "income"
                              ? "text-green-600"
                              : "text-red-600"
                          }
                        >
                          {transaction.type === "income" ? "+" : "-"}
                          {formatCurrency(transaction.amount)}
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {transaction.category?.name || "-"}
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {transaction.bank_account?.name || "-"}
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="performance" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Desempenho Mensal</CardTitle>
              <CardDescription>
                Análise do desempenho financeiro mensal
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground py-8 text-center">
                Dados de desempenho serão exibidos aqui
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="categories" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Categorias Principais</CardTitle>
              <CardDescription>
                Distribuição de gastos por categoria
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground py-8 text-center">
                Dados de categorias serão exibidos aqui
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="reports" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Relatórios</CardTitle>
              <CardDescription>
                Relatórios financeiros detalhados
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground py-8 text-center">
                Relatórios serão exibidos aqui
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
