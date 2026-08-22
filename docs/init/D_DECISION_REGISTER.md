# D — Decision Register

Classification: **SETTLED** | **PROVISIONAL** | **EXPERIMENTAL** | **OPEN**

Inherited design decisions D1–D14 from `docs/DESIGN_DECISIONS.md` remain in force unless superseded here.

| ID | Decision | Class | Notes |
|----|----------|-------|-------|
| D1 | Product identity: knowledge/causal graph + progressive planning, not orchestrator-first | SETTLED | |
| D2 | Git canonical for file/history content | SETTLED | |
| D3 | Graph is logical; SQLite/relational OK | SETTLED | |
| D4 | Hybrid retrieval | SETTLED | Semantic channel deferred in slice |
| D5 | Progressive context | SETTLED | |
| D6 | Agent claims not authoritative | SETTLED | |
| D7 | Multi-layer reviews | SETTLED | Only todo-level required in first slice |
| D8 | Discovery first-class | SETTLED | |
| D9 | User decisions first-class | SETTLED | |
| D10 | Advise; human authoritative | SETTLED | |
| D11 | Environment part of planning | SETTLED | Deferred from first slice implementation |
| D12 | Worktrees for concurrency | SETTLED | Deferred entirely until single-agent proven |
| D13 | Apache-2.0 open core | SETTLED | |
| D14 | Benchmark before scale | SETTLED | |
| DR-SLICE | First implementation is vertical experimental slice (Git index + structural graph + work/causal subset + exact/lexical/graph retrieval + context compiler + CLI + honesty review), not full roadmap | SETTLED | |
| DR-NOSSEM | No embedding/vector index in P0-X / X0 | PROVISIONAL | |
| DR-NOENV | Environment graph not required for P0-X / X0 | PROVISIONAL | |
| DR-NOIMP | Automated decision-impact engine deferred | PROVISIONAL | |
| DR-EVT | Thin append-only event log + materialized tables | PROVISIONAL | |
| DR-CLAIM | Explicit Claim entity | PROVISIONAL | |
| DR-CHURN | Plan-affecting churn controls | PROVISIONAL | Default N=5 |
| DR-FAIL | After 3 failed fix loops → cause investigation | PROVISIONAL | |
| DR-LANG | Implementation language = **Go** | SETTLED | User 2026-08-15 |
| DR-NAME | Product/CLI name = **`trace`** | SETTLED | User 2026-08-15 |
| DR-SURFACE | **Library + CLI only** until core validated; no daemon/HTTP in P0-X | SETTLED | User Q2=c |
| DR-AGENT | **CLI for P0-X/X0**; MCP only after query/context API validated via CLI | SETTLED | User Q4=b; supersedes MCP-first recommendation |
| DR-API | Canonical **Go library API** is the source of truth; CLI (then later MCP/HTTP) are adapters | SETTLED | **Supersedes** earlier “HTTP daemon first” reading of ARCHITECTURE for early phases; daemon remains a *later* deployment shape, not a v0 requirement |
| DR-BENCH | Synthetic fixture corpus first | SETTLED | User Q5=a |
| DR-SEED | Human-curated ground truth for scoring; agent seeding measured separately later | SETTLED | User Q6=a |
| DR-P0 | P0 does **not** close on docs alone; requires tiny experiment **P0-X** | SETTLED | User Q7=b |
| DR-ANLANG | First analyzers: TypeScript/JavaScript + Python | PROVISIONAL | |
| DR-DB | One SQLite file per bound project under `.trace/` | SETTLED | Aligns with DR-TRACEDIR |
| DR-PACK | Context packet: JSON (canonical) + Markdown render | PROVISIONAL | |
| DR-LICENSE | Apache-2.0 already present | SETTLED | |
| DR-GOMOD | Go module path = **`github.com/mrchatam/Trace`** | SETTLED | User Round 2 Q1=b |
| DR-TRACEDIR | Project store = **`.trace/`** in bound repo, gitignored | SETTLED | User Round 2 Q2=a; caches may split later |
| DR-P0X | P0-X pass = **7-point bar** (causal + files/symbols/imports + why + bounded context + GT match + several deterministic queries + **incremental file update**) | SETTLED | User Round 2 Q3=c |
| DR-INCREMENTAL | Ordinary edits must not require full-repo reanalysis; **full-rebuild-on-any-change is a forbidden default architecture**; P0-X #7 must demonstrate localized update | SETTLED | Aligns with G12 / DR-P0X #7; referenced historically on B_INITIAL_BOARD |
| DR-PARSE | Analyzer backend = **tree-sitter** from the beginning | SETTLED | User Round 2 Q4=a |
| DR-GIT | Git via **CLI subprocess** behind **VCS adapter interface** | SETTLED | User Round 2 Q5=a |
| DR-RISK | Early roadmap optimizes for **reducing wrong-product risk**, not implementation speed; no throwaway foundations | SETTLED | User directive Round 2 |
| DR-HARNESS | Perfect Planner–style board in-repo; Trace policies in agent-loop-protocol | SETTLED | 2026-08-15 |
| DR-IMPL-BOARD | Implementers may only set status + notes on the board | SETTLED | User |
| DR-FORWARD | done prompts/rows immutable (best-effort); fixes spawn forward | SETTLED | User |
| DR-REVIEW-LOOP | Implement→review per scope with spawn loops until high confidence | SETTLED | User |
| DR-HANDOFF | Closing phase VERIFY must scaffold next phase (planner + stub scopes + board) or record `no successor` | SETTLED | 2026-08-16; Phase 00 missed this — corrected by adding P01 planner after the fact |

## Supersession rule

Settled decisions change only via an explicit Decision record (once the system exists) or an update to this register with rationale during pre-implementation planning.
