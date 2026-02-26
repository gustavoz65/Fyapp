# Mapeamento de Dados do Usuário — FiNext

> Documento de referência para o fluxo de **Onboarding**.
> Contém todos os campos e tabelas que compõem o perfil completo de um usuário no sistema, sejam eles preenchidos no cadastro tradicional (email/senha) ou via Google.

---

## Visão Geral das Tabelas Relacionadas ao Usuário

```
users (núcleo)
  │
  ├── user_oauth_providers   → métodos de login vinculados
  ├── user_sessions          → sessões JWT ativas
  ├── user_settings          → preferências e notificações
  │
  ├── bank_accounts          → contas bancárias do usuário
  ├── categories             → categorias personalizadas (+ sistema)
  ├── transactions           → transações financeiras
  ├── recurring_transactions → transações recorrentes
  ├── budgets                → orçamentos por categoria
  ├── goals                  → metas financeiras
  │     └── goal_contributions
  └── notifications          → notificações do sistema
```

---

## 1. Tabela `users` — Perfil Principal

| Coluna | Tipo | Obrigatório | Preenchido em | Editável | Observação |
|--------|------|-------------|---------------|----------|------------|
| `id` | CHAR(36) UUID | Auto | Cadastro | ❌ | Gerado automaticamente |
| `email` | VARCHAR(255) | ✅ Sim | Cadastro | ❌ | Único no sistema |
| `password_hash` | VARCHAR(255) | Condicional | Cadastro / SetPassword | ❌ | Vazio para quem veio do Google |
| `first_name` | VARCHAR(100) | ✅ Sim | Cadastro / Google | ✅ | Vem do Google no social login |
| `last_name` | VARCHAR(100) | ✅ Sim | Cadastro / Google | ✅ | Vem do Google no social login |
| `phone` | VARCHAR(20) | ❌ Opcional | Onboarding | ✅ | Campo para onboarding |
| `avatar_url` | VARCHAR(500) | ❌ Opcional | Google (automático) | ✅ | URL da foto do Google |
| `preferred_currency` | VARCHAR(3) | Padrão: `BRL` | Onboarding | ✅ | Ex: `BRL`, `USD`, `EUR` |
| `preferred_language` | VARCHAR(5) | Padrão: `pt-BR` | Onboarding | ✅ | Ex: `pt-BR`, `en-US` |
| `timezone` | VARCHAR(50) | Padrão: `America/Sao_Paulo` | Onboarding | ✅ | Fuso horário do usuário |
| `role` | ENUM | Auto: `user` | Sistema | ❌ | `user`, `admin`, `premium` |
| `email_verified` | BOOLEAN | Auto | Sistema / Google | ❌ | `true` para Google automaticamente |
| `email_verified_at` | DATETIME | Auto | Sistema | ❌ | |
| `last_login_at` | DATETIME | Auto | Login | ❌ | Atualizado a cada login |
| `is_active` | BOOLEAN | Auto: `true` | Cadastro | ❌ | `false` = conta desativada |
| `created_at` | DATETIME | Auto | Cadastro | ❌ | |
| `updated_at` | DATETIME | Auto | Qualquer update | ❌ | |

### Campos relevantes para Onboarding
- `phone` — número de telefone (opcional)
- `preferred_currency` — moeda padrão (por ora fixo em BRL)
- `preferred_language` — idioma (por ora fixo em pt-BR)
- `timezone` — fuso horário (15 opções brasileiras disponíveis)

---

## 2. Tabela `user_oauth_providers` — Métodos de Login

> Criada automaticamente no primeiro login via Google. Um usuário pode ter múltiplos providers.

| Coluna | Tipo | Preenchido em | Observação |
|--------|------|---------------|------------|
| `id` | VARCHAR(36) UUID | Auto | |
| `user_id` | VARCHAR(36) | Social Login | FK → users.id |
| `provider` | VARCHAR(50) | Social Login | Ex: `google`, `facebook`, `github` |
| `provider_user_id` | VARCHAR(255) | Social Login | UID único no provider (Firebase UID) |
| `provider_email` | VARCHAR(255) | Social Login | Email da conta Google |
| `provider_name` | VARCHAR(255) | Social Login | Nome completo vindo do provider |
| `provider_avatar_url` | VARCHAR(500) | Social Login | URL da foto do Google |
| `is_primary` | BOOLEAN | Social Login | `true` se foi o método usado no cadastro |
| `metadata` | JSON | Social Login | Dados extras do provider |
| `last_login_at` | DATETIME | A cada login | Atualizado a cada uso do provider |
| `created_at` | DATETIME | Auto | |
| `updated_at` | DATETIME | Auto | |

### Constraints
- `UNIQUE (user_id, provider)` → um usuário tem no máximo 1 vínculo por provider
- `UNIQUE (provider, provider_user_id)` → um UID do provider só pode estar em 1 conta

---

## 3. Tabela `user_settings` — Preferências e Notificações

> Criada automaticamente junto com o usuário com valores padrão.

| Coluna | Tipo | Padrão | Editável | Observação |
|--------|------|--------|----------|------------|
| `notification_email` | BOOLEAN | `true` | ✅ | Receber notificações por email |
| `notification_push` | BOOLEAN | `true` | ✅ | Notificações push |
| `notification_sms` | BOOLEAN | `false` | ✅ | Notificações por SMS |
| `budget_alerts` | BOOLEAN | `true` | ✅ | Alertas de orçamento estourado |
| `bill_reminders` | BOOLEAN | `true` | ✅ | Lembretes de contas a pagar |
| `bill_reminder_days` | INT | `3` | ✅ | Dias de antecedência para lembrete |
| `weekly_summary` | BOOLEAN | `true` | ✅ | Resumo semanal |
| `monthly_report` | BOOLEAN | `true` | ✅ | Relatório mensal |
| `low_balance_alert` | BOOLEAN | `true` | ✅ | Alerta de saldo baixo |
| `low_balance_threshold` | DECIMAL(15,2) | `100.00` | ✅ | Valor mínimo para alerta de saldo |
| `allow_manual_transactions` | BOOLEAN | `true` | ✅ | Permitir lançamentos manuais |
| `theme` | ENUM | `system` | ✅ | `light`, `dark`, `system` |
| `dashboard_layout` | JSON | `null` | ✅ | Layout customizado do dashboard |

### Campos relevantes para Onboarding
- `low_balance_threshold` — valor que o usuário considera "saldo baixo"
- `bill_reminder_days` — com quantos dias de antecedência quer ser lembrado

---

## 4. Tabela `bank_accounts` — Primeira Conta (Onboarding Essencial)

> O usuário precisa criar ao menos uma conta para usar o sistema.

| Coluna | Tipo | Obrigatório | Observação |
|--------|------|-------------|------------|
| `name` | VARCHAR(100) | ✅ | Ex: "Nubank", "Poupança Itaú" |
| `bank_name` | VARCHAR(100) | ❌ | Nome do banco |
| `account_type` | ENUM | ✅ | `checking`, `savings`, `credit_card`, `investment`, `cash`, `other` |
| `initial_balance` | DECIMAL(15,2) | Padrão: `0.00` | Saldo inicial |
| `currency` | VARCHAR(3) | Padrão: `BRL` | Moeda |
| `color` | VARCHAR(7) | Padrão: `#10B981` | Cor no dashboard |
| `icon` | VARCHAR(50) | Padrão: `bank` | Ícone |
| `include_in_total` | BOOLEAN | Padrão: `true` | Incluir no saldo total |
| `credit_limit` | DECIMAL(15,2) | Condicional | Obrigatório para `credit_card` |
| `closing_day` | TINYINT (1-31) | Condicional | Dia de fechamento (cartão) |
| `due_day` | TINYINT (1-31) | Condicional | Dia de vencimento (cartão) |

---

## 5. Fluxo de Preenchimento por Tipo de Cadastro

### Cadastro Tradicional (email + senha)

```
Tela de Registro → users criado com:
  ✅ email          (digitado)
  ✅ password_hash  (hash da senha)
  ✅ first_name     (digitado)
  ✅ last_name      (digitado)
  ⬜ phone          (opcional no registro)
  ⬜ avatar_url     (vazio)
  ✅ preferred_currency = "BRL"
  ✅ preferred_language  = "pt-BR"
  ✅ timezone            = "America/Sao_Paulo"
  ✅ email_verified = false

user_settings criado automaticamente com valores padrão
```

### Cadastro via Google (Firebase)

```
Login Google → users criado com:
  ✅ email          (vindo do Firebase token)
  ⬜ password_hash  (vazio — sem senha)
  ✅ first_name     (extraído do displayName do Google)
  ✅ last_name      (extraído do displayName do Google)
  ⬜ phone          (vazio)
  ✅ avatar_url     (photoURL do Google)
  ✅ preferred_currency = "BRL"
  ✅ preferred_language  = "pt-BR"
  ✅ timezone            = "America/Sao_Paulo"
  ✅ email_verified = true  ← Firebase já verificou

user_oauth_providers criado:
  ✅ provider           = "google"
  ✅ provider_user_id   = Firebase UID
  ✅ provider_email     = email do Google
  ✅ provider_name      = nome do Google
  ✅ provider_avatar_url = foto do Google
  ✅ is_primary         = true

user_settings criado automaticamente com valores padrão
```

---

## 6. Campos que o Onboarding Deve Coletar

Com base no mapeamento acima, o onboarding deve coletar os campos ainda **vazios ou com padrão genérico** após o cadastro:

### Etapa 1 — Perfil (Personalização)
| Campo | Tabela | Endpoint de Update |
|-------|--------|--------------------|
| `phone` | `users` | `PUT /users/me` |
| `timezone` | `users` | `PUT /users/me` |
| `avatar_url` | `users` | `PUT /users/me` |

> Para usuários do Google: `first_name`, `last_name`, `avatar_url` já vêm preenchidos — pode pular ou pré-preencher.

### Etapa 2 — Primeira Conta Bancária
| Campo | Tabela | Endpoint de Criação |
|-------|--------|---------------------|
| `name` | `bank_accounts` | `POST /accounts` |
| `account_type` | `bank_accounts` | `POST /accounts` |
| `initial_balance` | `bank_accounts` | `POST /accounts` |
| `credit_limit` | `bank_accounts` | Só para `credit_card` |
| `closing_day` / `due_day` | `bank_accounts` | Só para `credit_card` |

### Etapa 3 — Preferências de Alerta (Opcional)
| Campo | Tabela | Endpoint de Update |
|-------|--------|--------------------|
| `low_balance_threshold` | `user_settings` | `PUT /users/settings` |
| `bill_reminder_days` | `user_settings` | `PUT /users/settings` |
| `budget_alerts` | `user_settings` | `PUT /users/settings` |
| `bill_reminders` | `user_settings` | `PUT /users/settings` |

---

## 7. Endpoints Relevantes para o Onboarding

```
# Ler dados do usuário atual (para pré-preencher formulários)
GET  /api/v1/users/me
GET  /api/v1/users/settings

# Atualizar perfil
PUT  /api/v1/users/me
     Body: { first_name, last_name, phone, avatar_url, timezone }

# Atualizar preferências
PUT  /api/v1/users/settings
     Body: { low_balance_threshold, bill_reminder_days, budget_alerts, ... }

# Criar primeira conta
POST /api/v1/accounts
     Body: { name, account_type, initial_balance, bank_name, currency,
             credit_limit*, closing_day*, due_day* }

# Verificar métodos de login (para mostrar opções de senha no onboarding)
GET  /api/v1/auth/social/providers
     Retorna: { providers: [...], has_password: bool }

# Definir senha (para quem veio do Google)
POST /api/v1/auth/set-password
     Body: { new_password, confirm_password }
```

---

## 8. Timezones Brasileiras Disponíveis

| Valor | Label |
|-------|-------|
| `America/Noronha` | Fernando de Noronha (UTC-2) |
| `America/Sao_Paulo` | São Paulo / Rio de Janeiro (UTC-3) |
| `America/Fortaleza` | Fortaleza / Recife (UTC-3) |
| `America/Recife` | Recife (UTC-3) |
| `America/Bahia` | Salvador / Bahia (UTC-3) |
| `America/Belem` | Belém / Macapá (UTC-3) |
| `America/Maceio` | Maceió (UTC-3) |
| `America/Santarem` | Santarém (UTC-3) |
| `America/Araguaina` | Araguaína (UTC-3) |
| `America/Cuiaba` | Cuiabá / Campo Grande (UTC-4) |
| `America/Porto_Velho` | Porto Velho (UTC-4) |
| `America/Boa_Vista` | Boa Vista (UTC-4) |
| `America/Manaus` | Manaus (UTC-4) |
| `America/Eirunepe` | Eirunepé (UTC-5) |
| `America/Rio_Branco` | Rio Branco / Acre (UTC-5) |