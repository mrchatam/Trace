# P34-S03-02 — Auto-port review

## Metadata
- id: P34-S03-02
- todo_ids: [P34-S03-02]
- role: reviewer
- skills: [code-review-and-quality, security-and-hardening]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of S03 auto free-port vs **L3/L4** and [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md). Confirm UA-increment `7432`–`7441`, `flag.Changed` strict `--addr`, shared httpapi helper, post-bind print/open. Small fixes OK; spawn forward (`02a`/`02b`) for structural gaps. **No** StaticDir/embed rework; **no** public-bind creep.

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L3, L4 + clarifications
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — auto-port table, T4–T7, T11
- [01-implement.md](01-implement.md) — locked defaults + exit criteria
- Diff focus: `internal/httpapi/server.go` (`ListenAndServe` / Options), `bind.go` / new helper, `addr_in_use.go` (+ tests), `cmd/trace/local_http.go`, `help.go`, `gui_test.go` / `serve_test.go` / httpapi listen tests

## Session start

Fresh context. Follow agent-loop-protocol. Re-read L3/L4 + PLAN auto-port section before judging. Do not reopen algorithm (UA-incr max 10) unless live contradicts locks.

## Checklist

### L3 / L4 / PLAN

- [ ] Default busy → next free loopback port (**T4** evidence: names + pass)
- [ ] Default free → still `:7432` first (**T11**)
- [ ] Two concurrent default `gui`/`serve` → distinct ports; gui open/hook = chosen addr (**T5**)
- [ ] Explicit `--addr` busy → **fail**, no hop; includes `--addr` == DefaultAddr string when Changed (**T6**)
- [ ] Auto exhausted (10 ports) → fail mentions range + `--addr` (**T7**)
- [ ] Detection is **`flag.Changed`**, not DefaultAddr string-equal
- [ ] Hop lives in **`internal/httpapi`**; `gui` and `serve` share it (no cmd-only fork)
- [ ] Printed URL and `gui` open use **post-bind** chosen addr; serve no longer prints pre-listen stale addr
- [ ] Help/usage: no “no auto free-port” for default; `--addr` documented as pin/strict
- [ ] L4: still one project root per process (no multi-root daemon)

### Security / Law 19

- [ ] Default host / refuse-remote / `--allow-remote` posture **unchanged** (no `0.0.0.0` default; hop does not widen host)
- [ ] No OS `:0` as Trace default happy path
- [ ] Law 19 — httpapi remains adapter; no business-logic fork in `web/`
- [ ] Full quickstart / AGENTS left to S04 (S03 may only help/usage/addr-in-use)

### Failure modes to spot-check

- [ ] Hopping on explicit `--addr` (including pinned DefaultAddr) — **must not**
- [ ] `OnListening` / open still using pre-hop `Options.Addr` after successful hop
- [ ] serve stderr claims listening before bind fails or before hop completes
- [ ] Exhaust path still uses old “pick `--addr` only” copy with no auto-range mention
- [ ] Tests green only because machine already free on 7432 while T4 never occupies it
- [ ] Parallel tests racing on `7432`–`7441`

## Verdict rule

- **PASS** (confidence medium+) → next **P34-S04-00**
- Blocker / missing T4–T7/T11 → spawn `P34-S03-02a` / `02b` immediately below this row; do not start S04

## Exit criteria

- [ ] Confidence medium+; Notes cite evidence (files + test names)
- [ ] Board status `done` / spawn recorded
- [ ] Next **P34-S04-00** (or spawned remediation)

## Todo updates

Status + notes on **P34-S03-02** only.

## Next

`P34-S04-00`
