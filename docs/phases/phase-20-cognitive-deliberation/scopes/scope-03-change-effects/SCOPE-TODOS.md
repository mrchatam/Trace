# S03 — Change + effects — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P20-S03-00 | scope planner | done — 017 schema + `change_paths` + comparison enum FINAL |
| 2 | P20-S03-01 | implementer | done |
| 3 | P20-S03-02 | reviewer | done — APPROVE high; next P20-S04-00 |

**Locks (do not re-debate in 01):** migration **017**; tables `changes` / `change_paths` / `effects`; paths = **table not JSON**; comparison `supported`\|`partially_supported`\|`contradicted`; Git SHA+path refs only (`ShowFile` at read); contradicted may Discovery **or** Hypothesis (not Discovery-as-hypothesis); no Regression this scope; library-only; compat ceiling **17** after 01.

**Feeds:** S04 (no tests/baseline columns on changes); S05 regression inputs; S06 apply `changes`/`effects` writes + bounded `recent_changes[]`.
