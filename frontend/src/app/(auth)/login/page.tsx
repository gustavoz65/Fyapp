"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { z } from "zod";
import { useAuth } from "@/providers/auth-provider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { toast } from "sonner";
import { loginSchema } from "@/lib/schemas";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const router = useRouter();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setIsLoading(true);

    try {
      // Validar com Zod
      const validated = loginSchema.parse({ email, password });
      await login(validated);
      router.push("/dashboard");
    } catch (error) {
      if (error instanceof z.ZodError) {
        toast.error(error.errors[0].message);
      } else {
        toast.error("Email ou senha invalidos");
      }
    } finally {
      setIsLoading(false);
    }
  }

  function handleGoogleLogin() {
    toast.info("Integracao com Google em breve!");
    // TODO: Implementar Firebase Google Auth
  }

  return (
    <div className="flex min-h-screen">
      {/* Left side - Branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary/10 via-primary/5 to-background items-center justify-center p-12">
        <div className="max-w-md space-y-6">
          <h1 className="text-6xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
            Cashing
          </h1>
          <p className="text-xl text-muted-foreground">
            Gerencie suas financas pessoais de forma inteligente e eficiente
          </p>
          <blockquote className="border-l-4 border-primary pl-4 italic text-muted-foreground">
            &quot;Esta plataforma transformou completamente a forma como gerencio meu dinheiro.&quot;
            <footer className="mt-2 text-sm font-semibold">- Usuario Satisfeito</footer>
          </blockquote>
        </div>
      </div>

      {/* Right side - Login Form */}
      <div className="flex w-full lg:w-1/2 items-center justify-center p-6">
        <div className="w-full max-w-md space-y-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold tracking-tight">Entrar na sua conta</h2>
            <p className="mt-2 text-sm text-muted-foreground">
              Digite seu email abaixo para entrar
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="space-y-2">
              <Input
                type="email"
                placeholder="nome@exemplo.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="h-11"
              />
            </div>
            <div className="space-y-2">
              <Input
                type="password"
                placeholder="Sua senha"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="h-11"
              />
            </div>
            <Button type="submit" className="w-full h-11" disabled={isLoading}>
              {isLoading ? "Entrando..." : "Entrar com Email"}
            </Button>
          </form>

          <div className="relative">
            <div className="absolute inset-0 flex items-center">
              <Separator />
            </div>
            <div className="relative flex justify-center text-xs uppercase">
              <span className="bg-background px-2 text-muted-foreground">
                Ou continue com
              </span>
            </div>
          </div>

          <Button
            type="button"
            variant="outline"
            className="w-full h-11"
            onClick={handleGoogleLogin}
          >
            <svg className="mr-2 h-4 w-4" viewBox="0 0 24 24">
              <path
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                fill="#4285F4"
              />
              <path
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                fill="#34A853"
              />
              <path
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                fill="#FBBC05"
              />
              <path
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                fill="#EA4335"
              />
            </svg>
            Google
          </Button>

          <p className="text-center text-sm text-muted-foreground">
            Nao tem uma conta?{" "}
            <Link href="/register" className="font-semibold text-primary hover:underline">
              Criar conta
            </Link>
          </p>

          <p className="text-center text-xs text-muted-foreground">
            Ao continuar, voce concorda com nossos{" "}
            <Link href="#" className="underline underline-offset-4 hover:text-primary">
              Termos de Servico
            </Link>{" "}
            e{" "}
            <Link href="#" className="underline underline-offset-4 hover:text-primary">
              Politica de Privacidade
            </Link>
            .
          </p>
        </div>
      </div>
    </div>
  );
}
