import { useState } from 'react'
import { MarkdownContent } from './MarkdownContent'

interface Props {
  value: string
  onChange: (v: string) => void
  placeholder?: string
  rows?: number
  className?: string
}

export function MarkdownEditor({ value, onChange, placeholder, rows = 8, className = '' }: Props) {
  const [tab, setTab] = useState<'write' | 'preview'>('write')

  return (
    <div className={`rounded-md border border-border overflow-hidden ${className}`}>
      <div className="flex border-b border-border bg-muted/30">
        <button
          type="button"
          onClick={() => setTab('write')}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            tab === 'write'
              ? 'bg-background text-foreground border-b-2 border-primary -mb-px'
              : 'text-muted-foreground hover:text-foreground'
          }`}
        >
          Написать
        </button>
        <button
          type="button"
          onClick={() => setTab('preview')}
          className={`px-4 py-2 text-sm font-medium transition-colors ${
            tab === 'preview'
              ? 'bg-background text-foreground border-b-2 border-primary -mb-px'
              : 'text-muted-foreground hover:text-foreground'
          }`}
        >
          Предпросмотр
        </button>
        <div className="ml-auto flex items-center px-3">
          <span className="text-xs text-muted-foreground">Поддерживается Markdown</span>
        </div>
      </div>

      {tab === 'write' ? (
        <textarea
          rows={rows}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
          className="w-full bg-background px-3 py-2.5 text-sm focus:outline-none resize-y font-mono"
        />
      ) : (
        <div className="min-h-[120px] px-3 py-2.5 bg-background">
          {value.trim() ? (
            <MarkdownContent>{value}</MarkdownContent>
          ) : (
            <p className="text-sm text-muted-foreground italic">Нечего предпросматривать</p>
          )}
        </div>
      )}
    </div>
  )
}
