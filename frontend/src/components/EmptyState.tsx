import { ReactNode } from 'react'
import { LucideIcon, FileQuestion } from 'lucide-react'

interface EmptyStateProps {
  title: string
  description?: string
  action?: ReactNode
  icon?: LucideIcon
}

export default function EmptyState({ title, description, action, icon: Icon = FileQuestion }: EmptyStateProps) {
  return (
    <div className="flex flex-col items-center justify-center py-12 px-4">
      <div className="bg-gray-100 rounded-full p-6 mb-4">
        <Icon className="w-12 h-12 text-gray-400" />
      </div>
      <h3 className="text-lg font-semibold text-gray-900 mb-2">{title}</h3>
      {description && <p className="text-gray-600 text-center max-w-md mb-6">{description}</p>}
      {action && <div>{action}</div>}
    </div>
  )
}
