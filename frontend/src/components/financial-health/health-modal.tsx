"use client";

import { useEffect, useMemo } from "react";
import { Activity, TrendingUp, TrendingDown, Target, PiggyBank, Repeat, AlertCircle } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Skeleton } from "@/components/ui/skeleton";
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/components/ui/chart";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, ResponsiveContainer } from "recharts";
import { cn } from "@/lib/utils";
import { useHealthModal } from "@/contexts/health-modal-context";
import { useHistoricalHealth } from "@/hooks/use-financial-health";
import type { FinancialHealthSnapshot, HealthScoreBreakdown } from "@/types";

interface HealthModalProps {
  currentScore: FinancialHealthSnapshot | null;
}

export function HealthModal({ currentScore }: HealthModalProps) {
  const { isOpen, closeModal } = useHealthModal();
  const { history, isLoading: isLoadingHistory, fetchHistory } = useHistoricalHealth(false);

  // Parse breakdown if available
  const breakdown: HealthScoreBreakdown | null = useMemo(() => {
    if (currentScore?.breakdown_json) {
      try {
        return JSON.parse(currentScore.breakdown_json);
      } catch {
        return null;
      }
    }
    return null;
  }, [currentScore]);

  // Prepare chart data
  const chartData = useMemo(() => {
    return history.map((snapshot) => ({
      month: new Date(snapshot.period).toLocaleDateString("pt-BR", {
        month: "short",
        year: "2-digit"
      }),
      score: parseFloat(snapshot.score_percentage),
    })).reverse();
  }, [history]);

  useEffect(() => {
    if (isOpen) {
      fetchHistory();
    }
  }, [isOpen, fetchHistory]);

  const scorePercentage = currentScore ? parseFloat(currentScore.score_percentage) : 0;

  const getScoreColor = (score: number) => {
    if (score < 40) return "text-red-600 dark:text-red-500";
    if (score < 70) return "text-yellow-600 dark:text-yellow-500";
    return "text-green-600 dark:text-green-500";
  };

  const getScoreStatus = (score: number) => {
    if (score < 40) return "Precisa de atenção";
    if (score < 70) return "Razoável";
    return "Excelente";
  };

  const getInsights = (score: number, breakdown?: HealthScoreBreakdown) => {
    const insights: string[] = [];

    if (score >= 70) {
      insights.push("Parabéns! Sua saúde financeira está excelente.");
      insights.push("Continue mantendo seus hábitos financeiros saudáveis.");
    } else if (score >= 40) {
      insights.push("Sua saúde financeira está razoável, mas há espaço para melhorias.");

      if (breakdown && currentScore) {
        const economyRate = parseFloat(currentScore.economy_rate);
        const budgetCompliance = parseFloat(currentScore.budget_compliance);

        if (economyRate < 3) {
          insights.push("Tente aumentar sua taxa de economia mensal.");
        }
        if (budgetCompliance < 3) {
          insights.push("Foque em manter seus gastos dentro do orçamento planejado.");
        }
      }
    } else {
      insights.push("Sua saúde financeira precisa de atenção urgente.");
      insights.push("Revise seus gastos e crie um orçamento realista.");
      insights.push("Considere reduzir despesas não essenciais.");
    }

    return insights;
  };

  const components = currentScore
    ? [
        {
          name: "Taxa de Economia",
          score: parseFloat(currentScore.economy_rate),
          weight: 25,
          icon: TrendingUp,
          description: "Capacidade de poupar mensalmente",
        },
        {
          name: "Cumprimento de Orçamento",
          score: parseFloat(currentScore.budget_compliance),
          weight: 25,
          icon: PiggyBank,
          description: "Aderência aos orçamentos definidos",
        },
        {
          name: "Progresso de Metas",
          score: parseFloat(currentScore.goals_progress),
          weight: 20,
          icon: Target,
          description: "Avanço em direção às suas metas",
        },
        {
          name: "Redução de Gastos",
          score: parseFloat(currentScore.spending_reduction),
          weight: 15,
          icon: TrendingDown,
          description: "Controle e redução de despesas",
        },
        {
          name: "Consistência",
          score: parseFloat(currentScore.consistency),
          weight: 15,
          icon: Repeat,
          description: "Estabilidade financeira ao longo do tempo",
        },
      ]
    : [];

  const insights = getInsights(scorePercentage, breakdown || undefined);

  return (
    <Dialog open={isOpen} onOpenChange={closeModal}>
      {currentScore && (
        <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              Saúde Financeira
            </DialogTitle>
            <DialogDescription>
              Acompanhe sua saúde financeira e receba insights personalizados
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-6">
            {/* Current Score Card */}
            <Card>
              <CardHeader>
                <CardTitle>Score Atual</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between">
                  <div>
                    <div className={cn("text-5xl font-bold tabular-nums", getScoreColor(scorePercentage))}>
                      {Math.round(scorePercentage)}%
                    </div>
                    <p className="text-muted-foreground mt-2">
                      {getScoreStatus(scorePercentage)}
                    </p>
                  </div>
                  <div className="relative w-32 h-32">
                    <svg className="w-full h-full -rotate-90" viewBox="0 0 120 120">
                      <circle
                        cx="60"
                        cy="60"
                        r="54"
                        className="stroke-muted fill-none"
                        strokeWidth="8"
                      />
                      <circle
                        cx="60"
                        cy="60"
                        r="54"
                        className={cn("fill-none transition-all", getScoreColor(scorePercentage).replace("text-", "stroke-"))}
                        strokeWidth="8"
                        strokeDasharray={2 * Math.PI * 54}
                        strokeDashoffset={2 * Math.PI * 54 * (1 - scorePercentage / 100)}
                        strokeLinecap="round"
                      />
                    </svg>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Score Breakdown */}
            <Card>
              <CardHeader>
                <CardTitle>Componentes do Score</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {components.map((component) => {
                    const componentPercentage = (component.score / 5) * 100;
                    return (
                      <div key={component.name} className="space-y-2">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            <component.icon className="h-4 w-4 text-muted-foreground" />
                            <div>
                              <p className="text-sm font-medium">{component.name}</p>
                              <p className="text-xs text-muted-foreground">{component.description}</p>
                            </div>
                          </div>
                          <div className="text-right">
                            <p className="text-sm font-bold">
                              {component.score.toFixed(1)}/5.0
                            </p>
                            <p className="text-xs text-muted-foreground">
                              Peso: {component.weight}%
                            </p>
                          </div>
                        </div>
                        <Progress value={componentPercentage} className="h-2" />
                      </div>
                    );
                  })}
                </div>
              </CardContent>
            </Card>

            {/* Historical Chart */}
            <Card>
              <CardHeader>
                <CardTitle>Evolução (Últimos 6 Meses)</CardTitle>
              </CardHeader>
              <CardContent>
                {isLoadingHistory ? (
                  <Skeleton className="h-[300px] w-full" />
                ) : chartData.length > 0 ? (
                  <ChartContainer
                    config={{
                      score: {
                        label: "Score",
                        color: "hsl(var(--primary))",
                      },
                    }}
                    className="h-[300px]"
                  >
                    <ResponsiveContainer width="100%" height="100%">
                      <LineChart data={chartData}>
                        <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
                        <XAxis
                          dataKey="month"
                          className="text-xs"
                          tick={{ fill: "hsl(var(--muted-foreground))" }}
                        />
                        <YAxis
                          domain={[0, 100]}
                          className="text-xs"
                          tick={{ fill: "hsl(var(--muted-foreground))" }}
                        />
                        <ChartTooltip content={<ChartTooltipContent />} />
                        <Line
                          type="monotone"
                          dataKey="score"
                          stroke="hsl(var(--primary))"
                          strokeWidth={2}
                          dot={{ fill: "hsl(var(--primary))", r: 4 }}
                          activeDot={{ r: 6 }}
                        />
                      </LineChart>
                    </ResponsiveContainer>
                  </ChartContainer>
                ) : (
                  <div className="h-[300px] flex items-center justify-center text-muted-foreground">
                    <p>Dados históricos não disponíveis</p>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Insights */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <AlertCircle className="h-5 w-5" />
                  Insights e Recomendações
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-2">
                  {insights.map((insight, index) => (
                    <li key={index} className="flex items-start gap-2">
                      <span className="text-primary mt-1">•</span>
                      <span className="text-sm">{insight}</span>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          </div>
        </DialogContent>
      )}
    </Dialog>
  );
}
