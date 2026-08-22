# P35-S01-01 — Plan

## Metadata
- id: P35-S01-01
- todo_ids: [P35-S01-01]
- role: implementer
- skills: [planning-and-task-breakdown, api-and-interface-design, test-driven-development]
- verification: automated
- mcps: []

## Objective

Author **only** `scopes/scope-01-plan/PLAN.md`: ranked selection-policy options, chosen default, Law 19 placement, file/API touch list, acceptance tests (incl. feet-seller all-DONE and >100-task/`limit` honesty). **No product code.** Do **not** re-open S00 root cause.

## References

- [INVESTIGATION.md](../scope-00-investigate/INVESTIGATION.md) — SoT for causes (S00-02 PASS)
- [00-PLANNER.md](00-PLANNER.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [INTAKE.md](../../INTAKE.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- `docs/rules/agent-loop-protocol.md` · Law 19

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Clarification needed only if a locked default below conflicts with live repo; otherwise proceed.

## Locked defaults (do not re-debate)

S00-02 PASS locked findings (treat as fact):

| Fact | Value |
|------|--------|
| Overview | `web/src/screens/Overview.tsx` `pickActiveTask` ~L17–20 → all-DONE ⇒ `tasks[0]` |
| Loop | `web/src/screens/Loop.tsx` ~L51–54 → no `?task_id=` ⇒ `items[0]` |
| Step1 | `33247e2d-aa10-4b25-b194-4b7afb5a6359` |
| Loop112 | `99d8fb92-65ac-462c-82c4-21bcf198c09e` |
| HTTP limit | `handleListTasks` ignores `limit`/`cursor`; `?limit=100` returns **123** |
| Library current-work | **None** (no HTTP/GUI durable focus; only agent `TRACE_TASK_ID`) |
| Red assert | exit **2** while default pick === Step1 |

| Item | Locked value |
|------|----------------|
| Artifact | **Only** `docs/phases/phase-35-active-task-selection/scopes/scope-01-plan/PLAN.md` |
| Dogfood | `/home/ali/Desktop/feet seller telegram app` — read-only evidence; **never** delete |
| Gate surfaces (must) | **Overview** + **Loop** share one pick policy |
| Related smell (optional mention) | Graph `prioritizeTaskSeeds` (`web/src/lib/overviewCompose.ts`) — same oldest-DONE odor; **not** a Phase-35 gate seed; may list as follow-on, not blocking acceptance |
| Store order | Oldest-first (`created_at ASC, id ASC`) — policy must not assume newest-first list |
| Client list | Overview/Loop call `listTasks({ limit: 100 })` via `web/src/api/ops.ts` |
| Out | Weakening `plan_missing` / PLAN-phase gates; TRACE_TASK_ID-only fix; delete/rewrite dogfood; hosted SaaS; divergent Overview vs Loop fallbacks |

### Placement options (rank in PLAN — closed set)

Law 19 + DESIGN-LOCKS: no existing library current-work → Phase 35 must still avoid business-logic forks in adapters.

| Rank id | Placement | When to choose |
|---------|-----------|----------------|
| **A** | Go library pure pick (+ thin HTTP if needed) + GUI adapters call it | Prefer if S02 vertical slice stays small (function + optional field/route + wire) |
| **B** | **One** shared TS helper under `web/src/lib/` (suggested name: `pickActiveTask.ts` / `pickCurrentTask.ts`); Overview + Loop import it; no duplicate `[0]` fallbacks | DESIGN-LOCKS “else” path — default lean if A does not fit one S02 implement row |
| **Reject** | Divergent Overview/Loop logic; TRACE_TASK_ID as sole fix; “just reverse list order” without shared policy; localStorage-only without shared pick | Violates DESIGN-LOCKS / Law 19 |

**Planner lean (normative until PLAN argues otherwise with sizing):** choose **B** for Phase 35 ship, with PLAN section “Future library promotion” noting A as Law-19 upgrade path — **unless** PLAN shows A fits S02 without scope creep. Do not invent a third placement.

### Policy semantics options (rank in PLAN — closed set)

URL / explicit `?task_id=` (and documented agent `TRACE_TASK_ID` for agents) **always win** when present.

When auto-picking from a task list, rank these (name them in PLAN):

| Id | Rule sketch |
|----|-------------|
| P1 | `IN_PROGRESS` first (unchanged) |
| P2 | Else first non-terminal (not `DONE`/`SKIPPED`) (unchanged) |
| P3a | Else **last** terminal / last by `created_at` on oldest-first list (`items[n-1]`) — fixes feet-seller |
| P3b | Else last **meaningful** (heuristic: title/body/plan-link) — only if PLAN defines cheap predicate |
| P3c | Else `planner.GetCurrentScope` / plan current — only if PLAN proves it applies when all DONE + no plan |
| P3d | Optional persist: localStorage last focused task id — **enhancement**, not sole fix |

**Planner lean:** **P1 → P2 → P3a** as default chosen policy. Reject keeping `tasks[0]` / `items[0]` as all-DONE fallback.

### Limit / pagination honesty (must cover in PLAN acceptance tests)

| Today | Risk |
|-------|------|
| Handler ignores `limit` | Dogfood still broken via `[0]` with full 123 |
| If limit honored at 100 | Loop112 at index **122** hidden from first page |

PLAN must: (1) state whether S02 **implements** HTTP `limit`/`cursor` or only documents honesty + client strategy (fetch enough / cursor to last page / dedicated current-work); (2) include an acceptance case for “>100 tasks + client `limit: 100`” so current work is not silently unreachable. **Lean:** selection fix is in-scope for S02; full OpenAPI pagination compliance is **in-scope only if** PLAN marks it required for honesty — otherwise schedule as explicit S02 optional or note for later phase, but **must not** leave “current work past page 1” unaddressed in acceptance criteria.

### Candidate S02 touch list (PLAN must finalize ⊆ this set)

| Layer | Paths |
|-------|--------|
| GUI must | `web/src/screens/Overview.tsx`, `web/src/screens/Loop.tsx` |
| GUI shared | new `web/src/lib/pickActiveTask.ts` (or chosen name) + unit test sibling |
| GUI optional | `web/src/lib/overviewCompose.ts` / Graph — only if PLAN pulls smell into scope |
| HTTP optional | `internal/httpapi/handlers_tasks.go`, `api/openapi.yaml`, `web/src/api/ops.ts` — if A or limit honesty |
| Library optional | new/existing `internal/…` pick helper — if A |
| Docs optional | `cmd/trace/AGENTS.md` — TRACE_TASK_ID alignment paragraph |
| Tests | Unit for pick; optional HTTP/integration from INVESTIGATION red script |

### Acceptance tests PLAN must specify (minimum)

1. **All-DONE + later tasks:** given oldest-first list `[Step1, …, Loop112]` all DONE → default pick **≠** Step1; prefer Loop112 under P3a.
2. **IN_PROGRESS wins:** mixed list → first `IN_PROGRESS`.
3. **Non-terminal wins over DONE:** no IN_PROGRESS → first PENDING/etc. before any DONE fallback.
4. **Explicit `task_id`:** Loop/URL override not overwritten by auto-pick.
5. **>100 / limit honesty:** document expected behavior when `limit=100` would hide index 122 — green strategy per PLAN.
6. **Red→green seed:** INVESTIGATION assert (`defaultPick ≠ Step1`) becomes automated test; today exit 2.

Fixture for narrative: feet-seller UUIDs above; CI should prefer **synthetic** in-repo fixtures (S02 will lock cmds) — PLAN must not require mutating dogfood.

## PLAN.md required headings

```markdown
# PLAN — Phase 35 active task selection

## Problem restatement (from INVESTIGATION)

## Options (ranked) + rejected

## Chosen policy (normative)

## Law 19 placement (library vs adapter)

## API / GUI / docs touch list

## Acceptance tests

## Agent TRACE_TASK_ID alignment

## Out of scope / non-goals

## Handoff to S02
```

Under **Options**, include both **placement** (A/B) and **semantics** (P1–P3*). Under **Handoff to S02**, list files, test names/cmds sketch, and what VERIFY (S03) should re-check on live feet-seller.

## Preflight / Plan

1. Re-read INVESTIGATION + DESIGN-LOCKS (do not re-curl unless PLAN needs a fresh cite).
2. Draft PLAN.md section-by-section against locked options above.
3. Self-check: Overview+Loop named; all-DONE + limit honesty in Acceptance tests; no product code.

## Role work

Write `PLAN.md` only. Update board row **P35-S01-01** status + Notes. Do not edit `done` S00 artifacts. Do not start S02 product work.

## Todo updates

Implementer: status + notes only on own board row.

## Exit criteria

- [ ] `PLAN.md` exists with all required headings
- [ ] Chosen policy + Law 19 placement explicit (A or B from locked set)
- [ ] Feet-seller all-DONE case + limit/>100 honesty covered in acceptance tests
- [ ] Touch list ⊆ candidate set; Overview + Loop mandatory
- [ ] No product code
- [ ] Board Notes cite PLAN path; next **P35-S02-00**

## Minimal todos

- [ ] Author `PLAN.md` per headings + locked defaults
- [ ] Rank placement A vs B and semantics P3*; state chosen + rejected
- [ ] Specify acceptance tests (incl. red-assert seed + limit honesty)
- [ ] Handoff section for S02-00 (files, tests, VERIFY expectations)
- [ ] Mark P35-S01-01 `done` with Notes

## Next

`P35-S02-00`
