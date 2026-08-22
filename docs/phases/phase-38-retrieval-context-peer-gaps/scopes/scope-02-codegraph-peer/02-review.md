# P38-S02-02 — Review Codegraph peer

## Metadata
- id: P38-S02-02
- todo_ids: [P38-S02-02]
- role: reviewer
- skills: [code-review-and-quality, research]
- mcps: [user-codegraph]
- verification: mixed

## Objective

Independent review of **`PEER-CG.md`** vs INVESTIGATION-INDEX, planner locks (P38-S02-00), TRACE-AUDIT baseline, and **peer repo mechanism evidence**. **APPROVE**, **REQUEST CHANGES**, or **SPAWN** (new S02 investigate row). **Fresh subagent** — do not share implementer session.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §2 verify/reject for H1, H5, H6, H7
- [TRACE-AUDIT.md](../scope-01-trace-audit/TRACE-AUDIT.md) — Trace baseline (do not require re-audit)
- [01-investigate.md](01-investigate.md) — investigation todos T0–T9 + locked defaults
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — saturation forward-fit; Law 6/7
- Peer root: [`similar projects/codegraph/`](../../../../../similar%20projects/codegraph/)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

---

## Checklist A — Hypothesis coverage (S02-owned)

For **each** of H1 (partial), H5, H6, H7:

- [ ] PEER-CG §3 row exists with verdict ∈ {supported, weakened, rejected, inconclusive}
- [ ] Verdict aligns with INVESTIGATION-INDEX **verified if / rejected if** (not invented criteria)
- [ ] **Peer mechanism cite:** at least one `similar projects/codegraph/…` **file:line** per hypothesis (not README-only for explore, watch, or MCP surface claims)
- [ ] Trace contrast cites TRACE-AUDIT or Trace `file:line` where comparing shipped behavior
- [ ] H1 marked **partial** — no over-claim that S02 alone closes full H1 matrix (S04 owns aggregate)

**Spot-check (minimum 3 — pick H7, H5, H6):**

- [ ] Peer anchor resolves in repo (`tools.ts` handleExplore / DEFAULT_MCP_TOOLS / `watcher.ts` debounce)
- [ ] Verdict consistent with cited mechanism (e.g. H7 rejected only if Trace compose proven equivalent with evidence)

---

## Checklist B — Peer mechanism gate (mandatory)

Major CG mechanism claims **must not** be README-only.

| Claim class | Required peer evidence |
|-------------|------------------------|
| `codegraph_explore` inputs/outputs | `src/mcp/tools.ts` schema + `handleExplore` (or dispatch case) file:line |
| Blast radius | `buildBlastRadiusSection` or equivalent — file:line |
| Single-tool default surface | `DEFAULT_MCP_TOOLS` / `getStaticTools` — file:line |
| Index + watch | `watcher.ts` + init/sync path or README **plus** code cite for debounce/sync |
| Output caps / progressive context | `MAX_OUTPUT_LENGTH`, budget helpers — file:line |
| Daemon (anti-pattern section) | `daemon.ts` or MCP proxy path — file:line |

- [ ] §6 Evidence appendix lists peer paths with line numbers for **T1, T2, T3, T7** classes
- [ ] Reviewer independently opens ≥2 peer cites (confirm line numbers still match)
- [ ] Optional live MCP (T8): if claimed, redacted capture in `$EV/` or marked skipped with reason

---

## Checklist C — Planner must-answer items (P38-S02-00)

Confirm PEER-CG addresses all five handoff questions:

- [ ] **Q1** Explore mechanism: inputs, outputs, blast radius (§1 + §2 row)
- [ ] **Q2** Index/watch vs Trace manual index (§2 + H5 verdict)
- [ ] **Q3** Single-tool UX vs Trace 16-tool discovery implications (§2 + H6 verdict — observation, not product decision)
- [ ] **Q4** P24 transfer (`trace_explore` / MCP consolidation) still **deferred** — cite EXTERNAL-RESEARCH + live grep or MCP inventory
- [ ] **Q5** Anti-patterns Trace must not copy: Law 6/7, daemon-as-P0 (§4)

---

## Checklist D — Hygiene & phase law

- [ ] No Go/TS/web product diff in investigate row
- [ ] No ranked build list / REMEDIATION-PLAN language
- [ ] §4 Anti-patterns include daemon + full-dump defaults + “graph-only replaces moat” rejection
- [ ] §5 Moat row present (CG lacks tasks/gates/evidence)
- [ ] §6 Spawn list empty **or** each item has trigger + owner scope
- [ ] Evidence under `experiments/runs/…-p38-s02-654/evidence/` with date + row id
- [ ] Benchmark/README numbers labeled observation if not re-run (T6)

---

## Checklist E — Saturation forward-fit (DESIGN-LOCKS)

S02 peer audit must feed S04 matrix + S05 gate:

- [ ] S02-owned H* verdicts evidence-backed enough for GAP-REGISTRY (no “TBD” without inconclusive + plan)
- [ ] Peer mechanism cites satisfy DESIGN-LOCKS “not README-only” for CG
- [ ] H1 partial boundary preserved — UA/GF leg still S03
- [ ] P24 deferral explicit so S06 does not treat CG transfer as silently shipped

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
| Blocker / high | Spawn `P38-S02-01a` implement + `P38-S02-01b` review **below this row**, OR fix ≤5 lines inline if trivial cite gap |
| Medium | Prefer spawn unless single-table row fix in PEER-CG |
| Low / nit | Note only; may APPROVE with listed residuals |

**SPAWN** when CG slice still shallow (e.g. explore handler not cited, watch mechanism README-only, H7 P24 deferral missing).

---

## Exit criteria

- [ ] Verdict + confidence in board Notes with checklist summary
- [ ] On APPROVE: no open blocker/high without pending follow-up row
- [ ] On APPROVE: next runnable **P38-S03-00** (or spawned S02 row first)

## Next

`P38-S03-00` on APPROVE (or spawned S02-01a first)
