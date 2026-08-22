import { Link } from 'react-router-dom'
import { dismissOrient } from '../lib/orientDismiss'
import { DEPTH, EXPAND_MAX_NODES, SEED_CAP, SEED_MAX_NODES, UI_CAP } from '../lib/overviewCompose'

type Props = {
  onDismiss: () => void
}

export function GraphOrientPanel({ onDismiss }: Props) {
  function handleDismiss() {
    dismissOrient()
    onDismiss()
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
        <button
          type="button"
          className="btn btn--ghost graph-orient__dismiss"
          aria-label="Dismiss orient panel"
          data-testid="graph-orient-dismiss"
          onClick={handleDismiss}
        >
          Dismiss
        </button>
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
        — every graph read passes <code className="mono">center</code> +{' '}
        <code className="mono">max_nodes</code> (Laws 6–7). Seeds ≤{SEED_CAP} · per-seed{' '}
        {SEED_MAX_NODES}/depth {DEPTH} · UI cap {UI_CAP} · expand ≤{EXPAND_MAX_NODES}. Truncation
        banners mean partial neighborhoods — not the full project.
      </p>
    </aside>
  )
}
