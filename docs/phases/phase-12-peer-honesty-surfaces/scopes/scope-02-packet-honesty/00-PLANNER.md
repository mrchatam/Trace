# P12-S02-00 — Packet honesty (FINAL) — RETRY 2026-08-17

## Metadata
- id: P12-S02-00
- todo_ids: [P12-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Re-inventory live `internal/compiler` and re-lock S02 implement/review prompts for **compiler packet honesty**: emission-time index-staleness banners + loud truncation/totals. Prior 2026-08-16 planner `done` is **invalid** (Notes claimed empty inventory / no product Go while types already shipped). **No product Go in this row** (read-only vs tree; thicken docs + board only).

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Laws 6–7 (bounded progressive packets), Law 18 (causal STALE)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A7; research ranks 2–3
- [SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — codegraph staleness; CBM loud truncation
- Live: `internal/compiler` (`packet.go`, `budget.go`, `compiler.go`, `index_honesty.go`), `internal/store` (`files.content_hash`, `ProjectRoot`), `cmd/trace/context.go`, `internal/mcp` (library marshal only)
- Depends: **S01 APPROVED** (P12-S01-02 / REVIEW-NOTES) — preserve JSON **`edge_provenance`** pass-through; do **not** overload causal `Item.Provenance` / `confidence`; do **not** regress S01 tests
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Prefer false-fresh over false-stale (A3); Law 18 causal STALE stays authoritative (emission banners ≠ mutating provenance `status`).

## Verdict on prior done (2026-08-16)

**Invalid.** Board Notes claimed pre-change inventory (“Budget has truncated but no totals / no emission staleness / no product Go”) while live tree already had SchemaVersion `"0.2"`, Budget totals/cap fields, `IndexHonesty` + `index_honesty.go`, and `compileAtDepth` wiring. Named locks `TestBudgetLoudTotals` / `TestCandidateCapSetsTruncated` / `TestIndexStaleBanner` were (and remain) **missing**. P12-S02-01 correctly stayed `pending` for remaining implement work.

## Live inventory (2026-08-17)

| Area | Present? | Evidence / gap |
|------|----------|----------------|
| `SchemaVersion` | **Present** | `packet.go` `SchemaVersion = "0.2"`; `compiler.go` sets `pkt.SchemaVersion`; capability test asserts `"0.2"` |
| Budget JSON `items_total` / `items_kept` / `candidates_capped` | **Present** | `Budget` fields in `packet.go`; set in `compiler.go` (`itemsTotal`, `itemsKept`, `candidatesCapped`) |
| `truncated` when kept&lt;total **or** candidates capped | **Present** | `compiler.go`: `if itemsKept < itemsTotal \|\| candidatesCapped { truncated = true }` |
| MaxCandidateHits (64) → `candidates_capped` | **Present** | `compiler.go` Layer-1 loop breaks at `MaxCandidateHits`; peeks remaining admits via `layer1AdmitKey` |
| MD budget `items=kept/total` + `candidates_capped=true` | **Present** | `RenderMarkdown` in `packet.go` |
| Emission `index_honesty` | **Present** | type `IndexHonesty`; `buildIndexHonesty` in `index_honesty.go`; wired `IndexHonesty: buildIndexHonesty(...)` in `compiler.go` |
| False-fresh on missing row / missing disk / I/O | **Present** | `index_honesty.go` `continue` on those errors; omit object when no stale |
| Stale path cap 8 + notice | **Present (minor residual)** | Cap via `maxIndexHonestyStalePaths`; notice const set. **Residual:** paths capped in encounter order then sorted among kept ≤8 — prefer **sort-then-cap** for deterministic lex first-8 (see lock delta) |
| Law 18 causal STALE untouched | **Present (by design)** | Honesty path never sets `Provenance.Status` |
| S01 `edge_provenance` | **Present** | Item / WhyTrace fields + `TestContextWhyTraceEdgeProvenance` |
| Migration size/mtime | **Not needed** | `files.content_hash` sufficient |
| Named `TestBudgetLoudTotals` | **Missing** | only S01 `TestContextWhyTraceEdgeProvenance` among honesty locks; no loud-totals assertion |
| Named `TestCandidateCapSetsTruncated` | **Missing** | — |
| Named `TestIndexStaleBanner` | **Missing** | no mutate-disk / false-fresh deleted-path coverage |

## Locked defaults (FINAL) — phase intent unchanged

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | **`internal/compiler`** primary; thin CLI/MCP only if marshal already covers (prefer **zero** adapter edits) |
| Migration | **None** — do not add size/mtime columns this scope |
| SchemaVersion | Packet `SchemaVersion` **`0.2`** (already shipped — do not re-bump inventively) |
| Loud truncation — Budget JSON | Keep existing fields; **`items_total`** (count fed to `trimToBudget`), **`items_kept`** (`len(kept)`), **`candidates_capped`** (`bool`, true when `MaxCandidateHits` stopped admission early). **`truncated=true`** whenever `items_kept < items_total` **OR** `candidates_capped` |
| Loud truncation — Markdown | Budget line must include `items=kept/total` and `truncated=…`; when capped, also `candidates_capped=true`. **No silent caps** |
| Staleness — trigger | Emission-time only: for each **kept** packet item with `entity_type=file`, resolve path via store (`GetFileByID`), read disk under `store.ProjectRoot()`, sha256 hex vs `files.content_hash`. Prefer **false-fresh**: missing indexed row, missing disk file, I/O/hash error → **do not** list as stale |
| Staleness — Packet field | `index_honesty` object (omitempty): `stale_paths` (unique rel paths, **sort then cap 8**), optional `notice` when non-empty. Omit entire object when no stale paths |
| Staleness — Markdown | When `index_honesty.stale_paths` non-empty, banner near packet header (before Items) naming paths; never emoji-required |
| Staleness ≠ Law 18 | Do **not** set causal `Provenance.Status=STALE` from disk drift; banners only |
| Symbol items | **Out of bar** this scope — only `file` entities |
| Skeletonization / session dedup | **Optional stretch** — not acceptance (A5) |
| Hash helper | `crypto/sha256` inside `internal/compiler` (do **not** import `internal/analyzers` from compiler) |
| S01 intact | Keep `edge_provenance` pass-through on Item / WhyTrace; never write enum into causal `Item.Provenance` / `confidence` |
| MCP | Library-only parity — **no** new MCP tool / menu (A6) |
| Forbidden | Silent truncation; daemon FS watcher as product; embeddings; Neo4j; full-rebuild indexer; research S03+ boarding; conflating index banner with causal STALE; board spawn by implementer; **blindly re-implementing already-present types** |
| Carry-forward | honesty A/B/C+G; E/F/ablation/H/compat; p0x; x0; Gate C `dry_run:false`; S01 provenance named tests |
| Named tests (min) | `TestBudgetLoudTotals` (truncated ⇒ `items_total`>`items_kept`, MD shows kept/total); `TestCandidateCapSetsTruncated` (force >64 Layer-1 admits ⇒ `candidates_capped` + `truncated`); `TestIndexStaleBanner` (index file, mutate disk bytes, context with file item ⇒ `index_honesty.stale_paths` contains path; false-fresh when path deleted); keep **`TestContextWhyTraceEdgeProvenance`** green |
| Verify | See `01-packet-honesty.md` |

### FINAL lock deltas vs 2026-08-16 draft

| Delta | Why |
|-------|-----|
| Inventory rewritten to **2026-08-17 present vs missing** | Prior inventory contradicted live tree |
| Implementer posture: **finish gaps / prove with named tests**; do not re-add SchemaVersion/`Budget`/`IndexHonesty` from scratch | Product types already shipped |
| `stale_paths`: clarify **sort-then-cap 8** (deterministic lex) | Live code caps in encounter order then sorts the truncated set — safer lock for stable banners/tests |
| Review checklist: planner may find pre-existing product Go; this row still must not author Go | Historical “no product Go” claim was the contradiction that forced retry |
| Named test names — **unchanged** | Still the acceptance bar for P12-S02-01 |

## Owns
| Item | Intent |
|------|--------|
| Staleness banners | Emission-time notice when indexed file view may lag disk |
| Loud truncation | Exact totals + truncated/capped signals — no silent caps |

## Explicit deferrals (not S02)
- size/mtime watermark columns + index pipeline fast-path (research rank 11 surface-hash)
- Daemon / FS watcher pending-sync queues (codegraph watcher path)
- Skeletonization / session dedup as required product
- Symbol-parent / whole-index scan staleness
- Research ranks 4+ / FUTURE S03–S05

## Remaining implement work (for P12-S02-01 — not this row)

1. Add named tests `TestBudgetLoudTotals`, `TestCandidateCapSetsTruncated`, `TestIndexStaleBanner` (plus keep S01 regression green).
2. If tests expose gaps: fix wiring only as needed (esp. sort-then-cap for `stale_paths`; any MD/JSON mismatch).
3. Do **not** rewrite working types unless a lock fails under test.
4. Run locked verify suite; board status + Notes only.

## Planner work (this row)
1. [x] Re-inventory live compiler (2026-08-17) — present vs missing
2. [x] Confirm / amend FINAL locks + document deltas
3. [x] Thicken `01-packet-honesty.md` + `02-scope-review.md` + SCOPE-TODOS for gap-focused implement
4. [x] Light Depends note for S03 VERIFY evidence table

## Exit
- [x] Thicken 01 + 02 + SCOPE-TODOS to FINAL (retry)
- [x] Board Notes cite inventory evidence; next **P12-S02-01**
- [x] Product Go — **not** this row
