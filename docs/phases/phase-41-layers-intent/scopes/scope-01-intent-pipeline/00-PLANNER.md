# P41-S01-00 — Scope planner (G9 intent pipeline)

## Metadata
- id: P41-S01-00
- todo_ids: [P41-S01-00]
- role: planner
- skills: [planning-and-task-breakdown, context-engineering, api-and-interface-design]
- mcps: [user-trace, user-codegraph]
- verification: automated

## Objective

Lock S01 **G9** against live repo: intent pipeline (G-009). Thicken `01-implement.md` + `02-review.md` with file targets, acceptance map, and implement-vs-doc-revise decision. **No product code in this row.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md) — P41-00 Q1 resolution (**implement bounded**)
- [REMEDIATION-PLAN §2 G9](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-009](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)
- [RETRIEVAL_AND_CONTEXT.md §3](../../../../RETRIEVAL_AND_CONTEXT.md) — aspirational pipeline diagram
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live anchors (verified 2026-08-22 P41-00):
  - `internal/retrieval/doc.go` — no intent stage; DR-NOSSEM semantic forbidden
  - `internal/retrieval/engine.go`, `search.go` — hybrid lookup entry points
  - `internal/compiler/compiler.go:158–165` — G1 `ContextOptions.Query` merge (prerequisite shipped)
  - `internal/compiler/packet.go:200–201` — "project intent" = Law 9 trust label only
  - Evidence: [h9-intent-pipeline.md](../../../../../experiments/runs/2026-08-22-p38-s01-651/evidence/h9-intent-pipeline.md)

## Session start

Follow agent-loop-protocol Session start. Unattended: INTAKE + P41-00 locks are authority.

## Locked defaults (FINAL — P41-00)

| Item | Value |
|------|-------|
| GAP ids | G-009 |
| Verdict | **Accept** per REMEDIATION-PLAN G9 |
| P41-00 decision | **Implement bounded rule-based intent** in `internal/retrieval/`; doc-revise only if S01-01 blocked |
| Scope boundary | Intent **stage before** retrieval channels — not LLM inference, not semantic/vector (G-004a) |
| Intent shape | `Intent{Keywords[], EntityHints[], Scope, Source}` — deterministic from task fields + optional query |
| Wiring | `ExtractIntent` → enriches `Search` / compile query tokenization — merges into moat, never replaces task packet |
| M-001 | Intent channel merges into task loop — never query-only replacement |
| Law 6–7 | Progressive caps on any enriched retrieval; no dump API |
| Doc update | Revise RETRIEVAL_AND_CONTEXT §3 to mark semantic leg **deferred (DR-NOSSEM)**; document shipped intent stage |
| Out | LLM intent extraction; embedding channel; query-only intent path; full §3 pipeline including semantic |

## Live repo gap (re-verified P41-00)

| Check | Shipped | Gap |
|-------|---------|-----|
| `internal/retrieval/` intent code | **Zero** grep matches | Need `ExtractIntent` stage |
| Packet intent field | Absent | Optional `intent` summary on packet or search opts |
| §3 semantic leg | Documented | Forbidden — document defer in §3 revise |
| G1 query merge | Shipped Phase 39 | Complementary — G9 precedes channel selection |

## Accept / reject (G9)

| Decision | Item |
|----------|------|
| **Accept** | `internal/retrieval/intent.go` — rule-based `ExtractIntent(task, query)` |
| **Accept** | Wire intent into `Search` query building (keyword boost / token expansion) |
| **Accept** | Optional intent summary on context packet or search response metadata |
| **Accept** | Tests G9-I1–I6 (see thickened `01-implement.md`) |
| **Accept** | §3 doc revise: semantic leg deferred; intent stage shipped |
| **Accept (fallback)** | Doc-revise only: mark §3 aspirational + supersede with ADR — if implement blocked |
| **Reject** | LLM/embedding intent extraction |
| **Reject** | Intent as query-only replacement (no task_id) |
| **Reject** | Silent no-op — must ship code **or** explicit doc supersede |

## Must lock for S01-01 (delivered in thickened 01-implement)

1. Touch-list: intent.go → search wiring → optional packet field → §3 doc → tests.
2. Six acceptance tests G9-I1–I6.
3. Clarify G9 vs G1 boundary in doc comment.

## Exit criteria

- [ ] `01-implement.md` + `02-review.md` runnable with file targets + G9 accept map
- [ ] SCOPE-TODOS updated
- [ ] Board row → `done` with Notes

## Next

`P41-S01-01`
