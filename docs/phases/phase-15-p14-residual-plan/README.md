# Phase 15 — P14 residual remediation plan (thin)

**Status:** **COMPLETE** (2026-08-17) — S01 MCP Assert **APPROVE high**; S02 VERIFY PASS + review **APPROVE**. **DR-HANDOFF closed = `no successor`** (no Phase 16 / S05 / plan simulate / D21+). Human-scheduled forward after Phase 14 (`no successor`) to disposition P14 VERIFY residuals.

## Why this phase exists

Phase 14 closed green with Notes-only residuals and historical DR-HANDOFF `no successor`. This phase is a **forward human reopen** to disposition those residuals — not a promotion of goals #2–#4 (S05 / `plan simulate` / D21+).

## Disposition matrix (FINAL — P15-00)

| ID | Residual | Disposition | One-line rationale | Blast radius |
|----|----------|-------------|--------------------|--------------|
| R1 | `AssertToolAllowed` library/CLI only — not on MCP dispatch | **fix** | MCP is an agent surface; fail-closed Assert with slug `mcp:<toolName>` closes the VERIFY honesty gap without new tools | `internal/mcp` (+ tests); reuse domain Assert / mig 013; **no** new MCP tools; **no** new mig |
| R2 | `allowContainsOut` late-upgrade may skip re-enqueue | **defer** | Rare multi-path incomplete contains-OUT; Gate F + primary asymmetry tests green | Notes / FUTURE only — no P15 implement |
| R3 | `similar projects/graphify` space-in-path FAIL on `go test ./...` | **wontfix** | Non-product reference clones; product bar is `./cmd\|internal\|evals` | Ops tip only: nest local `go.mod` or keep clones outside the Trace module |
| R4 | CGO0 analyzers FAIL OK | **wontfix** | tree-sitter bindings need CGO; product bar is CGO1 | None |

**Explicitly out:** research S05 / `plan simulate` / D21+ / ranks 7+ unless separately promoted.

## Scope order (FINAL)

| Scope | Focus | Outcome |
|-------|--------|---------|
| S00 / phase planner | Inventory + disposition matrix + spawn | **done** (P15-00) |
| S01 | Wire `AssertToolAllowed` into MCP tool dispatch (R1) | **APPROVE high** (P15-S01-02) |
| S02 | Phase VERIFY + DR-HANDOFF (`no successor`) | **PASS** / handoff **closed** = `no successor` |

## Out of scope unless promoted

- Fixing R2 / R3 / R4 in this phase
- Full MCP surface expansion / new tools / install·decide MCP dump
- Rewriting Phase 00–14 `done` history
- Re-opening closed DF-60…67 as if undone
- S05 supersession / `plan simulate` / D21+ / ranks 7+

## Assumptions (P15-00; unattended with prior orchestrator approval)

1. Disposition trades are not architecture-blocking beyond R1’s existing `mcp:<name>` slug convention.
2. Builtin MCP tools remain AUTO_ALLOWED on first resolve; Assert wire-up must not break the default nine-tool + `trace_version` path.
3. R2 stays deferred (no harden scope this phase).
4. R3/R4 are wontfix for the Trace product bar.
5. After S01+S02, default DR-HANDOFF = `no successor` unless VERIFY Notes explicitly promote.
6. Goals #2–#4 stay off-board.

## Parallel track (not board-blocking)

Optional dogfood under `experiments/`; feed new DF-* forward only.
