# S01 — Edge provenance — scope todos

**Depends-on:** P12-00 done. Owns research **rank 1** — structural import-edge `EXTRACTED|INFERRED|AMBIGUOUS` + Why/context `edge_provenance`.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks (mig 011 / imports.provenance) |
| 2 | 01-edge-provenance | implement | **pending** — next runnable **P12-S01-01** |
| 3 | 02-scope-review | review | pending |

## FINAL summary (P12-S01-00)
- Persist enum on **`imports.provenance`** (mig `011_import_edge_provenance.sql`); not causal `confidence`
- Analyzers: EXTRACTED (AST) / AMBIGUOUS (wildcard); INFERRED store-accepted + surfaced; call edges deferred
- Surface `edge_provenance` on Expand/Why + context WhyTrace (and structural Items)
- Named tests + carry-forward verify cmds in `01-edge-provenance.md`

## Reminders
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Forward-only board; implementers: status + Notes only
- No research S03+ (impact/install/supersession) in this scope
- Next after APPROVE: **P12-S02-00** (packets may cite `edge_provenance` if present)
