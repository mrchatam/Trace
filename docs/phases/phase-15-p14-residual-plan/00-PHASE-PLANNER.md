# P15-00 — Plan P14 residual remediation (FINAL)

## Metadata
- id: P15-00
- todo_ids: [P15-00]
- role: planner
- skills: [planning-and-task-breakdown, ask-questions-if-underspecified]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live repo + P14 VERIFY/REVIEW notes, produce a **FINAL disposition** for each P14 residual (fix / defer / wontfix) and, where fix is chosen, scaffold implement+review board rows + prompts. **No product Go in this row.**

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md)
- [phase README](README.md)
- P14 evidence:
  - [../phase-14-peer-impact-install-gates/scopes/scope-03-phase-verify/VERIFY-NOTES.md](../phase-14-peer-impact-install-gates/scopes/scope-03-phase-verify/VERIFY-NOTES.md)
  - [../phase-14-peer-impact-install-gates/scopes/scope-03-phase-verify/REVIEW-NOTES.md](../phase-14-peer-impact-install-gates/scopes/scope-03-phase-verify/REVIEW-NOTES.md)
  - [../phase-14-peer-impact-install-gates/scopes/scope-01-impact-walks/REVIEW-NOTES.md](../phase-14-peer-impact-install-gates/scopes/scope-01-impact-walks/REVIEW-NOTES.md)
  - [../phase-14-peer-impact-install-gates/DR-HANDOFF.md](../phase-14-peer-impact-install-gates/DR-HANDOFF.md)
- Live: `internal/retrieval/impact_walk.go`, `internal/domain/capability_decision.go`, MCP dispatch paths, `similar projects/graphify`
- [docs/TODO.md](../../TODO.md)

## Session start
Agent → clarify (grill if disposition trades are material) → Plan → execute (planner).

## Live residual confirmation (2026-08-17)

| ID | Still present? | Evidence |
|----|----------------|----------|
| R1 | **Yes** | `AssertToolAllowed` in `internal/domain/capability_decision.go`; **no** calls from `internal/mcp` (grep empty). Builtin slugs = `mcp:<toolName>` via `BuiltinMCPCapabilitySpecs`. |
| R2 | **Yes** | `impact_walk.go` ~110–125: flag upgrade without re-enqueue when `upgraded && existing.hop == d` fails on greater-hop rediscovery. |
| R3 | **Yes** | `go list ./...` → `malformed import path "…/similar projects/graphify/tests/fixtures"` (space). Product `./cmd\|internal\|evals` unaffected. |
| R4 | **Yes** | `CGO_ENABLED=0 go test ./internal/analyzers/...` → build failed (tree-sitter CGO bindings excluded). Product bar uses CGO1. |

## Disposition matrix (FINAL)

| ID | Disposition | Rationale | Blast radius |
|----|-------------|-----------|--------------|
| R1 | **fix** | Agent MCP surface should fail-closed like CLI; closes VERIFY honesty Note without new tools | `internal/mcp` (+ tests); slug `mcp:<Name>`; reuse domain Assert / mig 013 |
| R2 | **defer** | Rare multi-path edge; Gate F + primary asymmetry green | FUTURE Notes only |
| R3 | **wontfix** | Non-product reference clones; product bar package-scoped | Ops tip: nested `go.mod` or clones outside module |
| R4 | **wontfix** | Analyzers require CGO by design | None |

**Boarded:** S01 (R1 fix) → S02 VERIFY/close. **Not boarded:** R2/R3/R4 implement. **Out:** S05 / `plan simulate` / D21+ / ranks 7+.

## Planner work (this row)
1. [x] Re-read residual inventory R1–R4; confirm still present in live code
2. [x] For each residual: **fix** | **defer** | **wontfix** with one-line rationale + blast radius
3. [x] If any **fix**: create scope folders + `00-PLANNER`/`01`/`02` stubs + board rows immediately below this row (forward-only)
4. [x] Explicitly keep goals #2–#4 (S05 / plan simulate / D21+) **out**
5. [x] Update phase README + board Notes; mark this prompt **FINAL**

## Locked defaults (FINAL — phase)
| Item | Value |
|------|-------|
| Phase | Thin residual remediation — R1 fix only + VERIFY |
| History | Do not rewrite Phase 00–14 `done` prompts |
| Product Go | **Forbidden** on P15-00 |
| MCP | No new MCP tools; wire existing `AssertToolAllowed` into dispatch using `mcp:<toolName>` |
| DR-HANDOFF intent | After VERIFY: default **`no successor`** unless Notes promote |
| R2/R3/R4 | defer / wontfix / wontfix — no implement scopes |

## Exit criteria
- [x] Disposition matrix FINAL for R1–R4
- [x] Board updated (S01 + S02 rows spawned)
- [x] 00-PHASE-PLANNER marked **FINAL**
- [x] Notes point to next runnable row **P15-S01-00**

## Minimal todos
- [x] Inventory live residuals
- [x] Disposition matrix
- [x] Spawn S01 + S02
- [x] Board + README sync
