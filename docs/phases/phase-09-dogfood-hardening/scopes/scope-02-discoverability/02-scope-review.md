# P09 / S02 / 02 — Scope review (discoverability / DF-02/DF-04)

## Metadata
- id: P09-S02-02
- todo_ids: [P09-S02-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of discoverability: `trace tasks` JSON list + seed path resolve against `-C`. Reject daemon/MCP scope creep and S01 regressions. Spawn remediations on blocker/high. Fresh subagent — do not share the implementer session.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-discoverability.md](01-discoverability.md)
- Implementer board Notes on `P09-S02-01`
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-02 / DF-04
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (review).

## Focus checklist (locked — P09-S02-00)

- [ ] Command is **`trace tasks`** (not a colliding `status` top-level)
- [ ] Stdout JSON array objects include **id, title, work_state, goal_id** (`null` when unset)
- [ ] Empty project → `[]` exit 0
- [ ] `--goal` filters correctly; unfiltered uses `ListTasks`
- [ ] No new migration; G19 thin CLI (no domain fork / no SQL in MCP)
- [ ] Relative seed path resolves under project root from `-C`; absolute paths unchanged
- [ ] Help + `fixtures/x0/README.md` match live behavior
- [ ] p0x/x0 still use abs seed and PASS
- [ ] **No** new MCP tools / daemon / HTTP
- [ ] S01: Why/context with linked review still green (do not assume skip-reviews)
- [ ] Carry-forward: honesty A/B/C + Gate G, p0x, x0, `./...`

## Verify commands (re-run independently)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Role work

1. Diff claims (01 Notes + locks) vs repo evidence.
2. Findings by severity: blocker | high | medium | low | nit.
3. blocker/high: small inline fix **or** spawn `P09-S02-02a` / `02b` with full prompts immediately below this row.
4. Write [REVIEW-NOTES.md](REVIEW-NOTES.md); confidence medium/high with residuals listed.
5. Light-confirm S03 stubs still compatible (install-wire owns MCP config; tasks remain CLI).
6. Board status + Notes; forward-only.

## Exit criteria
- [ ] APPROVE (high, or medium with explicit residuals) **or** spawn with evidence
- [ ] REVIEW-NOTES under this scope folder
- [ ] Board status + Notes; next runnable is **P09-S03-00** when APPROVE (after this row)

## Out of scope
- Implementing S03 install-wire
- Re-running full dogfood A/B portfolio (optional smoke only)

## Minimal todos
- [ ] Diff claims vs repo; re-run locked tests
- [ ] Write REVIEW-NOTES; APPROVE or spawn; update board
