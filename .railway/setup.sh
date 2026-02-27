#!/bin/bash

# 🚂 Railway Setup Script
# Facilita a configuração inicial do projeto na Railway

set -e

echo "🚂 Railway Setup - Finext Beta 1"
echo "=================================="
echo ""

# Verificar se Railway CLI está instalado
if ! command -v railway &> /dev/null; then
    echo "❌ Railway CLI não está instalado!"
    echo ""
    echo "Instale com:"
    echo "  npm install -g @railway/cli"
    echo ""
    exit 1
fi

echo "✅ Railway CLI encontrado"
echo ""

# Login
echo "📝 Fazendo login na Railway..."
railway login

echo ""
echo "🔗 Linkando ao projeto Railway..."
echo "   (Selecione o projeto 'finext' quando solicitado)"
railway link

echo ""
echo "📊 Status do projeto:"
railway status

echo ""
echo "🎯 Service IDs:"
echo "=============="
echo ""
echo "Copie os IDs abaixo e adicione como Secrets no GitHub:"
echo "(GitHub → Settings → Secrets → Actions)"
echo ""

# Pegar service IDs
BACKEND_ID=$(railway status | grep -i "backend" | awk '{print $3}' | tr -d '()')
FRONTEND_ID=$(railway status | grep -i "frontend" | awk '{print $3}' | tr -d '()')

if [ -n "$BACKEND_ID" ]; then
    echo "RAILWAY_SERVICE_BACKEND=$BACKEND_ID"
else
    echo "⚠️  Service 'backend' não encontrado"
fi

if [ -n "$FRONTEND_ID" ]; then
    echo "RAILWAY_SERVICE_FRONTEND=$FRONTEND_ID"
else
    echo "⚠️  Service 'frontend' não encontrado"
fi

echo ""
echo "🔑 Railway Token:"
echo "================"
echo "1. Vá em: https://railway.app/account/tokens"
echo "2. Crie um novo token"
echo "3. Adicione como secret no GitHub:"
echo "   RAILWAY_TOKEN=<seu-token>"
echo ""

echo "✅ Setup completo!"
echo ""
echo "Próximos passos:"
echo "1. Adicionar secrets no GitHub (IDs e token acima)"
echo "2. Configurar variáveis de ambiente na Railway UI"
echo "3. Push para main → Deploy automático via GitHub Actions"
echo ""
echo "Documentação completa: .railway/README.md"
