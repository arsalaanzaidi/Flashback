import { describe, it, expect } from 'vitest'
import { getTypeConfig } from './typeConfig'

describe('getTypeConfig', () => {
  it('returns URL label', () => {
    expect(getTypeConfig('URL').label).toBe('URL')
  })

  it('returns mono=true for COLOR_CODE', () => {
    expect(getTypeConfig('COLOR_CODE').mono).toBe(true)
  })

  it('returns mono=false for EMAIL', () => {
    expect(getTypeConfig('EMAIL').mono).toBe(false)
  })

  it('returns IMAGE label', () => {
    expect(getTypeConfig('IMAGE').label).toBe('IMAGE')
  })

  it('falls back to TEXT for unknown / dropped types', () => {
    for (const t of ['UNKNOWN', 'JSON', 'CODE', 'JWT', 'YAML', 'API_KEY']) {
      expect(getTypeConfig(t).label).toBe('TEXT')
    }
  })
})
