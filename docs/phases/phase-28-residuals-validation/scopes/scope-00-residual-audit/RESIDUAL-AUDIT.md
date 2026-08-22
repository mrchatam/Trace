# Residual audit — Phase 28

**Date:** 2026-08-20  
**Git SHA:** unavailable (workspace has no `.git`)  
**Auditor row:** P28-S00-01

## Executive summary

Phases 26–27 shipped the P25 product stack (INT-01/02/05/06/07/09) and harness protocol v2 (INT-08/10 labels, `--strict --enforce` in `score.sh`). Phase 27 VERIFY closed with **high** confidence on honesty/enforce and P25-1/2/3a/3b **labels**, leaving expected thin-graph and dogfood residuals. **Six residuals close in Phase 28** (R1 Session-B dogfood; R2/R3/R8 hook failClosed; R4 honesty dedupe; R5 P25-4 attestation; R7 regression matrix). **R6 (FM matrix gaps)** partially closes via S01 tests and remains a measurement concern through S05 VERIFY. **Disposition update (P28-S02-02):** **R1 closed** — Session-B directed dogfood P25-3b PASS. **Disposition update (P28-S03-02):** **R2/R3/R8 closed** — Option A hook failClosed + INT-11 drift tests. **Disposition update (P28-S04-02):** **R4/R5 closed** — store BLOCKING orphan loop removed; P25-4 env attestation wired. **Disposition update (P28-S05-02):** **R7 closed**; **R6 partial/deferred** (non-blocking); Phase 28 DR-HANDOFF CLOSED — **no successor**.

---

## Residual register disposition (R1–R8)

| ID | Status | Evidence | Target scope | Notes |
|----|--------|----------|--------------|-------|
| R1 | **closed** (P28-S02-02) | Session-B: `SESSION-B-SCORE.txt` `p25 arm: directed`; P25-3b PASS (disc=1 dec=1); G2 `--strict --enforce` PASS; snapshot thin disc=0/dec=0; RESULTS `E02-SB`; reviewer re-score directed PASS 2026-08-20 | **S02** | Dogfood complete without `prepare.sh` wipe / post-mutation `--arm build`. P25-4 remains SKIP (manual attestation in NOTES; productize in S04/R5). |
| R2 | **closed** (P28-S03-02) | Option A: `CursorLoopGateHookScript` empty `TRACE_TASK_ID` + `enforce=strict` → deny JSON + exit 2; off/warn/missing/invalid → allow; `ParentOrchestratorRule` matches script | **S03** | Independent review re-ran `go test ./internal/install/...`; live script body spot-check PASS. |
| R3 | **closed** (P28-S03-02) | FM-05: script-level failClosed under strict+no-task; Cursor hooks.json `failClosed: false` retained (`cursorhook.go` L129); HookDriftNote documents policy | **S03** | Unit: `TestCursorLoopGateFailClosedStrictNoTask`, `TestCursorLoopGateAllowNonStrictNoTask`; Cursor entry failClosed stays false. |
| R4 | **closed** (P28-S04-02) | Store `ListDiscoveries` BLOCKING orphan loop removed from `collectExportGraphHonestyViolations`; sole source `CollectSeedDocumentHonestyViolations` (`seed_export_honesty.go` L43–48); `SeedEntity` still id/title/body only; regression `TestSeedExportStrictBlockingOrphanDiscoverySingleHonestyViolation` — exactly one violation, no `BLOCKING discovery` msg | **S04** | Independent review (P28-S04-02) re-ran cmd+internal tests PASS 2026-08-20. |
| R5 | **closed** (P28-S04-02) | `score.sh` L218–224: arm-matched `P25_ATTEST_BUILD=Y` / `P25_ATTEST_DIRECTED=Y` → pass; unset → skip; wrong-arm ignored; `PROTOCOL.md` L74–87 + `RUBRIC.md` env one-liner | **S04** | RESULTS.md remains human narrative; no RESULTS parser (deferred). |
| R6 | **partial/deferred** (P28-S05-02) | INTERVENTION-MATRIX §3 FM-01..10 residual gaps post INT-01..11; FM-03/05/06 largely closed; FM-07 warn-only by design; FM-01/02/04/08/09/10 remain behavioral/measurement gaps | **S00, S01, S05** | Non-blocking at Phase 28 close; future human theme only. |
| R7 | **closed** (P28-S05-02) | `TEST-MATRIX.md` M-01..M-16 (S01) + S05-01 full regression VERIFY (unit/cmd/install/matrix/score) PASS; evidence `experiments/runs/2026-08-20-p28-s05-01-verify/evidence/` | **S01, S05** | Consolidated matrix + VERIFY complete; M-16 is label smoke not P25-3b dogfood. |
| R8 | **closed** (P28-S03-02) | `hook_drift_test.go`: `TestHookDriftHooksJSONShape` (`command`/`matcher`/`failClosed:false`); `TestHookDriftAllowDenyPermissionJSON`; `TestHookDriftNoteNonEmpty` retained | **S03** | INT-11 automation shipped; reviewer confirmed PASS. |

---

## INT-01..11 implementation matrix

| INT | Theme | Implementation | Test coverage | Residual |
|-----|-------|----------------|---------------|----------|
| INT-01 | BLOCKING discovery → task promotion | **shipped** — `promote.go` L10–13 `PromoteBlockingDiscovery`; `apply.go` L497–511 `spawned_tasks[].discovery_id` path; `next.go` `promotion_candidates[]` | **unit** — `promote_test.go` (5 cases); `next_test.go` promotion_candidates; **gap:** no `apply_test.go` E2E promote-via-apply | **open** — FM-10 path requires agent invoke; S01 should add apply integration test |
| INT-02 | P19 saturation / hop budget recalibration | **shipped** — `deliberation/types.go` L11–13 `SaturationEmptyThreshold = 2`; `NextConsecutiveEmptyApplies` L88–98; migration `028_deliberation_consecutive_empty.sql` | **unit** — `saturation_reset_test.go` L23–60 threshold; `deliberation/saturation_test.go`; `store/deliberation_test.go` column round-trip | **closed** — D3 PASS (P26 VERIFY) |
| INT-03 | Default gap-pass prompt in install | **shipped** — `gappass.go` L8–11 `GapPassPrompt`; wired in `enforcement.go` L65, L83, L97 | **unit** — `enforcement_test.go` L467–476 GapPassPrompt assertions; harness P25-1 PASS (P27) | **closed** for install text; behavior validated in S02 |
| INT-04 | Parent orchestrator + hook gate | **shipped** — `ParentOrchestratorRule` Option A text; `CursorLoopGateHookScript` deny under strict+empty task | **unit** — `TestCursorLoopGateFailClosedStrictNoTask`; `TestCursorLoopGateAllowNonStrictNoTask`; install text + hook install | **closed** (R2/R3, P28-S03-02) |
| INT-05 | Gap-pass deliberation reset | **shipped** — `domain/deliberation.go` `ResetDeliberationState`; CLI `loop reset`; `apply.go` post-reset threshold | **unit/integration** — `saturation_reset_test.go` L94–119 reset sequence; `domain/deliberation_test.go` L98; `cmd/trace/loop_test.go` reset CLI | **closed** — D4 PASS (P26) |
| INT-06 | MCP trace_add ordering / nudge | **shipped** — `mcp/server.go` L13 description reorder; `gappass.go` L11 promotion nudge | **unit** — `mcp_test.go` L411 description strings | **closed** for product surface; agent compliance is dogfood |
| INT-07 | seed export graph honesty | **shipped** — `seed_export_honesty.go` L17–58; `seed.go` document-only `collectExportGraphHonestyViolations`; `--strict --enforce` path | **unit/integration** — `seed_export_honesty_test.go`; `enforce_test.go` thin-graph + `TestSeedExportStrictBlockingOrphanDiscoverySingleHonestyViolation` | **closed (R4)** — single honesty source; P28-S04-02 |
| INT-08 | Protocol v2 / score.sh enforce | **shipped** — `score.sh` L109–133 preflight export + strict enforce G2; L135–157 FM-07 warn; `--arm` L22–28 | **dogfood-only** — P27 VERIFY harness run; **gap:** no script unit test | **closed** for shipped behavior; S01 optional score.sh smoke |
| INT-09 | Sticky STOP reason UX parity | **shipped** — `select.go` persisted `StopReason`; status JSON `stop_reason` field | **integration** — `saturation_reset_test.go` L157–207 gate/export/status alignment | **closed** — D5 PASS (P26) |
| INT-10 | Two-session rubric P25-3a/3b | **shipped + dogfood** — `score.sh` L196–218 arm-specific labels; `RUBRIC.md` split; Session-B directed on E02 G1 | **dogfood** — P27 build P25-3a FAIL (expected); P28 S02 P25-3b PASS (disc=1 dec=1) | **closed (R1)** — Session-B evidence + independent review |
| INT-11 | Hook drift verification | **shipped** — `HookDriftNote` + hooks.json `failClosed: false`; script Option A | **unit** — `hook_drift_test.go` shape + allow/deny JSON; `TestHookDriftNoteNonEmpty` | **closed** (R8, P28-S03-02) |

---

## Dogfood residuals

| Topic | Status | Evidence |
|-------|--------|----------|
| **R1 Session-B** | **Done** | Directed gap on existing G1 (no prepare); evidence under `scope-02-session-b-dogfood/` + RESULTS `E02-SB`; P28-S02-02 APPROVE. |
| **P25-3b validation** | **PASS** | Live + re-score: discoveries=1 decisions=1; `p25 arm: directed`; G2 honesty enforce PASS. |
| **Arm isolation** | **Held** | Snapshot before mutation thin; no `prepare.sh`; no post-mutation `--arm build` claimed as Session-A; P27 E02 row intact beside `E02-SB`. |
| **G1 workspace state** | Session-A snapshot + Session-B rich | `SESSION-A-GRAPH-SNAPSHOT.json` disc=0/dec=0; live `runs/G1/trace/graph.json` disc=1/dec=1 with honesty links. |

---

## Harness residuals

| Topic | Status | Evidence |
|-------|--------|----------|
| **R2/R3 hook allow path** | **Closed** (P28-S03-02) | Option A deny under strict+empty task; default-off allow preserved; Cursor `failClosed: false`. |
| **P25-4 attestation** | **Closed** (P28-S04-02) | Env attestation arm-matched in `score.sh`; PROTOCOL documents path; RESULTS.md still human narrative. |
| **FM-05** | **Closed product-side** (P28-S03-02) | Script deny under strict+no-task; harness E3 measurement may still diverge — S05 VERIFY. |
| **FM-07 git-sparsity** | Warn-only (by design) | `score.sh` L135–157; P27 FM-07 PASS when SHA matches. |

---

## Product residuals

| Topic | Status | Evidence |
|-------|--------|----------|
| **R4 BLOCKING dup** | **Closed** (P28-S04-02) | Store BLOCKING re-check removed; document honesty sole orphan source; single-message regression PASS. |
| **R5 P25-4** | **Closed** (P28-S04-02) | See harness table — env attestation shipped; RESULTS parser not required. |
| **R6 FM gaps** | Partial | FM-03 closed; FM-02 partial; FM-05 product closed (S03); FM-09/10 dogfood partial (S02); remaining measurement → S05. |

---

## FM matrix §3 cross-check (post INT-01..11)

| FM-ID | Addressed by | Still open? | Phase 28 action |
|-------|--------------|-------------|-----------------|
| FM-01 | INT-01, INT-07, INT-08 | **Yes** — seed import pins roster; promotion requires agent | S01 apply promotion test; S02 dogfood |
| FM-02 | INT-07 | **Partial** — enforce blocks thin export; agents may skip writes pre-export | S02 directed gap; honesty already shipped |
| FM-03 | INT-02, INT-05, INT-09 | **Mostly closed** — consecutive-empty threshold + reset + reason parity | S01 regression; monitor in S05 |
| FM-04 | INT-03, INT-04 | **Yes** — parent can delegate to workers without graph | S03 hook deny helps parent path only |
| FM-05 | INT-04, INT-11 | **Mostly closed** — Option A script deny; Cursor `failClosed` stays false by design | Closed in S03; S05 may re-score |
| FM-06 | INT-08, INT-10 | **Protocol closed** — arm isolation in score.sh | S02 validates directed arm |
| FM-07 | INT-03, INT-08 | **Warn-only** — post-hoc commits possible | Document in S05; no product block |
| FM-08 | INT-01, INT-06 | **Partial** — tool ordering shipped; agent must choose task path | S01 + S02 |
| FM-09 | INT-03, INT-04, INT-02, INT-05, INT-10 | **Partial** — mode collapse needs directed session proof | **S02** primary validation |
| FM-10 | INT-01, INT-06, INT-05 | **Partial** — promotion API shipped; 0 discoveries in build-only export (P26/P27) | S01 apply E2E; S02 dogfood |

---

## Test coverage matrix (R7)

| Area | P25 theme | Existing tests | Gap | S01 seed |
|------|-----------|----------------|-----|----------|
| Promotion | P25-A | `internal/domain/promote_test.go` (5); `internal/loop/next_test.go` promotion_candidates | No `apply_test.go` path: BLOCKING discovery → `spawned_tasks[].discovery_id` → task row | `internal/loop/apply_promotion_test.go` — apply with discovery_id after seed import |
| Saturation threshold | P25-B | `internal/loop/saturation_reset_test.go` L23–60; `internal/deliberation/saturation_test.go` | Covered | Extend matrix doc only |
| Discoveries-only non-saturate | P25-B | `saturation_reset_test.go` L62–91 | Covered | — |
| Reset + re-apply | P25-B | `saturation_reset_test.go` L94–135; `domain/deliberation_test.go` | Covered | — |
| STOP reason parity | P25-B | `saturation_reset_test.go` L157–207 | Covered | — |
| Thin export enforce | P25-E | `cmd/trace/enforce_test.go` L425–442; `seed_export_honesty_test.go` | Covered | — |
| Orphan discovery/decision | P25-E | `seed_export_honesty_test.go` L46–87 | Covered | Add BLOCKING dup regression after S04 fix |
| Install gap pass + orchestrator | P25-C | `internal/install/enforcement_test.go` L467–491 | Text only | Assert three bodies include both prompts |
| Hook script content | P25-C | `enforcement_test.go` L316–337 (loop gate in script) | No assert on allow-without-task vs deny-when-strict | Defer deny behavior to S03; S01 may golden script substring |
| MCP description | P25-A | `internal/mcp/mcp_test.go` L411 | Covered | — |
| score.sh P25-3a/3b labels | P25-D | None automated | Script arm branching | Optional `evals/p28-regression/score_arm_test.sh` |
| Hook drift | P25-C | `enforcement_test.go` L494–502 | No JSON schema golden | S03 `hook_drift_test.go` |

**Inventory (2026-08-20):** 8× `internal/loop/*_test.go`, 2× `internal/deliberation/*_test.go`, 4× `internal/install/*_test.go`, 20× `cmd/trace/*_test.go`; dogfood-only: full E02 G1 session, P25-3b richness, hook deny under strict parent.

---

## Closure plan → S01..S04

| Scope | Tasks (actionable) | Files |
|-------|-------------------|-------|
| **S01** | Add apply promotion integration test (BLOCKING → spawned_tasks.discovery_id); document matrix in `TEST-MATRIX.md`; optional score.sh arm smoke; extend install text regression | `internal/loop/apply_promotion_test.go` (new), `docs/.../scope-01-integration-tests/TEST-MATRIX.md`, optionally `evals/p28-regression/` |
| **S02** | Run Session-B on E02 G1 with `PROMPT-G1-DIRECTED-GAP.md`; score `./score.sh G1 --p25 --arm directed`; record P25-3b + P25-4 directed attestation in `SESSION-B-NOTES.md` | `experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-DIRECTED-GAP.md`, `score.sh`, `runs/G1/` |
| **S03** | Deny hook when enforce=strict and no TRACE_TASK_ID (option A per S03 planner); set or document `failClosed`; add hook drift golden test | `internal/install/enforcement.go` L106–108, `cursorhook.go` L124, `internal/install/hook_drift_test.go` (new) |
| **S04** | Dedupe BLOCKING orphan honesty messages; automate P25-4 attestation (env `P25_ATTEST_BUILD=Y` / `P25_ATTEST_DIRECTED=Y` or RESULTS.md parser) | `internal/domain/seed_export_honesty.go`, `cmd/trace/seed.go` L158–170, `experiments/.../score.sh` L218, `PROTOCOL.md` |

---

## Explicit defers (with rationale)

| Item | Owner | Rationale |
|------|-------|-----------|
| Autonomous discovery→task spawn (no human gate) | Product owner / future phase | INTERVENTION-MATRIX §4 — Trace law requires human-approved backlog expansion |
| Full Graphiti episode / temporal invalidation DB | Future phase spike | INT-05 minimal reset sufficient for P25-B; multi-phase migration |
| Daemon / HTTP / hosted MCP on P0-X | Out of scope | Phase 28 README + project laws |
| Rewriting Phase 24–27 `done` board history | N/A | Forward-only protocol |
| E01 historical re-scoring | N/A | P27 AUDIT out of scope |
| Changing HopBudget from 12 | S03+ only with evidence | Locked in P26 PLAN architecture decision |

---

## Risks / open decisions

| Topic | Options | Recommendation for planners |
|-------|---------|----------------------------|
| S01 ∥ S02 parallel | Serial default vs parallel after S00 | **Parallel OK** after audit locks protocol; do not run `prepare.sh G1` in S02 |
| Hook deny scope | (A) strict + no task → deny; (B) parent-orchestrator-detected only | **A** locked in S03-00 planner |
| P25-4 automation | (A) env flags; (B) parse RESULTS.md; (C) keep manual with template | S04 implementer picks A or B |
| R4 dedupe strategy | (A) honesty skips discoveries already checked as BLOCKING in store pass; (B) remove store pass; (C) merge messages | **A or B** in S04 — prefer single source in `seed_export_honesty.go` |
| Build P25-3a semantics | Keep FAIL as expected baseline | Already documented — S05 VERIFY should not treat as regression |
| TEST-MATRIX vs RESIDUAL-AUDIT overlap | RESIDUAL-AUDIT seeds S01; TEST-MATRIX is deliverable | S01-00 planner consumes this doc § Test coverage matrix |

---

## Preflight record

All 22 paths from `01-residual-audit.md` preflight **PASS** (2026-08-20, P28-S00-01).

## Live anchor spot-checks

| Anchor | Location | Verified |
|--------|----------|----------|
| Hook Option A (was allow-without-task) | `internal/install/enforcement.go` empty-task branch | ✓ strict → deny; off/warn/missing → allow (P28-S03-02) |
| P25-3a/3b + P25-4 skip | `experiments/.../score.sh` L196–218 | ✓ arm branch + skip |
| BLOCKING dup (R4) | `cmd/trace/seed.go` L158–170 + `seed_export_honesty.go` L43–48 | ✓ overlapping orphan checks |
| Thin-graph honesty | `seed_export_honesty.go` L21–25 min count | ✓ |
| Saturation threshold | `internal/deliberation/types.go` L11–13 | ✓ `SaturationEmptyThreshold = 2` |
| Hook drift doc | `internal/install/cursorhook.go` L17–28 | ✓ `HookDriftNote` |
| Promotion API | `internal/domain/promote.go` L10–13 | ✓ |
| Install INT-03/04 text | `gappass.go`, `enforcement.go` L83, L97 | ✓ |

**Path correction applied:** No `internal/loop/saturation.go` — saturation policy in `internal/deliberation/types.go` + tests in `internal/loop/saturation_reset_test.go` + migration `028_deliberation_consecutive_empty.sql`.
