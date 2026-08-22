# P13 / S02 / 02 — Scope review (Packet honesty residuals) (FINAL)

## Metadata
- id: P13-S02-02
- todo_ids: [P13-S02-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S02 **DF-61, DF-62, DF-63, DF-65** vs sibling **00-PLANNER** FINAL locks and live evidence. Spawn forward on blocker/high. Fresh subagent — do not share implementer session.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — FINAL
- Sibling [01-packet-honesty-residuals.md](01-packet-honesty-residuals.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- P12 keepers; S01 REVIEW-NOTES (resolve reuse)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Checklist (FINAL)

| # | Check | Evidence |
|---|--------|----------|
| 1 | DF-61: `stale_total` + `stale_truncated` (or equivalent loud total); cap ≤8; MD shows total when truncated | named test + JSON/MD |
| 2 | DF-62: honesty universe = pre-trim file items; trim-dropped disk-stale file does **not** yield false-fresh null | named test |
| 3 | DF-63: when `candidates_capped`, `items_total` ≥ true L1-admissible universe (not post-cap ≤64 alone); MD `items=k/t` | named test |
| 4 | DF-65: TaskContext/ExpandContext carries import-hop `edge_provenance` via Expand on file seeds; **no** compiler path-join reimplementation | named test + Expand call site |
| 5 | P12 keepers green; Law 18 untouched; SchemaVersion `0.2`; false-fresh on I/O miss | tests + code path |
| 6 | No mig / analyzer rewrite / path-align / new MCP; G19 | diff |
| 7 | S01 Expand/Why import tests still green | `-run` import |
| 8 | Carry-forward + Gate C `dry_run:false` intact; dry-run ≠ C/H/checklist | verify cmds |

## Review loop
1. Compare 01 claims + board Notes to repo evidence.
2. Findings by severity: blocker | high | medium | low | nit.
3. blocker/high → small inline fix **or** spawn `02a`/`02b` immediately below this row.
4. Re-verify until no open blocker/high without pending follow-up.
5. Write [REVIEW-NOTES.md](./REVIEW-NOTES.md); Confidence medium|high with residuals listed (never silent).

## Verify (independent re-run)

```bash
CGO_ENABLED=0 go test ./internal/compiler/... ./internal/retrieval/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Todo updates
Status + Notes on this row; spawn rights if needed. Do **not** rewrite `done` prompts.

## Exit criteria
- [ ] Checklist 1–8 PASS or spawned
- [ ] APPROVE / REQUEST CHANGES + REVIEW-NOTES
- [ ] Confidence medium or high with residuals explicit
- [ ] Board Notes; next **P13-S03-00** (unless spawn)

## Next
**P13-S03-00** (on APPROVE)
