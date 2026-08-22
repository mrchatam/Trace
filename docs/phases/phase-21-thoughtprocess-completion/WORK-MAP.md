# Phase 21 — Work map (residuals + gaps + promoted MVP cuts)

Source: [`forward-p20-residuals.md`](../../TODO/forward-p20-residuals.md) (superseded by this phase), P20 verify notes, gap audit vs [`TRACE_THOUGHTPROCESS.md`](../../TRACE_THOUGHTPROCESS.md). Attribution: [`DECISION-LOG.md`](DECISION-LOG.md).

| Work ID | Kind | TRACE § | Scope | Board rows |
|---------|------|---------|-------|------------|
| W-01 | Residual FR-01 | §31, §29K | S01 portable P20 seed | S01-00…02 |
| W-02 | Residual FR-02 | §22, §23 | S02 retrieval | S02-00…02 |
| W-03 | Residual FR-03 | §23 | S02 retrieval | S02-00…02 |
| W-04 | Gap | §4, §6, §21, §31 | S03 full SelectNext cycle | S03-00…02 |
| W-05 | Gap | §13 | S04 baseline promotion | S04-00…02 |
| W-06 | Gap | §31 | S04 promotion block on regression | S04-00…02 |
| W-07 | Gap | §25 | S05 trace why + deliberation audit | S05-00…02 |
| W-08 | Gap | §17 | S05 historical hints in loop next | S05-00…02 |
| W-09 | Residual FR-04 | §29O | S06 transactional apply | S06-00…02 |
| W-10 | Residual FR-10 | §6 | S06 goal_id validation | S06-00…02 |
| W-11 | Residual FR-12 | §29Q | S06 verify floor / test location | S06-00…02 |
| W-12 | MVP cut D-01 | §16 | S07 thin experiments | S07-00…02 |
| W-13 | MVP cut D-02 | §18 | S07 minimal risk-adaptive hints | S07-00…02 |
| W-14 | Residual FR-11 | §31 | S08 live mini-eval + verify | S08-00…02 |
| W-15 | Optional D-18 | §5 | S03 PolicyInputs enrichment | S03-00…02 |

**Still out of Phase 21:** hosted MCP, daemon, HTTP, autonomous test runner, ML policies, graph DB, Requirement table, committing `.trace/` (see DECISION-LOG D-16–D-22).

**Schema note (P21-00 locked):** S01 **no mig** (export/import only; compat **19**). S04 **`020_baselines_promotion.sql`** (compat **20**). S07 **`021_experiments.sql`** (compat **21**). S02/S03/S05/S06 protocol-only. Do not break P20 keeper tests.
