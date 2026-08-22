/** Build Explore (/) href with optional neighborhood center. */
export type GraphViewMode = 'project' | 'neighborhood'

export function buildGraphHref(center?: string, mode: GraphViewMode = 'neighborhood'): string {
  if (!center?.trim()) return '/'
  const params = new URLSearchParams()
  params.set('center', center.trim())
  if (mode === 'neighborhood') params.set('mode', 'neighborhood')
  return `/?${params.toString()}`
}

export function parseGraphSearchParams(search: URLSearchParams): {
  center: string | null
  mode: GraphViewMode
} {
  const center = search.get('center')?.trim() || null
  if (center) return { center, mode: 'neighborhood' }
  const modeRaw = search.get('mode')
  return { center: null, mode: modeRaw === 'neighborhood' ? 'neighborhood' : 'project' }
}
