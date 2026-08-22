# P30-S01-01 — Plan implementer

## Metadata
- id: P30-S01-01
- todo_ids: [P30-S01-01]
- role: implementer
- skills: [planning-and-task-breakdown, documentation-and-adrs]
- verification: automated

## Objective

Docs-only: write `PLAN.md` from [`../scope-00-investigate/INVESTIGATION.md`](../scope-00-investigate/INVESTIGATION.md) and the **locked defaults** below (hygiene track locked by `P30-S00-02` PASS). **Do not invent a path-fix task.**

## Protocol

Follow [`docs/rules/agent-loop-protocol.md`](../../../../../rules/agent-loop-protocol.md) Session start. Clarification: none required — S00 verdict and T1–T4 are locked. Produce a short mental plan, then write `PLAN.md` only.

## Verdict taken (cite verbatim in PLAN.md)

From INVESTIGATION.md:

> **agent hygiene** — Trace never creates or opens `<projectRoot>/trace.db`; every executable store join is `<root>/.trace/trace.db`, and a disposable repro shows a **0-byte** root stub from `python3` `sqlite3.connect('trace.db')` after `trace init`, while CLI continues to use the live `.trace/` store.

Also record: INTAKE **confirmed**; **no Trace dual-store bug**; **no store-path change**.

## Locked defaults (S01-00 — do not renegotiate in PLAN.md)

| Lock | Value |
|------|-------|
| Canonical store | `<projectRoot>/.trace/trace.db` via `store.Open` / `OpenExisting` |
| Path redesign | **Forbidden** (S00 overturn path closed) |
| Silent / default delete of root `trace.db` | **Forbidden** |
| Optional documented delete flag | **Out of Phase 30** unless PLAN explicitly defers as future-only note (do not implement in S02) |
| Track | **T1–T4 only** (docs / warn / gitignore / tests) |
| HTTP / daemon / GUI | **Out of scope** (Phase 29 owns `trace serve`) |
| Green bar (S02) | `go test ./internal/...` PASS |

### T2 warn design (locked)

| Item | Value |
|------|-------|
| When | After project root is Abs-resolved, if `<root>/trace.db` exists as a **regular file**, while opening the canonical store |
| Where | `internal/store` open path (`Open` / `OpenExisting` / shared `openStore` — PLAN picks the single choke that covers CLI+MCP+serve without duplicating) |
| Frequency | **Once per successful open attempt** that sees the stray (acceptable: once per `Open`/`OpenExisting` call). Not a persistent on-disk flag. |
| Channel | **stderr** only (or package `io.Writer` injectable for tests); never stdout JSON pollution |
| Severity | Non-fatal — open **must still succeed**; never fail-closed on stray presence |
| Message intent | Root `trace.db` is **not** the Trace store; canonical is `.trace/trace.db`; do not open/create root stub (agents: use CLI/MCP) |
| Must not | Delete, rename, migrate, or open the root file; change `dbPath` join; treat stub as second store |

### T3 gitignore (locked)

| Item | Value |
|------|-------|
| Pattern | Root-only ignore: `/trace.db` (leading slash) |
| Must not | Un-ignore or commit `.trace/`; pattern must not target `.trace/trace.db` |
| Surfaces | This repo `.gitignore` **and** any install/scaffold guidance that tells consumers what to ignore (PLAN lists concrete files) |

### T1 docs surfaces (candidate set — PLAN narrows to minimal sufficient)

Prefer the smallest set that reaches operators + agents:

- `AGENTS.md` (stack / store one-liner)
- `docs/rules/project-rules.md` (Store row)
- `CONTRIBUTING.md` (`.trace/` local store section)
- `cmd/trace/help.go` (init / store path clarity if a one-liner fits)
- Install / agent-rule text under `internal/install/*` **if** a bundled rule/doc string is the durable agent surface (INVESTIGATION noted gap: install mentions `.trace/config.json` but not “never root `trace.db`”)

Do not rewrite Phase 29 HTTP docs for this hygiene.

### T4 tests (locked acceptance shape)

| Case | Expect |
|------|--------|
| Stub present + `Open` / `OpenExisting` | Warn observed (stderr or test writer); store path still `.trace/trace.db`; open succeeds |
| No stub | No warn; open succeeds |
| After warn | Live DB under `.trace/`; root stub untouched (still present; size unchanged) |

## Deliverable: `PLAN.md`

Path: `scopes/scope-01-plan/PLAN.md` (this folder).

### PLAN.md template (fill every section)

```markdown
# PLAN — Phase 30 stray root trace.db hygiene

**Date:** YYYY-MM-DD · **Author:** P30-S01-01 · **Depends:** INVESTIGATION.md + P30-S00-02 PASS

## 1. Verdict taken

Quote S00 one-liner. State: INTAKE confirmed; no dual-store bug; no path redesign.

## 2. Goals / non-goals

### Goals
- Operator + agent clarity: only `.trace/trace.db` is authoritative
- Ship T1–T4

### Non-goals (explicit)
- Store-path redesign / “path fix”
- Silent or default delete of root `trace.db`
- Treating root stub as a second store
- HTTP/daemon/GUI scope creep
- New phase for Trace dual-store bug (none found)

## 3. Tasks

| ID | Change | Files (likely) | Acceptance | Deps |
|----|--------|----------------|------------|------|
| T1 | … | … | … | — |
| T2 | … | … | … | — |
| T3 | … | … | … | T1 or parallel if file-disjoint |
| T4 | … | … | … | T2 |

## 4. Warn design (T2)

Document: trigger, once-per-open, stderr, never fail-closed, never deletes, never changes open path, proposed message text (draft OK).

## 5. Implementation order for S02

Numbered sequence (suggest: T1 docs → T3 gitignore → T2 warn → T4 tests, or T2+T4 before docs if tests drive API). Note parallel-safe file sets.

## 6. Test plan

Commands + named cases from locked T4 shape. Require `go test ./internal/...` green before S02-01 done.

## 7. Risks / residuals

Low-risk hygiene only. List anything deferred (e.g. documented delete flag = future, not S02).

## 8. S02 handoff

One paragraph: what S02-00 must lock; point at this PLAN as SoT.
```

## Acceptance hooks (PLAN.md is done when)

- [ ] Section 1 quotes INVESTIGATION verdict (agent hygiene + INTAKE confirmed)
- [ ] T1–T4 each have: ID, files likely touched, acceptance criteria, deps
- [ ] T2 section matches locked warn design (once-per-open, stderr, non-fatal, no delete, no path change)
- [ ] T3 specifies `/trace.db` root-only ignore + which files change
- [ ] Out-of-scope list includes path redesign, silent delete, HTTP/daemon, second-store treatment
- [ ] Explicit non-goal: **no “path fix” task**
- [ ] S02 can execute without re-opening S00

## Do not

- Implement product code in this row (no Go edits; no writing real `.gitignore` / AGENTS changes — that is S02)
- Write only a thin stub PLAN — use the template fully
- Change store path “just in case”
- Re-open S00 investigation as if the verdict were unsettled
- Auto-delete design as default behavior

## Exit criteria

- [ ] `PLAN.md` present and actionable for S02
- [ ] Board Notes point at `PLAN.md` + restate verdict one-liner
- [ ] Next: **P30-S02-00** (no separate S01 review row — S02 planner re-reads PLAN)

## Next

`P30-S02-00`
