"use client";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  changePassword,
  deactivateAccount,
  getLinkedProviders,
  getUserSettings,
  linkProvider,
  setPassword,
  unlinkProvider,
  updateUser,
  updateUserSettings,
} from "@/lib/auth";
import { auth, googleProvider } from "@/lib/firebase";
import { useAuth } from "@/providers/auth-provider";
import type {
  ChangePasswordRequest,
  LinkedProvider,
  ListProvidersResponse,
  UpdateUserRequest,
  UpdateUserSettingsRequest,
  UserSettings,
} from "@/types";
import { signInWithPopup } from "firebase/auth";
import { KeyRound, Link2, Link2Off, ShieldCheck } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

const BR_TIMEZONES = [
  { value: "America/Sao_Paulo", label: "Brasília / São Paulo (UTC-3)" },
  { value: "America/Bahia", label: "Bahia (UTC-3)" },
  { value: "America/Fortaleza", label: "Fortaleza (UTC-3)" },
  { value: "America/Recife", label: "Recife (UTC-3)" },
  { value: "America/Belem", label: "Belém (UTC-3)" },
  { value: "America/Maceio", label: "Maceió (UTC-3)" },
  { value: "America/Araguaina", label: "Araguaína (UTC-3)" },
  { value: "America/Santarem", label: "Santarém (UTC-3)" },
  { value: "America/Cuiaba", label: "Cuiabá (UTC-4)" },
  { value: "America/Campo_Grande", label: "Campo Grande (UTC-4)" },
  { value: "America/Manaus", label: "Manaus (UTC-4)" },
  { value: "America/Porto_Velho", label: "Porto Velho (UTC-4)" },
  { value: "America/Boa_Vista", label: "Boa Vista (UTC-4)" },
  { value: "America/Rio_Branco", label: "Rio Branco (UTC-5)" },
  { value: "America/Noronha", label: "Fernando de Noronha (UTC-2)" },
];

// Ícone SVG do Google
function GoogleIcon() {
  return (
    <svg className="h-5 w-5" viewBox="0 0 24 24">
      <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
      <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
      <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
      <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
    </svg>
  );
}

export default function SettingsPage() {
  const { user, refreshUser, logout } = useAuth();
  const [settings, setSettings] = useState<UserSettings | null>(null);
  const [providers, setProviders] = useState<ListProvidersResponse | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [linkingProvider, setLinkingProvider] = useState(false);

  const [profileForm, setProfileForm] = useState({
    first_name: "",
    last_name: "",
    phone: "",
    preferred_currency: "BRL",
    preferred_language: "pt-BR",
    timezone: "America/Sao_Paulo",
  });
  const [passwordForm, setPasswordForm] = useState({
    current_password: "",
    new_password: "",
    confirm_password: "",
  });
  const [firstPasswordForm, setFirstPasswordForm] = useState({
    new_password: "",
    confirm_password: "",
  });
  const [definePasswordOpen, setDefinePasswordOpen] = useState(false);

  const fetchData = useCallback(async () => {
    try {
      const [s, p] = await Promise.all([getUserSettings(), getLinkedProviders()]);
      setSettings(s);
      setProviders(p);
    } catch {
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    if (user) {
      setProfileForm({
        first_name: user.first_name,
        last_name: user.last_name,
        phone: user.phone || "",
        preferred_currency: user.preferred_currency,
        preferred_language: user.preferred_language,
        timezone: user.timezone,
      });
    }
    fetchData();
  }, [user, fetchData]);

  async function handleProfileSave(e: React.FormEvent) {
    e.preventDefault();
    try {
      const body: UpdateUserRequest = {
        first_name: profileForm.first_name,
        last_name: profileForm.last_name,
        phone: profileForm.phone || undefined,
        preferred_currency: profileForm.preferred_currency,
        preferred_language: profileForm.preferred_language,
        timezone: profileForm.timezone,
      };
      await updateUser(body);
      await refreshUser();
      toast.success("Perfil atualizado");
    } catch {
      toast.error("Erro ao atualizar perfil");
    }
  }

  async function handleSettingToggle(key: keyof UpdateUserSettingsRequest, value: boolean) {
    try {
      const updated = await updateUserSettings({ [key]: value });
      setSettings(updated);
    } catch {
      toast.error("Erro ao atualizar configuração");
    }
  }

  async function handlePasswordChange(e: React.FormEvent) {
    e.preventDefault();
    if (passwordForm.new_password !== passwordForm.confirm_password) {
      toast.error("Senhas não coincidem");
      return;
    }
    try {
      const body: ChangePasswordRequest = {
        current_password: passwordForm.current_password,
        new_password: passwordForm.new_password,
      };
      await changePassword(body);
      toast.success("Senha alterada com sucesso");
      setPasswordForm({ current_password: "", new_password: "", confirm_password: "" });
    } catch {
      toast.error("Erro ao alterar senha");
    }
  }

  async function handleDeactivate() {
    try {
      await deactivateAccount();
      logout();
    } catch {
      toast.error("Erro ao desativar conta");
    }
  }

  async function handleSetPassword(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (firstPasswordForm.new_password !== firstPasswordForm.confirm_password) {
      toast.error("As senhas não coincidem");
      return;
    }
    try {
      await setPassword({
        new_password: firstPasswordForm.new_password,
        confirm_password: firstPasswordForm.confirm_password,
      });
      toast.success("Senha definida com sucesso");
      setFirstPasswordForm({ new_password: "", confirm_password: "" });
      setDefinePasswordOpen(false);
      await fetchData();
    } catch {
      toast.error("Erro ao definir senha");
    }
  }

  async function handleLinkGoogle() {
    setLinkingProvider(true);
    try {
      const result = await signInWithPopup(auth, googleProvider);
      const idToken = await result.user.getIdToken();
      await linkProvider({ provider: "google", id_token: idToken });
      toast.success("Conta Google vinculada com sucesso");
      await fetchData();
    } catch (error) {
      console.error(error);
      toast.error("Erro ao vincular conta Google");
    } finally {
      setLinkingProvider(false);
    }
  }

  async function handleUnlinkGoogle() {
    setLinkingProvider(true);
    try {
      await unlinkProvider("google");
      toast.success("Conta Google desvinculada");
      await fetchData();
    } catch {
      toast.error("Erro ao desvincular conta Google");
    } finally {
      setLinkingProvider(false);
    }
  }

  const googleProvider_linked = providers?.providers.find(
    (p: LinkedProvider) => p.provider === "google"
  );

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold">Configurações</h1>
        <Skeleton className="h-100" />
      </div>
    );
  }

  return (
    <div className="space-y-8 pb-8">
      <div>
        <h1 className="text-4xl font-bold tracking-tight">Configurações</h1>
        <p className="text-muted-foreground mt-2">
          Gerencie suas preferências e configurações
        </p>
      </div>

      <Tabs defaultValue="profile">
        <TabsList>
          <TabsTrigger value="profile">Perfil</TabsTrigger>
          <TabsTrigger value="notifications">Notificações</TabsTrigger>
          <TabsTrigger value="security">Segurança</TabsTrigger>
        </TabsList>

        {/* ── PERFIL ── */}
        <TabsContent value="profile" className="space-y-4 mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Informações Pessoais</CardTitle>
              <CardDescription>Atualize seus dados</CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleProfileSave} className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label>Nome</Label>
                    <Input
                      value={profileForm.first_name}
                      onChange={(e) => setProfileForm({ ...profileForm, first_name: e.target.value })}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Sobrenome</Label>
                    <Input
                      value={profileForm.last_name}
                      onChange={(e) => setProfileForm({ ...profileForm, last_name: e.target.value })}
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label>Telefone</Label>
                  <Input
                    value={profileForm.phone}
                    onChange={(e) => setProfileForm({ ...profileForm, phone: e.target.value })}
                  />
                </div>
                <div className="grid grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label>Moeda</Label>
                    <Input value="BRL — Real Brasileiro" disabled className="bg-muted text-muted-foreground cursor-not-allowed" />
                  </div>
                  <div className="space-y-2">
                    <Label>Idioma</Label>
                    <Input value="Português (Brasil)" disabled className="bg-muted text-muted-foreground cursor-not-allowed" />
                  </div>
                  <div className="space-y-2">
                    <Label>Fuso Horário</Label>
                    <Select
                      value={profileForm.timezone}
                      onValueChange={(v) => setProfileForm({ ...profileForm, timezone: v })}
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {BR_TIMEZONES.map((tz) => (
                          <SelectItem key={tz.value} value={tz.value}>
                            {tz.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <Button type="submit">Salvar Perfil</Button>
              </form>
            </CardContent>
          </Card>

          {settings && (
            <Card>
              <CardHeader>
                <CardTitle>Transações</CardTitle>
                <CardDescription>Configure o comportamento das Transações</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium">Permitir Transações manuais</p>
                    <p className="text-xs text-muted-foreground">
                      Quando desativado, você só poderá ter Transações automáticas do banco
                    </p>
                  </div>
                  <Switch
                    checked={settings.allow_manual_transactions}
                    onCheckedChange={(checked) =>
                      handleSettingToggle("allow_manual_transactions" as keyof UpdateUserSettingsRequest, checked)
                    }
                  />
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* ── NOTIFICAÇÕES ── */}
        <TabsContent value="notifications" className="mt-4">
          {settings && (
            <Card>
              <CardHeader>
                <CardTitle>Preferências de Notificação</CardTitle>
                <CardDescription>Configure como deseja receber notificações</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {[
                  { key: "notification_email" as const, label: "Email", desc: "Receber notificações por email" },
                  { key: "notification_push" as const, label: "Push", desc: "Receber notificações push" },
                  { key: "notification_sms" as const, label: "SMS", desc: "Receber notificações por SMS" },
                  { key: "budget_alerts" as const, label: "Alertas de Orçamento", desc: "Quando atingir o limite do orçamento" },
                  { key: "bill_reminders" as const, label: "Lembretes de Contas", desc: "Quando uma conta estiver próxima do vencimento" },
                  { key: "weekly_summary" as const, label: "Resumo Semanal", desc: "Receber resumo semanal por email" },
                  { key: "monthly_report" as const, label: "Relatório Mensal", desc: "Receber relatório mensal por email" },
                  { key: "low_balance_alert" as const, label: "Alerta de Saldo Baixo", desc: "Quando o saldo estiver abaixo do limite" },
                ].map((item) => (
                  <div key={item.key} className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium">{item.label}</p>
                      <p className="text-xs text-muted-foreground">{item.desc}</p>
                    </div>
                    <Switch
                      checked={settings[item.key] as boolean}
                      onCheckedChange={(checked) => handleSettingToggle(item.key, checked)}
                    />
                  </div>
                ))}
              </CardContent>
            </Card>
          )}
        </TabsContent>

        {/* ── SEGURANÇA ── */}
        <TabsContent value="security" className="space-y-4 mt-4">

          {/* Métodos de Login */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <ShieldCheck className="h-5 w-5" />
                Métodos de Login
              </CardTitle>
              <CardDescription>
                Gerencie como você acessa sua conta
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Senha */}
              <div className="flex items-center justify-between py-2">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-full bg-muted">
                    <KeyRound className="h-4 w-4 text-muted-foreground" />
                  </div>
                  <div>
                    <p className="text-sm font-medium">Acesso por Senha</p>
                    <p className="text-xs text-muted-foreground">
                      {providers?.has_password
                        ? "Sua conta está protegida por senha"
                        : "Defina uma senha para ter outro método de acesso"}
                    </p>
                  </div>
                </div>
                {providers?.has_password ? (
                  <Badge variant="default">Ativo</Badge>
                ) : (
                  <Dialog open={definePasswordOpen} onOpenChange={setDefinePasswordOpen}>
                    <DialogTrigger asChild>
                      <Button variant="outline" size="sm">
                        <KeyRound className="h-4 w-4 mr-1" />
                        Definir Senha
                      </Button>
                    </DialogTrigger>
                    <DialogContent className="sm:max-w-md">
                      <DialogHeader>
                        <DialogTitle>Definir Senha</DialogTitle>
                        <DialogDescription>
                          Crie uma senha para acessar sua conta também por email e senha
                        </DialogDescription>
                      </DialogHeader>
                      <form onSubmit={handleSetPassword} className="space-y-4 pt-2">
                        <div className="space-y-2">
                          <Label>Nova Senha</Label>
                          <Input
                            type="password"
                            placeholder="Mínimo 8 caracteres"
                            value={firstPasswordForm.new_password}
                            onChange={(e) => setFirstPasswordForm({ ...firstPasswordForm, new_password: e.target.value })}
                            required
                          />
                        </div>
                        <div className="space-y-2">
                          <Label>Confirmar Senha</Label>
                          <Input
                            type="password"
                            placeholder="Repita a senha"
                            value={firstPasswordForm.confirm_password}
                            onChange={(e) => setFirstPasswordForm({ ...firstPasswordForm, confirm_password: e.target.value })}
                            required
                          />
                        </div>
                        <div className="flex justify-end gap-2 pt-2">
                          <Button type="button" variant="ghost" onClick={() => setDefinePasswordOpen(false)}>
                            Cancelar
                          </Button>
                          <Button type="submit">Definir Senha</Button>
                        </div>
                      </form>
                    </DialogContent>
                  </Dialog>
                )}
              </div>

              <Separator />

              {/* Google */}
              <div className="flex items-center justify-between py-2">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 items-center justify-center rounded-full bg-muted">
                    <GoogleIcon />
                  </div>
                  <div>
                    <p className="text-sm font-medium">Google</p>
                    {googleProvider_linked ? (
                      <p className="text-xs text-muted-foreground">{googleProvider_linked.email}</p>
                    ) : (
                      <p className="text-xs text-muted-foreground">Não vinculado</p>
                    )}
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  {googleProvider_linked ? (
                    <>
                      <Badge variant="default">Vinculado</Badge>
                      {providers?.has_password && (
                        <AlertDialog>
                          <AlertDialogTrigger asChild>
                            <Button variant="ghost" size="sm" disabled={linkingProvider}>
                              <Link2Off className="h-4 w-4 mr-1" />
                              Desvincular
                            </Button>
                          </AlertDialogTrigger>
                          <AlertDialogContent>
                            <AlertDialogHeader>
                              <AlertDialogTitle>Desvincular Google?</AlertDialogTitle>
                              <AlertDialogDescription>
                                Você não poderá mais fazer login com sua conta Google. Certifique-se de ter uma senha cadastrada.
                              </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                              <AlertDialogCancel>Cancelar</AlertDialogCancel>
                              <AlertDialogAction onClick={handleUnlinkGoogle}>
                                Sim, desvincular
                              </AlertDialogAction>
                            </AlertDialogFooter>
                          </AlertDialogContent>
                        </AlertDialog>
                      )}
                    </>
                  ) : (
                    <Button variant="outline" size="sm" onClick={handleLinkGoogle} disabled={linkingProvider}>
                      <Link2 className="h-4 w-4 mr-1" />
                      Vincular Google
                    </Button>
                  )}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Alterar Senha */}
          {providers?.has_password && (
            <Card>
              <CardHeader>
                <CardTitle>Alterar Senha</CardTitle>
              </CardHeader>
              <CardContent>
                <form onSubmit={handlePasswordChange} className="space-y-4">
                  <div className="space-y-2">
                    <Label>Senha Atual</Label>
                    <Input
                      type="password"
                      value={passwordForm.current_password}
                      onChange={(e) => setPasswordForm({ ...passwordForm, current_password: e.target.value })}
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Nova Senha</Label>
                    <Input
                      type="password"
                      value={passwordForm.new_password}
                      onChange={(e) => setPasswordForm({ ...passwordForm, new_password: e.target.value })}
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Confirmar Nova Senha</Label>
                    <Input
                      type="password"
                      value={passwordForm.confirm_password}
                      onChange={(e) => setPasswordForm({ ...passwordForm, confirm_password: e.target.value })}
                      required
                    />
                  </div>
                  <Button type="submit">Alterar Senha</Button>
                </form>
              </CardContent>
            </Card>
          )}

          <Separator />

          <Card className="border-destructive">
            <CardHeader>
              <CardTitle className="text-destructive">Zona de Perigo</CardTitle>
              <CardDescription>Ações irreversíveis</CardDescription>
            </CardHeader>
            <CardContent>
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive">Desativar Conta</Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Desativar sua conta?</AlertDialogTitle>
                    <AlertDialogDescription>
                      Sua conta será desativada e você perderá acesso. Essa ação pode ser irreversível.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancelar</AlertDialogCancel>
                    <AlertDialogAction
                      onClick={handleDeactivate}
                      className="bg-destructive text-destructive-foreground"
                    >
                      Sim, desativar
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
