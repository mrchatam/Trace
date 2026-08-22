# S03 — Phase VERIFY / peer-honesty closeout — scope todos

**Depends-on:** P12-S02-02 done. Owns S01+S02 named regressions + carry-forward + DR-HANDOFF.

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | 00-PLANNER | scope planner | **done** — FINAL locks + DR-HANDOFF = `no successor` |
| 2 | 01-verify | verify | **done** — PASS; VERIFY-NOTES + DR-HANDOFF start = `no successor` |
| 3 | 02-scope-review | review | **done** — APPROVE; DR-HANDOFF closed = `no successor`; Phase 12 complete |

## Locked evidence imports (P12-S03-00)

| Scope | Named regressions (must re-prove) |
|-------|-----------------------------------|
| S01 | `TestImportProvenanceRoundTrip` (store); `TestAnalyzerImportProvenanceExtracted` (analyzers, CGO1); `TestExpandImportEdgeProvenance` / `TestWhySurfacesEdgeProvenance` (retrieval); `TestContextWhyTraceEdgeProvenance` (compiler) |
| S02 | `TestBudgetLoudTotals`; `TestCandidateCapSetsTruncated`; `TestIndexStaleBanner` (compiler) — keep S01 `TestContextWhyTraceEdgeProvenance` green |

## Carry-forward
Honesty A/B/C+G; Gate E; Gate F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` N=3; Phase 11 DF surfaces via product pkgs; product `./...` (graphify space FAIL OK; CGO0 analyzers OK residual).

## Reminders
- Independent re-prove S01–S02 + carry-forward gates; close DR-HANDOFF (default **no successor**)
- Carry-forward gates stay green; Gate C `dry_run:false` untouched
- Dry-run ≠ Gate C / F / G / ablation / H / checklist
- Do **not** auto-board research impact/install/supersession (ranks 4+)
- Forward-only board; implementers: status + Notes only
- Spawn on fail: `P12-S03-01a` / `01b` / (`01c`) immediately below
- **DR-HANDOFF:** default **`no successor`** — **S03-01 starts** Notes; **S03-02 owns completion**
- Residuals OK into VERIFY: no provenance enum CHECK; synthetic context fixture; stale test does not pin exact lex-first-8 set
- Next after APPROVE: **(DR-HANDOFF close)** — no Phase 13 unless promoted
