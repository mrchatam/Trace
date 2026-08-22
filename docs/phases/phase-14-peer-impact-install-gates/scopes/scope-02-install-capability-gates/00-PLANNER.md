# P14-S02-00 — Install / capability gates (FINAL)

## Metadata
- id: P14-S02-00
- todo_ids: [P14-S02-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Inventory live `trace install` + capability catalog / selection. Lock **FINAL** defaults for marker-gated detect→install→uninstall registry (rank **4**) + graduated allowlist / durable tool-decision audit (rank **5**). Thicken `01`/`02`/SCOPE-TODOS. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Laws 9, 17
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A1–A8; ranks 4–5
- [SIMILAR-PROJECTS-REVIEW-2026-08-16.md](../../../../research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — CBM/codegraph install matrix; agentrq graduated audit
- [AGENT_ENVIRONMENT.md](../../../../AGENT_ENVIRONMENT.md)
- Live: Phase 06 capability, Phase 09/11 `install cursor`, `evals/capability`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Depends-on: S01 APPROVE (P14-S01-02). Material phase locks A1–A8 hold; inventory before FINAL. No clarification needed — A4 + ranks 4–5 already cut YOLO / 43-client fleets / MCP-as-P0.

## Depends (S01 — light)
After P14-S01-02 APPROVE: preserve `trace impact` including **`walk`** + named ImpactWalk tests + Gate F prelim. Install matrix work must not delete or fork impact CLI/domain planted report paths. S01 residual (non-blocking): late `allowContainsOut` upgrade after contains-UP at a greater hop may skip re-expand — do not “fix” in S02; optional S03 VERIFY spot-check only.

## Live inventory (2026-08-17)

| Area | Present? | Evidence / gap |
|------|----------|----------------|
| `trace install cursor` print / `--write` | **Present** | `cmd/trace/install.go` — upsert `$HOME/.cursor/mcp.json` + `.bak.<UTC>`; `--bin` / `--mcp-json`; DF-22/50 reload tip stderr; stdout JSON-only on print |
| Install detect / list targets | **Missing** | Hard-coded `cursor` only; unknown target → usage exit |
| Install uninstall / reverse | **Missing** | No remove of `mcpServers.trace`; no idempotent reverse |
| Multi-client registry + tiers | **Missing** | No STABLE\|CONDITIONAL\|OPT_IN; no marker gating |
| Capability catalog + task require | **Present** | mig **010**; `internal/domain/capability.go` + store; CLI `declare\|list\|require\|unrequire\|missing`; MCP `trace_capability` (9 tools total) |
| Builtin MCP specs (no auto-seed) | **Present** | `BuiltinMCPCapabilitySpecs()` → nine `mcp:trace_*` AVAILABLE specs |
| Graduated allowlist / tool-decision audit | **Missing** | No decision table; no AUTO_ALLOWED\|PENDING\|ALLOWED\|DENIED; chat/logs not a store |
| YOLO / AllowAll default | **Absent** (good) | Must stay absent (Laws 9/17; research §D) |
| Capability ablation Gate | **Present** | `evals/capability` — keep green; do not rewrite TP/FP planting |
| Install Cursor named tests | **Present** | `TestInstallCursor*` (print/write/reload tip) — keep green |
| Impact walk (S01) | **Present** | Must stay green; **do not touch** retrieval ImpactWalk / `trace impact walk` |
| Schema ceiling | **012** | Next additive mig for audit = **013** |
| MCP install / decide tools | **Absent** (good) | **Do not add** this scope (A4 CLI/rules-first) |

## Locked defaults (FINAL) — do not re-debate in P14-S02-01

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Research | Rank **4** marker-gated install matrix + detect→install→uninstall; rank **5** graduated allowlist + durable tool-decision audit (≠ chat) |
| Install home (library) | **`internal/install`** — registry + Detect / Install / Uninstall. Logic **not** forked in CLI (G19). Move Cursor upsert helpers here (or package-private under install); thin `cmd/trace/install.go` adapter only |
| Capability audit home | **`internal/domain` + `internal/store`** — additive mig **`013_capability_tool_decisions.sql`** (name free if intent clear). CLI extends `trace capability` |
| Tier vocabulary | Exact Trace strings: **`STABLE` \| `CONDITIONAL` \| `OPT_IN`** (peer STABLE\|CONDITIONAL\|OPT_IN) |
| Registry shape | In-code registry (no mig for targets). Each target: `id` (slug), `tier`, `Detect(ctx/env) → Detected bool + Reason`, `Install(opts)`, `Uninstall(opts)`. List via `Registry()` / `ListTargets()` |
| Cursor | **`STABLE`**. Retain print + `--write` behavior (stdio `trace-mcp` + `-C ${workspaceFolder}`; backup; reload tip). **Add** `Detect` (e.g. `$HOME/.cursor` exists **or** mcp.json path exists — document exact rule in 01) and **`Uninstall`** that removes only `mcpServers.trace`, leaves other servers, backups like write, **idempotent** if already absent |
| CONDITIONAL bar | Ship **≥1** non-Cursor CONDITIONAL target that **refuses Install** when its marker/platform/existing-config proof is absent (clear validation error), and **succeeds** when marker present. Marker = file/dir proof under project or home (implementer picks one real peer-shaped marker, e.g. `CLAUDE.md` / `.claude/` / Codex config path — **not** 43 clients). OPT_IN targets (if any) install **only** when explicitly named on CLI (`trace install <id>`), never from blank `detect`-driven mass write |
| detect CLI | `trace install detect` (or `trace install --detect`) → JSON list of targets with `id`, `tier`, `detected`, `reason`. **No** silent writes |
| install CLI | `trace install <target> [--write] …` — Cursor flags retained; CONDITIONAL without marker → fail closed (non-zero); print-without-write still OK for STABLE Cursor snippet |
| uninstall CLI | `trace install uninstall <target>` — idempotent reverse; Cursor removes `mcpServers.trace` only |
| Mass / fleet install | **Forbidden** — no “install all detected”, no 43-client matrix, no plugin YOLO flips |
| Allowlist graduation | **Narrow auto-allow → else human.** Auto-allow **exact** slug match against `BuiltinMCPCapabilitySpecs()` slugs only (`mcp:trace_why` … `mcp:trace_version`). **No** globs, **no** prefix wildcards beyond that fixed builtin set, **no** AllowAll / YOLO default |
| Decision statuses | **`AUTO_ALLOWED` \| `PENDING` \| `ALLOWED` \| `DENIED`** (store TEXT; normalize uppercase) |
| Resolve API | Domain e.g. `ResolveToolDecision(ctx, slug) (Decision, error)` — if persisted ALLOWED/DENIED → that; else if exact builtin MCP slug → `AUTO_ALLOWED` (may upsert durable AUTO_ALLOWED row on first resolve — lock: **persist** so audit ≠ ephemeral); else → `PENDING` (not an error by itself) |
| Fail-closed gate | `AssertToolAllowed(ctx, slug) error` — returns validation/deny error when decision is `PENDING` or `DENIED`; OK for `AUTO_ALLOWED`/`ALLOWED`. Wire **at least** into a library call site used by tests; **do not** weaken transition `allow_missing_caps` / DONE operator gates; **do not** change ablation planting |
| Human decide CLI | Extend capability: `trace capability decide --slug <slug> --decision ALLOWED\|DENIED [--reason <text>]` + `trace capability decisions` (list JSON). Actor default `cli` (not authorization) |
| Durable ≠ chat | SQLite decision rows are SoT; stderr/chat/MCP logs are **not** the audit record |
| MCP | **No** new MCP tools this scope (`trace_install_*`, `trace_capability_decide`, etc.). Nine tools + `trace_version` retained |
| Impact / S01 | **Untouched** — no edits to ImpactWalk / `trace impact walk` / planted Gate F domain as part of “install work” |
| Ablation / Cursor tests | Keep `evals/capability` + `TestInstallCursor*` green |
| Forbidden | YOLO/AllowAll defaults; actor-name allowlists as auth; daemon/HTTP/embeddings/Neo4j; full-rebuild indexer; MCP dump / new install MCP; 43-client fleets; boarding S05 / `plan simulate`; reopening DF-60…67; board spawn by implementer; regressing S01 ImpactWalk |
| Carry-forward | S01 ImpactWalk named tests + Gate F; honesty A/B/C+G; Gates E/F prelim/H; capability ablation; compat; p0x; x0; Gate C `dry_run:false`; product `./cmd\|internal\|evals`; existing `TestInstallCursor*` |
| Named tests (min) | See table below |

### Named tests (intent locked)

| Test | Intent |
|------|--------|
| `TestInstallDetectListsCursorStable` | `detect` includes `cursor` with tier `STABLE` and a non-empty reason |
| `TestInstallCursorUninstallIdempotent` | After `--write`, uninstall removes `mcpServers.trace`; second uninstall OK; other servers untouched |
| `TestInstallConditionalRefusesWithoutMarker` | CONDITIONAL target Install without marker → error / non-zero; no write |
| `TestInstallConditionalWritesWithMarker` | Same target with marker present → Install succeeds (write or in-memory fixture equivalent) |
| `TestInstallCursorPrintSnippet` / Write* / ReloadTip | **Keep** green (existing) |
| `TestCapabilityDecisionAutoAllowBuiltinMCP` | Resolve on a builtin `mcp:trace_*` slug → `AUTO_ALLOWED` (+ durable row) |
| `TestCapabilityDecisionUnknownPendingFailClosed` | Unknown slug → `PENDING`; `AssertToolAllowed` fails |
| `TestCapabilityDecisionHumanAllowPersists` | `decide ALLOWED` → Resolve returns `ALLOWED` across reopen/store |
| `TestCapabilityDecisionDenyBlocks` | `decide DENIED` → `AssertToolAllowed` fails |
| Capability ablation suite | **Keep** green (`evals/capability`) |
| S01 ImpactWalk + Gate F named tests | **Keep** green (carry-forward; do not “fix” in S02) |

## Owns
| Item | Intent |
|------|--------|
| Marker-gated install matrix | Write only when marker/platform/existing config proves activity (CONDITIONAL); STABLE Cursor retained |
| detect / install / uninstall | Idempotent reverse; Cursor full lifecycle; ≥1 CONDITIONAL proof target |
| Graduated allowlist + audit | Auto-allow narrow builtins → else human; durable decisions ≠ chat-as-sole-record |

## Explicit deferrals (not S02)
- YOLO / AllowAll defaults
- 43-client fleets / MCP as P0 architecture / new MCP install or decide tools
- Impact walk changes (S01) — including residual `allowContainsOut` re-enqueue
- Supersession / `plan simulate` / ranks 7+ / D21+
- Scout/Verify/Auditor tier skill policy (research rank 12)
- Actor OAuth / identity allowlists

## Depends note (S03)
S03 VERIFY must re-run: S02 named install + decision tests; existing `TestInstallCursor*`; capability ablation; **and** S01 ImpactWalk named tests + Gate F prelim. Optional spot-check only for S01 `allowContainsOut` residual — do not treat as S02 failure. See light note on S03 SCOPE-TODOS.

## Planner work (this row)
1. [x] Inventory install targets, markers, capability allow paths, audit persistence
2. [x] Lock FINAL: registry shape, STABLE\|CONDITIONAL\|OPT_IN, detect/install/uninstall, allowlist graduation, audit rows ≠ chat, named tests
3. [x] Thicken `01-install-capability-gates.md` + `02-scope-review.md` + SCOPE-TODOS to FINAL
4. [x] Light Depends note for S03 VERIFY

## Exit
- [x] 00-PLANNER marked **FINAL**
- [x] Board Notes; next **P14-S02-01**
- [x] Product Go — **not** this row
