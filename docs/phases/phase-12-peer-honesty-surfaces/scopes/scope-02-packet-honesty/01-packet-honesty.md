# P12 / S02 / 01 — Packet honesty (FINAL) — gap finish

## Metadata
- id: P12-S02-01
- todo_ids: [P12-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
**Finish** packet honesty to FINAL locks: primarily the **three named tests** (+ any small fixes they expose). Live tree **already has** SchemaVersion `0.2`, Budget totals/cap fields, MD loud line, and `index_honesty` emission helper. **Do not blindly re-implement** those types. Preserve S01 **`edge_provenance`**. **Stop if 00-PLANNER is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (2026-08-17 retry inventory required)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A7; A3 false-fresh
- Research ranks 2–3 — codegraph staleness; CBM truncation
- Laws 6–7 / 18 — bounded packets; causal STALE ≠ index banner
- S01 REVIEW-NOTES — do not regress `edge_provenance`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do not re-debate FINAL locks. Re-read **00-PLANNER live inventory (2026-08-17)** before coding.

## Already shipped (do not re-create from scratch)

| Surface | Where |
|---------|--------|
| `SchemaVersion = "0.2"` | `internal/compiler/packet.go` |
| `Budget.items_total` / `items_kept` / `candidates_capped` + `truncated` rules | `packet.go` + `compiler.go` `compileAtDepth` |
| MaxCandidateHits → `candidates_capped` | `compiler.go` Layer-1 loop + `layer1AdmitKey` peek |
| MD `items=kept/total` + optional `candidates_capped=true` | `RenderMarkdown` in `packet.go` |
| `IndexHonesty` + `buildIndexHonesty` (false-fresh) | `packet.go` + `index_honesty.go`; wired on Packet in `compiler.go` |
| S01 `edge_provenance` + regression | Item/WhyTrace + `TestContextWhyTraceEdgeProvenance` |

## Still missing / finish here

| Gap | Action |
|-----|--------|
| `TestBudgetLoudTotals` | Force trim (max_items and/or token limit) → assert `items_total` > `items_kept`, `truncated`, MD contains `items=` kept/total |
| `TestCandidateCapSetsTruncated` | Force >64 Layer-1 admits → `candidates_capped` + `truncated` (even if final kept ≤ MaxItems) |
| `TestIndexStaleBanner` | Index file item in packet, mutate disk bytes → `index_honesty.stale_paths` contains path; delete path → false-fresh (omit / no path); MD banner when stale set |
| sort-then-cap residual | If needed for stable tests: unique stale paths → **sort** → take first **8** (see 00 lock delta) |
| Verify suite | Carry-forward cmds below |

## Locked defaults (FINAL — do not renegotiate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/compiler` (+ CLI/MCP only if library JSON already insufficient — prefer zero adapter edits) |
| Migration | **None** |
| SchemaVersion | Packet **`0.2`** (already set) |
| Budget JSON | `items_total`, `items_kept`, `candidates_capped`; `truncated` if kept&lt;total **or** candidates capped |
| Markdown | Budget line: `items=kept/total` + truncated (+ candidates_capped when set) |
| Index honesty | `index_honesty` omitempty: `stale_paths` (unique, **sort then ≤8**) + optional `notice`; emission-time sha256(disk) vs `files.content_hash` for kept `file` items; **false-fresh** |
| Law 18 | Do **not** mutate causal `Provenance.Status` from disk drift |
| S01 intact | Keep `edge_provenance` pass-through; never write enum into causal `Item.Provenance` / `confidence` |
| Forbidden | Silent caps; daemon watcher product; embeddings; Neo4j; full-rebuild; board spawn; size/mtime migration; **rewriting working honesty types without a failing lock** |
| Named tests | `TestBudgetLoudTotals`; `TestCandidateCapSetsTruncated`; `TestIndexStaleBanner`; S01 `TestContextWhyTraceEdgeProvenance` stays green |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Tests (primary) | `internal/compiler/compiler_test.go` (+ focused `*_test.go` OK) | Add three named tests |
| Index honesty (only if residual) | `internal/compiler/index_honesty.go` | sort-then-cap if tests require |
| Packet / compiler | `packet.go` / `compiler.go` | Touch **only** if a named test proves a lock gap |
| CLI / MCP | Prefer **zero** edits | G19 — library marshal only |

## Role work
1. TDD: write the three named tests against **live** APIs first; expect red only where behavior truly missing.
2. Minimal product fixes only for failing locks (prefer not rewriting types that already match inventory).
3. Confirm `TestContextWhyTraceEdgeProvenance` stays green.
4. Run locked verify suite; board **status + Notes only** (no prompt/board spawn).

## Minimal todos
- [x] ~~Extend `Budget` + MD; bump SchemaVersion `0.2`~~ — **already shipped** (verify via tests; do not re-add)
- [x] ~~Wire `truncated` / `candidates_capped` / totals in `compileAtDepth`~~ — **already shipped** (prove with named tests)
- [x] ~~`index_honesty` helper + wire~~ — **already shipped** (prove + false-fresh case; fix sort-then-cap if needed)
- [ ] Named tests green: `TestBudgetLoudTotals`, `TestCandidateCapSetsTruncated`, `TestIndexStaleBanner`
- [ ] S01 `TestContextWhyTraceEdgeProvenance` green
- [ ] Carry-forward verify cmds below
- [ ] Board row P12-S02-01 Notes; next **P12-S02-02**

## Verify commands

```bash
CGO_ENABLED=0 go test ./internal/compiler/... ./internal/retrieval/... ./evals/honesty/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Named (min): `TestBudgetLoudTotals` `TestCandidateCapSetsTruncated` `TestIndexStaleBanner` `TestContextWhyTraceEdgeProvenance`

## Exit criteria
- [ ] FINAL locks met with named tests (product types may pre-exist — tests are the bar)
- [ ] Truncation is loud (totals + truncated/capped); no silent MaxCandidateHits
- [ ] Staleness banner when disk hash ≠ indexed hash; false-fresh preferred; causal STALE untouched
- [ ] S01 `edge_provenance` pass-through intact
- [ ] Carry-forward + Gate C `dry_run:false` untouched
- [ ] Board Notes; next **P12-S02-02**
