# P16 / S02 / 01 — Allowlist CHECK + CLI parity (DF-75, 77, 78)

## Metadata
- id: P16-S02-01
- todo_ids: [P16-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Fail-closed tool-decision statuses (CHECK + Resolve) and honor the same allowlist on CLI MCP-equivalent verbs, including unprefixed decide slugs. Board **status + Notes only**.

**Stop if sibling `00-PLANNER.md` is not FINAL.**

## Locked defaults (inherited)
| Item | Value |
|------|-------|
| Home | domain Assert + store CHECK; CLI thin G19 calls |
| Mig | **014** CHECK on `capability_tool_decisions.decision` (00 may confirm number) |
| CLI | `add` / `why` (and other MCP-equivalent verbs 00 lists) Assert `mcp:<tool>` |
| Ungated | `capability decide` / `decisions` |
| DF-78 | Normalize or hint; unprefixed DENIED gates MCP |
| Forbidden | YOLO/AllowAll; install MCP; DF-68/70–74; daemon |

## Named tests
See 00-PLANNER. Must include YOLO reject, garbage ≠ AUTO_ALLOWED, CLI add/why after DENIED fails, unprefixed slug gates MCP.

## Locked verify (minimum)
```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/mcp/... ./cmd/trace/... -count=1 -run 'TestCapabilityDecision|TestResolveGarbage|TestCLIHonorsDenied|TestDecideUnprefixed|TestMCPAssert'
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Exit criteria
- [ ] DF-75/77/78 named tests pass; CHECK live; CLI DENIED fail-closed
- [ ] Operator `decide` still works when MCP capability slug is DENIED
- [ ] Board Notes → **P16-S02-02**
