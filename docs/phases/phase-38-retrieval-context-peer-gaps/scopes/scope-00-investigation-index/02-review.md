# P38-S00-02 — Review investigation index

## Metadata
- id: P38-S00-02
- todo_ids: [P38-S00-02]
- role: reviewer
- skills: [code-review-and-quality, research]
- mcps: [user-trace]
- verification: manual

## Objective

Independent review of **`INVESTIGATION-INDEX.md`** vs INTAKE, DESIGN-LOCKS, PEER-FIXTURES, and planner locks from P38-S00-00. **APPROVE** or **REQUEST CHANGES** (spawn S00-01a fix row only if critical gaps).

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [INTAKE.md](../../INTAKE.md), [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md), [PEER-FIXTURES.md](../../PEER-FIXTURES.md)
- [01-investigate.md](01-investigate.md) — locked hypothesis register (authority)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). **Fresh subagent** — do not share implementer session.

---

## Checklist A — Hypothesis register completeness

For **each H1–H11**, confirm INVESTIGATION-INDEX §2 row includes:

- [ ] Statement matches INTAKE (no scope creep)
- [ ] Peers listed (primary: CG, UA, GF, CM, OH, SWE, Aider as applicable)
- [ ] Trace touch areas with repo paths (not vague package names)
- [ ] Method: Trace live / peer read / both
- [ ] Tools column (T / CG / GF / UA / —)
- [ ] Peer cite targets: resolvable paths under `similar projects/` or P24/web URLs
- [ ] **Verified criteria** — falsifiable ("live X shows Y")
- [ ] **Rejected criteria** — falsifiable ("gap closed if Z")
- [ ] Owner scope (S01 / S02 / S03 / S04 / combo)

Spot-check (minimum 3 hypotheses — pick H1, H5, H7):
- [ ] Peer file exists at cited path (Read or Glob)
- [ ] Trace anchor file exists (`internal/compiler/doc.go`, `internal/mcp/server.go`, etc.)

---

## Checklist B — Saturation prerequisites (DESIGN-LOCKS forward-fit)

Index must **prepare** S01–S05 to satisfy DESIGN-LOCKS saturation exit. Confirm INVESTIGATION-INDEX enables:

| DESIGN-LOCKS criterion | Index must provide |
|------------------------|-------------------|
| Every H1–H11 verified/rejected/deferred with evidence pointer | §2 evidence target column names `experiments/runs/…` paths per scope |
| Live Trace command per major gap claim | §4 lists CLI/MCP commands; S01 owns live captures |
| Each primary peer (CG, UA, GF) mechanism cite (not README-only) | §2 peer cite targets include source files (context-builder.ts, explore impl, graph.html/worked) |
| Cross-matrix Trace strengths (moat row) | §3 routes H10 + moat to S04; H10 owner includes S04 |
| Spawn list empty or deferred with trigger | §5 spawn rules + H12+ template |
| Reviewer confidence high before S06 | §1 links S05 gate; §6 non-goals block premature REMEDIATION-PLAN |

- [ ] All six rows above satisfied by index content

---

## Checklist C — Investigation hygiene

- [ ] No implement tasks disguised as investigation ("implement trace_explore in P38")
- [ ] §6 non-goals explicit: no product code, no GAP-REGISTRY, no REMEDIATION-PLAN, no build ranking
- [ ] P24 cited where relevant (H7, H10, transfer rows); not duplicated verbatim
- [ ] Evidence convention matches PEER-FIXTURES + 01-investigate locked paths
- [ ] Spawn rules align with DESIGN-LOCKS investigation row rules (dedicated row, not silent backlog)
- [ ] H5 notes lang-count reconciliation (INTAKE "3 langs" vs 4 shipped) — S01 must resolve
- [ ] H1 split across S01+S02+S03 documented; S04 owns matrix
- [ ] H11 owned by S04 (docs slice), not S01 alone

---

## Checklist D — Peer map & ownership

- [ ] §3 peer map matches SCOPE-TODOS routing table
- [ ] No orphan hypotheses (every H has owner)
- [ ] S01: H2, H3, H5, H6, H9, H10 (+ H8 partial)
- [ ] S02: H1 partial, H5, H6, H7
- [ ] S03: H1 partial, H4, H8
- [ ] S04: all (matrix), H11 primary
- [ ] S05: saturation over all

---

## Checklist E — Optional tools per hypothesis

Confirm §2 or §4 documents optional tooling from PEER-FIXTURES:

| Tool | Used for hypotheses |
|------|---------------------|
| Trace CLI (`trace context`, `trace search`, `trace why`) | H1–H3, H5–H6, H9–H10 |
| Trace MCP (`trace_context`, `trace_search`, etc.) | Same |
| Codegraph MCP `codegraph_explore` | H1, H5, H7 (if `.codegraph/` exists) |
| Graphify worked examples | H4, H8 |
| UA context-builder read | H1, H2 (partial) |

- [ ] Tool assignments present; optional tools marked optional (not required for P38 exit)

---

## Verdict

Document in board Notes:

```text
Verdict: APPROVE | REQUEST CHANGES
Confidence: high | medium | low
Findings: blocker | high | medium | low | nit counts
Residual risks: (if medium confidence)
```

| Severity | Action |
|----------|--------|
| Blocker / high | Spawn `P38-S00-01a` implement + `P38-S00-01b` review **below this row**, OR fix ≤5 lines inline if trivial doc gap |
| Medium | Prefer spawn unless single-table row fix |
| Low / nit | Note only; may APPROVE with listed residuals |

If **REQUEST CHANGES**: list specific §/H id fixes; do **not** implement in review row unless trivial doc fix ≤5 lines.

---

## Exit criteria

- [ ] Verdict + confidence in board Notes with evidence (checklist ticks summarized)
- [ ] On APPROVE: no open blocker/high without pending follow-up row
- [ ] On APPROVE: next runnable **P38-S01-00**

## Next

`P38-S01-00` (on APPROVE)
