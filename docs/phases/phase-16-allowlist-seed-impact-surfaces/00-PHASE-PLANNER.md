# P16-00 — Plan post-P15 surface hardening (FINAL)

## Metadata
- id: P16-00
- todo_ids: [P16-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live repo + post-P15 dogfood/hunt, produce **FINAL** locks, DF→scope map, non-goals, and runnable scope stubs for Phase 16. Address **all** still-open / deferred-product findings (not DF-76 alone). **No product Go in this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md)
- [phase README](README.md)
- Findings: [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md), [POST-P15-DOGFOOD.md](../../../experiments/POST-P15-DOGFOOD.md), [_bughunt/post-p15/POST-P15-BUGHUNT.md](../../../experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md), [BATCH-D21-D23.md](../../../experiments/BATCH-D21-D23.md)
- Pattern: [../phase-15-p14-residual-plan/](../phase-15-p14-residual-plan/), [../phase-14-peer-impact-install-gates/](../phase-14-peer-impact-install-gates/)
- Live: `internal/mcp/project.go`, `internal/store/open.go`, `internal/domain/capability_decision.go`, `cmd/trace/{install,seed,impact}.go`
- [docs/TODO.md](../../TODO.md)

## Session start
Agent → clarify (none material: user required all post-P15 DFs + suggested scope grouping) → Plan → execute (planner).

## Canonical open set (deduped 2026-08-17)

| ID | Status entering P16 | Action |
|----|---------------------|--------|
| DF-76 | open, high | fix S01 |
| DF-75, DF-77, DF-78 | open | fix S02 |
| DF-68 | open | fix S03 |
| DF-22, DF-37 | closed P11; residual ops | carry-forward S03 (tip keepers; no PID kill) |
| DF-70, DF-73 | open | fix S04 |
| DF-71, DF-72, DF-74 | open | fix S05 (DF-72 = thin MCP impact) |
| DF-67 | deferred residual | defer S06 (out of bar) |
| DF-36 | method | off-board |
| DF-66, R3, R4 | wontfix | do not reopen |
| R2 | P15 defer | remain deferred |

## Locked defaults (FINAL — phase)

| Item | Value |
|------|-------|
| Phase slug | `phase-16-allowlist-seed-impact-surfaces` |
| History | Do not rewrite Phase 00–15 `done` prompts; P15 `no successor` stays historical |
| Product Go | **Forbidden** on P16-00 |
| MCP critical path | No daemon/HTTP; thin MCP tools already exist — **S05 may add `trace_impact` only** |
| P14 A3 supersession | “No MCP impact” **superseded** for thin `trace_impact` (G19). Still **no** install/decide/plan/index MCP |
| DF-76 | MCP must not `store.Open` auto-mkdir a virgin `.trace/` via `project=`. Existing initialized roots OK. Per-store SoT stays. CLI Open auto-init unchanged |
| DF-75 | CHECK on `decision`; Resolve must **not** fall through garbage → AUTO_ALLOWED (fail-closed). Expect mig **014** |
| DF-77 | CLI MCP-equivalent verbs Assert the same `mcp:<tool>` slug. **Except** `capability decide` / `decisions` (operator ungated) |
| DF-78 | `decide` normalizes exact registered tool names (`trace_why` → `mcp:trace_why`) or rejects with hint; named test that unprefixed DENIED gates MCP |
| DF-68 | Pass `-C` / CLI `root` into `InstallOpts.ProjectRoot` for detect + claude (+ uninstall). Do not use process cwd when `-C` set. `cmdInstall` currently discards `root` |
| DF-70 | Seed switch accepts `discovery_mentions_task` + hyphen alias (DF-42 CLI already works) |
| DF-73 | Seed JSON v1 may import impact findings/alternatives (unknown keys stay rejected except the new allowed keys) |
| DF-71 | Task `context` / `why` include linked decision impact findings + `overall_class` (bounded; Law 6) |
| DF-74 | Impact report JSON snake_case (`id`, `impact_class`, `is_recommended`) like `tasks` |
| DF-72 | Thin MCP `trace_impact` mirroring CLI `finding`/`alternative`/`report`/`walk` subset; update `TestToolNamesRegistered` + **supersede** `TestImportBoundaryMCPNoPlanImpactIndexTools` to **allow** `trace_impact` only |
| DF-67 | **defer** — symbol honesty stays out of `index_honesty` bar |
| DF-22/37 | Keep print/write reload tip; no PID kill |
| YOLO / AllowAll | Forbidden |
| DR-HANDOFF intent | After VERIFY: default **`no successor`** unless Notes promote |
| Out | S05 supersession / plan simulate / D21+ / ranks 7+ / daemon / embeddings |

## Planner work (this row)
1. [x] Collect every still-open / deferred-product DF; dedupe
2. [x] FINAL disposition + DF→scope map (zero requested DFs unscheduled)
3. [x] Create scope folders + stub `00-PLANNER` / `01` / `02` / `SCOPE-TODOS` + `DR-HANDOFF`
4. [x] Board P16-* rows; mark this row **done**; next **P16-S01-00**
5. [x] Update AGENTS.md, PROJECT_DOCS_INDEX.md, DOGFOOD-FINDINGS schedule

## Exit criteria
- [x] Disposition matrix FINAL; DF→scope map complete
- [x] Board updated (S01–S06 rows spawned)
- [x] 00-PHASE-PLANNER marked **FINAL**
- [x] Notes point to next runnable row **P16-S01-00**
- [x] No product Go

## Minimal todos
- [x] Inventory live findings
- [x] Disposition + locks
- [x] Spawn S01–S06 stubs
- [x] Board + README + index + findings sync
