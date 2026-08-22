# P16-S03-02 REVIEW-NOTES — CLI vs MCP allowlist parity / DF-77

**Verdict:** APPROVE  
**Confidence:** high  
**Date:** 2026-08-17  
**Spawns:** none (`P16-S03-02a` / `02b` not inserted)  
**Next board:** P16-S04-00

Independent review (fresh subagent ≠ implementer). Claims from P16-S03-01 Notes re-verified against live code + locked verify cmds — not trusted alone. Sibling `00-PLANNER.md` is **FINAL**. Dual-slug isolation, write-set trim, unprefixed `add`→`cli:add`, and gating `capability decide` not re-opened.

## Checklist evidence

| # | Check | Result | Evidence |
|---|--------|--------|----------|
| 1 | Dual-slug: `mcp:trace_add` DENIED does **not** block `trace add` | PASS | `TestCLIAddSucceedsWhenMCPAddDenied`: decide `mcp:trace_add` DENIED → `trace add` exit 0 + goal persisted; durable `cli:add` AUTO_ALLOWED. Domain twin: `TestCapabilityDecisionAutoAllowBuiltinCLI` (`mcp:trace_add` DENIED leaves `cli:add` AUTO_ALLOWED). |
| 2 | `cli:add` DENIED fail-closes CLI add (no entity) | PASS | `TestCLIAddDeniedFailClosed`: decide `cli:add` DENIED → add non-zero, stderr contains `DENIED`, `ListGoals` empty. `failCLIDenied` → `assertCLICommand` → `AssertToolAllowed("cli:"+cmd)` after `store.Open` in `cmd/trace/add.go`. |
| 3 | Reverse isolation: `cli:add` DENIED does not block MCP `trace_add` | PASS | `TestCLIAddDeniedDoesNotBlockMCPAdd`: `cli:add` stays DENIED; CallTool `trace_add` succeeds; `mcp:trace_add` AUTO_ALLOWED. MCP helper unchanged: `assertMCPToolAllowed` still `"mcp:"+toolName`. |
| 4 | `why` gated; MCP why DENIED still allows CLI why | PASS | `TestCLIWhyDeniedFailClosed` (`cli:why` DENIED → non-zero + `DENIED` stderr). `TestCLIWhySucceedsWhenMCPWhyDenied` (`mcp:trace_why` DENIED → CLI why exit 0, stdout non-empty). `cmd/trace/why.go` Asserts after open. |
| 5 | Operator escape: `capability decide` works when `cli:add` DENIED; `init`/`install`/`help`/`version` never Assert | PASS | `TestUngatedCapabilityDecideWhenCLIAddDenied`: DENIED then decide ALLOWED exit 0. Grep: `failCLIDenied`/`assertCLICommand`/`AssertToolAllowed` absent from `capability.go`, `init.go`, `install.go`, `migrate.go`, `backup.go`, `restore.go`, `auth.go`, `help.go`. `cmdInstall` still drops `root`; `root.go` help/version/unknown never open store. |
| 6 | G19: CLI calls `domain.AssertToolAllowed("cli:"+cmd)`; `cli:reindex`→`cli:index`; unprefixed `add` is not a CLI slug | PASS | `cmd/trace/assert.go` concatenates only. Domain owns `canonicalizeToolSlug` fold `cli:reindex`→`cli:index` and `BuiltinCLICapabilitySpecs()`. `cmdIndex(..., command)` passes `index`/`reindex` from `root.go`. `TestCanonicalizeCLIReindexFoldsToIndex`; `TestCLIIndexAliasDenied` (both tokens fail). `TestUnprefixedAddDecideDoesNotGateCLI` in domain + CLI: slug `add` persists as `add`, `cli:add` still AUTO_ALLOWs. |
| 7 | Builtin CLI AUTO_ALLOWED; MCP specs still nine; helper `mcp:`+Name unchanged | PASS | `BuiltinCLICapabilitySpecs()` kind TOOL ×11 titles (`add`…`index`; not merged into MCP). Resolve AUTO_ALLOW reason `"builtin CLI command"`. `BuiltinMCPCapabilitySpecs()` still nine `mcp:` Names. `TestToolNamesRegistered` / `registerTools` those nine only. `internal/mcp/assert.go` production unchanged. |
| 8 | No shared slug / TTY skip / YOLO flags / mig 015 / new MCP tools / Assert on install | PASS | Independent slugs (`cli:` vs `mcp:`). No `isatty`/`IsTerminal`/`AllowAll`/`YOLO` in `cmd/trace`. No `015_*.sql`. Compat ceiling **14**. `install.go` ungated. |
| 9 | Carry-forward honesty/E–H/ablation/compat **14**/p0x/x0 + product pkgs; S01 virgin + S02 CHECK keepers | PASS | All `01` locked cmds re-run this session (see below). `TestMCPVirginProjectDoesNotMkdir` + `TestCanonicalizeCustomAndCLISlugsUnchanged` + `TestOpenCreatesDBAndMigratesIdempotent` (v**14**) in named `-run`. |

### Gated vs ungated (live vs FINAL)

**Gated** (Assert after open, before work): `add` `link` `transition` `review` (all six subs) `why` `context` `tasks` `seed` `impact` (`openImpact` + `cmdImpactWalk`) `plan` (`openPlanner`; `plan help`/`-h`/`--help`/unknown never open) `index`/`reindex`.

**Ungated:** help/version/empty; `init`; whole `capability` (`openDomain` has no Assert); migrate/backup/restore/auth; `install`; unknown usage.

## Locked verify (re-run 2026-08-17)

```text
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestCLIAddDeniedDoesNotBlockMCPAdd|TestUnprefixedAddDecideDoesNotGateCLI|TestCapabilityDecisionAutoAllowBuiltinCLI|TestCanonicalizeCLIReindexFoldsToIndex|TestCanonicalizeCustomAndCLISlugsUnchanged|TestCapabilityDecision|TestMCPAssert|TestMCPUnprefixed|TestToolNamesRegistered|TestMCPVirgin|TestOpenCreates'
→ ok mcp, domain, store

CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestCLIAddSucceedsWhenMCPAddDenied|TestCLIAddDeniedFailClosed|TestCLIWhySucceedsWhenMCPWhyDenied|TestCLIWhyDeniedFailClosed|TestUngatedCapabilityDecideWhenCLIAddDenied|TestCLIIndexAliasDenied|TestUnprefixedAddDecideDoesNotGateCLI'
→ ok cmd/trace

CGO_ENABLED=0 honesty A/B/C+G, replan E, impact F, ablation → ok
CGO_ENABLED=1 Gate H, compat 14, p0x, x0 → ok
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
→ all product pkgs ok
```

First CGO0 mcp attempt without `GOPROXY=off` hit sandbox 403 on `segmentio/encoding`; retry with locked `GOMODCACHE`+`GOPROXY=off` PASS (domain/store already ok on first run).

## Findings by severity

| Severity | Finding |
|----------|---------|
| blocker | *(none)* |
| high | *(none)* |
| medium | *(none)* |
| low | `BuiltinCLICapabilitySpecs()` 11 titles have no dedicated golden (unlike `TestBuiltinMCPCapabilitySpecs`). Live list matches FINAL. A typo would PENDING-fail that command (fail-closed) and break default CLI for that token. Isolation/DENIED/alias tests cover add/why/index. Do not spawn. |
| nit | `AssertToolAllowed` DENIED message uses the **input** slug, so `reindex` stderr is `tool cli:reindex is DENIED` after a `cli:index` DENIED (fail-closed still holds). `openDomain` in `impact.go` prefixes capability Open errors `impact:` — pre-existing helper location; S03 added `openImpact` without putting Assert on `openDomain`. |

## Find → refute (not reported as open)

| Proposed | Refute |
|----------|--------|
| Some gated cmds skip Assert | All store-opening paths in add/link/transition/review/why/context/tasks/seed/impact/plan/index call `failCLIDenied` after Open. `plan help` never opens. |
| `capability` / `install` / `init` Asserted | Grep empty for Assert helpers in those files. `openImpact` wraps Assert; `openDomain` (capability) does not. |
| Shared slug so MCP DENIED kills CLI | Named isolation tests + separate AUTO_ALLOW reasons. Canonicalize does not map `cli:`↔`mcp:`. |
| Unprefixed `add` folds to `cli:add` | FINAL rejected; `TestUnprefixedAddDecideDoesNotGateCLI` persists `add` and AUTO_ALLOWs `cli:add`. |
| Write-set trim left `why` ungated | `why.go` Asserts; `TestCLIWhyDeniedFailClosed` PASS. |
| `capability decide` gated (deadlock) | Whole catalog ungated; escape-hatch test PASS. |
| MCP helper / nine tools / mig 015 drifted | `assert.go` still `mcp:`+Name; `RegisteredToolNames` nine; no `015_*.sql`; compat EmbedExpected **14**. |
| TTY skip / YOLO flags | No matches in `cmd/trace`. |
| Missing per-command DENIED tests (plan/seed/…) | Locked named tests cover hunt + alias + escape; remaining gated cmds share the same helper after open. Residual coverage only — not a dual-slug defect. |
| `TestCLIAddDeniedFailClosed` does not require `exitFail` (2) | Implementation returns `exitFail`; test rejects `exitOK` + requires `DENIED` + no entity. Would still catch the hunt class. |
| Seed reads JSON before Assert | Assert is after `store.Open`, before Create* (FINAL: usage/parse with no store stays ungated). Not graph mutation. |
| `openDomain` `impact:` prefix is an S03 allowlist bug | Capability was already on `openDomain`; S03 correctly left it ungated. Out of DF-77. |

## Residuals for S04 / S06

- **S04 (DF-68):** `install` stays ungated — live-confirmed this review (`cmdInstall` drops `root`; no Assert). Do not add `cli:install`. `01` locked verify should use **CGO1** for `./cmd/trace/...` (R4: CGO0 tree-sitter). Thickened this row.
- **S06:** Import named isolation + `cli:` DENIED + alias + ungated decide + S01/S02 keepers (already on S06 SCOPE-TODOS). Claim DF-77 only when those pass. Runnable `-run` lines copied into S06 `01-verify.md`.
- Do **not** fail later rows for R2 defer / R3–R4 wontfix / CGO0 analyzer compile.

## Board

- P16-S03-02 → `done`
- Next runnable → **P16-S04-00**
- No `02a`/`02b` spawn
