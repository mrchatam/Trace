import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ErrorBanner, EmptyState } from '../components/ErrorBanner'
import { GateStrip } from '../components/GateStrip'
import { useApiToken } from '../context/AppChrome'
import {
  getLoopGate,
  getLoopStatus,
  listTasks,
  listTasksForPick,
  search,
  type LoopGateResponse,
  type LoopStatusResponse,
  type SearchItem,
  type TaskRow,
} from '../api/ops'
import { pickActiveTask } from '../lib/pickActiveTask'

export function Overview() {
  const token = useApiToken()
  const [tasks, setTasks] = useState<TaskRow[]>([])
  const [goals, setGoals] = useState<SearchItem[]>([])
  const [active, setActive] = useState<TaskRow | null>(null)
  const [status, setStatus] = useState<LoopStatusResponse | null>(null)
  const [gate, setGate] = useState<LoopGateResponse | null>(null)
  const [error, setError] = useState<unknown>(null)
  const [loading, setLoading] = useState(true)

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const [displayRes, pickItems, goalRes] = await Promise.all([
        listTasks({ limit: 100 }, { token }),
        listTasksForPick({}, { token }),
        search('goal', 10, { token }).catch(() => ({ items: [] as SearchItem[] })),
      ])
      setTasks(displayRes.items ?? [])
      setGoals((goalRes.items ?? []).filter((i) => i.kind === 'goal' || i.title.toLowerCase().includes('goal')))
      const pick = pickActiveTask(pickItems)
      setActive(pick)
      if (pick) {
        // Independent loads — status errors must not hide gate (envelope stays visible).
        const [stRes, gtRes] = await Promise.allSettled([
          getLoopStatus(pick.id, pick.goal_id ?? undefined, { token }),
          getLoopGate(pick.id, 'edit', { token }),
        ])
        setStatus(stRes.status === 'fulfilled' ? stRes.value : null)
        setGate(gtRes.status === 'fulfilled' ? gtRes.value : null)
        if (stRes.status === 'rejected') setError(stRes.reason)
        else if (gtRes.status === 'rejected') setError(gtRes.reason)
      } else {
        setStatus(null)
        setGate(null)
      }
    } catch (err) {
      setError(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void load()
  }, [token])

  const statusObj = status as Record<string, unknown> | null
  const advisories = collectAdvisories(statusObj, gate)

  return (
    <div>
      <h1 className="page-title">Overview</h1>
      <p className="page-lead">Summaries only — not a full plan tree or graph dump.</p>

      {error ? <ErrorBanner error={error} onRetry={() => void load()} /> : null}

      <div className="panel">
        <h2>Loop gate / status</h2>
        <GateStrip gate={gate} loading={loading} />
        {advisories.length > 0 ? (
          <div className="banner banner--warn" role="status" style={{ marginTop: '0.75rem' }}>
            <strong>Plan advisories</strong>
            <ul style={{ margin: '0.5rem 0 0', paddingLeft: '1.25rem' }}>
              {advisories.map((a) => (
                <li key={a.code}>
                  <span className="mono">{a.code}: </span>
                  {a.message}
                </li>
              ))}
            </ul>
          </div>
        ) : null}
        {statusObj ? (
          <div className="stack" style={{ marginTop: '0.75rem' }}>
            <div className="mono">
              saturated: {String(statusObj.saturated ?? '—')} · reason: {String(statusObj.reason ?? '—')}
            </div>
            {Array.isArray(statusObj.violations) && statusObj.violations.length > 0 ? (
              <div className="banner banner--warn" role="status">
                Status violations: {statusObj.violations.length}
              </div>
            ) : null}
          </div>
        ) : null}
        <div className="row" style={{ marginTop: '0.75rem' }}>
          <Link className="btn" to="/loop">
            Open Loop
          </Link>
          <button type="button" className="btn" onClick={() => void load()} disabled={loading}>
            Refresh
          </button>
        </div>
      </div>

      <div className="panel">
        <h2>Active task</h2>
        {loading ? (
          <p role="status">Loading…</p>
        ) : active ? (
          <div className="stack">
            <Link to={`/tasks/${active.id}`}>{active.title}</Link>
            <span className="pill">{active.work_state}</span>
            <span className="mono">{active.id}</span>
          </div>
        ) : (
          <EmptyState title="No tasks yet — add via CLI/MCP or Discoveries (S05)." />
        )}
      </div>

      <div className="panel">
        <h2>Goals (search)</h2>
        {goals.length === 0 ? (
          <EmptyState title="No goals matched — search may be empty." />
        ) : (
          <ul className="graph-list">
            {goals.map((g) => (
              <li key={g.id}>
                <span className="pill">{g.kind}</span> {g.title}
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="panel">
        <h2>Task count</h2>
        <p>{tasks.length} task(s) in store</p>
        <Link to="/tasks">Browse tasks</Link>
      </div>
    </div>
  )
}

type AdvisoryCopy = { code: string; message: string }

function collectAdvisories(
  status: Record<string, unknown> | null,
  gate: LoopGateResponse | null,
): AdvisoryCopy[] {
  const out: AdvisoryCopy[] = []
  const seen = new Set<string>()
  const add = (code: string, message: string) => {
    if (!code || seen.has(code)) return
    seen.add(code)
    out.push({ code, message: message || code })
  }

  if (status && Array.isArray(status.advisories)) {
    for (const raw of status.advisories) {
      if (!raw || typeof raw !== 'object') continue
      const a = raw as Record<string, unknown>
      add(String(a.code ?? ''), String(a.message ?? ''))
    }
  }

  const violations = Array.isArray(gate?.violations) ? gate.violations : []
  for (const v of violations) {
    const row = v as { reason_code?: string; message?: string }
    const code = row.reason_code ?? ''
    if (code === 'plan_missing' || code === 'goal_plan_gap_terminal_advisory') {
      add(code, row.message ?? code)
    }
  }

  return out
}
