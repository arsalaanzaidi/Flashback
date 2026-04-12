import { FC } from 'react'
import './ExpandTooltip.css'
import { store } from '../../wailsjs/go/models'
import { getTypeConfig } from '../utils/typeConfig'
import { formatAge } from '../utils/format'

interface ExpandTooltipProps {
  item: store.Item | null
  visible: boolean
}

export const ExpandTooltip: FC<ExpandTooltipProps> = ({ item, visible }) => {
  if (!item || !visible) return null

  const cfg = getTypeConfig(item.type, item.subtype)

  function renderBody() {
    if (item.type === 'IMAGE') {
      return (
        <div className="tooltip-image-wrap">
          {item.thumbBase64 ? (
            <img src={item.thumbBase64} alt="" className="tooltip-image" />
          ) : (
            <div className="tooltip-image-placeholder">
              {item.imagePath.split('/').pop()}
            </div>
          )}
        </div>
      )
    }

    if (item.type === 'COLOR' || item.type === 'COLOR_CODE') {
      return (
        <div className="tooltip-color-wrap">
          <div className="color-swatch" style={{ background: item.content }} />
          <div className="color-table">
            <div className="color-row">
              <span className="color-key">HEX</span>
              <span className="color-val">{item.content}</span>
            </div>
          </div>
        </div>
      )
    }

    return (
      <pre className={`tooltip-text ${cfg.mono ? 'mono' : ''}`}>
        {item.content}
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
