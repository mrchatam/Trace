# P34-S02-02 — Embed review

## Metadata
- id: P34-S02-02
- todo_ids: [P34-S02-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of S02 embed pipeline + StaticDir/docs tone vs **L1/L2** and [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md). Small fixes OK; spawn forward (`02a`/`02b`) for structural gaps. **No auto-port scope creep.**

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1, L2
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — pipeline A, StaticDir table, T1–T3, T10 seed
- [01-implement.md](01-implement.md) — locked defaults + exit criteria
- Diff focus: `scripts/embed-gui.sh`, `internal/httpapi/embed.go` (`go:generate`), `embeddist/**` + README, `static.go` (`placeholderHTML`), `httpapi_test.go`, optional root `Makefile`, minimal `cmd/trace/local_http.go` / `help.go` if touched

## Session start

Fresh context. Follow agent-loop-protocol. Re-read L1/L2 + PLAN handoff for S02 before judging.

## Checklist

### L1 / L2 / PLAN

- [ ] Consumer without `web/` gets **real** SPA (test **T1** evidence: names + pass)
- [ ] No required consumer `web/` in `embeddist/README.md` (and any S02-touched usage strings)
- [ ] Stub not default when embed full — embedded `index.html` lacks `Embedded GUI stub`; has `#root` + `/assets/` module script
- [ ] StaticDir order still disk-if-`index.html` → embed → placeholder; default path still `<root>/web/dist`; **no** Trace-module probe; **no** SPA copy under consumer `.trace/`
- [ ] Trace-checkout DX: disk wins when planted (**T2**)
- [ ] `StaticDir == project root` still refused (**T3**)
- [ ] Pipeline: `scripts/embed-gui.sh` + `//go:generate`; sync preserves short README; optional Makefile only (no phantom CI invented)
- [ ] `placeholderHTML` last-resort tone — does not teach consumer `web/` as primary

### Scope creep / Law 19

- [ ] Auto-port **not** silently shipped (no listen hop / `flag.Changed` / FormatAddrInUse auto-exhausted paths)
- [ ] Law 19 — no business-logic fork in `web/`; httpapi remains adapter
- [ ] Full quickstart / AGENTS flip left to S04 (S02 may only minimal honesty)

### Failure modes to spot-check

- [ ] Script deletes README forever or leaves stub `index.html` after “successful” sync
- [ ] Tests pass only because Trace checkout has disk `web/dist` (T1 must use temp root **without** `web/`)
- [ ] `go:generate` path broken when run from `internal/httpapi/`
- [ ] Binary/embed still stub while Notes claim full SPA

## Verdict rule

- **PASS** (confidence medium+) → next **P34-S03-00**
- Blocker / missing T1–T3 → spawn `P34-S02-02a` / `02b` immediately below this row; do not start S03

## Exit criteria

- [ ] Confidence medium+; Notes cite evidence (files + test names)
- [ ] Board status `done` / spawn recorded
- [ ] Next **P34-S03-00** (or spawned remediation)

## Todo updates

Status + notes on **P34-S02-02** only.

## Next

`P34-S03-00`
