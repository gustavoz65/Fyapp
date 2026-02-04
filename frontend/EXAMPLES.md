# 🎨 Exemplos de Código - Cashing

Este arquivo contém exemplos práticos de como usar os componentes e funcionalidades do Cashing.

## 📦 Importando Componentes

```typescript
import Button from '@/components/Button'
import Card from '@/components/Card'
import Modal from '@/components/Modal'
import Table from '@/components/Table'
import { useAuthStore } from '@/stores/authStore'
import { formatCurrency } from '@/lib/utils'
```

## 🔘 Usando Botões

```tsx
// Botão primário
<Button onClick={handleClick}>
  Salvar
</Button>

// Botão com loading
<Button loading={isLoading} variant="primary">
  Processando...
</Button>

// Botão de perigo
<Button variant="danger" size="sm">
  Excluir
</Button>

// Botão ghost
<Button variant="ghost">
  Cancelar
</Button>
```

## 📇 Usando Cards

```tsx
<Card 
  title="Título do Card"
  subtitle="Subtítulo opcional"
  headerAction={<Button size="sm">Ação</Button>}
  footer={<div>Footer content</div>}
>
  <p>Conteúdo do card</p>
</Card>
```

## 🔲 Usando Modais

```tsx
import { useState } from 'react'
import Modal from '@/components/Modal'

function MyComponent() {
  const [isOpen, setIsOpen] = useState(false)

  return (
    <>
      <Button onClick={() => setIsOpen(true)}>
        Abrir Modal
      </Button>

      <Modal
        isOpen={isOpen}
        onClose={() => setIsOpen(false)}
        title="Meu Modal"
        size="lg"
        footer={
          <>
            <Button variant="ghost" onClick={() => setIsOpen(false)}>
              Cancelar
            </Button>
            <Button onClick={handleSave}>
              Salvar
            </Button>
          </>
        }
      >
        <p>Conteúdo do modal</p>
      </Modal>
    </>
  )
}
```

## 📊 Usando Tabelas

```tsx
import Table from '@/components/Table'
import type { Transaction } from '@/types'

const columns = [
  {
    key: 'date',
    label: 'Data',
    render: (value: unknown) => formatDate(value as string),
  },
  {
    key: 'description',
    label: 'Descrição',
  },
  {
    key: 'amount',
    label: 'Valor',
    render: (value: unknown, row: Transaction) => (
      <span className={row.type === 'income' ? 'text-green-600' : 'text-red-600'}>
        {formatCurrency(value as number)}
      </span>
    ),
  },
]

<Table<Transaction>
  columns={columns}
  data={transactions}
  searchable
  sortable
  onRowClick={(row) => console.log(row)}
/>
```

## 📝 Usando Formulários

```tsx
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import Input from '@/components/Input'
import Select from '@/components/Select'

const schema = z.object({
  name: z.string().min(3, 'Mínimo 3 caracteres'),
  email: z.string().email('Email inválido'),
  type: z.enum(['individual', 'business']),
})

type FormData = z.infer<typeof schema>

function MyForm() {
  const { register, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
  })

  const onSubmit = (data: FormData) => {
    console.log(data)
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Input
        {...register('name')}
        label="Nome"
        error={errors.name?.message}
      />

      <Input
        {...register('email')}
        type="email"
        label="Email"
        error={errors.email?.message}
      />

      <Select
        {...register('type')}
        label="Tipo"
        options={[
          { value: 'individual', label: 'Pessoa Física' },
          { value: 'business', label: 'Pessoa Jurídica' },
        ]}
        error={errors.type?.message}
      />

      <Button type="submit">Enviar</Button>
    </form>
  )
}
```

## 🏪 Usando Zustand Stores

```tsx
import { useAuthStore } from '@/stores/authStore'
import { useSettingsStore } from '@/stores/settingsStore'

function MyComponent() {
  // Auth store
  const { user, login, logout } = useAuthStore()

  // Settings store
  const { currency, setCurrency } = useSettingsStore()

  const handleLogin = () => {
    login({
      id: '1',
      name: 'João Silva',
      email: 'joao@exemplo.com',
      type: 'individual',
      plan: 'pro',
      createdAt: new Date().toISOString(),
    })
  }

  return (
    <div>
      <p>Usuário: {user?.name}</p>
      <p>Moeda: {currency}</p>
      
      <Button onClick={handleLogin}>Login</Button>
      <Button onClick={logout}>Logout</Button>
      
      <Select
        value={currency}
        onChange={(e) => setCurrency(e.target.value as 'BRL' | 'USD' | 'EUR')}
        options={[
          { value: 'BRL', label: 'Real' },
          { value: 'USD', label: 'Dólar' },
          { value: 'EUR', label: 'Euro' },
        ]}
      />
    </div>
  )
}
```

## 🎨 Usando Utilitários

```tsx
import { formatCurrency, formatDate, formatPercent, cn } from '@/lib/utils'

// Formatar moeda
formatCurrency(1234.56) // "R$ 1.234,56"
formatCurrency(1234.56, 'USD') // "$1,234.56"

// Formatar data
formatDate('2024-01-27') // "27/01/2024"
formatDate('2024-01-27', 'MM/yyyy') // "01/2024"

// Formatar percentual
formatPercent(0.85) // "85.0%"

// Combinar classes com cn
<div className={cn(
  'base-class',
  isActive && 'active-class',
  'another-class'
)}>
```

## 📈 Criando Gráficos

```tsx
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import { formatCurrency } from '@/lib/utils'

const data = [
  { month: 'Jan', receitas: 5000, despesas: 3000 },
  { month: 'Fev', receitas: 6000, despesas: 4000 },
  { month: 'Mar', receitas: 7000, despesas: 3500 },
]

<ResponsiveContainer width="100%" height={300}>
  <LineChart data={data}>
    <CartesianGrid strokeDasharray="3 3" />
    <XAxis dataKey="month" />
    <YAxis />
    <Tooltip formatter={(value: number) => formatCurrency(value)} />
    <Legend />
    <Line type="monotone" dataKey="receitas" stroke="#10b981" strokeWidth={2} />
    <Line type="monotone" dataKey="despesas" stroke="#ef4444" strokeWidth={2} />
  </LineChart>
</ResponsiveContainer>
```

## 🎯 Progress Bar

```tsx
import ProgressBar from '@/components/ProgressBar'

// Básico
<ProgressBar value={75} max={100} />

// Com label e percentual
<ProgressBar
  value={budget.spent}
  max={budget.amount}
  color={percentage >= 80 ? 'warning' : 'success'}
  label="Orçamento de Alimentação"
  showPercent
/>
```

## 🏷️ Badges

```tsx
import Badge from '@/components/Badge'

<Badge variant="success">Ativo</Badge>
<Badge variant="warning">Pendente</Badge>
<Badge variant="error">Cancelado</Badge>
<Badge variant="info">Informação</Badge>
<Badge variant="default">Padrão</Badge>
```

## 🔄 Loading States

```tsx
import Skeleton, { SkeletonCard, SkeletonTable } from '@/components/Skeleton'

// Loading de card
{isLoading ? <SkeletonCard /> : <Card>...</Card>}

// Loading de tabela
{isLoading ? <SkeletonTable rows={5} columns={4} /> : <Table>...</Table>}

// Loading customizado
<Skeleton height="20px" width="200px" />
```

## 📭 Empty States

```tsx
import EmptyState from '@/components/EmptyState'
import { FileQuestion } from 'lucide-react'

<EmptyState
  icon={FileQuestion}
  title="Nenhuma transação encontrada"
  description="Você ainda não possui transações cadastradas. Clique no botão abaixo para adicionar sua primeira transação."
  action={<Button onClick={handleAdd}>Adicionar Transação</Button>}
/>
```

## 🎭 Animations com Framer Motion

```tsx
import { motion } from 'framer-motion'

<motion.div
  initial={{ opacity: 0, y: 20 }}
  animate={{ opacity: 1, y: 0 }}
  transition={{ duration: 0.5 }}
>
  <Card>Conteúdo animado</Card>
</motion.div>

// Lista animada
<motion.div
  initial={{ opacity: 0 }}
  animate={{ opacity: 1 }}
  transition={{ staggerChildren: 0.1 }}
>
  {items.map((item, i) => (
    <motion.div
      key={item.id}
      initial={{ x: -20, opacity: 0 }}
      animate={{ x: 0, opacity: 1 }}
      transition={{ delay: i * 0.1 }}
    >
      <Card>{item.name}</Card>
    </motion.div>
  ))}
</motion.div>
```

## 🛣️ Navegação Protegida

```tsx
import { Navigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore()

  if (!isAuthenticated) {
    return <Navigate to="/auth/login" replace />
  }

  return <>{children}</>
}
```

## 🔍 React Query

```tsx
import { useQuery, useMutation } from '@tanstack/react-query'

// Query
const { data, isLoading, error } = useQuery({
  queryKey: ['transactions'],
  queryFn: fetchTransactions,
})

// Mutation
const mutation = useMutation({
  mutationFn: createTransaction,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['transactions'] })
  },
})
```

## 🎨 Tailwind Custom Classes

```tsx
// Hover e focus states
<button className="bg-blue-500 hover:bg-blue-600 focus:ring-2 focus:ring-blue-300">
  Hover me
</button>

// Responsivo
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  {/* Cards */}
</div>

// Dark mode (quando implementado)
<div className="bg-white dark:bg-gray-800 text-gray-900 dark:text-white">
  {/* Content */}
</div>
```

## 📱 Responsive Patterns

```tsx
// Drawer mobile, sidebar desktop
<aside className="hidden lg:block">
  {/* Desktop sidebar */}
</aside>

<button className="lg:hidden" onClick={() => setDrawerOpen(true)}>
  <Menu />
</button>

// Grid responsivo
<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
  {/* Items */}
</div>
```

## 🎯 Trabalhando com Mock Data

```typescript
import { transactions } from '@/data/transactions'
import { categories } from '@/data/categories'
import { banks } from '@/data/banks'
import { budgets } from '@/data/budgets'
import { notifications } from '@/data/notifications'

// Filtrar transações do mês atual
const currentMonthTransactions = transactions.filter(t => {
  const tMonth = new Date(t.date).getMonth()
  const currentMonth = new Date().getMonth()
  return tMonth === currentMonth && t.status === 'completed'
})

// Calcular total de receitas
const totalIncome = transactions
  .filter(t => t.type === 'income' && t.status === 'completed')
  .reduce((acc, t) => acc + t.amount, 0)

// Agrupar por categoria
const byCategory = categories.map(cat => ({
  ...cat,
  total: transactions
    .filter(t => t.categoryId === cat.id)
    .reduce((acc, t) => acc + t.amount, 0),
}))
```

## 🔐 Auth Flow Completo

```tsx
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/stores/authStore'

function LoginPage() {
  const navigate = useNavigate()
  const { login } = useAuthStore()

  const handleLogin = async (credentials: LoginForm) => {
    try {
      // Mock login - substitua por chamada real à API
      const userData = await mockAuthApi.login(credentials)
      
      login(userData)
      navigate('/dashboard')
    } catch (error) {
      console.error('Erro no login:', error)
    }
  }

  return <LoginForm onSubmit={handleLogin} />
}
```

---

## 💡 Dicas de Performance

1. **Lazy Loading**: Use `React.lazy()` para rotas
```tsx
const Dashboard = lazy(() => import('./routes/protected/dashboard/Dashboard'))
```

2. **Memoização**: Use `useMemo` e `useCallback`
```tsx
const expensiveValue = useMemo(() => calculateTotal(data), [data])
const handleClick = useCallback(() => doSomething(), [])
```

3. **Virtual Lists**: Para listas grandes
```tsx
import { useVirtual } from 'react-virtual'
```

4. **Debounce**: Para inputs de busca
```tsx
import { debounce } from '@/lib/utils'
const debouncedSearch = debounce(handleSearch, 300)
```

---

Desenvolvido com ❤️ para Cashing
