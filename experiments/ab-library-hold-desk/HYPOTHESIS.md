# E03 hypotheses — full-stack regression

Maps to Phase 24–28 interventions after residual wave close.

## H-E03-1 (default gap pass — FM-09)

**Claim:** Installed GapPassPrompt + write-before-export nudges produce ≥1 discovery OR decision on **build-only** without human gap prompt.

**Pass:** P25-3a PASS on `--arm build`  
**Fail:** Thin graph (E01 Session A replay)

## H-E03-2 (install + hook — INT-03/04 + Option A)

**Claim:** Fresh `trace install cursor --write` + hook includes gap pass, Parent orchestrator, and deny-without-task under strict.

**Pass:** P25-1/2 PASS; optional smoke: empty TRACE_TASK_ID + strict → deny

## H-E03-3 (honesty — INT-07)

**Claim:** Thin export fails `--strict --enforce`; rich export passes.

**Pass:** G2 honesty PASS after agent writes disc/dec

## H-E03-4 (promotion — INT-01/06 — soft)

**Claim:** If discoveries exist, agent may promote via `--from-discovery` / `spawned_tasks` (not required for E03 PASS).

**Soft pass:** ≥1 discovery↔task link or new task beyond seed  
**Not a hard fail:** seed-only tasks with rich disc/dec still validates H-E03-1

## Baseline comparison

| Metric | E01 Sess A | E02 build | E02-SB | E03 build target |
|--------|------------|-----------|--------|------------------|
| Domain | incident | checkout | checkout | **library hold** |
| Human gap prompt | No | No | Yes | **No** |
| Discoveries | 0 | 0→1* | 1 | ≥1 (or decision) |
| Stack version | P23 | P25–P28 | P28 | **post-P28** |

\*E02 build was thin until Session-B / later dual-lane rich G1.
