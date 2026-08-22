# P38-S00-01 — Author investigation index

## Metadata
- id: P38-S00-01
- todo_ids: [P38-S00-01]
- role: implementer
- skills: [research, graphify, code-explorer, planning-and-task-breakdown]
- mcps: [user-trace, user-codegraph]
- verification: manual

## Objective

Author **`INVESTIGATION-INDEX.md`** in this scope folder. Register all H1–H11 from INTAKE with investigation plan — **not** remediation tasks. Use the **locked hypothesis register** below as authoritative content; expand into deliverable §§1–6 without re-debating methods or criteria.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner locks from P38-S00-00
- P24 [EXTERNAL-RESEARCH.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md) — extend, do not copy blindly

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P38-S00-00 — do not re-debate)

| Item | Value |
|------|-------|
| Output path | `scopes/scope-00-investigation-index/INVESTIGATION-INDEX.md` |
| Product edits | **Forbidden** |
| Mode | Investigate only — no implement suggestions as tasks |
| Hypothesis count | H1–H11 from INTAKE + optional H* spawn slots (empty template only) |
| Spawn rule | New H* or peer slice → new board row in S01–S03 (or S01–S04 Na/Nb), not silent backlog |
| Evidence root | `experiments/runs/YYYY-MM-DD-p38-<scope>-<row>/evidence/` |
| Saturation SoT | DESIGN-LOCKS S05 criteria — index must prepare S01–S03 to satisfy them |

## Deliverable shape (INVESTIGATION-INDEX.md)

### §1 Purpose
Investigation-only phase; exit loops only at S05 saturation gate. Link DESIGN-LOCKS + INTAKE.

### §2 Hypothesis register
Table columns: `H id | statement | peers | Trace areas | method | tools | peer cite targets | evidence target | verified criteria | rejected criteria | owner scope`

Copy locked register below; add spawn-slot rows (`H12+`) only if INTAKE/P24/conversation surfaced a gap — cite source.

### §3 Peer map
S01 Trace live · S02 Codegraph · S03 UA+Graphify · S04 cross-matrix · S05 saturation.

### §4 Investigation methods
Live Trace CLI/MCP · peer file:line read · optional Codegraph MCP · Graphify worked examples · UA context-builder tests. Evidence convention (§ locked below).

### §5 Spawn rules
When S01–S03 reviewer spawns: insert row immediately below review row; max 2 spawn cycles before S05 unless saturation reviewer extends.

### §6 Non-goals
No product code; no GAP-REGISTRY; no REMEDIATION-PLAN; no ranked build list; no semantic/embeddings spike in P38.

---

## Locked hypothesis register (copy into §2)

Methods: **Trace live** = CLI/MCP on Trace repo (or dogfood fixture); **Peer read** = file:line under `similar projects/`; **Both** = live Trace + peer mechanism cite.

Tools column: **T** = Trace MCP/CLI · **CG** = Codegraph MCP `codegraph_explore` · **GF** = Graphify skill/worked · **UA** = UA source read · **—** = peer read only.

| H | Statement (short) | Peers | Trace areas | Method | Tools | Peer cite targets | Evidence target | Verified if… | Rejected if… | Owner |
|---|---------------------|-------|-------------|--------|-------|-------------------|-----------------|--------------|--------------|-------|
| **H1** | No unified query+task context packet | UA, CG | `compiler`, MCP | Both | T, CG, UA | UA [`context-builder.ts`](../../../../similar%20projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts) L20–79; CG README agent workflow + `codegraph_explore`; Trace `trace context` / MCP `trace_context` | `experiments/runs/…-p38-s01-651/evidence/h1-trace-context.json`; S02/S03 peer cites | Trace packet lacks query-driven symbol/neighborhood assembly that peers provide in one orient step; side-by-side mechanism table shows structural gap | Trace `trace context` with task id already merges query (via search/expand) + task packet comparably to UA 1-hop or CG explore; or gap is harness/docs not product | S01+S02+S03 (matrix in S04) |
| **H2** | Compiler FTS uses task title only, not agent query | UA | `internal/compiler/compiler.go` | Trace live + Peer read | T, UA | Trace [`compiler.go`](../../../../../internal/compiler/compiler.go) ~L146 (title FTS); UA SearchEngine in context-builder; P24 EXTERNAL-RESEARCH UA row | `…-p38-s01-651/evidence/h2-compiler-fts.txt` | Live `trace context` / compiler path shows retrieval query ≠ agent question (only title tokens drive FTS assumptions) | Code path accepts explicit query param or expand uses search query from MCP input; documented and live-verified | S01 |
| **H3** | Layers 2–3 designed but not shipped | Aider | `compiler`, docs | Trace live + Peer read | T | [`internal/compiler/doc.go`](../../../../../internal/compiler/doc.go) L7; RETRIEVAL_AND_CONTEXT.md layer spec; Aider repo map docs (web/P24 AID row) | `…-p38-s01-651/evidence/h3-layers-packet.json` | Layer 2–3 fields absent from live packet JSON; docs describe them as future | Layers 2–3 present in live packet or explicitly deferred with shipped alternative meeting same intent | S01 |
| **H4** | No semantic retrieval (DR-NOSSEM) limits concept discovery | Graphify | `retrieval` | Peer read + Trace live | T, GF | [`internal/retrieval/doc.go`](../../../../../internal/retrieval/doc.go) L8–9; Graphify [`worked/example/README.md`](../../../../similar%20projects/graphify/worked/example/README.md) + GRAPH_REPORT EXTRACTED/INFERRED | `…-p38-s03-657/evidence/h4-semantic-contrast.md` | Concept-style query succeeds on Graphify graph labels but fails on Trace FTS-only (documented query pair) | Trace non-semantic channels + graph expand cover same queries adequately; or DR-NOSSEM makes semantic channel correctly out of scope (defer, not gap) | S03 (+ S04) |
| **H5** | Index: lang coverage + manual index vs watcher | CG | `analyzers`, `cmd/trace/index` | Both | T, CG | CG README index/watch; Trace analyzer registry + `trace index`; reconcile INTAKE "3 langs" vs 4 shipped (Go/JS/TS/Python) | `…-p38-s01-651/evidence/h5-index-langs.txt`; S02 watch cite | CG supports materially more langs OR auto-sync; Trace requires manual re-index without watch; lang gap confirmed with counts | Trace watcher exists or lang delta insignificant for Trace targets; manual index acceptable by law | S01 + S02 |
| **H6** | MCP 16 tools vs CG 1 → discovery paralysis | CG, CM | `internal/mcp` | Trace live + Peer read | T, CG | [`internal/mcp/server.go`](../../../../../internal/mcp/server.go) `RegisteredToolNames`; CG MCP single-tool tests/README; CM README tool list (P24 CM row) | `…-p38-s01-651/evidence/h6-mcp-tool-list.txt` | 16 read/write tools without ranked "start here"; peer doc/tooling shows simpler orient path; FM-08-style paralysis plausible | Tool descriptions + `trace_capability` adequately gate discovery; parity tests document intentional surface | S01 + S02 |
| **H7** | `trace_explore` unified read missing (P24 deferred) | CG, P24 | MCP + library | Both | T, CG | P24 EXTERNAL-RESEARCH CG row; CG explore implementation; Trace MCP inventory (no explore) | `…-p38-s02-654/evidence/h7-explore-gap.md` | No single Trace tool returns ranked symbols + paths + blast radius like `codegraph_explore`; P24 transfer still open | `trace_search`+`trace_why`+`trace_context` compose equivalently with evidence; or explicitly rejected as anti-pattern (tool sprawl consolidation) | S02 (+ S04) |
| **H8** | Explore GUI vs graph-first onboarding | Graphify, UA | `web/`, Phase 32–33 | Peer read + Trace live (GUI optional) | T, GF, UA | Graphify [`worked/rsl-siege-manager/graph.html`](../../../../similar%20projects/graphify/worked/rsl-siege-manager/graph.html); UA viewer/README; Trace `web/` Explore/Overview routes | `…-p38-s03-657/evidence/h8-onboarding-ux.md`; optional GUI screenshot in S01 | Peers offer committed graph artifact + human hook before implement; Trace GUI lacks comparable orient/onboarding | Trace GUI or install docs provide equivalent graph-first hook; or hook intentionally deferred (human-gated) | S03 (+ S01 partial T7) |
| **H9** | Intent extraction pipeline documented not implemented | RETRIEVAL_AND_CONTEXT.md | `retrieval`, `compiler` | Trace live + doc read | T | [`docs/RETRIEVAL_AND_CONTEXT.md`](../../../../../RETRIEVAL_AND_CONTEXT.md) intent sections; grep `intent` in `internal/retrieval`, `internal/compiler` | `…-p38-s01-651/evidence/h9-intent-pipeline.md` | Doc describes pipeline; no shipped code path; live trace shows absence | Pipeline partially shipped under different name; or doc marked aspirational and superseded | S01 |
| **H10** | Trace moat under-promoted vs peer orient-first | OH, SWE, CG | `install`, harness | Peer read + Trace live | T | Trace install/bootstrap docs; P24 OH/SWE/CG harness rows; `trace install` output | `…-p38-s01-651/evidence/h10-install-moat.md` | Install/skills emphasize tasks/gates/evidence less than peers emphasize graph orient | Install + MCP instructions lead with moat (task loop, gate, evidence) before code grep | S01 (+ S04 moat row) |
| **H11** | Trace+Codegraph stack undocumented | User dogfood | docs | Doc read + optional live | T, CG | CONTRIBUTING, README, AGENTS.md, install skills; user dogfood notes if any | `…-p38-s04-660/evidence/h11-stack-docs.md` | No documented recommended dual-index workflow | Doc exists recommending complementary stack with boundaries (Law 19, local-first) | S04 |

**Spawn slots (template only — do not invent gaps):**

| H* | Statement | Source | Owner | Notes |
|----|-----------|--------|-------|-------|
| H12+ | _(empty)_ | cite transcript / P24 / S01–S03 spawn | TBD | Reviewer inserts when investigation discovers uncovered slice |

---

## Evidence path convention (§4 in deliverable)

```
experiments/runs/YYYY-MM-DD-p38-<scope>-<row>/evidence/
```

| Scope | Example row slug | Example evidence files |
|-------|------------------|------------------------|
| S00 | `s00-648` | _(index authoring — minimal; cite locked register)_ |
| S01 | `s01-651` | `h2-compiler-fts.txt`, `h6-mcp-tool-list.txt`, CLI JSON captures |
| S02 | `s02-654` | `h7-explore-gap.md`, CG MCP stdout if run |
| S03 | `s03-657` | `h4-semantic-contrast.md`, `h8-onboarding-ux.md` |
| S04 | `s04-660` | `h11-stack-docs.md`, cross-matrix tables |

Each evidence file must include: **date**, **row id**, **command or file read**, **file:line cites** both sides where applicable. No secrets; no mutating consumer `.trace/` without row permission.

Optional live captures for index author (sanity only):
```bash
# Tool inventory
go test ./internal/mcp/ -run TestToolNamesRegistered -count=1 2>&1 | tee experiments/runs/$(date +%Y-%m-%d)-p38-s00-648/evidence/mcp-tool-count.txt
```

---

## Spawn rules (§5 in deliverable)

| Trigger | Action | Fold vs spawn |
|---------|--------|---------------|
| New hypothesis H12+ not in INTAKE | Insert investigate+review pair **below** the scope review that discovered it | **Spawn** — never fold into unrelated row |
| Single hypothesis needs >3 peer files or >2 live command classes | Dedicated investigate row in owning scope (S01/S02/S03) | **Spawn** when review scope would exceed ~7 todos |
| Minor cite fix / typo in peer path | Fix in place in INVESTIGATION-INDEX or scope artifact | **Fold** into S00-02 review trivial fix (≤5 lines) |
| Cross-peer comparison already on S04 matrix | Add row to GAP-REGISTRY only | **Fold** into S04-01 |
| Saturation reviewer sees duplicate work | Defer in SATURATION-NOTES with trigger | **No spawn** — document reject |
| CG/UA/GF mechanism already cited in S02/S03 | Reference existing evidence path | **Fold** — link, don't re-run |

Max investigation depth before S05: **2 spawn cycles** per scope (a/b suffix rows) unless S05-02 explicitly extends.

---

## Role work

1. Re-read INTAKE, PEER-FIXTURES, locked register above.
2. Optionally spot-check peer paths exist (CG, UA, Graphify under `similar projects/`).
3. Write `INVESTIGATION-INDEX.md` §§1–6 using locked content.
4. Board row → `done` with Notes: hypothesis count, spawn slots, evidence convention cited.

## Minimal todos (execute in order)

- [ ] Re-read INTAKE + PEER-FIXTURES + P24 EXTERNAL-RESEARCH §2 (UA/CG/GF/CM rows)
- [ ] Copy locked hypothesis register into §2; verify peer paths resolve
- [ ] Write §§1,3–6 (purpose, peer map, methods, spawn, non-goals)
- [ ] Add evidence convention + tool column to §4
- [ ] Self-check: no implement language ("add feature X in P38")
- [ ] Board row → `done` with Notes

## Exit criteria

- [ ] Every H1–H11 has method, tools, peer cites, verified/rejected criteria, owner scope
- [ ] Spawn rules + non-goals documented
- [ ] Evidence path convention with examples
- [ ] No implement tasks disguised as investigation
- [ ] `INVESTIGATION-INDEX.md` exists at scope path

## Next

`P38-S00-02`
