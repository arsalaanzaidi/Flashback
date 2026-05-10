import { FC } from 'react'
import './ClipboardItem.css'
import { store } from '../../wailsjs/go/models'
import { TypeBadge } from './TypeBadge'
import { formatAge, getPreview } from '../utils/format'
import { getTypeConfig } from '../utils/typeConfig'

interface ClipboardItemProps {
  item: store.Item
  selected: boolean
  justCopied: boolean
  isNew: boolean
  isDeleting: boolean
  staggerIndex: number  // -1 = no stagger, 0-7 = stagger delay
  onSelect: () => void
  onCopy: () => void
  onPin: () => void
  onDelete: () => void
}

export const ClipboardItem: FC<ClipboardItemProps> = ({
  item,
  selected,
  justCopied,
  isNew,
  isDeleting,
  staggerIndex,
  onSelect,
  onCopy,
  onPin,
  onDelete,
}) => {
  const cfg = getTypeConfig(item.type, item.subtype)

  // Stagger overrides isNew — both use fly-in, but stagger has a delay
  const staggerStyle = staggerIndex >= 0 && staggerIndex < 8
    ? { animation: `fly-in 180ms cubic-bezier(0.34,1.56,0.64,1) ${staggerIndex * 40}ms both` }
    : undefined

  const className = [
    'clipboard-item',
    selected ? 'selected' : '',
    isNew && staggerIndex < 0 ? 'is-new' : '',
    isDeleting ? 'is-deleting' : '',
    justCopied ? 'is-copied' : '',
  ].filter(Boolean).join(' ')

  return (
    <li
      role="listitem"
      className={className}
      style={staggerStyle}
      onMouseEnter={onSelect}
      onClick={onCopy}
    >
      <TypeBadge type={item.type} subtype={item.subtype} />

      {item.type === 'IMAGE' && item.thumbBase64 ? (
        <img className="item-thumb" src={item.thumbBase64} alt="" />
      ) : item.type === 'COLOR' || item.type === 'COLOR_CODE' ? (
        <span className="item-swatch" style={{ background: item.content }} />
      ) : null}

      <span className={`item-preview ${cfg.mono ? 'mono' : ''}`}>
        {getPreview(item.content)}
      </span>

      {selected ? (
        <div className="item-actions">
          {justCopied && (
            <span className="item-copied">
              ✓ Copied
            </span>
          )}
          <button
            className="action-btn"
            title="Pin"
            onClick={(e) => {
              e.stopPropagation()
              onPin()
            }}
          >
            📌
          </button>
          <button
            className="action-delete action-btn"
            title="Delete"
            onClick={(e) => {
              e.stopPropagation()
              onDelete()
            }}
          >
            ✕
          </button>
        </div>
      ) : (
        <span className="item-age">{formatAge(item.copiedAt)}</span>
      )}
    </li>
  )
}
