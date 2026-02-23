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
import { Button } from "@/components/ui/button";
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
  getUserSettings,
  updateUser,
  updateUserSettings,
} from "@/lib/auth";
import { useAuth } from "@/providers/auth-provider";
import type {
  ChangePasswordRequest,
  UpdateUserRequest,
  UpdateUserSettingsRequest,
  UserSettings,
} from "@/types";
import { useCallback, useEffect, useState } from "react";
import { toast } from "sonner";

export default function SettingsPage() {
  const { user, refreshUser, logout } = useAuth();
  const [settings, setSettings] = useState<UserSettings | null>(null);
  const [isLoading, setIsLoading] = useState(true);

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

  const fetchSettings = useCallback(async () => {
    try {
      const s = await getUserSettings();
      setSettings(s);
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
    fetchSettings();
  }, [user, fetchSettings]);

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

  async function handleSettingToggle(
    key: keyof UpdateUserSettingsRequest,
    value: boolean,
  ) {
    try {
      const body: UpdateUserSettingsRequest = { [key]: value };
      const updated = await updateUserSettings(body);
      setSettings(updated);
    } catch {
      toast.error("Erro ao atualizar configuração");
    }
  }

  async function handleThemeChange(theme: string) {
    try {
      const updated = await updateUserSettings({
        theme: theme as UpdateUserSettingsRequest["theme"],
      });
      setSettings(updated);
    } catch {
      toast.error("Erro ao atualizar tema");
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
      setPasswordForm({
        current_password: "",
        new_password: "",
        confirm_password: "",
      });
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

  if (isLoading) {
    return (
      <div className="space-y-6">
        <h1 className="text-3xl font-bold">Configurações</h1>
        <Skeleton className="h-[400px]" />
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
                      onChange={(e) =>
                        setProfileForm({
                          ...profileForm,
                          first_name: e.target.value,
                        })
                      }
                    />
                  </div>
                  <div className="space-y-2">
                    <Label>Sobrenome</Label>
                    <Input
                      value={profileForm.last_name}
                      onChange={(e) =>
                        setProfileForm({
                          ...profileForm,
                          last_name: e.target.value,
                        })
                      }
                    />
                  </div>
                </div>
                <div className="space-y-2">
                  <Label>Telefone</Label>
                  <Input
                    value={profileForm.phone}
                    onChange={(e) =>
                      setProfileForm({ ...profileForm, phone: e.target.value })
                    }
                  />
                </div>
                <div className="grid grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label>Moeda</Label>
                    <Select
                      value={profileForm.preferred_currency}
                      onValueChange={(v) =>
                        setProfileForm({
                          ...profileForm,
                          preferred_currency: v,
                        })
                      }
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="BRL">BRL - Real</SelectItem>
                        <SelectItem value="USD">USD - Dolar</SelectItem>
                        <SelectItem value="EUR">EUR - Euro</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label>Idioma</Label>
                    <Select
                      value={profileForm.preferred_language}
                      onValueChange={(v) =>
                        setProfileForm({
                          ...profileForm,
                          preferred_language: v,
                        })
                      }
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="pt-BR">Português</SelectItem>
                        <SelectItem value="en-US">English</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label>Timezone</Label>
                    <Select
                      value={profileForm.timezone}
                      onValueChange={(v) =>
                        setProfileForm({ ...profileForm, timezone: v })
                      }
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="America/Sao_Paulo">
                          São Paulo
                        </SelectItem>
                        <SelectItem value="America/New_York">
                          New York
                        </SelectItem>
                        <SelectItem value="Europe/London">London</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                <Button type="submit">Salvar Perfil</Button>
              </form>
            </CardContent>
          </Card>

          {settings && (
            <>
              <Card>
                <CardHeader>
                  <CardTitle>Transações</CardTitle>
                  <CardDescription>
                    Configure o comportamento das Transações
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm font-medium">
                        Permitir Transações manuais
                      </p>
                      <p className="text-xs text-muted-foreground">
                        Quando desativado, você só poderá ter Transações
                        automáticas do banco
                      </p>
                    </div>
                    <Switch
                      checked={settings.allow_manual_transactions}
                      onCheckedChange={(checked) =>
                        handleSettingToggle(
                          "allow_manual_transactions" as keyof UpdateUserSettingsRequest,
                          checked,
                        )
                      }
                    />
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle>Tema</CardTitle>
                </CardHeader>
                <CardContent>
                  <Select
                    value={settings.theme}
                    onValueChange={handleThemeChange}
                  >
                    <SelectTrigger className="w-[200px]">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="light">Claro</SelectItem>
                      <SelectItem value="dark">Escuro</SelectItem>
                      <SelectItem value="system">Sistema</SelectItem>
                    </SelectContent>
                  </Select>
                </CardContent>
              </Card>
            </>
          )}
        </TabsContent>

        <TabsContent value="notifications" className="mt-4">
          {settings && (
            <Card>
              <CardHeader>
                <CardTitle>Preferências de Notificação</CardTitle>
                <CardDescription>
                  Configure como deseja receber notificações
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                {[
                  {
                    key: "notification_email" as const,
                    label: "Email",
                    desc: "Receber notificações por email",
                  },
                  {
                    key: "notification_push" as const,
                    label: "Push",
                    desc: "Receber notificações push",
                  },
                  {
                    key: "notification_sms" as const,
                    label: "SMS",
                    desc: "Receber notificações por SMS",
                  },
                  {
                    key: "budget_alerts" as const,
                    label: "Alertas de Orçamento",
                    desc: "Quando atingir o limite do orçamento",
                  },
                  {
                    key: "bill_reminders" as const,
                    label: "Lembretes de Contas",
                    desc: "Quando uma conta estiver próxima do vencimento",
                  },
                  {
                    key: "weekly_summary" as const,
                    label: "Resumo Semanal",
                    desc: "Receber resumo semanal por email",
                  },
                  {
                    key: "monthly_report" as const,
                    label: "Relatório Mensal",
                    desc: "Receber relatorio mensal por email",
                  },
                  {
                    key: "low_balance_alert" as const,
                    label: "Alerta de Saldo Baixo",
                    desc: "Quando o saldo estiver abaixo do limite",
                  },
                ].map((item) => (
                  <div
                    key={item.key}
                    className="flex items-center justify-between"
                  >
                    <div>
                      <p className="text-sm font-medium">{item.label}</p>
                      <p className="text-xs text-muted-foreground">
                        {item.desc}
                      </p>
                    </div>
                    <Switch
                      checked={settings[item.key] as boolean}
                      onCheckedChange={(checked) =>
                        handleSettingToggle(item.key, checked)
                      }
                    />
                  </div>
                ))}
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="security" className="space-y-4 mt-4">
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
                    onChange={(e) =>
                      setPasswordForm({
                        ...passwordForm,
                        current_password: e.target.value,
                      })
                    }
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label>Nova Senha</Label>
                  <Input
                    type="password"
                    value={passwordForm.new_password}
                    onChange={(e) =>
                      setPasswordForm({
                        ...passwordForm,
                        new_password: e.target.value,
                      })
                    }
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label>Confirmar Nova Senha</Label>
                  <Input
                    type="password"
                    value={passwordForm.confirm_password}
                    onChange={(e) =>
                      setPasswordForm({
                        ...passwordForm,
                        confirm_password: e.target.value,
                      })
                    }
                    required
                  />
                </div>
                <Button type="submit">Alterar Senha</Button>
              </form>
            </CardContent>
          </Card>

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
                      Sua conta será desativada e você perderá acesso. Essa ação
                      pode ser irreversível.
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
