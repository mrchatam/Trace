# Scope S04 — Phase VERIFY / compat+security checklist

**Depends-on:** `P08-S03-02` done (production hardening reviewed). Planner finalizes checklist names before VERIFY runs.

**S03 carry-in (2026-08-16):** Checklist must cover `migrate status`, backup↔restore (no BLOBs; Abs rebind), local `access.token` fail-closed, S02 lock, S01 adapter API version; no new MCP tools from S03.

**S04-00 locks (2026-08-16 FINAL):**

| Item | Value |
|------|-------|
| Checklist | **`evals/compat`** / **`TestCompatibilitySecurityChecklist`** |
| Schema / metrics | **`schema-compat.json` v1** + temp **`metrics-compat.json`** (`dry_run:false`) |
| Harness ownership | **S04-01 creates** (S01–S03 left none) |
| DR-HANDOFF | **`no successor`** (`A_PROJECT_PLAN` ends at Phase 8) unless Notes promote |
| Spawn | On fail: `01a` / `01b` (+`01c`) |

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P08-S04-00 | planner | done | 2026-08-16: FINAL checklist names + no-successor; thickened 01+02 |
| P08-S04-01 | verify | done | 2026-08-16: VERIFY PASS; harness created; DR-HANDOFF started=`no successor` — see VERIFY-NOTES.md |
| P08-S04-02 | review | pending | Owns `no successor` handoff close |

## Checklist

- [x] P08-S04-00 planner
- [x] P08-S04-01 VERIFY + VERIFY-NOTES
- [ ] P08-S04-02 review + DR-HANDOFF complete

## Phase context

- Validation gate: compatibility + security checklist — **`evals/compat`** / `TestCompatibilitySecurityChecklist` + `schema-compat.json` v1
- Successor: **`no successor`** (`A_PROJECT_PLAN` ends at Phase 8) unless Notes promote
