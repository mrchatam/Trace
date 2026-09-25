# Phase 44 / S05 / VERIFY

## Metadata
- id: P44-S05-01
- todo_ids: [P44-S05-01]
- role: verify
- skills: [grinding-until-pass]
- mcps: [user-trace]
- verification: mixed

## Objective

Run VERIFY blocks **V1–V7** after S00–S04 reviews **APPROVE**. Aggregate prior PASS cites + live re-checks into **[`VERIFY-NOTES.md`](VERIFY-NOTES.md)** (+ evidence dir). Map to [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) §10 acceptance. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P44-S05-02**. **No product code.** Do **not** start S05-02 or invent a successor phase.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S05-00)
- [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) — D1–D5 + §10 acceptance map
- [`DR-HANDOFF.md`](../../DR-HANDOFF.md) — remains **OPEN** until S05-02
- [`SCOPE-TODOS.md`](SCOPE-TODOS.md)
- Prior PASS artifacts (cite in notes):
  - S00 [`RESEARCH.md`](../scope-00-intake-research/RESEARCH.md) + [`02-review`](../scope-00-intake-research/02-review.md) APPROVE
  - S01 [`01-DESIGN-LOCKS.md`](../../01-DESIGN-LOCKS.md) LOCKED + [`02-review`](../scope-01-design-locks/02-review.md) APPROVE
  - S02 board Notes (P44-S02-01/02) — migration 029, MVP rels, export round-trip
  - S03 board Notes (P44-S03-01/02/02a/02b) — CLI infer, H1 seed honesty, C1 CLI-only
  - S04 board Notes (P44-S04-01/02) — scope clustering, dashed inferred, browser smoke
- Live anchors:
  - Backend: `internal/retrieval/project_graph.go`, `internal/domain/scope.go`, `internal/domain/infer_scopes.go`, `internal/domain/seed_export.go`, `internal/mcp/tools_write.go`, `internal/store/migrations/029_graph_scopes.sql`
  - CLI: `cmd/trace/graph.go`, `cmd/trace/link.go`, `cmd/trace/help.go`
  - GUI (Law 19 adapter): `web/src/lib/graphLayout.ts`, `web/src/lib/overviewCompose.ts`, `web/src/pages/Graph.tsx`, `web/src/components/GraphOrientPanel.tsx`
  - Fixture: `web/testdata/p44-scope-smoke-seed.json`
  - Portable: `trace/graph.json`, [CONTRIBUTING.md](../../../../../CONTRIBUTING.md) Portable graph

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or change product bodies.

## Locked defaults (FINAL — S05-00)

| Item | Value |
|------|-------|
| Precondition | P44-00 … P44-S04-02 all `done`; S00–S04 reviews **APPROVE** (incl. S03-02b H1 closed) |
| Product / Go / TS / docs changes | **Forbidden** (evidence + notes only). Failures → spawn remediation from this row or leave FAIL for S05-02 to spawn |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p44-s05-01-verify/evidence/` |
| Notes artifact | `scopes/scope-05-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S05-02 closes |
| Successor | **Out of scope** — S05-02 only (lean default **`no successor`**) |
| Auth FE/BE fixture | `web/testdata/p44-scope-smoke-seed.json` — scopes `auth-fe` / `auth-be`, `scope_member` + `api_contract` |
| Caps (locked D5) | GUI default **500** nodes (`PROJECT_MAX_NODES`); overview **150** edges (`EDGE_OVERVIEW_MAX`); API hard cap **5000** (`MaxNeighborhoodNodes`) |
| Inference (locked D4) | **CLI only** — `trace graph infer`; no index-hook / MCP / `web/` write |
| Portable graph | Schema changed in S02/S03 — **`trace seed export -o trace/graph.json`** required evidence + import smoke |

### §10 acceptance map (must tick in VERIFY-NOTES)

| Block | §10 criterion | Primary evidence |
|-------|---------------|------------------|
| **V1** | Seeded auth FE/BE → ≥2 clusters + ≥1 cross-scope edge | Browser smoke on smoke seed + optional API spot-check |
| **V2** | `trace_link` / CLI creates `scope_member` + `api_contract` | MCP aliases + `trace link` + domain tests |
| **V3** | Inferred visually distinct; explicit wins | CLI infer tests + GUI dash/provenance + seed round-trip |
| **V4** | Laws 6–7: no unbounded graph; truncation honest | Retrieval tests + GUI cap constants + grep |
| **V5** | Law 19: walk/filter in `internal/retrieval`; GUI adapter | `rg` + layout uses API `scope_id` only |
| **V6** | Law 13 / M-001: rels in `entity_links`; goal/task loop preserved | No gate-rel infer; loop/gate paths unchanged |
| **V7** | Portable graph export/import if schema changed | `trace/graph.json` + seed round-trip tests |

### Aggregate evidence map (S00–S04 — cite in VERIFY-NOTES)

| Scope | Must cite / re-check |
|-------|----------------------|
| S00 | RESEARCH: hairball root cause; INTAKE auth FE/BE; D1–D5 gap matrices; caps 500/5000 record-only |
| S01 | LOCKED D1–D5; MVP rel taxonomy; CLI-only infer; GUI 500; seed-export allowlist obligation |
| S02 | Migration 029 `scopes`; MVP rels in `entity_links`; retrieval `scope` filter; MCP/CLI link aliases; export round-trip |
| S03 | `trace graph infer`; `INFERRED` + confidence 0.4; explicit wins; H1 seed `source_type`/`confidence` round-trip |
| S04 | Cluster-by-`scope_id`; tiered edge priority; dashed `provenance=inferred`; browser smoke ≥2 clusters |

### Fail vs residual (locked)

**Fail VERIFY for:**

- Smoke seed fails to show **≥2** scope clusters or **≥1** `api_contract` cross-scope edge in GUI/API
- MCP/CLI cannot create `scope_member` + `api_contract` (aliases missing or domain reject)
- Inferred edges not visually distinct (`provenance=inferred` dash) **or** explicit membership does not win over infer
- Unbounded graph path (missing `max_nodes`, cap >5000 accepted, default dump >500 in GUI without truncation honesty)
- `web/` performs scope membership walk or graph-walk logic fork (Law 19 violation)
- `trace graph infer` registered from index/watch/MCP/`web/` (D4 / C1 regression)
- Scope infer or rels mutate task loop / gate enforcement (M-001 regression)
- MVP scope rels missing from `seedExportLinkRels` **or** INFERRED provenance lost on seed round-trip
- Focused backend/GUI tests listed in blocks below **FAIL**

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES):

| Residual | Disposition |
|----------|-------------|
| `trace/graph.json` empty scopes (`omitempty`) on live Trace DB | Accept if export command + round-trip tests PASS |
| Synthesized goal→task SeedLinks omit provenance | Pre-existing; not INFERRED path (S03-02b residual) |
| `plan_scopes` ↔ graph slug auto-mapping quality | Deferred (DR-HANDOFF queue) |
| Scope filter pagination/tiles at 500 cap | Deferred |
| Codegraph complement edges | Out of scope — dual-stack doc only |
| Cosmetic orient legend wording | Accept if a11y substance present (dash + text) |

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir. Use `date +%Y-%m-%d` for the run folder name.

### Block 0 — Evidence dir + preflight

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p44-s05-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P44-S05-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=P44-S04-02 APPROVE; S00–S04 done; DESIGN-LOCKS §10"
} > "$EVID/00-run-metadata.txt"
test -f docs/phases/phase-44-scoped-graph-linking/01-DESIGN-LOCKS.md
test -f web/testdata/p44-scope-smoke-seed.json
```

**Pass:** `$EVID` exists; metadata cites S04-02 APPROVE + S00–S04 complete; smoke seed present.

### Block V1 — Auth FE/BE fixture smoke (GUI + optional API)

**Fixture:** `web/testdata/p44-scope-smoke-seed.json` (`auth-fe`, `auth-be`, `scope_member`, `api_contract`).

```bash
# Build web if needed for serve smoke
npm run build --prefix web 2>&1 | tee "$EVID/v1-web-build.txt"

TMP=$(mktemp -d)
cleanup() { rm -rf "$TMP"; }
trap cleanup EXIT
go run ./cmd/trace -C "$TMP" init 2>&1 | tee "$EVID/v1-init.txt"
go run ./cmd/trace -C "$TMP" seed import web/testdata/p44-scope-smoke-seed.json 2>&1 | tee "$EVID/v1-seed-import.txt"

# API spot-check (Law 19 server walk — optional but recommended)
go run ./cmd/trace -C "$TMP" context --format json 2>/dev/null | head -5 | tee "$EVID/v1-context-head.txt" || true
# Prefer HTTP project graph if serve available:
# go run ./cmd/trace -C "$TMP" serve --addr 127.0.0.1:17997 --static-dir web/dist --no-open &
# curl -s "http://127.0.0.1:17997/v1/project/graph?max_nodes=500" | tee "$EVID/v1-api-graph.json"
# jq '{scopes: [.nodes[]|.scope_id]|unique, api_contract: [.edges[]|select(.rel=="api_contract")]|length}' "$EVID/v1-api-graph.json"
```

**Browser smoke (required in VERIFY-NOTES — cite S04-01 pattern):**

1. `trace serve --static-dir web/dist --addr 127.0.0.1:<free-port> --no-open` with seeded `$TMP` store (or documented equivalent).
2. Open Explore; record URL, port, seed steps.
3. DOM evidence: **≥2** distinct scope cluster positions (`auth-fe` / `auth-be` or `distinctScopedPositions=2`); **≥1** `api_contract` edge visible; orient legend region visible.

**Pass:** seed import exit 0; browser/API shows **≥2 clusters** + **≥1 cross-scope** edge. If browser env blocked, API jq check + S04-02 board cite acceptable with explicit **WAIVED** reason for live browser.

### Block V2 — Agent link (`scope_member` + `api_contract`)

```bash
go test ./internal/domain/ -run 'TestScopeMVPRelsAndCRUD' -count=1 \
  2>&1 | tee "$EVID/v2-domain-mvp-rels.txt"
go test ./internal/mcp/ -run 'Test.*[Ll]ink' -count=1 \
  2>&1 | tee "$EVID/v2-mcp-link.txt"
go run ./cmd/trace link --help 2>&1 | tee "$EVID/v2-link-help.txt"
{
  echo "=== MCP trace_link aliases (scope_member, api_contract) ==="
  grep -n 'scope_member\|api_contract\|scope-member\|api-contract' internal/mcp/tools_write.go | head -20
  echo "=== MVP rel constants ==="
  grep -n 'RelScopeMember\|RelAPIContract' internal/domain/*.go | head -20
} 2>&1 | tee "$EVID/v2-grep-aliases.txt"
```

**Pass:** domain MVP rel test exit 0; MCP/CLI expose `scope_member` + `api_contract` aliases; no parallel relationship store (Law 13 — membership in `entity_links` only).

### Block V3 — Inference CLI-only + explicit wins + visual distinction

```bash
go test ./internal/domain/ -run 'TestInferScopes_' -count=1 \
  2>&1 | tee "$EVID/v3-infer-domain.txt"
go test ./cmd/trace/ -run 'TestGraphInferCLI_' -count=1 \
  2>&1 | tee "$EVID/v3-infer-cli.txt"
go test ./internal/domain/ -run 'TestINFERREDScopeMemberSeedRoundTrip' -count=1 \
  2>&1 | tee "$EVID/v3-inferred-seed.txt"
go run ./cmd/trace graph infer --help 2>&1 | tee "$EVID/v3-infer-help.txt"
{
  echo "=== D4 CLI-only: infer must NOT appear in index/watch/MCP/web ==="
  rg -n 'InferScopes|graph infer' internal/index internal/watch internal/mcp web/src || echo "PASS: no infer hooks outside cmd/domain"
} 2>&1 | tee "$EVID/v3-cli-only-rg.txt"
cd web && node --experimental-strip-types --test src/lib/graphLayout.test.ts 2>&1 | tee "../$EVID/v3-graphlayout-inferred.txt"
```

**Pass checklist:**

1. `TestInferScopes_ExplicitWins` PASS — explicit `scope_member` blocks INFERRED for same entity.
2. `TestInferScopes_NoIndexHookRegistration` PASS — no index-hook registration.
3. `trace graph infer` in CLI help; **no** infer write path in MCP/`web/`.
4. GUI: dashed inferred via `provenance=inferred` (unit test or browser cite from V1).
5. `TestINFERREDScopeMemberSeedRoundTrip` PASS — provenance survives export/import.

### Block V4 — Laws 6–7 (bounded graph + truncation honesty)

```bash
go test ./internal/retrieval/ -run 'TestProjectGraphTruncationHonesty|TestProjectGraphScopeEdgesAndFilter|TestProjectGraphIncludesGoalAndTask' -count=1 \
  2>&1 | tee "$EVID/v4-retrieval-caps.txt"
{
  echo "=== API hard cap ==="
  grep -n 'MaxNeighborhoodNodes\|5000' internal/retrieval/neighborhood.go internal/retrieval/project_graph.go
  echo "=== GUI caps ==="
  grep -n 'PROJECT_MAX_NODES\|EDGE_OVERVIEW_MAX\|500\|150' web/src/lib/overviewCompose.ts web/src/lib/graphLayout.ts
  echo "=== max_nodes required ==="
  grep -n 'max_nodes is required' internal/retrieval/project_graph.go
} 2>&1 | tee "$EVID/v4-cap-greps.txt"
cd web && node --experimental-strip-types --test src/lib/graphLayout.test.ts src/lib/overviewCompose.test.ts \
  2>&1 | tee "../$EVID/v4-web-cap-tests.txt"
```

**Pass:** retrieval rejects `max_nodes > 5000`; truncation honesty test PASS; GUI constants still **500** / **150**; no default full-graph dump API/GUI path without `max_nodes` + truncation flags.

### Block V5 — Law 19 (retrieval walk vs GUI adapter)

```bash
{
  echo "=== GUI must NOT walk scope_member for membership ==="
  rg -n 'scope_member|GetScope|ListScopes|ProjectGraph' web/src || true
  echo "=== Layout uses API scope_id only ==="
  rg -n 'resolveScopeClusterId|computeScopeCentroids|scope_id' web/src/lib/graphLayout.ts | head -30
  echo "=== Walk/filter in retrieval ==="
  rg -n 'Scope:|scope filter|scope_id' internal/retrieval/project_graph.go internal/retrieval/scope_graph_test.go | head -40
} 2>&1 | tee "$EVID/v5-law19-rg.txt"
go test ./internal/retrieval/ -run 'TestProjectGraphScopeEdgesAndFilter' -count=1 \
  2>&1 | tee "$EVID/v5-scope-filter-test.txt"
```

**Pass:** scope filter + graph walk live in `internal/retrieval`; `web/src` consumes API fields (`scope_id`, optional `scope` query) — **no** second source of truth / membership inference in `web/`.

### Block V6 — M-001 (task loop / gate unchanged; scope enriches only)

```bash
go test ./internal/domain/ -run 'TestInferScopes_NoGateRels' -count=1 \
  2>&1 | tee "$EVID/v6-no-gate-rels.txt"
{
  echo "=== goal_has_task synthesis still present ==="
  grep -n 'goal_has_task' internal/retrieval/graph_neighbors.go | head -10
  echo "=== infer must not touch loop/gate/install paths ==="
  rg -n 'InferScopes|scope_member' internal/loop internal/install cmd/trace/loop.go 2>/dev/null || echo "PASS: no scope infer in loop/gate CLI"
} 2>&1 | tee "$EVID/v6-m001-rg.txt"
```

**Pass:** infer does not create gate rels; goal-centric `goal_has_task` synthesis intact; scope MVP rels remain enrichment in `entity_links` (Law 13), not a replacement for Tasks → Loop → Gate → Review moat.

### Block V7 — Portable graph (`trace/graph.json`)

```bash
go test ./internal/domain/ -run 'TestSeedExportImportScopesRoundTrip|TestINFERREDScopeMemberSeedRoundTrip|TestNoINFERREDFromScopePaths' -count=1 \
  2>&1 | tee "$EVID/v7-seed-roundtrip.txt"
{
  echo "=== seedExportLinkRels includes MVP scope rels ==="
  grep -n 'seedExportLinkRels\|RelScopeMember\|RelAPIContract\|RelImplements\|RelBlocks' internal/domain/seed_export.go | head -20
} 2>&1 | tee "$EVID/v7-export-allowlist.txt"
go run ./cmd/trace seed export -o trace/graph.json 2>&1 | tee "$EVID/v7-export-cmd.txt"
test -f trace/graph.json && head -80 trace/graph.json | tee "$EVID/v7-graph-json-head.txt"
# Import smoke (temp store):
TMP=$(mktemp -d)
go run ./cmd/trace -C "$TMP" init && go run ./cmd/trace -C "$TMP" seed import trace/graph.json \
  2>&1 | tee "$EVID/v7-import-smoke.txt"
rm -rf "$TMP"
```

**Pass:** export allowlist includes MVP scope rels; round-trip tests PASS; `trace seed export -o trace/graph.json` succeeds; import smoke exit 0. Note in VERIFY-NOTES if live `trace/graph.json` has empty `scopes` array (omitempty) — not a FAIL if tests green.

### Block V8 — Backend regression floor (required)

```bash
go test ./internal/domain/ ./internal/store/ ./internal/retrieval/ ./internal/mcp/ ./internal/httpapi/ ./cmd/trace/ -count=1 \
  2>&1 | tee "$EVID/v8-backend-regression.txt"
```

**Pass:** all packages exit 0 (re-confirms S02/S03 test suites).

### Block V9 — Residuals + aggregate (required in notes)

1. Cite S00–S04 PASS board Notes / artifacts (paths above).
2. Tick §10 map V1–V7 in VERIFY-NOTES.
3. List non-blocking residuals from locked table + [`DR-HANDOFF.md`](../../DR-HANDOFF.md) queue.
4. Overall PASS only if V1–V7 + V8 green and no fail criteria tripped.

## Do not

- Change product Go/TS/CSS/docs in this row (spawn or FAIL instead)
- Close DR-HANDOFF (that is **P44-S05-02**)
- Claim/scaffold successor beyond citing lean **`no successor`** (S05-02 only)
- Reopen S02–S04 as “fix” without spawn
- Raise GUI 500 / overview 150 / API 5000 caps
- Fake PASS without evidence dir + VERIFY-NOTES

## VERIFY-NOTES.md template (required)

Write `docs/phases/phase-44-scoped-graph-linking/scopes/scope-05-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — P44-S05-01

**Date:** YYYY-MM-DD
**Git SHA:** …
**Overall:** PASS | FAIL
**Evidence:** experiments/runs/YYYY-MM-DD-p44-s05-01-verify/evidence/
**Precondition:** P44-S04-02 APPROVE; S00–S04 done; DESIGN-LOCKS §10

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS/FAIL | |
| V1 Auth FE/BE smoke | PASS/FAIL | clusters=…; api_contract=…; browser LIVE \| WAIVED … |
| V2 Agent link | PASS/FAIL | |
| V3 Inference | PASS/FAIL | |
| V4 Laws 6–7 | PASS/FAIL | |
| V5 Law 19 | PASS/FAIL | |
| V6 M-001 / Law 13 | PASS/FAIL | |
| V7 Portable graph | PASS/FAIL | |
| V8 Backend regression | PASS/FAIL | |
| V9 Residuals + aggregate | listed | |

## §10 acceptance (DESIGN-LOCKS)
- [ ] V1 seeded scopes → ≥2 clusters + cross-scope edge — …
- [ ] V2 trace_link / CLI scope_member + api_contract — …
- [ ] V3 inferred distinct; explicit wins — …
- [ ] V4 Laws 6–7 bounded + truncation honest — …
- [ ] V5 Law 19 retrieval walk; GUI adapter — …
- [ ] V6 M-001 loop/gate preserved — …
- [ ] V7 portable graph export/import — …

## Aggregate (S00–S04)
- S00 RESEARCH: …
- S01 LOCKED D1–D5: …
- S02 backend: …
- S03 infer + seed honesty: …
- S04 GUI: …

## Residuals (non-blocking)
- …

## Failures (if any)
- …

## DR-HANDOFF
remains OPEN — close owner **P44-S05-02**

## Next
P44-S05-02
```

## Todo updates

Status + notes on **P44-S05-01** only. Do not mark S05-02 done.

## Exit criteria

- [ ] `VERIFY-NOTES.md` with overall PASS/FAIL and blocks 0 + V1–V9
- [ ] Evidence dir under `experiments/runs/…-p44-s05-01-verify/evidence/`
- [ ] §10 V1–V7 explicitly ticked
- [ ] Board Notes summarize results + evidence path
- [ ] DR-HANDOFF still **OPEN**
- [ ] Next: **P44-S05-02**

## Next

`P44-S05-02`
