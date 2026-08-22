# Scope 02 — board map

**S02 embed + static defaults.** Serial: **P34-S02-00 → P34-S02-01 → P34-S02-02**. Do **not** implement auto-port (S03).

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 598 | P34-S02-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock implement defaults (this row) |
| 599 | P34-S02-01 | [01-implement.md](01-implement.md) | Implementer | `scripts/embed-gui.sh` + `go:generate` + sync real SPA into `embeddist` + tone + T1–T3 |
| 600 | P34-S02-02 | [02-review.md](02-review.md) | Reviewer | L1/L2 + PLAN checklist; no auto-port creep |

## Inputs (verified baseline 2026-08-21)

| Source | Fact |
|--------|------|
| [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) | Pipeline **A**; StaticDir opportunistic; T1–T3 (+ T10 seed); S02 handoff |
| L1 / L2 | Consumer = `.trace/` only; real SPA from binary embed; disk `web/dist` = contributor DX |
| Live `embeddist/` | **Stub only** — `index.html` contains `Embedded GUI stub`; README teaches two-artifact everyday |
| Live resolution | `static.go`: disk → embed → `placeholderHTML`; `server.go`: default `<root>/web/dist`, refuse == root |
| Live tooling | No Trace-root `Makefile` / `.github`; `scripts/` exists but no `embed-gui.sh`; `web/package.json` build → `web/dist` |
| Real SPA markers | Live `web/dist/index.html`: `<div id="root">` + `script type="module" … /assets/…` |

## Locked answers (P34-S02-00)

1. **Pipeline files:** `scripts/embed-gui.sh` + `//go:generate` near `embed.go` + sync `embeddist/*` + README rewrite; optional `make embed-gui`.
2. **StaticDir order:** unchanged disk → embed → placeholder; default path string unchanged.
3. **Tests:** T1 consumer no `web/` → real SPA; T2 disk wins; T3 refuse root (existing).
4. **Detection:** not `Embedded GUI stub`; require `#root` + `/assets/` module script.

## Out of this scope

- Auto free-port (S03); full docs/quickstart flip (S04); VERIFY (S05); CI workflow invention; Explore redesign.
