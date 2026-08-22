# P07 / S02 / 02 — Scope review (language plugins)

## Metadata
- id: P07-S02-02
- todo_ids: [P07-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S02 Go language adapter. Confirm claims match repo + tests. Forward-only.

## Session start
Agent → clarify → Plan → review (fresh ≠ S02-01).

## Locked claims to verify (from S02-00/01)
| Claim | Evidence |
|-------|----------|
| Language = Go only (one adapter) | `DetectLanguage` + extract; no extra langs |
| Module `github.com/tree-sitter/tree-sitter-go` **v0.25.0** | `go.mod` |
| Extension points = DetectLanguage + `extract` switch + `extract_go.go` | Diff vs prior |
| Adapter-shaped (no plugin registry / universal theater) | No new iface packages |
| CGO analyzers-only | `CGO_ENABLED=0` store/domain/vcs/gitcli still pass |
| Golden Go symbols+imports | `internal/analyzers` tests |
| S01 T0 walk not regressed | `cmd/trace` T0 tests + walk order |
| No Gate H threshold invent | No pass numbers / Gate H claim |
| Carry-forward bars green | honesty/replan/impact/capability/p0x/x0/`./...` |
| Gate C untouched | packs not rewritten; `dry_run:false` intact |

## Review focus
- Official grammar + go 1.24 floor respected?
- IndexFile path consistent; unsupported paths still skipped?
- Symbol/import kinds match locked vocabulary (`function`/`method`/`type`)?
- `vendor` still T0-skipped (Go-friendly)?
- No mig `011_*` unless Notes justify?

## Depends / blast radius
- On APPROVE: unlock **P07-S03-00**.
- Light note for S03: Gate H may optionally include tiny `.go` ladder fixtures; **thresholds still after measurements** — S02 must not invent them.

## Exit criteria
- [ ] REVIEW-NOTES.md APPROVE/REJECT + confidence
- [ ] blocker/high fixed or spawned
- [ ] Next runnable unlocked only on APPROVE

## Out of scope
- Implementing Gate H / Phase 07 VERIFY
- Adding further languages
