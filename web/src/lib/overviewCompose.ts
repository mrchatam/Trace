/**
 * Explore overview seed compose + merge (Laws 6–7).
 * Pure helpers — unit-testable without React.
 */

import type { BoundedGraph, SearchItem, TaskRow } from '../api/ops'

export const PROJECT_MAX_NODES = 500
export const SEED_TARGET = 6
export const SEED_CAP = 8
export const SEED_MAX_NODES = 40
export const UI_CAP = PROJECT_MAX_NODES
export const EXPAND_MAX_NODES = 50
export const DEPTH = 2

/** High-signal FTS fill queries (non-empty; empty q → 400). */
export const SEARCH_FILL_QUERIES = ['goal', 'capability', 'decision', 'discovery'] as const

const TERMINAL = new Set(['DONE', 'SKIPPED'])

function coerceGraphEdges(edges: unknown): BoundedGraph['edges'] {
  return Array.isArray(edges) ? edges : []
}

export type SeedCandidate = {
  id: string
  kind: string
  title: string
  work_state?: string
  source: 'task' | 'search'
}

export type GraphNodeMeta = {
  id: string
  kind: string
  title: string
  work_state?: string
  goal_id?: string
}

export type MergedOverview = {
  nodes: GraphNodeMeta[]
  edges: BoundedGraph['edges']
  truncated: boolean
  omitted: number
  max_nodes: number
  depth: number
  center: string
}

function isTerminal(ws: string | undefined | null): boolean {
  return ws != null && TERMINAL.has(ws)
}

/** Prefer IN_PROGRESS, then non-terminal; terminal only if needed for a useful minimum. */
export function prioritizeTaskSeeds(
  tasks: TaskRow[],
  limit: number,
): SeedCandidate[] {
  if (limit <= 0) return []
  const inProgress: TaskRow[] = []
  const nonTerminal: TaskRow[] = []
  const terminal: TaskRow[] = []
  for (const t of tasks) {
    if (t.work_state === 'IN_PROGRESS') inProgress.push(t)
    else if (isTerminal(t.work_state)) terminal.push(t)
    else nonTerminal.push(t)
  }
  const ordered = [...inProgress, ...nonTerminal]
  if (ordered.length === 0) {
    ordered.push(...terminal)
  }
  return ordered.slice(0, limit).map((t) => ({
    id: t.id,
    kind: 'task',
    title: t.title,
    work_state: t.work_state,
    source: 'task' as const,
  }))
}

/** Dedupe by id; stop at cap. */
export function dedupeSeeds(
  candidates: SeedCandidate[],
  cap: number = SEED_CAP,
): SeedCandidate[] {
  const seen = new Set<string>()
  const out: SeedCandidate[] = []
  for (const c of candidates) {
    if (!c.id || seen.has(c.id)) continue
    seen.add(c.id)
    out.push(c)
    if (out.length >= cap) break
  }
  return out
}

/**
 * Compose seed list: tasks first (priority order), then search hits to fill toward target (≤cap).
 * getProject is chrome-only today — do not invent a project node.
 */
export function composeSeedsFromParts(
  tasks: TaskRow[],
  searchHits: SearchItem[],
  opts: { target?: number; cap?: number } = {},
): SeedCandidate[] {
  const target = opts.target ?? SEED_TARGET
  const cap = opts.cap ?? SEED_CAP
  const fromTasks = prioritizeTaskSeeds(tasks, cap)
  const merged = dedupeSeeds(fromTasks, cap)
  if (merged.length >= target) return merged.slice(0, cap)
  const need = Math.min(cap - merged.length, target - merged.length)
  const extras = searchHits
    .filter((h) => h.id && !merged.some((s) => s.id === h.id))
    .slice(0, need)
    .map(
      (h): SeedCandidate => ({
        id: h.id,
        kind: h.kind,
        title: h.title,
        source: 'search',
      }),
    )
  return dedupeSeeds([...merged, ...extras], cap)
}

/** Collect search items across fill queries (caller runs sequential search). */
export function collectSearchFill(
  batches: SearchItem[][],
  excludeIds: Set<string>,
  need: number,
): SearchItem[] {
  if (need <= 0) return []
  const out: SearchItem[] = []
  const seen = new Set(excludeIds)
  for (const batch of batches) {
    for (const hit of batch) {
      if (!hit.id || seen.has(hit.id)) continue
      seen.add(hit.id)
      out.push(hit)
      if (out.length >= need) return out
    }
  }
  return out
}

function edgeDegree(id: string, edges: unknown): number {
  let d = 0
  for (const e of coerceGraphEdges(edges)) {
    if (e.from === id || e.to === id) d++
  }
  return d
}

/**
 * Merge parallel getGraph results; keep all seeds first, then neighbors by
 * seed proximity (edge to a seed) then degree. Trim to UI_CAP.
 */
export function mergeOverviewGraphs(
  graphs: BoundedGraph[],
  seedIds: string[],
  opts: { uiCap?: number; depth?: number } = {},
): MergedOverview {
  const uiCap = opts.uiCap ?? UI_CAP
  const depth = opts.depth ?? DEPTH
  const seedSet = new Set(seedIds)
  const nodeMap = new Map<string, GraphNodeMeta>()
  const edgeKey = new Set<string>()
  const edges: BoundedGraph['edges'] = []
  let anyTruncated = false
  let maxReported = 0

  for (const g of graphs) {
    if (g.truncated) anyTruncated = true
    maxReported = Math.max(maxReported, g.max_nodes ?? 0)
    const nodes = Array.isArray(g.nodes) ? g.nodes : []
    const graphEdges = coerceGraphEdges(g.edges)
    for (const n of nodes) {
      if (!nodeMap.has(n.id)) {
        nodeMap.set(n.id, { id: n.id, kind: n.kind, title: n.title })
      }
    }
    for (const e of graphEdges) {
      const k = `${e.from}|${e.rel}|${e.to}`
      if (edgeKey.has(k)) continue
      edgeKey.add(k)
      edges.push(e)
    }
  }

  const allIds = Array.from(nodeMap.keys())
  const seedPresent = seedIds.filter((id) => nodeMap.has(id))
  const nonSeeds = allIds.filter((id) => !seedSet.has(id))

  nonSeeds.sort((a, b) => {
    const aNear = edges.some(
      (e) =>
        (e.from === a && seedSet.has(e.to)) || (e.to === a && seedSet.has(e.from)),
    )
      ? 1
      : 0
    const bNear = edges.some(
      (e) =>
        (e.from === b && seedSet.has(e.to)) || (e.to === b && seedSet.has(e.from)),
    )
      ? 1
      : 0
    if (bNear !== aNear) return bNear - aNear
    return edgeDegree(b, edges) - edgeDegree(a, edges)
  })

  const keepIds = new Set<string>()
  for (const id of seedPresent) {
    if (keepIds.size >= uiCap) break
    keepIds.add(id)
  }
  for (const id of nonSeeds) {
    if (keepIds.size >= uiCap) break
    keepIds.add(id)
  }

  const omitted = Math.max(0, allIds.length - keepIds.size)
  const nodes = Array.from(keepIds).map((id) => nodeMap.get(id)!)
  const keptEdges = edges.filter((e) => keepIds.has(e.from) && keepIds.has(e.to))
  const center = seedPresent[0] ?? nodes[0]?.id ?? ''

  return {
    nodes,
    edges: keptEdges,
    truncated: anyTruncated || omitted > 0,
    omitted,
    max_nodes: Math.min(uiCap, maxReported || uiCap),
    depth,
    center,
  }
}

/** Attach known work_state from task seeds onto merged nodes. */
export function applySeedMeta(
  nodes: GraphNodeMeta[],
  seeds: SeedCandidate[],
): GraphNodeMeta[] {
  const byId = new Map(seeds.map((s) => [s.id, s]))
  return nodes.map((n) => {
    const s = byId.get(n.id)
    if (!s) return n
    return {
      ...n,
      kind: n.kind || s.kind,
      title: n.title || s.title,
      work_state: s.work_state ?? n.work_state,
    }
  })
}

/** Partial-fail: keep successful graphs only; normalize null JSON arrays. */
export function partitionSettledGraphs(
  results: PromiseSettledResult<BoundedGraph>[],
  seedIds: string[],
): { graphs: BoundedGraph[]; failedSeedIds: string[]; hardFail: boolean } {
  const graphs: BoundedGraph[] = []
  const failedSeedIds: string[] = []
  results.forEach((r, i) => {
    if (r.status === 'fulfilled') {
      graphs.push({
        ...r.value,
        nodes: Array.isArray(r.value.nodes) ? r.value.nodes : [],
        edges: coerceGraphEdges(r.value.edges),
      })
    } else failedSeedIds.push(seedIds[i] ?? `seed-${i}`)
  })
  return {
    graphs,
    failedSeedIds,
    hardFail: graphs.length === 0 && seedIds.length > 0,
  }
}
