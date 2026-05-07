import { useEffect, useRef, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { ArrowLeft, CheckCircle2, AlertTriangle, ExternalLink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { AttachmentsPanel } from '@/components/attachments/AttachmentsPanel'
import { MarkdownEditor } from '@/components/ui/MarkdownEditor'
import { createQuestion, getCategories, getTags, search } from '@/api'
import { useAuth } from '@/context/AuthContext'
import type { Category, Question, SearchResult, Tag } from '@/types'

const STATUS_LABEL: Record<string, string> = {
  OPEN: 'Открыт',
  CLOSED: 'Закрыт',
  DUPLICATE: 'Дубликат',
}

export function CreateQuestionPage() {
  const navigate = useNavigate()
  const { user } = useAuth()
  const [title, setTitle] = useState('')
  const [body, setBody] = useState('')
  const [categoryId, setCategoryId] = useState<number | undefined>()
  const [selectedTags, setSelectedTags] = useState<number[]>([])
  const [categories, setCategories] = useState<Category[]>([])
  const [tags, setTags] = useState<Tag[]>([])
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')
  const [created, setCreated] = useState<Question | null>(null)

  const [similar, setSimilar] = useState<SearchResult[]>([])
  const [similarDismissed, setSimilarDismissed] = useState(false)
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    if (!user) navigate('/login')
  }, [user])

  useEffect(() => {
    getCategories().then(setCategories).catch(() => {})
    getTags().then(setTags).catch(() => {})
  }, [])

  // Поиск похожих вопросов при вводе заголовка
  useEffect(() => {
    setSimilarDismissed(false)

    if (title.trim().length < 8) {
      setSimilar([])
      return
    }

    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(async () => {
      try {
        const res = await search({ q: title.trim() })
        const questions = (res.items ?? []).filter((r) => r.type === 'question').slice(0, 5)
        setSimilar(questions)
      } catch {
        setSimilar([])
      }
    }, 600)

    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current)
    }
  }, [title])

  const toggleTag = (id: number) => {
    setSelectedTags((prev) =>
      prev.includes(id) ? prev.filter((t) => t !== id) : [...prev, id]
    )
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!title.trim() || !body.trim()) {
      setError('Заголовок и текст вопроса обязательны')
      return
    }
    setError('')
    setSubmitting(true)
    try {
      const q = await createQuestion({
        title: title.trim(),
        body: body.trim(),
        category_id: categoryId,
        tag_ids: selectedTags.length > 0 ? selectedTags : undefined,
      })
      setCreated(q)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания вопроса')
    } finally {
      setSubmitting(false)
    }
  }

  if (created) {
    return (
      <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8 py-6">
        <div className="mb-6 flex items-center gap-2 text-green-700">
          <CheckCircle2 className="h-5 w-5" />
          <span className="font-semibold">Вопрос опубликован!</span>
        </div>

        <div className="rounded-lg border border-border bg-background p-5 mb-6">
          <p className="font-medium text-foreground mb-1">{created.title}</p>
          <p className="text-sm text-muted-foreground line-clamp-2">{created.body}</p>
        </div>

        <div className="mb-6">
          <p className="text-sm font-medium text-foreground mb-3">Прикрепить файлы (необязательно)</p>
          <AttachmentsPanel
            targetType="question"
            targetId={created.id}
            uploaderIdAllowed={created.author_id}
          />
        </div>

        <Button onClick={() => navigate(`/questions/${created.id}`)}>
          Перейти к вопросу
        </Button>
      </div>
    )
  }

  const showSimilar = similar.length > 0 && !similarDismissed

  return (
    <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8 py-6">
      <button
        onClick={() => navigate(-1)}
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground mb-5 transition-colors"
      >
        <ArrowLeft className="h-4 w-4" /> Назад
      </button>

      <h1 className="text-xl font-bold text-foreground mb-6">Задать вопрос</h1>

      <form onSubmit={handleSubmit} className="space-y-5">
        {/* Title */}
        <div>
          <label className="text-sm font-medium text-foreground block mb-1.5">
            Заголовок <span className="text-red-500">*</span>
          </label>
          <input
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Кратко опишите проблему или вопрос"
            className="w-full h-10 px-3 rounded-md border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-ring transition-colors"
            maxLength={200}
          />

          {/* Похожие вопросы */}
          {showSimilar && (
            <div className="mt-2 rounded-md border border-amber-200 bg-amber-50 text-sm overflow-hidden">
              <div className="flex items-start justify-between gap-2 px-3 py-2 bg-amber-100/70">
                <div className="flex items-center gap-1.5 text-amber-800 font-medium">
                  <AlertTriangle className="h-4 w-4 shrink-0" />
                  Похожие вопросы уже существуют — возможно, ответ уже есть:
                </div>
                <button
                  type="button"
                  onClick={() => setSimilarDismissed(true)}
                  className="text-amber-600 hover:text-amber-800 text-xs shrink-0 mt-0.5"
                >
                  Скрыть
                </button>
              </div>
              <ul className="divide-y divide-amber-100">
                {similar.map((q) => (
                  <li key={q.question_id}>
                    <Link
                      to={`/questions/${q.question_id}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center justify-between gap-2 px-3 py-2 hover:bg-amber-100/50 transition-colors group"
                    >
                      <span className="text-amber-900 group-hover:text-amber-700 line-clamp-1 flex-1">
                        {q.question_title}
                      </span>
                      <div className="flex items-center gap-2 shrink-0">
                        <span className={`text-xs px-1.5 py-0.5 rounded font-medium ${
                          q.status === 'CLOSED'
                            ? 'bg-green-100 text-green-700'
                            : 'bg-blue-100 text-blue-700'
                        }`}>
                          {STATUS_LABEL[q.status] ?? q.status}
                        </span>
                        <ExternalLink className="h-3.5 w-3.5 text-amber-500 group-hover:text-amber-700" />
                      </div>
                    </Link>
                  </li>
                ))}
              </ul>
              <div className="px-3 py-2 bg-amber-50 border-t border-amber-100 text-xs text-amber-700">
                Если ни один вопрос не подходит — продолжайте заполнять форму ниже.
              </div>
            </div>
          )}
        </div>

        {/* Body */}
        <div>
          <label className="text-sm font-medium text-foreground block mb-1.5">
            Описание <span className="text-red-500">*</span>
          </label>
          <MarkdownEditor
            rows={8}
            value={body}
            onChange={setBody}
            placeholder="Подробно опишите вопрос, что уже пробовали, какие ошибки возникают..."
          />
        </div>

        {/* Category */}
        {categories.length > 0 && (
          <div>
            <label className="text-sm font-medium text-foreground block mb-1.5">
              Категория
            </label>
            <select
              value={categoryId ?? ''}
              onChange={(e) => setCategoryId(e.target.value ? parseInt(e.target.value) : undefined)}
              className="w-full h-10 px-3 rounded-md border border-border bg-background text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            >
              <option value="">Без категории</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>{c.name}</option>
              ))}
            </select>
          </div>
        )}

        {/* Tags */}
        {tags.length > 0 && (
          <div>
            <label className="text-sm font-medium text-foreground block mb-1.5">
              Теги
            </label>
            <div className="flex flex-wrap gap-2">
              {tags.map((tag) => (
                <button
                  key={tag.id}
                  type="button"
                  onClick={() => toggleTag(tag.id)}
                  className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-medium border transition-colors ${
                    selectedTags.includes(tag.id)
                      ? 'bg-primary text-primary-foreground border-primary'
                      : 'bg-background text-muted-foreground border-border hover:border-primary hover:text-primary'
                  }`}
                >
                  {tag.name}
                </button>
              ))}
            </div>
          </div>
        )}

        {error && <p className="text-sm text-red-600">{error}</p>}

        <div className="flex items-center gap-3 pt-2">
          <Button type="submit" disabled={submitting || !title.trim() || !body.trim()}>
            {submitting ? 'Публикация...' : 'Опубликовать вопрос'}
          </Button>
          <Button type="button" variant="outline" onClick={() => navigate(-1)}>
            Отмена
          </Button>
        </div>
      </form>
    </div>
  )
}
