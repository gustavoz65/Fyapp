"use client";

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
import { api } from "@/lib/api";
import { useAuth } from "@/providers/auth-provider";
import {
  ArrowRight,
  Building2,
  CheckCircle2,
  CreditCard,
  PenLine,
  PiggyBank,
  Sparkles,
  Wallet,
} from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "sonner";

const ACCOUNT_TYPES = [
  { value: "checking", label: "Conta Corrente", icon: CreditCard },
  { value: "savings", label: "Poupança", icon: PiggyBank },
  { value: "cash", label: "Dinheiro em Espécie", icon: Wallet },
  { value: "other", label: "Outro", icon: Building2 },
];

const LEFT_PANEL_CONTENT = {
  1: {
    headline: "Sua jornada financeira começa agora",
    bullets: [
      { icon: Wallet, text: "Controle total das suas finanças" },
      { icon: Sparkles, text: "Insights inteligentes sobre seus gastos" },
      { icon: CheckCircle2, text: "Metas e orçamentos personalizados" },
    ],
  },
  2: {
    headline: "Vamos configurar tudo para você",
    bullets: [
      { icon: CreditCard, text: "Adicione quantas contas quiser" },
      { icon: Sparkles, text: "Categorização automática de transações" },
      { icon: PiggyBank, text: "Acompanhe seu saldo em tempo real" },
    ],
  },
  3: {
    headline: "Você está pronto para começar!",
    bullets: [
      { icon: CheckCircle2, text: "Conta criada e configurada" },
      { icon: Sparkles, text: "Adicione suas primeiras transações" },
      { icon: Wallet, text: "Comece a alcançar suas metas" },
    ],
  },
};

export default function OnboardingPage() {
  const { user, isAuthenticated, isLoading: authLoading } = useAuth();
  const router = useRouter();

  const [step, setStep] = useState<1 | 2 | 3>(1);
  const [accountForm, setAccountForm] = useState({
    name: "",
    account_type: "checking",
    initial_balance: "",
  });
  const [isLoading, setIsLoading] = useState(false);
  const [skipped, setSkipped] = useState(false);
  const [createdAccount, setCreatedAccount] = useState<{
    name: string;
    account_type: string;
    initial_balance: string;
  } | null>(null);

  useEffect(() => {
    if (authLoading) return;
    if (!isAuthenticated) {
      router.push("/login");
      return;
    }
    // Verifica via API se o usuário já tem contas — mais confiável que localStorage
    api.get<{ data?: unknown[] } | unknown[]>("/accounts").then((res) => {
      const accounts = Array.isArray(res) ? res : (res as { data?: unknown[] }).data ?? [];
      if (accounts.length > 0) {
        router.push("/dashboard");
      }
    }).catch(() => {
      // Se falhar a verificação, deixa no onboarding
    });
  }, [isAuthenticated, authLoading, router]);

  function completeOnboarding() {
    router.push("/dashboard");
  }

  async function handleCreateAccount(e: React.FormEvent) {
    e.preventDefault();

    if (!accountForm.name.trim()) {
      toast.error("Digite um nome para a conta");
      return;
    }

    setIsLoading(true);
    try {
      const rawBalance = accountForm.initial_balance.replace(",", ".");
      const balance = parseFloat(rawBalance) || 0;

      await api.post("/accounts", {
        name: accountForm.name.trim(),
        account_type: accountForm.account_type,
        initial_balance: balance,
      });

      setCreatedAccount({ ...accountForm });
      setSkipped(false);
      setStep(3);
    } catch {
      toast.error("Erro ao criar conta. Tente novamente.");
    } finally {
      setIsLoading(false);
    }
  }

  function handleSkip() {
    setSkipped(true);
    setCreatedAccount(null);
    setStep(3);
  }

  if (authLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    );
  }

  const panelContent = LEFT_PANEL_CONTENT[step];

  const accountTypeLabel =
    ACCOUNT_TYPES.find((t) => t.value === createdAccount?.account_type)
      ?.label ?? createdAccount?.account_type;

  return (
    <div className="flex min-h-screen">
      {/* Left panel — branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary/20 via-primary/10 to-background items-center justify-center p-12 relative overflow-hidden">
        <div className="absolute inset-0 opacity-5">
          <svg className="w-full h-full" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <pattern
                id="grid-onboarding"
                width="40"
                height="40"
                patternUnits="userSpaceOnUse"
              >
                <path
                  d="M 40 0 L 0 0 0 40"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="1"
                />
              </pattern>
            </defs>
            <rect width="100%" height="100%" fill="url(#grid-onboarding)" />
          </svg>
        </div>

        <div className="max-w-md space-y-8 relative z-10 transition-all duration-300">
          <div className="space-y-2">
            <h1 className="text-7xl font-bold bg-gradient-to-r from-primary via-accent to-primary bg-clip-text text-transparent">
              FiNext
            </h1>
            <p className="text-sm text-primary/70 font-medium tracking-wide uppercase">
              Financial Next Generation
            </p>
          </div>

          <p className="text-xl text-foreground/80 leading-relaxed">
            {panelContent.headline}
          </p>

          <div className="space-y-4 pt-4">
            {panelContent.bullets.map(({ icon: Icon, text }, i) => (
              <div key={i} className="flex items-start gap-3">
                <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center mt-0.5">
                  <Icon className="w-4 h-4 text-primary" />
                </div>
                <p className="text-foreground/70 pt-1">{text}</p>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Right panel — wizard */}
      <div className="flex w-full lg:w-1/2 items-center justify-center p-6">
        <div className="w-full max-w-md space-y-8">
          {/* Progress indicator */}
          <div className="space-y-2">
            <div className="flex items-center justify-between text-sm text-muted-foreground">
              <span>Configuração inicial</span>
              <span>Etapa {step} de 3</span>
            </div>
            <div className="w-full h-1.5 bg-muted rounded-full overflow-hidden">
              <div
                className="h-full bg-primary rounded-full transition-all duration-500 ease-out"
                style={{ width: `${(step / 3) * 100}%` }}
              />
            </div>
          </div>

          {/* Step 1 — Welcome */}
          {step === 1 && (
            <div className="space-y-6">
              <div className="space-y-2">
                <h2 className="text-3xl font-bold tracking-tight">
                  Bem-vindo, {user?.first_name}!
                </h2>
                <p className="text-muted-foreground">
                  Como você prefere adicionar suas finanças?
                </p>
              </div>

              <div className="space-y-4">
                {/* Card — Open Finance (disabled) */}
                <div className="relative border-2 border-muted rounded-xl p-5 opacity-60 cursor-not-allowed select-none">
                  <Badge className="absolute top-3 right-3 bg-amber-500/10 text-amber-600 border-amber-500/20 hover:bg-amber-500/10">
                    Em breve
                  </Badge>
                  <div className="flex items-start gap-4">
                    <div className="flex-shrink-0 w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center">
                      <Building2 className="w-6 h-6 text-primary" />
                    </div>
                    <div className="space-y-1 flex-1">
                      <h3 className="font-semibold text-base">
                        Conectar meu banco
                      </h3>
                      <p className="text-sm text-muted-foreground">
                        Sincronize automaticamente todas as transações
                      </p>
                      <ul className="mt-2 space-y-1">
                        {[
                          "Tudo sincronizado automaticamente",
                          "Economize 10 min/dia",
                          "7 dias grátis de Pro",
                        ].map((item) => (
                          <li
                            key={item}
                            className="text-xs text-muted-foreground flex items-center gap-1.5"
                          >
                            <span className="text-primary">✓</span>
                            {item}
                          </li>
                        ))}
                      </ul>
                    </div>
                  </div>
                </div>

                {/* Card — Manual (active) */}
                <button
                  type="button"
                  onClick={() => setStep(2)}
                  className="w-full text-left border-2 border-primary/40 rounded-xl p-5 hover:border-primary hover:shadow-md transition-all duration-200 cursor-pointer group"
                >
                  <div className="flex items-start gap-4">
                    <div className="flex-shrink-0 w-12 h-12 rounded-xl bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                      <PenLine className="w-6 h-6 text-primary" />
                    </div>
                    <div className="space-y-1 flex-1">
                      <h3 className="font-semibold text-base">
                        Adicionar manualmente
                      </h3>
                      <p className="text-sm text-muted-foreground">
                        Controle total, comece em 2 minutos
                      </p>
                      <ul className="mt-2 space-y-1">
                        {[
                          "Controle total das suas finanças",
                          "Rápido e simples de usar",
                          "Sempre gratuito",
                        ].map((item) => (
                          <li
                            key={item}
                            className="text-xs text-muted-foreground flex items-center gap-1.5"
                          >
                            <span className="text-primary">✓</span>
                            {item}
                          </li>
                        ))}
                      </ul>
                    </div>
                    <ArrowRight className="flex-shrink-0 w-5 h-5 text-muted-foreground group-hover:text-primary group-hover:translate-x-0.5 transition-all mt-3" />
                  </div>
                </button>
              </div>
            </div>
          )}

          {/* Step 2 — Create account */}
          {step === 2 && (
            <div className="space-y-6">
              <div className="space-y-2">
                <h2 className="text-3xl font-bold tracking-tight">
                  Sua primeira conta
                </h2>
                <p className="text-muted-foreground">
                  Não precisa ser perfeito — você pode editar depois.
                </p>
              </div>

              <form onSubmit={handleCreateAccount} className="space-y-5">
                <div className="space-y-2">
                  <Label htmlFor="account-name">Nome da conta</Label>
                  <Input
                    id="account-name"
                    placeholder="Ex: Nubank, Carteira, Bradesco"
                    value={accountForm.name}
                    onChange={(e) =>
                      setAccountForm((prev) => ({
                        ...prev,
                        name: e.target.value,
                      }))
                    }
                    className="h-11"
                    autoFocus
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="account-type">Tipo de conta</Label>
                  <Select
                    value={accountForm.account_type}
                    onValueChange={(value) =>
                      setAccountForm((prev) => ({
                        ...prev,
                        account_type: value,
                      }))
                    }
                  >
                    <SelectTrigger id="account-type" className="h-11">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {ACCOUNT_TYPES.map(({ value, label, icon: Icon }) => (
                        <SelectItem key={value} value={value}>
                          <div className="flex items-center gap-2">
                            <Icon className="w-4 h-4 text-muted-foreground" />
                            {label}
                          </div>
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="initial-balance">
                    Saldo inicial{" "}
                    <span className="text-muted-foreground font-normal">
                      (opcional)
                    </span>
                  </Label>
                  <div className="relative">
                    <span className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground text-sm font-medium">
                      R$
                    </span>
                    <Input
                      id="initial-balance"
                      placeholder="0,00"
                      value={accountForm.initial_balance}
                      onChange={(e) =>
                        setAccountForm((prev) => ({
                          ...prev,
                          initial_balance: e.target.value,
                        }))
                      }
                      className="h-11 pl-9"
                    />
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Quanto você tem nessa conta agora? (pode deixar em branco)
                  </p>
                </div>

                <Button
                  type="submit"
                  className="w-full h-11"
                  disabled={isLoading}
                >
                  {isLoading ? (
                    <span className="flex items-center gap-2">
                      <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                      Criando conta...
                    </span>
                  ) : (
                    <span className="flex items-center gap-2">
                      Criar conta
                      <ArrowRight className="w-4 h-4" />
                    </span>
                  )}
                </Button>
              </form>

              <div className="text-center">
                <button
                  type="button"
                  onClick={handleSkip}
                  className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                >
                  Pular por agora →
                </button>
              </div>
            </div>
          )}

          {/* Step 3 — Success */}
          {step === 3 && (
            <div className="space-y-6 text-center">
              <div className="flex justify-center">
                <div className="w-20 h-20 rounded-full bg-green-500/10 flex items-center justify-center animate-bounce">
                  <CheckCircle2 className="w-10 h-10 text-green-500" />
                </div>
              </div>

              <div className="space-y-2">
                <h2 className="text-3xl font-bold tracking-tight">
                  {skipped ? "Tudo configurado!" : "Conta criada com sucesso!"}
                </h2>
                <p className="text-muted-foreground">
                  Agora você pode adicionar suas transações, criar metas e
                  acompanhar suas finanças em tempo real.
                </p>
              </div>

              {createdAccount && !skipped && (
                <Card className="text-left border-primary/20 bg-primary/5">
                  <CardContent className="pt-4 pb-4 space-y-2">
                    <p className="text-xs text-muted-foreground uppercase tracking-wide font-medium">
                      Conta criada
                    </p>
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-semibold">{createdAccount.name}</p>
                        <p className="text-sm text-muted-foreground">
                          {accountTypeLabel}
                        </p>
                      </div>
                      {createdAccount.initial_balance &&
                        parseFloat(
                          createdAccount.initial_balance.replace(",", "."),
                        ) > 0 && (
                          <p className="font-semibold text-green-600">
                            R${" "}
                            {parseFloat(
                              createdAccount.initial_balance.replace(",", "."),
                            ).toLocaleString("pt-BR", {
                              minimumFractionDigits: 2,
                            })}
                          </p>
                        )}
                    </div>
                  </CardContent>
                </Card>
              )}

              <Button
                onClick={completeOnboarding}
                className="w-full h-11"
                size="lg"
              >
                <span className="flex items-center gap-2">
                  Ir para o Dashboard
                  <ArrowRight className="w-4 h-4" />
                </span>
              </Button>

              <p className="text-xs text-muted-foreground">
                Você pode adicionar mais contas a qualquer momento nas
                configurações.
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
