# P16-S01-00 — MCP project root / auto-init (FINAL)

## Metadata
- id: P16-S01-00
- todo_ids: [P16-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live `internal/mcp` + `store.Open`, lock **FINAL** defaults for DF-76: MCP CallTool must not auto-init a virgin dir into an AUTO_ALLOWED store. Thicken sibling `01`/`02`/SCOPE-TODOS. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md) — disposition
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — FINAL
- Live: `internal/mcp/project.go` `openStore` → `store.Open` (MkdirAll); `internal/store/open.go`
- Hunt: `experiments/_bughunt/post-p15/{mcp-deny,mcp-noinit,mcp-fresh}/` + [`POST-P15-BUGHUNT.md`](../../../../../experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-76
- Quality bar: [P15-S01-00](../../../phase-15-p14-residual-plan/scopes/scope-01-mcp-assert-dispatch/00-PLANNER.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute (planner). Phase lock: **no MCP auto-init**; per-store SoT HOLD (no session-global DENY). Grill only if OpenExisting vs bind-to-defaultRoot needs a tighter cut. **Unattended:** no architecture blockers; defaults below are FINAL (bind-to-defaultRoot **rejected** — see Assumptions).

## Live inventory (confirmed 2026-08-17)

| Area | Present? | Gap |
|------|----------|-----|
| MCP `openStore` | Yes — `internal/mcp/project.go` | Always `store.Open(abs)` after `resolveProject` |
| `store.Open` mkdir | Yes — `open.go` `os.MkdirAll(.trace)` then `sql.Open` creates `trace.db` | **This is the DF-76 create path.** CLI `trace init` and every CLI cmd still need it |
| `store.OpenExisting` | **Absent** (grep empty) | Add here; MCP must call it instead of `Open` |
| Sentinel for missing DB | **Absent** — store has `ErrLocked`, `ErrUnauthorized`, `ErrRestoreExists` only | Add `ErrNotInitialized` |
| Assert slug `mcp:`+Name | Yes — `assertMCPToolAllowed` at all nine tool entries | **Keep**; do not change |
| Per-store SoT | Yes — Assert reads the opened DB | Isolation **HOLD**; auto-init was the bypass |
| Hunt `mcp-noinit/` | Yes — now has `.trace/trace.db` + AUTO_ALLOWED `mcp:trace_add` + goal `bypass-noinit` | Live proof CallTool mkdir'd a virgin dir |
| Hunt `mcp-deny/` | Yes — nine builtins DENIED | Bound-root DENIED still holds **on that store** |
| Hunt `mcp-fresh/` | Yes — AUTO_ALLOWED after init | Positive control; not a bug |
| CLI `store.Open` | Yes — `cmd/trace/init.go` + add/why/… | **Must keep** mkdir (`TestInitCreatesDB`) |
| New MCP tools / daemon | Absent (good) | Do not add |

**Bug path (live):** `project=` (or virgin `defaultRoot`) → `resolveProject` Abs → `store.Open` MkdirAll → empty AUTO_ALLOWED DB → Assert never sees bound-root DENIED.

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Home | `store.OpenExisting` + MCP adapter `openStore`. G19: store owns existence; MCP does not fork mkdir policy. **Do not** remove CLI `store.Open` mkdir (`trace init` and other CLI cmds) |
| API | `func OpenExisting(projectRoot string) (*Store, error)` in `internal/store` (same Abs / no parent walk-up as `Open`) |
| Exists means | Regular file `<abs>/.trace/trace.db` **must exist** before any create. Missing file **or** `.trace/` dir without `trace.db` → fail-closed. Do not treat empty `.trace/` as initialized |
| Sentinel | `store.ErrNotInitialized` — wrap with `%w` so `errors.Is` works. Suggested text names missing `.trace/trace.db` |
| Create policy | `OpenExisting` **must not** create `.trace/` or `trace.db`. Do **not** implement as “Stat miss then `Open`” (that reintroduces mkdir). After a successful exists-check, sharing `Open` internals (lock, token, migrate, ensureProject) is OK. Creating `trace.lock` on an already-initialized store is OK |
| MCP | `Server.openStore` calls `store.OpenExisting` (not `Open`). Map `ErrNotInitialized` like `ErrLocked` / `ErrUnauthorized`: `fmt.Errorf("mcp: %w", err)` so CallTool returns **error** (not success content) |
| Fail-closed | Missing store → CallTool **error**; no AUTO_ALLOWED row; no goal/entity writes; virgin dir left without `.trace/` |
| Isolation | Per-store SoT **HOLD**. DENIED on initialized root A does **not** apply to initialized root B. S01 closes **auto-init**, not cross-root deny |
| `project=` vs defaultRoot | Both go through `openStore`. Virgin **either** path fails closed. **No** silent fallback to `defaultRoot` when `project=` is missing a store |
| Assert | Unchanged: `assertMCPToolAllowed` + slug `mcp:`+Name. Order stays openStore → Assert → tool body |
| Tests (named) | See table below |
| Verify cmds | See sibling `01` locked verify block |
| Forbidden | Session-global DENY; bind-to-defaultRoot; MCP-only Stat with no store API (rejected — duplicates the rule); new MCP tools; daemon/HTTP; YOLO/AllowAll; changing Assert slug `mcp:`+Name; new mig; editing CLI `Open` mkdir; rewriting Phase 00–15 `done` history |
| Carry-forward | honesty A/B/C+G; Gates E/F/H; ablation; compat **13** (S02 owns 14); p0x; x0; product pkgs `./cmd\|internal\|evals` |

## Named tests (required)

| Test | Package | Intent |
|------|---------|--------|
| `TestMCPVirginProjectDoesNotMkdir` | `internal/mcp` | Virgin `t.TempDir()` as `ProjectRoot` (and a second case: initialized defaultRoot + `project=` virgin override). CallTool (`trace_add` and `trace_version`) returns error (`errors.Is` `ErrNotInitialized` or error string wrapping it). **No** `<dir>/.trace/` created (and no `trace.db`). Empty `.trace/` without db also fails closed and does not create `trace.db` |
| `TestMCPAssertDeniedBlocksCallTool` | `internal/mcp` | **Keeper** — bound-root DENIED still blocks CallTool |
| `TestMCPAssertBuiltinAutoAllowedSucceeds` | `internal/mcp` | **Keeper** — initialized store still AUTO_ALLOWS builtins |
| `TestMCPInitializedOtherRootIsolated` | `internal/mcp` | Root A: `DecideTool` DENIED `mcp:trace_add`. Root B: `store.Open` (init, no deny). Server `ProjectRoot=A`. `callAdd` with `Project=B` **succeeds** (B’s own allowlist). `callAdd` without override (A) still **DENIED**. B must not inherit A’s DENIED rows |
| `TestOpenExistingMissingReturnsErrNotInitialized` | `internal/store` | Virgin dir: `OpenExisting` errors `ErrNotInitialized`; no `.trace/` mkdir |
| `TestOpenExistingEmptyTraceDir` | `internal/store` | Pre-create `.trace/` only (no db): `OpenExisting` errors; still no `trace.db` |
| `TestOpenCreatesDBAndMigratesIdempotent` | `internal/store` | **Keeper** — `Open` still mkdir+migrate |
| `TestInitCreatesDB` | `cmd/trace` | **Keeper** — CLI init still creates `.trace/trace.db` |
| `TestToolNamesRegistered` | `internal/mcp` | **Keeper** — still exactly nine tools |

## Owns
| Item | Intent |
|------|--------|
| DF-76 | MCP CallTool cannot mint a fresh AUTO_ALLOWED store by pointing `project=` (or a virgin `-C`/cwd) at a dir with no `trace.db` |
| CLI init | Unchanged mkdir via `store.Open` |

## Explicit deferrals
- Session-global DENY across all `project=` roots (phase HOLD)
- DF-75 / DF-78 (S02 CHECK + `mcp:` prefix)
- DF-77 CLI allowlist (S03)
- R2 `allowContainsOut`; R3 graphify space; R4 CGO0 analyzers
- S05 / plan simulate / D21+

## Assumptions (unattended)
1. **OpenExisting vs bind-to-defaultRoot:** missing store is **fail-closed CallTool error**, not a silent bind to `defaultRoot`. Fallback would hide a mis-pointed `project=` and mix allowlists. Phase lock already chose fail-closed.
2. **MCP-only exists-check** without a store API is rejected: one sentinel + store tests; `openStore` stays a one-line switch to `OpenExisting`.
3. Stat-then-`Open` on the **exists** path is OK (TOCTOU if the file vanishes between Stat and open is a low residual — not DF-76). Stat-miss **must not** call `Open`.
4. `trace_version` has no `project` arg — it uses `openStore("")` → defaultRoot/cwd; virgin defaultRoot fails closed (same API).
5. S06 VERIFY imports the named virgin + isolation tests + P15 Assert keepers; S02 does not need S01 product (board order still sequential).

## Planner work
1. [x] Confirm live `openStore` / `Open` mkdir path + hunt `mcp-noinit` auto-created DB
2. [x] Lock FINAL API `store.OpenExisting` + `ErrNotInitialized`; reject bind-to-defaultRoot
3. [x] Named tests (virgin no mkdir; DENIED keeper; isolation)
4. [x] Thicken `01-mcp-project-root.md` + `02-scope-review.md` + SCOPE-TODOS
5. [x] Board Notes → next **P16-S01-01**; mark this prompt **FINAL**

## Exit criteria
- [x] 00-PLANNER **FINAL**
- [x] 01/02 runnable with locked defaults
- [x] No product Go
- [x] Next board row **P16-S01-01**

## Minimal todos
- [x] Inventory MCP openStore + store.Open mkdir
- [x] FINAL locks + named tests
- [x] Thicken 01/02/SCOPE-TODOS
- [x] Board sync
