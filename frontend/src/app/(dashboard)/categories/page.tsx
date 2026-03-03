"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
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
import { api } from "@/lib/api";
import type {
  Category,
  CreateCategoryRequest,
  UpdateCategoryRequest,
} from "@/types";
import { Loader2, Pencil, Plus, Tags, Trash2 } from "lucide-react";
import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

function getApiErrorMessage(error: unknown, fallback: string): string {
  if (error && typeof error === "object" && "message" in error) {
    return (error as { message: string }).message;
  }
  return fallback;
}

export default function CategoriesPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<Category | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const pendingDeleteRef = useRef<{ id: string; timer: ReturnType<typeof setTimeout> } | null>(null);

  const [form, setForm] = useState({
    name: "",
    description: "",
    type: "expense" as string,
    color: "#3b82f6",
    icon: "tag",
  });

  const fetchCategories = useCallback(async () => {
    try {
      const data = await api.get<Category[]>("/categories");
      setCategories(data || []);
    } catch {
      // empty
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  function openCreate() {
    setEditing(null);
    setForm({ name: "", description: "", type: "expense", color: "#3b82f6", icon: "tag" });
    setDialogOpen(true);
  }

  function openEdit(cat: Category) {
    setEditing(cat);
    setForm({
      name: cat.name,
      description: cat.description || "",
      type: cat.type,
      color: cat.color,
      icon: cat.icon,
    });
    setDialogOpen(true);
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      if (editing) {
        const body: UpdateCategoryRequest = {
          name: form.name,
          description: form.description || undefined,
          color: form.color,
          icon: form.icon,
        };
        await api.put(`/categories/${editing.id}`, body);
        toast.success("Categoria atualizada com sucesso");
      } else {
        const body: CreateCategoryRequest = {
          name: form.name,
          description: form.description || undefined,
          type: form.type as CreateCategoryRequest["type"],
          color: form.color,
          icon: form.icon,
        };
        await api.post("/categories", body);
        toast.success("Categoria criada com sucesso");
      }
      setDialogOpen(false);
      fetchCategories();
    } catch (error) {
      toast.error(getApiErrorMessage(error, "Erro ao salvar categoria. Tente novamente."));
    } finally {
      setIsSubmitting(false);
    }
  }

  function handleDelete(id: string) {
    const catToDelete = categories.find((c) => c.id === id);
    if (!catToDelete) return;

    setCategories((prev) => prev.filter((c) => c.id !== id));

    if (pendingDeleteRef.current) {
      clearTimeout(pendingDeleteRef.current.timer);
    }

    const timer = setTimeout(async () => {
      try {
        await api.delete(`/categories/${id}`);
        pendingDeleteRef.current = null;
        fetchCategories();
      } catch {
        setCategories((prev) => [...prev, catToDelete]);
        toast.error("Erro ao remover categoria. Tente novamente.");
      }
    }, 5000);

    pendingDeleteRef.current = { id, timer };

    toast("Categoria removida", {
      description: catToDelete.name,
      action: {
        label: "Desfazer",
        onClick: () => {
          if (pendingDeleteRef.current?.id === id) {
            clearTimeout(pendingDeleteRef.current.timer);
            pendingDeleteRef.current = null;
            setCategories((prev) => [...prev, catToDelete]);
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
            <Skeleton className="h-4 w-64" />
          </div>
          <Skeleton className="h-10 w-40" />
        </div>
        <Skeleton className="h-6 w-24" />
        <div className="grid gap-4 md:grid-cols-3 lg:grid-cols-4">
          {["a", "b", "c", "d"].map((k) => (
            <Skeleton key={k} className="h-24" />
          ))}
        </div>
      </div>
    );
  }

  const incomeCategories = categories.filter((c) => c.type === "income");
  const expenseCategories = categories.filter((c) => c.type === "expense");

  return (
    <div className="space-y-8 pb-8">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight">Categorias</h1>
          <p className="text-muted-foreground mt-2">
            Organize suas transações por categoria
          </p>
        </div>
        <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
          <DialogTrigger asChild>
            <Button size="lg" onClick={openCreate}>
              <Plus className="h-4 w-4 mr-2" />
              Nova Categoria
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>
                {editing ? "Editar Categoria" : "Nova Categoria"}
              </DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label>Nome</Label>
                <Input
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Ex: Alimentação, Transporte..."
                  required
                />
              </div>
              <div className="space-y-2">
                <Label>Descrição</Label>
                <Input
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                  placeholder="Descrição opcional"
                />
              </div>
              {!editing && (
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
                  <span className="text-sm text-muted-foreground">Escolha uma cor para identificar a categoria</span>
                </div>
              </div>
              <div className="flex gap-2 pt-2">
                <Button type="button" variant="outline" className="flex-1" onClick={() => setDialogOpen(false)}>
                  Cancelar
                </Button>
                <Button type="submit" className="flex-1" disabled={isSubmitting}>
                  {isSubmitting && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
                  {editing ? "Salvar alterações" : "Criar categoria"}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {[
        { title: "Receitas", items: incomeCategories, type: "income" },
        { title: "Despesas", items: expenseCategories, type: "expense" },
      ].map((section) => (
        <div key={section.title} className="space-y-4">
          <h2 className="text-2xl font-semibold tracking-tight">
            {section.title}
          </h2>
          {section.items.length === 0 ? (
            <Card className="border-dashed">
              <CardContent className="py-10 flex flex-col items-center gap-3 text-center">
                <Tags className="h-8 w-8 text-muted-foreground" />
                <div>
                  <p className="text-sm font-medium">Nenhuma categoria de {section.title.toLowerCase()}</p>
                  <p className="text-xs text-muted-foreground mt-1">Crie uma categoria para organizar suas transações</p>
                </div>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => {
                    setForm({ name: "", description: "", type: section.type, color: "#3b82f6", icon: "tag" });
                    setEditing(null);
                    setDialogOpen(true);
                  }}
                >
                  <Plus className="h-3 w-3 mr-1" />
                  Criar categoria
                </Button>
              </CardContent>
            </Card>
          ) : (
            <div className="grid gap-4 md:grid-cols-3 lg:grid-cols-4">
              {section.items.map((cat) => (
                <Card key={cat.id} className="hover:shadow-md transition-shadow">
                  <CardContent className="flex items-center justify-between p-5">
                    <div className="flex items-center gap-3">
                      <div
                        className="flex h-8 w-8 items-center justify-center rounded-lg"
                        style={{ backgroundColor: cat.color + "20", color: cat.color }}
                      >
                        <Tags className="h-4 w-4" />
                      </div>
                      <div>
                        <p className="text-sm font-medium">{cat.name}</p>
                        {cat.is_system && (
                          <Badge variant="outline" className="text-xs">
                            Sistema
                          </Badge>
                        )}
                      </div>
                    </div>
                    {!cat.is_system && (
                      <div className="flex gap-1">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          onClick={() => openEdit(cat)}
                        >
                          <Pencil className="h-3 w-3" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 text-destructive hover:text-destructive"
                          onClick={() => handleDelete(cat.id)}
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              ))}
            </div>
          )}
        </div>
      ))}
    </div>
  );
}
