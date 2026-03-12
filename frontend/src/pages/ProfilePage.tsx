import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { MessageSquare, HelpCircle, Star, Calendar, CheckCircle2 } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Avatar } from '@/components/ui/avatar'
import { formatDate, formatRelative } from '@/lib/utils'
import { getUser } from '@/api'
import type { PublicUser } from '@/types'

const roleLabel: Record<string, string> = {
  user: 'Пользователь',
  specialist: 'Специалист АО ВНИИЖТ',
  admin: 'Администратор',
  guest: 'Гость',
}

const roleVariant: Record<string, 'default' | 'verified' | 'secondary'> = {
  user: 'secondary',
  specialist: 'verified',
  admin: 'default',
  guest: 'secondary',
}

export function ProfilePage() {
  const { id } = useParams<{ id: string }>()
  const [user, setUser] = useState<PublicUser | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!id) return
    getUser(parseInt(id))
      .then(setUser)
      .catch(() => setUser(null))
      .finally(() => setLoading(false))
  }, [id])

  if (loading) {
    return (
      <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 py-6">
        <div className="animate-pulse h-48 rounded-lg border border-border bg-muted" />
      </div>
    )
  }

  if (!user) {
    return (
      <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 py-6 text-center text-muted-foreground py-16">
        Пользователь не найден
      </div>
    )
  }

  return (
    <div className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 py-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Sidebar: user info */}
        <div className="space-y-4">
          <div className="rounded-lg border border-border bg-background p-5 text-center">
            <Avatar name={user.username} size="lg" className="mx-auto mb-3" />
            <h1 className="text-lg font-bold text-foreground">{user.username}</h1>
            <div className="mt-2">
              <Badge variant={roleVariant[user.role]}>{roleLabel[user.role]}</Badge>
            </div>
          </div>

          {/* Stats */}
          <div className="rounded-lg border border-border bg-background p-4 space-y-3">
            <h3 className="text-sm font-semibold text-muted-foreground uppercase tracking-wide">Статистика</h3>
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Star className="h-4 w-4 text-amber-500" />
                  Репутация
                </div>
                <span className="font-bold text-foreground">{user.reputation}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <HelpCircle className="h-4 w-4 text-blue-500" />
                  Вопросов
                </div>
                <span className="font-semibold">{user.question_count}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <MessageSquare className="h-4 w-4 text-green-500" />
                  Ответов
                </div>
                <span className="font-semibold">{user.answer_count}</span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Calendar className="h-4 w-4" />
                  В системе с
                </div>
                <span className="text-sm">{formatDate(user.created_at)}</span>
              </div>
            </div>
          </div>

          {user.role === 'SPECIALIST' && (
            <div className="rounded-lg border border-green-200 bg-green-50/50 p-4">
              <div className="flex items-center gap-1.5 text-green-700 font-medium text-sm mb-1">
                <CheckCircle2 className="h-4 w-4" />
                Специалист АО ВНИИЖТ
              </div>
              <p className="text-xs text-green-600">
                Верифицирует ответы и обрабатывает обращения пользователей
              </p>
            </div>
          )}
        </div>

        {/* Main: activity */}
        <div className="md:col-span-2 space-y-5">
          {/* Recent questions */}
          <div>
            <h2 className="text-base font-semibold mb-3">Последние вопросы</h2>
            {user.recent_questions.length === 0 ? (
              <div className="rounded-lg border border-border bg-background p-6 text-center text-muted-foreground text-sm">
                Вопросов пока нет
              </div>
            ) : (
              <div className="space-y-2">
                {user.recent_questions.map((q) => (
                  <Link
                    key={q.id}
                    to={`/questions/${q.id}`}
                    className="flex items-center justify-between rounded-lg border border-border bg-background px-4 py-3 hover:border-primary/40 transition-colors"
                  >
                    <span className="text-sm font-medium line-clamp-1">{q.title}</span>
                    <span className="text-xs text-muted-foreground ml-3 shrink-0">{formatRelative(q.created_at)}</span>
                  </Link>
                ))}
              </div>
            )}
          </div>

          {/* Recent answers */}
          <div>
            <h2 className="text-base font-semibold mb-3">Последние ответы</h2>
            {user.recent_answers.length === 0 ? (
              <div className="rounded-lg border border-border bg-background p-6 text-center text-muted-foreground text-sm">
                Ответов пока нет
              </div>
            ) : (
              <div className="space-y-2">
                {user.recent_answers.map((a) => (
                  <Link
                    key={a.id}
                    to={`/questions/${a.question_id}`}
                    className="flex items-center justify-between rounded-lg border border-border bg-background px-4 py-3 hover:border-primary/40 transition-colors"
                  >
                    <span className="text-sm text-muted-foreground line-clamp-1">{a.snippet}</span>
                    <span className="text-xs text-muted-foreground ml-3 shrink-0">{formatRelative(a.created_at)}</span>
                  </Link>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
