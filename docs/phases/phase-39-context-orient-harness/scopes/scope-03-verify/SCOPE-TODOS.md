# Scope 03 — board map (VERIFY)

**S03 VERIFY + DR-HANDOFF**. Serial: **P39-S03-00 → P39-S03-01 → P39-S03-02**.

| Order | Board ID | Prompt | Role | Status |
|------:|----------|--------|------|--------|
| 681 | P39-S03-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | **done** (S03-00) |
| 682 | P39-S03-01 | [01-verify.md](01-verify.md) | Verify | pending |
| 683 | P39-S03-02 | [02-dr-handoff.md](02-dr-handoff.md) | Closer | pending |

## Precondition (locked P39-S03-00)

| Review row | Verdict | Blocks gated |
|------------|---------|--------------|
| P39-S00-02 | **APPROVE** (high) | Block 0 |
| P39-S01-02 | **APPROVE** (high) | Block 1 |
| P39-S02-02 | **APPROVE** (high) | Block 2 |

## Verify blocks (P39-S03-00 — locked P39-00)

| Block | Check | Evidence prefix |
|-------|-------|-----------------|
| 0 | G1 T1–T6 + T1-MCP + S00-02 APPROVE | `$EVID/00-g1-*` |
| 1 | G3 G3-A1–A6 + 16 tools + S01-02 APPROVE | `$EVID/01-g3-*` |
| 2 | G4 G4-D1–D8 docs-only + S02-02 APPROVE | `$EVID/02-g4-*` |
| 3 | M-001 moat preserved | `$EVID/03-moat-*` |
| 4 | Laws 6–7 / 19 | `$EVID/04-law*` |
| 5 | Successor Phase 40+ (G5/G2 queue) documented | `$EVID/05-successor-*` |
| 6 | Graph export if entities changed | `$EVID/06-graph-*` |

**Evidence dir:** `experiments/runs/YYYY-MM-DD-p39-s03-01-verify/evidence/`

**Notes artifact (S03-01):** `scopes/scope-03-verify/VERIFY-NOTES.md`

**Review artifact (S03-02):** `scopes/scope-03-verify/REVIEW-NOTES.md`

## Accept maps (for VERIFY-NOTES)

### G1 (Block 0)

T1 `TestG1QueryHitMerged` · T2 `TestG1TaskMoatPreserved` · T3 `TestG1TitleFTSStillRunsWithQuery` · T4 `TestG1QueryExpandDedupe` · T5 `TestG1QueryCapHonesty` · T6 `TestG1QuerySearchFailOpen` · T1-MCP `TestMCPContextQueryMerged`

### G3 (Block 1)

G3-A1–A6 per S01 implement; `TestToolNamesRegistered` = 16

### G4 (Block 2)

G4-D1–D8 per S02 implement; touch-list `CONTRIBUTING.md` + `AGENTS.md` only

## Successor (locked — never TBD at close)

**Phase 40+ — Read surface & retrieval depth**

- **Entry:** G5 GUI graph orient + G2 unified `trace_explore`
- **Secondary:** G6, G7 (REMEDIATION-PLAN rank)
- **Scaffold owner:** P39-S03-02
- **Human gate:** Promote **P40-00**

See [DR-HANDOFF.md](../../DR-HANDOFF.md) forward note + [02-dr-handoff.md](02-dr-handoff.md) Phase 40+ scaffold table.
