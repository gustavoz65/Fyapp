# 📚 Guia de Arquitetura - Cashing API

## 🏗️ Estrutura do Projeto (Clean Architecture)

O projeto segue uma arquitetura em camadas bem definida:

```
Handler (HTTP) → Service (Regras de Negócio) → Repository (Banco de Dados)
```

### 1. **Model** (`internal/model/`)
- Define as estruturas de dados (tabelas do banco)
- DTOs (Data Transfer Objects) para requests/responses da API

### 2. **Repository** (`internal/repository/`)
- Acessa o banco de dados diretamente
- Executa queries SQL (CRUD básico)
- Exemplo: `userRepo.Create()`, `userRepo.FindByEmail()`

### 3. **Service** (`internal/service/`)
- Contém a lógica de negócio
- Valida dados, aplica regras
- Chama os repositories
- Exemplo: `authService.Login()` valida senha e gera token JWT

### 4. **Handler** (`internal/handler/`)
- Recebe requests HTTP
- Valida input do usuário
- Chama os services
- Retorna responses JSON

### 5. **Router** (`internal/router/`)
- Define todas as rotas da API
- Configura middlewares (autenticação, CORS, logs)

---

## 📁 Módulos do Sistema

### 🔐 **1. AUTH (Autenticação)**
**Propósito**: Login, registro, tokens JWT

**Rotas Públicas** (não precisa estar logado):
- `POST /api/v1/auth/register` - Registrar novo usuário
- `POST /api/v1/auth/login` - Fazer login
- `POST /api/v1/auth/refresh` - Renovar access token
- `POST /api/v1/auth/logout` - Fazer logout

**Rotas Autenticadas**:
- `POST /api/v1/auth/change-password` - Mudar senha

**Fluxo**:
1. Usuário envia email/senha
2. `AuthHandler` recebe → valida input
3. `AuthService` verifica se usuário existe e senha está correta
4. Gera JWT token (access + refresh)
5. Retorna tokens

---

### 👤 **2. USERS (Usuários)**
**Propósito**: Gerenciar perfil do usuário

**Rotas**:
- `GET /api/v1/users/me` - Ver meu perfil
- `PUT /api/v1/users/me` - Atualizar meu perfil
- `DELETE /api/v1/users/me` - Desativar minha conta
- `GET /api/v1/users/settings` - Ver configurações
- `PUT /api/v1/users/settings` - Atualizar configurações

**Fluxo**:
- `UserHandler` → `UserService` → `UserRepository` → MySQL

---

### 🏷️ **3. CATEGORIES (Categorias)**
**Propósito**: Organizar transações por categoria (ex: Alimentação, Transporte)

**Rotas**:
- `GET /api/v1/categories` - Listar todas
- `GET /api/v1/categories/:id` - Ver uma categoria
- `POST /api/v1/categories` - Criar categoria
- `PUT /api/v1/categories/:id` - Atualizar categoria
- `DELETE /api/v1/categories/:id` - Deletar categoria

**Exemplo de uso**:
```json
POST /api/v1/categories
{
  "name": "Alimentação",
  "type": "expense",
  "icon": "🍔",
  "color": "#FF5733"
}
```

---

### 🏦 **4. BANK ACCOUNTS (Contas Bancárias)**
**Propósito**: Gerenciar contas bancárias do usuário

**Rotas**:
- `GET /api/v1/accounts` - Listar todas minhas contas
- `GET /api/v1/accounts/balance` - Saldo total de todas as contas
- `GET /api/v1/accounts/:id` - Ver uma conta
- `POST /api/v1/accounts` - Criar conta
- `PUT /api/v1/accounts/:id` - Atualizar conta
- `DELETE /api/v1/accounts/:id` - Deletar conta

**Exemplo**:
```json
POST /api/v1/accounts
{
  "name": "Nubank",
  "type": "checking",
  "balance": 1500.00,
  "currency": "BRL"
}
```

---

### 💸 **5. TRANSACTIONS (Transações)**
**Propósito**: Registrar receitas e despesas

**Rotas**:
- `GET /api/v1/transactions` - Listar todas
- `GET /api/v1/transactions/upcoming` - Ver próximas transações (recorrentes)
- `GET /api/v1/transactions/:id` - Ver uma transação
- `POST /api/v1/transactions` - Criar transação
- `PUT /api/v1/transactions/:id` - Atualizar transação
- `DELETE /api/v1/transactions/:id` - Deletar transação
- `PATCH /api/v1/transactions/:id/pay` - Marcar como paga

**Exemplo**:
```json
POST /api/v1/transactions
{
  "description": "Almoço no restaurante",
  "amount": 45.50,
  "type": "expense",
  "category_id": 1,
  "account_id": 1,
  "date": "2026-02-04",
  "is_recurring": false
}
```

---

### 💰 **6. BUDGETS (Orçamentos)**
**Propósito**: Definir limites de gastos por categoria/mês

**Rotas**:
- `GET /api/v1/budgets` - Listar todos
- `GET /api/v1/budgets/summary` - Resumo dos orçamentos (quanto gastou vs limite)
- `GET /api/v1/budgets/:id` - Ver um orçamento
- `POST /api/v1/budgets` - Criar orçamento
- `PUT /api/v1/budgets/:id` - Atualizar orçamento
- `DELETE /api/v1/budgets/:id` - Deletar orçamento

**Exemplo**:
```json
POST /api/v1/budgets
{
  "name": "Alimentação Fevereiro",
  "category_id": 1,
  "amount": 800.00,
  "period": "monthly",
  "start_date": "2026-02-01",
  "alert_threshold": 80
}
```

---

### 🎯 **7. GOALS (Metas Financeiras)**
**Propósito**: Definir objetivos de economia (ex: guardar R$ 5000 para viagem)

**Rotas**:
- `GET /api/v1/goals` - Listar todas
- `GET /api/v1/goals/summary` - Resumo das metas
- `GET /api/v1/goals/:id` - Ver uma meta
- `POST /api/v1/goals` - Criar meta
- `PUT /api/v1/goals/:id` - Atualizar meta
- `DELETE /api/v1/goals/:id` - Deletar meta
- `POST /api/v1/goals/:id/contributions` - Adicionar contribuição
- `GET /api/v1/goals/:id/contributions` - Listar contribuições

**Exemplo**:
```json
POST /api/v1/goals
{
  "name": "Viagem para Europa",
  "target_amount": 10000.00,
  "current_amount": 0,
  "deadline": "2026-12-31",
  "priority": "high"
}
```

---

### 📊 **8. DASHBOARD (Painel de Controle)**
**Propósito**: Visão geral das finanças

**Rotas**:
- `GET /api/v1/dashboard` - Resumo geral
- `GET /api/v1/dashboard/cash-flow` - Fluxo de caixa (entradas vs saídas)
- `GET /api/v1/dashboard/income-expense` - Receitas vs Despesas
- `GET /api/v1/dashboard/monthly-comparison` - Comparação mensal
- `GET /api/v1/dashboard/account-balances` - Saldos das contas

---

### 🔔 **9. NOTIFICATIONS (Notificações)**
**Propósito**: Alertas do sistema (orçamento estourado, conta a pagar, etc)

**Rotas**:
- `GET /api/v1/notifications` - Listar todas
- `GET /api/v1/notifications/unread` - Apenas não lidas
- `GET /api/v1/notifications/unread/count` - Quantidade de não lidas
- `PATCH /api/v1/notifications/:id/read` - Marcar como lida
- `PATCH /api/v1/notifications/read-all` - Marcar todas como lidas
- `DELETE /api/v1/notifications/:id` - Deletar notificação

---

## 🔐 Autenticação

### Como funciona:
1. **Login**: Recebe email/senha → retorna `access_token` e `refresh_token`
2. **Access Token**: Válido por 15 minutos, usado em todas as rotas protegidas
3. **Refresh Token**: Válido por 7 dias, usado para renovar o access token
4. **Header**: Todas as rotas autenticadas precisam do header:
   ```
   Authorization: Bearer {access_token}
   ```

### Middleware de Autenticação:
- Arquivo: `internal/middleware/auth.go`
- Valida o JWT token
- Extrai o `user_id` do token
- Adiciona no contexto da request

---

## 🧪 Como Testar as Rotas

### Opção 1: cURL

```bash
# 1. Registrar usuário
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Gustavo",
    "email": "gustavo@example.com",
    "password": "senha123"
  }'

# 2. Fazer login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "gustavo@example.com",
    "password": "senha123"
  }'

# Resposta: { "access_token": "eyJhbG...", "refresh_token": "eyJhbG..." }

# 3. Usar token em rota protegida
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer {SEU_ACCESS_TOKEN}"
```

### Opção 2: Postman / Insomnia

1. Importe as rotas do arquivo `router.go`
2. Configure a variável `{{baseUrl}}` = `http://localhost:8080`
3. Após login, salve o token em variável
4. Use `{{token}}` no header Authorization

### Opção 3: REST Client (VS Code Extension)

Crie um arquivo `.http`:

```http
### Registrar
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json

{
  "name": "Gustavo",
  "email": "gustavo@example.com",
  "password": "senha123"
}

### Login
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json

{
  "email": "gustavo@example.com",
  "password": "senha123"
}

### Ver perfil
GET http://localhost:8080/api/v1/users/me
Authorization: Bearer {COLE_SEU_TOKEN_AQUI}
```

---

## 📂 Estrutura de Pastas

```
backend/
├── cmd/
│   └── go-boilerplate/     # Arquivo main.go (inicia o servidor)
├── internal/
│   ├── config/             # Configurações (carrega .env)
│   ├── database/           # Conexão com MySQL
│   ├── handler/            # Controllers HTTP (recebe requests)
│   ├── middleware/         # Auth, CORS, Logger, Recovery
│   ├── model/              # Estruturas de dados (User, Transaction, etc)
│   ├── repository/         # Acesso ao banco de dados
│   ├── router/             # Definição de rotas
│   ├── service/            # Lógica de negócio
│   └── validation/         # Validações customizadas
├── migrations/             # Scripts SQL para criar tabelas
├── .env                    # Variáveis de ambiente
└── Taskfile.yml           # Comandos do projeto (task run, task migrate, etc)
```

---

## 🚀 Comandos Úteis

```bash
# Rodar servidor
task run

# Rodar migrações
task migrate-up

# Reverter migrações
task migrate-down

# Ver logs do MySQL
docker logs cashing-mysql

# Ver logs do Redis
docker logs cashing-redis
```

---

## 💡 Dicas

1. **Sempre use o Postman/Insomnia** para testar - é mais fácil que cURL
2. **Comece testando na ordem**:
   - ✅ Health check (`/health`)
   - ✅ Register (`/auth/register`)
   - ✅ Login (`/auth/login`)
   - ✅ Get profile (`/users/me`)
   - ✅ Criar categoria, conta, transação, etc

3. **Erros comuns**:
   - `401 Unauthorized` → Token inválido ou expirado
   - `400 Bad Request` → JSON mal formatado ou campo obrigatório faltando
   - `404 Not Found` → Recurso não existe (ex: categoria com ID 999)
   - `500 Internal Server Error` → Erro no servidor (veja os logs)

4. **Logs**: Todos os erros aparecem no terminal onde você rodou `task run`

---

## 📝 Próximos Passos

1. ✅ **Testar rotas básicas** (auth, users)
2. ✅ **Criar uma categoria** de teste
3. ✅ **Criar uma conta bancária**
4. ✅ **Registrar uma transação**
5. ✅ **Ver o dashboard**

Qualquer dúvida, me chama! 🚀
