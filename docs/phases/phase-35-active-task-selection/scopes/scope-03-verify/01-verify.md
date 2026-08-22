# P35-S03-01 — VERIFY

## Metadata
- id: P35-S03-01
- todo_ids: [P35-S03-01]
- role: verify
- skills: [test-driven-development, diagnosing-bugs, grinding-until-pass]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Run Phase 35 VERIFY floor after S00–S02 (**PLAN** live re-check + DESIGN-LOCKS must-test). Aggregate prior PASS cites + **live feet-seller** into **`VERIFY-NOTES.md`** (+ evidence dir). Prove Overview / Loop / gate do **not** default to Step 1 when all tasks are DONE. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P35-S03-02**. **No product code.** Do **not** mutate/delete feet-seller data. Do **not** start S03-02 or invent a successor phase.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S03-00)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S03-02
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — P1→P2→P3a; acceptance #1–6; VERIFY handoff
- Prior PASS artifacts (cite in notes):
  - S00 [`INVESTIGATION.md`](../scope-00-investigate/INVESTIGATION.md) + board Notes
  - S01 [`PLAN.md`](../scope-01-plan/PLAN.md) + board Notes
  - S02 board Notes (P35-S02-01/02) — `pickActiveTask` + `listTasksForPick` + Overview/Loop wire; review **PASS** 6/6
- Live fixture: `/home/ali/Desktop/feet seller telegram app` (**read-only**)
- Code anchors: `web/src/lib/pickActiveTask.ts`, `web/src/api/ops.ts` (`listTasksForPick`), `web/src/screens/Overview.tsx`, `web/src/screens/Loop.tsx`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or change product bodies.

## Locked defaults (FINAL — S03-00)

| Item | Value |
|------|-------|
| Precondition | P35-00 … P35-S02-02 all `done`; S02 review **PASS** (6/6 unit) |
| Product / Go / TS / OpenAPI changes | **Forbidden** (evidence + notes only). Failures → spawn remediation from this row or leave FAIL for S03-02 to spawn |
| Dogfood | **Read-only** — no delete, no status rewrite, no seed wipe |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p35-s03-01-verify/evidence/` |
| Notes artifact | `scopes/scope-03-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S03-02 closes |
| Successor | **Out of scope** — S03-02 only (lean default **no successor**) |
| Fixture path | `/home/ali/Desktop/feet seller telegram app` |
| Step 1 UUID | `33247e2d-aa10-4b25-b194-4b7afb5a6359` — must **NOT** be default pick |
| Loop 112 UUID | `99d8fb92-65ac-462c-82c4-21bcf198c09e` — expected P3a all-DONE pick |
| GUI launch | `trace gui -C "/home/ali/Desktop/feet seller telegram app"` (or `go run ./cmd/trace gui -C "…"` from Trace checkout) |
| Unit cmd | `cd web && node --experimental-strip-types --test src/lib/pickActiveTask.test.ts` |

### Narrative UUIDs (never mutate)

| Label | Full id | Role |
|-------|---------|------|
| Step1 | `33247e2d-aa10-4b25-b194-4b7afb5a6359` | Index 0 / oldest — **anti-pick** |
| Loop112 | `99d8fb92-65ac-462c-82c4-21bcf198c09e` | Index 122 / last DONE — **expect pick** under P3a |

### DESIGN-LOCKS + PLAN acceptance map (must tick in VERIFY-NOTES)

| Lock / case | Acceptance |
|-------------|------------|
| Must-fix | Default pick ≠ first `listTasks` row when later DONE work exists |
| Must-test automated | Unit suite green (acceptance #1–6 covered by S02; re-run here) |
| Must-test live | Overview active + Loop URL bind ≠ Step1; prefer Loop112 under all-DONE |
| Override | Explicit `?task_id=` wins (live evidence — **not** code-review-only) |
| Limit honesty | Document: today HTTP may still ignore `limit`; client `listTasksForPick` pages until exhausted |
| Out of scope | No `plan_missing` weaken; no feet-seller delete; no hosted SaaS |

### Aggregate evidence map (S00–S02 — cite in VERIFY-NOTES)

| Scope | Must cite / re-check |
|-------|----------------------|
| S00 | INVESTIGATION: 123 all DONE; Overview/Loop `[0]` root cause; HTTP `?limit=100` → 123 (handler ignores) |
| S01 | PLAN placement **B**; semantics **P1→P2→P3a**; fetch-for-pick honesty |
| S02 | `pickActiveTask.ts` + test 6/6; `listTasksForPick`; Overview+Loop wired; no `tasks[0]`/`items[0]` auto-pick; review PASS |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Unit suite non-zero exit
- Live Overview or Loop (no `?task_id=`) binds Step1 (`33247e2d-…`) as implied current / gate seed
- Live all-DONE pick is neither Loop112 nor another **last-list** task that is clearly not Step1 (prefer exact Loop112; if fixture drifted, record actual last DONE id + prove ≠ Step1)
- Explicit `?task_id=` overwritten by auto-pick on Loop (live)
- Dogfood mutated / deleted
- Product code shipped in this row
- `plan_missing` / PLAN-phase gates weakened
- VERIFY-NOTES missing or evidence dir absent after claimed PASS

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES — S02-02 fold-ins):

| Residual | Disposition |
|----------|-------------|
| **URL override** was code-review-only in S02 | **Must** produce live evidence here (Block 3); residual closed if live PASS |
| Display `listTasks({ limit: 100 })` vs pick full list | Accept if bind/gate uses `listTasksForPick` / pick path; note Loop `<select>` may omit P3a id **if** HTTP starts truncating — URL/gate bind still correct |
| `listTasksForPick` no max-pages guard | Accept as known low risk; note unbounded page loop if pathological `next_cursor` |
| HTTP still ignores `limit`/`cursor` (full 123 one page) | Accept — honesty satisfied by client page-through ready when pagination lands |
| No vitest / no React component test for override | Accept if live Block 3 PASS |
| Graph / `overviewCompose` smell | Out of phase |
| Placement A (Go current-work API) | Deferred — not required for PASS |
| TRACE_TASK_ID docs paragraph optional | Accept if absent |

## Locked verify command floor

Run from Trace repo root unless noted. Tee outputs into evidence dir. Use `date +%Y-%m-%d` for the run folder name.

### Block 0 — Evidence dir + preflight

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p35-s03-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P35-S03-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=P35-S02-02 PASS; S00–S02 done; PLAN P3a floor"
  echo "fixture=/home/ali/Desktop/feet seller telegram app"
  echo "step1=33247e2d-aa10-4b25-b194-4b7afb5a6359"
  echo "loop112=99d8fb92-65ac-462c-82c4-21bcf198c09e"
} > "$EVID/00-run-metadata.txt"
test -d "/home/ali/Desktop/feet seller telegram app/.trace"
```

**Pass:** `$EVID` exists; metadata cites S02-02 PASS; fixture `.trace/` present (do not write into it beyond normal GUI read).

### Block 1 — Automated unit suite (S02 re-check)

```bash
cd /home/ali/Desktop/Trace/web
node --experimental-strip-types --test src/lib/pickActiveTask.test.ts \
  2>&1 | tee "../experiments/runs/$(date +%Y-%m-%d)-p35-s03-01-verify/evidence/01-pickActiveTask-unit.txt"
```

**Pass:** exit 0; **6/6** (or current suite count ≥6) green. Cite command + exit in VERIFY-NOTES.

### Block 2 — Live dogfood: default pick ≠ Step1 (Overview + Loop)

Read-only. Launch GUI against feet-seller (prefer built `trace` on PATH, else `go run ./cmd/trace`):

```bash
# Example — do not background forever without noting PID; stop when screenshots/notes captured
trace gui -C "/home/ali/Desktop/feet seller telegram app"
# or: go run ./cmd/trace gui -C "/home/ali/Desktop/feet seller telegram app"
```

**Checklist (record id + title for each):**

1. Open **Overview** with **no** `?task_id=` in the URL.
2. Record **bound / active** task id + title (GateStrip / active panel — whatever Overview shows as current).
3. Assert active id **≠** `33247e2d-aa10-4b25-b194-4b7afb5a6359`.
4. Prefer active id **=** `99d8fb92-65ac-462c-82c4-21bcf198c09e` (Loop 112) under all-DONE P3a.
5. Open **Loop** with **no** `task_id` query (or clear it once and let auto-pick run).
6. After load, URL should gain `?task_id=<picked>` via replace — record that id + select label/title.
7. Assert Loop bind **≠** Step1; prefer Loop112.

Optional API sanity (read-only curl against the printed local URL / token if exposed — do not invent auth bypass):

```bash
# If serve prints base URL, optional:
# curl -sS "$BASE/v1/tasks?limit=100" | … count items; confirm first/last ids match INVESTIGATION
```

Capture: screenshot or copy-paste of Overview active + Loop URL into `$EVID/02-overview-loop-bind.txt` (and optional PNGs under `$EVID/`).

**Pass:** both surfaces ≠ Step1; prefer exact Loop112; notes include id + title for Overview **and** Loop.

### Block 3 — Live `?task_id=` override (S02-02 residual — mandatory evidence)

S02 review accepted override by **code review only**. VERIFY must close that gap with live proof:

1. Navigate to Loop with an explicit non-P3a id, e.g. Step1:
   - `/loop?task_id=33247e2d-aa10-4b25-b194-4b7afb5a6359`
2. Wait for load; confirm URL still shows that `task_id` (not replaced by Loop112).
3. Optionally pick another mid-list DONE id if Step1 triggers confusing gate UI — any explicit id that is **not** the auto-pick must stick.
4. Record before/after URL in `$EVID/03-task-id-override.txt`.

**Pass:** explicit `?task_id=` not overwritten by `pickActiveTask` / `setParams` auto path.

### Block 4 — Residual documentation (do not invent product work)

In VERIFY-NOTES, explicitly address:

| Residual | What to write |
|----------|----------------|
| Display vs pick truncation | Overview still uses `listTasks({ limit: 100 })` for display + `listTasksForPick` for pick; if HTTP later truncates, `<select>` options may omit last id while bind/gate remain correct |
| `listTasksForPick` max-pages | No max-pages guard today; page-until-`!next_cursor`; note as low residual |
| HTTP pagination future | Client ready; OpenAPI/`handlers_tasks.go` pagination still deferred |

**Pass:** residuals section present; none promoted to FAIL unless they cause Block 2/3 failure.

### Block 5 — WRITE VERIFY-NOTES.md

Create `docs/phases/phase-35-active-task-selection/scopes/scope-03-verify/VERIFY-NOTES.md` with at least:

```markdown
# VERIFY-NOTES — Phase 35 / S03-01

**Date:** …
**Overall:** PASS | FAIL
**Git SHA:** …
**Evidence:** experiments/runs/…-p35-s03-01-verify/evidence/

## Precondition cites
- S00 / S01 / S02 board Notes + review PASS

## Block results
| Block | Result | Evidence file |
|-------|--------|---------------|
| 0 preflight | | |
| 1 unit | | |
| 2 live Overview+Loop | | |
| 3 URL override | | |
| 4 residuals | | |

## Live binds (required)
| Surface | Bound task id | Title | ≠ Step1? | = Loop112? |
|---------|---------------|-------|----------|------------|
| Overview | | | | |
| Loop (no initial task_id) | | | | |

## Override spot-check
| Explicit task_id | URL after load | Overwritten? |
|------------------|----------------|--------------|
| … | … | no |

## Residuals (S02-02 fold-ins)
- …

## Failures / spawns
- none | spawn id …

## DR-HANDOFF
Still OPEN — close owner P35-S03-02. Successor decision deferred to S03-02 (default no successor).
```

## Exit criteria

- [ ] Blocks 0–4 executed; Block 5 `VERIFY-NOTES.md` written
- [ ] Unit green; live Overview + Loop ≠ Step1; prefer Loop112
- [ ] Live `?task_id=` override evidenced
- [ ] Residuals documented (display truncation; max-pages; HTTP pagination)
- [ ] Dogfood untouched; no product code; DR-HANDOFF still OPEN
- [ ] Board Notes cite evidence path + bound ids; status `done` or `failed` with spawn

## Minimal todos

- [ ] Block 0 evidence dir + metadata
- [ ] Block 1 unit re-run → tee
- [ ] Block 2 live Overview + Loop binds → record id/title
- [ ] Block 3 live URL override → tee
- [ ] Block 4 residuals section
- [ ] Block 5 VERIFY-NOTES.md
- [ ] Board status + Notes; hand off **P35-S03-02** (do not close DR)

## Next

`P35-S03-02`
