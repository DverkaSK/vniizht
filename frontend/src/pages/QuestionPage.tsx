import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, CheckCircle2, ChevronDown, ChevronUp, MessageSquare } from 'lucide-react'

import {
  createAnswer,
  createAnswerVote,
  createComment,
  deleteAnswerVote,
  getAnswers,
  getComments,
  getQuestion,
  unverifyAnswer,
  updateAnswerVote,
  verifyAnswer,
} from '@/api'
import { AttachmentsPanel } from '@/components/attachments/AttachmentsPanel'
import { Avatar } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/context/AuthContext'
import { formatRelative } from '@/lib/utils'
import type { Answer, Comment, Question, User, VoteValue } from '@/types'

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
                    </div>
                    <p className="text-sm text-foreground">{comment.body}</p>
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
}: {
  answer: AnswerWithComments
  user: User | null
  canVerify: boolean
  votePending: boolean
  verifyPending: boolean
  onVote: (answer: AnswerWithComments, value: VoteValue) => Promise<void>
  onToggleVerify: (answer: AnswerWithComments) => Promise<void>
}) {
  const voteBlocked = !user || user.id === answer.author_id || votePending

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
          <div className="prose prose-sm mb-4 max-w-none text-foreground">
            {answer.body.split('\n\n').map((paragraph, index) => {
              if (paragraph.startsWith('```')) {
                const code = paragraph.replace(/```[\w]*\n?/, '').replace(/```$/, '')
                return (
                  <pre key={index} className="my-3 overflow-x-auto rounded-md bg-gray-900 p-3 text-sm text-green-400">
                    <code>{code}</code>
                  </pre>
                )
              }

              return (
                <p key={index} className="mb-2 whitespace-pre-wrap text-sm leading-relaxed last:mb-0">
                  {paragraph}
                </p>
              )
            })}
          </div>

          <div className="mb-3 flex flex-wrap items-center justify-between gap-2">
            <div className="flex items-center gap-2">
              <Avatar name={answer.author_username} size="sm" />
              <Link to={`/users/${answer.author_id}`} className="text-sm font-medium hover:text-primary">
                {answer.author_username}
              </Link>
            </div>
            <span className="text-xs text-muted-foreground">{formatRelative(answer.created_at)}</span>
          </div>

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

  const canVerifyAnswers = user?.role === 'SPECIALIST' || user?.role === 'ADMIN'

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

      <div className="mb-6 rounded-lg border border-border bg-background p-6">
        <div className="mb-3 flex flex-wrap items-center gap-2">
          <Badge variant={statusVariant[question.status]}>{statusLabel[question.status]}</Badge>
          {question.has_verified && <Badge variant="verified">✓ Есть достоверный ответ</Badge>}
        </div>

        <h1 className="mb-3 text-xl font-bold leading-tight text-foreground">
          {question.title}
        </h1>

        <p className="mb-4 whitespace-pre-wrap text-sm leading-relaxed text-foreground">
          {question.body}
        </p>

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

        <div className="mt-4 flex flex-wrap items-center justify-between gap-3 border-t border-border pt-4">
          <div className="flex items-center gap-2">
            <Avatar name={question.author_username} size="sm" />
            <div>
              <Link to={`/users/${question.author_id}`} className="text-sm font-medium hover:text-primary">
                {question.author_username}
              </Link>
              <p className="text-xs text-muted-foreground">{formatRelative(question.created_at)}</p>
            </div>
          </div>
          <span className="text-xs text-muted-foreground">{question.view_count} просмотров</span>
        </div>
      </div>

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
              <textarea
                rows={5}
                placeholder="Поделитесь знаниями по данному вопросу..."
                value={answerText}
                onChange={(e) => setAnswerText(e.target.value)}
                className="w-full resize-y rounded-md border border-border bg-background px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
              />
              <div className="flex justify-end">
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
