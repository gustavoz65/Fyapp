# Configuração de Email com Resend

Este documento descreve como configurar o sistema de envio de emails transacionais usando a plataforma Resend.

## 📋 Checklist de Setup

### 1. Criar Conta no Resend

- Acesse [https://resend.com/](https://resend.com/)
- Realize o cadastro com seu email
- Ative sua conta através do link enviado por email

### 2. Obter API Key

1. Faça login em [https://resend.com/](https://resend.com/)
2. Acesse "API Keys" no dashboard
3. Clique em "Create API Key"
4. Selecione o projeto (ou deixe como padrão)
5. Copie a chave gerada (formato: `re_xxxxx`)

### 3. Configurar Variáveis de Ambiente

Adicione as seguintes variáveis no arquivo `.env`:

```env
# Resend API Key (obrigatório para enviar emails)
FINEXT_INTEGRATION_RESEND_API_KEY=re_xxxxx

# Nome do remetente (padrão: Finext)
FINEXT_INTEGRATION_SENDER_NAME=Finext

# Email do remetente (padrão: noreply@resend.dev - sandbox)
FINEXT_INTEGRATION_SENDER_EMAIL=noreply@resend.dev

# Caminho base dos templates HTML (padrão: templates/emails)
FINEXT_INTEGRATION_TEMPLATES_PATH=templates/emails
```

### 4. Verificar Domínio no Resend (Produção)

⚠️ **IMPORTANTE PARA PRODUÇÃO**

O email `noreply@resend.dev` funciona apenas para testes/sandbox.

Para enviar emails de produção com um domínio personalizado:

1. Acesse "Domains" em [https://resend.com/domains](https://resend.com/domains)
2. Clique em "Add Domain"
3. Insira seu domínio (ex: noreply@seudominio.com)
4. Siga as instruções para verificar o domínio via DNS (SPF, DKIM, DMARC)
5. Após verificação, atualize `FINEXT_INTEGRATION_SENDER_EMAIL` com seu domínio

### 5. Criar Templates de Email

Os templates devem estar em `backend/templates/emails/` como arquivos `.html`.

Templates atualmente suportados:

- **welcome.html** - Email de boas-vindas enviado ao registrar um novo usuário
  - Variáveis: `{{.FirstName}}` - Primeiro nome do usuário

Exemplo de estrutura de template:

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Título do Email</title>
</head>
<body>
    <h1>Bem-vindo {{.FirstName}}!</h1>
    <p>Seu email foi verificado com sucesso.</p>
</body>
</html>
```

## 🚀 Testando a Configuração

### Teste em Desenvolvimento

```bash
# 1. Verifique se o servidor inicia sem erros
cd backend
go run cmd/main.go

# 2. Observe o log ao iniciar:
# - Se `FINEXT_INTEGRATION_RESEND_API_KEY` está vazio: 
#   "Resend API Key is not configured - email notifications will not be sent"
# - Se está configurado:
#   "Email service (Resend) initialized successfully"

# 3. Faça um registro de novo usuário
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@example.com",
    "password": "SenhaSegura123!",
    "first_name": "João",
    "last_name": "Silva"
  }'

# 4. Verifique se a tarefa de email foi enfileirada:
# - Logs mostrarão: "welcome email task enqueued"
# - O email será processado pelo job worker

# 5. Verifique se o email foi enviado:
# - No sandbox: verifique a email de teste registrada no Resend
# - Nos logs do Resend dashboard: Activity > Email Logs
```

### Troubleshooting

| Problema | Causa | Solução |
|----------|-------|--------|
| "Resend API Key is not configured" | Variável de env vazia | Configure `FINEXT_INTEGRATION_RESEND_API_KEY` |
| Failed to parse email template | Template não existe | Verifique existência de `templates/emails/welcome.html` |
| Email não é enviado | Job worker não está rodando | Verifique Redis conectado e job server iniciado |
| "401 Unauthorized" do Resend | API Key inválida | Gere nova chave em [https://resend.com/api-keys](https://resend.com/api-keys) |
| Email marcado como spam | Domínio não verificado | Verifique DKIM/SPF em [https://resend.com/domains](https://resend.com/domains) |

## 📧 Tipos de Email Implementados

### 1. Email de Boas-vindas (Welcome)

- **Disparado:** Ao registrar novo usuário
- **Template:** `templates/emails/welcome.html`
- **Variáveis:** `FirstName`
- **Job Queue:** `default`
- **Retry:** Máximo 3 tentativas com timeout de 30s

Exemplo de fila:

```go
task, err := job.NewWelcomeEmailTask(user.Email, user.FirstName)
if err == nil {
    _, err = jobService.Client.Enqueue(task)
}
```

## 🔐 Segurança

### Variáveis de Ambiente

- **NÃO** commite `FINEXT_INTEGRATION_RESEND_API_KEY` em controle de versão
- Use `.env.local` ou arquivo seguro em produção
- Para Docker, passe como variável de ambiente segura (secrets)

### Domínio de Produção

- Sempre majore verificação de domínio via DKIM/SPF/DMARC
- Monitore taxa de bounce/spam no Resend dashboard
- Implemente unsubscribe links em emails (recomendado)

## 📚 Referências

- [Resend Documentation](https://resend.com/docs)
- [Resend API Reference](https://resend.com/docs/api-reference/emails/send)
- [Domain Verification Guide](https://resend.com/docs/dashboard/domains)
- [Rate Limits](https://resend.com/docs/rate-limiting)

## 🛠️ Configurações Avançadas

### Adicionar Novo Template de Email

1. Criar arquivo em `backend/templates/emails/{nome}.html`
2. Adicionar constante em `backend/internal/lib/email/templates.go`:
   ```go
   const (
       TemplateWelcome Template = "welcome"
       TemplateForgotPassword Template = "forgot_password"  // novo
   )
   ```
3. Criar método helper em `backend/internal/lib/email/client.go`:
   ```go
   func (c *Client) SendForgotPasswordEmail(to, resetLink string) error {
       data := map[string]string{
           "ResetLink": resetLink,
       }
       return c.SendEmail(to, "Redefina sua senha", TemplateForgotPassword, data)
   }
   ```
4. Criar task em `backend/internal/lib/utils/job/email_tasks.go`
5. Registrar handler em `backend/internal/lib/utils/job/job.go`
6. Enfileirar na service quando necessário

## 💡 Dicas

- Use templates responsive para melhor visualização em mobile
- Inclua link de unsubscribe para manter reputação de domínio
- Teste emails com [Resend Email Preview](https://resend.com/emails/preview)
- Monitore taxa de entrega no Resend dashboard
- Configure alertas para taxa alta de bounce/spam
