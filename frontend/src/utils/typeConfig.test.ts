import { describe, it, expect } from 'vitest'
import { getTypeConfig } from './typeConfig'

describe('getTypeConfig', () => {
  it('returns green/mono config for CODE with subtype as label', () => {
    const cfg = getTypeConfig('CODE', 'go')
    expect(cfg.label).toBe('GO')
    expect(cfg.color).toBe('#4ade80')
    expect(cfg.mono).toBe(true)
  })

  it('returns "CODE" label when CODE subtype is empty', () => {
    expect(getTypeConfig('CODE', '').label).toBe('CODE')
  })

  it('returns blue config for URL', () => {
    const cfg = getTypeConfig('URL')
    expect(cfg.color).toBe('#60a5fa')
    expect(cfg.mono).toBe(false)
  })

  it('returns amber config for API_KEY', () => {
    expect(getTypeConfig('API_KEY').color).toBe('#fbbf24')
  })

  it('returns grey/TEXT for unknown type', () => {
    const cfg = getTypeConfig('UNKNOWN')
    expect(cfg.label).toBe('TEXT')
    expect(cfg.color).toBe('#888')
  })

  it('returns mono=true for JSON', () => {
    expect(getTypeConfig('JSON').mono).toBe(true)
  })

  it('returns mono=false for EMAIL', () => {
    expect(getTypeConfig('EMAIL').mono).toBe(false)
  })
})
