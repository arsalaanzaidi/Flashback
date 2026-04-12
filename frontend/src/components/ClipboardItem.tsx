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
  onSelect: () => void
  onCopy: () => void
  onPin: () => void
  onDelete: () => void
}

export const ClipboardItem: FC<ClipboardItemProps> = ({
  item,
  selected,
  justCopied,
  onSelect,
  onCopy,
  onPin,
  onDelete,
}) => {
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
      ) : item.type === 'COLOR' || item.type === 'COLOR_CODE' ? (
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
