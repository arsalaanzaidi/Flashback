import { useRef, useEffect, KeyboardEvent } from 'react'
import './SearchBar.css'

interface SearchBarProps {
  value: string
  onChange: (q: string) => void
  onEscape: () => void
}

export function SearchBar({ value, onChange, onEscape }: SearchBarProps) {
  const ref = useRef<HTMLInputElement>(null)

  useEffect(() => { ref.current?.focus() }, [])

  function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
    if (e.key === 'Escape') {
      if (value !== '') {
        onChange('')
      } else {
        onEscape()
      }
    }
  }

  return (
    <div className="search-bar">
      <span className="search-icon">⌕</span>
      <input
        ref={ref}
        id="search-input"
        type="text"
        className="search-input"
        placeholder="Search clipboard history..."
        value={value}
        onChange={e => onChange(e.target.value)}
        onKeyDown={handleKeyDown}
      />
      <kbd className="search-hint">⌘K</kbd>
    </div>
  )
}
