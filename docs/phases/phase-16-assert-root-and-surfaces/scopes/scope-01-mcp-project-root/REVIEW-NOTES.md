# P16-S01-02 REVIEW-NOTES — MCP project root / DF-76

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none (`P16-S01-02a` / `02b` not inserted)  
**Next board:** P16-S02-00

Independent review (fresh subagent ≠ implementer). Claims from P16-S01-01 Notes re-verified against live code + locked verify cmds — not trusted alone. Sibling `00-PLANNER.md` is **FINAL**. Bind-to-defaultRoot vs OpenExisting not re-opened.

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | Virgin `project=` / missing `.trace/trace.db` → CallTool error; **no** auto-mkdir AUTO_ALLOWED DB | PASS | `openStore` calls `store.OpenExisting` only (`internal/mcp/project.go`). Production MCP has **zero** `store.Open` call sites (tests only). `TestMCPVirginProjectDoesNotMkdir` subtests: virgin `defaultRoot` (`callAdd` + `callVersion`); initialized bound root + `AddInput.Project` virgin override (errors; would succeed if silent bind-to-defaultRoot). `assertNoTraceDir` on both virgin dirs. |
| 2 | Empty `.trace/` without db also fail-closed (no `trace.db` created) | PASS | `TestOpenExistingEmptyTraceDir`: `errors.Is(ErrNotInitialized)`; no `trace.db`. MCP subcase `empty trace dir no db` in `TestMCPVirginProjectDoesNotMkdir` + `assertNoTraceDB`. |
| 3 | Bound-root DENIED still fail-closed | PASS | `TestMCPAssertDeniedBlocksCallTool`: planted `mcp:trace_why` DENIED → CallTool error contains `DENIED`. |
| 4 | Initialized other root stays isolated (no session-global DENY) | PASS | `TestMCPInitializedOtherRootIsolated`: A DENIED `mcp:trace_add`; B `store.Open` only; server bound to A; `callAdd` `Project=B` succeeds; unbound add on A still DENIED; B’s row is not DENIED. |
| 5 | CLI `store.Open` / `trace init` still creates `.trace/` | PASS | `cmd/trace/init.go` still `store.Open`. `TestOpenCreatesDBAndMigratesIdempotent` + `TestInitCreatesDB` (CGO1) PASS. CLI cmds still `store.Open` (add/why/index/…). |
| 6 | Sentinel `ErrNotInitialized`; MCP wraps `%w`; `errors.Is` in tests | PASS | `var ErrNotInitialized` in `internal/store/lock.go`. `OpenExisting` wraps `%w`. MCP `fmt.Errorf("mcp: %w", err)` for Locked/Unauthorized/NotInitialized. Store tests `errors.Is`; MCP `mustNotInitialized` uses `errors.Is` or wrap string. |
| 7 | No bind-to-defaultRoot; no new MCP tools; G19; no daemon; Assert slug unchanged | PASS | Virgin override errors instead of writing bound root. `TestToolNamesRegistered` still exactly nine; `registerTools` those nine only; no `trace_install`/`trace_decide`. Assert still `assertMCPToolAllowed` → `mcp:`+Name after successful open. Schema still through `013_*` (no 014). `TestImportBoundaryNoLibraryImportsMCP` in product bar. |
| 8 | Carry-forward honesty/E–H/ablation/compat **13**/p0x/x0 + product pkgs | PASS | All `01` locked cmds re-run this session (see below). Compat ceiling still **13**. |
| 9 | R2/R3/R4 / DF-75/77/78 not falsely claimed fixed | PASS | Implementer Notes claim DF-76 only; CGO0 `cmd/trace` residual attributed to R4. Live: `allowContainsOut` still in `impact_walk.go`; no mig 014 CHECK; no CLI `cli:` Assert; YOLO still persistable at schema 013. |

## Locked verify (re-run 2026-08-17)

```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestMCPVirgin|TestMCPInitialized|TestOpenExisting|TestOpenCreates|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision'
→ ok mcp, domain, store

CGO_ENABLED=0 go test ./cmd/trace/... -count=1 -run 'TestInitCreatesDB|TestInitFailClosedWhenStoreLocked'
→ FAIL build (tree-sitter CGO) — R4 deferred; same residual as implementer Notes

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInitCreatesDB|TestInitFailClosedWhenStoreLocked'
→ ok cmd/trace  (keepers proven)

CGO_ENABLED=0 honesty A/B/C+G, replan E, impact F, ablation → ok
CGO_ENABLED=1 Gate H, compat 13, p0x, x0 → ok
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
→ all product pkgs ok
```

## Findings by severity

| Severity | Finding |
|----------|---------|
| blocker | *(none)* |
| high | *(none)* |
| medium | *(none)* |
| low | Locked `01` verify lists `CGO_ENABLED=0` for `./cmd/trace/...` keepers; that package cannot compile without CGO (P15 R4 wontfix). Keepers are green under CGO1. Do not spawn. |
| nit | `OpenExisting` shares `openStore` internals after Stat (MkdirAll no-op on existing `.trace/`). FINAL allows this. Stat-then-open TOCTOU if the file vanishes is planner assumption 3 (not DF-76). |

## Find → refute (not reported as open)

| Proposed | Refute |
|----------|--------|
| `OpenExisting` still `MkdirAll` via shared `openStore` → auto-init hole | Miss path returns before `openStore`; named tests prove no `.trace/` / no `trace.db`. Exists path sharing internals is FINAL. |
| Tests call handler seams (`CallAdd`) not SDK CallTool | Same functions registered in `registerTools`; handler `return nil, nil, err` is CallTool error. P15 quality bar used the same seams. |
| Isolation does not pre-read B’s allowlist | `callAdd` on B succeeding while A stays DENIED **is** the HOLD. Post-add B row is not DENIED. |
| Empty `.trace/` MCP subcase skips `trace_version` | Same `openStore`; store named test + add write path cover fail-closed / no db. |
| CGO0 `cmd/trace` verify fail is S01 product defect | R4 deferred; CLI mkdir keepers PASS CGO1; product bar CGO1 green. |

## Residuals for S02 / S06

- S02 does **not** need S01 product (board still sequential). No S02 prompt thicken.
- S06 must import: `TestMCPVirginProjectDoesNotMkdir`, `TestMCPInitializedOtherRootIsolated`, `TestOpenExistingMissingReturnsErrNotInitialized`, `TestOpenExistingEmptyTraceDir`, plus P15 keepers `TestMCPAssertDeniedBlocksCallTool`, `TestMCPAssertBuiltinAutoAllowedSucceeds`, `TestToolNamesRegistered`. (S06 SCOPE-TODOS + S06-00 intent thickened this review.)
- Do **not** fail later rows for R2 defer / R3–R4 wontfix / CGO0 `cmd/trace` compile.
- Compat ceiling remains **13** until S02.

## Board

- P16-S01-02 → `done`
- Next runnable → **P16-S02-00**
- No `02a`/`02b` spawn
