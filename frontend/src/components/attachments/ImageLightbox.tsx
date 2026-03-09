import { useEffect, useRef, useState } from 'react'
import { X } from 'lucide-react'

interface Props {
  src: string
  filename: string
  onClose: () => void
}

// Минимум MARGIN пикселей картинки должно оставаться в видимой области
const MARGIN = 80

export function ImageLightbox({ src, filename, onClose }: Props) {
  const [scale, setScale] = useState(1)
  const [pos, setPos] = useState({ x: 0, y: 0 })
  const [dragging, setDragging] = useState(false)

  const overlayRef = useRef<HTMLDivElement>(null)
  const imgRef = useRef<HTMLImageElement>(null)
  const lastMouse = useRef({ x: 0, y: 0 })
  // Храним scale в ref, чтобы не было stale closure в wheel-обработчике
  const scaleRef = useRef(1)

  // Зажимает позицию так, чтобы картинка не уходила за экран полностью
  const clamp = (x: number, y: number, s: number) => {
    const el = imgRef.current
    if (!el) return { x, y }
    const vw = window.innerWidth
    const vh = window.innerHeight
    // offsetWidth/Height — размер до CSS transform (то, что нам нужно)
    const hw = (el.offsetWidth * s) / 2
    const hh = (el.offsetHeight * s) / 2
    return {
      x: Math.min(Math.max(x, MARGIN - vw / 2 - hw), vw / 2 - MARGIN + hw),
      y: Math.min(Math.max(y, MARGIN - vh / 2 - hh), vh / 2 - MARGIN + hh),
    }
  }

  // Escape
  useEffect(() => {
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [onClose])

  // non-passive wheel
  useEffect(() => {
    const el = overlayRef.current
    if (!el) return
    const handler = (e: WheelEvent) => {
      e.preventDefault()
      const factor = e.deltaY < 0 ? 1.15 : 1 / 1.15
      const newScale = Math.min(Math.max(scaleRef.current * factor, 0.1), 15)
      scaleRef.current = newScale
      setScale(newScale)
      // После изменения масштаба сразу пережимаем позицию
      setPos(p => clamp(p.x, p.y, newScale))
    }
    el.addEventListener('wheel', handler, { passive: false })
    return () => el.removeEventListener('wheel', handler)
  }, [])

  const onMouseDown = (e: React.MouseEvent) => {
    e.preventDefault()
    setDragging(true)
    lastMouse.current = { x: e.clientX, y: e.clientY }
  }

  const onMouseMove = (e: React.MouseEvent) => {
    if (!dragging) return
    const dx = e.clientX - lastMouse.current.x
    const dy = e.clientY - lastMouse.current.y
    lastMouse.current = { x: e.clientX, y: e.clientY }
    setPos(p => clamp(p.x + dx, p.y + dy, scaleRef.current))
  }

  const onMouseUp = () => setDragging(false)

  return (
    <div
      ref={overlayRef}
      className="fixed inset-0 z-50 bg-black/88 flex items-center justify-center select-none"
      onMouseMove={onMouseMove}
      onMouseUp={onMouseUp}
      onMouseLeave={onMouseUp}
    >
      <div className="absolute inset-0" onClick={onClose} />

      <button
        onClick={onClose}
        className="absolute top-4 right-4 z-10 text-white/70 hover:text-white transition-colors"
        title="Закрыть (Esc)"
      >
        <X className="h-6 w-6" />
      </button>

      <p className="absolute top-4 left-4 z-10 text-white/60 text-sm truncate max-w-[60vw]">
        {filename}
      </p>

      <p className="absolute bottom-4 left-1/2 -translate-x-1/2 z-10 text-white/40 text-xs pointer-events-none">
        Колесо — масштаб · Перетащи — перемещение · Esc — закрыть
      </p>

      <img
        ref={imgRef}
        src={src}
        alt={filename}
        draggable={false}
        onClick={e => e.stopPropagation()}
        onMouseDown={onMouseDown}
        style={{
          transform: `translate(${pos.x}px, ${pos.y}px) scale(${scale})`,
          transformOrigin: 'center center',
          cursor: dragging ? 'grabbing' : 'grab',
          maxWidth: '90vw',
          maxHeight: '90vh',
          objectFit: 'contain',
          position: 'relative',
          zIndex: 10,
        }}
      />
    </div>
  )
}
