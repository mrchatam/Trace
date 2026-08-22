# P02 / S01 / 02 — Scope review (X0 Gate C)

## Metadata
- id: P02-S01-02
- todo_ids: [P02-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: mixed

## Objective
Independent review of S01 Gate C evaluation deliverables. Confirm evidence matches runs; scoring is fair; kill criteria applied honestly; Phase 01 dry-run was **not** treated as Gate C pass. Spawn remediations on blocker/high. Produce/confirm issue-list pointer usable by S02.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop)
- Sibling [01-gate-c-eval.md](01-gate-c-eval.md) + board Notes
- Expected artifacts: `GATE-C-NOTES.md`, `docs/verification/gate-c-x0/`, `evals/x0` Gate C harness
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md)
- Live regressions: `evals/x0` dry-run, `evals/p0x`, `evals/honesty`

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → review). Fresh subagent ≠ implementer.

## Review focus (checklist)

### Honesty of Gate C claim
- [ ] Verdict in `GATE-C-NOTES.md` is exactly one of Go / No-Go / Iterate
- [ ] Evidence table present (N, mean accuracy, pins) — not vibes
- [ ] Phase 01 `dry_run:true` explicitly **not** used as pass
- [ ] Kill criteria check filled; No-Go/Iterate if G1≤B0 with non-trivial seeding

### Fairness & protocol
- [ ] Same query bank for B0 and G1; ≥5 queries with GT keys
- [ ] B0 tools/path never invoke `why`/`context`
- [ ] G1 may use CLI why/context; tools_used reflects that
- [ ] N≥3 runs per condition in persisted metrics (`dry_run: false`)
- [ ] Scoring grades match locked definitions (accuracy, critical_miss)
- [ ] No contaminated GT / easier G1 questions

### Artifacts & reproducibility
- [ ] `docs/verification/gate-c-x0/metrics-b0.json` + `metrics-g1.json` schema-valid
- [ ] Pins: model, fixture identity, trace version, date (pins.md or Notes)
- [ ] Dry-run test still green; schema bump justified if any

### Bars / laws
- [ ] `evals/p0x` + `evals/honesty` still PASS; not merged into X0
- [ ] No MCP-required creep; no daemon/HTTP/embeddings
- [ ] G19: no library import of `cmd/trace`
- [ ] `CGO_ENABLED=1 go test ./evals/x0/... ./evals/p0x/... ./evals/honesty/...` + `./...` PASS

### S02 handoff
- [ ] Issue list in GATE-C-NOTES is concrete + measurement-driven (or empty with reason)
- [ ] If needed, thicken **upcoming** S02 `01-slice-hardening.md` Depends/input pointer only
- [ ] Do not rewrite Phase 00/01 or P02-00 / S01-00/01 `done` history

## Issue-list shape (for S02 — enforce in Notes)

Each issue entry should look like:

```text
- id: GC-NN
  severity: blocker|high|medium|low
  metric: <what failed / delta>
  evidence: <metrics path or GATE-C section>
  proposed_fix_surface: <package/CLI area or "defer">
  defer: true|false
```

## Exit criteria
- [ ] `REVIEW-NOTES.md` with severity findings
- [ ] blocker/high fixed inline or spawned (`02a`/`02b`)
- [ ] Confidence medium or high with residuals listed
- [ ] S02 upcoming prompts note issue-list pointer if needed
- [ ] Board status + Notes updated

## Minimal todos
- [ ] Compare 01 claims + GATE-C-NOTES vs repo evidence
- [ ] Re-run x0 + p0x + honesty + `./...`
- [ ] Fix or spawn; re-verify
- [ ] Board update
