"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { z } from "zod";
import { Loader2, Pencil, Plus, Target, Trash2, TrendingUp } from "lucide-react";
import { api } from "@/lib/api";
import type { Goal, GoalContribution, GoalSummary, CreateGoalRequest, CreateGoalContributionRequest, UpdateGoalRequest } from "@/types";
import { formatCurrency, formatPercentage, formatDate, getGoalStatusLabel } from "@/lib/format";
import { createGoalSchema } from "@/lib/schemas";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

function getApiErrorMessage(error: unknown, fallback: string): string {
  if (error && typeof error === "object" && "message" in error) {
    return (error as { message: string }).message;
  }
  return fallback;
}

export default function GoalsPage() {
  const [goals, setGoals] = useState<Goal[]>([]);
  const [summary, setSummary] = useState<GoalSummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [contribDialogOpen, setContribDialogOpen] = useState(false);
  const [editing, setEditing] = useState<Goal | null>(null);
  const [selectedGoal, setSelectedGoal] = useState<Goal | null>(null);
  const [contributions, setContributions] = useState<GoalContribution[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isContributing, setIsContributing] = useState(false);
  const pendingDeleteRef = useRef<{ id: string; timer: ReturnType<typeof setTimeout> } | null>(null);

  const [form, setForm] = useState({
    name: "", description: "", target_amount: "", target_date: "", color: "#3b82f6", priority: "3",
  });
  const [contribForm, setContribForm] = useState({ amount: "", note: "" });

  const fetchData = useCallback(async () => {
    try {
      const [g, s] = await Promise.all([
        api.get<Goal[]>("/goals"),
        api.get<GoalSummary>("/goals/summary"),
      ]);
      setGoals(g || []);
      setSummary(s || null);
    } catch {
      // empty
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => { fetchData(); }, [fetchData]);

  function openCreate() {
    setEditing(null);
    setForm({ name: "", description: "", target_amount: "", target_date: "", color: "#3b82f6", priority: "3" });
    setDialogOpen(true);
  }

  function openEdit(goal: Goal) {
    setEditing(goal);
    setForm({
      name: goal.name, description: goal.description || "",
      target_amount: goal.target_amount, target_date: goal.target_date?.split("T")[0] || "",
      color: goal.color, priority: String(goal.priority),
    });
    setDialogOpen(true);
  }

  async function openContributions(goal: Goal) {
    setSelectedGoal(goal);
    setContribForm({ amount: "", note: "" });
    try {
      const c = await api.get<GoalContribution[]>(`/goals/${goal.id}/contributions`);
      setContributions(c || []);
    } catch {
      setContributions([]);
    }
    setContribDialogOpen(true);
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      if (editing) {
        const body: UpdateGoalRequest = {
          name: form.name,
          description: form.description || undefined,
          target_amount: form.target_amount,
          color: form.color,
          priority: parseInt(form.priority),
        };
        if (form.target_date) body.target_date = new Date(form.target_date).toISOString();
        await api.put(`/goals/${editing.id}`, body);
        toast.success("Meta atualizada com sucesso");
      } else {
        const validated = createGoalSchema.parse({
          name: form.name,
          description: form.description || "",
          target_amount: form.target_amount,
          target_date: form.target_date || "",
          color: form.color,
          priority: parseInt(form.priority),
        });

        const body: CreateGoalRequest = {
          name: validated.name,
          description: validated.description || undefined,
          target_amount: validated.target_amount,
          color: validated.color,
          priority: validated.priority,
        };
        if (validated.target_date) body.target_date = new Date(validated.target_date).toISOString();
        await api.post("/goals", body);
        toast.success("Meta criada com sucesso");
      }
      setDialogOpen(false);
      fetchData();
    } catch (error) {
      if (error instanceof z.ZodError) {
        toast.error(error.issues[0].message);
      } else {
        toast.error(getApiErrorMessage(error, "Erro ao salvar meta. Tente novamente."));
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handleContribution(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!selectedGoal) return;
    setIsContributing(true);
    try {
      const body: CreateGoalContributionRequest = {
        amount: contribForm.amount,
        note: contribForm.note || undefined,
      };
      await api.post(`/goals/${selectedGoal.id}/contributions`, body);
      toast.success(`${formatCurrency(contribForm.amount)} adicionado à meta "${selectedGoal.name}"`);
      setContribDialogOpen(false);
      fetchData();
    } catch (error) {
      toast.error(getApiErrorMessage(error, "Erro ao adicionar contribuição. Tente novamente."));
    } finally {
      setIsContributing(false);
    }
  }

  function handleDelete(id: string) {
    const goalToDelete = goals.find((g) => g.id === id);
    if (!goalToDelete) return;

    setGoals((prev) => prev.filter((g) => g.id !== id));

    if (pendingDeleteRef.current) {
      clearTimeout(pendingDeleteRef.current.timer);
    }

    const timer = setTimeout(async () => {
      try {
        await api.delete(`/goals/${id}`);
        pendingDeleteRef.current = null;
        fetchData();
      } catch {
        setGoals((prev) => [...prev, goalToDelete]);
        toast.error("Erro ao remover meta. Tente novamente.");
      }
    }, 5000);

    pendingDeleteRef.current = { id, timer };

    toast("Meta removida", {
      description: goalToDelete.name,
      action: {
        label: "Desfazer",
        onClick: () => {
          if (pendingDeleteRef.current?.id === id) {
            clearTimeout(pendingDeleteRef.current.timer);
            pendingDeleteRef.current = null;
            setGoals((prev) => [...prev, goalToDelete]);
            toast.success("Ação desfeita com sucesso");
          }
        },
      },
      duration: 5000,
    });
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="space-y-2">
            <Skeleton className="h-9 w-56" />
            <Skeleton className="h-4 w-64" />
          </div>
          <Skeleton className="h-10 w-32" />
        </div>
        <div className="grid gap-6 md:grid-cols-3">
          <Skeleton className="h-24" />
          <Skeleton className="h-24" />
          <Skeleton className="h-24" />
        </div>
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          <Skeleton className="h-56" />
          <Skeleton className="h-56" />
          <Skeleton className="h-56" />
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8 pb-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Metas Financeiras</h1>
          <p className="text-muted-foreground mt-2">
            Defina e acompanhe suas metas de economia
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button size="lg" onClick={openCreate}>
              <Plus className="h-4 w-4 mr-2" />Nova Meta
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editing ? "Editar Meta" : "Nova Meta"}</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label>Nome da meta</Label>
                <Input
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Ex: Viagem para Europa, Fundo de emergência..."
                  required
                />
              </div>
              <div className="space-y-2">
                <Label>Descrição</Label>
                <Input
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  placeholder="Descreva sua meta (opcional)"
                />
              </div>
              <div className="space-y-2">
                <Label>Valor alvo (R$)</Label>
                <Input
                  type="number"
                  step="0.01"
                  min="0.01"
                  value={form.target_amount}
                  onChange={(e) => setForm({ ...form, target_amount: e.target.value })}
                  placeholder="0,00"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label>Data alvo</Label>
                <Input
                  type="date"
                  value={form.target_date}
                  onChange={(e) => setForm({ ...form, target_date: e.target.value })}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Cor</Label>
                  <Input
                    type="color"
                    value={form.color}
                    onChange={(e) => setForm({ ...form, color: e.target.value })}
                    className="h-10 w-full p-1 cursor-pointer"
                  />
                </div>
                <div className="space-y-2">
                  <Label>Prioridade (1 = alta, 5 = baixa)</Label>
                  <Input
                    type="number"
                    min="1"
                    max="5"
                    value={form.priority}
                    onChange={(e) => setForm({ ...form, priority: e.target.value })}
                  />
                </div>
              </div>
              <div className="flex gap-2 pt-2">
                <Button type="button" variant="outline" className="flex-1" onClick={() => setDialogOpen(false)}>
                  Cancelar
                </Button>
                <Button type="submit" className="flex-1" disabled={isSubmitting}>
                  {isSubmitting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                  {editing ? "Salvar alterações" : "Criar meta"}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {summary && (
        <div className="grid gap-6 md:grid-cols-3">
          <Card className="hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">Total em Metas</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold tracking-tight">{formatCurrency(summary.total_target_amount)}</p>
            </CardContent>
          </Card>
          <Card className="hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">Total Economizado</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold tracking-tight text-green-600 dark:text-green-500">{formatCurrency(summary.total_saved_amount)}</p>
            </CardContent>
          </Card>
          <Card className="hover:shadow-md transition-shadow">
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-muted-foreground">Progresso Geral</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-bold tracking-tight">{formatPercentage(summary.overall_progress)}</p>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Dialog de contribuição */}
      <Dialog open={contribDialogOpen} onOpenChange={setContribDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Contribuir para &quot;{selectedGoal?.name}&quot;</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleContribution} className="space-y-4">
            <div className="space-y-2">
              <Label>Valor a depositar (R$)</Label>
              <Input
                type="number"
                step="0.01"
                min="0.01"
                value={contribForm.amount}
                onChange={(e) => setContribForm({ ...contribForm, amount: e.target.value })}
                placeholder="0,00"
                required
              />
            </div>
            <div className="space-y-2">
              <Label>Observação</Label>
              <Input
                value={contribForm.note}
                onChange={(e) => setContribForm({ ...contribForm, note: e.target.value })}
                placeholder="Ex: Economias de fevereiro (opcional)"
              />
            </div>
            <div className="flex gap-2 pt-2">
              <Button type="button" variant="outline" className="flex-1" onClick={() => setContribDialogOpen(false)}>
                Cancelar
              </Button>
              <Button type="submit" className="flex-1" disabled={isContributing}>
                {isContributing && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                Adicionar contribuição
              </Button>
            </div>
          </form>
          {contributions.length > 0 && (
            <div className="mt-2 space-y-2 border-t pt-4">
              <h4 className="text-sm font-medium">Histórico de contribuições</h4>
              <div className="space-y-1 max-h-40 overflow-y-auto">
                {contributions.map((c) => (
                  <div key={c.id} className="flex justify-between text-sm py-1">
                    <span className="text-muted-foreground">
                      {formatDate(c.contribution_date)}
                      {c.note && ` — ${c.note}`}
                    </span>
                    <span className="font-medium text-green-600">+{formatCurrency(c.amount)}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {goals.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="py-16 flex flex-col items-center gap-4 text-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Target className="h-8 w-8 text-muted-foreground" />
            </div>
            <div>
              <h3 className="font-semibold text-lg">Nenhuma meta criada</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Defina metas financeiras e acompanhe seu progresso para realizá-las
              </p>
            </div>
            <Button onClick={openCreate}>
              <Plus className="h-4 w-4 mr-2" />
              Criar primeira meta
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {goals.map((goal) => {
            const current = Number.parseFloat(goal.current_amount);
            const target = Number.parseFloat(goal.target_amount);
            const percentage = target > 0 ? (current / target) * 100 : 0;

            return (
              <Card key={goal.id} className="hover:shadow-md transition-shadow">
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <div
                        className="h-3 w-3 rounded-full shrink-0"
                        style={{ backgroundColor: goal.color }}
                      />
                      <CardTitle className="text-sm font-medium">{goal.name}</CardTitle>
                    </div>
                    <Badge variant={goal.status === "completed" ? "default" : goal.status === "cancelled" ? "destructive" : "secondary"}>
                      {getGoalStatusLabel(goal.status)}
                    </Badge>
                  </div>
                  {goal.description && (
                    <p className="text-xs text-muted-foreground mt-1">{goal.description}</p>
                  )}
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex justify-between text-sm">
                    <span className="font-medium">{formatCurrency(goal.current_amount)}</span>
                    <span className="text-muted-foreground">de {formatCurrency(goal.target_amount)}</span>
                  </div>
                  <Progress value={Math.min(percentage, 100)} className="h-2" style={{ "--progress-color": goal.color } as React.CSSProperties} />
                  <p className="text-xs text-muted-foreground">
                    {formatPercentage(percentage)} concluído
                    {goal.target_date && ` · Prazo: ${formatDate(goal.target_date)}`}
                  </p>
                  <div className="flex gap-2 flex-wrap">
                    <Button variant="ghost" size="sm" onClick={() => openContributions(goal)}>
                      <TrendingUp className="h-3 w-3 mr-1" />Contribuir
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => openEdit(goal)}>
                      <Pencil className="h-3 w-3 mr-1" />Editar
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="text-destructive hover:text-destructive"
                      onClick={() => handleDelete(goal.id)}
                    >
                      <Trash2 className="h-3 w-3" />
                    </Button>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}
    </div>
  );
}
