"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { z } from "zod";
import { Landmark, Loader2, Pencil, Plus, Trash2 } from "lucide-react";
import { api } from "@/lib/api";
import type { BankAccount, CreateBankAccountRequest, UpdateBankAccountRequest } from "@/types";
import { formatCurrency, getAccountTypeLabel } from "@/lib/format";
import { createAccountSchema } from "@/lib/schemas";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

const accountTypes = [
  { value: "checking", label: "Conta Corrente" },
  { value: "savings", label: "Poupança" },
  { value: "credit_card", label: "Cartão de Crédito" },
  { value: "investment", label: "Investimento" },
  { value: "cash", label: "Dinheiro" },
  { value: "other", label: "Outro" },
];

function getApiErrorMessage(error: unknown, fallback: string): string {
  if (error && typeof error === "object" && "message" in error) {
    return (error as { message: string }).message;
  }
  return fallback;
}

export default function AccountsPage() {
  const [accounts, setAccounts] = useState<BankAccount[]>([]);
  const [totalBalance, setTotalBalance] = useState("0");
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<BankAccount | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const pendingDeleteRef = useRef<{ id: string; timer: ReturnType<typeof setTimeout> } | null>(null);

  const [form, setForm] = useState({
    name: "", bank_name: "", account_type: "checking" as string,
    initial_balance: "0", color: "#3b82f6", icon: "landmark",
  });

  const fetchAccounts = useCallback(async () => {
    try {
      const [accs, balance] = await Promise.all([
        api.get<BankAccount[]>("/accounts"),
        api.get<{ total_balance: string }>("/accounts/balance"),
      ]);
      setAccounts(accs || []);
      setTotalBalance(balance?.total_balance || "0");
    } catch {
      // empty
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => { fetchAccounts(); }, [fetchAccounts]);

  function openCreate() {
    setEditing(null);
    setForm({ name: "", bank_name: "", account_type: "checking", initial_balance: "0", color: "#3b82f6", icon: "landmark" });
    setDialogOpen(true);
  }

  function openEdit(account: BankAccount) {
    setEditing(account);
    setForm({
      name: account.name, bank_name: account.bank_name || "",
      account_type: account.account_type, initial_balance: account.initial_balance,
      color: account.color, icon: account.icon,
    });
    setDialogOpen(true);
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      if (editing) {
        const body: UpdateBankAccountRequest = {
          name: form.name, bank_name: form.bank_name || undefined,
          color: form.color, icon: form.icon,
        };
        await api.put(`/accounts/${editing.id}`, body);
        toast.success("Conta atualizada com sucesso");
      } else {
        const validated = createAccountSchema.parse({
          name: form.name,
          bank_name: form.bank_name || undefined,
          account_type: form.account_type,
          initial_balance: form.initial_balance,
          color: form.color,
          icon: form.icon,
        });

        const body: CreateBankAccountRequest = {
          name: validated.name,
          bank_name: validated.bank_name,
          account_type: validated.account_type as CreateBankAccountRequest["account_type"],
          initial_balance: validated.initial_balance,
          color: validated.color,
          icon: validated.icon,
        };
        await api.post("/accounts", body);
        toast.success("Conta criada com sucesso");
      }
      setDialogOpen(false);
      fetchAccounts();
    } catch (error) {
      if (error instanceof z.ZodError) {
        toast.error(error.issues[0].message);
      } else {
        toast.error(getApiErrorMessage(error, "Erro ao salvar conta. Tente novamente."));
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  function handleDelete(id: string) {
    const accountToDelete = accounts.find((a) => a.id === id);
    if (!accountToDelete) return;

    // Remove optimisticamente
    setAccounts((prev) => prev.filter((a) => a.id !== id));

    // Cancela delete pendente anterior
    if (pendingDeleteRef.current) {
      clearTimeout(pendingDeleteRef.current.timer);
    }

    const timer = setTimeout(async () => {
      try {
        await api.delete(`/accounts/${id}`);
        pendingDeleteRef.current = null;
        fetchAccounts();
      } catch {
        setAccounts((prev) => [...prev, accountToDelete]);
        toast.error("Erro ao remover conta. Tente novamente.");
      }
    }, 5000);

    pendingDeleteRef.current = { id, timer };

    toast("Conta removida", {
      description: accountToDelete.name,
      action: {
        label: "Desfazer",
        onClick: () => {
          if (pendingDeleteRef.current?.id === id) {
            clearTimeout(pendingDeleteRef.current.timer);
            pendingDeleteRef.current = null;
            setAccounts((prev) => [...prev, accountToDelete]);
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
            <Skeleton className="h-9 w-48" />
            <Skeleton className="h-4 w-64" />
          </div>
          <Skeleton className="h-10 w-32" />
        </div>
        <Skeleton className="h-28" />
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => <Skeleton key={i} className="h-40" />)}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8 pb-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Contas Bancárias</h1>
          <p className="text-muted-foreground mt-2">
            Gerencie suas contas e acompanhe seus saldos
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button size="lg" onClick={openCreate}>
              <Plus className="h-4 w-4 mr-2" />Nova Conta
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editing ? "Editar Conta" : "Nova Conta"}</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label>Nome da conta</Label>
                <Input
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Ex: Nubank, Bradesco..."
                  required
                />
              </div>
              <div className="space-y-2">
                <Label>Banco</Label>
                <Input
                  value={form.bank_name}
                  onChange={(e) => setForm({ ...form, bank_name: e.target.value })}
                  placeholder="Nome do banco (opcional)"
                />
              </div>
              {!editing && (
                <>
                  <div className="space-y-2">
                    <Label>Tipo de conta</Label>
                    <Select value={form.account_type} onValueChange={(v) => setForm({ ...form, account_type: v })}>
                      <SelectTrigger><SelectValue /></SelectTrigger>
                      <SelectContent>
                        {accountTypes.map((t) => (
                          <SelectItem key={t.value} value={t.value}>{t.label}</SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label>Saldo Inicial</Label>
                    <Input
                      type="number"
                      step="0.01"
                      value={form.initial_balance}
                      onChange={(e) => setForm({ ...form, initial_balance: e.target.value })}
                      placeholder="0,00"
                      required
                    />
                  </div>
                </>
              )}
              <div className="space-y-2">
                <Label>Cor de identificação</Label>
                <div className="flex items-center gap-3">
                  <Input
                    type="color"
                    value={form.color}
                    onChange={(e) => setForm({ ...form, color: e.target.value })}
                    className="h-10 w-16 p-1 cursor-pointer"
                  />
                  <span className="text-sm text-muted-foreground">Escolha uma cor para identificar a conta</span>
                </div>
              </div>
              <div className="flex gap-2 pt-2">
                <Button type="button" variant="outline" className="flex-1" onClick={() => setDialogOpen(false)}>
                  Cancelar
                </Button>
                <Button type="submit" className="flex-1" disabled={isSubmitting}>
                  {isSubmitting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                  {editing ? "Salvar alterações" : "Criar conta"}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card className="border-2">
        <CardHeader>
          <CardTitle className="text-sm font-medium text-muted-foreground">Saldo Total</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-4xl font-bold tracking-tight">{formatCurrency(totalBalance)}</p>
          <p className="text-sm text-muted-foreground mt-2">{accounts.length} conta{accounts.length !== 1 ? "s" : ""} ativa{accounts.length !== 1 ? "s" : ""}</p>
        </CardContent>
      </Card>

      {accounts.length === 0 ? (
        <Card className="border-dashed">
          <CardContent className="py-16 flex flex-col items-center gap-4 text-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-muted">
              <Landmark className="h-8 w-8 text-muted-foreground" />
            </div>
            <div>
              <h3 className="font-semibold text-lg">Nenhuma conta cadastrada</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Adicione sua primeira conta para começar a controlar suas finanças
              </p>
            </div>
            <Button onClick={openCreate}>
              <Plus className="h-4 w-4 mr-2" />
              Adicionar primeira conta
            </Button>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {accounts.map((account) => (
            <Card key={account.id} className="hover:shadow-md transition-shadow">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <div className="flex items-center gap-2">
                  <div className="flex h-8 w-8 items-center justify-center rounded-lg" style={{ backgroundColor: account.color + "20", color: account.color }}>
                    <Landmark className="h-4 w-4" />
                  </div>
                  <div>
                    <CardTitle className="text-sm font-medium">{account.name}</CardTitle>
                    {account.bank_name && <p className="text-xs text-muted-foreground">{account.bank_name}</p>}
                  </div>
                </div>
                <Badge variant="outline">{getAccountTypeLabel(account.account_type)}</Badge>
              </CardHeader>
              <CardContent>
                <p className="text-2xl font-bold">{formatCurrency(account.current_balance)}</p>
                <div className="mt-3 flex gap-2">
                  <Button variant="ghost" size="sm" onClick={() => openEdit(account)}>
                    <Pencil className="h-3 w-3 mr-1" />Editar
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="text-destructive hover:text-destructive"
                    onClick={() => handleDelete(account.id)}
                  >
                    <Trash2 className="h-3 w-3 mr-1" />Remover
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
