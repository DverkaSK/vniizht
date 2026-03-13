import type {
  User, PublicUser, Question, QuestionListResponse,
  Answer, Comment, Attachment, SearchResponse, Category, Tag, VoteValue
} from './types'

const BASE = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(BASE + path, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || res.statusText)
  }
  if (res.status === 204) return undefined as unknown as T
  return res.json()
}

// ── Аутентификация ────────────────────────────────────────────

export const login = (username: string, password: string) =>
  request<void>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  })

export const logout = () =>
  request<void>('/auth/logout', { method: 'POST' })

// ── Профиль ───────────────────────────────────────────────────

export const getMe = () => request<User>('/users/me')

export const getUser = (id: number) => request<PublicUser>(`/users/${id}`)

// ── Вопросы ───────────────────────────────────────────────────

export const getQuestions = (params?: {
  page?: number
  status?: string
  category_id?: number
}) => {
  const qs = new URLSearchParams()
  if (params?.page) qs.set('page', String(params.page))
  if (params?.status) qs.set('status', params.status)
  if (params?.category_id) qs.set('category_id', String(params.category_id))
  return request<QuestionListResponse>(`/questions?${qs}`)
}

export const getQuestion = (id: number) => request<Question>(`/questions/${id}`)

export const createQuestion = (data: {
  title: string
  body: string
  category_id?: number
  tag_ids?: number[]
}) => request<Question>('/questions', { method: 'POST', body: JSON.stringify(data) })

export const updateQuestion = (id: number, data: {
  title: string
  body: string
  tag_ids?: number[]
}) => request<void>(`/questions/${id}`, { method: 'PATCH', body: JSON.stringify(data) })

export const closeQuestion = (id: number) =>
  request<void>(`/questions/${id}/close`, { method: 'PATCH', body: '{}' })

// ── Ответы ────────────────────────────────────────────────────

export const getAnswers = (questionId: number) =>
  request<Answer[]>(`/questions/${questionId}/answers`)

export const createAnswer = (questionId: number, body: string) =>
  request<Answer>(`/questions/${questionId}/answers`, {
    method: 'POST',
    body: JSON.stringify({ body }),
  })

export const updateAnswer = (id: number, body: string) =>
  request<void>(`/answers/${id}`, { method: 'PATCH', body: JSON.stringify({ body }) })

export const deleteAnswer = (id: number) =>
  request<void>(`/answers/${id}`, { method: 'DELETE' })

export const createAnswerVote = (id: number, value: VoteValue) =>
  request<void>(`/answers/${id}/vote`, {
    method: 'POST',
    body: JSON.stringify({ value }),
  })

export const updateAnswerVote = (id: number, value: VoteValue) =>
  request<void>(`/answers/${id}/vote`, {
    method: 'PATCH',
    body: JSON.stringify({ value }),
  })

export const deleteAnswerVote = (id: number) =>
  request<void>(`/answers/${id}/vote`, { method: 'DELETE' })

export const verifyAnswer = (id: number) =>
  request<void>(`/answers/${id}/verify`, { method: 'PATCH', body: '{}' })

export const unverifyAnswer = (id: number) =>
  request<void>(`/answers/${id}/verify`, { method: 'DELETE' })

// ── Комментарии ───────────────────────────────────────────────

export const getComments = (answerId: number) =>
  request<Comment[]>(`/answers/${answerId}/comments`)

export const createComment = (answerId: number, body: string) =>
  request<Comment>(`/answers/${answerId}/comments`, {
    method: 'POST',
    body: JSON.stringify({ body }),
  })

export const updateComment = (id: number, body: string) =>
  request<void>(`/comments/${id}`, { method: 'PATCH', body: JSON.stringify({ body }) })

export const deleteComment = (id: number) =>
  request<void>(`/comments/${id}`, { method: 'DELETE' })

// ── Вложения ──────────────────────────────────────────────────

export const getAttachments = (targetType: string, targetId: number) =>
  request<Attachment[]>(`/attachments?target_type=${targetType.toUpperCase()}&target_id=${targetId}`)

export const uploadAttachment = (file: File, targetType: string, targetId: number) => {
  const form = new FormData()
  form.append('file', file)
  form.append('target_type', targetType.toUpperCase())
  form.append('target_id', String(targetId))
  return fetch(BASE + '/attachments', {
    method: 'POST',
    credentials: 'include',
    body: form,
  }).then(async res => {
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: res.statusText }))
      throw new Error(err.error || res.statusText)
    }
    return res.json() as Promise<Attachment>
  })
}

export const deleteAttachment = (id: number) =>
  request<void>(`/attachments/${id}`, { method: 'DELETE' })

export const attachmentUrl = (id: number) => `${BASE}/attachments/${id}`

// ── Поиск ─────────────────────────────────────────────────────

export const search = (params: {
  q: string
  category_id?: number
  tag_id?: number
  status?: string
}) => {
  const qs = new URLSearchParams({ q: params.q })
  if (params.category_id) qs.set('category_id', String(params.category_id))
  if (params.tag_id) qs.set('tag_id', String(params.tag_id))
  if (params.status) qs.set('status', params.status)
  return request<SearchResponse>(`/search?${qs}`)
}

// ── Категории и теги ──────────────────────────────────────────

export const getCategories = () => request<Category[]>('/categories')

export const getTags = () => request<Tag[]>('/tags')

// ── Вспомогательные ───────────────────────────────────────────

export const formatDate = (iso: string) =>
  new Date(iso).toLocaleDateString('ru-RU', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
