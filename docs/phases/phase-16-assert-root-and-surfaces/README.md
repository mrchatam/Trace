# Phase 16 — Assert root & surfaces (thin)

**Status:** **complete** (2026-08-17) — S01–S05 APPROVE (DF-76/75/78/77/68/70…74 + thin `trace_impact`). S06 VERIFY PASS + review APPROVE. **DR-HANDOFF closed = `no successor`.** Human-scheduled forward after Phase 15 (`no successor`) to close **all new post-P15 open findings** (severity-agnostic). Phase 17 also complete (independently queued — **not** this successor). Next runnable: **none** unless human promotes follow-on.

## Why this phase exists

Phase 14/15 closed green with historical DR-HANDOFF **`no successor`**. Post-close dogfood + adversarial hunt then filed open product DFs: MCP allowlist fail-open / `project=` auto-init (**DF-75…78**), install `-C` vs cwd (**DF-68**), and seed/impact packet polish (**DF-70…74**). This phase is a **forward human reopen** — not a rewrite of P15 close, not a promotion of goals #2–#4 (S05 / `plan simulate` / D21+ ladder).

Findings SoT: [`experiments/DOGFOOD-FINDINGS.md`](../../../experiments/DOGFOOD-FINDINGS.md). Sources: [`POST-P15-DOGFOOD.md`](../../../experiments/POST-P15-DOGFOOD.md), [`POST-P15-BUGHUNT.md`](../../../experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md), [`BATCH-D21-D23.md`](../../../experiments/BATCH-D21-D23.md).

## Disposition matrix (live — P16-00 + DF-72 forward)

P16-00 prompt history deferred DF-72. Live lock: **fix** in S05 — [`DF-72-FORWARD.md`](DF-72-FORWARD.md).

| ID | Residual / finding | Disposition | One-line rationale | Home |
|----|--------------------|-------------|--------------------|------|
| DF-76 | MCP `project=` auto-init virgin dir AUTO_ALLOWED | **fix** | Highest: CallTool must not mkdir a fresh AUTO_ALLOWED store; DENIED on bound root stays per-store | **S01** |
| DF-75 | No CHECK on `capability_tool_decisions`; YOLO → AUTO_ALLOWED | **fix** | Fail-open enum; same class as DF-64 | **S02** |
| DF-78 | Unprefixed `decide --slug trace_why` does not gate `mcp:trace_why` | **fix** | Canonicalize registered MCP names to `mcp:` prefix | **S02** |
| DF-77 | CLI ignores MCP allowlist | **fix** (dual-slug) | Fail-closed CLI via **`cli:`** slugs; MCP DENIED does **not** deny CLI (CLI-first) | **S03** |
| DF-68 | `install -C` checks cwd not `-C` | **fix** | Pass parsed `-C` root into CONDITIONAL marker detect | **S04** |
| DF-22 / DF-37 | Cursor MCP reload still manual | **carry-forward** | Keep print/write reload tip; **no** PID kill | **S04** keepers |
| DF-70 | seed import rejects `discovery_mentions_task` | **fix** | Wire existing DF-42 link API into seed switch | **S05** |
| DF-71 | context/why omit impact findings | **fix** | Compiler packet includes findings / `overall_class` (G19) | **S05** |
| DF-73 | seed cannot import impact findings | **fix** | Seed JSON v1 `findings`/`alternatives` via existing impact domain | **S05** |
| DF-74 | impact report JSON PascalCase | **fix** | snake_case like DF-32 / tasks JSON | **S05** |
| DF-72 | MCP has no impact tool | **fix** (thin adapter) | User-required-all-findings: G19 `trace_impact` only. **Supersedes** P14 A3 for this tool. See [`DF-72-FORWARD.md`](DF-72-FORWARD.md) | **S05** |
| DF-67 | Symbol-entity staleness out of `index_honesty` bar | **defer** | Remain out of P16 bar; VERIFY reconfirm | **S06** residual |
| DF-36 | Self-dogfood method caveat | **off-board** | Experiment-only | experiments |
| P14 R2 | `allowContainsOut` late-upgrade | **defer** (hold) | P15 disposition; not a new DF | Notes only |
| P15 R3/R4 | graphify space / CGO0 analyzers | **wontfix** (hold) | Product bar = `./cmd\|internal\|evals` CGO1 | Notes only |

## Scope order (locked at P16-00)

| Scope | Focus | DFs |
|-------|--------|-----|
| S00 / phase planner | Inventory + disposition + spawn | **done** (P16-00) |
| S01 | MCP project root / auto-init / DENIED isolation | DF-76 (**high**) |
| S02 | Tool-decision enum CHECK + slug prefix + YOLO fail-closed | DF-75, DF-78 |
| S03 | CLI vs MCP allowlist parity (`cli:` vs `mcp:`) | DF-77 |
| S04 | install `-C` vs cwd for CONDITIONAL markers | DF-68 |
| S05 | seed / impact packet + thin `trace_impact` | DF-70, DF-71, DF-72, DF-73, DF-74 |
| S06 | Phase VERIFY + DR-HANDOFF | named S01–S05 + carry-forward; DF-67 reconfirm |

## Out of scope unless promoted

- Product MCP daemon / HTTP / embeddings / new MCP **install** or **decide** tools; `trace_plan` / `trace_index` MCP
- P14 R2 `allowContainsOut`; P15 R3/R4
- Research S05 supersession / `plan simulate` / D21+ ladder / ranks 7+
- Rewriting Phase 00–15 `done` history (P15 historical `no successor` intact)
- Re-opening closed DF-60…66 or claiming P15 R1 Assert undone; DF-67 stays deferred
- Session-global DENY across all `project=` roots (per-store SoT **HOLD**; S01 closes **auto-init**, not isolation)
- PID-kill / auto-reload of live Cursor MCP (DF-22/37)

## Independently queued (not this phase — now complete)

- **Phase 17** portable graph git — [`../phase-17-portable-graph-git/`](../phase-17-portable-graph-git/) — human-scheduled after P16 rows; **complete** 2026-08-17 (DR-HANDOFF **`no successor`**); this phase’s VERIFY DR-HANDOFF remains **`no successor`**

## Assumptions (P16-00 + DF-72 forward)

1. Human cut is **all new post-P15 open DFs**; DF-72 is a **thin MCP adapter** (not a daemon). P16-00 prompt history deferred DF-72; upcoming S05 **supersedes** that defer — [`DF-72-FORWARD.md`](DF-72-FORWARD.md).
2. MCP auto-init fix is **adapter-side** (`openStore` / `OpenExisting`) — CLI `store.Open` mkdir for `trace init` stays.
3. Dual-slug CLI gating preserves CLI-first: denying `mcp:trace_add` does not block `trace add` (S03).
4. S02 CHECK (mig **014**, ceiling **14**) lands before S03 CLI Assert.
5. VERIFY default DR-HANDOFF = **`no successor`** unless Notes explicitly promote.
6. Goals #2–#4 stay off-board.
7. Thin MCP tools already exist post-P10; `trace_impact` is G19, not “MCP on the P0-X critical path.”

## Parallel track (not board-blocking)

Optional dogfood under `experiments/`; feed new DF-* **forward** only (next free **DF-79**).
