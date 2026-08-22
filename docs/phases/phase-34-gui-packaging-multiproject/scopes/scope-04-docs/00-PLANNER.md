# P34-S04-00 — Scope planner (docs)

## Metadata
- id: P34-S04-00
- todo_ids: [P34-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-guidelines]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock S04 to flip consumer-facing docs/help/AGENTS: **no consumer `web/`**, embed primary, **auto-port**, multi-project Just Works. Residual tests per PLAN. No Explore UI work.

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1–L3
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — docs touch list
- Live: `docs/gui-quickstart.md`, `README.md`, `AGENTS.md`, `web/README.md`, `cmd/trace/help.go`, `internal/httpapi/embeddist/README.md`

## Session start

Follow agent-loop-protocol. Prefer PLAN touch list.

## Locked defaults

| Item | Value |
|------|-------|
| Primary story | `trace gui` from any Trace-initialized repo — SPA from binary; no `web/` required |
| Multi-project | Second `trace gui` auto-picks free port; print + open correct URL |
| `--addr` | Document pin + fail-if-busy |
| Contributor DX | Trace checkout may still build `web/dist` for local GUI iteration — clearly labeled **contributor**, not consumer |
| Residual tests | Any PLAN matrix rows deferred from S02/S03 |
| Out | VERIFY (S05); changing L1–L4 |
| S01-00 seed (prefer PLAN) | Touch list from RESEARCH audit: `docs/gui-quickstart.md` (primary), `cmd/trace/help.go`, `local_http.go` usage, `web/README.md` (contributor), root `README.md` if needed, `AGENTS.md` optional, polish `embeddist/README.md`; assert T8 (no consumer two-artifact / no “no auto free-port” for default) |
| S03-00 seed | **S03 owns** port help/usage + `FormatAddrInUseMessage` (auto vs strict). S04 verifies T8 residual: quickstart/AGENTS/`web/README` still may lag until this scope; do not re-implement hop — document multi-project auto-port + `--addr` pin |
| S03-02 note | Auto-port shipped: UA-incr `7432`–`7441` in `httpapi.ListenAndServe`; explicit pin via CLI `fs.Visit("addr")` (stdlib has no `Flag.Changed`; Visit ≡ PLAN intent). Docs may say “`--addr` set on cmdline” / pin-strict — do not invent a `flag.Changed` API in user-facing copy. |

## Planner gate

- [x] `01-implement.md` + `02-review.md` + `SCOPE-TODOS.md` ready

## Exit criteria

- [x] Next **P34-S04-01**

## Todo updates

Status + notes on **P34-S04-00** only.

## Next

`P34-S04-01`
