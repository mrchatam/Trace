# DF-88 decision — keep seed exclude (clone PENDING expected)

P17-00 / [DF-84-FORWARD.md](../phase-17-portable-graph-git/DF-84-FORWARD.md) locked **default export exclude** of reviews, transitions, and task `work_state`. D40 scored clone PENDING as a D5 honesty deduct (**DF-88**). This file is the Phase 18 lock. It does **not** rewrite P17 `done` prompts.

## Grill vs laws (2026-08-17, unattended; human preference: do not explode portable-graph identity)

| Option | What clone gets | Laws | Reject / keep |
|--------|-----------------|------|----------------|
| **A. Keep exclude** (chosen) | Entities + links + plan tree; tasks **PENDING**; reviews absent | Law 1 (git = source; UUIDs = identity); Law 6 (no process dump); Law 11 (DONE is an explicit transition, not JSON replay); Law 9 (P17 user lock stands) | **Keep** |
| B. Default-export reviews + `work_state` | Clone looks DONE/SKIPPED | Mixes process history into portable identity; clone DONE without clone-side operator (DF-44 / DF-28 class) | **Reject** |
| C. CONDITIONAL `--include-reviews` / `--include-work-state` | Opt-in replay | Same identity explosion when used; extra flags this phase; easy to become default in agent habit | **Reject this phase** |

**Chosen: A.** Document clone PENDING as **expected**. `TestSeedExportOmitsDeniedSurfaces` remains a **fail bar** (must keep omitting). Do not add include flags. Do not steal DF-84.

## Clone honesty (what to write in S02)

- After `seed import`, tasks are **PENDING** unless the clone operator transitions them.
- `why` / `plan show` work from links + plan tree **without** reviews/`work_state`.
- Live DONE/SKIPPED in the exporter’s `.trace/` is **local process**, not git identity.
- D40 D5 deduct for “clone not DONE” is a **scoring confusion**, not a product bug, once docs/help say this.

## Do not

- Silently reverse P17 export exclude
- Treat `exported_at_commit` as identity
- Board hosted MCP as the way to “sync DONE”
