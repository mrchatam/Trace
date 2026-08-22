# P14 / S01 / 01 — Impact walks (FINAL)

## Metadata
- id: P14-S01-01
- todo_ids: [P14-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement multi-seed impact BFS + contains-asymmetric radius per sibling **00-PLANNER FINAL** locks. Wire into existing `trace impact` surface via new `walk` subcommand. Keep planted Gate F + Expand context semantics green. **Stop if 00-PLANNER is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A3/A5/A6
- Research rank 6 — CBM / codegraph impact walks
- Live: `internal/retrieval/expand.go` (+ `resolveImportedFile`), `internal/domain/impact.go`, `cmd/trace/impact.go`, `evals/impact`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do not re-debate FINAL locks.

## Locked defaults (FINAL — do not renegotiate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Algorithm home | **`internal/retrieval`** — ImpactWalk (name free); **do not** alter Expand bidirectional file↔symbol for context/Why |
| CLI | **`trace impact walk`** — thin G19 adapter; keep `finding`/`alternative`/`report` |
| Seeds | `file`\|`symbol` only; ≥1; multi-seed = one BFS; seeds **excluded** from blast |
| Depth | 1..2; default **2**; reject outside range |
| Deps | **Incoming** imports only (reuse P13 resolve to match targets) |
| Contains | File→symbols OUT OK; symbol→file may enter blast for dep walk; **no** sibling symbols via contains climb |
| Loudness | Cap **64**; `blast_total` / `blast_kept` / `truncated`; hop ASC then type/id ASC before cap |
| Hop risk | `hop_risk = float64(hop)` (monotonic) |
| Migration / MCP / `internal/impact` | **None** |
| Domain planted APIs | Untouched (Gate F substrate) |
| Forbidden | Silent caps; Expand “fix”; outgoing-as-blast; auto-plant findings; new MCP; daemon/HTTP/embeddings; full-rebuild; board spawn |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Retrieval | `internal/retrieval/` new walk file + `types.go` as needed | `ImpactWalk` API + result type + hop_risk |
| Retrieval tests | `internal/retrieval/*_test.go` | Named ImpactWalk tests (CGO0 fixtures OK) |
| Store (optional) | `internal/store/file_graph.go` (+ test) | Reverse importer helper if cleaner than retrieval-only scan |
| CLI | `cmd/trace/impact.go` (+ `cli_test.go` / impact test) | `walk` subcommand + help; JSON stdout |
| Domain / mig / analyzers / compiler / MCP | Prefer **zero** | Do not touch unless Notes prove need |

### CLI sketch (informative — names may vary)

```text
trace impact walk --seed file:<uuid> [--seed symbol:<uuid> ...] [--depth 1|2]
# JSON: ok, seeds, blast[{entity_type,entity_id,path?,title?,hop,hop_risk,edge_provenance?}],
#        blast_total, blast_kept, truncated, depth
```

## Role work
1. TDD named tests first (multi-seed exclusion → contains asymmetry → incoming hop → truncation → hop_risk).
2. Implement walk; reuse `resolveImportedFile` for importer matching; leave Expand alone.
3. Thin `trace impact walk`; prove existing impact subcommands + Gate F still green.
4. Run locked verify; board **status + Notes only** (no prompt/board spawn).

## Minimal todos
- [ ] ImpactWalk API + seed validation (file\|symbol, depth 1..2)
- [ ] Multi-seed BFS + seed exclusion + deterministic order
- [ ] Contains asymmetry (no sibling climb)
- [ ] Incoming import hops + `edge_provenance` when present
- [ ] Loud truncation fields (cap 64)
- [ ] `hop_risk` monotonic with hop
- [ ] `trace impact walk` CLI JSON
- [ ] Named tests green; Gate F prelim green
- [ ] Carry-forward verify cmds below
- [ ] Board row P14-S01-01 Notes; next **P14-S01-02**

## Named tests (must ship)

| Test | Intent |
|------|--------|
| `TestImpactWalkMultiSeedExcludeSeeds` | Multi-seed one walk; seeds absent from blast |
| `TestImpactWalkContainsAsymmetryNoSiblings` | Symbol seed → no sibling symbols via contains |
| `TestImpactWalkIncomingImportHop` | Importer of seed file appears; provenance when set |
| `TestImpactWalkLoudTruncation` | &gt;64 nodes ⇒ truncated + totals |
| `TestImpactWalkHopRiskIncreases` | Deeper hop ⇒ ≥ hop_risk |
| `TestPlantedImpactConflictsGateFPrelim` | **Keep** |

## Verify commands

```bash
# Walk + retrieval (CGO0 OK for store fixtures)
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/store/... ./evals/impact/... ./evals/honesty/... -count=1

# CLI + carry-forward (CGO1)
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/impact/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1

# Product packages
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Also keep green (as available in tree): P12/P13 provenance + packet honesty named tests under `internal/retrieval` / `internal/compiler`; Gate C artifacts `dry_run:false` **untouched** (do not re-run Mode-B as a fix).

## Exit criteria
- [ ] FINAL locks met with named ImpactWalk tests above
- [ ] `trace impact walk` returns loud totals; seeds excluded; contains asymmetry holds
- [ ] Gate F prelim + existing `trace impact` finding/alternative/report green
- [ ] Expand/Why context semantics not regressively “fixed” by removing symbol→file
- [ ] Carry-forward gates green; Gate C `dry_run:false` untouched
- [ ] Board Notes; next **P14-S01-02**
