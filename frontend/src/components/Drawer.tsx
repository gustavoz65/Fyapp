import { Fragment, ReactNode } from 'react'
import { Transition } from '@headlessui/react'
import { X } from 'lucide-react'
import { cn } from '@/lib/utils'

interface DrawerProps {
  isOpen: boolean
  onClose: () => void
  title: string
  children: ReactNode
  side?: 'left' | 'right' | 'top' | 'bottom'
}

export default function Drawer({ isOpen, onClose, title, children, side = 'right' }: DrawerProps) {
  const sideClasses = {
    left: 'left-0 top-0 h-full',
    right: 'right-0 top-0 h-full',
    top: 'top-0 left-0 w-full',
    bottom: 'bottom-0 left-0 w-full',
  }

  const translateClasses = {
    left: { from: '-translate-x-full', to: 'translate-x-0' },
    right: { from: 'translate-x-full', to: 'translate-x-0' },
    top: { from: '-translate-y-full', to: 'translate-y-0' },
    bottom: { from: 'translate-y-full', to: 'translate-y-0' },
  }

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <div className="relative z-50">
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black bg-opacity-50" onClick={onClose} />
        </Transition.Child>

        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom={translateClasses[side].from}
          enterTo={translateClasses[side].to}
          leave="ease-in duration-200"
          leaveFrom={translateClasses[side].to}
          leaveTo={translateClasses[side].from}
        >
          <div
            className={cn(
              'fixed bg-white shadow-xl',
              sideClasses[side],
              (side === 'left' || side === 'right') && 'w-96 max-w-full',
              (side === 'top' || side === 'bottom') && 'h-96 max-h-full'
            )}
          >
            <div className="flex items-center justify-between p-6 border-b border-gray-200">
              <h3 className="text-xl font-semibold text-gray-900">{title}</h3>
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 transition-colors"
                aria-label="Fechar"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-6 overflow-y-auto" style={{ maxHeight: 'calc(100% - 80px)' }}>
              {children}
            </div>
          </div>
        </Transition.Child>
      </div>
    </Transition>
  )
}
