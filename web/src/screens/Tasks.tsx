import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ErrorBanner, EmptyState } from '../components/ErrorBanner'
import { useApiToken } from '../context/AppChrome'
import { listTasks, type TaskRow } from '../api/ops'

export function Tasks() {
  const token = useApiToken()
  const [items, setItems] = useState<TaskRow[]>([])
  const [filter, setFilter] = useState('')
  const [error, setError] = useState<unknown>(null)
  const [loading, setLoading] = useState(true)

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const res = await listTasks(
        { limit: 200, work_state: filter || undefined },
        { token },
      )
      setItems(res.items ?? [])
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void load()
  }, [token, filter])

  return (
    <div>
      <h1 className="page-title">Tasks</h1>
      <p className="page-lead">Browse work and open a task for detail / light transition.</p>

      <div className="row" style={{ marginBottom: '1rem' }}>
        <div className="field">
          <label htmlFor="ws-filter">work_state</label>
          <select
            id="ws-filter"
            value={filter}
            onChange={(e) => setFilter(e.target.value)}
          >
            <option value="">All</option>
            <option value="PENDING">PENDING</option>
            <option value="IN_PROGRESS">IN_PROGRESS</option>
            <option value="AWAITING_REVIEW">AWAITING_REVIEW</option>
            <option value="BLOCKED">BLOCKED</option>
            <option value="FAILED">FAILED</option>
            <option value="DONE">DONE</option>
            <option value="STALE">STALE</option>
            <option value="SKIPPED">SKIPPED</option>
          </select>
        </div>
        <button type="button" className="btn" onClick={() => void load()} disabled={loading}>
          Refresh
        </button>
      </div>

      {error ? <ErrorBanner error={error} onRetry={() => void load()} /> : null}

      {loading ? (
        <p role="status">Loading tasks…</p>
      ) : items.length === 0 ? (
        <EmptyState title={filter ? 'No matches for this filter.' : 'No tasks in store.'}>
          {filter ? (
            <button type="button" className="btn" onClick={() => setFilter('')}>
              Clear filter
            </button>
          ) : null}
        </EmptyState>
      ) : (
        <div className="panel" style={{ padding: 0 }}>
          <table className="table">
            <thead>
              <tr>
                <th scope="col">Title</th>
                <th scope="col">State</th>
                <th scope="col">ID</th>
              </tr>
            </thead>
            <tbody>
              {items.map((t) => (
                <tr key={t.id}>
                  <td>
                    <Link to={`/tasks/${t.id}`}>{t.title}</Link>
                  </td>
                  <td>
                    <span className="pill">{t.work_state}</span>
                  </td>
                  <td className="mono">{t.id}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
