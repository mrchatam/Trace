# P34-S01-00 — Scope planner (plan)

## Metadata
- id: P34-S01-00
- todo_ids: [P34-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Finalize S01 so implementer authors **`PLAN.md`**: release embed pipeline, StaticDir defaults, auto-port behavior + flags, test matrix — from S00 `RESEARCH.md`. **No product code.** Do **not** author `PLAN.md` in this row.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md)
- [Phase 34 README](../../README.md)

## Session start

Follow agent-loop-protocol Session start. Prefer RESEARCH leans; do not reopen L1–L4.

## Locked defaults

| Item | Value |
|------|-------|
| Artifact of S01-01 | `scopes/scope-01-plan/PLAN.md` |
| Board | **P34-S01-00 → P34-S01-01** only (no separate S01 review row — S01-00 gate + S02-00 may spot-check PLAN) |
| Product code | **Forbidden** on S01 |
| Sequence | S02 (embed) then S03 (auto-port) — PLAN must not reverse board order |
| Must lock in PLAN | Embed pipeline steps; StaticDir consumer vs Trace-checkout; auto-port + `--addr`; test matrix; help/docs touch list for S04 |
| RESEARCH lean (S00 PASS — do not re-debate) | Embed **(A)** `go:embed` + `go generate`/make/CI populate `embeddist`; StaticDir = opportunistic disk-if-`index.html` → real embed → placeholder; auto-port = UA-increment `7432`…`7441` (max 10), host `127.0.0.1`, **shared** `gui`+`serve`; explicit `--addr` via **`flag.Changed`** = strict fail; reject `.trace/` SPA copy |

## Must answer (handoff to 01)

1. Exact release/dev commands to populate embeddist with real SPA? (seed: `web` build → sync into `internal/httpapi/embeddist` → `go build`; name `go generate` and/or `make embed-gui` / CI release step; S05 VERIFY fails if stub shipped when full SPA intended)
2. Default StaticDir resolution table (consumer / Trace repo / `--static-dir`)? (seed: keep `<root>/web/dist` as candidate only; no Trace-module probe; refuse StaticDir == root)
3. Auto-port: which commands (`gui`, `serve`, both), algorithm, max attempts? (seed: both; start `7432`; +1 on EADDRINUSE; max 10; loopback; `flag.Changed` not DefaultAddr string-equal; serve print moves **post-bind**)
4. Test matrix rows for S02/S03/S04/S05? (seed: consumer temp no `web/` → real SPA marker; default busy → `:7433` + printed/opened URL; `--addr` busy → fail; Trace checkout with `web/dist` → disk wins)
5. Docs files S04 must rewrite? (seed audit table in RESEARCH — gui-quickstart, help.go, local_http usage, web/README, embeddist README polish)

## Planner gate

- [x] `01-plan.md` runnable (metadata, exit criteria, PLAN template)
- [x] `SCOPE-TODOS.md` accurate
- [x] Do **not** write `PLAN.md` in this planner row

## Exit criteria

- [x] Plan implementer prompt locked; next **P34-S01-01**

## Todo updates

Status + notes on **P34-S01-00** only.

## Next

`P34-S01-01`
