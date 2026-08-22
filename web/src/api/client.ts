/** API client — relative `/v1` only (Law 19). */

export type ApiErrorBody = {
  error: {
    code: string
    message: string
    details?: unknown
  }
}

export class ApiError extends Error {
  readonly status: number
  readonly code: string
  readonly details: unknown

  constructor(status: number, code: string, message: string, details?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.details = details
  }
}

export type RequestOptions = {
  method?: string
  query?: Record<string, string | number | boolean | undefined | null>
  body?: unknown
  token?: string | null
  signal?: AbortSignal
}

function buildUrl(path: string, query?: RequestOptions['query']): string {
  const base = path.startsWith('/v1') ? path : `/v1${path.startsWith('/') ? path : `/${path}`}`
  if (!query) return base
  const params = new URLSearchParams()
  for (const [k, v] of Object.entries(query)) {
    if (v === undefined || v === null || v === '') continue
    params.set(k, String(v))
  }
  const qs = params.toString()
  return qs ? `${base}?${qs}` : base
}

export async function apiFetch<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    Accept: 'application/json',
  }
  if (opts.body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }
  if (opts.token) {
    headers.Authorization = `Bearer ${opts.token}`
  }

  const res = await fetch(buildUrl(path, opts.query), {
    method: opts.method ?? (opts.body !== undefined ? 'POST' : 'GET'),
    headers,
    body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    signal: opts.signal,
  })

  const text = await res.text()
  let data: unknown = undefined
  if (text) {
    try {
      data = JSON.parse(text) as unknown
    } catch {
      throw new ApiError(res.status, 'BAD_RESPONSE', 'Response was not JSON')
    }
  }

  if (!res.ok) {
    const env = data as ApiErrorBody | undefined
    const code = env?.error?.code ?? 'HTTP_ERROR'
    const message = env?.error?.message ?? (res.statusText || `HTTP ${res.status}`)
    throw new ApiError(res.status, code, message, env?.error?.details)
  }

  return data as T
}

export function isApiError(err: unknown): err is ApiError {
  return err instanceof ApiError
}
