# P09 / S01 / 02 — Scope review (retrieval review / DF-01)

## Metadata
- id: P09-S01-02
- todo_ids: [P09-S01-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of the DF-01 fix: `review` is a first-class ExactLookup entity so why/context work after reviews exist. Reject fail-soft skips. Spawn remediations on blocker/high. Fresh subagent — do not share the implementer session.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-retrieval-review.md](01-retrieval-review.md)
- Implementer board Notes on `P09-S01-01`
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-01
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (review).

## Focus checklist (locked — P09-S01-00)

- [ ] `lookupEntity` has `case "review"` calling `store.GetReview` (no new mig)
- [ ] Hit: Title from review title; Excerpt = Result if set else body excerpt; ReasonCode from caller
- [ ] Expand link path: `review_judges_task` / `review_judges_scope` reason codes mapped (not silent GraphNeighbor-only unless Notes justify)
- [ ] **No** fail-soft skip of unknown `"review"` (hard error paths for other unknown types may remain)
- [ ] Regression exists and fails closed without the case; passes with plant like honesty/D07
- [ ] Why on linked task succeeds; context/Expand includes review neighbor
- [ ] DONE/review promotion policy untouched; no CLI/MCP scope creep unless justified
- [ ] Carry-forward green: honesty A/B/C + Gate G, p0x, x0, `./...`
- [ ] DF-01 / DF-09 consequence addressed for multi-step review workflows

## Verify commands (re-run independently)

```bash
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Role work

1. Diff claims (01 Notes + locks) vs repo evidence.
2. Findings by severity: blocker | high | medium | low | nit.
3. blocker/high: small inline fix **or** spawn `P09-S01-02a` / `02b` with full prompts immediately below this row.
4. Write [REVIEW-NOTES.md](REVIEW-NOTES.md); confidence medium/high with residuals listed.
5. Light-confirm S02 stubs still compatible (discoverability does not depend on skipping reviews).
6. Board status + Notes; forward-only.

## Exit criteria
- [ ] APPROVE (high, or medium with explicit residuals) **or** spawn with evidence
- [ ] REVIEW-NOTES under this scope folder
- [ ] Board status + Notes; next runnable is **P09-S02-00** when APPROVE (after this row)

## Out of scope
- Implementing S02 `trace tasks` / S03 install-wire
- Re-running dogfood A/B portfolio (optional smoke only)

## Minimal todos
- [ ] Diff claims vs repo; re-run locked tests
- [ ] Write REVIEW-NOTES; APPROVE or spawn; update board
