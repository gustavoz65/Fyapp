"use client";

import { z } from "zod";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { api } from "@/lib/api";
import {
  formatCurrency,
  formatDate,
  getTransactionSourceLabel,
  getTransactionTypeLabel,
} from "@/lib/format";
import { createTransactionSchema } from "@/lib/schemas";
import type {
  BankAccount,
  Category,
  CreateTransactionRequest,
  PaginatedResponse,
  Transaction,
} from "@/types";
import { ArrowDownLeft, ArrowUpRight, Check, Plus, Sparkles } from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

interface CategorySuggestion {
  category_id: string;
  category_name: string;
  confidence: number;
  match_count: number;
}

export default function TransactionsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [accounts, setAccounts] = useState<BankAccount[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [filterType, setFilterType] = useState<string>("all");
  const [filterPaid, setFilterPaid] = useState<string>("all");
  const [filterSource, setFilterSource] = useState<string>("all");
  const [categorySuggestions, setCategorySuggestions] = useState<CategorySuggestion[]>([]);
  const suggestDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const initialFormState = {
    bank_account_id: "",
    category_id: "",
    type: "expense" as string,
    amount: "",
    description: "",
    transaction_date: new Date().toISOString().split("T")[0],
    due_date: "",
    is_paid: false,
    auto_pay: false,
  };

  const [form, setForm] = useState(initialFormState);

  const fetchData = useCallback(async () => {
    try {
      let endpoint = `/transactions?page=${page}&page_size=20`;
      if (filterType !== "all") endpoint += `&type=${filterType}`;
      if (filterPaid !== "all") endpoint += `&is_paid=${filterPaid === "paid"}`;
      if (filterSource !== "all") endpoint += `&source=${filterSource}`;

      const [txData, accs, cats] = await Promise.all([
        api.get<PaginatedResponse<Transaction>>(endpoint),
        api.get<BankAccount[]>("/accounts"),
        api.get<Category[]>("/categories"),
      ]);
      setTransactions(txData?.data || []);
      setTotalPages(txData?.total_pages || 1);
      setAccounts(accs || []);
      setCategories(cats || []);
    } catch {
      // empty
    } finally {
      setIsLoading(false);
    }
  }, [page, filterType, filterPaid, filterSource]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  function resetForm() {
    setForm({
      ...initialFormState,
      transaction_date: new Date().toISOString().split("T")[0],
    });
    setCategorySuggestions([]);
  }

  function handleDescriptionChange(value: string) {
    setForm((prev) => ({ ...prev, description: value }));

    if (suggestDebounceRef.current) {
      clearTimeout(suggestDebounceRef.current);
    }

    if (value.trim().length < 3) {
      setCategorySuggestions([]);
      return;
    }

    suggestDebounceRef.current = setTimeout(async () => {
      try {
        const result = await api.get<{ suggestions: CategorySuggestion[] }>(
          `/categories/suggest?description=${encodeURIComponent(value)}`
        );
        const suggestions = result?.suggestions || [];
        setCategorySuggestions(suggestions);

        if (suggestions.length > 0 && !form.category_id) {
          setForm((prev) => ({
            ...prev,
            category_id: suggestions[0].category_id,
          }));
        }
      } catch {
        // silencia erros de sugestão
      }
    }, 500);
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      const validated = createTransactionSchema.parse({
        bank_account_id: form.bank_account_id,
        category_id: form.category_id || "",
        type: form.type,
        amount: form.amount,
        description: form.description,
        notes: form.notes || "",
        transaction_date: form.transaction_date,
        due_date: form.due_date || "",
        is_paid: form.is_paid,
      });

      const body: CreateTransactionRequest & { auto_pay: boolean } = {
        bank_account_id: validated.bank_account_id,
        category_id: validated.category_id || undefined,
        type: validated.type as CreateTransactionRequest["type"],
        amount: validated.amount,
        description: validated.description,
        transaction_date: new Date(validated.transaction_date).toISOString(),
        due_date: validated.due_date
          ? new Date(validated.due_date).toISOString()
          : undefined,
        is_paid: validated.is_paid,
        auto_pay: form.auto_pay,
      };
      await api.post("/transactions", body);
      toast.success("Transacao criada");
      resetForm();
      setDialogOpen(false);
      fetchData();
    } catch (error) {
      if (error instanceof z.ZodError) {
        toast.error(error.issues[0].message);
      } else {
        toast.error("Erro ao criar transacao");
      }
    }
  }

  async function markAsPaid(id: string) {
    try {
      await api.patch(`/transactions/${id}/pay`);
      toast.success("Marcado como pago");
      fetchData();
    } catch {
      toast.error("Erro ao marcar como pago");
    }
  }

  function isSuggested(categoryId: string): boolean {
    return categorySuggestions.some((s) => s.category_id === categoryId);
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold">Transações</h1>
        <Skeleton className="h-[400px]" />
      </div>
    );
  }

  return (
    <div className="space-y-8 pb-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Transações</h1>
          <p className="text-muted-foreground mt-2">
            Gerencie todas as suas transações
          </p>
        </div>
        <Dialog
          open={dialogOpen}
          onOpenChange={(open) => {
            setDialogOpen(open);
            if (!open) resetForm();
          }}
        >
          <DialogTrigger asChild>
            <Button size="lg">
              <Plus className="h-4 w-4 mr-2" />
              Nova Transacao
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>Nova Transacao</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label>Tipo</Label>
                <Select
                  value={form.type}
                  onValueChange={(v) => setForm({ ...form, type: v })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="income">Receita</SelectItem>
                    <SelectItem value="expense">Despesa</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Conta</Label>
                <Select
                  value={form.bank_account_id}
                  onValueChange={(v) =>
                    setForm({ ...form, bank_account_id: v })
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione..." />
                  </SelectTrigger>
                  <SelectContent>
                    {accounts.map((a) => (
                      <SelectItem key={a.id} value={a.id}>
                        {a.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Descricao</Label>
                <Input
                  value={form.description}
                  onChange={(e) => handleDescriptionChange(e.target.value)}
                  required
                />
              </div>
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <Label>Categoria</Label>
                  {categorySuggestions.length > 0 && (
                    <span className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Sparkles className="h-3 w-3 text-yellow-500" />
                      sugestão automática
                    </span>
                  )}
                </div>
                <Select
                  value={form.category_id}
                  onValueChange={(v) => setForm({ ...form, category_id: v })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Opcional" />
                  </SelectTrigger>
                  <SelectContent>
                    {categorySuggestions.length > 0 && (
                      <>
                        <div className="px-2 py-1 text-xs font-medium text-muted-foreground flex items-center gap-1">
                          <Sparkles className="h-3 w-3 text-yellow-500" />
                          Sugestões
                        </div>
                        {categorySuggestions.map((s) => (
                          <SelectItem
                            key={`sug-${s.category_id}`}
                            value={s.category_id}
                          >
                            ✨ {s.category_name}
                          </SelectItem>
                        ))}
                        <div className="px-2 py-1 text-xs font-medium text-muted-foreground border-t mt-1 pt-1">
                          Todas as categorias
                        </div>
                      </>
                    )}
                    {categories
                      .filter((c) => c.type === form.type)
                      .map((c) => (
                        <SelectItem key={c.id} value={c.id}>
                          {isSuggested(c.id) ? `✨ ${c.name}` : c.name}
                        </SelectItem>
                      ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Valor</Label>
                <Input
                  type="number"
                  step="0.01"
                  min="0.01"
                  value={form.amount}
                  onChange={(e) => setForm({ ...form, amount: e.target.value })}
                  required
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Data</Label>
                  <Input
                    type="date"
                    value={form.transaction_date}
                    onChange={(e) =>
                      setForm({ ...form, transaction_date: e.target.value })
                    }
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label>Vencimento</Label>
                  <Input
                    type="date"
                    value={form.due_date}
                    onChange={(e) =>
                      setForm({ ...form, due_date: e.target.value })
                    }
                  />
                </div>
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="is_paid"
                  checked={form.is_paid}
                  onCheckedChange={(checked) =>
                    setForm({ ...form, is_paid: checked === true })
                  }
                />
                <Label
                  htmlFor="is_paid"
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Marcar como pago
                </Label>
              </div>
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="auto_pay"
                  checked={form.auto_pay}
                  onCheckedChange={(checked) =>
                    setForm({ ...form, auto_pay: checked === true })
                  }
                />
                <Label
                  htmlFor="auto_pay"
                  className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                >
                  Pagamento automático no vencimento
                </Label>
              </div>
              <Button type="submit" className="w-full">
                Criar Transacao
              </Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-wrap gap-3">
            <Select
              value={filterType}
              onValueChange={(v) => {
                setFilterType(v);
                setPage(1);
              }}
            >
              <SelectTrigger className="w-[150px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos tipos</SelectItem>
                <SelectItem value="income">Receita</SelectItem>
                <SelectItem value="expense">Despesa</SelectItem>
              </SelectContent>
            </Select>
            <Select
              value={filterPaid}
              onValueChange={(v) => {
                setFilterPaid(v);
                setPage(1);
              }}
            >
              <SelectTrigger className="w-[150px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos status</SelectItem>
                <SelectItem value="paid">Pago</SelectItem>
                <SelectItem value="unpaid">Pendente</SelectItem>
              </SelectContent>
            </Select>
            <Select
              value={filterSource}
              onValueChange={(v) => {
                setFilterSource(v);
                setPage(1);
              }}
            >
              <SelectTrigger className="w-[150px]">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todas origens</SelectItem>
                <SelectItem value="manual">Manual</SelectItem>
                <SelectItem value="bank_sync">Banco</SelectItem>
                <SelectItem value="recurring">Recorrente</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Data</TableHead>
                <TableHead>Descricao</TableHead>
                <TableHead>Tipo</TableHead>
                <TableHead>Categoria</TableHead>
                <TableHead>Origem</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Valor</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={8}
                    className="text-center text-muted-foreground"
                  >
                    Nenhuma transacao encontrada
                  </TableCell>
                </TableRow>
              ) : (
                transactions.map((tx) => (
                  <TableRow key={tx.id}>
                    <TableCell className="text-sm">
                      {formatDate(tx.transaction_date)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {tx.type === "income" ? (
                          <ArrowDownLeft className="h-4 w-4 text-green-500" />
                        ) : (
                          <ArrowUpRight className="h-4 w-4 text-red-500" />
                        )}
                        <span className="text-sm font-medium">
                          {tx.description}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={tx.type === "income" ? "default" : "secondary"}
                      >
                        {getTransactionTypeLabel(tx.type)}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {tx.category?.name || "-"}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          tx.source === "manual"
                            ? "outline"
                            : tx.source === "bank_sync"
                              ? "default"
                              : "secondary"
                        }
                      >
                        {getTransactionSourceLabel(tx.source)}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant={tx.is_paid ? "default" : "destructive"}>
                        {tx.is_paid ? "Pago" : "Pendente"}
                      </Badge>
                    </TableCell>
                    <TableCell
                      className={`text-right font-semibold ${tx.type === "income" ? "text-green-500" : "text-red-500"}`}
                    >
                      {tx.type === "income" ? "+" : "-"}
                      {formatCurrency(tx.amount)}
                    </TableCell>
                    <TableCell>
                      {!tx.is_paid && (
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => markAsPaid(tx.id)}
                        >
                          <Check className="h-4 w-4" />
                        </Button>
                      )}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-2 mt-4">
              <Button
                variant="outline"
                size="sm"
                disabled={page <= 1}
                onClick={() => setPage(page - 1)}
              >
                Anterior
              </Button>
              <span className="text-sm text-muted-foreground">
                Pagina {page} de {totalPages}
              </span>
              <Button
                variant="outline"
                size="sm"
                disabled={page >= totalPages}
                onClick={() => setPage(page + 1)}
              >
                Proxima
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
