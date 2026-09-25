# Phase 44 / S01 / Lock design

## Metadata
- id: P44-S01-01
- todo_ids: [P44-S01-01]
- role: implementer
- skills: [documentation-and-adrs, domain-modeling]
- verification: automated

## Objective

Finalize [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — scope model, rel taxonomy, cap policy, non-goals, portable-graph wording. Set document status to **LOCKED (S01)**. Close D1–D5 with no TBD. Apply RESEARCH.md §7 wording refinements only. **No product code.** Do **not** silently reverse D1–D5 (reversal → REVIEW spawn).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`../scope-00-intake-research/RESEARCH.md`](../scope-00-intake-research/RESEARCH.md) — **APPROVE** (P44-S00-02); inherit §4–§7
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — light-locked D1–D5 (P44-00)
- [`INTAKE.md`](../../INTAKE.md)
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol Session start. Unattended authority: light locks + RESEARCH APPROVE. If a wording change would flip a D1–D5 choice, **stop** and note spawn — do not edit the lock row to the opposite option.

## Locked defaults (inherit — do not flip)

| Item | Value |
|------|-------|
| Target | `docs/phases/phase-44-scoped-graph-linking/01-DESIGN-LOCKS.md` |
| Status line | Change from “Light-locked pending S01 APPROVE …” → **`LOCKED (S01)`** + date `2026-08-23` (or session date) |
| D1 | **A + thin scope records**; reject C as primary |
| D2 | MVP: `scope_member`, `api_contract`, `implements`, `blocks` |
| D3 | `scope` query param + bounded N-hop; optional `scope_id` / `edge.provenance`; no tiles |
| D4 | Inference **CLI only** |
| D5 | GUI **500** stays; prefer filter over raise |
| Optional rels | `same_feature_front_back`, `depends_on`, `relates_to` — document optional; do not block S02 MVP |
| Product code | **Forbidden** |

## Preflight

1. Confirm RESEARCH.md APPROVE (board P44-S00-02 Notes).
2. Re-read RESEARCH §5 (D1–D5 gaps) and §7 (open risks — wording only).
3. Confirm current `01-DESIGN-LOCKS.md` status is still light-locked (not already LOCKED).

---

## Exact doc edits (apply in order)

### E0 — Status banner (top of `01-DESIGN-LOCKS.md`)

Replace the opening status paragraph with:

```markdown
**Status:** LOCKED (S01) — 2026-08-23. Implementers must not contradict locked rows without a forward REVIEW spawn. Wording refinements closed via P44-S01-01; D1–D5 choices unchanged from P44-00 light locks.
```

### E1 — D1 wording (plan_scopes ≠ graph scopes) — §2 + §9

In **§2 D1** and **§9 D1 row**, ensure these sentences appear (merge into existing bullets; do not flip A+thin):

- Persist lightweight **`scope` records**: `id`, `slug`, `title`, `kind` (`feature` \| `layer` \| `business`).
- Membership via **`scope_member`** in existing `entity_links` (Law 13 — no parallel relationship store).
- **`plan_scopes` are planner hierarchy rows**, not feature/layer/business cartography. “Reuse/map where titles align” means **soft title/body alignment for S03 inference** — **not** a hard FK, and **not** “every `plan_scope` becomes a graph scope node” in S02 MVP.
- Reject **C** (tag-column primary). Pure A (synthetic-only, no record table) only if seed-export cannot host thin records — default remains **A + thin records** (RESEARCH did not prove export blocker).

### E2 — D2 name collisions — §3.1 / §3.3 / §9

Add an explicit **collision** note (new short subsection under §3 or bold callout in D2 lock):

| Proposed D2 string | Must stay distinct from (existing) | Why |
|--------------------|------------------------------------|-----|
| `implements` | `change_implements_decision` | Task/decision peer edge ≠ change→decision side-effect |
| `blocks` | `uncertainty_blocks_task` | Task↔task scope edge ≠ uncertainty gate |

Require: domain consts, MCP/CLI aliases, seed allowlists, and docs use the **exact** D2 strings; do not alias D2 to the longer causal names.

**MVP set unchanged:** `scope_member`, `api_contract`, `implements`, `blocks`.

Optional (same pass if cheap): `same_feature_front_back`, `depends_on`, `relates_to` — presentation tokens in GUI may already list some of these; that is **not** evidence they persist at HEAD (see E5).

### E3 — D3 provenance + filter — §5.3 / D3 lock

Confirm and leave as locked intent:

- `GET /v1/graph?mode=project&scope=<slug|id>&max_nodes=…` + bounded N-hop (default depth 1–2), hard-capped by `max_nodes`; reject missing `max_nodes`; **no tiles**.
- Optional `scope_id` on `GraphNode` (from `scope_member`).
- Optional `edge.provenance` on wire: `explicit` \| `inferred`, **mapped from** `entity_links.source_type` (`USER_ASSERTED`/`IMPORTED` → explicit; `INFERRED` → inferred).
- **Do not** require UI to mutate stored `rel` strings with an `inferred_` prefix; dashed styling / provenance field only (align §3.4).

### E4 — Portable graph / seed-export — §8

Thicken §8 so S02 cannot miss CONTRIBUTING policy:

- Expanding D2 MVP rels **and** thin `scope` records requires: migration (if new table/entity) + **`seedExportLinkRels` / export allowlist** update + `trace seed export -o trace/graph.json` before PRs that change Trace entities.
- Today export allowlist is **four causal rels only** (`seed_export.go` ~341–346) — S02 must add MVP scope rels explicitly.
- INFERRED links export with honest `source_type`; re-import idempotent.
- Do **not** reopen D1–D5 while stating export expectations.

### E5 — Aspirational GUI tokens — §6

Add one bullet under §6 UI consequences:

- `EDGE_PRIORITY_RELS` may already list `blocks` / `depends_on` / `relates_to` — treat as **presentation prep**, not proof that scope semantics exist until S02 persists MVP rels. S04 elevates real MVP rels after they exist.

### E6 — D4 / D5 — no change of choice

- **D4:** Opt-in CLI only (`trace graph infer` or documented equivalent); no post-`trace index` hook; no silent daemon.
- **D5:** Keep GUI `PROJECT_MAX_NODES = 500`; API hard max 5000 unchanged; prefer scope filter; edge overview 150 stays presentation policy.

### E7 — §9 table → closed (not “open”)

Retitle §9 from “Open decisions — light-locked for S01” to **“Closed decisions (LOCKED)”**. Each row: final choice + one-sentence rationale (keep existing light-lock text; add collision/plan_scopes clarifiers only where needed). Remove any “pending S01” language.

### E8 — §10 Acceptance map

Keep V-style checks; ensure they still match VERIFY intent:

- [ ] Seeded auth FE/BE scopes → ≥2 clusters + ≥1 cross-scope edge (GUI smoke)
- [ ] `trace_link` / CLI can create `scope_member` + `api_contract`
- [ ] Inferred edges visually distinct; explicit wins
- [ ] Laws 6–7: no unbounded graph; truncation honest
- [ ] Law 19: walk/filter in `internal/retrieval`; GUI adapter only
- [ ] Law 13 / M-001: membership+MVP rels in `entity_links`; goal/task loop preserved

---

## Acceptance checklist (locker self-check)

Mark all before board `done`:

| ID | Check |
|----|-------|
| A1 | Status line is **LOCKED (S01)** with date |
| A2 | D1 final = A + thin records; C rejected; plan_scopes soft-map wording present |
| A3 | D2 MVP four rels named; collision table/note for `implements`/`blocks` present |
| A4 | D3 `scope` filter + optional `scope_id` / `edge.provenance` mapping stated; no tiles |
| A5 | D4 CLI-only; D5 GUI 500 unchanged |
| A6 | §8 seed-export allowlist + portable `graph.json` expectations explicit |
| A7 | §6 notes aspirational GUI tokens vs persisted semantics |
| A8 | §9 titled closed/LOCKED; no TBD on D1–D5 |
| A9 | §10 still aligned with VERIFY / INTAKE clusters |
| A10 | Diff is docs-only; no silent D1–D5 reversal |

## File targets (write)

| Path | Action |
|------|--------|
| `docs/phases/phase-44-scoped-graph-linking/01-DESIGN-LOCKS.md` | Apply E0–E8 |
| `SCOPE-TODOS.md` | Mark S01-1 / S01-D done when locks written |

Optional (only if still thin after P44-S01-00): leave S02–S04 deep thicken to those planners.

## Exit criteria

- [ ] `01-DESIGN-LOCKS.md` status **LOCKED (S01)**; D1–D5 closed; A1–A10 PASS
- [ ] Acceptance map §10 still aligned with VERIFY
- [ ] Board `done` with Notes listing edits applied; next **P44-S01-02**

## Minimal todos

- [ ] Apply E0–E8 to 01-DESIGN-LOCKS
- [ ] Self-check A1–A10 vs RESEARCH + INTAKE
- [ ] SCOPE-TODOS + board Notes

## Next

`P44-S01-02`
