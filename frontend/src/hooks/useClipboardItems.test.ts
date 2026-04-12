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
