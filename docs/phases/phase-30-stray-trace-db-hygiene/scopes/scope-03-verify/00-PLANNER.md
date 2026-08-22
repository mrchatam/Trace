# P30-S03-00 — Scope planner (VERIFY)

## Metadata
- id: P30-S03-00
- todo_ids: [P30-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- verification: automated

## Objective

Lock VERIFY floor and thicken `01-verify.md` + `02-dr-handoff.md`.

## Floor

- `go test` for store/open + install (and any new packages touched)
- Repro: init in tmp; python create stub; Trace command still hits `.trace/`; warn observed if planned
- Confirm `trace init` does not create root `trace.db`
- Evidence dir: `experiments/runs/YYYY-MM-DD-p30-s03-01-verify/`

## Successor lean

Default **no successor** unless S00 found a larger store-path defect (document in DR-HANDOFF).

## Exit criteria

- [x] Verify + handoff prompts runnable (`01-verify.md` + `02-dr-handoff.md` thickened 2026-08-21)
- [x] Next: **P30-S03-01**

## Next

`P30-S03-01` → `P30-S03-02`
