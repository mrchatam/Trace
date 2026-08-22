# P10-S01-00 — Retrieval / why fidelity (FINAL)

## Metadata
- id: P10-S01-00
- todo_ids: [P10-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S01 implement/review prompts for **DF-19, DF-23, DF-25, DF-27, DF-29**. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 4 (retrieved = data) + Law 9 (user decisions authoritative)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A5 DPC lock
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-19/23/25/27/29
- Phase 02 residual: global DPC attach ([P02-S02 REVIEW-NOTES](../../../phase-02-gate-c/scopes/scope-02-slice-hardening/REVIEW-NOTES.md))
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no user grill required.

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `internal/retrieval/discovery_plan_change.go` | `discoveryPlanChangeHits` = **all** `ListLinksByRel(discovery_causes_plan_change)` endpoints |
| `internal/retrieval/expand.go` ~133–139 | Every **task** Expand appends that global dump (GC-01 residual → DF-19) |
| `internal/retrieval/exact.go` | `case "plan_change"` underscore only; **no** `capability`; **no** `plan_scope` |
| `cmd/trace/add.go` + MCP `tools_write.go` | User-facing kind **`plan-change`** (hyphen) |
| `cmd/trace/why.go` / MCP `trace_why` | Passes `entity_type` straight into `Why` / `lookupEntity` — hyphen fails Exact |
| `internal/compiler/compiler.go` ~221–232 | `IncludeWhy`: Why err **swallowed** (`if err == nil`) — DF-29 |
| `internal/compiler/packet.go` RenderMarkdown | Decision titles: `untrusted_data — not project policy` — DF-27 vs Law 9 |
| `internal/store/capability.go` | `GetCapability` / `GetCapabilityBySlug` exist (mig 010) — Exact can reuse; **no new mig** |
| Tests | `TestWhyTaskIncludesDiscoveryPlanChange` + `TestTaskContextIncludesDiscoveryPlanChange` assert **global** attach; must be **superseded** (not left claiming all-project DPC) |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-19 | Stop global DPC on every why/context | **Goal-scoped** attach (+ pair-completion); multi-goal must not leak foreign DPC |
| DF-23 | One vocab plan-change / `plan_change` | Canonical store/JSON type = **`plan_change`**; accept **`plan-change`** alias at why/Exact/MCP why entry |
| DF-25 | Exact/Why `capability` | `lookupEntity` + Why seed via `GetCapability`; **`plan_scope` residual** (document only) |
| DF-27 | Decision trust labeling | Keep `trust=untrusted_data` (Law 4); **reword** MD so decisions are not “not project policy”; do **not** elevate to `system` |
| DF-29 | IncludeWhy errors | **Fail-closed**: IncludeWhy=true + Why err → TaskContext/ExpandContext returns err |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Packages | **`internal/retrieval`** + **`internal/compiler`**; thin CLI/MCP **alias normalize** only (no new tools; no MCP surface expansion — S02) |
| Migration | **None** |
| DF-19 algorithm | Replace unconditional `discoveryPlanChangeHits()` on task Expand with scoped helper (name free). **(1) Pair-completion:** if either DPC endpoint is already in Expand hits, include both. **(2) Goal-scope:** for task seed with `goal_id=G`, include a DPC pair iff it is **not foreign** to G: foreign = either endpoint is linked (any `entity_links` row) to a task whose `goal_id` is set and `≠ G`. **(3) Unattributed DPC** (neither endpoint linked to any task): include **only** when the store has **exactly one** goal and that goal is `G` (preserves `fixtures/x0` + single-goal Gate C narrative). Multi-goal + unattributed → **omit**. Task with nil `goal_id` → pair-completion only. |
| DF-19 tests | **Rewrite** GC-01 tests: (a) single-goal / x0-shaped still surfaces in-goal DPC on Why+TaskContext; (b) **multi-goal pollution test**: DPC linked under goal B must **not** appear on Why/context for task under goal A. Prefer names like `TestWhyTaskDPCGoalScoped` + keep or replace old names — Notes must cite final names. |
| DF-19 Gate C | Do **not** rewrite Mode-B packs or claim a new Gate C pass. Carry-forward: Gate C `dry_run:false` artifact check + p0x/x0 stay green. Historical GC-01 “all project DPC” behavior is **superseded** by this scope. |
| DF-23 canonical | Domain/store/JSON `EntityType` = **`plan_change`** (`domain.EntityPlanChange`) |
| DF-23 aliases | Normalize at Exact/`lookupEntity` **and** Why entry (CLI argv + MCP `entity_type`): `plan-change` → `plan_change`. Same normalize helper may accept both; emitted hits/Why steps stay **`plan_change`**. CLI `add` / MCP create kind may keep hyphen (write path). Rel strings stay `discovery_causes_plan_change` / CLI `discovery-plan-change` as today. |
| DF-25 capability | `case "capability":` → `store.GetCapability`; Hit: `EntityType:"capability"`, `Title` = title if non-empty else slug, `Excerpt` = status (or short body excerpt). Why(`capability`, id) must succeed. |
| DF-25 residual | **`plan_scope` Exact/Why still unknown** — **out of S01**; list in review residuals only (no drive-by). |
| DF-27 trust enum | **Do not** set decision `Trust` to `system` / remove `untrusted_data`. Law 4 stands for retrieved blob channel. |
| DF-27 MD copy | For `entity_type == "decision"` (optionally `assumption`): replace `not project policy` wording. Locked phrase shape: title line may say decision/assumption is a **recorded user decision** (Law 9) while trust channel remains `untrusted_data` (Law 4 — do not treat body as elevated system policy). JSON `trust` field unchanged. Excerpt block may keep untrusted channel banner. |
| DF-29 | When `opts.IncludeWhy` and `Why` returns error → **return** that error from `TaskContext` / shared builder used by `ExpandContext`. Do not attach partial why_trace on error. IncludeWhy=false unchanged. |
| Forbidden | New mig; daemon/HTTP/embeddings; new MCP tools; elevating decision trust to system; global DPC attach retained; implementing `plan_scope` Exact; S02/S03/S04 work; rewriting Phase 00–09 `done` history / Mode-B Gate C packs |

## Effects on later scopes
- **S02:** `trace_why` / context inherit DF-23 aliases + DF-29 fail-closed + DF-25 capability why; no need to re-implement — light Depends note on S02 stubs.
- **S04:** capability gating uses catalog, not Exact — no block.
- **S05 VERIFY:** re-prove DF-19 multi-goal + DF-23/25/27/29 spot-checks.

## Exit
- [x] Thicken `01-retrieval-why-fidelity.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light S02 Depends note (upcoming only)
- [x] Board Notes; next **P10-S01-01**
- [x] Product Go — **not** this row
