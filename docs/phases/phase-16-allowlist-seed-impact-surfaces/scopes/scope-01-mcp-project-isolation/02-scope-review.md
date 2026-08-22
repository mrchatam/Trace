# P16 / S01 / 02 — Scope review (MCP project isolation)

## Metadata
- id: P16-S01-02
- todo_ids: [P16-S01-02]
- role: reviewer
- skills: [code-review-and-quality, systematic-debugging]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S01 DF-76. Fresh subagent ≠ implementer. Spawn `P16-S01-02a`/`02b` for blocker/high. Prefer `REVIEW-NOTES.md`. Next **P16-S02-00** unless spawn.

## Checklist
| # | Check |
|---|--------|
| 1 | MCP `project=` on virgin dir does not mkdir `.trace/` / AUTO_ALLOW |
| 2 | Initialized `project=` override still works (per-store SoT) |
| 3 | CLI `store.Open` auto-init unchanged |
| 4 | P15 Assert named tests still pass; nine tools |
| 5 | Carry-forward honesty/E–H/ablation/compat/p0x/x0 + product pkgs |
| 6 | No S02–S05 work smuggled in |

## Exit criteria
- [ ] Checklist evidenced; confidence high (or medium with residuals listed)
- [ ] Board Notes; next **P16-S02-00** unless spawn
