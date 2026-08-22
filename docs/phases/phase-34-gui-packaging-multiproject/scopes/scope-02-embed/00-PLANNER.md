# P34-S02-00 — Scope planner (embed + static defaults)

## Metadata
- id: P34-S02-00
- todo_ids: [P34-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock S02 implementer to ship **real SPA into embed** + **StaticDir resolution** so consumer roots without `web/` serve the Explore SPA (not stub). Follow S01 `PLAN.md`. **No auto-port** in this scope (S03).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1, L2
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md)
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md)
- Live: `internal/httpapi/static.go`, `embed.go`, `embeddist/`, `server.go`, `cmd/trace/local_http.go`, `Makefile` / CI if any, `web/package.json`

## Session start

Follow agent-loop-protocol Session start. Prefer PLAN; do not reopen L1–L4; do not implement auto-port.

## Locked defaults

| Item | Value |
|------|-------|
| Ship | Build pipeline so release/dev binary embeds **real** `web/dist`; update resolution so consumer without `web/` gets embed SPA |
| Trace-checkout DX | Disk `web/dist` may still win **in Trace repo** when present (per PLAN) — never require consumer `web/` |
| Explicit `--static-dir` | Honored; still refuse StaticDir == project root |
| Stub | Last-resort only; tests must prove consumer path is **not** stub when embed full |
| Out | Auto free-port (S03); full docs flip (S04) |
| Law 19 | httpapi adapter only — no business-logic fork |
| S01-00 seed (prefer PLAN) | Pipeline **(A):** `cd web && npm ci && npm run build` → sync `web/dist` → `embeddist` (keep short README) → `go build`; primary `scripts/embed-gui.sh` + `//go:generate` (optional `make embed-gui`). StaticDir: keep `<root>/web/dist` candidate; opportunistic disk-if-`index.html` → embed → placeholder; no Trace-module probe. SPA markers: `#root` + `/assets/` module script; stub phrase `Embedded GUI stub` |

## Must answer (into 01-implement)

1. Exact files/scripts changed for embed pipeline? (seed: `scripts/embed-gui.sh`, `//go:generate` near `embed.go`, optional Makefile `embed-gui`, `embeddist/*` + README rewrite)
2. Resolution order after change? (seed: unchanged disk→embed→placeholder; default path string stays)
3. Tests proving consumer temp dir without `web/` gets real SPA? (PLAN T1; also T2 disk-wins, T3 refuse root)
4. How to detect stub vs real SPA in tests (e.g. not placeholder title / has `#root` / assets)?

## Planner gate

- [x] `01-implement.md` runnable with locked defaults + exit criteria
- [x] `02-review.md` checklist vs L1/L2 + PLAN
- [x] `SCOPE-TODOS.md` accurate

## Exit criteria

- [x] Implementer locked; next **P34-S02-01**

## Todo updates

Status + notes on **P34-S02-00** only.

## Next

`P34-S02-01`
