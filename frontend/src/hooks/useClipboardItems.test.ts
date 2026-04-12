import { renderHook, waitFor, act } from '@testing-library/react'
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

  it('adds id to deletingIds immediately and removes from items after 200ms', async () => {
    const { DeleteItem } = await import('../../wailsjs/go/main/App')
    vi.mocked(DeleteItem).mockResolvedValue(undefined)
    vi.mocked(GetItems).mockResolvedValue([makeItem('a'), makeItem('b')])

    const { result } = renderHook(() => useClipboardItems())
    await waitFor(() => expect(result.current.items).toHaveLength(2))

    vi.useFakeTimers()
    try {
      act(() => { result.current.deleteItem('a') })

      // immediately: id is in deletingIds, still in items
      expect(result.current.deletingIds.has('a')).toBe(true)
      expect(result.current.items.find(i => i.id === 'a')).toBeDefined()
      expect(vi.mocked(DeleteItem)).not.toHaveBeenCalled()

      // after 200ms: item removed from list
      await act(async () => { await vi.runAllTimersAsync() })
      expect(vi.mocked(DeleteItem)).toHaveBeenCalledWith('a')
      expect(result.current.items.find(i => i.id === 'a')).toBeUndefined()
      expect(result.current.deletingIds.has('a')).toBe(false)
    } finally {
      vi.useRealTimers()
    }
  })

  it('adds new item id to newItemIds when clipboard:new-item fires and removes it after 600ms', async () => {
    const { EventsOn } = await import('../../wailsjs/runtime/runtime')

    let newItemCallback: ((item: store.Item) => void) | null = null
    vi.mocked(EventsOn).mockImplementation((event, cb) => {
      if (event === 'clipboard:new-item') newItemCallback = cb as (item: store.Item) => void
      return vi.fn()
    })

    vi.mocked(GetItems).mockResolvedValue([])

    const { result } = renderHook(() => useClipboardItems())
    await waitFor(() => expect(result.current.loading).toBe(false))

    vi.useFakeTimers()
    try {
      const newItem = makeItem('new-1')
      act(() => { newItemCallback!(newItem) })

      expect(result.current.newItemIds.has('new-1')).toBe(true)

      act(() => { vi.advanceTimersByTime(600) })
      expect(result.current.newItemIds.has('new-1')).toBe(false)
    } finally {
      vi.useRealTimers()
    }
  })
})
