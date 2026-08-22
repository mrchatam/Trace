# P27-S03-00 — Scope planner (VERIFY)

## Metadata
- id: P27-S03-00
- todo_ids: [P27-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Lock verify floor and DR-HANDOFF close policy for Phase 27 after S00–S02 APPROVE/done. No product feature work beyond thickening VERIFY/DR-HANDOFF prompts; harness `--enforce` upgrade is VERIFY work owned by **S03-01**.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- [Phase 27 README](../../README.md)
- Pattern: [phase-26 scope-05-verify](../../../phase-26-loop-implementation/scopes/scope-05-verify/)
- Prior: S01 protocol v2 (warn-only `--strict`); S02 product honesty (thin blocked on `--strict --enforce`); harness enforce deferred here

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended orchestrator runs may execute after prior prompt approval.

## Locked defaults (FINAL — S03-00)

| Item | Value |
|------|-------|
| Precondition | P27-S01-02 + P27-S02-02 both `done` (APPROVE) |
| Binary | Rebuild `bin/trace` from repo HEAD before harness/product demos |
| Unit tests | `go test ./internal/... -count=1` — **must PASS** |
| Cmd honesty tests | `go test ./cmd/trace/... -run 'SeedExport\|Enforce\|Strict' -count=1` — **must PASS** |
| **Harness `--enforce` upgrade** | **S03-01 owns** — change `score.sh` T02 from warn-only `--strict` to **`--strict --enforce` on both arms** (`build` and `directed`) |
| Expected build arm | **P25-3a FAIL** on thin graph (expected); enforce may change G2 **WARN→FAIL** on honesty |
| Expected directed arm | Label **P25-3b**; PASS only if rich graph (Session-B dogfood) — thin FAIL OK for phase close |
| Product thin demo | P26-equivalent thin → `--strict --enforce` blocked; `--strict` alone warns + writes |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p27-s03-01-verify/` (+ `evidence/` subdir) |
| Notes artifact | `scopes/scope-03-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** until S03-02 |
| Successor | Decided **only** on P27-S03-02 (never TBD) |
| Default successor | **`no successor`** when VERIFY PASS (INT-07/08/10 delivered) |
| Alt successor | **Phase 28** only if regressions / blocking residual needing new phase work / human promotes |
| Human-gated | P25-4 attestation; optional Session-B dogfood — **non-blocking** residuals |
| Residual note | S02-02 low: BLOCKING store check can duplicate orphan-discovery message — document, do not FAIL VERIFY |

### Harness enforce policy detail (AUDIT option B)

S01 locked warn-only until S02 honesty existed. S02 shipped product rules and deferred harness. **S03-01** must:

1. Edit `experiments/ab-p25-gap-pass-validation/score.sh` T02: run `seed export … --strict --enforce`.
2. On non-zero exit / honesty stderr: **FAIL G2** (or equivalent named check) — do **not** leave as WARN-only.
3. Apply on **both** `--arm build` and `--arm directed` (not directed-only).
4. Keep FM-07 **warn-only** (git SHA drift never fails G2 alone).
5. Keep P25-3a/3b graph-count checks (fail on thin) — enforce is additive honesty alignment.

## Planner gate

- [x] `01-verify.md` runnable with locked defaults
- [x] `02-dr-handoff.md` includes successor decision table
- [x] `SCOPE-TODOS.md` current

## Exit criteria

- [x] Verify floor + enforce policy locked
- [x] Sibling prompts thickened
- [x] Board Notes → next **P27-S03-01**

## Next

`P27-S03-01`
