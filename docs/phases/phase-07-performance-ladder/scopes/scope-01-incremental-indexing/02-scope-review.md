# P07 / S01 / 02 — Scope review (incremental indexing)

## Metadata
- id: P07-S01-02
- todo_ids: [P07-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S01 incremental indexing / ignore tiers. Confirm claims match repo + tests. Severity-tag findings; small fixes or spawn `a`/`b` pairs. Forward-only.

## Session start
Agent → clarify → Plan → review (fresh ≠ S01-01).

## Review focus
- File-local incremental preserved (no full-rebuild architecture)?
- T0 locks match implement Notes: dir basenames (`node_modules`/`vendor`/`__pycache__`/`.venv`/`venv`/`dist`/`.next`/`target`/`coverage` + existing `.git`/`.trace`); `.min.js`/`.min.mjs`/`.min.cjs`; path-segment rule on walk **and** explicit argv?
- Walk order: T0 → DetectLanguage → T0 file/path → gitignore?
- Sibling isolation (`TestIndexIncrementalIsolation`) + new T0 tests green?
- **No** `011_*` by default (or justified Notes)?
- Optional `evals/perf` seed only — **no** Gate H pass / threshold invent?
- Carry-forward bars green (honesty A/B/C; Gate G/E/F; capability ablation; p0x; x0; Gate C `dry_run:false`)?
- No daemon/HTTP/embeddings; Gate C artifacts intact?

## Depends / blast radius
- On APPROVE: unlock **P07-S02-00**. Confirm S02 will consume T0 walk (do not regress). Light-note S03 if timing/fixture seeds appeared for Gate H.

## Exit criteria
- [ ] REVIEW-NOTES.md with APPROVE/REJECT + confidence
- [ ] blocker/high fixed or spawned
- [ ] TODO.md updated; next runnable unlocked only on APPROVE

## Out of scope
- Implementing S02 language plugins
- Declaring Gate H / Phase 07 complete
