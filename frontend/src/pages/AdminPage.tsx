import { useEffect, useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Upload, Download, CheckCircle, AlertCircle,
  Plus, Pencil, Trash2, X, Check, Users, Layers, Tag, Database,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/context/AuthContext'
import * as api from '@/api'
import type { Category, Tag as TagType } from '@/types'

type AdminUser = {
  id: number
  username: string
  email: string
  role: string
  is_active: boolean
  reputation: number
}

type Tab = 'users' | 'categories' | 'tags' | 'data'

const ROLES = ['USER', 'SPECIALIST', 'ADMIN']
const ROLE_LABELS: Record<string, string> = {
  USER: 'Пользователь',
  SPECIALIST: 'Специалист',
  ADMIN: 'Администратор',
}

export function AdminPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [tab, setTab] = useState<Tab>('users')

  if (!user || user.role !== 'ADMIN') {
    navigate('/')
    return null
  }

  return (
    <div className="mx-auto max-w-5xl px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">Панель администратора</h1>

      <div className="flex gap-1 border-b border-border mb-6">
        {([
          ['users', 'Пользователи', Users],
          ['categories', 'Категории', Layers],
          ['tags', 'Теги', Tag],
          ['data', 'Импорт / Экспорт', Database],
        ] as [Tab, string, React.ElementType][]).map(([key, label, Icon]) => (
          <button
            key={key}
            onClick={() => setTab(key)}
            className={`flex items-center gap-1.5 px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
              tab === key
                ? 'border-primary text-primary'
                : 'border-transparent text-muted-foreground hover:text-foreground'
            }`}
          >
            <Icon className="h-4 w-4" />
            {label}
          </button>
        ))}
      </div>

      {tab === 'users' && <UsersTab />}
      {tab === 'categories' && <CategoriesTab />}
      {tab === 'tags' && <TagsTab />}
      {tab === 'data' && <DataTab />}
    </div>
  )
}

// ─── Users Tab ───────────────────────────────────────────────────────────────

function UsersTab() {
  const [users, setUsers] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)

  const [newUsername, setNewUsername] = useState('')
  const [newEmail, setNewEmail] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [newRole, setNewRole] = useState('USER')
  const [creating2, setCreating2] = useState(false)

  const load = () => {
    setLoading(true)
    api.adminListUsers()
      .then((u) => setUsers(u as unknown as AdminUser[]))
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(load, [])

  const handleRoleChange = async (id: number, role: string) => {
    try {
      await api.adminChangeRole(id, role)
      setUsers((prev) => prev.map((u) => u.id === id ? { ...u, role } : u))
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    }
  }

  const handleToggleActive = async (u: AdminUser) => {
    try {
      await api.adminSetUserActive(u.id, !u.is_active)
      setUsers((prev) => prev.map((x) => x.id === u.id ? { ...x, is_active: !x.is_active } : x))
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setCreating2(true)
    try {
      const created = await api.adminCreateUser({ username: newUsername, email: newEmail, password: newPassword, role: newRole })
      setUsers((prev) => [...prev, created as unknown as AdminUser])
      setNewUsername(''); setNewEmail(''); setNewPassword(''); setNewRole('USER')
      setCreating(false)
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setCreating2(false)
    }
  }

  if (loading) return <div className="text-sm text-muted-foreground">Загрузка...</div>
  if (error) return <ErrorBox message={error} />

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">Всего пользователей: {users.length}</p>
        <Button size="sm" onClick={() => setCreating((v) => !v)}>
          <Plus className="h-4 w-4" /> Создать пользователя
        </Button>
      </div>

      {creating && (
        <form onSubmit={handleCreate} className="rounded-lg border border-border p-4 space-y-3 bg-accent/30">
          <p className="font-medium text-sm">Новый пользователь</p>
          <div className="grid grid-cols-2 gap-3">
            <input
              required placeholder="Логин" value={newUsername}
              onChange={(e) => setNewUsername(e.target.value)}
              className="rounded-md border border-border bg-background px-3 py-1.5 text-sm w-full"
            />
            <input
              required placeholder="Email" type="email" value={newEmail}
              onChange={(e) => setNewEmail(e.target.value)}
              className="rounded-md border border-border bg-background px-3 py-1.5 text-sm w-full"
            />
            <input
              required placeholder="Пароль" type="password" value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              className="rounded-md border border-border bg-background px-3 py-1.5 text-sm w-full"
            />
            <select
              value={newRole} onChange={(e) => setNewRole(e.target.value)}
              className="rounded-md border border-border bg-background px-3 py-1.5 text-sm w-full"
            >
              {ROLES.map((r) => <option key={r} value={r}>{ROLE_LABELS[r]}</option>)}
            </select>
          </div>
          <div className="flex gap-2">
            <Button type="submit" size="sm" disabled={creating2}>
              {creating2 ? 'Создаётся...' : 'Создать'}
            </Button>
            <Button type="button" size="sm" variant="outline" onClick={() => setCreating(false)}>
              Отмена
            </Button>
          </div>
        </form>
      )}

      <div className="rounded-lg border border-border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-accent/50">
            <tr>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">ID</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Логин</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Email</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Роль</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Репутация</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Активен</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {users.map((u) => (
              <tr key={u.id} className={!u.is_active ? 'opacity-50' : ''}>
                <td className="px-3 py-2 text-muted-foreground">{u.id}</td>
                <td className="px-3 py-2 font-medium">{u.username}</td>
                <td className="px-3 py-2 text-muted-foreground">{u.email}</td>
                <td className="px-3 py-2">
                  <select
                    value={u.role}
                    onChange={(e) => handleRoleChange(u.id, e.target.value)}
                    className="rounded border border-border bg-background px-2 py-0.5 text-xs"
                  >
                    {ROLES.map((r) => <option key={r} value={r}>{ROLE_LABELS[r]}</option>)}
                  </select>
                </td>
                <td className="px-3 py-2 text-muted-foreground">{u.reputation}</td>
                <td className="px-3 py-2">
                  <button
                    onClick={() => handleToggleActive(u)}
                    className={`inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium transition-colors ${
                      u.is_active
                        ? 'bg-green-100 text-green-700 hover:bg-green-200'
                        : 'bg-red-100 text-red-700 hover:bg-red-200'
                    }`}
                  >
                    {u.is_active ? 'Да' : 'Нет'}
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

// ─── Categories Tab ──────────────────────────────────────────────────────────

function CategoriesTab() {
  const [categories, setCategories] = useState<Category[]>([])
  const [specialists, setSpecialists] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [showCreate, setShowCreate] = useState(false)

  const [editName, setEditName] = useState('')
  const [editDesc, setEditDesc] = useState('')
  const [editSpecialist, setEditSpecialist] = useState<number | ''>('')

  const [newName, setNewName] = useState('')
  const [newDesc, setNewDesc] = useState('')
  const [newSpecialist, setNewSpecialist] = useState<number | ''>('')
  const [saving, setSaving] = useState(false)

  const load = () => {
    setLoading(true)
    Promise.all([api.getCategories(), api.adminListUsers()])
      .then(([cats, users]) => {
        setCategories(cats)
        setSpecialists((users as unknown as AdminUser[]).filter((u) => u.role === 'SPECIALIST' || u.role === 'ADMIN'))
      })
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(load, [])

  const startEdit = (c: Category) => {
    setEditingId(c.id)
    setEditName(c.name)
    setEditDesc(c.description)
    setEditSpecialist(c.specialist_id ?? '')
  }

  const saveEdit = async () => {
    if (!editingId) return
    setSaving(true)
    try {
      await api.adminUpdateCategory(editingId, {
        name: editName,
        description: editDesc,
        specialist_id: editSpecialist !== '' ? Number(editSpecialist) : null,
      })
      setCategories((prev) => prev.map((c) =>
        c.id === editingId
          ? { ...c, name: editName, description: editDesc, specialist_id: editSpecialist !== '' ? Number(editSpecialist) : undefined }
          : c
      ))
      setEditingId(null)
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить категорию? Вопросы останутся, но потеряют привязку к категории.')) return
    try {
      await api.adminDeleteCategory(id)
      setCategories((prev) => prev.filter((c) => c.id !== id))
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    try {
      const created = await api.adminCreateCategory({
        name: newName,
        description: newDesc,
        specialist_id: newSpecialist !== '' ? Number(newSpecialist) : null,
      })
      setCategories((prev) => [...prev, created])
      setNewName(''); setNewDesc(''); setNewSpecialist('')
      setShowCreate(false)
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const specialistName = (id?: number) =>
    id ? (specialists.find((u) => u.id === id)?.username ?? `#${id}`) : '—'

  if (loading) return <div className="text-sm text-muted-foreground">Загрузка...</div>
  if (error) return <ErrorBox message={error} />

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <p className="text-sm text-muted-foreground">Всего категорий: {categories.length}</p>
        <Button size="sm" onClick={() => setShowCreate((v) => !v)}>
          <Plus className="h-4 w-4" /> Создать категорию
        </Button>
      </div>

      {showCreate && (
        <form onSubmit={handleCreate} className="rounded-lg border border-border p-4 space-y-3 bg-accent/30">
          <p className="font-medium text-sm">Новая категория</p>
          <input
            required placeholder="Название" value={newName}
            onChange={(e) => setNewName(e.target.value)}
            className="rounded-md border border-border bg-background px-3 py-1.5 text-sm w-full"
          />
          <input
            placeholder="Описание (необязательно)" value={newDesc}
            onChange={(e) => setNewDesc(e.target.value)}
            className="rounded-md border border-border bg-background px-3 py-1.5 text-sm w-full"
          />
          <select
            value={newSpecialist}
            onChange={(e) => setNewSpecialist(e.target.value === '' ? '' : Number(e.target.value))}
            className="rounded-md border border-border bg-background px-3 py-1.5 text-sm w-full"
          >
            <option value="">— Специалист не назначен —</option>
            {specialists.map((u) => (
              <option key={u.id} value={u.id}>{u.username} ({ROLE_LABELS[u.role]})</option>
            ))}
          </select>
          <div className="flex gap-2">
            <Button type="submit" size="sm" disabled={saving}>Создать</Button>
            <Button type="button" size="sm" variant="outline" onClick={() => setShowCreate(false)}>Отмена</Button>
          </div>
        </form>
      )}

      <div className="rounded-lg border border-border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-accent/50">
            <tr>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">ID</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Название</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Описание</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Специалист</th>
              <th className="px-3 py-2"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {categories.map((c) => (
              <tr key={c.id}>
                <td className="px-3 py-2 text-muted-foreground">{c.id}</td>
                {editingId === c.id ? (
                  <>
                    <td className="px-3 py-2">
                      <input
                        value={editName} onChange={(e) => setEditName(e.target.value)}
                        className="rounded border border-border bg-background px-2 py-1 text-sm w-full"
                      />
                    </td>
                    <td className="px-3 py-2">
                      <input
                        value={editDesc} onChange={(e) => setEditDesc(e.target.value)}
                        className="rounded border border-border bg-background px-2 py-1 text-sm w-full"
                      />
                    </td>
                    <td className="px-3 py-2">
                      <select
                        value={editSpecialist}
                        onChange={(e) => setEditSpecialist(e.target.value === '' ? '' : Number(e.target.value))}
                        className="rounded border border-border bg-background px-2 py-1 text-sm w-full"
                      >
                        <option value="">— не назначен —</option>
                        {specialists.map((u) => (
                          <option key={u.id} value={u.id}>{u.username}</option>
                        ))}
                      </select>
                    </td>
                    <td className="px-3 py-2">
                      <div className="flex gap-1 justify-end">
                        <button onClick={saveEdit} disabled={saving} className="p-1 rounded hover:bg-accent text-green-600">
                          <Check className="h-4 w-4" />
                        </button>
                        <button onClick={() => setEditingId(null)} className="p-1 rounded hover:bg-accent text-muted-foreground">
                          <X className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </>
                ) : (
                  <>
                    <td className="px-3 py-2 font-medium">{c.name}</td>
                    <td className="px-3 py-2 text-muted-foreground max-w-xs truncate">{c.description || '—'}</td>
                    <td className="px-3 py-2 text-muted-foreground">{specialistName(c.specialist_id)}</td>
                    <td className="px-3 py-2">
                      <div className="flex gap-1 justify-end">
                        <button onClick={() => startEdit(c)} className="p-1 rounded hover:bg-accent text-muted-foreground">
                          <Pencil className="h-3.5 w-3.5" />
                        </button>
                        <button onClick={() => handleDelete(c.id)} className="p-1 rounded hover:bg-accent text-destructive">
                          <Trash2 className="h-3.5 w-3.5" />
                        </button>
                      </div>
                    </td>
                  </>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

// ─── Tags Tab ────────────────────────────────────────────────────────────────

function TagsTab() {
  const [tags, setTags] = useState<TagType[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editName, setEditName] = useState('')
  const [newName, setNewName] = useState('')
  const [saving, setSaving] = useState(false)

  const load = () => {
    setLoading(true)
    api.getTags()
      .then(setTags)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }

  useEffect(load, [])

  const startEdit = (t: TagType) => {
    setEditingId(t.id)
    setEditName(t.name)
  }

  const saveEdit = async () => {
    if (!editingId) return
    setSaving(true)
    try {
      await api.adminUpdateTag(editingId, editName)
      setTags((prev) => prev.map((t) => t.id === editingId ? { ...t, name: editName } : t))
      setEditingId(null)
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Удалить тег? Он будет убран со всех вопросов.')) return
    try {
      await api.adminDeleteTag(id)
      setTags((prev) => prev.filter((t) => t.id !== id))
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newName.trim()) return
    setSaving(true)
    try {
      const created = await api.adminCreateTag(newName.trim())
      setTags((prev) => [...prev, created])
      setNewName('')
    } catch (e) {
      alert(e instanceof Error ? e.message : 'Ошибка')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <div className="text-sm text-muted-foreground">Загрузка...</div>
  if (error) return <ErrorBox message={error} />

  return (
    <div className="space-y-4">
      <form onSubmit={handleCreate} className="flex gap-2">
        <input
          placeholder="Название нового тега"
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          className="flex-1 rounded-md border border-border bg-background px-3 py-1.5 text-sm"
        />
        <Button type="submit" size="sm" disabled={saving || !newName.trim()}>
          <Plus className="h-4 w-4" /> Добавить
        </Button>
      </form>

      <div className="rounded-lg border border-border overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-accent/50">
            <tr>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">ID</th>
              <th className="text-left px-3 py-2 font-medium text-muted-foreground">Название</th>
              <th className="px-3 py-2"></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border">
            {tags.map((t) => (
              <tr key={t.id}>
                <td className="px-3 py-2 text-muted-foreground">{t.id}</td>
                {editingId === t.id ? (
                  <>
                    <td className="px-3 py-2">
                      <input
                        value={editName}
                        onChange={(e) => setEditName(e.target.value)}
                        onKeyDown={(e) => { if (e.key === 'Enter') saveEdit(); if (e.key === 'Escape') setEditingId(null) }}
                        autoFocus
                        className="rounded border border-border bg-background px-2 py-1 text-sm w-full"
                      />
                    </td>
                    <td className="px-3 py-2">
                      <div className="flex gap-1 justify-end">
                        <button onClick={saveEdit} disabled={saving} className="p-1 rounded hover:bg-accent text-green-600">
                          <Check className="h-4 w-4" />
                        </button>
                        <button onClick={() => setEditingId(null)} className="p-1 rounded hover:bg-accent text-muted-foreground">
                          <X className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </>
                ) : (
                  <>
                    <td className="px-3 py-2">
                      <span className="inline-flex items-center rounded-full bg-accent border border-border px-2.5 py-0.5 text-xs font-medium">
                        #{t.name}
                      </span>
                    </td>
                    <td className="px-3 py-2">
                      <div className="flex gap-1 justify-end">
                        <button onClick={() => startEdit(t)} className="p-1 rounded hover:bg-accent text-muted-foreground">
                          <Pencil className="h-3.5 w-3.5" />
                        </button>
                        <button onClick={() => handleDelete(t.id)} className="p-1 rounded hover:bg-accent text-destructive">
                          <Trash2 className="h-3.5 w-3.5" />
                        </button>
                      </div>
                    </td>
                  </>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

// ─── Data Tab (import/export) ────────────────────────────────────────────────

function DataTab() {
  const fileRef = useRef<HTMLInputElement>(null)
  const [importing, setImporting] = useState(false)
  const [importResult, setImportResult] = useState<{ count: number } | null>(null)
  const [error, setError] = useState<string | null>(null)

  // Экспорт за период
  const today = new Date()
  const firstOfMonth = new Date(today.getFullYear(), today.getMonth(), 1)
  const fmt = (d: Date) => d.toISOString().slice(0, 10)

  const [periodFrom, setPeriodFrom] = useState(fmt(firstOfMonth))
  const [periodTo, setPeriodTo] = useState(fmt(today))
  const [periodFormat, setPeriodFormat] = useState<'csv' | 'json'>('csv')
  const [exportingPeriod, setExportingPeriod] = useState(false)
  const [periodError, setPeriodError] = useState<string | null>(null)

  const handleImport = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setError(null)
    setImportResult(null)
    setImporting(true)
    try {
      const text = await file.text()
      const json = JSON.parse(text)
      const entries = Array.isArray(json) ? json : json.training_data
      if (!Array.isArray(entries)) throw new Error('Ожидается массив или объект с полем training_data')
      const result = await api.importData(entries)
      setImportResult({ count: result.imported })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка импорта')
    } finally {
      setImporting(false)
      if (fileRef.current) fileRef.current.value = ''
    }
  }

  const handleExportPeriod = async () => {
    setPeriodError(null)
    setExportingPeriod(true)
    try {
      const res = await api.exportPeriod(periodFrom, periodTo, periodFormat)
      if (!res.ok) {
        const body = await res.json().catch(() => ({ error: res.statusText }))
        throw new Error(body.error || res.statusText)
      }
      const blob = await res.blob()
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      const from = periodFrom || 'начало'
      const to = periodTo || 'конец'
      a.download = `questions_${from}_${to}.${periodFormat}`
      a.click()
      URL.revokeObjectURL(url)
    } catch (err) {
      setPeriodError(err instanceof Error ? err.message : 'Ошибка экспорта')
    } finally {
      setExportingPeriod(false)
    }
  }

return (
    <div className="grid gap-6 max-w-2xl">
      <div className="rounded-lg border border-border p-6">
        <h2 className="text-lg font-semibold mb-1">Импорт вопросов</h2>
        <p className="text-sm text-muted-foreground mb-4">
          Загрузите JSON-файл в формате training_data. Дубликаты по заголовку пропускаются.
        </p>
        <input
          ref={fileRef}
          type="file"
          accept=".json,application/json"
          className="hidden"
          onChange={handleImport}
        />
        <Button onClick={() => fileRef.current?.click()} disabled={importing} className="gap-2">
          <Upload className="h-4 w-4" />
          {importing ? 'Импортируется...' : 'Выбрать файл и импортировать'}
        </Button>
        {importResult && (
          <div className="mt-4 flex items-center gap-2 text-sm text-green-700 bg-green-50 border border-green-200 rounded-md px-3 py-2">
            <CheckCircle className="h-4 w-4 shrink-0" />
            Импортировано записей: <strong>{importResult.count}</strong>
          </div>
        )}
      </div>

      <div className="rounded-lg border border-border p-6">
        <h2 className="text-lg font-semibold mb-1">Экспорт вопросов за период</h2>
        <p className="text-sm text-muted-foreground mb-4">
          Выгрузите вопросы за выбранный диапазон дат. CSV — удобно для Excel, JSON — тот же формат что и полный экспорт.
        </p>
        <div className="flex flex-wrap items-end gap-3">
          <div className="flex flex-col gap-1">
            <label className="text-xs font-medium text-muted-foreground">С</label>
            <input
              type="date"
              value={periodFrom}
              onChange={(e) => setPeriodFrom(e.target.value)}
              className="rounded-md border border-border bg-background px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            />
          </div>
          <div className="flex flex-col gap-1">
            <label className="text-xs font-medium text-muted-foreground">По</label>
            <input
              type="date"
              value={periodTo}
              onChange={(e) => setPeriodTo(e.target.value)}
              className="rounded-md border border-border bg-background px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            />
          </div>
          <div className="flex flex-col gap-1">
            <label className="text-xs font-medium text-muted-foreground">Формат</label>
            <div className="flex rounded-md border border-border overflow-hidden text-sm h-[34px]">
              {(['csv', 'json'] as const).map((f) => (
                <button
                  key={f}
                  type="button"
                  onClick={() => setPeriodFormat(f)}
                  className={`px-3 font-medium transition-colors ${
                    periodFormat === f
                      ? 'bg-primary text-primary-foreground'
                      : 'bg-background text-muted-foreground hover:bg-muted'
                  }`}
                >
                  {f.toUpperCase()}
                </button>
              ))}
            </div>
          </div>
          <Button
            variant="outline"
            onClick={handleExportPeriod}
            disabled={exportingPeriod}
            className="gap-2"
          >
            <Download className="h-4 w-4" />
            {exportingPeriod ? 'Формируется...' : 'Скачать'}
          </Button>
        </div>
        {periodError && <div className="mt-3"><ErrorBox message={periodError} /></div>}
      </div>

      {error && <ErrorBox message={error} />}
    </div>
  )
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

function ErrorBox({ message }: { message: string }) {
  return (
    <div className="flex items-center gap-2 text-sm text-red-700 bg-red-50 border border-red-200 rounded-md px-3 py-2">
      <AlertCircle className="h-4 w-4 shrink-0" />
      {message}
    </div>
  )
}
