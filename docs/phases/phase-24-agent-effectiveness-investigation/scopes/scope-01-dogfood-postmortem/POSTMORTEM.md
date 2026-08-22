# E01 dogfood post-mortem — Session A vs Session B

Phase 24 / Scope 01 / P24-S01-01 deliverable. Evidence from required reads + command block (2026-08-20).

---

## §1 Runs reviewed

### E01 Session A — Build (required)

**Date:** 2026-08-20. **Workspace:** `experiments/ab-incident-tracker/runs/G1/`. **Trace binary:** `/home/ali/Desktop/Trace/bin/trace` ([PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) L5).

**Harness:** Multitask parent orchestrator + worker subagents. Prompts: [PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) (enforcement arm), [SUBAGENT-DELEGATION.md](../../../../../experiments/ab-incident-tracker/prompts/SUBAGENT-DELEGATION.md) (packet format), feature parity bar [PROMPT-B0.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-B0.md). Parent told to own Trace graph; workers gate on `TRACE_TASK_ID`. Cursor hook + rule installed under G1 (`.cursor/hooks/trace-loop-gate.sh`, `.cursor/rules/trace-enforcement.mdc`); config `{ "enforce": "strict" }` ([`.trace/config.json`](../../../../../experiments/ab-incident-tracker/runs/G1/.trace/config.json)).

**What happened:** Parent delegated seed tasks `…0010`–`…0050` ([gt.json](../../../../../experiments/ab-incident-tracker/runs/G1/seed/gt.json) L10–40) via Multitask. Full incident tracker shipped with tests passing; scored verdict: seed-only graph, 0 decisions, 0 discoveries, G1 product **≡ B0**, loop STOP / saturated ([experiments/RESULTS.md](../../../../../experiments/RESULTS.md) E01 row). **Git sparsity:** init commit `f70aaea` contains only 4 starter files (`.gitignore`, `README.md`, `cmd/incidentd/main.go`, `go.mod`) — no `internal/` at anchor SHA (`git show f70aaea --stat`). Session A product evidence is RESULTS + B0 parity claim, not intermediate commits.

### E01 Session B — Directed gap (required)

**Date:** 2026-08-20 (same workspace, continuation). **Harness:** Single-agent session; human prompt: gap analysis + plan + fix **using Trace** ([INVESTIGATION.md](../../INVESTIGATION.md) L14; [FINDINGS.md](../../FINDINGS.md) L29–42).

**What happened:** Agent recorded **7 discoveries**, **2 decisions**, **4 evidence** in terminal `trace/graph.json`; linked 2 discoveries to verify task `…0050` via `discovery_mentions_task`; fixed unassign, assignee filter, public `started_at`, store/web tests; filled [VERIFY.md](../../../../../experiments/ab-incident-tracker/runs/G1/VERIFY.md) manual table (L35–51). Product **diverged from B0** (9 differing paths under `internal/` at HEAD — `diff -rq` command output). Git range `704e2ff`…`a37e7c0` (6 commits after first gap-closure commit). Still **5 seed tasks only**; 0 uncertainties; all deliberation states STOP.

### D44 — Phase 23 / ab-cms-fullstack (optional)

**Not deeply reviewed.** Skim only: [Phase 23 README](../../../phase-23-enforcement-choke-points/README.md) L11 notes D44 showed Trace wins on planning quality when G1 agents used loop + seed, but B0/G-prompt arms could still implement prematurely. No D44 row in [experiments/RESULTS.md](../../../../../experiments/RESULTS.md).

### D45 (optional)

**Not reviewed** — no D45 row in [experiments/RESULTS.md](../../../../../experiments/RESULTS.md).

---

## §2 Two-mode comparison table

| Dimension | Session A (build) | Session B (directed gap) | Evidence pointer |
|-----------|-----------------|--------------------------|----------------|
| Tasks in graph | 5 seed only (`…0010`–`…0050`) | 5 seed only — no new task IDs | [graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L10–40; python count `tasks 5` |
| Discoveries / decisions / evidence | 0 / 0 / 0 (Session A period) | 7 / 2 / 4 | [graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L42–137; [RESULTS.md](../../../../../experiments/RESULTS.md) E01 row |
| Uncertainties | 0 | 0 | [graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) — no `uncertainties` key; python `uncertainties 0` |
| Loop / gate on verify `…0050` | STOP / `p19_saturated` (export); live gate N/A if DB empty | STOP / `p19_saturated`; hop_count=1; live CLI `task_not_found` | [graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L344–353; live `$TRACE_BIN loop gate --task …0050 --for edit` → `reason_code":"task_not_found"` |
| Product vs B0 | **≡ B0** (scored parity) | **Diverged** — unassign, filters, `started_at`, `store_test.go`, template deltas | [RESULTS.md](../../../../../experiments/RESULTS.md); `diff -rq …/B0/internal …/G1/internal` (9 files differ at HEAD) |
| Orchestrator pattern | Multitask parent + SUBAGENT-DELEGATION workers | Single-agent directed gap session | [PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) L81–90; [SUBAGENT-DELEGATION.md](../../../../../experiments/ab-incident-tracker/prompts/SUBAGENT-DELEGATION.md) |
| Tests / VERIFY | `go test ./...` PASS (build) | PASS + manual table complete | [VERIFY.md](../../../../../experiments/ab-incident-tracker/runs/G1/VERIFY.md) L35–51; evidence entity L117–121 |
| SPEC / planning docs | Produced during build (in working tree) | Updated for gap deltas (unassign, filters, redirect) | [SPEC.md](../../../../../experiments/ab-incident-tracker/runs/G1/SPEC.md) L57, L76–79; [PLANNING-MATRIX.md](../../../../../experiments/ab-incident-tracker/runs/G1/PLANNING-MATRIX.md) |
| Git boundary | Anchor `f70aaea` (starter only in git) | `704e2ff`…`a37e7c0` (7 commits total log) | `git -C runs/G1 log --oneline`; [704e2ff stat](../../../../../experiments/ab-incident-tracker/runs/G1) via `git show 704e2ff --stat` |
| Harness enforcement | Hook + strict config present | Same harness; agent recorded graph entities | [AGENTS.md](../../../../../experiments/ab-incident-tracker/runs/G1/AGENTS.md); [trace-loop-gate.sh](../../../../../experiments/ab-incident-tracker/runs/G1/.cursor/hooks/trace-loop-gate.sh) |

---

## §3 Failure mode matrix (FM-01..FM-10)

| FM | Name | Session A | Session B | Both | Status | Evidence |
|----|------|-----------|-----------|------|--------|----------|
| FM-01 | Seed anchoring | Y | Y | Y | **confirmed** | [gt.json](../../../../../experiments/ab-incident-tracker/runs/G1/seed/gt.json) fixed UUIDs; [graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L10–40 — 5 tasks, all seed IDs |
| FM-02 | Graph thin export | Y | partial | partial | **partial** | A: 0 decisions/discoveries ([RESULTS.md](../../../../../experiments/RESULTS.md)); B: 2 decisions + 4 evidence but 0 uncertainties ([graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L42–137) |
| FM-03 | Loop saturation | Y | Y | Y | **confirmed** | All tasks `current_phase: STOP`, `stop_reason: p19_saturated` ([graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L299–353); coding continued ([VERIFY.md](../../../../../experiments/ab-incident-tracker/runs/G1/VERIFY.md) PASS) |
| FM-04 | Orchestrator bypass | Y | N | N | **confirmed (A)** | [PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) L81–90 vs workers-only Trace; B single-agent used Trace for discoveries |
| FM-05 | Enforcement optional | partial | partial | partial | **partial** | Hook allows when no `TRACE_TASK_ID` ([trace-loop-gate.sh](../../../../../experiments/ab-incident-tracker/runs/G1/.cursor/hooks/trace-loop-gate.sh) L5–7); strict config set but product edits proceeded |
| FM-06 | Cross-arm leakage (G1≡B0) | Y | N | N | **confirmed (A)** | [RESULTS.md](../../../../../experiments/RESULTS.md) Session A ≡ B0; HEAD `diff -rq` shows 9 internal diffs (Session B divergence) |
| FM-07 | Post-hoc planning | Y | partial | partial | **partial** | SPEC/PLANNING-MATRIX appear in gap-closure commit `704e2ff` with full product (`git show 704e2ff --stat`); graph export timestamps post-code |
| FM-08 | Tool surface gap | Y | Y | Y | **confirmed** | [PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) L75–77 says prefer `trace add` task; B recorded discoveries without task promotion |
| FM-09 | Mode-dependent effectiveness | Y | partial | partial | **confirmed** | [INVESTIGATION.md](../../INVESTIGATION.md) L26–27; B rich graph only after human directed gap |
| FM-10 | Discovery without task promotion | N | Y | N | **confirmed (B)** | 7 discoveries, 2× `discovery_mentions_task` → `…0050` ([graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L165–172); tasks[] unchanged |

### Command outputs (cited)

```text
# git -C runs/G1 log --oneline
a37e7c0 Export trace/graph.json after loop 1 review PASS.
a6d91ee Document env vars and add assignee filter/logout tests.
4f7e61e Re-export trace/graph.json at HEAD after graph commit.
44d4657 Export trace/graph.json after form POST test fix and gap discoveries.
337d2ca Fix form POST tests to capture redirects via noRedirectFrom.
704e2ff Close G1 gap analysis: unassign, filters, public started_at, tests.
f70aaea init G1 starter (E01 incidentops)

# python graph counts
tasks 5 | discoveries 7 | decisions 2 | evidence 4 | uncertainties 0
0010–0050: STOP p19_saturated

# diff -rq B0/internal G1/internal (HEAD) — Session B divergence
9 files differ (domain, store, web, templates; G1-only store_test.go)

# git diff f70aaea..704e2ff --stat
25 files changed, 3902 insertions (full product + docs + graph land in gap-closure commit)

# live loop gate (2026-08-20 S01 run)
loop status: task "…0050": sql: no rows in result set
loop gate --for edit: reason_code task_not_found (exit 1)
Note: exported graph deliberation shows p19_saturated, not hop_budget_exceeded (hop_count=1 on …0050).
FINDINGS/README hop_budget label is imprecise vs export; S02 should audit SelectNext reason codes.
```

---

## §4 Open questions (handoff to S02/S04)

- **Why 7 discoveries but 0 new tasks?** Session B used discovery entities + in-place fixes on seed tasks; `discovery_mentions_task` links point at verify `…0050` rather than spawning work items. Seed UUID table in [SUBAGENT-DELEGATION.md](../../../../../experiments/ab-incident-tracker/prompts/SUBAGENT-DELEGATION.md) L44–52 may anchor closed roster despite L5 disclaimer.
- **Was `trace add` (task type) invoked in Session B?** Export shows only seed task IDs ([graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L10–40). No non-seed task UUIDs. `.trace/trace.db` entity audit deferred — `sqlite3` unavailable in S01 environment.
- **Why verify gate blocked after green tests?** Export: all tasks STOP with `p19_saturated` ([graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L344–353). Live CLI returned `task_not_found` (SQLite out of sync with export). Reconcile `hop_budget_exceeded` vs `p19_saturated` in S02 (`internal/deliberation`, hop_count thresholds).
- **Directed gap: prompt hack vs product entry point?** Session B required explicit human “gap analysis + fix with Trace” ([INVESTIGATION.md](../../INVESTIGATION.md) L14). Trace MCP/CLI surface was **sufficient to record** gaps once directed; default build prompt did not produce equivalent graph richness ([RESULTS.md](../../../../../experiments/RESULTS.md)).
- **Git sparsity:** Session A product not reconstructable from `f70aaea` alone — experiment scoring relied on working-tree parity with B0. Protocol fix for E02+?
- **Deliberation reset:** Gap fixes did not transition verify task out of STOP — should gap closure reset hop_count or reopen EXECUTE phase?

---

## Must answer (exit gate)

### 1. Why 7 discoveries but 0 new tasks?

Agents treated each gap as **discovery + direct code fix** under existing seed tasks, not as backlog expansion. [PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) L75–77 recommends `trace add` task after discovery, but Session B human prompt emphasized closing gaps, not replanning the task tree. Two discoveries link to verify task `…0050` ([graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L165–172) instead of new UUIDs. FM-01 seed anchoring + FM-10 discovery-without-promotion explain the pattern.

### 2. Did Session B call `trace add` for tasks?

**No evidence of new tasks.** `tasks[]` length = 5; every `id` matches [gt.json](../../../../../experiments/ab-incident-tracker/runs/G1/seed/gt.json) seed UUIDs. Discoveries use new UUIDs (e.g. `5f65db6d-…`) — so `trace add` (or MCP equivalent) was used for **discoveries**, not **tasks**. DB audit blocked (no sqlite3).

### 3. Why verify task `…0050` gate-blocked after PASS tests?

**Export evidence:** `current_phase: STOP`, `stop_reason: p19_saturated`, `hop_count: 1` ([graph.json](../../../../../experiments/ab-incident-tracker/runs/G1/trace/graph.json) L344–353). Gap work and test PASS ([VERIFY.md](../../../../../experiments/ab-incident-tracker/runs/G1/VERIFY.md) L52) did not clear deliberation state or transition task to DONE. **Live CLI (S01):** `task_not_found` — `.trace/` SQLite missing seed tasks at investigation time; gate could not run against exported state. Prior docs citing `hop_budget_exceeded` ([FINDINGS.md](../../FINDINGS.md) L39) do not match export (`hop_count=1`); likely conflation of STOP family — S02 to confirm SelectNext mapping.

### 4. FM session split (Mode A only / Mode B only / both)

| Scope | FMs |
|-------|-----|
| **Session A only** | FM-04 (orchestrator bypass), FM-06 (G1≡B0) |
| **Session B only** | FM-10 (discovery without task promotion) |
| **Both** | FM-01, FM-03, FM-08 |
| **A strong / B partial** | FM-02, FM-05, FM-07, FM-09 |
| **Neither primary** | — |

### 5. Directed gap — prompt hack or sufficient entry point?

**Hybrid:** Trace product surface **can** record discoveries, decisions, and evidence when the human names gap analysis ([INVESTIGATION.md](../../INVESTIGATION.md) working-as-intended table L32–36). That is **not** the default build path (Session A thin graph). So directed gap is partly a **prompt/harness entry point** (FM-09), not a missing capability — but product defaults (loop STOP, seed anchoring, no auto task spawn) prevent build mode from reaching parity without human steering.
