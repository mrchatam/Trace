# P14 / S01 / 02 — Scope review (Impact walks) FINAL

## Metadata
- id: P14-S01-02
- todo_ids: [P14-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S01 impact walks. Fresh subagent. Compare claims + **00-PLANNER FINAL** to live code/tests. Spawn `02a`/`02b` for blocker/high. **Stop if 00-PLANNER is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-impact-walks.md](01-impact-walks.md) — **FINAL**
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A3/A5/A6
- Research rank 6
- Live: `internal/retrieval`, `cmd/trace/impact.go`, `evals/impact`, Expand neighbors
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Do not share the implementer’s session. Re-run verify; do not trust Notes alone.

## Checklist (FINAL)

| # | Check | How to evidence |
|---|--------|-----------------|
| 1 | Multi-seed BFS + seed exclusion per FINAL | `TestImpactWalkMultiSeedExcludeSeeds` + read walk code (one `seen`, seeds out) |
| 2 | Contains asymmetry per FINAL (no sibling climb) | `TestImpactWalkContainsAsymmetryNoSiblings` + file↔symbol edge rules |
| 3 | Incoming deps only; outgoing not used as blast | Code review of neighbor selection; `TestImpactWalkIncomingImportHop` |
| 4 | Loud truncation/totals; no silent caps | `TestImpactWalkLoudTruncation` — `blast_total`/`blast_kept`/`truncated` |
| 5 | Hop risk monotonic | `TestImpactWalkHopRiskIncreases` |
| 6 | Existing impact / Gate F prelim not weakened | `TestPlantedImpactConflictsGateFPrelim` + `finding`/`alternative`/`report` still work |
| 7 | Expand/Why not “fixed” by removing bidirectional contains for context | Expand still symbol→file for Why/context; ImpactWalk is separate |
| 8 | G19 — no domain fork in CLI; no new MCP impact menu | CLI thin; no `trace_impact` MCP tool added |
| 9 | No daemon/HTTP/embeddings/Neo4j/full-rebuild; no mig; no `internal/impact` | Diff + schema dir |
| 10 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + Gate C `dry_run:false` | Re-run 01 verify cmds; Gate C artifacts untouched |
| 11 | Board Notes accurate; planner row had no product Go | Diff scoped to S01-01; P14-S01-00 Notes docs-only |
| 12 | DR-NOIMP — no commercial auto-engine / auto-plant findings from walk | Walk is structural report only |

## Review procedure
1. Read FINAL locks in 00 + 01.
2. Diff implement vs inventory gaps; map each named test to code.
3. Fresh verify suite from 01.
4. Decide APPROVE / spawn `P14-S01-02a`/`02b` for blocker/high.
5. Write `REVIEW-NOTES.md` (optional but preferred) + board Notes; next **P14-S02-00** unless spawn.

## Exit criteria
- [x] Checklist evidenced; confidence high (or medium with residuals)
- [x] Board status + Notes; next **P14-S02-00** (unless spawn)
- [x] No rewrite of done P14-S01-00/01 history — spawn forward only
