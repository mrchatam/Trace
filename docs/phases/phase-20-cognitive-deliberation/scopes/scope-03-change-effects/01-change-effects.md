# P20-S03-01 — Implement change + effects

## Metadata
- id: P20-S03-01
- todo_ids: [P20-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Implement S03-00 FINAL: Change as first-class object (Git SHA + path refs), expected vs actual effects, comparison enum, contradiction hooks to Discovery/Hypothesis/reconsideration. **No blobs. Library-only.**

## Session start
Follow agent-loop-protocol. Unattended: execute after S03-00 is `done`. Board edits: **status + notes only**. Do not re-debate locks. **No CLI / MCP / loop / SelectNext / gitcli edits.**

## Locked defaults (from S03-00 FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Migration | **`017_changes_effects.sql`** (after 016). Additive. No ALTER on `files` / `vcs_*` / S02 tables. |
| Tables | `changes`, `change_paths`, `effects` |
| Paths | **Child table** — not JSON on `changes` |
| Comparison | `''` \| `supported` \| `partially_supported` \| `contradicted` — unknown **fail closed** |
| Git read | `ResolveChangePath` → `vcs.Repository.ShowFile` (iface only; tests: `vcs.Fake`) |
| Contradiction | Optional Discovery and/or Hypothesis link; FIRED `contradicted_effect` if `change_implements_decision`; **no** Regression; **no** auto-replan |
| Compat | Ceiling **17**; forbid `018+` |
| Out | CLI, MCP, loop apply keys, seed export, FTS, vcs index dual-write, raw CoT, blobs/patches |

Full SQL, rels, API signatures, status transitions, and 8192-byte caps: [00-PLANNER.md](00-PLANNER.md) FINAL section. Copy them; do not invent columns.

## Requirements

1. Store CRUD matching locked SQL (CHECK enums, unique `(change_id, path)` and `(change_id, dimension)`, indexes).
2. Domain wrappers with fail-closed `ErrValidation`:
   - `CreateChange` / `RecordChangeCommit` / `SupersedeChange`
   - `RecordExpectedEffect` / `RecordActualEffect`
   - `GetChange` / `ListChangesByTaskID` / `ListChangePaths` / `ListEffectsByChangeID`
   - `ResolveChangePath(ctx, changeID, path, repo vcs.Repository)`
3. Entity/rel/event constants on `internal/domain/service.go` (`change`, `effect`, listed rels, `change.recorded`, `effect.compared`, `effect.contradicted`).
4. Create requires TaskID + ≥1 normalized path. Empty SHA → `OPEN`; SHA present → `RECORDED`.
5. Actual requires existing expected dimension + comparison in the 3-value enum.
6. Contradicted: event; optional PLAN_AFFECTING Discovery; Hypothesis via existing `CreateHypothesis` (not Discovery-as-hypothesis); linked decisions get FIRED `contradicted_effect`.
7. Bump embed/compat tests that still expect **16** / forbid **017**.
8. Extend `TestNoSourceContentColumns` spirit: new tables must not grow blob/patch/diff/content columns.

## Named tests (must exist and pass)

See 00-PLANNER list (14 names). Minimum proofs:

- Create with SHA + paths; pragma/scan shows no blob columns; paths have no content field
- Missing TaskID or zero paths fail closed
- Expected → actual `supported` sets comparison; status becomes `COMPARED` when all dimensions compared
- Actual without expected dimension fails; unknown comparison fails
- Contradicted + `CreateHypothesis` inserts hypothesis + `hypothesis_explains_effect`, **zero** extra Discovery unless `EmitDiscovery`
- `EmitDiscovery` → PLAN_AFFECTING + `discovery_from_contradicted_effect`
- Linked decision: FIRED reconsideration trigger `contradicted_effect`; Decision row remains
- No `regressions` table insert; `SelectNext` / `ApplyDeliberationTransition` not called
- Parent `parent_change_id` chain; parent not auto-SUPERSEDED
- `ResolveChangePath` returns `vcs.Fake` bytes; SQLite has no copy of those bytes
- Empty `git_commit` cannot resolve
- Reason/expected/actual > 8192 bytes fail closed

## Likely touch points

- `internal/store/schema/017_changes_effects.sql` (**new**)
- `internal/store/changes.go` + `changes_test.go` (**new**)
- `internal/domain/changes.go` + `changes_test.go` (**new**)
- `internal/domain/service.go` (constants only)
- `evals/compat/compat_test.go` + `evals/compat/doc.go` (16→17, 017 present, forbid 018+)
- `internal/store/production_hardening_test.go`
- `internal/store/deliberation_test.go` (EmbedExpected 17)
- `internal/store/store_test.go` (`TestOpenCreatesDBAndMigratesIdempotent` include 17)

Do **not** touch: `internal/loop`, `cmd/trace`, `internal/deliberation/select.go`, `internal/mcp`, `internal/gitcli`.

## Proof commands

```bash
go test ./internal/store/ -count=1 -run 'TestOpenCreatesDBAndMigratesIdempotent|TestMigrationStatusReportsEmbedMax|TestDeliberationStateRoundtrip|TestChange|TestNoSourceContentColumns'
go test ./internal/domain/ -count=1 -run 'TestCreateChange|TestRecordExpected|TestRecordActual|TestUnknownEffect|TestContradicted|TestParentChange|TestResolveChange|TestOversizedEffect'
go test ./internal/deliberation/...
go test ./cmd/trace -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory|TestLoopStatusSaturatedByZeroDeltaAndMaxIteration|TestHelpIncludesLoopNext'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

P19 Loop keepers must stay green (this row must not change loop).

## Todo updates
Status + notes only on `P20-S03-01`. Next after green: `P20-S03-02`.

## Exit criteria

- [ ] Migration 017 applied; EmbedExpected/compat ceiling **17**
- [ ] All 14 named tests green
- [ ] No source blobs / patches in SQLite (Law 1)
- [ ] Schema matches S03-00 locks (`change_paths` table, not JSON)
- [ ] Contradiction does not create Regression or auto-replan
- [ ] No daemon / CoT blobs / hosted MCP / CLI in this row
