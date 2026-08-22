# P13 / S03 / 02 — Scope review (Provenance schema) (FINAL)

## Metadata
- id: P13-S03-02
- todo_ids: [P13-S03-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S03 **DF-64 / DF-66 / DF-67** vs sibling **00-PLANNER** FINAL locks and live evidence. Confirm no silent drop of residuals. Spawn forward on blocker/high. Fresh subagent — do not share implementer session.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — FINAL
- Sibling [01-provenance-schema.md](01-provenance-schema.md)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- P12 edge-provenance REVIEW-NOTES; S01/S02 keepers
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute.

## Checklist (FINAL)

| # | Check | Evidence |
|---|--------|----------|
| 1 | DF-64: write rejects garbage; empty→EXTRACTED; read normalize empty→EXTRACTED | named tests + `ReplaceFileImports` |
| 2 | DF-64: mig **012** CHECK on `imports.provenance`; heal on migrate; embed ceiling **12** (no 013+) | schema file + compat |
| 3 | DF-64: empty cannot omitempty-hide on structural hop after normalize | Expand/Why or store+JSON path |
| 4 | DF-66: **wontfix** product analyzer/CLI setter documented; Law 5 store-fixture tests still green | `ANALYZER_CONTRIBUTION.md` + named tests |
| 5 | DF-67: **no** symbol honesty impl; residual explicit for VERIFY (`symstale/`) | diff + Notes |
| 6 | P12 provenance named tests green; packet SchemaVersion still `0.2`; G19 no new MCP/CLI provenance | tests + diff |
| 7 | Prefer zero analyzer/retrieval/compiler churn beyond necessity | diff |
| 8 | Carry-forward + Gate C `dry_run:false` intact; dry-run ≠ C/H/checklist | verify cmds |

## Review loop
1. Compare 01 claims + board Notes to repo evidence.
2. Findings by severity: blocker | high | medium | low | nit.
3. blocker/high → small inline fix **or** spawn `02a`/`02b` immediately below this row.
4. Re-verify until no open blocker/high without pending follow-up.
5. Write [REVIEW-NOTES.md](./REVIEW-NOTES.md); Confidence medium|high with residuals listed (never silent) — must include DF-66 wontfix + DF-67 out-of-bar.

## Verify (independent re-run)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/retrieval/... ./internal/compiler/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Todo updates
Status + Notes on this row; spawn rights if needed. Do **not** rewrite `done` prompts. On APPROVE, light-thicken S04 Depends if residuals need VERIFY callouts (upcoming only).

## Exit criteria
- [ ] Checklist 1–8 PASS or spawned
- [ ] APPROVE / REQUEST CHANGES + REVIEW-NOTES
- [ ] Confidence medium or high; DF-66/67 residuals explicit
- [ ] Board Notes; next **P13-S04-00** (unless spawn)

## Next
**P13-S04-00** (on APPROVE)
