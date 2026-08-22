# P11 / S06 / 01 — MCP / install reload UX

## Metadata
- id: P11-S06-01
- todo_ids: [P11-S06-01]
- role: implementer
- skills: [incremental-implementation, test-driven-development, systematic-debugging]
- mcps: [Shell, Read, Write, Grep, Glob]
- verification: automated

## Objective
Implement **DF-22, DF-37, DF-50** per sibling **00-PLANNER** FINAL locks (2026-08-16). Print + `--write` emit identical stderr reload tip; help/README cover both; no Cursor process control. **No new migration. No new MCP tools. Gate C `dry_run:false` untouched.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) — locks FINAL
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-22, DF-37, DF-50
- [experiments/_post_p10/bughunt/adv_install2/RESULTS.txt](../../../../../experiments/_post_p10/bughunt/adv_install2/RESULTS.txt)
- [experiments/POST-P10-MCP.md](../../../../../experiments/POST-P10-MCP.md)
- Phase 10 S02: [../../../phase-10-integrity-surfaces/scopes/scope-02-mcp-parity-install/00-PLANNER.md](../../../phase-10-integrity-surfaces/scopes/scope-02-mcp-parity-install/00-PLANNER.md)
- [phase README](../../README.md)
- Live: `cmd/trace/{install,help,install_test}.go`; README Install; `internal/mcp` nine tools / `trace_version` (inherit)
- Prior: P11-S05 no install/reload coupling
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol: Agent → clarify → Plan → execute. Locks are FINAL — do not re-debate. If 00-PLANNER is still DRAFT, stop and return to planner.

## Locked defaults (FINAL — P11-S06-00)

| Item | Value |
|------|-------|
| Module | `github.com/mrchatam/Trace` |
| Packages | Thin **`cmd/trace`** (`install.go`, `help.go`, tests); optional README + light PROTOCOL |
| Migration | **None** |
| DF-50 | Print success → same stderr tip as `--write`; stdout JSON-only |
| DF-22 | Tip on all successful install paths; help/README print+write; keep `trace_version` |
| DF-37 | Tip + docs only; no PID kill / daemon / new tools |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false`; P10 nine tools; P11-S01…S05 |
| Forbidden | Mig; force Cursor restart; new MCP tools; daemon/HTTP; board spawn; rewrite `done` history |

## Extension points (exact files)

| File | Work |
|------|------|
| `cmd/trace/install.go` | Shared tip helper; call from print success **and** `--write` success |
| `cmd/trace/install_test.go` | Assert print stderr tip; keep write tip tests |
| `cmd/trace/help.go` | Tip prose: after print or `--write` |
| `README.md` (Install / Cursor MCP) | Tip applies to print and write; keep abs `--bin` + `trace_version` |
| `experiments/ab-simple/PROTOCOL.md` (optional) | One-line reload/`trace_version` if assert needs it |

## Role work

1. TDD: print-only `install cursor` → stderr contains `reload` + `trace-mcp`; stdout still pretty `mcpServers.trace` JSON.
2. Extract shared tip emitter; wire print + write success paths (identical text).
3. Update `help.go` + README so tip is not write-only.
4. Confirm write tests still green; nine MCP tools / `trace_version` untouched.
5. Run locked verify suite; board **status + Notes only** (cite test names + DF-22/37/50).

## Algorithm sketch (non-normative — locks above win)

```text
printInstallCursorReloadTip(stderr):
  Fprint tip: rebuild trace-mcp, prefer abs --bin, reload/restart Cursor MCP…

install cursor (no --write):
  Encode snippet → stdout
  tip → stderr
  exit 0

install cursor --write:
  upsert mcp.json (+ backup line if any)
  tip → stderr   # same helper
  exit 0
```

## Verify commands (locked)

```bash
CGO_ENABLED=0 go test ./internal/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/compat/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

Prefer named asserts: `TestInstallCursorPrintReloadTip` (or extended `TestInstallCursorPrintSnippet`), `TestInstallCursorWriteMergeBackup`, `TestInstallCursorWriteCreateMissing` (or equiv).

## Exit criteria
- [ ] DF-22, DF-37, DF-50 green per FINAL locks
- [ ] Carry-forward suite green; Gate C `dry_run:false` untouched
- [ ] Board Notes ready for **P11-S06-02**

## Out of scope
- Other Phase 11 scopes; daemon/HTTP/embeddings; forcing Cursor restart; rewriting `done` history
