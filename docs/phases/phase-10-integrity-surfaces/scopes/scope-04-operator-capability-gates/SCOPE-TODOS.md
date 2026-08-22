# S04 — Operator + capability gates — scope todos

**Depends-on:** P10-S03-02 done. Owns DF-17, DF-18, DF-24, DF-26, DF-31.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks; 01+02 thickened |
| 2 | 01-operator-capability-gates | implement | **done** — DF-17/18/24/26/31 |
| 3 | 02-scope-review | review | **done** — APPROVE high; no spawns |

## Locked reminders (from P10-S04-00)
- **DF-17:** `AllowOperatorDone` / `--as-operator` / `as_operator` — never trust `Actor` string
- **DF-18:** leave DONE → invalidate linked PASS to `UNCERTAIN`
- **DF-24:** fail-closed missing caps on **all** transitions; `--allow-missing-caps`
- **DF-26:** keep hatch; loud WARNING on use
- **DF-31:** usable `capability missing` without `--task`
- Honesty Path C **supersedes** to set operator flag; Gate G hatch unchanged
- MCP `trace_transition` must call same domain gates (G19)
- ab-operator-gate PROBE is acceptance evidence shape (agent DONE without `as_operator` fails)
- No new mig (default); no daemon/HTTP

## Next
**P10-S05-00** VERIFY planner.
