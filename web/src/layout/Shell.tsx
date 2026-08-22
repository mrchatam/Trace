import { Link, Outlet, useLocation } from 'react-router-dom'
import { useAppChrome } from '../context/AppChrome'
import { ProjectGate } from '../screens/ProjectGate'
import { Nav } from './Nav'

export function Shell() {
  const { project, health, version, storeReady, loading, chromeError } = useAppChrome()
  const location = useLocation()

  const blocked = !loading && (!storeReady || !!chromeError)
  const allowSettings = location.pathname === '/settings'

  return (
    <div className="app-shell">
      <header className="chrome">
        <Link to="/" className="chrome__brand">
          Trace
        </Link>
        <span className="chrome__project" title={project?.root ?? ''}>
          {project?.root ?? (loading ? '…' : 'no project')}
        </span>
        <div className="chrome__strip" aria-live="polite">
          {health?.ok ? (
            <span className="pill pill--ok">health ok</span>
          ) : (
            <span className="pill pill--bad">health down</span>
          )}
          {version ? (
            <span className="pill" title={`api ${version.api_version}`}>
              {version.trace_version || version.api_version}
            </span>
          ) : null}
          <span className="pill">{storeReady ? 'store ready' : 'store blocked'}</span>
        </div>
      </header>

      {blocked && !allowSettings ? (
        <ProjectGate />
      ) : (
        <div className="app-body">
          <Nav storeReady={storeReady} />
          <main className="main">
            <Outlet />
          </main>
          <Nav storeReady={storeReady} className="bottom-nav" />
        </div>
      )}
    </div>
  )
}
