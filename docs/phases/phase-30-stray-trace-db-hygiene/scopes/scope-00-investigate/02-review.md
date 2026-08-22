# P30-S00-02 — Investigation review

## Metadata
- id: P30-S00-02
- todo_ids: [P30-S00-02]
- role: reviewer
- skills: [code-review-and-quality, diagnosing-bugs]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of `INVESTIGATION.md` and S00-01 evidence. Fresh session — do not share the investigate implementer’s context as authority. Do not treat INTAKE as proven; judge whether S00-01’s evidence supports its verdict.

## Session start

Follow agent-loop-protocol Session start. Re-read board row Notes for P30-S00-01.

## Preflight

- [ ] `scopes/scope-00-investigate/INVESTIGATION.md` exists
- [ ] No unexpected product code changes attributed to S00 (git status / diff skim)
- [ ] Phase still Phase 30; next after this row should be **P30-S01-00** if PASS
- [ ] Spot-check live anchors still match citations (open.go / mcp/project.go / httpapi/server.go / help.go)

## Checklist

### Verdict & coverage
- [ ] Verdict explicit: **Trace bug** | **agent hygiene** | **both**
- [ ] All five must-answer questions from `01-investigate.md` are answered (create root? accidental open? stub repro? open.go join cites? P29 HTTP root write/open?)
- [ ] INTAKE relation stated (confirmed / partial / overturned)

### Evidence quality
- [ ] Store path in code cited with file:line (`open.go` join under `.trace/` — expect ~L15–16 constants; OpenExisting ~L57; openStore ~L77/L92)
- [ ] Search coverage: store, cmd/trace, mcp, httpapi (install optional)
- [ ] Repro of 0-byte (or empty) root stub confirmed **or** refuted with commands + observed size
- [ ] CLI `-C` / MCP `project` / HTTP `store.Open(s.root)` path risk addressed honestly
- [ ] HTTP section distinguishes: opens `.trace/trace.db` via `store.Open` vs static-dir refuse of project root (not a root-db creator)
- [ ] Recommendations for S01 are suggestions only (no silent path redesign; path fix only if Trace bug proven)

### Process
- [ ] No product code in this scope
- [ ] Implementer did not delete dogfood artifacts

## Severity → spawn

| Severity | Action |
|----------|--------|
| blocker / high | Insert `P30-S00-02a` (fix) + `P30-S00-02b` (re-review) **immediately below** this row on `docs/TODO/phase-30.md`; thicken prompts under `scopes/scope-00-investigate/` |
| medium / low | Prefer inline doc fix in INVESTIGATION.md if clear; else note for S01 |
| nit | Notes only |

**Spawn policy:** `02a`/`02b` only on blocker/high (missing verdict, wrong path claim without cite, unreproduced stub treated as fact, or product code slipped into S00). Do not spawn for preference nits.

## Exit criteria

- [ ] Review PASS or FAIL with reason in board Notes
- [ ] If PASS: thicken **upcoming** S01 prompts lightly if verdict changes the default plan (path fix first vs hygiene-only)
- [ ] Next runnable **P30-S01-00** (or `P30-S00-02a` if spawned)

## Todo updates

Status + notes on **P30-S00-02**; may thicken upcoming S01 only.

## Next

`P30-S01-00` (unless 02a/02b spawned)
