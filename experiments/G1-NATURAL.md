# Natural G1 prompts

G1 arm prompts should read like a **product stakeholder ask**, not a Trace CLI tutorial.

## Do

- State the vague product requirement and feature list.
- Say “use Trace throughout” and point to `runs/G1/AGENTS.md` + installed cursor rules.
- Put operator-only env (`TRACE_BIN`, `TRACE_TASK_ID`, `TRACE_PROJECT_ROOT`) in a **Before you start (operator)** block — not in the agent’s first-turn mental model.
- Put enforcement mechanics in a short **Enforcement (mandatory)** section (gate before edit, `--enforce` on DONE/export).

## Do not

- Paste long `trace loop next` cookbooks into the product-facing requirement.
- Put Trace policy in `project/README.md` (B0 and G1 share the same product tree until prepare copies it).

## Why

Cookbook prompts inflate G1 wins — agents follow steps without internalizing deliberation. Natural prompts test whether standing harness (rules, hook, seed) is enough.
