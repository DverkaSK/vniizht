import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Search, FileText, MessageSquare } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { formatRelative } from '@/lib/utils'
import { search } from '@/api'
import type { SearchResult } from '@/types'

const statusLabel: Record<string, string> = {
  open: 'Открыт',
  closed: 'Закрыт',
  duplicate: 'Дубликат',
}
const statusVariant: Record<string, 'status-open' | 'status-progress' | 'status-closed'> = {
  open: 'status-open',
  closed: 'status-closed',
  duplicate: 'status-closed',
}

export function SearchPage() {
  const [searchParams] = useSearchParams()
  const query = searchParams.get('q') ?? ''
  const [results, setResults] = useState<SearchResult[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!query.trim()) {
      setResults([])
      setTotal(0)
      return
    }
    setLoading(true)
    search({ q: query })
      .then((res) => {
        setResults(res.items ?? [])
        setTotal(res.total)
      })
      .catch(() => setResults([]))
      .finally(() => setLoading(false))
  }, [query])

  return (
    <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 py-6">
      <div className="flex items-center gap-3 mb-6">
        <Search className="h-5 w-5 text-muted-foreground shrink-0" />
        <h1 className="text-xl font-bold text-foreground">
          {query ? `Результаты поиска: «${query}»` : 'Поиск'}
        </h1>
      </div>

      {!query.trim() ? (
        <div className="text-center py-16 text-muted-foreground">
          <p>Введите запрос в строку поиска</p>
        </div>
      ) : loading ? (
        <div className="space-y-3">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="rounded-lg border border-border bg-background p-4 animate-pulse h-20" />
          ))}
        </div>
      ) : results.length === 0 ? (
        <div className="text-center py-16 text-muted-foreground">
          <p className="text-lg font-medium mb-1">Ничего не найдено</p>
          <p className="text-sm">Попробуйте изменить запрос</p>
        </div>
      ) : (
        <>
          <p className="text-sm text-muted-foreground mb-4">
            Найдено: <span className="font-medium text-foreground">{total}</span> результатов
          </p>
          <div className="space-y-3">
            {results.map((result, i) => (
              <Link
                key={i}
                to={result.type === 'answer'
                  ? `/questions/${result.question_id}#answer-${result.answer_id}`
                  : `/questions/${result.question_id}`
                }
                className="block rounded-lg border border-border bg-background p-4 hover:border-primary/40 hover:shadow-sm transition-all"
              >
                <div className="flex items-start gap-3">
                  <div className="shrink-0 mt-0.5">
                    {result.type === 'question' ? (
                      <FileText className="h-4 w-4 text-blue-500" />
                    ) : (
                      <MessageSquare className="h-4 w-4 text-green-500" />
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1 flex-wrap">
                      <span className="text-xs text-muted-foreground">
                        {result.type === 'question' ? 'Вопрос' : 'Ответ'}
                      </span>
                      <Badge variant={statusVariant[result.status]} className="text-[10px] py-0">
                        {statusLabel[result.status]}
                      </Badge>
                    </div>
                    <h3 className="text-sm font-semibold text-foreground mb-1 line-clamp-1">
                      {result.question_title}
                    </h3>
                    {result.snippet && (
                      <p className="text-xs text-muted-foreground line-clamp-2">{result.snippet}</p>
                    )}
                    <p className="text-xs text-muted-foreground mt-1">{formatRelative(result.created_at)}</p>
                  </div>
                </div>
              </Link>
            ))}
          </div>
        </>
      )}
    </div>
  )
}
