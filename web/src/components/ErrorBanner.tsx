import type { ReactNode } from 'react'
import { ApiError } from '../api/client'

type Props = {
  error: unknown
  onRetry?: () => void
  settingsHref?: string
}

export function ErrorBanner({ error, onRetry, settingsHref = '/settings' }: Props) {
  if (!error) return null

  let code = 'ERROR'
  let message = 'Something went wrong'
  if (error instanceof ApiError) {
    code = error.code
    message = error.message
  } else if (error instanceof Error) {
    message = error.message
  }

  const unauthorized = code === 'UNAUTHORIZED'

  return (
    <div className="banner banner--error" role="alert">
      <strong>{code}</strong>
      <div>{message}</div>
      <div className="row" style={{ marginTop: '0.5rem' }}>
        {unauthorized ? (
          <a className="btn btn--primary" href={settingsHref}>
            Open Settings (token)
          </a>
        ) : null}
        {onRetry ? (
          <button type="button" className="btn" onClick={onRetry}>
            Retry
          </button>
        ) : null}
      </div>
    </div>
  )
}

export function EmptyState({
  title,
  children,
}: {
  title: string
  children?: ReactNode
}) {
  return (
    <div className="empty" role="status">
      <div>{title}</div>
      {children ? <div style={{ marginTop: '0.5rem' }}>{children}</div> : null}
    </div>
  )
}
