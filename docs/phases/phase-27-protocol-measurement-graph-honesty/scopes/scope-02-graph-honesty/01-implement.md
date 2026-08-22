# P27-S02-01 — Graph honesty implementation (INT-07)

## Metadata
- id: P27-S02-01
- todo_ids: [P27-S02-01, S02-T01, S02-T02, S02-T03, S02-T04, S02-T05, S02-T06]
- role: implementer
- skills: [incremental-implementation, test-driven-development, source-driven-development]
- mcps: [user-codegraph, user-trace]
- verification: mixed
- hooks: []

## Objective

Implement **INT-07 graph honesty** in product Go: extend `trace seed export --strict` so thin graphs and orphan causal entities fail under `--strict --enforce`, per **locked S02 planner defaults** (P27-S02-00). Align product enforcement threshold with harness **P25-3** (`discoveries≥1 OR decisions≥1`) while S01 keeps harness `--strict` warn-only until S03 VERIFY enables `--enforce`.

## References

- [00-PLANNER.md](00-PLANNER.md) — planner row (locked defaults source)
- [../scope-00-investigation/AUDIT.md](../scope-00-investigation/AUDIT.md) — S02-T01..T06 seeds, Risks table L126–136, Phase 26 thin-graph root cause
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- Live code (preflight read):
  - [`cmd/trace/seed.go`](../../../../../cmd/trace/seed.go) — `collectExportStructuralViolations`, `collectExportViolations`, strict path L178–195
  - [`internal/loop/gate.go`](../../../../../internal/loop/gate.go) — `GateForExport` ≡ `evaluateDone` L60–61, L259–261 STOP/ORIENT allow
  - [`internal/loop/gate_test.go`](../../../../../internal/loop/gate_test.go) — `TestEvaluateGate_Export_SameAsDone` L362–382
  - [`cmd/trace/enforce_test.go`](../../../../../cmd/trace/enforce_test.go) — strict/enforce fixtures L298–418
- Phase 26 baseline fixture:
  - [`experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-export-snippet.json`](../../../../../../experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-export-snippet.json) — `discoveries=0 decisions=0`, verify task STOP `p19_saturated`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P27-S02-00 — do not re-debate)

| Item | Locked value | Rationale (AUDIT Risks) |
|------|--------------|-------------------------|
| **Scope boundary** | Product Go only: `cmd/trace/seed.go`, `cmd/trace/enforce_test.go`, `internal/domain/*` (honesty helpers/tests), optional `internal/loop/gate_test.go` doc-only | INT-08/10 harness = **S01 done**; no `experiments/**` edits on this row |
| **Thin-graph strict rule** | **(A)+(B) hybrid in `seed.go` strict path:** (1) document-level min count `len(discoveries)≥1 OR len(decisions)≥1`; (2) when `discoveries>0`, each discovery must have `discovery_mentions_task` link in export `links[]`; when `decisions>0`, each decision must have `decision_affects_task` link | Matches P25-3 / P25-3a harness semantics; catches P26 build-only thin export. **Not** eval-rules.json-driven (defer **C**) — no `trace/eval-rules.json` body evaluation in S02 |
| **Export vs done gate split** | **(A) pre-structural + document honesty in `seed.go`; keep `GateForExport` ≡ `evaluateDone`** | Minimal correct diff: thin graph is document-level (`discoveries=0 decisions=0`), not fixed by splitting `evaluateExport`. STOP/ORIENT per-task allow stays in gate; honesty layer runs before gate append |
| **BLOCKING uncertainty rule** | **(A) store-backed promotion check:** for each discovery with `severity=BLOCKING` in store, export must include `discovery_mentions_task` link from that discovery ID | Severity absent from `SeedEntity`; promotion link is verifiable and matches P25-A gap-pass contract. **Do not** require separate uncertainty row (findings-based **B** deferred) |
| **Done-task gate skip** | **(A) keep skip** — `collectExportViolations` continues skipping `done/skipped/stale` tasks L121–124 | Honesty catches thin graph without re-evaluating completed tasks; avoids gate sweep churn |
| **Strict vs enforce behavior** | Honesty violations **always** print `seed export strict: …` on `--strict`; **only `--enforce`** returns `exitGateBlocked` and skips write (existing pattern L193–195) | S01 harness warn-only uses `--strict` without `--enforce`; build-only G1 will WARN after S02 |
| **Harness interaction** | **S02-01 does NOT edit `score.sh`** — upgrading harness to `--strict --enforce` is **S03 VERIFY** after S02-02 APPROVE | S01 locked warn-only until product rules exist; VERIFY row owns end-to-end P25-3 + enforce alignment |
| **Test fixtures** | **P26 export snippet** (`p26-export-snippet.json`) as thin-graph baseline; copy to `cmd/trace/testdata/` or embed in `enforce_test.go` | Reproduces Phase 26 `discoveries=0 decisions=0` STOP export that wrongly passed `--strict --enforce` pre-S02 |
| **Clean fixture update** | **`setupCleanFullCycleFixture` must gain ≥1 decision (or discovery) + required link** so `TestSeedExportStrictCleanAllowsWrite` still passes under min-count rule | Current clean fixture has zero discoveries/decisions |

## Hard boundary

```text
ALLOWED:  cmd/trace/seed.go, cmd/trace/enforce_test.go,
          internal/domain/seed_export_honesty.go (new helper OK),
          internal/domain/seed_export_honesty_test.go,
          internal/domain/seed_eval_rules_test.go (honesty tests only),
          cmd/trace/testdata/p26-export-snippet.json (fixture copy),
          internal/loop/gate_test.go (comment/doc only if export≡done unchanged)

FORBIDDEN: experiments/ab-p25-gap-pass-validation/** (S03 VERIFY),
           internal/loop/gate.go evaluateExport split (unless reviewer spawn),
           eval-rules.json body evaluation engine,
           daemon/HTTP/MCP product surfaces
```

If harness changes seem required, **stop** and note on board — schedule under S03 VERIFY.

## Preflight

Run from repo root before editing:

```bash
cd /home/ali/Desktop/Trace

# Product paths
test -f cmd/trace/seed.go
test -f internal/loop/gate.go
test -f internal/loop/gate_test.go
test -f cmd/trace/enforce_test.go
test -f internal/domain/seed_export.go

# Phase 26 thin-graph anchor (read-only)
test -f experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-export-snippet.json

# Build
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
test -x bin/trace

# Baseline tests (expect PASS pre-change)
go test ./internal/loop/... -count=1 -run TestEvaluateGate_Export_SameAsDone
```

If any path missing, mark row `blocked` with path in Notes.

## Files to change

| Path | Task IDs | Change summary |
|------|----------|----------------|
| `cmd/trace/seed.go` | T01, T02, T03 | Add `collectExportGraphHonestyViolations(ctx, st, doc)`; call from strict path after structural checks, before/with gate violations; store lookup for BLOCKING severity |
| `internal/domain/seed_export_honesty.go` | T01, T03 | Pure helpers: min-count, orphan discovery/decision link checks on `SeedDocument` |
| `internal/domain/seed_export_honesty_test.go` | T01, T04 | Unit tests for link/count rules on fixture JSON |
| `cmd/trace/enforce_test.go` | T05 | Thin-graph test from P26 snippet; enforce block + strict-warn-only paths; update clean fixture |
| `cmd/trace/testdata/p26-export-snippet.json` | T05 | Copy of Phase 26 export snippet (or testdata subset) |
| `internal/domain/seed_eval_rules_test.go` | T04 | **Optional:** table-driven honesty invariant test (hardcoded threshold, not eval-rules body) — pointer-only eval-rules export unchanged |
| `internal/loop/gate_test.go` | T06 | **No behavior change expected** — confirm `TestEvaluateGate_Export_SameAsDone` still passes unchanged |

**Do not modify** `internal/loop/gate.go` dispatch (`GateForExport` stays `evaluateDone`).

## Per-task acceptance criteria

### S02-T01 — `collectExportGraphHonestyViolations`

- [ ] New function (signature may be `(ctx, st, doc) []exportViolation` in `seed.go` or domain helper called from `seed.go`)
- [ ] **Min count:** if `len(doc.Discoveries)==0 && len(doc.Decisions)==0` → violation message contains `graph honesty` and `discoveries=0 decisions=0` (or equivalent stable substring)
- [ ] **Orphan discoveries:** for each discovery ID, require `links[]` entry with `rel=discovery_mentions_task` and `from` (or `from_id`) matching discovery ID
- [ ] **Orphan decisions:** for each decision ID, require `rel=decision_affects_task` link
- [ ] Violations use empty `TaskID` for document-level rules (consistent with structural violations)
- [ ] Called from `cmdSeedExport` when `*strict` **after** `collectExportStructuralViolations`, **before or merged with** gate violations

### S02-T02 — Export vs done gate (no split)

- [ ] `GateForExport` still dispatches to `evaluateDone` — **no** new `evaluateExport` function
- [ ] Document honesty runs in `seed.go` strict path; does not depend on STOP/ORIENT short-circuit in gate
- [ ] `--strict --enforce` on P26-equivalent thin graph **blocked by honesty**, not by gate sweep alone

### S02-T03 — Done-task gate skip (unchanged)

- [ ] `collectExportViolations` loop still skips `WorkStateDone`, `WorkStateSkipped`, `WorkStateStale` (L121–124)
- [ ] No new `--strict`-only full-task sweep in gate.go
- [ ] Verify-task STOP scenario covered by document min-count, not by forcing done-task evaluation

### S02-T04 — Tests / eval-rules boundary

- [ ] Domain unit test: P26 snippet JSON → honesty helper returns ≥1 violation for min-count
- [ ] Domain unit test: discovery without `discovery_mentions_task` → violation
- [ ] Domain unit test: decision without `decision_affects_task` → violation
- [ ] **Do not** implement eval-rules.json body parser; `TestSeedExportIncludesEvalRulesPath` unchanged

### S02-T05 — `enforce_test.go` integration

- [ ] New test: import or seed workspace from `p26-export-snippet.json` (or equivalent DB state) → `seed export --strict --enforce` → `exitGateBlocked`, no output file written
- [ ] Same fixture with `--strict` only (no `--enforce`) → `exitOK`, stderr contains honesty violation line, file **written**
- [ ] `TestSeedExportStrictCleanAllowsWrite` passes after clean fixture gets ≥1 decision + link
- [ ] Existing verification-debt / regression strict tests still pass

### S02-T06 — Gate test parity

- [ ] `TestEvaluateGate_Export_SameAsDone` passes **without modification** (export ≡ done unchanged)
- [ ] If test fails, **stop** — do not split gate without reviewer spawn; note on board

## Implementation order

1. Domain honesty helpers + unit tests (T01, T04)
2. Wire into `cmdSeedExport` strict path + BLOCKING store check (T01–T03)
3. P26 fixture + enforce tests + clean fixture update (T05)
4. Full test run (T06)

## Test commands (exit criteria evidence)

```bash
cd /home/ali/Desktop/Trace
CGO_ENABLED=1 go build -o bin/trace ./cmd/trace

# Primary gate
go test ./internal/... -count=1

# CLI enforce / export strict
go test ./cmd/trace/... -count=1 -run 'SeedExport|Enforce|Strict'

# Spot: thin graph must block enforce (manual after fixture wired)
# bin/trace -C <thin-fixture-dir> seed export -o /tmp/out.json --strict --enforce
# expect exitGateBlocked
```

Record PASS output summary in board Notes.

## Role work

Implement S02-T01→T06 in order. Self-check exit criteria before marking row `done`. **Board: status + notes only.**

## Exit criteria

- [ ] All locked defaults implemented (table above)
- [ ] P26-equivalent thin graph fails `--strict --enforce`
- [ ] Clean full-cycle fixture still passes `--strict --enforce`
- [ ] `go test ./internal/...` PASS
- [ ] `go test ./cmd/trace/... -run 'SeedExport|Strict'` PASS
- [ ] No edits under `experiments/`
- [ ] Next runnable: **P27-S02-02**

## Minimal todos

- [ ] S02-T01: honesty collector + wire strict path
- [ ] S02-T02: confirm no gate.go split
- [ ] S02-T03: confirm gate skip unchanged
- [ ] S02-T04: domain unit tests
- [ ] S02-T05: enforce_test + P26 fixture + clean fixture
- [ ] S02-T06: gate_test parity check
- [ ] Run test commands; update board row

## Next

`P27-S02-02`
