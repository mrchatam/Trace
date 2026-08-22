# P16-S03-00 — CLI vs MCP allowlist parity (FINAL)

## Metadata
- id: P16-S03-00
- todo_ids: [P16-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live `cmd/trace` dispatch + S02 CHECK/`mcp:` canonicalize, lock **FINAL** dual-slug design for **DF-77**: fail-closed CLI allowlist **without** making MCP DENIED imply CLI DENIED. Thicken sibling `01`/`02`/SCOPE-TODOS. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md) — disposition
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — dual-slug lock
- Live: `cmd/trace/root.go` `run()` switch; **zero** `AssertToolAllowed` in `cmd/trace`; MCP `internal/mcp/assert.go` still `"mcp:"+toolName`; domain `ResolveToolDecision` AUTO_ALLOWs **MCP builtins only** (`isBuiltinMCPSlug`)
- Hunt: `experiments/_bughunt/post-p15/` — write-up `mcp-cli/` + results `mcp_cli_add.*`; `cap-decisions/` CLI why (`cap_cli_why.code` / `cap_cli_ungated.verdict`); [`POST-P15-BUGHUNT.md`](../../../../../experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-77
- Quality bar: [P16-S02-00](../scope-02-tool-decision-enum/00-PLANNER.md) FINAL
- S02 live: `TestCanonicalizeCustomAndCLISlugsUnchanged`; `assertMCPToolAllowed` unchanged; compat **14**
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (planner). Depends-on: **P16-S02-02 APPROVE** (board, already true). Phase locks below. **Unattended:** no architecture blockers; defaults below are FINAL. Do **not** make MCP DENIED imply CLI DENIED. **No SwitchMode** (orchestrator). **No product Go.**

## Depends (from S02 — live P16-S02-02 APPROVE)

S02 FINAL reserved the **`cli:` prefix** and did **not** implement CLI Assert. Live after S02:

- Canonicalize maps exact MCP **Name** / `mcp:`+Name → `mcp:`+Name only (`BuiltinMCPCapabilitySpecs`). `cli:add` stays `cli:add` (`TestCanonicalizeCustomAndCLISlugsUnchanged`).
- `decide --slug trace_why DENIED` now persists **`mcp:trace_why`** (DF-78). S03 must **not** treat unprefixed MCP Names as CLI commands; do not strip `cli:` into `mcp:`; do **not** canonicalize unprefixed `add`/`why` → `cli:add`/`cli:why`.
- Enum CHECK + YOLO→PENDING fail-closed is the shared store (`014_*`). S03 Asserts `cli:<command>` independently of `mcp:<tool>` — **no second decision table**, no mig **015**.
- CLI still has **zero** `AssertToolAllowed` (DF-77 still open). MCP `assertMCPToolAllowed` still `"mcp:"+toolName`.
- G19: slug policy stays in domain; S03 CLI adapter = `domain.AssertToolAllowed("cli:"+command)` — do not fork canonicalize in `cmd/trace`.
- Resolve AUTO_ALLOWs **MCP slugs only** today. Without `BuiltinCLICapabilitySpecs`, first `Assert("cli:add")` would PENDING fail-closed and break default `trace add`. S03 **must** AUTO_ALLOW exact builtin CLI slugs.

## Live inventory (confirmed 2026-08-17)

| Area | Present? | Gap |
|------|----------|-----|
| CLI dispatch | Yes — `cmd/trace/root.go` `run()` after `parseGlobal` | Tokens below; no Assert |
| CLI `AssertToolAllowed` | **Absent** (grep empty in `cmd/trace`) | **DF-77** |
| MCP Assert | Yes — `assertMCPToolAllowed` → `"mcp:"+toolName` at all nine tools | **Keep** P15 helper contract |
| Builtin AUTO_ALLOW | `isBuiltinMCPSlug` exact `mcp:`+Name only | `cli:add` is **not** builtin → would PENDING if Asserted today |
| Canonicalize | MCP Name / `mcp:`+Name only; `cli:add` unchanged | Need **only** `cli:reindex` → `cli:index` fold; **no** unprefixed command fold |
| Hunt `mcp-cli/` dir | **Absent** in tree (not committed / cleaned) | Results remain: `results/mcp_cli_add.code` = **0**; `mcp_cli_add.json` `ok:true` goal |
| Hunt `cap-decisions/` | Yes — `mcp:trace_why` DENIED then CLI why | `cap_cli_why.code`; `cap_cli_ungated.verdict` **PASS** (CLI why exit 0) |
| Shared table | mig **014** CHECK; ceiling **14** | Do not add 015 / second enum table |
| Nine MCP tools | Yes — `TestToolNamesRegistered` | Do not add MCP tools |

**Bug path DF-77 (live):** `decide --slug mcp:trace_add DENIED` (or `mcp:trace_why`) gates CallTool only. `trace add` / `trace why` never Assert → exit 0. Dual-slug **keeps** that isolation (MCP DENIED ≠ CLI DENIED) and **adds** independent `cli:` Assert so Shell lockdown is `decide --slug cli:add DENIED`.

### Live `run()` tokens (`cmd/trace/root.go`)

| Token(s) | Handler | Notes |
|----------|---------|-------|
| `help`, `-h`, `--help`; empty args | `printHelp` | No store |
| `version`, `--version`, `-version` | print `version` | No store |
| `init` | `cmdInit` | `store.Open` mkdir (keep) |
| `index`, `reindex` | `cmdIndex` | Alias same handler |
| `add` | `cmdAdd` | Hunt repro write |
| `link` | `cmdLink` | |
| `transition` | `cmdTransition` | |
| `review` | `cmdReview` | subs: create\|set\|get\|show\|list\|residual |
| `impact` | `cmdImpact` | |
| `capability` | `cmdCapability` | declare\|list\|require\|unrequire\|missing\|**decide**\|**decisions** |
| `plan` | `cmdPlan` | |
| `seed` | `cmdSeed` | |
| `tasks` | `cmdTasks` | |
| `why` | `cmdWhy` | Hunt repro read |
| `context` | `cmdContext` | |
| `migrate` | `cmdMigrate` | |
| `backup` | `cmdBackup` | |
| `restore` | `cmdRestore` | |
| `auth` | `cmdAuth` | |
| `install` | `cmdInstall` | S04 owns `-C`; **ungated** here |
| unknown | usage | No store |

`why` / `context` / `tasks` / `index` / some `plan` paths open `store` **without** `domain.New` today — S03 may `domain.New(st)` **only** to Assert (G19; no business-logic fork).

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Design | **Dual slug families.** MCP stays `mcp:<toolName>`. Gated CLI Asserts `cli:<command>` (e.g. `cli:add`, `cli:why`). Same `capability_tool_decisions` table (014 CHECK). **No** shared slug |
| Isolation | `decide --slug mcp:trace_add DENIED` does **not** block `trace add`. `decide --slug cli:add DENIED` does **not** block MCP `trace_add`. Operator Shell lockdown = `cli:` DENIED |
| Builtins | Exact `BuiltinCLICapabilitySpecs()` slugs AUTO_ALLOWED on first Resolve (reason `"builtin CLI command"`), same graduated model as MCP. Default `trace add` stays green |
| Specs | **New** `BuiltinCLICapabilitySpecs()` — kind `TOOL`, slug `cli:`+canonical command, title = command. **Do not** put CLI slugs in `BuiltinMCPCapabilitySpecs()` (nine MCP keepers stay nine) |
| Resolve | AUTO_ALLOW if `isBuiltinMCPSlug` **or** `isBuiltinCLISlug`. Unknown `cli:…` still PENDING fail-closed |
| Canonicalize | Keep S02 MCP fold. **Add only** `cli:reindex` → `cli:index`. Unprefixed `add`/`why`/`index`/`reindex` **unchanged** (not CLI slugs). `trace_add` stays MCP Name → `mcp:trace_add`. Do not strip `cli:` into `mcp:` |
| G19 | CLI helper concatenates `"cli:"+command` then `domain.AssertToolAllowed`. Alias fold + builtin list live in **domain**. No Assert logic fork; no TTY/isatty skip |
| Call site | Thin `assertCLICommand(ctx, svc, command)` in `cmd/trace` (MCP-shaped helper). Call **after** successful store open, **before** domain/retrieval/index work. Usage/flag parse with **no** store open stays ungated (no graph access) |
| Parent slug | Subcommands share the parent: `cli:review`, `cli:impact`, `cli:plan`, `cli:seed` — not `cli:review.create` |
| MCP Assert | **Unchanged** helper `mcp:`+Name; nine tools; no new MCP tools |
| Compat | Ceiling **14**; **no** mig 015+ |
| Forbidden | Shared slug so MCP DENIED kills CLI; TTY heuristics; YOLO/AllowAll flags; new MCP tools; daemon/HTTP; session-global DENY; rewriting Phase 00–15 `done` history |
| Carry-forward | honesty A/B/C+G; Gates E/F/H; ablation; compat **14**; p0x; x0; product pkgs `./cmd\|internal\|evals`. S01 virgin + S02 CHECK/canonicalize keepers stay green |

### Gated vs ungated (locked — do **not** trim to write-set)

Hunt repro is `add` **and** `why`. Trimming why/context/tasks would leave Shell-read paths with no `cli:` Assert (operator could not `decide --slug cli:why DENIED`). Full min list stands.

**Ungated** (never call Assert — operator escape + store lifecycle):

| Command | Why ungated |
|---------|-------------|
| `help` / `-h` / `--help` / empty args | No store |
| `version` / `--version` / `-version` | No store; not MCP `trace_version` |
| `init` | Must create `.trace/` so later Assert has a DB |
| `capability` **all** subs (`declare\|list\|require\|unrequire\|missing\|decide\|decisions`) | Escape hatch: `decide`/`decisions` must remain reachable if `cli:add` is DENIED. Do not explode sub-slugs; catalog ops stay operator |
| `migrate`, `backup`, `restore`, `auth` | Store lifecycle / token |
| `install` | Marker-gated installer; S04 `-C`; not an agent graph path |
| unknown command | Usage only |

**Gated** (Assert `cli:<canonical>` after open):

| CLI invocation | Canonical slug | Notes |
|----------------|----------------|-------|
| `trace add …` | `cli:add` | Hunt write |
| `trace link …` | `cli:link` | |
| `trace transition …` | `cli:transition` | |
| `trace review …` | `cli:review` | All subs |
| `trace why …` | `cli:why` | Hunt read — **gated** |
| `trace context …` | `cli:context` | |
| `trace tasks …` | `cli:tasks` | |
| `trace seed …` | `cli:seed` | |
| `trace impact …` | `cli:impact` | All subs |
| `trace plan …` | `cli:plan` | All store-opening subs; `plan help` with no open stays help |
| `trace index …` | `cli:index` | |
| `trace reindex …` | `cli:index` | Adapter passes `cli:reindex`; domain folds to `cli:index`. **Not** a separate builtin |

`BuiltinCLICapabilitySpecs()` titles (11): `add`, `link`, `transition`, `review`, `why`, `context`, `tasks`, `seed`, `impact`, `plan`, `index`.

### Dual-slug map (MCP twins vs CLI — independent)

| MCP slug (unchanged) | CLI slug | Isolation |
|----------------------|----------|-----------|
| `mcp:trace_add` | `cli:add` | DENY one ≠ DENY the other |
| `mcp:trace_link` | `cli:link` | |
| `mcp:trace_transition` | `cli:transition` | |
| `mcp:trace_review` | `cli:review` | |
| `mcp:trace_why` | `cli:why` | Hunt: MCP DENIED + CLI why still OK |
| `mcp:trace_context` | `cli:context` | |
| `mcp:trace_tasks` | `cli:tasks` | |
| `mcp:trace_capability` | *(CLI `capability` ungated)* | MCP DENIED must not block `trace capability decide` |
| `mcp:trace_version` | *(CLI `version` ungated)* | |
| — | `cli:seed` / `cli:impact` / `cli:plan` / `cli:index` | CLI-only agent paths |

## Named tests (required)

| Test | Package | Intent |
|------|---------|--------|
| `TestCLIAddSucceedsWhenMCPAddDenied` | `cmd/trace` | Hunt class: `DecideTool` `mcp:trace_add` DENIED → `trace add` exit 0 + entity; durable `cli:add` AUTO_ALLOWED. **Must not** fail because MCP is DENIED |
| `TestCLIAddDeniedFailClosed` | `cmd/trace` | `decide --slug cli:add DENIED` → `trace add` non-zero, stderr DENIED, **no** entity |
| `TestCLIWhySucceedsWhenMCPWhyDenied` | `cmd/trace` | Hunt `cap-decisions`: `mcp:trace_why` DENIED → `trace why` exit 0 |
| `TestCLIWhyDeniedFailClosed` | `cmd/trace` | `cli:why` DENIED → `trace why` fail-closed (proves why stayed gated) |
| `TestCLIAddDeniedDoesNotBlockMCPAdd` | `internal/mcp` | Reverse isolation: `cli:add` DENIED → CallTool `trace_add` still AUTO_ALLOWs / succeeds (fresh builtin) |
| `TestUngatedCapabilityDecideWhenCLIAddDenied` | `cmd/trace` | `cli:add` DENIED → `trace capability decide --slug cli:add --decision ALLOWED` still exit 0 (escape hatch) |
| `TestUnprefixedAddDecideDoesNotGateCLI` | `internal/domain` or `cmd/trace` | `DecideTool` slug `add` DENIED persists `add` (not `cli:add`); `AssertToolAllowed("cli:add")` still AUTO_ALLOWs |
| `TestCapabilityDecisionAutoAllowBuiltinCLI` | `internal/domain` | Fresh DB: Resolve/Assert `cli:add` AUTO_ALLOWED durable; `mcp:trace_add` DENIED does not change `cli:add` |
| `TestCanonicalizeCLIReindexFoldsToIndex` | `internal/domain` | `cli:reindex` Resolve/Decide/Assert ≡ `cli:index`; `decide cli:index DENIED` gates both |
| `TestCLIIndexAliasDenied` | `cmd/trace` | `cli:index` DENIED → both `index` and `reindex` fail-closed |
| `TestMCPAssertDeniedBlocksCallTool` | `internal/mcp` | **Keeper** — MCP path unchanged |
| `TestMCPAssertBuiltinAutoAllowedSucceeds` | `internal/mcp` | **Keeper** |
| `TestMCPUnprefixedDecideGatesCallTool` | `internal/mcp` | **S02 keeper** |
| `TestToolNamesRegistered` | `internal/mcp` | **Keeper** — still exactly nine |
| `TestCanonicalizeCustomAndCLISlugsUnchanged` | `internal/domain` | **Keeper** — `cli:add` / `tool:custom-allow` / globs still not rewritten to `mcp:` |
| `TestCapabilityDecisionAutoAllowBuiltinMCP` | `internal/domain` | **Keeper** |
| `TestMCPVirginProjectDoesNotMkdir` | `internal/mcp` | **S01 keeper** |
| `TestOpenCreatesDBAndMigratesIdempotent` | `internal/store` | **Keeper** — version **14** |

TDD: named isolation + `cli:` DENIED tests first (red: CLI never Asserts), then builtin CLI specs + helper call sites (green). Existing `cmd/trace` tests that `add` after `init` must stay green via AUTO_ALLOWED.

## Owns

| Item | Intent |
|------|--------|
| DF-77 | CLI gated commands Assert `cli:<command>`; MCP DENIED ≠ CLI DENIED; default CLI builtins AUTO_ALLOWED |
| Dual-slug table | Canonical list + ungated escape hatch |
| Alias | `reindex` → `cli:index` |

## Explicit deferrals

- DF-68 install `-C` (**S04**) — `install` stays ungated; do not Assert `cmdInstall`
- Unprefixed command names as CLI slugs (`add` → `cli:add`) — **rejected** (S02 residual: do not treat `add` / `trace_why` as CLI slugs)
- Per-subcommand slugs (`cli:review.create`)
- Gating `capability declare\|require` separately
- Trimming why/context/tasks to write-set
- R2 `allowContainsOut`; R3 graphify space; R4 CGO0 analyzers
- S05 / plan simulate / D21+; new MCP tools / install / decide MCP / daemon
- Session-global DENY across `project=` roots

## Assumptions (unattended)

1. **Do not trim:** hunt used `why`; Shell agents read `why`/`context`/`tasks`. Write-only would leave those with no `cli:` lever.
2. **Whole `capability` ungated:** so `decide` cannot deadlock behind `cli:capability` DENIED. Catalog declare/list/require are operator, not DF-77 hunt.
3. **No unprefixed CLI canonicalize:** `decide --slug add DENIED` is a custom slug; operator lockdown is `cli:add`. Matches S02 “do not treat `add` as a CLI slug.”
4. **`cli:reindex` fold only:** one alias exists in `run()`; denying `cli:index` must cover both tokens.
5. **AUTO_ALLOW CLI builtins:** CLI-first default green; operator DENIES after. PENDING-by-default would break `trace add`.
6. **Assert after open:** usage-without-store stays usage (exit 1); DENIED is operational fail (exit 2) with no mutation.
7. **CGO:** named `cmd/trace` tests run **CGO1** (R4: CGO0 `./cmd/trace/...` tree-sitter). Domain/MCP named tests stay CGO0.
8. S06 VERIFY imports the named isolation + `cli:` DENIED + keepers + compat **14**.

## Effects on later scopes

- **S04:** `install` remains ungated. Do not add Assert to `cmdInstall`. S03 does not change `-C` / ProjectRoot.
- **S06:** Import named tests in the table above (isolation + DENIED + alias + ungated decide + S01/S02 keepers). Claim DF-77 only when those pass.

## Planner work

1. [x] Inventory live `run()` tokens + zero CLI Assert + hunt `mcp_cli_add` / `cap_cli_why`
2. [x] Lock gated vs ungated + dual-slug map; **no trim**; `cli:reindex`→`cli:index` only
3. [x] Lock `BuiltinCLICapabilitySpecs` + Resolve AUTO_ALLOW + named tests
4. [x] Thicken `01-cli-mcp-allowlist-parity.md` + `02-scope-review.md` + SCOPE-TODOS
5. [x] Light S04 Depends (install ungated); S06 named-test pointer
6. [x] Board Notes → next **P16-S03-01**; mark this prompt **FINAL**

## Exit criteria
- [x] 00-PLANNER **FINAL** with dual-slug table
- [x] 01/02 runnable with locked defaults
- [x] No product Go
- [x] Next board row **P16-S03-01**

## Minimal todos
- [x] Inventory live CLI / hunt DF-77 / S02 keepers
- [x] FINAL locks + named tests
- [x] Thicken 01/02/SCOPE-TODOS + S04/S06 pointers
- [x] Board sync

## Next
**P16-S03-01** (implement DF-77 dual-slug CLI Assert).
