# P14 / S02 / 01 — Install / capability gates (FINAL)

## Metadata
- id: P14-S02-01
- todo_ids: [P14-S02-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement marker-gated install registry + graduated allowlist / durable audit per sibling **00-PLANNER FINAL**. Keep Cursor install + capability ablation + S01 ImpactWalk green. **Stop if 00-PLANNER is still DRAFT.**

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL** (required)
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) — A4
- Research ranks 4–5 — CBM/codegraph install; agentrq audit (adapt, reject AllowAll)
- Live: `cmd/trace/install.go`, `internal/domain/capability.go`, mig `010`, `evals/capability`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Do not re-debate FINAL locks.

## Locked defaults (FINAL — do not renegotiate)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Install library | **`internal/install`** — registry + Detect/Install/Uninstall; thin CLI only |
| Tiers | **`STABLE` \| `CONDITIONAL` \| `OPT_IN`** |
| Cursor | **STABLE** — keep print/`--write`/backup/reload tip; add Detect + Uninstall (`mcpServers.trace` only, idempotent) |
| CONDITIONAL | ≥1 non-Cursor target; refuse Install without marker; succeed with marker |
| detect | `trace install detect` → JSON targets (`id`,`tier`,`detected`,`reason`); no writes |
| uninstall | `trace install uninstall <target>` |
| Audit mig | **`013_*_capability_tool_decisions.sql`** (additive; ceiling today = 012) |
| Decision statuses | **`AUTO_ALLOWED` \| `PENDING` \| `ALLOWED` \| `DENIED`** |
| Auto-allow | Exact match **only** `BuiltinMCPCapabilitySpecs()` slugs; persist AUTO_ALLOWED on resolve |
| APIs | `ResolveToolDecision` + `AssertToolAllowed` (fail-closed on PENDING/DENIED) |
| Human CLI | `trace capability decide --slug … --decision ALLOWED\|DENIED` + `decisions` list |
| MCP / ImpactWalk | **No** new MCP tools; **do not** edit S01 ImpactWalk / `trace impact walk` |
| Forbidden | YOLO/AllowAll; mass install-all; 43-client fleets; daemon/HTTP/embeddings; board spawn |

## Extension points / files likely touched

| Layer | Path | Change |
|-------|------|--------|
| Install lib | `internal/install/` (new) | Registry, Cursor target, ≥1 CONDITIONAL target, Detect/Install/Uninstall |
| Install tests | `internal/install/*_test.go` | Named install tests (temp dirs; no real $HOME writes required) |
| CLI install | `cmd/trace/install.go` (+ tests) | detect / uninstall / dispatch; keep Cursor flags; help |
| Store | `internal/store/schema/013_*.sql` + store CRUD | Decision table + Get/Upsert/List |
| Domain | `internal/domain/capability*.go` (+ tests) | Resolve/Assert/Decide; auto-allow builtins |
| CLI capability | `cmd/trace/capability.go` (+ tests) | `decide` + `decisions`; help |
| Retrieval / impact / MCP | Prefer **zero** | Do not touch ImpactWalk or add MCP tools |
| Ablation | `evals/capability` | Expect **unchanged** green |

### CLI sketch (informative)

```text
trace install detect
# JSON: [{id,tier,detected,reason}, …]

trace install cursor [--write] [--bin path] [--mcp-json path]   # retained
trace install uninstall cursor [--mcp-json path]

trace install <conditional-id> [--write] …   # fails closed without marker

trace capability decide --slug mcp:trace_why --decision ALLOWED --reason "…"
trace capability decisions
# JSON: [{id,slug,decision,reason,actor,created_at,updated_at}, …]
```

## Role work
1. TDD named tests first (install detect/uninstall/conditional → then decision auto/pending/allow/deny).
2. Extract Cursor write path into `internal/install`; add Detect + Uninstall; register ≥1 CONDITIONAL proof target.
3. Add mig 013 + domain Resolve/Assert/Decide; extend capability CLI.
4. Prove existing `TestInstallCursor*` + ablation + **S01 ImpactWalk + Gate F** still green (do not edit walk code).
5. Run locked verify; board **status + Notes only**.

## Minimal todos
- [ ] `internal/install` registry + Cursor STABLE Detect/Install/Uninstall
- [ ] ≥1 CONDITIONAL target with marker gate (refuse / write tests)
- [ ] `trace install detect` + `uninstall` CLI (+ help)
- [ ] Mig 013 + store decision CRUD
- [ ] `ResolveToolDecision` / `AssertToolAllowed` / human decide
- [ ] Auto-allow builtin MCP slugs only; persist AUTO_ALLOWED
- [ ] `trace capability decide` + `decisions`
- [ ] Named tests green; keep `TestInstallCursor*` + ablation
- [ ] Carry-forward verify (incl. S01 ImpactWalk + Gate F) — no walk edits
- [ ] Board row P14-S02-01 Notes; next **P14-S02-02**

## Named tests (must ship)

| Test | Intent |
|------|--------|
| `TestInstallDetectListsCursorStable` | detect lists cursor STABLE + reason |
| `TestInstallCursorUninstallIdempotent` | uninstall removes trace entry; idempotent; others kept |
| `TestInstallConditionalRefusesWithoutMarker` | no marker → no write / error |
| `TestInstallConditionalWritesWithMarker` | marker → install OK |
| `TestInstallCursorPrintSnippet` / Write* / ReloadTip | **Keep** |
| `TestCapabilityDecisionAutoAllowBuiltinMCP` | builtin slug → AUTO_ALLOWED durable |
| `TestCapabilityDecisionUnknownPendingFailClosed` | unknown → PENDING; Assert fails |
| `TestCapabilityDecisionHumanAllowPersists` | ALLOWED survives reopen |
| `TestCapabilityDecisionDenyBlocks` | DENIED → Assert fails |
| Ablation + S01 ImpactWalk + Gate F | **Keep** |

## Verify commands

```bash
# Install + capability domain/store (CGO0 OK)
CGO_ENABLED=0 go test ./internal/install/... ./internal/domain/... ./internal/store/... ./evals/capability/... ./evals/honesty/... -count=1

# CLI + S01 impact carry-forward + gates (CGO1)
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/retrieval/... ./evals/impact/... ./evals/capability/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1

# Product packages
CGO_ENABLED=1 go test ./cmd/... ./internal/... ./evals/... -count=1
```

Also confirm Gate C artifacts under `docs/verification/gate-c-x0/` remain `dry_run:false` (do not rewrite Mode-B packs).

## Exit
- [ ] Named tests + verify green; board Notes; next **P14-S02-02**
- [ ] No product MCP/daemon; no ImpactWalk edits; no YOLO
