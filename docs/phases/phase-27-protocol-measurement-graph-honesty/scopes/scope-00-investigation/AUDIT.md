# S00 Investigation AUDIT — INT-07/08/10

**Date:** 2026-08-20  
**Row:** P27-S00-01  
**Git SHA:** unavailable (workspace has no `.git`)

## Executive summary

Phase 26 verify closed on **P25-2** (Parent orchestrator in installed rules) but left **P25-3 FAIL** on build-only G1: `discoveries=0 decisions=0` in the exported graph. This is **expected** per `RUBRIC.md` — the harness correctly measures that the agent did not record causal graph richness during a build-only session, and it is **not** a P25-C install regression (P25-1/2 PASS). The failure exposes a **measurement vs product** split: **INT-08/INT-10** (protocol/scoring) can detect thin graphs but cannot force agents to write discoveries; **INT-07** (`seed export --strict`) was intended to catch dishonest exports but today only enforces structural shape + `GateForExport` (parity with `GateForDone`), which **allows** export when deliberation is in STOP/ORIENT — exactly the build-only G1 state. **S01** should fix protocol v2: automate export, separate build-only vs directed-gap scoring arms, and harden `score.sh`. **S02** should extend `--strict` with graph-honesty rules (thin-graph counts, discovery→task linkage, BLOCKING→uncertainty) so export enforcement complements harness scoring.

## Phase 26 residual

### Evidence

| Artifact | Path | Key finding |
|----------|------|-------------|
| Score output | `experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-s05-score.txt` | L7: `decisions=0 discoveries=0`; L11–12: P25-1/2 PASS; L13: **P25-3 FAIL**; L16: VERDICT FAIL (1 check) |
| Export snippet | `experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-export-snippet.json` | L42–55: empty `decisions[]` and `discoveries[]`; L123–136: `deliberation_states` STOP `p19_saturated`; goals=1 tasks=5 (structurally valid) |
| Spot checks | `experiments/runs/2026-08-20-p26-s05-01-verify/evidence/spot-checks.txt` | L11–14: harness P25-3 FAIL documented as RUBRIC-expected on build-only G1; D3–D5 saturation/reset PASS |
| VERIFY notes | `docs/phases/phase-26-loop-implementation/scopes/scope-05-verify/VERIFY-NOTES.md` | L36: P25-3 FAIL non-blocking; L38: manual export step required before score |

### Repro

1. `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace`
2. `cd experiments/ab-p25-gap-pass-validation && ./prepare.sh G1`
3. Run build-only G1 session (`prompts/PROMPT-G1-BUILD.md`) — no gap-analysis follow-up
4. `trace seed export -o runs/G1/trace/graph.json` — **NOT** automated by `prepare.sh`
5. `./score.sh G1 --p25` → expect **P25-3 FAIL**: `discoveries=0 decisions=0`

### Interpretation

Build-only G1 is designed to test whether **installed rules alone** (gap pass in `trace-enforcement.mdc`) cause agents to record discoveries/decisions without a human “run gap analysis” prompt. Phase 26 evidence shows install wiring works (P25-1/2) but agent behavior did not produce graph richness (P25-3). `RUBRIC.md` L33–38 and verdict matrix L45–46 explicitly treat P25-3 FAIL on build-only as **“Install OK; behavior unchanged”** — it does **not** block Phase 26 closure (P25-2 was the closure signal). The residual is a **deferred measurement/product gap** (P25-D/E), not a regression in Phase 25 install bundle.

---

## File-by-file findings

### Harness (`experiments/ab-p25-gap-pass-validation/`)

| Path | Lines | Current behavior | Gap |
|------|-------|------------------|-----|
| `PROTOCOL.md` | L5, L54–56, L59–61 | One build session per arm; **no** follow-up gap prompt on G1; optional E02-Session-B noted separately | Export step omitted from prepare flow; Session-B not a first-class scored arm (INT-08, INT-10) |
| `RUBRIC.md` | L27–29, L33–38, L43–48 | P25-3: `discoveries≥1 OR decisions≥1`; expected failures table documents build-only thin graph | Same P25-3 rule applies to build-only and directed-gap; no separate pass row for Session-B (INT-10) |
| `score.sh` | L39–57, L82–96, L106–128 | `count_kind` handles legacy `entities[]` and seed v1 plural arrays; P25-3 checks file counts only | No `trace seed export --strict` invocation; no DB-vs-export drift check; G2 = file exists only (INT-08) |
| `prepare.sh` | L96–106 | G1: `init`, `index`, `seed import`, config, `install cursor` — **no export** | Operator must manually export before G2/P25-3; protocol drift risk (INT-08) |
| `prompts/PROMPT-G1-BUILD.md` | L45–47 | Instructs `seed export -o trace/graph.json --strict --enforce` | Agent may export thin graph successfully because product strict gate does not fail on zero discoveries (INT-07) |

### Product export / gate

| Path | Lines | Current behavior | Gap |
|------|-------|------------------|-----|
| `cmd/trace/seed.go` | L67–80, L83–129, L136–148, L178–195 | `--strict` runs `collectExportStructuralViolations` (version, goals, tasks presence) + `collectExportViolations` → `GateForExport`; `--enforce` fails closed | Structural check does not count discoveries/decisions; skips done/skipped/stale tasks in gate sweep (L121–124) (INT-07) |
| `internal/domain/seed_export.go` | L55–62, L347–349, L417–423, L588–597, L710–722 | Exports all entity arrays; omits `work_state`; includes `deliberation_states`, `uncertainties` when present | Faithfully exports empty `discoveries[]`/`decisions[]`; no export-time richness validation (INT-07) |
| `internal/loop/gate.go` | L60–61, L227–265 | `GateForExport` dispatches to same `evaluateDone` as `GateForDone` | `evaluateDone` allows export when phase is ORIENT or STOP (L259–261); build-only G1 verify task in STOP passes export gate despite thin graph (INT-07) |
| `internal/domain/seed_eval_rules_test.go` | L14–86 | Tests `eval_rules_path` pointer export/import only | **No** strict-export or thin-graph coverage; `cmd/trace/enforce_test.go` covers verification debt/regression but not graph richness (INT-07) |

### Supporting test evidence

| Path | Lines | Relevance |
|------|-------|-----------|
| `internal/loop/gate_test.go` | L362–382 | `TestEvaluateGate_Export_SameAsDone` — explicitly locks export≡done parity |
| `cmd/trace/enforce_test.go` | L298–335, L397–418 | Strict blocks on `verification_incomplete` / open regression; clean fixture passes `--strict --enforce` with no discovery count requirement |

---

## INT mapping table

| INT | Path | Current behavior | Gap | S01/S02 task seed |
|-----|------|------------------|-----|-------------------|
| INT-07 | `cmd/trace/seed.go` L67–80 | Structural strict: version=1, goals/tasks non-nil only | No min discovery/decision count; no orphan-discovery rule | **S02-T01:** Add `collectExportGraphHonestyViolations(doc)` called from `cmdSeedExport` strict path |
| INT-07 | `internal/loop/gate.go` L60–61, L259–261 | `GateForExport` = `evaluateDone`; STOP/ORIENT → allowed | Thin graph with `p19_saturated` STOP exports cleanly under `--strict --enforce` | **S02-T02:** Split export gate from done gate OR add export-only richness checks before `evaluateDone` short-circuit |
| INT-07 | `internal/domain/seed_export.go` L417–423 | Exports discoveries/decisions as stored (may be empty) | Export shape honest but semantically thin; no linkage validation | **S02-T03:** Validate `discovery_mentions_task` / `decision_affects_task` links when discoveries/decisions non-empty |
| INT-07 | `internal/domain/seed_eval_rules_test.go` | Eval-rules path tests only | No regression tests for thin-graph strict failure | **S02-T04:** Add fixture mirroring P26 export snippet (`discoveries=0 decisions=0`) expecting strict violation when policy enabled |
| INT-07 | `cmd/trace/enforce_test.go` | Tests verification debt, regression block | No thin-graph enforce test | **S02-T05:** Extend enforce tests for new honesty rules (warn vs enforce modes) |
| INT-08 | `prepare.sh` L96–106 | Automates import + install; no export | Manual export before score; operator error omits G2 | **S01-T01:** Add optional post-session `seed export` step or `score.sh` preflight that exports if missing/stale |
| INT-08 | `score.sh` L82–83, L118–128 | Reads committed `trace/graph.json`; counts entities | Does not run `trace seed export --strict`; no DB freshness check | **S01-T02:** Pre-score: `trace seed export -o … --strict` (warn) or `--strict --enforce` (fail G2) |
| INT-08 | `score.sh` L39–57 | Dual-shape `count_kind` | No git-sparsity / commit-order guard for FM-07 | **S01-T03:** Compare `exported_at_commit` vs workspace HEAD; flag post-hoc planning commits |
| INT-08 | `PROTOCOL.md` | Arm paths isolated; export not in §1 Prepare | Protocol omits export as required step | **S01-T04:** Document mandatory export + import-before-gate in Prepare section |
| INT-10 | `RUBRIC.md` L27–29, L43–48 | Single P25-3 threshold; verdict matrix expects build-only FAIL | Build-only FAIL is correct but conflated with “behavior unchanged” vs Session-B PASS | **S01-T05:** Add `P25-3a` (build-only baseline) vs `P25-3b` (Session-B richness) rows |
| INT-10 | `PROTOCOL.md` L59–61 | E02-Session-B optional, separate notes row | Not isolated in `score.sh`; no `--session-b` flag | **S01-T06:** Add `score.sh G1 --p25 --arm build\|directed` with separate rubric columns |
| INT-10 | `prompts/PROMPT-G1-BUILD.md` L47 | Build-only; follow installed rules for gap pass | Scoring does not distinguish “rules present” from “rules executed” | **S01-T07:** Session-B prompt + rubric link; P25-4 attestation formalized per arm |

---

## S01 task seeds (protocol v2 — INT-08 + INT-10)

| ID | File | Behavior change |
|----|------|-----------------|
| S01-T01 | `prepare.sh` | After G1 install block: run `trace seed export -o trace/graph.json` OR emit explicit “export required before score” guard |
| S01-T02 | `score.sh` | Preflight: if `graph.json` missing, fail G2 with repro hint; optional `--export-strict` runs `trace seed export --strict` before counts |
| S01-T03 | `score.sh` | Read `exported_at_commit` from export; warn/fail if behind `git rev-parse HEAD` (FM-07 git-sparsity) |
| S01-T04 | `PROTOCOL.md` | §1 Prepare: add step 4 export; §3 G1: separate build-only score from optional Session-B |
| S01-T05 | `RUBRIC.md` | Split P25-3 into build-only expectation (FAIL = documented) vs directed-gap (PASS required for P25-C validation) |
| S01-T06 | `score.sh` | Add `--arm build\|directed` or `--session-b` flag; isolate E02-Session-B from build-only P25-3 |
| S01-T07 | `PROTOCOL.md` + `RUBRIC.md` | Formalize P25-4 operator attestation per arm; document invalidation if gap prompt sent before score |

## S02 task seeds (graph honesty — INT-07)

| ID | File | Behavior change |
|----|------|-----------------|
| S02-T01 | `cmd/trace/seed.go` | New `collectExportGraphHonestyViolations(doc SeedDocument) []exportViolation` — options: min discoveries/decisions, orphan discoveries without `discovery_mentions_task` |
| S02-T02 | `internal/loop/gate.go` | Either separate `evaluateExport` from `evaluateDone`, or add pre-check that rejects STOP-with-zero-causal-entities when `--strict` |
| S02-T03 | `cmd/trace/seed.go` L120–124 | Revisit skipping done/skipped/stale tasks in export gate sweep — verify task in STOP may evade gate |
| S02-T04 | `internal/domain/seed_eval_rules_test.go` + `trace/eval-rules.json` | Wire eval-rules invariants for export honesty (pointer-only path exists; body evaluation TBD) |
| S02-T05 | `cmd/trace/enforce_test.go` | Thin-graph fixture: export with `discoveries=0 decisions=0` → strict stderr + enforce block when policy on |
| S02-T06 | `internal/loop/gate_test.go` | Update `TestEvaluateGate_Export_SameAsDone` if export gate diverges from done gate |

---

## Delta from Phase 24 CODEBASE-AUDIT

| P24 finding | Still valid? | Phase 26 change | Notes |
|-------------|--------------|-----------------|-------|
| FM-02 thin export / strict honesty (`CODEBASE-AUDIT.md` L18) | **Yes** | `--strict` ships but does not catch `discoveries=0 decisions=0` | P26 export snippet confirms; INT-07 still open |
| FM-03 P19 saturation sticky STOP (L19) | **Partial** | P25-B: D3–D5 saturation reset + `loop reset` (VERIFY-NOTES L25–27) | First empty apply no longer sticky-saturates; STOP still allows export (gate.go L259–261) |
| FM-06 arm isolation / G1≡B0 (L22) | **Yes** | E02 harness isolates `runs/B0` vs `runs/G1` workspaces | Protocol scoring still single P25-3 column; INT-08/10 deferred |
| FM-09 mode collapse (L25) | **Partial** | P25-C gap pass install (P25-1 PASS); P25-2 orchestrator rule | Build-only still thin graph (P25-3 FAIL); human directed-gap not automated |
| FM-10 discovery without task promotion (L26) | **Partial** | P25-A: `PromoteBlockingDiscovery`, `spawned_tasks[].discovery_id`, MCP reorder (Phase 26 S02) | Agent must still invoke promotion; 0 discoveries in P26 build-only export |
| Export/DB drift (P24 §3, L55–59) | **Yes** | `prepare.sh` imports seed; export manual | `score.sh` reads file only; no import-before-gate |
| `GateForExport` same as done (implied P24 FM-02) | **Yes** | `gate_test.go` L362–382 locks parity | Confirmed root cause for thin-graph `--strict --enforce` pass |
| Deliberation reset absence (P24 §3 L63) | **Closed** | `loop reset` API + D4 PASS | Stale for S02; not INT-07 scope |
| P24 open question: “Does `--strict --enforce` catch Session A thin graph?” (L99) | **Answered: No** | P26 evidence + code audit | Export allowed when STOP + zero discoveries |

---

## Risks / open decisions (for S01/S02 planners)

| Topic | Options (do not lock) | Owner scope |
|-------|----------------------|-------------|
| P25-3 threshold | (A) keep `discoveries≥1 OR decisions≥1` for directed-gap only; (B) add `plan_changes≥1`; (C) require `deliberation_states` hop_count>0 after gap pass | S01 |
| Build-only P25-3 semantics | (A) document FAIL as expected baseline; (B) downgrade to SKIP for build arm; (C) separate rubric ID P25-3a/3b | S01 |
| Export enforcement in harness | (A) warn-only `--strict` in score.sh; (B) `--strict --enforce` fails G2; (C) enforce only on Session-B arm | S01 + S02 |
| Thin-graph strict rule | (A) min count discoveries OR decisions; (B) require link when discoveries>0; (C) eval-rules.json driven | S02 |
| Export vs done gate split | (A) keep parity, add pre-structural honesty; (B) separate `evaluateExport`; (C) export gate only on active tasks | S02 |
| BLOCKING uncertainty rule | (A) fail if BLOCKING discovery exists without uncertainty row; (B) fail if uncertainty missing on export when findings reference BLOCKING | S02 |
| Done-task gate skip | (A) keep skip; (B) always evaluate verify task; (C) evaluate all tasks when `--strict` | S02 |

---

## Out of scope (this audit)

- Product implementation (S01/S02 rows own code changes)
- Re-scoring E01 historical runs
- Daemon / HTTP / MCP on P0-X path
- P27-S00-02 independent review (next board row)
- Changing Phase 26 closure verdict or rewriting `done` board history

---

## Preflight record

All 12 paths from `01-investigation.md` preflight verified **PASS** (2026-08-20).
