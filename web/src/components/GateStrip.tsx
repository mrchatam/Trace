import type { LoopGateResponse } from '../api/ops'

type GateLike = LoopGateResponse & {
  allowed?: boolean
  violations?: Array<{ reason_code?: string; message?: string; recommended_phase?: string }>
  reason_code?: string
}

type Props = {
  gate: GateLike | null
  loading?: boolean
}

export function GateStrip({ gate, loading }: Props) {
  if (loading) {
    return (
      <div className="banner banner--info" role="status">
        Loading gate…
      </div>
    )
  }
  if (!gate) {
    return (
      <div className="banner banner--info" role="status">
        No gate data — pick a task to evaluate.
      </div>
    )
  }

  const allowed = Boolean(gate.allowed)
  const violations = Array.isArray(gate.violations) ? gate.violations : []

  if (allowed && violations.length === 0) {
    return (
      <div className="banner banner--ok" role="status">
        Gate OK
      </div>
    )
  }

  return (
    <div className={`banner ${allowed ? 'banner--warn' : 'banner--error'}`} role="status">
      <strong>{allowed ? 'Gate warnings' : 'Gate blocked'}</strong>
      {gate.reason_code ? <div className="mono">{gate.reason_code}</div> : null}
      {violations.length > 0 ? (
        <ul style={{ margin: '0.5rem 0 0', paddingLeft: '1.25rem' }}>
          {violations.map((v, i) => (
            <li key={i}>
              {v.reason_code ? <span className="mono">{v.reason_code}: </span> : null}
              {v.message ?? 'violation'}
              {v.recommended_phase ? ` → ${v.recommended_phase}` : null}
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  )
}
