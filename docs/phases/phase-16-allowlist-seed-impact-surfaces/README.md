# Phase 16 — Allowlist, seed, and impact surfaces (post-P15 DFs)

**Status:** **SUPERSEDED** (2026-08-17) — do **not** run. Canonical: [`../phase-16-assert-root-and-surfaces/`](../phase-16-assert-root-and-surfaces/). See [`SUPERSEDED.md`](SUPERSEDED.md).

## Why this phase exists

Phase 15 closed MCP Assert (R1) with historical DR-HANDOFF `no successor`. Live dogfood + adversarial hunt then filed **DF-68**, **DF-70…74**, **DF-75…78** off-board. This phase is a **forward human reopen** to schedule those findings plus remaining deferred-product residuals that are still not closed/wontfix.

It does **not** reopen Phase 00–15 `done` history, does **not** board S05 supersession / `plan simulate` / D21+ ladder / ranks 7+, and does **not** introduce daemon/HTTP/embeddings.

## Disposition matrix (FINAL — P16-00)

| ID | Sev | Finding | Disposition | Scope |
|----|-----|---------|-------------|-------|
| **DF-76** | high | MCP `project=` + `store.Open` auto-mkdir → fresh AUTO_ALLOWED DB (DENIED bypass) | **fix** | S01 |
| **DF-75** | med | No CHECK on `capability_tool_decisions.decision`; YOLO stores; Resolve fall-through AUTO_ALLOWED | **fix** | S02 |
| **DF-77** | med | Allowlist is MCP-only; CLI `add`/`why` ignore DENIED | **fix** | S02 |
| **DF-78** | med | Unprefixed `decide --slug trace_why` does not gate `mcp:trace_why` | **fix** | S02 |
| **DF-68** | low | `install -C` uses process cwd for Claude CONDITIONAL marker | **fix** | S03 |
| **DF-22 / DF-37** | ops | Cursor MCP reload still manual (P11 closed product; residual ops) | **carry-forward** — keep reload tip; **no** PID kill | S03 lock |
| **DF-70** | med | Seed import rejects `discovery_mentions_task` | **fix** | S04 |
| **DF-73** | low | Seed cannot import impact findings/alternatives | **fix** | S04 |
| **DF-71** | med | `context` / `why` omit impact findings / `overall_class` | **fix** | S05 |
| **DF-72** | med | MCP has no impact tool (P14 A3 non-goal) | **fix** — thin `trace_impact` adapter; **supersedes** P14 “no MCP impact” for this tool only | S05 |
| **DF-74** | low | `impact report` JSON PascalCase vs tasks snake_case | **fix** | S05 |
| **DF-67** | low | Symbol-entity staleness out of `index_honesty` bar | **defer** (remain out of P16 bar; VERIFY reconfirm) | S06 residual |
| **DF-36** | method | Self-dogfood caveat | **off-board** | experiments only |
| **R2** | residual | `allowContainsOut` late-upgrade (P15 defer) | **defer** — not boarded | Notes |
| **R3 / R4** | ops | graphify space / CGO0 analyzers | **wontfix** (P15) | Notes |

Zero requested product DFs left unscheduled: DF-68, 70–78 are **fix**; DF-22/37 **carry-forward** in S03; DF-67 **defer** with explicit lock; DF-36 method-only.

## Scope order (FINAL)

| Scope | Focus | DF IDs |
|-------|--------|--------|
| S00 / phase planner | Inventory + locks + spawn | **done** (P16-00) |
| S01 | MCP project-root + auto-init / DENIED isolation | DF-76 |
| S02 | Allowlist CHECK + YOLO fail-closed + slug prefix + CLI parity | DF-75, DF-77, DF-78 |
| S03 | `install -C` vs cwd (+ DF-22/37 tip keepers) | DF-68; DF-22/37 residual |
| S04 | Seed import rels + impact findings | DF-70, DF-73 |
| S05 | context/why impact packet + snake_case JSON + thin MCP `trace_impact` | DF-71, DF-72, DF-74 |
| S06 | Phase VERIFY + DR-HANDOFF | all P16 named + carry-forward; DF-67 reconfirm |

## Out of scope unless promoted

- Daemon / always-on HTTP / embeddings / Neo4j / full-rebuild indexer
- New MCP **install** or **decide** tools; `trace_plan` / `trace_index` MCP
- Rewriting Phase 00–15 `done` history (P15 `no successor` stays historical)
- S05 supersession / `plan simulate` / D21+ / ranks 7+
- Re-opening closed DF-60…66 as if undone; claiming R2/R3/R4 fixed
- PID-kill / auto-reload of live Cursor MCP (DF-22/37)
- Symbol-entity honesty product bar (DF-67 remains deferred)

## Assumptions (P16-00; unattended with orchestrator approval)

1. Thin MCP tools already exist post-P10; adding `trace_impact` is G19 adapter work, not “MCP on the P0-X critical path.”
2. P14 A3 “no MCP impact” is **superseded for this thin tool only** — not a full impact/install/decide dump.
3. Per-store allowlist remains correct; DF-76 is auto-init / virgin `project=`, not “inherit DENIED across projects.”
4. CLI `capability decide` / `decisions` stay operator-ungated so DENIED cannot lock out allowlist repair.
5. After S01–S06, default DR-HANDOFF = **`no successor`** unless VERIFY Notes explicitly promote.
6. Compat ceiling may move **13 → 14** if S02 lands a CHECK migration.

## Parallel track (not board-blocking)

Optional dogfood under `experiments/` stays off-board except the DF IDs scheduled here.
