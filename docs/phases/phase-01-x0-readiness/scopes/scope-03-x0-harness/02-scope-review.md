# P01 / S03 / 02 — Scope review (Agent X0 harness (CLI))

## Metadata
- id: P01-S03-02
- todo_ids: [P01-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (Agent X0 harness). Findings by severity; small fixes or spawn `a`/`b` pairs. Forward-only. May thicken **upcoming** S04/S05 if metrics paths or CLI-vs-MCP boundaries matter for VERIFY.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop)
- Sibling [01-x0-harness.md](01-x0-harness.md) + board Notes
- [docs/init/I_BENCHMARK_PLAN.md](../../../../init/I_BENCHMARK_PLAN.md) — Experiment X0
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) — Phase 1 gate ≠ Gate C
- Live: `evals/x0`, `evals/p0x`, `evals/honesty`, `fixtures/x0`, `cmd/trace` why/context

## Session start
Agent → clarify → Plan → review (fresh subagent).

## Review focus (checklist)

### Conditions & CLI
- [ ] **B0** metrics exist; harness does **not** call `trace why` / `trace context` for B0
- [ ] **G1** metrics exist; harness **does** call `why task <TaskUUID>` and `context <TaskUUID> --format json` (exit 0)
- [ ] Conditions labeled correctly in JSON (`condition`: `B0` / `G1`)

### Schema & dry-run bar
- [ ] Committed `evals/x0/schema.json` (schema_version 1); tests validate emits
- [ ] Temp `metrics-b0.json` + `metrics-g1.json` (or equivalent) with `dry_run: true`
- [ ] No Gate C scoring / “G1 beats B0” product claim in Notes or code comments as exit criteria
- [ ] Honesty **not** merged into X0 metrics (H5 stays `evals/honesty`)

### Corpus / hygiene
- [ ] Uses `fixtures/x0` + **absolute** seed path; stable Task UUID `22222222-…`
- [ ] No committed `.trace/` under fixtures/evals; no committed live metrics
- [ ] `evals/p0x` package untouched (still 7/7 path)

### Laws / regression
- [ ] No MCP requirement creep (DR-AGENT: X0 CLI-capable)
- [ ] No daemon/HTTP/embeddings
- [ ] G19: no library import of `cmd/trace`
- [ ] `CGO_ENABLED=1 go test ./evals/x0/...` PASS
- [ ] `CGO_ENABLED=1 go test ./evals/p0x/... ./evals/honesty/...` + `./...` PASS

### Cross-scope (upcoming only)
- [ ] If metrics paths or “X0 without MCP” need clarifying for S04/S05, thicken **upcoming** prompts only
- [ ] Do not rewrite S01/S02 `done` history

## Exit criteria
- [ ] Findings recorded (REVIEW-NOTES.md preferred)
- [ ] blocker/high fixed or spawned
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated

## Minimal todos
- [ ] Compare 01 claims + Notes vs repo evidence
- [ ] Run x0 + p0x + honesty + `./...` verification commands
- [ ] Fix or spawn
- [ ] Re-verify → board update
