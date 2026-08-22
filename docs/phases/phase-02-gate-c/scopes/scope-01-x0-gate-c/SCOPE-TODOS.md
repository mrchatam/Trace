# Scope S01 — X0 Gate C evaluation

**Depends-on:** Phase 01 complete (`P01-S05-02` done); `P02-00` done; this scope’s `00-PLANNER` done.

**Bar:** Documented Gate C evidence (pass/fail/iterate) with B0 vs G1 agent runs (**N≥3**), understanding primary, kill criteria applied. **Phase 01 dry-run ≠ Gate C pass.**

**Locked by P02-S01-00:**
- Instrument: extend `evals/x0`; keep dry-run green
- Evidence: `GATE-C-NOTES.md` + `docs/verification/gate-c-x0/`
- Scoring: same ≥5 query bank; mean `understanding_accuracy`; critical_miss defined
- Kill: G1 mean ≤ B0 mean **and** non-trivial seeding (`fixtures/x0` seed = non-trivial)
- MCP optional (CLI path); honesty + p0x stay separate/green

**Issue-list out:** `GATE-C-NOTES.md` issue list (GC-NN shape) → S02 backlog after S01-02.

- [x] P02-S01-00 planner
- [ ] P02-S01-01 implement
- [ ] P02-S01-02 review
