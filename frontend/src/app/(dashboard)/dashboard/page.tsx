"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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

type Period =
  | "this_month"
  | "last_month"
  | "last_3_months"
  | "last_6_months"
  | "all_time"
  | "custom";

interface DateRange {
  start: string; // formato YYYY-MM-DD
  end: string;
}

function toDateStr(d: Date): string {
  return d.toISOString().slice(0, 10);
}

function getDateRange(period: Period, custom: DateRange): DateRange {
  const now = new Date();
  const y = now.getFullYear();
  const m = now.getMonth();

  if (period === "custom") return custom;
  if (period === "all_time") {
    return { start: "2000-01-01", end: toDateStr(now) };
  }
  if (period === "this_month") {
    const start = new Date(y, m, 1);
    const end = new Date(y, m + 1, 0);
    return { start: toDateStr(start), end: toDateStr(end) };
  }
  if (period === "last_month") {
    const start = new Date(y, m - 1, 1);
    const end = new Date(y, m, 0);
    return { start: toDateStr(start), end: toDateStr(end) };
  }
  if (period === "last_3_months") {
    const start = new Date(y, m - 2, 1);
    const end = new Date(y, m + 1, 0);
    return { start: toDateStr(start), end: toDateStr(end) };
  }
  const start = new Date(y, m - 5, 1);
  const end = new Date(y, m + 1, 0);
  return { start: toDateStr(start), end: toDateStr(end) };
}

const PERIOD_OPTIONS: { value: Period; label: string }[] = [
  { value: "this_month", label: "Este mês" },
  { value: "last_month", label: "Mês passado" },
  { value: "last_3_months", label: "3 meses" },
  { value: "last_6_months", label: "6 meses" },
  { value: "all_time", label: "Todo período" },
  { value: "custom", label: "Personalizado" },
];

function periodComparisonLabel(period: Period): string {
  return period === "this_month" || period === "last_month"
    ? "vs mês anterior"
    : "vs período anterior";
}

type HealthStatus = "good" | "warning" | "critical";

const statusColor: Record<HealthStatus, string> = {
  good: "text-green-500",
  warning: "text-yellow-500",
  critical: "text-red-500",
};

function TrendIcon({
  status,
  trend,
}: {
  status: HealthStatus;
  trend: "up" | "down";
}) {
  const cls = `h-3 w-3 mr-1 ${statusColor[status]}`;
  const showUp = status === "good" || (status === "warning" && trend === "up");
  return showUp ? <TrendingUp className={cls} /> : <TrendingDown className={cls} />;
}

export default function DashboardPage() {
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [monthlyData, setMonthlyData] = useState<MonthlyIncomeExpense[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const [period, setPeriod] = useState<Period>("this_month");
  const [customRange, setCustomRange] = useState<DateRange>({
    start: toDateStr(new Date(new Date().getFullYear(), new Date().getMonth(), 1)),
    end: toDateStr(new Date()),
  });

  const fetchData = useCallback(
    async (p: Period, custom: DateRange) => {
      setIsLoading(true);
      try {
        const { start, end } = getDateRange(p, custom);
        const [summaryData, monthly] = await Promise.all([
          api.get<DashboardSummary>(
            `/dashboard?start_date=${start}&end_date=${end}`,
          ),
          api.get<MonthlyIncomeExpense[]>(
            "/dashboard/monthly-comparison?months=6",
          ),
        ]);
        setSummary(summaryData);
        setMonthlyData(monthly ?? []);
      } catch {
        // ignora erro — componente exibe estado vazio
      } finally {
        setIsLoading(false);
      }
    },
    [],
  );

  useEffect(() => {
    fetchData(period, customRange);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function handlePeriodChange(p: Period) {
    setPeriod(p);
    if (p !== "custom") {
      fetchData(p, customRange);
    }
  }

  function handleCustomApply() {
    fetchData("custom", customRange);
  }

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

  const totalBalance = Number.parseFloat(data.total_balance || "0");
  const isNegativeBalance = totalBalance < 0;

  const incomeChange = Number.parseFloat(data.income_change || "0");
  const expenseChange = Number.parseFloat(data.expense_change || "0");
  const goalsProgress = Number.parseFloat(data.goals_progress || "0");

  let balanceStatus: HealthStatus = "warning";
  if (isNegativeBalance) balanceStatus = "critical";
  else if (incomeChange >= 0) balanceStatus = "good";

  let incomeStatus: HealthStatus = "critical";
  if (incomeChange > 0) incomeStatus = "good";
  else if (incomeChange > -30) incomeStatus = "warning";

  let expenseStatus: HealthStatus = "good";
  if (expenseChange >= 30) expenseStatus = "critical";
  else if (expenseChange >= 0) expenseStatus = "warning";

  let goalsStatus: HealthStatus = "critical";
  if (goalsProgress > 50) goalsStatus = "good";
  else if (goalsProgress > 20) goalsStatus = "warning";

  const compLabel = periodComparisonLabel(period);

  const metrics = [
    {
      title: "Saldo Total",
      value: formatCurrency(data.total_balance),
      change: `${incomeChange > 0 ? "+" : ""}${incomeChange.toFixed(1)}%`,
      trend: incomeChange >= 0 ? ("up" as const) : ("down" as const),
      icon: Wallet,
      isNegative: isNegativeBalance,
      status: balanceStatus,
    },
    {
      title: "Receitas do Período",
      value: formatCurrency(data.month_income),
      change: data.income_change
        ? `${incomeChange > 0 ? "+" : ""}${incomeChange.toFixed(1)}%`
        : "+0%",
      trend: incomeChange >= 0 ? ("up" as const) : ("down" as const),
      icon: TrendingUpIcon,
      status: incomeStatus,
    },
    {
      title: "Despesas do Período",
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

  const chartData = monthlyData
    .map((item) => ({
      month: new Date(item.month).toLocaleDateString("pt-BR", {
        month: "short",
      }),
      income: Number.parseFloat(item.income || "0"),
      expense: Number.parseFloat(item.expense || "0"),
    }))
    .reverse();

  const recentTransactions = data.recent_transactions || [];

  return (
    <div className="space-y-8 pb-8">
      {/* Header + seletor de período */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground mt-2">
            Visão geral das suas finanças
          </p>
        </div>

        {/* Seletor de período */}
        <div className="flex flex-col gap-2 sm:items-end">
          <div className="flex flex-wrap gap-1">
            {PERIOD_OPTIONS.map((opt) => (
              <Button
                key={opt.value}
                size="sm"
                variant={period === opt.value ? "default" : "outline"}
                onClick={() => handlePeriodChange(opt.value)}
                className="shrink-0"
              >
                {opt.label}
              </Button>
            ))}
          </div>

          {period === "custom" && (
            <div className="flex flex-wrap items-end gap-2 mt-1">
              <div className="space-y-1">
                <Label className="text-xs">De</Label>
                <Input
                  type="date"
                  className="h-8 w-36 text-xs"
                  value={customRange.start}
                  onChange={(e) =>
                    setCustomRange((r) => ({ ...r, start: e.target.value }))
                  }
                />
              </div>
              <div className="space-y-1">
                <Label className="text-xs">Até</Label>
                <Input
                  type="date"
                  className="h-8 w-36 text-xs"
                  value={customRange.end}
                  onChange={(e) =>
                    setCustomRange((r) => ({ ...r, end: e.target.value }))
                  }
                />
              </div>
              <Button size="sm" onClick={handleCustomApply}>
                Aplicar
              </Button>
            </div>
          )}
        </div>
      </div>

      {/* Metric Cards */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        {metrics.map((metric) => {
          const Icon = metric.icon;
          const isNegative = "isNegative" in metric && metric.isNegative;
          return (
            <Card
              key={metric.title}
              className={`hover:shadow-md transition-shadow border ${isNegative ? "border-red-500 bg-red-50 dark:bg-red-950/20" : "border-border"}`}
            >
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-semibold">
                  {metric.title}
                </CardTitle>
                <Icon
                  className={`h-5 w-5 ${isNegative ? "text-red-500" : "text-muted-foreground"}`}
                />
              </CardHeader>
              <CardContent>
                <div
                  className={`text-2xl font-bold ${isNegative ? "text-red-600 dark:text-red-400" : "text-foreground"}`}
                >
                  {metric.value}
                </div>
                <div className="flex items-center text-xs text-muted-foreground mt-1">
                  <TrendIcon status={metric.status} trend={metric.trend} />
                  <span className={statusColor[metric.status]}>
                    {metric.change}
                  </span>
                  <span className="ml-1">{compLabel}</span>
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
                  <stop offset="5%" stopColor="#00A859" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#00A859" stopOpacity={0} />
                </linearGradient>
                <linearGradient id="colorExpense" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#EF4444" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#EF4444" stopOpacity={0} />
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
                stroke="#00A859"
                fillOpacity={1}
                fill="url(#colorIncome)"
                strokeWidth={2}
                name="Receitas"
              />
              <Area
                type="monotone"
                dataKey="expense"
                stroke="#EF4444"
                fillOpacity={1}
                fill="url(#colorExpense)"
                strokeWidth={2}
                name="Despesas"
              />
            </AreaChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Tabs */}
      <Tabs defaultValue="overview" className="space-y-4">
        <TabsList>
          <TabsTrigger value="overview">Visão Geral</TabsTrigger>
          <TabsTrigger value="performance">Desempenho</TabsTrigger>
          <TabsTrigger value="categories">Categorias</TabsTrigger>
          <TabsTrigger value="reports">Relatórios</TabsTrigger>
        </TabsList>

        {/* Transações Recentes */}
        <TabsContent value="overview" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Transações Recentes</CardTitle>
              <CardDescription>
                Últimas movimentações no período selecionado
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="overflow-x-auto">
                <Table>
                  <TableHeader>
                    <TableRow>
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
                          colSpan={6}
                          className="text-center text-muted-foreground py-8"
                        >
                          Nenhuma transação no período selecionado
                        </TableCell>
                      </TableRow>
                    ) : (
                      recentTransactions.map((transaction) => (
                        <TableRow
                          key={transaction.id}
                          className="hover:bg-muted/50"
                        >
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
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Desempenho */}
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

        {/* Categorias */}
        <TabsContent value="categories" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Categorias Principais</CardTitle>
              <CardDescription>
                Onde você mais gastou no período selecionado
              </CardDescription>
            </CardHeader>
            <CardContent>
              {!data.top_categories || data.top_categories.length === 0 ? (
                <div className="text-sm text-muted-foreground py-8 text-center">
                  Nenhuma despesa categorizada no período selecionado
                </div>
              ) : (
                <div className="space-y-5">
                  {data.top_categories.map((cat, i) => {
                    const colors = [
                      "#ef4444",
                      "#f97316",
                      "#eab308",
                      "#22c55e",
                      "#06b6d4",
                      "#6366f1",
                      "#a855f7",
                      "#ec4899",
                    ];
                    const color = colors[i % colors.length];
                    const pct = Math.min(
                      Number.parseFloat(cat.percentage || "0"),
                      100,
                    );
                    return (
                      <div
                        key={cat.category_id || cat.category_name}
                        className="space-y-1.5"
                      >
                        <div className="flex items-center justify-between text-sm">
                          <div className="flex items-center gap-2">
                            <span
                              className="inline-block h-2.5 w-2.5 rounded-full shrink-0"
                              style={{ backgroundColor: color }}
                            />
                            <span className="font-medium">
                              {cat.category_name}
                            </span>
                            <span className="text-xs text-muted-foreground">
                              {cat.count}{" "}
                              {cat.count === 1 ? "transação" : "transações"}
                            </span>
                          </div>
                          <div className="flex items-center gap-3 shrink-0">
                            <span className="text-xs text-muted-foreground">
                              {pct.toFixed(1)}%
                            </span>
                            <span className="font-semibold text-red-600 dark:text-red-400">
                              {formatCurrency(cat.amount)}
                            </span>
                          </div>
                        </div>
                        <div className="w-full bg-muted rounded-full h-1.5">
                          <div
                            className="h-1.5 rounded-full transition-all"
                            style={{
                              width: `${pct}%`,
                              backgroundColor: color,
                            }}
                          />
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        {/* Relatórios */}
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
