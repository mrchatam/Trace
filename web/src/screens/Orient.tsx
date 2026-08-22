import { useEffect, useMemo, useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { ErrorBanner, EmptyState } from '../components/ErrorBanner'
import { GraphOrientPanel } from '../components/GraphOrientPanel'
import { useApiToken } from '../context/AppChrome'
import { listTasks, search, type SearchItem, type TaskRow } from '../api/ops'
import { buildGraphHref } from '../lib/graphNavigate'
import { EXPAND_MAX_NODES, UI_CAP } from '../lib/overviewCompose'

type KindFilter = 'all' | string

export function Orient() {
  const token = useApiToken()
  const [tasks, setTasks] = useState<TaskRow[]>([])
  const [searchItems, setSearchItems] = useState<SearchItem[]>([])
  const [q, setQ] = useState('')
  const [kind, setKind] = useState<KindFilter>('all')
  const [maxNodes, setMaxNodes] = useState(EXPAND_MAX_NODES)
  const [error, setError] = useState<unknown>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    void (async () => {
      setLoading(true)
      setError(null)
      try {
        const res = await listTasks({ limit: 100 }, { token })
        setTasks(res.items ?? [])
      } catch (err) {
        setError(err)
      } finally {
        setLoading(false)
      }
    })()
  }, [token])

  const filteredSearch = useMemo(() => {
    if (kind === 'all') return searchItems
    return searchItems.filter((i) => i.kind === kind)
  }, [searchItems, kind])

  const kinds = useMemo(() => {
    const s = new Set(searchItems.map((i) => i.kind))
    if (tasks.length > 0) s.add('task')
    return Array.from(s).sort()
  }, [searchItems, tasks.length])

  async function runSearch(e: FormEvent) {
    e.preventDefault()
    if (!q.trim()) return
    setError(null)
    try {
      const res = await search(q.trim(), 40, { token })
      setSearchItems(res.items ?? [])
    } catch (err) {
      setError(err)
    }
  }

  function graphLink(entityId: string) {
    const href = buildGraphHref(entityId)
    const budget = Math.min(UI_CAP, Math.max(1, maxNodes))
    return budget === EXPAND_MAX_NODES ? href : `${href}&max_nodes=${budget}`
  }

  return (
    <div>
      <h1 className="page-title">Orient</h1>
      <p className="page-lead">
        Search entities or pick a task, then open a bounded neighborhood on Explore. Explore stays
        graph-first — orientation and pickers live here.
      </p>

      <GraphOrientPanel persistent />

      {error ? <ErrorBanner error={error} onRetry={() => setError(null)} /> : null}

      <div className="panel" data-testid="graph-manual-center">
        <h2>Find a center</h2>
        <form className="filters" onSubmit={(e) => void runSearch(e)}>
          <div className="field" style={{ flex: 1 }}>
            <label htmlFor="orient-q">Search (q)</label>
            <input
              id="orient-q"
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder="search…"
            />
          </div>
          <div className="field">
            <label htmlFor="orient-kind">Kind filter</label>
            <select id="orient-kind" value={kind} onChange={(e) => setKind(e.target.value)}>
              <option value="all">all</option>
              {kinds.map((k) => (
                <option key={k} value={k}>
                  {k}
                </option>
              ))}
            </select>
          </div>
          <div className="field">
            <label htmlFor="orient-budget-input">max_nodes</label>
            <input
              id="orient-budget-input"
              type="number"
              min={1}
              max={UI_CAP}
              value={maxNodes}
              onChange={(e) =>
                setMaxNodes(Math.min(UI_CAP, Math.max(1, Number(e.target.value) || EXPAND_MAX_NODES)))
              }
            />
          </div>
          <button type="submit" className="btn">
            Search
          </button>
        </form>

        {filteredSearch.length > 0 ? (
          <ul className="graph-list">
            {filteredSearch.map((h) => (
              <li key={h.id}>
                <Link className="btn btn--ghost" to={graphLink(h.id)} data-testid={`orient-pick-${h.id}`}>
                  {h.kind}: {h.title}
                </Link>
                <span className="mono" style={{ marginLeft: '0.5rem' }}>
                  {h.id}
                </span>
              </li>
            ))}
          </ul>
        ) : null}

        <div style={{ marginTop: '1rem' }}>
          <h3 style={{ fontSize: '0.9rem', margin: '0 0 0.5rem' }}>Or pick a task</h3>
          {loading ? (
            <p role="status">Loading tasks…</p>
          ) : tasks.length === 0 ? (
            <EmptyState title="No tasks — open Tasks or search above.">
              <Link className="btn btn--primary" to="/tasks">
                Open Tasks
              </Link>
            </EmptyState>
          ) : (
            <ul className="graph-list" data-testid="graph-task-list">
              {tasks.slice(0, 50).map((t) => (
                <li key={t.id}>
                  <Link
                    className="btn btn--ghost"
                    to={graphLink(t.id)}
                    data-testid={`graph-pick-task-${t.id}`}
                  >
                    {t.title}
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>

      <p>
        <Link to="/">Open Explore graph</Link>
        {' · '}
        <Link to="/tasks">Tasks</Link>
        {' · '}
        <Link to="/loop">Loop</Link>
      </p>
    </div>
  )
}
