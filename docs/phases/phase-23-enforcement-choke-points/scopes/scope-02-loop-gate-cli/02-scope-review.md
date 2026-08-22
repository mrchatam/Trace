# P23-S02-02 — Review loop gate CLI

## Metadata
- id: P23-S02-02
- todo_ids: [P23-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- verification: automated

## Objective
Independent review: exit-code contract (0/1/2), `trace.loop.gate.v1` JSON stability for harness hooks, thin adapter over S01 — no policy fork in `cmd/trace`.

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- S02-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- S02-01 deliverable: [01-loop-gate-cli.md](./01-loop-gate-cli.md)
- S01 library: `internal/loop/gate.go` (must be unchanged by S02)

## Session start
Follow agent-loop-protocol. Fresh reviewer context. Board edits: **status + notes only**.

## Keeper tests (must re-run — all green)

```bash
go test ./cmd/trace -run 'TestLoopGate|TestLoopNext|TestLoopApply|TestLoopStatus|TestHelpIncludesLoop'
go test ./internal/loop/... -run Gate
```

## Evidence to collect

| Check | Evidence |
|-------|----------|
| Thin adapter | `grep EvaluateGate` in `cmd/trace/loop.go` — single call site in `cmdLoopGate` |
| No policy fork | No `SelectNext` / `BuildPolicyInputs` in `cmd/trace/loop.go` gate path |
| Exit 0 | Allowed path returns 0 + JSON with `allowed=true` |
| Exit 1 | Blocked path returns **1** (not global `exitUsage` constant if that equals 1 — verify blocked uses literal policy exit) |
| Exit 2 | Missing `--task`, bad `--for`, store errors → 2, no valid gate JSON on stdout |
| Schema version | Every exit 0/1 stdout has `schema_version: trace.loop.gate.v1` |
| Top-level lift | Blocked responses duplicate `recommended_phase` / `reason_code` at envelope root from `violations[0]` |
| Empty violations | Allowed responses have `"violations": []` not null |
| Stderr on block | Exit 1 writes one-line hint (violation message) to stderr |
| Default `--for` | Omitted flag behaves as `edit` |
| Help | `trace help` lists gate subcommand + exit semantics |
| S01 untouched | `internal/loop/gate.go` diff empty for S02 scope |
| Named tests | All 14 tests from 01 prompt present and passing |
| No MCP | No new MCP tool registrations |

## Review checklist

- [ ] **Blocker:** Duplicate SelectNext / policy logic in `cmd/trace`
- [ ] **Blocker:** Blocked gate returns exit 2 or exit 0 (wrong harness contract)
- [ ] **Blocker:** Exit 0/1 path missing JSON or wrong `schema_version`
- [ ] **Blocker:** Missing named CLI tests from 01 prompt
- [ ] **Blocker:** Loop keeper tests regressed
- [ ] **High:** Top-level `recommended_phase` / `reason_code` not lifted on block
- [ ] **High:** Allowed response omits violations array or uses null
- [ ] **High:** Default `--for` not `edit`
- [ ] **High:** Exit 2 still emits parseable gate JSON (confuses hooks)
- [ ] **Medium:** Help text missing exit-code line
- [ ] **Medium:** Changes to `internal/loop/gate.go` (policy belongs in S01)
- [ ] **Low:** Cursor/harness-specific code in cmd
- [ ] **Nit:** Inconsistent JSON formatting vs `loop next`

## S05 handoff verification

Confirm harness install can call without changes:

```bash
trace loop gate --task "$TRACE_TASK_ID" --for edit
# exit 0 → parse stdout JSON, proceed with edit
# exit 1 → parse stdout JSON, surface recommended_phase to agent
# exit 2 → do not parse stdout; surface stderr error
```

JSON field parity with ENFORCEMENT.md status `violations[]` element (same shape as S01 `Violation`).

## Spawn policy

- **blocker/high:** inline fix if ≤10 lines and zero policy change; else spawn `P23-S02-02a` implement + `02b` review immediately below this row
- **medium:** prefer spawn unless trivial typo
- Do not rewrite S02-00/S02-01 `done` prompts

## Exit criteria

- [ ] No open blocker/high without pending forward row
- [ ] Confidence **medium** or **high** with command output in Notes
- [ ] Residual risks listed if medium (e.g. S04 status violations reuse deferred)
- [ ] APPROVE or spawn documented on board

## Minimal todos

- [ ] Re-run keeper tests; paste pass summary in Notes
- [ ] Walk exit-code table against `cmdLoopGate` line-by-line
- [ ] Verify all 14 named tests exist in `loop_test.go`
- [ ] Confirm `internal/loop/gate.go` diff empty for S02
- [ ] Spot-check blocked JSON against S02-00 example shape
- [ ] Set row done with confidence + residuals
