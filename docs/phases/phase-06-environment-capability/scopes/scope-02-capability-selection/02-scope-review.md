# P06 / S02 / 02 — Scope review (capability selection)

## Metadata
- id: P06-S02-02
- todo_ids: [P06-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of P06-S02-01 against S02-00 locks. APPROVE with evidence or spawn. Write [REVIEW-NOTES.md](REVIEW-NOTES.md).

## Session start
Agent → clarify if needed → Plan → execute (review).

## Review focus
- Ablation harness **`evals/capability`** / **`TestPlantedCapabilitySelectionAblation`** / **`schema-capability.json`** v1 + temp **`metrics-capability.json`**
- Planted tallies TP=3 / FN=0 / FP=0 / TN=1 → precision=1.0 / recall=1.0 (Pos UNAVAILABLE / UNKNOWN / selection filter; Neg clean AVAILABLE)
- S01 hooks consumed correctly (Upsert/Require/Missing + packet required+missing; mig `010` only — no schema fork)
- No catalog dump in packet `required_capabilities`
- Carry-forward Gate F/G/E / honesty A/B/C / p0x / x0 / Gate C `dry_run:false`
- No commercial multi-model theater; no weakening of prior gates
- G19: harness uses library APIs (no `cmd/trace` scrape); no product Go outside `evals/capability`

## Re-prove (must cite in REVIEW-NOTES)

```bash
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation
CGO_ENABLED=0 go test ./evals/honesty/... ./evals/replan/... ./evals/impact/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./... -count=1
```

## Light S03 note
Upcoming VERIFY re-prove command locked:  
`CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation`

## Exit criteria
- [ ] APPROVE or spawns
- [ ] REVIEW-NOTES.md written
- [ ] Light S03 Depends / re-prove command confirmed
