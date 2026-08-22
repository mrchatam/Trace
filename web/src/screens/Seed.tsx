import { useEffect, useState, type FormEvent } from 'react'
import { ConfirmDialog } from '../components/ConfirmDialog'
import { ErrorBanner } from '../components/ErrorBanner'
import { useApiToken } from '../context/AppChrome'
import {
  getSeedStatus,
  postSeedExport,
  postSeedImport,
  type SeedJobStatus,
  type SeedStatus,
} from '../api/ops'

export function Seed() {
  const token = useApiToken()
  const [status, setStatus] = useState<SeedStatus | null>(null)
  const [error, setError] = useState<unknown>(null)
  const [jobError, setJobError] = useState<unknown>(null)
  const [job, setJob] = useState<SeedJobStatus | null>(null)
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)

  const [exportPath, setExportPath] = useState('trace/graph.json')
  const [exportStrict, setExportStrict] = useState(false)
  const [exportTaskId, setExportTaskId] = useState('')
  const [importPath, setImportPath] = useState('trace/graph.json')
  const [confirmImport, setConfirmImport] = useState(false)

  async function load() {
    setLoading(true)
    setError(null)
    try {
      setStatus(await getSeedStatus({ token }))
    } catch (err) {
      setError(err)
      setStatus(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void load()
  }, [token])

  async function onExport(e: FormEvent) {
    e.preventDefault()
    setBusy(true)
    setJobError(null)
    setJob(null)
    try {
      const res = await postSeedExport(
        {
          output_path: exportPath.trim() || undefined,
          strict: exportStrict || undefined,
          task_id: exportTaskId.trim() || undefined,
        },
        { token },
      )
      setJob(res)
      await load()
    } catch (err) {
      // Path escape / strict 501 — show envelope honestly
      setJobError(err)
    } finally {
      setBusy(false)
    }
  }

  async function doImport() {
    setBusy(true)
    setJobError(null)
    setJob(null)
    try {
      const res = await postSeedImport({ input_path: importPath.trim() }, { token })
      setJob(res)
      setConfirmImport(false)
      await load()
    } catch (err) {
      setJobError(err)
      setConfirmImport(false)
    } finally {
      setBusy(false)
    }
  }

  return (
    <div>
      <h1 className="page-title">Seed</h1>
      <p className="page-lead">
        Export/import with project-relative paths. Responses are status/summary only — not a full
        graph dump.
      </p>

      <div className="banner banner--warn" role="note">
        Prefer project-relative paths (e.g. <code className="mono">trace/graph.json</code>). Path
        escape → <code className="mono">VALIDATION_ERROR</code>. Export{' '}
        <code className="mono">strict</code>/<code className="mono">task_id</code> → API{' '}
        <code className="mono">501 NOT_IMPLEMENTED</code> (shown honestly, not silently ignored).
      </div>

      {error ? <ErrorBanner error={error} onRetry={() => void load()} /> : null}
      {jobError ? <ErrorBanner error={jobError} /> : null}

      <div className="panel">
        <h2>Status</h2>
        <div className="row" style={{ marginBottom: '0.75rem' }}>
          <button type="button" className="btn" onClick={() => void load()} disabled={loading}>
            {loading ? 'Refreshing…' : 'Refresh'}
          </button>
        </div>
        {status ? (
          <dl className="stack">
            <div>
              <dt>ready</dt>
              <dd>
                <span className={`pill ${status.ready ? 'pill--ok' : 'pill--bad'}`}>
                  {String(status.ready)}
                </span>
              </dd>
            </div>
            <div>
              <dt>last_export_at</dt>
              <dd className="mono">{status.last_export_at ?? '—'}</dd>
            </div>
            <div>
              <dt>last_import_at</dt>
              <dd className="mono">{status.last_import_at ?? '—'}</dd>
            </div>
            {status.notes ? (
              <div>
                <dt>notes</dt>
                <dd>{status.notes}</dd>
              </div>
            ) : null}
          </dl>
        ) : !loading ? (
          <p>No status yet.</p>
        ) : null}
      </div>

      <div className="panel">
        <h2>Export</h2>
        <form className="stack" onSubmit={(e) => void onExport(e)}>
          <div className="field">
            <label htmlFor="export-path">output_path (project-relative)</label>
            <input
              id="export-path"
              value={exportPath}
              onChange={(e) => setExportPath(e.target.value)}
              placeholder="trace/graph.json"
            />
          </div>
          <div className="row">
            <label className="row" style={{ gap: '0.35rem' }}>
              <input
                type="checkbox"
                checked={exportStrict}
                onChange={(e) => setExportStrict(e.target.checked)}
              />
              strict (expect 501)
            </label>
          </div>
          <div className="field">
            <label htmlFor="export-task">task_id (optional; expect 501 with API)</label>
            <input
              id="export-task"
              value={exportTaskId}
              onChange={(e) => setExportTaskId(e.target.value)}
              placeholder="uuid — surfaces NOT_IMPLEMENTED"
            />
          </div>
          <button type="submit" className="btn btn--primary" disabled={busy}>
            Export
          </button>
        </form>
      </div>

      <div className="panel">
        <h2>Import</h2>
        <form
          className="stack"
          onSubmit={(e) => {
            e.preventDefault()
            setConfirmImport(true)
          }}
        >
          <div className="field">
            <label htmlFor="import-path">input_path (project-relative)</label>
            <input
              id="import-path"
              value={importPath}
              onChange={(e) => setImportPath(e.target.value)}
              required
            />
          </div>
          <button type="submit" className="btn" disabled={busy || !importPath.trim()}>
            Import…
          </button>
        </form>
      </div>

      {job ? (
        <div className="panel">
          <h2>Last job</h2>
          <pre className="drawer">{JSON.stringify(job, null, 2)}</pre>
        </div>
      ) : null}

      <ConfirmDialog
        open={confirmImport}
        title="Import seed?"
        confirmLabel="Import"
        danger
        busy={busy}
        onCancel={() => setConfirmImport(false)}
        onConfirm={() => void doImport()}
      >
        Imports from <code className="mono">{importPath}</code> under the bound project root. Path
        escape is rejected by the API.
      </ConfirmDialog>
    </div>
  )
}
