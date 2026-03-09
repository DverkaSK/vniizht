import { cn } from '@/lib/utils'

interface AvatarProps {
  name: string
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

function getInitials(name: string) {
  return name
    .split(' ')
    .slice(0, 2)
    .map((n) => n[0])
    .join('')
    .toUpperCase()
}

function getColor(name: string) {
  const colors = [
    'bg-blue-500', 'bg-violet-500', 'bg-emerald-500',
    'bg-amber-500', 'bg-rose-500', 'bg-cyan-500',
  ]
  const index = name.charCodeAt(0) % colors.length
  return colors[index]
}

export function Avatar({ name, size = 'md', className }: AvatarProps) {
  return (
    <div
      className={cn(
        'rounded-full flex items-center justify-center text-white font-semibold shrink-0',
        getColor(name),
        {
          'h-7 w-7 text-xs': size === 'sm',
          'h-9 w-9 text-sm': size === 'md',
          'h-14 w-14 text-lg': size === 'lg',
        },
        className
      )}
    >
      {getInitials(name)}
    </div>
  )
}
