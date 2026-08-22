# P28-S04-00 — Scope planner (product polish)

## Metadata
- id: P28-S04-00
- todo_ids: [P28-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated
- hooks: []

## Objective

Plan low-risk product/harness polish from Phase 27 residuals (R4, R5, partial R6). Lock file targets and exit criteria for honesty dedupe + P25-4 attestation automation.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — R4/R5
- [Phase 27 VERIFY-NOTES](../../../phase-27-protocol-measurement-graph-honesty/scopes/scope-03-verify/VERIFY-NOTES.md)
- [Phase 28 README](../../README.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked tasks

| ID | Fix | Files |
|----|-----|-------|
| R4 | Deduplicate BLOCKING orphan message | `internal/domain/seed_export_honesty.go`, `cmd/trace/seed.go` |
| R5 | P25-4 attestation: env/flag in score.sh or PROTOCOL template | `score.sh`, `PROTOCOL.md` |
| R6 (partial) | FM-02: optional stderr hint when export passes strict but graph thin (warn only) | honesty path — optional |

## Scope boundary

- Small only — no new INT scope
- `go test ./internal/... ./cmd/trace/...` must stay green
- No daemon/HTTP

## Planner gate

- [ ] `RESIDUAL-AUDIT.md` R4/R5 seeds present
- [ ] `01-implement.md` + `02-review.md` runnable

## Exit criteria

- [ ] Implementer prompt locked for fresh subagent
- [ ] Board row P28-S04-00 Notes cite file locks
- [ ] Next runnable **P28-S04-01**

## Todo updates

Status + notes on **P28-S04-00** only.

## Next

`P28-S04-01`
