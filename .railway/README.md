# 🚂 Railway - Configuração Beta 1

## ⚠️ IMPORTANTE: Docker e Volumes

**NÃO precisa configurar Docker Image nem Volume na Railway!**

A Railway faz o build automaticamente:
- **Backend**: Detecta `go.mod` e usa o Dockerfile para build multi-stage
- **Frontend**: Detecta `package.json` e usa o Dockerfile Next.js

Apenas configure:
1. **Root Directory**: `/backend` ou `/frontend`
2. **Variáveis de Ambiente**: Veja seção abaixo
3. **Wait for CI**: Ative nas settings (garante deploy só após CI passar)

---

## 📋 Service IDs (para GitHub Secrets)

Para pegar os IDs dos services:

```bash
# Instalar Railway CLI
npm install -g @railway/cli

# Login
railway login

# Link ao projeto
railway link

# Ver services
railway status

# Vai mostrar algo como:
# Service: backend (abc123def456)
# Service: frontend (xyz789uvw012)
```

Adicione esses IDs como secrets no GitHub:
- `RAILWAY_SERVICE_BACKEND` = `abc123def456`
- `RAILWAY_SERVICE_FRONTEND` = `xyz789uvw012`

---

## 🖥️ Localhost vs Produção - Como o Sistema Identifica Ambiente

### Backend (Go) ✅ JÁ ESTÁ PREPARADO

O backend usa variável de ambiente `FINEXT_PRIMARY_ENV` para identificar o ambiente:

```env
# Desenvolvimento (local)
FINEXT_PRIMARY_ENV=development

# Produção (Railway/Fly.io)
FINEXT_PRIMARY_ENV=production
```

**Como funciona:**
- Em **development**: CORS permite `http://localhost:3000` e `http://localhost:4000` (fallback)
- Em **production**: CORS usa `FINEXT_SERVER_CORS_ALLOWED_ORIGINS` (veja seção CORS abaixo)

**Código**: [backend/internal/config/config.go:28](../backend/internal/config/config.go#L28)

### Frontend (Next.js) ⚠️ PRECISA AJUSTAR

**PROBLEMA ENCONTRADO**: [frontend/src/lib/api.ts:3](../frontend/src/lib/api.ts#L3) está **HARDCODED** em localhost:

```typescript
const API_BASE_URL = "http://localhost:3000/api/v1"; // ❌ HARDCODED!
```

**SOLUÇÃO**: Usar variável de ambiente `NEXT_PUBLIC_API_URL`:

```typescript
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3000/api/v1";
```

**Configuração**:
- **Local** (.env.local): `NEXT_PUBLIC_API_URL=http://localhost:3000/api/v1`
- **Railway**: `NEXT_PUBLIC_API_URL=${{Backend.RAILWAY_PUBLIC_DOMAIN}}/api/v1`
- **Produção**: `NEXT_PUBLIC_API_URL=https://api.finext.com.br/api/v1`

### Outras Referências a Localhost

1. **[frontend/next.config.ts:37](../frontend/next.config.ts#L37)** - CSP com `http://localhost:3000`
   - ⚠️ Em produção, CSP deve incluir domínio real do backend
   - Solução: Usar variável de ambiente para construir CSP dinamicamente

2. **[backend/internal/middleware/cors.go:26](../backend/internal/middleware/cors.go#L26)** - Default localhost
   - ✅ Já está correto! É apenas fallback quando `CORS_ALLOWED_ORIGINS` não está definido
   - Em produção, sempre configure `FINEXT_SERVER_CORS_ALLOWED_ORIGINS`

---

## 🚨 CORS - LEIA COM ATENÇÃO (evita problemas futuros)

### Por que CORS vai dar problema?

CORS (Cross-Origin Resource Sharing) bloqueia requisições entre domínios diferentes por segurança.

**Cenário de erro comum:**
```
Frontend: https://finext-frontend.up.railway.app
Backend:  https://finext-backend.up.railway.app

❌ ERRO: "Access to fetch at 'https://finext-backend.up.railway.app/api/v1/auth/login'
from origin 'https://finext-frontend.up.railway.app' has been blocked by CORS policy"
```

### Como Configurar CORS Corretamente

#### 1. Backend - Variáveis de Ambiente OBRIGATÓRIAS

Configure na Railway (Settings → Variables):

```env
# Ambiente
FINEXT_PRIMARY_ENV=production

# CORS - URL do Frontend (SEM barra no final!)
FINEXT_SERVER_CORS_ALLOWED_ORIGINS=https://finext-frontend.up.railway.app

# Ou múltiplas origens separadas por vírgula:
FINEXT_SERVER_CORS_ALLOWED_ORIGINS=https://finext-frontend.up.railway.app,https://www.finext.com.br
```

**⚠️ IMPORTANTE:**
- NÃO use `*` em produção (quebra credentials/cookies)
- NÃO coloque barra `/` no final da URL
- Use HTTPS em produção (nunca HTTP)

#### 2. Frontend - Apontar para Backend Correto

Configure na Railway (Settings → Variables):

```env
# Backend API - Use referência Railway
NEXT_PUBLIC_API_URL=${{Backend.RAILWAY_PUBLIC_DOMAIN}}/api/v1

# Ou manualmente (se preferir):
NEXT_PUBLIC_API_URL=https://finext-backend.up.railway.app/api/v1
```

#### 3. Verificar CSP (Content Security Policy)

O Next.js precisa permitir conexões ao backend no CSP.

**Atualmente** em [frontend/next.config.ts:37](../frontend/next.config.ts#L37):
```typescript
"connect-src 'self' http://localhost:3000 ..." // ❌ Localhost hardcoded
```

**Em produção**, adicione o domínio do backend:
```typescript
"connect-src 'self' https://finext-backend.up.railway.app ..."
```

### Testando CORS

Depois de configurar, teste no browser console:

```javascript
// Deve funcionar SEM erros de CORS
fetch('https://finext-backend.up.railway.app/api/v1/health')
  .then(r => r.json())
  .then(console.log)
```

### CORS no Código

**Backend** usa middleware customizado: [backend/internal/middleware/cors.go](../backend/internal/middleware/cors.go)

**Configuração atual:**
- ✅ Permite credenciais (cookies)
- ✅ Métodos: GET, POST, PUT, PATCH, DELETE, OPTIONS
- ✅ Headers: Authorization, Content-Type, etc
- ✅ MaxAge: 24h (cacheia preflight)

**Logs de CORS:**
- Se aparecer "Origin not allowed" nos logs → Configure `FINEXT_SERVER_CORS_ALLOWED_ORIGINS`
- Se aparecer "Method not allowed" → Verifique se método HTTP está na lista permitida

---

## 🔐 Variáveis de Ambiente

### Backend Service

Configure na Railway UI (Settings → Variables):

**⚠️ IMPORTANTE**: Todas as variáveis do backend usam prefixo `FINEXT_`

```env
# ========================================
# AMBIENTE E SERVIDOR
# ========================================
FINEXT_PRIMARY_ENV=production
FINEXT_SERVER_PORT=3000

# CORS - URL(s) do Frontend (SEM barra no final!)
# 🚨 OBRIGATÓRIO EM PRODUÇÃO - Evita erros de CORS!
FINEXT_SERVER_CORS_ALLOWED_ORIGINS=https://finext-frontend.up.railway.app
# Ou use referência Railway:
# FINEXT_SERVER_CORS_ALLOWED_ORIGINS=${{Frontend.RAILWAY_PUBLIC_DOMAIN}}
# Ou múltiplas origens:
# FINEXT_SERVER_CORS_ALLOWED_ORIGINS=https://finext-frontend.up.railway.app,https://www.finext.com.br

# ========================================
# DATABASE (MySQL)
# ========================================
# Opção 1: Railway MySQL addon (Recomendado para Beta)
FINEXT_DATABASE_HOST=${{MySQL.MYSQLHOST}}
FINEXT_DATABASE_PORT=${{MySQL.MYSQLPORT}}
FINEXT_DATABASE_USER=${{MySQL.MYSQLUSER}}
FINEXT_DATABASE_PASSWORD=${{MySQL.MYSQLPASSWORD}}
FINEXT_DATABASE_DB_NAME=${{MySQL.MYSQLDATABASE}}

# Opção 2: Database Externo (Hostinger - Futuro)
# FINEXT_DATABASE_HOST=seu-host.hostinger.com
# FINEXT_DATABASE_PORT=3306
# FINEXT_DATABASE_USER=seu-usuario
# FINEXT_DATABASE_PASSWORD=sua-senha-segura
# FINEXT_DATABASE_DB_NAME=finext_db
# FINEXT_DATABASE_SSL_MODE=require

# ========================================
# REDIS
# ========================================
# Opção 1: Railway Redis addon
FINEXT_REDIS_ADDRESS=${{Redis.REDIS_URL}}
# Railway Redis normalmente já inclui password na URL

# Opção 2: Redis externo
# FINEXT_REDIS_ADDRESS=seu-host.hostinger.com:6379
# FINEXT_REDIS_PASSWORD=sua-senha-redis
# FINEXT_REDIS_DB=0

# ========================================
# AUTENTICAÇÃO (JWT)
# ========================================
# 🚨 GERE UM SECRET FORTE! Use: openssl rand -base64 32
FINEXT_AUTH_SECRET_KEY=SUA-CHAVE-SUPER-SEGURA-AQUI-MINIMO-32-CARACTERES
FINEXT_AUTH_ACCESS_TOKEN_DURATION=15
FINEXT_AUTH_REFRESH_TOKEN_DURATION=168
FINEXT_AUTH_ISSUER=finext-api

# ========================================
# FIREBASE
# ========================================
FINEXT_FIREBASE_ENABLED=true
FINEXT_FIREBASE_PROJECT_ID=seu-projeto-firebase-id

# Opção 1: JSON direto (Recomendado - mais simples)
FINEXT_FIREBASE_SERVICE_ACCOUNT_JSON={"type":"service_account","project_id":"seu-projeto","private_key_id":"...","private_key":"-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n","client_email":"..."}

# Opção 2: Path para arquivo (usar Railway Config File)
# FINEXT_FIREBASE_SERVICE_ACCOUNT_PATH=/app/firebase-key.json

# ========================================
# INTEGRAÇÕES
# ========================================
# Resend (Email)
FINEXT_INTEGRATION_RESEND_API_KEY=re_sua_chave_resend_aqui

# ========================================
# OBSERVABILITY (Opcional)
# ========================================
FINEXT_OBSERVABILITY_SERVICE_NAME=finext-backend-beta
FINEXT_OBSERVABILITY_ENVIRONMENT=production
FINEXT_OBSERVABILITY_LOGGING_LEVEL=info
FINEXT_OBSERVABILITY_LOGGING_FORMAT=json

# New Relic (opcional)
# FINEXT_OBSERVABILITY_NEW_RELIC_LICENSE_KEY=sua-chave-newrelic
# FINEXT_OBSERVABILITY_NEW_RELIC_APP_LOG_FORWARDING_ENABLED=true
# FINEXT_OBSERVABILITY_NEW_RELIC_DISTRIBUTED_TRACING_ENABLED=true
```

**📝 Notas:**
- **CORS**: Se não configurar `FINEXT_SERVER_CORS_ALLOWED_ORIGINS`, vai dar erro 403/CORS em produção!
- **Firebase JSON**: Copie o JSON completo do arquivo de credenciais do Firebase Console
- **JWT Secret**: NUNCA use o mesmo secret em dev e produção
- **Database**: Railway addons setam variáveis automaticamente com `${{Service.VAR}}`

### Frontend Service

Configure na Railway UI (Settings → Variables):

```env
# ========================================
# NEXT.JS
# ========================================
NODE_ENV=production
PORT=3000

# ========================================
# BACKEND API
# ========================================
# 🚨 OBRIGATÓRIO - Frontend precisa saber onde está o backend!
# ⚠️ ATENÇÃO: Inclua /api/v1 no final!
NEXT_PUBLIC_API_URL=${{Backend.RAILWAY_PUBLIC_DOMAIN}}/api/v1
# Ou manualmente:
# NEXT_PUBLIC_API_URL=https://finext-backend.up.railway.app/api/v1

# ========================================
# FIREBASE (Credenciais Públicas - Client-side)
# ========================================
# 📋 Pegue no Firebase Console → Project Settings → General → Your apps
NEXT_PUBLIC_FIREBASE_API_KEY=AIzaSyBexampleKeyHere123456789
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=seu-projeto.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=seu-projeto-id
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=seu-projeto.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=123456789012
NEXT_PUBLIC_FIREBASE_APP_ID=1:123456789012:web:abc123def456
NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID=G-XXXXXXXXXX
```

**📝 Notas:**
- **NEXT_PUBLIC_API_URL**: DEVE incluir `/api/v1` no final (é o path base da API)
- **Firebase**: Essas credenciais são públicas (client-side), não são secretas
- **NODE_ENV**: Railway seta automaticamente para `production`

---

## 🔗 Service References

Railway permite referenciar variáveis entre services:

- `${{Backend.RAILWAY_PUBLIC_DOMAIN}}` - URL pública do backend
- `${{Frontend.RAILWAY_PUBLIC_DOMAIN}}` - URL pública do frontend
- `${{MySQL.MYSQLHOST}}` - Host do MySQL addon
- `${{Redis.REDIS_URL}}` - URL do Redis addon

Essas referências são atualizadas automaticamente!

---

## 📦 Addons Recomendados

Na Railway, adicione:

1. **MySQL** (ou use Hostinger no futuro)
   - New → Database → MySQL
   - Conecta automaticamente via `${{MySQL.*}}`

2. **Redis** (ou use Hostinger no futuro)
   - New → Database → Redis
   - Conecta automaticamente via `${{Redis.REDIS_URL}}`

---

## 🚀 Deploy Manual (se precisar)

```bash
# Login
railway login

# Link ao projeto
railway link

# Deploy backend
cd backend
railway up --service backend

# Deploy frontend
cd ../frontend
railway up --service frontend

# Ver logs
railway logs --service backend
railway logs --service frontend
```

---

## 🔄 Fluxo de Deploy

```
1. Push para main
   ↓
2. GitHub Actions roda CI
   ├─ Backend CI (lint, tests, vuln scan)
   └─ Frontend CI (lint, type check, build)
   ↓
3. Se CI passar → Deploy
   ├─ railway up backend
   └─ railway up frontend
   ↓
4. Railway faz build com Dockerfile
   ├─ Backend: Dockerfile multi-stage
   └─ Frontend: Dockerfile Next.js
   ↓
5. Deploy completo! 🎉
```

---

## 🛠️ Troubleshooting

### Erro: "No service found"
```bash
railway link  # Re-link ao projeto
railway status  # Ver services disponíveis
```

### Erro: "Unauthorized"
```bash
railway logout
railway login
# Ou regere o token na Railway
```

### Build falha no Railway
- Verifique logs: Railway → Service → Deployments → Ver logs
- Verifique Dockerfile
- Verifique variáveis de ambiente

### Backend não conecta no DB
- Verifique variáveis `DB_*`
- Verifique se MySQL addon está criado
- Teste conexão: `railway run --service backend printenv | grep DB`

---

## 💰 Custos Estimados (Beta)

Railway Free Tier:
- $5/mês de crédito grátis
- Depois: ~$0.01/hora por service

Estimativa para Beta:
- Backend: ~$7/mês
- Frontend: ~$7/mês
- MySQL: ~$5/mês
- Redis: ~$3/mês
**Total: ~$22/mês** (primeiros $5 grátis)

---

## 🎯 Próximos Passos (Produção)

Quando migrar para produção:

1. **Backend → Fly.io**
   - Criar workflow `.github/workflows/deploy-fly.yml`
   - Usar `fly deploy` via GitHub Actions

2. **Bancos → Hostinger**
   - Atualizar variáveis `DB_HOST`, `REDIS_URL`
   - Configurar conexões seguras (SSL)

3. **Frontend → ?**
   - Vercel (recomendado para Next.js)
   - Ou Fly.io também

4. **Cassandra → Hostinger**
   - Adicionar driver Cassandra no backend
   - Variáveis `FINEXT_CASSANDRA_*`

---

## 📚 Resumo Rápido - Checklist de Deploy

### ✅ Configuração Inicial

- [ ] Criar projeto na Railway
- [ ] Criar service **backend** (Root Directory: `/backend`)
- [ ] Criar service **frontend** (Root Directory: `/frontend`)
- [ ] Adicionar MySQL addon
- [ ] Adicionar Redis addon
- [ ] Ativar "Wait for CI" nas settings de cada service

### ✅ Variáveis de Ambiente Backend (OBRIGATÓRIAS)

```env
FINEXT_PRIMARY_ENV=production
FINEXT_SERVER_PORT=3000
FINEXT_SERVER_CORS_ALLOWED_ORIGINS=${{Frontend.RAILWAY_PUBLIC_DOMAIN}}
FINEXT_DATABASE_HOST=${{MySQL.MYSQLHOST}}
FINEXT_DATABASE_PORT=${{MySQL.MYSQLPORT}}
FINEXT_DATABASE_USER=${{MySQL.MYSQLUSER}}
FINEXT_DATABASE_PASSWORD=${{MySQL.MYSQLPASSWORD}}
FINEXT_DATABASE_DB_NAME=${{MySQL.MYSQLDATABASE}}
FINEXT_REDIS_ADDRESS=${{Redis.REDIS_URL}}
FINEXT_AUTH_SECRET_KEY=<GERE_COM_OPENSSL>
FINEXT_FIREBASE_ENABLED=true
FINEXT_FIREBASE_PROJECT_ID=<SEU_PROJETO_ID>
FINEXT_FIREBASE_SERVICE_ACCOUNT_JSON=<JSON_COMPLETO>
FINEXT_INTEGRATION_RESEND_API_KEY=<SUA_CHAVE_RESEND>
```

### ✅ Variáveis de Ambiente Frontend (OBRIGATÓRIAS)

```env
NODE_ENV=production
NEXT_PUBLIC_API_URL=${{Backend.RAILWAY_PUBLIC_DOMAIN}}/api/v1
NEXT_PUBLIC_FIREBASE_API_KEY=<SUA_CHAVE>
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=<SEU_PROJETO>.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=<SEU_PROJETO_ID>
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=<SEU_PROJETO>.appspot.com
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=<SEU_SENDER_ID>
NEXT_PUBLIC_FIREBASE_APP_ID=<SEU_APP_ID>
NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID=<SEU_MEASUREMENT_ID>
```

### ✅ GitHub Secrets (para Actions)

- [ ] `RAILWAY_TOKEN` - Token de API da Railway
- [ ] `RAILWAY_SERVICE_BACKEND` - Service ID do backend
- [ ] `RAILWAY_SERVICE_FRONTEND` - Service ID do frontend

### ⚠️ Problemas Comuns e Soluções

| Problema | Causa | Solução |
|----------|-------|---------|
| **CORS Error** | `FINEXT_SERVER_CORS_ALLOWED_ORIGINS` não configurado | Adicione URL do frontend (sem `/` no final) |
| **502 Bad Gateway** | Backend não iniciou | Verifique logs e variáveis de ambiente |
| **Cannot connect to database** | Variáveis `FINEXT_DATABASE_*` incorretas | Verifique se MySQL addon está conectado |
| **Firebase auth failed** | `FINEXT_FIREBASE_SERVICE_ACCOUNT_JSON` inválido | Copie JSON completo do Firebase Console |
| **Frontend 404 on API** | `NEXT_PUBLIC_API_URL` sem `/api/v1` | Adicione `/api/v1` no final da URL |
| **Build failed** | Root Directory errado | Configure `/backend` ou `/frontend` |

### 🔧 Comandos Úteis

```bash
# Gerar JWT Secret
openssl rand -base64 32

# Ver variáveis configuradas
railway variables --service backend
railway variables --service frontend

# Ver logs em tempo real
railway logs --service backend
railway logs --service frontend

# Forçar redeploy
railway up --service backend --detach
railway up --service frontend --detach
```

---

## 📞 Suporte

- **Railway Docs**: https://docs.railway.app
- **Railway Discord**: https://discord.gg/railway
- **Logs de Deploy**: Railway Dashboard → Service → Deployments → View Logs
