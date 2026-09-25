# Phase 44 / S03 / Review — Inference

## Metadata
- id: P44-S03-02
- todo_ids: [P44-S03-02]
- role: reviewer
- skills: [code-review]
- verification: automated

## Objective

Independent review of S03 inference for **honesty** (INFERRED provenance, explicit wins, fail-closed conflicts), **D4 CLI-only trigger**, and **caps / non-creep** (no daemon, no index-hook, no gate-rel inference, no embeddings). Fresh subagent — do not share implementer session. Spawn forward on blocker/high.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-implement.md`](01-implement.md) — rule table R-PATH / R-TITLE / R-PLAN; acceptance tests
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §3.4, §4, **D4**, §7 non-goals
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- Board Notes on **P44-S03-01**

## Session start

Fresh subagent. Follow agent-loop-protocol. Re-read P44-S03-01 Notes before judging. Re-verify S02 APPROVE still stands (inference must not regress scope CRUD / export).

## Review rubric

### H1 — Honesty & provenance

| Check | Pass |
|-------|------|
| Every infer-written edge has `source_type=INFERRED` | Yes |
| Confidence is the documented low constant (≤0.4 or implement Notes const) and **below** explicit defaults | Yes |
| Stored `rel` strings are **not** prefixed/mutated with `inferred_` | Yes |
| Wire may expose `edge.provenance=inferred` via existing map — no parallel provenance store | Yes |
| Seed export/import leaves `INFERRED` honest (no silent upgrade to `USER_ASSERTED`) | Yes |

### H2 — Explicit wins & fail-closed

| Check | Pass |
|-------|------|
| Existing `USER_ASSERTED` / `IMPORTED` links never deleted or overwritten by infer | Yes |
| Inserts use `InsertLinkOrIgnore` (or equivalent DO NOTHING on endpoint UNIQUE) | Yes |
| Conflicting multi-scope candidates → skip (fail-closed), not arbitrary pick | Yes |
| Tests cover: ≥1 rule happy path **and** explicit-wins **and** conflict skip | Yes |

### C1 — D4 trigger (CLI-only) — **blocker if fail**

| Check | Pass |
|-------|------|
| User-facing command is opt-in CLI (`trace graph infer` or documented equiv. in help) | Yes |
| **No** invocation from `trace index` / index watch / post-index hooks | Yes |
| **No** always-on daemon / background infer loop | Yes |
| **No** new MCP/HTTP write tool that auto-runs inference on link/add/index | Yes |
| Grep or test evidence cited in Notes for “no index-hook” | Yes |

### C2 — Caps & scope discipline

| Check | Pass |
|-------|------|
| Inference does not bypass graph `max_nodes` / raise GUI 500 / API 5000 | Yes |
| No unbounded “dump all inferred edges” API as default | Yes |
| `--dry-run` (or equivalent) performs **zero** writes when claimed | Yes |
| Re-run is idempotent (`inserted=0` when unchanged) | Yes |

### C3 — Non-goals / creep

| Check | Pass |
|-------|------|
| No gate-rel auto-infer (`review_judges_task`, etc.) | Yes |
| No embeddings / vector scope fill (DR-NOSSEM) | Yes |
| No hard FK `plan_scopes` → graph scopes; no “every plan_scope becomes a scope node” | Yes |
| MVP rules stay path / title / plan_scope soft-align; no silent OpenAPI/`api_contract` mass invent unless Notes explicitly scoped it out | Yes |
| No `web/` business logic for inference (Law 19) | Yes |
| M-001: goal/task loop / review gate paths untouched | Yes |

---

## Evidence checklist (reviewer must collect)

### E1 — Command & help
- [ ] `trace help` (or `--help` on graph) documents opt-in infer and that index does not run it.
- [ ] Running infer on a fixture DB produces expected INFERRED `scope_member` rows **or** dry-run lists them.

### E2 — Tests
- [ ] Re-run named tests from P44-S03-01 Notes (or `go test` packages listed).
- [ ] Confirm explicit-wins + conflict + idempotency/dry-run covered.

### E3 — Anti-creep spot-check
- [ ] `rg -n 'InferScope|graph infer|INFERRED' cmd/trace/index.go cmd/trace/index_watch.go` → no infer registration (or equivalent proof).
- [ ] `rg` MCP tools: no new infer write tool.
- [ ] No new embedding imports in infer package.

### E4 — Regression
- [ ] Spot-check S02: scopes CRUD / MVP rels / `TestNoINFERREDFromScopePaths` still meaningful (S03 CLI path is the **only** INFERRED writer intended).

---

## Verdict protocol

| Outcome | When |
|---------|------|
| **APPROVE** | H1–H2, C1–C3 all PASS; E1–E4 evidenced; confidence medium/high; residuals listed |
| **Spawn** | blocker/high (esp. C1 index-hook/daemon, explicit overwrite, gate-rel infer) → insert `P44-S03-02a` / `02b` immediately below this row |
| **Reject-as-notes** | Do not silent-fix product architecture in review beyond small typos; prefer spawn |

## Exit criteria

- [ ] APPROVE or spawn with board rows
- [ ] SCOPE-TODOS S03-2 / S03-C updated
- [ ] Next **P44-S04-00** (only after APPROVE)

## Minimal todos

- [ ] Evidence review against rubric H/C + E1–E4
- [ ] Board Notes (APPROVE/spawn + confidence)
- [ ] Next **P44-S04-00** if APPROVE

## Next

`P44-S04-00`
