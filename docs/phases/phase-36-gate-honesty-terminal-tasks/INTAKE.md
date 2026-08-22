# INTAKE — plan_missing everywhere (feet-seller dogfood)

**Human 2026-08-22.** Project: `/home/ali/Desktop/feet seller telegram app`

## Symptom

After Phase 35, GUI lands on Loop 112. Clicking **any** task shows identical gate strip:

- **Gate blocked** / **`plan_missing`**
- `done blocked: recommended phase PLAN (plan_missing) → PLAN`

User request: **not a GUI patch** — find out **why**, decide if this is Trace vs agent misuse, fix **fundamentally**.

## Light facts (2026-08-22 investigation)

| Fact | Value |
|------|--------|
| Tasks | **123**, all **`DONE`** |
| Goals | **1** — `353b12a4-57dd-4d68-8379-b2024e064733` |
| Progressive planner | **Empty** — 0 phases, 0 scopes, 0 deep plans, 0 goal_plan_state |
| `plan_changes` in graph | **11** (rich titles: MVP plan, PIR loop pivot, grilling locked, …) |
| Discoveries / decisions | **45** / **24** |
| Reviews | **~120 PASS** (linked to tasks) |
| Deliberation states | **~0** exported / minimal loop history |
| `.trace/config.json` | **Missing** → enforce **`off`** |
| `trace install` artifacts | **None** — no project `AGENTS.md`, no `.cursor/rules/trace-enforcement.mdc` |
| MCP plan tools | **None** — `cmd/trace/plan.go` help: *"No MCP plan tools."* |
| `PlanExists` gate input | Requires `current_scope_id` **and** `current_deep_plan` only (`internal/loop/policy.go`) |

## Root-cause hypothesis (S00 must prove/disprove)

### A — Trace product: two planning systems, one gate

Trace has **two** planning representations:

1. **Causal graph:** `plan-change` entities (agents *can* create via `trace add` / MCP `trace_add`)
2. **Progressive planner:** `trace plan create-coarse|set-current|deep` → `plan_phases`, `plan_scopes`, `scope_deep_plans`

The deliberation loop **`PlanExists`** reads **only (2)**. Feet-seller has **11 plan-changes** but **never initialized (2)**. Gate is **correct per current code** but **wrong per agent mental model** ("we planned extensively").

### B — Harness: enforcement optional, install missing

Phase 23 enforcement is **opt-in** (`enforce: off` default; `--enforce` flag required on transition/export). Feet-seller never ran `trace install`, never set `enforce`, never wired hooks. Agents completed 123 tasks via domain transition (reviews + `as_operator`) **without ever passing** `loop gate --for edit`. Gates surfaced only when human opened GUI.

### C — Agent workflow: MCP path cannot satisfy PLAN gate

Primary agent surface is MCP. MCP exposes `trace_add kind=plan-change` but **not** `create-coarse` / `set-current` / `deep`. Agents following MCP naturally record planning as plan-changes, not progressive planner rows. Install rules say "run loop gate before edit" but **do not document** how to bootstrap `trace plan create-coarse`.

### D — Data model: one mega-goal

123 tasks under **one goal** without coarse structure amplifies goal-level `plan_missing` — every task inherits the same gate outcome.

## Verdict sketch (pending S00 formal write-up)

| Layer | Share of blame | Notes |
|-------|----------------|-------|
| Trace product design | **Primary** | Disconnected plan signals; MCP gap; gate checks store agents cannot populate via MCP |
| Agent misuse | **Secondary** | Could have run CLI `trace plan create-coarse` — but not discoverable from MCP/install path |
| Harness | **Secondary** | No install in consumer repo; enforce off → 123 DONE without gate ever blocking |
| GUI | **Tertiary** | Misleading copy on terminal DONE tasks; not the root cause |

## Desired outcome (fundamental)

1. **Single honest planning contract** — agents know what satisfies `PlanExists`; plan-changes and progressive planner aligned or bridged.
2. **Agent-complete surface** — MCP (or loop apply) can bootstrap coarse plan; install docs spell mandatory bootstrap on new goal.
3. **Dogfood recovery** — explicit path for feet-seller-like repos (backfill / import / acknowledge legacy).
4. **Gate display honesty** — terminal tasks don't scream "done blocked" when work is finished (secondary UX once root cause fixed).
5. **Evidence** — S03 live feet-seller + synthetic greenfield project showing correct agent workflow end-to-end.

## Not acceptable as "done"

- GUI-only hide of `plan_missing` on DONE tasks without addressing planning model
- "Agents should have known" without MCP/install/product fix
- Weakening PLAN gate for active work
