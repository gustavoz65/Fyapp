"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { Switch } from "@/components/ui/switch";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { api } from "@/lib/api";
import { formatCurrency, formatDate } from "@/lib/format";
import type { BankAccount, Category, RecurringTransaction } from "@/types";
import { CalendarClock, Loader2, Plus, Power, PowerOff, Trash2 } from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

function getApiErrorMessage(error: unknown, fallback: string): string {
  if (error && typeof error === "object" && "message" in error) {
    return (error as { message: string }).message;
  }
  return fallback;
}

const FREQUENCY_LABELS: Record<string, string> = {
  daily: "Diária",
  weekly: "Semanal",
  biweekly: "Quinzenal",
  monthly: "Mensal",
  quarterly: "Trimestral",
  yearly: "Anual",
};

export default function RecurringPage() {
  const [recurrings, setRecurrings] = useState<RecurringTransaction[]>([]);
  const [accounts, setAccounts] = useState<BankAccount[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [togglingId, setTogglingId] = useState<string | null>(null);
  const pendingDeleteRef = useRef<{ id: string; timer: ReturnType<typeof setTimeout> } | null>(null);

  const initialFormState = {
    bank_account_id: "",
    category_id: "",
    type: "expense",
    amount: "",
    description: "",
    frequency: "monthly",
    day_of_month: "1",
    auto_confirm: false,
  };

  const [form, setForm] = useState(initialFormState);

  const fetchData = useCallback(async () => {
    try {
      const [recs, accs, cats] = await Promise.all([
        api.get<RecurringTransaction[]>("/recurring-transactions"),
        api.get<BankAccount[]>("/accounts"),
        api.get<Category[]>("/categories"),
      ]);
      setRecurrings(recs || []);
      setAccounts(accs || []);
      setCategories(cats || []);
    } catch {
      toast.error("Erro ao carregar dados. Verifique sua conexão.");
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!form.bank_account_id || !form.description || !form.amount) {
      toast.error("Preencha todos os campos obrigatórios: descrição, conta e valor");
      return;
    }

    setIsSubmitting(true);
    try {
      await api.post("/recurring-transactions", {
        bank_account_id: form.bank_account_id,
        category_id: form.category_id || null,
        type: form.type,
        amount: Number.parseFloat(form.amount),
        description: form.description,
        frequency: form.frequency,
        day_of_month: parseInt(form.day_of_month),
        auto_confirm: form.auto_confirm,
      });

      toast.success("Recorrência criada com sucesso");
      setDialogOpen(false);
      setForm(initialFormState);
      fetchData();
    } catch (error) {
      toast.error(getApiErrorMessage(error, "Erro ao criar recorrência. Tente novamente."));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleToggle = async (id: string, currentStatus: boolean) => {
    setTogglingId(id);
    try {
      await api.patch(`/recurring-transactions/${id}/toggle`, {
        is_active: !currentStatus,
      });
      toast.success(currentStatus ? "Recorrência pausada" : "Recorrência ativada");
      fetchData();
    } catch (error) {
      toast.error(getApiErrorMessage(error, "Erro ao atualizar recorrência."));
    } finally {
      setTogglingId(null);
    }
  };

  const handleDelete = (id: string) => {
    const recToDelete = recurrings.find((r) => r.id === id);
    if (!recToDelete) return;

    setRecurrings((prev) => prev.filter((r) => r.id !== id));

    if (pendingDeleteRef.current) {
      clearTimeout(pendingDeleteRef.current.timer);
    }

    const timer = setTimeout(async () => {
      try {
        await api.delete(`/recurring-transactions/${id}`);
        pendingDeleteRef.current = null;
        fetchData();
      } catch {
        setRecurrings((prev) => [...prev, recToDelete]);
        toast.error("Erro ao excluir recorrência. Tente novamente.");
      }
    }, 5000);

    pendingDeleteRef.current = { id, timer };

    toast("Recorrência excluída", {
      description: recToDelete.description,
      action: {
        label: "Desfazer",
        onClick: () => {
          if (pendingDeleteRef.current?.id === id) {
            clearTimeout(pendingDeleteRef.current.timer);
            pendingDeleteRef.current = null;
            setRecurrings((prev) => [...prev, recToDelete]);
            toast.success("Ação desfeita com sucesso");
          }
        },
      },
      duration: 5000,
    });
  };

  if (isLoading) {
    return (
      <div className="space-y-6 pb-8">
        <div className="flex items-center justify-between">
          <div className="space-y-2">
            <Skeleton className="h-9 w-64" />
            <Skeleton className="h-4 w-32" />
          </div>
          <Skeleton className="h-10 w-40" />
        </div>
        <Card>
          <CardContent className="p-0">
            <div className="space-y-0">
              {["a", "b", "c", "d"].map((k) => (
                <div key={k} className="flex items-center gap-4 p-4 border-b last:border-0">
                  <Skeleton className="h-4 w-40" />
                  <Skeleton className="h-4 w-24" />
                  <Skeleton className="h-4 w-20" />
                  <Skeleton className="h-4 w-24 ml-auto" />
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  const activeCount = recurrings.filter((r) => r.is_active).length;

  return (
    <div className="space-y-8 pb-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Transações Recorrentes</h1>
          <p className="text-muted-foreground mt-1">
            {activeCount} ativa{activeCount !== 1 ? "s" : ""}
            {recurrings.length > activeCount && ` · ${recurrings.length - activeCount} pausada${recurrings.length - activeCount !== 1 ? "s" : ""}`}
          </p>
        </div>

        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button size="lg">
              <Plus className="w-4 h-4 mr-2" />
              Nova Recorrência
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>Nova Recorrência</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <Label>Descrição *</Label>
                <Input
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  placeholder="Ex: Aluguel, Netflix, Academia..."
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Tipo</Label>
                  <Select value={form.type} onValueChange={(v) => setForm({ ...form, type: v })}>
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="expense">Despesa</SelectItem>
                      <SelectItem value="income">Receita</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label>Valor (R$) *</Label>
                  <Input
                    type="number"
                    step="0.01"
                    value={form.amount}
                    onChange={(e) => setForm({ ...form, amount: e.target.value })}
                    placeholder="0,00"
                  />
                </div>
              </div>

              <div>
                <Label>Conta *</Label>
                <Select
                  value={form.bank_account_id}
                  onValueChange={(v) => setForm({ ...form, bank_account_id: v })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione a conta" />
                  </SelectTrigger>
                  <SelectContent>
                    {accounts.map((acc) => (
                      <SelectItem key={acc.id} value={acc.id}>
                        {acc.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label>Categoria</Label>
                <Select
                  value={form.category_id}
                  onValueChange={(v) => setForm({ ...form, category_id: v })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione (opcional)" />
                  </SelectTrigger>
                  <SelectContent>
                    {categories.map((cat) => (
                      <SelectItem key={cat.id} value={cat.id}>
                        {cat.icon} {cat.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Frequência</Label>
                  <Select
                    value={form.frequency}
                    onValueChange={(v) => setForm({ ...form, frequency: v })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {Object.entries(FREQUENCY_LABELS).map(([value, label]) => (
                        <SelectItem key={value} value={value}>{label}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label>Dia do mês</Label>
                  <Input
                    type="number"
                    min="1"
                    max="31"
                    value={form.day_of_month}
                    onChange={(e) => setForm({ ...form, day_of_month: e.target.value })}
                  />
                </div>
              </div>

              <div className="flex items-center justify-between rounded-lg border p-3">
                <div>
                  <p className="text-sm font-medium">Confirmar automaticamente</p>
                  <p className="text-xs text-muted-foreground">Lançar transação sem necessidade de confirmação manual</p>
                </div>
                <Switch
                  checked={form.auto_confirm}
                  onCheckedChange={(checked) => setForm({ ...form, auto_confirm: checked })}
                />
              </div>

              <div className="flex gap-2 pt-2">
                <Button type="button" variant="outline" className="flex-1" onClick={() => setDialogOpen(false)}>
                  Cancelar
                </Button>
                <Button type="submit" className="flex-1" disabled={isSubmitting}>
                  {isSubmitting && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
                  Criar recorrência
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {recurrings.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="py-16 flex flex-col items-center gap-4 text-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <CalendarClock className="h-8 w-8 text-muted-foreground" />
            </div>
            <div>
              <h3 className="font-semibold text-lg">Nenhuma recorrência cadastrada</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Automatize lançamentos fixos como aluguel, salário, assinaturas e contas mensais
              </p>
            </div>
            <Button onClick={() => setDialogOpen(true)}>
              <Plus className="h-4 w-4 mr-2" />
              Criar primeira recorrência
            </Button>
          </CardContent>
        </Card>
      ) : (
        <Card>
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Descrição</TableHead>
                  <TableHead>Valor</TableHead>
                  <TableHead>Frequência</TableHead>
                  <TableHead>Próxima</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-25"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {recurrings.map((rec) => (
                  <TableRow key={rec.id}>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {rec.category?.icon && (
                          <span className="text-lg">{rec.category.icon}</span>
                        )}
                        <div>
                          <span className="font-medium">{rec.description}</span>
                          {rec.bank_account && (
                            <p className="text-xs text-muted-foreground">{rec.bank_account.name}</p>
                          )}
                        </div>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className={rec.type === "income" ? "text-green-600 font-medium" : "text-red-600 font-medium"}>
                        {rec.type === "income" ? "+" : "-"}
                        {formatCurrency(rec.amount)}
                      </span>
                    </TableCell>
                    <TableCell>{FREQUENCY_LABELS[rec.frequency] || rec.frequency}</TableCell>
                    <TableCell>{formatDate(rec.next_occurrence)}</TableCell>
                    <TableCell>
                      <Badge variant={rec.is_active ? "default" : "secondary"}>
                        {rec.is_active ? "Ativa" : "Pausada"}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Button
                          size="icon"
                          variant="ghost"
                          disabled={togglingId === rec.id}
                          onClick={() => handleToggle(rec.id, rec.is_active)}
                          title={rec.is_active ? "Pausar" : "Ativar"}
                        >
                          {togglingId === rec.id ? (
                            <Loader2 className="w-4 h-4 animate-spin" />
                          ) : rec.is_active ? (
                            <PowerOff className="w-4 h-4" />
                          ) : (
                            <Power className="w-4 h-4" />
                          )}
                        </Button>
                        <Button
                          size="icon"
                          variant="ghost"
                          onClick={() => handleDelete(rec.id)}
                          title="Excluir"
                        >
                          <Trash2 className="w-4 h-4 text-destructive" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
