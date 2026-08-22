# P35-S00-01 — Investigate

## Metadata
- id: P35-S00-01
- todo_ids: [P35-S00-01]
- role: implementer
- skills: [diagnosing-bugs, test-driven-development, planning-and-task-breakdown]
- mcps: []
- agents: []
- verification: automated
- hooks: []

## Objective

Author **only** `scopes/scope-00-investigate/INVESTIGATION.md`: live feet-seller repro, root-cause cites (Overview `pickActiveTask`, Loop default, `listTasks`/`limit` honesty), library “current work” search result, and a **red-capable** loop sketch for S02. **No product code.** Do not edit DESIGN-LOCKS, INTAKE (except if human asks), OpenAPI, CLI, or `web/`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law **19**
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [Phase 35 README](../../README.md)
- [00-PLANNER.md](00-PLANNER.md) — locked defaults + must-answer set
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- Live Trace (planner-verified 2026-08-21 — re-verify line numbers if drifted):
  - `web/src/screens/Overview.tsx` L17–20 `pickActiveTask`; L38 `listTasks({ limit: 100 })`; L44–50 gate/status on pick
  - `web/src/screens/Loop.tsx` L51–54 no `?task_id=` ⇒ `setParams({ task_id: res.items[0].id })`
  - `web/src/api/ops.ts` L41–45 `listTasks` forwards `limit`/`cursor` query
  - `internal/httpapi/handlers_tasks.go` L18–48 `handleListTasks` — reads `goal_id`/`work_state` only; **does not read `limit`/`cursor`**; returns full filtered `items`
  - `api/openapi.yaml` `/v1/tasks` documents optional `limit`/`cursor` (contract vs handler gap)
  - `internal/store/helpers.go` `ListTasks` / `ListTasksByGoalID` — `ORDER BY created_at ASC, id ASC`
  - `web/src/lib/overviewCompose.ts` + `overviewCompose.test.ts` — seed prioritization (related; confirm whether Explore/Overview graph seeds differ from `pickActiveTask`)
  - `cmd/trace/AGENTS.md` — `TRACE_TASK_ID` agent bind (no durable GUI current-work)
  - `web/src/context/AppChrome.tsx` — localStorage theme/token only (no task preference key)
- Fixture: `/home/ali/Desktop/feet seller telegram app`

## Session start

Follow agent-loop-protocol Session start. Unattended: prefer DESIGN-LOCKS/INTAKE; proceed without waiting for plan confirmation. **Do not mutate** feet-seller `.trace/` beyond read/query (no delete/reset/transition).

## Locked defaults (planner FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Artifact | `…/scopes/scope-00-investigate/INVESTIGATION.md` **only** |
| Product / CLI / web / OpenAPI edits | **Forbidden** |
| Fixture | Read-only dogfood path above |
| Primary symptom | Gate/Overview/Loop imply **Step 1** while human expects **~task 123** (Loop 112) |
| Secondary | `plan_missing` on DONE — note only; not the must-fix |
| Sequence | This row → `P35-S00-02`; S01 waits for S00-02 PASS |
| Feedback loop | Prefer CLI + HTTP JSON asserts; GUI screenshot/notes secondary |

## Planner baseline (re-confirm, do not invent)

| Fact | Value |
|------|--------|
| Task count | **123**, all **`DONE`** (CLI spot-check 2026-08-21) |
| Index 0 / Step 1 | `33247e2d-aa10-4b25-b194-4b7afb5a6359` — Step 1: Market and feature research |
| Index 122 / Loop 112 | `99d8fb92-65ac-462c-82c4-21bcf198c09e` — Loop 112: Entitlements polish + RESUME STOP |
| Store order | oldest-first (`created_at ASC, id ASC`) |

## Minimal todos

1. **CLI baseline** — `trace -C "<feet-seller>" tasks` → count, first id/title, last id/title, all DONE (or document exceptions).
2. **HTTP limit honesty** — start `trace serve -C "<feet-seller>"` (loopback), `GET /v1/tasks?limit=100` and without limit; record `items.length` and first id. Cite handler if length ≠ 100.
3. **Selection path cites** — confirm/refute Overview L17–20 + Loop L51–54 still match behavior; note if `overviewCompose.prioritizeTaskSeeds` is a separate surface.
4. **Simulate pick (no GUI required)** — given all-DONE list oldest-first, show that `pickActiveTask` and `items[0]` both yield Step 1 id (logic proof + optional GUI open Overview/Loop).
5. **Library “current work” search** — grep/codegraph for current-work / last-touched / plan-current / focused task APIs in `internal/` + HTTP routes; state **found** or **none**. Note `TRACE_TASK_ID` + no localStorage task key.
6. **Red loop outline** — one agent-runnable assert that fails on “defaults to Step 1” (see sketch below); paste a command you actually ran or a precise script outline S02 can turn into a test.
7. **Write INVESTIGATION.md** with required headings; board Notes ≤3 sentences.

## Must answer (all required in INVESTIGATION.md)

1. **Repro:** Commands + (if feasible) `trace gui` / HTTP: which task id/title Overview and Loop bind on open.
2. **Ordering:** Confirm store/API list order (oldest-first?).
3. **`limit` honesty:** With 123 tasks, does `/v1/tasks?limit=100` return 100 or 123? Cite handler behavior. Separate: would client `limit: 100` hide task 123 **if** pagination were honored?
4. **Root causes:** File:line for Overview pick + Loop default; refute or confirm each INTAKE bullet (1–5).
5. **Library gap:** Any existing “current work” / last-touched / plan-current API? If none, say so explicitly.
6. **Red loop:** One agent-runnable assert outline that fails while bug exists (CLI/HTTP preferred).
7. **Rejected:** List investigation rejects from SCOPE-TODOS.

## Suggested red-capable loop (outline for S02 — refine & run in 01)

Prefer a **pure logic / HTTP** assert over GUI:

```bash
# A) HTTP: prove list[0] is Step 1 while last is Loop 112 (symptom substrate)
#    After: trace serve -C "/home/ali/Desktop/feet seller telegram app" …
#    curl -s "http://127.0.0.1:$PORT/v1/tasks?limit=100" | jq '{n:(.items|length), first:.items[0].id, last:.items[-1].id}'
#    FAIL condition for selection bug (document as red until pick changes):
#      first == 33247e2d-… AND (simulated Overview/Loop pick == first)

# B) Unit-shaped (no product edit in S00): paste Overview pickActiveTask into a throwaway
#    node assert, OR describe exact Go/TS test S02 should add:
#      given 123 DONE tasks oldest-first, defaultPick(tasks).id !== step1Id
#      (today that assert goes RED)

# C) Optional GUI: open /overview and /loop without ?task_id=; note bound title in Notes
```

Red = assert that **default bound task ≠ Step 1** when later DONE tasks exist (DESIGN-LOCKS must-test). Green only after S02 policy change.

## INVESTIGATION.md template (required headings)

```markdown
# INVESTIGATION — Phase 35 S00

## Verdict (1 paragraph)

## Repro steps + evidence

## Root causes (with cites)

## limit / pagination honesty

## Library / API “current work” search

## Red-capable feedback loop (for S02)

## Secondary notes (plan_missing, etc.)

## Rejected alternatives

## Handoff to S01
```

## Exit criteria

- [ ] `INVESTIGATION.md` exists with all required headings filled
- [ ] Feet-seller evidence includes task count + first/last ids (or equivalent)
- [ ] At least one Overview and one Loop cite with line-level pointers
- [ ] `limit` behavior stated with evidence (HTTP response length **and** handler cite)
- [ ] Each INTAKE “likely cause” confirmed or refuted
- [ ] No product code changes
- [ ] Board Notes summarize verdict in ≤3 sentences
- [ ] Next is **P35-S00-02**

## Todo updates

Status + notes on **P35-S00-01** only.

## Next

`P35-S00-02`
