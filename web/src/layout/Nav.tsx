import { NavLink } from 'react-router-dom'

const LINKS = [
  { to: '/', label: 'Explore', end: true },
  { to: '/overview', label: 'Overview' },
  { to: '/tasks', label: 'Tasks' },
  { to: '/loop', label: 'Loop' },
  { to: '/discoveries', label: 'Discoveries' },
  { to: '/reviews', label: 'Reviews' },
  { to: '/seed', label: 'Seed' },
  { to: '/settings', label: 'Settings' },
] as const

type Props = {
  storeReady: boolean
  className?: string
}

export function Nav({ storeReady, className = 'nav' }: Props) {
  return (
    <nav className={className} aria-label="Primary">
      {LINKS.map((link) => {
        const needsStore = link.to !== '/settings'
        const disabled = needsStore && !storeReady
        return (
          <NavLink
            key={link.to}
            to={link.to}
            end={'end' in link ? link.end : false}
            aria-disabled={disabled || undefined}
            tabIndex={disabled ? -1 : undefined}
            onClick={(e) => {
              if (disabled) e.preventDefault()
            }}
          >
            {link.label}
          </NavLink>
        )
      })}
    </nav>
  )
}
