import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { BarChart3, ArrowLeft, CheckCircle } from 'lucide-react'
import Button from '@/components/Button'
import Input from '@/components/Input'

const resetSchema = z.object({
  email: z.string().email('Email inválido'),
})

type ResetForm = z.infer<typeof resetSchema>

export default function ResetPassword() {
  const [loading, setLoading] = useState(false)
  const [sent, setSent] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResetForm>({
    resolver: zodResolver(resetSchema),
  })

  const onSubmit = async () => {
    setLoading(true)
    setTimeout(() => {
      setLoading(false)
      setSent(true)
    }, 1000)
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-white flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <Link to="/" className="inline-flex items-center gap-2 mb-4">
            <BarChart3 className="w-10 h-10 text-primary-600" />
            <span className="text-3xl font-bold text-gray-900">Cashing</span>
          </Link>
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Recuperar senha</h1>
          <p className="text-gray-600">
            {sent
              ? 'Enviamos um link para seu email'
              : 'Digite seu email para receber o link de recuperação'}
          </p>
        </div>

        <div className="bg-white rounded-2xl shadow-xl p-8">
          {sent ? (
            <div className="text-center">
              <CheckCircle className="w-16 h-16 text-green-600 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Email enviado!</h3>
              <p className="text-gray-600 mb-6">
                Verifique sua caixa de entrada e siga as instruções para redefinir sua senha.
              </p>
              <Link to="/auth/login">
                <Button className="w-full">Voltar para o Login</Button>
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <Input
                {...register('email')}
                type="email"
                label="Email"
                placeholder="seu@email.com"
                error={errors.email?.message}
                disabled={loading}
              />

              <Button type="submit" className="w-full" loading={loading}>
                Enviar Link de Recuperação
              </Button>

              <Link to="/auth/login" className="flex items-center justify-center gap-2 text-sm text-gray-600 hover:text-gray-900">
                <ArrowLeft className="w-4 h-4" />
                Voltar para o login
              </Link>
            </form>
          )}
        </div>
      </div>
    </div>
  )
}
