# INTAKE — Concept & index

**Promoted 2026-08-22** — P42-00 complete; human promotion satisfied.

## Trigger

Phase 41 delivered G8+G9 entry wave. [REMEDIATION-PLAN](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) ranks **G6 + G7** as Phase 42+ entry themes (secondary queue promoted).

## Human locks

| Lock | Value |
|------|-------|
| Phase 41 | **CLOSED** — G8+G9 delivered; do not reopen implement rows |
| Phase 42+ mode | **Implement** G6+G7 entry |
| M-001 moat | Progressive task packet, loop+gate+review, harness, `trace_why`, plan tree, Laws 6–7 — **non-negotiable** |
| G-004a vector | **Out** — permanent defer (DR-NOSSEM) |

## Entry themes

| Theme | GAP ids | Scope | Deliverable sketch |
|-------|---------|-------|-------------------|
| **G6** Non-semantic concept retrieval | G-004b | S00 | Graph-label/summary channel without vector |
| **G7** Index freshness & langs | G-005 | S01 | Analyzer/lang policy + index honesty |

## In scope (P42 default)

- G6: Graph-label retrieval channel under DR-NOSSEM; law review desk-check at S00-00 before S00-01 build
- G7: Language tier policy (5 langs frozen); git-hook primary freshness; optional foreground `trace index watch`

## Out of scope (P42 default)

- G-004a semantic/vector channel
- Hosted SaaS / daemon defaults
- Product dual-index default
- Rewriting Phase 41 artifacts
- LLM-based concept extraction
- Tier-2 language adapters (rust/java/…) without human-promoted board row

## Rejects preserved

G-004a vector, product dual-index default, bundled MCP, CG 1-tool-only facade, query-only replaces task packet, full-graph dump defaults, always-on daemon defaults.

## Evidence pointers (read-only)

| Artifact | Path |
|----------|------|
| REMEDIATION-PLAN §2 G6/G7 | [§2](../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md) |
| GAP-REGISTRY G-004b/G-005 | [GAP-REGISTRY.md](../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md) |
| Phase 41 VERIFY-NOTES | [VERIFY-NOTES.md](../phase-41-layers-intent/scopes/scope-02-verify/VERIFY-NOTES.md) |
| G8 shipped | `internal/compiler/compiler.go` — `MaxLayer`; `internal/retrieval/layer_enrich.go` |
| G9 shipped | `internal/retrieval/intent.go`; `docs/RETRIEVAL_AND_CONTEXT.md` §3 |
| Live langs (P42-00) | `internal/analyzers/detect.go` — 5 IDs; `language_adapter.go` static table |
| Index freshness (P42-00) | `internal/install/githook.go`; `cmd/trace/index_status.go` |
| Peer G6 pattern | REMEDIATION-PLAN — GF EXTRACTED/INFERRED edges; MP BM25 text leg |
| Peer G7 pattern | REMEDIATION-PLAN — CG watcher debounce study (not daemon stack) |

## Open questions (resolved at P42-00)

| # | Question | Resolution (P42-00) |
|---|----------|---------------------|
| 1 | G6 law review gate — desk-check vs live spike? | **Desk-check at S00-00** → `LAW-REVIEW-NOTES.md`; implement at S00-01 |
| 2 | G7 lang policy — which langs beyond go/js/ts/tsx/py? | **Tier-1 frozen (5 langs)**; Tier-2 defer per lang; Tier-3 path-only (.md/.json/.yaml/.toml) |
| 3 | G7 watch/hook — git-hook only vs optional file watcher? | **Git-hook primary** (`trace install git-hook`); optional **foreground** `trace index watch`; **no always-on daemon** |

## Next runnable

**P42-S00-00** (G6 scope planner + law review desk-check)
