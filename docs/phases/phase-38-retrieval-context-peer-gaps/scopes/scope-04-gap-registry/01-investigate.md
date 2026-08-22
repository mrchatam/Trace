# P38-S04-01 — Author gap registry

## Metadata
- id: P38-S04-01
- todo_ids: [P38-S04-01]
- role: implementer
- skills: [analyst, research, planning-and-task-breakdown]
- mcps: [user-trace]
- verification: mixed
- hooks: none

## Objective

Synthesize S01–S03 into **`GAP-REGISTRY.md`** — single SoT for evidence-backed gaps, non-gaps, and deferred ideas. Cross-matrix **Trace | CG | UA | GF | MP | moat row**. **No REMEDIATION-PLAN. No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §2 H1–H11 verify/reject authority
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — saturation forward-fit; S05 spawn rules
- [INTAKE.md](../../INTAKE.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner locks from P38-S04-00
- **Required inputs (read-only — do not re-audit from scratch):**
  - [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md)
  - [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md)
  - [PEER-UA-GF.md](../scope-03-ua-graphify-peer/PEER-UA-GF.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults (P38-S04-00 — do not re-debate)

| Item | Value |
|------|-------|
| Output path | `scopes/scope-04-gap-registry/GAP-REGISTRY.md` |
| Product edits | **Forbidden** (Go/TS/web) |
| Method | **Synthesis** of approved S01–S03 artifacts + H11 doc read — not full re-audit |
| Evidence root | `experiments/runs/YYYY-MM-DD-p38-s04-660/evidence/` (H11 + cross-matrix notes only) |
| Matrix columns (LOCK) | **Trace \| CG \| UA \| GF \| MP \| moat row** |
| Gap IDs | **G-001…G-011** (1:1 with H1–H11 unless dedupe note); moat **M-001** |
| Verdict vocabulary | `gap` \| `non-gap` \| `defer` (+ DR-NOSSEM / law sub-tags) |
| Severity column | **Investigation confidence only** (`high` \| `medium` \| `low`) — **not** build priority (S06) |
| Dual-side evidence | Every **gap** row: Trace cite **and** ≥1 peer cite (or explicit “peer N/A” with reason) |
| Spawn | Incomplete registry → §6 spawn list with S05 triggers; may spawn back to S01–S03 |
| Non-goals | No ranked remediation; no G1–Gn themes (S06); no product code |

## Must answer (planner handoff — embed in GAP-REGISTRY.md)

1. **Unified gap IDs (G-001…G-011)** linked to H* with 1:1 mapping table in §2 preamble.
2. **Severity:** investigation confidence only — cite lowest confidence across contributing scopes per row.
3. **Moat row (M-001):** Trace strengths peers lack — merge TRACE-AUDIT §5 + PEER-CG §5 + PEER-UA-GF §5; dedupe.
4. **Spawn triggers for S05:** explicit list when registry incomplete (missing dual-side evidence, H11 inconclusive, H7 compose-equivalence untested, etc.).

---

## Investigation todos (run in order; do not skip)

### T0 — Preflight + evidence folder

```bash
EV=experiments/runs/$(date +%Y-%m-%d)-p38-s04-660/evidence
mkdir -p "$EV"
```

- Confirm S01–S03 artifacts exist and are APPROVED (board rows 652, 655, 658).
- Record date + row id in every new evidence file header.
- Copy **no** large JSON blobs — link existing `$EV/` paths from S01–S03.

### T1 — Ingest S01 TRACE-AUDIT (Trace column seed)

**Read:** TRACE-AUDIT §§1–6.

**Extract per H* owned or partial in S01:**

| H | TRACE-AUDIT verdict | Confidence | Primary evidence path |
|---|---------------------|------------|----------------------|
| H1 partial | confirmed gap | high | `…-s01-651/evidence/h1-trace-partial.md` |
| H2 | confirmed gap | high | `h2-compiler-fts.txt` |
| H3 | confirmed gap | high | `h3-layers-designed-vs-shipped.md` |
| H5 | confirmed gap | high | `h5-index-langs.txt` |
| H6 | confirmed gap | high | `h6-mcp-surface.md` |
| H8 partial → S03 resolved | inconclusive → see S03 | — | defer to T3 |
| H9 | confirmed gap | high | `h9-intent-pipeline.md` |
| H10 | confirmed gap | medium | `h10-install-moat.md` |

**Verdict target:** `$EV/t1-trace-column-seed.md` — table: H \| Trace state \| evidence link \| file:line anchor.

**Rule:** Do **not** re-run live CLI unless spot-checking a contested claim; link S01 evidence.

### T2 — Ingest S02 PEER-CG (CG column seed)

**Read:** PEER-CG §§1–6.

**Extract CG column for H1 partial, H5, H6, H7:**

| H | PEER-CG verdict | CG mechanism cite | Trace contrast |
|---|-----------------|-------------------|----------------|
| H1 partial | supported | `tools.ts` L1168–1170 explore query-only | TRACE-AUDIT H1 |
| H5 | supported | `watcher.ts` L68–69; 29 extractors | TRACE-AUDIT H5 |
| H6 | supported | `tools.ts` L1275–1286 single default tool | TRACE-AUDIT H6 |
| H7 | supported | `handleExplore` L3193+; P24 still deferred | no `trace_explore` in Go |

**Verdict target:** `$EV/t2-cg-column-seed.md`.

**Also capture:** PEER-CG §4 anti-patterns (not for Trace) → GAP-REGISTRY §5 defer/reject seeds.

### T3 — Ingest S03 PEER-UA-GF (UA · GF · MP column seeds)

**Read:** PEER-UA-GF §§1–7.

**Extract peer columns:**

| H | UA | GF | MP | S03 verdict |
|---|----|----|-----|-------------|
| H1 partial | context-builder L25–79 | — | wake_up L404–431 | supported |
| H4 | optional SearchEngine | EXTRACTED/INFERRED L289–370 | `_hybrid_rank` L276–329 | supported (DR-NOSSEM split) |
| H6 slice | — | — | 35/44 MCP tools | supported (contrast) |
| H8 | onboard-builder L7+ | graph.html + worked | onboarding.py + wake_up | supported (upgrades S01 inconclusive) |
| H9 contrast | — | — | fact_checker L55–78 | supported |

**Verdict target:** `$EV/t3-ua-gf-mp-column-seed.md`.

**H8 resolution:** Document upgrade from TRACE-AUDIT inconclusive → S03 supported with dual-side cite.

**H4 DR-NOSSEM split:** Record vector/embedding leg as **law defer**; label/summary graph channel as **product gap** if evidence supports.

### T4 — Assign unified gap IDs G-001…G-011

**Lock 1:1 mapping (planner — do not renumber without reviewer spawn):**

| Gap ID | H* | Theme (short) | Expected verdict seed |
|--------|-----|---------------|----------------------|
| **G-001** | H1 | Unified query+task orient packet | gap |
| **G-002** | H2 | Compiler FTS title-only | gap |
| **G-003** | H3 | Layers 2–3 designed not shipped | gap (explicit defer in doc.go) |
| **G-004** | H4 | Semantic/concept retrieval | gap + **defer** sub-row for vector leg |
| **G-005** | H5 | Index langs + manual vs watcher | gap |
| **G-006** | H6 | MCP discovery surface (16 / 1 / 44) | gap |
| **G-007** | H7 | `trace_explore` unified read missing | gap |
| **G-008** | H8 | Graph-first onboarding hook | gap |
| **G-009** | H9 | Intent pipeline doc-only | gap |
| **G-010** | H10 | Moat under-promoted in install/harness | gap |
| **G-011** | H11 | Trace+CG dual-stack undocumented | gap or inconclusive → T9 |
| **M-001** | moat | Trace strengths peers lack | **non-gap** (moat row) |

**Verdict target:** `$EV/t4-gap-id-registry.md` — G-ID \| H* \| owner scopes \| contributing artifacts.

### T5 — Cross-matrix rows G-001–G-004 (H1–H4)

For each gap row, fill **all six matrix columns** plus verdict, severity, law fit, dual-side evidence:

**Matrix row template (repeat per G-ID):**

```markdown
| G-00N | H* | Trace state | CG | UA | GF | MP | Verdict | Severity | Law fit | Evidence |
```

**G-001 (H1) — minimum cells:**

- **Trace:** No `query` on `trace_context`; compose search+context (TRACE-AUDIT; `compiler.go`, MCP schema)
- **CG:** `codegraph_explore` query-only orient (PEER-CG §2)
- **UA:** `buildChatContext(query, taskId)` 1-hop (PEER-UA-GF §1)
- **GF:** N/A or “graph orient separate from task packet” — justify
- **MP:** `wake_up()` L0+L1 packet (PEER-UA-GF §3)
- **Severity:** high (all scopes agree structural gap)

**G-004 (H4) — DR-NOSSEM nuance:**

- Split **vector/embedding** → `defer` (law); **label/summary/concept via graph** → `gap` if GF/MP evidence + Trace FTS fail
- Trace: `doc.go` L8–9 DR-NOSSEM
- GF: EXTRACTED/INFERRED mechanism
- MP: hybrid rank (BM25+vector) — note vector leg deferred for Trace

**Verdict target:** `$EV/t5-matrix-h1-h4.md`.

### T6 — Cross-matrix rows G-005–G-007 (H5–H7)

**G-005 (H5):** Trace 5 langs manual index vs CG watcher + 29 extractors.

**G-006 (H6):** Three-way contrast — Trace 16 tools vs CG 1 default vs MP 35–44; discovery paralysis vs over-surface; **no** “copy MP tool count.”

**G-007 (H7):** P24 consolidation still deferred; compose-equivalence **open question** — if not tested, add spawn trigger (do not claim non-gap without evidence).

**Verdict target:** `$EV/t6-matrix-h5-h7.md`.

### T7 — Cross-matrix rows G-008–G-010 (H8–H10)

**G-008 (H8):** Merge S01 partial + S03 supported — peers have committed graph.html / onboard hooks / wake_up identity story.

**G-009 (H9):** Trace doc-only vs MP fact_checker shipped — UA/GF N/A OK with explicit note.

**G-010 (H10):** Install/harness messaging — Trace moat exists (M-001) but **under-promoted** vs peer orient-first READMEs.

**Verdict target:** `$EV/t7-matrix-h8-h10.md`.

### T8 — G-011 (H11) Trace+CG stack documentation

**Scope:** Doc read only — S04 owns H11.

**Read:**

- [`CONTRIBUTING.md`](../../../../../CONTRIBUTING.md) — portable graph, clone import
- [`README.md`](../../../../../README.md), [`AGENTS.md`](../../../../../AGENTS.md)
- [`docs/init/`](../../../../init/) if dual-tool workflow mentioned
- PEER-FIXTURES — complementary stack guidance

**Questions:**

- Is there a documented recommended Trace + Codegraph dual-index workflow with Law 19 boundaries?
- User dogfood notes in phase docs?

**Live (optional):** `grep -r 'codegraph\|\.codegraph' docs/ AGENTS.md README.md CONTRIBUTING.md`

**Verdict target:** `$EV/h11-stack-docs.md` — gap \| non-gap \| inconclusive + doc cites.

### T9 — Moat row M-001 synthesis

**Merge non-gap seeds (dedupe overlapping bullets):**

| Source | Strengths cited |
|--------|-----------------|
| TRACE-AUDIT §5 | Task loop, gate, evidence, progressive L0–1 packet, local-first, layer honesty |
| PEER-CG §5 | CG lacks task/gate/evidence — read-only graph |
| PEER-UA-GF §5 | UA/GF/MP lack task loop, plan tree, enforcement |

**Matrix:** One **M-001** row — for each peer column, state what they **lack** vs Trace.

**Verdict target:** `$EV/t9-moat-row-m001.md`.

**Rule:** M-001 is **non-gap** — Trace strengths; distinct from G-010 (under-promotion gap).

### T10 — Non-gaps, deferrals, peer-weaker explicit list

**Non-gaps (peers weaker — not Trace gaps):**

- CG: no task loop, gates, evidence (PEER-CG §5)
- UA/GF/MP: no enforcement / plan tree (PEER-UA-GF §5)
- CG daemon always-on (anti-pattern — PEER-CG §4) — Trace correctly avoids

**Deferrals (not gaps — law/product policy):**

- DR-NOSSEM embedding channel (H4 vector leg)
- Layers 2–3 explicit defer (`doc.go` L7) — may be **gap** (G-003) vs **defer** — use INVESTIGATION-INDEX verify/reject
- P24 `trace_explore` consolidation — **gap** (G-007) until rejected with compose evidence

**Verdict target:** `$EV/t10-non-gaps-deferrals.md`.

### T11 — Spawn list + S05 triggers

**Populate §6 when any of:**

| Trigger | Condition | Owner spawn |
|---------|-----------|-------------|
| H11 inconclusive | Doc read insufficient | S04-01a or doc slice |
| H7 compose-equivalence | No evidence that search+why+context ≈ explore | S01 or S02 live compose test |
| Dual-side missing | Gap row lacks Trace **or** peer cite | Back to owning scope |
| New H12+ | Uncovered peer slice | S01–S03 per INVESTIGATION-INDEX §5 |
| MP mechanism stale | file:line drift | S03-01a |

**Empty spawn list OK** only if every G-001…G-011 has verdict + dual-side evidence + severity.

**Verdict target:** `$EV/t11-spawn-triggers.md`.

### T12 — Synthesize GAP-REGISTRY.md

Author deliverable using T0–T11.

---

## Deliverable shape (GAP-REGISTRY.md)

### §1 Purpose

Investigation-only SoT; links S01–S03 + INVESTIGATION-INDEX; feeds S05 saturation gate. **Not** REMEDIATION-PLAN.

### §2 Gap register (cross-matrix)

**Preamble:** G-001…G-011 ↔ H1–H11 mapping table.

**Main table — columns (LOCK):**

| Gap ID | H* | Trace | CG | UA | GF | MP | Verdict | Severity | Law fit | Evidence (dual-side) |

- **Verdict:** gap \| non-gap \| defer
- **Severity:** investigation confidence only — **not** build priority
- **Evidence:** `file:line` or `$EV/` path on **both** Trace and peer side for every **gap** row

Optional sub-table for G-004 DR-NOSSEM split (vector defer vs label gap).

### §3 Moat row (M-001)

Trace strengths peers lack — full matrix row + bullet expansion.

### §4 Non-gaps (peers weaker)

Explicit list where peers are **weaker** — not Trace backlog.

### §5 Deferred (law / policy / explicit design defer)

DR-NOSSEM, layer defer honesty, anti-patterns rejected (PEER-CG §4).

### §6 Spawn list → S05

Triggers, owner scope, fold vs spawn — empty OK with justification.

### §7 Evidence index

Links to S01 `…-s01-651/`, S02 `…-s02-654/`, S03 `…-s03-657/`, S04 `…-s04-660/` — no duplicate captures.

---

## Exit criteria

- [ ] GAP-REGISTRY.md §§1–7 complete
- [ ] G-001…G-011 each mapped to H* with verdict + severity + law fit
- [ ] Matrix columns **Trace | CG | UA | GF | MP** populated (N/A justified)
- [ ] M-001 moat row present
- [ ] Every **gap** row has dual-side evidence pointer
- [ ] §6 spawn list explicit (empty OK with rationale)
- [ ] No REMEDIATION-PLAN / ranked G1–Gn / product code
- [ ] Board row P38-S04-01 → `done` with Notes (confidence)

## Next

`P38-S04-02`
