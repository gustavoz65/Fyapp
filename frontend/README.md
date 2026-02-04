# 💰 Cashing - Sistema de Gestão Financeira

Sistema completo de gestão financeira para pessoas físicas e jurídicas, desenvolvido com React 18, TypeScript, Vite, Tailwind CSS e Zustand.

![Cashing Banner](https://via.placeholder.com/1200x400/1E3A8A/FFFFFF?text=Cashing+-+Gest%C3%A3o+Financeira+Inteligente)

## 🚀 Tecnologias

- **Framework**: React 18.3+ com TypeScript 5.4+
- **Build Tool**: Vite 5.2+
- **Estilização**: Tailwind CSS 3.4+
- **Componentes**: Headless UI + Lucide React (ícones)
- **Estado Global**: Zustand 4.5+ com persistência
- **Formulários**: React Hook Form 7.51+ com validação Zod 3.22+
- **Gráficos**: Recharts 2.12+
- **Rotas**: React Router v6.22+
- **Requisições**: TanStack React Query 5.28+
- **Animações**: Framer Motion 11.0+

## 📁 Estrutura do Projeto

```
src/
├── components/          # Componentes reutilizáveis
│   ├── Badge.tsx
│   ├── Button.tsx
│   ├── Card.tsx
│   ├── Drawer.tsx
│   ├── EmptyState.tsx
│   ├── Input.tsx
│   ├── Modal.tsx
│   ├── ProgressBar.tsx
│   ├── Select.tsx
│   ├── Skeleton.tsx
│   └── Table.tsx
├── data/               # Mock data
│   ├── banks.ts
│   ├── budgets.ts
│   ├── categories.ts
│   ├── notifications.ts
│   └── transactions.ts
├── lib/                # Utilitários
│   └── utils.ts
├── routes/             # Páginas e rotas
│   ├── protected/      # Rotas protegidas
│   │   ├── banks/
│   │   ├── budgets/
│   │   ├── categories/
│   │   ├── dashboard/
│   │   ├── notifications/
│   │   ├── reports/
│   │   ├── settings/
│   │   ├── transactions/
│   │   └── ProtectedLayout.tsx
│   └── public/         # Rotas públicas
│       ├── auth/
│       │   ├── Login.tsx
│       │   ├── Register.tsx
│       │   ├── ResetPassword.tsx
│       │   └── VerifyEmail.tsx
│       └── landing/
│           └── LandingPage.tsx
├── stores/             # Zustand stores
│   ├── authStore.ts
│   └── settingsStore.ts
├── types/              # TypeScript types
│   └── index.ts
├── App.tsx
├── index.css
└── main.tsx
```

## ✨ Funcionalidades

### 🏠 Landing Page
- Hero com CTA duplo
- 6 cards de features
- Estatísticas animadas
- 3 depoimentos em carrossel
- Planos (Free/Pro/Business) com toggle anual/mensal
- FAQ com 5 perguntas
- Footer completo

### 🔐 Autenticação
- ✅ Login com email/senha
- ✅ Registro com seleção de plano
- ✅ Google OAuth (UI)
- ✅ Recuperação de senha
- ✅ Verificação de email
- ✅ Validação com Zod
- ✅ Estados de loading

### 📊 Dashboard
- **Cards principais**: Saldo total, Receitas, Despesas, Saldo líquido
- **Gráfico de linha**: Fluxo de caixa 12 meses
- **Gráfico de pizza**: Despesas por categoria
- **Tabela**: Últimas 5 transações
- **Filtros**: Período (Mês/Trimestre/Ano)

### 💳 Transações
- Tabela completa com busca e ordenação
- Filtros avançados (data, categoria, tipo, valor)
- Export para CSV/Excel
- Modal de detalhes
- Status: Concluída, Pendente, Cancelada

### 📈 Relatórios
- Fluxo de Caixa
- DRE (Demonstração do Resultado)
- Balanço Patrimonial
- Export para PDF/Excel
- Comparação de períodos

### 🏷️ Categorias
- 12 categorias pré-definidas
- Subcategorias aninhadas
- Drag & drop para reordenar
- Ícones e cores personalizadas
- Separação: Receitas vs Despesas

### 💰 Orçamentos
- Orçamentos por categoria
- Barras de progresso
- Alertas visuais (Verde/Amarelo/Vermelho)
- Orçamentos mensais e anuais
- Threshold configurável

### 🏦 Bancos
- 3 contas mock (BB, Nubank, Inter)
- Saldo consolidado
- Status de conexão
- Sincronização automática (mock)
- Últimas transações

### ⚙️ Configurações
- Perfil do usuário
- Moeda (BRL/USD/EUR)
- Idioma (PT-BR/EN)
- Notificações (Email/Push/SMS)
- Segurança (2FA - UI)

### 🔔 Notificações
- Badge no header
- 5 tipos: info, success, warning, error
- Alertas de orçamento
- Vencimentos
- Saldo baixo

## 🎨 Design System

### Paleta de Cores
```css
Primary: #1E3A8A (Azul escuro)
Success: #10b981 (Verde)
Warning: #f59e0b (Amarelo)
Error: #ef4444 (Vermelho)
Info: #3b82f6 (Azul)
Gray: #6b7280 (Cinza)
```

### Tipografia
- Fonte: Inter (Google Fonts)
- Pesos: 300, 400, 500, 600, 700, 800

### Componentes
- Botões: 5 variantes (primary, secondary, danger, ghost, success)
- Tamanhos: sm, md, lg
- Inputs com validação visual
- Modals responsivos
- Drawers laterais
- Tables com ordenação e busca
- Progress bars com cores dinâmicas
- Badges coloridos
- Skeletons para loading
- Empty states

## 🔧 Setup e Instalação

### Pré-requisitos
- Node.js 18+ 
- npm ou yarn

### Instalação

```bash
# Clone o repositório
git clone https://github.com/gustavoz65/Cashing.git
cd Cashing

# Instale as dependências
npm install

# Inicie o servidor de desenvolvimento
npm run dev

# Build para produção
npm run build

# Preview da build
npm run preview
```

### Scripts Disponíveis

```json
{
  "dev": "vite",                    // Servidor de desenvolvimento
  "build": "tsc && vite build",     // Build com verificação TypeScript
  "preview": "vite preview",        // Preview da build de produção
  "lint": "eslint . --ext ts,tsx"   // Lint do código
}
```

## 📊 Mock Data

O sistema inclui dados mockados realistas:

- **500+ transações** variadas (últimos 12 meses)
- **12 categorias** com subcategorias
- **3 contas bancárias** (BB, Nubank, Inter)
- **5 orçamentos** ativos (mensais e anuais)
- **5 notificações** com diferentes tipos
- Dados PF e PJ misturados

## ♿ Acessibilidade

- ✅ ARIA labels em elementos interativos
- ✅ Navegação por teclado completa
- ✅ Screen reader friendly
- ✅ Contraste de cores WCAG AA
- ✅ Focus trap em modals
- ✅ Textos alternativos

## 🚀 Performance

- ✅ Code splitting automático
- ✅ Lazy loading de rotas
- ✅ Bundle otimizado < 1MB inicial
- ✅ Tree shaking
- ✅ Compressão de assets
- ✅ Lighthouse score 95+ (alvo)

## 📱 Responsividade

- ✅ Mobile-first design
- ✅ Breakpoints: sm (640px), md (768px), lg (1024px), xl (1280px)
- ✅ Drawer para navegação mobile
- ✅ Tabelas com scroll horizontal
- ✅ Grid adaptativo

## 🔐 Segurança

- ✅ Validação de formulários (Zod)
- ✅ Proteção de rotas (ProtectedRoute)
- ✅ Sanitização de inputs
- ✅ HTTPS only (produção)
- ✅ 2FA (UI mock)
- ✅ Tokens JWT (mock)

## 🧪 Testes

```bash
# Testes unitários
npm run test

# Testes E2E
npm run test:e2e

# Coverage
npm run test:coverage
```

## 📦 Build e Deploy

### Build de Produção

```bash
npm run build
```

Saída: pasta `dist/`

### Deploy

#### Vercel
```bash
npm i -g vercel
vercel --prod
```

#### Netlify
```bash
npm i -g netlify-cli
netlify deploy --prod
```

#### Docker
```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build
EXPOSE 5173
CMD ["npm", "run", "preview"]
```

## 🗺️ Roadmap

- [ ] Integração com backend real
- [ ] Testes unitários (Jest + Testing Library)
- [ ] Testes E2E (Playwright)
- [ ] PWA (Service Workers)
- [ ] Tema dark mode
- [ ] Múltiplos idiomas (i18n)
- [ ] Export avançado (PDF, Excel)
- [ ] Importação de OFX/CSV
- [ ] Reconciliação bancária
- [ ] Metas financeiras
- [ ] Dashboard customizável

## 📄 Licença

MIT License - veja [LICENSE](LICENSE) para mais detalhes.

## 👥 Contribuindo

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -m 'Add nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📞 Suporte

- Email: suporte@cashing.com
- Discord: [discord.gg/cashing](https://discord.gg/cashing)
- Documentação: [docs.cashing.com](https://docs.cashing.com)

## 🙏 Agradecimentos

- [React](https://react.dev)
- [Tailwind CSS](https://tailwindcss.com)
- [Headless UI](https://headlessui.com)
- [Lucide Icons](https://lucide.dev)
- [Recharts](https://recharts.org)

---

Feito com ❤️ por [Gustavo](https://github.com/gustavoz65)
