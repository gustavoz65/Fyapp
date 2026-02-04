import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { BarChart3 } from 'lucide-react'
import Button from '@/components/Button'
import Input from '@/components/Input'
import Select from '@/components/Select'
import { useAuthStore } from '@/stores/authStore'

const registerSchema = z.object({
  name: z.string().min(3, 'Nome deve ter no mínimo 3 caracteres'),
  email: z.string().email('Email inválido'),
  password: z.string().min(6, 'Senha deve ter no mínimo 6 caracteres'),
  confirmPassword: z.string(),
  type: z.enum(['individual', 'business']),
  plan: z.enum(['free', 'pro', 'business']),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'As senhas não coincidem',
  path: ['confirmPassword'],
})

type RegisterForm = z.infer<typeof registerSchema>

export default function Register() {
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { login } = useAuthStore()

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterForm>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      type: 'individual',
      plan: 'free',
    },
  })

  const onSubmit = async (data: RegisterForm) => {
    setLoading(true)
    setTimeout(() => {
      login({
        id: '1',
        name: data.name,
        email: data.email,
        type: data.type,
        plan: data.plan,
        createdAt: new Date().toISOString(),
      })
      navigate('/auth/verify-email')
    }, 1000)
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-white flex items-center justify-center p-4">
      <div className="w-full max-w-2xl">
        <div className="text-center mb-8">
          <Link to="/" className="inline-flex items-center gap-2 mb-4">
            <BarChart3 className="w-10 h-10 text-primary-600" />
            <span className="text-3xl font-bold text-gray-900">Cashing</span>
          </Link>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Criar sua conta</h1>
          <p className="text-gray-600">Comece a gerenciar suas finanças hoje mesmo</p>
        </div>

        <div className="bg-white rounded-2xl shadow-xl p-8">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="grid md:grid-cols-2 gap-6">
              <Select
                {...register('type')}
                label="Tipo de Conta"
                options={[
                  { value: 'individual', label: 'Pessoa Física' },
                  { value: 'business', label: 'Pessoa Jurídica' },
                ]}
                error={errors.type?.message}
                disabled={loading}
              />

              <Select
                {...register('plan')}
                label="Plano"
                options={[
                  { value: 'free', label: 'Free - R$ 0/mês' },
                  { value: 'pro', label: 'Pro - R$ 49/mês' },
                  { value: 'business', label: 'Business - R$ 149/mês' },
                ]}
                error={errors.plan?.message}
                disabled={loading}
              />
            </div>

            <Input
              {...register('name')}
              label="Nome Completo"
              placeholder="João Silva"
              error={errors.name?.message}
              disabled={loading}
            />

            <Input
              {...register('email')}
              type="email"
              label="Email"
              placeholder="seu@email.com"
              error={errors.email?.message}
              disabled={loading}
            />

            <div className="grid md:grid-cols-2 gap-6">
              <Input
                {...register('password')}
                type="password"
                label="Senha"
                placeholder="••••••••"
                error={errors.password?.message}
                disabled={loading}
              />

              <Input
                {...register('confirmPassword')}
                type="password"
                label="Confirmar Senha"
                placeholder="••••••••"
                error={errors.confirmPassword?.message}
                disabled={loading}
              />
            </div>

            <label className="flex items-start gap-3 cursor-pointer">
              <input
                type="checkbox"
                required
                className="mt-1 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
              />
              <span className="text-sm text-gray-600">
                Concordo com os{' '}
                <a href="#" className="text-primary-600 hover:text-primary-700">
                  Termos de Serviço
                </a>{' '}
                e{' '}
                <a href="#" className="text-primary-600 hover:text-primary-700">
                  Política de Privacidade
                </a>
              </span>
            </label>

            <Button type="submit" className="w-full" loading={loading}>
              Criar Conta
            </Button>
          </form>

          <p className="mt-6 text-center text-sm text-gray-600">
            Já tem uma conta?{' '}
            <Link to="/auth/login" className="text-primary-600 hover:text-primary-700 font-medium">
              Entrar
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
