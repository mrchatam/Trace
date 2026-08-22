# P14-S01-00 — Impact walks (FINAL)

## Metadata
- id: P14-S01-00
- todo_ids: [P14-S01-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Inventory live impact domain (`internal/` impact + `trace impact` + Gate F prelim). Lock **FINAL** defaults for multi-seed BFS + depth-bounded contains asymmetry (research rank **6**). Thicken `01`/`02`/SCOPE-TODOS. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A8; rank 6
- [SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — CBM multi-seed BFS; codegraph contains asymmetry
- [DECISION_IMPACT.md](../../../../DECISION_IMPACT.md)
- [TRACE-GOALS-PROGRESS-2026-08-17.md](../../../../research/TRACE-GOALS-PROGRESS-2026-08-17.md) — H4 partial
- Live: Phase 05 impact packages, `evals/impact`, `cmd/trace` impact, `internal/retrieval` Expand
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material phase locks A1–A8 hold; inventory before FINAL.

## Live inventory (2026-08-17)

| Area | Present? | Evidence / gap |
|------|----------|----------------|
| Decision impact classes + alternatives | **Present** | `internal/domain/impact.go` + `internal/store/impact.go` + mig **009**; bands SAFE…REVERSAL; fail-closed `ImpactReport` |
| CLI `trace impact` | **Present** | `cmd/trace/impact.go` — `finding` / `alternative` / `report` only — **no walk/blast subcommand** |
| Gate F prelim | **Present** | `evals/impact` `TestPlantedImpactConflictsGateFPrelim` (planted P/R); keep green |
| Multi-seed structural BFS | **Missing** | No ImpactWalk / BlastRadius API; Expand is retrieval neighborhood, not impact blast |
| Seed exclusion from blast | **Missing** | Expand keeps seeds in output at distance 0 |
| Contains asymmetry | **Missing / opposite of Expand** | Expand `file`→all symbols **and** `symbol`→parent file (bidirectional). Impact must **not** climb contains-up into siblings |
| Incoming import deps | **Missing as walk** | Expand walks **outgoing** imports (`ListImportsByPath` → resolve). No `ListImportersOf*` helper; reverse scan needed for “deps in” |
| Import path resolve | **Present (P13)** | `resolveImportedFile` in retrieval — **reuse** for matching importers; do not reimplement |
| Edge provenance on imports | **Present (P12)** | Copy onto blast nodes when edge is an import hop |
| Loud truncation / totals | **Missing on impact** | Compiler Budget has `items_kept`/`items_total`/`truncated` — impact walk needs analogous fields |
| Hop risk | **Missing** | No per-node hop_risk on structural neighbors |
| Causal / decision seeds in walk | **Out of bar** | Decision impact stays `ImpactReport`; structural walk = `file`\|`symbol` only |
| MCP impact tool | **Absent** (good) | No `trace_impact` MCP — **do not add** (A3) |
| Migration for walks | **Not needed** | Compute-time over `files`/`symbols`/`imports` |
| `internal/impact` package | **Absent** (keep) | Extend retrieval + thin CLI; do not invent second impact stack |
| DR-NOIMP | **In force** | Walks are structural honesty — **not** commercial auto-classifier / plan mutation / auto-plant findings |

## Locked defaults (FINAL) — do not re-debate in P14-S01-01

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Research | Rank **6** — multi-seed BFS (CBM) + contains asymmetry (codegraph) |
| Home (algorithm) | **`internal/retrieval`** — new ImpactWalk API (name free: `ImpactWalk` / `WalkImpact` / `BlastRadius`). Peer of Expand; **do not** change Expand bidirectional semantics for context/Why |
| Home (CLI) | **`trace impact walk`** thin adapter (G19) — open store + `retrieval.New`; mirror `why`/`context` pattern. Keep existing `finding`/`alternative`/`report` green |
| Domain planted impact | **Untouched** — `ImpactReport` / findings / alternatives / mig 009 stay Gate F substrate |
| Seeds | ≥1 seed; each `{entity_type, entity_id}` with `entity_type` ∈ **`file`\|`symbol` only**. Reject empty / unknown types. Multi-seed = **one** BFS (shared `seen` + frontier); all seeds enter at distance **0** |
| Seed exclusion | Result set **excludes** every seed `(type,id)`. Seeds never appear as blast hits |
| Depth | Integer **1..2** (same ceiling as Expand). Reject &lt;1 or &gt;2. Default when CLI omits: **2** |
| Dep edges (“walk deps in”) | **Incoming import only** — files whose resolved `imported_path` targets the current file path (reuse P13 `resolveImportedFile` from importer side). **Do not** walk outgoing imports as blast edges |
| Reverse importer lookup | Add store and/or retrieval helper (name free) — e.g. scan imports + resolve, or indexed query. Prefer small store helper + retrieval walk; **no** schema migration |
| Contains edges | Trace **file → symbols** = contains-OUT; **symbol → file** = contains-UP |
| Contains asymmetry | (1) From a **file** node: may expand **contains-OUT** to its symbols. (2) From a **symbol** seed/node: may place the **containing file** into the blast (one contains-UP into the result) so incoming deps of that file can be walked — but **must not** enqueue contains-OUT from that file to **sibling** symbols. (3) Never treat sibling symbols as impacted solely via contains climb |
| Hop risk | Each blast node exposes `hop` (= BFS distance) and `hop_risk` (monotonic non-decreasing with hop; lock: `hop_risk = float64(hop)` is enough). Seeds are hop 0 and excluded from results |
| Loud truncation | Cap blast result list at **`MaxImpactBlast = 64`** (align MaxCandidateHits). JSON (and any MD if added): **`blast_total`** (pre-cap unique non-seed nodes discovered), **`blast_kept`** (`len(returned)`), **`truncated`** true when `blast_kept < blast_total`. Deterministic order before cap: hop ASC, then entity_type ASC, entity_id ASC. **No silent caps** |
| Edge provenance | When a blast hop is via import, copy `imports.provenance` onto the hit (`edge_provenance`). Contains hops omit / empty provenance |
| Result shape | Typed result (e.g. `ImpactWalkResult`) with `Seeds`, `Blast []Hit` (or dedicated blast node type reusing Hit fields), `BlastTotal`, `BlastKept`, `Truncated`, `Depth`. Fail-closed: missing seed entity → error (not empty success) |
| Packages touched | Prefer **`internal/retrieval`** (+ tests) + thin **`cmd/trace/impact.go`** (+ CLI test). Store helper OK in `internal/store/file_graph.go` if reverse lookup needs SQL. **No** mig; **no** `internal/impact`; **no** analyzer rewrite; **no** compiler packet changes required |
| Migration | **None** |
| MCP | **No** new impact tool / menu (A3, A6 loudness ≠ new MCP) |
| Gate F | Keep `TestPlantedImpactConflictsGateFPrelim` green; do not repurpose planted harness into graph walk |
| Forbidden | Silent truncation; changing Expand to “fix” context by removing symbol→file; outgoing-import-as-blast; auto-planting findings / plan mutation; commercial engine UX (DR-NOIMP); daemon/HTTP/embeddings/Neo4j; full-rebuild indexer; new MCP impact; boarding S05 / `plan simulate`; reopening DF-60…67; board spawn by implementer |
| Carry-forward | honesty A/B/C+G; Gates E/F prelim/H; ablation; compat; p0x; x0; Gate C `dry_run:false`; P12/P13 provenance + packet honesty named tests; product `./cmd|internal|evals` |
| Named tests (min) | See table below |

### Named tests (intent locked)

| Test | Intent |
|------|--------|
| `TestImpactWalkMultiSeedExcludeSeeds` | Two file seeds → one walk; neither seed in blast; shared neighbors appear once at min hop |
| `TestImpactWalkContainsAsymmetryNoSiblings` | Symbol seed in file with ≥2 symbols → blast may include containing file + importers; **must not** include sibling symbol via contains |
| `TestImpactWalkIncomingImportHop` | B imports A (analyzer-shaped or store fixture + P13 resolve) → walk seed A includes B at hop≥1 with `edge_provenance` when import row has provenance |
| `TestImpactWalkLoudTruncation` | Force &gt;64 blast nodes → `truncated=true`, `blast_kept=64`, `blast_total` &gt; kept |
| `TestImpactWalkHopRiskIncreases` | Node at hop 2 has `hop_risk` ≥ hop-1 neighbor’s `hop_risk` |
| `TestPlantedImpactConflictsGateFPrelim` | **Keep** green |

## Owns
| Item | Intent |
|------|--------|
| Multi-seed impact BFS | One walk over N seeds; exclude seeds from blast; hop risk |
| Contains asymmetry | Expand containers out via contains; never climb contains-up into siblings |
| Loud honesty | Truncation/totals align with Phase 12 packet loudness |

## Explicit deferrals (not S01)
- Install matrix / allowlist audit (S02)
- Supersession / episodes (research S05)
- `plan simulate`; commercial full impact engine (DR-NOIMP)
- New MCP impact tool menu
- Causal/decision/task seeds in structural walk
- Outgoing-dep / “what does this depend on” mode
- Call-graph / symbol-level edges beyond file↔symbol contains
- Auto-linking walk results into `decision_impact_findings`

## Depends note (S02)
S02 install/capability work must **not** regress `trace impact` (`finding`/`alternative`/`report` **and** new `walk`). Carry-forward verify for S02/S03 includes S01 named ImpactWalk tests + Gate F prelim. See light note on S02 SCOPE-TODOS.

## Planner work (this row)
1. [x] Inventory live impact walk APIs, seed model, edge types, truncation, Gate F tests
2. [x] Lock FINAL: multi-seed entry, seed exclusion, hop/depth rules, contains asymmetry, loud truncation, packages/migrations, named tests
3. [x] Thicken `01-impact-walks.md` + `02-scope-review.md` + SCOPE-TODOS to FINAL
4. [x] Light Depends note for S02 (install must not regress impact CLI)

## Exit
- [x] 00-PLANNER marked **FINAL**
- [x] Board Notes; next **P14-S01-01**
- [x] Product Go — **not** this row
