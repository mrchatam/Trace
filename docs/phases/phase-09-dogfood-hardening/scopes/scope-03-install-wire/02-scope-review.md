# P09 / S03 / 02 — Scope review (install-wire / DF-03/DF-05)

## Metadata
- id: P09-S03-02
- todo_ids: [P09-S03-02]
- role: reviewer
- skills: [code-review-and-quality, silent-failure-hunter]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of Cursor install-wire: `trace install cursor` print/merge + DF-05 docs. Reject MCP list-tasks / daemon scope creep and S01/S02 regressions. Spawn remediations on blocker/high. Fresh subagent — do not share the implementer session.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-install-wire.md](01-install-wire.md)
- Implementer board Notes on `P09-S03-01`
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-03 / DF-05
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (review).

## Focus checklist (locked — P09-S03-00)

- [x] Command is **`trace install cursor`** (not a competing `mcp snippet` top-level)
- [x] Default print: merge-ready `{"mcpServers":{"trace":…}}` with `args: ["-C","${workspaceFolder}"]`
- [x] `--write` upserts only `trace`; other `mcpServers` preserved
- [x] Backup created before overwrite of existing file; invalid JSON fail-closed
- [x] Default `command` is `trace-mcp`; `--bin` overrides; `--mcp-json` for path override
- [x] Help + README + `experiments/ab-simple/PROTOCOL.md` document install + **DF-05** (open run folder as workspace)
- [x] Thin `cmd/trace` only (G19); **no** new MCP tools / daemon / HTTP / mig
- [x] S02 `trace tasks` still green; S01 Why/context-with-review still green
- [x] Carry-forward: honesty A/B/C + Gate G, p0x, x0, `./...`

## Verify commands (re-run independently)

```bash
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Role work

1. Diff claims (01 Notes + locks) vs repo evidence.
2. Findings by severity: blocker | high | medium | low | nit.
3. blocker/high: small inline fix **or** spawn `P09-S03-02a` / `02b` with full prompts immediately below this row.
4. Write [REVIEW-NOTES.md](REVIEW-NOTES.md); confidence medium/high with residuals listed.
5. Light-confirm S04 stubs still compatible (VERIFY will include install spot-check / DF-05 docs).
6. Board status + Notes; forward-only.

## Exit criteria
- [x] APPROVE (high, or medium with explicit residuals) **or** spawn with evidence
- [x] REVIEW-NOTES under this scope folder
- [x] Board status + Notes; next runnable is **P09-S04-00** when APPROVE (after this row)

## Out of scope
- Implementing S04 VERIFY / DR-HANDOFF
- Re-running full dogfood A/B portfolio (optional smoke only)

## Minimal todos
- [x] Diff claims vs repo; re-run locked tests
- [x] Write REVIEW-NOTES; APPROVE or spawn; update board
