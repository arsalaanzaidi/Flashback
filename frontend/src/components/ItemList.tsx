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
