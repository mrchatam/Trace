# P11 / S04 / 02 — Scope review (Capability upsert + hatch vs caps)

## Metadata
- id: P11-S04-02
- todo_ids: [P11-S04-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S04 (**DF-41, DF-51**). Fresh subagent. Compare claims + locks to live code/tests. Spawn `02a`/`02b` for blocker/high. Do not rewrite prior `done` history.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- Sibling [01-capability-upsert-hatch.md](01-capability-upsert-hatch.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-41, DF-51
- Live: `internal/domain/{capability,task_state}.go`; `cmd/trace/{capability,transition,help}.go`; `internal/mcp/{tools_parity,tools_write}.go`
- Prior: P10 DF-24/26; P11-S02 independence; P11-S03 no coupling

## Checklist (must all pass for APPROVE)

| # | Check |
|---|--------|
| 1 | DF-41: empty-ID re-declare same slug → same id + updated fields; explicit different-id slug clash still fails |
| 2 | DF-51: `AllowDoneWithoutReview` does **not** bypass missing-caps; check order caps→DONE; WARNING/docs mention missing-caps override |
| 3 | Gate G hatch + DF-24 fail-closed retained; no mig; no hatch→auto-missing-caps |
| 4 | G19 — no domain fork in CLI/MCP adapters |
| 5 | No forbidden architecture (daemon/HTTP/full-rebuild) |
| 6 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + P11-S01–S03 + Gate C `dry_run:false` |
| 7 | Board Notes accurate; planner row had no product Go |

## Verify (independent — re-run)

```bash
CGO_ENABLED=0 go test ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Prefer named asserts: `TestUpsertCapabilityBySlugUpdatesExisting`, duplicate-id clash, `TestAllowDoneDoesNotBypassMissingCaps`, `TestAllowDoneWarnsOnStderr`, `TestTransitionAllowDoneEmitsWarning` (or equiv locked names).

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals)
- [x] Board status + Notes; next **P11-S05-00** (unless spawn)
- [x] Write [REVIEW-NOTES.md](REVIEW-NOTES.md) on APPROVE / spawn
