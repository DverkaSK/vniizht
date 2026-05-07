import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, CheckCircle2, ChevronDown, ChevronUp, Clock, MessageSquare, Pencil, Trash2, X } from 'lucide-react'

import {
  adminListUsers,
  assignSpecialist,
  closeQuestion,
  createAnswer,
  createAnswerVote,
  createComment,
  deleteAnswer,
  deleteComment,
  deleteAnswerVote,
  deleteQuestion,
  getAnswerHistory,
  getAnswers,
  getComments,
  getQuestion,
  getQuestionHistory,
  getTags,
  markQuestionDuplicate,
  search,
  unverifyAnswer,
  updateAnswer,
  updateAnswerVote,
  updateComment,
  updateQuestion,
  verifyAnswer,
} from '@/api'
import { AttachmentsPanel } from '@/components/attachments/AttachmentsPanel'
import { Avatar } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { MarkdownContent } from '@/components/ui/MarkdownContent'
import { MarkdownEditor } from '@/components/ui/MarkdownEditor'
import { useAuth } from '@/context/AuthContext'
import { formatDateTime, formatRelative, isEdited } from '@/lib/utils'
import type { Answer, AnswerHistoryEntry, Comment, Question, QuestionHistoryEntry, SearchResult, Tag, User, VoteValue } from '@/types'

function EditedBadge({ createdAt, updatedAt }: { createdAt: string; updatedAt: string }) {
  if (!isEdited(createdAt, updatedAt)) return null
  return (
    <span className="relative group inline-flex items-center">
      <Pencil className="h-3 w-3 text-muted-foreground cursor-default" />
      <span className="pointer-events-none absolute bottom-full left-1/2 -translate-x-1/2 mb-1.5 whitespace-nowrap rounded bg-popover border border-border px-2 py-1 text-xs text-foreground shadow-md opacity-0 group-hover:opacity-100 transition-opacity z-10">
        Изменено {formatDateTime(updatedAt)}
      </span>
    </span>
  )
}

type HistoryEntry = (QuestionHistoryEntry | AnswerHistoryEntry) & { title?: string }

function HistoryModal({ entries, onClose }: { entries: HistoryEntry[]; onClose: () => void }) {
  const [expanded, setExpanded] = useState<number | null>(null)

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" onClick={onClose}>
      <div
        className="w-full max-w-2xl max-h-[80vh] overflow-hidden rounded-lg border border-border bg-background shadow-lg flex flex-col"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 className="text-base font-semibold">История редактирования</h2>
          <button onClick={onClose} className="rounded-md p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
            <X className="h-4 w-4" />
          </button>
        </div>

        {entries.length === 0 ? (
          <div className="flex flex-col items-center justify-center gap-2 py-12 text-muted-foreground">
            <Clock className="h-8 w-8 opacity-40" />
            <p className="text-sm">История изменений пуста</p>
          </div>
        ) : (
          <div className="overflow-y-auto">
            {entries.map((entry, idx) => (
              <div key={entry.id} className="border-b border-border last:border-0">
                <button
                  className="flex w-full items-center gap-3 px-5 py-3.5 text-left hover:bg-muted/50 transition-colors"
                  onClick={() => setExpanded(expanded === idx ? null : idx)}
                >
                  <Avatar name={entry.editor_username} size="sm" />
                  <div className="flex-1 min-w-0">
                    <span className="text-sm font-medium">{entry.editor_username}</span>
                    {entry.title && (
                      <span className="ml-2 text-sm text-muted-foreground truncate">— {entry.title}</span>
                    )}
                  </div>
                  <span className="shrink-0 text-xs text-muted-foreground">{formatDateTime(entry.edited_at)}</span>
                  <ChevronDown className={`h-4 w-4 shrink-0 text-muted-foreground transition-transform ${expanded === idx ? 'rotate-180' : ''}`} />
                </button>

                {expanded === idx && (
                  <div className="bg-muted/30 px-5 pb-4 pt-1">
                    {entry.title && (
                      <p className="mb-2 text-sm font-semibold text-foreground">{entry.title}</p>
                    )}
                    <div className="rounded-md border border-border bg-background p-3 text-sm text-foreground overflow-y-auto max-h-60">
                      <MarkdownContent>{entry.body}</MarkdownContent>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}

function DuplicateModal({
  currentId,
  onConfirm,
  onClose,
}: {
  currentId: number
  onConfirm: (original: SearchResult) => Promise<void>
  onClose: () => void
}) {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<SearchResult[]>([])
  const [selected, setSelected] = useState<SearchResult | null>(null)
  const [searching, setSearching] = useState(false)
  const [confirming, setConfirming] = useState(false)

  useEffect(() => {
    if (query.trim().length < 3) { setResults([]); return }
    const timer = setTimeout(async () => {
      setSearching(true)
      try {
        const res = await search({ q: query })
        setResults(
          res.items
            .filter((r) => r.type === 'question' && r.question_id !== currentId)
            .slice(0, 6)
        )
      } catch {
        setResults([])
      } finally {
        setSearching(false)
      }
    }, 400)
    return () => clearTimeout(timer)
  }, [query, currentId])

  const handleConfirm = async () => {
    if (!selected) return
    setConfirming(true)
    try {
      await onConfirm(selected)
    } finally {
      setConfirming(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" onClick={onClose}>
      <div
        className="w-full max-w-lg rounded-lg border border-border bg-background shadow-lg"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 className="text-base font-semibold">Пометить как дубликат</h2>
          <button onClick={onClose} className="rounded-md p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="p-4 space-y-3">
          <p className="text-sm text-muted-foreground">Найдите оригинальный вопрос, дубликатом которого является текущий:</p>
          <input
            autoFocus
            type="text"
            placeholder="Поиск по названию..."
            value={query}
            onChange={(e) => { setQuery(e.target.value); setSelected(null) }}
            className="w-full rounded-md border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
          />

          {searching && <p className="text-xs text-muted-foreground">Поиск...</p>}

          {results.length > 0 && (
            <div className="rounded-md border border-border divide-y divide-border max-h-60 overflow-y-auto">
              {results.map((r) => (
                <button
                  key={r.question_id}
                  onClick={() => setSelected(r)}
                  className={`w-full px-3 py-2.5 text-left text-sm transition-colors ${
                    selected?.question_id === r.question_id
                      ? 'bg-primary/10 text-primary'
                      : 'hover:bg-muted'
                  }`}
                >
                  <span className="font-medium line-clamp-1">{r.question_title}</span>
                  <span className="text-xs text-muted-foreground ml-1">#{r.question_id}</span>
                </button>
              ))}
            </div>
          )}

          {query.trim().length >= 3 && !searching && results.length === 0 && (
            <p className="text-sm text-muted-foreground text-center py-2">Ничего не найдено</p>
          )}

          {selected && (
            <div className="rounded-md bg-amber-50 border border-amber-200 px-3 py-2 text-sm">
              <span className="text-amber-800">Выбрано: </span>
              <span className="font-medium text-amber-900">{selected.question_title}</span>
            </div>
          )}
        </div>

        <div className="flex justify-end gap-2 border-t border-border px-5 py-3">
          <Button variant="outline" size="sm" onClick={onClose}>Отмена</Button>
          <Button
            size="sm"
            disabled={!selected || confirming}
            onClick={handleConfirm}
          >
            {confirming ? 'Сохранение...' : 'Пометить как дубликат'}
          </Button>
        </div>
      </div>
    </div>
  )
}

function AssignModal({
  currentSpecialistId,
  onConfirm,
  onClose,
}: {
  currentSpecialistId?: number
  onConfirm: (specialistId: number | null) => Promise<void>
  onClose: () => void
}) {
  const [users, setUsers] = useState<User[]>([])
  const [selected, setSelected] = useState<number | null>(currentSpecialistId ?? null)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    adminListUsers()
      .then((all) => setUsers(all.filter((u) => u.role === 'SPECIALIST' || u.role === 'ADMIN')))
      .catch(() => {})
  }, [])

  const handleSave = async () => {
    setSaving(true)
    try { await onConfirm(selected) } finally { setSaving(false) }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" onClick={onClose}>
      <div
        className="w-full max-w-sm rounded-lg border border-border bg-background shadow-lg"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 className="text-base font-semibold">Назначить специалиста</h2>
          <button onClick={onClose} className="rounded-md p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors">
            <X className="h-4 w-4" />
          </button>
        </div>

        <div className="p-4 space-y-1 max-h-72 overflow-y-auto">
          <button
            onClick={() => setSelected(null)}
            className={`w-full rounded-md px-3 py-2 text-left text-sm transition-colors ${selected === null ? 'bg-primary/10 text-primary font-medium' : 'hover:bg-muted text-muted-foreground'}`}
          >
            — Без специалиста
          </button>
          {users.map((u) => (
            <button
              key={u.id}
              onClick={() => setSelected(u.id)}
              className={`w-full rounded-md px-3 py-2 text-left text-sm transition-colors ${selected === u.id ? 'bg-primary/10 text-primary font-medium' : 'hover:bg-muted text-foreground'}`}
            >
              {u.username}
              <span className="ml-2 text-xs text-muted-foreground">{u.role}</span>
            </button>
          ))}
        </div>

        <div className="flex justify-end gap-2 border-t border-border px-5 py-3">
          <Button variant="outline" size="sm" onClick={onClose}>Отмена</Button>
          <Button size="sm" disabled={saving} onClick={handleSave}>
            {saving ? 'Сохранение...' : 'Назначить'}
          </Button>
        </div>
      </div>
    </div>
  )
}

const statusLabel: Record<string, string> = {
  OPEN: 'Открыт',
  CLOSED: 'Закрыт',
  DUPLICATE: 'Дубликат',
}

const statusVariant: Record<string, 'status-open' | 'status-progress' | 'status-closed'> = {
  OPEN: 'status-open',
  CLOSED: 'status-closed',
  DUPLICATE: 'status-closed',
}

type AnswerWithComments = Answer & { comments: Comment[] }

async function loadAnswersWithComments(questionId: number): Promise<AnswerWithComments[]> {
  const answers = await getAnswers(questionId)
  return Promise.all(
    answers.map(async (answer) => {
      const comments = await getComments(answer.id).catch(() => [] as Comment[])
      return { ...answer, comments }
    })
  )
}

function CommentThread({ answerId, initialComments }: { answerId: number; initialComments: Comment[] }) {
  const { user } = useAuth()
  const [comments, setComments] = useState<Comment[]>(initialComments)
  const [show, setShow] = useState(true)
  const [text, setText] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editText, setEditText] = useState('')
  const [editSaving, setEditSaving] = useState(false)
  const [deleteConfirmId, setDeleteConfirmId] = useState<number | null>(null)
  const [deletingId, setDeletingId] = useState<number | null>(null)

  useEffect(() => {
    setComments(initialComments)
  }, [initialComments])

  const handleSubmit = async () => {
    if (!text.trim()) return
    setSubmitting(true)
    try {
      const comment = await createComment(answerId, text.trim())
      setComments((prev) => [...prev, comment])
      setText('')
    } catch {
      // ignore
    } finally {
      setSubmitting(false)
    }
  }

  const startEdit = (comment: Comment) => {
    setEditingId(comment.id)
    setEditText(comment.body)
  }

  const handleEditSave = async (commentId: number) => {
    if (!editText.trim()) return
    setEditSaving(true)
    try {
      const updated = await updateComment(commentId, editText.trim())
      setComments((prev) => prev.map((c) => c.id === commentId ? updated : c))
      setEditingId(null)
    } catch {
      // ignore
    } finally {
      setEditSaving(false)
    }
  }

  const handleDelete = async (commentId: number) => {
    setDeletingId(commentId)
    try {
      await deleteComment(commentId)
      setComments((prev) => prev.filter((c) => c.id !== commentId))
      setDeleteConfirmId(null)
    } catch {
      // ignore
    } finally {
      setDeletingId(null)
    }
  }

  return (
    <div className="mt-3 border-t border-border pt-3">
      {comments.length > 0 && (
        <>
          <button
            onClick={() => setShow(!show)}
            className="mb-2 flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <MessageSquare className="h-3.5 w-3.5" />
            {comments.length} {comments.length === 1 ? 'комментарий' : 'комментариев'}
            {show ? <ChevronUp className="h-3.5 w-3.5" /> : <ChevronDown className="h-3.5 w-3.5" />}
          </button>

          {show && (
            <div className="mb-3 ml-2 space-y-2">
              {comments.map((comment) => (
                <div key={comment.id} className="flex gap-2 text-sm">
                  <Avatar name={comment.author_username} size="sm" />
                  <div className="flex-1 rounded-lg bg-muted/50 px-3 py-2">
                    <div className="mb-0.5 flex items-center gap-2">
                      <Link to={`/users/${comment.author_id}`} className="text-xs font-medium hover:text-primary">
                        {comment.author_username}
                      </Link>
                      <span className="text-xs text-muted-foreground">{formatRelative(comment.created_at)}</span>
                      <EditedBadge createdAt={comment.created_at} updatedAt={comment.updated_at} />
                      {user && user.id === comment.author_id && editingId !== comment.id && (
                        <div className="ml-auto flex items-center gap-1">
                          <button
                            onClick={() => startEdit(comment)}
                            className="text-muted-foreground hover:text-foreground transition-colors"
                            title="Редактировать"
                          >
                            <Pencil className="h-3 w-3" />
                          </button>
                          {deleteConfirmId === comment.id ? (
                            <>
                              <button
                                onClick={() => handleDelete(comment.id)}
                                disabled={deletingId === comment.id}
                                className="text-xs text-red-600 hover:underline"
                              >
                                {deletingId === comment.id ? '...' : 'Да'}
                              </button>
                              <button
                                onClick={() => setDeleteConfirmId(null)}
                                className="text-xs text-muted-foreground hover:underline"
                              >
                                Нет
                              </button>
                            </>
                          ) : (
                            <button
                              onClick={() => setDeleteConfirmId(comment.id)}
                              className="text-muted-foreground hover:text-red-600 transition-colors"
                              title="Удалить"
                            >
                              <Trash2 className="h-3 w-3" />
                            </button>
                          )}
                        </div>
                      )}
                    </div>
                    {editingId === comment.id ? (
                      <div className="flex gap-2 mt-1">
                        <input
                          autoFocus
                          type="text"
                          value={editText}
                          onChange={(e) => setEditText(e.target.value)}
                          onKeyDown={(e) => {
                            if (e.key === 'Enter') handleEditSave(comment.id)
                            if (e.key === 'Escape') setEditingId(null)
                          }}
                          className="flex-1 rounded-md border border-border bg-background px-2 py-1 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                        />
                        <Button size="sm" disabled={editSaving || !editText.trim()} onClick={() => handleEditSave(comment.id)}>
                          {editSaving ? '...' : 'Сохранить'}
                        </Button>
                        <Button size="sm" variant="outline" onClick={() => setEditingId(null)}>Отмена</Button>
                      </div>
                    ) : (
                      <p className="text-sm text-foreground">{comment.body}</p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      )}

      {user && (
        <div className="mt-2 flex gap-2">
          <Avatar name={user.username} size="sm" />
          <div className="flex flex-1 gap-2">
            <input
              type="text"
              placeholder="Добавить комментарий..."
              value={text}
              onChange={(e) => setText(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSubmit()}
              className="flex-1 rounded-md border border-border bg-background px-3 py-1.5 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
            />
            {text.trim() && (
              <Button size="sm" onClick={handleSubmit} disabled={submitting}>
                Отправить
              </Button>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

function AnswerBlock({
  answer,
  user,
  canVerify,
  votePending,
  verifyPending,
  onVote,
  onToggleVerify,
  onDelete,
  onUpdated,
}: {
  answer: AnswerWithComments
  user: User | null
  canVerify: boolean
  votePending: boolean
  verifyPending: boolean
  onVote: (answer: AnswerWithComments, value: VoteValue) => Promise<void>
  onToggleVerify: (answer: AnswerWithComments) => Promise<void>
  onDelete: (id: number) => void
  onUpdated: () => void
}) {
  const voteBlocked = !user || user.id === answer.author_id || votePending
  const [historyEntries, setHistoryEntries] = useState<AnswerHistoryEntry[] | null>(null)
  const [historyLoading, setHistoryLoading] = useState(false)
  const [editing, setEditing] = useState(false)
  const [editBody, setEditBody] = useState('')
  const [editSaving, setEditSaving] = useState(false)
  const [deleteConfirm, setDeleteConfirm] = useState(false)
  const [deleting, setDeleting] = useState(false)

  const canEdit = !!user && (user.id === answer.author_id || user.role === 'ADMIN' || user.role === 'SPECIALIST')
  const canDelete = !!user && (user.id === answer.author_id || user.role === 'ADMIN')

  const openHistory = async () => {
    if (historyEntries !== null) { setHistoryEntries(null); return }
    setHistoryLoading(true)
    try {
      const h = await getAnswerHistory(answer.id)
      setHistoryEntries(h)
    } catch {
      setHistoryEntries([])
    } finally {
      setHistoryLoading(false)
    }
  }

  const startEdit = () => {
    setEditBody(answer.body)
    setEditing(true)
  }

  const handleEditSave = async () => {
    if (!editBody.trim()) return
    setEditSaving(true)
    try {
      await updateAnswer(answer.id, editBody.trim())
      setEditing(false)
      onUpdated()
    } catch {
      // ignore
    } finally {
      setEditSaving(false)
    }
  }

  const handleDelete = async () => {
    setDeleting(true)
    try {
      await deleteAnswer(answer.id)
      onDelete(answer.id)
    } catch {
      setDeleting(false)
      setDeleteConfirm(false)
    }
  }

  return (
    <div className={`rounded-lg border p-5 ${answer.is_verified ? 'border-green-300 bg-green-50/30' : 'border-border bg-background'}`}>
      {answer.is_verified && (
        <div className="mb-3 flex items-center gap-1.5 text-sm font-medium text-green-700">
          <CheckCircle2 className="h-4 w-4" />
          Достоверный ответ, подтвержден специалистом
        </div>
      )}

      <div className="flex gap-4">
        <div className="shrink-0 pt-1">
          <div className="flex flex-col items-center gap-1">
            <button
              type="button"
              aria-label="Голос вверх"
              disabled={voteBlocked}
              onClick={() => onVote(answer, 'UP')}
              className={`rounded-md p-1.5 transition-colors ${
                answer.current_user_vote === 'UP'
                  ? 'bg-green-50 text-green-600'
                  : voteBlocked
                    ? 'text-muted-foreground/50'
                    : 'text-muted-foreground hover:bg-green-50 hover:text-green-600'
              }`}
            >
              <ChevronUp className="h-5 w-5" />
            </button>

            <span className={`text-lg font-semibold tabular-nums ${answer.vote_score > 0 ? 'text-green-600' : answer.vote_score < 0 ? 'text-red-500' : 'text-muted-foreground'}`}>
              {answer.vote_score}
            </span>

            <button
              type="button"
              aria-label="Голос вниз"
              disabled={voteBlocked}
              onClick={() => onVote(answer, 'DOWN')}
              className={`rounded-md p-1.5 transition-colors ${
                answer.current_user_vote === 'DOWN'
                  ? 'bg-red-50 text-red-500'
                  : voteBlocked
                    ? 'text-muted-foreground/50'
                    : 'text-muted-foreground hover:bg-red-50 hover:text-red-500'
              }`}
            >
              <ChevronDown className="h-5 w-5" />
            </button>

            <span className="text-[10px] text-muted-foreground">голосов</span>
          </div>
        </div>

        <div className="min-w-0 flex-1">
          {editing ? (
            <div className="mb-4 space-y-2">
              <MarkdownEditor rows={6} value={editBody} onChange={setEditBody} placeholder="Текст ответа" />
              <div className="flex gap-2">
                <Button size="sm" disabled={editSaving || !editBody.trim()} onClick={handleEditSave}>
                  {editSaving ? 'Сохранение...' : 'Сохранить'}
                </Button>
                <Button size="sm" variant="outline" onClick={() => setEditing(false)}>Отмена</Button>
              </div>
            </div>
          ) : (
            <MarkdownContent className="mb-4">{answer.body}</MarkdownContent>
          )}

          <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
            <div className="flex items-center gap-2">
              <Avatar name={answer.author_username} size="sm" />
              <Link to={`/users/${answer.author_id}`} className="text-sm font-medium hover:text-primary">
                {answer.author_username}
              </Link>
            </div>
            <div className="flex items-center gap-1.5">
              <span className="text-xs text-muted-foreground">{formatRelative(answer.created_at)}</span>
              <EditedBadge createdAt={answer.created_at} updatedAt={answer.updated_at} />
              {isEdited(answer.created_at, answer.updated_at) && (
                <button
                  onClick={openHistory}
                  disabled={historyLoading}
                  className="flex items-center gap-1 rounded px-1.5 py-0.5 text-xs text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
                  title="История изменений"
                >
                  <Clock className="h-3 w-3" />
                  {historyLoading ? '...' : 'История'}
                </button>
              )}
              {canEdit && !editing && (
                <button
                  onClick={startEdit}
                  className="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
                  title="Редактировать ответ"
                >
                  <Pencil className="h-3.5 w-3.5" />
                </button>
              )}
              {canDelete && !editing && (
                deleteConfirm ? (
                  <div className="flex items-center gap-1.5">
                    <span className="text-xs text-red-600">Удалить?</span>
                    <Button size="sm" variant="destructive" disabled={deleting} onClick={handleDelete}>
                      {deleting ? '...' : 'Да'}
                    </Button>
                    <Button size="sm" variant="outline" onClick={() => setDeleteConfirm(false)}>Нет</Button>
                  </div>
                ) : (
                  <button
                    onClick={() => setDeleteConfirm(true)}
                    className="rounded p-1 text-muted-foreground hover:bg-red-50 hover:text-red-600 transition-colors"
                    title="Удалить ответ"
                  >
                    <Trash2 className="h-3.5 w-3.5" />
                  </button>
                )
              )}
            </div>
          </div>
          {historyEntries !== null && (
            <HistoryModal
              entries={historyEntries as HistoryEntry[]}
              onClose={() => setHistoryEntries(null)}
            />
          )}

          {canVerify && (
            <div className="mb-3">
              <Button
                size="sm"
                variant={answer.is_verified ? 'secondary' : 'default'}
                disabled={verifyPending}
                onClick={() => onToggleVerify(answer)}
              >
                {answer.is_verified ? 'Снять подтверждение' : 'Подтвердить ответ'}
              </Button>
            </div>
          )}

          <AttachmentsPanel
            targetType="answer"
            targetId={answer.id}
            uploaderIdAllowed={answer.author_id}
          />

          <CommentThread answerId={answer.id} initialComments={answer.comments} />
        </div>
      </div>
    </div>
  )
}

export function QuestionPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const { user } = useAuth()

  const [question, setQuestion] = useState<Question | null>(null)
  const [answers, setAnswers] = useState<AnswerWithComments[]>([])
  const [loading, setLoading] = useState(true)
  const [answerText, setAnswerText] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [votePendingId, setVotePendingId] = useState<number | null>(null)
  const [verifyPendingId, setVerifyPendingId] = useState<number | null>(null)

  const [editing, setEditing] = useState(false)
  const [editTitle, setEditTitle] = useState('')
  const [editBody, setEditBody] = useState('')
  const [editSaving, setEditSaving] = useState(false)
  const [allTags, setAllTags] = useState<Tag[]>([])
  const [editTagIds, setEditTagIds] = useState<number[]>([])
  const [deleteConfirm, setDeleteConfirm] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [questionHistory, setQuestionHistory] = useState<QuestionHistoryEntry[] | null>(null)
  const [questionHistoryLoading, setQuestionHistoryLoading] = useState(false)
  const [showDuplicateModal, setShowDuplicateModal] = useState(false)
  const [showAssignModal, setShowAssignModal] = useState(false)

  const openQuestionHistory = async () => {
    if (questionHistory !== null) { setQuestionHistory(null); return }
    if (!id) return
    setQuestionHistoryLoading(true)
    try {
      const h = await getQuestionHistory(parseInt(id, 10))
      setQuestionHistory(h)
    } catch {
      setQuestionHistory([])
    } finally {
      setQuestionHistoryLoading(false)
    }
  }

  const canVerifyAnswers = user?.role === 'SPECIALIST' || user?.role === 'ADMIN'
  const canCloseQuestion = question !== null && question.status === 'OPEN' && (
    user?.id === question.author_id ||
    user?.role === 'SPECIALIST' ||
    user?.role === 'ADMIN'
  )
  const canEditQuestion = question !== null && question.status === 'OPEN' && (
    user?.id === question.author_id || user?.role === 'ADMIN'
  )
  const canDeleteQuestion = question !== null && (
    user?.id === question.author_id || user?.role === 'ADMIN' || user?.role === 'SPECIALIST'
  )
  const canMarkDuplicate = question !== null && question.status === 'OPEN' && (
    user?.id === question.author_id ||
    user?.role === 'SPECIALIST' ||
    user?.role === 'ADMIN'
  )

  const loadPage = async (questionId: number) => {
    const [loadedQuestion, loadedAnswers] = await Promise.all([
      getQuestion(questionId),
      loadAnswersWithComments(questionId),
    ])
    setQuestion(loadedQuestion)
    setAnswers(loadedAnswers)
  }

  useEffect(() => {
    if (!id) return

    const questionId = parseInt(id, 10)
    setLoading(true)
    loadPage(questionId)
      .catch(() => navigate('/'))
      .finally(() => setLoading(false))
  }, [id, navigate, user?.id])

  useEffect(() => {
    getTags().then(setAllTags).catch(() => {})
  }, [])

  const startEdit = () => {
    if (!question) return
    setEditTitle(question.title)
    setEditBody(question.body)
    setEditTagIds(question.tags.map((t) => t.id))
    setEditing(true)
  }

  const handleEditSave = async () => {
    if (!id || !editTitle.trim() || !editBody.trim()) return
    setEditSaving(true)
    try {
      await updateQuestion(parseInt(id, 10), { title: editTitle.trim(), body: editBody.trim(), tag_ids: editTagIds })
      await loadPage(parseInt(id, 10))
      setEditing(false)
    } catch {
      // ignore
    } finally {
      setEditSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!id) return
    setDeleting(true)
    try {
      await deleteQuestion(parseInt(id, 10))
      navigate('/')
    } catch {
      setDeleting(false)
      setDeleteConfirm(false)
    }
  }

  const handleAnswerSubmit = async () => {
    if (!id || !answerText.trim()) return

    setSubmitting(true)
    try {
      const answer = await createAnswer(parseInt(id, 10), answerText.trim())
      setAnswers((prev) => [...prev, { ...answer, comments: [] }])
      setAnswerText('')
      setQuestion((prev) => prev ? { ...prev, answer_count: prev.answer_count + 1 } : prev)
    } catch {
      // ignore
    } finally {
      setSubmitting(false)
    }
  }

  const handleVote = async (answer: AnswerWithComments, value: VoteValue) => {
    if (!id || !user || user.id === answer.author_id) return

    setVotePendingId(answer.id)
    try {
      if (answer.current_user_vote === value) {
        await deleteAnswerVote(answer.id)
      } else if (answer.current_user_vote) {
        await updateAnswerVote(answer.id, value)
      } else {
        await createAnswerVote(answer.id, value)
      }
      await loadPage(parseInt(id, 10))
    } catch {
      // ignore
    } finally {
      setVotePendingId(null)
    }
  }

  const handleCloseQuestion = async () => {
    if (!id) return
    try {
      await closeQuestion(parseInt(id, 10))
      setQuestion((prev) => prev ? { ...prev, status: 'CLOSED' } : prev)
    } catch {
      // ignore
    }
  }

  const handleAssignSpecialist = async (specialistId: number | null) => {
    if (!id) return
    await assignSpecialist(parseInt(id, 10), specialistId)
    await loadPage(parseInt(id, 10))
    setShowAssignModal(false)
  }

  const handleMarkDuplicate = async (original: SearchResult) => {
    if (!id) return
    await markQuestionDuplicate(parseInt(id, 10), original.question_id)
    await loadPage(parseInt(id, 10))
    setShowDuplicateModal(false)
  }

  const handleToggleVerify = async (answer: AnswerWithComments) => {
    if (!id || !canVerifyAnswers) return

    setVerifyPendingId(answer.id)
    try {
      if (answer.is_verified) {
        await unverifyAnswer(answer.id)
      } else {
        await verifyAnswer(answer.id)
      }
      await loadPage(parseInt(id, 10))
    } catch {
      // ignore
    } finally {
      setVerifyPendingId(null)
    }
  }

  if (loading) {
    return (
      <div className="mx-auto max-w-4xl space-y-4 px-4 py-6 sm:px-6 lg:px-8">
        <div className="h-8 w-32 animate-pulse rounded bg-muted" />
        <div className="h-48 animate-pulse rounded-lg border border-border bg-background p-6" />
      </div>
    )
  }

  if (!question) return null

  return (
    <div className="mx-auto max-w-4xl px-4 py-6 sm:px-6 lg:px-8">
      <Link to="/" className="mb-5 inline-flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground">
        <ArrowLeft className="h-4 w-4" /> Все вопросы
      </Link>

      {question.duplicate_of && (
        <div className="mb-4 flex items-center gap-2 rounded-lg border border-amber-200 bg-amber-50 px-4 py-3 text-sm text-amber-800">
          <span>Этот вопрос помечен как дубликат.</span>
          <Link
            to={`/questions/${question.duplicate_of}`}
            className="font-medium underline hover:text-amber-900"
          >
            Перейти к оригинальному вопросу →
          </Link>
        </div>
      )}

      <div className="mb-6 rounded-lg border border-border bg-background p-6">
        <div className="mb-3 flex flex-wrap items-center gap-2">
          <Badge variant={statusVariant[question.status]}>{statusLabel[question.status]}</Badge>
          {question.has_verified && <Badge variant="verified">✓ Есть достоверный ответ</Badge>}
          {question.specialist_username && (
            <span className="inline-flex items-center gap-1 rounded-full border border-blue-200 bg-blue-50 px-2.5 py-0.5 text-xs font-medium text-blue-700">
              Специалист: {question.specialist_username}
            </span>
          )}
        </div>

        {editing ? (
          <div className="space-y-3">
            <input
              className="w-full rounded-md border border-border bg-background px-3 py-2 text-base font-bold focus:outline-none focus:ring-2 focus:ring-ring"
              value={editTitle}
              onChange={(e) => setEditTitle(e.target.value)}
              placeholder="Заголовок вопроса"
            />
            <MarkdownEditor
              rows={6}
              value={editBody}
              onChange={setEditBody}
              placeholder="Текст вопроса"
            />
            {allTags.length > 0 && (
              <div>
                <p className="mb-1.5 text-xs font-medium text-muted-foreground">Теги</p>
                <div className="flex flex-wrap gap-1.5">
                  {allTags.map((tag) => {
                    const selected = editTagIds.includes(tag.id)
                    return (
                      <button
                        key={tag.id}
                        type="button"
                        onClick={() => setEditTagIds(selected ? editTagIds.filter((t) => t !== tag.id) : [...editTagIds, tag.id])}
                        className={`rounded-full border px-2.5 py-0.5 text-xs font-medium transition-colors ${selected ? 'border-primary bg-primary text-primary-foreground' : 'border-border bg-accent text-muted-foreground hover:border-primary hover:text-primary'}`}
                      >
                        {tag.name}
                      </button>
                    )
                  })}
                </div>
              </div>
            )}
            <div className="flex gap-2">
              <Button size="sm" disabled={editSaving || !editTitle.trim() || !editBody.trim()} onClick={handleEditSave}>
                {editSaving ? 'Сохранение...' : 'Сохранить'}
              </Button>
              <Button size="sm" variant="outline" onClick={() => setEditing(false)}>Отмена</Button>
            </div>
          </div>
        ) : (
          <>
            <h1 className="mb-3 text-xl font-bold leading-tight text-foreground">
              {question.title}
            </h1>

            <MarkdownContent className="mb-4">{question.body}</MarkdownContent>

            {question.tags.length > 0 && (
              <div className="mb-4 flex flex-wrap gap-1.5">
                {question.tags.map((tag) => (
                  <Link
                    key={tag.id}
                    to={`/?tag=${tag.id}`}
                    className="inline-flex items-center rounded-full border border-border bg-accent px-2.5 py-0.5 text-xs font-medium text-muted-foreground transition-colors hover:border-primary hover:text-primary"
                  >
                    {tag.name}
                  </Link>
                ))}
              </div>
            )}

            <AttachmentsPanel
              targetType="question"
              targetId={question.id}
              uploaderIdAllowed={question.author_id}
            />
          </>
        )}

        <div className="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-border pt-4">
          <div className="flex items-center gap-2">
            <Avatar name={question.author_username} size="sm" />
            <div>
              <Link to={`/users/${question.author_id}`} className="text-sm font-medium hover:text-primary">
                {question.author_username}
              </Link>
              <div className="flex items-center gap-1.5">
                <p className="text-xs text-muted-foreground">{formatRelative(question.created_at)}</p>
                <EditedBadge createdAt={question.created_at} updatedAt={question.updated_at} />
                {isEdited(question.created_at, question.updated_at) && (
                  <button
                    onClick={openQuestionHistory}
                    disabled={questionHistoryLoading}
                    className="flex items-center gap-1 rounded px-1.5 py-0.5 text-xs text-muted-foreground hover:bg-muted hover:text-foreground transition-colors"
                    title="История изменений"
                  >
                    <Clock className="h-3 w-3" />
                    {questionHistoryLoading ? '...' : 'История'}
                  </button>
                )}
              </div>
            </div>
          </div>
          <div className="flex items-center gap-2 flex-wrap">
            <span className="text-xs text-muted-foreground">{question.view_count} просмотров</span>
            {canEditQuestion && !editing && (
              <Button size="sm" variant="outline" onClick={startEdit} className="gap-1.5">
                <Pencil className="h-3.5 w-3.5" /> Редактировать
              </Button>
            )}
            {canDeleteQuestion && !editing && (
              deleteConfirm ? (
                <div className="flex items-center gap-2">
                  <span className="text-xs text-red-600">Удалить вопрос?</span>
                  <Button size="sm" variant="destructive" disabled={deleting} onClick={handleDelete}>
                    {deleting ? 'Удаление...' : 'Да, удалить'}
                  </Button>
                  <Button size="sm" variant="outline" onClick={() => setDeleteConfirm(false)}>Отмена</Button>
                </div>
              ) : (
                <Button size="sm" variant="outline" onClick={() => setDeleteConfirm(true)} className="gap-1.5 text-red-600 hover:border-red-300 hover:bg-red-50">
                  <Trash2 className="h-3.5 w-3.5" /> Удалить
                </Button>
              )
            )}
            {user?.role === 'ADMIN' && !editing && (
              <Button size="sm" variant="outline" onClick={() => setShowAssignModal(true)}>
                Назначить специалиста
              </Button>
            )}
            {canMarkDuplicate && !editing && (
              <Button size="sm" variant="outline" onClick={() => setShowDuplicateModal(true)}>
                Дубликат
              </Button>
            )}
            {canCloseQuestion && !editing && (
              <Button size="sm" variant="secondary" onClick={handleCloseQuestion}>
                Закрыть вопрос
              </Button>
            )}
          </div>
        </div>
      </div>

      {questionHistory !== null && (
        <HistoryModal
          entries={questionHistory as HistoryEntry[]}
          onClose={() => setQuestionHistory(null)}
        />
      )}
      {showDuplicateModal && question && (
        <DuplicateModal
          currentId={question.id}
          onConfirm={handleMarkDuplicate}
          onClose={() => setShowDuplicateModal(false)}
        />
      )}
      {showAssignModal && question && (
        <AssignModal
          currentSpecialistId={question.specialist_id}
          onConfirm={handleAssignSpecialist}
          onClose={() => setShowAssignModal(false)}
        />
      )}

      <div className="mb-6">
        <h2 className="mb-3 text-base font-semibold">
          {question.answer_count} {question.answer_count === 1 ? 'ответ' : 'ответов'}
        </h2>
        <div className="space-y-4">
          {answers.map((answer) => (
            <AnswerBlock
              key={answer.id}
              answer={answer}
              user={user}
              canVerify={canVerifyAnswers}
              votePending={votePendingId === answer.id}
              verifyPending={verifyPendingId === answer.id}
              onVote={handleVote}
              onToggleVerify={handleToggleVerify}
              onDelete={(deletedId) => {
                setAnswers((prev) => prev.filter((a) => a.id !== deletedId))
                setQuestion((prev) => prev ? { ...prev, answer_count: prev.answer_count - 1 } : prev)
              }}
              onUpdated={() => id && loadPage(parseInt(id, 10))}
            />
          ))}
        </div>
      </div>

      {user && question.status === 'OPEN' && (
        <div className="rounded-lg border border-border bg-background p-5">
          <h3 className="mb-3 text-base font-semibold">Ваш ответ</h3>
          <div className="flex gap-3">
            <Avatar name={user.username} size="md" />
            <div className="flex-1 space-y-3">
              <MarkdownEditor
                rows={6}
                placeholder="Поделитесь знаниями по данному вопросу..."
                value={answerText}
                onChange={setAnswerText}
              />
              <div className="flex items-center justify-between gap-3">
                <span className="text-xs text-muted-foreground">Прикрепить файлы можно после публикации ответа</span>
                <Button size="sm" disabled={!answerText.trim() || submitting} onClick={handleAnswerSubmit}>
                  {submitting ? 'Публикация...' : 'Опубликовать ответ'}
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
