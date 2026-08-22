# P34-S04-02 — Docs review

## Metadata
- id: P34-S04-02
- todo_ids: [P34-S04-02]
- role: reviewer
- skills: [code-review-and-quality, writing-guidelines]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of S04 docs/help vs **L1–L3** (and L4 multi-process story) + PLAN **T8**. Spot-check that consumers are **not** taught to ship `web/` or manual free-port for default bind. Small doc fixes OK; spawn `02a`/`02b` for structural gaps. **No** Explore craft; **no** re-open of auto-port/embed implementation unless docs claim false behavior.

## References

- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1–L4
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — touch list + T8
- [01-implement.md](01-implement.md) — locked defaults + live debt map
- Live: `docs/gui-quickstart.md`, `web/README.md`, `AGENTS.md`, root `README.md`, `cmd/trace/help.go`, `cmd/trace/local_http.go` usage, `internal/httpapi/embeddist/README.md`, `FormatAddrInUseMessage`, `placeholderHTML`

## Session start

Fresh context. Re-read locks + PLAN T8; **do not** trust implementer Notes alone — open the files.

## Checklist (must all hold for PASS)

### Consumer story (L1 / L2)

- [ ] Primary path = `trace gui` in app repo with **`.trace/` only** — no required consumer `web/`
- [ ] SPA source taught as **binary embed**; disk `web/dist` labeled **contributor / override** when mentioned
- [ ] No primary teaching of “build `web/dist` in your app” / two-artifact everyday / “embedded stub” as normal UX
- [ ] `embeddist/README.md` + `placeholderHTML` still do not teach consumer `web/` as primary (spot-check)

### Auto-port / multi-project (L3 / L4)

- [ ] Default busy → next free loopback port documented (`7432`–`7441` or equivalent max-10 wording)
- [ ] Second project / concurrent default **does not** require manual `--addr` as the happy path
- [ ] `--addr` documented as **pin / fail-if-busy** (cmdline set); no fake `flag.Changed` API in user copy
- [ ] One process = one root; multi-project = multiple processes (L4)
- [ ] Quickstart sample stderr / hints match live `FormatAddrInUseMessage` (or exhausted-range) — not P32 “pick a free port yourself” as default

### PATH / help / AGENTS

- [ ] PATH install ≠ `trace install` preserved where discussed
- [ ] `help.go` / usage still accurate (auto-port + embed); no regression to “no auto free-port”
- [ ] `AGENTS.md` next-runnable / Phase 33 “until P34” language not contradicting shipped embed+auto-port
- [ ] Root `README.md` does not imply consumer must own `web/`

### T8 + residuals

- [ ] T8: greppable consumer surfaces free of forbidden “no auto free-port” / consumer two-artifact primary teaching
- [ ] Residual product matrix (T1–T7/T10/T11) correctly treated as already done — no false “tests missing” without evidence
- [ ] Diff is docs/help-string scope — no silent hop/embed product rewrites

### PASS gate

- [ ] Confidence **medium+**
- [ ] PASS → next **P34-S05-00**; FAIL → spawn or Notes with concrete file:line fixes

## Failure modes (reject / spawn)

- [ ] Quickstart still says Trace does **not** auto-pick a free port
- [ ] Quickstart secondary path still presents `cd web && npm run build` as normal consumer setup
- [ ] `web/README.md` still “no auto free-port” without contributor framing
- [ ] Docs claim auto-hop on explicit `--addr`
- [ ] Docs invent OS `:0` or public bind as Trace default
- [ ] Implementer changed product listen/embed code “to make docs true” instead of documenting shipped behavior
- [ ] AGENTS still points agents at long-closed `P34-S01-01` as next runnable

## Role work

1. Spot-check every PLAN touch-list path (table in PLAN + 01-implement debt map).
2. Grep for: `no auto free-port`, `does not auto-pick`, `two-artifact`, `embedded stub`, `Pick a free port yourself` in consumer-facing docs.
3. Compare multi-project section examples to shipped S03 behavior (UA-incr, Visit/pin).
4. Small doc fixes in-place OK; structural gap → spawn `P34-S04-02a` / `02b` on board.
5. Board Notes: PASS/FAIL, confidence, files checked, any fixes.

## Exit criteria

- [ ] Checklist complete; confidence medium+
- [ ] Next **P34-S05-00** on PASS

## Todo updates

Status + notes on **P34-S04-02** only.

## Next

`P34-S05-00`
