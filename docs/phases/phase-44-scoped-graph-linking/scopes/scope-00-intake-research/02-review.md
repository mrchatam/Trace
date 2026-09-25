# Phase 44 / S00 / Review

## Metadata
- id: P44-S00-02
- todo_ids: [P44-S00-02]
- role: reviewer
- skills: [code-review, research]
- mcps: [user-trace]
- verification: automated
- agents: []

## Objective

Independent review of [`RESEARCH.md`](RESEARCH.md) against [`INTAKE.md`](../../INTAKE.md) user examples and P44-00 light locks in [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md). Fresh context — do not share the researcher session. Docs-only; spawn forward if blocker/high gaps.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`INTAKE.md`](../../INTAKE.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md)
- [`01-research.md`](01-research.md) — required sections + acceptance
- [`RESEARCH.md`](RESEARCH.md) — artifact under review
- Board Notes on `P44-S00-00` / `P44-S00-01`

## Session start

Follow agent-loop-protocol Session start. Do not share implementer session. Unattended: INTAKE + light locks are the APPROVE authority (not the researcher's narrative alone).

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Confidence **medium** or **high**; residuals listed explicitly — never silent |
| Spawn | Insert `P44-S00-02a` (implement) / `02b` (review) immediately below this row if blocker/high |
| Product code | No (docs-only scope) |
| Lock reopen | Reject any RESEARCH that reverses D1–D5; spawn wording fix instead of editing locks here |

## Preflight

1. Confirm `RESEARCH.md` exists at this folder.
2. Spot-check 2–3 file:line cites against HEAD (center, MCP allowlist, caps).
3. Read INTAKE example table + D1–D5 summary before scoring.

## APPROVE rubric vs INTAKE

Score each row **PASS** / **FAIL**. APPROVE requires **all PASS** (or FAIL remediated via spawn still pending is **not** APPROVE).

### A. INTAKE example coverage

| ID | Check | PASS when |
|----|-------|-----------|
| A1 | Frontend auth | RESEARCH §4 names intra gaps (form/session/guard class) **and** cross-scope need for **`api_contract`** (or equivalent D2 cross rel) |
| A2 | Backend auth | RESEARCH §4 names intra gaps (handler/middleware/token class) **and** FE↔BE contract gap |
| A3 | Business scope | RESEARCH §4 addresses marketing/design/operator-style cluster **and** “not flatten to one goal” (center/hairball interaction) |
| A4 | Success criteria echo | RESEARCH ties gaps to clusters + typed edges + agent-create/infer + Laws 6–7 (may be in §1/§5/§6) |

### B. Hairball & live accuracy

| ID | Check | PASS when |
|----|-------|-----------|
| B1 | Hairball mechanism | §3 cites live files for center + `goal_has_task` fan-out (+ preferably GUI edge priority/cap) |
| B2 | MCP allowlist | States **5** `trace_link` aliases; no claim that scope rels are creatable today |
| B3 | Domain rel gap | Notes absence of D2 MVP set (`scope_member`, `api_contract`, `implements`, `blocks`) as first-class scope rels |
| B4 | plan_scopes | Notes planner `plan_scopes` exist but are **not** in graph walk |
| B5 | Caps honesty | Records GUI **500** / API **5000**; does not recommend silent raise |

### C. Design locks (non-contradiction)

| ID | Check | PASS when |
|----|-------|-----------|
| C1 | D1 | Gap = need thin scope records + `scope_member`; does not push tag-only (C) as primary |
| C2 | D2 | MVP rel set named; optional rels not required to block S02 |
| C3 | D3 | Notes missing `scope` filter / `scope_id` on `GraphNode` today |
| C4 | D4 | Affirms CLI-only inference; rejects daemon/index-hook as P44 path |
| C5 | D5 | Prefers scope filter over raising GUI 500 |

### D. Completeness & hygiene

| ID | Check | PASS when |
|----|-------|-----------|
| D1 | Seven sections | §1–§7 present and substantive |
| D2 | Name collisions | Mentions or safely distinguishes `implements`/`blocks` vs existing `change_implements_decision` / `uncertainty_blocks_task` (or flags as S01 wording risk) |
| D3 | Non-goals | §6 confirms DR-NOSSEM + no silent inference daemon |
| D4 | No product code | Researcher Notes / diff show docs-only for S00 |

## Findings format

```text
## Findings
- severity: blocker | high | medium | low | nit
- id: …
- evidence: …
- action: inline fix | spawn P44-S00-02a/02b | accept residual
```

| Severity | Action |
|----------|--------|
| blocker / high | Fix RESEARCH inline if trivial **or** spawn 02a/02b with full prompts |
| medium | Prefer spawn unless one-paragraph fix |
| low / nit | Notes OK; do not block APPROVE |

## Verdict block (required in board Notes)

```text
VERDICT: APPROVE | REJECT
Confidence: high | medium | low
INTAKE A1–A4: PASS/FAIL …
Locks C1–C5: PASS/FAIL …
Residuals: (none | list)
Next: P44-S01-00
```

**APPROVE** only if confidence ≥ medium **and** A1–A4 + B1–B5 + C1–C5 + D1–D3 all PASS (D4 verified).

## Exit criteria

- [ ] Findings by severity; blocker/high fixed or spawned
- [ ] APPROVE (or explicit REJECT + spawn) with confidence + evidence Notes
- [ ] `SCOPE-TODOS.md` S00-2 → done on APPROVE
- [ ] Next **P44-S01-00**

## Minimal todos

- [ ] Compare RESEARCH.md to INTAKE + locks using rubric tables
- [ ] Spot-check cites vs HEAD
- [ ] APPROVE or spawn
- [ ] Board Notes only

## Todo updates

Reviewer: status + notes; may spawn upcoming rows; **do not** edit `done` P44-00 / P44-S00-00 / P44-S00-01 prompts.

## Next

`P44-S01-00`
