# X0 evaluator ground-truth map

**Evaluator / harness only.** Do **not** treat this file (or `fixtures/x0/seed/gt.json`) as an answer oracle for live Gate C B0/G1 agents. Harness code may import the seed JSON by absolute path for seeding and grading.

SoT for IDs and links remains `fixtures/x0/seed/gt.json` (seed JSON v1).

## UUID map

| Kind | UUID | Title / meaning |
|------|------|-----------------|
| **goal** | `11111111-1111-1111-1111-111111111111` | Ship greeter + math demo |
| **task** | `22222222-2222-2222-2222-222222222222` | Wire greeting to arithmetic helpers (`goal_id` → goal) |
| **decision** | `33333333-3333-3333-3333-333333333333` | Prefer TypeScript greeter surface |
| **discovery** | `44444444-4444-4444-4444-444444444444` | math_util lacks a clamp helper |
| **plan_change** | `55555555-5555-5555-5555-555555555555` | Add clamp helper to math_util |

## Causal story

1. **Goal** “Ship greeter + math demo” owns the **task** via `goal_id` + `goal_has_task`.
2. **Decision** “Prefer TypeScript greeter surface” **affects** the task (`decision_affects_task`).
3. **Discovery** that math_util lacks clamp **causes** the **plan_change** to add clamp (`discovery_causes_plan_change`).
4. Seed **transitions** the task `PENDING → IN_PROGRESS` with a non-empty reason (via `TransitionTask`, not CreateTask DONE shortcuts).
