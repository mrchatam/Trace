# P14-S02-02 — Scope review notes (Install / capability gates)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none

## Checklist evidence

| # | Result | Evidence |
|---|--------|----------|
| 1 | PASS | `internal/install` registry (`cursor` STABLE, `claude` CONDITIONAL); `TestInstallDetectListsCursorStable`; `TestInstallCursorUninstallIdempotent` |
| 2 | PASS | `TestInstallConditionalRefusesWithoutMarker` + `TestInstallConditionalWritesWithMarker` (`.claude/` \| `CLAUDE.md`) |
| 3 | PASS | `trace install detect` JSON list; no `install-all` / mass-write in product Go; CLI hard-fails unknown targets |
| 4 | PASS | Uninstall removes only `mcpServers.trace`; sibling `other` kept; second uninstall + missing path OK |
| 5 | PASS | All `TestInstallCursor*` (print/bin/write/reload tip/invalid JSON) PASS via thin `cmd/trace` → `internal/install` |
| 6 | PASS | `isBuiltinMCPSlug` exact `BuiltinMCPCapabilitySpecs()` only; no AllowAll/globs in product paths |
| 7 | PASS | Named decision tests: AUTO_ALLOWED durable; unknown → PENDING + Assert fail; ALLOWED persists reopen; DENIED blocks |
| 8 | PASS | Mig `013_capability_tool_decisions.sql` + store CRUD; CLI `capability decide`/`decisions` |
| 9 | PASS | MCP still nine + `trace_version` (`RegisteredToolNames` / `TestToolNamesRegistered`); no install/decide MCP tools |
| 10 | PASS | `impact_walk.go` mtime **before** S02 install files; named ImpactWalk + `TestPlantedImpactConflictsGateFPrelim` PASS; Expand untouched |
| 11 | PASS | `evals/capability` PASS under CGO0/CGO1/product |
| 12 | PASS | No daemon/HTTP/embeddings/Neo4j/full-rebuild; no YOLO/AllowAll in Trace product tree |
| 13 | PASS | Fresh verify (below); Gate C `metrics-*.json` still `dry_run: false` |
| 14 | PASS | P14-S02-00 Notes docs-only; P14-S02-01 Notes match live install + decisions |

## Independent verify (fresh)

```text
CGO_ENABLED=0 go test ./internal/install/... ./internal/domain/... ./internal/store/... ./evals/capability/... ./evals/honesty/... -count=1  → PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/retrieval/... ./evals/impact/... ./evals/capability/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1  → PASS
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1  → PASS
```

(Sandbox default `GOMODCACHE` hit proxy 403 on `segmentio/encoding`; home module cache + `GOPROXY=off` recovers — env, not product defect.)

## Findings

### blocker / high
None.

### medium
1. **`AssertToolAllowed` is library + tests only** — Not wired into MCP tool dispatch (correct for A4 CLI/rules-first / no new MCP). S03 should treat “audit exists” ≠ “every MCP call is gated” unless a future phase boards that wire-up.
2. **S01 residual `allowContainsOut` late-upgrade** — Still open from P14-S01-02; optional S03 spot-check only (non-blocking for S02).

### low
3. **CLI usage strings hardcode `cursor|claude`** — Registry is the SoT; adding OPT_IN later needs help/usage sync.
4. **No dedicated CLI tests for `install detect` / `capability decide|decisions`** — Library + Cursor CLI keepers cover bar; thin adapters only.

### nit
5. Claude CONDITIONAL writes `.claude/trace-mcp.json` (peer-shaped proof), not a full Claude Desktop schema — acceptable for ≥1 marker-gated target lock.

## Board / next
- Mark **P14-S02-02** done; next runnable **P14-S03-00**.
- No rewrite of P14-S02-00/01 prompt history; **no spawns**.
