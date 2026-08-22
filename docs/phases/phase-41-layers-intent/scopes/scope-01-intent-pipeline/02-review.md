# P41-S01-02 — Review (G9 intent pipeline)

## Metadata
- id: P41-S01-02
- todo_ids: [P41-S01-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter, security-and-hardening]
- mcps: [user-trace, user-codegraph]
- verification: mixed

## Objective

Fresh independent review of S01-01 G9 deliverable vs REMEDIATION-PLAN G9, GAP-REGISTRY G-009, M-001 moat, Laws 6–7/19.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — G9-I1–I6 acceptance map + touch-list
- [00-PLANNER.md](00-PLANNER.md) — locks + G9 vs G1 boundary
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [REMEDIATION-PLAN G9](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-009](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- Pre-ship baseline: [h9-intent-pipeline.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h9-intent-pipeline.md)
- P41-S01-00 live re-verify (2026-08-22): zero intent in `internal/retrieval/`; G8 shipped (`MaxLayer`, `layer_enrich.go`); G1 merge at `compiler.go:165–172`; `SearchOptions` intent-free pre-ship

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high → spawn 02a/02b below this row |
| Implement vs doc-revise | **Implement expected** — doc-revise acceptable only with ADR + §3 supersede |
| G9 vs G1 | Intent precedes channels; G1 merge unchanged |
| G9 vs G8 | Intent enriches FTS pre-Search; G8 layer admission post-candidate — verify no coupling regression |

## Review checklist

### A — G9 gap closure

- [ ] `internal/retrieval/intent.go` exists with `ExtractIntent`, `Intent`, `IntentInput`
- [ ] `Search` uses intent when `SearchOptions.Intent != nil` (`search.go`); nil preserves legacy raw-query path
- [ ] Compiler passes `IntentInput` at title + query Search sites (`compiler.go` ~157, ~168)
- [ ] §3 revised: intent stage shipped; semantic leg deferred (DR-NOSSEM); §2 Semantic cross-ref if added
- [ ] No zero-code silent defer — code **or** explicit doc supersede + ADR

### B — M-001 moat

- [ ] Intent requires task context on compile path — `TaskID` moat unchanged
- [ ] No query-only intent API without task (MCP direct Search with raw `q` only is OK — no task moat bypass on compile/explore)
- [ ] Task packet / loop path unchanged as primary moat
- [ ] Intent enriches retrieval — does not replace task UUID + gates

### C — Laws 6–7

- [ ] No dump API introduced
- [ ] `IntentSummary` capped if on packet (~256 chars per implement lock)
- [ ] Search limits unchanged by default (32 default, 64 hard cap in `search.go:13–17`)
- [ ] Keyword/token caps on `ExtractIntent` output (bounded — no unbounded Body dump into FTS)

### D — Law 19 + DR-NOSSEM

- [ ] Intent logic in `internal/retrieval/` library (not duplicated in `web/`)
- [ ] No semantic/vector channel (G9-I5)
- [ ] No LLM external calls; no new embedding imports
- [ ] `retrieval/doc.go` documents intent stage; no `semantic_match` reason_code

### E — G9 vs G1 boundary

- [ ] G1 `ContextOptions.Query` merge behavior preserved (`TestG1*` green)
- [ ] Title FTS still runs when query set (`TestG1TitleFTSStillRunsWithQuery`)
- [ ] Intent is pre-channel extraction, not duplicate post-merge
- [ ] `packet.go:200–201` Law 9 banner unchanged (not conflated with `IntentSummary`)

### F — G8 regression (G8 shipped S00)

- [ ] Default `MaxLayer=1` unchanged (`TestTaskContextAndBudgets` / G8-L1)
- [ ] G8-L2–L7 green; intent wiring did not break layer admission or trim order

### G — Tests

- [ ] G9-I1–I6 evidenced green (or doc-revise path documented)
- [ ] G1/G8 regression green
- [ ] `go test ./internal/retrieval/... ./internal/compiler/... ./internal/mcp/... ./cmd/trace/... -count=1` passes

### H — Rejects

- [ ] No embedding/vector intent
- [ ] No query-only intent on compile/explore paths
- [ ] §3 not left aspirational without supersede if code shipped
- [ ] No silent no-op `ExtractIntent` returning empty always

## Exit criteria

- [ ] APPROVE or REJECT with confidence in board Notes
- [ ] Zero open blocker/high without pending spawn
- [ ] Next **P41-S02-00** on APPROVE

## Next

`P41-S02-00` (on APPROVE)
