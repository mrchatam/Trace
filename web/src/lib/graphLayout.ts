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
/** At/above: full node cards when compact overview is active. */
export const FULL_NODE_MIN_ZOOM = 0.75
/** Below: hide edges entirely (too dense at overview). */
export const EDGE_HIDE_MAX_ZOOM = 0.25
/** Below + many edges: render a sampled subset. */
export const EDGE_SAMPLE_MIN_ZOOM = 0.45
export const EDGE_SAMPLE_MIN_COUNT = 100
export const EDGE_SAMPLE_RATIO = 3
export const PROJECT_FIT_PADDING = 0.2

export type NodeLod = 'minimal' | 'dot' | 'full'

export type LayoutEdge = { from: string; to: string }

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
        .distance(78)
        .strength(0.38),
    )
    .force('charge', forceManyBody().strength(-150).distanceMax(440))
    .force('center', forceCenter(width / 2, height / 2).strength(0.1))
    .force('collide', forceCollide(24))

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
  if (!compactOverview) {
    if (zoom < LOD_MINIMAL_MAX_ZOOM) return 'minimal'
    if (zoom < LOD_DOT_MAX_ZOOM) return 'dot'
    return 'full'
  }
  if (zoom >= FULL_NODE_MIN_ZOOM) return 'full'
  if (zoom >= LOD_MINIMAL_MAX_ZOOM) return 'dot'
  return 'minimal'
}

export function shouldShowNodeLabel(lod: NodeLod): boolean {
  return lod === 'full'
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
