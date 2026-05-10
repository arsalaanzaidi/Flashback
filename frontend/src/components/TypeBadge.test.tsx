import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import { TypeBadge } from './TypeBadge'

describe('TypeBadge', () => {
  it('renders URL label', () => {
    render(<TypeBadge type="URL" />)
    expect(screen.getByText(/URL/)).toBeInTheDocument()
  })
  it('renders EMAIL label', () => {
    render(<TypeBadge type="EMAIL" />)
    expect(screen.getByText(/EMAIL/)).toBeInTheDocument()
  })
  it('falls back to TEXT for unknown / dropped type', () => {
    render(<TypeBadge type="UNKNOWN" />)
    expect(screen.getByText(/TEXT/)).toBeInTheDocument()
  })
})
