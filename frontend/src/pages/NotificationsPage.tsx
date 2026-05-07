import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Bell, Check, CheckCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/context/AuthContext'
import * as api from '@/api'
import { formatDate } from '@/api'
import type { Notification } from '@/types'

const NOTIF_LABELS: Record<string, string> = {
  NEW_ANSWER: 'Новый ответ на ваш вопрос',
  ANSWER_VERIFIED: 'Ваш ответ подтверждён как достоверный',
  QUESTION_ASSIGNED: 'Вам назначен вопрос',
  QUESTION_CLOSED: 'Ваш вопрос закрыт',
  NEW_COMMENT: 'Комментарий к вашему ответу',
}

const NOTIF_DESCRIPTIONS: Record<string, (p: Notification['payload']) => string> = {
  NEW_ANSWER: (p) => `${p.actor_username} ответил на вопрос «${p.question_title}»`,
  ANSWER_VERIFIED: (p) => `${p.actor_username} отметил ваш ответ как достоверный в вопросе «${p.question_title}»`,
  QUESTION_ASSIGNED: (p) => `${p.actor_username} задал вопрос «${p.question_title}», который назначен вам`,
  QUESTION_CLOSED: (p) => `${p.actor_username} закрыл вопрос «${p.question_title}»`,
  NEW_COMMENT: (p) => `${p.actor_username} оставил комментарий к вашему ответу`,
}

export function NotificationsPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [hasMore, setHasMore] = useState(false)
  const [marking, setMarking] = useState(false)

  if (!user) {
    navigate('/login')
    return null
  }

  const load = (p: number, append = false) => {
    setLoading(true)
    api.getNotifications(p)
      .then((items) => {
        setNotifications((prev) => append ? [...prev, ...items] : items)
        setHasMore(items.length === 20)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }

  useEffect(() => { load(1) }, [])

  const handleMarkRead = async (n: Notification) => {
    if (!n.is_read) {
      await api.markNotificationRead(n.id).catch(() => {})
      setNotifications((prev) => prev.map((x) => x.id === n.id ? { ...x, is_read: true } : x))
    }
    if (n.payload.question_id) navigate(`/questions/${n.payload.question_id}`)
  }

  const handleMarkAll = async () => {
    setMarking(true)
    await api.markAllNotificationsRead().catch(() => {})
    setNotifications((prev) => prev.map((x) => ({ ...x, is_read: true })))
    setMarking(false)
  }

  const loadMore = () => {
    const next = page + 1
    setPage(next)
    load(next, true)
  }

  const unreadCount = notifications.filter((n) => !n.is_read).length

  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div className="flex items-center gap-2">
          <Bell className="h-5 w-5 text-muted-foreground" />
          <h1 className="text-xl font-bold">Уведомления</h1>
          {unreadCount > 0 && (
            <span className="rounded-full bg-destructive text-destructive-foreground text-xs font-bold px-2 py-0.5">
              {unreadCount}
            </span>
          )}
        </div>
        {unreadCount > 0 && (
          <Button variant="outline" size="sm" onClick={handleMarkAll} disabled={marking}>
            <CheckCheck className="h-4 w-4" />
            Прочитать все
          </Button>
        )}
      </div>

      {loading && notifications.length === 0 ? (
        <div className="space-y-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <div key={i} className="h-16 rounded-lg border border-border bg-background animate-pulse" />
          ))}
        </div>
      ) : notifications.length === 0 ? (
        <div className="text-center py-16 text-muted-foreground">
          <Bell className="h-10 w-10 mx-auto mb-3 opacity-30" />
          <p className="font-medium">Нет уведомлений</p>
          <p className="text-sm mt-1">Здесь будут появляться уведомления об ответах и комментариях</p>
        </div>
      ) : (
        <div className="space-y-1">
          {notifications.map((n) => (
            <button
              key={n.id}
              onClick={() => handleMarkRead(n)}
              className={`w-full text-left rounded-lg border border-border px-4 py-3 hover:bg-accent transition-colors ${
                !n.is_read ? 'bg-primary/5 border-primary/20' : 'bg-background'
              }`}
            >
              <div className="flex items-start gap-3">
                {!n.is_read ? (
                  <span className="mt-1.5 h-2 w-2 rounded-full bg-primary shrink-0" />
                ) : (
                  <Check className="mt-1 h-4 w-4 text-muted-foreground/40 shrink-0" />
                )}
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-foreground">{NOTIF_LABELS[n.type] ?? n.type}</p>
                  <p className="text-sm text-muted-foreground mt-0.5 truncate">
                    {NOTIF_DESCRIPTIONS[n.type]?.(n.payload) ?? ''}
                  </p>
                  <p className="text-xs text-muted-foreground mt-1">{formatDate(n.created_at)}</p>
                </div>
              </div>
            </button>
          ))}
        </div>
      )}

      {hasMore && (
        <div className="mt-4 text-center">
          <Button variant="outline" size="sm" onClick={loadMore} disabled={loading}>
            {loading ? 'Загрузка...' : 'Показать ещё'}
          </Button>
        </div>
      )}
    </div>
  )
}
