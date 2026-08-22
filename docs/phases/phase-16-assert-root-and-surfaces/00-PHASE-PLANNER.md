# P16-00 — Plan post-P15 open findings (FINAL)

> **Forward addendum (after this row `done`, 2026-08-17):** Live DF-72 lock is **fix** (thin MCP `trace_impact` in S05). P14 A3 is superseded for that tool only. See [`DF-72-FORWARD.md`](DF-72-FORWARD.md). Disposition tables below are **historical** (this planner deferred DF-72) and must not be treated as the upcoming S05 SoT.

## Metadata
- id: P16-00
- todo_ids: [P16-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live repo + post-P15 dogfood/hunt/combo SoT, produce a **FINAL disposition** for every new open DF and scaffold implement+review board rows + prompts. **No product Go in this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md) — Laws 9, 15, 17, 19 (G19)
- [phase README](README.md)
- Findings SoT: [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md)
- [experiments/POST-P15-DOGFOOD.md](../../../experiments/POST-P15-DOGFOOD.md) — DF-68
- [experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md](../../../experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md) — DF-75…78
- [experiments/BATCH-D21-D23.md](../../../experiments/BATCH-D21-D23.md) — DF-70…74
- P14 A3 (no MCP impact): [../phase-14-peer-impact-install-gates/00-PHASE-PLANNER.md](../phase-14-peer-impact-install-gates/00-PHASE-PLANNER.md)
- P15 close: [../phase-15-p14-residual-plan/DR-HANDOFF.md](../phase-15-p14-residual-plan/DR-HANDOFF.md) — historical `no successor` (intact)
- Live: `internal/mcp/project.go` `openStore`→`store.Open` MkdirAll; `internal/domain/capability_decision.go` Resolve fall-through; `internal/store/schema/013_capability_tool_decisions.sql` (no CHECK); `cmd/trace/install.go` `ProjectRoot: cwd`; `cmd/trace/seed.go` link switch; `TestImportBoundaryMCPNoPlanImpactIndexTools`
- [docs/TODO.md](../../TODO.md)

## Session start
Agent → clarify (grill if DF-72 / CLI-first trades are material) → Plan → execute (planner). **Unattended:** human already cut all-open-DFs + DF-72 lock; defaults below are FINAL.

## Live confirmation (2026-08-17)

| ID | Still present? | Evidence |
|----|----------------|----------|
| DF-76 | **Yes** | `store.Open` MkdirAll `.trace/`; MCP `openStore` always Open; Assert is per opened DB |
| DF-75 | **Yes** | mig 013 no CHECK; `ResolveToolDecision` unknown persisted status falls through to builtin AUTO_ALLOWED upsert |
| DF-78 | **Yes** | `DecideTool` persists slug as given; Assert uses `mcp:`+Name only |
| DF-77 | **Yes** | CLI `add`/`why` never call `AssertToolAllowed` |
| DF-68 | **Yes** | `cmdInstall` receives `-C` root but `cmdInstallClaude`/`Detect`/`Uninstall` set `ProjectRoot: cwd` |
| DF-70 | **Yes** | `cmd/trace/seed.go` switch omits `RelDiscoveryMentionsTask` / hyphen alias |
| DF-71 | **Yes** | `internal/compiler` has no impact findings / `overall_class` |
| DF-72 | **Yes (non-goal)** | Nine MCP tools; `TestImportBoundaryMCPNoPlanImpactIndexTools` forbids `trace_impact`; P14 A3 |
| DF-73 | **Yes** | `seedDocument` has no findings/alternatives keys; unknown top-level rejected |
| DF-74 | **Yes** | `trace impact report` JSON uses PascalCase struct fields |
| P14 R2 / P15 R3/R4 | Hold | Not new DFs; do not board |

**Count boarded fix: 9** (DF-76, 75, 78, 77, 68, 70, 71, 73, 74). **Deferred:** DF-72. Unused gap: DF-69. Next free: **DF-79**.

## Disposition matrix (FINAL)

See [README.md](README.md). **Boarded:** S01→S05 fix + S06 VERIFY. **Not boarded:** DF-72 implement; R2/R3/R4; goals #2–#4.

## Locked defaults (FINAL — phase)

| Item | Value |
|------|-------|
| Phase | Thin post-P15 DF closeout — assert root + surfaces |
| History | Do not rewrite Phase 00–15 `done` prompts; P15 `no successor` intact as history |
| Product Go | **Forbidden** on P16-00 |
| MCP tools | Still **nine** + `trace_version`; **no** `trace_impact` / install / decide MCP; no daemon/HTTP |
| DF-72 | **defer / non-goal** — P14 A3 + import-boundary keeper hold |
| S01 | MCP must not auto-init; `OpenExisting` (or equivalent) — missing `.trace/` → CallTool error; per-store SoT HOLD (no session-global DENY) |
| S02 | mig **014** CHECK on decision enum; Resolve unknown → fail-closed (no AUTO_ALLOWED overwrite); canonicalize registered tool names to `mcp:` |
| S03 | Dual slug: MCP `mcp:<tool>`; CLI `cli:<command>` AUTO_ALLOWED independently; MCP DENIED ≠ CLI DENIED; do not Assert `capability decide` / `init` / `install` / migrate/backup/auth |
| S04 | Thread `cmdInstall(root)` into `InstallOpts.ProjectRoot`; Cursor STABLE home detect unchanged |
| S05 | Seed link + findings/alternatives; compiler impact fields; impact JSON snake_case; **no** MCP impact tool |
| Compat | Ceiling **14** after S02 (no 015+ from later scopes unless spawned) |
| DR-HANDOFF intent | After VERIFY: default **`no successor`** unless Notes promote |
| G19 | Adapters call library; no domain fork in CLI/MCP |

## Scope order (locked)
1. **S01 mcp-project-root** — DF-76
2. **S02 tool-decision-enum** — DF-75, DF-78
3. **S03 cli-mcp-allowlist-parity** — DF-77 (depends on S02 CHECK + prefix)
4. **S04 install-project-root** — DF-68
5. **S05 seed-impact-packet** — DF-70, 71, 73, 74 (DF-72 Notes-only)
6. **S06 VERIFY** — named S01–S05 + carry-forward + DR-HANDOFF

S04/S05 do not depend on S01–S03 product-wise; **board order stays sequential** (protocol). Scope planners must not start a lower-order pending row.

## Non-goals
- Product Go in this planner row
- `trace_impact` MCP / product daemon / new MCP tool menu
- Session-wide DENY across all project roots
- Boarding S05 research / `plan simulate` / D21+ / ranks 7+
- Fixing P14 R2 / P15 R3/R4
- Claiming Mode-B Gate C / new Gate C from this phase

## Planner work (this row)
1. [x] Deduplicate DF IDs vs SoT (DF-68, 70–78); confirm live gaps
2. [x] Disposition matrix FINAL (incl. DF-72 defer)
3. [x] Create scope folders + `00-PLANNER`/`01`/`02` stubs + SCOPE-TODOS + DR-HANDOFF stub
4. [x] Board P16-* rows; P16-00 done; next **P16-S01-00**
5. [x] AGENTS.md + PROJECT_DOCS_INDEX + DOGFOOD-FINDINGS scheduled → Phase 16
6. [x] Keep goals #2–#4 **out**; P15 historical handoff intact

## Exit criteria
- [x] Disposition matrix FINAL for DF-68, 70–78 + R2/R3/R4 hold
- [x] Board updated (S01–S06 rows spawned)
- [x] 00-PHASE-PLANNER marked **FINAL**
- [x] Notes point to next runnable row **P16-S01-00**
- [ ] Product Go — **not** this row

## Minimal todos
- [x] Inventory live DFs
- [x] Disposition matrix
- [x] Spawn S01–S06 stubs
- [x] Board + README + index + findings sync

## Next
**P16-S01-00** (scope planner for MCP project root / DF-76).
