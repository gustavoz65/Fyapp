import { CheckCircle, AlertCircle, Info, XCircle } from 'lucide-react'
import Card from '@/components/Card'
import Badge from '@/components/Badge'
import { notifications } from '@/data/notifications'

export default function Notifications() {
  const iconMap = {
    success: CheckCircle,
    warning: AlertCircle,
    info: Info,
    error: XCircle,
  }

  return (
    <div className="space-y-6">
      <Card title="Notificações" subtitle={`${notifications.filter(n => !n.read).length} não lidas`}>
        <div className="space-y-3">
          {notifications.map((notification) => {
            const Icon = iconMap[notification.type]
            return (
              <div
                key={notification.id}
                className={`p-4 rounded-lg border ${
                  notification.read ? 'bg-white border-gray-200' : 'bg-primary-50 border-primary-200'
                }`}
              >
                <div className="flex items-start gap-3">
                  <Icon className={`w-5 h-5 mt-0.5 ${
                    notification.type === 'success' ? 'text-green-600' :
                    notification.type === 'warning' ? 'text-yellow-600' :
                    notification.type === 'error' ? 'text-red-600' : 'text-blue-600'
                  }`} />
                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <h4 className="font-semibold text-gray-900">{notification.title}</h4>
                      {!notification.read && <Badge variant="info">Nova</Badge>}
                    </div>
                    <p className="text-sm text-gray-600 mb-2">{notification.message}</p>
                    <p className="text-xs text-gray-500">
                      {new Date(notification.date).toLocaleString('pt-BR')}
                    </p>
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      </Card>
    </div>
  )
}
