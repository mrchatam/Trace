import { Link } from 'react-router-dom'
import { dismissOrient } from '../lib/orientDismiss'
import { DEPTH, EXPAND_MAX_NODES, API_MAX_NODES, PROJECT_MAX_NODES } from '../lib/overviewCompose'

type Props = {
  onDismiss?: () => void
  persistent?: boolean
}

export function GraphOrientPanel({ onDismiss, persistent }: Props) {
  function handleDismiss() {
    dismissOrient()
    onDismiss?.()
  }

  return (
    <aside
      className="graph-orient"
      data-testid="graph-orient-panel"
      role="region"
      aria-labelledby="graph-orient-title"
    >
      <div className="graph-orient__header">
        <h2 id="graph-orient-title" className="graph-orient__title">
          Explore — orient first
        </h2>
        {!persistent ? (
          <button
            type="button"
            className="btn btn--ghost graph-orient__dismiss"
            aria-label="Dismiss orient panel"
            data-testid="graph-orient-dismiss"
            onClick={handleDismiss}
          >
            Dismiss
          </button>
        ) : null}
      </div>
      <p className="graph-orient__lead">
        Trace is <strong>moat-first</strong>: pick work on Tasks, run the deliberation loop, pass the
        gate before edits, and review before done. This Explore graph is your{' '}
        <strong>orient entry</strong> — not a replacement for the task loop.
      </p>
      <ol className="graph-orient__steps">
        <li>
          <Link to="/tasks">Tasks</Link> — pick or create work
        </li>
        <li>
          <Link to="/loop">Loop</Link> — load bounded context and deliberate
        </li>
        <li>Gate — required before product edits</li>
        <li>Review — evidence-backed completion</li>
      </ol>
      <p className="graph-orient__budget">
        <span className="graph-orient__confidence" data-testid="graph-orient-confidence">
          Confidence: high within budget
        </span>{' '}
        — Explore opens on the full project graph (<code className="mono">mode=project</code>, default cap{' '}
        {PROJECT_MAX_NODES}, max {API_MAX_NODES}) with explicit <code className="mono">max_nodes</code> (Laws 6–7).
        Double-click or “Use as center” for a bounded neighborhood (depth {DEPTH}, expand ≤
        {EXPAND_MAX_NODES}). Truncation banners show how many entities were omitted.
      </p>
      <section className="graph-orient__legend" aria-labelledby="graph-orient-legend-title">
        <h3 id="graph-orient-legend-title" className="graph-orient__legend-title">
          Graph legend
        </h3>
        <ul className="graph-orient__legend-list">
          <li>
            <strong>Scopes</strong> — nodes cluster by server <code className="mono">scope_id</code>{' '}
            (feature/layer/business cartography). Missing membership stays ungrouped at center; the
            browser does not infer scope from edges.
          </li>
          <li>
            <strong>Inferred edges</strong> — dashed stroke means{' '}
            <code className="mono">provenance=inferred</code> (from CLI{' '}
            <code className="mono">trace graph infer</code>); solid edges are explicit. Explicit
            membership wins over inference.
          </li>
          <li>
            <strong>Edge budget</strong> — overview shows up to 150 priority edges; MVP scope rels (
            <code className="mono">api_contract</code>, etc.) rank above goal→task when present.
          </li>
        </ul>
      </section>
    </aside>
  )
}
