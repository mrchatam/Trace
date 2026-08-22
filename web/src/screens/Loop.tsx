import { useEffect, useMemo, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { ConfirmDialog } from '../components/ConfirmDialog'
import { ErrorBanner, EmptyState } from '../components/ErrorBanner'
import { GateStrip } from '../components/GateStrip'
import { useApiToken } from '../context/AppChrome'
import {
  getLoopGate,
  getLoopNext,
  getLoopStatus,
  listTasks,
  listTasksForPick,
  postLoopApply,
  postLoopReset,
  type LoopApplyEnvelope,
  type LoopGateResponse,
  type LoopNextResponse,
  type LoopStatusResponse,
  type TaskRow,
} from '../api/ops'
import { pickActiveTask } from '../lib/pickActiveTask'

function draftEmptyApply(taskId: string, goalId: string): LoopApplyEnvelope {
  return {
    schema_version: 'trace.loop.apply.v1',
    apply_id: crypto.randomUUID(),
    seed: { task_id: taskId, goal_id: goalId || undefined },
    writes: {},
  }
}

export function Loop() {
  const token = useApiToken()
  const [params, setParams] = useSearchParams()
  const [tasks, setTasks] = useState<TaskRow[]>([])
  const [status, setStatus] = useState<LoopStatusResponse | null>(null)
  const [gate, setGate] = useState<LoopGateResponse | null>(null)
  const [nextPkt, setNextPkt] = useState<LoopNextResponse | null>(null)
  const [applyDraft, setApplyDraft] = useState('')
  const [applyResult, setApplyResult] = useState<unknown>(null)
  const [error, setError] = useState<unknown>(null)
  const [writeError, setWriteError] = useState<unknown>(null)
  const [loading, setLoading] = useState(true)
  const [busy, setBusy] = useState(false)
  const [gateFor, setGateFor] = useState('edit')
  const [confirm, setConfirm] = useState<'apply' | 'reset' | null>(null)

  const taskId = params.get('task_id') ?? ''

  useEffect(() => {
    void (async () => {
      try {
        const [displayRes, pickItems] = await Promise.all([
          listTasks({ limit: 100 }, { token }),
          listTasksForPick({}, { token }),
        ])
        setTasks(displayRes.items ?? [])
        // Explicit ?task_id= wins — never overwrite when present.
        if (!taskId) {
          const pick = pickActiveTask(pickItems)
          if (pick) setParams({ task_id: pick.id }, { replace: true })
        }
      } catch (err) {
        setError(err)
      }
    })()
  }, [token])

  async function load() {
    if (!taskId) {
      setStatus(null)
      setGate(null)
      setLoading(false)
      return
    }
    setLoading(true)
    setError(null)
    // Load independently — status failure must not blank a successful gate.
    const [stRes, gtRes] = await Promise.allSettled([
      getLoopStatus(taskId, undefined, { token }),
      getLoopGate(taskId, gateFor, { token }),
    ])
    if (stRes.status === 'fulfilled') {
      setStatus(stRes.value)
    } else {
      setStatus(null)
      setError(stRes.reason)
    }
    if (gtRes.status === 'fulfilled') {
      setGate(gtRes.value)
    } else {
      setGate(null)
      if (stRes.status === 'fulfilled') setError(gtRes.reason)
    }
    setLoading(false)
  }

  useEffect(() => {
    void load()
  }, [taskId, gateFor, token])

  const goalId = useMemo(() => {
    const st = status as Record<string, unknown> | null
    const seed = st?.seed as Record<string, unknown> | undefined
    if (typeof seed?.goal_id === 'string') return seed.goal_id
    const next = nextPkt as Record<string, unknown> | null
    const nseed = next?.seed as Record<string, unknown> | undefined
    if (typeof nseed?.goal_id === 'string') return nseed.goal_id
    const task = tasks.find((t) => t.id === taskId)
    return task?.goal_id ?? ''
  }, [status, nextPkt, tasks, taskId])

  async function fetchNext() {
    if (!taskId) return
    setBusy(true)
    setWriteError(null)
    try {
      const pkt = await getLoopNext(taskId, { token })
      setNextPkt(pkt)
      const draft = draftEmptyApply(taskId, goalId || ((pkt as Record<string, unknown>).seed as { goal_id?: string })?.goal_id || '')
      setApplyDraft(JSON.stringify(draft, null, 2))
    } catch (err) {
      setWriteError(err)
    } finally {
      setBusy(false)
    }
  }

  async function doApply() {
    setBusy(true)
    setWriteError(null)
    setApplyResult(null)
    try {
      const body = JSON.parse(applyDraft) as LoopApplyEnvelope
      const res = await postLoopApply(body, { token })
      setApplyResult(res)
      setConfirm(null)
      await load()
    } catch (err) {
      setWriteError(err)
      setConfirm(null)
    } finally {
      setBusy(false)
    }
  }

  async function doReset() {
    if (!taskId) return
    setBusy(true)
    setWriteError(null)
    try {
      await postLoopReset(taskId, { token })
      setConfirm(null)
      setNextPkt(null)
      setApplyResult(null)
      await load()
    } catch (err) {
      setWriteError(err)
      setConfirm(null)
    } finally {
      setBusy(false)
    }
  }

  const statusObj = status as Record<string, unknown> | null

  return (
    <div>
      <h1 className="page-title">Loop</h1>
      <p className="page-lead">
        Status, gate, next packet, apply, and reset. Apply/reset require confirmation (Enter /
        Escape).
      </p>

      <div className="row" style={{ marginBottom: '1rem' }}>
        <div className="field">
          <label htmlFor="loop-task">task_id</label>
          <select
            id="loop-task"
            value={taskId}
            onChange={(e) => setParams(e.target.value ? { task_id: e.target.value } : {})}
          >
            <option value="">Select task…</option>
            {tasks.map((t) => (
              <option key={t.id} value={t.id}>
                {t.title} ({t.work_state})
              </option>
            ))}
          </select>
        </div>
        <div className="field">
          <label htmlFor="gate-for">gate for</label>
          <select id="gate-for" value={gateFor} onChange={(e) => setGateFor(e.target.value)}>
            <option value="orient">orient</option>
            <option value="edit">edit</option>
            <option value="execute">execute</option>
            <option value="done">done</option>
            <option value="export">export</option>
          </select>
        </div>
        <button type="button" className="btn" onClick={() => void load()} disabled={loading}>
          Refresh
        </button>
      </div>

      {error ? <ErrorBanner error={error} onRetry={() => void load()} /> : null}
      {writeError ? <ErrorBanner error={writeError} /> : null}

      {!taskId ? (
        <EmptyState title="Pick a task to load loop status and gate.">
          <Link to="/tasks">Browse tasks</Link>
        </EmptyState>
      ) : (
        <>
          <div className="panel">
            <h2>Gate</h2>
            <GateStrip gate={gate} loading={loading} />
          </div>
          <div className="panel">
            <h2>Status</h2>
            {loading ? (
              <p role="status">Loading…</p>
            ) : statusObj ? (
              <pre className="drawer" style={{ maxHeight: '16rem' }}>
                {JSON.stringify(statusObj, null, 2)}
              </pre>
            ) : (
              <EmptyState title="No status payload (gate may still be available)." />
            )}
          </div>

          <div className="panel">
            <h2>Write console</h2>
            <div className="row" style={{ marginBottom: '0.75rem' }}>
              <button type="button" className="btn btn--primary" disabled={busy} onClick={() => void fetchNext()}>
                Fetch next
              </button>
              <button
                type="button"
                className="btn"
                disabled={busy || !applyDraft}
                onClick={() => setConfirm('apply')}
              >
                Apply…
              </button>
              <button
                type="button"
                className="btn btn--danger"
                disabled={busy}
                onClick={() => setConfirm('reset')}
              >
                Reset…
              </button>
            </div>
            {nextPkt ? (
              <div className="drawer" style={{ marginBottom: '0.75rem', maxHeight: '14rem' }}>
                <strong>Next packet</strong>
                <pre>{JSON.stringify(nextPkt, null, 2)}</pre>
              </div>
            ) : null}
            <div className="field">
              <label htmlFor="apply-draft">Apply envelope (JSON)</label>
              <textarea
                id="apply-draft"
                rows={12}
                value={applyDraft}
                onChange={(e) => setApplyDraft(e.target.value)}
                placeholder="Fetch next to seed a draft, or paste an envelope"
                style={{ fontFamily: 'var(--font-mono)', fontSize: '0.75rem' }}
              />
            </div>
            {applyResult ? (
              <div className="drawer" style={{ marginTop: '0.75rem' }}>
                <strong>Apply result</strong>
                <pre>{JSON.stringify(applyResult, null, 2)}</pre>
              </div>
            ) : null}
          </div>
        </>
      )}

      <ConfirmDialog
        open={confirm === 'apply'}
        title="Apply loop envelope?"
        confirmLabel="Apply"
        busy={busy}
        onCancel={() => setConfirm(null)}
        onConfirm={() => void doApply()}
      >
        Posts the JSON envelope to <code className="mono">POST /v1/loop/apply</code>. Empty writes
        still count toward saturation.
      </ConfirmDialog>

      <ConfirmDialog
        open={confirm === 'reset'}
        title="Reset deliberation?"
        confirmLabel="Reset"
        danger
        busy={busy}
        onCancel={() => setConfirm(null)}
        onConfirm={() => void doReset()}
      >
        Resets deliberation state for task <code className="mono">{taskId}</code>.
      </ConfirmDialog>
    </div>
  )
}
