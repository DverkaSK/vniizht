import { useEffect, useRef, useState } from 'react'
import { Paperclip, Download, Trash2, Loader2, FileText, Image, Archive } from 'lucide-react'
import { getAttachments, uploadAttachment, deleteAttachment, attachmentUrl } from '@/api'
import { useAuth } from '@/context/AuthContext'
import type { Attachment } from '@/types'
import { ImageLightbox } from './ImageLightbox'

interface Props {
  targetType: 'question' | 'answer' | 'comment'
  targetId: number
  uploaderIdAllowed?: number // ID владельца объекта (может прикреплять файлы)
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) return `${bytes} Б`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} КБ`
  return `${(bytes / (1024 * 1024)).toFixed(1)} МБ`
}

function fileIcon(mimeType: string) {
  if (mimeType.startsWith('image/')) return <Image className="h-4 w-4 text-blue-500" />
  if (mimeType.includes('zip') || mimeType.includes('archive') || mimeType.includes('tar') || mimeType.includes('gz'))
    return <Archive className="h-4 w-4 text-amber-500" />
  return <FileText className="h-4 w-4 text-muted-foreground" />
}

export function AttachmentsPanel({ targetType, targetId, uploaderIdAllowed }: Props) {
  const { user } = useAuth()
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [attachments, setAttachments] = useState<Attachment[]>([])
  const [uploading, setUploading] = useState(false)
  const [error, setError] = useState('')
  const [lightbox, setLightbox] = useState<Attachment | null>(null)

  useEffect(() => {
    if (!targetId) return
    getAttachments(targetType, targetId)
      .then(setAttachments)
      .catch(() => {})
  }, [targetType, targetId])

  const canUpload = user && (
    user.id === uploaderIdAllowed ||
    user.role === 'admin' ||
    user.role === 'specialist'
  )

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    e.target.value = ''

    if (file.size > 50 * 1024 * 1024) {
      setError('Файл превышает 50 МБ')
      return
    }

    setError('')
    setUploading(true)
    try {
      const att = await uploadAttachment(file, targetType, targetId)
      setAttachments((prev) => [...prev, att])
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка загрузки')
    } finally {
      setUploading(false)
    }
  }

  const handleDelete = async (id: number) => {
    try {
      await deleteAttachment(id)
      setAttachments((prev) => prev.filter((a) => a.id !== id))
    } catch {
      setError('Ошибка удаления')
    }
  }

  const canDelete = (att: Attachment) =>
    user && (att.uploader_id === user.id || user.role === 'admin' || user.role === 'specialist')

  const isImage = (mime: string) => mime.startsWith('image/')

  if (!canUpload && attachments.length === 0) return null

  return (
    <div className="mt-3">
      {/* Image previews */}
      {attachments.some(a => isImage(a.mime_type)) && (
        <div className="flex flex-wrap gap-2 mb-2">
          {attachments.filter(a => isImage(a.mime_type)).map((att) => (
            <button
              key={att.id}
              type="button"
              onClick={() => setLightbox(att)}
              className="relative group rounded-md overflow-hidden border border-border hover:border-primary transition-colors"
              title={att.filename}
            >
              <img
                src={attachmentUrl(att.id)}
                alt={att.filename}
                className="h-20 w-20 object-cover"
              />
              <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors flex items-center justify-center">
                <Image className="h-5 w-5 text-white opacity-0 group-hover:opacity-100 transition-opacity" />
              </div>
            </button>
          ))}
        </div>
      )}

      {/* Attachment list */}
      {attachments.length > 0 && (
        <div className="space-y-1.5 mb-2">
          {attachments.map((att) => (
            <div
              key={att.id}
              className="flex items-center gap-2 rounded-md border border-border bg-muted/30 px-3 py-2 text-sm"
            >
              {fileIcon(att.mime_type)}
              {isImage(att.mime_type) ? (
                <button
                  type="button"
                  onClick={() => setLightbox(att)}
                  className="flex-1 min-w-0 font-medium hover:text-primary transition-colors truncate text-left"
                >
                  {att.filename}
                </button>
              ) : (
                <a
                  href={attachmentUrl(att.id)}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex-1 min-w-0 font-medium hover:text-primary transition-colors truncate"
                >
                  {att.filename}
                </a>
              )}
              <span className="text-xs text-muted-foreground shrink-0">
                {formatBytes(att.size_bytes)}
              </span>
              <a
                href={attachmentUrl(att.id)}
                target="_blank"
                rel="noopener noreferrer"
                className="shrink-0 text-muted-foreground hover:text-foreground transition-colors"
                title="Скачать"
              >
                <Download className="h-3.5 w-3.5" />
              </a>
              {canDelete(att) && (
                <button
                  onClick={() => handleDelete(att.id)}
                  className="shrink-0 text-muted-foreground hover:text-red-500 transition-colors"
                  title="Удалить"
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </button>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Upload button */}
      {canUpload && (
        <>
          <input
            ref={fileInputRef}
            type="file"
            className="hidden"
            onChange={handleFileChange}
          />
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            disabled={uploading}
            className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground border border-dashed border-border hover:border-primary rounded-md px-3 py-1.5 transition-colors disabled:opacity-50"
          >
            {uploading ? (
              <Loader2 className="h-3.5 w-3.5 animate-spin" />
            ) : (
              <Paperclip className="h-3.5 w-3.5" />
            )}
            {uploading ? 'Загрузка...' : 'Прикрепить файл'}
          </button>
          {error && <p className="text-xs text-red-600 mt-1">{error}</p>}
        </>
      )}

      {lightbox && (
        <ImageLightbox
          src={attachmentUrl(lightbox.id)}
          filename={lightbox.filename}
          onClose={() => setLightbox(null)}
        />
      )}
    </div>
  )
}
