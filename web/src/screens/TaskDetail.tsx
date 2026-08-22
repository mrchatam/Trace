import { useEffect, useState, type FormEvent } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ErrorBanner } from '../components/ErrorBanner'
import { CopyButton } from '../components/CopyButton'
import { GateStrip } from '../components/GateStrip'
import { useApiToken } from '../context/AppChrome'
import {
  createReview,
  createTransition,
  getContext,
  getLoopGate,
  getTask,
  getWhy,
  type LoopGateResponse,
  type TaskRow,
} from '../api/ops'

/** Suggested next states for UX chrome — API remains SoT (Law 19). */
const ALL_STATES = [
  'PENDING',
  'IN_PROGRESS',
  'AWAITING_REVIEW',
  'BLOCKED',
  'FAILED',
  'DONE',
  'SKIPPED',
  'STALE',
] as const

const SUGGESTED: Record<string, string[]> = {
  PENDING: ['IN_PROGRESS', 'BLOCKED', 'SKIPPED'],
  IN_PROGRESS: ['AWAITING_REVIEW', 'BLOCKED', 'FAILED', 'DONE', 'PENDING'],
  AWAITING_REVIEW: ['DONE', 'IN_PROGRESS', 'FAILED', 'BLOCKED'],
  BLOCKED: ['PENDING', 'IN_PROGRESS', 'SKIPPED', 'FAILED'],
  FAILED: ['PENDING', 'IN_PROGRESS', 'SKIPPED'],
  DONE: ['STALE', 'PENDING'],
  SKIPPED: ['PENDING'],
  STALE: ['PENDING'],
}

export function TaskDetail() {
  const { taskId = '' } = useParams()
  const token = useApiToken()
  const [task, setTask] = useState<TaskRow | null>(null)
  const [gate, setGate] = useState<LoopGateResponse | null>(null)
  const [gateLoading, setGateLoading] = useState(false)
  const [error, setError] = useState<unknown>(null)
  const [txError, setTxError] = useState<unknown>(null)
  const [txOk, setTxOk] = useState<string | null>(null)
  const [toState, setToState] = useState('')
  const [reason, setReason] = useState('gui transition')
  const [allowDone, setAllowDone] = useState(false)
  const [busy, setBusy] = useState(false)
  const [drawer, setDrawer] = useState<'context' | 'why' | null>(null)
  const [drawerData, setDrawerData] = useState<unknown>(null)
  const [drawerErr, setDrawerErr] = useState<unknown>(null)
  const [reviewTitle, setReviewTitle] = useState('')
  const [reviewBody, setReviewBody] = useState('')
  const [reviewMsg, setReviewMsg] = useState<string | null>(null)
  const [reviewErr, setReviewErr] = useState<unknown>(null)

  async function load() {
    setError(null)
    try {
      const t = await getTask(taskId, { token })
      setTask(t)
      const next = SUGGESTED[t.work_state]?.[0] ?? ''
      setToState(next)
      setReviewTitle(`Review: ${t.title}`)
    } catch (err) {
      setError(err)
      setTask(null)
    }
  }

  async function loadGate() {
    if (!taskId) return
    setGateLoading(true)
    try {
      setGate(await getLoopGate(taskId, 'done', { token }))
    } catch {
      setGate(null)
    } finally {
      setGateLoading(false)
    }
  }

  useEffect(() => {
    if (taskId) {
      void load()
      void loadGate()
    }
  }, [taskId, token])

  async function onTransition(e: FormEvent) {
    e.preventDefault()
    if (!task) return
    setBusy(true)
    setTxError(null)
    setTxOk(null)
    try {
      const res = await createTransition(
        {
          task_id: task.id,
          to_state: toState,
          reason,
          actor: 'gui',
          allow_done: toState === 'DONE' ? allowDone : undefined,
        },
        { token },
      )
      setTxOk(`Now ${res.work_state}${res.warning ? ` — ${res.warning}` : ''}`)
      await load()
      await loadGate()
    } catch (err) {
      // Envelope only — no SPA DONE/review policy
      setTxError(err)
    } finally {
      setBusy(false)
    }
  }

  async function openContext() {
    setDrawer('context')
    setDrawerErr(null)
    try {
      setDrawerData(await getContext(taskId, 1, { token }))
    } catch (err) {
      setDrawerErr(err)
      setDrawerData(null)
    }
  }

  async function openWhy() {
    setDrawer('why')
    setDrawerErr(null)
    try {
      setDrawerData(await getWhy('task', taskId, { token }))
    } catch (err) {
      setDrawerErr(err)
      setDrawerData(null)
    }
  }

  async function onCreateReview(e: FormEvent) {
    e.preventDefault()
    setReviewErr(null)
    setReviewMsg(null)
    setBusy(true)
    try {
      const r = await createReview(
        { title: reviewTitle.trim(), body: reviewBody.trim() || undefined, task_id: taskId },
        { token },
      )
      setReviewMsg(`Created review ${r.id}`)
      setReviewBody('')
    } catch (err) {
      setReviewErr(err)
    } finally {
      setBusy(false)
    }
  }

  const envLine = `TRACE_TASK_ID=${taskId}`
  const suggestions = task ? SUGGESTED[task.work_state] ?? [] : []
  const targetingDone = toState === 'DONE'

  return (
    <div>
      <p>
        <Link to="/tasks">← Tasks</Link>
        {' · '}
        <Link to={`/loop?task_id=${encodeURIComponent(taskId)}`}>Loop</Link>
        {' · '}
        <Link to="/reviews">Reviews</Link>
      </p>
      <h1 className="page-title">{task?.title ?? 'Task'}</h1>

      {error ? <ErrorBanner error={error} onRetry={() => void load()} /> : null}

      {task ? (
        <>
          <div className="panel">
            <div className="stack">
              <div>
                <span className="pill">{task.work_state}</span>
              </div>
              <div className="mono">{task.id}</div>
              {task.goal_id ? <div className="mono">goal: {task.goal_id}</div> : null}
              {task.body ? <p>{task.body}</p> : null}
              <div className="row">
                <code className="mono">{envLine}</code>
                <CopyButton text={envLine} label="Copy TRACE_TASK_ID" />
              </div>
            </div>
          </div>

          <div className="panel">
            <h2>DONE gate</h2>
            <p className="page-lead" style={{ marginBottom: '0.75rem' }}>
              Gate strip is advisory chrome — transition denials still come from the API envelope
              only (Law 19).
            </p>
            <GateStrip gate={gate} loading={gateLoading} />
            {Array.isArray(gate?.violations) &&
            gate.violations.some(
              (v) => (v as { reason_code?: string }).reason_code === 'goal_plan_gap_terminal_advisory',
            ) ? (
              <p className="page-lead" style={{ marginTop: '0.75rem' }}>
                Goal lacks a progressive plan; work is already terminal. Bootstrap via{' '}
                <code className="mono">trace plan bootstrap --goal</code> or MCP{' '}
                <code className="mono">trace_plan</code> when ready — not a blocker for finished
                tasks.
              </p>
            ) : null}
            {targetingDone && gate && (gate as { allowed?: boolean }).allowed === false ? (
              <div className="banner banner--warn" role="status">
                Gate for <code className="mono">done</code> is blocked — you can still attempt the
                transition; the server decides.
              </div>
            ) : null}
          </div>

          <div className="panel">
            <h2>Transition</h2>
            <p className="page-lead" style={{ marginBottom: '0.75rem' }}>
              Suggested states are UX only. Warnings from{' '}
              <code className="mono">TransitionResult.warning</code> are shown as returned.
            </p>
            {txError ? <ErrorBanner error={txError} /> : null}
            {txOk ? (
              <div className="banner banner--ok" role="status">
                {txOk}
              </div>
            ) : null}
            <form className="stack" onSubmit={(e) => void onTransition(e)}>
              <div className="field">
                <label htmlFor="to-state">to_state</label>
                <select
                  id="to-state"
                  value={toState}
                  onChange={(e) => setToState(e.target.value)}
                  required
                >
                  <optgroup label="Suggested">
                    {suggestions.map((s) => (
                      <option key={s} value={s}>
                        {s}
                      </option>
                    ))}
                  </optgroup>
                  <optgroup label="All">
                    {ALL_STATES.filter((s) => !suggestions.includes(s)).map((s) => (
                      <option key={s} value={s}>
                        {s}
                      </option>
                    ))}
                  </optgroup>
                </select>
              </div>
              <div className="field">
                <label htmlFor="reason">reason</label>
                <input
                  id="reason"
                  value={reason}
                  onChange={(e) => setReason(e.target.value)}
                  required
                />
              </div>
              {targetingDone ? (
                <label className="row" style={{ gap: '0.35rem' }}>
                  <input
                    type="checkbox"
                    checked={allowDone}
                    onChange={(e) => setAllowDone(e.target.checked)}
                  />
                  pass allow_done to API (server still decides)
                </label>
              ) : null}
              <button type="submit" className="btn btn--primary" disabled={busy || !toState}>
                {busy ? 'Applying…' : 'Transition'}
              </button>
            </form>
          </div>

          <div className="panel">
            <h2>Create review</h2>
            {reviewErr ? <ErrorBanner error={reviewErr} /> : null}
            {reviewMsg ? (
              <div className="banner banner--ok" role="status">
                {reviewMsg} — <Link to="/reviews">View reviews</Link>
              </div>
            ) : null}
            <form className="stack" onSubmit={(e) => void onCreateReview(e)}>
              <div className="field">
                <label htmlFor="review-title">title</label>
                <input
                  id="review-title"
                  value={reviewTitle}
                  onChange={(e) => setReviewTitle(e.target.value)}
                  required
                />
              </div>
              <div className="field">
                <label htmlFor="review-body">body</label>
                <textarea
                  id="review-body"
                  rows={3}
                  value={reviewBody}
                  onChange={(e) => setReviewBody(e.target.value)}
                />
              </div>
              <button type="submit" className="btn" disabled={busy || !reviewTitle.trim()}>
                Create review
              </button>
            </form>
          </div>

          <div className="panel">
            <h2>Bounded packets</h2>
            <div className="row">
              <button type="button" className="btn" onClick={() => void openContext()}>
                Context
              </button>
              <button type="button" className="btn" onClick={() => void openWhy()}>
                Why
              </button>
              {drawer ? (
                <button type="button" className="btn btn--ghost" onClick={() => setDrawer(null)}>
                  Close
                </button>
              ) : null}
            </div>
            {drawer ? (
              <div className="drawer">
                <strong>{drawer}</strong>
                {drawerErr ? <ErrorBanner error={drawerErr} /> : null}
                {drawerData ? <pre>{JSON.stringify(drawerData, null, 2)}</pre> : null}
              </div>
            ) : null}
          </div>
        </>
      ) : null}
    </div>
  )
}
