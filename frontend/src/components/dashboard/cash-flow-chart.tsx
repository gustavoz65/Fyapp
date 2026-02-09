"use client";

import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ChartContainer, ChartTooltip, ChartTooltipContent, type ChartConfig } from "@/components/ui/chart";
import type { MonthlyIncomeExpense } from "@/types";

const chartConfig = {
  income: { label: "Receita", color: "var(--chart-2)" },
  expense: { label: "Despesa", color: "var(--chart-5)" },
} satisfies ChartConfig;

interface CashFlowChartProps {
  data: MonthlyIncomeExpense[];
}

export function CashFlowChart({ data }: CashFlowChartProps) {
  const chartData = data.map((item) => ({
    month: item.month.slice(0, 3),
    income: parseFloat(item.income),
    expense: parseFloat(item.expense),
  }));

  return (
    <Card>
      <CardHeader>
        <CardTitle>Fluxo de Caixa</CardTitle>
        <CardDescription>Receita vs Despesa por mes</CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="h-[300px] w-full">
          <BarChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" vertical={false} />
            <XAxis dataKey="month" tickLine={false} axisLine={false} />
            <YAxis tickLine={false} axisLine={false} tickFormatter={(v) => `R$${v}`} />
            <ChartTooltip content={<ChartTooltipContent />} />
            <Bar dataKey="income" fill="var(--color-income)" radius={[4, 4, 0, 0]} />
            <Bar dataKey="expense" fill="var(--color-expense)" radius={[4, 4, 0, 0]} />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
