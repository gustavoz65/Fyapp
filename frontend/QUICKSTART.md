# Cashing - Quick Start Guide

## 🚀 Início Rápido

### 1. Instalação

```bash
npm install
```

### 2. Desenvolvimento

```bash
npm run dev
```

Acesse: `http://localhost:5173`

### 3. Credenciais de Teste

Para acessar o sistema, use qualquer email e senha (mock):

- **Email**: teste@exemplo.com
- **Senha**: 123456

## 📋 Comandos Disponíveis

```bash
npm run dev      # Inicia servidor de desenvolvimento
npm run build    # Build para produção
npm run preview  # Preview da build
npm run lint     # Lint do código
```

## 🎯 Navegação Rápida

### Páginas Públicas
- `/` - Landing Page
- `/auth/login` - Login
- `/auth/register` - Registro
- `/auth/reset-password` - Recuperar senha

### Páginas Protegidas (após login)
- `/dashboard` - Dashboard principal
- `/dashboard/transactions` - Lista de transações
- `/dashboard/reports` - Relatórios financeiros
- `/dashboard/categories` - Gestão de categorias
- `/dashboard/budgets` - Orçamentos
- `/dashboard/banks` - Contas bancárias
- `/dashboard/settings` - Configurações
- `/dashboard/notifications` - Notificações

## 🎨 Design Tokens

### Cores Principais
```css
Primary: #1E3A8A
Success: #10b981
Warning: #f59e0b
Error: #ef4444
Info: #3b82f6
```

### Breakpoints
```css
sm: 640px
md: 768px
lg: 1024px
xl: 1280px
```

## 📦 Componentes Disponíveis

- `Button` - Botões (5 variantes, 3 tamanhos)
- `Card` - Cards com header/footer
- `Modal` - Modais responsivos
- `Drawer` - Drawers laterais
- `Table` - Tabelas com ordenação e busca
- `Input` - Inputs com validação
- `Select` - Selects estilizados
- `Badge` - Badges coloridos
- `ProgressBar` - Barras de progresso
- `Skeleton` - Loading states
- `EmptyState` - Estados vazios

## 🔧 Estrutura de Dados

### Mock Data Disponível

#### Transações (500+)
```typescript
import { transactions } from '@/data/transactions'
```

#### Categorias (12)
```typescript
import { categories } from '@/data/categories'
```

#### Bancos (3)
```typescript
import { banks } from '@/data/banks'
```

#### Orçamentos (5)
```typescript
import { budgets } from '@/data/budgets'
```

#### Notificações (5)
```typescript
import { notifications } from '@/data/notifications'
```

## 🛠️ Customização

### Adicionar Nova Cor

1. Edite `tailwind.config.js`:
```javascript
theme: {
  extend: {
    colors: {
      custom: {
        500: '#YOUR_COLOR',
      },
    },
  },
}
```

2. Use: `bg-custom-500`, `text-custom-500`, etc.

### Adicionar Nova Rota

1. Crie o componente em `src/routes/`
2. Adicione em `src/App.tsx`:
```tsx
<Route path="/nova-rota" element={<NovoComponente />} />
```

### Adicionar Novo Store Zustand

```typescript
import { create } from 'zustand'

export const useMyStore = create((set) => ({
  value: 0,
  setValue: (value: number) => set({ value }),
}))
```

## 🐛 Troubleshooting

### Erro de dependências
```bash
rm -rf node_modules package-lock.json
npm install
```

### Erro de build
```bash
npm run lint
npx tsc --noEmit
```

### Porta em uso
```bash
# Altere a porta em vite.config.ts
server: {
  port: 3000,
}
```

## 📚 Recursos Úteis

- [React Docs](https://react.dev)
- [TypeScript Docs](https://www.typescriptlang.org/docs/)
- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [Zustand Docs](https://docs.pmnd.rs/zustand)
- [React Router Docs](https://reactrouter.com)
- [Recharts Docs](https://recharts.org)

## 📝 Próximos Passos

1. ✅ Explorar a Landing Page
2. ✅ Fazer login no sistema
3. ✅ Navegar pelo dashboard
4. ✅ Visualizar transações e relatórios
5. ✅ Testar criação de categorias e orçamentos
6. ⏳ Integrar com backend real
7. ⏳ Adicionar testes
8. ⏳ Deploy em produção

## 💡 Dicas

- Use `Ctrl + K` para busca rápida (a implementar)
- Todas as rotas protegidas requerem autenticação
- Mock data é resetado ao recarregar a página
- Estado persiste em localStorage (auth, settings)

## 🤝 Suporte

Encontrou um bug ou tem uma sugestão?
- Abra uma issue no GitHub
- Entre em contato: suporte@cashing.com

---

Desenvolvido com ❤️ usando React + TypeScript
