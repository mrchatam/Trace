# P16 / S01 / 01 — MCP project isolation (DF-76)

## Metadata
- id: P16-S01-01
- todo_ids: [P16-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Fail-closed MCP store open when `.trace/` does not already exist so `project=` cannot mint a fresh AUTO_ALLOWED DB. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is still a stub (not FINAL).**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — required FINAL
- [phase README](../../README.md)
- Live: `internal/mcp/project.go`, `internal/store/open.go`, `internal/mcp/mcp_test.go`

## Session start
Follow agent-loop-protocol. Do not re-debate FINAL locks.

## Locked defaults (inherited — 00-PLANNER may tighten)
| Item | Value |
|------|-------|
| Home | `internal/mcp` thin adapter (G19). Prefer MCP-only exist-check; CLI Open still mkdir |
| Fail-closed | Virgin dir + `project=` → CallTool error; **no** `.trace/` created |
| Isolation | DENIED on store A does not apply to initialized store B; virgin C must not become B |
| Forbidden | New MCP tools; install/decide MCP; YOLO; ImpactWalk; rewriting done history; DF-75/77/78 (S02) |

## Named tests (required — 00 may rename)
| Test | Intent |
|------|--------|
| `TestMCPVirginProjectDoesNotCreateStore` | CallTool `project=` on empty dir → error; no `.trace/` |
| `TestMCPProjectOverrideDeniedDoesNotEscapeViaFreshRoot` | DENIED on workspace store; virgin override does not succeed / AUTO_ALLOW |

## Locked verify (minimum)
```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/store/... ./internal/domain/... -count=1 -run 'TestMCPVirgin|TestMCPProject|TestMCPAssert|TestToolNamesRegistered'
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Exit criteria
- [ ] Named DF-76 tests pass; no virgin auto-init from MCP
- [ ] Nine tools + `trace_version` unchanged; P15 Assert keepers green
- [ ] Locked verify PASS; board Notes → **P16-S01-02**
