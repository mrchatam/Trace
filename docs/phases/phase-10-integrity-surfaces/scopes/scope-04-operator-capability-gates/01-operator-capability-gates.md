# P10 / S04 / 01 — Operator + capability gates

## Metadata
- id: P10-S04-01
- todo_ids: [P10-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-17, DF-18, DF-24, DF-26, DF-31** per sibling **00-PLANNER** FINAL locks (2026-08-16). Enforce operator DONE flag; invalidate sticky PASS on reopen; fail-closed missing-capability gate; louder `--allow-done`; usable `capability missing` without `--task`. Inherit S01–S03 (do not re-implement). Keep carry-forward gates green. **No new migration** (default).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-17/18/24/26/31
- [experiments/ab-operator-gate/results/G1/PROBE-AGENT-DONE.md](../../../../../experiments/ab-operator-gate/results/G1/PROBE-AGENT-DONE.md)
- [phase README](../../README.md)
- Live: `internal/domain/task_state.go`, `service.go` (`TransitionOptions`), `capability.go`, `review.go`; `cmd/trace/transition.go`, `capability.go`, `seed.go`, `help.go`; `internal/mcp/tools_write.go`; `evals/honesty/`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate. **Honesty Path C supersession** (operator flag) is intentional and documented in 00-PLANNER.

## Locked defaults (FINAL — P10-S04-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/domain` + thin CLI/MCP adapters only |
| Migration | **None** |
| DF-17 | `AllowOperatorDone` / `--as-operator` / MCP `as_operator`; Actor string ≠ auth; →DONE = hatch **OR** (PASS ∧ operator flag) |
| DF-18 | From DONE → PENDING\|STALE: linked PASS reviews → `UNCERTAIN` before/with state change |
| DF-24 | Fail-closed on **any** transition if `MissingCapabilities` nonempty; override `AllowMissingCapabilities` / `--allow-missing-caps` / `allow_missing_caps` |
| DF-26 | Loud CLI stderr WARNING + MCP success `"warning"` when hatch used; keep Gate G hatch |
| DF-31 | Clear usage + `trace tasks` hint; MCP clear `task`/`task_id` error |
| Honesty | A/B unchanged; Path C sets `AllowOperatorDone: true`; Gate G hatch unchanged |
| Carry-forward | honesty A/B/C + G; E/F; ablation; H; compat; p0x; x0; S01–S03 tests; Gate C `dry_run:false` |
| Forbidden | New mig; daemon/HTTP/embeddings; Actor allowlists; removing hatch; Mode-B rewrite; board spawn |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/domain/service.go` | Extend `TransitionOptions` with `AllowOperatorDone`, `AllowMissingCapabilities` |
| `internal/domain/task_state.go` | Cap gate; DONE = hatch \|\| (PASS ∧ operator); reopen invalidate PASSes; payload may record flags |
| `internal/domain/review.go` | Reuse `SetReviewResult` or small helper for invalidate-on-reopen |
| `internal/domain/*_test.go` | DF-17/18/24 cases (deny/allow/reopen/cap) |
| `cmd/trace/transition.go` | Flags `--as-operator`, `--allow-missing-caps`; DF-26 stderr on `--allow-done` success |
| `cmd/trace/capability.go` | DF-31 usage/hint when `--task` missing |
| `cmd/trace/seed.go` | Wire `as_operator` (+ optional `allow_missing_caps` if useful) |
| `cmd/trace/help.go` | Louder DONE / allow-done / as-operator / missing-caps notes |
| `cmd/trace/*_test.go` | Flag + warning + capability-missing usage |
| `internal/mcp/tools_write.go` (+ server descriptions) | `as_operator`, `allow_missing_caps`; DF-26 `warning` on hatch; tool text |
| `internal/mcp/tools_parity.go` | DF-31 missing-without-task error clarity |
| `internal/mcp/*_test.go` | Transition flags + warning; capability missing param error |
| `evals/honesty/honesty_test.go` (+ doc if needed) | Path C `AllowOperatorDone: true`; Gate G still hatch-only; update reject-reason asserts if copy changes |

## Role work

1. TDD domain: PASS without operator flag → DONE rejected; with flag → OK; hatch → OK.
2. TDD DF-18: DONE→PENDING (or STALE), old PASS invalidated, DONE rejected until new PASS + operator.
3. TDD DF-24: require UNAVAILABLE/UNKNOWN cap → transition blocked; override succeeds.
4. Wire CLI/MCP/seed flags; DF-26 warning; DF-31 usage.
5. Update honesty Path C + any domain tests that DONE after PASS without operator.
6. Run locked verify suite; board **status + Notes only** (cite test names + DF checklist).

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: MCP/CLI DONE after PASS without `as_operator` fails; reopen then DONE fails until new PASS; missing caps block; `--allow-done` stderr warning; `trace capability missing` without `--task` prints tasks hint.

## Exit criteria
- [ ] DF-17: operator flag enforced; Actor spoof insufficient; hatch still works
- [ ] DF-18: leave DONE invalidates linked PASS; sticky PASS gone
- [ ] DF-24: missing caps fail-closed on transitions; override wired
- [ ] DF-26: loud warning on hatch success (CLI + MCP)
- [ ] DF-31: usable missing-`--task` / missing-`task_id` errors
- [ ] Honesty A/B/C + Gate G green (Path C uses `AllowOperatorDone`)
- [ ] No new mig; G19 intact; nine MCP tools retained
- [ ] Carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P10-S04-02**

## Out of scope
- S05 VERIFY aggregation / DR-HANDOFF
- DF-28 handoff SoT; authn; `DONE→IN_PROGRESS` edge
- Re-opening S01–S03 product work

## Todo updates
Implementer: **status + notes only**. Record test names + DF checklist evidence. No spawning; no rewriting upcoming prompts.

## Minimal todos
- [ ] DF-17 domain + CLI/MCP/seed + tests; honesty Path C flag
- [ ] DF-18 invalidate-on-reopen + tests
- [ ] DF-24 cap gate + override + tests
- [ ] DF-26 louder hatch warning
- [ ] DF-31 capability missing UX
- [ ] Locked verify suite; board Notes
