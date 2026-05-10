import { useState, useEffect, useCallback, useRef } from 'react'
import './App.css'
import { SearchBar } from './components/SearchBar'
import { ItemList } from './components/ItemList'
import { ExpandTooltip } from './components/ExpandTooltip'
import { SettingsPanel } from './components/SettingsPanel'
import { useClipboardItems } from './hooks/useClipboardItems'
import { ErrorToast } from './components/ErrorToast'
import { SearchItems } from '../wailsjs/go/main/App'
import { EventsOn, WindowHide } from '../wailsjs/runtime/runtime'
import { store } from '../wailsjs/go/models'

export default function App() {
  const { items, loading, copyItem, pinItem, deleteItem, newItemIds, deletingIds, error, clearError } = useClipboardItems()
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<store.Item[]>([])
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [copiedId, setCopiedId] = useState<string | null>(null)
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [isOpening, setIsOpening] = useState(false)
  const [tooltipVisible, setTooltipVisible] = useState(false)
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const openingTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const itemsRef = useRef(items)
  useEffect(() => { itemsRef.current = items }, [items])

  const displayItems = query ? results : items

  // Auto-select first item whenever list changes
  useEffect(() => {
    if (displayItems.length > 0 && !displayItems.find(i => i.id === selectedId)) {
      setSelectedId(displayItems[0].id)
    }
    setTooltipVisible(false)
  }, [displayItems])

  // Keep the selected item in view as the user navigates
  useEffect(() => {
    if (!selectedId) return
    const el = document.querySelector<HTMLElement>('.clipboard-item.selected')
    el?.scrollIntoView({ block: 'nearest' })
  }, [selectedId])

  // Debounced FTS search
  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current)
    if (!query.trim()) { setResults([]); return }
    debounceRef.current = setTimeout(() => {
      SearchItems(query).then(r => {
        setResults(r ?? [])
        if (r?.length) setSelectedId(r[0].id)
      })
    }, 150)
    return () => { if (debounceRef.current) clearTimeout(debounceRef.current) }
  }, [query])

  // Cleanup openingTimerRef on unmount
  useEffect(() => {
    return () => {
      if (openingTimerRef.current) clearTimeout(openingTimerRef.current)
    }
  }, [])

  // Reset on window show (hotkey re-opens panel)
  useEffect(() => {
    const off = EventsOn('wails:window-show', () => {
      setQuery('')
      setSelectedId(itemsRef.current[0]?.id ?? null)
      setIsOpening(true)
      if (openingTimerRef.current) clearTimeout(openingTimerRef.current)
      openingTimerRef.current = setTimeout(() => setIsOpening(false), 500)
      document.getElementById('search-input')?.focus()
    })
    return off
  }, [])

  // Hide on window blur (click away)
  useEffect(() => {
    const handler = () => WindowHide()
    window.addEventListener('blur', handler)
    return () => window.removeEventListener('blur', handler)
  }, [])

  const handleCopy = useCallback(async (id: string) => {
    await copyItem(id)
    setCopiedId(id)
    setTimeout(() => setCopiedId(null), 1200)
  }, [copyItem])

  const selectedItem = displayItems.find(i => i.id === selectedId) ?? null

  // Global keyboard handler
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      if (settingsOpen) {
        if (e.key === 'Escape' || (e.key === ',' && e.metaKey)) setSettingsOpen(false)
        return
      }

      const idx = displayItems.findIndex(i => i.id === selectedId)

      switch (true) {
        case e.key === 'Escape':
          if (tooltipVisible) { setTooltipVisible(false) }
          else if (query) { setQuery('') }
          else { WindowHide() }
          break
        case e.key === 'ArrowLeft':
          e.preventDefault()
          setTooltipVisible(true)
          break
        case e.key === 'ArrowRight':
          e.preventDefault()
          setTooltipVisible(false)
          break
        case e.key === 'ArrowDown':
          e.preventDefault()
          setTooltipVisible(false)
          setSelectedId(displayItems[Math.min(idx + 1, displayItems.length - 1)]?.id ?? selectedId)
          break
        case e.key === 'ArrowUp':
          e.preventDefault()
          setTooltipVisible(false)
          setSelectedId(displayItems[Math.max(idx - 1, 0)]?.id ?? selectedId)
          break
        case e.key === 'Enter' && !!selectedId:
          handleCopy(selectedId!)
          break
        case e.key === 'k' && e.metaKey:
          e.preventDefault()
          document.getElementById('search-input')?.focus()
          break
        case e.key === 'p' && e.metaKey && !!selectedId: {
          e.preventDefault()
          const item = displayItems.find(i => i.id === selectedId)
          if (item) pinItem(selectedId!, !item.pinned)
          break
        }
        case e.key === 'Backspace' && e.metaKey && !!selectedId:
          e.preventDefault()
          deleteItem(selectedId!)
          setSelectedId(displayItems[Math.max(idx - 1, 0)]?.id ?? null)
          break
        case e.key === ',' && e.metaKey:
          e.preventDefault()
          setSettingsOpen(true)
          break
      }
    }

    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [displayItems, selectedId, settingsOpen, tooltipVisible, query, handleCopy, pinItem, deleteItem])

  return (
    <div className="app">
      <ExpandTooltip item={selectedItem} visible={tooltipVisible && !settingsOpen && !loading} />
      <div className={`panel${isOpening ? ' is-opening' : ''}`}>
        <SearchBar value={query} onChange={setQuery} onEscape={WindowHide} />
        <div className="panel-body">
          {loading
            ? <div className="panel-loading">Loading…</div>
            : <ItemList
                items={displayItems}
                selectedId={selectedId}
                copiedId={copiedId}
                isOpening={isOpening}
                newItemIds={newItemIds}
                deletingIds={deletingIds}
                searchActive={!!query}
                onSelect={setSelectedId}
                onCopy={handleCopy}
                onPin={pinItem}
                onDelete={deleteItem}
              />
          }
        </div>
        <ErrorToast message={error} onDismiss={clearError} />
        <div className="panel-footer">
          <span className="footer-count">{displayItems.length} items</span>
          <div className="footer-hints">
            <span>↑↓ navigate</span>
            <span>↵ copy</span>
            <button className="footer-gear" onClick={() => setSettingsOpen(true)}>⚙</button>
          </div>
        </div>
        <SettingsPanel open={settingsOpen} onClose={() => setSettingsOpen(false)} />
      </div>
    </div>
  )
}
