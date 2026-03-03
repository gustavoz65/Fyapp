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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
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
import { Progress } from "@/components/ui/progress";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { api, getApiBaseUrl, getWebSocketUrl } from "@/lib/api";
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
import {
  ArrowDownLeft,
  ArrowUpRight,
  Check,
  Loader2,
  MoreVertical,
  Pencil,
  Plus,
  ReceiptText,
  Sparkles,
  Upload,
} from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

interface CategorySuggestion {
  category_id: string;
  category_name: string;
  confidence: number;
  match_count: number;
}

function getApiErrorMessage(error: unknown, fallback: string): string {
  if (error && typeof error === "object" && "message" in error) {
    return (error as { message: string }).message;
  }
  return fallback;
}

export default function TransactionsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [accounts, setAccounts] = useState<BankAccount[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isImporting, setIsImporting] = useState(false);
  const [importProgress, setImportProgress] = useState<number>(0);
  const [importMessage, setImportMessage] = useState<string>("");
  const [markingPaidId, setMarkingPaidId] = useState<string | null>(null);
  const [editingCategoryId, setEditingCategoryId] = useState<string | null>(null);
  const [editCategoryDialogOpen, setEditCategoryDialogOpen] = useState(false);
  const [selectedCategoryForEdit, setSelectedCategoryForEdit] = useState<string>("");
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
    notes: "",
    transaction_date: new Date().toISOString().split("T")[0],
    due_date: "",
    is_paid: false,
    auto_pay: false,
  };

  const [form, setForm] = useState(initialFormState);
  const [importForm, setImportForm] = useState({
    bank_account_id: "",
    bank_type: "generic",
    file: null as File | null,
  });

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

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setIsSubmitting(true);
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
        notes: validated.notes || undefined,
        transaction_date: new Date(validated.transaction_date).toISOString(),
        due_date: validated.due_date
          ? new Date(validated.due_date).toISOString()
          : undefined,
        is_paid: validated.is_paid,
        auto_pay: form.auto_pay,
      };
      await api.post("/transactions", body);
      toast.success("Transação criada com sucesso");
      resetForm();
      setDialogOpen(false);
      fetchData();
    } catch (error) {
      if (error instanceof z.ZodError) {
        toast.error(error.issues[0].message);
      } else {
        toast.error(getApiErrorMessage(error, "Erro ao criar transação. Tente novamente."));
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  async function markAsPaid(id: string) {
    setMarkingPaidId(id);
    try {
      await api.patch(`/transactions/${id}/pay`);
      toast.success("Transação marcada como paga");
      fetchData();
    } catch (error) {
      toast.error(getApiErrorMessage(error, "Erro ao marcar como pago."));
    } finally {
      setMarkingPaidId(null);
    }
  }

  async function handleUpdateCategory(txId: string, categoryId: string) {
    try {
      await api.put(`/transactions/${txId}`, {
        category_id: categoryId || undefined,
      });
      toast.success("Categoria atualizada com sucesso");
      setEditCategoryDialogOpen(false);
      setEditingCategoryId(null);
      setSelectedCategoryForEdit("");
      fetchData();
    } catch (error) {
      toast.error(getApiErrorMessage(error, "Erro ao atualizar categoria"));
    }
  }

  function openEditCategoryDialog(tx: Transaction) {
    setEditingCategoryId(tx.id);
    setSelectedCategoryForEdit(tx.category?.id || "");
    setEditCategoryDialogOpen(true);
  }

  async function handleImport(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!importForm.file || !importForm.bank_account_id) {
      toast.error("Selecione a conta e o arquivo CSV");
      return;
    }

    setIsImporting(true);
    setImportProgress(0);
    setImportMessage("Enviando arquivo...");

    try {
      const formData = new FormData();
      formData.append("file", importForm.file);
      formData.append("bank_account_id", importForm.bank_account_id);
      formData.append("bank_type", importForm.bank_type);

      const response = await fetch(`${getApiBaseUrl()}/api/v1/transactions/import`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${localStorage.getItem("access_token")}`,
        },
        body: formData,
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Erro ao importar");
      }

      const result = await response.json();
      const jobId = result.job_id;

      setImportMessage("Processando transações...");

      // Connect to WebSocket for real-time progress
      const wsUrl = `${getWebSocketUrl()}/ws/import-progress?job_id=${jobId}`;
      const ws = new WebSocket(wsUrl);

      ws.onopen = () => {
        console.log("WebSocket connected");
      };

      ws.onmessage = (event) => {
        const status = JSON.parse(event.data);

        if (status.error) {
          toast.error(status.error);
          ws.close();
          setIsImporting(false);
          return;
        }

        setImportProgress(status.progress || 0);
        setImportMessage(status.message || "Processando...");

        if (status.status === "completed") {
          toast.success(
            `Importação concluída! ${status.imported || 0} novas, ${status.duplicates || 0} duplicadas.`
          );
          ws.close();
          setImportDialogOpen(false);
          setImportForm({ bank_account_id: "", bank_type: "generic", file: null });
          setIsImporting(false);
          setImportProgress(0);
          setImportMessage("");
          fetchData();
        } else if (status.status === "failed") {
          toast.error(status.message || "Erro ao importar transações");
          ws.close();
          setIsImporting(false);
          setImportProgress(0);
          setImportMessage("");
        }
      };

      ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        toast.error("Erro na conexão WebSocket. Usando polling...");
        ws.close();
        // Fallback to polling
        pollImportStatus(jobId);
      };

      ws.onclose = () => {
        console.log("WebSocket disconnected");
      };

    } catch (error) {
      toast.error(getApiErrorMessage(error, "Erro ao importar transações"));
      setIsImporting(false);
      setImportProgress(0);
      setImportMessage("");
    }
  }

  // Fallback polling function if WebSocket fails
  async function pollImportStatus(jobId: string) {
    const interval = setInterval(async () => {
      try {
        const response = await fetch(`${getApiBaseUrl()}/api/v1/transactions/import/${jobId}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("access_token")}`,
          },
        });

        if (!response.ok) {
          clearInterval(interval);
          setIsImporting(false);
          toast.error("Job não encontrado");
          return;
        }

        const status = await response.json();
        setImportProgress(status.progress || 0);
        setImportMessage(status.message || "Processando...");

        if (status.status === "completed") {
          clearInterval(interval);
          toast.success(
            `Importação concluída! ${status.imported || 0} novas, ${status.duplicates || 0} duplicadas.`
          );
          setImportDialogOpen(false);
          setImportForm({ bank_account_id: "", bank_type: "generic", file: null });
          setIsImporting(false);
          setImportProgress(0);
          setImportMessage("");
          fetchData();
        } else if (status.status === "failed") {
          clearInterval(interval);
          toast.error(status.message || "Erro ao importar transações");
          setIsImporting(false);
          setImportProgress(0);
          setImportMessage("");
        }
      } catch (error) {
        clearInterval(interval);
        toast.error("Erro ao consultar status da importação");
        setIsImporting(false);
      }
    }, 1000); // Poll every second
  }

  function isSuggested(categoryId: string): boolean {
    return categorySuggestions.some((s) => s.category_id === categoryId);
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="space-y-2">
            <Skeleton className="h-9 w-40" />
            <Skeleton className="h-4 w-56" />
          </div>
          <Skeleton className="h-10 w-40" />
        </div>
        <Card>
          <CardHeader>
            <div className="flex gap-3">
              <Skeleton className="h-9 w-36" />
              <Skeleton className="h-9 w-36" />
              <Skeleton className="h-9 w-36" />
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {["a", "b", "c", "d", "e"].map((k) => (
                <Skeleton key={k} className="h-12 w-full" />
              ))}
            </div>
          </CardContent>
        </Card>
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
        <div className="flex gap-2">
          <Dialog open={importDialogOpen} onOpenChange={setImportDialogOpen}>
            <DialogTrigger asChild>
              <Button size="lg" variant="outline">
                <Upload className="h-4 w-4 mr-2" />
                Importar Extrato
              </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md">
              <DialogHeader>
                <DialogTitle>Importar Extrato Bancário</DialogTitle>
              </DialogHeader>
              {isImporting && (
                <div className="space-y-2">
                  <div className="flex items-center justify-between text-sm">
                    <span className="text-muted-foreground">{importMessage}</span>
                    <span className="font-medium">{importProgress}%</span>
                  </div>
                  <Progress value={importProgress} className="h-2" />
                </div>
              )}
              <form onSubmit={handleImport} className="space-y-4">
                <div className="space-y-2">
                  <Label>Conta *</Label>
                  <Select
                    value={importForm.bank_account_id}
                    onValueChange={(v) =>
                      setImportForm({ ...importForm, bank_account_id: v })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="Selecione a conta..." />
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
                  <Label>Tipo de Banco</Label>
                  <Select
                    value={importForm.bank_type}
                    onValueChange={(v) =>
                      setImportForm({ ...importForm, bank_type: v })
                    }
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="generic">Genérico (CSV padrão)</SelectItem>
                      <SelectItem value="nubank">Nubank</SelectItem>
                      <SelectItem value="inter">Inter</SelectItem>
                      <SelectItem value="itau">Itaú</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <div className="space-y-2">
                  <Label>Arquivo CSV *</Label>
                  <Input
                    type="file"
                    accept=".csv"
                    onChange={(e) =>
                      setImportForm({
                        ...importForm,
                        file: e.target.files?.[0] || null,
                      })
                    }
                    required
                  />
                  <p className="text-xs text-muted-foreground">
                    Formatos aceitos: CSV (máx. 10MB)
                  </p>
                </div>
                <div className="flex gap-2 pt-2">
                  <Button
                    type="button"
                    variant="outline"
                    className="flex-1"
                    onClick={() => {
                      setImportDialogOpen(false);
                      setImportForm({
                        bank_account_id: "",
                        bank_type: "generic",
                        file: null,
                      });
                    }}
                  >
                    Cancelar
                  </Button>
                  <Button type="submit" className="flex-1" disabled={isImporting}>
                    {isImporting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                    {isImporting ? "Importando..." : "Importar"}
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
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
                Nova Transação
              </Button>
            </DialogTrigger>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>Nova Transação</DialogTitle>
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
                <Label>Conta *</Label>
                <Select
                  value={form.bank_account_id}
                  onValueChange={(v) => setForm({ ...form, bank_account_id: v })}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione a conta..." />
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
                <Label>Descrição *</Label>
                <Input
                  value={form.description}
                  onChange={(e) => handleDescriptionChange(e.target.value)}
                  placeholder="Ex: Supermercado, Salário..."
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
                    <SelectValue placeholder="Selecione (opcional)" />
                  </SelectTrigger>
                  <SelectContent>
                    {categorySuggestions.length > 0 && (
                      <>
                        <div className="px-2 py-1 text-xs font-medium text-muted-foreground flex items-center gap-1">
                          <Sparkles className="h-3 w-3 text-yellow-500" />
                          Sugestões
                        </div>
                        {categorySuggestions.map((s) => (
                          <SelectItem key={`sug-${s.category_id}`} value={s.category_id}>
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
                <Label>Valor (R$) *</Label>
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
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Data da transação</Label>
                  <Input
                    type="date"
                    value={form.transaction_date}
                    onChange={(e) => setForm({ ...form, transaction_date: e.target.value })}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label>Data de vencimento</Label>
                  <Input
                    type="date"
                    value={form.due_date}
                    onChange={(e) => setForm({ ...form, due_date: e.target.value })}
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
                <Label htmlFor="is_paid" className="text-sm font-medium leading-none">
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
                <Label htmlFor="auto_pay" className="text-sm font-medium leading-none">
                  Pagamento automático no vencimento
                </Label>
              </div>
              <div className="flex gap-2 pt-2">
                <Button
                  type="button"
                  variant="outline"
                  className="flex-1"
                  onClick={() => {
                    setDialogOpen(false);
                    resetForm();
                  }}
                >
                  Cancelar
                </Button>
                <Button type="submit" className="flex-1" disabled={isSubmitting}>
                  {isSubmitting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                  Criar transação
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-wrap gap-3">
            <Select
              value={filterType}
              onValueChange={(v) => { setFilterType(v); setPage(1); }}
            >
              <SelectTrigger className="w-36">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os tipos</SelectItem>
                <SelectItem value="income">Receita</SelectItem>
                <SelectItem value="expense">Despesa</SelectItem>
              </SelectContent>
            </Select>
            <Select
              value={filterPaid}
              onValueChange={(v) => { setFilterPaid(v); setPage(1); }}
            >
              <SelectTrigger className="w-36">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos os status</SelectItem>
                <SelectItem value="paid">Pago</SelectItem>
                <SelectItem value="unpaid">Pendente</SelectItem>
              </SelectContent>
            </Select>
            <Select
              value={filterSource}
              onValueChange={(v) => { setFilterSource(v); setPage(1); }}
            >
              <SelectTrigger className="w-36">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todas as origens</SelectItem>
                <SelectItem value="manual">Manual</SelectItem>
                <SelectItem value="bank_sync">Banco</SelectItem>
                <SelectItem value="recurring">Recorrente</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardHeader>
        <CardContent>
          {transactions.length === 0 ? (
            <div className="py-16 flex flex-col items-center gap-4 text-center">
              <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
                <ReceiptText className="h-8 w-8 text-muted-foreground" />
              </div>
              <div>
                <h3 className="font-semibold text-lg">Nenhuma transação encontrada</h3>
                <p className="text-sm text-muted-foreground mt-1">
                  {filterType !== "all" || filterPaid !== "all" || filterSource !== "all"
                    ? "Tente ajustar os filtros para ver mais resultados"
                    : "Registre sua primeira transação para começar a acompanhar suas finanças"}
                </p>
              </div>
              {filterType === "all" && filterPaid === "all" && filterSource === "all" && (
                <Button onClick={() => setDialogOpen(true)}>
                  <Plus className="h-4 w-4 mr-2" />
                  Adicionar primeira transação
                </Button>
              )}
            </div>
          ) : (
            <>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Data</TableHead>
                    <TableHead>Descrição</TableHead>
                    <TableHead>Tipo</TableHead>
                    <TableHead>Categoria</TableHead>
                    <TableHead>Origem</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Valor</TableHead>
                    <TableHead></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {transactions.map((tx) => (
                    <TableRow key={tx.id}>
                      <TableCell className="text-sm">
                        {formatDate(tx.transaction_date)}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          {tx.type === "income" ? (
                            <ArrowDownLeft className="h-4 w-4 text-green-500 shrink-0" />
                          ) : (
                            <ArrowUpRight className="h-4 w-4 text-red-500 shrink-0" />
                          )}
                          <span className="text-sm font-medium">{tx.description}</span>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant={tx.type === "income" ? "default" : "secondary"}>
                          {getTransactionTypeLabel(tx.type)}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        {tx.category?.name || "—"}
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
                        className={`text-right font-semibold ${tx.type === "income" ? "text-green-600" : "text-red-600"}`}
                      >
                        {tx.type === "income" ? "+" : "-"}
                        {formatCurrency(tx.amount)}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-1">
                          {!tx.is_paid && (
                            <Button
                              variant="ghost"
                              size="sm"
                              disabled={markingPaidId === tx.id}
                              onClick={() => markAsPaid(tx.id)}
                              title="Marcar como pago"
                            >
                              {markingPaidId === tx.id ? (
                                <Loader2 className="h-4 w-4 animate-spin" />
                              ) : (
                                <Check className="h-4 w-4" />
                              )}
                            </Button>
                          )}
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="sm">
                                <MoreVertical className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem onClick={() => openEditCategoryDialog(tx)}>
                                <Pencil className="h-4 w-4 mr-2" />
                                Editar Categoria
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
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
                    Página {page} de {totalPages}
                  </span>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={page >= totalPages}
                    onClick={() => setPage(page + 1)}
                  >
                    Próxima
                  </Button>
                </div>
              )}
            </>
          )}
        </CardContent>
      </Card>

      {/* Dialog para editar categoria */}
      <Dialog open={editCategoryDialogOpen} onOpenChange={setEditCategoryDialogOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Editar Categoria</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Categoria</Label>
              <Select
                value={selectedCategoryForEdit}
                onValueChange={setSelectedCategoryForEdit}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Selecione uma categoria..." />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Sem categoria</SelectItem>
                  {categories.map((c) => (
                    <SelectItem key={c.id} value={c.id}>
                      {c.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            <div className="flex gap-2">
              <Button
                type="button"
                variant="outline"
                className="flex-1"
                onClick={() => {
                  setEditCategoryDialogOpen(false);
                  setEditingCategoryId(null);
                  setSelectedCategoryForEdit("");
                }}
              >
                Cancelar
              </Button>
              <Button
                type="button"
                className="flex-1"
                onClick={() => {
                  if (editingCategoryId) {
                    handleUpdateCategory(editingCategoryId, selectedCategoryForEdit);
                  }
                }}
              >
                Salvar
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
