# P16 / S03 / 01 — Install `-C` (DF-68)

## Metadata
- id: P16-S03-01
- todo_ids: [P16-S03-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Honor `-C` for Claude marker detect/install/uninstall. Keep Cursor reload tips. Board **status + Notes only**.

**Stop if `00-PLANNER.md` is not FINAL.**

## Locked defaults
| Item | Value |
|------|-------|
| Home | `cmd/trace/install.go` pass `root`; `internal/install` already uses `ProjectRoot` |
| DF-22/37 | Do not remove print/write reload tip; no PID kill |
| Forbidden | New MCP tools; YOLO; DF-70+ |

## Named tests
- `TestInstallClaudeHonorsDashC` — from a cwd **without** marker, `-C` project **with** marker: detect/install behave as if under project
- Keep `TestInstallCursor*` tip assertions

## Locked verify (minimum)
```text
CGO_ENABLED=0 go test ./internal/install/... ./cmd/trace/... -count=1 -run 'TestInstall'
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Exit criteria
- [ ] DF-68 named tests pass; detect reason uses `-C` root not `.` alone
- [ ] Reload tip keepers green; board → **P16-S03-02**
