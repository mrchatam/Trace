# S05 — Harness install rules — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P23-S05-00 | scope planner | **done 2026-08-20** — locked paths, markers, cursor-hook contract, env vars; thickened 01+02 |
| 2 | P23-S05-01 | implementer | **done 2026-08-20** — rules, AGENTS.md, cursor-hook; 22 keeper tests green |
| 3 | P23-S05-02 | reviewer | **done 2026-08-20** — APPROVE (high confidence); git-hook unchanged; S06 unblocked |

**Depends on:** P23-S04-02 done (status violations + config). **Blocks:** S06 verify install smoke.

**S05-00 locks (summary):** extend cursor/claude; new `cursor-hook` (preToolUse → `trace loop gate --for edit`); AGENTS.md `# begin-trace-enforcement` block; git-hook **unchanged**; no auto `.trace/config.json`; `TRACE_TASK_ID` env contract.
