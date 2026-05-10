export interface TypeConfig {
  label: string
  icon: string
  color: string
  bg: string
  mono: boolean
}

export function getTypeConfig(type: string, _subtype = ''): TypeConfig {
  switch (type) {
    case 'URL':       return { label: 'URL',     icon: '🔗', color: '#60a5fa', bg: '#0f1e3a', mono: false }
    case 'EMAIL':     return { label: 'EMAIL',   icon: '📧', color: '#f87171', bg: '#2a0a0a', mono: false }
    case 'IP':        return { label: 'IP',      icon: '🌐', color: '#f87171', bg: '#2a0a0a', mono: false }
    case 'IMAGE':     return { label: 'IMAGE',   icon: '🖼', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'PDF':       return { label: 'PDF',     icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'RICH_TEXT': return { label: 'RTF',     icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'HTML':      return { label: 'HTML',    icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: true  }
    case 'COLOR':     return { label: 'COLOR',   icon: '🎨', color: '#facc15', bg: '#1a1a0a', mono: true  }
    case 'COLOR_CODE':return { label: 'COLOR',   icon: '🎨', color: '#facc15', bg: '#1a1a0a', mono: true  }
    case 'FILE_REF':  return { label: 'FILE',    icon: '📁', color: '#34d399', bg: '#0a2a1a', mono: false }
    case 'FILE_PATH': return { label: 'PATH',    icon: '📁', color: '#34d399', bg: '#0a2a1a', mono: true  }
    case 'UUID':      return { label: 'UUID',    icon: '🔢', color: '#888',    bg: '#222',    mono: true  }
    default:          return { label: 'TEXT',    icon: '📝', color: '#888',    bg: '#222',    mono: false }
  }
}
