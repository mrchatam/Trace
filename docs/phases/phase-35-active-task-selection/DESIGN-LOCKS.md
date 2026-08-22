# DESIGN-LOCKS — Phase 35

**Human-promoted 2026-08-21.** Light clarifications from P35-00.

| Lock | Value |
|------|-------|
| Dogfood fixture | `/home/ali/Desktop/feet seller telegram app` (123 DONE tasks; Step 1 vs Loop 112) |
| Theme | **Active / current task selection** for GUI + loop gate display (+ agent-facing clarity) |
| Must fix | Default pick must **not** be “first row of `listTasks`” when that is stale DONE work while later tasks exist |
| Must test | Automated + live dogfood: Overview/Loop gate seed ≠ Step 1 when only DONE exist and later tasks exist; pagination/`limit` honesty when >100 tasks (client requests `limit: 100`; S00 verifies whether HTTP truncates today) |
| Selection home | Prefer **library/API** “current work” if one exists (Law 19); else **one** explicit GUI policy shared by Overview/Loop (extract helper — no divergent forks) |
| Agent story | Document `TRACE_TASK_ID` expectation; GUI default must not contradict “current work” |
| Out of scope | Weakening `plan_missing` gate for true PLAN-phase work; hosted SaaS; deleting feet-seller data |

## Success sketch

Opening `trace gui -C "<feet-seller>"` does not present Step 1 as the implied current task/gate when work has progressed to task ~123. User/agent can see and bind the sensible current task (last meaningful / current plan scope / explicit pick).
