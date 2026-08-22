# P13 / S02 / 01 — Packet honesty residuals (FINAL)

## Metadata
- id: P13-S02-01
- todo_ids: [P13-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-61, DF-62, DF-63, DF-65** per sibling **00-PLANNER** FINAL locks. Loud stale totals, honesty under trim, admit-universe `items_total` under cap, and context import-hop `edge_provenance` via Expand on file seeds (reuse S01 resolve). **Stop if 00-PLANNER is still DRAFT.** Keep P12 named honesty tests green. **No** mig / analyzer rewrite / path-align hook.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- DOGFOOD-FINDINGS DF-61…63, DF-65
- Repros: `experiments/_bughunt/post-p12/{stalecap,staledrop,candcap,prov}/`
- Live: `internal/compiler/{packet.go,index_honesty.go,compiler.go,budget.go}`; Expand file→import in `internal/retrieval/expand.go`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do not re-debate FINAL locks.

## Locked defaults (FINAL — do not renegotiate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Home | `internal/compiler` primary; call existing `Retriever.Expand` for DF-65 file seeds |
| SchemaVersion | Keep **`0.2`** (additive JSON only) |
| DF-61 | `IndexHonesty.stale_total` + `stale_truncated`; keep sort-then-cap 8; MD shows total when truncated |
| DF-62 | Honesty universe = **pre-trim** file items (post Layer-1 / MaxCandidateHits), not `kept` alone; false-fresh on I/O miss unchanged |
| DF-63 | `items_total` = layer0 + unique L1-admissible hits in **full** Expand+FTS list (not MaxCandidateHits-truncated); pipeline admit cap unchanged; `candidates_capped` unchanged |
| DF-65 | After task Expand + FTS, Expand **file-typed** candidate hits (depth 1), merge, preserve `edge_provenance`; **reuse S01** resolve — no compiler path-join |
| Law 18 | Never set causal `Provenance.Status=STALE` from disk drift |
| Migration / analyzers / path-align | **None** |
| MCP / CLI | Prefer **zero** adapter edits |
| Forbidden | Silent totals; fake provenance; new MCP; daemon; full-rebuild; board spawn; invent SchemaVersion `0.3` |
| Keep green | `TestBudgetLoudTotals`, `TestCandidateCapSetsTruncated`, `TestIndexStaleBanner`, `TestContextWhyTraceEdgeProvenance`, S01 Expand/Why import tests |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Compiler | `internal/compiler/packet.go` | `IndexHonesty` fields; MD banner total signal |
| Compiler | `internal/compiler/index_honesty.go` | Pre-trim universe; stale_total / truncated before path cap |
| Compiler | `internal/compiler/compiler.go` | DF-63 universe count; DF-65 Expand file seeds + merge; pass pre-trim items to honesty |
| Tests | `internal/compiler/compiler_test.go` (+ optional focused `*_test.go`) | One named test (or suite) per DF; keep P12 names |
| Retrieval | Prefer **zero**; only if Expand merge helper is clearly cleaner | **No** resolve reimplementation |

## Role work
1. TDD per DF locks (61→62→63→65 or grouped honesty residual suite).
2. Wire DF-61/62 in `buildIndexHonesty` + MD; DF-63 in `compileAtDepth` totals; DF-65 Expand file seeds via existing Retriever.
3. Prove P12 keepers + S01 import Expand still green.
4. Run locked verify; board **status + Notes only** (no prompt/board spawn).

## Minimal todos
- [ ] DF-61: `stale_total` / `stale_truncated` + MD + named test (>8 stale)
- [ ] DF-62: honesty over pre-trim file items + named test (trim drops stale file → banner remains)
- [ ] DF-63: admit-universe `items_total` when capped + named test (universe ≫64)
- [ ] DF-65: Expand file seeds in `compileAtDepth` + named test (`edge_provenance` on context hop)
- [ ] P12 keepers + S01 import tests green
- [ ] Carry-forward verify cmds below
- [ ] Board row P13-S02-01 Notes; next **P13-S02-02**

## Named tests (intent locked)

| Test | Intent |
|------|--------|
| `TestBudgetLoudTotals` | **Keep** |
| `TestCandidateCapSetsTruncated` | **Keep** |
| `TestIndexStaleBanner` | **Keep** |
| `TestContextWhyTraceEdgeProvenance` | **Keep** |
| New DF-61 | >8 stale → cap 8 + `stale_total` + `stale_truncated` + MD total |
| New DF-62 | Trim drops disk-stale file → `index_honesty` non-null (not false-fresh null) |
| New DF-63 | Cap → `items_total` reflects full admissible universe; MD `items=k/t` |
| New DF-65 | Context packet admits import neighbor with `edge_provenance` (analyzer-shaped `./` + S01 resolve) |

## Verify commands

```bash
# S02 + P12 honesty (CGO0)
CGO_ENABLED=0 go test ./internal/compiler/... ./internal/retrieval/... ./evals/honesty/... -count=1

# Named focus (adjust -run to match new names)
CGO_ENABLED=0 go test ./internal/compiler/... -count=1 -run 'TestBudgetLoudTotals|TestCandidateCapSetsTruncated|TestIndexStaleBanner|TestContextWhyTraceEdgeProvenance|TestStale|TestCand|TestIndex|TestContext.*[Pp]rov|TestHonesty'

# S01 resolve still green
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestExpand.*[Ii]mport|TestWhy.*[Ii]mport'

# Carry-forward (CGO1)
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1

# Product packages
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Do **not** treat `go test ./...` space-path FAIL (known graphify) as S02 failure. Gate C `dry_run:false` untouched. Dry-run ≠ Gate C / ≠ H / ≠ checklist.

## Exit criteria
- [ ] FINAL locks met for DF-61/62/63/65
- [ ] P12 keepers + S01 import tests green
- [ ] No mig / analyzer rewrite / path-align; SchemaVersion stays `0.2`
- [ ] Carry-forward gates green; Gate C `dry_run:false` untouched
- [ ] Board Notes; next **P13-S02-02**

## Next
**P13-S02-02**
