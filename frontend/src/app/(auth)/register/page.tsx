"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { z } from "zod";
import { useAuth } from "@/providers/auth-provider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { toast } from "sonner";
import { registerSchema } from "@/lib/schemas";

export default function RegisterPage() {
  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const { register } = useAuth();
  const router = useRouter();

  function updateField(field: string, value: string) {
    setForm((prev) => ({ ...prev, [field]: value }));
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    if (form.password !== form.confirmPassword) {
      toast.error("As senhas nao coincidem");
      return;
    }

    setIsLoading(true);
    try {
      const validated = registerSchema.parse({
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        password: form.password,
      });

      await register(validated);
      router.push("/dashboard");
    } catch (error) {
      if (error instanceof z.ZodError) {
        toast.error(error.issues[0].message);
      } else {
        toast.error("Erro ao criar conta. Tente novamente.");
      }
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="flex min-h-screen">
      {/* Left side - Branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary/20 via-primary/10 to-background items-center justify-center p-12 relative overflow-hidden">
        {/* Decorative background pattern */}
        <div className="absolute inset-0 opacity-5">
          <svg className="w-full h-full" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <pattern id="grid-register" width="40" height="40" patternUnits="userSpaceOnUse">
                <path d="M 40 0 L 0 0 0 40" fill="none" stroke="currentColor" strokeWidth="1"/>
              </pattern>
            </defs>
            <rect width="100%" height="100%" fill="url(#grid-register)" />
          </svg>
        </div>

        <div className="max-w-md space-y-8 relative z-10">
          <div className="space-y-2">
            <h1 className="text-7xl font-bold bg-gradient-to-r from-primary via-accent to-primary bg-clip-text text-transparent">
              FiNext
            </h1>
            <p className="text-sm text-primary/70 font-medium tracking-wide uppercase">
              Financial Next Generation
            </p>
          </div>
          <p className="text-xl text-foreground/80 leading-relaxed">
            Comece hoje a transformar sua relação com o dinheiro. Controle total, insights inteligentes e segurança garantida.
          </p>
          <div className="space-y-4 pt-4">
            <div className="flex items-start gap-3">
              <div className="w-2 h-2 bg-primary rounded-full mt-2"></div>
              <p className="text-foreground/70">Controle total das suas finanças em tempo real</p>
            </div>
            <div className="flex items-start gap-3">
              <div className="w-2 h-2 bg-primary rounded-full mt-2"></div>
              <p className="text-foreground/70">Metas e orçamentos personalizados automaticamente</p>
            </div>
            <div className="flex items-start gap-3">
              <div className="w-2 h-2 bg-primary rounded-full mt-2"></div>
              <p className="text-foreground/70">Relatórios e insights inteligentes</p>
            </div>
          </div>
        </div>
      </div>

      {/* Right side - Register Form */}
      <div className="flex w-full lg:w-1/2 items-center justify-center p-6">
        <div className="w-full max-w-md space-y-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold tracking-tight">Criar Conta</h2>
            <p className="mt-2 text-sm text-muted-foreground">
              Comece a gerenciar suas finanças agora
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="first_name">Nome</Label>
                <Input
                  id="first_name"
                  placeholder="Seu nome"
                  value={form.first_name}
                  onChange={(e) => updateField("first_name", e.target.value)}
                  required
                  className="h-11"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="last_name">Sobrenome</Label>
                <Input
                  id="last_name"
                  placeholder="Seu sobrenome"
                  value={form.last_name}
                  onChange={(e) => updateField("last_name", e.target.value)}
                  required
                  className="h-11"
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="seu@email.com"
                value={form.email}
                onChange={(e) => updateField("email", e.target.value)}
                required
                className="h-11"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Senha</Label>
              <Input
                id="password"
                type="password"
                placeholder="Minimo 8 caracteres"
                value={form.password}
                onChange={(e) => updateField("password", e.target.value)}
                required
                className="h-11"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">Confirmar Senha</Label>
              <Input
                id="confirmPassword"
                type="password"
                placeholder="Repita a senha"
                value={form.confirmPassword}
                onChange={(e) => updateField("confirmPassword", e.target.value)}
                required
                className="h-11"
              />
            </div>
            <Button type="submit" className="w-full h-11" disabled={isLoading}>
              {isLoading ? "Criando conta..." : "Criar Conta"}
            </Button>
          </form>

          <p className="text-center text-sm text-muted-foreground">
            Já tem conta?{" "}
            <Link href="/login" className="font-semibold text-primary hover:underline">
              Entrar
            </Link>
          </p>

          <p className="text-center text-xs text-muted-foreground">
            Ao criar sua conta, você concorda com nossos{" "}
            <Link href="#" className="underline underline-offset-4 hover:text-primary">
              Termos de Serviço
            </Link>{" "}
            e{" "}
            <Link href="#" className="underline underline-offset-4 hover:text-primary">
              Política de Privacidade
            </Link>
            .
          </p>
        </div>
      </div>
    </div>
  );
}
