# P23-S06-02 — Phase 23 verify review + DR-HANDOFF

## Metadata
- id: P23-S06-02
- todo_ids: [P23-S06-02]
- role: reviewer
- skills: [code-review-and-quality, documentation-and-adrs, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective

Independent **fresh-session** review of S06-01 evidence; re-run locked verify floor (do **not** trust Notes alone); **close DR-HANDOFF** with explicit successor decision (**never TBD**). Phase 23 complete when this row is `done`.

## References

- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) — Phase handoff (mandatory)
- [00-PLANNER.md](00-PLANNER.md) — FINAL evidence bar + DR-HANDOFF policy
- [01-verify.md](01-verify.md) — locked floor Blocks A–G + CLI smoke
- [ENFORCEMENT.md](../../ENFORCEMENT.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md)
- Pattern: [P21 S08-02 review](../../../phase-21-thoughtprocess-completion/scopes/scope-08-phase-verify/02-scope-review.md)

## Session start

Follow agent-loop-protocol. **Fresh session** — reviewer must not be the S06-01 implementer. Board edits: **status + notes only** on prior rows. Unattended: execute review loop until blocker/high clear or spawned forward.

## Locked DR-HANDOFF close policy (FINAL — S06-00)

| Field | Locked value |
|-------|--------------|
| Who gathers evidence | **S06-01** — floor + archive; DR-HANDOFF stays **OPEN** |
| Who closes | **S06-02 only** |
| Status on pass | `DR-HANDOFF.md` → **CLOSED** |
| Successor decision | **`no successor`** — default unless **human operator** has already named a forward phase before this row runs |
| Rationale | Phase 23 MVP delivers gate library, gate CLI, opt-in DONE/export enforce, status violations, config modes, harness install — without daemon or hosted MCP. Forward work (MCP gate tool, full harness parity, auto-strict CI) requires human promotion as Phase 22 → Phase 23 did. |
| Must not | Leave `Successor decision: TBD`, `later`, or empty; rewrite Phase 22 historical `no successor`; auto-scaffold Phase 24 without human promotion |
| Phase complete | **Yes** when this row `done` + DR-HANDOFF CLOSED + all Phase 23 board rows done |
| Next runnable after close | **none** (board queue empty until human promotes) |

### DR-HANDOFF scope checklist (tick on APPROVE)

From [DR-HANDOFF.md](../../DR-HANDOFF.md):

- [ ] S01: `domain.PrematureImplementation` + gate evaluator (PolicyInputs/SelectNext reuse)
- [ ] S02: `trace loop gate` CLI + `trace.loop.gate.v1`
- [ ] S03: `--enforce` on DONE transition + `seed export --strict`
- [ ] S04: loop status `violations[]` + `.trace/config.json` enforce modes
- [ ] S05: harness install rules (cursor/claude + optional cursor-hook)
- [ ] S06: VERIFY evidence + successor decision

### Residuals to list on close (non-blocking)

| Topic | Disposition |
|-------|-------------|
| MCP wrapper for `trace loop gate` | Optional; stdout CLI sufficient for hooks |
| Non-Cursor harness full hook parity | Claude rules shipped; other IDEs deferred |
| Auto-strict CI without explicit `--enforce` | Human promotion; config `strict` ≠ auto-block |
| Multi-violation top-level lift | S02 locked `len==1` only |
| Export-honesty beyond GateForExport | S03 residual |
| Multi-task export full-scan perf | S03 residual |
| Env override for enforce mode | S04 residual |

### DR-HANDOFF.md update template (on APPROVE)

```markdown
**Status:** **CLOSED**

| Field | Value |
|-------|-------|
| Closed | YYYY-MM-DD |
| Successor decision | **no successor** |
| Phase 23 outcome | Gate library + CLI; opt-in DONE/export enforce; status violations[]; config off/warn/strict; harness install (cursor/claude/cursor-hook); schema **027**; MCP **15** |
| Residuals (non-blocking) | MCP gate wrapper optional; non-Cursor harness adapters; auto-strict CI; multi-violation lift |
| Forward (human queue) | Phase 24+ only if human promotes |
```

If human has **named a forward phase** before S06-02: replace successor line with `Phase NN — <name>` and cite promotion evidence in board Notes.

If verify bar **fails**: keep DR-HANDOFF **OPEN**, spawn `P23-S06-02a` implement + `02b` review immediately below; successor stays **`no successor` intent** until pass.

## Review focus

Confirm independently:

1. **S01 library (Block A)** — 18 named gate tests PASS; no SelectNext fork outside `gate.go`.
2. **S02 CLI (Block B)** — 14 gate CLI tests; exit 0/1/2 contract; `trace.loop.gate.v1`.
3. **S03 enforce (Block C)** — two-layer DONE preserved; export enforce=no write.
4. **S04 status/config (Block D)** — parity test PASS; config fail-closed off; no auto-enforce.
5. **S05 install (Block E)** — rules + hook; git-hook unchanged; no config auto-write.
6. **P19/P20 keepers (Block F)** — loop next/apply/status green.
7. **Compat (Block G)** — ceiling **027**; no 028+.
8. **CLI evidence** — archived; ordinary CLI not test-only helpers.
9. **ENFORCEMENT.md bar** — all § Evidence bar claims mapped in S06-01 Notes.
10. **No scope creep** — no daemon, no hosted MCP, no new MCP tools beyond **15**.

## Evidence to collect

| Check | Evidence |
|-------|----------|
| Gate blocks edit | CLI #1 or `TestLoopGateBlockedExitOne` |
| Gate allows clean | `TestLoopGateAllowedExitZero` |
| Status ↔ gate parity | `TestLoopStatusViolationsMatchGateEdit`, CLI #2 |
| DONE `--enforce` block | `TestTransitionDoneEnforceBlocksVerificationDebt`, CLI #3 |
| Export enforce no write | `TestSeedExportStrictEnforceNoWriteOnViolation`, CLI #4 |
| Config default off | `TestTraceConfigEnforceDefaultOff`, `TestLoadEnforceModeMissingFile` |
| No auto-enforce | plain DONE/export with config strict unchanged |
| Install rules present | `TestInstallCursorIncludesLoopGateRule`, CLI #5 |
| git-hook unchanged | `TestInstallGitHookUnchanged`, CLI #6 |
| Further rules scope | grep status / DONE / export in rules fixtures |
| cursor-hook gate call | `TestInstallCursorHookCallsGate` |
| MCP count | **15** tools — no gate MCP added |
| Schema | **027** max |
| Named tests | All from 01-verify Blocks A–E present + passing |
| P19 keepers | Block F minimum tests PASS |

## Review checklist

### Blockers

- [ ] **Blocker:** SelectNext / policy fork outside `internal/loop/gate.go`
- [ ] **Blocker:** Config `strict` auto-enables transition `--enforce` or export `--strict --enforce`
- [ ] **Blocker:** `trace.loop.status.v1` schema string changed (non-additive)
- [ ] **Blocker:** Gate CLI exit **1** uses global `exitUsage` instead of gate block code
- [ ] **Blocker:** Export `--enforce` writes file on violation
- [ ] **Blocker:** git-hook fragment changed (`# begin-trace` → enforcement markers)
- [ ] **Blocker:** Install auto-writes `.trace/config.json` with strict
- [ ] **Blocker:** Missing named tests from 01-verify floor
- [ ] **Blocker:** P19 loop keeper tests regressed
- [ ] **Blocker:** Compat ceiling ≠ **027** or 028+ present
- [ ] **Blocker:** Daemon / hosted MCP / HTTP enforcement introduced

### High

- [ ] **High:** Status violations disagree with gate for same fixture
- [ ] **High:** Two-layer DONE broken (enforce bypasses or replaces domain review incorrectly)
- [ ] **High:** Hook fail-closed when `trace` missing (should fail-open per S05)
- [ ] **High:** cursor-hook calls wrong `--for` or wrong CLI
- [ ] **High:** S06-01 evidence missing CLI files #1–#6

### Medium / low

- [ ] **Medium:** Help missing gate/enforce/config/install docs
- [ ] **Medium:** AGENTS.md markers collide with git-hook markers
- [ ] **Low:** stderr prefix inconsistency across gate/status/transition
- [ ] **Nit:** Hook script mode ≠ 0755

## Locked re-verify commands (reviewer session — minimum)

```bash
# S01 library
go test ./internal/loop/... ./internal/domain/... -count=1 -run 'TestEvaluateGate|TestPrematureImplementation_Code'

# S02 gate CLI
go test ./cmd/trace -count=1 -run 'TestLoopGate|TestHelpIncludesLoopGate'

# S03 enforce
go test ./cmd/trace -count=1 -run 'TestTransitionDoneEnforce|TestSeedExportStrict'

# S04 status + config
go test ./internal/config/... -count=1
go test ./cmd/trace -count=1 -run 'TestLoopStatusViolations|TestTraceConfig'

# S05 install
go test ./internal/install/... -count=1 -run 'TestInstallGitHookUnchanged|Enforcement|CursorHook'

# P19 keepers (minimum)
go test ./cmd/trace -count=1 -run 'TestLoopNextPacketShape|TestLoopApplyMalformedInputFailsClosed|TestLoopApplyReplayAndStatusFlow|TestLoopStatusInsufficientHistory'

# Compat
CGO_ENABLED=1 go test ./evals/compat/... -count=1 -run TestCompatibilitySecurityChecklist
```

Broader re-run of full Blocks A–G from [01-verify.md](01-verify.md) encouraged; **minimum** above mandatory for APPROVE.

## Two-layer DONE walkthrough (verify at review)

1. Task with verification debt + review PASS + `--as-operator`:
   - Without `--enforce` → transition may succeed (domain escape hatch).
   - With `--enforce` → exit **1**, task stays non-DONE, stderr violation.
2. Config `{ "enforce": "strict" }` does **not** change either path without explicit `--enforce`.

## Default-off walkthrough (verify at review)

1. No `.trace/config.json` → status exit **0**, violations in JSON when blocked, **no** stderr hints.
2. `{"enforce":"warn"}` → status exit **0**, stderr hints on violation.
3. Install `--write` → **no** `.trace/config.json` created.

## Spawn policy

| Severity | Action |
|----------|--------|
| blocker / high | Small inline fix **or** spawn `P23-S06-02a` + `02b` immediately below |
| medium | Prefer spawn unless trivial one-liner |
| low / nit | List in REVIEW-NOTES; do not block close |

Do not rewrite S06-00/S06-01 `done` prompts.

## Evidence artifacts

- Read S06-01 `experiments/runs/…-p23-s06-01-verify/evidence/*` + optional `VERIFY-NOTES.md`
- Write **`REVIEW-NOTES.md`** in this scope folder (recommended)
- Update [DR-HANDOFF.md](../../DR-HANDOFF.md) on APPROVE
- Update board Notes: confidence, successor, residuals, schema **027**, MCP **15**
- Update [`docs/TODO.md`](../../../../TODO.md) orchestrator paste if phase complete (status line only)

## Exit criteria

- [ ] Independent re-verify floor PASS (minimum commands above + compat)
- [ ] S06-01 evidence reviewed (CLI files exist; Blocks A–G cited)
- [ ] ENFORCEMENT.md evidence bar met
- [ ] No open blocker/high without pending spawn
- [ ] **`DR-HANDOFF.md` CLOSED** with successor **`no successor`** (or human-named forward phase)
- [ ] DR-HANDOFF scope checklist all ticked
- [ ] Phase 23 all rows `done` in board Notes
- [ ] Confidence **high** (or **medium** with explicit residuals)
- [ ] Residuals listed (MCP wrapper, harness adapters, auto-strict CI)

## Forbidden

- Leaving successor **TBD** when row is `done`
- Closing DR-HANDOFF on S06-01 evidence without independent re-run
- Rewriting Phase 22 `done` history or Phase 22 DR-HANDOFF
- Implementing Phase 24 scaffold as part of review
- Adding migration 028+ or gate MCP tool during review fix

## Next

**none** after close (human promotes forward work)
