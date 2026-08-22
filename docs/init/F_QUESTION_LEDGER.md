# F — Question Ledger

Template fields: QUESTION · WHY IT MATTERS · CURRENT ANSWER · CONFIDENCE · EVIDENCE · WHO/WHAT CAN RESOLVE · RESOLUTION NEEDED BY · STATUS

Statuses: `OPEN` | `INVESTIGATING` | `ANSWERED` | `EXPERIMENT_REQUIRED` | `DEFERRED` | `ACCEPTED_RISK`

---

## Grill-Me Round 1 — RESOLVED (2026-08-15)

### Q-LANG → Go · Q-DAEMON → library+CLI only · Q-NAME → `trace` · Q-AGENT → CLI then MCP · Q-BENCH → synthetic · Q-SEED → human GT · Q-P0-DONE → answered by P00-S09-01 + P00-S09-02 (Phase 00 closed)

(See prior resolution log; unchanged.)

---

## Grill-Me Round 2 — RESOLVED (2026-08-15)

### Q-GOMOD

- **QUESTION:** Go module path?
- **WHY IT MATTERS:** Import paths; avoid placeholder migration.
- **CURRENT ANSWER:** **`github.com/mrchatam/Trace`** — real path now; GitHub repo already exists.
- **CONFIDENCE:** High.
- **EVIDENCE:** User Round 2 Q1=b.
- **STATUS:** ANSWERED

### Q-TRACEDIR

- **QUESTION:** Where does project state live?
- **CURRENT ANSWER:** **`.trace/` inside the bound repository**, created/maintained as **gitignored** project-local state. Large caches/indexes may split out later if scale requires.
- **CONFIDENCE:** High.
- **EVIDENCE:** User Round 2 Q2=a.
- **STATUS:** ANSWERED

### Q-P0X-PASS

- **QUESTION:** Exact pass criteria to close P0 after P0-X?
- **WHY IT MATTERS:** Must validate real foundation (causal + structural + incremental), not mere CRUD.
- **CURRENT ANSWER:** P0-X is complete when **all** of the following hold:

  1. Goal / Task / Decision / Discovery round-trip with provenance.
  2. Files + **minimal symbols/imports** can be represented.
  3. `trace why` returns a causal chain with evidence/reason codes.
  4. `trace context` produces **bounded** task-specific context.
  5. Human-seeded graph matches fixture ground truth.
  6. Deterministic tests pass **several** understanding queries without an LLM.
  7. **Incremental update** of a changed file works **without rebuilding the entire fixture graph**.

- **CONFIDENCE:** High.
- **EVIDENCE:** User Round 2 Q3=c + explicit 7-point bar. **Independent S09 VERIFY (2026-08-15):** `CGO_ENABLED=1 go test ./evals/p0x/... -count=1` and `./...` PASS; 7/7 evidence in [`docs/phases/phase-00-foundation/scopes/scope-09-phase-verify/VERIFY-NOTES.md`](../phases/phase-00-foundation/scopes/scope-09-phase-verify/VERIFY-NOTES.md).
- **STATUS:** ANSWERED (criteria locked + bar **met** by P00-S09-01)

### Q-P0-DONE

- **QUESTION:** When is roadmap P0 closable (before agent Gate C / X0)?
- **WHY IT MATTERS:** Separates foundation proof (P0-X) from agent understanding Gate C.
- **CURRENT ANSWER:** **P0 closable** — DR-P0X 7/7 independently re-proven (S09 VERIFY) and phase review (`P00-S09-02`) closed; Gate C remains Phase 01+.
- **CONFIDENCE:** High for foundation bar; Gate C still experiment-required.
- **EVIDENCE:** P00-S09-01 VERIFY-NOTES (7/7) + P00-S09-02 REVIEW-NOTES (fresh harness re-run APPROVE high) 2026-08-15.
- **STATUS:** ANSWERED — Phase 00 / P0 foundation bar closed; Phase 01 planner may unblock `P01-00`

### Q-ANALYZER

- **QUESTION:** TS/JS + Python structural analysis approach?
- **CURRENT ANSWER:** **tree-sitter from the beginning.** Code structure is fundamental; do not ship a temporary file-only architecture intended for immediate replacement.
- **CONFIDENCE:** High.
- **EVIDENCE:** User Round 2 Q4=a.
- **STATUS:** ANSWERED

### Q-GITLIB

- **QUESTION:** Git access strategy?
- **CURRENT ANSWER:** **`git` CLI subprocess** initially for complete/battle-tested behavior with low effort. **Abstract behind a VCS adapter interface** so the implementation can be replaced or supplemented later if benchmarks justify it.
- **CONFIDENCE:** High.
- **EVIDENCE:** User Round 2 Q5=a.
- **STATUS:** ANSWERED

---

## Planning principle (user directive, 2026-08-15)

> Do **not** optimize the early roadmap primarily for speed of implementation. Optimize early stages around **reducing the risk of building the wrong product**.

Implications recorded as **DR-RISK** and reinforced in laws/board: incremental indexing is a **P0-X architectural requirement**, not a later optimization; tree-sitter and minimal symbols/imports are in the foundation bar; full-rebuild-on-any-change architectures are forbidden.

---

## Grill-Me Round 3 frontier

No further questions block `T001`. Optional later (non-blocking):

### Q-UNDERSTAND-N

- **QUESTION:** How many “several” understanding queries for P0-X item 6?
- **CURRENT ANSWER:** **≥5** distinct deterministic queries — locked in S08 as: why-task, why-decision, decision-constraint, import-or-symbol-neighbor, context-boundedness (`docs/phases/.../scope-08-fixture-p0x/01-fixture.md`).
- **STATUS:** ANSWERED (2026-08-15 / P00-S08-00)

### Q-FIXTURE-LANG

- **QUESTION:** Synthetic fixture primary language mix?
- **CURRENT ANSWER:** Small **TypeScript/JavaScript + Python** polyglot under `fixtures/x0` (both analyzers exercised). Locked in S08 planner.
- **STATUS:** ANSWERED (2026-08-15 / P00-S08-00)

---

## Additional important questions

### Q-SCHEMA-CLAIM / Q-EVENTS-DEPTH / Q-IGNORE / Q-PROMOTE / Q-CONCURRENCY

Unchanged; provisional answers stand.

### Q-MULTI-PROJECT

- **CURRENT ANSWER:** One SQLite per project under **`.trace/`** (Q-TRACEDIR settled).
- **STATUS:** ANSWERED

### Q-H1-METRIC

- **STATUS:** EXPERIMENT_REQUIRED (X0 agent rubric; distinct from P0-X)

### Q-VECTOR / Q-UI / Q-COMMERCIAL

- **STATUS:** DEFERRED

### Q-INJECTION

- **STATUS:** OPEN (refine during T008)

### Q-DOGFOOD / Q-MCP-WHEN

- **STATUS:** ANSWERED (intent)

---

## Resolution log

| Date | ID | Change |
|------|-----|--------|
| 2026-08-15 | * | Ledger created |
| 2026-08-15 | Round 1 | Go, CLI-only, trace, CLI→MCP, synthetic, human GT, P0 until P0-X |
| 2026-08-15 | Q-GOMOD | ANSWERED → `github.com/mrchatam/Trace` |
| 2026-08-15 | Q-TRACEDIR | ANSWERED → `.trace/` in-repo, gitignored |
| 2026-08-15 | Q-P0X-PASS | ANSWERED → 7-point bar incl. structural + incremental |
| 2026-08-15 | Q-ANALYZER | ANSWERED → tree-sitter from day one |
| 2026-08-15 | Q-GITLIB | ANSWERED → git CLI + VCS adapter abstraction |
| 2026-08-15 | DR-RISK | Early stages optimize for wrong-product risk, not impl speed |
| 2026-08-15 | * | Round 2 closed; T001 unblocked |
| 2026-08-15 | Q-P0X-PASS / Q-P0-DONE | Independent P00-S09-01 VERIFY: 7/7 PASS; P0 closable pending P00-S09-02; evidence VERIFY-NOTES.md |
| 2026-08-15 | Q-P0-DONE | P00-S09-02 APPROVE high; Phase 00 complete / P0 closable; evidence REVIEW-NOTES.md |
