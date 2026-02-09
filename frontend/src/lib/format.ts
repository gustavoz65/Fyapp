import { format, formatDistanceToNow, parseISO } from "date-fns";
import { ptBR } from "date-fns/locale";

export function formatCurrency(value: string | number, currency = "BRL"): string {
  const num = typeof value === "string" ? parseFloat(value) : value;
  return new Intl.NumberFormat("pt-BR", {
    style: "currency",
    currency,
  }).format(num);
}

export function formatDate(dateStr: string): string {
  return format(parseISO(dateStr), "dd/MM/yyyy", { locale: ptBR });
}

export function formatDateTime(dateStr: string): string {
  return format(parseISO(dateStr), "dd/MM/yyyy HH:mm", { locale: ptBR });
}

export function formatRelativeDate(dateStr: string): string {
  return formatDistanceToNow(parseISO(dateStr), { addSuffix: true, locale: ptBR });
}

export function formatPercentage(value: string | number): string {
  const num = typeof value === "string" ? parseFloat(value) : value;
  return `${num.toFixed(1)}%`;
}

export function getAccountTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    checking: "Conta Corrente",
    savings: "Poupanca",
    credit_card: "Cartao de Credito",
    investment: "Investimento",
    cash: "Dinheiro",
    other: "Outro",
  };
  return labels[type] || type;
}

export function getTransactionTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    income: "Receita",
    expense: "Despesa",
    transfer: "Transferencia",
  };
  return labels[type] || type;
}

export function getBudgetPeriodLabel(type: string): string {
  const labels: Record<string, string> = {
    monthly: "Mensal",
    quarterly: "Trimestral",
    yearly: "Anual",
    custom: "Personalizado",
  };
  return labels[type] || type;
}

export function getGoalStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    in_progress: "Em Andamento",
    completed: "Concluida",
    cancelled: "Cancelada",
  };
  return labels[status] || status;
}

export function getNotificationTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    budget_alert: "Alerta de Orcamento",
    bill_reminder: "Lembrete de Conta",
    goal_achieved: "Meta Alcancada",
    low_balance: "Saldo Baixo",
    transaction_alert: "Alerta de Transacao",
    system: "Sistema",
  };
  return labels[type] || type;
}
