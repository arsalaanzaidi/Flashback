import { FC } from 'react'
import { getTypeConfig } from '../utils/typeConfig'
import './TypeBadge.css'

export const TypeBadge: FC<{ type: string; subtype?: string }> = ({ type, subtype = '' }) => {
  const cfg = getTypeConfig(type, subtype)
  return <span className="type-badge">{cfg.label}</span>
}
