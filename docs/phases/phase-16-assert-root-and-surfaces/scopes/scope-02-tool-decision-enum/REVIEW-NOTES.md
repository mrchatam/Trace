# P16-S02-02 REVIEW-NOTES — Tool-decision enum / DF-75 + DF-78

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none (`P16-S02-02a` / `02b` not inserted)  
**Next board:** P16-S03-00

Independent review (fresh subagent ≠ implementer). Claims from P16-S02-01 Notes re-verified against live code + locked verify cmds — not trusted alone. Sibling `00-PLANNER.md` is **FINAL**. Heal→PENDING vs DENIED, glob matching, and CLI Assert (DF-77 / S03) not re-opened.

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | Mig **014** CHECK rejects YOLO/garbage; four enums remain valid | PASS | `014_capability_tool_decision_enum.sql` `CHECK (decision IN ('AUTO_ALLOWED','PENDING','ALLOWED','DENIED'))`. `TestCapabilityToolDecisionCheckRejectsYOLO`: raw INSERT YOLO + empty fail; Upsert YOLO + whitespace fail; four enums upsert. Store `UpsertCapabilityToolDecision` Go `switch` rejects non-enum (no silent coerce). 001–013 not rewritten (`013_*` still unconstrained; 014 rebuilds). |
| 2 | Existing YOLO heals to **PENDING** (not AUTO_ALLOWED, not dropped) | PASS | Copy CASE: empty/NULL/unknown/`YOLO` → `PENDING`; valid four pass through. `TestCapabilityToolDecisionMigrateHealsYOLOToPending`: plant 013-shaped YOLO on `mcp:trace_add` → Open/014 → row exists, `PENDING`, not `AUTO_ALLOWED`. |
| 3 | Resolve does **not** upsert AUTO_ALLOWED over unknown/healed YOLO | PASS | `ResolveToolDecision` switch has `default` returning durable PENDING with **no** upsert; `isBuiltinMCPSlug` AUTO_ALLOWED only after `sql.ErrNoRows`. `TestResolveYOLOBuiltinDoesNotAutoAllow`: after heal, Resolve PENDING, store stays PENDING, Assert fails. |
| 4 | Unprefixed registered Name canonicalizes to `mcp:`+Name and gates MCP | PASS | `canonicalizeToolSlug`: exact `spec.Title` or `spec.Slug` from `BuiltinMCPCapabilitySpecs()` → `spec.Slug`. Used by Decide **and** Resolve. `TestDecideUnprefixedMCPNameCanonicalizes`: `trace_why` DENIED persists `mcp:trace_why`, no leftover `trace_why` row; Resolve/Assert on both forms DENIED. `TestMCPUnprefixedDecideGatesCallTool`: CallTool `trace_why` errors with DENIED; store not AUTO_ALLOWED. |
| 5 | Exact match only (no globs); custom slugs unchanged; `cli:` not rewritten | PASS | No `LIKE`/glob in SQL or Go (`==` / `IN` list). `TestCanonicalizeCustomAndCLISlugsUnchanged`: `tool:custom-allow`, `cli:add`, `trace_*`, `mcp:trace_*`, `trace_why_extra` persist as-given; `MCP:trace_why` / `Trace_Why` do not map. |
| 6 | Footgun fold: unprefixed DENIED wins over prefixed AUTO_ALLOWED | PASS | 014 fold table prio DENIED=4 > PENDING=3 > ALLOWED=2 > AUTO_ALLOWED=1; one canonical row. `TestMigrateUnprefixedDeniedFoldsOverAutoAllowed`: dual plant → one `mcp:trace_why` DENIED; unprefixed dropped. |
| 7 | Compat ceiling **14**; no 015+; P15 Assert helper still `mcp:`+Name | PASS | No `015_*.sql`. `evals/compat` allow 014 / forbid 015+ / `EmbedExpected==14`. `TestOpenCreatesDBAndMigratesIdempotent` versions include **14**. `assertMCPToolAllowed` still `AssertToolAllowed(ctx, "mcp:"+toolName)`. `TestToolNamesRegistered` still exactly nine (`registerTools` those nine only). |
| 8 | DF-77 **not** implemented; no new MCP tools; YOLO/AllowAll absent | PASS | `cmd/trace` has **zero** `AssertToolAllowed` / `cli:add` Assert. No `trace_install`/`trace_decide` tools. Product Go has no YOLO/AllowAll flags (hits only in tests + 014 comment + unrelated `similar projects/`). CLI `cmdCapabilityDecide` still thin → `DecideTool` (domain canonicalize; no CLI slug-policy fork). |
| 9 | Carry-forward honesty/E–H/ablation/compat **14**/p0x/x0 + product pkgs; S01 virgin keeper | PASS | All `01` locked cmds re-run this session (see below). `TestMCPVirginProjectDoesNotMkdir` in named `-run`. Compat ceiling **14**. |

## Locked verify (re-run 2026-08-17)

```text
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestMCPUnprefixed|TestMCPVirgin|TestOpenExisting|TestOpenCreates|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision|TestCapabilityToolDecision|TestResolveYOLO|TestDecideUnprefixed|TestCanonicalize|TestMigrateUnprefixed'
→ ok mcp, domain, store

CGO_ENABLED=0 honesty A/B/C+G, replan E, impact F, ablation → ok
CGO_ENABLED=1 Gate H, compat 14, p0x, x0 → ok
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
→ all product pkgs ok
```

## Findings by severity

| Severity | Finding |
|----------|---------|
| blocker | *(none)* |
| high | *(none)* |
| medium | *(none)* |
| low | Mig 014 SQL `IN` list of nine unprefixed Names duplicates `BuiltinMCPCapabilitySpecs()` (SQL cannot call Go). Live lists match `TestToolNamesRegistered` / specs Titles. Drift if a tenth tool is added later — S03+ must not add MCP tools; S06 keepers stay nine. Do not spawn. |
| nit | `capability decide` usage still examples `mcp:trace_why` only; unprefixed Names now work via domain. FINAL said optional usage note is OK, not required. |

## Find → refute (not reported as open)

| Proposed | Refute |
|----------|--------|
| SQL hardcoded names violate planner “no second nine-name table” | Mig SQL cannot invoke Go. Domain canonicalize **does** use `BuiltinMCPCapabilitySpecs()` only. Duplicate is migrate-fold only; currently identical; residual **low**. |
| Store Upsert does not canonicalize unprefixed slugs | FINAL home is domain Decide+Resolve. CLI `cmdCapabilityDecide` → `DecideTool`. MCP Assert uses `mcp:`+Name. Hunt path is `decide --slug`. Raw store upsert of unprefixed is out of DF-75/78. |
| Resolve `default` does not UPDATE YOLO→PENDING at runtime | Migrate already heals; CHECK blocks new garbage. Lock is no AUTO_ALLOWED upsert; default returns PENDING without falling through to builtin upsert. Named Resolve test proves store stays PENDING. |
| Dual-row fold test does not cover single unprefixed DENIED migrate | Fold test requires unprefixed Name rewrite onto the same slug as `mcp:`+Name; that CASE is the single-row rewrite. Runtime Decide covers new writes. |
| Handler seams (`callWhy`) ≠ SDK CallTool | Same functions registered in `registerTools`; P15/S01 quality bar used the same seams. |
| Heal→PENDING is fail-open vs DENIED | FINAL lock; user instructed not to reopen. Assert already fail-closes PENDING. |
| Globs might still match via prefix | Exact `==` / SQL `IN`; tests include `trace_*`, `mcp:trace_*`, `trace_why_extra`. |
| DF-77 already fixed because unprefixed decide gates MCP | That is DF-78. CLI still has zero Assert. Isolation MCP DENIED ≠ CLI DENIED remains S03. |

## Residuals for S03 / S06

- **S03 (DF-77):** `cli:` reserved and proven unchanged. Unprefixed MCP **Names** now persist as `mcp:`+Name — S03 must Assert `cli:<command>` only; must **not** treat `trace_why` / `add` as CLI slugs. Shared CHECK table; no second enum table. G19: CLI calls `domain.AssertToolAllowed`, no slug-policy fork. MCP helper stays `mcp:`+Name.
- **S06:** Import named CHECK/heal/canonicalize + unprefixed-decide tests (already on S06 SCOPE-TODOS) + keepers `TestMCPAssert*` / `TestToolNamesRegistered` / `TestOpenCreatesDBAndMigratesIdempotent` (v14) / `TestMCPVirginProjectDoesNotMkdir` + compat **14**. Do **not** claim DF-77 fixed.
- Do **not** fail later rows for R2 defer / R3–R4 wontfix / CGO0 analyzer compile.

## Board

- P16-S02-02 → `done`
- Next runnable → **P16-S03-00**
- No `02a`/`02b` spawn
