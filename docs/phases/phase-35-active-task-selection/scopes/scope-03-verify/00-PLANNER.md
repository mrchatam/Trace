# P35-S03-00 — Scope planner (VERIFY)

## Metadata
- id: P35-S03-00
- todo_ids: [P35-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated

## Objective

Lock VERIFY floor: unit/e2e green + live `trace gui -C "<feet-seller>"` does **not** default gate/Overview/Loop to Step 1. Evidence paths under `experiments/runs/`. Thicken 01/02. **No product features** in this row.

## Floor (seed)

1. Automated tests from S02 still green (`node --experimental-strip-types --test src/lib/pickActiveTask.test.ts`)
2. Live dogfood: implied current task ≠ Step 1 (`33247e2d-…`) when opening Overview/Loop without `task_id`; prefer Loop112 (`99d8fb92-…`) under all-DONE
3. Spot-check `?task_id=` override preserved
4. Notes capture bound task id + title (Overview + Loop); read-only dogfood
5. DR-HANDOFF close owner = S03-02; default successor **no successor** unless VERIFY finds a thin follow-on
6. Note residuals: display `limit: 100` vs pick completeness; HTTP pagination future

## Next

`P35-S03-01`
