import { Link } from 'react-router-dom'
import { ErrorBanner } from '../components/ErrorBanner'
import { CopyButton } from '../components/CopyButton'
import { useAppChrome } from '../context/AppChrome'

export function ProjectGate() {
  const { project, chromeError, refresh, loading } = useAppChrome()

  return (
    <div className="gate-blocked">
      <div className="gate-blocked__card">
        <h1 className="page-title">Project not ready</h1>
        <p className="page-lead">
          No Trace store at this root, or the API is unreachable. Run{' '}
          <code className="mono">trace init</code> in the project, then restart{' '}
          <code className="mono">trace serve</code>.
        </p>

        {project ? (
          <div className="panel">
            <div className="row">
              <span className="mono">{project.root}</span>
              <CopyButton text={project.root} label="Copy root" />
            </div>
            <div style={{ marginTop: '0.5rem' }}>
              store_ready: <strong>{String(project.store_ready)}</strong>
              {project.store_path ? (
                <div className="mono" style={{ marginTop: '0.25rem' }}>
                  {project.store_path}
                </div>
              ) : null}
            </div>
          </div>
        ) : null}

        {chromeError ? <ErrorBanner error={chromeError} onRetry={() => void refresh()} /> : null}

        <div className="row">
          <button type="button" className="btn btn--primary" onClick={() => void refresh()} disabled={loading}>
            {loading ? 'Checking…' : 'Retry health / project'}
          </button>
          <Link className="btn" to="/settings">
            Settings
          </Link>
        </div>
      </div>
    </div>
  )
}
