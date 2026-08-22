# P04 / S02 / 02 — Scope review (honesty escape-rate)

## Metadata
- id: P04-S02-02
- todo_ids: [P04-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of P04-S02-01 against S02-00 locks. APPROVE with evidence or spawn forward remediations. Reject weakening Paths A/B/C, missing S01 consumption, or inventing Gate G without harness evidence.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling [01-honesty-escape-rate.md](01-honesty-escape-rate.md) + [00-PLANNER.md](00-PLANNER.md)
- S01 APPROVE: [../scope-01-scope-review-layer/REVIEW-NOTES.md](../scope-01-scope-review-layer/REVIEW-NOTES.md)
- Live: `evals/honesty/`

## Session start
Agent → clarify if needed → Plan → execute (review).

## Review focus
- Named test **`TestHonestyEscapeRateGateGPrelim`** exists under **`evals/honesty`** (not a new package)
- Committed **`schema-gate-g.json`** v1; temp **`metrics-gate-g.json`** written + validated in test
- Escape formula locked: escapes=1, caught=2, attempts=3; hatch only on escape case
- **`TestHonestyFailClosedPlantedClaim`** Paths A/B/C intact; never sets `AllowDoneWithoutReview`
- S01 hooks: `LinkReviewScope` / `review_judges_scope`; `CountOpenResidualsByScope`; OPEN `POLICY_EXCEPTION` residual
- No product Go outside harness; no daemon/HTTP/embeddings; VerifiedFact absent
- Gate E / Gate C `dry_run:false` / p0x / x0 / `./...` intact
- S03 VERIFY Depends notes still name this test + schema path

## Required re-verify commands

```bash
CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./... -count=1
```

Confirm Gate C metrics under `docs/verification/gate-c-x0/` still `dry_run:false` N=3; means G1 > B0 (do not rewrite packs).

## Exit criteria
- [ ] Verdict + confidence + [REVIEW-NOTES.md](REVIEW-NOTES.md)
- [ ] No open blocker/high without spawn
- [ ] Board status + Notes; next **P04-S03-00** on APPROVE

## Minimal todos
- [ ] Diff claims vs evidence
- [ ] Re-run required tests + Gate C artifact spot-check
- [ ] Write REVIEW-NOTES; mark done or spawn
