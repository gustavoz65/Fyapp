"use client";

import { useCallback, useEffect, useState } from "react";
import { z } from "zod";
import { Plus, Pencil, Trash2, TrendingUp } from "lucide-react";
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
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

export default function GoalsPage() {
  const [goals, setGoals] = useState<Goal[]>([]);
  const [summary, setSummary] = useState<GoalSummary | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [contribDialogOpen, setContribDialogOpen] = useState(false);
  const [editing, setEditing] = useState<Goal | null>(null);
  const [selectedGoal, setSelectedGoal] = useState<Goal | null>(null);
  const [contributions, setContributions] = useState<GoalContribution[]>([]);

  const [form, setForm] = useState({ name: "", description: "", target_amount: "", target_date: "", color: "#3b82f6", priority: "3" });
  const [contribForm, setContribForm] = useState({ amount: "", note: "" });

  const fetchData = useCallback(async () => {
    try {
      const [g, s] = await Promise.all([
        api.get<Goal[]>("/goals"),
        api.get<GoalSummary>("/goals/summary"),
      ]);
      setGoals(g || []);
      setSummary(s || null);
    } catch {} finally { setIsLoading(false); }
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
    } catch { setContributions([]); }
    setContribDialogOpen(true);
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      if (editing) {
        const body: UpdateGoalRequest = { name: form.name, description: form.description || undefined, target_amount: form.target_amount, color: form.color, priority: parseInt(form.priority) };
        if (form.target_date) body.target_date = new Date(form.target_date).toISOString();
        await api.put(`/goals/${editing.id}`, body);
        toast.success("Meta atualizada");
      } else {
        // Validar com Zod
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
        toast.success("Meta criada");
      }
      setDialogOpen(false);
      fetchData();
    } catch (error) {
      if (error instanceof z.ZodError) {
        toast.error(error.issues[0].message);
      } else {
        toast.error("Erro ao salvar meta");
      }
    }
  }

  async function handleContribution(e: React.FormEvent) {
    e.preventDefault();
    if (!selectedGoal) return;
    try {
      const body: CreateGoalContributionRequest = { amount: contribForm.amount, note: contribForm.note || undefined };
      await api.post(`/goals/${selectedGoal.id}/contributions`, body);
      toast.success("Contribuicao adicionada");
      setContribDialogOpen(false);
      fetchData();
    } catch { toast.error("Erro ao adicionar contribuicao"); }
  }

  async function handleDelete(id: string) {
    try {
      await api.delete(`/goals/${id}`);
      toast.success("Meta removida");
      fetchData();
    } catch { toast.error("Erro ao remover"); }
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold">Metas</h1>
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => <Skeleton key={i} className="h-56" />)}
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
            <Button size="lg" onClick={openCreate}><Plus className="h-4 w-4 mr-2" />Nova Meta</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader><DialogTitle>{editing ? "Editar Meta" : "Nova Meta"}</DialogTitle></DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2"><Label>Nome</Label><Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required /></div>
              <div className="space-y-2"><Label>Descricao</Label><Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></div>
              <div className="space-y-2"><Label>Valor Alvo</Label><Input type="number" step="0.01" min="0.01" value={form.target_amount} onChange={(e) => setForm({ ...form, target_amount: e.target.value })} required /></div>
              <div className="space-y-2"><Label>Data Alvo</Label><Input type="date" value={form.target_date} onChange={(e) => setForm({ ...form, target_date: e.target.value })} /></div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2"><Label>Cor</Label><Input type="color" value={form.color} onChange={(e) => setForm({ ...form, color: e.target.value })} /></div>
                <div className="space-y-2"><Label>Prioridade (1-5)</Label><Input type="number" min="1" max="5" value={form.priority} onChange={(e) => setForm({ ...form, priority: e.target.value })} /></div>
              </div>
              <Button type="submit" className="w-full">{editing ? "Salvar" : "Criar"}</Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {summary && (
        <div className="grid gap-6 md:grid-cols-3">
          <Card className="hover:shadow-md transition-shadow"><CardHeader className="pb-3"><CardTitle className="text-sm font-medium text-muted-foreground">Total em Metas</CardTitle></CardHeader><CardContent><p className="text-3xl font-bold tracking-tight">{formatCurrency(summary.total_target_amount)}</p></CardContent></Card>
          <Card className="hover:shadow-md transition-shadow"><CardHeader className="pb-3"><CardTitle className="text-sm font-medium text-muted-foreground">Total Economizado</CardTitle></CardHeader><CardContent><p className="text-3xl font-bold tracking-tight text-green-600 dark:text-green-500">{formatCurrency(summary.total_saved_amount)}</p></CardContent></Card>
          <Card className="hover:shadow-md transition-shadow"><CardHeader className="pb-3"><CardTitle className="text-sm font-medium text-muted-foreground">Progresso Geral</CardTitle></CardHeader><CardContent><p className="text-3xl font-bold tracking-tight">{formatPercentage(summary.overall_progress)}</p></CardContent></Card>
        </div>
      )}

      <Dialog open={contribDialogOpen} onOpenChange={setContribDialogOpen}>
        <DialogContent>
          <DialogHeader><DialogTitle>Contribuir para &quot;{selectedGoal?.name}&quot;</DialogTitle></DialogHeader>
          <form onSubmit={handleContribution} className="space-y-4">
            <div className="space-y-2"><Label>Valor</Label><Input type="number" step="0.01" min="0.01" value={contribForm.amount} onChange={(e) => setContribForm({ ...contribForm, amount: e.target.value })} required /></div>
            <div className="space-y-2"><Label>Nota</Label><Input value={contribForm.note} onChange={(e) => setContribForm({ ...contribForm, note: e.target.value })} /></div>
            <Button type="submit" className="w-full">Adicionar Contribuicao</Button>
          </form>
          {contributions.length > 0 && (
            <div className="mt-4 space-y-2">
              <h4 className="text-sm font-medium">Historico</h4>
              {contributions.map((c) => (
                <div key={c.id} className="flex justify-between text-sm">
                  <span>{formatDate(c.contribution_date)}{c.note && ` - ${c.note}`}</span>
                  <span className="font-medium text-green-500">+{formatCurrency(c.amount)}</span>
                </div>
              ))}
            </div>
          )}
        </DialogContent>
      </Dialog>

      {goals.length === 0 ? (
        <Card><CardContent className="py-8 text-center text-muted-foreground">Nenhuma meta criada</CardContent></Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {goals.map((goal) => {
            const current = parseFloat(goal.current_amount);
            const target = parseFloat(goal.target_amount);
            const percentage = target > 0 ? (current / target) * 100 : 0;

            return (
              <Card key={goal.id}>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-sm font-medium">{goal.name}</CardTitle>
                    <Badge variant={goal.status === "completed" ? "default" : goal.status === "cancelled" ? "destructive" : "secondary"}>
                      {getGoalStatusLabel(goal.status)}
                    </Badge>
                  </div>
                  {goal.description && <p className="text-xs text-muted-foreground">{goal.description}</p>}
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex justify-between text-sm">
                    <span>{formatCurrency(goal.current_amount)}</span>
                    <span className="text-muted-foreground">{formatCurrency(goal.target_amount)}</span>
                  </div>
                  <Progress value={Math.min(percentage, 100)} className="h-2" />
                  <p className="text-xs text-muted-foreground">{formatPercentage(percentage)} concluido</p>
                  {goal.target_date && <p className="text-xs text-muted-foreground">Prazo: {formatDate(goal.target_date)}</p>}
                  <div className="flex gap-2">
                    <Button variant="ghost" size="sm" onClick={() => openContributions(goal)}>
                      <TrendingUp className="h-3 w-3 mr-1" />Contribuir
                    </Button>
                    <Button variant="ghost" size="sm" onClick={() => openEdit(goal)}>
                      <Pencil className="h-3 w-3 mr-1" />Editar
                    </Button>
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button variant="ghost" size="sm" className="text-destructive"><Trash2 className="h-3 w-3" /></Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader><AlertDialogTitle>Remover meta?</AlertDialogTitle><AlertDialogDescription>Essa acao nao pode ser desfeita.</AlertDialogDescription></AlertDialogHeader>
                        <AlertDialogFooter><AlertDialogCancel>Cancelar</AlertDialogCancel><AlertDialogAction onClick={() => handleDelete(goal.id)}>Remover</AlertDialogAction></AlertDialogFooter>
                      </AlertDialogContent>
                    </AlertDialog>
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
