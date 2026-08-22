# Scope 05 — board map

**S05 VERIFY + handoff.** Serial: **P34-S05-00 → P34-S05-01 → P34-S05-02**. After S04-02 PASS.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 607 | P34-S05-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock VERIFY floor (PLAN T9; L1–L4) |
| 608 | P34-S05-01 | [01-verify.md](01-verify.md) | Verify | `VERIFY-NOTES.md` + evidence dir; L1–L4; DR stays OPEN |
| 609 | P34-S05-02 | [02-dr-handoff.md](02-dr-handoff.md) | Close handoff | Close DR-HANDOFF; successor lean (**no successor** default) |

## Locked by S05-00 (2026-08-21)

| Lock | Value |
|------|-------|
| VERIFY = PLAN **T9** | T1 real SPA + T4/T5 concurrent auto-port + T8 docs + T10 stub-fail |
| VERIFY artifact | `scopes/scope-05-verify/VERIFY-NOTES.md` |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p34-s05-01-verify/evidence/` |
| Close owner | **P34-S05-02** (S05-01 leaves DR-HANDOFF **OPEN**) |
| Successor lean | Default **`no successor`** |
| L1–L4 | Must tick in VERIFY-NOTES (`.trace/` only; real SPA; auto-port; one root/process) |
| Port range | `7432`–`7441` max 10; `--addr` pin-strict (`fs.Visit` / AddrExplicit) |
| Stub-fail | Shipped `embeddist/index.html` must not contain `Embedded GUI stub`; markers `#root` / `/assets/` |
| Product code in S05 | **Forbidden** (verify/handoff notes only) |

## Carried from S02–S04 reviews (PASS)

| Item | Disposition for VERIFY |
|------|------------------------|
| T1–T3 static embed | Re-run focused tests + embeddist marker grep |
| T4/T7/T11 httpapi auto-port | Re-run `TestListenAutoPort_*` |
| T5/T6 concurrent + pin | Re-run `TestGuiServeConcurrent*` + explicit busy |
| T8 docs | Re-grep quickstart/help/AGENTS/web.README/embeddist |
| Contributor `web/` DX | Residual OK if labeled |
| Optional CI embed-gui | Out of phase — do not fail |

## Command floor (S05-01) — summary

0. Evidence dir + metadata
1. Static embed tests + embeddist stub-fail markers (T1–T3/T10)
2. `TestListenAutoPort_*` (T4/T7/T11)
3. Concurrent + explicit pin CLI tests (T5/T6)
4. `go test ./internal/httpapi/ ./cmd/trace/ -p 1` + gui/serve `--help`
5. Docs T8 greps (forbidden + positive)
6. Optional live consumer-temp smoke **or** waive
7. Residuals + S00–S04 aggregate → handoff

## Out of this scope

- New product features (spawn remediation rows instead).
- Re-implementing embed pipeline / listen hop / docs rewrite.
- Closing DR-HANDOFF from S05-01 (owner is S05-02).
- Inventing Phase 35 unless thin follow-on exception fires.
- Explore craft redesign; hosted SaaS; public bind default; SPA under consumer `.trace/`.
