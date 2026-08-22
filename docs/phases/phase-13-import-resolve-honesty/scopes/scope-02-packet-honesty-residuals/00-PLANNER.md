# P13-S02-00 — Packet honesty residuals (FINAL)

## Metadata
- id: P13-S02-00
- todo_ids: [P13-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S02 for **DF-61, DF-62, DF-63, DF-65** (packet/index honesty residuals after P12). **No product Go in this row.**

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Laws 6–7, 18
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A4; S02 home
- DOGFOOD-FINDINGS DF-61…63, DF-65; [`_bughunt/post-p12/POST-P12-BUGHUNT.md`](../../../../../experiments/_bughunt/post-p12/POST-P12-BUGHUNT.md)
- Repros: `_bughunt/post-p12/{stalecap,staledrop,candcap,prov}/`
- Live: `internal/compiler/{packet.go,index_honesty.go,compiler.go,budget.go}` + `internal/retrieval/expand.go` (S01 resolve)
- Depends: **P13-S01-02 APPROVE** (DF-60 resolve shipped — reuse for DF-65; do not re-join paths)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Prefer false-fresh on I/O miss; Law 18 causal STALE untouched.

## Live inventory (2026-08-17, post–S01 APPROVE)

| Area | Present? | Gap vs DF |
|------|----------|-----------|
| SchemaVersion `0.2` + Budget loud fields | **Present** | Keep; additive JSON only (no inventive schema bump) |
| `candidates_capped` flag + MD | **Present** | Flag holds; **DF-63** `items_total` = post-cap pipeline len (≤64), understates admit universe |
| `buildIndexHonesty(kept)` + sort-then-cap 8 | **Present** | Cap silent — **DF-61** no `stale_total` / truncated signal; universe = kept only — **DF-62** trim/FTS-omit → `index_honesty: null` while disk-stale files existed pre-trim |
| False-fresh on missing row / disk / I/O | **Present** | Keep (A4) |
| Law 18 | **Present** | Honesty never sets `Provenance.Status=STALE` |
| P12 named tests | **Present** | Keep green: `TestBudgetLoudTotals` / `TestCandidateCapSetsTruncated` / `TestIndexStaleBanner` / `TestContextWhyTraceEdgeProvenance` |
| S01 `resolveImportedFile` + Expand file→import | **Present** | Resolve works when Expand seeds are **files** |
| Task `compileAtDepth` | **Gap DF-65** | Expand(task) then **append FTS**; FTS file hits never Expand'd → import hops / `edge_provenance` absent on `context` even after DF-60 |

### Fixture map (bughunt → lock)

| DF | Fixture | Observed defect |
|----|---------|-----------------|
| DF-61 | `stalecap/` (13 stale → 8 listed) | Silent path cap |
| DF-62 | `staledrop/` | Honesty only on kept → false-fresh under trim |
| DF-63 | `candcap/` (80 links → MD `32/64`) | Post-cap `items_total` |
| DF-65 | `prov/` (+ S01 resolve) | Context FTS file without import-hop provenance |

## Owns (phase lock)

| DF | Intent |
|----|--------|
| DF-61 | `stale_paths` cap ≤8 must expose **total** / truncated signal (not silent drop) |
| DF-62 | Disk-stale files omitted by trim must not yield false-fresh `index_honesty: null` when they were in the pre-trim honesty universe |
| DF-63 | When `candidates_capped`, `items_total` must not understate admit universe as post-cap length alone |
| DF-65 | Task `context` carries import-hop `edge_provenance` when hops are relevant (not FTS-only silence) — **reuse S01 resolve via Expand on file seeds** |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | **`internal/compiler` primary**; DF-65 may call existing `Retriever.Expand` on file hits (prefer **zero** new retrieval resolve logic; thin retrieval edit only if Expand merge helper is cleaner — **no** path-join reimplementation) |
| Migration | **None** |
| SchemaVersion | Keep packet **`0.2`** (additive fields OK; do not invent `0.3`) |
| Law 18 | Causal `Provenance.Status=STALE` **untouched** by index honesty |
| False-fresh | Missing row / missing disk / I/O/hash error → **omit** path (never false-stale) |
| Symbol / whole-index scan | **Out of bar** (DF-67 → S03); honesty universe = **file** items only |
| MCP / CLI | Library marshal only — **prefer zero** adapter edits; no new MCP tools (G19) |
| Forbidden | Silent omit of totals; fake provenance; path-align product hook; analyzer rewrite; daemon watcher; full-rebuild; board spawn by implementer; weakening P12 named locks |
| Carry-forward | P12 honesty named tests; honesty A/B/C+G; E/F/ablation/H/compat; p0x; x0; Gate C `dry_run:false`; S01 Expand/Why import tests |

### DF-61 — stale list honesty

| Lock | Value |
|------|-------|
| Cap | Keep `maxIndexHonestyStalePaths = 8`; **sort-then-cap** lex (already) |
| JSON | `IndexHonesty` gains **`stale_total`** (`int`, full unique stale count **before** path cap) and **`stale_truncated`** (`bool`, true when `stale_total > len(stale_paths)`) |
| Markdown | When truncated, banner must include total signal (e.g. `stale_total=N` or `showing K of N`) — not paths alone |
| Omit | Still omit entire `index_honesty` when stale_total==0 |

### DF-62 — honesty under trim

| Lock | Value |
|------|-------|
| Honesty universe | File items in the **pre-trim** pipeline (`items` after Layer-1 admit / MaxCandidateHits, **before** `trimToBudget`) — **not** `kept` alone |
| Behavior | Disk-stale file dropped by MaxItems/token trim still contributes to `stale_paths` / `stale_total` |
| Non-goal | Do **not** scan entire index / un-admitted FTS misses outside the pre-trim pipeline (bounded residual OK to Note) |

### DF-63 — admit-universe totals

| Lock | Value |
|------|-------|
| `items_kept` | Unchanged — `len(kept)` after trim |
| `items_total` | Count = Layer-0 items already on the packet **plus** number of **unique Layer-1–admissible** hits in the full Expand+FTS candidate list (via `layer1AdmitKey` / same eligibility), **not** truncated by `MaxCandidateHits` |
| Cap behavior | Pipeline still admits ≤ `MaxCandidateHits` into trim input; `candidates_capped=true` when more admissible hits remain |
| When capped | Expect `items_total` **>** len(pre-trim pipeline items) (e.g. 80 decisions → total ≈ layer0+80, not 64) |
| Markdown | Existing `items=kept/total` + `candidates_capped=true` — total must reflect universe |
| New Budget fields | **Prefer none** — fix `items_total` meaning under cap (P12 tests do not assert total==64) |

### DF-65 — context import hops (post–S01)

| Lock | Value |
|------|-------|
| Root cause | `compileAtDepth` Expand(task) then appends FTS; FTS **file** hits never become Expand seeds → no file→import neighbors / `edge_provenance` |
| Fix home | `internal/compiler` `compileAtDepth`: after task Expand + FTS (or equivalent), **Expand file-typed hits** (depth 1) via existing `Retriever.Expand`, merge into candidates, preserve `Hit.EdgeProvenance` → Item |
| Resolve | **Reuse S01** `resolveImportedFile` inside Expand — do not re-join `./` / extensions in compiler |
| Bound | Prefer Expand only file hits already present in Expand∪FTS candidate pool (not whole-repo walk); MaxCandidateHits still applies |
| Depth | Must work for `TaskContext` (depth 1) and `ExpandContext` depth 2 — dogfood used `--depth 2` |
| Non-goal | Do not require every FTS file to keep hops after trim; assert at least one admitted import-neighbor Item (or WhyTrace if IncludeWhy) carries `edge_provenance` when fixture hop is relevant |

## Named tests (intent locked)

| Test | Intent |
|------|--------|
| `TestBudgetLoudTotals` | **Keep** |
| `TestCandidateCapSetsTruncated` | **Keep** (flag + truncated); extend or add sibling for DF-63 |
| `TestIndexStaleBanner` | **Keep** (kept-file stale + false-fresh delete + Law 18) |
| `TestContextWhyTraceEdgeProvenance` | **Keep** |
| New DF-61 | >8 disk-stale files in honesty universe → `len(stale_paths)==8`, `stale_total` > 8, `stale_truncated==true`, MD shows total |
| New DF-62 | Force MaxItems trim to drop a disk-stale file item → `index_honesty` still non-null listing that path (or total>0); not null false-fresh |
| New DF-63 | >64 Layer-1 admits → `candidates_capped` + `items_total` ≥ true admissible universe (≫64), MD `items=kept/total` uses that total |
| New DF-65 | Task title/FTS admits importer file; indexed relative import (analyzer-shaped `./util`) → TaskContext/ExpandContext packet Item for neighbor includes `edge_provenance` (EXTRACTED fixture) |

Names free except P12 keepers. Combined suite file OK.

## Explicit deferrals (not S02)
- DF-64 / DF-66 / DF-67 → **S03**
- Whole-index / symbol-entity staleness
- New clarifying experiments; optional `ab-import-resolve` dogfood (not board blocker)

## Planner work (this row)
1. [x] Inventory live compiler/packet/index_honesty vs bughunt
2. [x] Lock FINAL APIs/tests (DF-61/62/63/65)
3. [x] Thicken 01 / 02 / SCOPE-TODOS
4. [x] Light S03 Depends note (schema stays additive `0.2`)
5. [x] Board P13-S02-00 → done; next **P13-S02-01**
6. [ ] Product Go — **not** this row

## Exit
- [x] FINAL locks; thicken 01/02/SCOPE-TODOS; next **P13-S02-01**
