# P38-S04-00 — Scope planner (cross-matrix gap registry)

## Metadata
- id: P38-S04-00
- todo_ids: [P38-S04-00]
- role: planner
- skills: [planning-and-task-breakdown, analyst]
- verification: automated

## Objective

Lock S04: synthesize S01–S03 into **`GAP-REGISTRY.md`**. Single SoT for evidence-backed gaps, non-gaps, deferred ideas. **No REMEDIATION-PLAN yet. No product code.**

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- S01–S03 artifacts (required inputs)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Inputs (required)

- `TRACE-AUDIT.md`, `PEER-CG.md`, `PEER-UA-GF.md` (§1 UA · §2 GF · §3 MP), `INVESTIGATION-INDEX.md`

## Locked defaults

| Item | Value |
|------|-------|
| Artifact | `scopes/scope-04-gap-registry/GAP-REGISTRY.md` |
| Method | Synthesis of APPROVED S01–S03 + H11 doc read |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-s04-660/evidence/` |
| Gap IDs | G-001…G-011 (1:1 H*); moat M-001 |
| Matrix columns | Trace \| CG \| UA \| GF \| MP \| moat row |
| Verdicts | gap \| non-gap \| defer |
| Severity | Investigation confidence only — not S06 build priority |
| Dual-side evidence | Required per gap row (Trace + peer cite) |
| Product edits | **Forbidden** |

## Matrix columns (planner lock — upcoming S04-01)

Include peer columns: **Trace** | **CG** | **UA** | **GF** | **MP** (Mempalace — human-added 2026-08-22) | moat row.

## Must answer for 01

1. Unified gap IDs (G-001…) linked to H*.
2. Severity: investigation confidence only (not build priority — that's S06).
3. Moat row: Trace strengths peers lack.
4. Spawn triggers for S05 if registry incomplete.

## Planner gate

- [x] `01-investigate.md` matrix columns locked (Trace | CG | UA | GF | MP | moat; T0–T12)
- [x] `02-review.md` requires dual-side evidence per gap row (Checklists A–E)
- [x] SCOPE-TODOS IDs 659–661 + G-001…G-011 registry

## Exit criteria

- [x] S04-01/02 prompts runnable alone (spot-check: S01–S03 verdict seeds align G-ID map)
- [x] Board `P38-S04-00` → `done`

## Next

`P38-S04-01`
