# P18-00 — Plan D40 residuals: FTS + clone honesty (FINAL)

## Metadata
- id: P18-00
- todo_ids: [P18-00]
- role: planner
- skills: [planning-and-task-breakdown, documentation-and-adrs, writing-for-agents, grilling]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Against live repo + D40 / DOGFOOD-FINDINGS, lock a **FINAL** thin Phase 18: fix **DF-87** (context FTS MATCH syntax), **document DF-88** (keep P17 exclude; clone PENDING expected), golden **DF-89** (Go handler methods). **No product Go on this row.** Do **not** edit Phase 17 product Go or rewrite P17 `done` rows. Do **not** steal closed DFs. Next free after DF-88 was **DF-89** (this phase owns it).

## References
- [docs/rules/agent-loop-protocol.md](../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../init/G_PROJECT_LAWS.md) — Laws 1, 6, 7, 9, 11, 13
- [phase README](README.md)
- [DF-88-DECISION.md](DF-88-DECISION.md)
- P17 exclude SoT: [../phase-17-portable-graph-git/DF-84-FORWARD.md](../phase-17-portable-graph-git/DF-84-FORWARD.md) (history; do not rewrite)
- Findings: [experiments/DOGFOOD-FINDINGS.md](../../../experiments/DOGFOOD-FINDINGS.md)
- D40: `experiments/runs/2026-08-17-ab-compare/`
- Live: `internal/store/fts.go` `sanitizeFTSQuery`; `internal/compiler/compiler.go` TaskContext Search; `internal/analyzers/extract_go.go`; `cmd/trace/seed.go` export exclude
- [docs/TODO.md](../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). **Unattended:** human scheduled Phase 18 + DF-88 grill preference (do not explode portable-graph identity). Defaults below are **FINAL**.

## Live confirmation (2026-08-17)

See [README.md](README.md) investigation table. Count boarded: **DF-87 fix**, **DF-88 wontfix+docs**, **DF-89 golden**. Next free after this phase: **DF-90**.

## Disposition matrix (FINAL)

See [README.md](README.md). **Boarded:** S01→S03 + S04 VERIFY. **Not boarded:** hosted MCP, DF-86, CANCELLED state, app 400/429, harness rsync/stdio EOF, why-by-name CLI.

## Locked defaults (FINAL — phase)

| Item | Value |
|------|-------|
| Phase | Thin D40 residuals: FTS sanitizer + clone-PENDING honesty + Go method golden |
| History | Do not rewrite Phase 00–17 `done` prompts; P17 DR-HANDOFF stays `no successor` |
| Product Go | **Forbidden** on P18-00 |
| DF-87 fix | `sanitizeFTSQuery`: punctuation (incl. `/`) → token separators; quote remaining FTS5 tokens; keep `MATCH ?`. TaskContext on `GET /notes` / `GET /notes/search` must return a packet. Search error must **not** abort the packet (Expand-only fallback) |
| DF-87 not | “Add SQL parameters” (already bound); slash-only escape; phrase-quote the whole title as the only strategy |
| DF-88 | **Keep exclude.** No `--include-reviews`. Docs+help: clone tasks import PENDING. [DF-88-DECISION.md](DF-88-DECISION.md) |
| DF-89 | Analyzer golden for `method:Search` + `method:SearchCursor` on handler/store-shaped Go. Fix extract **only if** golden red. No why-by-name CLI |
| MCP | **No** new tools |
| G19 | FTS sanitizer in `internal/store` (library); compiler uses Search; CLI/MCP unchanged except S02 help strings |
| DR-HANDOFF intent | After VERIFY: default **`no successor`** |
| Forbidden | Daemon; hosted MCP; reversing P17 export exclude; DF-86 hook; kitchen-sink harness |

## Scope order (locked)

1. **S01 fts-query-sanitize** — DF-87
2. **S02 clone-pending-honesty** — DF-88 (docs + help + omit-test keeper)
3. **S03 go-method-extract** — DF-89
4. **S04 VERIFY** — named S01–S03 + carry-forward; DR-HANDOFF `no successor`

Board order is sequential (protocol).

## Non-goals
- Product Go in this planner row
- CONDITIONAL review/`work_state` export
- Hosted MCP; DF-86; CANCELLED work_state
- Claiming P17 close was wrong
- Auto-boarding research S05 / `plan simulate` / D21+

## Planner work (this row)
1. [x] Investigate FTS, seed exclude, Go extract, experiments FAIL/SKIP
2. [x] Disposition matrix FINAL (GO)
3. [x] Create scope folders + `00-PLANNER`/`01`/`02` stubs + SCOPE-TODOS + DR-HANDOFF stub + DF-88-DECISION
4. [x] Board P18-* **after** Phase 17 table, **before** Later developments; P18-00 done; next **P18-S01-00**
5. [x] AGENTS.md current focus; PROJECT_DOCS_INDEX; DOGFOOD-FINDINGS scheduled → Phase 18

## Exit criteria
- [x] Disposition matrix FINAL for DF-87/88/89
- [x] Board updated after P17 block
- [x] 00-PHASE-PLANNER marked **FINAL**
- [x] Notes: next **runnable** = **P18-S01-00**
- [ ] Product Go — **not** this row

## Minimal todos
- [x] Investigate + lock DF-87 root cause and DF-88 decision
- [x] Spawn S01–S04 stubs
- [x] Board + README + index + findings + AGENTS

## Next
Orchestrator: **P18-S02-00** (S01-02 **APPROVE**). Then S02→S05.

## Forward (post-FINAL refine, 2026-08-17)

Spawn-forward **S05 rebuild-binaries** after S04 VERIFY (`P18-S05-00` / `01` / `02`). Lesson: stale `bin/trace-mcp` caused experiment SKIP. Rebuild `bin/trace` (`CGO_ENABLED=1`) + `bin/trace-mcp` (`CGO_ENABLED=0`). Confirm `trace-mcp -h` lists 10 tools including `trace_impact`. DR-HANDOFF default remains **`no successor`**; **close moves to P18-S05-02**. S04 VERIFY is product evidence only. S05 may correct README/`help.go` MCP build lines to CGO0. Run order after S04: **S05** (not a new product phase; ≠ research S05). This section does not reopen DF-88 or rewrite the FINAL locks above.
