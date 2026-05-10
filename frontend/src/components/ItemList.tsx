import { useEffect, useRef } from 'react'
import './ItemList.css'
import { store } from '../../wailsjs/go/models'
import { ClipboardItem } from './ClipboardItem'

interface ItemListProps {
  items: store.Item[]
  selectedId: string | null
  copiedId: string | null
  isOpening: boolean
  newItemIds: Set<string>
  deletingIds: Set<string>
  searchActive: boolean
  onSelect: (id: string) => void
  onCopy: (id: string) => void
  onPin: (id: string, pinned: boolean) => void
  onDelete: (id: string) => void
}

export function ItemList({ items, selectedId, copiedId, isOpening, newItemIds, deletingIds, searchActive, onSelect, onCopy, onPin, onDelete }: ItemListProps) {
  const pinned = items.filter(i => i.pinned)
  const recent = items.filter(i => !i.pinned)

  // build a flat index map for stagger (across both groups)
  const ordered = [...pinned, ...recent]
  const indexMap = new Map(ordered.map((item, i) => [item.id, i]))

  const listRef = useRef<HTMLDivElement>(null)
  useRubberband(listRef)

  return (
    <div className="item-list" ref={listRef}>
      {pinned.length > 0 && (
        <>
          <div className="section-label">Pinned</div>
          <ul className="item-group">
            {pinned.map(item => (
              <ClipboardItem
                key={item.id}
                item={item}
                selected={selectedId === item.id}
                justCopied={copiedId === item.id}
                isNew={newItemIds.has(item.id)}
                isDeleting={deletingIds.has(item.id)}
                staggerIndex={(isOpening || searchActive) ? (indexMap.get(item.id) ?? -1) : -1}
                onSelect={() => onSelect(item.id)}
                onCopy={() => onCopy(item.id)}
                onPin={() => onPin(item.id, !item.pinned)}
                onDelete={() => onDelete(item.id)}
              />
            ))}
          </ul>
          <div className="section-divider" />
        </>
      )}
      <div className="section-label">Recent</div>
      <ul className="item-group">
        {recent.map(item => (
          <ClipboardItem
            key={item.id}
            item={item}
            selected={selectedId === item.id}
            justCopied={copiedId === item.id}
            isNew={newItemIds.has(item.id)}
            isDeleting={deletingIds.has(item.id)}
            staggerIndex={(isOpening || searchActive) ? (indexMap.get(item.id) ?? -1) : -1}
            onSelect={() => onSelect(item.id)}
            onCopy={() => onCopy(item.id)}
            onPin={() => onPin(item.id, !item.pinned)}
            onDelete={() => onDelete(item.id)}
          />
        ))}
      </ul>
    </div>
  )
}

// Rubber-band overscroll. Native macOS bounce on inner scroll containers is
// unreliable across WebView engines, so we synthesize it: at the top/bottom
// boundary, wheel deltas translate the list with diminishing resistance, then
// snap back via requestAnimationFrame easing once the user stops scrolling.
function useRubberband(ref: React.RefObject<HTMLDivElement>) {
  useEffect(() => {
    const el = ref.current
    if (!el) return

    const MAX = 80
    const SNAP_DELAY = 80
    const SNAP_DURATION = 280

    let overscroll = 0
    let snapTimer: ReturnType<typeof setTimeout> | null = null
    let rafId = 0

    const apply = (v: number) => {
      el.style.transform = v === 0 ? '' : `translateY(${v}px)`
    }

    const snapBack = () => {
      if (rafId) cancelAnimationFrame(rafId)
      const start = performance.now()
      const startVal = overscroll
      const tick = (now: number) => {
        const t = Math.min(1, (now - start) / SNAP_DURATION)
        const eased = 1 - Math.pow(1 - t, 3) // ease-out-cubic
        overscroll = startVal * (1 - eased)
        apply(overscroll)
        if (t < 1) {
          rafId = requestAnimationFrame(tick)
        } else {
          overscroll = 0
          apply(0)
        }
      }
      rafId = requestAnimationFrame(tick)
    }

    const onWheel = (e: WheelEvent) => {
      const atTop = el.scrollTop <= 0
      const atBottom = Math.ceil(el.scrollTop + el.clientHeight) >= el.scrollHeight

      const overscrollAt = (sign: 1 | -1) => {
        e.preventDefault()
        if (rafId) cancelAnimationFrame(rafId)
        const resistance = 0.35 * (1 - Math.min(Math.abs(overscroll), MAX) / MAX)
        const next = overscroll - e.deltaY * resistance
        overscroll = sign === 1 ? Math.min(MAX, next) : Math.max(-MAX, next)
        apply(overscroll)
        if (snapTimer) clearTimeout(snapTimer)
        snapTimer = setTimeout(snapBack, SNAP_DELAY)
      }

      if (atTop && e.deltaY < 0) {
        overscrollAt(1)
      } else if (atBottom && e.deltaY > 0) {
        overscrollAt(-1)
      } else if (overscroll !== 0) {
        // direction reversed mid-bounce: snap back immediately
        if (snapTimer) clearTimeout(snapTimer)
        snapBack()
      }
    }

    el.addEventListener('wheel', onWheel, { passive: false })
    return () => {
      el.removeEventListener('wheel', onWheel)
      if (snapTimer) clearTimeout(snapTimer)
      if (rafId) cancelAnimationFrame(rafId)
    }
  }, [ref])
}
