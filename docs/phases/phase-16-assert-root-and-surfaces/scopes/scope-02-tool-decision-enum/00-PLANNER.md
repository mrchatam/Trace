# P16-S02-00 — Tool-decision enum + slug prefix (FINAL)

## Metadata
- id: P16-S02-00
- todo_ids: [P16-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live mig 013 / `ResolveToolDecision` / `DecideTool`, lock **FINAL** defaults for **DF-75** (CHECK + YOLO fail-closed) and **DF-78** (`mcp:` canonicalize). Thicken sibling `01`/`02`/SCOPE-TODOS. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md) — disposition
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — FINAL
- Live: `internal/store/schema/013_capability_tool_decisions.sql` (no CHECK); `internal/domain/capability_decision.go` Resolve fall-through; `cmd/trace/capability.go` decide (thin); `internal/mcp/assert.go` `mcp:`+Name
- Pattern: [`012_import_provenance_enum.sql`](../../../../../internal/store/schema/012_import_provenance_enum.sql) (P13 DF-64 rebuild + CHECK)
- Hunt: `experiments/_bughunt/post-p15/{cap-decisions,mcp-yolo,mcp-footgun}/` + [`POST-P15-BUGHUNT.md`](../../../../../experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md) + `results/{cap_yolo,cap_nocheck,cap_prefix_report,mcp_yolo_after,mcp_footgun_*}`
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-75, DF-78
- Quality bar: [P16-S01-00](../scope-01-mcp-project-root/00-PLANNER.md); [P15-S01-00](../../../phase-15-p14-residual-plan/scopes/scope-01-mcp-assert-dispatch/00-PLANNER.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (planner). Depends-on: **P16-S01-02 APPROVE** (board, already true). Phase locks below. **Unattended:** no architecture blockers; defaults below are FINAL.

## Live inventory (confirmed 2026-08-17)

| Area | Present? | Gap |
|------|----------|-----|
| Mig **013** `capability_tool_decisions` | Yes — `id/slug/decision/reason/actor/timestamps`; UNIQUE(slug); index on decision | **No CHECK** on `decision` (**DF-75**). Raw SQLite accepts `YOLO` |
| Schema files | `001`…`013` only (embed `schema/*.sql`) | No `014_*`. Compat ceiling **13**, forbids `014+` |
| `NormalizeToolDecision` / `DecideTool` | Yes — human ALLOWED\|DENIED; unknown → `ErrValidation` | CLI/domain write path already rejects YOLO; **store Upsert + SQL** do not |
| `ResolveToolDecision` | Yes — persisted AUTO_ALLOWED\|PENDING\|ALLOWED\|DENIED returned | **No `default`:** unknown persisted status **falls through** to `isBuiltinMCPSlug` → upsert **AUTO_ALLOWED** (hunt `mcp-yolo/`) |
| Builtin slugs | `BuiltinMCPCapabilitySpecs()` → `mcp:` + nine Names | `isBuiltinMCPSlug` exact-matches **prefixed** spec.Slug only |
| `DecideTool` slug | Persists **as given** (`cmdCapabilityDecide` → domain) | `decide --slug trace_why` writes `trace_why`, not `mcp:trace_why` (**DF-78**) |
| MCP Assert | `assertMCPToolAllowed` → `AssertToolAllowed(ctx, "mcp:"+toolName)` | **Keep** P15 helper contract; do not change |
| Hunt `cap-decisions/` / `cap_nocheck.verdict` | FAIL | YOLO insert accepted |
| Hunt `mcp-yolo/` / `mcp_yolo_after.json` | `mcp:trace_add` became **AUTO_ALLOWED** after Resolve; add succeeded | Fail-open overwrite |
| Hunt `mcp-footgun/` | `trace_why` DENIED **and** `mcp:trace_why` AUTO_ALLOWED coexist | Unprefixed decide does not gate MCP |
| `cli:` Assert | Absent | **S03** (DF-77). S02 must not implement CLI gating |
| P13 DF-64 pattern | Mig **012** rebuild `imports` + CHECK; heal `''`/unknown → `EXTRACTED` on copy | Mirror rebuild+CHECK; **do not** heal privilege-like YOLO → AUTO_ALLOWED |

**Bug path DF-75 (live):** raw `INSERT … decision='YOLO'` on a builtin slug → Resolve switch misses → upsert AUTO_ALLOWED → Assert passes.

**Bug path DF-78 (live):** `decide --slug trace_why DENIED` persists a non-gating row; first MCP `trace_why` AUTO_ALLOWS `mcp:trace_why`.

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Home | `internal/store` mig **014** + `internal/domain` Resolve/Decide canonicalize. G19: CLI stays thin (no slug policy fork in `cmd/trace`). MCP Assert helper **unchanged** |
| Migration | **`014_capability_tool_decision_enum.sql`**. Do **not** rewrite 001–013. Rebuild `capability_tool_decisions` (012 pattern) with `CHECK (decision IN ('AUTO_ALLOWED','PENDING','ALLOWED','DENIED'))`. Restore `idx_capability_tool_decisions_decision` |
| Compat | Ceiling **13→14** (`evals/compat`: allow `014_*`; **forbid 015+**). `TestOpenCreatesDBAndMigratesIdempotent` must see version **14** |
| YOLO / garbage **write** (post-014) | CHECK rejects. Domain `NormalizeToolDecision` / `DecideTool` already error. Store `UpsertCapabilityToolDecision` must also reject non-enum (Go guard and/or CHECK) — no silent coerce |
| YOLO / garbage **migrate heal** | On 014 copy: empty / unknown / `YOLO` → **`PENDING`**. **Never** heal to `AUTO_ALLOWED`. **Never DROP** the row (drop + builtin Resolve would AUTO_ALLOW — fail-open) |
| Resolve unknown persisted | **Fail-closed.** `default` of the status switch: treat as PENDING (durable row, do **not** upsert AUTO_ALLOWED). Do **not** fall through to `isBuiltinMCPSlug` AUTO_ALLOWED |
| Slug canonicalize (Decide **and** Resolve) | Trim space. If slug is exact registered MCP **Name** (`BuiltinMCPCapabilitySpecs` Title / Name) **or** exact spec.Slug `mcp:`+Name → persist/lookup **`mcp:`+Name**. Exact match only — **no globs**, no `LIKE`, no prefix-of-prefix |
| Unprefixed example | `trace_why` → `mcp:trace_why`. Already-prefixed `mcp:trace_why` unchanged |
| Custom slugs | Unchanged (`tool:custom-allow`, unknown names, `trace_why_extra`) |
| `cli:` prefix | **Reserved, not implemented here.** `cli:add` must **not** canonicalize to `mcp:…`. S03 owns CLI Assert |
| Case | Exact (specs are `mcp:` + lowercase Names). `MCP:trace_why` / `Trace_Why` do **not** match |
| Migrate slug fold (DF-78 existing rows) | Unprefixed exact Name → `mcp:`+Name. If **both** `trace_why` and `mcp:trace_why` exist (hunt footgun): **fail-closed fold** into the canonical slug then drop the unprefixed row. Priority: **DENIED > PENDING (incl. healed garbage) > ALLOWED > AUTO_ALLOWED** |
| MCP Assert | **Unchanged:** `assertMCPToolAllowed` still `"mcp:"+toolName`. Do not change P15 helper contract or call sites |
| CLI decide | Prefer **zero** `cmd/trace` edits — domain canonicalize is enough. Optional usage-string note is OK, not required |
| Tests (named) | See table below |
| Verify cmds | See sibling `01` locked verify block |
| Forbidden | YOLO/AllowAll product flags; new MCP tools; daemon/HTTP; changing P15 Assert helper contract; DF-77 CLI Assert (S03); mig **015+**; rewriting Phase 00–15 `done` history; session-global DENY |
| Carry-forward | honesty A/B/C+G; Gates E/F/H; ablation; compat **14**; p0x; x0; product pkgs `./cmd\|internal\|evals`. S01 virgin/isolation keepers stay green |

### YOLO heal vs fail-closed (locked — do not reopen)

| Layer | Policy |
|-------|--------|
| Migrate 014 | Heal invalid → **PENDING** (Open must succeed on hunt DBs). Not AUTO_ALLOWED. Not DROP |
| Runtime Resolve | Unknown status → fail-closed PENDING path; **no** builtin AUTO_ALLOWED upsert |
| Runtime write | Reject (CHECK + domain + store). Human decide remains ALLOWED\|DENIED only |
| Why not heal→DENIED? | PENDING is “unknown/awaiting human” (same as non-builtin); DENIED would permanently lock builtins until ALLOWED. Assert already fail-closes PENDING |
| Why not fail the migration? | Hunt/dogfood stores with YOLO must still Open; CHECK copy would abort without a CASE heal |

## Named tests (required)

| Test | Package | Intent |
|------|---------|--------|
| `TestCapabilityToolDecisionCheckRejectsYOLO` | `internal/store` | After Open (schema 14), `INSERT` / `UpsertCapabilityToolDecision` with `YOLO` (and empty) **errors**. Valid four enums still upsert |
| `TestCapabilityToolDecisionMigrateHealsYOLOToPending` | `internal/store` | 013-shaped table (no CHECK) + builtin slug `YOLO` → apply 014 → row is **PENDING**, not AUTO_ALLOWED, row **not** dropped. Fixture: apply 001–013 then 014 SQL (or equivalent); do not rely on post-014 INSERT |
| `TestResolveYOLOBuiltinDoesNotAutoAllow` | `internal/domain` | Same hunt class: YOLO on `mcp:trace_add` then Open/014 → `ResolveToolDecision` returns PENDING (durable); **store row stays PENDING**; `AssertToolAllowed` fails. Must **not** upsert AUTO_ALLOWED |
| `TestDecideUnprefixedMCPNameCanonicalizes` | `internal/domain` | `DecideTool` slug `trace_why` DENIED persists **`mcp:trace_why`** (no leftover gating-inert `trace_why` row). `Resolve`/`Assert` on `mcp:trace_why` DENIED. `Resolve("trace_why")` same canonical slug |
| `TestCanonicalizeCustomAndCLISlugsUnchanged` | `internal/domain` | `tool:custom-allow` stays; `cli:add` stays (`cli:` reserved). `trace_*` / `mcp:trace_*` globs do **not** map to builtins |
| `TestMigrateUnprefixedDeniedFoldsOverAutoAllowed` | `internal/store` | Hunt footgun: rows `trace_why` DENIED **and** `mcp:trace_why` AUTO_ALLOWED → after 014, **one** canonical `mcp:trace_why` **DENIED** |
| `TestMCPUnprefixedDecideGatesCallTool` | `internal/mcp` | `DecideTool` `trace_why` DENIED then CallTool `trace_why` **errors** (DF-78). Must not AUTO_ALLOW `mcp:trace_why` |
| `TestMCPAssertDeniedBlocksCallTool` | `internal/mcp` | **Keeper** — prefixed `mcp:trace_why` DENIED still blocks |
| `TestMCPAssertBuiltinAutoAllowedSucceeds` | `internal/mcp` | **Keeper** |
| `TestToolNamesRegistered` | `internal/mcp` | **Keeper** — still exactly nine |
| `TestCapabilityDecisionAutoAllowBuiltinMCP` | `internal/domain` | **Keeper** — first resolve still AUTO_ALLOWED on `mcp:trace_why` |
| `TestCapabilityDecisionUnknownPendingFailClosed` | `internal/domain` | **Keeper** |
| `TestOpenCreatesDBAndMigratesIdempotent` | `internal/store` | **Keeper** — applied versions include **14** |
| `TestMCPVirginProjectDoesNotMkdir` | `internal/mcp` | **S01 keeper** — do not regress DF-76 |

TDD: named CHECK/heal/canonicalize tests first (red on live 013), then mig 014 + Resolve default + canonicalize (green).

## Owns

| Item | Intent |
|------|--------|
| DF-75 | Enum CHECK; YOLO cannot persist; existing YOLO heals to PENDING; Resolve never AUTO_ALLOWs over garbage |
| DF-78 | Unprefixed registered MCP Name gates `mcp:`+Name (Decide+Resolve+migrate fold) |
| Compat | Ceiling **14**; no 015+ |

## Explicit deferrals

- DF-77 CLI `cli:` Assert (**S03**) — prefix reserved only
- DF-76 already APPROVE (S01) — keepers only
- R2 `allowContainsOut`; R3 graphify space; R4 CGO0 analyzers
- S05 / plan simulate / D21+
- New MCP tools / install / decide MCP / daemon

## Assumptions (unattended)

1. **Heal→PENDING not DENIED:** fail-closed for Assert; operator can `decide ALLOWED` without treating YOLO as a hostile deny. Hunt close is “no AUTO_ALLOW overwrite.”
2. **Drop-on-migrate rejected:** deleting YOLO on a builtin reopens AUTO_ALLOWED.
3. **Canonicalize in domain** (Resolve+Decide) so CLI and any future caller share one rule; MCP helper stays `mcp:`+Name.
4. **Exact Name list** comes from `BuiltinMCPCapabilitySpecs()` — do not maintain a second hardcoded nine-name table (and do not add tools).
5. **Footgun fold is required** to close existing hunt DBs where unprefixed DENIED and prefixed AUTO_ALLOWED coexist; runtime canonicalize alone would still AUTO_ALLOW `mcp:trace_why`.
6. S03 must treat `cli:` as a **different family**; S02 must not rewrite `cli:*` to `mcp:`.
7. S06 VERIFY imports the named CHECK/heal/canonicalize + MCP unprefixed-decide tests + compat **14**.

## Effects on later scopes

- **S03:** `cli:` prefix reserved; do not canonicalize CLI slugs to `mcp:`. Dual-slug design unchanged. Needs S02 CHECK so garbage cannot fail-open either family.
- **S06:** Import named tests above + compat 14. Do not claim DF-77 fixed.

## Planner work

1. [x] Confirm live schema / Resolve fall-through / Decide as-given + 012 CHECK rebuild pattern
2. [x] Lock FINAL mig 014 shape, YOLO heal→PENDING vs Resolve fail-closed, named tests
3. [x] Lock slug canonicalize + migrate fold; `cli:` reserved
4. [x] Thicken `01-tool-decision-enum.md` + `02-scope-review.md` + SCOPE-TODOS
5. [x] Light S03 Depends (`cli:` reserved); S06 named-test pointer
6. [x] Board Notes → next **P16-S02-01**; mark this prompt **FINAL**

## Exit criteria
- [x] 00-PLANNER **FINAL**
- [x] 01/02 runnable with locked defaults
- [x] No product Go
- [x] Next board row **P16-S02-01**

## Minimal todos
- [x] Inventory live 013 / Resolve / Decide / hunt YOLO+footgun
- [x] FINAL locks + named tests
- [x] Thicken 01/02/SCOPE-TODOS + S03 Depends
- [x] Board sync
