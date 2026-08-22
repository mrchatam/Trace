# P11-S04-00 — Capability upsert + hatch vs caps (FINAL)

## Metadata
- id: P11-S04-00
- todo_ids: [P11-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S04 implement/review prompts for **DF-41, DF-51**. Confirm live inventory; lock APIs/tests. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-41, DF-51
- [experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md](../../../../../experiments/_bughunt/post-p10/POST-P10-BUGHUNT.md) — DF-41
- [experiments/_post_p10/bughunt/adv_hatch_caps/](../../../../../experiments/_post_p10/bughunt/adv_hatch_caps/) — DF-51
- [experiments/_post_p10/bughunt/adv_cap_done/](../../../../../experiments/_post_p10/bughunt/adv_cap_done/) — DF-51
- Phase 10 S04 FINAL: [../../../phase-10-integrity-surfaces/scopes/scope-04-operator-capability-gates/00-PLANNER.md](../../../phase-10-integrity-surfaces/scopes/scope-04-operator-capability-gates/00-PLANNER.md) — DF-24/26 Gate G
- P11-S02 FINAL: hatch↔caps independence; WARNING may mention missing-caps here
- Live: `internal/domain/{capability,task_state}.go`; `internal/store/capability.go`; `cmd/trace/{capability,transition,help}.go`; `internal/mcp/{tools_parity,tools_write}.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no grill (phase homes + live inventory + S02 independence note do not conflict).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| Store `UpsertCapability` | `ON CONFLICT(id)` only; empty ID → new UUID then INSERT; **UNIQUE(slug)** fails on re-declare without `--id` |
| Domain `UpsertCapability` | Passes `in.ID` through; `GetCapabilityBySlug` exists but **not** used for empty-ID upsert |
| Tests | `TestUpsertCapabilityGetAndReject` asserts **duplicate slug with new id must fail** — DF-41 flips empty-ID path to update-by-slug; **keep** explicit different-id conflict fail |
| CLI/MCP declare | Thin adapters; `--id` / `id` optional — agents re-declare status by slug alone today → UNIQUE error (DF-41) |
| `TransitionTask` order | Missing-caps gate **before** →DONE hatch/PASS/operator (DF-24) |
| Hatch vs caps | `AllowDoneWithoutReview` does **not** set/bypass `AllowMissingCapabilities` — intentional independence (S02 forward note) |
| Dogfood DF-51 | `allow-done` alone + missing caps → exit **2** (no WARNING); `allow-done`+`allow-missing-caps` → OK with WARNING that mentions only PASS/`as_operator` |
| Migration | None needed — `UNIQUE(slug)` already present |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-41 | Re-declare same slug without `--id` fails UNIQUE | Empty-ID upsert resolves existing row **by slug** and updates; stable id; keep id-keyed conflict when explicit different id |
| DF-51 | Hatch footgun vs missing-caps | **Keep** independence (hatch ≠ missing-caps override); thicken WARNING + help/MCP so agents know `--allow-missing-caps` is still required |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF home | **DF-41, DF-51** only (DF-24/26 Gate G remain as shipped; do not collapse flags) |
| Packages | **`internal/domain`** (`UpsertCapability` slug resolve; hatch WARNING strings stay in thin adapters but domain policy for caps gate unchanged); thin **`cmd/trace`** (`capability`, `transition`, `help`); thin **`internal/mcp`** (`trace_capability` declare path inherits domain; `trace_transition` warning + schema/desc). **G19** — no business logic in adapters |
| Migration | **None** |
| DF-41 rule | When `CapabilityInput.ID` is empty: if `GetCapabilityBySlug(slug)` hits → reuse that **id** then store upsert-by-id. If miss → allocate new id (store behavior OK). Same slug must yield **same id** across re-declares; status/title/body/kind updates apply |
| DF-41 conflict | Nonempty `ID` that would create a **second** row for an existing slug (different id) → **still fail** (UNIQUE / assertable error). Do not silently retarget another id when caller supplied an explicit id |
| DF-41 surface | CLI `capability declare` + MCP `trace_capability` action=declare inherit domain — no adapter fork. Optional usage/help note that omit `--id` updates by slug |
| DF-51 independence | **Retain:** `AllowDoneWithoutReview` does **not** imply or bypass `AllowMissingCapabilities`. Check order unchanged: DF-24 caps → then →DONE hatch/PASS/operator. Gate G hatch retained for review/operator only |
| DF-51 WARNING | On **successful** transition with `AllowDoneWithoutReview`: CLI stderr WARNING + MCP `warning` must mention (1) escape / Review PASS / `--as-operator`/`as_operator` bypass **and** (2) missing capabilities still need `--allow-missing-caps` / `allow_missing_caps` (assertable phrases). Prefer extending existing WARNING strings — do not remove loud hatch UX (DF-26) |
| DF-51 docs | Help + transition usage and/or MCP `allow_done` schema/desc: state clearly that `--allow-done` / `allow_done` does **not** bypass the missing-capability gate |
| DF-51 non-goals | Do **not** make hatch auto-set `AllowMissingCapabilities`; do not remove separate override flag; do not weaken DF-24 fail-closed default |
| Tests (required) | (1) **`TestUpsertCapabilityBySlugUpdatesExisting`** (or equiv): declare slug → re-declare same slug without id, flip status → same id + new status. (2) Keep explicit different-id duplicate-slug **fail**. (3) **`TestAllowDoneDoesNotBypassMissingCaps`** (or equiv): missing required cap + `AllowDoneWithoutReview` only → reject. (4) Hatch + `AllowMissingCapabilities` → DONE OK. (5) Update **`TestAllowDoneWarnsOnStderr`** / **`TestTransitionAllowDoneEmitsWarning`** (or equiv) so WARNING asserts missing-caps / `allow-missing-caps` / `allow_missing_caps` wording. (6) Optional: help/MCP schema phrase assert |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` untouched; P10 DF-17/18/24/26; P11-S01 DF-40; P11-S02 DF-43/44; P11-S03 DF-47 |
| Forbidden | New mig; collapsing hatch into missing-caps override; multi-writer; daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–10 / P11-S01–S03 `done` history; S05+ product work |

## Effects on later scopes
- **S05** (retrieval/trust/DPC): no capability-upsert/hatch coupling — serial after S04 review only. Light Depends note on S05 stubs.
- **S08 VERIFY:** include DF-41 slug re-declare same-id update + DF-51 hatch≠missing-caps + WARNING mentions both review hatch and missing-caps override in evidence table.

## Exit
- [x] Thicken `01-capability-upsert-hatch.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends note (S05)
- [x] Board Notes; next **P11-S04-01**
- [x] Product Go — **not** this row
