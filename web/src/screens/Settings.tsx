import { useState, type FormEvent } from 'react'
import { useAppChrome, type Theme } from '../context/AppChrome'

export function Settings() {
  const { theme, setTheme, token, setToken, version, project, health, refresh } = useAppChrome()
  const [draft, setDraft] = useState(token)

  function saveToken(e: FormEvent) {
    e.preventDefault()
    setToken(draft.trim())
    void refresh()
  }

  return (
    <div>
      <h1 className="page-title">Settings</h1>
      <p className="page-lead">Client chrome — theme, bearer token, and serve identity display.</p>

      <div className="panel">
        <h2>Theme</h2>
        <div className="row" role="group" aria-label="Theme">
          {(['light', 'dark'] as Theme[]).map((t) => (
            <button
              key={t}
              type="button"
              className={`btn ${theme === t ? 'btn--primary' : ''}`}
              aria-pressed={theme === t}
              onClick={() => setTheme(t)}
            >
              {t}
            </button>
          ))}
        </div>
      </div>

      <div className="panel">
        <h2>Bearer token</h2>
        <p className="page-lead">
          Loopback trust: token not required on 127.0.0.1. Paste a token when serving with
          --allow-remote.
        </p>
        <form className="stack" onSubmit={saveToken}>
          <div className="field">
            <label htmlFor="token">Authorization Bearer</label>
            <input
              id="token"
              type="password"
              autoComplete="off"
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              placeholder="optional"
            />
          </div>
          <div className="row">
            <button type="submit" className="btn btn--primary">
              Save token
            </button>
            <button
              type="button"
              className="btn"
              onClick={() => {
                setDraft('')
                setToken('')
                void refresh()
              }}
            >
              Clear
            </button>
          </div>
        </form>
        <p style={{ marginTop: '0.75rem' }}>
          {token ? (
            <span className="pill pill--ok">token present</span>
          ) : (
            <span className="pill">token not set (loopback-trust)</span>
          )}
        </p>
      </div>

      <div className="panel">
        <h2>Server / project</h2>
        <dl className="stack">
          <div>
            <dt>API base</dt>
            <dd className="mono">/v1 (same-origin)</dd>
          </div>
          <div>
            <dt>Health</dt>
            <dd>{health?.ok ? 'ok' : 'down'}</dd>
          </div>
          <div>
            <dt>API version</dt>
            <dd className="mono">{version?.api_version ?? '—'}</dd>
          </div>
          <div>
            <dt>Trace version</dt>
            <dd className="mono">{version?.trace_version ?? '—'}</dd>
          </div>
          <div>
            <dt>Project root</dt>
            <dd className="mono">{project?.root ?? '—'}</dd>
          </div>
          <div>
            <dt>Store ready</dt>
            <dd>{project ? String(project.store_ready) : '—'}</dd>
          </div>
          <div>
            <dt>Store path</dt>
            <dd className="mono">{project?.store_path ?? '—'}</dd>
          </div>
        </dl>
        <button type="button" className="btn" onClick={() => void refresh()} style={{ marginTop: '0.75rem' }}>
          Refresh
        </button>
      </div>
    </div>
  )
}
