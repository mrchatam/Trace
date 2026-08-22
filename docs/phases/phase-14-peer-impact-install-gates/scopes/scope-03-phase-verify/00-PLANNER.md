# P14-S03-00 — Phase VERIFY / peer-impact-install-gates (FINAL)

## Metadata
- id: P14-S03-00
- todo_ids: [P14-S03-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock Phase 14 VERIFY evidence: **S01 ImpactWalk named regressions** + **S02 install/capability-gate named regressions** + **carry-forward gates** + product pkgs. Decide **DR-HANDOFF** = **`no successor`**. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A8; VERIFY default `no successor`
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: Phase 12/13 VERIFY
- Sibling REVIEW-NOTES: [S01](../scope-01-impact-walks/REVIEW-NOTES.md), [S02](../scope-02-install-capability-gates/REVIEW-NOTES.md)
- Goals sequence: [TRACE-GOALS-PROGRESS-2026-08-17.md](../../../../research/TRACE-GOALS-PROGRESS-2026-08-17.md) — #2 S05 / #3 plan simulate stay off-board
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Depends-on: S01 + S02 APPROVE. Material locks A1–A8 hold — grill only if they conflict.

## Depends-on (S01–S02 — landed)

| Scope | Board | Locks imported |
|-------|-------|----------------|
| S01 | **APPROVED** high (P14-S01-02) | `retrieval.ImpactWalk` multi-seed BFS; seed exclusion; depth 1..2; incoming imports + P13 resolve; contains asymmetry; hop_risk; cap 64 loud totals; `trace impact walk`; no mig/MCP/`internal/impact`. Named: `TestImpactWalkMultiSeedExcludeSeeds`, `TestImpactWalkContainsAsymmetryNoSiblings`, `TestImpactWalkIncomingImportHop`, `TestImpactWalkLoudTruncation`, `TestImpactWalkHopRiskIncreases` + Gate F `TestPlantedImpactConflictsGateFPrelim`. Residual (non-blocking): `allowContainsOut` late-upgrade |
| S02 | **APPROVED** high (P14-S02-02) | `internal/install` registry STABLE\|CONDITIONAL\|OPT_IN; Cursor STABLE Detect/Install/Uninstall; CONDITIONAL `claude` marker-gated; mig **013** decisions AUTO_ALLOWED\|PENDING\|ALLOWED\|DENIED; Assert fail-closed; CLI detect/install/uninstall + capability decide/decisions; no new MCP; ImpactWalk untouched. Named: `TestInstallDetectListsCursorStable`, `TestInstallCursorUninstallIdempotent`, `TestInstallConditionalRefusesWithoutMarker`, `TestInstallConditionalWritesWithMarker`, `TestCapabilityDecision*` + `TestInstallCursor*` keepers + ablation. Residuals: Assert ≠ MCP dispatch (by design); S01 `allowContainsOut` |

## Live residuals → DR-HANDOFF decision (2026-08-17)

| Bucket | Items | Phase implication |
|--------|-------|-------------------|
| Product gaps scheduled in Phase 14 | S01 impact walks (rank 6) + S02 install/capability gates (ranks 4–5) | Closed by S01–S02 APPROVE high — VERIFY must **re-prove named tests** |
| Explicit residual OK into VERIFY | `AssertToolAllowed` library/CLI only (not MCP dispatch); optional S01 `allowContainsOut` spot-check; MCP sandbox `GOMODCACHE`/`GOPROXY=off` env note | Forward notes only — **not** a successor phase |
| Goals sequence #2–#4 | S05 supersession / `plan simulate` / D21+ | Stay off-board — **not** Phase 15 unless Notes + human promote |
| Parallel dogfood (not board-blocking) | `experiments/` ladders | Stay in `experiments/` — **not** boarded |
| Known `./...` nit | `similar projects/graphify` space-in-path FAIL; CGO0 analyzers FAIL (tree-sitter) | Pre-existing non-product / expected — VERIFY records **product pkgs PASS** |

**DR-HANDOFF = `no successor`.** No Notes or APPROVE residuals justify scaffolding Phase 15 / S05 / `plan simulate` / D21+. Reopen only with explicit human promotion + scaffold (same posture as Phase 10–13 historical closes / Phase 14 forward reopen).

## Planner work
1. [x] Import S01/S02 named tests + carry-forward command set
2. [x] Lock DR-HANDOFF default **`no successor`** (S05 / plan simulate / D21+ stay off-board unless Notes + human)
3. [x] Thicken `01-verify.md` + `02-scope-review.md` + SCOPE-TODOS to FINAL
4. [x] Stamp DR-HANDOFF ownership (S03-01 starts / S03-02 completes)

## Locked defaults (FINAL — this row)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Phase gate | Phase 14 peer-impact-install-gates closeout — S01+S02 named regressions — **not** a new planted eval gate |
| S01 home | Multi-seed ImpactWalk + contains asymmetry + `trace impact walk`; Gate F planted harness kept |
| S02 home | Install registry + CONDITIONAL marker gates + mig **013** tool decisions + Assert fail-closed; CLI-first (A4) |
| Migration | **None** from VERIFY — mig 013 already landed in S02; compat ceiling **13** (no 014+) |
| S01 named | `TestImpactWalkMultiSeedExcludeSeeds`; `TestImpactWalkContainsAsymmetryNoSiblings`; `TestImpactWalkIncomingImportHop`; `TestImpactWalkLoudTruncation`; `TestImpactWalkHopRiskIncreases`; keep `TestPlantedImpactConflictsGateFPrelim` |
| S02 named | `TestInstallDetectListsCursorStable`; `TestInstallCursorUninstallIdempotent`; `TestInstallConditionalRefusesWithoutMarker`; `TestInstallConditionalWritesWithMarker`; `TestCapabilityDecisionAutoAllowBuiltinMCP`; `TestCapabilityDecisionUnknownPendingFailClosed`; `TestCapabilityDecisionHumanAllowPersists`; `TestCapabilityDecisionDenyBlocks`; keep `TestInstallCursor*` + `TestPlantedCapabilitySelectionAblation` |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat **13**; p0x; x0; Gate C `dry_run:false` N=3 |
| Dry-run ≠ | Gate C / Gate F / Gate G / ablation / Gate H / checklist |
| Full bar | `CGO_ENABLED=1 go test ./... -count=1` — **product pkgs PASS**; known FAIL only `similar projects/graphify` space (non-product); CGO0 analyzers FAIL OK residual |
| Allowed Go on VERIFY | **None** for features — re-run + evidence docs only; spawn remediation if fail |
| MCP | Still **nine** tools + `trace_version`; no install/decide MCP; Assert ≠ MCP dispatch (honesty Note) |
| Optional strong evidence | `TestImpactWalkCLI`; Gate C artifact inspect; G19; optional `allowContainsOut` code spot-check — non-blocking |
| Spawn | On fail: `01a` implement / `01b` review (+`01c` re-VERIFY if needed) immediately below |
| DR-HANDOFF | **`no successor`** — **S03-01 starts** Notes; **S03-02 owns completion**. Do **not** scaffold Phase 15 / S05 / plan simulate / D21+ without explicit promotion |
| Forbidden | Claiming “every MCP call gated”; inventing Phase 15; Mode-B Gate C rewrite; daemon/HTTP/embeddings; full-rebuild indexer; rewriting Phase 00–13 `done` history; claiming Phase 13 historical handoff was wrong; YOLO/AllowAll; new MCP install/decide tools |

### Locked verify command set (FINAL)

```bash
# --- S01 ImpactWalk ---
CGO_ENABLED=0 go test ./internal/retrieval/... -count=1 -run 'TestImpactWalk'
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim

# --- S02 install registry + capability decisions ---
CGO_ENABLED=0 go test ./internal/install/... -count=1 -run 'TestInstallDetectListsCursorStable|TestInstallCursorUninstallIdempotent|TestInstallConditional'
CGO_ENABLED=0 go test ./internal/domain/... -count=1 -run 'TestCapabilityDecision'
CGO_ENABLED=1 go test ./cmd/trace/... -count=1 -run 'TestInstallCursor|TestInstallUsage|TestImpactWalkCLI'

# Capability ablation keep
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Honesty: Paths A/B/C + Gate G (CGO-free)
CGO_ENABLED=0 go test ./evals/honesty/... -count=1
CGO_ENABLED=0 go test ./evals/honesty/... -count=1 -run 'TestHonestyFailClosedPlantedClaim|TestHonestyEscapeRateGateGPrelim'

# Gate E / F / capability ablation carry-forward (F also under S01)
CGO_ENABLED=0 go test ./evals/replan/... -count=1 -run TestPlantedDiscoveryReplan
CGO_ENABLED=0 go test ./evals/impact/... -count=1 -run TestPlantedImpactConflictsGateFPrelim
CGO_ENABLED=0 go test ./evals/capability/... -count=1 -run TestPlantedCapabilitySelectionAblation

# Gate H + compat (compat covers mig 013 ceiling)
CGO_ENABLED=1 go test ./evals/perf/... -count=1 -run TestPlantedPerfLadderGateH
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist

# P0-X + X0
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... -count=1

# Supporting surfaces (optional strong evidence)
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./internal/planner/... ./internal/compiler/... ./internal/retrieval/... ./internal/install/... -count=1
# If mcp proxy 403: GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off …
CGO_ENABLED=0 go test ./internal/mcp/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... -count=1

# Full regression bar (product pkgs; graphify space FAIL is known residual)
# Prefer: GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off if sandbox module proxy 403s
CGO_ENABLED=1 go test ./... -count=1
```

Optional (strong evidence, **not** substitutes for package PASS / not Mode-B Gate C):

```bash
# Gate C artifact inspect (jq/grep OK): dry_run:false, N=3 — do not re-score
# G19: library packages do not import cmd/trace or cmd/trace-mcp
# Honesty Notes: AssertToolAllowed not on MCP dispatch (by design A4)
# Optional non-blocking: spot-check allowContainsOut late-upgrade residual in impact_walk.go
# Goals #2–#4 (S05 / plan simulate / D21+): stay off-board unless Notes explicitly promote
```

## Exit criteria
- [x] `01-verify.md` + `02-scope-review.md` runnable (thickened)
- [x] VERIFY commands + DR-HANDOFF locked (`no successor`)
- [x] SCOPE-TODOS + board Notes; next `P14-S03-01`
- [x] Product Go — **not** this row

## Out of scope
- Running VERIFY (S03-01)
- Product Go / new MCP tools / daemon / mig
- Scaffolding Phase 15 / S05 / plan simulate / D21+ without explicit promotion
- Closing parallel dogfood experiments
- Claiming Phase 13 historical handoff was wrong
- Wiring Assert into MCP dispatch
