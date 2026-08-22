# Scope S05 — Phase 01 VERIFY

**Depends-on:** S01–S04 done (all APPROVED).

**Bar:** Honesty green; X0 dry-run B0+G1 metrics (`dry_run:true`); MCP documented (six tools, stdio, G19); p0x 7/7 regression. **Not** Gate C.

**Commands (VERIFY):**
- `CGO_ENABLED=0 go test ./evals/honesty/... -count=1`
- `CGO_ENABLED=1 go test ./evals/x0/... -count=1`
- `CGO_ENABLED=1 go test ./evals/p0x/... -count=1`
- `CGO_ENABLED=1 go test ./... -count=1`

**S04 MCP (live):** `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review` via `cmd/trace-mcp` (`go-sdk` v1.4.0 stdio). X0 stays CLI-only (MCP not required to re-run).

**DR-HANDOFF:**
- **S05-01** creates `docs/phases/phase-02-gate-c/` (README + `00-PHASE-PLANNER` + scope stubs) + board `P02-00`.
- **S05-02** owns completion check — refuse `done` until handoff runnable; may finish missing pieces.
- Do not leave Phase 02 README-only / blocked-until-noticed.

**Board rights:** VERIFY **may spawn** remediations (`01a`/`01b`/`01c`); S05-02 closes phase only after handoff complete.

- [x] P01-S05-00 planner — 2026-08-16: thickened 01-verify (commands, evidence table, MCP checklist, spawn/closeout) + 02-scope-review (DR-HANDOFF completion); SCOPE-TODOS synced; no product Go; **next P01-S05-01**
- [x] P01-S05-01 verify — 2026-08-16: Phase 01 VERIFY PASS (not Gate C); honesty+x0+p0x+`./...`+MCP checklist green; VERIFY-NOTES written; Phase 02 scaffold + `P02-00` board started; **next P01-S05-02**
- [x] P01-S05-02 review — 2026-08-16: APPROVE high; fresh suite re-run matches VERIFY-NOTES; DR-HANDOFF complete; Phase 01 complete; next **P02-00**; [REVIEW-NOTES.md](./REVIEW-NOTES.md)
