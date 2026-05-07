import { useEffect, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { Plus, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { QuestionCard } from '@/components/questions/QuestionCard'
import { Sidebar } from '@/components/layout/Sidebar'
import { getQuestions, getCategories, getTags } from '@/api'
import type { Question, Category, Tag } from '@/types'
import { useAuth } from '@/context/AuthContext'

type StatusFilter = '' | 'OPEN' | 'CLOSED'

export function HomePage() {
  const { user } = useAuth()
  const [searchParams, setSearchParams] = useSearchParams()
  const [questions, setQuestions] = useState<Question[]>([])
  const [total, setTotal] = useState(0)
  const [categories, setCategories] = useState<Category[]>([])
  const [tags, setTags] = useState<Tag[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('')
  const [myQuestionsOnly, setMyQuestionsOnly] = useState(false)

  const isSpecialist = user?.role === 'SPECIALIST' || user?.role === 'ADMIN'

  const categoryParam = searchParams.get('category')
  const tagParam = searchParams.get('tag')
  const categoryId = categoryParam ? parseInt(categoryParam) : undefined
  const tagId = tagParam ? parseInt(tagParam) : undefined

  const activeCategory = categories.find((c) => c.id === categoryId)
  const activeTag = tags.find((t) => t.id === tagId)

  useEffect(() => {
    getCategories().then(setCategories).catch(() => {})
    getTags().then(setTags).catch(() => {})
  }, [])

  useEffect(() => {
    setPage(1)
  }, [categoryId, tagId, statusFilter, myQuestionsOnly])

  useEffect(() => {
    setLoading(true)
    getQuestions({
      page,
      status: statusFilter || undefined,
      category_id: categoryId,
      tag_id: tagId,
      assigned_specialist_id: myQuestionsOnly && user ? user.id : undefined,
    })
      .then((res) => {
        setQuestions(res.items ?? [])
        setTotal(res.total)
      })
      .catch(() => setQuestions([]))
      .finally(() => setLoading(false))
  }, [page, categoryId, tagId, statusFilter, myQuestionsOnly, user?.id])

  const categoryMap = Object.fromEntries(categories.map((c) => [c.id, c.name]))

  const limit = 20
  const totalPages = Math.ceil(total / limit)

  return (
    <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-6">
      <div className="flex gap-8">
        {/* Sidebar */}
        <div className="hidden lg:block">
          <Sidebar />
        </div>

        {/* Main content */}
        <div className="flex-1 min-w-0">
          {/* Header */}
          <div className="flex items-center justify-between mb-4 gap-3 flex-wrap">
            <div className="flex items-center gap-2 flex-wrap">
              <h1 className="text-xl font-bold text-foreground">База знаний</h1>
              {activeCategory && (
                <span className="flex items-center gap-1 rounded-full bg-primary/10 text-primary px-2.5 py-0.5 text-sm font-medium">
                  {activeCategory.name}
                  <button onClick={() => setSearchParams({})} className="hover:opacity-70">
                    <X className="h-3.5 w-3.5" />
                  </button>
                </span>
              )}
              {activeTag && (
                <span className="flex items-center gap-1 rounded-full bg-accent border border-border px-2.5 py-0.5 text-sm font-medium text-foreground">
                  #{activeTag.name}
                  <button onClick={() => setSearchParams({})} className="hover:opacity-70">
                    <X className="h-3.5 w-3.5" />
                  </button>
                </span>
              )}
            </div>
            {user && (
              <Link to="/questions/new">
                <Button size="sm">
                  <Plus className="h-4 w-4" />
                  Задать вопрос
                </Button>
              </Link>
            )}
          </div>

          {/* Filters */}
          <div className="flex items-center gap-3 mb-4 flex-wrap">
            <div className="flex rounded-lg border border-border overflow-hidden text-sm">
              {([['', 'Все'], ['OPEN', 'Открытые'], ['CLOSED', 'Закрытые']] as [StatusFilter, string][]).map(([val, label]) => (
                <button
                  key={val}
                  onClick={() => setStatusFilter(val)}
                  className={`px-3 py-1.5 transition-colors ${
                    statusFilter === val
                      ? 'bg-primary text-primary-foreground'
                      : 'bg-background hover:bg-accent text-foreground'
                  }`}
                >
                  {label}
                </button>
              ))}
            </div>
            {isSpecialist && (
              <button
                onClick={() => setMyQuestionsOnly((v) => !v)}
                className={`rounded-lg border px-3 py-1.5 text-sm transition-colors ${
                  myQuestionsOnly
                    ? 'border-primary bg-primary text-primary-foreground'
                    : 'border-border bg-background text-foreground hover:bg-accent'
                }`}
              >
                Мои вопросы
              </button>
            )}
          </div>

          {/* Count */}
          <p className="text-sm text-muted-foreground mb-3">
            Найдено: <span className="font-medium text-foreground">{total}</span> вопросов
          </p>

          {/* Questions list */}
          {loading ? (
            <div className="space-y-3">
              {Array.from({ length: 5 }).map((_, i) => (
                <div key={i} className="rounded-lg border border-border bg-background p-4 animate-pulse h-24" />
              ))}
            </div>
          ) : questions.length === 0 ? (
            <div className="text-center py-16 text-muted-foreground">
              <p className="text-lg font-medium mb-1">Вопросов не найдено</p>
              <p className="text-sm">Попробуйте изменить фильтры или задайте вопрос первым</p>
            </div>
          ) : (
            <div className="space-y-3">
              {questions.map((q) => (
                <QuestionCard key={q.id} question={q} categoryName={q.category_id ? categoryMap[q.category_id] : undefined} />
              ))}
            </div>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-center gap-2 mt-6">
              <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(p => p - 1)}>
                Назад
              </Button>
              <span className="text-sm text-muted-foreground">{page} / {totalPages}</span>
              <Button variant="outline" size="sm" disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}>
                Вперёд
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
