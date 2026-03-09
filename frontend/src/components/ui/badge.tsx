import { cn } from '@/lib/utils'

interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'secondary' | 'outline' | 'verified' | 'status-open' | 'status-progress' | 'status-closed'
}

export function Badge({ className, variant = 'default', ...props }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors',
        {
          'bg-primary text-primary-foreground': variant === 'default',
          'bg-secondary text-secondary-foreground': variant === 'secondary',
          'border border-border text-foreground': variant === 'outline',
          'bg-green-100 text-green-800 border border-green-200': variant === 'verified',
          'bg-blue-100 text-blue-700 border border-blue-200': variant === 'status-open',
          'bg-amber-100 text-amber-700 border border-amber-200': variant === 'status-progress',
          'bg-gray-100 text-gray-600 border border-gray-200': variant === 'status-closed',
        },
        className
      )}
      {...props}
    />
  )
}