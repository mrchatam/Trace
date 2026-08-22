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

export const COMPACT_NODE_THRESHOLD = 30
export const EDGE_LABEL_MAX = 60
export const EDGE_LABEL_MIN_ZOOM = 0.55
export const FULL_NODE_MIN_ZOOM = 0.75
export const PROJECT_FIT_PADDING = 0.2

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

export function shouldShowFullNode(
  compactOverview: boolean,
  zoom: number,
  nodeId: string,
  selectedId: string | null,
  centerId: string,
): boolean {
  if (!compactOverview) return true
  if (nodeId === selectedId || nodeId === centerId) return true
  return zoom >= FULL_NODE_MIN_ZOOM
}

export function shouldShowEdgeLabels(edgeCount: number, zoom: number): boolean {
  if (edgeCount > EDGE_LABEL_MAX) return false
  return zoom >= EDGE_LABEL_MIN_ZOOM
}
