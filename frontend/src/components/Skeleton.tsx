interface SkeletonProps {
  className?: string
  width?: string
  height?: string
}

export default function Skeleton({ className, width, height }: SkeletonProps) {
  return (
    <div
      className={`animate-pulse bg-gray-200 rounded ${className || ''}`}
      style={{ width, height }}
    />
  )
}

export function SkeletonCard() {
  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
      <Skeleton height="24px" width="40%" className="mb-4" />
      <Skeleton height="16px" width="60%" className="mb-2" />
      <Skeleton height="16px" width="80%" />
    </div>
  )
}

export function SkeletonTable({ rows = 5, columns = 4 }: { rows?: number; columns?: number }) {
  return (
    <div className="space-y-3">
      <div className="flex gap-4">
        {Array.from({ length: columns }).map((_, idx) => (
          <Skeleton key={idx} height="20px" width={`${100 / columns}%`} />
        ))}
      </div>
      {Array.from({ length: rows }).map((_, rowIdx) => (
        <div key={rowIdx} className="flex gap-4">
          {Array.from({ length: columns }).map((_, colIdx) => (
            <Skeleton key={colIdx} height="16px" width={`${100 / columns}%`} />
          ))}
        </div>
      ))}
    </div>
  )
}
