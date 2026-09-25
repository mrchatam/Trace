# Phase 44 / S02 / Review — Backend links

## Metadata
- id: P44-S02-02
- todo_ids: [P44-S02-02]
- role: reviewer
- skills: [code-review, security-and-hardening]
- verification: automated

## Objective

Independent review of S02 backend + **portable export round-trip** against **LOCKED** [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) (D1–D3, §3, §8). Spawn forward if Law 6 / 13 / 19 or M-001 regressions. Fresh subagent — do not share implementer session.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-implement.md`](01-implement.md) — acceptance A1–A6
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — LOCKED (S01); P44-S01-02 APPROVE
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- CONTRIBUTING — `trace seed export -o trace/graph.json`

## Session start

Fresh subagent. Follow agent-loop-protocol. Re-read board Notes on **P44-S02-01** before judging.

## Review rubric

### R1 — Domain & Law 13

| Check | Pass |
|-------|------|
| MVP rels exact: `scope_member`, `api_contract`, `implements`, `blocks` | Yes |
| Consts/aliases **not** colliding with `change_implements_decision` / `uncertainty_blocks_task` | Yes |
| Membership + MVP edges only in `entity_links` (no parallel rel store) | Yes |
| Thin `scopes` table (or equivalent) — **not** hijacking `plan_scopes` as graph scopes | Yes |
| `plan_scopes` not hard-FK’d / not auto-injected as graph scope nodes | Yes |

### R2 — Law 6 / D3 API

| Check | Pass |
|-------|------|
| `max_nodes` still required; hard max ≤5000 unchanged | Yes |
| `scope` filter bounded by `max_nodes` (N-hop capped); no tiles | Yes |
| Optional `scope_id` / `edge.provenance` match locks (provenance from `source_type`, no `inferred_` rel mutation) | Yes |
| GUI 500 / edge 150 **not** raised in this scope (S02 should not touch GUI) | Yes |

### R3 — Law 19 / MCP

| Check | Pass |
|-------|------|
| Walk/filter logic in `internal/retrieval` (+ domain/store); HTTP/MCP thin | Yes |
| No business logic fork in `web/` from S02 | Yes |
| `trace_link` allowlist extended — **no** new write tool | Yes |
| CLI mirror matches MCP aliases | Yes |

### R4 — Tests & M-001

| Check | Pass |
|-------|------|
| A1–A4 from `01-implement.md` evidenced (tests or Notes with commands) | Yes |
| Causal goal/task / review gate paths not broken (spot-check) | Yes |
| No S02-authored `INFERRED` writers | Yes |

---

## Export round-trip checklist (mandatory)

Run or verify evidence for each item. Prefer a **temp DB** / test harness over mutating the developer’s live `.trace/` without cleanup.

### E0 — Preconditions

- [ ] Schema migration **029** (or shipped equivalent) present under `internal/store/schema/` and applied on open.
- [ ] `seedExportLinkRels` (or successor allowlist) includes **all four** MVP rels by exact string.
- [ ] `SeedDocument` exports thin scopes (dedicated array — not only `plan_scopes`).

### E1 — Seed fixture (minimal)

Create (via domain/CLI/test) at least:

| Artifact | Intent |
|----------|--------|
| Scope A | e.g. slug `auth-fe`, kind `feature` |
| Scope B | e.g. slug `auth-be`, kind `feature` |
| Tasks T1, T2 | Members of A / B |
| Link `scope_member` | T1→A, T2→B |
| Link `api_contract` | T1→T2 (cross-scope) |
| Optional | one `implements` and/or `blocks` if cheap |

### E2 — Export

- [ ] `trace seed export` (or domain `Export`) includes both scopes in scopes array.
- [ ] Exported `links` contain MVP rels with correct from/to ids (or from_id/to_id).
- [ ] Causal-only allowlist regression: existing four causal export rels still present if used.
- [ ] If PR changed Trace entities: committed `trace/graph.json` refreshed (`trace seed export -o trace/graph.json`).

### E3 — Import into empty store

- [ ] Fresh store + import restores scope rows (id/slug/title/kind).
- [ ] Restored `entity_links` for MVP rels; endpoint types resolve (`to_type=scope` for membership).
- [ ] Duplicate import is **idempotent** (no multiply; no hard error on same endpoints) — align DF-81 style.

### E4 — Graph honesty after import

- [ ] `ProjectGraph` (sufficient `max_nodes`) shows MVP edges between restored nodes.
- [ ] `scope=<slug|id>` filter returns members (+ bounded hops) without dumping whole project.
- [ ] `source_type` on imported links is honest (`IMPORTED` or documented import default) — not silently `INFERRED`.

### E5 — Negative / hygiene

- [ ] Rel named `implements` did **not** write `change_implements_decision`.
- [ ] Rel named `blocks` did **not** write `uncertainty_blocks_task`.
- [ ] Optional rels absent from MVP are OK if not shipped; if shipped, same export allowlist coverage.
- [ ] No secrets / `.trace/trace.db` committed; only portable `trace/graph.json` as required.

---

## Verdict rules

| Outcome | When |
|---------|------|
| **APPROVE** | R1–R4 pass; E0–E5 pass or Notes cite equivalent automated tests with file:line |
| **Spawn** | Law 6/13/19 fail, missing export allowlist, plan_scopes hijack, collision alias, or missing A1–A3 evidence — insert `P44-S02-02a` / `02b` immediately below this row |

## Exit criteria

- [ ] APPROVE or spawn with reason
- [ ] Board Notes with evidence cites
- [ ] Next **P44-S03-00** (on APPROVE)

## Minimal todos

- [ ] Spot-check HEAD vs P44-S02-01 Notes
- [ ] Execute / verify export round-trip checklist E0–E5
- [ ] Rubric R1–R4
- [ ] APPROVE / spawn
- [ ] Board Notes

## Next

`P44-S03-00`
