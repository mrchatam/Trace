# P23-S06-00 — Phase 23 verify planner

## Metadata
- id: P23-S06-00
- todo_ids: [P23-S06-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- verification: automated

## Objective
Lock VERIFY bar for enforcement choke points + harness install. Own DR-HANDOFF close policy. **No product Go this row.**

## References
- [ENFORCEMENT.md](../../ENFORCEMENT.md) — product SoT + evidence bar (§ Evidence bar)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — **OPEN** until S06-02
- S01–S05 deliverables: scopes `scope-01` … `scope-05` (all implemented; S05-02 review pending)
- P19 loop keepers: `cmd/trace/loop_test.go` (`TestLoopNextPacketShape`, …)
- P20 deliberation floor: inherited via loop/deliberation/domain keepers (no P20 verify re-run required — P23 delta only)
- Compat ceiling: `evals/compat/compat_test.go` — **027** (`027_harness_agents.sql`; no 028+)
- Evidence pattern: `experiments/runs/YYYY-MM-DD-p23-s06-01-verify/evidence/`

## Live inventory (confirmed post S01–S05, 2026-08-20)

| Surface | Location | Landed (S01–S05) | S06 verify proof |
|---------|----------|------------------|------------------|
| Gate evaluator | `internal/loop/gate.go` | S01 — `EvaluateGate`, GateFor table | 17 named unit tests |
| Domain gate error | `internal/domain/gate_errors.go` | S01 — `PrematureImplementation.Code()` | `TestPrematureImplementation_Code` |
| Gate CLI | `cmd/trace/loop.go` `cmdLoopGate` | S02 — `trace.loop.gate.v1`, exit 0/1/2 | 14 named CLI tests |
| DONE enforce | `cmd/trace/transition.go` | S03 — `--enforce` before `TransitionTask` | 8 named transition tests |
| Export strict | `cmd/trace/seed.go` | S03 — `--strict`/`--enforce`/`--task` | 9 named export tests |
| Status violations | `internal/loop/apply.go` | S04 — additive `violations[]` | parity + schema tests |
| Config enforce | `internal/config/enforce.go` | S04 — off\|warn\|strict fail-closed off | 4 loader + 7 CLI tests |
| Harness install | `internal/install/enforcement.go`, `cursorhook.go` | S05 — cursor/claude/cursor-hook + AGENTS.md | 16 named install tests |
| git-hook | `internal/install/githook.go` | **Unchanged** P22 post-commit | `TestInstallGitHookUnchanged` |
| MCP catalog | `internal/mcp/` | **15** tools (P22) — no gate MCP | compat grep / no new registrations |
| Schema max | `internal/store/schema/` | **027** — no S01–S04 SQL | `TestCompatibilitySecurityChecklist` |
| `.trace/config.json` | gitignored | file-only; default off | config loader + status tests |
| Daemon / hosted MCP | — | **None** | review checklist |

### P23 delta → verify proof (FINAL)

| Scope | Delta | Verify proof |
|-------|-------|--------------|
| S01 | Shared gate library (SelectNext reuse) | 17 `TestEvaluateGate_*` + `TestPrematureImplementation_Code` |
| S02 | `trace loop gate` CLI + envelope | 14 `TestLoopGate*` + help |
| S03 | `--enforce` DONE + export strict | 17 `TestTransitionDoneEnforce*` + `TestSeedExportStrict*` |
| S04 | Status `violations[]` + config modes | 7 status + 7 config + 4 loader tests |
| S05 | Harness rules + cursor-hook | 16 enforcement/hook tests + git-hook unchanged |
| S06 | Phase verify + DR-HANDOFF | Blocks A–G + CLI smoke + compat **027** |

## FINAL verify floor (S06-01)

S06-01 must run and report PASS/FAIL for **every** block in [01-verify.md](01-verify.md). `-count=1` recommended.

| Block | Scope | Anchor |
|-------|-------|--------|
| **A** | S01 gate library | `go test ./internal/loop/... ./internal/domain/... -run 'Gate\|Premature'` |
| **B** | S02 gate CLI | `go test ./cmd/trace -run 'TestLoopGate\|TestHelpIncludesLoopGate'` |
| **C** | S03 enforce DONE/export | `go test ./cmd/trace -run 'TestTransitionDoneEnforce\|TestSeedExportStrict'` |
| **D** | S04 status + config | `go test ./internal/config/...`; loop status + `TestTraceConfig*` |
| **E** | S05 harness install | `go test ./internal/install/... -run 'TestInstall\|Enforcement\|CursorHook\|GitHook'` |
| **F** | P19/P20 loop keepers | `go test ./cmd/trace -run 'TestLoopNext\|TestLoopApply\|TestLoopStatus'` |
| **G** | Compat ceiling | `CGO_ENABLED=1 go test ./evals/compat/... -run TestCompatibilitySecurityChecklist` |

Preflight: all named tests exist (grep at S06-01 start). If absent, row **`failed`** — do not invent replacements.

## ENFORCEMENT.md evidence bar (must map in S06-01 Notes)

| ENFORCEMENT § | Claim | S06 evidence anchor |
|---------------|-------|---------------------|
| Pre-edit gate | Exit 1 on blocking uncertainty before `--for edit` | `TestLoopGateBlockedExitOne`, CLI smoke #1 |
| Status parity | `violations[]` matches gate for same task | `TestLoopStatusViolationsMatchGateEdit`, CLI smoke #2 |
| DONE enforce | `--enforce` fails with verification debt | `TestTransitionDoneEnforceBlocksVerificationDebt`, CLI smoke #3 |
| Export strict | `--strict --enforce` no write on violation | `TestSeedExportStrictEnforceNoWriteOnViolation`, CLI smoke #4 |
| Default off | Config missing → off; no auto-enforce | `TestTraceConfigEnforceDefaultOff`, `TestLoadEnforceModeMissingFile` |
| Harness install | Rules present after `trace install cursor --write` | `TestInstallCursorIncludesLoopGateRule`, CLI smoke #5 |
| P19/P20 keepers | Loop next/apply/status green | Block F |
| No daemon/MCP gate | stdout CLI only | review checklist + no new MCP tools |

## CLI smoke evidence (S06-01 archive)

Minimum artifacts under `experiments/runs/YYYY-MM-DD-p23-s06-01-verify/evidence/`:

| # | Artifact | Proves |
|---|----------|--------|
| 1 | `01-gate-blocked-edit.json` | Gate exit **1** + `trace.loop.gate.v1` + `recommended_phase` on blocking uncertainty |
| 2 | `02-status-gate-parity.json` | Status `violations[]` equals gate `--for edit` for same task |
| 3 | `03-transition-done-enforce-block.txt` | `--to DONE --enforce` exit **1** on verification debt |
| 4 | `04-export-strict-enforce-no-write.txt` | `--strict --enforce -o` exit **1**; file absent |
| 5 | `05-install-cursor-rules.txt` | `trace install cursor --write` → rules file + AGENTS.md markers |
| 6 | `06-git-hook-unchanged.txt` | post-commit fragment still `# begin-trace` (not `-enforcement`) |
| 7 | `99-run-metadata.txt` | git SHA, go version, Blocks A–G PASS/FAIL summary |

Optional: `VERIFY-NOTES.md` in this scope folder (recommended).

## DR-HANDOFF close policy (S06-02 only — locked here)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S06-01** — runs floor, archives evidence; **does not** close DR-HANDOFF |
| Who closes | **S06-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Successor decision | **`no successor`** — default unless **human operator** has already named a forward phase before S06-02 runs |
| Rationale | Phase 23 MVP (gate CLI, opt-in enforce, status violations, config, harness install) complete. Residuals (MCP gate wrapper, non-Cursor harnesses, auto-strict CI) are forward human promotion — not this close's job unless pre-named. |
| Must not | Leave `Successor decision: TBD`, `later`, or empty after S06-02 `done`; rewrite Phase 22 historical `no successor`; auto-scaffold Phase 24 without human promotion |
| Phase complete | **Yes** when S06-02 `done` + DR-HANDOFF CLOSED + all Phase 23 board rows done |
| DR-HANDOFF scope checklist | S06-02 ticks all boxes in [DR-HANDOFF.md](../../DR-HANDOFF.md) scope checklist |

### Residuals to list on close (non-blocking)

| Topic | Disposition |
|-------|-------------|
| MCP wrapper for `trace loop gate` | Optional forward; stdout CLI sufficient for hooks |
| Non-Cursor harness adapters | Claude rules shipped; full hook parity deferred |
| Auto-strict CI without explicit `--enforce` | Human promotion; config `strict` ≠ auto-block (S04 lock) |
| Multi-violation top-level lift | S02 locked `len==1` only |
| Export-honesty beyond GateForExport | S03 residual; forward if needed |
| Multi-task export full-scan perf | S03 residual |

If verify bar **fails** at S06-02: keep DR-HANDOFF **OPEN**, spawn fix row; successor stays **`no successor` intent** until pass.

## Locked defaults

| Item | Value |
|------|-------|
| Schema max | **027** (`027_harness_agents.sql`; re-lock at verify via compat) |
| Compat ceiling | **27** — forbid 028+ |
| MCP tools | **15** (P22 catalog unchanged) |
| Loop schemas | `trace.loop.gate.v1` (new S02); `trace.loop.status.v1` additive `violations[]` only |
| Successor | **`no successor`** default at S06-02 — never TBD after close |
| Product bar | [ENFORCEMENT.md](../../ENFORCEMENT.md) completion checklist |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p23-s06-01-verify/evidence/` |
| Out-of-scope S06-00 | Product Go; DR-HANDOFF close; Phase 24 scaffold |

## Named tests (S06-01 keeper summary)

| Scope | Count | Anchor |
|-------|------:|--------|
| S01 library | 18 | 17 `TestEvaluateGate_*` + `TestPrematureImplementation_Code` |
| S02 CLI gate | 14 | `TestLoopGateAllowedExitZero` … `TestLoopGateExecuteAllowedWhenPending` |
| S03 enforce | 17 | 8 transition + 9 export strict |
| S04 status/config | 18 | 7 status + 7 CLI config + 4 loader |
| S05 install | 16+ | enforcement + cursor-hook + `TestInstallGitHookUnchanged` |
| P19/P20 keepers | ~10+ | `TestLoopNext*`, `TestLoopApply*`, `TestLoopStatus*` (non-gate subset) |
| Compat | 1 | `TestCompatibilitySecurityChecklist` |

Full `-run` patterns in [01-verify.md](01-verify.md).

## Touch files

- [01-verify.md](01-verify.md) — thickened command floor + evidence spec
- [02-scope-review.md](02-scope-review.md) — thickened review + DR-HANDOFF template
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — S06-00 done
- `experiments/runs/` — S06-01 evidence (not this row)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — S06-02 closes

## Planner work

1. [x] Live inventory: S01–S05 surfaces landed; schema **027**; MCP **15**.
2. [x] Lock Blocks A–G verify floor in `01-verify.md`.
3. [x] Lock ENFORCEMENT.md evidence bar → test/CLI mapping.
4. [x] Lock CLI smoke artifacts (#1–#7).
5. [x] Lock DR-HANDOFF close: S06-01 open, S06-02 closes; **`no successor`** default; never TBD.
6. [x] Thicken `02-scope-review.md` with evidence table, re-verify cmds, DR-HANDOFF template.
7. [x] Update `SCOPE-TODOS.md`.
8. [x] No product Go.

## Exit criteria

- [x] 01-verify + 02-review thickened enough for S06-01/02 alone
- [x] Verify floor + DR-HANDOFF policy locked
- [x] No product Go

## Next

**P23-S06-01** after this row is `done`.
