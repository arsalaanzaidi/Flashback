# Clipboard Manager Frontend — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the complete React frontend for the clipboard-manager Wails app — search bar, item list, expand tooltip, and settings panel.

**Architecture:** Single-view React app, no routing. All state lives in `App.tsx` and flows down via props. Wails bindings and events are encapsulated in `useClipboardItems`. Components are presentational and independently testable.

**Tech Stack:** React 18, TypeScript, Vite 3, Vitest, @testing-library/react, plain CSS

---

## File Map

| File | Action | Purpose |
|------|--------|---------|
| `frontend/vite.config.ts` | Modify | Add vitest test block |
| `frontend/package.json` | Modify | Add test deps + scripts |
| `frontend/src/test-utils/setup.ts` | Create | @testing-library/jest-dom setup |
| `frontend/src/utils/typeConfig.ts` | Create | Badge colour/icon/label for 23 types |
| `frontend/src/utils/typeConfig.test.ts` | Create | |
| `frontend/src/utils/format.ts` | Create | Timestamp formatting, content preview |
| `frontend/src/utils/format.test.ts` | Create | |
| `frontend/src/components/TypeBadge.tsx` | Create | Type pill component |
| `frontend/src/components/TypeBadge.css` | Create | |
| `frontend/src/components/TypeBadge.test.tsx` | Create | |
| `frontend/src/components/SearchBar.tsx` | Create | Search input with debounce + Escape |
| `frontend/src/components/SearchBar.css` | Create | |
| `frontend/src/components/SearchBar.test.tsx` | Create | |
| `frontend/src/components/ClipboardItem.tsx` | Create | Single row: badge, preview, actions |
| `frontend/src/components/ClipboardItem.css` | Create | |
| `frontend/src/components/ClipboardItem.test.tsx` | Create | |
| `frontend/src/components/ItemList.tsx` | Create | Pinned + Recent sections |
| `frontend/src/components/ItemList.css` | Create | |
| `frontend/src/components/ExpandTooltip.tsx` | Create | Full-content overlay to left of panel |
| `frontend/src/components/ExpandTooltip.css` | Create | |
| `frontend/src/components/SettingsPanel.tsx` | Create | Settings slide-in |
| `frontend/src/components/SettingsPanel.css` | Create | |
| `frontend/src/hooks/useClipboardItems.ts` | Create | Wails data + events + mutations |
| `frontend/src/hooks/useClipboardItems.test.ts` | Create | |
| `frontend/src/App.tsx` | Replace | Root: wires all components + keyboard nav |
| `frontend/src/App.css` | Replace | Layout + panel styles |
| `frontend/src/style.css` | Replace | Global dark theme base |

---

### Task 1: Bootstrap — Vitest + folder structure

**Files:**
- Modify: `frontend/package.json`
- Modify: `frontend/vite.config.ts`
- Create: `frontend/src/test-utils/setup.ts`

- [ ] **Step 1: Install test dependencies**

```bash
cd /Users/arsalaanabbas.zaidi/Projects/clipboard-manager/frontend
npm install -D vitest @testing-library/react @testing-library/user-event @testing-library/jest-dom jsdom
```

Expected: packages install cleanly.

- [ ] **Step 2: Update `frontend/vite.config.ts`**

```typescript
/// <reference types="vitest" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test-utils/setup.ts'],
  },
})
```

- [ ] **Step 3: Create `frontend/src/test-utils/setup.ts`**

```typescript
import '@testing-library/jest-dom'
```

- [ ] **Step 4: Add test scripts to `frontend/package.json`**

In the `"scripts"` block add:
```json
"test": "vitest run",
"test:watch": "vitest"
```

- [ ] **Step 5: Create directories and run**

```bash
mkdir -p frontend/src/utils frontend/src/components frontend/src/hooks frontend/src/test-utils
cd frontend && npm test
```

Expected: `No test files found` or `0 tests` — confirms Vitest works.

- [ ] **Step 6: Commit**

```bash
git add frontend/package.json frontend/vite.config.ts frontend/src/test-utils/
git commit -m "feat: add vitest + testing-library, scaffold frontend dirs"
```

---

### Task 2: typeConfig utility

**Files:**
- Create: `frontend/src/utils/typeConfig.ts`
- Create: `frontend/src/utils/typeConfig.test.ts`

- [ ] **Step 1: Write failing tests — `frontend/src/utils/typeConfig.test.ts`**

```typescript
import { describe, it, expect } from 'vitest'
import { getTypeConfig } from './typeConfig'

describe('getTypeConfig', () => {
  it('returns green/mono config for CODE with subtype as label', () => {
    const cfg = getTypeConfig('CODE', 'go')
    expect(cfg.label).toBe('GO')
    expect(cfg.color).toBe('#4ade80')
    expect(cfg.mono).toBe(true)
  })

  it('returns "CODE" label when CODE subtype is empty', () => {
    expect(getTypeConfig('CODE', '').label).toBe('CODE')
  })

  it('returns blue config for URL', () => {
    const cfg = getTypeConfig('URL')
    expect(cfg.color).toBe('#60a5fa')
    expect(cfg.mono).toBe(false)
  })

  it('returns amber config for API_KEY', () => {
    expect(getTypeConfig('API_KEY').color).toBe('#fbbf24')
  })

  it('returns grey/TEXT for unknown type', () => {
    const cfg = getTypeConfig('UNKNOWN')
    expect(cfg.label).toBe('TEXT')
    expect(cfg.color).toBe('#888')
  })

  it('returns mono=true for JSON', () => {
    expect(getTypeConfig('JSON').mono).toBe(true)
  })

  it('returns mono=false for EMAIL', () => {
    expect(getTypeConfig('EMAIL').mono).toBe(false)
  })
})
```

- [ ] **Step 2: Run — confirm fail**

```bash
cd frontend && npm test
```

Expected: `Cannot find module './typeConfig'`

- [ ] **Step 3: Implement `frontend/src/utils/typeConfig.ts`**

```typescript
export interface TypeConfig {
  label: string
  icon: string
  color: string
  bg: string
  mono: boolean
}

export function getTypeConfig(type: string, subtype = ''): TypeConfig {
  switch (type) {
    case 'URL':       return { label: 'URL',     icon: '🔗', color: '#60a5fa', bg: '#0f1e3a', mono: false }
    case 'EMAIL':     return { label: 'EMAIL',   icon: '📧', color: '#f87171', bg: '#2a0a0a', mono: false }
    case 'PHONE':     return { label: 'PHONE',   icon: '📱', color: '#f87171', bg: '#2a0a0a', mono: false }
    case 'IP':        return { label: 'IP',      icon: '🌐', color: '#f87171', bg: '#2a0a0a', mono: false }
    case 'IMAGE':     return { label: 'IMAGE',   icon: '🖼', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'PDF':       return { label: 'PDF',     icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'RICH_TEXT': return { label: 'RTF',     icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'HTML':      return { label: 'HTML',    icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: true  }
    case 'COLOR':     return { label: 'COLOR',   icon: '🎨', color: '#facc15', bg: '#1a1a0a', mono: true  }
    case 'COLOR_CODE':return { label: 'COLOR',   icon: '🎨', color: '#facc15', bg: '#1a1a0a', mono: true  }
    case 'FILE_REF':  return { label: 'FILE',    icon: '📁', color: '#34d399', bg: '#0a2a1a', mono: false }
    case 'FILE_PATH': return { label: 'PATH',    icon: '📁', color: '#34d399', bg: '#0a2a1a', mono: true  }
    case 'HASH':      return { label: 'HASH',    icon: '#',  color: '#fbbf24', bg: '#3a2f0d', mono: true  }
    case 'JWT':       return { label: 'JWT',     icon: '🔑', color: '#fbbf24', bg: '#3a2f0d', mono: true  }
    case 'SSH_KEY':   return { label: 'SSH KEY', icon: '🔏', color: '#fbbf24', bg: '#3a2f0d', mono: true  }
    case 'API_KEY':   return { label: 'API KEY', icon: '🔑', color: '#fbbf24', bg: '#3a2f0d', mono: true  }
    case 'JSON':      return { label: 'JSON',    icon: '{}', color: '#4ade80', bg: '#0f2a1a', mono: true  }
    case 'XML':       return { label: 'XML',     icon: '<>', color: '#4ade80', bg: '#0f2a1a', mono: true  }
    case 'YAML':      return { label: 'YAML',    icon: '—',  color: '#4ade80', bg: '#0f2a1a', mono: true  }
    case 'SQL':       return { label: 'SQL',     icon: '🗃', color: '#4ade80', bg: '#0f2a1a', mono: true  }
    case 'CODE':      return {
      label: subtype ? subtype.toUpperCase() : 'CODE',
      icon: '💻', color: '#4ade80', bg: '#14290a', mono: true,
    }
    case 'UUID':      return { label: 'UUID',    icon: '🔢', color: '#888',    bg: '#222',    mono: true  }
    case 'BASE64':    return { label: 'B64',     icon: '📦', color: '#888',    bg: '#222',    mono: true  }
    case 'MARKDOWN':  return { label: 'MD',      icon: '📝', color: '#888',    bg: '#222',    mono: false }
    default:          return { label: 'TEXT',    icon: '📝', color: '#888',    bg: '#222',    mono: false }
  }
}
```

- [ ] **Step 4: Run — confirm pass**

```bash
cd frontend && npm test
```

Expected: `7 tests passed`

- [ ] **Step 5: Commit**

```bash
git add frontend/src/utils/typeConfig.ts frontend/src/utils/typeConfig.test.ts
git commit -m "feat: typeConfig utility — badge colors for all 23 content types"
```

---

### Task 3: format utility

**Files:**
- Create: `frontend/src/utils/format.ts`
- Create: `frontend/src/utils/format.test.ts`

- [ ] **Step 1: Write failing tests — `frontend/src/utils/format.test.ts`**

```typescript
import { describe, it, expect, beforeAll, afterAll, vi } from 'vitest'
import { formatAge, getPreview } from './format'

describe('formatAge', () => {
  const now = new Date('2026-04-12T12:00:00.000Z').getTime()

  beforeAll(() => { vi.setSystemTime(now) })
  afterAll(() => { vi.useRealTimers() })

  it('returns "just now" for < 60s', () => {
    expect(formatAge(now - 30_000)).toBe('just now')
  })
  it('returns minutes for < 1h', () => {
    expect(formatAge(now - 5 * 60_000)).toBe('5m')
  })
  it('returns hours for < 1d', () => {
    expect(formatAge(now - 3 * 3_600_000)).toBe('3h')
  })
  it('returns days for < 30d', () => {
    expect(formatAge(now - 2 * 86_400_000)).toBe('2d')
  })
  it('returns months for >= 30d', () => {
    expect(formatAge(now - 35 * 86_400_000)).toBe('1mo')
  })
})

describe('getPreview', () => {
  it('returns short content unchanged', () => {
    expect(getPreview('hello world')).toBe('hello world')
  })
  it('truncates at 80 chars with ellipsis', () => {
    const result = getPreview('a'.repeat(100))
    expect(result).toHaveLength(81)
    expect(result.endsWith('…')).toBe(true)
  })
  it('collapses internal whitespace', () => {
    expect(getPreview('hello   \n   world')).toBe('hello world')
  })
})
```

- [ ] **Step 2: Run — confirm fail**

```bash
cd frontend && npm test
```

Expected: `Cannot find module './format'`

- [ ] **Step 3: Implement `frontend/src/utils/format.ts`**

```typescript
export function formatAge(ms: number): string {
  const diff = Math.floor((Date.now() - ms) / 1000)
  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h`
  const days = Math.floor(diff / 86400)
  if (days < 30) return `${days}d`
  return `${Math.floor(days / 30)}mo`
}

export function getPreview(content: string): string {
  const text = content.replace(/\s+/g, ' ').trim()
  return text.length > 80 ? text.slice(0, 80) + '…' : text
}
```

- [ ] **Step 4: Run — confirm pass**

```bash
cd frontend && npm test
```

Expected: all tests pass

- [ ] **Step 5: Commit**

```bash
git add frontend/src/utils/format.ts frontend/src/utils/format.test.ts
git commit -m "feat: format utilities (age, content preview)"
```

---

### Task 4: TypeBadge component

**Files:**
- Create: `frontend/src/components/TypeBadge.tsx`
- Create: `frontend/src/components/TypeBadge.css`
- Create: `frontend/src/components/TypeBadge.test.tsx`

- [ ] **Step 1: Write failing tests — `frontend/src/components/TypeBadge.test.tsx`**

```tsx
import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { TypeBadge } from './TypeBadge'

describe('TypeBadge', () => {
  it('renders URL label', () => {
    render(<TypeBadge type="URL" />)
    expect(screen.getByText(/URL/)).toBeInTheDocument()
  })
  it('renders CODE subtype as uppercase label', () => {
    render(<TypeBadge type="CODE" subtype="python" />)
    expect(screen.getByText(/PYTHON/)).toBeInTheDocument()
  })
  it('falls back to TEXT for unknown type', () => {
    render(<TypeBadge type="UNKNOWN" />)
    expect(screen.getByText(/TEXT/)).toBeInTheDocument()
  })
})
```

- [ ] **Step 2: Run — confirm fail**

```bash
cd frontend && npm test
```

Expected: `Cannot find module './TypeBadge'`

- [ ] **Step 3: Implement `frontend/src/components/TypeBadge.tsx`**

```tsx
import './TypeBadge.css'
import { getTypeConfig } from '../utils/typeConfig'

interface TypeBadgeProps {
  type: string
  subtype?: string
}

export function TypeBadge({ type, subtype = '' }: TypeBadgeProps) {
  const cfg = getTypeConfig(type, subtype)
  return (
    <span
      className="type-badge"
      style={{ color: cfg.color, backgroundColor: cfg.bg }}
    >
      {cfg.icon} {cfg.label}
    </span>
  )
}
```

- [ ] **Step 4: Create `frontend/src/components/TypeBadge.css`**

```css
.type-badge {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 9px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
  white-space: nowrap;
  flex-shrink: 0;
}
```

- [ ] **Step 5: Run — confirm pass**

```bash
cd frontend && npm test
```

Expected: `3 tests passed`

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/TypeBadge.tsx frontend/src/components/TypeBadge.css frontend/src/components/TypeBadge.test.tsx
git commit -m "feat: TypeBadge component"
```

---

### Task 5: SearchBar component

**Files:**
- Create: `frontend/src/components/SearchBar.tsx`
- Create: `frontend/src/components/SearchBar.css`
- Create: `frontend/src/components/SearchBar.test.tsx`

- [ ] **Step 1: Write failing tests — `frontend/src/components/SearchBar.test.tsx`**

```tsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { SearchBar } from './SearchBar'

describe('SearchBar', () => {
  it('renders placeholder text', () => {
    render(<SearchBar value="" onChange={vi.fn()} onEscape={vi.fn()} />)
    expect(screen.getByPlaceholderText('Search clipboard history...')).toBeInTheDocument()
  })

  it('calls onChange when user types', async () => {
    const onChange = vi.fn()
    render(<SearchBar value="" onChange={onChange} onEscape={vi.fn()} />)
    await userEvent.type(screen.getByRole('textbox'), 'hello')
    expect(onChange).toHaveBeenCalledWith('hello')
  })

  it('calls onEscape when Escape pressed on empty input', async () => {
    const onEscape = vi.fn()
    render(<SearchBar value="" onChange={vi.fn()} onEscape={onEscape} />)
    await userEvent.keyboard('{Escape}')
    expect(onEscape).toHaveBeenCalled()
  })

  it('clears input (calls onChange("")) on Escape when non-empty, does not call onEscape', async () => {
    const onChange = vi.fn()
    const onEscape = vi.fn()
    render(<SearchBar value="hello" onChange={onChange} onEscape={onEscape} />)
    await userEvent.keyboard('{Escape}')
    expect(onChange).toHaveBeenCalledWith('')
    expect(onEscape).not.toHaveBeenCalled()
  })
})
```

- [ ] **Step 2: Run — confirm fail**

```bash
cd frontend && npm test
```

Expected: `Cannot find module './SearchBar'`

- [ ] **Step 3: Implement `frontend/src/components/SearchBar.tsx`**

```tsx
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
```

- [ ] **Step 4: Create `frontend/src/components/SearchBar.css`**

```css
.search-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-bottom: 1px solid #222;
  background: #161616;
  flex-shrink: 0;
}

.search-icon {
  color: #555;
  font-size: 14px;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #ddd;
  font-size: 13px;
  font-family: -apple-system, sans-serif;
  caret-color: #4ade80;
}

.search-input::placeholder {
  color: #444;
}

.search-hint {
  background: #222;
  color: #555;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  border: none;
  flex-shrink: 0;
}
```

- [ ] **Step 5: Run — confirm pass**

```bash
cd frontend && npm test
```

Expected: `4 tests passed`

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/SearchBar.tsx frontend/src/components/SearchBar.css frontend/src/components/SearchBar.test.tsx
git commit -m "feat: SearchBar component with Escape handling and auto-focus"
```

---

### Task 6: ClipboardItem component

**Files:**
- Create: `frontend/src/components/ClipboardItem.tsx`
- Create: `frontend/src/components/ClipboardItem.css`
- Create: `frontend/src/components/ClipboardItem.test.tsx`

- [ ] **Step 1: Write failing tests — `frontend/src/components/ClipboardItem.test.tsx`**

```tsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { ClipboardItem } from './ClipboardItem'
import { store } from '../../wailsjs/go/models'

function makeItem(overrides: Partial<store.Item> = {}): store.Item {
  const item = new store.Item()
  item.id = 'test-id'
  item.content = 'https://example.com'
  item.type = 'URL'
  item.subtype = ''
  item.pinned = false
  item.copiedAt = Date.now() - 60_000
  item.createdAt = Date.now() - 60_000
  item.charCount = 19
  item.imagePath = ''
  item.thumbBase64 = ''
  return Object.assign(item, overrides)
}

const defaultProps = {
  selected: false,
  justCopied: false,
  onSelect: vi.fn(),
  onCopy: vi.fn(),
  onPin: vi.fn(),
  onDelete: vi.fn(),
}

describe('ClipboardItem', () => {
  it('renders the type badge', () => {
    render(<ClipboardItem item={makeItem()} {...defaultProps} />)
    expect(screen.getByText(/URL/)).toBeInTheDocument()
  })

  it('renders content preview', () => {
    render(<ClipboardItem item={makeItem({ content: 'https://example.com' })} {...defaultProps} />)
    expect(screen.getByText('https://example.com')).toBeInTheDocument()
  })

  it('shows "✓ Copied" when justCopied and selected', () => {
    render(<ClipboardItem item={makeItem()} {...defaultProps} selected={true} justCopied={true} />)
    expect(screen.getByText(/Copied/)).toBeInTheDocument()
  })

  it('calls onCopy when clicked', async () => {
    const onCopy = vi.fn()
    render(<ClipboardItem item={makeItem()} {...defaultProps} onCopy={onCopy} />)
    await userEvent.click(screen.getByRole('listitem'))
    expect(onCopy).toHaveBeenCalled()
  })

  it('calls onPin when pin button clicked on selected item', async () => {
    const onPin = vi.fn()
    render(<ClipboardItem item={makeItem()} {...defaultProps} selected={true} onPin={onPin} />)
    await userEvent.click(screen.getByTitle(/pin/i))
    expect(onPin).toHaveBeenCalled()
  })
})
```

- [ ] **Step 2: Run — confirm fail**

```bash
cd frontend && npm test
```

Expected: `Cannot find module './ClipboardItem'`

- [ ] **Step 3: Implement `frontend/src/components/ClipboardItem.tsx`**

```tsx
import './ClipboardItem.css'
import { store } from '../../wailsjs/go/models'
import { TypeBadge } from './TypeBadge'
import { formatAge, getPreview } from '../utils/format'
import { getTypeConfig } from '../utils/typeConfig'

interface ClipboardItemProps {
  item: store.Item
  selected: boolean
  justCopied: boolean
  onSelect: () => void
  onCopy: () => void
  onPin: () => void
  onDelete: () => void
}

export function ClipboardItem({
  item, selected, justCopied, onSelect, onCopy, onPin, onDelete,
}: ClipboardItemProps) {
  const cfg = getTypeConfig(item.type, item.subtype)

  return (
    <li
      role="listitem"
      className={`clipboard-item ${selected ? 'selected' : ''}`}
      style={{
        borderLeftColor: selected ? cfg.color : 'transparent',
        background: selected ? `${cfg.bg}55` : '',
      }}
      onMouseEnter={onSelect}
      onClick={onCopy}
    >
      <TypeBadge type={item.type} subtype={item.subtype} />

      {item.type === 'IMAGE' && item.thumbBase64 ? (
        <img className="item-thumb" src={item.thumbBase64} alt="" />
      ) : (item.type === 'COLOR' || item.type === 'COLOR_CODE') ? (
        <span className="item-swatch" style={{ background: item.content }} />
      ) : null}

      <span className={`item-preview ${cfg.mono ? 'mono' : ''}`}>
        {getPreview(item.content)}
      </span>

      {selected ? (
        <div className="item-actions">
          {justCopied && (
            <span className="item-copied" style={{ color: cfg.color, background: cfg.bg }}>
              ✓ Copied
            </span>
          )}
          <button className="action-btn" title="Pin" onClick={e => { e.stopPropagation(); onPin() }}>
            📌
          </button>
          <button className="action-btn action-delete" title="Delete" onClick={e => { e.stopPropagation(); onDelete() }}>
            ✕
          </button>
        </div>
      ) : (
        <span className="item-age">{formatAge(item.copiedAt)}</span>
      )}
    </li>
  )
}
```

- [ ] **Step 4: Create `frontend/src/components/ClipboardItem.css`**

```css
.clipboard-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  margin: 0 8px 3px;
  border-radius: 6px;
  border-left: 2px solid transparent;
  cursor: pointer;
  list-style: none;
  min-width: 0;
  transition: background 80ms ease;
}

.item-preview {
  flex: 1;
  font-size: 12px;
  color: #aaa;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.item-preview.mono { font-family: 'SF Mono', 'Menlo', monospace; }
.clipboard-item.selected .item-preview { color: #ddd; }

.item-thumb {
  width: 36px;
  height: 24px;
  object-fit: cover;
  border-radius: 3px;
  flex-shrink: 0;
}

.item-swatch {
  width: 16px;
  height: 16px;
  border-radius: 3px;
  flex-shrink: 0;
}

.item-age {
  font-size: 10px;
  color: #444;
  white-space: nowrap;
  flex-shrink: 0;
}

.item-actions {
  display: flex;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
}

.item-copied {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}

.action-btn {
  background: #2a2a2a;
  border: none;
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 10px;
  color: #aaa;
  cursor: pointer;
}

.action-btn:hover { background: #333; }
.action-delete { color: #f87171; }
```

- [ ] **Step 5: Run — confirm pass**

```bash
cd frontend && npm test
```

Expected: all tests pass

- [ ] **Step 6: Commit**

```bash
git add frontend/src/components/ClipboardItem.tsx frontend/src/components/ClipboardItem.css frontend/src/components/ClipboardItem.test.tsx
git commit -m "feat: ClipboardItem — badge, preview, thumb, copy/pin/delete actions"
```

---

### Task 7: ItemList component

**Files:**
- Create: `frontend/src/components/ItemList.tsx`
- Create: `frontend/src/components/ItemList.css`

No dedicated tests — pure composition of already-tested components with no logic.

- [ ] **Step 1: Implement `frontend/src/components/ItemList.tsx`**

```tsx
import './ItemList.css'
import { store } from '../../wailsjs/go/models'
import { ClipboardItem } from './ClipboardItem'

interface ItemListProps {
  items: store.Item[]
  selectedId: string | null
  copiedId: string | null
  onSelect: (id: string) => void
  onCopy: (id: string) => void
  onPin: (id: string, pinned: boolean) => void
  onDelete: (id: string) => void
}

export function ItemList({ items, selectedId, copiedId, onSelect, onCopy, onPin, onDelete }: ItemListProps) {
  const pinned = items.filter(i => i.pinned)
  const recent = items.filter(i => !i.pinned)

  return (
    <div className="item-list">
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
```

- [ ] **Step 2: Create `frontend/src/components/ItemList.css`**

```css
.item-list {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
  min-height: 0;
}

.item-list::-webkit-scrollbar { width: 4px; }
.item-list::-webkit-scrollbar-track { background: transparent; }
.item-list::-webkit-scrollbar-thumb { background: #333; border-radius: 2px; }

.section-label {
  padding: 4px 14px;
  font-size: 10px;
  font-weight: 600;
  letter-spacing: 0.8px;
  text-transform: uppercase;
  color: #555;
}

.section-divider {
  height: 1px;
  background: #222;
  margin: 6px 14px;
}

.item-group { margin: 0; padding: 0; }
```

- [ ] **Step 3: Run tests — verify nothing broken**

```bash
cd frontend && npm test
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/ItemList.tsx frontend/src/components/ItemList.css
git commit -m "feat: ItemList with pinned/recent sections"
```

---

### Task 8: ExpandTooltip component

**Files:**
- Create: `frontend/src/components/ExpandTooltip.tsx`
- Create: `frontend/src/components/ExpandTooltip.css`

- [ ] **Step 1: Implement `frontend/src/components/ExpandTooltip.tsx`**

```tsx
import './ExpandTooltip.css'
import { store } from '../../wailsjs/go/models'
import { getTypeConfig } from '../utils/typeConfig'
import { formatAge } from '../utils/format'

interface ExpandTooltipProps {
  item: store.Item | null
  visible: boolean
}

export function ExpandTooltip({ item, visible }: ExpandTooltipProps) {
  if (!item || !visible) return null
  const cfg = getTypeConfig(item.type, item.subtype)

  function renderBody() {
    if (item!.type === 'IMAGE') {
      return (
        <div className="tooltip-image-wrap">
          {item!.thumbBase64
            ? <img src={item!.thumbBase64} alt="" className="tooltip-image" />
            : <div className="tooltip-image-placeholder">{item!.imagePath.split('/').pop()}</div>}
        </div>
      )
    }
    if (item!.type === 'COLOR' || item!.type === 'COLOR_CODE') {
      return (
        <div className="tooltip-color-wrap">
          <div className="color-swatch" style={{ background: item!.content }} />
          <div className="color-table">
            <div className="color-row"><span className="color-key">HEX</span><span className="color-val">{item!.content}</span></div>
          </div>
        </div>
      )
    }
    return (
      <pre className={`tooltip-text ${cfg.mono ? 'mono' : ''}`}>
        {item!.content}
      </pre>
    )
  }

  return (
    <div className="expand-tooltip">
      <div className="tooltip-header">
        <span className="tooltip-badge" style={{ color: cfg.color, background: cfg.bg }}>
          {cfg.icon} {cfg.label}
        </span>
        <span className="tooltip-meta">full content</span>
      </div>
      <div className="tooltip-body">{renderBody()}</div>
      <div className="tooltip-footer">
        <span>{item.charCount > 0 ? `${item.charCount} chars` : ''}</span>
        <span>{formatAge(item.copiedAt)}</span>
      </div>
    </div>
  )
}
```

- [ ] **Step 2: Create `frontend/src/components/ExpandTooltip.css`**

```css
.expand-tooltip {
  position: absolute;
  /* right edge flush with panel left edge: panel is 360px centered in 900px window */
  right: calc(50% + 181px);
  top: 52px;
  width: 260px;
  background: #141414;
  border: 1px solid #2a2a2a;
  border-right: none;
  border-radius: 12px 0 0 12px;
  overflow: hidden;
  box-shadow: -16px 8px 40px rgba(0, 0, 0, 0.6);
  animation: tooltip-in 100ms ease;
  z-index: 10;
  font-family: -apple-system, sans-serif;
}

@keyframes tooltip-in {
  from { opacity: 0; transform: translateX(6px); }
  to   { opacity: 1; transform: translateX(0); }
}

.tooltip-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: #161616;
  border-bottom: 1px solid #222;
}

.tooltip-badge {
  font-size: 9px;
  font-weight: 700;
  padding: 2px 6px;
  border-radius: 4px;
}

.tooltip-meta { font-size: 10px; color: #444; }

.tooltip-body {
  padding: 12px;
  max-height: 200px;
  overflow-y: auto;
}

.tooltip-text {
  margin: 0;
  color: #ccc;
  font-size: 11px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: inherit;
}

.tooltip-text.mono { font-family: 'SF Mono', 'Menlo', monospace; }

.tooltip-image-wrap { background: #111; border-radius: 4px; overflow: hidden; }
.tooltip-image { width: 100%; display: block; }
.tooltip-image-placeholder {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255,255,255,0.2);
  font-size: 11px;
}

.tooltip-color-wrap { display: flex; flex-direction: column; gap: 8px; }
.color-swatch { height: 56px; border-radius: 6px; }
.color-table { display: flex; flex-direction: column; gap: 4px; }
.color-row { display: flex; justify-content: space-between; }
.color-key { color: #555; font-size: 10px; font-family: monospace; }
.color-val { color: #ddd; font-size: 11px; font-family: monospace; }

.tooltip-footer {
  display: flex;
  justify-content: space-between;
  padding: 6px 12px;
  background: #161616;
  border-top: 1px solid #222;
  font-size: 10px;
  color: #444;
}
```

- [ ] **Step 3: Run tests — verify nothing broken**

```bash
cd frontend && npm test
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/ExpandTooltip.tsx frontend/src/components/ExpandTooltip.css
git commit -m "feat: ExpandTooltip — full content overlay for selected item"
```

---

### Task 9: SettingsPanel component

**Files:**
- Create: `frontend/src/components/SettingsPanel.tsx`
- Create: `frontend/src/components/SettingsPanel.css`

- [ ] **Step 1: Implement `frontend/src/components/SettingsPanel.tsx`**

```tsx
import { useState, useEffect } from 'react'
import './SettingsPanel.css'
import { GetSettings, SaveSettings } from '../../wailsjs/go/main/App'
import { store } from '../../wailsjs/go/models'

interface SettingsPanelProps {
  open: boolean
  onClose: () => void
}

export function SettingsPanel({ open, onClose }: SettingsPanelProps) {
  const [settings, setSettings] = useState<store.Settings | null>(null)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (open) GetSettings().then(setSettings)
  }, [open])

  async function handleSave() {
    if (!settings) return
    setSaving(true)
    try { await SaveSettings(settings); onClose() }
    finally { setSaving(false) }
  }

  if (!open) return null

  return (
    <div className="settings-panel">
      <div className="settings-header">
        <span className="settings-title">Settings</span>
        <button className="settings-close" onClick={onClose}>✕</button>
      </div>

      {settings ? (
        <div className="settings-body">
          <div className="settings-field">
            <label className="settings-label">Retention</label>
            <div className="settings-row">
              <select
                className="settings-select"
                value={settings.retentionMode}
                onChange={e => setSettings({ ...settings, retentionMode: e.target.value })}
              >
                <option value="unlimited">Unlimited</option>
                <option value="count">Last N items</option>
                <option value="days">Last N days</option>
              </select>
              {settings.retentionMode !== 'unlimited' && (
                <input
                  type="number"
                  className="settings-number"
                  min={1}
                  value={settings.retentionValue}
                  onChange={e => setSettings({ ...settings, retentionValue: parseInt(e.target.value) || 1 })}
                />
              )}
            </div>
          </div>

          <div className="settings-field">
            <label className="settings-label">Launch at login</label>
            <input
              type="checkbox"
              className="settings-checkbox"
              checked={settings.launchAtLogin}
              onChange={e => setSettings({ ...settings, launchAtLogin: e.target.checked })}
            />
          </div>

          <div className="settings-field">
            <label className="settings-label">Global shortcut</label>
            <kbd className="settings-kbd">{settings.globalShortcut}</kbd>
          </div>
        </div>
      ) : (
        <div className="settings-loading">Loading…</div>
      )}

      <div className="settings-footer">
        <button className="btn-cancel" onClick={onClose}>Cancel</button>
        <button className="btn-save" onClick={handleSave} disabled={saving}>
          {saving ? 'Saving…' : 'Save'}
        </button>
      </div>
    </div>
  )
}
```

- [ ] **Step 2: Create `frontend/src/components/SettingsPanel.css`**

```css
.settings-panel {
  position: absolute;
  inset: 0;
  background: #141414;
  border-radius: 14px;
  display: flex;
  flex-direction: column;
  animation: slide-in 120ms ease;
  z-index: 20;
  font-family: -apple-system, sans-serif;
}

@keyframes slide-in {
  from { transform: translateX(20px); opacity: 0; }
  to   { transform: translateX(0);    opacity: 1; }
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid #222;
  background: #161616;
}

.settings-title { font-size: 13px; font-weight: 600; color: #ddd; }

.settings-close {
  background: none;
  border: none;
  color: #555;
  font-size: 13px;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
}
.settings-close:hover { background: #2a2a2a; color: #aaa; }

.settings-body { flex: 1; padding: 16px; display: flex; flex-direction: column; gap: 20px; }

.settings-field { display: flex; flex-direction: column; gap: 8px; }

.settings-label {
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  color: #666;
}

.settings-row { display: flex; gap: 8px; align-items: center; }

.settings-select {
  flex: 1;
  background: #222;
  border: 1px solid #333;
  border-radius: 6px;
  color: #ddd;
  font-size: 12px;
  padding: 6px 10px;
  outline: none;
}

.settings-number {
  width: 70px;
  background: #222;
  border: 1px solid #333;
  border-radius: 6px;
  color: #ddd;
  font-size: 12px;
  padding: 6px 10px;
  outline: none;
  text-align: center;
}

.settings-checkbox { width: 16px; height: 16px; accent-color: #4ade80; }

.settings-kbd {
  background: #222;
  border: none;
  border-radius: 6px;
  color: #aaa;
  font-size: 12px;
  padding: 4px 10px;
  font-family: monospace;
}

.settings-loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #555;
  font-size: 13px;
}

.settings-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid #222;
  background: #161616;
}

.btn-cancel {
  background: #222;
  border: none;
  border-radius: 6px;
  color: #aaa;
  font-size: 12px;
  padding: 6px 14px;
  cursor: pointer;
}
.btn-cancel:hover { background: #2a2a2a; }

.btn-save {
  background: #14290a;
  border: none;
  border-radius: 6px;
  color: #4ade80;
  font-size: 12px;
  font-weight: 600;
  padding: 6px 14px;
  cursor: pointer;
}
.btn-save:hover:not(:disabled) { background: #1c3a10; }
.btn-save:disabled { opacity: 0.5; cursor: default; }
```

- [ ] **Step 3: Run tests — verify nothing broken**

```bash
cd frontend && npm test
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/SettingsPanel.tsx frontend/src/components/SettingsPanel.css
git commit -m "feat: SettingsPanel — retention, launch at login, shortcut display"
```

---

### Task 10: useClipboardItems hook

**Files:**
- Create: `frontend/src/hooks/useClipboardItems.ts`
- Create: `frontend/src/hooks/useClipboardItems.test.ts`

- [ ] **Step 1: Write failing tests — `frontend/src/hooks/useClipboardItems.test.ts`**

```typescript
import { renderHook, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { store } from '../../wailsjs/go/models'

vi.mock('../../wailsjs/go/main/App', () => ({
  GetItems: vi.fn(),
  CopyToClipboard: vi.fn(),
  PinItem: vi.fn(),
  DeleteItem: vi.fn(),
}))

vi.mock('../../wailsjs/runtime/runtime', () => ({
  EventsOn: vi.fn(() => vi.fn()),
}))

import { useClipboardItems } from './useClipboardItems'
import { GetItems } from '../../wailsjs/go/main/App'

function makeItem(id: string): store.Item {
  const item = new store.Item()
  item.id = id
  item.content = `content-${id}`
  item.contentHash = `hash-${id}`
  item.type = 'TEXT'
  item.subtype = ''
  item.pinned = false
  item.copiedAt = Date.now()
  item.createdAt = Date.now()
  item.charCount = 10
  item.imagePath = ''
  item.thumbBase64 = ''
  return item
}

describe('useClipboardItems', () => {
  beforeEach(() => { vi.clearAllMocks() })

  it('starts with empty items and loading=true', () => {
    vi.mocked(GetItems).mockReturnValue(new Promise(() => {}))
    const { result } = renderHook(() => useClipboardItems())
    expect(result.current.items).toHaveLength(0)
    expect(result.current.loading).toBe(true)
  })

  it('populates items after fetch resolves', async () => {
    vi.mocked(GetItems).mockResolvedValue([makeItem('1'), makeItem('2')])
    const { result } = renderHook(() => useClipboardItems())
    await waitFor(() => expect(result.current.loading).toBe(false))
    expect(result.current.items).toHaveLength(2)
  })
})
```

- [ ] **Step 2: Run — confirm fail**

```bash
cd frontend && npm test
```

Expected: `Cannot find module './useClipboardItems'`

- [ ] **Step 3: Implement `frontend/src/hooks/useClipboardItems.ts`**

```typescript
import { useState, useEffect, useCallback } from 'react'
import { store } from '../../wailsjs/go/models'
import { GetItems, CopyToClipboard, PinItem, DeleteItem } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export interface ClipboardState {
  items: store.Item[]
  loading: boolean
  copyItem: (id: string) => Promise<void>
  pinItem: (id: string, pinned: boolean) => Promise<void>
  deleteItem: (id: string) => Promise<void>
}

export function useClipboardItems(): ClipboardState {
  const [items, setItems] = useState<store.Item[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    GetItems(50, 0).then(fetched => {
      setItems(fetched ?? [])
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    const offNew = EventsOn('clipboard:new-item', (item: store.Item) => {
      setItems(prev => [item, ...prev.filter(i => i.contentHash !== item.contentHash)])
    })
    const offUpdate = EventsOn('clipboard:type-updated', (data: { id: string; type: string; subtype: string }) => {
      setItems(prev => prev.map(i => i.id === data.id ? { ...i, type: data.type, subtype: data.subtype } as store.Item : i))
    })
    return () => { offNew(); offUpdate() }
  }, [])

  const copyItem = useCallback((id: string) => CopyToClipboard(id), [])

  const pinItem = useCallback(async (id: string, pinned: boolean) => {
    await PinItem(id, pinned)
    setItems(prev => prev.map(i => i.id === id ? { ...i, pinned } as store.Item : i))
  }, [])

  const deleteItem = useCallback(async (id: string) => {
    await DeleteItem(id)
    setItems(prev => prev.filter(i => i.id !== id))
  }, [])

  return { items, loading, copyItem, pinItem, deleteItem }
}
```

- [ ] **Step 4: Run — confirm pass**

```bash
cd frontend && npm test
```

Expected: all tests pass

- [ ] **Step 5: Commit**

```bash
git add frontend/src/hooks/useClipboardItems.ts frontend/src/hooks/useClipboardItems.test.ts
git commit -m "feat: useClipboardItems hook — fetch, events, mutations"
```

---

### Task 11: Wire App.tsx

**Files:**
- Replace: `frontend/src/App.tsx`
- Replace: `frontend/src/App.css`
- Replace: `frontend/src/style.css`

- [ ] **Step 1: Replace `frontend/src/App.tsx`**

```tsx
import { useState, useEffect, useCallback, useRef } from 'react'
import './App.css'
import { SearchBar } from './components/SearchBar'
import { ItemList } from './components/ItemList'
import { ExpandTooltip } from './components/ExpandTooltip'
import { SettingsPanel } from './components/SettingsPanel'
import { useClipboardItems } from './hooks/useClipboardItems'
import { SearchItems } from '../wailsjs/go/main/App'
import { EventsOn } from '../wailsjs/runtime/runtime'
import { store } from '../wailsjs/go/models'

export default function App() {
  const { items, loading, copyItem, pinItem, deleteItem } = useClipboardItems()
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<store.Item[]>([])
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [copiedId, setCopiedId] = useState<string | null>(null)
  const [settingsOpen, setSettingsOpen] = useState(false)
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const displayItems = query ? results : items

  // Auto-select first item whenever list changes
  useEffect(() => {
    if (displayItems.length > 0 && !displayItems.find(i => i.id === selectedId)) {
      setSelectedId(displayItems[0].id)
    }
  }, [displayItems])

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

  // Reset on window show (hotkey re-opens panel)
  useEffect(() => {
    const off = EventsOn('wails:window-show', () => {
      setQuery('')
      setSelectedId(items[0]?.id ?? null)
      document.getElementById('search-input')?.focus()
    })
    return off
  }, [items])

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
        case e.key === 'ArrowDown':
          e.preventDefault()
          setSelectedId(displayItems[Math.min(idx + 1, displayItems.length - 1)]?.id ?? selectedId)
          break
        case e.key === 'ArrowUp':
          e.preventDefault()
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
  }, [displayItems, selectedId, settingsOpen, handleCopy, pinItem, deleteItem])

  return (
    <div className="app">
      <ExpandTooltip item={selectedItem} visible={!settingsOpen && !loading} />
      <div className="panel">
        <SearchBar value={query} onChange={setQuery} onEscape={() => {}} />
        <div className="panel-body">
          {loading
            ? <div className="panel-loading">Loading…</div>
            : <ItemList
                items={displayItems}
                selectedId={selectedId}
                copiedId={copiedId}
                onSelect={setSelectedId}
                onCopy={handleCopy}
                onPin={pinItem}
                onDelete={deleteItem}
              />
          }
        </div>
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
```

- [ ] **Step 2: Replace `frontend/src/App.css`**

```css
* { box-sizing: border-box; margin: 0; padding: 0; }

.app {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;   /* positioning context for ExpandTooltip */
  background: transparent;
}

.panel {
  width: 360px;
  height: 520px;
  background: #141414;
  border-radius: 14px;
  border: 1px solid #2a2a2a;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;   /* positioning context for SettingsPanel */
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.8);
}

.panel-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.panel-loading {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #555;
  font-size: 13px;
  font-family: -apple-system, sans-serif;
}

.panel-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 14px;
  border-top: 1px solid #222;
  background: #161616;
  flex-shrink: 0;
}

.footer-count { font-size: 11px; color: #444; font-family: -apple-system, sans-serif; }

.footer-hints {
  display: flex;
  gap: 10px;
  align-items: center;
  font-size: 11px;
  color: #555;
  font-family: -apple-system, sans-serif;
}

.footer-gear {
  background: none;
  border: none;
  color: #555;
  font-size: 13px;
  cursor: pointer;
  padding: 0 2px;
}
.footer-gear:hover { color: #888; }
```

- [ ] **Step 3: Replace `frontend/src/style.css`**

```css
:root { color-scheme: dark; }

html, body, #root {
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: transparent;
}

::-webkit-scrollbar { width: 4px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: #333; border-radius: 2px; }
```

- [ ] **Step 4: Run tests — verify all pass**

```bash
cd frontend && npm test
```

Expected: all tests pass

- [ ] **Step 5: Commit**

```bash
git add frontend/src/App.tsx frontend/src/App.css frontend/src/style.css
git commit -m "feat: wire App.tsx — search, list, tooltip, settings, keyboard nav"
```

---

### Task 12: Build and smoke test

- [ ] **Step 1: TypeScript build check**

```bash
cd /Users/arsalaanabbas.zaidi/Projects/clipboard-manager/frontend && npm run build
```

Expected: `dist/` generated with no TS errors. Fix any errors before continuing — common issues:
- `store.Item` spread — cast result as `store.Item` e.g. `{ ...i, pinned } as store.Item`
- Import path off by one `../` — verify `wailsjs/` is reachable from each file's location

- [ ] **Step 2: Wails dev server**

```bash
cd /Users/arsalaanabbas.zaidi/Projects/clipboard-manager && wails dev
```

Expected: app compiles and window is hidden (StartHidden=true). Trigger with `option+space`. Verify:
- Items populate from DB
- Search filters in real-time
- `↑↓` navigates, highlighted row changes
- `↵` flashes "✓ Copied", item written to NSPasteboard
- Expand tooltip appears to the left of the panel
- `⌘,` opens settings
- `Esc` dismisses settings / clears search

- [ ] **Step 3: Final commit**

```bash
cd /Users/arsalaanabbas.zaidi/Projects/clipboard-manager
git add -A
git commit -m "feat: clipboard manager frontend v1 complete"
```

---

## Spec Coverage

| Spec requirement | Task |
|---|---|
| Search bar (FTS5, ⌘K focus) | 5, 11 |
| Pinned + Recent sections | 7 |
| 23 type badges, colour-coded | 2, 4 |
| Image thumbnails inline | 6 |
| Keyboard nav + Enter-to-copy | 11 |
| Pin / delete on hover | 6 |
| Live updates via Wails events | 10 |
| Settings panel (⌘,) | 9, 11 |
| Expand tooltip full content | 8 |
| Auto-focus search on window show | 11 |
| ✓ Copied flash (1.2s) | 11 |
