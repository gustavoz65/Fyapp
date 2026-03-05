# 📚 Fy - Documentação de Produto

Bem-vindo à documentação estratégica do Fy! Aqui você encontra tudo sobre produto, strategy, UX, conversão e negócios.

---

## 📑 Índice de Documentos

### 💎 [PRODUCT_STRATEGY.md](./PRODUCT_STRATEGY.md)

**Estratégia completa de produto, pricing e conversão**

Contém:

- ✅ Estratégia de Pricing (Free, Pro, Ultra)
- ✅ Projeções Financeiras (Mês 3, 6, 9, 12)
- ✅ Alavancadores de Conversão (Free→Pro, Pro→Ultra)
- ✅ Onboarding Strategy (maximiza conversão)
- ✅ Funil de Conversão Completo
- ✅ Integração Open Finance (quando implementar)
- ✅ Pitch para Pluggy (negociação)
- ✅ Unit Economics (LTV, CAC, Margem)
- ✅ Métricas de Sucesso (KPIs)

**Recomendado para:**

- Entender modelo de negócio
- Preparar pitch para investidores/parceiros
- Planejar roadmap de produto
- Definir estratégia de pricing

---

### 🧠 [UX_PSYCHOLOGY.md](./UX_PSYCHOLOGY.md)

**Psicologia comportamental aplicada a design de produto**

Contém:

- ✅ 8 Princípios de Persuasão (Cialdini + Fogg)
- ✅ Como empresas bilionárias fazem (Netflix, Spotify, Notion, Duolingo)
- ✅ Behavior Model (Motivação × Ability × Trigger)
- ✅ Dark Patterns (o que NÃO fazer)
- ✅ Copywriting que Converte (PAS, AIDA)
- ✅ A/B Testing (otimização contínua)
- ✅ Framework prático de implementação
- ✅ Tools e Recursos

**Recomendado para:**

- Entender por que UX decisions funcionam
- Otimizar conversão
- Criar copy persuasivo
- Evitar dark patterns

---

## 🎯 Quick Reference

### Pricing Atual

| Plano | Mensal   | Anual     | Open Finance |
| ----- | -------- | --------- | ------------ |
| Free  | R$ 0     | R$ 0      | ❌           |
| Pro   | R$ 59,99 | R$ 499,90 | ✅           |
| Ultra | R$ 99,99 | R$ 799,90 | ✅           |

---

### Projeções (MRR)

| Mês    | Meta MRR  | Status       |
| ------ | --------- | ------------ |
| Mês 3  | R$ 1.500  | Beta         |
| Mês 6  | R$ 6.500  | Cobre Pluggy |
| Mês 9  | R$ 16.600 | Crescimento  |
| Mês 12 | R$ 42.000 | Escala       |

---

### Conversão Esperada

| Métrica               | Meta   | Benchmark |
| --------------------- | ------ | --------- |
| Onboarding Completion | 80%    | 40-60%    |
| Free→Paid             | 20-25% | 2-5%      |
| Trial→Pro             | 40-50% | 25-40%    |
| Pro→Ultra             | 25-30% | 15-20%    |

---

### Unit Economics

**LTV (18 meses):**

- Pro: R$ 1.080
- Ultra: R$ 1.800

**CAC:**

- Orgânico: R$ 0-20
- Conteúdo: R$ 30-50
- Ads: R$ 80-120

**LTV/CAC Ratio:**

- Orgânico: 54x ✅
- Ads: 10x ✅
- Benchmark: >3x

**Margem Bruta:** 85% (benchmark: 70-80%)

---

## 🚀 Roadmap de Implementação

### Fase 1: Beta SEM Open Finance (Mês 1-2)

**Objetivo:** Validar produto e conversão

**Tarefas:**

- [ ] Deploy na Railway
- [ ] Onboarding simplificado (4 etapas)
- [ ] Sistema de trials (7 dias grátis)
- [ ] Gamificação ("Complete perfil")
- [ ] Setup analytics
- [ ] Convidar 50 early adopters

**Meta:** 50 usuários, 10 pagantes (20% conversão)

---

### Fase 2: Decisão Open Finance (Mês 3)

**Objetivo:** Validar se vale a pena Pluggy

**Tarefas:**

- [ ] Medir conversão real
- [ ] Coletar feedback usuários
- [ ] Calcular LTV real
- [ ] **DECISÃO:** Implementar Open Finance?
  - Se SIM → Negociar com Pluggy
  - Se NÃO → Continuar manual, focar growth

**Meta:** 100 usuários, 20 pagantes, dados reais

---

### Fase 3: Escala (Mês 4+)

**Objetivo:** Crescer para 1.000+ usuários

**Tarefas:**

- [ ] Implementar Open Finance (se validado)
- [ ] Trial de 7 dias Open Finance
- [ ] Otimizar conversão
- [ ] Marketing (SEO, conteúdo, ads)
- [ ] Atingir R$ 16k+ MRR

**Meta:** 1.000+ usuários, 200+ pagantes

---

## 📊 Métricas Chave (KPIs)

### Métricas Primárias (North Star)

1. **MRR** (Monthly Recurring Revenue)
   - O que mede: Receita recorrente mensal
   - Por que importa: Saúde do negócio
   - Meta Mês 6: R$ 6.500

2. **Conversão Free→Paid**
   - O que mede: % de usuários que viram pagantes
   - Por que importa: Eficiência de aquisição
   - Meta: 20-25%

3. **Trial→Pro Conversion**
   - O que mede: % de trials que viram pagantes
   - Por que importa: ROI do Open Finance
   - Meta: 40-50%

---

### Métricas Secundárias

4. **Onboarding Completion Rate**
   - Meta: 80% (benchmark: 40-60%)

5. **DAU/MAU Ratio**
   - Meta: 30%+ (benchmark: 20%)

6. **Churn Rate**
   - Meta: <5%/mês (benchmark: 5-7%)

7. **NPS** (Net Promoter Score)
   - Meta: 50+ (benchmark: 30-40)

---

## 🧪 Frameworks Usados

### 1. BJ Fogg Behavior Model

```
Behavior = Motivation × Ability × Trigger
```

**Aplicação:**

- Motivação: Benefícios claros ("economize 10 min/dia")
- Ability: Fácil de usar (1 clique, sem fricção)
- Trigger: Momento certo (após 10 transações manuais)

---

### 2. Hook Model (Nir Eyal)

```
Trigger → Action → Variable Reward → Investment
```

**Aplicação:**

- Trigger: "Você adicionou 10 transações manualmente"
- Action: "Experimentar Pro 7 dias grátis"
- Reward: "IA categorizou 47 transações para você!"
- Investment: Dados no sistema, histórico, metas criadas

---

### 3. Jobs to be Done

**Job principal do Fy:**

> "Quando eu recebo meu salário, eu quero saber se vai sobrar no final do mês, para que eu possa dormir tranquilo sem medo de passar aperto."

**Jobs secundários:**

- Economizar para objetivo (viagem, carro)
- Sair das dívidas
- Entender para onde vai o dinheiro
- Parar de estourar o cartão

---

## 🎯 Decisões de Produto Chave

### 1. Onboarding Simplificado (4 etapas vs 7)

**Decisão:** Reduzir de 7 para 4 etapas obrigatórias

**Razão:**

- Menos fricção = +80% completion (vs 40% com 7 etapas)
- Quick win mais rápido (<2 min)
- Resto via gamificação pós-onboarding

**Trade-off:**

- ❌ Menos dados iniciais
- ✅ Mais usuários completam e veem valor

**Vencedor:** Simplicidade

---

### 2. Trial SEM pedir cartão

**Decisão:** Trial de 7 dias sem pedir cartão de crédito

**Razão:**

- Reciprocidade (dá primeiro, pede depois)
- Baixa fricção = +40% aceitação
- Conversão trial→paid: 40-50% vs 25% com cartão

**Trade-off:**

- ❌ Risco de "trial abuse" (usuários fazem múltiplos trials)
- ✅ Conversão muito maior compensa

**Vencedor:** Trial sem cartão

---

### 3. Open Finance como Trial (não free tier)

**Decisão:** Open Finance só em trial Pro (7 dias), não no Free

**Razão:**

- Open Finance custa R$ 2.500/mês (Pluggy)
- Não podemos dar de graça
- Trial = taste do valor, converte bem

**Trade-off:**

- ❌ Menos pessoas experimentam Open Finance
- ✅ Só quem tem intenção real experimenta (menos custo desperdiçado)

**Vencedor:** Trial pago

---

### 4. Preço Premium (R$59,99 vs R$19,99)

**Decisão:** Posicionar como premium (R$59,99) vs low-cost (R$19,99)

**Razão:**

- Premium permite margem alta (85%)
- Valor entregue justifica (Open Finance, IA)
- Low-cost atrai tire-kickers (não pagam anyway)

**Trade-off:**

- ❌ Menos volume de usuários
- ✅ Mais receita por usuário (LTV R$1.080 vs R$340)

**Vencedor:** Premium

---

## 🔗 Links Úteis

### Internal Docs

- [Product Strategy](./PRODUCT_STRATEGY.md)
- [UX Psychology](./UX_PSYCHOLOGY.md)
- [Railway Deployment](../.railway/README.md)
- [API Documentation](../backend/README.md)

### External Resources

- [BJ Fogg Behavior Model](https://behaviormodel.org/)
- [Hooked Model](https://www.nirandfar.com/hooked/)
- [Growth Loops (Reforge)](https://www.reforge.com/growth-loops)
- [Cialdini's 6 Principles](https://www.influenceatwork.com/)

---

## 📞 Contato

**Product Questions:** gustavo@Fy.com.br

**Última atualização:** 2026-02-27
