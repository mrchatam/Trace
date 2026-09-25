# Phase 44 / S03 / Implement inference

## Metadata
- id: P44-S03-01
- todo_ids: [P44-S03-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: [user-trace]
- verification: automated

## Objective

Ship an **opt-in CLI** inference pass that proposes `scope_member` (and only soft optional peers if cheap) with `source_type=INFERRED`, lower confidence than explicit, **fail-closed** on conflicts, and **never** deletes or overwrites `USER_ASSERTED` / `IMPORTED`. **Prerequisite:** S02 APPROVE (P44-S02-02) + S01 locks APPROVE. **No** GUI (`web/`). **No** MCP write tool for inference. **No** index-hook / daemon.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §3.4, §4, **D4** — LOCKED
- [`../scope-00-intake-research/RESEARCH.md`](../scope-00-intake-research/RESEARCH.md) §5 D4, §7 risks 2 & 5
- [`../scope-02-backend-links/01-implement.md`](../scope-02-backend-links/01-implement.md) — scopes + MVP rels already landed
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- Live: `internal/store/links.go` (`InsertLinkOrIgnore`, `source_type`, `confidence`); `internal/store/plan_hierarchy.go` (`plan_scopes`); `internal/domain/scope.go` (membership APIs); `internal/retrieval/graph_neighbors.go` (`provenanceFromSourceType`)

## Session start

Follow agent-loop-protocol Session start. Block if P44-S02-02 is not APPROVE.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Trigger | **CLI only:** `trace graph infer` (preferred) or documented alias under `cmd/trace/` — **must** appear in `trace help` |
| Out of band | **No** call from `cmd/trace/index.go` / watch / post-index promote; **no** silent daemon; **no** HTTP/MCP auto-infer endpoint |
| Provenance | Every written edge: `source_type=INFERRED` |
| Confidence | Documented constant **`InferredLinkConfidence = 0.4`** (must be **lower than** typical explicit `0.9`/`1.0`) |
| Insert path | `store.InsertLinkOrIgnore` (or domain wrapper that uses it) — `ON CONFLICT … DO NOTHING` |
| Explicit wins | Existing row with same endpoints (`from_type,from_id,rel,to_type,to_id`) is **never** replaced; never `DELETE` / never bump confidence of `USER_ASSERTED`/`IMPORTED` |
| Fail-closed | Ambiguous / conflicting rule outputs → **skip** (emit neither), log/report reason; do not pick arbitrarily |
| Rels in scope | MVP: **`scope_member` only** for rules below. Optional `relates_to` between co-members **only if cheap** and still `INFERRED` — do **not** invent `api_contract`/`implements`/`blocks`/`review_judges_task` from rules |
| plan_scopes | Soft title/body **token alignment** → candidate membership to an **existing** graph `scopes` row — **not** hard FK; **not** auto-create a graph scope per plan_scope |
| Wire | Do **not** mutate stored `rel` with `inferred_` prefix; retrieval already maps `INFERRED` → `edge.provenance=inferred` |
| Caps | Inference writes links only; any graph walk still honors `max_nodes` (do not add unbounded dump CLI) |
| DR-NOSSEM | No embeddings / vector similarity |
| Portable | INFERRED links export with honest `source_type` (S02 allowlist already includes MVP rels); re-run infer is idempotent |

### Preflight anchors @ HEAD (re-verify if drift)

| Claim | Location |
|-------|----------|
| `InsertLinkOrIgnore` + endpoint UNIQUE | `internal/store/links.go:23–54` |
| Thin scopes + `scope_member` APIs | `internal/domain/scope.go`, migration `029_graph_scopes.sql` |
| `plan_scopes` planner rows (soft-map only) | `internal/store/plan_hierarchy.go` |
| Wire provenance map | `internal/retrieval/graph_neighbors.go:9–19` |
| No `trace graph` CLI yet | `cmd/trace/` — add `graph` parent + `infer` subcommand |
| Index must stay hook-free for infer | `cmd/trace/index.go` (post-index only watermark / VCS promote today) |
| S02 must not create INFERRED | `TestNoINFERREDFromScopePaths` — keep; S03 adds **CLI-only** writers |

---

## Rule table (MVP)

Implement these three rules as a deterministic, ordered pass. Each candidate is a proposed `entity_links` row with `rel=scope_member`, `to_type=scope`, `to_id=<scopes.id>`, `source_type=INFERRED`, `confidence=InferredLinkConfidence`.

| ID | Rule | Input evidence | Match → target scope | Output |
|----|------|----------------|----------------------|--------|
| **R-PATH** | Path → layer/feature | Repo-relative paths tied to an entity (prefer: indexed `files` paths linked via changes / evidence / body path hints; document exact signal chosen in Notes). Globs: `web/**`, `web/src/**` → frontend layer; `internal/**`, `cmd/**` → backend layer. Optional feature refine: path contains `auth` / `billing` **and** a matching `scopes.slug` exists (e.g. `auth-fe`, `auth-be`) | Resolve to an **existing** `scopes` row (`kind=layer` or `feature`). **Do not** create scopes in infer MVP unless a tiny helper is already used by tests — prefer seed/create scopes in test fixtures | `scope_member` entity→scope |
| **R-TITLE** | Title tokens | Entity `title` (tasks first; decisions/discoveries if cheap) case-folded | Token map (locked starters): `auth`/`login`/`session` → prefer feature scopes whose slug/title contains `auth`; `billing`/`payment` → `billing`; `gui`/`overview`/`frontend` → FE layer/feature if present. Require **existing** scope match; no invent | `scope_member` |
| **R-PLAN** | plan_scope co-membership | Active `plan_scopes` rows (`ListPlanScopesByGoal` / by phase). Soft-align: normalize plan_scope `title` (+ optional `body` first line) against graph `scopes.title`/`slug` (substring / token overlap — **no** embedding) | If alignment score is unique → candidate `scope_member` for entities already tied to that plan subtree (tasks under goal/phase). If **zero or ≥2** equally good graph scopes → **skip** (fail-closed) | `scope_member` |

### Rule ordering & conflict policy (fail-closed)

1. Collect candidates from R-PATH, R-TITLE, R-PLAN independently (same entity may get multiple candidates).
2. **Per entity, per `rel=scope_member`:**
   - If an explicit (`USER_ASSERTED` / `IMPORTED`) `scope_member` already exists for **any** scope → **do not add** further `INFERRED` membership for that entity in this MVP (explicit wins; avoid dual-membership hairball). Document in Notes if you choose the narrower rule “only skip when endpoints collide” — default for S03 is **skip all INFERRED membership when any explicit membership exists**.
   - If multiple INFERRED candidates point to **different** scopes → **emit none** for that entity (fail-closed); report `conflict` in dry-run / summary.
   - If multiple rules agree on the **same** scope → one insert (idempotent).
3. Never delete or update existing links. `InsertLinkOrIgnore` no-ops when endpoints already exist (including prior INFERRED).
4. **Hard forbid:** inferring gate/causal rels (`review_judges_task`, `decision_affects_task`, `goal_has_task`, `uncertainty_blocks_task`, `change_implements_decision`, …).

### Out of MVP (do not implement in S03)

| Item | Why |
|------|-----|
| OpenAPI tag → `api_contract` | Named in locks as future hint; not in planner rule trio |
| Auto-create thin scopes from plan_scopes | Violates D1 soft-map / RESEARCH risk 2 |
| LLM-only fill | Rejected in §4 |
| Index / watch hook | D4 |

---

## CLI surface

```text
trace graph infer [--dry-run] [-C <root>]
```

| Flag / behavior | Spec |
|-----------------|------|
| Default | Run rules; insert via `InsertLinkOrIgnore`; print summary counts: `inserted`, `skipped_existing`, `skipped_conflict`, `skipped_explicit` |
| `--dry-run` | Compute candidates + conflict skips; **zero** writes; print table or JSON lines of proposed edges |
| Exit | `0` on success (including “nothing to do”); non-zero only on store/validation errors |
| Help | Document under `trace help` that inference is opt-in and does not run on `trace index` |

Library entrypoint (suggested): `domain.InferScopes(ctx, InferOptions{DryRun bool}) (InferReport, error)` called only from CLI — keeps Law 19 (no logic in `web/`).

---

## Idempotency & dry-run (cheap — required)

| Property | How |
|----------|-----|
| Idempotent re-run | Second `trace graph infer` inserts `0` new rows when graph unchanged (`InsertLinkOrIgnore`) |
| Dry-run | `--dry-run` shares the same candidate pipeline; skip `InsertLinkOrIgnore` |
| Export honesty | INFERRED rows round-trip with `source_type=INFERRED` (no special-case wipe on import) |

---

## File touch-list (ordered)

### Slice A — Library + constants
- Prefer `internal/domain/` (e.g. `infer_scopes.go`) or small `internal/infer/` imported by domain/CLI — **not** `web/`, **not** MCP tool.
- Consts: `SourceTypeInferred = "INFERRED"`, `InferredLinkConfidence = 0.4`.
- Rule runners + conflict resolver + `InferReport`.

### Slice B — CLI
- `cmd/trace/graph.go` (or `graph_infer.go`): parent `graph` + `infer`.
- Wire into root command dispatch + help text.
- **Grep proof:** no reference from `index.go` / `index_watch.go`.

### Slice C — Tests (minimum)
| Test | Assert |
|------|--------|
| Path or title rule happy path | ≥1 `scope_member` with `source_type=INFERRED`, confidence `0.4` |
| Explicit wins | Pre-seed `USER_ASSERTED` membership → infer inserts **0** for that entity (or no overwrite of source_type/confidence) |
| Conflict fail-closed | Two rules → two different scopes → **0** inserts for that entity |
| Idempotent | Run twice → second `inserted=0` |
| Dry-run | `--dry-run` (or API DryRun) leaves link count unchanged |
| No gate rels | Infer never writes `review_judges_task` / causal MVP-forbidden set |
| No index hook | Test or `rg` in CI-style assert: `infer` not called from index package path |

### Slice D — Docs
- Help string + brief CONTRIBUTING / command note if export behavior needs a one-liner (INFERRED already honest).

---

## Exit criteria

- [ ] `trace graph infer` (or locked documented equiv.) works; documented in help
- [ ] Tests: ≥1 rule path; explicit wins; conflict fail-closed; idempotent and/or dry-run
- [ ] Grep/tests prove **no** index-hook / daemon registration
- [ ] No `web/` edits; no new MCP write tool; no gate-rel inference
- [ ] Board `done` with Notes; next **P44-S03-02**

## Minimal todos

- [ ] Slice A — infer library + rule table + fail-closed
- [ ] Slice B — CLI `graph infer` + `--dry-run`
- [ ] Slice C — tests (rule + explicit-wins + conflict + idempotency/dry-run)
- [ ] Slice D — help / Notes
- [ ] Board Notes → next **P44-S03-02**

## Next

`P44-S03-02`
