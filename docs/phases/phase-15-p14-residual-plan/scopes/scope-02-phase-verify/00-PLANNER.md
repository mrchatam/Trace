# P15-S02-00 — Phase VERIFY / P14 residual remediation (FINAL)

## Metadata
- id: P15-S02-00
- todo_ids: [P15-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 15 VERIFY evidence: **S01 MCP Assert named regressions** + **carry-forward gates** + product pkgs. Decide **DR-HANDOFF** = **`no successor`**. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- [phase README](../../README.md) — disposition matrix
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — FINAL
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- S01 REVIEW-NOTES: [../scope-01-mcp-assert-dispatch/REVIEW-NOTES.md](../scope-01-mcp-assert-dispatch/REVIEW-NOTES.md) — **APPROVE high**
- Pattern: Phase 14 S03 VERIFY
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Depends-on: S01 APPROVE (landed). Grill only if DR-HANDOFF should promote a successor — default remains **`no successor`**.

## Depends-on (S01 — landed)

| Scope | Board | Locks imported |
|-------|-------|----------------|
| S01 | **APPROVE high** (P15-S01-02) — [REVIEW-NOTES.md](../scope-01-mcp-assert-dispatch/REVIEW-NOTES.md) | MCP Assert on every tool with slug `mcp:<Name>` (incl. `trace_version` store+Assert); no new tools/mig; named `TestMCPAssertDeniedBlocksCallTool` + `TestMCPAssertBuiltinAutoAllowedSucceeds` + `TestToolNamesRegistered`; R2 defer / R3–R4 wontfix unchanged |

## Live residuals → DR-HANDOFF decision (2026-08-17)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gap scheduled in Phase 15 | R1 MCP Assert wire-up | Closed by S01 APPROVE high — VERIFY must **re-prove named tests** |
| Explicit residual OK into VERIFY | R2 `allowContainsOut` late-upgrade (**defer**); R3 graphify space-in-path (**wontfix**); R4 CGO0 analyzers FAIL (**wontfix**) | **Do not fail VERIFY** for these |
| Goals sequence #2–#4 | S05 supersession / `plan simulate` / D21+ | Stay off-board — **not** Phase 16 unless Notes + human promote |
| Parallel dogfood (not board-blocking) | `experiments/` | Stay in `experiments/` — **not** boarded |
| Product bar | `./cmd\|internal\|evals` with CGO1 | Prefer product pkgs over full-module `./...` when graphify space FAIL present |

**DR-HANDOFF = `no successor`.** No Notes or APPROVE residuals justify scaffolding Phase 16 / S05 / `plan simulate` / D21+. Reopen only with explicit human promotion + scaffold (same posture as Phase 10–14 historical closes).

## Planner work
1. [x] Import S01 APPROVE evidence into 01-verify checklist
2. [x] Lock verify commands + DR-HANDOFF = `no successor`
3. [x] Thicken 01-verify / 02-scope-review / SCOPE-TODOS
4. [x] Mark this prompt **FINAL**; next **P15-S02-01**

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Phase 15 P14 residual remediation closeout — S01 MCP Assert named regressions — **not** a new planted eval gate |
| S01 home | `assertMCPToolAllowed` → `AssertToolAllowed(ctx, "mcp:"+name)` at all nine tool entries incl. `toolVersion` openStore+Assert |
| Migration | **None** from VERIFY — mig 013 already landed; compat ceiling **13** (no 014+) |
| S01 named | `TestMCPAssertDeniedBlocksCallTool`; `TestMCPAssertBuiltinAutoAllowedSucceeds`; `TestToolNamesRegistered`; keep `TestBuiltinMCPCapabilitySpecs` |
| Carry-forward | honesty A/B/C+G; Gates E/F/H; ablation; compat **13**; p0x; x0; product `./cmd\|internal\|evals` |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1` — **product pkgs PASS**; R3 graphify space FAIL on full `./...` OK; R4 CGO0 analyzers FAIL OK |
| Residuals OK | R2 deferred; R3/R4 wontfix — **do not fail VERIFY** |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; spawn remediation if fail |
| MCP | Still **nine** tools + `trace_version`; no install/decide MCP; Assert **is** on MCP dispatch (R1 closed) |
| Optional strong evidence | Grep `assertMCPToolAllowed` call sites; Gate C artifact inspect; G19 — non-blocking |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | **`no successor`** — **S02-01 starts** Notes; **S02-02 owns completion**. Do **not** scaffold Phase 16 / S05 / plan simulate / D21+ without explicit promotion |
| Forbidden | Product features on VERIFY; claiming R2/R3/R4 fixed; inventing Phase 16; Mode-B Gate C rewrite; daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–14 `done` history; YOLO/AllowAll; new MCP install/decide tools |

### Locked verify command set (FINAL)

```bash
# --- S01 MCP Assert named ---
CGO_ENABLED=0 go test ./internal/mcp/... ./internal/domain/... ./internal/store/... -count=1 -run 'TestToolNamesRegistered|TestMCPAssert|TestBuiltinMCPCapabilitySpecs|TestCapabilityDecision'

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E / F / capability ablation
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Gate H + compat (compat covers mig 013 ceiling)
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# P0-X + X0
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1

# Product regression bar (prefer over full-module ./... when R3 graphify present)
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Optional (strong evidence, **not** substitutes for package PASS / not Mode-B Gate C):

```bash
# Grep: assertMCPToolAllowed at all nine tool entries (incl. toolVersion)
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# Spot-check only (non-blocking): R2 allowContainsOut still present in impact_walk.go
# Goals #2–#4 (S05 / plan simulate / D21+): stay off-board unless Notes explicitly promote
# Do NOT fail for R3 graphify space FAIL or R4 CGO0 analyzers FAIL
```

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + DR-HANDOFF locked (`no successor`)
- [x] SCOPE-TODOS + board Notes; next `P15-S02-01`
- [x] Product Go — **not** this row

## Out of scope
- Running VERIFY (S02-01)
- Product Go / new MCP tools / daemon / mig
- Scaffolding Phase 16 / S05 / plan simulate / D21+ without explicit promotion
- Fixing R2 / R3 / R4
- Closing parallel dogfood experiments
- Claiming Phase 14 historical handoff was wrong
