import { useEffect, useMemo, useState, type ReactNode } from 'react'
import { Link } from 'react-router-dom'
import { ErrorBanner } from './ErrorBanner'
import { GateStrip } from './GateStrip'
import {
  getContext,
  getEntity,
  getImpact,
  getLoopGate,
  getLoopStatus,
  getTask,
  getWhy,
  listReviews,
  type BoundedGraph,
  type ContextPacket,
  type EntitySummary,
  type LoopGateResponse,
  type LoopStatusResponse,
  type ReviewSummary,
  type TaskRow,
  type WhyPacket,
} from '../api/ops'

type SectionState<T> = {
  loading: boolean
  error: unknown
  data: T | null
}

function emptySection<T>(): SectionState<T> {
  return { loading: false, error: null, data: null }
}

function formatPacketValue(v: unknown): string {
  if (typeof v === 'string' || typeof v === 'number' || typeof v === 'boolean') {
    const s = String(v)
    return s.length > 220 ? `${s.slice(0, 220)}…` : s
  }
  if (v == null) return '—'
  if (Array.isArray(v)) {
    if (v.length === 0) return '0 item(s)'
    const first = v[0]
    if (typeof first === 'string' || typeof first === 'number') {
      const preview = v
        .slice(0, 3)
        .map((x) => String(x))
        .join(', ')
      return v.length > 3 ? `${v.length} item(s): ${preview}…` : `${v.length} item(s): ${preview}`
    }
    return `${v.length} item(s)`
  }
  if (typeof v === 'object') {
    const keys = Object.keys(v as object)
    if (keys.length === 0) return '{ }'
    const preview = keys.slice(0, 4).join(', ')
    return keys.length > 4 ? `{${keys.length} keys: ${preview}…}` : `{${keys.length} keys: ${preview}}`
  }
  return String(v).slice(0, 160)
}

function PacketView({ data, label }: { data: unknown; label: string }) {
  if (data == null) return null
  const entries =
    typeof data === 'object' && data !== null && !Array.isArray(data)
      ? Object.entries(data as Record<string, unknown>).slice(0, 24)
      : null

  return (
    <div className="inspector-packet">
      {entries && entries.length > 0 ? (
        <dl className="inspector-dl">
          {entries.map(([k, v]) => (
            <div key={k} className="inspector-dl__row">
              <dt>{k}</dt>
              <dd>{formatPacketValue(v)}</dd>
            </div>
          ))}
        </dl>
      ) : (
        <details>
          <summary>Structured view unavailable — show raw</summary>
          <pre className="mono inspector-pre">{JSON.stringify(data, null, 2).slice(0, 1200)}</pre>
        </details>
      )}
      <details>
        <summary>Raw {label}</summary>
        <pre className="mono inspector-pre">{JSON.stringify(data, null, 2).slice(0, 4000)}</pre>
      </details>
    </div>
  )
}

function Section({
  title,
  testId,
  loading,
  error,
  onRetry,
  children,
  collapsedNote,
}: {
  title: string
  testId: string
  loading?: boolean
  error?: unknown
  onRetry?: () => void
  children?: ReactNode
  collapsedNote?: string
}) {
  return (
    <section className="inspector-section" data-section={testId} aria-labelledby={`insp-${testId}`}>
      <h3 id={`insp-${testId}`} className="inspector-section__title">
        {title}
      </h3>
      {collapsedNote ? <p className="inspector-section__note">{collapsedNote}</p> : null}
      {loading ? (
        <p role="status" className="inspector-section__status">
          Loading…
        </p>
      ) : null}
      {error ? <ErrorBanner error={error} onRetry={onRetry} /> : null}
      {!loading && !error ? children : null}
    </section>
  )
}

type Props = {
  selectedId: string
  center: string
  graph: BoundedGraph | null
  token?: string
  onUseAsCenter: (id: string) => void
  onClose?: () => void
}

export function Inspector({ selectedId, center, graph, token, onUseAsCenter, onClose }: Props) {
  const [entity, setEntity] = useState<SectionState<EntitySummary>>(emptySection())
  const [task, setTask] = useState<SectionState<TaskRow>>(emptySection())
  const [why, setWhy] = useState<SectionState<WhyPacket>>(emptySection())
  const [context, setContext] = useState<SectionState<ContextPacket>>(emptySection())
  const [contextDepth, setContextDepth] = useState<1 | 2>(1)
  const [impact, setImpact] = useState<SectionState<Record<string, unknown>>>(emptySection())
  const [reviews, setReviews] = useState<SectionState<ReviewSummary[]>>(emptySection())
  const [loopStatus, setLoopStatus] = useState<SectionState<LoopStatusResponse>>(emptySection())
  const [loopGate, setLoopGate] = useState<SectionState<LoopGateResponse>>(emptySection())
  const [reloadKey, setReloadKey] = useState(0)

  const nodeKind = useMemo(() => {
    const fromGraph = graph?.nodes.find((n) => n.id === selectedId)?.kind
    return fromGraph ?? entity.data?.kind ?? null
  }, [graph, selectedId, entity.data?.kind])

  const isTask = nodeKind === 'task' || task.data != null

  const neighborLinks = useMemo(() => {
    if (!graph) return []
    const byId = new Map(graph.nodes.map((n) => [n.id, n]))
    return (graph.edges ?? [])
      .filter((e) => e.from === selectedId || e.to === selectedId)
      .map((e, i) => {
        const otherId = e.from === selectedId ? e.to : e.from
        const other = byId.get(otherId)
        return {
          key: `${e.from}-${e.rel}-${e.to}-${i}`,
          rel: e.rel,
          otherId,
          otherTitle: other?.title ?? otherId,
          otherKind: other?.kind ?? 'entity',
          direction: e.from === selectedId ? 'out' : 'in',
        }
      })
  }, [graph, selectedId])

  useEffect(() => {
    let cancelled = false
    const ac = new AbortController()
    const opt = { token, signal: ac.signal }
    const aborted = () => cancelled || ac.signal.aborted

    setEntity({ loading: true, error: null, data: null })
    setTask(emptySection())
    setWhy({ loading: true, error: null, data: null })
    setContext(emptySection())
    setImpact(emptySection())
    setReviews(emptySection())
    setLoopStatus(emptySection())
    setLoopGate(emptySection())

    void (async () => {
      let kindHint = graph?.nodes.find((n) => n.id === selectedId)?.kind ?? null

      try {
        const entityRow = await getEntity(selectedId, opt)
        if (aborted()) return
        setEntity({ loading: false, error: null, data: entityRow })
        kindHint = entityRow.kind
      } catch (err) {
        if (aborted()) return
        setEntity({ loading: false, error: err, data: null })
      }

      const treatAsTask = kindHint === 'task'
      if (treatAsTask) {
        setTask({ loading: true, error: null, data: null })
        try {
          const t = await getTask(selectedId, opt)
          if (aborted()) return
          setTask({ loading: false, error: null, data: t })
        } catch (err) {
          if (aborted()) return
          setTask({ loading: false, error: err, data: null })
        }
      }

      const entityType = kindHint ?? 'entity'
      try {
        const w = await getWhy(entityType, selectedId, opt)
        if (aborted()) return
        setWhy({ loading: false, error: null, data: w })
      } catch (err) {
        if (aborted()) return
        setWhy({ loading: false, error: err, data: null })
      }

      if (!treatAsTask) return

      setContext({ loading: true, error: null, data: null })
      setImpact({ loading: true, error: null, data: null })
      setReviews({ loading: true, error: null, data: null })
      setLoopStatus({ loading: true, error: null, data: null })
      setLoopGate({ loading: true, error: null, data: null })

      await Promise.all([
        getContext(selectedId, contextDepth, opt)
          .then((c) => {
            if (!aborted()) setContext({ loading: false, error: null, data: c })
          })
          .catch((err) => {
            if (!aborted()) setContext({ loading: false, error: err, data: null })
          }),
        getImpact(selectedId, opt)
          .then((imp) => {
            if (!aborted()) setImpact({ loading: false, error: null, data: imp })
          })
          .catch((err) => {
            if (!aborted()) setImpact({ loading: false, error: err, data: null })
          }),
        listReviews({ task_id: selectedId }, opt)
          .then((r) => {
            if (!aborted()) setReviews({ loading: false, error: null, data: r.items ?? [] })
          })
          .catch((err) => {
            if (!aborted()) setReviews({ loading: false, error: err, data: null })
          }),
        getLoopStatus(selectedId, undefined, opt)
          .then((s) => {
            if (!aborted()) setLoopStatus({ loading: false, error: null, data: s })
          })
          .catch((err) => {
            if (!aborted()) setLoopStatus({ loading: false, error: err, data: null })
          }),
        getLoopGate(selectedId, 'edit', opt)
          .then((g) => {
            if (!aborted()) setLoopGate({ loading: false, error: null, data: g })
          })
          .catch((err) => {
            if (!aborted()) setLoopGate({ loading: false, error: err, data: null })
          }),
      ])
    })()

    return () => {
      cancelled = true
      ac.abort()
    }
  }, [selectedId, token, graph, reloadKey, contextDepth])

  const summaryTitle =
    task.data?.title ?? entity.data?.title ?? graph?.nodes.find((n) => n.id === selectedId)?.title
  const summaryKind = nodeKind ?? '—'
  const workState = task.data?.work_state ?? entity.data?.work_state

  return (
    <aside
      className="graph-inspector"
      data-testid="graph-inspector"
      aria-label="Node inspector"
    >
      <div className="graph-inspector__chrome row">
        <strong>Inspector</strong>
        <div className="row" style={{ marginLeft: 'auto', gap: '0.35rem' }}>
          {selectedId !== center ? (
            <button
              type="button"
              className="btn btn--primary"
              data-testid="inspector-use-center"
              onClick={() => onUseAsCenter(selectedId)}
            >
              Use as center
            </button>
          ) : (
            <span className="pill">Center</span>
          )}
          {onClose ? (
            <button type="button" className="btn btn--ghost" onClick={onClose} aria-label="Close inspector">
              Close
            </button>
          ) : null}
        </div>
      </div>

      <div className="graph-inspector__body" key={selectedId}>
      <Section
        title="Summary"
        testId="summary"
        loading={entity.loading || (isTask && task.loading)}
        error={entity.error ?? (isTask ? task.error : null)}
        onRetry={() => setReloadKey((k) => k + 1)}
      >
        <div className="stack" style={{ gap: '0.35rem' }}>
          <div className="row">
            <span className="pill">{summaryKind}</span>
            {workState ? <span className="pill">{workState}</span> : null}
          </div>
          <div>{summaryTitle ?? 'Untitled'}</div>
          <div className="mono inspector-summary-id" data-testid="inspector-selected-id">
            {selectedId}
          </div>
          {isTask ? (
            <Link to={`/tasks/${encodeURIComponent(selectedId)}`}>Open task</Link>
          ) : null}
        </div>
      </Section>

      <Section
        title="Why"
        testId="why"
        loading={why.loading}
        error={why.error}
        onRetry={() => setReloadKey((k) => k + 1)}
      >
        <PacketView data={why.data} label="why" />
      </Section>

      {isTask ? (
        <Section
          title="Context"
          testId="context"
          loading={context.loading}
          error={context.error}
          onRetry={() => setReloadKey((k) => k + 1)}
        >
          <div className="field" style={{ marginBottom: '0.75rem' }}>
            <label htmlFor="inspector-context-depth">Depth</label>
            <select
              id="inspector-context-depth"
              value={contextDepth}
              onChange={(e) => setContextDepth(Number(e.target.value) === 2 ? 2 : 1)}
            >
              <option value={1}>1</option>
              <option value={2}>2</option>
            </select>
          </div>
          <PacketView data={context.data} label="context" />
        </Section>
      ) : (
        <Section title="Context" testId="context" collapsedNote="Task only — omitted for this selection." />
      )}

      {isTask ? (
        <Section
          title="Impact"
          testId="impact"
          loading={impact.loading}
          error={impact.error}
          onRetry={() => setReloadKey((k) => k + 1)}
        >
          {impact.data && Object.keys(impact.data).length === 0 ? (
            <p className="inspector-section__note">No impact rows for this task.</p>
          ) : (
            <PacketView data={impact.data} label="impact" />
          )}
        </Section>
      ) : (
        <Section title="Impact" testId="impact" collapsedNote="Task only — omitted for this selection." />
      )}

      {isTask ? (
        <Section
          title="Reviews"
          testId="reviews"
          loading={reviews.loading}
          error={reviews.error}
          onRetry={() => setReloadKey((k) => k + 1)}
        >
          {(reviews.data?.length ?? 0) === 0 ? (
            <p className="inspector-section__note">No reviews for this task.</p>
          ) : (
            <ul className="graph-list">
              {reviews.data!.map((r) => (
                <li key={r.id}>
                  <Link to={`/reviews/${encodeURIComponent(r.id)}`}>{r.title}</Link>
                  {r.result ? <span className="pill" style={{ marginLeft: '0.5rem' }}>{r.result}</span> : null}
                </li>
              ))}
            </ul>
          )}
        </Section>
      ) : (
        <Section title="Reviews" testId="reviews" collapsedNote="Task only — omitted for this selection." />
      )}

      <Section title="Links" testId="links">
        {neighborLinks.length === 0 ? (
          <p className="inspector-section__note">No edges in the current budgeted neighborhood.</p>
        ) : (
          <ul className="graph-list">
            {neighborLinks.map((l) => (
              <li key={l.key}>
                <span className="mono">{l.direction === 'out' ? '→' : '←'}</span>{' '}
                <span className="pill">{l.rel}</span>{' '}
                {l.otherKind === 'task' ? (
                  <Link to={`/tasks/${encodeURIComponent(l.otherId)}`}>
                    {l.otherKind}: {l.otherTitle}
                  </Link>
                ) : (
                  <span>
                    {l.otherKind}: {l.otherTitle}
                  </span>
                )}
                <span className="mono" style={{ marginLeft: '0.35rem' }}>
                  {l.otherId}
                </span>
              </li>
            ))}
          </ul>
        )}
      </Section>

      {isTask ? (
        <section className="inspector-section" data-section="loop" aria-labelledby="insp-loop">
          <h3 id="insp-loop" className="inspector-section__title">
            Loop
          </h3>
          <GateStrip gate={loopGate.data} loading={loopGate.loading} />
          {loopStatus.error || loopGate.error ? (
            <ErrorBanner error={loopStatus.error ?? loopGate.error} onRetry={() => setReloadKey((k) => k + 1)} />
          ) : null}
          {loopStatus.data ? <PacketView data={loopStatus.data} label="loop status" /> : null}
          <p>
            <Link to={`/loop?task_id=${encodeURIComponent(selectedId)}`}>Open Loop</Link>
          </p>
        </section>
      ) : (
        <Section title="Loop" testId="loop" collapsedNote="Task only — omitted for this selection." />
      )}
      </div>
    </aside>
  )
}
