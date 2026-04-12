export interface TypeConfig {
  label: string
  icon: string
  color: string
  bg: string
  mono: boolean
}

export function getTypeConfig(type: string, subtype = ''): TypeConfig {
  switch (type) {
    case 'URL':       return { label: 'URL',     icon: '🔗', color: '#60a5fa', bg: '#0f1e3a', mono: false }
    case 'EMAIL':     return { label: 'EMAIL',   icon: '📧', color: '#f87171', bg: '#2a0a0a', mono: false }
    case 'PHONE':     return { label: 'PHONE',   icon: '📱', color: '#f87171', bg: '#2a0a0a', mono: false }
    case 'IP':        return { label: 'IP',      icon: '🌐', color: '#f87171', bg: '#2a0a0a', mono: false }
    case 'IMAGE':     return { label: 'IMAGE',   icon: '🖼', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'PDF':       return { label: 'PDF',     icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'RICH_TEXT': return { label: 'RTF',     icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: false }
    case 'HTML':      return { label: 'HTML',    icon: '📄', color: '#c084fc', bg: '#2a1a4a', mono: true  }
    case 'COLOR':     return { label: 'COLOR',   icon: '🎨', color: '#facc15', bg: '#1a1a0a', mono: true  }
    case 'COLOR_CODE':return { label: 'COLOR',   icon: '🎨', color: '#facc15', bg: '#1a1a0a', mono: true  }
    case 'FILE_REF':  return { label: 'FILE',    icon: '📁', color: '#34d399', bg: '#0a2a1a', mono: false }
    case 'FILE_PATH': return { label: 'PATH',    icon: '📁', color: '#34d399', bg: '#0a2a1a', mono: true  }
    case 'HASH':      return { label: 'HASH',    icon: '#',  color: '#fbbf24', bg: '#3a2f0d', mono: true  }
    case 'JWT':       return { label: 'JWT',     icon: '🔑', color: '#fbbf24', bg: '#3a2f0d', mono: true  }
    case 'SSH_KEY':   return { label: 'SSH KEY', icon: '🔏', color: '#fbbf24', bg: '#3a2f0d', mono: true  }
    case 'API_KEY':   return { label: 'API KEY', icon: '🔑', color: '#fbbf24', bg: '#3a2f0d', mono: true  }
    case 'JSON':      return { label: 'JSON',    icon: '{}', color: '#4ade80', bg: '#0f2a1a', mono: true  }
    case 'XML':       return { label: 'XML',     icon: '<>', color: '#4ade80', bg: '#0f2a1a', mono: true  }
    case 'YAML':      return { label: 'YAML',    icon: '—',  color: '#4ade80', bg: '#0f2a1a', mono: true  }
    case 'SQL':       return { label: 'SQL',     icon: '🗃', color: '#4ade80', bg: '#0f2a1a', mono: true  }
    case 'CODE':      return {
      label: subtype ? subtype.toUpperCase() : 'CODE',
      icon: '💻', color: '#4ade80', bg: '#14290a', mono: true,
    }
    case 'UUID':      return { label: 'UUID',    icon: '🔢', color: '#888',    bg: '#222',    mono: true  }
    case 'BASE64':    return { label: 'B64',     icon: '📦', color: '#888',    bg: '#222',    mono: true  }
    case 'MARKDOWN':  return { label: 'MD',      icon: '📝', color: '#888',    bg: '#222',    mono: false }
    default:          return { label: 'TEXT',    icon: '📝', color: '#888',    bg: '#222',    mono: false }
  }
}
