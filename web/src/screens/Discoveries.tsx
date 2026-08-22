import { useEffect, useMemo, useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { ErrorBanner, EmptyState } from '../components/ErrorBanner'
import { useApiToken } from '../context/AppChrome'
import {
  createEntity,
  createLink,
  createTransition,
  getEntity,
  listTasks,
  search,
  type EntitySummary,
  type SearchItem,
  type TaskRow,
} from '../api/ops'

const KIND_OPTIONS = ['all', 'discovery', 'decision', 'task', 'goal', 'assumption'] as const

export function Discoveries() {
  const token = useApiToken()
  const [q, setQ] = useState('discovery')
  const [kind, setKind] = useState<(typeof KIND_OPTIONS)[number]>('all')
  const [items, setItems] = useState<SearchItem[]>([])
  const [selected, setSelected] = useState<EntitySummary | null>(null)
  const [tasks, setTasks] = useState<TaskRow[]>([])
  const [error, setError] = useState<unknown>(null)
  const [actionError, setActionError] = useState<unknown>(null)
  const [actionOk, setActionOk] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [busy, setBusy] = useState(false)

  // Create form
  const [createKind, setCreateKind] = useState<'discovery' | 'decision'>('discovery')
  const [createTitle, setCreateTitle] = useState('')
  const [createBody, setCreateBody] = useState('')

  // Promote form
  const [promoteGoalId, setPromoteGoalId] = useState('')
  const [promoteTitle, setPromoteTitle] = useState('')

  async function runSearch(e?: FormEvent) {
    e?.preventDefault()
    setLoading(true)
    setError(null)
    setSelected(null)
    try {
      const res = await search(q.trim() || 'discovery', 40, { token })
      const all = res.items ?? []
      setItems(all)
    } catch (err) {
      setError(err)
      setItems([])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void runSearch()
    void listTasks({ limit: 50 }, { token })
      .then((r) => setTasks(r.items ?? []))
      .catch(() => setTasks([]))
  }, [token])

  const filtered = useMemo(() => {
    if (kind === 'all') {
      const preferred = items.filter(
        (i) => i.kind === 'discovery' || i.kind === 'decision' || q.includes(i.kind),
      )
      return preferred.length > 0 ? preferred : items
    }
    return items.filter((i) => i.kind === kind)
  }, [items, kind, q])

  async function openDetail(id: string) {
    setActionError(null)
    setActionOk(null)
    try {
      const ent = await getEntity(id, { token })
      setSelected(ent)
      setPromoteTitle(ent.title)
    } catch (err) {
      setError(err)
      setSelected(null)
    }
  }

  async function onCreate(e: FormEvent) {
    e.preventDefault()
    setBusy(true)
    setActionError(null)
    setActionOk(null)
    try {
      const created = await createEntity(
        { kind: createKind, title: createTitle.trim(), body: createBody.trim() || undefined },
        { token },
      )
      setActionOk(`Created ${created.kind} ${created.id}`)
      setCreateTitle('')
      setCreateBody('')
      await runSearch()
      setSelected(created)
      setPromoteTitle(created.title)
    } catch (err) {
      setActionError(err)
    } finally {
      setBusy(false)
    }
  }

  async function onPromote(e: FormEvent) {
    e.preventDefault()
    if (!selected) return
    setBusy(true)
    setActionError(null)
    setActionOk(null)
    try {
      const task = await createEntity(
        {
          kind: 'task',
          title: promoteTitle.trim() || selected.title,
          goal_id: promoteGoalId || undefined,
          body: selected.body ?? undefined,
        },
        { token },
      )
      try {
        await createLink(
          {
            rel: 'discovery-mentions-task',
            from: selected.id,
            to: task.id,
            source_type: 'gui',
          },
          { token },
        )
      } catch (linkErr) {
        // Non-discovery kinds may reject this rel — surface envelope, keep task.
        setActionError(linkErr)
        setActionOk(`Task created ${task.id} (link failed — see envelope)`)
        setBusy(false)
        return
      }
      try {
        await createTransition(
          {
            task_id: task.id,
            to_state: 'PENDING',
            reason: 'gui promote from discovery',
            actor: 'gui',
          },
          { token },
        )
      } catch (transErr) {
        // Task + link still succeed; surface transition deny honestly (do not swallow).
        setActionError(transErr)
        setActionOk(`Promoted → task ${task.id} (transition denied — see envelope)`)
        await runSearch()
        return
      }
      setActionOk(`Promoted → task ${task.id}`)
      await runSearch()
    } catch (err) {
      setActionError(err)
    } finally {
      setBusy(false)
    }
  }

  const goals = useMemo(() => {
    const g = new Map<string, string>()
    for (const t of tasks) {
      if (t.goal_id) g.set(t.goal_id, t.goal_id)
    }
    return Array.from(g.keys())
  }, [tasks])

  return (
    <div>
      <h1 className="page-title">Discoveries</h1>
      <p className="page-lead">
        Search, create discoveries/decisions, and promote to a linked task via library-backed{' '}
        <code className="mono">/v1</code> ops.
      </p>

      <form className="filters" onSubmit={(e) => void runSearch(e)}>
        <div className="field" style={{ flex: 1 }}>
          <label htmlFor="disc-q">Search (q)</label>
          <input
            id="disc-q"
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="discovery | decision | …"
          />
        </div>
        <div className="field">
          <label htmlFor="disc-kind">Kind filter</label>
          <select
            id="disc-kind"
            value={kind}
            onChange={(e) => setKind(e.target.value as (typeof KIND_OPTIONS)[number])}
          >
            {KIND_OPTIONS.map((k) => (
              <option key={k} value={k}>
                {k}
              </option>
            ))}
          </select>
        </div>
        <button type="submit" className="btn btn--primary" disabled={loading}>
          {loading ? 'Searching…' : 'Search'}
        </button>
      </form>

      {error ? <ErrorBanner error={error} /> : null}
      {actionError ? <ErrorBanner error={actionError} /> : null}
      {actionOk ? (
        <div className="banner banner--ok" role="status">
          {actionOk}
        </div>
      ) : null}

      <div className="panel">
        <h2>Create</h2>
        <form className="stack" onSubmit={(e) => void onCreate(e)}>
          <div className="row">
            <div className="field">
              <label htmlFor="create-kind">kind</label>
              <select
                id="create-kind"
                value={createKind}
                onChange={(e) => setCreateKind(e.target.value as 'discovery' | 'decision')}
              >
                <option value="discovery">discovery</option>
                <option value="decision">decision</option>
              </select>
            </div>
            <div className="field" style={{ flex: 1 }}>
              <label htmlFor="create-title">title</label>
              <input
                id="create-title"
                value={createTitle}
                onChange={(e) => setCreateTitle(e.target.value)}
                required
              />
            </div>
          </div>
          <div className="field">
            <label htmlFor="create-body">body</label>
            <textarea
              id="create-body"
              rows={3}
              value={createBody}
              onChange={(e) => setCreateBody(e.target.value)}
            />
          </div>
          <button type="submit" className="btn btn--primary" disabled={busy || !createTitle.trim()}>
            Create
          </button>
        </form>
      </div>

      <div className="panel">
        <h2>Results</h2>
        {filtered.length === 0 ? (
          <EmptyState title="No discoveries or decisions yet.">
            Create one above, or search with a different query.
          </EmptyState>
        ) : (
          <ul className="graph-list">
            {filtered.map((i) => (
              <li key={i.id}>
                <button type="button" className="btn btn--ghost" onClick={() => void openDetail(i.id)}>
                  <span className="pill">{i.kind}</span> {i.title}
                </button>
                {i.snippet ? <div className="page-lead">{i.snippet}</div> : null}
              </li>
            ))}
          </ul>
        )}
      </div>

      {selected ? (
        <div className="panel">
          <h2>Detail</h2>
          <div className="stack">
            <div>
              <span className="pill">{selected.kind}</span> {selected.title}
            </div>
            <div className="mono">{selected.id}</div>
            {selected.work_state ? <span className="pill">{selected.work_state}</span> : null}
            {selected.body ? <p>{selected.body}</p> : null}
          </div>

          {selected.kind === 'discovery' ? (
            <form className="stack" style={{ marginTop: '1rem' }} onSubmit={(e) => void onPromote(e)}>
              <h3 style={{ margin: 0, fontSize: '0.95rem' }}>Promote to task</h3>
              <p className="page-lead" style={{ margin: 0 }}>
                createEntity(task) → createLink(discovery-mentions-task) → optional transition.
                Denials show the API envelope only.
              </p>
              <div className="field">
                <label htmlFor="promote-title">task title</label>
                <input
                  id="promote-title"
                  value={promoteTitle}
                  onChange={(e) => setPromoteTitle(e.target.value)}
                  required
                />
              </div>
              <div className="field">
                <label htmlFor="promote-goal">goal_id (optional)</label>
                <select
                  id="promote-goal"
                  value={promoteGoalId}
                  onChange={(e) => setPromoteGoalId(e.target.value)}
                >
                  <option value="">—</option>
                  {goals.map((g) => (
                    <option key={g} value={g}>
                      {g}
                    </option>
                  ))}
                </select>
              </div>
              <div className="row">
                <button type="submit" className="btn btn--primary" disabled={busy}>
                  Promote
                </button>
                <Link className="btn" to="/tasks">
                  Open tasks
                </Link>
              </div>
            </form>
          ) : null}
        </div>
      ) : null}
    </div>
  )
}
