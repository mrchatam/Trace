# FM-10 / FR-P28-07 — notes (P28-S06-13)

**Date:** 2026-08-20  
**Status:** measured (promotion API used in scripted loops; acceptance met; no product code)  
**D1:** auto-spawn **not** implemented (human gate preserved)

## Objective

Ensure the promotion API is exercised in live/scripted loops; measure promotion rate; document residual risk for build-only 0-discovery exports. M-01 already covers apply promote — left unchanged.

## Acceptance

Directed/build dogfood shows **≥1 discovery linked to task** **or** spawned task from BLOCKING — **met** via both branches:

| Path | Evidence |
|------|----------|
| discovery↔task (dogfood) | G1 export: `discovery_mentions_task` `c4d9a046-…` → `e0200000-…0010` |
| BLOCKING → spawned/promoted task (scripted) | M-01 apply `spawned_tasks[].discovery_id`; CLI `--from-discovery` demo; domain promote suite |

## Scripted promotion (API usage)

Re-ran 2026-08-20 (`GOPROXY=direct`):

| Suite | Tests | Result |
|-------|-------|--------|
| M-01 apply E2E | `TestLoopApplyPromotesBlockingDiscoveryViaSpawnedTask` | **PASS** |
| Apply no-auto-spawn | `TestLoopApplyDiscoveriesOnlyDoNotCreateTasks` | **PASS** (discoveries alone do not create tasks) |
| Domain BLOCKING promote | `TestPromoteBlockingDiscoveryCreatesAndLinksTask`, `IsIdempotent`, `FailsClosed`, `PreservesSeedWorkState`, `AfterImport` | **PASS** (5/5) |
| CLI promote | `TestAddTaskFromDiscoveryCLI` | **PASS** |

Also matched by `-run 'Promote'` in domain: 5 VCS/baseline promote tests + 5 discovery promote — all PASS (10 total under that filter).

### Temp-project CLI demo (live scripted)

```text
trace init
trace add goal --title 'FM10 demo goal' --id aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa
trace add discovery --title 'BLOCKING gap for FM10' --severity BLOCKING
  → id=ddf433fc-3159-4d61-834f-007167ed5278
trace add task --from-discovery ddf433fc-3159-4d61-834f-007167ed5278 --goal-id aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa
  → task id == discovery id (PROMOTE_ID_MATCH=ok)
seed export → links: discovery_mentions_task from=to=ddf433fc-… ; tasks=1 discoveries=1
```

No `prepare.sh`. Temp dir only (not G1 mutation).

## Dogfood / G1 (read-only)

| Source | Finding |
|--------|---------|
| [`SESSION-B-NOTES.md`](../scope-02-session-b-dogfood/SESSION-B-NOTES.md) | Session-B added discovery + `discovery-mentions-task` → seed task `…0010` (PLAN_AFFECTING link-to-existing; **not** BLOCKING `--from-discovery` spawn) |
| [`runs/G1/trace/graph.json`](../../../../../experiments/ab-p25-gap-pass-validation/runs/G1/trace/graph.json) | Still has `discovery_mentions_task` + disc=1; tasks remain seed roster (5) |
| [`SESSION-A-GRAPH-SNAPSHOT.json`](../scope-02-session-b-dogfood/SESSION-A-GRAPH-SNAPSHOT.json) | Thin baseline disc=0/dec=0 |

## Promotion rate

| Context | Rate / interpretation |
|---------|------------------------|
| Scripted promote-path tests (M-01 + 5 domain BLOCKING + CLI FromDiscovery) | **100% PASS** (7/7 primary promote asserts; discoveries-only control PASS) |
| Live temp CLI BLOCKING→task | **1/1** promote succeeded |
| Session-B dogfood BLOCKING→spawn via promote API | **0** this wave (linked existing task instead) |
| Build-only / Session-A thin (disc=0) | **N/A — zero opportunity** (no discoveries → no candidates → no promotes) |

## Residual risk (build-only 0-discovery exports)

- Build-only thin graphs (`discoveries=0 decisions=0`) **cannot** exercise promotion; P25-3a FAIL on thin is expected richness failure, not a missing promote API.
- Agents must **write discoveries first** (FM-02) then **explicitly promote** (FM-01 candidates / FM-08 nudge / `--from-discovery` or `spawned_tasks[].discovery_id`).
- Without that sequence, live promotion rate stays 0 even though the API is shipped and green in tests.
- Auto-spawn on import/add discovery remains **out of scope** (FR-P28-D1).

## Human gate (FR-P28-D1)

No auto-spawn. Expansion requires explicit:

- `trace add task --from-discovery <discovery_id>`, or
- `trace loop apply` with `spawned_tasks[].discovery_id`

## Product code

**None.** M-01 sufficient; no assert strengthen needed.

## Exit

Acceptance hint met. Next runnable: **P28-S06-14** (FM-10 review).
