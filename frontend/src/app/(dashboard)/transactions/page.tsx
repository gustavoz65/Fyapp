"use client";

import { useCallback, useEffect, useState } from "react";
import { Plus, Check, ArrowUpRight, ArrowDownLeft } from "lucide-react";
import { api } from "@/lib/api";
import type { Transaction, BankAccount, Category, PaginatedResponse, CreateTransactionRequest } from "@/types";
import { formatCurrency, formatDate, getTransactionTypeLabel } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

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

  const [form, setForm] = useState({
    bank_account_id: "", category_id: "", type: "expense" as string,
    amount: "", description: "", transaction_date: new Date().toISOString().split("T")[0],
    due_date: "", is_paid: true,
  });

  const fetchData = useCallback(async () => {
    try {
      let endpoint = `/transactions?page=${page}&page_size=20`;
      if (filterType !== "all") endpoint += `&type=${filterType}`;
      if (filterPaid !== "all") endpoint += `&is_paid=${filterPaid === "paid"}`;

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
  }, [page, filterType, filterPaid]);

  useEffect(() => { fetchData(); }, [fetchData]);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      const body: CreateTransactionRequest = {
        bank_account_id: form.bank_account_id,
        category_id: form.category_id || undefined,
        type: form.type as CreateTransactionRequest["type"],
        amount: form.amount,
        description: form.description,
        transaction_date: new Date(form.transaction_date).toISOString(),
        due_date: form.due_date ? new Date(form.due_date).toISOString() : undefined,
        is_paid: form.is_paid,
      };
      await api.post("/transactions", body);
      toast.success("Transacao criada");
      setDialogOpen(false);
      fetchData();
    } catch {
      toast.error("Erro ao criar transacao");
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

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold">Transacoes</h1>
        <Skeleton className="h-[400px]" />
      </div>
    );
  }

  return (
    <div className="space-y-8 pb-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Transacoes</h1>
          <p className="text-muted-foreground mt-2">
            Gerencie todas as suas transacoes
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button size="lg"><Plus className="h-4 w-4 mr-2" />Nova Transacao</Button>
          </DialogTrigger>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>Nova Transacao</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label>Tipo</Label>
                <Select value={form.type} onValueChange={(v) => setForm({ ...form, type: v })}>
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="income">Receita</SelectItem>
                    <SelectItem value="expense">Despesa</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Conta</Label>
                <Select value={form.bank_account_id} onValueChange={(v) => setForm({ ...form, bank_account_id: v })}>
                  <SelectTrigger><SelectValue placeholder="Selecione..." /></SelectTrigger>
                  <SelectContent>
                    {accounts.map((a) => <SelectItem key={a.id} value={a.id}>{a.name}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Categoria</Label>
                <Select value={form.category_id} onValueChange={(v) => setForm({ ...form, category_id: v })}>
                  <SelectTrigger><SelectValue placeholder="Opcional" /></SelectTrigger>
                  <SelectContent>
                    {categories.filter((c) => c.type === form.type).map((c) => (
                      <SelectItem key={c.id} value={c.id}>{c.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label>Descricao</Label>
                <Input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} required />
              </div>
              <div className="space-y-2">
                <Label>Valor</Label>
                <Input type="number" step="0.01" min="0.01" value={form.amount} onChange={(e) => setForm({ ...form, amount: e.target.value })} required />
              </div>
              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label>Data</Label>
                  <Input type="date" value={form.transaction_date} onChange={(e) => setForm({ ...form, transaction_date: e.target.value })} required />
                </div>
                <div className="space-y-2">
                  <Label>Vencimento</Label>
                  <Input type="date" value={form.due_date} onChange={(e) => setForm({ ...form, due_date: e.target.value })} />
                </div>
              </div>
              <Button type="submit" className="w-full">Criar Transacao</Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-wrap gap-3">
            <Select value={filterType} onValueChange={(v) => { setFilterType(v); setPage(1); }}>
              <SelectTrigger className="w-[150px]"><SelectValue /></SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos tipos</SelectItem>
                <SelectItem value="income">Receita</SelectItem>
                <SelectItem value="expense">Despesa</SelectItem>
              </SelectContent>
            </Select>
            <Select value={filterPaid} onValueChange={(v) => { setFilterPaid(v); setPage(1); }}>
              <SelectTrigger className="w-[150px]"><SelectValue /></SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Todos status</SelectItem>
                <SelectItem value="paid">Pago</SelectItem>
                <SelectItem value="unpaid">Pendente</SelectItem>
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
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Valor</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {transactions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center text-muted-foreground">
                    Nenhuma transacao encontrada
                  </TableCell>
                </TableRow>
              ) : (
                transactions.map((tx) => (
                  <TableRow key={tx.id}>
                    <TableCell className="text-sm">{formatDate(tx.transaction_date)}</TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {tx.type === "income" ? (
                          <ArrowDownLeft className="h-4 w-4 text-green-500" />
                        ) : (
                          <ArrowUpRight className="h-4 w-4 text-red-500" />
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
                      {tx.category?.name || "-"}
                    </TableCell>
                    <TableCell>
                      <Badge variant={tx.is_paid ? "default" : "destructive"}>
                        {tx.is_paid ? "Pago" : "Pendente"}
                      </Badge>
                    </TableCell>
                    <TableCell className={`text-right font-semibold ${tx.type === "income" ? "text-green-500" : "text-red-500"}`}>
                      {tx.type === "income" ? "+" : "-"}{formatCurrency(tx.amount)}
                    </TableCell>
                    <TableCell>
                      {!tx.is_paid && (
                        <Button variant="ghost" size="sm" onClick={() => markAsPaid(tx.id)}>
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
              <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>
                Anterior
              </Button>
              <span className="text-sm text-muted-foreground">
                Pagina {page} de {totalPages}
              </span>
              <Button variant="outline" size="sm" disabled={page >= totalPages} onClick={() => setPage(page + 1)}>
                Proxima
              </Button>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
