# 🧠 UX Psychology - The Dark Art of Product Design

> Como usar psicologia comportamental para criar produtos viciantes (de forma ética)

**TL;DR:** UX não é só "deixar bonito" - é **psicologia aplicada para fazer o usuário QUERER o que você tá vendendo**. Este documento revela os segredos sujos do design de produto.

---

## 🎭 A Real sobre UX

### UX é estratégia de negócio disfarçada de design

**Toda decisão tem motivo psicológico:**

| Decisão de Design                    | Psicologia Aplicada | Resultado                            |
| ------------------------------------ | ------------------- | ------------------------------------ |
| Trial 7 dias grátis SEM pedir cartão | Reciprocidade       | 40-50% conversão vs 5-10% com cartão |
| Progress bar no onboarding           | Efeito Zeigarnik    | 80% completion vs 40% sem            |
| "Recomendado" na opção               | Viés de autoridade  | 60% escolhem opção marcada           |
| "Complete perfil 65%"                | Loop de dopamina    | +40% engagement                      |
| "Experimente GRÁTIS"                 | Aversão à perda     | Baixa taxa de cancelamento           |

---

## 💰 Como Empresas Bilionárias Fazem Isso

### Netflix

**❌ Ruim:** "Assine agora por R$39,90/mês"

**✅ Netflix:**

```
1 MÊS GRÁTIS (letras garrafais)
cancele quando quiser (letra pequena)
```

**Psicologia:** Trial grátis + facilidade de cancelar = baixa fricção
**Resultado:** Após 1 mês assistindo, usuário não cancela (aversão à perda)

---

### Duolingo

**❌ Ruim:** "Aprenda inglês"

**✅ Duolingo:**

- Streak de 47 dias 🔥 (vai perder tudo se parar?)
- "Você está no top 10% dos usuários!" (status social)
- Owl triste se você não fizer a lição (culpa)

**Psicologia:** Gamificação + FOMO + culpa
**Resultado:** 500M+ usuários ativos, conversão altíssima para Plus

---

### Spotify

**❌ Ruim:** "Premium por R$21,90/mês"

**✅ Spotify:**

- Free tier com ads (dor intencional)
- "Remova os anúncios" (alívio da dor)
- "Baixe suas músicas" (novo benefício)
- Trial de 3 meses por R$0,99 (preço âncora)

**Psicologia:** Freemium + dor intencional + alívio
**Resultado:** 200M+ pagantes

---

### Notion

**❌ Ruim:** "Cadastre-se"

**✅ Notion:**

- Templates prontos (quick win imediato)
- Workspace vazio já começa com exemplo (não assusta)
- "Convide seu time" (efeito rede)
- Plano Free generoso (lock-in)

**Psicologia:** Redução de fricção + valor imediato
**Resultado:** Crescimento viral, $10B valuation

---

## 🧪 8 Princípios de Persuasão

### 1. Reciprocidade (Robert Cialdini)

**Conceito:**

- Humanos sentem necessidade de retribuir favores
- Você dá algo → pessoa sente-se em dívida

**Aplicação:**

- Trial grátis SEM pedir cartão
- Free tier generoso
- Conteúdo grátis (blog, ebooks)

**Evidência:**

- Dropbox: 2GB grátis → 500M usuários
- Spotify: 3 meses grátis → 200M pagantes
- **Fy:** 7 dias grátis → esperamos 40-50% conversão

**Como implementar:**

```typescript
// Trial sem cartão
const startTrial = async () => {
  // ❌ Ruim: Pedir cartão antes
  if (!hasCreditCard) {
    showPaymentForm();
    return;
  }

  // ✅ Bom: Dar acesso imediato
  await activateProFeatures();
  setTrialEndDate(+7 days);
  // Cartão só é pedido no dia 7
};
```

---

### 2. Escassez e Urgência

**Conceito:**

- Humanos valorizam mais o que é raro/limitado
- FOMO (Fear of Missing Out)

**Aplicação:**

- "Seu trial acaba em 2 dias!"
- "Últimas 24h para manter funcionalidade"
- "Oferta válida até amanhã"

**⚠️ Cuidado:**

- **Ética:** Usar escassez REAL, não fake
- ❌ "Apenas 2 vagas!" (mentira)
- ✅ "Trial acaba em 2 dias" (verdade)

**Evidência:**

- Booking.com: "2 pessoas vendo este hotel" → +33% conversão
- Amazon: "Só restam 3 unidades" → +12% vendas

---

### 3. Prova Social

**Conceito:**

- Humanos seguem comportamento de grupo
- "Se muitos fazem, deve ser bom"

**Aplicação:**

- "Junte-se a 10.000+ usuários"
- "95% economizam R$200/mês"
- Depoimentos de clientes
- "Top 10% dos usuários"

**Exemplo visual:**

```
┌─────────────────────────────────────┐
│ 👥 Junte-se a 10.243 usuários       │
│    que já economizaram R$1.2M       │
│                                     │
│ ⭐⭐⭐⭐⭐ 4.8/5 (2.341 avaliações)  │
└─────────────────────────────────────┘
```

**Evidência:**

- Yelp reviews: +18% conversão vs sem reviews
- Amazon ratings: produtos com 50+ reviews vendem 4x mais

---

### 4. Comprometimento e Consistência

**Conceito:**

- Pessoa que dá passo pequeno → mais provável dar próximo
- "Já comecei, vou terminar"
- Sunk cost (já investi, não vou desperdiçar)

**Aplicação:**

- Onboarding gradual (pequenos passos)
- "Complete perfil 65%"
- Metas criadas pelo próprio usuário
- Histórico de conquistas

**Foot-in-the-door technique:**

```
Passo 1: "Qual seu nome?"           ← Fácil, todo mundo faz
Passo 2: "Qual seu email?"          ← Já começou, vai continuar
Passo 3: "Qual sua renda?"          ← Já investiu tempo, não quer desperdiçar
Passo 4: "Conectar banco?"          ← Já está comprometido

vs

Tudo de uma vez: "Preencha 20 campos" ← Desiste
```

**Evidência:**

- LinkedIn: "Complete seu perfil" → +40% engagement
- Duolingo: "Streak de 47 dias" → +70% retenção

---

### 5. Ancoragem (Pricing)

**Conceito:**

- Primeira informação "ancora" percepção
- Comparações subsequentes são relativas à âncora

**Aplicação no Pricing:**

```
❌ Ruim (sem contexto):
Pro: R$59,99/mês

✅ Bom (com âncora):
Ultra: R$99,99/mês (riscado)
Pro:   R$59,99/mês ← "Melhor valor!" ⭐
Free:  R$0/mês

→ R$59,99 parece BARATO depois de ver R$99,99
```

**Goldilocks Effect:**

- Plano caro (âncora)
- **Plano médio (60% escolhem)** ← "nem muito caro, nem muito barato"
- Plano grátis (comparação)

**Evidência:**

- Apple: iPhone Pro Max R$9k faz Pro R$7k parecer "razoável"
- SaaS: 60% escolhem plano do meio

---

### 6. Aversão à Perda (Loss Aversion)

**Conceito:**

- Perder algo DOI 2x mais que ganhar
- "Não quero perder o que já tenho"

**Aplicação:**

```
┌─────────────────────────────────────┐
│ 🎉 Seu trial está acabando!         │
│                                     │
│ Você vai PERDER:                    │
│ ✗ 156 transações sincronizadas      │
│ ✗ Categorização automática          │
│ ✗ 12 alertas inteligentes           │
│                                     │
│ [Manter Pro - R$59,99/mês]          │
│ [Perder tudo e voltar para Free]    │
└─────────────────────────────────────┘
```

**Framing negativo funciona melhor:**

- ❌ "Ganhe categorização automática" (6% conversão)
- ✅ "Não perca categorização automática" (14% conversão)

**Evidência:**

- Kahneman & Tversky: Perda pesa 2.25x mais que ganho equivalente
- Experimento: "Economize R$5" (40%) vs "Não perca R$5" (70%)

---

### 7. Efeito Zeigarnik (Tarefas Incompletas)

**Conceito:**

- Humanos odeiam deixar coisas incompletas
- Tendência obsessiva de completar tarefas iniciadas

**Aplicação:**

```
┌─────────────────────────────────────┐
│ 👤 Gustavo                    65% ✓ │
│ [Complete para desbloquear          │
│  relatórios avançados]              │
│                                     │
│ Faltam:                             │
│ ☐ Adicionar endereço          +10%  │
│ ☐ Configurar alertas          +15%  │
│ ☐ Conectar segundo banco      +10%  │
└─────────────────────────────────────┘
```

**Por que funciona:**

- Cérebro vê tarefa incompleta → tensão
- Completar → dopamina liberada → alívio
- Cria loop viciante

**Evidência:**

- LinkedIn: Progress bar → +20% profile completion
- Angry Birds: "3 estrelas" → +60% re-tentativas

---

### 8. Gamificação (Loop de Dopamina)

**Conceito:**

- Dopamina é liberada quando atingimos objetivo
- Cria vício positivo (quer sentir de novo)
- Progresso visível = motivação

**Aplicação:**

```typescript
// Sistema de badges
const badges = [
  { id: "first-week", title: "Primeira semana!", trigger: 7 },
  { id: "categorizer", title: "100 transações categorizadas", trigger: 100 },
  {
    id: "goal-setter",
    title: "Primeira meta atingida",
    trigger: "goal_completed",
  },
  { id: "streak-7", title: "7 dias consecutivos", trigger: "streak_7" },
  { id: "top-10", title: "Top 10% mais organizado", trigger: "percentile_90" },
];

// Leaderboard (cuidado: pode desmotivar quem está atrás)
const stats = {
  you: { rank: 127, score: 450 },
  average: { score: 320 },
  percentile: 85, // "Você está no top 15%!"
};
```

**Elementos de gamificação:**

- ✅ Progress bars
- ✅ Badges/conquistas
- ✅ Streaks (consecutivos)
- ✅ Levels/ranking
- ⚠️ Leaderboards (competição pode estressar)

**Evidência:**

- Duolingo: Gamificação → 70% retenção D7
- Fitbit: Steps tracking → +35% exercício

---

## 🎯 BJ Fogg Behavior Model

### Fórmula do Comportamento

```
Behavior = Motivation × Ability × Trigger

Para comportamento acontecer:
  Motivation (quer fazer) +
  Ability (é fácil fazer) +
  Trigger (lembrete no momento certo)
= AÇÃO
```

### Aplicação no Fy

#### Comportamento: Aceitar trial Open Finance

| Componente      | Implementação                      | Resultado              |
| --------------- | ---------------------------------- | ---------------------- |
| **Motivation**  | "Economize 10 min/dia"             | ALTO (benefício claro) |
| **Ability**     | 1 clique, sem cartão               | ALTO (super fácil)     |
| **Trigger**     | No onboarding, opção "recomendada" | Momento certo          |
| **→ Conversão** |                                    | **40-50%** ✅          |

---

#### Comportamento: Converter trial→Pro

| Componente      | Implementação               | Resultado                |
| --------------- | --------------------------- | ------------------------ |
| **Motivation**  | "Vai perder 156 transações" | ALTO (aversão à perda)   |
| **Ability**     | 1 clique, já tem cartão     | ALTO (sem fricção)       |
| **Trigger**     | Dia 7, notificação + email  | Momento certo (urgência) |
| **→ Conversão** |                             | **45-50%** ✅            |

---

#### Comportamento: Completar perfil

| Componente      | Implementação                   | Resultado                |
| --------------- | ------------------------------- | ------------------------ |
| **Motivation**  | "+10% = desbloqueia relatórios" | MÉDIO (benefício futuro) |
| **Ability**     | 1 campo de cada vez             | ALTO (fácil)             |
| **Trigger**     | Dashboard header sempre visível | Sempre presente          |
| **→ Conversão** |                                 | **60-70%** ✅            |

---

### Matriz de Priorização (Fogg)

```
         Alta Motivação
              │
   Difícil    │    Fácil
──────────────┼──────────────
              │
   ❌ Não    │    ✅ SIM!
   funciona  │    (Sweet Spot)
              │
         Baixa Motivação
```

**Regra:** Quando motivação é baixa, ability precisa ser MUITO alta

**Exemplo:**

- Motivação baixa: "Adicionar endereço" (não vê benefício imediato)
- Ability precisa ser MUITO alta: 1 campo, auto-complete, skippable
- Se pedir 10 campos → 0% completam

---

## ❌ Dark Patterns (NÃO usar)

### O que são Dark Patterns?

**Definição:** Truques de design que manipulam usuário contra seu interesse

**Por que NÃO usar:**

- ❌ Ética: não é certo
- ❌ Long-term: destrói confiança
- ❌ Churn: usuários cancelam logo depois
- ❌ Reputação: bad word-of-mouth
- ❌ Legal: pode ser ilegal (LGPD, GDPR)

---

### Exemplos de Dark Patterns (EVITAR)

#### 1. Roach Motel (Fácil entrar, difícil sair)

**❌ Dark:**

- Assinar: 1 clique
- Cancelar: ligar, falar com atendente, esperar 30min

**✅ Fy:**

- Assinar: 1 clique
- Cancelar: 1 clique em Settings → "Cancelar assinatura"

---

#### 2. Hidden Costs (Custos escondidos)

**❌ Dark:**

- "Grátis!" (esconde que cobra depois)
- Preço final só aparece no checkout

**✅ Fy:**

- Preço claro em todas páginas
- Trial: "Grátis por 7 dias, depois R$59,99/mês"

---

#### 3. Trick Questions (Perguntas confusas)

**❌ Dark:**

```
□ Não me envie promoções
□ Não quero descontos
```

(Dupla negativa confunde → usuário erra)

**✅ Fy:**

```
□ Quero receber promoções
□ Quero receber descontos
```

(Claro e direto)

---

#### 4. Confirmshaming (Guilt trip)

**❌ Dark:**

- "Não obrigado, eu não me importo com minhas finanças" (guilt)
- "Não quero economizar dinheiro"

**✅ Fy:**

- "Continuar Free" (neutro)
- "Não, obrigado" (respeitoso)

---

#### 5. Fake Urgency (Urgência falsa)

**❌ Dark:**

- "Apenas 2 vagas restantes!" (mentira)
- "Oferta termina em 3h!" (reset todo dia)

**✅ Fy:**

- "Seu trial acaba em 2 dias" (verdade)
- "Oferta válida até 31/12" (data real)

---

#### 6. Bait and Switch (Isca e troca)

**❌ Dark:**

- Anunciar recurso grátis
- No último momento: "Ah, isso é pago"

**✅ Fy:**

- Se recurso é Pro, avisar ANTES do usuário tentar usar
- "Este recurso está disponível no plano Pro"

---

#### 7. Forced Continuity (Continuidade forçada)

**❌ Dark:**

- Trial acaba → cobra cartão SEM avisar
- Dificulta cancelamento antes do fim do trial

**✅ Fy:**

- Avisar 3 dias antes: "Trial acaba em 3 dias"
- Email dia 7: "Trial acabou, quer continuar?"
- Permitir cancelar DURANTE o trial

---

#### 8. Sneak into Basket (Adicionar sem consentimento)

**❌ Dark:**

- Checkbox pré-marcado: "Adicionar seguro R$19,90"
- Usuário não vê → paga sem querer

**✅ Fy:**

- Checkboxes sempre desmarcadas por padrão
- Opt-in explícito

---

## 🎨 Copywriting que Converte

### Fórmulas Clássicas

#### PAS (Problem - Agitate - Solve)

**Estrutura:**

1. **Problem:** Identifica dor do usuário
2. **Agitate:** Intensifica a dor
3. **Solve:** Apresenta solução

**Exemplo:**

```
[Problem]
Você gasta 15 minutos todo dia adicionando transações manualmente?

[Agitate]
São 7.5 horas por mês fazendo trabalho repetitivo.
Tempo que você poderia estar com a família ou fazendo o que gosta.

[Solve]
Com Fy Pro, suas transações aparecem automaticamente.
Zero trabalho manual. Zero estresse.

[CTA]
Experimente grátis por 7 dias →
```

---

#### AIDA (Attention - Interest - Desire - Action)

**Estrutura:**

1. **Attention:** Chama atenção
2. **Interest:** Cria interesse
3. **Desire:** Cria desejo
4. **Action:** CTA

**Exemplo:**

```
[Attention]
🚨 Você está gastando R$1.200/mês em delivery

[Interest]
Isso é 30% a mais que mês passado.
No ritmo atual, são R$14.400/ano.

[Desire]
E se você pudesse receber alertas ANTES de estourar o orçamento?
Com Fy Pro, você seria avisado: "Você vai estourar em 2 dias"

[Action]
Ative alertas inteligentes →
```

---

### Palavras que Convertem

**Poderosas:**

- ✅ Você (personaliza)
- ✅ Grátis (todo mundo ama)
- ✅ Porque (justifica)
- ✅ Instantaneamente (urgência)
- ✅ Novo (novidade)

**Fracas:**

- ❌ Nós/nosso (foco errado)
- ❌ Talvez (incerteza)
- ❌ Eventualmente (não urgente)
- ❌ Recursos (features não benefícios)

---

### Benefícios vs Features

**❌ Features (o que é):**

- "Categorização automática com IA"
- "Sync Open Banking"
- "Relatórios customizáveis"

**✅ Benefícios (o que ganha):**

- "Economize 15 min/dia - IA categoriza para você"
- "Suas transações aparecem sozinhas - zero digitação"
- "Veja exatamente onde seu dinheiro vai"

**Regra:** Features dizem O QUE é, benefícios dizem POR QUE importa

---

## 📊 A/B Testing (Otimização Contínua)

### O que testar

**1. Headlines/CTAs:**

```
A: "Cadastre-se"          vs  B: "Comece grátis"
A: "Experimente agora"    vs  B: "7 dias grátis"
A: "Saiba mais"           vs  B: "Economize 10 min/dia"
```

**2. Cores de botões:**

```
A: Verde (seguro)         vs  B: Laranja (urgência)
A: Azul (confiança)       vs  B: Vermelho (ação)
```

**3. Timing de gatilhos:**

```
A: Mostrar upgrade após 10 transações
B: Mostrar upgrade após 7 dias
C: Mostrar upgrade após estourar meta
```

**4. Pricing display:**

```
A: R$59,99/mês
B: R$1,99/dia (parece mais barato)
C: R$59,99/mês (economize R$220/ano vs mensal)
```

---

### Como testar

**Ferramentas:**

- Google Optimize (grátis)
- Optimizely (pago, mais robusto)
- VWO (Visual Website Optimizer)
- Amplitude (analytics + experiments)

**Processo:**

1. **Hipótese:** "Botão laranja vai converter 15% a mais"
2. **Teste:** 50% vê verde (A), 50% vê laranja (B)
3. **Métrica:** Taxa de cliques
4. **Duração:** Até significância estatística (>95% confidence)
5. **Decisão:** Implementar vencedor

**Cuidado:**

- ⚠️ Não testar 10 coisas ao mesmo tempo (impossível saber o que funcionou)
- ⚠️ Esperar significância estatística (não concluir com 50 usuários)
- ⚠️ Testar uma variável por vez (isola o efeito)

---

## 🧪 Framework Prático de Implementação

### Checklist de Conversão

Quando criar qualquer fluxo de conversão, validar:

#### 1. Motivação está clara?

- [ ] Benefício é óbvio em <3 segundos?
- [ ] Benefício é específico (não genérico)?
- [ ] Usa números concretos? ("economize 10 min", não "economize tempo")

#### 2. Ability está maximizada?

- [ ] Menos de 3 cliques?
- [ ] Menos de 3 campos de formulário?
- [ ] Tem skip/later option?
- [ ] Funciona mobile? (60% do tráfego)

#### 3. Trigger está no momento certo?

- [ ] Aparece no momento de dor/necessidade?
- [ ] Não interrompe fluxo crítico?
- [ ] Pode ser facilmente dispensado?

#### 4. Sem dark patterns?

- [ ] É fácil cancelar quanto assinar?
- [ ] Preços claros?
- [ ] Sem confirmshaming?
- [ ] Sem fake urgency?

---

## 📚 Recursos e Referências

### Livros Essenciais

1. **"Hooked"** - Nir Eyal
   - Como criar produtos viciantes
   - Hook Model: Trigger → Action → Reward → Investment

2. **"Influence"** - Robert Cialdini
   - 6 princípios de persuasão
   - Bíblia do marketing

3. **"Don't Make Me Think"** - Steve Krug
   - UX básico, usabilidade
   - Simplicidade acima de tudo

4. **"The Mom Test"** - Rob Fitzpatrick
   - Como validar ideias sem viés
   - Entrevistar usuários corretamente

5. **"Atomic Habits"** - James Clear
   - Como criar hábitos
   - Aplicável a product design

---

### Frameworks e Modelos

1. **Fogg Behavior Model**
   - B = MAT (Motivation × Ability × Trigger)

2. **Hook Model** (Nir Eyal)
   - Trigger → Action → Variable Reward → Investment

3. **Jobs to be Done** (Clayton Christensen)
   - Usuário "contrata" produto para fazer job
   - Foco no job, não no produto

4. **Growth Loops** (Reforge)
   - Loop viral que se auto-alimenta
   - Exemplo: Dropbox (convida amigo = mais espaço)

---

### Tools e Software

**Analytics:**

- Amplitude (eventos, funnels, cohorts)
- Mixpanel (similar Amplitude)
- Google Analytics 4 (básico, grátis)

**A/B Testing:**

- Google Optimize (grátis)
- Optimizely (pago)
- VWO (Visual Website Optimizer)

**Heatmaps/Session Recording:**

- Hotjar (grátis até 35 sessões/dia)
- FullStory (pago, mais robusto)
- Microsoft Clarity (grátis, ilimitado!)

**User Feedback:**

- Typeform (forms bonitos)
- SurveyMonkey (surveys)
- UserTesting.com (testar com usuários reais)

---

## 🎯 Conclusão

**UX Psychology NÃO é manipulação** - é **design intencional** para:

1. ✅ Reduzir fricção (facilitar vida do usuário)
2. ✅ Comunicar valor (mostrar benefícios claramente)
3. ✅ Criar hábito (usuário usa, vê valor, continua)
4. ✅ Converter no momento certo (quando usuário está pronto)

**É win-win:**

- Usuário: Resolve problema, economiza tempo/dinheiro
- Empresa: Constrói negócio sustentável, ganha dinheiro

**Seria errado se:**

- ❌ Cobrasse sem entregar valor (scam)
- ❌ Escondesse custos (dark pattern)
- ❌ Tornasse impossível cancelar (dark UX)

**Mas quando bem feito:**

- ✅ Ética: entrega valor real
- ✅ Long-term: constrói confiança
- ✅ Crescimento: word-of-mouth positivo
- ✅ Sustentável: negócio cresce saudável

---

**Agora você sabe o jogo. Use sem dó! 😈**

Bem-vindo ao lado sombrio do product design! 🎭

---

**Documento mantido por:** Product Team
**Feedback:** gustavo@Fy.com.br
