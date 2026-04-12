import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { ClipboardItem } from './ClipboardItem'
import { store } from '../../wailsjs/go/models'

function makeItem(overrides: Partial<store.Item> = {}): store.Item {
  const item = new store.Item()
  item.id = 'test-id'
  item.content = 'https://example.com'
  item.type = 'URL'
  item.subtype = ''
  item.pinned = false
  item.copiedAt = Date.now() - 60_000
  item.createdAt = Date.now() - 60_000
  item.charCount = 19
  item.imagePath = ''
  item.thumbBase64 = ''
  return Object.assign(item, overrides)
}

const defaultProps = {
  selected: false,
  justCopied: false,
  isNew: false,
  isDeleting: false,
  staggerIndex: -1,
  onSelect: vi.fn(),
  onCopy: vi.fn(),
  onPin: vi.fn(),
  onDelete: vi.fn(),
}

describe('ClipboardItem', () => {
  it('renders the type badge', () => {
    render(<ClipboardItem item={makeItem()} {...defaultProps} />)
    expect(screen.getByText(/URL/)).toBeInTheDocument()
  })

  it('renders content preview', () => {
    render(<ClipboardItem item={makeItem({ content: 'https://example.com' })} {...defaultProps} />)
    expect(screen.getByText('https://example.com')).toBeInTheDocument()
  })

  it('shows "✓ Copied" when justCopied and selected', () => {
    render(<ClipboardItem item={makeItem()} {...defaultProps} selected={true} justCopied={true} />)
    expect(screen.getByText(/Copied/)).toBeInTheDocument()
  })

  it('calls onCopy when clicked', async () => {
    const onCopy = vi.fn()
    render(<ClipboardItem item={makeItem()} {...defaultProps} onCopy={onCopy} />)
    await userEvent.click(screen.getByRole('listitem'))
    expect(onCopy).toHaveBeenCalled()
  })

  it('calls onPin when pin button clicked on selected item', async () => {
    const onPin = vi.fn()
    render(<ClipboardItem item={makeItem()} {...defaultProps} selected={true} onPin={onPin} />)
    await userEvent.click(screen.getByTitle(/pin/i))
    expect(onPin).toHaveBeenCalled()
  })

  it('applies is-new class when isNew is true and staggerIndex is -1', () => {
    const { container } = render(
      <ClipboardItem item={makeItem()} {...defaultProps} isNew={true} staggerIndex={-1} />
    )
    expect(container.querySelector('.is-new')).toBeInTheDocument()
  })

  it('does NOT apply is-new class when isNew is true but staggerIndex >= 0', () => {
    const { container } = render(
      <ClipboardItem item={makeItem()} {...defaultProps} isNew={true} staggerIndex={0} />
    )
    expect(container.querySelector('.is-new')).not.toBeInTheDocument()
  })

  it('applies is-deleting class when isDeleting is true', () => {
    const { container } = render(
      <ClipboardItem item={makeItem()} {...defaultProps} isDeleting={true} />
    )
    expect(container.querySelector('.is-deleting')).toBeInTheDocument()
  })

  it('applies is-copied class when justCopied is true', () => {
    const { container } = render(
      <ClipboardItem item={makeItem()} {...defaultProps} justCopied={true} />
    )
    expect(container.querySelector('.is-copied')).toBeInTheDocument()
  })
})
