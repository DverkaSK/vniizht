import { useRef, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Upload, Download, CheckCircle, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/context/AuthContext'
import * as api from '@/api'

export function AdminPage() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const fileRef = useRef<HTMLInputElement>(null)

  const [importing, setImporting] = useState(false)
  const [exporting, setExporting] = useState(false)
  const [importResult, setImportResult] = useState<{ count: number } | null>(null)
  const [error, setError] = useState<string | null>(null)

  if (!user || user.role !== 'ADMIN') {
    navigate('/')
    return null
  }

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

  const handleExport = async () => {
    setError(null)
    setExporting(true)
    try {
      const res = await api.exportData()
      if (!res.ok) {
        const body = await res.json().catch(() => ({ error: res.statusText }))
        throw new Error(body.error || res.statusText)
      }
      const blob = await res.blob()
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = 'training_data.json'
      a.click()
      URL.revokeObjectURL(url)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка экспорта')
    } finally {
      setExporting(false)
    }
  }

  return (
    <div className="mx-auto max-w-3xl px-4 py-10">
      <h1 className="text-2xl font-bold mb-8">Панель администратора</h1>

      <div className="grid gap-6">
        {/* Import */}
        <div className="rounded-lg border border-border p-6">
          <h2 className="text-lg font-semibold mb-1">Импорт вопросов</h2>
          <p className="text-sm text-muted-foreground mb-4">
            Загрузите JSON-файл в формате training_data (массив или объект с полем training_data).
            Дубликаты по заголовку пропускаются автоматически.
          </p>
          <input
            ref={fileRef}
            type="file"
            accept=".json,application/json"
            className="hidden"
            onChange={handleImport}
          />
          <Button
            onClick={() => fileRef.current?.click()}
            disabled={importing}
            className="gap-2"
          >
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

        {/* Export */}
        <div className="rounded-lg border border-border p-6">
          <h2 className="text-lg font-semibold mb-1">Экспорт вопросов</h2>
          <p className="text-sm text-muted-foreground mb-4">
            Скачайте все вопросы с верифицированными ответами в формате training_data.json.
          </p>
          <Button
            variant="outline"
            onClick={handleExport}
            disabled={exporting}
            className="gap-2"
          >
            <Download className="h-4 w-4" />
            {exporting ? 'Формируется...' : 'Скачать training_data.json'}
          </Button>
        </div>

        {error && (
          <div className="flex items-center gap-2 text-sm text-red-700 bg-red-50 border border-red-200 rounded-md px-3 py-2">
            <AlertCircle className="h-4 w-4 shrink-0" />
            {error}
          </div>
        )}
      </div>
    </div>
  )
}
