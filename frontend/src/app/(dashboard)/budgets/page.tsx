"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { z } from "zod";
import { Loader2, Pencil, PiggyBank, Plus, Trash2 } from "lucide-react";
import { api } from "@/lib/api";
import type { Budget, Category, CreateBudgetRequest, UpdateBudgetRequest } from "@/types";
import { formatCurrency, formatPercentage, getBudgetPeriodLabel } from "@/lib/format";
import { createBudgetSchema } from "@/lib/schemas";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

function getApiErrorMessage(error: unknown, fallback: string): string {
  if (error && typeof error === "object" && "message" in error) {
    return (error as { message: string }).message;
  }
  return fallback;
}

export default function BudgetsPage() {
  const [budgets, setBudgets] = useState<Budget[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<Budget | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const pendingDeleteRef = useRef<{ id: string; timer: ReturnType<typeof setTimeout> } | null>(null);

  const [form, setForm] = useState({
    name: "", category_id: "", amount: "", period_type: "monthly",
    start_date: new Date().toISOString().split("T")[0],
    end_date: "", alert_threshold: "80",
  });

  const fetchData = useCallback(async () => {
    try {
      const [b, c] = await Promise.all([
        api.get<Budget[]>("/budgets"),
        api.get<Category[]>("/categories"),
      ]);
      setBudgets(b || []);
      setCategories(c || []);
    } catch {
      // empty
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => { fetchData(); }, [fetchData]);

  function openCreate() {
    setEditing(null);
    const now = new Date();
    const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    setForm({
      name: "", category_id: "", amount: "", period_type: "monthly",
      start_date: now.toISOString().split("T")[0],
      end_date: endOfMonth.toISOString().split("T")[0],
      alert_threshold: "80",
    });
    setDialogOpen(true);
  }

  function openEdit(budget: Budget) {
    setEditing(budget);
    setForm({
      name: budget.name, category_id: budget.category_id || "",
      amount: budget.amount, period_type: budget.period_type,
      start_date: budget.start_date.split("T")[0],
      end_date: budget.end_date.split("T")[0],
      alert_threshold: budget.alert_threshold,
    });
    setDialogOpen(true);
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      if (editing) {
        const body: UpdateBudgetRequest = {
          name: form.name,
          amount: form.amount,
          alert_threshold: form.alert_threshold,
        };
        await api.put(`/budgets/${editing.id}`, body);
        toast.success("Orçamento atualizado com sucesso");
      } else {
        const validated = createBudgetSchema.parse({
          name: form.name,
          category_id: form.category_id || "",
          amount: form.amount,
          period_type: form.period_type,
          start_date: form.start_date,
          end_date: form.end_date,
          alert_threshold: form.alert_threshold || "",
        });

        const body: CreateBudgetRequest = {
          name: validated.name,
          category_id: validated.category_id || undefined,
          amount: validated.amount,
          period_type: validated.period_type as CreateBudgetRequest["period_type"],
          start_date: new Date(validated.start_date).toISOString(),
          end_date: new Date(validated.end_date).toISOString(),
          alert_threshold: validated.alert_threshold,
        };
        await api.post("/budgets", body);
        toast.success("Orçamento criado com sucesso");
      }
      setDialogOpen(false);
      fetchData();
    } catch (error) {
      if (error instanceof z.ZodError) {
        toast.error(error.issues[0].message);
      } else {
        toast.error(getApiErrorMessage(error, "Erro ao salvar orçamento. Tente novamente."));
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  function handleDelete(id: string) {
    const budgetToDelete = budgets.find((b) => b.id === id);
    if (!budgetToDelete) return;

    setBudgets((prev) => prev.filter((b) => b.id !== id));

    if (pendingDeleteRef.current) {
      clearTimeout(pendingDeleteRef.current.timer);
    }

    const timer = setTimeout(async () => {
      try {
        await api.delete(`/budgets/${id}`);
        pendingDeleteRef.current = null;
        fetchData();
      } catch {
        setBudgets((prev) => [...prev, budgetToDelete]);
        toast.error("Erro ao remover orçamento. Tente novamente.");
      }
    }, 5000);

    pendingDeleteRef.current = { id, timer };

    toast("Orçamento removido", {
      description: budgetToDelete.name,
      action: {
        label: "Desfazer",
        onClick: () => {
          if (pendingDeleteRef.current?.id === id) {
            clearTimeout(pendingDeleteRef.current.timer);
            pendingDeleteRef.current = null;
            setBudgets((prev) => [...prev, budgetToDelete]);
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
            <Skeleton className="h-9 w-40" />
            <Skeleton className="h-4 w-56" />
          </div>
          <Skeleton className="h-10 w-36" />
        </div>
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {["a", "b", "c"].map((i) => <Skeleton key={i} className="h-48" />)}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8 pb-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Orçamentos</h1>
          <p className="text-muted-foreground mt-2">
            Controle seus gastos por categoria
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button size="lg" onClick={openCreate}>
              <Plus className="h-4 w-4 mr-2" />Novo Orçamento
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editing ? "Editar Orçamento" : "Novo Orçamento"}</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label>Nome</Label>
                <Input
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Ex: Alimentação do mês"
                  required
                />
              </div>
              {!editing && (
                <>
                  <div className="space-y-2">
                    <Label>Categoria</Label>
                    <Select value={form.category_id} onValueChange={(v) => setForm({ ...form, category_id: v })}>
                      <SelectTrigger><SelectValue placeholder="Selecione (opcional)" /></SelectTrigger>
                      <SelectContent>
                        {categories.filter((c) => c.type === "expense").map((c) => (
                          <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label>Período</Label>
                    <Select value={form.period_type} onValueChange={(v) => setForm({ ...form, period_type: v })}>
                      <SelectTrigger><SelectValue /></SelectTrigger>
                      <SelectContent>
                        <SelectItem value="monthly">Mensal</SelectItem>
                        <SelectItem value="quarterly">Trimestral</SelectItem>
                        <SelectItem value="yearly">Anual</SelectItem>
                        <SelectItem value="custom">Personalizado</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <Label>Data de início</Label>
                      <Input type="date" value={form.start_date} onChange={(e) => setForm({ ...form, start_date: e.target.value })} required />
                    </div>
                    <div className="space-y-2">
                      <Label>Data de fim</Label>
                      <Input type="date" value={form.end_date} onChange={(e) => setForm({ ...form, end_date: e.target.value })} required />
                    </div>
                  </div>
                </>
              )}
              <div className="space-y-2">
                <Label>Valor limite (R$)</Label>
                <Input
                  type="number"
                  step="0.01"
                  min="0.01"
                  value={form.amount}
                  onChange={(e) => setForm({ ...form, amount: e.target.value })}
                  placeholder="0,00"
                  required
                />
              </div>
              <div className="space-y-2">
                <Label>Alertar quando atingir (%)</Label>
                <Input
                  type="number"
                  min="1"
                  max="100"
                  value={form.alert_threshold}
                  onChange={(e) => setForm({ ...form, alert_threshold: e.target.value })}
                  placeholder="80"
                />
              </div>
              <div className="flex gap-2 pt-2">
                <Button type="button" variant="outline" className="flex-1" onClick={() => setDialogOpen(false)}>
                  Cancelar
                </Button>
                <Button type="submit" className="flex-1" disabled={isSubmitting}>
                  {isSubmitting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                  {editing ? "Salvar alterações" : "Criar orçamento"}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {budgets.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="py-16 flex flex-col items-center gap-4 text-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <PiggyBank className="h-8 w-8 text-muted-foreground" />
            </div>
            <div>
              <h3 className="font-semibold text-lg">Nenhum orçamento criado</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Defina limites de gastos por categoria e receba alertas quando estiver se aproximando
              </p>
            </div>
            <Button onClick={openCreate}>
              <Plus className="h-4 w-4 mr-2" />
              Criar primeiro orçamento
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {budgets.map((budget) => {
            const spent = Number.parseFloat(budget.spent_amount);
            const total = Number.parseFloat(budget.amount);
            const percentage = total > 0 ? (spent / total) * 100 : 0;
            const isOver = percentage > 100;
            const isNearLimit = percentage >= Number.parseFloat(budget.alert_threshold);

            return (
              <Card key={budget.id} className={isOver ? "border-destructive" : isNearLimit ? "border-yellow-500" : ""}>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-sm font-medium">{budget.name}</CardTitle>
                    <Badge variant={isOver ? "destructive" : isNearLimit ? "secondary" : "default" }>
                      {getBudgetPeriodLabel(budget.period_type)}
                    </Badge>
                  </div>
                  {budget.category?.name && (
                    <p className="text-xs text-muted-foreground">{budget.category.name}</p>
                  )}
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex justify-between text-sm">
                    <span>Gasto: <strong>{formatCurrency(budget.spent_amount)}</strong></span>
                    <span className="text-muted-foreground">Limite: {formatCurrency(budget.amount)}</span>
                  </div>
                  <Progress value={Math.min(percentage, 100)} className={`h-2 ${isOver ? "[&>div]:bg-destructive" : isNearLimit ? "[&>div]:bg-yellow-500" : ""}`} />
                  <p className={`text-xs ${isOver ? "text-destructive font-medium" : isNearLimit ? "text-yellow-600 dark:text-yellow-400" : "text-muted-foreground"}`}>
                    {formatPercentage(percentage)} utilizado
                    {isOver && " — Acima do limite!"}
                    {!isOver && isNearLimit && " — Próximo do limite"}
                  </p>
                  <div className="flex gap-2">
                    <Button variant="ghost" size="sm" onClick={() => openEdit(budget)}>
                      <Pencil className="h-3 w-3 mr-1" />Editar
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="text-destructive hover:text-destructive"
                      onClick={() => handleDelete(budget.id)}
                    >
                      <Trash2 className="h-3 w-3 mr-1" />Remover
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
