# P28-S06-01 — FM-01 / FR-P28-01 implementer

## Metadata
- id: P28-S06-01
- todo_ids: [P28-S06-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [user-codegraph, user-trace]
- verification: mixed
- hooks: []

## Objective

**FR-P28-01 / FM-01:** Reduce seed-import roster pin so BLOCKING discoveries can expand the executable backlog without relying solely on agent memory. Prefer guided promotion UX / `loop next` `promotion_candidates` surfacing / optional harness nudge post-import.

**Do not** ship autonomous spawn without human gate (FR-P28-D1). **No** daemon/HTTP.

## References

- [00-PLANNER.md](00-PLANNER.md)
- [RESIDUAL-AUDIT.md](../scope-00-residual-audit/RESIDUAL-AUDIT.md) — FM-01
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) §3 FM-01
- [TEST-MATRIX.md](../scope-01-integration-tests/TEST-MATRIX.md) — M-01 promotion
- Live: `internal/domain/promote.go`, `internal/loop/next.go`, `internal/loop/apply.go`
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Acceptance hint

Dogfood or integration: after seed import with orphan BLOCKING discovery, agent/path yields a new task (or explicit decline) without inventing UUIDs; document remaining human-gate if auto-spawn stays deferred.

## Preflight

```bash
cd /home/ali/Desktop/Trace
test -f docs/phases/phase-28-residuals-validation/scopes/scope-06-r6-fm-residuals/00-PLANNER.md
grep -n 'promotion_candidates\|PromoteBlocking' internal/loop/next.go internal/domain/promote.go | head
GOPROXY=direct go test ./internal/domain/ ./internal/loop/ -count=1 -run 'Promote|Next|Apply' 
```

## Suggested work (pick minimal path that meets acceptance)

1. Survey live promotion_candidates + seed-import pin behavior.
2. Product and/or harness: surface candidates post-import / strengthen nudge; optional CLI/JSON clarity.
3. Add or extend test/dogfood evidence under this scope folder (`FM01-NOTES.md` if dogfood).
4. Keep `go test ./internal/... ./cmd/trace/...` green if product touched.

## Out of scope

- FR-P28-D1 auto-spawn; Multitask rewrite; FM-02/04/07/08/09/10 (later rows); rewriting S00–S05 history

## Exit criteria

- [ ] Acceptance hint met with evidence in board Notes
- [ ] No daemon/HTTP; human gate preserved
- [ ] Next runnable **P28-S06-02**

## Todo updates

Status + notes on **P28-S06-01** only.

## Next

`P28-S06-02`
