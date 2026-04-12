import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useState } from 'react'
import { SearchBar } from './SearchBar'

function SearchBarWrapper() {
  const [value, setValue] = useState('')
  const onEscape = vi.fn()
  return <SearchBar value={value} onChange={setValue} onEscape={onEscape} />
}

describe('SearchBar', () => {
  it('renders placeholder text', () => {
    render(<SearchBar value="" onChange={vi.fn()} onEscape={vi.fn()} />)
    expect(screen.getByPlaceholderText('Search clipboard history...')).toBeInTheDocument()
  })

  it('calls onChange when user types', async () => {
    const onChange = vi.fn()
    render(<SearchBar value="" onChange={onChange} onEscape={vi.fn()} />)
    await userEvent.type(screen.getByRole('textbox'), 'h')
    expect(onChange).toHaveBeenCalledWith('h')
  })

  it('calls onEscape when Escape pressed on empty input', async () => {
    const onEscape = vi.fn()
    render(<SearchBar value="" onChange={vi.fn()} onEscape={onEscape} />)
    await userEvent.keyboard('{Escape}')
    expect(onEscape).toHaveBeenCalled()
  })

  it('clears input (calls onChange("")) on Escape when non-empty, does not call onEscape', async () => {
    const onChange = vi.fn()
    const onEscape = vi.fn()
    render(<SearchBar value="hello" onChange={onChange} onEscape={onEscape} />)
    await userEvent.keyboard('{Escape}')
    expect(onChange).toHaveBeenCalledWith('')
    expect(onEscape).not.toHaveBeenCalled()
  })
})
