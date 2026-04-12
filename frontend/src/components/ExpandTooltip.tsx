import { FC, useState, useEffect } from 'react'
import './ExpandTooltip.css'
import { store } from '../../wailsjs/go/models'
import { getTypeConfig } from '../utils/typeConfig'
import { formatAge } from '../utils/format'
import { GetImageBase64 } from '../../wailsjs/go/main/App'

interface ExpandTooltipProps {
  item: store.Item | null
  visible: boolean
}

const MAX_W = 320
const MAX_H = 420

function computeDispDims(w: number, h: number) {
  const scale = Math.min(MAX_W / w, MAX_H / h, 1)
  return { w: Math.round(w * scale), h: Math.round(h * scale) }
}

export const ExpandTooltip: FC<ExpandTooltipProps> = ({ item, visible }) => {
  const [imgSrc, setImgSrc] = useState<string | null>(null)
  const [naturalDims, setNaturalDims] = useState<{ w: number; h: number } | null>(null)

  useEffect(() => {
    if (!visible || !item || item.type !== 'IMAGE') {
      setImgSrc(null)
      setNaturalDims(null)
      return
    }
    setNaturalDims(null)
    let cancelled = false
    GetImageBase64(item.imagePath)
      .then(src => { if (!cancelled) setImgSrc(src) })
      .catch(() => { if (!cancelled) setImgSrc(null) })
    return () => { cancelled = true }
  }, [visible, item])

  if (!item || !visible) return null

  const cfg = getTypeConfig(item.type, item.subtype)
  const dispDims = naturalDims ? computeDispDims(naturalDims.w, naturalDims.h) : null

  const tooltipStyle = dispDims ? { width: dispDims.w } : undefined
  const bodyStyle = item.type === 'IMAGE'
    ? { padding: 0, maxHeight: 'none' as const }
    : undefined

  function renderBody(it: store.Item) {
    if (it.type === 'IMAGE') {
      const src = imgSrc ?? it.thumbBase64
      return src ? (
        <img
          src={src}
          alt=""
          style={{
            display: 'block',
            width: dispDims ? dispDims.w : 260,
            height: dispDims ? dispDims.h : 'auto',
          }}
          onLoad={e => {
            if (src === imgSrc) {
              const img = e.currentTarget
              if (img.naturalWidth > 0 && img.naturalHeight > 0) {
                setNaturalDims({ w: img.naturalWidth, h: img.naturalHeight })
              }
            }
          }}
        />
      ) : (
        <div className="tooltip-image-placeholder">
          {it.imagePath.split('/').pop()}
        </div>
      )
    }

    if (it.type === 'COLOR' || it.type === 'COLOR_CODE') {
      return (
        <div className="tooltip-color-wrap">
          <div className="color-swatch" style={{ background: it.content }} />
          <div className="color-table">
            <div className="color-row">
              <span className="color-key">HEX</span>
              <span className="color-val">{it.content}</span>
            </div>
          </div>
        </div>
      )
    }

    return (
      <pre className={`tooltip-text ${cfg.mono ? 'mono' : ''}`}>
        {it.content}
      </pre>
    )
  }

  return (
    <div className="expand-tooltip" style={tooltipStyle}>
      <div className="tooltip-header">
        <span className="tooltip-badge" style={{ color: cfg.color, background: cfg.bg }}>
          {cfg.icon} {cfg.label}
        </span>
        <span className="tooltip-meta">full content</span>
      </div>
      <div className="tooltip-body" style={bodyStyle}>{renderBody(item)}</div>
      <div className="tooltip-footer">
        <span>{item.charCount > 0 ? `${item.charCount} chars` : ''}</span>
        <span>{formatAge(item.copiedAt)}</span>
      </div>
    </div>
  )
}
