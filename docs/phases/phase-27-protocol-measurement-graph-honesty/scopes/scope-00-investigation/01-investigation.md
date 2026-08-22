# P27-S00-01 — Protocol + export path investigation

## Metadata
- id: P27-S00-01
- todo_ids: [P27-S00-01]
- role: implementer
- skills: [code-explorer, spec-miner, investigator]
- mcps: [user-codegraph]
- verification: mixed
- hooks: []

## Objective

Audit live `experiments/` protocol, `score.sh`, and `seed export` paths for INT-07/08/10 gaps. Produce `AUDIT.md` mapping Phase 26 P25-3 residual (`discoveries=0 decisions=0` on build-only G1) to **concrete file targets** for S01 (protocol v2 / INT-08+10) and S02 (graph honesty / INT-07). **No product code** on this row.

## References

- [00-PLANNER.md](00-PLANNER.md) — locked defaults + file audit table
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [Phase 26 VERIFY-NOTES](../../../phase-26-loop-implementation/scopes/scope-05-verify/VERIFY-NOTES.md)
- [Phase 24 CODEBASE-AUDIT](../../../phase-24-agent-effectiveness-investigation/scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md)
- [INTERVENTION-MATRIX.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-04-intervention-matrix/INTERVENTION-MATRIX.md) — INT-07/08/10 + P25-D/E themes
- Phase 26 evidence:
  - [`experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-s05-score.txt`](../../../../../../experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-s05-score.txt)
  - [`p26-export-snippet.json`](../../../../../../experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-export-snippet.json)
  - [`spot-checks.txt`](../../../../../../experiments/runs/2026-08-20-p26-s05-01-verify/evidence/spot-checks.txt)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| Deliverable | `scopes/scope-00-investigation/AUDIT.md` |
| Product code | **No** — audit artifact only |
| Harness root | `experiments/ab-p25-gap-pass-validation/` (E02/E03 verify arm) |
| Residual anchor | P25-3 FAIL on build-only G1 (`discoveries=0 decisions=0`); P25-1/2 PASS; **not** a P25-C regression |
| INT themes | INT-07 (export `--strict`), INT-08 (protocol v2 / score.sh), INT-10 (two-session rubric) |
| Threshold numbers | Document **options only** in AUDIT; S01/S02 planners pick — do not lock |

## Preflight

Run from repo root. All paths must exist (planner verified 2026-08-20):

```bash
cd /home/ali/Desktop/Trace

# Harness + prompts
test -f experiments/ab-p25-gap-pass-validation/PROTOCOL.md
test -f experiments/ab-p25-gap-pass-validation/RUBRIC.md
test -f experiments/ab-p25-gap-pass-validation/score.sh
test -f experiments/ab-p25-gap-pass-validation/prepare.sh
test -f experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-BUILD.md

# Product export / gate
test -f cmd/trace/seed.go
test -f internal/domain/seed_export.go
test -f internal/loop/gate.go
test -f internal/domain/seed_eval_rules_test.go

# Phase 26 verify evidence
test -f experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-s05-score.txt
test -f experiments/runs/2026-08-20-p26-s05-01-verify/evidence/p26-export-snippet.json
test -f experiments/runs/2026-08-20-p26-s05-01-verify/evidence/spot-checks.txt
```

If any path is missing, **stop** and mark row `blocked` with the missing path in Notes.

## Files to audit (minimum)

| Path | What to find | INT |
|------|--------------|-----|
| `experiments/ab-p25-gap-pass-validation/PROTOCOL.md` | Build-only vs directed-gap session modes; arm isolation; export step omission | INT-08, INT-10 |
| `experiments/ab-p25-gap-pass-validation/RUBRIC.md` | P25-3 pass criteria (`discoveries≥1 OR decisions≥1`); verdict matrix; expected build-only FAIL | INT-10 |
| `experiments/ab-p25-gap-pass-validation/score.sh` | `count_kind` dual export shapes (L45–57); P25-3 check (L119–128); G2 export presence; no `--strict` enforcement | INT-08 |
| `experiments/ab-p25-gap-pass-validation/prepare.sh` | `seed import` automated; **no** `seed export` — operator manual step | INT-08 |
| `experiments/ab-p25-gap-pass-validation/prompts/PROMPT-G1-BUILD.md` | Build-only wording; no gap-analysis instruction | INT-10 |
| `cmd/trace/seed.go` | `--strict` / `--enforce` flags (L136–148); `collectExportStructuralViolations` (L67–80); `collectExportViolations` → `GateForExport` (L83–129) | INT-07 |
| `internal/domain/seed_export.go` | `BuildSeedDocument` shape; `deliberation_states`, `uncertainties` export; what default export omits | INT-07 |
| `internal/loop/gate.go` | `GateForExport` shares `evaluateDone` with `GateForDone` (L60–61) — does export gate catch thin graph? | INT-07 |
| `internal/domain/seed_eval_rules_test.go` | Existing strict/export test coverage gaps | INT-07 |
| `experiments/runs/2026-08-20-p26-s05-01-verify/evidence/` | P25-3 FAIL line; export snippet counts | anchor |

Use `user-codegraph` `codegraph_explore` with `projectPath: /home/ali/Desktop/Trace` for call paths (`seed export` → `GateForExport` → `evaluateDone`) when line refs are unclear.

## Role work

### 1. Phase 26 residual (anchor)

Read VERIFY-NOTES + evidence files. Confirm:

- **Closure:** P25-2 PASS (Parent orchestrator in installed rules).
- **Residual:** P25-3 FAIL — `graph counts: decisions=0 discoveries=0` in `p26-s05-score.txt` L8–13.
- **Manual export gap:** VERIFY-NOTES L38 — `prepare.sh` does not export; operator must run `trace seed export -o runs/G1/trace/graph.json` before `score.sh G1 --p25`.
- **Not P25-C regression:** RUBRIC expected failure on build-only; install checks P25-1/2 passed.

**Repro steps** (document verbatim in AUDIT):

```text
1. CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
2. cd experiments/ab-p25-gap-pass-validation && ./prepare.sh G1
3. Run build-only G1 session (PROMPT-G1-BUILD.md) — no gap-analysis follow-up
4. trace seed export -o runs/G1/trace/graph.json   # NOT automated by prepare.sh
5. ./score.sh G1 --p25
   → expect P25-3 FAIL: discoveries=0 decisions=0
```

### 2. INT mapping (with line refs)

| INT | Theme | Audit questions | Likely scope |
|-----|-------|-----------------|--------------|
| INT-07 | Graph export honesty (`--strict`) | Does `--strict` today only run structural + `GateForExport` (done-gate parity)? Does it fail on `discoveries=0 decisions=0` thin graph? Are discovery→task / BLOCKING→uncertainty rules implemented? | **S02** |
| INT-08 | Protocol v2 | Is export operator-only? Does score.sh enforce import-before-gate? Git-sparsity / export-vs-DB drift handling? | **S01** |
| INT-10 | Two-session rubric | Are build-only vs directed-gap arms scored separately in PROTOCOL + RUBRIC? Is E02-Session-B isolated from build-only P25-3? | **S01** |

For each INT: list **path → current behavior → gap → proposed S01/S02 task seed** (file + function, not vague themes).

### 3. Delta from Phase 24 CODEBASE-AUDIT

Compare P24 FM-02/06/09 rows to post-Phase-26 state. Note what Phase 26 **closed** (P25-A/B loop policy, promotion, saturation reset) vs what remains open for INT-07/08/10. Cite P24 line refs where still valid; flag stale refs.

### 4. Write `AUDIT.md`

Use this template (fill all sections):

```markdown
# S00 Investigation AUDIT — INT-07/08/10

**Date:** YYYY-MM-DD  
**Row:** P27-S00-01  
**Git SHA:** (if available)

## Executive summary

(1 paragraph: P25-3 residual meaning, which INTs block measurement vs product behavior, recommended S01/S02 split.)

## Phase 26 residual

### Evidence

| Artifact | Path | Key finding |
|----------|------|-------------|
| Score output | `experiments/runs/.../p26-s05-score.txt` | … |
| Export snippet | `.../p26-export-snippet.json` | … |
| Spot checks | `.../spot-checks.txt` | … |

### Repro

(numbered steps from Role work §1)

### Interpretation

(Why P25-3 FAIL is expected on build-only; why it does not block Phase 26 closure.)

## INT mapping table

| INT | Path | Current behavior | Gap | S01/S02 task seed |
|-----|------|------------------|-----|-------------------|
| INT-07 | `cmd/trace/seed.go` | … | … | S02-T… |
| INT-08 | `score.sh` | … | … | S01-T… |
| INT-10 | `RUBRIC.md` | … | … | S01-T… |

(Add rows per file; minimum one row per INT.)

## Delta from Phase 24 CODEBASE-AUDIT

| P24 finding | Still valid? | Phase 26 change | Notes |
|-------------|--------------|-----------------|-------|

## Risks / open decisions (for S01/S02 planners)

| Topic | Options (do not lock) | Owner scope |
|-------|----------------------|-------------|
| P25-3 threshold | discoveries≥1 OR decisions≥1 (current) vs stricter | S01 |
| Export enforcement | warn-only `--strict` vs `--strict --enforce` in harness | S02 |
| … | … | … |

## Out of scope (this audit)

- Product implementation (S01/S02 rows)
- Re-scoring E01 historical runs
- Daemon / HTTP
```

## Exit criteria

- [ ] All preflight paths verified (or row `blocked` with missing path)
- [ ] INT-07, INT-08, INT-10 each have ≥1 mapping row with live path + line ref
- [ ] P25-3 repro documented with Phase 26 evidence citations
- [ ] `AUDIT.md` complete per template above in this scope folder
- [ ] S01/S02 task seeds are **actionable** (file + behavior gap), not theme-only bullets
- [ ] **No product code** changes (`git diff` shows only `AUDIT.md` if anything)

## Minimal todos

- [ ] Run preflight; record pass/fail in board Notes
- [ ] Read Phase 26 VERIFY-NOTES + evidence trio
- [ ] Audit each file in “Files to audit” table
- [ ] Cross-read P24 CODEBASE-AUDIT + INTERVENTION-MATRIX INT rows
- [ ] Write `AUDIT.md` from template
- [ ] Self-check exit criteria; update board row P27-S00-01

## Todo updates

Status + notes on **P27-S00-01** only.

## Next

`P27-S00-02`
