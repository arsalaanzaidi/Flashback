import { describe, it, expect } from 'vitest'
import { getTypeConfig } from './typeConfig'

describe('getTypeConfig', () => {
  it('returns blue config for URL', () => {
    const cfg = getTypeConfig('URL')
    expect(cfg.color).toBe('#60a5fa')
    expect(cfg.mono).toBe(false)
  })

  it('returns mono=true for COLOR_CODE', () => {
    expect(getTypeConfig('COLOR_CODE').mono).toBe(true)
  })

  it('returns mono=false for EMAIL', () => {
    expect(getTypeConfig('EMAIL').mono).toBe(false)
  })

  it('returns purple family for IMAGE', () => {
    expect(getTypeConfig('IMAGE').color).toBe('#c084fc')
  })

  it('falls back to TEXT for unknown / dropped types', () => {
    for (const t of ['UNKNOWN', 'JSON', 'CODE', 'JWT', 'YAML', 'API_KEY']) {
      const cfg = getTypeConfig(t)
      expect(cfg.label).toBe('TEXT')
      expect(cfg.color).toBe('#888')
    }
  })
})
