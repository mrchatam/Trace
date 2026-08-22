# P24-S02-01 — Codebase loop audit

## Metadata
- id: P24-S02-01
- todo_ids: [P24-S02-01]
- role: implementer (investigator)
- skills: [code-explorer, debugging-and-error-recovery, documentation-and-adrs, graphify]
- mcps: [Read, Glob, Grep, Shell, user-codegraph, user-trace]
- agents: [explore, investigator]
- verification: manual (file:line citations)
- hooks: none

## Objective

Trace **why the product allows or encourages** E01 failure modes (FM-01..FM-10). Write **`CODEBASE-AUDIT.md`** mapping each FM to concrete mechanisms in the live repo — file:line, agent-visible symptom, and change lever (product vs harness vs protocol).

**Investigation only** — no product Go commits. Optional ≤30-line repro comment block inside `CODEBASE-AUDIT.md` only.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [INVESTIGATION.md](../../INVESTIGATION.md) — FM taxonomy; investigation questions A–D
- [FINDINGS.md](../../FINDINGS.md) — two-mode model; S02 codebase section owner
- S01 deliverable: [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md) — §3 FM matrix, §4 open questions, Must answer
- S02-00 locks: [00-PLANNER.md](./00-PLANNER.md)
- Phase 23 enforcement baseline: [ENFORCEMENT.md](../../../phase-23-enforcement-choke-points/ENFORCEMENT.md)

## Session start

Follow [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md): Agent mode → clarify if FM scope or path inventory unclear → Plan mode → execute. **Do not re-litigate Session A vs B** — use POSTMORTEM as evidence; this scope explains **product code paths**.

## Locked defaults

| Item | Value |
|------|-------|
| Output file | `scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md` |
| FM coverage minimum | **≥8** of FM-01..FM-10 with file:line (target all 10) |
| Table columns | FM-ID \| mechanism \| file:line \| agent-visible symptom \| change lever |
| Change lever taxonomy | `product` \| `harness` \| `protocol` \| `experiment` (per INVESTIGATION intervention categories) |
| Product Go | **Forbidden** — docs + optional comment block in audit only |
| E01 workspace | `experiments/ab-incident-tracker/runs/G1/` |
| Verify task UUID | `e0100000-0000-4000-8000-000000000050` |
| Trace binary (live gate) | `/home/ali/Desktop/Trace/bin/trace` or `go run ./cmd/trace` from repo root |
| Hop budget constant | `deliberation.HopBudget = 12` (`internal/deliberation/types.go`) |

## Codebase audit map (locked inventory — verified 2026-08-20)

Read these packages **in order** when tracing FM mechanisms. Paths confirmed in live repo.

| Area | Primary paths | FM relevance |
|------|---------------|--------------|
| **Loop CLI** | `cmd/trace/loop.go` | Agent entry for `loop next\|apply\|status\|gate`; wires to `internal/loop` |
| **Loop core** | `internal/loop/next.go`, `apply.go`, `apply_writes.go`, `policy.go`, `gate.go`, `deliberation_packet.go`, `recommend.go` | FM-03 saturation/gate; FM-07 apply path; FM-10 spawn vs discovery |
| **Deliberation policy** | `internal/deliberation/select.go`, `types.go` | FM-03 SelectNext reason codes (`hop_budget_exceeded` vs `p19_saturated`); hop increment |
| **Deliberation persistence** | `internal/domain/deliberation.go`, `internal/store/deliberation.go` | FM-03 STOP stickiness; deliberation reset (S01 residual) |
| **P19 saturation signal** | `internal/loop/policy.go` (`p19SaturatedFromLastStep`, `BuildPolicyInputs`) | FM-03 why `p19_saturated` fires with hop_count=1 |
| **Task add (CLI)** | `cmd/trace/add.go` | FM-01/08/10 — `trace add task` vs `discovery` kinds |
| **Task create (domain)** | `internal/domain/create.go` (`CreateTask`, discovery/decision helpers) | FM-10 what creates tasks vs cognitive entities |
| **Loop apply writes** | `internal/loop/apply.go`, `apply_writes.go` | FM-10 `spawned_tasks[]` path vs standalone `trace add` |
| **MCP entry** | `cmd/trace-mcp/main.go` | Thin stdio server only |
| **MCP handlers** | `internal/mcp/server.go` (registration + descriptions), `tools_write.go` (`trace_add`), `tools_loop.go` (`trace_loop`), `tools_parity.go` (`trace_tasks`), `tools_context.go` (`trace_context`) | FM-08 tool discoverability; FM-02 context packet shape |
| **Install / enforcement** | `internal/install/enforcement.go`, `cursor.go`, `cursorhook.go`, `agents.go` | FM-04/05 — gate-only rules; no gap-pass / add-task nudge |
| **Seed export/import** | `cmd/trace/seed.go`, `internal/domain/seed_export.go`, `seed_import.go` | FM-01 seed UUID anchoring; FM-02 export honesty; export vs SQLite sync |
| **Gate evaluation** | `internal/loop/gate.go` (`EvaluateGate`, `evaluateEdit`) | FM-03 agent-visible block on `--for edit`; `task_not_found` when DB empty |
| **Planner** | `internal/planner/` (via `BuildPolicyInputs` / `GetPlan`) | FM-03 `PlanExists`, execute_pending gating |
| **Compiler / context** | `internal/compiler/`, `internal/retrieval/` (via `trace_context`, loop next packet) | FM-02 what agents see in context vs export |

**Note:** MCP tool implementations live under `internal/mcp/`, not `cmd/trace-mcp/` (entrypoint only).

## S01 review residuals (must address in audit)

These were forwarded from P24-S01-02; each needs a **mechanism row** in CODEBASE-AUDIT.md (may map to FM-03, FM-02, FM-10):

| Residual | Audit question | Start paths |
|----------|----------------|-------------|
| **SelectNext reason-code mapping** | Why E01 export shows `p19_saturated` at hop_count=1 while docs cite `hop_budget_exceeded`? Document first-match order in `select.go` L7–12 vs `policy.go` saturation inputs | `internal/deliberation/select.go`, `select_test.go`, `internal/loop/policy.go` L103–109 |
| **Export vs SQLite sync** | Why live `loop gate` returned `task_not_found` while `trace/graph.json` had deliberation? Trace import/export round-trip and when `.trace/trace.db` diverges from committed export | `cmd/trace/seed.go`, `internal/domain/seed_import.go`, `seed_export.go`; G1 `.trace/` vs `trace/graph.json` |
| **Deliberation reset after gap pass** | Is there any API/path to reset `hop_count` or transition out of STOP after gap fixes? If not, cite absence | `internal/domain/deliberation.go`, `ApplyTransition`, `loop apply` — grep `HopCount`, `Stopped`, `InitialState` |
| **FM→code paths** | Every FM-01..FM-10 row in POSTMORTEM §3 must link to ≥1 code mechanism (or label `protocol`/`experiment` with rationale) | Full inventory table above |

## Required reads (mandatory)

Read before writing audit rows. Record path + line in every claim.

### Phase 24 SoT + S01 handoff

| # | Path | Extract |
|---|------|---------|
| 1 | [POSTMORTEM.md](../scope-01-dogfood-postmortem/POSTMORTEM.md) | §3 FM matrix; §4 open questions; command block (reason codes, graph counts) |
| 2 | [FINDINGS.md](../../FINDINGS.md) | Two-mode model; FM per-session status — do not flatten |
| 3 | [INVESTIGATION.md](../../INVESTIGATION.md) | Investigation questions A–D; intervention categories |

### Loop / deliberation (FM-03 core)

| # | Path | Extract |
|---|------|---------|
| 4 | `internal/deliberation/select.go` | `SelectNext` first-match; hop vs p19 order |
| 5 | `internal/deliberation/types.go` | `HopBudget`, reason codes, `ApplyTransition` hop increment rules |
| 6 | `internal/loop/policy.go` | `BuildPolicyInputs`, `p19SaturatedFromLastStep` |
| 7 | `internal/loop/gate.go` | `EvaluateGate`, `evaluateEdit` — agent-visible violations |
| 8 | `internal/loop/next.go` | Next packet `deliberation` section — what agents see before edit |
| 9 | `internal/loop/apply.go` | Apply envelope, plan_changes, spawned_tasks, discovery links |

### Task creation / promotion (FM-01, FM-08, FM-10)

| # | Path | Extract |
|---|------|---------|
| 10 | `cmd/trace/add.go` | Kinds: task vs discovery; flags |
| 11 | `internal/mcp/tools_write.go` | `trace_add` kinds + error messages |
| 12 | `internal/mcp/server.go` | Tool descriptions for add, link, loop, tasks |
| 13 | `internal/mcp/tools_loop.go` | `trace_loop` next/apply/status parity with CLI |

### Seed / export (FM-01, FM-02)

| # | Path | Extract |
|---|------|---------|
| 14 | `cmd/trace/seed.go` | import allowed keys; export `--strict --enforce` |
| 15 | `internal/domain/seed_export.go` | `BuildSeedDocument`, omitted fields |
| 16 | `internal/domain/seed_import.go` | `ImportSeedDocument`, deliberation upsert |
| 17 | `experiments/.../runs/G1/seed/gt.json` | Fixed UUID anchoring (FM-01 evidence) |

### Harness (FM-04, FM-05)

| # | Path | Extract |
|---|------|---------|
| 18 | `internal/install/enforcement.go` | Installed rules text — gap pass? add task? |
| 19 | `experiments/.../runs/G1/.cursor/hooks/trace-loop-gate.sh` | Hook bypass when no TRACE_TASK_ID |
| 20 | [PROMPT-G1-ENFORCE.md](../../../../../experiments/ab-incident-tracker/prompts/PROMPT-G1-ENFORCE.md) L75–77 | Harness recommends `trace add` task — product support? |

## Investigation tasks (by theme)

Map findings to FM IDs in CODEBASE-AUDIT.md.

1. **Saturation / STOP / hop_budget (FM-03)** — Trace `SelectNext`: `hop_count >= HopBudget` → `hop_budget_exceeded` **before** `p19_saturated` check. Explain E01 export (`hop_count=1`, `p19_saturated`). Trace `p19SaturatedFromLastStep` (`NewPlanChanges == 0 && NewSpawnedTasks == 0`).
2. **Discovery → task gap (FM-10)** — After `trace add discovery`, is there an apply path to spawn tasks? Contrast `loop apply` `spawned_tasks[]` vs agents using `trace_link discovery-mentions-task` only.
3. **Task creation defaults (FM-01, FM-08)** — When does loop/apply create tasks vs expect manual `trace add task`? Does MCP `trace_add` description prioritize discovery over task?
4. **Seed anchoring (FM-01)** — Does seed import pin task UUIDs in UX (`trace_tasks`, context packet task list)? Import idempotency by ID?
5. **Export honesty (FM-02)** — Session A vs B decision counts; are uncertainties omitted from export by design? `--strict` gate behavior.
6. **MCP vs CLI parity (FM-08)** — Tool names, descriptions, discoverability for loop vs add vs discovery vs link.
7. **Install rules (FM-04, FM-05)** — Do installed rules mention gap pass, orchestrator ownership, or `trace add task` — or only gate?
8. **Deliberation reset (FM-03, S01 residual)** — Is `hop_count` ever cleared on new session / gap pass / transition DONE? Cite code or document absence.
9. **Orchestrator / mode (FM-04, FM-09)** — Product vs harness: what is **not** in Trace code (Multitask split) vs what install could add.
10. **Cross-arm / post-hoc (FM-06, FM-07)** — Label protocol/experiment levers where code is N/A; cite export timestamps / apply ordering if product-adjacent.

## Commands to run (paste summaries into CODEBASE-AUDIT evidence appendix)

Run from repo root unless noted.

```bash
# Build trace if needed
go build -o bin/trace ./cmd/trace

# Reason-code mapping — read select.go + run unit tests
go test ./internal/deliberation/... -run 'SelectNext|HopBudget' -count=1

# Gate semantics spot-check (G1 — may task_not_found if DB empty; document both outcomes)
export TRACE_BIN=$PWD/bin/trace
export TRACE_PROJECT_ROOT=$PWD/experiments/ab-incident-tracker/runs/G1
$TRACE_BIN loop status --task e0100000-0000-4000-8000-000000000050
$TRACE_BIN loop gate --task e0100000-0000-4000-8000-000000000050 --for edit

# MCP tool registration strings (FM-08)
grep -n 'Name:\|Description:' internal/mcp/server.go | head -40

# Saturation computation
grep -n 'p19Saturated\|P19Saturated\|MaxIterationsReached' internal/loop/*.go internal/deliberation/*.go

# Task spawn paths
grep -n 'spawned_tasks\|CreateTask\|discovery-mentions' internal/loop/*.go cmd/trace/add.go internal/mcp/tools_write.go

# Deliberation reset search
grep -rn 'HopCount\s*=' internal/domain/deliberation.go internal/loop/ internal/store/deliberation.go

# Seed import shape
grep -n 'seedImportAllowedKeys\|deliberation_states' cmd/trace/seed.go internal/domain/seed_*.go

# Optional: codegraph explore (if index present)
# codegraph_explore on SelectNext, EvaluateGate, toolAdd
```

## Deliverable

### `CODEBASE-AUDIT.md` — locked template

Create in this scope folder.

#### §1 Summary

2–4 sentences: top product mechanisms explaining Mode A thin graph + Mode B discovery-without-task-promotion + persistent STOP.

#### §2 FM mechanism table (required)

| FM-ID | mechanism | file:line | agent-visible symptom | change lever |
|-------|-----------|-----------|----------------------|--------------|
| FM-01 | e.g. seed import preserves task UUIDs | `internal/domain/seed_import.go:L…` | Agent sees same 5 tasks in `trace_tasks`; no prompt to expand | product / harness |
| FM-03 | e.g. p19 saturated when no plan changes on apply | `internal/loop/policy.go:L103–109` | `loop gate --for edit` blocked; STOP in export | product |
| … | … | … | … | … |

Rules:
- **≥8 rows** (FM-01..FM-10); prefer all 10.
- **file:line** must be specific (function or branch).
- **agent-visible symptom** = what Cursor agent sees (CLI output, MCP JSON, hook stderr, context packet field).
- **change lever** = one of product \| harness \| protocol \| experiment — no vague “improve UX”.

#### §3 S01 residual reconciliation

Dedicated subsection:

- SelectNext: `hop_budget_exceeded` vs `p19_saturated` with hop_count=1
- Export vs `.trace/` SQLite divergence (G1 evidence)
- Deliberation reset: exists or documented gap

#### §4 Cross-cutting observations

Bullets: apply vs add split; MCP description gaps; install enforcement scope; optional FM-06/07 protocol notes.

#### §5 Optional repro snippet

≤30 lines pseudocode or shell — comment block only; illustrates one FM mechanism.

#### §6 Open questions → S04

Bullets deferring intervention ranking (not fixing here).

### FINDINGS.md update

Set codebase audit row:

```markdown
| Codebase audit | S02 | **draft** — [CODEBASE-AUDIT.md](scopes/scope-02-codebase-loop-audit/CODEBASE-AUDIT.md) |
```

Add 3–5 bullet summary under a `## Codebase audit (S02 draft)` section — link to CODEBASE-AUDIT §2; do not duplicate full table.

## FM-ID → mechanism mapping guidance

Use POSTMORTEM §3 status as hypothesis; **confirm or refute** in code.

| FM | Name | Likely mechanisms (start grep/read here) |
|----|------|-------------------------------------------|
| FM-01 | Seed anchoring | `seed_import.go` task upsert by ID; `gt.json` fixed UUIDs; `trace_tasks` list; loop next `tasks[]` |
| FM-02 | Graph thin export | `seed_export.go` field omission; uncertainties export path; strict gate |
| FM-03 | Loop saturation | `select.go` STOP branches; `policy.go` p19; `gate.go` premature_implementation; hop sticky STOP |
| FM-04 | Orchestrator bypass | **Harness** — install rules lack parent-graph ownership; hook skips without TRACE_TASK_ID |
| FM-05 | Enforcement optional | `enforcement.go` text; `.trace/config.json` enforce off default; hook script L5–7 |
| FM-06 | Cross-arm leakage | **Protocol** — identical starter; not loop code; note absence of arm isolation in product |
| FM-07 | Post-hoc planning | Loop apply ordering; export timestamps; plan_critiqued set on apply_writes |
| FM-08 | Tool surface gap | MCP `trace_add` description lists discovery; `trace_link` for mentions-task; no “promote to task” tool |
| FM-09 | Mode-dependent effectiveness | Composite — cite 2+ mechanisms; human prompt not in product |
| FM-10 | Discovery without task promotion | `trace add discovery` + `trace_link discovery-mentions-task`; `spawned_tasks` only via loop apply |

## Forbidden

- Product Go changes (including “tiny fixes”)
- Rewriting POSTMORTEM.md or S01 board history
- Intervention ranking (defer to S04 INTERVENTION-MATRIX.md)
- Flattening Session A/B in FINDINGS

## Todo updates

Per agent-loop-protocol: edit **only** board row **P24-S02-01** — status + notes when done.

## Minimal todos

- [ ] Read **Required reads** tables (SoT + loop + task + seed + harness)
- [ ] Walk **Codebase audit map** inventory; confirm paths still exist
- [ ] Run **Commands to run**; capture outputs for §3 residuals
- [ ] Create `CODEBASE-AUDIT.md` §1–§6 per locked template
- [ ] Map **≥8 FMs** with file:line + change lever
- [ ] Reconcile S01 residuals in §3
- [ ] Update [FINDINGS.md](../../FINDINGS.md) codebase section
- [ ] Self-check exit criteria; mark board row done

## Exit criteria

- [ ] `CODEBASE-AUDIT.md` exists with §2 table (FM-ID \| mechanism \| file:line \| agent-visible symptom \| change lever)
- [ ] **≥8** FM-01..FM-10 rows with specific file:line citations
- [ ] S01 residuals addressed in §3 (reason codes, export/DB sync, deliberation reset)
- [ ] Change levers tagged product \| harness \| protocol \| experiment — no vague levers
- [ ] [FINDINGS.md](../../FINDINGS.md) codebase audit row → draft + summary bullets
- [ ] No product Go in diff
- [ ] Board row P24-S02-01: status=done, Notes only

## Next

**P24-S02-02**
