# S01 — Import path resolve — scope todos

**Depends-on:** P13-00 done. Owns **DF-60**.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks (resolve-time Expand helper) |
| 2 | 01-import-path-resolve | implement | **done** — DF-60 resolve helper + Expand + tests |
| 3 | 02-scope-review | review | **done** — APPROVE high; [REVIEW-NOTES.md](./REVIEW-NOTES.md); next **P13-S02-00** |

## FINAL summary (P13-S01-00)
- **Home:** `internal/retrieval` resolve helper + Expand import loop (A3 resolve-time)
- **Strategy:** relative `./`/`../` → join `Dir(importer)` + Clean/NormalizePath → try extensions (`.js/.jsx/.mjs/.cjs/.ts/.tsx/.go/.py`) + `index.*`; first indexed hit; bare modules skip
- **Keep green:** root `./util.js`→`util.js`; P12 exact-path provenance tests
- **No** analyzer rewrite / mig / path-align product hook
- Named tests + verify cmds in `01-import-path-resolve.md`

## Reminders
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- Optional dogfood `experiments/ab-import-resolve/` — not board blocker
- Next after APPROVE: **P13-S02-00** (DF-65 context hops **depend** on this resolve)
