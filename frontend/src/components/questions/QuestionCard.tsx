import { Link } from 'react-router-dom'
import { MessageSquare, CheckCircle2, Clock } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Avatar } from '@/components/ui/avatar'
import { formatRelative } from '@/lib/utils'
import type { Question } from '@/types'

interface QuestionCardProps {
  question: Question
  categoryName?: string
}

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

export function QuestionCard({ question, categoryName }: QuestionCardProps) {
  return (
    <article className="group rounded-lg border border-border bg-background p-4 hover:border-primary/40 hover:shadow-sm transition-all">
      <div className="flex gap-4">
        {/* Stats column */}
        <div className="hidden sm:flex flex-col items-center gap-3 shrink-0 w-14 pt-0.5">
          <div className="text-center">
            <div className={`text-lg font-bold ${question.answer_count > 0 ? 'text-foreground' : 'text-muted-foreground'}`}>
              {question.answer_count}
            </div>
            <div className="text-[11px] text-muted-foreground leading-tight">
              {question.answer_count === 1 ? 'ответ' : 'ответов'}
            </div>
          </div>
          {question.has_verified && (
            <div className="flex flex-col items-center">
              <CheckCircle2 className="h-5 w-5 text-green-500" />
              <span className="text-[10px] text-green-600 leading-tight mt-0.5">верифицирован</span>
            </div>
          )}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start gap-2 flex-wrap mb-1.5">
            <Badge variant={statusVariant[question.status]}>
              {statusLabel[question.status]}
            </Badge>
            {question.has_verified && (
              <Badge variant="verified" className="sm:hidden">
                ✓ Верифицирован
              </Badge>
            )}
          </div>

          <Link to={`/questions/${question.id}`}>
            <h2 className="text-base font-semibold text-foreground group-hover:text-primary transition-colors line-clamp-2 mb-1.5">
              {question.title}
            </h2>
          </Link>

          <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
            {question.body}
          </p>

          {/* Tags */}
          {question.tags.length > 0 && (
            <div className="flex flex-wrap gap-1.5 mb-3">
              {question.tags.map((tag) => (
                <Link
                  key={tag.id}
                  to={`/?tag=${tag.id}`}
                  className="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium bg-accent text-muted-foreground border border-border hover:border-primary hover:text-primary transition-colors"
                >
                  {tag.name}
                </Link>
              ))}
            </div>
          )}

          {/* Footer */}
          <div className="flex items-center justify-between flex-wrap gap-2">
            <div className="flex items-center gap-1.5">
              <Avatar name={question.author_username} size="sm" />
              <Link
                to={`/users/${question.author_id}`}
                className="text-sm font-medium hover:text-primary transition-colors"
              >
                {question.author_username}
              </Link>
            </div>

            <div className="flex items-center gap-3 text-xs text-muted-foreground">
              <span className="sm:hidden flex items-center gap-1">
                <MessageSquare className="h-3.5 w-3.5" />
                {question.answer_count}
              </span>
              {categoryName && (
                <span className="hidden sm:block">{categoryName}</span>
              )}
              <span className="flex items-center gap-1">
                <Clock className="h-3.5 w-3.5" />
                {formatRelative(question.created_at)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </article>
  )
}
