# P34-S03-00 — Scope planner (auto free-port)

## Metadata
- id: P34-S03-00
- todo_ids: [P34-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Lock S03 to ship **L3 auto free-port** for default bind on `trace gui` / `serve` (per PLAN). Explicit `--addr` remains fail-if-busy. Print + open **chosen** URL. Overturn P32/P33 no-auto-port for happy path only.

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L3, L4
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md)
- Live: `internal/httpapi/addr_in_use.go`, `bind.go`, `cmd/trace/local_http.go`, `gui_test.go`, serve tests

## Session start

Follow agent-loop-protocol. Prefer PLAN algorithm; do not reopen L3.

## Locked defaults

| Item | Value |
|------|-------|
| Default bind busy | Auto try next free **loopback** port (algorithm from PLAN) |
| Explicit `--addr` | **Strict** — fail if busy (friendly message OK) |
| After listen | Print URL; `gui` opens browser to that URL (unless `--no-open`) |
| Commands | Prefer **both** `gui` and `serve` share bind helper |
| Loopback | Do not weaken refuse-remote / default host |
| Out | Embed/StaticDir (S02 done); full docs flip (S04) |
| S01-00 seed (prefer PLAN) | UA-increment start **7432**, host **127.0.0.1**, max **10** (`7432`–`7441`); detect explicit `--addr` via **`flag.Changed`** (not DefaultAddr string-equal); shared `gui`+`serve`; move **serve** listen print to **post-bind**; exhausted → fail + `--addr` tip |

## Must answer (into 01)

1. Exact algorithm + max attempts from PLAN? (seed: +1 on EADDRINUSE, max 10)
2. Shared helper location (`httpapi` vs `cmd/trace`)?
3. Test: two concurrent default `gui` → distinct ports + correct open URL? (PLAN T4–T7)

## Planner gate

- [ ] `01-implement.md` + `02-review.md` + `SCOPE-TODOS.md` ready

## Exit criteria

- [ ] Next **P34-S03-01**

## Todo updates

Status + notes on **P34-S03-00** only.

## Next

`P34-S03-01`
