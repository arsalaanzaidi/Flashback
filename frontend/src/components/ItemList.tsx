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
