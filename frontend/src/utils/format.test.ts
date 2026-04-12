import { describe, it, expect, beforeAll, afterAll, vi } from 'vitest'
import { formatAge, getPreview } from './format'

describe('formatAge', () => {
  const now = new Date('2026-04-12T12:00:00.000Z').getTime()

  beforeAll(() => { vi.setSystemTime(now) })
  afterAll(() => { vi.useRealTimers() })

  it('returns "just now" for < 60s', () => {
    expect(formatAge(now - 30_000)).toBe('just now')
  })
  it('returns minutes for < 1h', () => {
    expect(formatAge(now - 5 * 60_000)).toBe('5m')
  })
  it('returns hours for < 1d', () => {
    expect(formatAge(now - 3 * 3_600_000)).toBe('3h')
  })
  it('returns days for < 30d', () => {
    expect(formatAge(now - 2 * 86_400_000)).toBe('2d')
  })
  it('returns months for >= 30d', () => {
    expect(formatAge(now - 35 * 86_400_000)).toBe('1mo')
  })
})

describe('getPreview', () => {
  it('returns short content unchanged', () => {
    expect(getPreview('hello world')).toBe('hello world')
  })
  it('truncates at 80 chars with ellipsis', () => {
    const result = getPreview('a'.repeat(100))
    expect(result).toHaveLength(81)
    expect(result.endsWith('…')).toBe(true)
  })
  it('collapses internal whitespace', () => {
    expect(getPreview('hello   \n   world')).toBe('hello world')
  })
})
