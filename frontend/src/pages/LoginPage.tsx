import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Eye, EyeOff, LogIn } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/context/AuthContext'

export function LoginPage() {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [showPassword, setShowPassword] = useState(false)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(username, password)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка входа')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-muted/30 flex items-center justify-center px-4">
      <div className="w-full max-w-sm">
        {/* Logo */}
        <div className="text-center mb-8">
          <div className="h-12 w-12 rounded-xl bg-primary flex items-center justify-center mx-auto mb-3">
            <span className="text-white font-bold text-xl">Э</span>
          </div>
          <h1 className="text-xl font-bold text-foreground">ИСС АПК «ЭЛЬБРУС»</h1>
          <p className="text-sm text-muted-foreground mt-1">База знаний технической поддержки</p>
        </div>

        {/* Card */}
        <div className="bg-background rounded-xl border border-border shadow-sm p-6 space-y-4">
          <h2 className="text-base font-semibold text-center">Вход в систему</h2>

          <form onSubmit={handleSubmit} className="space-y-3">
            <div>
              <label className="text-sm font-medium text-foreground block mb-1.5">
                Логин
              </label>
              <input
                type="text"
                placeholder="Введите логин"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                className="w-full h-10 px-3 rounded-md border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-ring transition-colors"
                autoComplete="username"
              />
            </div>

            <div>
              <label className="text-sm font-medium text-foreground block mb-1.5">
                Пароль
              </label>
              <div className="relative">
                <input
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Введите пароль"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full h-10 pl-3 pr-10 rounded-md border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-ring transition-colors"
                  autoComplete="current-password"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
                >
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
            </div>

            {error && (
              <p className="text-sm text-red-600">{error}</p>
            )}

            <Button type="submit" className="w-full" disabled={!username || !password || loading}>
              <LogIn className="h-4 w-4" />
              {loading ? 'Вход...' : 'Войти'}
            </Button>
          </form>

          <p className="text-xs text-center text-muted-foreground">
            Нет учётной записи?{' '}
            <span className="text-foreground font-medium">
              Обратитесь к администратору системы
            </span>
          </p>
        </div>

        {/* Guest link */}
        <p className="text-center mt-4 text-sm">
          <Link to="/" className="text-primary hover:underline">
            Продолжить без входа →
          </Link>
        </p>
      </div>
    </div>
  )
}
