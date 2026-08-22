# P22-S04-03 — Implement: regression↔change attribution

## Metadata
- id: P22-S04-03
- todo_ids: [P22-S04-03]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Grep, Write]
- verification: automated

## Objective

**Identify which change is associated with a regression** (**C16**), with evidence-backed **`caused`** when possible. Board: status + notes only.

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md).

## References

- [00-PLANNER.md](00-PLANNER.md) — FINAL locks
- Live: `internal/domain/regressions.go` — create **`correlated`**; `SetRegressionAttributionCaused` fail-closed; no change link today
- S05 will expose CLI query — implement **domain/store helpers only** this row

## Live baseline

| Present | Absent |
|---------|--------|
| `regressions` + attribution ladder | `regression_associated_change` link |
| `RecordRegressionFromContradictedEffect` → effect has `change_id` | auto-link regression→change |
| `SetRegressionAttributionCaused` + evidence links | `ListRegressionsByChangeID` / `RegressionsForChange` |
| Schema **024** (from S04-01); compat **24** | **025+** |

## Locked defaults

| Item | Value |
|------|-------|
| Migration | **None** — use **`entity_links`** only (no ALTER `regressions`) |
| Compat | Stays **24** (forbid **025+**) |
| Link rel | **`regression_associated_change`**: `from_type=regression`, `to_type=change`, `to_id=change_id` |
| Constant | `RelRegressionAssociatedChange` in `internal/domain/service.go` |
| Default create | Still **`correlated`** on insert (P20 / Law 5) |
| Auto-link | `RecordRegressionFromContradictedEffect` inserts associated_change → `eff.ChangeID` (idempotent) |
| Manual API | **`AssociateRegressionWithChange(regressionID, changeID)`** — validates entities; InsertLinkOrIgnore |
| Caused path | **`SetRegressionAttributionCaused` unchanged**; test proves caused row **also** has change link when association set before caused |
| Query | **`ListRegressionsByChangeID(changeID)`** + alias **`RegressionsForChange`** |
| Evaluation regressions | No auto change link from evaluation source — agent/hook calls `AssociateRegressionWithChange` when known |
| CLI | **Out of scope** — S05-03 |

## Requirements

1. Add rel constant + seed import allowlist if needed.
2. Auto-associate on contradicted-effect regression create.
3. `AssociateRegressionWithChange` domain API.
4. `ListRegressionsByChangeID` store/domain query (join links or list links→regression).
5. Named tests including new + keepers.
6. Do **not** auto-set `caused` from association or correlation.

## Touch files

- `internal/domain/regressions.go`, `regressions_test.go`
- `internal/domain/service.go` (rel constant)
- `internal/store/regressions.go` (list-by-change query if SQL-heavy)
- `internal/domain/seed_import.go` (allow new rel if links exported)

## Named tests

| Test | Proves |
|------|--------|
| `TestRegressionLinkedToChangeCausedWithEvidence` | associate → hypothesized → confirmed → caused + link + evidence |
| `TestRegressionAutoLinkedFromContradictedEffect` | effect regression has associated_change without caused |
| `TestListRegressionsByChangeID` | query round-trip |
| `TestSetAttributionCausedFailClosedWithoutEvidence` | keeper |
| `TestCorrelationAndContradictionNeverAutoSetCaused` | keeper |
| `TestSetAttributionCausedRequiresConfirmedHypothesisAndEvidence` | keeper |

```bash
go test ./internal/domain/... -count=1 -run 'TestRegression|TestSetAttributionCaused|TestCorrelationAndContradiction|TestListRegressionsByChange'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
ls internal/store/schema/*.sql | wc -l  # still 24
```

## Exit criteria

- [ ] C16 true (change associated; caused when evidence policy satisfied)
- [ ] Named tests PASS; compat **24**; no **025+**
- [ ] Checklist C16 **not** boxed until S04-04 review
- [ ] Board Notes

## Minimal todos

- [ ] Link rel + auto-associate from effect
- [ ] Associate + list APIs
- [ ] Tests + keepers
- [ ] Board notes
