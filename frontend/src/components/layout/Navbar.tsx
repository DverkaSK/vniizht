import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Bell, Search, ChevronDown, LogOut, User, Menu, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Avatar } from '@/components/ui/avatar'
import { useAuth } from '@/context/AuthContext'

export function Navbar() {
  const { user, logout } = useAuth()
  const [searchQuery, setSearchQuery] = useState('')
  const [menuOpen, setMenuOpen] = useState(false)
  const [userMenuOpen, setUserMenuOpen] = useState(false)
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

  return (
    <header className="sticky top-0 z-40 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-14 items-center gap-4">
          {/* Logo */}
          <Link to="/" className="flex items-center gap-2 shrink-0">
            <div className="h-8 w-8 rounded-md bg-primary flex items-center justify-center">
              <span className="text-white font-bold text-sm">Э</span>
            </div>
            <span className="hidden sm:block font-semibold text-foreground">
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

                <Button variant="ghost" size="icon" className="relative">
                  <Bell className="h-4 w-4" />
                </Button>

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
