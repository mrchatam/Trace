/**
 * Force-directed layout + graph view policy (Laws 6–7 unchanged).
 * Pure helpers — unit-testable without React.
 */

import {
  forceCenter,
  forceCollide,
  forceLink,
  forceManyBody,
  forceSimulation,
  type SimulationNodeDatum,
} from 'd3-force'

/** Maps-style LOD zoom thresholds (overview → detail). */
export const COMPACT_NODE_THRESHOLD = 30
export const EDGE_LABEL_MAX = 60
export const EDGE_LABEL_MIN_ZOOM = 0.55
/** Below: 6px dots, no labels, sampled edges. */
export const LOD_MINIMAL_MAX_ZOOM = 0.35
/** Below: colored dots only (no card, no label). */
export const LOD_DOT_MAX_ZOOM = 0.75
/** At/above: full node cards for all visible nodes (Maps-style: focus-only until very high zoom). */
export const FULL_CARD_MIN_ZOOM = 1.25
/** @deprecated Use FULL_CARD_MIN_ZOOM */
export const FULL_NODE_MIN_ZOOM = FULL_CARD_MIN_ZOOM
/** Below: hide edges entirely (too dense at overview). */
export const EDGE_HIDE_MAX_ZOOM = 0.25
/** Below + many edges: render a sampled subset. */
export const EDGE_SAMPLE_MIN_ZOOM = 0.45
export const EDGE_SAMPLE_MIN_COUNT = 100
export const EDGE_SAMPLE_RATIO = 3
export const PROJECT_FIT_PADDING = 0.2

export type NodeLod = 'minimal' | 'dot' | 'full'

export type LayoutEdge = { from: string; to: string }

/** Node metadata for semantic (kind-lane) layout. */
export type SemanticLayoutNode = {
  id: string
  kind: string
  goal_id?: string | null
  work_state?: string
  title?: string
}

/** Kind swimlane order (aligned with backend projectGraphKindOrder). */
export const KIND_LANE_ORDER = [
  'goal',
  'task',
  'decision',
  'assumption',
  'discovery',
  'plan_change',
  'claim',
  'evidence',
  'review',
  'capability',
  'change',
  'regression',
] as const

export const SEMANTIC_LAYOUT_LABEL = 'by kind'

const UNGROUPED_GOAL = '__ungrouped__'

const TASK_STATUS_RANK: Record<string, number> = {
  PENDING: 0,
  IN_PROGRESS: 1,
  BLOCKED: 2,
  FAIL: 3,
  DONE: 4,
  SKIPPED: 5,
}

/** Horizontal distance between kind-lane centers (px). */
export const SEMANTIC_LANE_WIDTH = 280
/** Vertical gap between nodes in the same goal/kind cell (px). */
export const SEMANTIC_NODE_SPACING = 52
export const SEMANTIC_GOAL_GAP = 48
export const SEMANTIC_PADDING = 80
/** Max nodes stacked in one cell before wrapping to an overflow sub-column. */
export const SEMANTIC_MAX_NODES_PER_CELL = 24
const SEMANTIC_OVERFLOW_COLUMN_OFFSET = 72
/** Minimum distance between distinct node positions (px). */
export const SEMANTIC_MIN_NODE_DISTANCE = 20

/** Map API kind (plan_change) to CSS data-kind token (plan-change). */
export function kindCssKey(kind: string): string {
  return kind.replace(/_/g, '-')
}

function kindLaneIndex(kind: string): number {
  const idx = KIND_LANE_ORDER.indexOf(kind as (typeof KIND_LANE_ORDER)[number])
  return idx >= 0 ? idx : KIND_LANE_ORDER.length
}

function taskStatusRank(workState: string | undefined): number {
  if (!workState) return 99
  return TASK_STATUS_RANK[workState] ?? 50
}

function compareNodesInCell(a: SemanticLayoutNode, b: SemanticLayoutNode): number {
  if (a.kind === 'task' && b.kind === 'task') {
    const sr = taskStatusRank(a.work_state) - taskStatusRank(b.work_state)
    if (sr !== 0) return sr
  }
  const ta = (a.title ?? a.id).toLowerCase()
  const tb = (b.title ?? b.id).toLowerCase()
  if (ta !== tb) return ta.localeCompare(tb)
  return a.id.localeCompare(b.id)
}

/** Assign each node to a goal band via direct goal_id and BFS from goals. */
export function inferGoalBands(
  nodes: readonly SemanticLayoutNode[],
  edges: readonly LayoutEdge[],
): Map<string, string> {
  const byId = new Map(nodes.map((n) => [n.id, n]))
  const assigned = new Map<string, string>()

  for (const n of nodes) {
    if (n.kind === 'goal') assigned.set(n.id, n.id)
    else if (n.goal_id) assigned.set(n.id, n.goal_id)
  }

  const adj = new Map<string, string[]>()
  for (const e of edges) {
    if (!byId.has(e.from) || !byId.has(e.to)) continue
    const fromList = adj.get(e.from) ?? []
    fromList.push(e.to)
    adj.set(e.from, fromList)
    const toList = adj.get(e.to) ?? []
    toList.push(e.from)
    adj.set(e.to, toList)
  }

  const goals = nodes
    .filter((n) => n.kind === 'goal')
    .sort((a, b) => (a.title ?? a.id).localeCompare(b.title ?? b.id))

  for (const g of goals) {
    const queue = [g.id]
    const visited = new Set<string>()
    while (queue.length > 0) {
      const cur = queue.shift()!
      if (visited.has(cur)) continue
      visited.add(cur)
      if (cur !== g.id && !assigned.has(cur)) assigned.set(cur, g.id)
      for (const nb of adj.get(cur) ?? []) {
        if (!visited.has(nb)) queue.push(nb)
      }
    }
  }

  for (const n of nodes) {
    if (!assigned.has(n.id)) assigned.set(n.id, UNGROUPED_GOAL)
  }

  return assigned
}

function cellStackPosition(
  baseX: number,
  baseY: number,
  indexInCell: number,
): { x: number; y: number } {
  const overflowCol = Math.floor(indexInCell / SEMANTIC_MAX_NODES_PER_CELL)
  const rowInCol = indexInCell % SEMANTIC_MAX_NODES_PER_CELL
  return {
    x: baseX + overflowCol * SEMANTIC_OVERFLOW_COLUMN_OFFSET,
    y: baseY + rowInCol * SEMANTIC_NODE_SPACING,
  }
}

/** Deterministic swimlane layout: X = entity kind, Y = goal cluster (+ status within task lane). */
export function computeSemanticLayout(
  nodes: readonly SemanticLayoutNode[],
  edges: readonly LayoutEdge[],
  _opts: { width?: number; height?: number } = {},
): Map<string, { x: number; y: number }> {
  const positions = new Map<string, { x: number; y: number }>()
  if (nodes.length === 0) return positions

  const goalBands = inferGoalBands(nodes, edges)
  const bandKeys = [...new Set(goalBands.values())].sort((a, b) => {
    if (a === UNGROUPED_GOAL) return 1
    if (b === UNGROUPED_GOAL) return -1
    const na = nodes.find((n) => n.id === a)
    const nb = nodes.find((n) => n.id === b)
    return (na?.title ?? a).localeCompare(nb?.title ?? b)
  })

  const bandIndex = new Map(bandKeys.map((k, i) => [k, i]))
  const cells = new Map<string, SemanticLayoutNode[]>()

  for (const n of nodes) {
    const band = goalBands.get(n.id) ?? UNGROUPED_GOAL
    const lane = kindLaneIndex(n.kind)
    const key = `${band}\x00${lane}`
    const list = cells.get(key) ?? []
    list.push(n)
    cells.set(key, list)
  }

  for (const list of cells.values()) {
    list.sort(compareNodesInCell)
  }

  const bandHeights = bandKeys.map((band) => {
    let maxRows = 1
    for (const [key, list] of cells) {
      if (!key.startsWith(`${band}\x00`)) continue
      const rows = Math.min(list.length, SEMANTIC_MAX_NODES_PER_CELL)
      if (rows > maxRows) maxRows = rows
    }
    return maxRows * SEMANTIC_NODE_SPACING + SEMANTIC_GOAL_GAP
  })

  const bandBaseY: number[] = []
  let yCursor = SEMANTIC_PADDING
  for (let i = 0; i < bandKeys.length; i++) {
    bandBaseY.push(yCursor)
    yCursor += bandHeights[i]!
  }

  for (const n of nodes) {
    const band = goalBands.get(n.id) ?? UNGROUPED_GOAL
    const lane = kindLaneIndex(n.kind)
    const key = `${band}\x00${lane}`
    const list = cells.get(key) ?? []
    const idx = list.findIndex((item) => item.id === n.id)
    const bi = bandIndex.get(band) ?? 0
    const baseX = SEMANTIC_PADDING + lane * SEMANTIC_LANE_WIDTH
    positions.set(n.id, cellStackPosition(baseX, bandBaseY[bi]!, idx))
  }

  // No global scale-to-fit — React Flow fitView handles large graphs at overview zoom.
  return positions
}

type SimNode = SimulationNodeDatum & { id: string }

/** Synchronous d3-force layout for project overview graphs. */
export function computeForceLayout(
  nodeIds: readonly string[],
  edges: readonly LayoutEdge[],
  opts: { width?: number; height?: number; iterations?: number } = {},
): Map<string, { x: number; y: number }> {
  const width = opts.width ?? 900
  const height = opts.height ?? 700
  const iterations = opts.iterations ?? 280
  const idSet = new Set(nodeIds)

  const nodes: SimNode[] = nodeIds.map((id, i) => {
    const angle = (2 * Math.PI * i) / Math.max(nodeIds.length, 1)
    return {
      id,
      x: width / 2 + Math.cos(angle) * 48,
      y: height / 2 + Math.sin(angle) * 48,
    }
  })

  const links = edges
    .filter((e) => idSet.has(e.from) && idSet.has(e.to))
    .map((e) => ({ source: e.from, target: e.to }))

  const simulation = forceSimulation(nodes)
    .force(
      'link',
      forceLink(links)
        .id((d) => (d as SimNode).id)
        .distance(96)
        .strength(0.38),
    )
    .force('charge', forceManyBody().strength(-180).distanceMax(480))
    .force('center', forceCenter(width / 2, height / 2).strength(0.1))
    .force('collide', forceCollide(32))

  simulation.stop()
  for (let i = 0; i < iterations; i++) simulation.tick()

  return new Map(nodes.map((n) => [n.id, { x: n.x ?? width / 2, y: n.y ?? height / 2 }]))
}

export function countNodesByKind(nodes: readonly { kind: string }[]): Map<string, number> {
  const counts = new Map<string, number>()
  for (const n of nodes) {
    counts.set(n.kind, (counts.get(n.kind) ?? 0) + 1)
  }
  return counts
}

export function filterNodesByKinds<T extends { kind: string }>(
  nodes: readonly T[],
  enabledKinds: ReadonlySet<string>,
): T[] {
  return nodes.filter((n) => enabledKinds.has(n.kind))
}

export function filterEdgesByNodeIds<T extends { from: string; to: string }>(
  edges: readonly T[],
  nodeIds: ReadonlySet<string>,
): T[] {
  return edges.filter((e) => nodeIds.has(e.from) && nodeIds.has(e.to))
}

export function shouldUseCompactNodes(
  nodeCount: number,
  layoutMode: 'project' | 'neighborhood',
): boolean {
  return layoutMode === 'project' && nodeCount >= COMPACT_NODE_THRESHOLD
}

/** Resolve node detail level from zoom + focus (Maps-style LOD). */
export function getNodeLod(
  compactOverview: boolean,
  zoom: number,
  nodeId: string,
  selectedId: string | null,
  centerId: string,
  hoveredId: string | null = null,
): NodeLod {
  const focused =
    nodeId === selectedId || nodeId === centerId || nodeId === hoveredId
  if (focused) return 'full'
  if (compactOverview && zoom < FULL_CARD_MIN_ZOOM) {
    if (zoom < LOD_MINIMAL_MAX_ZOOM) return 'minimal'
    return 'dot'
  }
  if (zoom >= FULL_CARD_MIN_ZOOM) return 'full'
  if (zoom >= LOD_MINIMAL_MAX_ZOOM) return 'dot'
  return 'minimal'
}

export function shouldShowNodeLabel(lod: NodeLod): boolean {
  return lod === 'full' || lod === 'dot'
}

/** One-line label for dot LOD (keeps dense clusters readable). */
export function truncateNodeLabel(label: string, maxLen = 28): string {
  const trimmed = label.trim()
  if (trimmed.length <= maxLen) return trimmed
  return `${trimmed.slice(0, maxLen - 1)}…`
}

export function shouldShowFullNode(
  compactOverview: boolean,
  zoom: number,
  nodeId: string,
  selectedId: string | null,
  centerId: string,
  hoveredId: string | null = null,
): boolean {
  return getNodeLod(compactOverview, zoom, nodeId, selectedId, centerId, hoveredId) === 'full'
}

export function shouldRenderEdges(zoom: number): boolean {
  return zoom >= EDGE_HIDE_MAX_ZOOM
}

/** Thin edges at mid zoom; hide entirely when zoomed out past EDGE_HIDE_MAX_ZOOM. */
export function edgeStrokeOpacity(zoom: number, edgeCount: number): number {
  if (zoom < EDGE_HIDE_MAX_ZOOM) return 0
  if (edgeCount > EDGE_SAMPLE_MIN_COUNT && zoom < EDGE_SAMPLE_MIN_ZOOM) return 0.18
  if (zoom < LOD_MINIMAL_MAX_ZOOM) return 0.28
  return 0.42
}

/** Sample edges at low zoom when count is high; always keep focus-adjacent edges. */
export function filterEdgesForLod<T extends { from: string; to: string }>(
  edges: readonly T[],
  zoom: number,
  focusIds: ReadonlySet<string>,
): T[] {
  if (!shouldRenderEdges(zoom)) return []
  if (zoom >= EDGE_SAMPLE_MIN_ZOOM || edges.length <= EDGE_SAMPLE_MIN_COUNT) {
    return [...edges]
  }
  return edges.filter(
    (e, i) => focusIds.has(e.from) || focusIds.has(e.to) || i % EDGE_SAMPLE_RATIO === 0,
  )
}

export function shouldShowEdgeLabels(edgeCount: number, zoom: number): boolean {
  if (edgeCount > EDGE_LABEL_MAX) return false
  return zoom >= EDGE_LABEL_MIN_ZOOM
}

/** Fixed radial layout for neighborhood overview (no force sim). */
export function computeOverviewPositions(
  nodeIds: readonly string[],
  seedIds: ReadonlySet<string>,
): Map<string, { x: number; y: number }> {
  const seeds = nodeIds.filter((id) => seedIds.has(id))
  const others = nodeIds.filter((id) => !seedIds.has(id))
  const seedCount = Math.max(seeds.length, 1)
  const positions = new Map<string, { x: number; y: number }>()

  seeds.forEach((id, i) => {
    const angle = (2 * Math.PI * i) / seedCount - Math.PI / 2
    const radius = seeds.length === 1 ? 0 : 220
    positions.set(id, {
      x: 400 + radius * Math.cos(angle),
      y: 280 + radius * Math.sin(angle),
    })
  })

  others.forEach((id, i) => {
    const seedIdx = i % seedCount
    const baseAngle = (2 * Math.PI * seedIdx) / seedCount - Math.PI / 2
    const local = Math.floor(i / seedCount)
    const angle = baseAngle + (local % 5) * 0.35 - 0.7
    const radius = 320 + (local % 4) * 36
    positions.set(id, {
      x: 400 + radius * Math.cos(angle),
      y: 280 + radius * Math.sin(angle),
    })
  })

  return positions
}
