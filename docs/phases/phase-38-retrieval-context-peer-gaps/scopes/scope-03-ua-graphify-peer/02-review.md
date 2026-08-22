# P38-S03-02 — Review UA + Graphify + Mempalace peers

## Metadata
- id: P38-S03-02
- todo_ids: [P38-S03-02]
- role: reviewer
- skills: [code-review-and-quality, research]
- mcps: [user-codegraph]
- verification: mixed

## Objective

Independent review of **`PEER-UA-GF.md`** vs INVESTIGATION-INDEX, planner locks (P38-S03-00), TRACE-AUDIT baseline, and **peer repo mechanism evidence** for UA, Graphify, and Mempalace. **APPROVE**, **REQUEST CHANGES**, or **SPAWN** (new S03 investigate row). **Fresh subagent** — do not share implementer session.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §2 verify/reject for H1, H4, H6, H8, H9 (MP contrast)
- [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md) — Trace baseline (do not require re-audit)
- [PEER-CG.md](../scope-02-codegraph-peer/PEER-CG.md) — CG leg (do not re-litigate)
- [01-investigate.md](01-investigate.md) — investigation todos T0–T11 + locked defaults
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — saturation forward-fit; DR-NOSSEM
- Peer roots: [UA](../../../../similar%20projects/Understand-Anything/), [Graphify](../../../../similar%20projects/graphify/), [Mempalace](../../../../similar%20projects/mempalace/)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

---

## Checklist A — Hypothesis coverage (S03-owned)

For **each** of H1 (partial), H4, H8:

- [ ] PEER-UA-GF §4 row exists with verdict ∈ {supported, weakened, rejected, inconclusive}
- [ ] Verdict aligns with INVESTIGATION-INDEX **verified if / rejected if**
- [ ] **Peer mechanism cite:** at least one `similar projects/…` **file:line** per hypothesis (not README-only for search, edges, graph UX, hybrid rank)
- [ ] Trace contrast cites TRACE-AUDIT or Trace `file:line` where comparing shipped behavior
- [ ] H1 marked **partial** — full matrix deferred S04

**Mempalace contrast slices (required in §3 + §4):**

- [ ] **H4:** `searcher.py` `_hybrid_rank` / `search_memories` cited — not website-only
- [ ] **H6 (MP slice):** `service.py` READ/WRITE tool sets or `mcp_server.py` TOOLS — file:line
- [ ] **H9 (MP contrast):** `fact_checker.py` mechanism cited; labeled contrast-only (S01 owns Trace H9 live)

**Spot-check (minimum 4 — pick H1-UA, H4-GF, H4-MP, H8-GF):**

- [ ] Peer anchor resolves in repo (UA context-builder, GF validate/symbol_resolution, MP searcher, GF graph.html or html.py)
- [ ] Verdict consistent with cited mechanism

---

## Checklist B — Peer mechanism gate (mandatory)

Major mechanism claims **must not** be README-only.

| Claim class | Required peer evidence |
|-------------|------------------------|
| UA query+task packet | `context-builder.ts` buildChatContext + ChatContext interface — file:line |
| UA SearchEngine | `packages/core/src/search.ts` Fuse keys + search() — file:line |
| GF EXTRACTED/INFERRED | `validate.py` + `symbol_resolution.py` (or ingest path) — file:line |
| GF onboarding UX | `exporters/html.py` and/or worked `graph.html` RAW_EDGES confidence — file:line |
| MP hybrid search | `searcher.py` `_hybrid_rank`, `search_memories` — file:line |
| MP memory stack | `layers.py` wake_up / L0–L3 — file:line |
| MP MCP surface | `mcp_server.py` TOOLS or `service.py` READ_TOOLS — file:line |
| MP fact check (H9 contrast) | `fact_checker.py` check_text — file:line |

- [ ] §6 Evidence appendix lists peer paths with line numbers for **T1–T4, T6–T9**
- [ ] Reviewer independently opens ≥3 peer cites (confirm line numbers still match)
- [ ] If peer clone **missing** from workspace: document gap in Notes — **do not fake cites**

---

## Checklist C — Planner must-answer items (P38-S03-00)

Confirm PEER-UA-GF addresses all eight handoff questions:

- [ ] **Q1** UA context-builder packet vs Trace compiler (§1 + §4 H1)
- [ ] **Q2** UA SearchEngine mechanism (§1)
- [ ] **Q3** Graphify EXTRACTED/INFERRED vs H4 (§2 + §4 H4)
- [ ] **Q4** Graphify graph.html / worked examples vs H8 (§2 + §4 H8)
- [ ] **Q5** Mempalace hybrid search vs Trace FTS-only / DR-NOSSEM (§3 + §4 H4)
- [ ] **Q6** Mempalace MCP + memory stack vs Trace (§3 + §4 H1/H6/H8)
- [ ] **Q7** Mempalace fact_checker vs Trace intent (§3 + §4 H9 contrast)
- [ ] **Q8** Moat seed — peers lack task loop/gates/evidence (§5)

---

## Checklist D — Hygiene & phase law

- [ ] No Go/TS/web product diff in investigate row
- [ ] No ranked build list / REMEDIATION-PLAN language
- [ ] DR-NOSSEM acknowledged for H4 — semantic gap may be **law deferral**, not bug
- [ ] §5 Moat row present (UA/GF/MP lack tasks/gates/evidence)
- [ ] §7 Spawn list empty **or** each item has trigger + owner
- [ ] Evidence under `experiments/runs/…-p38-s03-657/evidence/` with date + row id
- [ ] Mempalace §3 present (human-added peer 2026-08-22)

---

## Checklist E — Saturation forward-fit (DESIGN-LOCKS)

S03 peer audit must feed S04 matrix + S05 gate:

- [ ] S03-owned H* verdicts evidence-backed enough for GAP-REGISTRY
- [ ] Peer mechanism cites satisfy DESIGN-LOCKS “not README-only” for UA, GF, **MP**
- [ ] H1 partial boundary preserved — CG leg in PEER-CG; S04 aggregates
- [ ] H12+ not invented without uncovered slice — mempalace mapped to existing H* unless spawn justified

---

## Verdict

Document in board Notes:

```text
Verdict: APPROVE | REQUEST CHANGES | SPAWN
Confidence: high | medium | low
Findings: blocker | high | medium | low | nit counts
Residual risks: (if medium confidence)
```

| Severity | Action |
|----------|--------|
| Blocker / high | Spawn `P38-S03-01a` implement + `P38-S03-01b` review **below this row**, OR fix ≤5 lines inline if trivial cite gap |
| Medium | Prefer spawn unless single-table row fix in PEER-UA-GF |
| Low / nit | Note only; may APPROVE with listed residuals |

**SPAWN** when any peer slice still shallow (e.g. MP hybrid not cited, GF edges README-only, H8 graph.html not opened, mempalace missing with fake cites).

---

## Exit criteria

- [ ] Verdict recorded in board Notes with confidence
- [ ] No open blocker/high without pending spawn row
- [ ] Board row P38-S03-02 → `done`

## Next

**P38-S04-00** on APPROVE.
