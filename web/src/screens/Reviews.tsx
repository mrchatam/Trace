import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ErrorBanner, EmptyState } from '../components/ErrorBanner'
import { useApiToken } from '../context/AppChrome'
import { listReviews, type ReviewSummary } from '../api/ops'

export function Reviews() {
  const token = useApiToken()
  const [items, setItems] = useState<ReviewSummary[]>([])
  const [error, setError] = useState<unknown>(null)
  const [loading, setLoading] = useState(true)
  const [taskFilter, setTaskFilter] = useState('')

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const res = await listReviews(
        taskFilter.trim() ? { task_id: taskFilter.trim() } : {},
        { token },
      )
      setItems(res.items ?? [])
    } catch (err) {
      setError(err)
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void load()
  }, [token])

  return (
    <div>
      <h1 className="page-title">Reviews</h1>
      <p className="page-lead">List and open reviews. Create from a task detail when needed.</p>

      <div className="row" style={{ marginBottom: '1rem' }}>
        <div className="field" style={{ flex: 1 }}>
          <label htmlFor="rev-task">Filter by task_id</label>
          <input
            id="rev-task"
            value={taskFilter}
            onChange={(e) => setTaskFilter(e.target.value)}
            placeholder="optional uuid"
          />
        </div>
        <button type="button" className="btn" onClick={() => void load()} disabled={loading}>
          {loading ? 'Loading…' : 'Refresh'}
        </button>
      </div>

      {error ? <ErrorBanner error={error} onRetry={() => void load()} /> : null}

      <div className="panel">
        {items.length === 0 && !loading ? (
          <EmptyState title="No reviews yet.">
            Create one from a <Link to="/tasks">task detail</Link>.
          </EmptyState>
        ) : (
          <table className="table">
            <thead>
              <tr>
                <th>Title</th>
                <th>Result</th>
                <th>Task</th>
              </tr>
            </thead>
            <tbody>
              {items.map((r) => (
                <tr key={r.id}>
                  <td>
                    <Link to={`/reviews/${encodeURIComponent(r.id)}`}>{r.title}</Link>
                  </td>
                  <td>
                    {r.result ? <span className="pill">{r.result}</span> : '—'}
                  </td>
                  <td className="mono">{r.task_id ?? '—'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
