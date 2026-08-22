# INVESTIGATION-INDEX — Phase 38 retrieval & context peer gaps

**Author:** P38-S00-01 (2026-08-22)  
**Status:** Investigation plan only — no product changes, no remediation tasks.

---

## §1 Purpose

Phase 38 is an **investigation-only** phase. Its job is to test hypotheses (H1–H11) about Trace retrieval, context assembly, MCP ergonomics, and onboarding against peer projects — then aggregate evidence through S04 and exit investigation loops only at the **S05 saturation gate**.

This index is the **single routing document** for all downstream investigate rows (S01–S04). It registers hypotheses, assigns owner scopes, locks investigation methods, and defines spawn rules. It does **not** rank fixes, propose build tasks, or pre-judge outcomes.

**Authoritative inputs:**

| Document | Role |
|----------|------|
| [INTAKE.md](../../INTAKE.md) | Hypothesis backlog source; human locks |
| [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) | Saturation exit criteria (S05), investigation row rules, law fit |
| [PEER-FIXTURES.md](../../PEER-FIXTURES.md) | Local clone paths, optional tooling, evidence convention |
| [SCOPE-TODOS.md](SCOPE-TODOS.md) | Hypothesis → scope routing, tool matrix |
| P24 [EXTERNAL-RESEARCH.md](../../../phase-24-agent-effectiveness-investigation/scopes/scope-03-external-research/EXTERNAL-RESEARCH.md) | Prior peer table — **extend with live re-verify**, do not copy blindly |

**Phase flow:**

```text
S00 (this index) → S01 Trace live → S02 Codegraph → S03 UA+Graphify+Mempalace
  → S04 cross-matrix → S05 saturation gate → S06 remediation PLAN → S07 VERIFY
```

Investigation may **spawn** additional rows in S01–S04 (or `a`/`b` suffix cycles) until S05 reviewer signs confident exit. **No implement** in P38; successor phase owns build.

---

## §2 Hypothesis register

### Methods (locked)

| Method | Meaning |
|--------|---------|
| **Trace live** | CLI/MCP on Trace repo (or dogfood fixture) |
| **Peer read** | File:line under `similar projects/` |
| **Both** | Live Trace + peer mechanism cite |

### Tools column (locked)

| Code | Tool |
|------|------|
| **T** | Trace MCP/CLI |
| **CG** | Codegraph MCP `codegraph_explore` |
| **GF** | Graphify skill / worked examples |
| **UA** | UA source read |
| **MP** | Mempalace source read (hybrid search, layers, MCP) |
| **—** | Peer read only |

### H1–H11 (locked — copy from P38-S00-00 planner)

| H | Statement (short) | Peers | Trace areas | Method | Tools | Peer cite targets | Evidence target | Verified if… | Rejected if… | Owner |
|---|-------------------|-------|-------------|--------|-------|-------------------|-----------------|--------------|--------------|-------|
| **H1** | No unified query+task context packet | UA, CG, MP | `compiler`, MCP | Both | T, CG, UA, MP | UA [`context-builder.ts`](../../../../similar%20projects/Understand-Anything/understand-anything-plugin/src/context-builder.ts) L20–79; CG README agent workflow + `codegraph_explore`; MP [`layers.py`](../../../../similar%20projects/mempalace/mempalace/layers.py) L404–431 `wake_up`; Trace `trace context` / MCP `trace_context` | `experiments/runs/…-p38-s01-651/evidence/h1-trace-context.json`; S02/S03 peer cites | Trace packet lacks query-driven symbol/neighborhood assembly that peers provide in one orient step; side-by-side mechanism table shows structural gap | Trace `trace context` with task id already merges query (via search/expand) + task packet comparably to UA 1-hop or CG explore; or gap is harness/docs not product | S01+S02+S03 (matrix in S04) |
| **H2** | Compiler FTS uses task title only, not agent query | UA | `internal/compiler/compiler.go` | Trace live + Peer read | T, UA | Trace [`compiler.go`](../../../../../internal/compiler/compiler.go) ~L146 (title FTS); UA SearchEngine in context-builder; P24 EXTERNAL-RESEARCH UA row | `…-p38-s01-651/evidence/h2-compiler-fts.txt` | Live `trace context` / compiler path shows retrieval query ≠ agent question (only title tokens drive FTS assumptions) | Code path accepts explicit query param or expand uses search query from MCP input; documented and live-verified | S01 |
| **H3** | Layers 2–3 designed but not shipped | Aider | `compiler`, docs | Trace live + Peer read | T | [`internal/compiler/doc.go`](../../../../../internal/compiler/doc.go) L7; RETRIEVAL_AND_CONTEXT.md layer spec; Aider repo map docs (web/P24 AID row) | `…-p38-s01-651/evidence/h3-layers-packet.json` | Layer 2–3 fields absent from live packet JSON; docs describe them as future | Layers 2–3 present in live packet or explicitly deferred with shipped alternative meeting same intent | S01 |
| **H4** | No semantic retrieval (DR-NOSSEM) limits concept discovery | Graphify, MP | `retrieval` | Peer read + Trace live | T, GF, MP | [`internal/retrieval/doc.go`](../../../../../internal/retrieval/doc.go) L8–9; Graphify [`validate.py`](../../../../similar%20projects/graphify/graphify/validate.py) + worked GRAPH_REPORT EXTRACTED/INFERRED; MP [`searcher.py`](../../../../similar%20projects/mempalace/mempalace/searcher.py) L276–329 `_hybrid_rank` | `…-p38-s03-657/evidence/h4-semantic-contrast.md` | Concept-style query succeeds on Graphify graph labels or MP vector+BM25 but fails on Trace FTS-only (documented query pair) | Trace non-semantic channels + graph expand cover same queries adequately; or DR-NOSSEM makes semantic channel correctly out of scope (defer, not gap) | S03 (+ S04) |
| **H5** | Index: lang coverage + manual index vs watcher | CG | `analyzers`, `cmd/trace/index` | Both | T, CG | CG README index/watch; Trace analyzer registry + `trace index`; reconcile INTAKE "3 langs" vs 4 shipped (Go/JS/TS/Python) | `…-p38-s01-651/evidence/h5-index-langs.txt`; S02 watch cite | CG supports materially more langs OR auto-sync; Trace requires manual re-index without watch; lang gap confirmed with counts | Trace watcher exists or lang delta insignificant for Trace targets; manual index acceptable by law | S01 + S02 |
| **H6** | MCP 16 tools vs CG 1 → discovery paralysis | CG, CM, MP | `internal/mcp` | Trace live + Peer read | T, CG, MP | [`internal/mcp/server.go`](../../../../../internal/mcp/server.go) `RegisteredToolNames`; CG MCP single-tool tests/README; MP [`service.py`](../../../../similar%20projects/mempalace/mempalace/service.py) L60–102 READ/WRITE sets + [`mcp_server.py`](../../../../similar%20projects/mempalace/mempalace/mcp_server.py) TOOLS; CM README (P24 CM row) | `…-p38-s01-651/evidence/h6-mcp-tool-list.txt`; S03 `h6-mp-mcp-surface.md` | 16 read/write tools without ranked "start here"; peer doc/tooling shows simpler orient path; FM-08-style paralysis plausible | Tool descriptions + `trace_capability` adequately gate discovery; parity tests document intentional surface | S01 + S02 (+ S03 MP slice) |
| **H7** | `trace_explore` unified read missing (P24 deferred) | CG, P24 | MCP + library | Both | T, CG | P24 EXTERNAL-RESEARCH CG row; CG explore implementation; Trace MCP inventory (no explore) | `…-p38-s02-654/evidence/h7-explore-gap.md` | No single Trace tool returns ranked symbols + paths + blast radius like `codegraph_explore`; P24 transfer still open | `trace_search`+`trace_why`+`trace_context` compose equivalently with evidence; or explicitly rejected as anti-pattern (tool sprawl consolidation) | S02 (+ S04) |
| **H8** | Explore GUI vs graph-first onboarding | Graphify, UA, MP | `web/`, Phase 32–33 | Peer read + Trace live (GUI optional) | T, GF, UA, MP | Graphify [`worked/rsl-siege-manager/graph.html`](../../../../similar%20projects/graphify/worked/rsl-siege-manager/graph.html); UA [`onboard-builder.ts`](../../../../similar%20projects/Understand-Anything/understand-anything-plugin/src/onboard-builder.ts); MP [`onboarding.py`](../../../../similar%20projects/mempalace/mempalace/onboarding.py) + [`layers.py`](../../../../similar%20projects/mempalace/mempalace/layers.py) wake-up; Trace `web/` routes | `…-p38-s03-657/evidence/h8-onboarding-ux.md`; optional GUI screenshot in S01 | Peers offer committed graph/memory artifact + human hook before implement; Trace GUI lacks comparable orient/onboarding | Trace GUI or install docs provide equivalent graph-first hook; or hook intentionally deferred (human-gated) | S03 (+ S01 partial T7) |
| **H9** | Intent extraction pipeline documented not implemented | RETRIEVAL_AND_CONTEXT.md, MP | `retrieval`, `compiler` | Trace live + doc read + peer contrast | T, MP | [`docs/RETRIEVAL_AND_CONTEXT.md`](../../../../RETRIEVAL_AND_CONTEXT.md) intent sections; grep `intent` in `internal/retrieval`, `internal/compiler`; MP [`fact_checker.py`](../../../../similar%20projects/mempalace/mempalace/fact_checker.py) L55–78 (contrast only) | `…-p38-s01-651/evidence/h9-intent-pipeline.md`; S03 `h9-mp-fact-check-contrast.md` | Doc describes pipeline; no shipped code path; live trace shows absence | Pipeline partially shipped under different name; or doc marked aspirational and superseded | S01 (+ S03 MP contrast) |
| **H10** | Trace moat under-promoted vs peer orient-first | OH, SWE, CG | `install`, harness | Peer read + Trace live | T | Trace install/bootstrap docs; P24 OH/SWE/CG harness rows; `trace install` output | `…-p38-s01-651/evidence/h10-install-moat.md` | Install/skills emphasize tasks/gates/evidence less than peers emphasize graph orient | Install + MCP instructions lead with moat (task loop, gate, evidence) before code grep | S01 (+ S04 moat row) |
| **H11** | Trace+Codegraph stack undocumented | User dogfood | docs | Doc read + optional live | T, CG | CONTRIBUTING, README, AGENTS.md, install skills; user dogfood notes if any | `…-p38-s04-660/evidence/h11-stack-docs.md` | No documented recommended dual-index workflow | Doc exists recommending complementary stack with boundaries (Law 19, local-first) | S04 |

### Spawn slots (template only — do not invent gaps)

| H* | Statement | Source | Owner | Notes |
|----|-----------|--------|-------|-------|
| H12+ | _(empty)_ | human-added peer mempalace 2026-08-22 mapped to H1/H4/H6/H8/H9 — no new H12 at S03 planner | S03 | Spawn only if S03-01 finds uncovered slice |

---

## §3 Peer map

Investigation scopes and their primary peer focus:

| Scope | Board rows | Artifact | Primary peers | Hypotheses (primary) |
|-------|------------|----------|---------------|----------------------|
| **S01** Trace live audit | P38-S01-* | `TRACE-AUDIT.md` | Trace CLI/MCP (live) | H2, H3, H5, H6, H9, H10 (+ H1 partial, H8 GUI optional) |
| **S02** Codegraph peer | P38-S02-* | `PEER-CG.md` | Codegraph (`similar projects/codegraph/`) | H1 partial, H5, H6, H7 |
| **S03** UA + Graphify + Mempalace | P38-S03-* | `PEER-UA-GF.md` | UA, Graphify worked examples, **Mempalace** | H1 partial, H4, H8 (+ H6/H9 MP slices) |
| **S04** Cross-matrix | P38-S04-* | `GAP-REGISTRY.md` | All peers + Trace strengths | H11; aggregate H1–H11 verdicts |
| **S05** Saturation gate | P38-S05-* | `SATURATION-NOTES.md` | — | Exit investigation loops; spawn or APPROVE |
| **S06** Remediation plan | P38-S06-* | `REMEDIATION-PLAN.md` | — | **After S05 only** — plan, not build |

**Peer fixture paths** (read-only): see [PEER-FIXTURES.md](../../PEER-FIXTURES.md).

**Hypothesis routing summary:**

| Hypotheses | Primary scope | Method emphasis |
|------------|---------------|-----------------|
| H2, H3, H5, H6, H9, H10 | S01 | Trace live + doc read |
| H1 (partial), H5, H6, H7 | S02 | Peer read + optional CG MCP |
| H1 (partial), H4, H8 (+ H6/H9 MP) | S03 | Peer read + worked examples + mempalace source |
| H11 | S04 | Doc read |
| All | S04 cross-matrix; S05 saturation | Aggregate + gate |

---

## §4 Investigation methods

### 4.1 Method classes

| Class | When to use | Examples |
|-------|-------------|----------|
| **Trace live (CLI)** | Verify shipped behavior, capture packet JSON | `trace context`, `trace search`, `trace why`, `trace index`, `trace install` |
| **Trace live (MCP)** | Same via harness; compare tool surface | `trace_context`, `trace_search`, `trace_capability`, `trace_plan` |
| **Peer file:line read** | Mechanism cite without running peer stack | UA `context-builder.ts`, CG explore source, Graphify `GRAPH_REPORT.md`, Mempalace `searcher.py` |
| **Codegraph MCP (optional)** | Live orient compare when `.codegraph/` exists | `codegraph_explore` with `projectPath` |
| **Graphify worked examples** | Concept/semantic contrast, graph artifact UX | `similar projects/graphify/worked/*/GRAPH_REPORT.md`, `graph.html` |
| **UA context-builder read** | Query-driven assembly pattern | SearchEngine usage, 1-hop neighborhood assembly |
| **Doc read** | Design intent vs shipped | `RETRIEVAL_AND_CONTEXT.md`, `internal/*/doc.go`, install/bootstrap docs |
| **GUI optional** | H8 onboarding compare | Trace `web/` Explore/Overview routes; screenshot evidence |

### 4.2 Evidence convention

All investigate rows store evidence under:

```text
experiments/runs/YYYY-MM-DD-p38-<scope>-<row>/evidence/
```

| Scope | Example row slug | Example evidence files |
|-------|------------------|------------------------|
| S00 | `s00-648` | Index authoring — minimal; cite locked register (`mcp-tool-count.txt` sanity) |
| S01 | `s01-651` | `h2-compiler-fts.txt`, `h6-mcp-tool-list.txt`, CLI JSON captures |
| S02 | `s02-654` | `h7-explore-gap.md`, CG MCP stdout if run |
| S03 | `s03-657` | `h4-semantic-contrast.md`, `h8-onboarding-ux.md` |
| S04 | `s04-660` | `h11-stack-docs.md`, cross-matrix tables |

**Each evidence file must include:**

- **Date** and **board row id**
- **Command or file read** (reproducible)
- **File:line cites** on both Trace and peer sides where applicable
- No secrets; no mutating consumer `.trace/` without row permission

### 4.3 Tool column reference

See §2 tools column: **T** · **CG** · **GF** · **UA** · **MP** · **—**

Optional tool matrix by hypothesis: [SCOPE-TODOS.md](SCOPE-TODOS.md).

### 4.4 Saturation forward-fit

S01–S03 investigate rows must produce evidence sufficient for [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) S05 checklist:

- Every H1–H11 → verified / rejected / deferred with evidence pointer
- At least one live Trace command trace per major gap claim
- Each primary peer (CG, UA, GF, MP) → mechanism cite (not README-only)
- Cross-matrix includes Trace strengths (moat row)
- Spawn list empty or explicitly deferred with trigger

---

## §5 Spawn rules

When S01–S03 reviewers discover work not covered by this register:

| Trigger | Action | Fold vs spawn |
|---------|--------|---------------|
| New hypothesis H12+ not in INTAKE | Insert investigate+review pair **below** the scope review that discovered it | **Spawn** — never fold into unrelated row |
| Single hypothesis needs >3 peer files or >2 live command classes | Dedicated investigate row in owning scope (S01/S02/S03) | **Spawn** when review scope would exceed ~7 todos |
| Minor cite fix / typo in peer path | Fix in place in INVESTIGATION-INDEX or scope artifact | **Fold** into S00-02 review trivial fix (≤5 lines) |
| Cross-peer comparison already on S04 matrix | Add row to GAP-REGISTRY only | **Fold** into S04-01 |
| Saturation reviewer sees duplicate work | Defer in SATURATION-NOTES with trigger | **No spawn** — document reject |
| CG/UA/GF/MP mechanism already cited in S02/S03 | Reference existing evidence path | **Fold** — link, don't re-run |

**Depth limit:** Max **2 spawn cycles** per scope (`a`/`b` suffix rows) before S05 unless S05-02 explicitly extends.

**Board insertion:** New rows go **immediately below** the review row that spawned them. Forward-only — do not rewrite `done` history.

---

## §6 Non-goals

Phase 38 investigation explicitly excludes:

| Non-goal | Rationale |
|----------|-----------|
| **Product code changes** | Go/TS implement forbidden; evidence captures only |
| **GAP-REGISTRY content in S00** | S04 owns aggregate registry |
| **REMEDIATION-PLAN content** | S06 after S05 saturation APPROVE only |
| **Ranked build list** | "Potential improvement" ≠ accept for build until REMEDIATION-PLAN |
| **Semantic/embeddings spike** | DR-NOSSEM noted; no embedding experiments in P38 |
| **Reopening Phase 36/37** | P37 closed; P38 is separate investigation |
| **Hosted SaaS / daemon defaults** | Law fit; local-first only |
| **Full peer monorepo scans** | Read mechanism cites + worked examples; Graphify/UA per PEER-FIXTURES |
| **Implement tasks disguised as investigation** | Rows cite and compare only; fix proposals belong in S06 plan |

---

## Appendix — S00 sanity evidence

| File | Row | Purpose |
|------|-----|---------|
| [`experiments/runs/2026-08-22-p38-s00-648/evidence/mcp-tool-count.txt`](../../../../../experiments/runs/2026-08-22-p38-s00-648/evidence/mcp-tool-count.txt) | P38-S00-01 | Optional MCP tool inventory sanity (`TestToolNamesRegistered`) |

**Peer path spot-check (P38-S00-01):** UA `context-builder.ts`, Graphify `worked/rsl-siege-manager/graph.html`, Codegraph clone, Trace `internal/compiler/doc.go`, `internal/retrieval/doc.go`, `internal/mcp/server.go` — all resolve.
