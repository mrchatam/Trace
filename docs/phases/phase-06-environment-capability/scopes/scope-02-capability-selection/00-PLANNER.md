# P06 / S02 / 00-PLANNER — Capability selection

## Metadata
- id: P06-S02-00
- todo_ids: [P06-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Read, Glob, Grep, Write]
- verification: automated

## Objective
Finalize sibling `01-capability-selection.md` for **task capability selection + ablation harness** consuming S01 surface. Lock eval path/name and exit criteria. No product code.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [docs/init/A_PROJECT_PLAN.md](../../../../init/A_PROJECT_PLAN.md) Phase 6 — capability-selection ablation
- [docs/EVALUATION.md](../../../../EVALUATION.md) H7
- Prior planted patterns: `evals/impact` Gate F; `evals/honesty` Gate G; `evals/replan` Gate E
- S01 APPROVE: [../scope-01-capability-surface/REVIEW-NOTES.md](../scope-01-capability-surface/REVIEW-NOTES.md)

## Session start
Agent → clarify if needed → Plan → execute (planner).

## Depends
S01 capability surface APPROVE (`P06-S01-02`).

### Expected S01 hooks (live-confirmed 2026-08-16 after P06-S01-02 APPROVE)
| Hook | Use in ablation plant |
|------|------------------------|
| `UpsertCapability` / `ListCapabilities` | Seed AVAILABLE vs UNAVAILABLE/UNKNOWN capabilities |
| `RequireCapability` / `ListRequiredCapabilities` | Plant required set for a task |
| `MissingCapabilities(taskID)` | Score positive probes when required absent / non-AVAILABLE |
| Compiler packet `required_capabilities` + `missing_capabilities` | Assert packet lists only required / warns missing (no catalog dump) |
| `BuiltinMCPCapabilitySpecs()` | Optional plant of six `mcp:trace_*` MCP-kind caps |
| Mig **`010_capability_surface`** | Tables `capabilities` + `task_capability_requirements`; no S02 schema fork |
| CLI `trace capability` | Optional for humans; **harness must call library APIs** (G19) |

**Live confirm:** hooks match shipped `internal/domain` + `internal/compiler` + mig 010; no hook rename/divergence from S01-00 stubs.

## Phase ablation lock (P06-S02-00 FINAL — 2026-08-16)
| Item | Value |
|------|-------|
| Package | **`evals/capability`** (new; do not overload honesty/replan/impact/x0/p0x) |
| Named test | **`TestPlantedCapabilitySelectionAblation`** |
| Schema | **`evals/capability/schema-capability.json`** v1 (`schema_version` const 1) |
| Metrics artifact | Temp **`metrics-capability.json`** (schema-validated in test) |
| Evidence shape | 4 planted probes → TP=3 / FN=0 / FP=0 / TN=1 → **precision=1.0**, **recall=1.0** |
| Probes | Pos-1 UNAVAILABLE missing; Pos-2 UNKNOWN missing; Pos-3 selection filter (no catalog dump); Neg-1 clean AVAILABLE |
| Product Go | **None outside `evals/capability`** (+ schema). Selection = plant `RequireCapability` + assert packet/missing; S01 attach already filters to required-only |
| Non-claim | Not commercial multi-model theater; not Gate C/F/G/E substitute |
| Carry-forward | Honesty A/B/C; Gate G; Gate E; Gate F; p0x; x0; Gate C `dry_run:false` intact |
| Re-prove (S03) | `CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation` |

## Planner work
- [x] Lock selection APIs + missing-capability warnings vs live S01
- [x] Lock ablation harness path/names under `evals/capability`
- [x] Thicken `01-capability-selection.md`
- [x] Thicken `02-scope-review.md`
- [x] Light S03 Depends note (VERIFY re-prove command); sync board

## Locked defaults (FINAL — P06-S02-00)
| Item | Value |
|------|-------|
| Validation | Capability-selection ablation — **`TestPlantedCapabilitySelectionAblation`** + P/R tallies |
| S01 consume | Capability surface APIs only (no schema fork) |
| Carry-forward | Gate F/G/E + honesty + p0x/x0 + Gate C intact |

## Exit criteria
- [x] `01-capability-selection.md` runnable alone
- [x] Ablation harness path locked
- [x] No product Go

## Out of scope
- Phase VERIFY (S03); commercial multi-model capability theater; ontology expansion beyond S01
