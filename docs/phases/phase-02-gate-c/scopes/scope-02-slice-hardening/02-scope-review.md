# P02 / S02 / 02 — Scope review (slice hardening)

## Metadata
- id: P02-S02-02
- todo_ids: [P02-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of P02-S02-01 hardening against the **S01 GC-NN issue list**. Confirm GC-01/GC-02 addressed (or explicitly re-deferred with reason), GC-03/GC-04 still deferred, and honesty + p0x + x0 bars intact. No scope creep into Phase 03 progressive planner. Fresh session ≠ implementer.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling [01-slice-hardening.md](01-slice-hardening.md) locked defaults + acceptance
- Issue SoT: [`../scope-01-x0-gate-c/GATE-C-NOTES.md`](../scope-01-x0-gate-c/GATE-C-NOTES.md)
- S01 review: [`../scope-01-x0-gate-c/REVIEW-NOTES.md`](../scope-01-x0-gate-c/REVIEW-NOTES.md)

## Session start
Follow agent-loop-protocol Session start (Agent → clarify → Plan → review).

## Review focus (checklist)

| Focus | Pass bar |
|-------|----------|
| GC-01 | `Why`/`TaskContext` on task seed surfaces discovery + plan_change with `discovery_causes_plan_change`; regression tests exist; Expand max depth still 1..2; no fake task↔discovery GT edge unless spawned later |
| GC-02 | Agent-facing `fixtures/x0/README.md` lacks UUID oracle; evaluator map / seed-SoT documented; `pins.md` hash + Agent brief updated; guard test preferred |
| Deferrals | GC-03 (live-model packs) and GC-04 (N/variance) remain deferred — no silent implement |
| Bars | Fresh `CGO_ENABLED=1 go test ./evals/honesty/... ./evals/p0x/... ./evals/x0/... ./... -count=1` PASS |
| Honesty | Paths A/B/C unchanged (no escape hatch in honesty proof) |
| P0-X | 7/7 still green |
| Scope creep | No daemon/HTTP/embeddings; no MCP-required Gate C; no progressive planner product; Mode-B packs not rewritten to fake q3 pass |
| Verdict | Gate C **Go** text not silently altered; this row is slice harden, not a new Gate C run |

## Findings severity
blocker | high → inline fix or spawn `02a`/`02b` immediately below this row.  
medium → prefer spawn unless trivial.  
Re-verify until no open blocker/high without pending follow-up.

## Exit criteria
- [x] `REVIEW-NOTES.md` written (verdict + confidence medium/high + evidence)
- [x] Each claimed fix maps to GC-01 or GC-02 (or documented re-defer)
- [x] Bars intact (honesty, p0x, x0, no daemon/embeddings)
- [x] Spawns inserted on board if needed; else none
- [x] Board status + Notes

## Minimal todos
- [x] Claims vs repo evidence (GC-01/02 + deferrals)
- [x] Re-run honesty / p0x / x0 / `./...`
- [x] REVIEW-NOTES + board update
