export interface TypeConfig {
  label: string
  mono: boolean
}

export function getTypeConfig(type: string, _subtype = ''): TypeConfig {
  switch (type) {
    case 'URL':       return { label: 'URL',   mono: false }
    case 'EMAIL':     return { label: 'EMAIL', mono: false }
    case 'IP':        return { label: 'IP',    mono: false }
    case 'IMAGE':     return { label: 'IMAGE', mono: false }
    case 'PDF':       return { label: 'PDF',   mono: false }
    case 'RICH_TEXT': return { label: 'RTF',   mono: false }
    case 'HTML':      return { label: 'HTML',  mono: true  }
    case 'COLOR':     return { label: 'COLOR', mono: true  }
    case 'COLOR_CODE':return { label: 'COLOR', mono: true  }
    case 'FILE_REF':  return { label: 'FILE',  mono: false }
    case 'FILE_PATH': return { label: 'PATH',  mono: true  }
    case 'UUID':      return { label: 'UUID',  mono: true  }
    default:          return { label: 'TEXT',  mono: false }
  }
}
