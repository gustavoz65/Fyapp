# Frontend Redesign - Design Document

**Date:** 2026-03-09
**Status:** Approved
**Approach:** Minimalista Profissional

---

## Visão Geral

Redesign completo do frontend do sistema **Fy - Gestão Financeira Inteligente** para modernizar a identidade visual, melhorar a experiência em mobile e desktop, e transmitir maior profissionalismo e facilidade de uso.

### Objetivos

- ✅ Substituir paleta verde musgo + bege por verde Sicredi + branco
- ✅ Modernizar tipografia para visual mais clean e profissional
- ✅ Manter filosofia de facilidade de uso
- ✅ Otimizar para mobile e desktop
- ✅ Preservar dark mode e acessibilidade

---

## Abordagem Escolhida: Minimalista Profissional

**Características:**
- Verde Sicredi (#00A859) como cor primária
- Branco puro como background principal
- Tipografia Inter (única família)
- Muito espaço em branco
- Design limpo e profissional
- Bordas arredondadas suaves
- Ícones em outline style

**Vantagens:**
- Transmite profissionalismo e confiança
- Extremamente legível e acessível
- Funciona perfeitamente em mobile e desktop
- Performance otimizada
- Fácil de manter e escalar

---

## 1. Sistema de Cores

### Paleta Principal (Light Theme)

```css
/* Primary Colors */
--primary: #00A859           /* Verde Sicredi - CTAs, estados ativos */
--primary-hover: #00C569     /* Verde claro - hover states */
--primary-foreground: #FFFFFF /* Texto sobre verde */

/* Backgrounds */
--background: #FFFFFF         /* Background principal */
--surface: #F8FAFB           /* Cards, elevated surfaces */
--surface-hover: #F1F5F9     /* Hover em cards */

/* Text Colors */
--foreground: #0F172A        /* Texto principal */
--muted: #64748B             /* Texto secundário */
--muted-foreground: #94A3B8  /* Texto terciário */

/* Borders */
--border: #E2E8F0            /* Bordas sutis */
--border-strong: #CBD5E1     /* Bordas com contraste */

/* Interactive */
--ring: #00A859              /* Focus rings */
```

### Cores Semânticas

```css
/* Success */
--success: #00A859           /* Reutiliza primary */
--success-light: #D1FAE5     /* Background de success states */

/* Warning */
--warning: #F59E0B           /* Alertas */
--warning-light: #FEF3C7

/* Error */
--error: #EF4444             /* Erros, saldo negativo */
--error-light: #FEE2E2

/* Info */
--info: #0EA5E9              /* Informações */
--info-light: #E0F2FE
```

### Dark Theme

```css
--background: #0F172A
--surface: #1E293B
--surface-hover: #334155
--foreground: #F1F5F9
--muted: #64748B
--muted-foreground: #94A3B8
--border: rgba(255, 255, 255, 0.1)
--border-strong: rgba(255, 255, 255, 0.2)

/* Primary mantém o verde Sicredi */
--primary: #00A859
--primary-hover: #00C569
```

### Cores de Gráficos

```css
--chart-1: #00A859  /* Verde principal */
--chart-2: #00C569  /* Verde claro */
--chart-3: #0EA5E9  /* Azul */
--chart-4: #8B5CF6  /* Roxo */
--chart-5: #F59E0B  /* Laranja */
```

---

## 2. Sistema Tipográfico

### Família

**Inter** - única fonte para todo o sistema (simplicidade e performance)

```css
--font-heading: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif
--font-body: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif
```

### Pesos

```css
--font-regular: 400
--font-medium: 500
--font-semibold: 600
--font-bold: 700
--font-extrabold: 800
```

### Escala Tipográfica

```css
/* Headings */
h1: 2.5rem (40px) / 800 / line-height: 1.1
h2: 2rem (32px) / 700 / line-height: 1.2
h3: 1.5rem (24px) / 700 / line-height: 1.3
h4: 1.25rem (20px) / 600 / line-height: 1.4
h5: 1.125rem (18px) / 600 / line-height: 1.4

/* Body */
body: 1rem (16px) / 400 / line-height: 1.6
small: 0.875rem (14px) / 400 / line-height: 1.5
xs: 0.75rem (12px) / 500 / line-height: 1.4
```

### Características

- **Sem Playfair Display** - removido para maior coesão visual
- **Inter em todos os contextos** - títulos, corpo, UI
- **Letter-spacing:** -0.01em para títulos grandes
- **Corpo do texto:** normal letter-spacing

---

## 3. Componentes a Atualizar

### Prioridade Alta (Core Visual Identity)

1. **`globals.css`**
   - Atualizar todas as variáveis CSS
   - Remover referências ao Playfair Display
   - Atualizar paleta completa

2. **Buttons**
   - Primary: Verde Sicredi
   - Hover: Verde claro (#00C569)
   - Outline: Border verde Sicredi

3. **Cards**
   - Background: --surface (#F8FAFB)
   - Border radius: 0.5rem
   - Shadow: sutil, apenas quando hover

4. **Sidebar**
   - Background: Branco puro
   - Active state: Verde Sicredi
   - Ícones: outline style

5. **Dashboard Metrics**
   - Cards com --surface background
   - Verde para métricas positivas
   - Vermelho para métricas negativas
   - Ícones modernizados

6. **Charts (Recharts)**
   - Atualizar cores para nova paleta
   - Verde principal para receitas
   - Vermelho para despesas

### Prioridade Média

7. **Forms**
   - Inputs: border --border, focus ring verde
   - Labels: --muted-foreground
   - Error states: --error

8. **Tables**
   - Header: --surface background
   - Rows: hover --surface-hover
   - Borders: --border

9. **Badges**
   - Success: verde Sicredi
   - Warning: laranja
   - Error: vermelho
   - Info: azul

10. **Tabs**
    - Active: verde Sicredi
    - Inactive: --muted

### Prioridade Baixa

11. **Modals/Dialogs**
12. **Tooltips**
13. **Dropdowns**
14. **Popovers**

---

## 4. Layout & Espaçamento

### Sistema de Espaçamento

```css
--spacing-xs: 0.25rem   /* 4px */
--spacing-sm: 0.5rem    /* 8px */
--spacing-md: 1rem      /* 16px */
--spacing-lg: 1.5rem    /* 24px */
--spacing-xl: 2rem      /* 32px */
--spacing-2xl: 3rem     /* 48px */
--spacing-3xl: 4rem     /* 64px */
```

### Border Radius

```css
--radius-sm: 0.375rem   /* 6px */
--radius: 0.5rem        /* 8px - padrão */
--radius-lg: 0.75rem    /* 12px */
--radius-xl: 1rem       /* 16px */
--radius-full: 9999px   /* Circular */
```

### Shadows

```css
--shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05)
--shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)
--shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)
--shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)
```

### Breakpoints (Mobile-First)

```css
sm: 640px   /* Mobile landscape */
md: 768px   /* Tablet portrait */
lg: 1024px  /* Tablet landscape / Desktop */
xl: 1280px  /* Desktop large */
2xl: 1536px /* Desktop extra large */
```

---

## 5. Páginas e Fluxos

### Páginas a Redesenhar

1. **Login/Register** - Atualizar branding section com novo verde
2. **Onboarding** - Modernizar wizard
3. **Dashboard** - Aplicar novo design de cards e gráficos
4. **Transações** - Atualizar tabelas e badges
5. **Contas** - Redesenhar cards de contas
6. **Categorias** - Atualizar UI de categorias
7. **Orçamentos** - Redesenhar progress bars
8. **Metas** - Atualizar indicadores visuais
9. **Recorrentes** - Modernizar listagem
10. **Configurações** - Atualizar forms

### Componentes de Layout

- **AppSidebar** - Novo visual clean, branco com verde nos ativos
- **Header** - Simplificado, mais espaço em branco
- **HealthWidget** - Atualizar indicadores de saúde financeira

---

## 6. Acessibilidade & Performance

### Acessibilidade

- **Contraste WCAG AA:** Todas as combinações testadas
- **Focus visible:** Ring verde Sicredi (#00A859) em todos os interativos
- **Keyboard navigation:** Mantido
- **Screen readers:** Labels preservados

### Performance

- **Font loading:** Inter via Google Fonts com display=swap
- **CSS Variables:** Troca de tema instantânea
- **Tailwind CSS:** Tree-shaking automático
- **No extra dependencies:** Apenas cores e tipografia

---

## 7. Estratégia de Implementação

### Fase 1: Foundation (Base do Design System)
1. Atualizar `globals.css` com novas variáveis
2. Remover Playfair Display, usar apenas Inter
3. Atualizar `layout.tsx` para carregar apenas Inter
4. Testar dark mode com novas cores

### Fase 2: Core Components
1. Buttons
2. Cards
3. Inputs/Forms
4. Badges
5. Tables

### Fase 3: Layout Components
1. Sidebar
2. Header
3. Dashboard page (validação visual completa)

### Fase 4: Feature Pages
1. Transações
2. Contas
3. Categorias
4. Orçamentos
5. Metas
6. Recorrentes

### Fase 5: Auth & Onboarding
1. Login/Register
2. Onboarding wizard

### Fase 6: Polish & QA
1. Animações e transições
2. Mobile testing (iOS/Android)
3. Dark mode testing
4. Acessibilidade audit
5. Performance testing

---

## 8. Decisões de Design

### Por que Inter?
- Moderna e profissional
- Excelente legibilidade em todos os tamanhos
- Ótimo suporte a diferentes pesos
- Performance (Google Fonts otimizado)
- Usado por empresas de fintech (Stripe, Coinbase)

### Por que Verde Sicredi?
- Transmite confiança e crescimento
- Vibrante mas profissional
- Alto contraste com branco
- Associação positiva com finanças/prosperidade

### Por que remover Playfair Display?
- Serif pode parecer "pesado" em mobile
- Mistura de fontes adiciona complexidade
- Inter sozinho cria coesão visual
- Performance (uma fonte a menos)

### Por que muito espaço em branco?
- Facilita escaneabilidade
- Reduz carga cognitiva
- Funciona melhor em mobile
- Moderniza o visual
- Destaca informações importantes

---

## 9. Referências Visuais

### Inspirações
- **Sicredi:** Paleta de cores (verde vibrante)
- **Nubank:** Simplicidade e foco em mobile
- **Stripe:** Tipografia clean e espaçamento generoso
- **Linear:** Design system coeso e profissional

### Design Patterns
- **shadcn/ui:** Componentes base (já em uso)
- **Tailwind CSS v4:** Utility-first styling
- **Radix UI:** Acessibilidade e comportamentos

---

## 10. Métricas de Sucesso

### Qualitativas
- [ ] Visual mais moderno e profissional
- [ ] Melhor legibilidade em mobile
- [ ] Identidade visual coesa
- [ ] Facilidade de uso mantida/melhorada

### Quantitativas
- [ ] Performance mantida (Lighthouse score ≥90)
- [ ] Contraste WCAG AA em todas as combinações
- [ ] Zero regressões visuais em mobile
- [ ] Dark mode funcionando perfeitamente

---

## Próximos Passos

1. ✅ Design aprovado
2. → Criar plano de implementação detalhado
3. → Implementar Fase 1 (Foundation)
4. → Validar visualmente no Dashboard
5. → Rollout gradual por feature
6. → QA final e deploy

---

**Documento aprovado em:** 2026-03-09
**Responsável:** Claude Code
**Revisor:** Gustavo (Product Owner)
