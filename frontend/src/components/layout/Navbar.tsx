import { useEffect, useRef, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Bell, Search, ChevronDown, LogOut, User, Menu, X, Shield, Check, Settings } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Avatar } from '@/components/ui/avatar'
import { useAuth } from '@/context/AuthContext'
import * as api from '@/api'
import type { Notification } from '@/types'
import { formatDate } from '@/api'

const NOTIF_LABELS: Record<string, string> = {
  NEW_ANSWER: 'Новый ответ на ваш вопрос',
  ANSWER_VERIFIED: 'Ваш ответ подтверждён',
  QUESTION_ASSIGNED: 'Вам назначен вопрос',
  QUESTION_CLOSED: 'Ваш вопрос закрыт',
  NEW_COMMENT: 'Комментарий к вашему ответу',
}

export function Navbar() {
  const { user, logout } = useAuth()
  const [searchQuery, setSearchQuery] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const [userMenuOpen, setUserMenuOpen] = useState(false)
  const [bellOpen, setBellOpen] = useState(false)
  const [unreadCount, setUnreadCount] = useState(0)
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [notifsLoading, setNotifsLoading] = useState(false)
  const bellRef = useRef<HTMLDivElement>(null)
  const navigate = useNavigate()

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      navigate(`/search?q=${encodeURIComponent(searchQuery.trim())}`)
    }
  }

  const handleLogout = async () => {
    setUserMenuOpen(false)
    await logout()
    navigate('/login')
  }

  useEffect(() => {
    if (!user) return
    const fetchCount = () => api.getUnreadCount().then((r) => setUnreadCount(r.count)).catch(() => {})
    fetchCount()
    const timer = setInterval(fetchCount, 30_000)
    return () => clearInterval(timer)
  }, [user])

  const openBell = async () => {
    if (bellOpen) { setBellOpen(false); return }
    setBellOpen(true)
    setNotifsLoading(true)
    try {
      const items = await api.getNotifications(1)
      setNotifications(items.slice(0, 7))
    } catch { /* ignore */ } finally {
      setNotifsLoading(false)
    }
  }

  const handleMarkRead = async (n: Notification) => {
    if (!n.is_read) {
      await api.markNotificationRead(n.id).catch(() => {})
      setNotifications((prev) => prev.map((x) => x.id === n.id ? { ...x, is_read: true } : x))
      setUnreadCount((c) => Math.max(0, c - 1))
    }
    setBellOpen(false)
    if (n.payload.question_id) navigate(`/questions/${n.payload.question_id}`)
  }

  const handleMarkAll = async () => {
    await api.markAllNotificationsRead().catch(() => {})
    setNotifications((prev) => prev.map((x) => ({ ...x, is_read: true })))
    setUnreadCount(0)
  }

  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-14 items-center gap-4">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2 shrink-0">
            <span className="font-semibold text-foreground">
              ИСС АПК <span className="text-primary">«ЭЛЬБРУС»</span>
            </span>
          </Link>

          {/* Search */}
          <form onSubmit={handleSearch} className="flex-1 max-w-xl">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <input
                type="text"
                placeholder="Поиск по базе знаний..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full h-9 pl-9 pr-4 rounded-md border border-border bg-background text-sm placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring transition-colors"
              />
            </div>
          </form>

          {/* Desktop right */}
          <div className="hidden sm:flex items-center gap-2 ml-auto">
            {user ? (
              <>
                <Button variant="default" size="sm" onClick={() => navigate('/questions/new')}>
                  Задать вопрос
                </Button>

                {/* Bell */}
                <div ref={bellRef} className="relative">
                  <Button variant="ghost" size="icon" className="relative" onClick={openBell}>
                    <Bell className="h-4 w-4" />
                    {unreadCount > 0 && (
                      <span className="absolute -top-0.5 -right-0.5 h-4 w-4 rounded-full bg-destructive text-destructive-foreground text-[10px] font-bold flex items-center justify-center">
                        {unreadCount > 9 ? '9+' : unreadCount}
                      </span>
                    )}
                  </Button>

                  {bellOpen && (
                    <>
                      <div className="fixed inset-0 z-30" onClick={() => setBellOpen(false)} />
                      <div className="absolute right-0 top-full mt-1 z-40 w-80 rounded-md border border-border bg-background shadow-lg">
                        <div className="flex items-center justify-between px-3 py-2 border-b border-border">
                          <span className="text-sm font-medium">Уведомления</span>
                          <div className="flex items-center gap-2">
                            {unreadCount > 0 && (
                              <button onClick={handleMarkAll} className="text-xs text-muted-foreground hover:text-foreground flex items-center gap-1">
                                <Check className="h-3 w-3" /> Прочитать все
                              </button>
                            )}
                            <Link to="/notifications" className="text-xs text-primary hover:underline" onClick={() => setBellOpen(false)}>
                              Все
                            </Link>
                          </div>
                        </div>
                        <div className="max-h-80 overflow-y-auto divide-y divide-border">
                          {notifsLoading ? (
                            <div className="px-3 py-6 text-center text-sm text-muted-foreground">Загрузка...</div>
                          ) : notifications.length === 0 ? (
                            <div className="px-3 py-6 text-center text-sm text-muted-foreground">Нет уведомлений</div>
                          ) : notifications.map((n) => (
                            <button
                              key={n.id}
                              onClick={() => handleMarkRead(n)}
                              className={`w-full text-left px-3 py-2.5 hover:bg-accent transition-colors ${!n.is_read ? 'bg-primary/5' : ''}`}
                            >
                              <div className="flex items-start gap-2">
                                {!n.is_read && <span className="mt-1.5 h-2 w-2 rounded-full bg-primary shrink-0" />}
                                <div className={!n.is_read ? '' : 'ml-4'}>
                                  <p className="text-xs font-medium text-foreground">{NOTIF_LABELS[n.type] ?? n.type}</p>
                                  {n.payload.question_title && (
                                    <p className="text-xs text-muted-foreground truncate max-w-[220px]">{n.payload.question_title}</p>
                                  )}
                                  {n.payload.actor_username && (
                                    <p className="text-xs text-muted-foreground">от {n.payload.actor_username}</p>
                                  )}
                                  <p className="text-[10px] text-muted-foreground mt-0.5">{formatDate(n.created_at)}</p>
                                </div>
                              </div>
                            </button>
                          ))}
                        </div>
                      </div>
                    </>
                  )}
                </div>

                {/* User menu */}
                <div className="relative">
                  <button
                    onClick={() => setUserMenuOpen(!userMenuOpen)}
                    className="flex items-center gap-2 rounded-md px-2 py-1 hover:bg-accent transition-colors"
                  >
                    <Avatar name={user.username} size="sm" />
                    <span className="text-sm font-medium hidden lg:block">{user.username}</span>
                    <ChevronDown className="h-3.5 w-3.5 text-muted-foreground" />
                  </button>

                  {userMenuOpen && (
                    <>
                      <div className="fixed inset-0 z-30" onClick={() => setUserMenuOpen(false)} />
                      <div className="absolute right-0 top-full mt-1 z-40 w-48 rounded-md border border-border bg-background shadow-lg py-1">
                        <div className="px-3 py-2 border-b border-border">
                          <p className="text-sm font-medium">{user.username}</p>
                          <p className="text-xs text-muted-foreground capitalize">{user.role.toLowerCase()}</p>
                        </div>
                        <Link
                          to={`/users/${user.id}`}
                          className="flex items-center gap-2 px-3 py-2 text-sm hover:bg-accent transition-colors"
                          onClick={() => setUserMenuOpen(false)}
                        >
                          <User className="h-4 w-4" /> Профиль
                        </Link>
                        <Link
                          to="/notifications"
                          className="flex items-center gap-2 px-3 py-2 text-sm hover:bg-accent transition-colors"
                          onClick={() => setUserMenuOpen(false)}
                        >
                          <Bell className="h-4 w-4" /> Уведомления
                          {unreadCount > 0 && (
                            <span className="ml-auto rounded-full bg-destructive text-destructive-foreground text-[10px] font-bold px-1.5 py-0.5">
                              {unreadCount}
                            </span>
                          )}
                        </Link>
                        <Link
                          to="/settings"
                          className="flex items-center gap-2 px-3 py-2 text-sm hover:bg-accent transition-colors"
                          onClick={() => setUserMenuOpen(false)}
                        >
                          <Settings className="h-4 w-4" /> Настройки
                        </Link>
                        {user.role === 'ADMIN' && (
                          <Link
                            to="/admin"
                            className="flex items-center gap-2 px-3 py-2 text-sm hover:bg-accent transition-colors"
                            onClick={() => setUserMenuOpen(false)}
                          >
                            <Shield className="h-4 w-4" /> Администрирование
                          </Link>
                        )}
                        <div className="border-t border-border mt-1">
                          <button
                            onClick={handleLogout}
                            className="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
                          >
                            <LogOut className="h-4 w-4" /> Выйти
                          </button>
                        </div>
                      </div>
                    </>
                  )}
                </div>
              </>
            ) : (
              <Button variant="default" size="sm" onClick={() => navigate('/login')}>
                Войти
              </Button>
            )}
          </div>

          {/* Mobile menu button */}
          <Button
            variant="ghost"
            size="icon"
            className="sm:hidden ml-auto"
            onClick={() => setMenuOpen(!menuOpen)}
          >
            {menuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </Button>
        </div>
      </div>

      {/* Mobile menu */}
      {menuOpen && (
        <div className="sm:hidden border-t border-border bg-background px-4 py-3 space-y-2">
          {user ? (
            <>
              <Button variant="default" size="sm" className="w-full" onClick={() => { navigate('/questions/new'); setMenuOpen(false) }}>
                Задать вопрос
              </Button>
              <Link to={`/users/${user.id}`} className="flex items-center gap-2 py-2 text-sm" onClick={() => setMenuOpen(false)}>
                <User className="h-4 w-4" /> Профиль ({user.username})
              </Link>
              <Link to="/notifications" className="flex items-center gap-2 py-2 text-sm" onClick={() => setMenuOpen(false)}>
                <Bell className="h-4 w-4" /> Уведомления
                {unreadCount > 0 && (
                  <span className="ml-1 rounded-full bg-destructive text-destructive-foreground text-[10px] font-bold px-1.5">
                    {unreadCount}
                  </span>
                )}
              </Link>
              <Link to="/settings" className="flex items-center gap-2 py-2 text-sm" onClick={() => setMenuOpen(false)}>
                <Settings className="h-4 w-4" /> Настройки
              </Link>
              {user.role === 'ADMIN' && (
                <Link to="/admin" className="flex items-center gap-2 py-2 text-sm" onClick={() => setMenuOpen(false)}>
                  <Shield className="h-4 w-4" /> Администрирование
                </Link>
              )}
              <button onClick={handleLogout} className="flex items-center gap-2 py-2 text-sm text-red-600">
                <LogOut className="h-4 w-4" /> Выйти
              </button>
            </>
          ) : (
            <Button variant="default" size="sm" className="w-full" onClick={() => { navigate('/login'); setMenuOpen(false) }}>
              Войти
            </Button>
          )}
        </div>
      )}
    </header>
  )
}
