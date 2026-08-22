import { apiFetch, type RequestOptions } from './client'
import type { components } from './schema'

type Schemas = components['schemas']

export type HealthResponse = Schemas['HealthResponse']
export type VersionResponse = Schemas['VersionResponse']
export type ProjectResponse = Schemas['ProjectResponse']
export type TaskRow = Schemas['TaskRow']
export type TaskListResponse = Schemas['TaskListResponse']
export type SearchResponse = Schemas['SearchResponse']
export type SearchItem = SearchResponse['items'][number]
export type EntitySummary = Schemas['EntitySummary']
export type SeedStatus = Schemas['SeedStatus']
export type BoundedGraph = Schemas['BoundedGraph']
export type TransitionResult = Schemas['TransitionResult']
export type CreateTransitionRequest = Schemas['CreateTransitionRequest']
export type LoopStatusResponse = Schemas['LoopStatusResponse']
export type LoopGateResponse = Schemas['LoopGateResponse']
export type ContextPacket = Schemas['ContextPacket']
export type WhyPacket = Schemas['WhyPacket']

type TokenOpt = Pick<RequestOptions, 'token' | 'signal'>

/** getHealth */
export function getHealth(opt: TokenOpt = {}) {
  return apiFetch<HealthResponse>('/v1/health', opt)
}

/** getVersion */
export function getVersion(opt: TokenOpt = {}) {
  return apiFetch<VersionResponse>('/v1/version', opt)
}

/** getProject */
export function getProject(opt: TokenOpt = {}) {
  return apiFetch<ProjectResponse>('/v1/project', opt)
}

/** listTasks */
export function listTasks(
  query: { goal_id?: string; work_state?: string; cursor?: string; limit?: number } = {},
  opt: TokenOpt = {},
) {
  return apiFetch<TaskListResponse>('/v1/tasks', { ...opt, query })
}

/**
 * Fetch all task pages for auto-pick (P3a must see the true last item).
 * Pages with limit 100 while next_cursor is set; concatenates items.
 * Display lists may still call listTasks({ limit: 100 }) separately.
 */
export async function listTasksForPick(
  query: { goal_id?: string; work_state?: string } = {},
  opt: TokenOpt = {},
): Promise<TaskRow[]> {
  const items: TaskRow[] = []
  let cursor: string | undefined
  for (;;) {
    const res = await listTasks({ ...query, limit: 100, cursor }, opt)
    items.push(...(res.items ?? []))
    const next = res.next_cursor
    if (!next) break
    cursor = next
  }
  return items
}

/** getTask */
export function getTask(taskId: string, opt: TokenOpt = {}) {
  return apiFetch<TaskRow>(`/v1/tasks/${encodeURIComponent(taskId)}`, opt)
}

/** getLoopStatus — requires task_id */
export function getLoopStatus(taskId: string, goalId?: string, opt: TokenOpt = {}) {
  return apiFetch<LoopStatusResponse>('/v1/loop/status', {
    ...opt,
    query: { task_id: taskId, goal_id: goalId },
  })
}

/** getLoopGate — requires task_id */
export function getLoopGate(taskId: string, gateFor = 'edit', opt: TokenOpt = {}) {
  return apiFetch<LoopGateResponse>('/v1/loop/gate', {
    ...opt,
    query: { task_id: taskId, for: gateFor },
  })
}

/** search — response is { items: [...] }, never CLI hits */
export function search(q: string, limit = 20, opt: TokenOpt = {}) {
  return apiFetch<SearchResponse>('/v1/search', { ...opt, query: { q, limit } })
}

/** getEntity */
export function getEntity(entityId: string, opt: TokenOpt = {}) {
  return apiFetch<EntitySummary>(`/v1/entities/${encodeURIComponent(entityId)}`, opt)
}

/** getSeedStatus */
export function getSeedStatus(opt: TokenOpt = {}) {
  return apiFetch<SeedStatus>('/v1/seed/status', opt)
}

function coerceGraphArray<T>(value: unknown): T[] {
  return Array.isArray(value) ? value : []
}

/** getGraph — center + max_nodes required */
export function normalizeBoundedGraph(g: BoundedGraph): BoundedGraph {
  return {
    ...g,
    nodes: coerceGraphArray(g.nodes),
    edges: coerceGraphArray(g.edges),
  }
}

export async function getGraph(center: string, maxNodes: number, depth?: number, opt: TokenOpt = {}) {
  const raw = await apiFetch<BoundedGraph>('/v1/graph', {
    ...opt,
    query: { center, max_nodes: maxNodes, depth },
  })
  return normalizeBoundedGraph(raw)
}

/** getProjectGraph — mode=project; bounded full-project graph (Explore default). */
export async function getProjectGraph(maxNodes: number, opt: TokenOpt = {}) {
  const raw = await apiFetch<BoundedGraph>('/v1/graph', {
    ...opt,
    query: { mode: 'project', max_nodes: maxNodes },
  })
  return normalizeBoundedGraph(raw)
}

/** createTransition — denials are API envelope only */
export function createTransition(body: CreateTransitionRequest, opt: TokenOpt = {}) {
  return apiFetch<TransitionResult>('/v1/transitions', { ...opt, method: 'POST', body })
}

/** getContext */
export function getContext(taskId: string, depth: 1 | 2 = 1, opt: TokenOpt = {}) {
  return apiFetch<ContextPacket>('/v1/context', {
    ...opt,
    query: { task_id: taskId, depth, format: 'json' },
  })
}

/** getWhy */
export function getWhy(entityType: string, id: string, opt: TokenOpt = {}) {
  return apiFetch<WhyPacket>('/v1/why', {
    ...opt,
    query: { entity_type: entityType, id },
  })
}

/** getImpact — GET /v1/impact (OpenAPI). Prefer task_id when selection is a task (UX-IA). */
export function getImpact(taskId: string, opt: TokenOpt = {}) {
  return apiFetch<Record<string, unknown>>('/v1/impact', {
    ...opt,
    query: { task_id: taskId },
  })
}

export type LoopNextResponse = Schemas['LoopNextResponse']
export type LoopApplyEnvelope = Schemas['LoopApplyEnvelope']
export type LoopApplyResult = Schemas['LoopApplyResult']
export type LoopResetResponse = Schemas['LoopResetResponse']
export type CreateEntityRequest = Schemas['CreateEntityRequest']
export type CreateLinkRequest = Schemas['CreateLinkRequest']
export type LinkSummary = Schemas['LinkSummary']
export type SeedExportRequest = Schemas['SeedExportRequest']
export type SeedImportRequest = Schemas['SeedImportRequest']
export type SeedJobStatus = Schemas['SeedJobStatus']
export type ReviewSummary = Schemas['ReviewSummary']
export type ReviewListResponse = { items: ReviewSummary[] }
export type CreateReviewRequest = {
  title: string
  body?: string
  task_id?: string
}

/** getLoopNext — requires task_id */
export function getLoopNext(taskId: string, opt: TokenOpt = {}) {
  return apiFetch<LoopNextResponse>('/v1/loop/next', {
    ...opt,
    query: { task_id: taskId },
  })
}

/** postLoopApply */
export function postLoopApply(body: LoopApplyEnvelope, opt: TokenOpt = {}) {
  return apiFetch<LoopApplyResult>('/v1/loop/apply', { ...opt, method: 'POST', body })
}

/** postLoopReset */
export function postLoopReset(taskId: string, opt: TokenOpt = {}) {
  return apiFetch<LoopResetResponse>('/v1/loop/reset', {
    ...opt,
    method: 'POST',
    body: { task_id: taskId },
  })
}

/** createEntity */
export function createEntity(body: CreateEntityRequest, opt: TokenOpt = {}) {
  return apiFetch<EntitySummary>('/v1/entities', { ...opt, method: 'POST', body })
}

/** createLink */
export function createLink(body: CreateLinkRequest, opt: TokenOpt = {}) {
  return apiFetch<LinkSummary>('/v1/links', { ...opt, method: 'POST', body })
}

/** postSeedExport */
export function postSeedExport(body: SeedExportRequest = {}, opt: TokenOpt = {}) {
  return apiFetch<SeedJobStatus>('/v1/seed/export', { ...opt, method: 'POST', body })
}

/** postSeedImport */
export function postSeedImport(body: SeedImportRequest, opt: TokenOpt = {}) {
  return apiFetch<SeedJobStatus>('/v1/seed/import', { ...opt, method: 'POST', body })
}

/** listReviews */
export function listReviews(query: { task_id?: string } = {}, opt: TokenOpt = {}) {
  return apiFetch<ReviewListResponse>('/v1/reviews', { ...opt, query })
}

/** getReview */
export function getReview(reviewId: string, opt: TokenOpt = {}) {
  return apiFetch<ReviewSummary>(`/v1/reviews/${encodeURIComponent(reviewId)}`, opt)
}

/** createReview */
export function createReview(body: CreateReviewRequest, opt: TokenOpt = {}) {
  return apiFetch<ReviewSummary>('/v1/reviews', { ...opt, method: 'POST', body })
}
