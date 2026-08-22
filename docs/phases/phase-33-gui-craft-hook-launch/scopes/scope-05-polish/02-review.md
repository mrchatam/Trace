# P33-S05-02 — Polish review

## Metadata
- id: P33-S05-02
- todo_ids: [P33-S05-02]
- role: reviewer
- skills: [code-review-and-quality, writing-guidelines]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review: docs primary story = **`trace gui`**; residuals closed or tracked for VERIFY. Spawn `P33-S05-02a`/`02b` only on blocker/high. Next **P33-S06-00**.

## References

- [00-PLANNER.md](00-PLANNER.md) · [01-implement.md](01-implement.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) Theme C
- Live: `docs/gui-quickstart.md`, `README.md`, `web/README.md`, `AGENTS.md`
- If touched: `internal/httpapi/addr_in_use.go`, `web/src/screens/Graph.tsx`
- Prior: S02/S04 REVIEW residuals

## Session start

Follow agent-loop-protocol. Fresh session vs implementer. Do not rewrite done S01–S04 prompts.

## Checklist

### Docs primary story
- [ ] `gui-quickstart` title + lead walkthrough use **`trace gui`** (not `./bin/trace serve` as first path)
- [ ] PATH install (`go install …/cmd/trace@…`) present; **≠** `trace install`
- [ ] `serve` / `./bin/trace serve` demoted to secondary/scripting — **not** deleted
- [ ] Multi-project / ports still fail-on-conflict; **no** auto-port claim
- [ ] No false hosted SaaS / always-on daemon claims
- [ ] `README.md` + `web/README.md` + `AGENTS.md` no longer teach serve as primary GUI launch
- [ ] Related / help pointers mention `trace gui` (serve OK alongside)

### Residuals (S02–S04)
- [ ] EmptyState Tasks CTA: inline primary verified or polish landed; else deferred in Notes
- [ ] Optional Explore **canvas** screenshot linked or deferred to VERIFY
- [ ] Optional addr-in-use dual-word (`gui`\|`serve`): landed + tests/docs match, **or** deferred with reason
- [ ] Optional craft literacy one-liner (chroma strip + labels) present or N/A

### Non-regression
- [ ] **No** palette/token rewrite undoing S04 forest-moss + kind chroma strip (`tokens.css` / chroma rules untouched unless trivial unrelated typo)
- [ ] No compose / route / budget changes
- [ ] Law 19 / loopback defaults unchanged
- [ ] If `FormatAddrInUseMessage` changed: `go test ./internal/httpapi/ -run 'FormatAddrInUse|IsAddrInUse' -count=1` green

## Artifact

Write short [`REVIEW.md`](REVIEW.md) in this folder: verdict PASS/FAIL, confidence, checklist evidence, findings table, thickenings for S06 if needed.

## Exit criteria

- [ ] Confidence medium/high; no blocker/high (or spawn remediation)
- [ ] Next **P33-S06-00**

## Todo updates

Status + notes on **P33-S05-02** only. Reviewer may thicken upcoming S06 prompts only.
