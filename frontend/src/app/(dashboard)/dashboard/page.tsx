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
            Visao geral das suas financas
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

  const metrics = [
    {
      title: "Saldo Total",
      value: formatCurrency(data.total_balance),
      change: "+12.5%",
      trend: "up" as const,
      icon: Wallet,
    },
    {
      title: "Receitas do Mes",
      value: formatCurrency(data.month_income),
      change: data.income_change
        ? `${parseFloat(data.income_change) > 0 ? "+" : ""}${parseFloat(data.income_change).toFixed(1)}%`
        : "+0%",
      trend:
        parseFloat(data.income_change || "0") > 0
          ? ("up" as const)
          : ("down" as const),
      icon: TrendingUpIcon,
    },
    {
      title: "Despesas do Mes",
      value: formatCurrency(data.month_expense),
      change: data.expense_change
        ? `${parseFloat(data.expense_change) > 0 ? "+" : ""}${parseFloat(data.expense_change).toFixed(1)}%`
        : "+0%",
      trend:
        parseFloat(data.expense_change || "0") > 0
          ? ("down" as const)
          : ("up" as const),
      icon: CreditCard,
    },
    {
      title: "Metas Ativas",
      value: data.active_goals.toString(),
      change: `${data.goals_progress}% concluido`,
      trend:
        parseFloat(data.goals_progress || "0") > 50
          ? ("up" as const)
          : ("down" as const),
      icon: DollarSign,
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

  // Sample table data - in production this would come from API
  const tableData = [
    {
      id: "1",
      header: "Supermercado",
      section: "Despesa",
      status: "Pago",
      target: formatCurrency(500),
      limit: formatCurrency(600),
      reviewer: "Voce",
    },
    {
      id: "2",
      header: "Salario",
      section: "Receita",
      status: "Recebido",
      target: formatCurrency(5000),
      limit: "-",
      reviewer: "Sistema",
    },
    {
      id: "3",
      header: "Aluguel",
      section: "Despesa",
      status: "Pendente",
      target: formatCurrency(1200),
      limit: formatCurrency(1200),
      reviewer: "Voce",
    },
    {
      id: "4",
      header: "Freelance",
      section: "Receita",
      status: "Recebido",
      target: formatCurrency(1500),
      limit: "-",
      reviewer: "Sistema",
    },
  ];

  return (
    <div className="space-y-8 pb-8">
      <div>
        <h1 className="text-4xl font-bold tracking-tight">Dashboard</h1>
        <p className="text-muted-foreground mt-2">
          Visao geral das suas financas
        </p>
      </div>

      {/* Metric Cards */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        {metrics.map((metric, index) => {
          const Icon = metric.icon;
          return (
            <Card key={index} className="hover:shadow-md transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">
                  {metric.title}
                </CardTitle>
                <Icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{metric.value}</div>
                <div className="flex items-center text-xs text-muted-foreground mt-1">
                  {metric.trend === "up" ? (
                    <TrendingUp className="h-3 w-3 mr-1 text-green-500" />
                  ) : (
                    <TrendingDown className="h-3 w-3 mr-1 text-red-500" />
                  )}
                  <span
                    className={
                      metric.trend === "up" ? "text-green-500" : "text-red-500"
                    }
                  >
                    {metric.change}
                  </span>
                  <span className="ml-1">do mes anterior</span>
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
            Receitas e despesas dos ultimos 6 meses
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ResponsiveContainer width="100%" height={350}>
            <AreaChart data={chartData}>
              <defs>
                <linearGradient id="colorIncome" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="5%"
                    stopColor="hsl(var(--primary))"
                    stopOpacity={0.3}
                  />
                  <stop
                    offset="95%"
                    stopColor="hsl(var(--primary))"
                    stopOpacity={0}
                  />
                </linearGradient>
                <linearGradient id="colorExpense" x1="0" y1="0" x2="0" y2="1">
                  <stop
                    offset="5%"
                    stopColor="hsl(var(--destructive))"
                    stopOpacity={0.3}
                  />
                  <stop
                    offset="95%"
                    stopColor="hsl(var(--destructive))"
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
                stroke="hsl(var(--primary))"
                fillOpacity={1}
                fill="url(#colorIncome)"
                strokeWidth={2}
                name="Receitas"
              />
              <Area
                type="monotone"
                dataKey="expense"
                stroke="hsl(var(--destructive))"
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
          <TabsTrigger value="overview">Visao Geral</TabsTrigger>
          <TabsTrigger value="performance">Desempenho</TabsTrigger>
          <TabsTrigger value="categories">Categorias</TabsTrigger>
          <TabsTrigger value="reports">Relatorios</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Transações Recentes</CardTitle>
              <CardDescription>
                Ultimas movimentacoes financeiras
              </CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-12">
                      <Checkbox />
                    </TableHead>
                    <TableHead>Descricao</TableHead>
                    <TableHead>Tipo</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Valor</TableHead>
                    <TableHead>Limite</TableHead>
                    <TableHead>Responsavel</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {tableData.map((row) => (
                    <TableRow key={row.id} className="hover:bg-muted/50">
                      <TableCell>
                        <Checkbox />
                      </TableCell>
                      <TableCell className="font-medium">
                        {row.header}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={
                            row.section === "Receita" ? "default" : "secondary"
                          }
                        >
                          {row.section}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={
                            row.status === "Pago" || row.status === "Recebido"
                              ? "outline"
                              : "secondary"
                          }
                        >
                          {row.status}
                        </Badge>
                      </TableCell>
                      <TableCell>{row.target}</TableCell>
                      <TableCell className="text-muted-foreground">
                        {row.limit}
                      </TableCell>
                      <TableCell className="text-muted-foreground">
                        {row.reviewer}
                      </TableCell>
                    </TableRow>
                  ))}
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
                Analise do desempenho financeiro mensal
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground py-8 text-center">
                Dados de desempenho serao exibidos aqui
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="categories" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Categorias Principais</CardTitle>
              <CardDescription>
                Distribuicao de gastos por categoria
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground py-8 text-center">
                Dados de categorias serao exibidos aqui
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="reports" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Relatorios</CardTitle>
              <CardDescription>
                Relatorios financeiros detalhados
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="text-sm text-muted-foreground py-8 text-center">
                Relatorios serao exibidos aqui
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
