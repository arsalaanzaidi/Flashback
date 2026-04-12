import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { TypeBadge } from './TypeBadge'

describe('TypeBadge', () => {
  it('renders URL label', () => {
    render(<TypeBadge type="URL" />)
    expect(screen.getByText(/URL/)).toBeInTheDocument()
  })
  it('renders CODE subtype as uppercase label', () => {
    render(<TypeBadge type="CODE" subtype="python" />)
    expect(screen.getByText(/PYTHON/)).toBeInTheDocument()
  })
  it('falls back to TEXT for unknown type', () => {
    render(<TypeBadge type="UNKNOWN" />)
    expect(screen.getByText(/TEXT/)).toBeInTheDocument()
  })
})
