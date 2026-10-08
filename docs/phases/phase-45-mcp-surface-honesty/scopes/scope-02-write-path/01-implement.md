# Phase 45 / S02 / Implement — Write-path validation

## Row `P45-S02-01` — **done** 2026-10-04. Issues [#150](https://github.com/mrchatam/Trace/issues/150), [#151](https://github.com/mrchatam/Trace/issues/151), [#154](https://github.com/mrchatam/Trace/issues/154), [#162](https://github.com/mrchatam/Trace/issues/162).

## 1. Typed not-found — `internal/store/errors.go` (new)

```go
var ErrNotFound = errors.New("store: entity not found")

type NotFoundError struct{ Entity, ID string }
func (e *NotFoundError) Error() string      // "goal 3f2…: not found"
func (e *NotFoundError) Unwrap() error      // → ErrNotFound
func (e *NotFoundError) Is(target error) bool

func notFound(entity, id string) error      // internal constructor
func IsNotFound(err error) bool             // works through fmt.Errorf("%w") chains
```

**`errors.Is(err, sql.ErrNoRows)` still returns true** for every converted getter — `Is` matches that target explicitly. This is what let ~15 getters be converted in one pass without auditing every caller of `sql.ErrNoRows` across the codebase.

Converted: `entities.go`, `entities_causal.go`, `capability.go`, `changes.go`, `cognitive.go`, `deliberation.go`, `scopes.go`.

Consumers updated to prefer `store.IsNotFound`: `internal/httpapi/helpers.go` and `internal/domain/scope.go` (`isNotFoundErr`), legacy string sniff retained as fallback.

## 2. Validation primitives

| File | Symbol | Purpose |
|------|--------|---------|
| `internal/domain/task_state.go` | `NormalizeWorkState(ws) (string, error)` | Closed 8-value vocabulary; **fail-closed** on anything else |
| `internal/domain/entity_lookup.go` (new) | `AddKind*` consts + `EntityAddKinds` | Closed create-kind vocabulary (goal, task, decision, assumption, discovery, plan-change, claim, evidence) |
| ″ | `NormalizeAddKind(kind) (string, error)` | Validates against `EntityAddKinds` |
| ″ | `ValidateEntityID(field, id) (string, error)` | UUID check; `field` is named in the error so `id` and `goal_id` are distinguishable |
| ″ | `(*Service).EntityExists(ctx, kind, id) (bool, error)` | Kind→getter dispatch via `addKindGetters()` |
| `internal/domain/scope.go` | `refuseSelfLink(rel, fromID, toID) error` | `:92` — rejects `from == to` |

## 3. `trace_add` — `internal/mcp/tools_write.go:69`

Input gains:

```go
ID      string `json:"id,omitempty"        jsonschema:"optional UUID; must not already exist unless replace=true"`
Replace bool   `json:"replace,omitempty"  jsonschema:"allow overwriting an existing entity with the same id (default false; refuse rather than silently overwrite)"`
```

Order of operations (`tools_write.go:74–120`):

1. `NormalizeAddKind(in.Kind)` — unknown kind refused
2. `ValidateEntityID("id", in.ID)` — non-UUID refused (#151)
3. if `id` supplied → `EntityExists`; if it exists and `!replace` → **refuse with the remedy in the message**; if `replace` → proceed, `created = false`
4. `goal_id`, when supplied → `ValidateEntityID("goal_id", …)` **then** `GetGoal` — an unresolvable goal is refused (#151)

Response (`tools_write.go:205`) gains **`created`**:

```json
{"ok": true, "type": "task", "id": "…", "created": false}
```

## 4. `trace_link` — self-reference (#162)

`internal/mcp/tools_write.go` rejects `from == to` **before any store lookup**. Independently, `refuseSelfLink` guards both `LinkBlocks` (`scope.go:239`) and `LinkAPIContract` (`scope.go:161`) so the rule holds for every adapter, not just MCP.

## 5. `trace_tasks` — filter validation (#154)

`internal/mcp/tools_parity.go:95,102` routes both `work_state` and `work_states[i]` through `NormalizeWorkState`:

```go
return nil, nil, fmt.Errorf("trace_tasks: work_states[%d]: %w", i, err)  // tools_parity.go:95
return nil, nil, fmt.Errorf("trace_tasks: work_state: %w", err)          // tools_parity.go:102
```

The index in the plural form is the whole point — `work_states[3]: "TYPO"` tells the caller exactly which element to fix.

## 6. Acceptance — all PASS

`internal/mcp/tools_write_test.go` (9 tests + 5 helpers):

| Test | Asserts |
|------|---------|
| `TestAddRejectsExistingIDWithoutReplace` | **#150** — refuses, no overwrite |
| `TestAddReplaceTrueReportsCreatedFalse` | `created:false`; opt-in path works |
| `TestAddRejectsNonUUIDID` | **#151** |
| `TestAddRejectsUnresolvableGoalID` | **#151** — well-formed UUID that resolves to nothing |
| `TestAddAcceptsResolvableGoalID` | Happy path still works |
| `TestTasksRejectsInvalidWorkState` | **#154** |
| `TestTasksAcceptsValidWorkState` | Valid values still accepted |
| `TestLinkRejectsSelfReference` | **#162** at the MCP boundary |
| `TestLinkBlocksDomainRejectsSelfLink` | **#162** at the domain boundary (adapter-independent) |

Helpers: `writeFixture`, `addServer`, `addCreated`, `seedGoal`, `seedTask`.

## 7. Boundaries honoured

- **No FK migration** for `tasks.goal_id` — stated non-goal. Validation at the write path; the schema is untouched.
- **No cascade** on `replace=true`.
- The guard lives in `internal/domain`, so `trace link` from the CLI and HTTP is covered too.

## Next

`P45-S02-02`