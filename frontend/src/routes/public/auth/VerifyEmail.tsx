import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { BarChart3, Mail, CheckCircle } from 'lucide-react'
import Button from '@/components/Button'

export default function VerifyEmail() {
  const [verified, setVerified] = useState(false)
  const navigate = useNavigate()

  useEffect(() => {
    const timer = setTimeout(() => {
      setVerified(true)
    }, 2000)
    return () => clearTimeout(timer)
  }, [])

  const handleContinue = () => {
    navigate('/dashboard')
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 to-white flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <Link to="/" className="inline-flex items-center gap-2 mb-4">
            <BarChart3 className="w-10 h-10 text-primary-600" />
            <span className="text-3xl font-bold text-gray-900">Cashing</span>
          </Link>
        </div>

        <div className="bg-white rounded-2xl shadow-xl p-8 text-center">
          {verified ? (
            <>
              <CheckCircle className="w-20 h-20 text-green-600 mx-auto mb-6" />
              <h1 className="text-2xl font-bold text-gray-900 mb-2">Email Verificado!</h1>
              <p className="text-gray-600 mb-8">
                Sua conta foi verificada com sucesso. Agora você pode acessar todas as funcionalidades do Cashing.
              </p>
              <Button onClick={handleContinue} className="w-full">
                Ir para o Dashboard
              </Button>
            </>
          ) : (
            <>
              <Mail className="w-20 h-20 text-primary-600 mx-auto mb-6" />
              <h1 className="text-2xl font-bold text-gray-900 mb-2">Verificando seu email...</h1>
              <p className="text-gray-600 mb-8">
                Por favor, aguarde enquanto verificamos sua conta.
              </p>
              <div className="flex justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  )
}
