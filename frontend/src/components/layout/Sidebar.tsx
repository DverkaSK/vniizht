import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Layers, Tag } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { Category, Tag as TagType } from '@/types'
import { getCategories, getTags } from '@/api'

export function Sidebar() {
  const [searchParams] = useSearchParams()
  const activeCategory = searchParams.get('category')
  const activeTag = searchParams.get('tag')
  const [categories, setCategories] = useState<Category[]>([])
  const [tags, setTags] = useState<TagType[]>([])

  useEffect(() => {
    getCategories().then(setCategories).catch(() => {})
    getTags().then(setTags).catch(() => {})
  }, [])

  return (
    <aside className="w-56 shrink-0 space-y-6">
      {/* Categories */}
      <div>
        <div className="flex items-center gap-2 mb-2 px-1">
          <Layers className="h-4 w-4 text-muted-foreground" />
          <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">
            Категории
          </h3>
        </div>
        <nav className="space-y-0.5">
          <Link
            to="/"
            className={cn(
              'flex items-center justify-between px-3 py-1.5 rounded-md text-sm transition-colors',
              !activeCategory && !activeTag
                ? 'bg-primary text-primary-foreground font-medium'
                : 'hover:bg-accent text-foreground'
            )}
          >
            Все вопросы
          </Link>
          {categories.map((cat) => (
            <Link
              key={cat.id}
              to={`/?category=${cat.id}`}
              className={cn(
                'flex items-center justify-between px-3 py-1.5 rounded-md text-sm transition-colors',
                activeCategory === String(cat.id)
                  ? 'bg-primary text-primary-foreground font-medium'
                  : 'hover:bg-accent text-foreground'
              )}
            >
              <span className="truncate">{cat.name}</span>
            </Link>
          ))}
        </nav>
      </div>

      {/* Tags */}
      {tags.length > 0 && (
        <div>
          <div className="flex items-center gap-2 mb-2 px-1">
            <Tag className="h-4 w-4 text-muted-foreground" />
            <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">
              Теги
            </h3>
          </div>
          <div className="flex flex-wrap gap-1.5 px-1">
            {tags.slice(0, 8).map((tag) => (
              <Link
                key={tag.id}
                to={`/?tag=${tag.id}`}
                className={cn(
                  'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium transition-colors border',
                  activeTag === String(tag.id)
                    ? 'bg-primary text-primary-foreground border-primary'
                    : 'bg-background text-muted-foreground border-border hover:border-primary hover:text-primary'
                )}
              >
                {tag.name}
              </Link>
            ))}
          </div>
        </div>
      )}
    </aside>
  )
}
