# P38-S01-02 — Review Trace audit

## Metadata
- id: P38-S01-02
- todo_ids: [P38-S01-02]
- role: reviewer
- skills: [code-review-and-quality, research]
- mcps: [user-trace]
- verification: mixed

## Objective

Independent review of **`TRACE-AUDIT.md`** vs INVESTIGATION-INDEX, planner locks (P38-S01-00), and live repo evidence. **APPROVE**, **REQUEST CHANGES**, or **SPAWN** (new S01 investigate row). **Fresh subagent** — do not share implementer session.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INVESTIGATION-INDEX.md](../scope-00-investigation-index/INVESTIGATION-INDEX.md) — §2 verify/reject criteria
- [01-investigate.md](01-investigate.md) — investigation todos + locked defaults
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — saturation forward-fit

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

---

## Checklist A — Hypothesis coverage (S01-owned)

For **each** of H2, H3, H5, H6, H9, H10, H1 (partial), H8 (partial):

- [ ] TRACE-AUDIT §2 row exists with verdict ∈ {confirmed gap, not a gap, inconclusive}
- [ ] Verdict aligns with INVESTIGATION-INDEX **verified if / rejected if** (not invented criteria)
- [ ] Evidence cite: `file:line` **or** reproducible command **or** `$EV/` artifact path
- [ ] Partial hypotheses (H1, H8) explicitly scoped — no over-claim vs S02/S03

Spot-check (minimum 3 — pick **H2, H3, H6**):

- [ ] Re-run or read evidence file for cited command (e.g. `trace context`, MCP tool list test)
- [ ] Code anchor resolves (`compiler.go` FTS line, `RegisteredToolNames`, `doc.go` L7)

---

## Checklist B — Live evidence gate (mandatory)

Major claims **must not** be docs-only.

| Claim class | Required live evidence |
|-------------|------------------------|
| H2 FTS uses task title | `trace context` JSON **and** `compiler.go` cite; optional `trace search` contrast |
| H3 Layers 2–3 absent | Live packet JSON with max `layer` inspected |
| H5 Index manual/stale | `trace index status` output in `$EV/` or §3 log |
| H6 MCP surface / friction | `TestToolNamesRegistered` or MCP `trace_version` + ≥1 read tool invoked |
| H9 Intent pipeline absent | Grep output **or** explicit zero-match in retrieval/ |
| H10 Install moat | `trace install detect` output or install target doc cite |
| H1 partial | Live `trace_context` / CLI `trace context` capture |
| H8 partial | GUI route cite **or** optional screenshot |

- [ ] §3 Live command log lists commands with exit status
- [ ] At least **one** live Trace CLI **or** MCP call cited for **each** of H2, H3, H5, H6, H9, H10
- [ ] Reviewer independently spot-checks ≥2 commands (re-run or verify `$EV/` hash/timestamp)

---

## Checklist C — Planner must-answer items

Confirm TRACE-AUDIT addresses all six planner handoff questions:

- [ ] **Q1** Per-hypothesis verdict + evidence (§2)
- [ ] **Q2** Layer 0–1 vs 2–3 designed vs shipped (§4)
- [ ] **Q3** FTS/query: compiler input source documented (§4 — not just “search exists”)
- [ ] **Q4** MCP tool list + discovery friction notes (§2 H6 + optional table)
- [ ] **Q5** Index langs, manual vs auto, freshness (§2 H5)
- [ ] **Q6** Install/harness moat surfacing (§2 H10 + §5 non-gaps)

---

## Checklist D — Hygiene & phase law

- [ ] No Go/TS/web product diff in investigate row
- [ ] No ranked build list / REMEDIATION-PLAN language
- [ ] “Potential improvement (unranked)” only — clearly separated from verdicts
- [ ] §5 Non-gaps present (moat strengths — tasks, gates, evidence)
- [ ] §6 Spawn list empty **or** each item has trigger + owner scope
- [ ] Evidence under `experiments/runs/…-p38-s01-651/evidence/` with date + row id
- [ ] H5 lang count reconciled (INTAKE “3 langs” vs 5 ids / 4 families)

---

## Checklist E — Saturation forward-fit (DESIGN-LOCKS)

S01 audit must feed S05 gate:

- [ ] S01-owned H* verdicts evidence-backed enough for S04 matrix (no “TBD” without inconclusive + plan)
- [ ] Live Trace command trace exists per **major gap** claim
- [ ] Open questions routed to §6 spawn list, not silent backlog
- [ ] H1/H8 partial boundaries preserved for S02/S03

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
| Blocker / high | Spawn `P38-S01-01a` implement + `P38-S01-01b` review **below this row**, OR fix ≤5 lines inline if trivial doc gap |
| Medium | Prefer spawn unless single-table row fix in TRACE-AUDIT |
| Low / nit | Note only; may APPROVE with listed residuals |

**SPAWN** when audit slice unbounded (e.g. entire MCP handler audit) — scoped prompt referencing specific H* only.

---

## Exit criteria

- [ ] Verdict + confidence in board Notes with checklist summary
- [ ] On APPROVE: no open blocker/high without pending follow-up row
- [ ] On APPROVE: next runnable **P38-S02-00** (or spawned S01 row first)

## Next

`P38-S02-00` on APPROVE (or spawned S01-01a first)
