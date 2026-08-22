# E01 — Cursor Multitask mode

When the **parent agent orchestrates subagents**, follow this split.

## Root cause of past failures

| Failure | Why | E01 fix |
|---------|-----|---------|
| Wrong experiment path | Parent workspace = Trace repo; subagent searched familiar paths | Open **run folder** first; `pwd` gate; workspace guard rule |
| Orchestrator skipped Trace | Parent delegated coding; only workers touched `.trace/` | G1 prompt: orchestrator **must** run loop on `…0010` before subagents |
| Post-hoc graph | Workers marked DONE without decisions/uncertainties | Rubric G3/G4/G6; `--enforce` on export |

## Recommended split

### Orchestrator (G1)

1. `pwd` — confirm `runs/G1`
2. `trace loop next` / `apply` on planning task `e0100000-0000-4000-8000-000000000010`
3. Write **SPEC.md**, **PLANNING-MATRIX.md**, record **≥3 decisions** and **≥1 resolved uncertainty**
4. Transition planning task toward DONE with `--enforce` when gates allow
5. Delegate **any** Trace task (seed or newly added) using [SUBAGENT-DELEGATION.md](prompts/SUBAGENT-DELEGATION.md). The seed’s four implementation rows are examples, not a cap. Extra agents (reviewer, tester, schema) are allowed.

### Subagents

1. Receive absolute path + `TASK_UUID` + env block
2. `trace context` for that task
3. `loop gate --for edit` before product edits
4. Implement slice; evidence + `--enforce` DONE for **their** task only

### SQLite lock

Serialize Trace **writes** from one coordinator at a time, or run **one implementation subagent** at a time.

## B0 Multitask

No Trace. Orchestrator still must verify workspace and pass `/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/B0` in every delegation.

## Full Trace, not a four-agent roster

Seed tasks `…0020`–`…0050` are a starter backlog. Agents may `trace add` more tasks, run reviews, record discoveries/plan_changes, and spawn more workers. Restricting G1 to four implementers under-tests Trace.

## Answer: orchestrator + subagents both using Trace?

**Yes for G1.** Orchestrator owns planning graph; subagents own gated implementation. Subagents-only Trace produces fragmented plans and wrong-path searches.
