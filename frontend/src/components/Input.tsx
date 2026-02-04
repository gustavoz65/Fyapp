import { forwardRef, InputHTMLAttributes } from 'react'
import { cn } from '@/lib/utils'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
  helperText?: string
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, helperText, className, ...props }, ref) => {
    return (
      <div className="w-full">
        {label && (
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            {label}
            {props.required && <span className="text-red-500 ml-1">*</span>}
          </label>
        )}
        <input
          ref={ref}
          className={cn(
            'w-full px-4 py-2 border border-gray-300 dark:border-gray-700 rounded-lg',
            'bg-white dark:bg-[#0f0f0f] text-gray-900 dark:text-gray-100',
            'focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors',
            'placeholder:text-gray-400 dark:placeholder:text-gray-500',
            error && 'border-red-500 focus:ring-red-500',
            props.disabled && 'bg-gray-100 dark:bg-gray-800 cursor-not-allowed opacity-50',
            className
          )}
          {...props}
        />
        {error && <p className="mt-1 text-sm text-red-600 dark:text-red-400">{error}</p>}
        {helperText && !error && <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{helperText}</p>}
      </div>
    )
  }
)

Input.displayName = 'Input'

export default Input
