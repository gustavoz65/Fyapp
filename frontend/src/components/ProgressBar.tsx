import { cn } from '@/lib/utils'

interface ProgressBarProps {
  value: number
  max?: number
  color?: 'primary' | 'success' | 'warning' | 'danger'
  label?: string
  showPercent?: boolean
}

export default function ProgressBar({ value, max = 100, color = 'primary', label, showPercent }: ProgressBarProps) {
  const percent = (value / max) * 100
  const cappedPercent = Math.min(percent, 100)

  const colors = {
    primary: 'bg-primary-600',
    success: 'bg-green-600',
    warning: 'bg-yellow-600',
    danger: 'bg-red-600',
  }

  const bgColors = {
    primary: 'bg-primary-100',
    success: 'bg-green-100',
    warning: 'bg-yellow-100',
    danger: 'bg-red-100',
  }

  return (
    <div className="w-full">
      {(label || showPercent) && (
        <div className="flex justify-between items-center mb-2">
          {label && <span className="text-sm font-medium text-gray-700">{label}</span>}
          {showPercent && (
            <span className="text-sm font-medium text-gray-700">{cappedPercent.toFixed(0)}%</span>
          )}
        </div>
      )}
      <div className={cn('w-full h-2 rounded-full overflow-hidden', bgColors[color])}>
        <div
          className={cn('h-full rounded-full transition-all duration-300', colors[color])}
          style={{ width: `${cappedPercent}%` }}
        />
      </div>
    </div>
  )
}
