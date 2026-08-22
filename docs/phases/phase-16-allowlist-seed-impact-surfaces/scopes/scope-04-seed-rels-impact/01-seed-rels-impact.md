# P16 / S04 / 01 — Seed rels + impact import (DF-70, 73)

## Metadata
- id: P16-S04-01
- todo_ids: [P16-S04-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Seed import accepts `discovery_mentions_task` and impact findings/alternatives. Board **status + Notes only**.

**Stop if `00-PLANNER.md` is not FINAL.**

## Locked defaults
| Item | Value |
|------|-------|
| Home | `cmd/trace/seed.go` + tests; domain link/impact APIs already exist |
| DF-70 | Same rel vocabulary as CLI DF-42 |
| DF-73 | New allowed top-level keys only; still reject garbage keys |
| Forbidden | MCP tools; new rel types beyond DF-42; YOLO |

## Named tests
- `TestSeedImportDiscoveryMentionsTask`
- `TestSeedImportImpactFindings` (and alternatives if 00 includes them)

## Locked verify (minimum)
```text
CGO_ENABLED=0 go test ./cmd/trace/... ./internal/domain/... -count=1 -run 'TestSeedImport|TestLinkDiscoveryMentionsTask'
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

## Exit criteria
- [ ] DF-70/73 named tests pass; unknown keys still fail-closed
- [ ] Board → **P16-S04-02**
