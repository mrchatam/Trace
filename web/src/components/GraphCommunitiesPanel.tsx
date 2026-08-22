import { kindCssKey } from '../lib/graphLayout'

type GraphCommunitiesPanelProps = {
  counts: Map<string, number>
  enabledKinds: Set<string>
  onToggle: (kind: string) => void
  onSelectAll: () => void
  onClearAll: () => void
}

export function GraphCommunitiesPanel({
  counts,
  enabledKinds,
  onToggle,
  onSelectAll,
  onClearAll,
}: GraphCommunitiesPanelProps) {
  const kinds = Array.from(counts.keys()).sort()

  return (
    <aside className="graph-communities" data-testid="graph-communities-panel" aria-label="Communities">
      <div className="graph-communities__header">
        <h3 className="graph-communities__title">Communities</h3>
        <div className="graph-communities__actions">
          <button
            type="button"
            className="btn btn--ghost graph-communities__action"
            data-testid="graph-communities-all"
            onClick={onSelectAll}
          >
            All
          </button>
          <button
            type="button"
            className="btn btn--ghost graph-communities__action"
            data-testid="graph-communities-none"
            onClick={onClearAll}
          >
            None
          </button>
        </div>
      </div>
      <ul className="graph-communities__list">
        {kinds.map((kind) => {
          const count = counts.get(kind) ?? 0
          const inputId = `graph-community-${kind}`
          return (
            <li key={kind}>
              <label className="graph-communities__row" htmlFor={inputId}>
                <input
                  id={inputId}
                  type="checkbox"
                  checked={enabledKinds.has(kind)}
                  data-testid={`graph-community-${kind}`}
                  onChange={() => onToggle(kind)}
                />
                <span className="graph-communities__dot" data-kind={kindCssKey(kind)} aria-hidden />
                <span className="graph-communities__kind">{kind}</span>
                <span className="graph-communities__count">{count}</span>
              </label>
            </li>
          )
        })}
      </ul>
    </aside>
  )
}
