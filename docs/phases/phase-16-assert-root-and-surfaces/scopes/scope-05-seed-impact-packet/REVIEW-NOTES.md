# P16-S05-02 REVIEW-NOTES — seed / impact packet + thin MCP (DF-70…74)

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none (`P16-S05-02a` / `02b` not inserted)  
**Next board:** P16-S06-00

Independent review (fresh subagent ≠ implementer). Claims from P16-S05-01 Notes re-verified against live code + locked verify cmds — not trusted alone. Sibling `00-PLANNER.md` is **FINAL**. Daemon, install/decide MCP, S04 `-C`, P14 walk/R2, and packet schema 0.3 were not re-opened.

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | DF-70 seed link rel + hyphen alias → `LinkDiscoveryMentionsTask` | PASS | `cmd/trace/seed.go` switch: `domain.RelDiscoveryMentionsTask` (`discovery_mentions_task`) **and** `"discovery-mentions-task"` → `svc.LinkDiscoveryMentionsTask`. `TestSeedImportDiscoveryMentionsTask` underscore + hyphen write store `discovery_mentions_task`; `unknown_rel` → `exitUsage` + `unknown link rel`. Keeper `TestLinkDiscoveryMentionsTaskCLI` PASS. |
| 2 | DF-71 why/context include findings / `overall_class` when present | PASS | Packet `Impact []domain.DecisionImpact` `json:"impact,omitempty"`; `SchemaVersion` still `"0.2"`. `attachTaskImpact` after trim calls `ImpactSummariesForTask` → `ImpactReport` (no forked rollup). MD `## Impact` after Items, before Why/Capabilities; shows `overall_class` + finding `impact_class`. Why CLI + MCP `toolWhy` merge `ImpactSummariesForWhySeed` (task neighborhood / decision iff findings). `TestContextIncludesImpactOverallClass` + `TestWhyIncludesImpactOverallClass` (task + decision subtests) PASS. MCP context uses `compiler.TaskContext` (inherits attach). |
| 3 | DF-73 seed `findings`/`alternatives` via existing impact domain | PASS | Allowlist adds `"findings"`, `"alternatives"` only. Order: entities → links → findings → alternatives → transitions. Calls `AddImpactFinding` / `AddDecisionAlternative` (input key `recommended`). Summary ints `findings` / `alternatives`. `TestSeedImportImpactFindings`: round-trip snake_case report + `is_recommended`; top-level `impact_findings` still `unknown top-level key`. No new rels/tables; `decision_impact_findings` / `decision_alternatives` already in Open keeper table list. |
| 4 | DF-74 snake_case impact report JSON | PASS | `store.DecisionImpactFinding` / `DecisionAlternative` json tags: `id`, `decision_id`, `impact_class`, `is_recommended`, …. CLI `cmdImpactReport` / finding+alt list encode those structs. `TestImpactReportJSONSnakeCase` requires `impact_class` / `is_recommended` / `overall_class`; forbids `"ImpactClass"` / `"IsRecommended"` / `"ID"`. MCP report inherits tags. Walk/blast types **not** retagged this scope. |
| 5 | DF-72 thin `trace_impact` MCP; Assert `mcp:trace_impact`; G19 | PASS | `internal/mcp/tools_impact.go` + `registerTools` `trace_impact` before `trace_version`. Annotations match `trace_capability` (`ReadOnlyHint=false`, `DestructiveHint=false`, `OpenWorldHint=false`). `assertMCPToolAllowed(ctx, st, "trace_impact")` after first `openStore`, before action switch (same close-then-reopen as `toolCapability`). Actions `finding`\|`alternative`\|`report`\|`walk`; ops `add`\|`list`\|`recommend` as locked. Adapter calls `AddImpactFinding` / `ImpactReport` / `retrieval.ImpactWalk` only. `TestMCPTraceImpactReport` AUTO_ALLOWED `mcp:trace_impact`; `TestMCPImpactDeniedBlocksCallTool` fail-closed. CLI still `failCLIDenied(..., "impact")` → `cli:impact` (no `cli:trace_impact`). |
| 6 | Boundary keeper **allows** `trace_impact`; still forbids plan/index/install/decide MCP | PASS | `TestImportBoundaryMCPNoPlanImpactIndexTools`: `trace_impact` must be registered; fatals on `trace_plan` / `trace_index` / `trace_install` / `trace_decide` / names containing `install` or `decide`. Live `registerTools` has none of the forbidden names. |
| 7 | Catalog ten including version; specs updated; Gate F keeper green | PASS | `RegisteredToolNames` length **10**: why, context, add, link, transition, review, tasks, capability, **`trace_impact`**, version. `TestToolNamesRegistered` + `TestBuiltinMCPCapabilitySpecs` (ten `mcp:` slugs including `mcp:trace_impact`). `BuiltinCLICapabilitySpecs` still 11 titles with `impact` only (no `cli:trace_impact`). `TestPlantedImpactConflictsGateFPrelim` PASS (ImpactReport rollup unforked). |
| 8 | Carry-forward + product pkgs PASS; no 015; S04 DashC keeper | PASS | All `01` locked cmds re-run this session (see below). Compat EmbedExpected **14**; `TestOpenCreatesDBAndMigratesIdempotent` migrations 1…**14**; no `015_*.sql`; `014_*.sql` Name `IN` list **untouched** (nine historical Names). `TestInstallClaudeDashCRefuseCitesProjectRoot` PASS. |

### Catalog after S05 (live)

| # | Name | Spec slug |
|--:|------|-----------|
| 1 | `trace_why` | `mcp:trace_why` |
| 2 | `trace_context` | `mcp:trace_context` |
| 3 | `trace_add` | `mcp:trace_add` |
| 4 | `trace_link` | `mcp:trace_link` |
| 5 | `trace_transition` | `mcp:trace_transition` |
| 6 | `trace_review` | `mcp:trace_review` |
| 7 | `trace_tasks` | `mcp:trace_tasks` |
| 8 | `trace_capability` | `mcp:trace_capability` |
| 9 | **`trace_impact`** | **`mcp:trace_impact`** |
| 10 | `trace_version` | `mcp:trace_version` |

**Not registered:** `trace_plan`, `trace_index`, `trace_install`, `trace_decide`.

## Locked verify (re-run 2026-08-17)

```text
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/compiler/... ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestContextIncludesImpactOverallClass|TestToolNamesRegistered|TestBuiltinMCPCapabilitySpecs|TestImportBoundaryMCPNoPlanImpactIndexTools|TestMCPTraceImpactReport|TestMCPImpactDeniedBlocksCallTool|TestMCPVirgin|TestOpenCreates'
→ ok compiler, mcp, domain, store
  PASS: TestContextIncludesImpactOverallClass
  PASS: TestToolNamesRegistered (10)
  PASS: TestBuiltinMCPCapabilitySpecs (10)
  PASS: TestImportBoundaryMCPNoPlanImpactIndexTools
  PASS: TestMCPTraceImpactReport
  PASS: TestMCPImpactDeniedBlocksCallTool
  PASS: TestMCPVirginProjectDoesNotMkdir
  PASS: TestOpenCreatesDBAndMigratesIdempotent (v1…14)

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestSeedImportDiscoveryMentionsTask|TestSeedImportImpactFindings|TestImpactReportJSONSnakeCase|TestWhyIncludesImpactOverallClass|TestCLIAddDeniedFailClosed|TestLinkDiscoveryMentionsTaskCLI|TestInstallClaudeDashCRefuseCitesProjectRoot'
→ ok cmd/trace
  PASS: TestSeedImportDiscoveryMentionsTask (underscore/hyphen/unknown_rel)
  PASS: TestSeedImportImpactFindings
  PASS: TestImpactReportJSONSnakeCase
  PASS: TestWhyIncludesImpactOverallClass (task/decision)
  PASS: TestCLIAddDeniedFailClosed
  PASS: TestLinkDiscoveryMentionsTaskCLI
  PASS: TestInstallClaudeDashCRefuseCitesProjectRoot

CGO_ENABLED=0 honesty A/B/C+G, replan E, impact F, ablation → ok
CGO_ENABLED=1 Gate H, compat 14, p0x, x0 → ok
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
→ all product pkgs ok
```

Named `-run` re-listed with `-v` this session (all PASS). Product suite used locked `GOMODCACHE`+`GOPROXY=off`.

## Findings by severity

| Severity | Finding |
|----------|---------|
| blocker | *(none)* |
| high | *(none)* |
| medium | *(none)* |
| low | `attachTaskImpact` returns silently if `ImpactSummariesForTask` errors (`if err != nil \|\| len == 0`). Matches `attachTaskCapabilities`. Why CLI/MCP **fail-closed** on the same helper. Named DF-71 tests cover the happy path. Do not spawn. |
| nit | `TestMCPImpactDeniedBlocksCallTool` plants DENIED then calls `action=report` (read). Fail-closed is proven; “no domain write” is vacuously true for report. Optional finding-add DENIED coverage was not a substitute lock. Seed summary always emits `findings`/`alternatives` ints (0 if absent) — planner asked for counts. |

## Find → refute (not reported as open)

| Proposed | Refute |
|----------|--------|
| Seed still omits mentions-task | Switch has both tokens → `LinkDiscoveryMentionsTask`. Named underscore+hyphen store `discovery_mentions_task`. |
| Packet schema bumped to 0.3 / alternatives on packet | `SchemaVersion = "0.2"`; `DecisionImpact` has findings only. Named test asserts `"0.2"`. |
| Forked OverallClass rollup | Compiler/why call `domain.ImpactReport` / `ImpactSummariesForTask`; Gate F keeper still plants ImpactReport. |
| Nested report still PascalCase | Store json tags snake_case; named report test forbids `ImpactClass` / `IsRecommended` / `"ID"`. |
| Boundary still forbids `trace_impact` / catalog 9 or 11 | Keeper **requires** `trace_impact`; `RegisteredToolNames` length 10, impact before version. |
| `cli:trace_impact` added / CLI specs changed | `BuiltinCLICapabilitySpecs` titles unchanged (`impact` only). CLI Assert slug remains `cli:impact`. |
| install/decide/plan/index MCP | Grep-clean in `registerTools`; boundary fatals those names. |
| mig 015 / edited 014 Name IN list | No `015_*.sql`; 014 still nine historical Names (`trace_why`…`trace_version`, no `trace_impact`). Live canonicalize uses `BuiltinMCPCapabilitySpecs` (new tool AUTO_ALLOWs). Compat **14**. |
| P14 walk / R2 rewritten | MCP walk parses seeds + `ImpactWalk`; library depth clamp 1..2 unchanged. Gate F green. Walk JSON keys match CLI wrapper. |
| Assert skipped | `toolImpact` Asserts after `openStore` before dispatch; DENIED report CallTool fails; AUTO_ALLOWED row after success. Pattern = `toolCapability`. |
| G19 fork / daemon / HTTP | Adapter encodes maps; domain/retrieval own logic. No daemon. |
| Why BFS changed | `cmd/trace/why.go` + `toolWhy` call `eng.Why` then merge impact. Retrieval Why signature unchanged. |
| Double-open skips second Assert | Same as `trace_capability`. First Assert gates the CallTool; DENIED test would fail if skipped. |
| S04 `-C` / `cli:install` / PID kill reopened | DashC keeper PASS; this scope did not thread `InstallOpts.ProjectRoot` into seed/impact/compiler. |
| New entity_links rels | Seed uses existing `RelDiscoveryMentionsTask`. No new rel constants. |

## Residuals for S06

- **S06:** Import named tests above + catalog **10** (`TestToolNamesRegistered` / `TestBuiltinMCPCapabilitySpecs`) + boundary keeper + Gate F + S04 DashC refuse keeper + `TestCLIAddDeniedFailClosed` / `TestLinkDiscoveryMentionsTaskCLI`. Claim DF-70…74 only when those pass. Still **no** install/decide/plan/index MCP. Runnable `-run` lines already on S06 `01-verify.md` (CLI keepers copied this review).
- `attachTaskImpact` error-swallow (capability-style) — do **not** fail VERIFY.
- 014 SQL nine-Name `IN` list remains historical — do **not** edit done 014.
- DF-67 / P14 R2 / P15 R3–R4 / DF-22/37 ops remain non-blocking.

## Board

- P16-S05-02 → `done`
- Next runnable → **P16-S06-00**
- No `02a`/`02b` spawn
