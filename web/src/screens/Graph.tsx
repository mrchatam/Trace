import { useCallback, useEffect, useMemo, useRef, useState, type FormEvent, type MouseEvent } from 'react'
import {
  Background,
  Controls,
  Handle,
  Position,
  ReactFlow,
  ReactFlowProvider,
  useEdgesState,
  useNodesState,
  type Edge,
  type Node,
  type NodeMouseHandler,
  type NodeProps,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import { Link } from 'react-router-dom'
import { ErrorBanner, EmptyState } from '../components/ErrorBanner'
import { GraphOrientPanel } from '../components/GraphOrientPanel'
import { dismissOrient, isOrientDismissed } from '../lib/orientDismiss'
import { Inspector } from '../components/Inspector'
import { useApiToken } from '../context/AppChrome'
import {
  getGraph,
  getProject,
  listTasks,
  search,
  type BoundedGraph,
  type ProjectResponse,
  type SearchItem,
  type TaskRow,
} from '../api/ops'
import {
  DEPTH,
  EXPAND_MAX_NODES,
  SEARCH_FILL_QUERIES,
  SEED_CAP,
  SEED_MAX_NODES,
  SEED_TARGET,
  UI_CAP,
  applySeedMeta,
  collectSearchFill,
  composeSeedsFromParts,
  mergeOverviewGraphs,
  partitionSettledGraphs,
  type GraphNodeMeta,
  type SeedCandidate,
} from '../lib/overviewCompose'

type KindFilter = 'all' | string

type GraphNodeData = {
  label: string
  kind: string
  work_state?: string
  isSeed?: boolean
}

function GraphNodeView({ data, id, selected }: NodeProps & { data: GraphNodeData }) {
  return (
    <div
      className="graph-node__inner"
      data-testid={`graph-canvas-node-${id}`}
      data-kind={data.kind}
      data-state={data.work_state || undefined}
      tabIndex={selected ? 0 : -1}
      role="button"
      aria-pressed={selected}
      aria-label={`${data.kind}: ${data.label}`}
    >
      <Handle type="target" position={Position.Left} style={{ opacity: 0 }} />
      <span className="graph-node__kind">{data.kind}</span>
      <span className="graph-node__label">{data.label}</span>
      {data.work_state ? (
        <span className="graph-node__state" data-state={data.work_state}>
          {data.work_state}
        </span>
      ) : null}
      <Handle type="source" position={Position.Right} style={{ opacity: 0 }} />
    </div>
  )
}

const nodeTypes = { graphNode: GraphNodeView }

function layoutOverview(
  nodesMeta: GraphNodeMeta[],
  edgesMeta: BoundedGraph['edges'],
  seedIds: Set<string>,
  centerId: string,
  selectedId: string | null,
): { nodes: Node[]; edges: Edge[] } {
  const seeds = nodesMeta.filter((n) => seedIds.has(n.id))
  const others = nodesMeta.filter((n) => !seedIds.has(n.id))
  const seedCount = Math.max(seeds.length, 1)
  const nodes: Node[] = []

  seeds.forEach((node, i) => {
    const angle = (2 * Math.PI * i) / seedCount - Math.PI / 2
    const radius = seeds.length === 1 ? 0 : 220
    const classes = ['graph-node', 'graph-node--seed']
    if (node.id === centerId) classes.push('graph-node--center')
    if (selectedId && node.id === selectedId) classes.push('graph-node--selected')
    nodes.push({
      id: node.id,
      type: 'graphNode',
      position: {
        x: 400 + radius * Math.cos(angle),
        y: 280 + radius * Math.sin(angle),
      },
      data: {
        label: node.title,
        kind: node.kind,
        work_state: node.work_state,
        isSeed: true,
      },
      className: classes.join(' '),
      selected: selectedId === node.id,
      focusable: true,
    })
  })

  others.forEach((node, i) => {
    const seedIdx = i % seedCount
    const baseAngle = (2 * Math.PI * seedIdx) / seedCount - Math.PI / 2
    const local = Math.floor(i / seedCount)
    const angle = baseAngle + (local % 5) * 0.35 - 0.7
    const radius = 320 + (local % 4) * 36
    const classes = ['graph-node']
    if (node.id === centerId) classes.push('graph-node--center')
    if (selectedId && node.id === selectedId) classes.push('graph-node--selected')
    if (seedIds.has(node.id)) classes.push('graph-node--seed')
    nodes.push({
      id: node.id,
      type: 'graphNode',
      position: {
        x: 400 + radius * Math.cos(angle),
        y: 280 + radius * Math.sin(angle),
      },
      data: {
        label: node.title,
        kind: node.kind,
        work_state: node.work_state,
        isSeed: seedIds.has(node.id),
      },
      className: classes.join(' '),
      selected: selectedId === node.id,
      focusable: true,
    })
  })

  const edges: Edge[] = (edgesMeta ?? []).map((e, i) => ({
    id: `${e.from}-${e.rel}-${e.to}-${i}`,
    source: e.from,
    target: e.to,
    label: e.rel,
  }))
  return { nodes, edges }
}

function GraphCanvas({
  nodesMeta,
  edgesMeta,
  seedIds,
  center,
  selectedId,
  onSelect,
  onExpand,
}: {
  nodesMeta: GraphNodeMeta[]
  edgesMeta: BoundedGraph['edges']
  seedIds: Set<string>
  center: string
  selectedId: string | null
  onSelect: (id: string) => void
  onExpand: (id: string) => void
}) {
  const laid = useMemo(
    () => layoutOverview(nodesMeta, edgesMeta, seedIds, center, selectedId),
    [nodesMeta, edgesMeta, seedIds, center, selectedId],
  )
  const [nodes, setNodes, onNodesChange] = useNodesState(laid.nodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(laid.edges)

  useEffect(() => {
    setNodes(laid.nodes)
    setEdges(laid.edges)
  }, [laid, setNodes, setEdges])

  const onNodeClick: NodeMouseHandler = useCallback(
    (_: MouseEvent, node: Node) => {
      onSelect(node.id)
    },
    [onSelect],
  )

  const onNodeDoubleClick: NodeMouseHandler = useCallback(
    (_: MouseEvent, node: Node) => {
      onExpand(node.id)
    },
    [onExpand],
  )

  return (
    <div className="graph-canvas" data-testid="graph-canvas">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={onNodeClick}
        onNodeDoubleClick={onNodeDoubleClick}
        fitView
        nodesConnectable={false}
        nodesFocusable
        elementsSelectable
        proOptions={{ hideAttribution: true }}
      >
        <Background />
        <Controls />
      </ReactFlow>
    </div>
  )
}

function toBoundedGraph(
  nodes: GraphNodeMeta[],
  edges: BoundedGraph['edges'],
  center: string,
  maxNodes: number,
  truncated: boolean,
): BoundedGraph {
  return {
    center,
    max_nodes: maxNodes,
    truncated,
    nodes: nodes.map((n) => ({ id: n.id, kind: n.kind, title: n.title })),
    edges: edges ?? [],
  }
}

export function Graph() {
  const token = useApiToken()
  const [center, setCenter] = useState('')
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [maxNodes, setMaxNodes] = useState(EXPAND_MAX_NODES)
  const [tasks, setTasks] = useState<TaskRow[]>([])
  const [searchItems, setSearchItems] = useState<SearchItem[]>([])
  const [q, setQ] = useState('')
  const [kind, setKind] = useState<KindFilter>('all')
  const [nodesMeta, setNodesMeta] = useState<GraphNodeMeta[]>([])
  const [edgesMeta, setEdgesMeta] = useState<BoundedGraph['edges']>([])
  const [graphTruncated, setGraphTruncated] = useState(false)
  const [omitted, setOmitted] = useState(0)
  const [seedIds, setSeedIds] = useState<string[]>([])
  const [project, setProject] = useState<ProjectResponse | null>(null)
  const [error, setError] = useState<unknown>(null)
  const [partialFailed, setPartialFailed] = useState<string[]>([])
  const [overviewLoading, setOverviewLoading] = useState(true)
  const [expandLoading, setExpandLoading] = useState(false)
  const [noSeeds, setNoSeeds] = useState(false)
  const [orientDismissed, setOrientDismissed] = useState(() => isOrientDismissed())
  const seedsRef = useRef<SeedCandidate[]>([])
  const bootstrapped = useRef(false)

  const seedSet = useMemo(() => new Set(seedIds), [seedIds])
  const selectionOutOfView =
    selectedId != null && !nodesMeta.some((n) => n.id === selectedId)

  const graph: BoundedGraph | null = useMemo(() => {
    if (nodesMeta.length === 0) return null
    return toBoundedGraph(nodesMeta, edgesMeta, center, Math.min(UI_CAP, maxNodes), graphTruncated)
  }, [nodesMeta, edgesMeta, center, maxNodes, graphTruncated])

  const filteredSearch = useMemo(() => {
    if (kind === 'all') return searchItems
    return searchItems.filter((i) => i.kind === kind)
  }, [searchItems, kind])

  const kinds = useMemo(() => {
    const s = new Set([
      ...searchItems.map((i) => i.kind),
      ...nodesMeta.map((n) => n.kind),
    ])
    return Array.from(s).sort()
  }, [searchItems, nodesMeta])

  const visibleNodes = useMemo(() => {
    if (kind === 'all') return nodesMeta
    return nodesMeta.filter((n) => n.kind === kind)
  }, [nodesMeta, kind])

  const visibleEdges = useMemo(() => {
    const meta = edgesMeta ?? []
    if (kind === 'all') return meta
    const ids = new Set(visibleNodes.map((n) => n.id))
    return meta.filter((e) => ids.has(e.from) && ids.has(e.to))
  }, [edgesMeta, visibleNodes, kind])

  const applyMerged = useCallback(
    (
      merged: ReturnType<typeof mergeOverviewGraphs>,
      seeds: SeedCandidate[],
      nextCenter?: string,
    ) => {
      const withMeta = applySeedMeta(merged.nodes, seeds)
      setNodesMeta(withMeta)
      setEdgesMeta(merged.edges)
      setGraphTruncated(merged.truncated)
      setOmitted(merged.omitted)
      setSeedIds(seeds.map((s) => s.id))
      seedsRef.current = seeds
      const c = nextCenter ?? merged.center
      setCenter(c)
      setMaxNodes(Math.min(UI_CAP, merged.max_nodes || SEED_MAX_NODES))
    },
    [],
  )

  const loadOverview = useCallback(async () => {
    setOverviewLoading(true)
    setError(null)
    setPartialFailed([])
    setNoSeeds(false)
    setNodesMeta([])
    setEdgesMeta([])
    try {
      const [projSettled, taskRes] = await Promise.all([
        getProject({ token }).catch(() => null),
        listTasks({ limit: 50 }, { token }),
      ])
      if (projSettled) setProject(projSettled)
      const taskItems = taskRes.items ?? []
      setTasks(taskItems)

      const taskSeeds = composeSeedsFromParts(taskItems, [], {
        target: SEED_TARGET,
        cap: SEED_CAP,
      })
      const need = Math.max(0, SEED_TARGET - taskSeeds.length)
      const batches: SearchItem[][] = []
      if (need > 0) {
        for (const qFill of SEARCH_FILL_QUERIES) {
          if (collectSearchFill(batches, new Set(taskSeeds.map((s) => s.id)), need).length >= need) {
            break
          }
          try {
            const res = await search(qFill, 20, { token })
            batches.push(res.items ?? [])
          } catch {
            batches.push([])
          }
        }
      }
      const fillHits = collectSearchFill(
        batches,
        new Set(taskSeeds.map((s) => s.id)),
        Math.max(need, SEED_CAP - taskSeeds.length),
      )
      const seeds = composeSeedsFromParts(taskItems, fillHits, {
        target: SEED_TARGET,
        cap: SEED_CAP,
      })

      if (seeds.length === 0) {
        setNoSeeds(true)
        setOverviewLoading(false)
        return
      }

      const results = await Promise.allSettled(
        seeds.map((s) => getGraph(s.id, SEED_MAX_NODES, DEPTH, { token })),
      )
      const { graphs, failedSeedIds, hardFail } = partitionSettledGraphs(
        results,
        seeds.map((s) => s.id),
      )
      if (hardFail) {
        setError(new Error(`Overview failed for all ${seeds.length} seeds`))
        setPartialFailed(failedSeedIds)
        setOverviewLoading(false)
        return
      }
      if (failedSeedIds.length > 0) setPartialFailed(failedSeedIds)

      const merged = mergeOverviewGraphs(
        graphs,
        seeds.map((s) => s.id),
        { uiCap: UI_CAP, depth: DEPTH },
      )
      applyMerged(merged, seeds)
    } catch (err) {
      setError(err)
      setNodesMeta([])
      setEdgesMeta([])
    } finally {
      setOverviewLoading(false)
    }
  }, [token, applyMerged])

  useEffect(() => {
    if (bootstrapped.current) return
    bootstrapped.current = true
    void loadOverview()
  }, [loadOverview])

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

  async function loadGraph(centerId: string, budget = EXPAND_MAX_NODES) {
    const capped = Math.min(Math.max(1, budget), EXPAND_MAX_NODES, UI_CAP)
    setMaxNodes(capped)
    setExpandLoading(true)
    setError(null)
    try {
      const g = await getGraph(centerId, capped, DEPTH, { token })
      const seeds = seedsRef.current
      const withMeta = applySeedMeta(
        g.nodes.map((n) => ({ id: n.id, kind: n.kind, title: n.title })),
        seeds,
      )
      setNodesMeta(withMeta)
      setEdgesMeta(g.edges ?? [])
      setGraphTruncated(g.truncated)
      setOmitted(0)
      setCenter(centerId)
      setSelectedId(centerId)
      // Prefer replace neighborhood; keep seed pin set for styling
      setSeedIds(seeds.map((s) => s.id))
    } catch (err) {
      setError(err)
    } finally {
      setExpandLoading(false)
    }
  }

  function onSelect(id: string) {
    setSelectedId(id)
  }

  function onExpand(id: string) {
    void loadGraph(id)
  }

  function handleDismissOrient() {
    dismissOrient()
    setOrientDismissed(true)
  }

  const loading = overviewLoading || expandLoading
  const showCanvas = graph != null && visibleNodes.length > 0

  return (
    <div>
      <h1 className="page-title">Graph</h1>
      {!orientDismissed ? <GraphOrientPanel onDismiss={handleDismissOrient} /> : null}
      <p className="page-lead">
        Project overview graph (seed-composed). Click to inspect; double-click or “Use as center” to
        re-center. Seeds ≤{SEED_CAP} · per-seed {SEED_MAX_NODES}/depth {DEPTH} · UI cap {UI_CAP} ·
        expand ≤{EXPAND_MAX_NODES}.
      </p>

      {project ? (
        <p className="mono" data-testid="graph-project-chrome" style={{ fontSize: '0.8rem' }}>
          root={project.root}
          {project.store_ready != null ? ` · store_ready=${String(project.store_ready)}` : ''}
        </p>
      ) : null}

      <div className="banner banner--info" role="note">
        <span className="confidence-label" data-testid="graph-law-banner-confidence">
          Confidence: high — bounded read
        </span>{' '}
        · Always passes <code className="mono">center</code> +{' '}
        <code className="mono">max_nodes</code> — not a full-project dump (Laws 6–7).
      </div>

      <details className="panel" data-testid="graph-manual-center">
        <summary>Manual center (secondary)</summary>
        <form className="filters" onSubmit={(e) => void runSearch(e)}>
          <div className="field" style={{ flex: 1 }}>
            <label htmlFor="graph-q">Search (q)</label>
            <input
              id="graph-q"
              value={q}
              onChange={(e) => setQ(e.target.value)}
              placeholder="search…"
            />
          </div>
          <div className="field">
            <label htmlFor="graph-kind">Kind filter</label>
            <select id="graph-kind" value={kind} onChange={(e) => setKind(e.target.value)}>
              <option value="all">all</option>
              {kinds.map((k) => (
                <option key={k} value={k}>
                  {k}
                </option>
              ))}
            </select>
          </div>
          <div className="field">
            <label htmlFor="graph-budget-input">max_nodes</label>
            <input
              id="graph-budget-input"
              type="number"
              min={1}
              max={UI_CAP}
              value={maxNodes}
              onChange={(e) =>
                setMaxNodes(
                  Math.min(UI_CAP, Math.max(1, Number(e.target.value) || EXPAND_MAX_NODES)),
                )
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
                <button
                  type="button"
                  className="btn btn--ghost"
                  onClick={() => void loadGraph(h.id)}
                >
                  {h.kind}: {h.title}
                </button>
                <span className="mono" style={{ marginLeft: '0.5rem' }}>
                  {h.id}
                </span>
              </li>
            ))}
          </ul>
        ) : null}

        <div style={{ marginTop: '1rem' }}>
          <h3 style={{ fontSize: '0.9rem', margin: '0 0 0.5rem' }}>Or pick a task</h3>
          {tasks.length === 0 ? (
            <EmptyState title="No tasks — open Tasks or search above." />
          ) : (
            <ul className="graph-list" data-testid="graph-task-list">
              {tasks.slice(0, 20).map((t) => (
                <li key={t.id}>
                  <button
                    type="button"
                    className="btn btn--ghost"
                    data-testid={`graph-pick-task-${t.id}`}
                    onClick={() => void loadGraph(t.id)}
                  >
                    {t.title}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      </details>

      {error ? (
        <div>
          <ErrorBanner error={error} />
          <button type="button" className="btn" onClick={() => void loadOverview()}>
            Retry
          </button>
        </div>
      ) : null}

      {partialFailed.length > 0 && !error ? (
        <div className="banner banner--warn" role="status" data-testid="graph-partial-fail">
          Partial overview: {partialFailed.length} seed(s) failed (
          {partialFailed.slice(0, 3).join(', ')}
          {partialFailed.length > 3 ? '…' : ''}). Showing successful neighborhoods.{' '}
          <button type="button" className="btn btn--ghost" onClick={() => void loadOverview()}>
            Retry failed seeds
          </button>
        </div>
      ) : null}

      <div className="panel">
        <h2>Neighborhood</h2>
        {overviewLoading ? (
          <p role="status" data-testid="graph-overview-loading">
            Building project overview…
          </p>
        ) : null}
        {expandLoading && !overviewLoading ? (
          <p role="status">Loading neighborhood…</p>
        ) : null}

        {noSeeds && !overviewLoading ? (
          <EmptyState title="Nothing to seed yet. Add tasks or entities, or pick a center to explore a neighborhood.">
            <Link className="btn btn--primary" to="/tasks" data-testid="explore-empty-open-tasks">
              Open Tasks
            </Link>
          </EmptyState>
        ) : null}

        {!overviewLoading && !noSeeds && !showCanvas && !error && !loading ? (
          <EmptyState title="Overview unavailable — retry or pick a center below." />
        ) : null}

        {showCanvas ? (
          <div className="stack">
            {graphTruncated ? (
              <div className="banner banner--warn" role="status">
                <span className="confidence-label" data-testid="graph-truncation-confidence">
                  Confidence: partial — budget exceeded
                </span>{' '}
                · Truncated: neighborhood exceeded budget
                {omitted > 0 ? ` (+${omitted} omitted)` : ''} (max_nodes≤{UI_CAP}). Budgeted
                neighborhood — not the full project graph.
              </div>
            ) : null}
            <div className="mono" data-testid="graph-budget-line">
              <span className="confidence-label" data-testid="graph-budget-confidence">
                {graphTruncated ? 'Confidence: partial' : 'Confidence: high within budget'}
              </span>{' '}
              · seeds={seedIds.length}/{SEED_CAP} · center=
              <span data-testid="graph-center-id">{center}</span> · nodes={visibleNodes.length}/
              {UI_CAP} · edges={visibleEdges.length} · depth={DEPTH} · max_nodes=
              <span id="graph-budget">{maxNodes}</span>
            </div>
            {selectionOutOfView && selectedId ? (
              <div className="banner banner--info" role="status" data-testid="graph-selection-out">
                Selected entity still in inspector but not in current view.
              </div>
            ) : null}
            <ul
              className="graph-list"
              data-testid="graph-node-list"
              aria-label="Neighborhood nodes"
            >
              {visibleNodes.map((n) => (
                <li key={n.id}>
                  <button
                    type="button"
                    className="btn btn--ghost"
                    data-testid={`graph-select-node-${n.id}`}
                    data-kind={n.kind}
                    data-state={n.work_state || undefined}
                    aria-pressed={selectedId === n.id}
                    onClick={() => onSelect(n.id)}
                  >
                    {n.kind}: {n.title}
                    {n.work_state ? ` (${n.work_state})` : ''}
                  </button>
                  {n.id !== center ? (
                    <button
                      type="button"
                      className="btn btn--ghost"
                      data-testid={`graph-expand-node-${n.id}`}
                      onClick={() => onExpand(n.id)}
                    >
                      Use as center
                    </button>
                  ) : (
                    <span className="pill" style={{ marginLeft: '0.35rem' }}>
                      center
                    </span>
                  )}
                  {seedSet.has(n.id) ? (
                    <span className="pill" style={{ marginLeft: '0.35rem' }}>
                      seed
                    </span>
                  ) : null}
                </li>
              ))}
            </ul>
            <div className="graph-shell">
              <ReactFlowProvider>
                <GraphCanvas
                  nodesMeta={visibleNodes}
                  edgesMeta={visibleEdges}
                  seedIds={seedSet}
                  center={center}
                  selectedId={selectedId}
                  onSelect={onSelect}
                  onExpand={onExpand}
                />
              </ReactFlowProvider>
              {selectedId ? (
                <Inspector
                  selectedId={selectedId}
                  center={center}
                  graph={graph}
                  token={token}
                  onUseAsCenter={onExpand}
                  onClose={() => setSelectedId(null)}
                />
              ) : (
                <div className="graph-inspector graph-inspector--empty" role="status">
                  <p className="graph-inspector__empty-title">No node selected</p>
                  <p className="graph-inspector__empty-lead">
                    Click a node on the canvas or in the list to inspect.
                  </p>
                </div>
              )}
            </div>
          </div>
        ) : null}
      </div>

      <p>
        <Link to="/tasks">Tasks</Link>
        {' · '}
        <Link to="/overview">Overview</Link>
      </p>
    </div>
  )
}
