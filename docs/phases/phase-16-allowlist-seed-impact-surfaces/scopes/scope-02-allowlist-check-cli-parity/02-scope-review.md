# P16 / S02 / 02 — Scope review (allowlist CHECK + CLI parity)

## Metadata
- id: P16-S02-02
- todo_ids: [P16-S02-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S02 DF-75/77/78. Spawn `P16-S02-02a`/`02b` on blocker/high. Next **P16-S03-00** unless spawn.

## Checklist
| # | Check |
|---|--------|
| 1 | CHECK rejects YOLO / garbage writes |
| 2 | Resolve never fail-opens builtins on garbage |
| 3 | CLI add/why honor DENIED `mcp:` slug |
| 4 | Unprefixed decide slug gates MCP (or equivalent FINAL lock) |
| 5 | `capability decide` remains operator-ungated |
| 6 | Compat ceiling matches mig; no YOLO default |
| 7 | Carry-forward + product pkgs; S01 keepers still green |

## Exit criteria
- [ ] Checklist evidenced; REVIEW-NOTES; next **P16-S03-00** unless spawn
