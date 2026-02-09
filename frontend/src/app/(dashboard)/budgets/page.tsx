"use client";

import { useCallback, useEffect, useState } from "react";
import { Plus, Pencil, Trash2 } from "lucide-react";
import { api } from "@/lib/api";
import type { Budget, Category, CreateBudgetRequest, UpdateBudgetRequest } from "@/types";
import { formatCurrency, formatPercentage, getBudgetPeriodLabel } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

export default function BudgetsPage() {
  const [budgets, setBudgets] = useState<Budget[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<Budget | null>(null);

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
    } catch {} finally { setIsLoading(false); }
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

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      if (editing) {
        const body: UpdateBudgetRequest = { name: form.name, amount: form.amount, alert_threshold: form.alert_threshold };
        await api.put(`/budgets/${editing.id}`, body);
        toast.success("Orcamento atualizado");
      } else {
        const body: CreateBudgetRequest = {
          name: form.name, category_id: form.category_id || undefined,
          amount: form.amount, period_type: form.period_type as CreateBudgetRequest["period_type"],
          start_date: new Date(form.start_date).toISOString(),
          end_date: new Date(form.end_date).toISOString(),
          alert_threshold: form.alert_threshold,
        };
        await api.post("/budgets", body);
        toast.success("Orcamento criado");
      }
      setDialogOpen(false);
      fetchData();
    } catch { toast.error("Erro ao salvar orcamento"); }
  }

  async function handleDelete(id: string) {
    try {
      await api.delete(`/budgets/${id}`);
      toast.success("Orcamento removido");
      fetchData();
    } catch { toast.error("Erro ao remover"); }
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold">Orcamentos</h1>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => <Skeleton key={i} className="h-48" />)}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Orcamentos</h1>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />Novo Orcamento</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editing ? "Editar Orcamento" : "Novo Orcamento"}</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label>Nome</Label>
                <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
              </div>
              {!editing && (
                <>
                  <div className="space-y-2">
                    <Label>Categoria</Label>
                    <Select value={form.category_id} onValueChange={(v) => setForm({ ...form, category_id: v })}>
                      <SelectTrigger><SelectValue placeholder="Opcional" /></SelectTrigger>
                      <SelectContent>
                        {categories.filter((c) => c.type === "expense").map((c) => (
                          <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label>Periodo</Label>
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
                      <Label>Inicio</Label>
                      <Input type="date" value={form.start_date} onChange={(e) => setForm({ ...form, start_date: e.target.value })} required />
                    </div>
                    <div className="space-y-2">
                      <Label>Fim</Label>
                      <Input type="date" value={form.end_date} onChange={(e) => setForm({ ...form, end_date: e.target.value })} required />
                    </div>
                  </div>
                </>
              )}
              <div className="space-y-2">
                <Label>Valor Limite</Label>
                <Input type="number" step="0.01" min="0.01" value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} required />
              </div>
              <div className="space-y-2">
                <Label>Alerta em (%)</Label>
                <Input type="number" min="1" max="100" value={form.alert_threshold} onChange={(e) => setForm({ ...form, alert_threshold: e.target.value })} />
              </div>
              <Button type="submit" className="w-full">{editing ? "Salvar" : "Criar"}</Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {budgets.length === 0 ? (
        <Card><CardContent className="py-8 text-center text-muted-foreground">Nenhum orcamento criado</CardContent></Card>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {budgets.map((budget) => {
            const spent = parseFloat(budget.spent_amount);
            const total = parseFloat(budget.amount);
            const percentage = total > 0 ? (spent / total) * 100 : 0;
            const isOver = percentage > 100;

            return (
              <Card key={budget.id}>
                <CardHeader className="pb-2">
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-sm font-medium">{budget.name}</CardTitle>
                    <Badge variant={isOver ? "destructive" : percentage > 80 ? "secondary" : "default"}>
                      {getBudgetPeriodLabel(budget.period_type)}
                    </Badge>
                  </div>
                  {budget.category?.name && (
                    <p className="text-xs text-muted-foreground">{budget.category.name}</p>
                  )}
                </CardHeader>
                <CardContent className="space-y-3">
                  <div className="flex justify-between text-sm">
                    <span>Gasto: {formatCurrency(budget.spent_amount)}</span>
                    <span>Limite: {formatCurrency(budget.amount)}</span>
                  </div>
                  <Progress value={Math.min(percentage, 100)} className={`h-2 ${isOver ? "[&>div]:bg-destructive" : ""}`} />
                  <p className={`text-xs ${isOver ? "text-destructive" : "text-muted-foreground"}`}>
                    {formatPercentage(percentage)} utilizado
                    {isOver && " - Acima do limite!"}
                  </p>
                  <div className="flex gap-2">
                    <Button variant="ghost" size="sm" onClick={() => openEdit(budget)}>
                      <Pencil className="h-3 w-3 mr-1" />Editar
                    </Button>
                    <AlertDialog>
                      <AlertDialogTrigger asChild>
                        <Button variant="ghost" size="sm" className="text-destructive">
                          <Trash2 className="h-3 w-3 mr-1" />Remover
                        </Button>
                      </AlertDialogTrigger>
                      <AlertDialogContent>
                        <AlertDialogHeader>
                          <AlertDialogTitle>Remover orcamento?</AlertDialogTitle>
                          <AlertDialogDescription>Essa acao nao pode ser desfeita.</AlertDialogDescription>
                        </AlertDialogHeader>
                        <AlertDialogFooter>
                          <AlertDialogCancel>Cancelar</AlertDialogCancel>
                          <AlertDialogAction onClick={() => handleDelete(budget.id)}>Remover</AlertDialogAction>
                        </AlertDialogFooter>
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
