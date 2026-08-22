# P10 / S01 / 01 — Retrieval / why fidelity

## Metadata
- id: P10-S01-01
- todo_ids: [P10-S01-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-19, DF-23, DF-25, DF-27, DF-29** per sibling **00-PLANNER** FINAL locks (2026-08-16). Goal-scope DPC attach; unify `plan-change`/`plan_change` at why/Exact; Exact/Why for `capability`; clarify decision MD trust labels without elevating blobs; fail-closed IncludeWhy. Keep carry-forward gates green. **No new migration. No new MCP tools.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law 4 + Law 9
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-19/23/25/27/29
- [phase README](../../README.md)
- Live: `internal/retrieval/{expand,discovery_plan_change,exact,why,types}.go`; `internal/compiler/{compiler,packet}.go`; thin `cmd/trace/why.go` + `internal/mcp/tools_why.go` alias only
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate.

## Locked defaults (FINAL — P10-S01-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | `internal/retrieval` + `internal/compiler`; thin CLI/MCP type normalize only |
| Migration | **None** |
| DF-19 | Goal-scoped DPC + pair-completion; see 00-PLANNER algorithm (foreign-task filter + single-goal unattributed fallback) |
| DF-19 tests | Single-goal still surfaces DPC; **multi-goal** task A must not see DPC tied to goal B tasks; supersede global-attach assertions |
| DF-23 | Canonical `plan_change`; accept `plan-change` at Exact/Why/CLI why/MCP why; emit underscore |
| DF-25 | `lookupEntity` `capability` via `GetCapability`; **`plan_scope` residual — do not implement** |
| DF-27 | Keep `trust=untrusted_data`; reword decision (opt. assumption) MD title away from “not project policy”; never TrustSystem for bodies |
| DF-29 | IncludeWhy=true → Why error propagates from TaskContext/ExpandContext |
| Carry-forward | honesty A/B/C + Gate G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` intact; `./...` |
| Forbidden | New mig; daemon/HTTP/embeddings; new MCP tools; global DPC dump; elevating decision trust; `plan_scope` Exact; Mode-B pack rewrite; board spawn |

## Extension points (exact files)

| File | Work |
|------|------|
| `internal/retrieval/discovery_plan_change.go` | Scoped helper (replace/repurpose global list); pair-completion + goal foreign filter + single-goal unattributed rule |
| `internal/retrieval/expand.go` | Call scoped helper for task Expand (pass seed task id / hits as needed) |
| `internal/retrieval/exact.go` | Normalize type aliases; `case "capability"`; keep `plan_change` |
| `internal/retrieval/why.go` (or Exact entry) | Normalize entityType before lookup |
| `internal/retrieval/types.go` | Optional: document alias consts / helper |
| `internal/retrieval/*_test.go` | DF-19 single+multi goal; DF-23 alias; DF-25 capability Exact/Why |
| `internal/compiler/compiler.go` | DF-29 fail-closed IncludeWhy |
| `internal/compiler/packet.go` | DF-27 decision(/assumption) MD title copy |
| `internal/compiler/*_test.go` | DF-19 context; DF-27 string assert; DF-29 error path |
| `cmd/trace/why.go` and/or shared normalize | Accept `plan-change` argv (if not fully handled in retrieval) |
| `internal/mcp/tools_why.go` | Same alias (prefer retrieval normalize so one place) |

## Role work

1. TDD DF-19 multi-goal pollution (fail under current global attach) + keep single-goal/x0-shaped green path.
2. Implement scoped DPC; update/replace old GC-01 tests.
3. DF-23 normalize + tests (`why plan-change <id>` / Exact).
4. DF-25 `capability` Exact+Why (+ test).
5. DF-27 MD copy + test (`RenderMarkdown` / TaskContext markdown).
6. DF-29 IncludeWhy error propagation + test.
7. Run locked verify suite; board **status + Notes only** (cite test names).

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/retrieval/... ./internal/compiler/... ./internal/domain/... ./internal/store/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Spot-checks: multi-goal Why on task A omits foreign DPC; `Exact`/`Why` accept `plan-change` and `capability`; IncludeWhy with forced Why failure returns err; decision markdown lacks “not project policy”.

## Exit criteria
- [ ] DF-19: no global all-project DPC on every task; goal-scope + pair-completion per locks; multi-goal regression green
- [ ] DF-23: `plan-change` alias works on why/Exact/MCP why; emitted type `plan_change`
- [ ] DF-25: Exact/Why support `capability`; `plan_scope` still residual (noted in Notes)
- [ ] DF-27: decision MD clarifies Law 9 vs Law 4; `trust` remains `untrusted_data`; no system elevation
- [ ] DF-29: IncludeWhy surfaces Why errors (fail-closed)
- [ ] No new mig; no new MCP tools
- [ ] Carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P10-S01-02**

## Out of scope
- S02 MCP tasks/capability tools / install reload / JSON case (DF-21/22/32)
- S03 index GC (DF-20)
- S04 operator/capability transition gates (DF-17/18/24/26/31)
- `plan_scope` Exact/Why; DF-28 handoff SoT; Mode-B Gate C pack rewrite

## Todo updates
Implementer: **status + notes only**. Record test names + DF checklist evidence. No spawning; no rewriting upcoming prompts.

## Minimal todos
- [ ] DF-19 scoped DPC + tests (single + multi-goal)
- [ ] DF-23 alias normalize + tests
- [ ] DF-25 capability Exact/Why + test; note plan_scope residual
- [ ] DF-27 decision MD label + test
- [ ] DF-29 IncludeWhy fail-closed + test
- [ ] Locked verify suite; board Notes
