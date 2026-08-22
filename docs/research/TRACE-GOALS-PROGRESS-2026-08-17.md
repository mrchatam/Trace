# Trace goals vs progress — 2026-08-17

Read-only snapshot for post–Phase 13 product roadmap planning. Originally **did not** schedule Phase 14; **human-scheduled 2026-08-17:** thin Phase 14 boarded as [`docs/phases/phase-14-peer-impact-install-gates/`](../phases/phase-14-peer-impact-install-gates/) from §4 step **#1** only (research ranks 4–6 / FUTURE S03–S04). Phase 13 (import-resolve / DF-60…67) is now **complete** and **excluded** from the remaining backlog below.

**Numbering note:** `docs/ROADMAP.md` uses conceptual **P0–P20**. The execution board uses **Phase 00–14**. They are *not* 1:1 (board Phase 13 ≠ roadmap P13 “plan simulate”; board Phase 14 ≠ roadmap P14).

---

## 1. Original north star

Trace is a **local-first, versioned project knowledge graph** for AI coding agents: not a Git/IDE/agent replacement, but the missing layer that records *why* code exists, which goals/decisions/tasks caused it, how changes affect other work, what the project believes (with provenance), how discoveries should replan future work, what evidence supports “done,” and which skills/tools/MCPs a task needs. Planning is **progressive refinement**, not exhaustive up-front prediction. The first validation question (README) is whether an agent can understand an unfamiliar repo, plan a bounded task, adapt on discovery, and decide better with the graph than with raw repository contents alone.

---

## 2. Capability / hypothesis matrix

Status key: **done** | **partial** | **not started** | **deferred**. Evidence cites board phases, eval gates, and dogfood IDs where relevant.

### Hypotheses (H1–H7)

| Goal | Status | Evidence |
|------|--------|----------|
| **H1 — Project understanding** | **done** (first slice) | P0-X 7/7 (`evals/p0x`); Gate C **Go** (G1 mean 0.800 > B0 0.000, `dry_run:false`, N=3 — Phase 02); dogfood D01/D03/D04b. Residual: richer corpora / Mode-B q3 / Gate C live refresh = polish, not core miss. |
| **H2 — Progressive planning** | **partial** | Phase 03: coarse goal→phase→scope + deep-plan current+1 (`internal/planner`, mig 006); CLI `trace plan`. Dogfood D02/D10/D11b/D20. **Gate D** (progressive vs static upfront) has **no** named comparative harness yet — feature exists, gate incomplete. |
| **H3 — Discovery-driven replanning** | **done** (mini) | Phase 03 S02: severity + `ApplyDiscoveryReplan` + churn N=5; **Gate E** green (`evals/replan` `TestPlantedDiscoveryReplan`); dogfood D05. Combo dogfood D21 still planned. |
| **H4 — Decision impact** | **partial** | Phase 05: impact classes/alternatives + `ImpactReport` + `trace impact`; **Gate F prelim** green (`evals/impact`, P/R=1.0 planted). Dogfood D06. Full commercial engine blocked (DR-NOIMP); multi-seed/contains walks = research ranks 6 (deferred); `plan simulate` = roadmap P13 **not started**. |
| **H5 — Evidence-driven review** | **partial → strong** | Phase 01 fail-closed DONE; Phase 04 scope review + residuals; **Gate G prelim** green (`TestHonestyEscapeRateGateGPrelim`); Phases 09–11 honesty/operator fixes; dogfood D07/D07b/D32. **VerifiedFact** promotion still **out**. Phase reviews as a third formal layer remain thinner than roadmap P9. |
| **H6 — Progressive context** | **partial** | Phase 00: Exact+FTS5+graph Expand≤2 + compiler budgets; Phase 12 packet honesty (staleness/truncation); dogfood D04→D04b. No embeddings; RRF over FTS+graph = research rank 14 deferred; surgical “explore packet” / skeletonization = ranks 7 deferred. Token-efficiency Gate-style ablation vs success not a named published gate. |
| **H7 — Capability-aware planning** | **partial** | Phase 06 capability catalog + packet required/missing + ablation green (`evals/capability`); Phase 10/11 transition gates; dogfood D08. Graduated allowlist audit + marker-gated multi-client install matrix = research ranks 4–5 deferred. |

### Roadmap pillars (product surfaces)

| Pillar | Status | Evidence |
|--------|--------|----------|
| **Local-first graph + `.trace/` SQLite** | **done** | Phase 00 store; Git CLI VCS; no cloud SoT. |
| **Git history substrate** | **done** | Phase 00 S03; Gate A operationalized via P0-X history queries. |
| **Structural code graph (TS/JS/Py + Go)** | **partial** | Phase 00 analyzers + Phase 07 Go adapter; file-local incremental. Import resolve honesty = **Phase 13** (in flight — not counted as remaining product backlog here). Edge provenance enum = Phase 12 S01. Call-resolve two-pass = research rank 8 deferred. |
| **Work / causal model** | **done** (core) | Goals/tasks/decisions/assumptions/discoveries/plan_changes/reviews/evidence + links; provenance ACTIVE\|STALE\|SUPERSEDED. |
| **Hybrid retrieval + context compiler** | **partial** | Exact+FTS+graph; Layer 0–1 packets; Why reasons. Semantic/embeddings channel explicitly **not** default (Law 13). |
| **Progressive planner** | **partial** | See H2/H3; horizon intentionally minimal (current+1). |
| **Honesty / review / evidence** | **partial** | See H5; multi-agent reviewer dogfood D31 planned. |
| **Decision impact** | **partial** | See H4. |
| **Environment / capability graph** | **partial** | See H7; skills/rules/hooks as rich ontology = intentionally thin. |
| **Indexing / performance** | **partial** | T0 ignore tiers; Gate H green for smoke/~1k/~10k (`evals/perf`); **100k/1M deferred**. GC-03/04 (Gate C follow-ups) still deferred historically unless re-promoted. |
| **CLI** | **done** (primary surface) | `init`/`index`/`reindex`/`add`/`link`/`transition`/`seed`/`why`/`context`/`review`/`plan`/`impact`/`capability`/`tasks`/`install` (+ polish from P09–P11). |
| **MCP** | **partial** | Thin stdio adapter (why/context/add/link/transition/review/tasks/capability*); G19 library-only. Full plan/impact/index MCP **deferred** by design. |
| **Install (Cursor)** | **partial** | Phase 09 `trace install cursor` print/`--write`; reload tips Phase 11. Multi-client registry / conditional matrix = research rank 4 deferred. |
| **Plugins / worktrees / hardening** | **partial** | Phase 08: versioned analyzer contribution surface, worktree bind, migrate/backup/local auth checklist (`evals/compat`). Not a swarm runtime. |
| **Dogfood → product loop** | **done** as process | Phases 09–11 closed DF backlog waves; Phase 12 peer honesty; Phase 13 owns DF-60…67. |

### Milestone gates (ROADMAP §3)

| Gate | Status | Notes |
|------|--------|-------|
| **A** history | **done** | P0-X / Git index joins. |
| **B** why | **done** | Seeded causal why + Gate C understanding lift. |
| **C** graph > raw | **done** | Phase 02 Go; dry-run ≠ pass. |
| **D** progressive > static | **not started** (as named gate) | Planner shipped; comparative eval missing. |
| **E** discovery→plan | **done** | Mini-eval green. |
| **F** impact | **partial** | Prelim planted P/R; not commercial engine. |
| **G** review honesty | **partial** | Prelim escape-rate; VerifiedFact out. |
| **H** large repos | **partial** | ≤~10k planted ladder; 100k/1M out. |

---

## 3. What remains (product backlog)

Ordered for roadmap thinking. **Excludes** DF-60…67 and any Phase 13 implementation work.

### (A) Deferred research ranks still unboarded

From [`SIMILAR-PROJECTS-REVIEW-2026-08-16.md`](SIMILAR-PROJECTS-REVIEW-2026-08-16.md) FUTURE S03–S06 (ranks **4–20** product-accept lane; Phase 12 only boarded ranks 1–3):

| Priority band | Items (research rank) |
|---------------|------------------------|
| High product value | **4** marker-gated install matrix + detect/install/uninstall; **5** graduated capability allowlist + durable tool-decision audit; **6** multi-seed impact BFS + contains-asymmetric radius |
| Medium | **7** surgical context packet + skeletonization/session dedup; **8** two-pass call resolve tagged INFERRED; **9–10** episode/evidence pointers + contradict→PlanChange/STALE; **11** surface-hash skip for dependents; **14** RRF FTS+graph (no embeddings); **15** single-active-task claim; **18** deterministic entity identity cascade; **19** Discovery→PlanChange templates |
| Dogfood-first / later | **12** Scout/Verify/Auditor policy; **13** Graph-vs-Explorer A/B methodology; **17** report mix % / verify-INFERRED prompts; **20** community clustering (heuristics only) |

Rejected ranks (21–22 and §D) stay rejected — see §5.

### (B) Roadmap items never started (or only stubbed)

Mapped to `docs/ROADMAP.md` conceptual IDs (not board phase numbers):

| Roadmap theme | Status | Note |
|---------------|--------|------|
| **P13 Hypothetical plan branches** (`plan simulate`, adopt/discard) | **not started** | Explicitly deferred through Phase 05+; **≠** board Phase 13. |
| **P14 Forward-state / reversal model** (beyond supersede stubs) | **partial / thin** | Supersede-not-delete + STALE exist; full reversal taxonomy / rollback UX incomplete. |
| **P9 Phase-level review** as first-class layer | **partial** | Task + scope review live; phase review strategy thinner than design doc. |
| **P16 scale** 100k–1M LOC ladders | **deferred** | Gate H stopped at ~10k planted. |
| **P17 Multi-agent concurrency** (beyond worktree bind) | **partial** | Phase 08 worktrees; no swarm; dogfood D33 planned. |
| **P18 Ecosystem** (VCS/LLM/viz plugins beyond analyzers) | **partial** | Analyzer contribution path only. |
| **P19 Evaluation / research release** (public benchmark suite, full baseline ladder) | **partial** | Harnesses exist; published research package not cut. |
| **P20 Production hardening** (observability, stable plugin APIs, compatibility guarantees) | **partial** | Phase 08 checklist; not “production release.” |
| **VerifiedFact** promotion engine | **deferred** | Cost/risk; residuals instead. |
| **Embeddings / semantic SoT** | **deferred** | Law 13 until measured need. |
| **HTTP daemon / canonical API hub** | **deferred / rejected as P0** | CLI+MCP adapters only; roadmap P6 daemon shape not adopted. |
| **GC-03 / GC-04** (Gate C Mode-B residuals) | **deferred** | Unless measurement re-promotes. |

### (C) Dogfood ladder gaps

From [`experiments/LADDER.md`](../../experiments/LADDER.md) / [`RESULTS.md`](../../experiments/RESULTS.md) — still **planned** / **deferred** (not owned by Phase 13):

| ID | Gap |
|----|-----|
| **D12** | Shared domain invariant stressed by two features |
| **D21** | Combo replan + honesty |
| **D22** | Combo context + impact before destructive edit |
| **D23** | Combo full slice (why+context+plan+decision+review) |
| **D31** | Separate reviewer agent |
| **D33** | Parallel sibling agents / colliding decisions |
| **D40–D41** | Real mini / multi-day external validity |
| **D42** | Gate C live packs on richer corpora (harness-adjacent) |

Tier 0–1 isolations and many Tier 2/3 rungs already scored (**G1>pass / B0 fail** portfolio). Method notes DF-06/07/13/34/36 remain experiment-only.

### (D) Polish / release

- Named **Gate D** ablation (progressive vs static).
- Stronger **Gate F/G** beyond planted prelim (cost-aware multi-model optional later — not required to keep honesty).
- Gate H **100k/1M** optional when hardware/CI budget exists.
- MCP surface expansion only where dogfood proves CLI friction (plan/impact/index) — keep thin.
- Public packaging: versioned CLI, contributor docs, benchmark publication (P19), Apache-2.0 core narrative (D13).
- Manual Cursor MCP reload remains operator friction (mitigated by tips; not solvable solely in-repo).

---

## 4. Recommended sequence after P13

Prefer **thin phases**; promote from research ranks + ladder gaps; do **not** assume a single mega “Phase 14.”

1. **Peer impact + install gates (thin)** — Research S03–S04: multi-seed/contains impact walks + marker-gated install/allowlist audit. Closes H4/H7 depth without daemon/embeddings. **Boarded 2026-08-17** as [`phase-14-peer-impact-install-gates`](../phases/phase-14-peer-impact-install-gates/).
2. **Causal supersession / episodes (thin)** — Research S05: contradict→STALE/PlanChange + evidence pointers. Strengthens H3/H5 honesty under change.
3. **`plan simulate` / plan branches (roadmap P13)** — Only after impact walks are trustworthy; keep adopt/discard fail-closed; human authority (D10).
4. **Dogfood Tier-2/3 batch** — D21–D23 then D31/D33 off-board in parallel; feed DF-* forward only. Optional Gate D harness as eval-only scope.
5. **Scale + release cut** — Gate H next ladder rung *or* P19 benchmark publish + packaging; pick by whether “large repo” or “research credibility” is the next risk.

If capacity is tiny: do **(1)** then dogfood **D21/D22** before simulate.

---

## 5. Non-goals still rejected

Still rejected (ROADMAP §1, `G_PROJECT_LAWS`, peer review §D, AGENTS hard boundaries):

- Cloud SaaS / hosted control plane as core identity  
- Large dashboard / Notion-like project manager / Obsidian clone as product  
- Swarm / unrestricted autonomous multi-agent frameworks  
- Proprietary graph DB (Neo4j/FalkorDB/etc.) as SoT — keep `.trace/` SQLite  
- Embeddings / vector store / in-binary embed models as **default** retrieval SoT  
- Full-rebuild-on-any-change indexer  
- Daemon / always-on HTTP / SSE as P0 architecture (MCP remains thin adapter)  
- Universal language megastore / Hybrid LSP matrix / Rust napi kernel rewrite without measured Gate need  
- YOLO / AllowAll capability defaults  
- LLM similarity edges as verified facts  
- Auto-summarize-everything / auto-generate entire backlog via LLM  
- Rewriting closed board history (spawn forward only)

---

## Appendix — Board phase delivery map (00–12 closed; 13 in flight)

| Board phase | Delivered (one line) |
|-------------|----------------------|
| **00** | Scaffold, SQLite, Git index, tree-sitter JS/TS/Py, causal domain, FTS+Why+context, CLI, P0-X 7/7 |
| **01** | Review→DONE, honesty Paths A/B/C, X0 dry-run harness, MCP stdio |
| **02** | Gate C Go + GC-01/02 harden |
| **03** | Progressive planner + discovery replan; Gate E |
| **04** | Scope review + residuals; Gate G prelim |
| **05** | Impact classes/report; Gate F prelim |
| **06** | Capability surface + selection ablation |
| **07** | Ignore tiers, Go analyzer, Gate H ≤~10k |
| **08** | Plugin surface, worktrees, hardening checklist; historical `no successor` |
| **09** | Dogfood UX: why-after-review, `tasks`, install |
| **10** | Integrity: DPC scoping, index GC, operator/caps, MCP parity |
| **11** | Residual DF wave (18 DFs); historical `no successor` |
| **12** | Peer honesty: edge provenance + packet staleness/truncation; historical `no successor` |
| **13** | *(in progress)* Import resolve + honesty residuals DF-60…67 — **out of scope for this remaining backlog** |

---

## Sources

README, `docs/ROADMAP.md`, `docs/EVALUATION.md`, `docs/DESIGN_DECISIONS.md`, `docs/init/A_PROJECT_PLAN.md`, `docs/init/G_PROJECT_LAWS.md`, `AGENTS.md`, `docs/TODO.md` (status skim), phase READMEs 00–13, `docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md`, `experiments/{LADDER,RESULTS,DOGFOOD-FINDINGS}.md`.
