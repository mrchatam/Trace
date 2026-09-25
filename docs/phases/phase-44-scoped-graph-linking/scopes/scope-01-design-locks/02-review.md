# Phase 44 / S01 / Review — Design locks

## Metadata
- id: P44-S01-02
- todo_ids: [P44-S01-02]
- role: reviewer
- skills: [code-review, domain-modeling]
- verification: automated

## Objective

Independent lock review before the S02 implement wave. Confirm [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) is **LOCKED**, D1–D5 closed without silent reversal of P44-00 light locks, and Laws **6–7 / 13 / 19** plus **M-001** are respected. Block S02 start if locks incomplete. Docs-only; spawn forward if blocker/high gaps.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — artifact under review
- [`01-lock.md`](01-lock.md) — required edits E0–E8 + acceptance A1–A10
- [`../scope-00-intake-research/RESEARCH.md`](../scope-00-intake-research/RESEARCH.md) — APPROVE baseline
- [`INTAKE.md`](../../INTAKE.md)
- Phase README moat charter (M-001 / Law 13)

## Session start

Follow agent-loop-protocol. **Fresh subagent** — do not share the locker (P44-S01-01) session. Unattended authority: P44-00 light locks + RESEARCH APPROVE + Laws below.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Confidence **medium** or **high**; residuals explicit — never silent |
| Spawn | Insert `P44-S01-02a` / `02b` immediately below this row if blocker/high |
| Product code | No |
| Reversal | Any D1–D5 flip vs light lock without REVIEW spawn reason → **FAIL** / spawn |

## Preflight

1. Confirm `01-DESIGN-LOCKS.md` status line says **LOCKED (S01)**.
2. Diff mentally (or `git diff`) against light-lock intent: D1 A+thin, D2 four MVP rels, D3 scope filter, D4 CLI, D5 GUI 500.
3. Spot-check RESEARCH §7 risks are addressed as **wording**, not choice flips.
4. Confirm INTAKE auth FE/BE + business remain implementable under locked taxonomy.

---

## Independent lock rubric

Score each row **PASS** / **FAIL**. APPROVE requires **all PASS** (or FAIL remediated via spawn still pending ≠ APPROVE).

### L. Document hygiene

| ID | Check | PASS when |
|----|-------|-----------|
| L1 | Status | Banner is **LOCKED (S01)** with date; no “light-locked pending S01” |
| L2 | D1–D5 closed | §9 (or equivalent) lists each decision with final choice + rationale; **no TBD** |
| L3 | No silent reverse | Choices match P44-00 light locks (A+thin; four MVP rels; scope filter; CLI infer; GUI 500) |
| L4 | RESEARCH §7 | Name collisions, plan_scopes≠graph scopes, seed-export, aspirational GUI tokens, optional `edge.provenance` are **worded** in locks (or explicitly deferred with reason that does not reopen D*) |

### A. Laws 6–7 (bounded graph)

| ID | Check | PASS when |
|----|-------|-----------|
| A1 | Law 6 | Locks keep `max_nodes` required; API hard max **5000**; reject unbounded / no full dump default |
| A2 | Law 7 | Scope filter is **bounded N-hop** + `max_nodes`; progressive expand preserved; **no tiles** in P44 |
| A3 | D5 | GUI default **500** retained; prefer scope filter over raising cap; raise only via VERIFY evidence |

### B. Law 13 (no parallel relationship store)

| ID | Check | PASS when |
|----|-------|-----------|
| B1 | Store | Membership + MVP rels live in **`entity_links`** (+ thin scope **records**, not a second edge table / graph DB) |
| B2 | Reject C | Tag-column primary rejected; no “scopes[] JSON as primary graph model” |
| B3 | plan_scopes | Soft title map only; **no** hard FK mandate that forces planner table to double as graph edge store |

### C. Law 19 (adapters)

| ID | Check | PASS when |
|----|-------|-----------|
| C1 | Walk/filter | Scope filter / membership resolution owned by **`internal/retrieval`** (library), not `web/` |
| C2 | GUI | S04 consequences are layout/legend/priority only — no client-side inventing of membership |
| C3 | HTTP/MCP | Documented as thin adapters extending allowlists/OpenAPI — no forked business rules |

### D. M-001 moat

| ID | Check | PASS when |
|----|-------|-----------|
| D1 | Loop | Locks state scope cartography **enriches** orientation; does **not** replace Tasks → Loop → Gate → Review / board order |
| D2 | Causal rels | Preserved rels (`goal_has_task`, `review_judges_task`, decision/discovery causal set) remain authoritative for gates |
| D3 | Collisions | `implements` / `blocks` documented as **distinct** from `change_implements_decision` / `uncertainty_blocks_task` |

### E. INTAKE + non-goals

| ID | Check | PASS when |
|----|-------|-----------|
| E1 | Auth FE/BE | MVP set can express `scope_member` clusters + `api_contract` cross edge |
| E2 | Business | `kind=business` (or equivalent) allowed under D1; not flattened solely to first goal |
| E3 | DR-NOSSEM | No embedding auto-cluster as P44 path |
| E4 | D4 honesty | No silent index-hook / daemon; CLI opt-in only |
| E5 | Portable | Seed-export allowlist expansion + `trace/graph.json` called out for S02 |

### F. Provenance honesty (Law 5 adjacent / D3)

| ID | Check | PASS when |
|----|-------|-----------|
| F1 | Store | Inference uses `source_type=INFERRED`; explicit wins; never deletes explicit |
| F2 | Wire | Optional `edge.provenance` maps from `source_type`; UI must not rename stored `rel` for inference |

---

## Findings format

```text
## Findings
- severity: blocker | high | medium | low | nit
- id: …
- evidence: …
- action: inline fix | spawn P44-S01-02a/02b | accept residual
```

| Severity | Action |
|----------|--------|
| blocker / high | Inline fix if trivial **or** spawn 02a/02b with full prompts; **do not** start S02 |
| medium | Prefer spawn unless one-paragraph doc fix |
| low / nit | Notes OK; do not block APPROVE |

## Verdict

```text
## Verdict
- result: APPROVE | REJECT
- confidence: high | medium | low
- rubric: L_/A_/B_/C_/D_/E_/F_ (list FAIL ids if any)
- residuals: …
- next: P44-S02-00 | spawned row id
```

**APPROVE** ⇒ all rubric rows PASS (or only low/nit residuals).  
**REJECT** ⇒ any FAIL on L/A/B/C/D/E/F without completed remediation spawn.

## Exit criteria

- [ ] Verdict APPROVE or spawn remediation with confidence medium/high
- [ ] Evidence: status LOCKED, D1–D5 closed, Laws 6–7/13/19 + M-001 cited
- [ ] Board Notes; next **P44-S02-00** (only on APPROVE)

## Minimal todos

- [ ] Diff locks vs light-lock + RESEARCH §7
- [ ] Score rubric L–F
- [ ] APPROVE / spawn + board Notes

## Next

`P44-S02-00`
