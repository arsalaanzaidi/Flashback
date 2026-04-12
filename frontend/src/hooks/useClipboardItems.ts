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
