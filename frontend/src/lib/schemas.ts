import { z } from "zod";


export const loginSchema = z.object({
  email: z
    .string()
    .min(1, "Email é obrigatório")
    .email("Email inválido")
    .max(255, "Email muito longo"),
  password: z
    .string()
    .min(1, "Senha é obrigatória")
    .min(8, "Senha deve ter no mínimo 8 caracteres")
    .max(72, "Senha muito longa"),
});

export const registerSchema = z.object({
  first_name: z
    .string()
    .min(2, "Nome deve ter no mínimo 2 caracteres")
    .max(100, "Nome muito longo")
    .regex(/^[a-zA-ZÀ-ÿ\s'-]+$/, "Nome contém caracteres inválidos"),
  last_name: z
    .string()
    .min(2, "Sobrenome deve ter no mínimo 2 caracteres")
    .max(100, "Sobrenome muito longo")
    .regex(/^[a-zA-ZÀ-ÿ\s'-]+$/, "Sobrenome contém caracteres inválidos"),
  email: z
    .string()
    .min(1, "Email é obrigatório")
    .email("Email inválido")
    .max(255, "Email muito longo"),
  password: z
    .string()
    .min(8, "Senha deve ter no mínimo 8 caracteres")
    .max(72, "Senha muito longa")
    .regex(/[A-Z]/, "Senha deve conter ao menos uma letra maiúscula")
    .regex(/[a-z]/, "Senha deve conter ao menos uma letra minúscula")
    .regex(/[0-9]/, "Senha deve conter ao menos um número"),
  phone: z.string().max(20, "Telefone muito longo").optional(),
});


export const createTransactionSchema = z.object({
  bank_account_id: z.string().uuid("ID de conta inválido"),
  category_id: z.string().uuid("ID de categoria inválido").optional().or(z.literal("")),
  type: z.enum(["income", "expense"], {
    errorMap: () => ({ message: "Tipo deve ser 'Receita' ou 'Despesa'" }),
  }),
  amount: z
    .string()
    .refine((val) => !isNaN(parseFloat(val)), "Valor inválido")
    .refine((val) => parseFloat(val) > 0, "Valor deve ser maior que zero")
    .refine((val) => parseFloat(val) <= 999999999.99, "Valor muito alto"),
  description: z
    .string()
    .min(2, "Descrição deve ter no mínimo 2 caracteres")
    .max(255, "Descrição muito longa")
    .regex(/^[^<>]*$/, "Descrição contém caracteres proibidos (< ou >)"),
  notes: z
    .string()
    .max(1000, "Notas muito longas")
    .optional()
    .or(z.literal("")),
  transaction_date: z
    .string()
    .refine((val) => !isNaN(Date.parse(val)), "Data inválida"),
  due_date: z
    .string()
    .refine((val) => !isNaN(Date.parse(val)), "Data de vencimento inválida")
    .optional()
    .or(z.literal("")),
  is_paid: z.boolean(),
  tags: z
    .array(z.string().max(50, "Tag muito longa"))
    .max(10, "Máximo 10 tags permitidas")
    .optional(),
  total_installments: z
    .number()
    .int("Parcelas deve ser número inteiro")
    .min(2, "Mínimo 2 parcelas")
    .max(120, "Máximo 120 parcelas")
    .optional(),
});


export const createAccountSchema = z.object({
  name: z
    .string()
    .min(2, "Nome deve ter no mínimo 2 caracteres")
    .max(100, "Nome muito longo"),
  bank_name: z.string().max(100, "Nome do banco muito longo").optional(),
  bank_code: z.string().max(10, "Código do banco muito longo").optional(),
  account_type: z.enum(
    ["checking", "savings", "credit_card", "investment", "cash", "other"],
    { errorMap: () => ({ message: "Tipo de conta inválido" }) }
  ),
  account_number: z.string().max(50, "Número da conta muito longo").optional(),
  agency: z.string().max(20, "Agência muito longa").optional(),
  initial_balance: z
    .string()
    .refine((val) => !isNaN(parseFloat(val)), "Saldo inicial inválido"),
  credit_limit: z
    .string()
    .refine((val) => !isNaN(parseFloat(val)), "Limite de crédito inválido")
    .optional()
    .or(z.literal("")),
  closing_day: z
    .number()
    .int()
    .min(1, "Dia de fechamento deve ser entre 1 e 31")
    .max(31, "Dia de fechamento deve ser entre 1 e 31")
    .optional(),
  due_day: z
    .number()
    .int()
    .min(1, "Dia de vencimento deve ser entre 1 e 31")
    .max(31, "Dia de vencimento deve ser entre 1 e 31")
    .optional(),
  currency: z.string().length(3, "Moeda deve ter 3 caracteres").optional(),
  color: z
    .string()
    .regex(/^#[0-9A-Fa-f]{6}$/, "Cor deve ser hexadecimal válida (ex: #FF0000)")
    .optional(),
  icon: z.string().max(50, "Ícone muito longo").optional(),
  include_in_total: z.boolean().optional(),
});


export const createBudgetSchema = z.object({
  category_id: z.string().uuid("ID de categoria inválido").optional().or(z.literal("")),
  name: z
    .string()
    .min(2, "Nome deve ter no mínimo 2 caracteres")
    .max(100, "Nome muito longo"),
  amount: z
    .string()
    .refine((val) => !isNaN(parseFloat(val)), "Valor inválido")
    .refine((val) => parseFloat(val) > 0, "Valor deve ser maior que zero"),
  period_type: z.enum(["monthly", "quarterly", "yearly", "custom"], {
    errorMap: () => ({ message: "Período inválido" }),
  }),
  start_date: z.string().refine((val) => !isNaN(Date.parse(val)), "Data inicial inválida"),
  end_date: z.string().refine((val) => !isNaN(Date.parse(val)), "Data final inválida"),
  alert_threshold: z
    .string()
    .refine(
      (val) =>
        !isNaN(parseFloat(val)) &&
        parseFloat(val) > 0 &&
        parseFloat(val) <= 100,
      "Limite de alerta deve estar entre 0 e 100"
    )
    .optional()
    .or(z.literal("")),
});


export const createGoalSchema = z.object({
  name: z
    .string()
    .min(2, "Nome deve ter no mínimo 2 caracteres")
    .max(100, "Nome muito longo"),
  description: z
    .string()
    .max(500, "Descrição muito longa")
    .optional()
    .or(z.literal("")),
  target_amount: z
    .string()
    .refine((val) => !isNaN(parseFloat(val)), "Valor da meta inválido")
    .refine((val) => parseFloat(val) > 0, "Valor da meta deve ser maior que zero"),
  target_date: z
    .string()
    .refine((val) => !isNaN(Date.parse(val)), "Data da meta inválida")
    .optional()
    .or(z.literal("")),
  icon: z.string().max(50, "Ícone muito longo").optional(),
  color: z
    .string()
    .regex(/^#[0-9A-Fa-f]{6}$/, "Cor deve ser hexadecimal válida")
    .optional(),
  priority: z
    .number()
    .int()
    .min(1, "Prioridade deve ser entre 1 e 5")
    .max(5, "Prioridade deve ser entre 1 e 5")
    .optional(),
});


export type LoginInput = z.infer<typeof loginSchema>;
export type RegisterInput = z.infer<typeof registerSchema>;
export type CreateTransactionInput = z.infer<typeof createTransactionSchema>;
export type CreateAccountInput = z.infer<typeof createAccountSchema>;
export type CreateBudgetInput = z.infer<typeof createBudgetSchema>;
export type CreateGoalInput = z.infer<typeof createGoalSchema>;
