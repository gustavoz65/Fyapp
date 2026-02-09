"use client";

import { useCallback, useEffect, useState } from "react";
import { Plus, Pencil, Trash2, Landmark } from "lucide-react";
import { api } from "@/lib/api";
import type { BankAccount, CreateBankAccountRequest, UpdateBankAccountRequest } from "@/types";
import { formatCurrency, getAccountTypeLabel } from "@/lib/format";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle, AlertDialogTrigger } from "@/components/ui/alert-dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

const accountTypes = [
  { value: "checking", label: "Conta Corrente" },
  { value: "savings", label: "Poupanca" },
  { value: "credit_card", label: "Cartao de Credito" },
  { value: "investment", label: "Investimento" },
  { value: "cash", label: "Dinheiro" },
  { value: "other", label: "Outro" },
];

export default function AccountsPage() {
  const [accounts, setAccounts] = useState<BankAccount[]>([]);
  const [totalBalance, setTotalBalance] = useState("0");
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<BankAccount | null>(null);
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

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    try {
      if (editing) {
        const body: UpdateBankAccountRequest = {
          name: form.name, bank_name: form.bank_name || undefined,
          color: form.color, icon: form.icon,
        };
        await api.put(`/accounts/${editing.id}`, body);
        toast.success("Conta atualizada");
      } else {
        const body: CreateBankAccountRequest = {
          name: form.name, bank_name: form.bank_name || undefined,
          account_type: form.account_type as CreateBankAccountRequest["account_type"],
          initial_balance: form.initial_balance, color: form.color, icon: form.icon,
        };
        await api.post("/accounts", body);
        toast.success("Conta criada");
      }
      setDialogOpen(false);
      fetchAccounts();
    } catch {
      toast.error("Erro ao salvar conta");
    }
  }

  async function handleDelete(id: string) {
    try {
      await api.delete(`/accounts/${id}`);
      toast.success("Conta removida");
      fetchAccounts();
    } catch {
      toast.error("Erro ao remover conta");
    }
  }

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-3xl font-bold">Contas</h1>
        </div>
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => <Skeleton key={i} className="h-40" />)}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Contas</h1>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={openCreate}><Plus className="h-4 w-4 mr-2" />Nova Conta</Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{editing ? "Editar Conta" : "Nova Conta"}</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label>Nome</Label>
                <Input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
              </div>
              <div className="space-y-2">
                <Label>Banco</Label>
                <Input value={form.bank_name} onChange={(e) => setForm({ ...form, bank_name: e.target.value })} />
              </div>
              {!editing && (
                <>
                  <div className="space-y-2">
                    <Label>Tipo</Label>
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
                    <Input type="number" step="0.01" value={form.initial_balance} onChange={(e) => setForm({ ...form, initial_balance: e.target.value })} required />
                  </div>
                </>
              )}
              <div className="space-y-2">
                <Label>Cor</Label>
                <Input type="color" value={form.color} onChange={(e) => setForm({ ...form, color: e.target.value })} />
              </div>
              <Button type="submit" className="w-full">{editing ? "Salvar" : "Criar"}</Button>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Saldo Total</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-bold">{formatCurrency(totalBalance)}</p>
        </CardContent>
      </Card>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {accounts.map((account) => (
          <Card key={account.id}>
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
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button variant="ghost" size="sm" className="text-destructive">
                      <Trash2 className="h-3 w-3 mr-1" />Remover
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Remover conta?</AlertDialogTitle>
                      <AlertDialogDescription>
                        Essa acao nao pode ser desfeita. A conta &quot;{account.name}&quot; sera removida.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancelar</AlertDialogCancel>
                      <AlertDialogAction onClick={() => handleDelete(account.id)}>Remover</AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
