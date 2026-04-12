import { useState, useEffect, useCallback, useRef } from 'react'
import { store } from '../../wailsjs/go/models'
import { GetItems, CopyToClipboard, PinItem, DeleteItem } from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'

export interface ClipboardState {
  items: store.Item[]
  loading: boolean
  copyItem: (id: string) => Promise<void>
  pinItem: (id: string, pinned: boolean) => Promise<void>
  deleteItem: (id: string) => Promise<void>
  newItemIds: Set<string>
  deletingIds: Set<string>
}

export function useClipboardItems(): ClipboardState {
  const [items, setItems] = useState<store.Item[]>([])
  const [loading, setLoading] = useState(true)
  const [newItemIds, setNewItemIds] = useState<Set<string>>(new Set())
  const [deletingIds, setDeletingIds] = useState<Set<string>>(new Set())
  const deletingInFlight = useRef(new Set<string>())

  useEffect(() => {
    GetItems(50, 0).then(fetched => {
      setItems(fetched ?? [])
      setLoading(false)
    })
  }, [])

  useEffect(() => {
    const timeouts = new Set<ReturnType<typeof setTimeout>>()

    const offNew = EventsOn('clipboard:new-item', (item: store.Item) => {
      setItems(prev => [item, ...prev.filter(i => i.contentHash !== item.contentHash)])
      setNewItemIds(prev => new Set(prev).add(item.id))
      const t = setTimeout(() => {
        timeouts.delete(t)
        setNewItemIds(prev => {
          const next = new Set(prev)
          next.delete(item.id)
          return next
        })
      }, 600)
      timeouts.add(t)
    })

    const offUpdate = EventsOn('clipboard:type-updated', (data: { id: string; type: string; subtype: string }) => {
      setItems(prev => prev.map(i => i.id === data.id ? { ...i, type: data.type, subtype: data.subtype } as store.Item : i))
    })

    return () => {
      offNew()
      offUpdate()
      timeouts.forEach(clearTimeout)
    }
  }, [])

  const copyItem = useCallback((id: string) => CopyToClipboard(id), [])

  const pinItem = useCallback(async (id: string, pinned: boolean) => {
    await PinItem(id, pinned)
    setItems(prev => prev.map(i => i.id === id ? { ...i, pinned } as store.Item : i))
  }, [])

  const deleteItem = useCallback(async (id: string) => {
    if (deletingInFlight.current.has(id)) return
    deletingInFlight.current.add(id)
    setDeletingIds(prev => new Set(prev).add(id))
    await new Promise<void>(resolve => setTimeout(resolve, 200))
    try {
      await DeleteItem(id)
      setItems(prev => prev.filter(i => i.id !== id))
    } finally {
      deletingInFlight.current.delete(id)
      setDeletingIds(prev => {
        const next = new Set(prev)
        next.delete(id)
        return next
      })
    }
  }, [])

  return { items, loading, copyItem, pinItem, deleteItem, newItemIds, deletingIds }
}
