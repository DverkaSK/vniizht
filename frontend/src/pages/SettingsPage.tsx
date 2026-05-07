import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Bell, Mail, Save } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/context/AuthContext'
import * as api from '@/api'
import type { NotificationPreference, NotificationType } from '@/types'

const ALL_TYPES: NotificationType[] = [
  'NEW_ANSWER',
  'ANSWER_VERIFIED',
  'QUESTION_ASSIGNED',
  'QUESTION_CLOSED',
  'NEW_COMMENT',
]

const TYPE_LABELS: Record<NotificationType, { title: string; desc: string }> = {
  NEW_ANSWER: {
    title: 'Новый ответ',
    desc: 'Кто-то ответил на ваш вопрос',
  },
  ANSWER_VERIFIED: {
    title: 'Ответ подтверждён',
    desc: 'Специалист отметил ваш ответ как достоверный',
  },
  QUESTION_ASSIGNED: {
    title: 'Назначен вопрос',
    desc: 'Вам назначен новый вопрос (для специалистов)',
  },
  QUESTION_CLOSED: {
    title: 'Вопрос закрыт',
    desc: 'Ваш вопрос был закрыт',
  },
  NEW_COMMENT: {
    title: 'Новый комментарий',
    desc: 'Кто-то прокомментировал ваш ответ',
  },
}

export function SettingsPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [prefs, setPrefs] = useState<Record<NotificationType, boolean>>(() =>
    Object.fromEntries(ALL_TYPES.map((t) => [t, true])) as Record<NotificationType, boolean>
  )
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)
  const [error, setError] = useState<string | null>(null)

  if (!user) {
    navigate('/login')
    return null
  }

  useEffect(() => {
    api.getNotificationPreferences()
      .then((items) => {
        const map = { ...prefs }
        for (const p of items) map[p.type] = p.email_enabled
        setPrefs(map)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const handleSave = async () => {
    setSaving(true)
    setError(null)
    setSaved(false)
    try {
      const payload: NotificationPreference[] = ALL_TYPES.map((t) => ({
        type: t,
        email_enabled: prefs[t],
      }))
      await api.saveNotificationPreferences(payload)
      setSaved(true)
      setTimeout(() => setSaved(false), 3000)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Ошибка сохранения')
    } finally {
      setSaving(false)
    }
  }

  const toggle = (type: NotificationType) =>
    setPrefs((prev) => ({ ...prev, [type]: !prev[type] }))

  return (
    <div className="mx-auto max-w-2xl px-4 py-8">
      <h1 className="text-xl font-bold mb-6">Настройки</h1>

      <div className="rounded-lg border border-border overflow-hidden">
        <div className="flex items-center gap-2 px-4 py-3 bg-accent/40 border-b border-border">
          <Bell className="h-4 w-4 text-muted-foreground" />
          <h2 className="font-semibold text-sm">Уведомления по email</h2>
        </div>

        <div className="divide-y divide-border">
          {loading ? (
            <div className="px-4 py-6 text-sm text-muted-foreground">Загрузка...</div>
          ) : (
            ALL_TYPES.map((type) => {
              const { title, desc } = TYPE_LABELS[type]
              return (
                <div key={type} className="flex items-center justify-between px-4 py-3 gap-4">
                  <div className="flex items-start gap-3 min-w-0">
                    <Mail className="h-4 w-4 text-muted-foreground mt-0.5 shrink-0" />
                    <div>
                      <p className="text-sm font-medium">{title}</p>
                      <p className="text-xs text-muted-foreground">{desc}</p>
                    </div>
                  </div>
                  <button
                    role="switch"
                    aria-checked={prefs[type]}
                    onClick={() => toggle(type)}
                    className={`relative inline-flex h-5 w-9 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors focus-visible:outline-none ${
                      prefs[type] ? 'bg-primary' : 'bg-input'
                    }`}
                  >
                    <span
                      className={`pointer-events-none inline-block h-4 w-4 rounded-full bg-white shadow-lg ring-0 transition-transform ${
                        prefs[type] ? 'translate-x-4' : 'translate-x-0'
                      }`}
                    />
                  </button>
                </div>
              )
            })
          )}
        </div>

        <div className="px-4 py-3 bg-accent/20 border-t border-border flex items-center gap-3">
          <Button size="sm" onClick={handleSave} disabled={saving || loading}>
            <Save className="h-4 w-4" />
            {saving ? 'Сохранение...' : 'Сохранить'}
          </Button>
          {saved && <span className="text-sm text-green-600">Настройки сохранены</span>}
          {error && <span className="text-sm text-destructive">{error}</span>}
        </div>
      </div>

      <p className="mt-4 text-xs text-muted-foreground">
        Внутрисистемные уведомления (колокольчик) приходят всегда, независимо от этих настроек.
      </p>
    </div>
  )
}
