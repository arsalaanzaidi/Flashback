interface ErrorToastProps {
  message: string | null
  onDismiss: () => void
}

export function ErrorToast({ message, onDismiss }: ErrorToastProps) {
  return (
    <div className={`error-toast${message ? ' error-toast--visible' : ''}`} role="alert" aria-hidden={!message}>
      <span className="error-toast__msg">{message}</span>
      <button type="button" className="error-toast__dismiss" onClick={onDismiss} aria-label="Dismiss">×</button>
    </div>
  )
}
