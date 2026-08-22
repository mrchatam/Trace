# Phase 01 — S0 completion & X0 readiness

## Goal

Finish the vertical slice beyond P0-X so **agent** baseline-vs-graph Experiment X0 can run via CLI: real Claim→Evidence→Review promotion (replace DONE stub-only), an honesty demo, an `evals/x0` harness (B0 vs G1), and MCP only **after** S01–S03.

P0 is closed (P0-X **7/7**, 2026-08-15). **P0 close ≠ Gate C.** Gate C measurement lives in **Phase 02** after this phase’s X0 harness exists.

## Prior phase outcomes (live — carry forward)

| Item | Live value |
|------|------------|
| Module | `github.com/mrchatam/Trace` (`go.mod` go 1.24.0) |
| Layout | `internal/{store,vcs,gitcli,analyzers,domain,retrieval,compiler}` + thin `cmd/trace` |
| Store | `.trace/trace.db`, `modernc.org/sqlite`, embed migrations `001`–`005` (`005` = `reviews.result`) |
| CLI | stdlib argv; G19 library-only; init/index/reindex/add/link/transition/seed/why/context/**review** |
| Analyzers | tree-sitter JS/TS + Python; file-local incremental |
| Fixture / P0-X | `fixtures/x0` + `evals/p0x` — **7/7 PASS** (keep green; do not replace with agent eval) |
| Claim/Evidence | Real path: `CreateClaim` / `CreateEvidence` / `LinkClaimEvidence` (`claim_has_evidence`); `claim.go` (was stub) |
| DONE policy | Linked Review `result=PASS` via `review_judges_task` **or** explicit `AllowDoneWithoutReview` / `--allow-done` / seed `allow_done`; **EvidenceIDs alone insufficient** |
| Review API | `CreateReview` / `SetReviewResult` / `LinkReviewTask` / `GetReview`; events `entity.created` + `review.result` |
| MCP / daemon / HTTP | Still absent (correct until S04) |

## Locked phase defaults (do not weaken)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| X0 conditions | **B0** = agent + ordinary repo tools; **G1** = agent + `trace` CLI `context`/`why` (may still read repo) |
| X0 corpus | Start with synthetic `fixtures/x0` (+ seed); no new OSS required this phase |
| X0 harness | `evals/x0` (new); keep `evals/p0x` unchanged and green |
| MCP | Thin adapter over library; **after** S01–S03; no business logic in MCP package |
| Daemon / HTTP | Still forbidden as primary surface |
| Embeddings | Still forbidden |
| Honesty demo | Prefer deterministic automated fail-closed demo under **`evals/honesty`** (Paths A/B/C); `verification: mixed` only if scripted claim is insufficient |
| Review policy | Every scope: `00-PLANNER` → `01` → `02-review` before next scope’s implement runs in board order |
| P0-X bar | Closed — do **not** reopen as Phase 01 exit; regression-keep green |
| Gate C | **Out** — Phase 02; this phase only needs X0 **dry-run metrics** for B0+G1 |
| DR-HANDOFF | Before S05 VERIFY final review is `done`, scaffold Phase 02 (planner + stubs + board) — do not repeat Phase 00→01 gap |

## Scope run order

| Scope | Theme | Board IDs | Folder |
|-------|--------|-----------|--------|
| S01 | Claim / Evidence / Review promotion (replace DONE stub) | P01-S01-00/01/02 | `scopes/scope-01-claim-review/` |
| S02 | Honesty demo (false claim rejected without evidence) | P01-S02-00/01/02 | `scopes/scope-02-honesty-demo/` |
| S03 | Agent X0 harness — B0 vs G1 via CLI (`evals/x0`) | P01-S03-00/01/02 | `scopes/scope-03-x0-harness/` |
| S04 | MCP adapter (thin; library-only; after S01–S03) | P01-S04-00/01/02 | `scopes/scope-04-mcp/` |
| S05 | Phase verify — X0 dry-run metrics; readiness for Gate C + Phase 02 handoff | P01-S05-00/01/02 | `scopes/scope-05-phase-verify/` |

Each scope folder has `00-PLANNER.md`, `01-*.md`, `02-scope-review.md`, and `SCOPE-TODOS.md`. Scope planners thicken `01-*`; this phase planner only light-locks order and cross-scope assumptions.

## Cross-scope blast radius

- **S01** changes DONE promotion → thickens S02 honesty expectations and CLI `transition` docs in **upcoming** S02/S03 only.
- **S03** harness must not break `evals/p0x` (keep both).
- **S04** must not fork domain logic into the MCP package.
- **S05** owns DR-HANDOFF for Phase 02 before VERIFY closes.

## Phase rules

- Run `00-PHASE-PLANNER` first, then scopes in order.
- Each scope: `00-PLANNER` → `01-implement` → `02-review` → spawns until confidence high.
- Do not start Phase 02 until S05 VERIFY is `done` **and** Phase 02 scaffold exists (DR-HANDOFF).
- Forward-only: do not rewrite Phase 00 `done` prompts; supersede behavior in new code + Notes.

## Out of scope (this phase)

- Gate C scoring / Go-No-Go report (Phase 02)
- Progressive planner / discovery replanning product (Phase 03+)
- Daemon / always-on HTTP server
- Embeddings / env graph / impact engine
- Rewriting P0-X deterministic harness as the agent eval (keep `evals/p0x`; add `evals/x0`)
- Multi-model adversarial review; VerifiedFact research suite beyond minimal promotion path

## References

- [`docs/init/A_PROJECT_PLAN.md`](../../init/A_PROJECT_PLAN.md) Phase 1–2
- [`docs/init/I_BENCHMARK_PLAN.md`](../../init/I_BENCHMARK_PLAN.md) Experiment X0
- [`docs/init/H_VERIFICATION_STRATEGY.md`](../../init/H_VERIFICATION_STRATEGY.md) authority matrix
- [`docs/init/D_DECISION_REGISTER.md`](../../init/D_DECISION_REGISTER.md) DR-AGENT, DR-SURFACE, DR-HANDOFF
- [`docs/TODO.md`](../../TODO.md)
- Phase 00 VERIFY: [`../phase-00-foundation/scopes/scope-09-phase-verify/VERIFY-NOTES.md`](../phase-00-foundation/scopes/scope-09-phase-verify/VERIFY-NOTES.md)
