# P09-S02-02 — Scope review notes (2026-08-16)

Independent review of DF-02/DF-04 discoverability vs `00-PLANNER.md` / `01-discoverability.md` locks + `P09-S02-01` Notes. Fresh session; claims re-verified in-repo (no implementer session shared).

## Verdict

**APPROVE** — no blocker / high findings. Confidence: **high**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Command is `trace tasks` (not colliding top-level `status`) | Pass — `root.go` `case "tasks"`; help documents `tasks [--goal]` |
| JSON array objects: `id`, `title`, `work_state`, `goal_id` (`null` when unset) | Pass — `taskListRow` + `*string` GoalID; CLI test asserts fields + null |
| Empty project → `[]` exit 0 | Pass — `TestTasksListAfterSeed` empty init; `make([]taskListRow, 0, …)` avoids JSON `null` |
| `--goal` → `ListTasksByGoalID`; unfiltered → `ListTasks` | Pass — `tasks.go` + CLI filter / unknown-goal `[]` |
| `ListTasks` in store; **no** new migration | Pass — `helpers.go`; no `011_*`; schema still through `010_capability_surface.sql` |
| G19 thin CLI (store only; no domain fork / no SQL in MCP) | Pass — `tasks.go` imports `store` only |
| Relative seed resolve under `-C` root; absolute unchanged | Pass — `resolveSeedPath`; `TestSeedImportRelativePathAgainstC` (cwd≠project + abs regression) |
| Help + `fixtures/x0/README.md` match live behavior | Pass |
| p0x/x0 abs seed still PASS | Pass — independent CGO1 re-run |
| **No** new MCP tools / daemon / HTTP | Pass — MCP still six tools (`trace_why`/`context`/`add`/`link`/`transition`/`review`); no `trace_tasks` |
| S01 Why/context with linked review not regressed | Pass — `TestWhyAndContextWithLinkedReview` in `./...` |
| Carry-forward: honesty A/B/C + Gate G, p0x, x0, `./...` | Pass (independent re-run) |

## Verify (independent re-run)

```text
CGO_ENABLED=0 go test ./internal/store/... ./evals/honesty/... -count=1                                    PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1         PASS
CGO_ENABLED=1 go test ./... -count=1                                                                      PASS
```

## Findings

None at blocker / high / medium.

### Low (no spawn)

- Store-layer empty `ListTasks` may return a nil Go slice (`len==0`); CLI always reifies to JSON `[]` via `make(..., 0, …)`. Covered at CLI; store test only asserts `len==0`.

### Nit

- Absolute seed paths are returned as given (no second `Abs`); matches lock “absolute unchanged.”

## Spawns

None.

## Residuals

- MCP list-tasks remains out of scope (S03 install-wire owns MCP config; keep discoverability CLI-primary unless a future row promotes).
- S01 residual unchanged: `plan_scope` ExactLookup still out; scope-only review expand path untested.

## S03 compatibility

S03 stubs remain compatible: install-wire should print/merge Cursor MCP config only; **do not** invent `trace_tasks` MCP unless promoted. Light Depends note added on S03 `01-install-wire.md`.

## Next

**P09-S03-00** (install-wire scope planner).
