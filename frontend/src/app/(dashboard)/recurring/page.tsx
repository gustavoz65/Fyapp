"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
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
import { Plus, Power, PowerOff, Trash2 } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

export default function RecurringPage() {
  const [recurrings, setRecurrings] = useState<RecurringTransaction[]>([]);
  const [accounts, setAccounts] = useState<BankAccount[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);

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
      toast.error("Erro ao carregar dados");
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!form.bank_account_id || !form.description || !form.amount) {
      toast.error("Preencha todos os campos obrigatórios");
      return;
    }

    try {
      await api.post("/recurring-transactions", {
        bank_account_id: form.bank_account_id,
        category_id: form.category_id || null,
        type: form.type,
        amount: parseFloat(form.amount),
        description: form.description,
        frequency: form.frequency,
        day_of_month: parseInt(form.day_of_month),
        auto_confirm: form.auto_confirm,
      });

      toast.success("Recorrência criada");
      setDialogOpen(false);
      setForm(initialFormState);
      fetchData();
    } catch {
      toast.error("Erro ao criar recorrência");
    }
  };

  const handleToggle = async (id: string, currentStatus: boolean) => {
    try {
      await api.patch(`/recurring-transactions/${id}/toggle`, {
        is_active: !currentStatus,
      });
      toast.success(currentStatus ? "Pausada" : "Ativada");
      fetchData();
    } catch {
      toast.error("Erro ao atualizar");
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Excluir esta recorrência?")) return;

    try {
      await api.delete(`/recurring-transactions/${id}`);
      toast.success("Excluída");
      fetchData();
    } catch {
      toast.error("Erro ao excluir");
    }
  };

  const getFrequencyLabel = (freq: string) => {
    const labels: Record<string, string> = {
      daily: "Diária",
      weekly: "Semanal",
      biweekly: "Quinzenal",
      monthly: "Mensal",
      quarterly: "Trimestral",
      yearly: "Anual",
    };
    return labels[freq] || freq;
  };

  if (isLoading) {
    return <div className="p-8">Carregando...</div>;
  }

  return (
    <div className="p-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Transações Recorrentes</h1>
          <p className="text-muted-foreground mt-1">
            {recurrings.filter((r) => r.is_active).length} ativas
          </p>
        </div>

        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button>
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
                <Label>Descrição</Label>
                <Input
                  value={form.description}
                  onChange={(e) =>
                    setForm({ ...form, description: e.target.value })
                  }
                  placeholder="Ex: Aluguel"
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label>Tipo</Label>
                  <Select
                    value={form.type}
                    onValueChange={(v) => setForm({ ...form, type: v })}
                  >
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
                  <Label>Valor</Label>
                  <Input
                    type="number"
                    step="0.01"
                    value={form.amount}
                    onChange={(e) =>
                      setForm({ ...form, amount: e.target.value })
                    }
                    placeholder="0.00"
                  />
                </div>
              </div>

              <div>
                <Label>Conta</Label>
                <Select
                  value={form.bank_account_id}
                  onValueChange={(v) =>
                    setForm({ ...form, bank_account_id: v })
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione" />
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
                    <SelectValue placeholder="Selecione" />
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
                      <SelectItem value="daily">Diária</SelectItem>
                      <SelectItem value="weekly">Semanal</SelectItem>
                      <SelectItem value="biweekly">Quinzenal</SelectItem>
                      <SelectItem value="monthly">Mensal</SelectItem>
                      <SelectItem value="quarterly">Trimestral</SelectItem>
                      <SelectItem value="yearly">Anual</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label>Dia do Mês</Label>
                  <Input
                    type="number"
                    min="1"
                    max="31"
                    value={form.day_of_month}
                    onChange={(e) =>
                      setForm({ ...form, day_of_month: e.target.value })
                    }
                  />
                </div>
              </div>

              <div className="flex items-center justify-between">
                <Label>Confirmar automaticamente</Label>
                <Switch
                  checked={form.auto_confirm}
                  onCheckedChange={(checked) =>
                    setForm({ ...form, auto_confirm: checked })
                  }
                />
              </div>

              <div className="flex gap-2 pt-4">
                <Button
                  type="button"
                  variant="outline"
                  className="flex-1"
                  onClick={() => setDialogOpen(false)}
                >
                  Cancelar
                </Button>
                <Button type="submit" className="flex-1">
                  Criar
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

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
                <TableHead className="w-[100px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {recurrings.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center text-muted-foreground py-8">
                    Nenhuma recorrência cadastrada
                  </TableCell>
                </TableRow>
              ) : (
                recurrings.map((rec) => (
                  <TableRow key={rec.id}>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {rec.category?.icon && (
                          <span className="text-lg">{rec.category.icon}</span>
                        )}
                        <span className="font-medium">{rec.description}</span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span
                        className={
                          rec.type === "income"
                            ? "text-green-600"
                            : "text-red-600"
                        }
                      >
                        {rec.type === "income" ? "+" : "-"}
                        {formatCurrency(rec.amount)}
                      </span>
                    </TableCell>
                    <TableCell>{getFrequencyLabel(rec.frequency)}</TableCell>
                    <TableCell>{formatDate(rec.next_occurrence)}</TableCell>
                    <TableCell>
                      <Badge variant={rec.is_active ? "default" : "secondary"}>
                        {rec.is_active ? "Ativa" : "Pausada"}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <Button
                          size="icon"
                          variant="ghost"
                          onClick={() => handleToggle(rec.id, rec.is_active)}
                        >
                          {rec.is_active ? (
                            <PowerOff className="w-4 h-4" />
                          ) : (
                            <Power className="w-4 h-4" />
                          )}
                        </Button>
                        <Button
                          size="icon"
                          variant="ghost"
                          onClick={() => handleDelete(rec.id)}
                        >
                          <Trash2 className="w-4 h-4 text-red-600" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
