# P11-S06-00 — MCP / install reload UX (FINAL)

## Metadata
- id: P11-S06-00
- todo_ids: [P11-S06-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Finalize S06 implement/review prompts for **DF-22, DF-37, DF-50**. Confirm live inventory; lock APIs/tests. **No product Go in this row.**

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [../../../../rules/project-rules.md](../../../../rules/project-rules.md)
- [../../../../rules/skills-map.md](../../../../rules/skills-map.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — G19
- [phase README](../../README.md)
- [00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md)
- [experiments/DOGFOOD-FINDINGS.md](../../../../../experiments/DOGFOOD-FINDINGS.md) — DF-22, DF-37, DF-50
- [experiments/_post_p10/BUGHUNT.md](../../../../../experiments/_post_p10/BUGHUNT.md) — DF-22 mitigated-when-reloaded; DF-50
- [experiments/_post_p10/bughunt/adv_install2/RESULTS.txt](../../../../../experiments/_post_p10/bughunt/adv_install2/RESULTS.txt) — print stderr empty vs write tip
- [experiments/POST-P10-MCP.md](../../../../../experiments/POST-P10-MCP.md) — DF-37 catalog 6→9 mid-session
- Phase 10 S02 FINAL: [../../../phase-10-integrity-surfaces/scopes/scope-02-mcp-parity-install/00-PLANNER.md](../../../phase-10-integrity-surfaces/scopes/scope-02-mcp-parity-install/00-PLANNER.md) — DF-22 tip on `--write` + `trace_version` + README
- Live: `cmd/trace/{install,help,install_test}.go`; README Install section; `internal/mcp` nine tools incl. `trace_version`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (planner). Material locks below — no grill (ops+docs tip parity; cannot force Cursor restart — A1–A7 do not conflict).

## Live confirmation (2026-08-16)

| Surface | Finding |
|---------|---------|
| `cmdInstallCursor` print path | Pretty JSON on **stdout**; **stderr empty** (DF-50; adv_install2) |
| `cmdInstallCursor --write` | Upsert + optional `.bak.<UTC>`; stderr tip present (P10 DF-22) |
| Tip string (write today) | `install: tip: rebuild trace-mcp, prefer absolute --bin, then reload/restart Cursor MCP (or reload window) so the stdio process is not stale` |
| `help.go` | Tip prose only under **“After --write:”** — print path not mentioned |
| README Install | Tip after build / `install --write`; documents `trace_version` |
| MCP binary | **Nine** tools incl. `trace_version` (P10 DF-21/22) — binary OK; live Cursor catalog may lag (DF-22/37) |
| Force restart | **Impossible** from Trace — Cursor owns stdio MCP lifetime |
| Migration | **None** — install is files+stderr only |

## Owns (FINAL)

| DF | Intent | Lock summary |
|----|--------|--------------|
| DF-50 | Print-only omits reload tip | Print success path emits **same** stderr tip as `--write`; stdout stays JSON-only |
| DF-22 | Stale MCP after install/rebuild (residual ops) | Tip parity on **all** successful `install cursor` exits (print + write); keep README + `trace_version`; help covers print and write |
| DF-37 | Live Cursor may advertise pre-P10 tool set until reload | Ops+docs only — tip + light reload/`trace_version` guidance; **no** process kill / daemon / catalog API |

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| DF home | **DF-22, DF-37, DF-50** only (P10 DF-21 nine tools + DF-32 snake_case remain as shipped; do not reopen) |
| Packages | Thin **`cmd/trace`** only (`install.go`, `help.go`, `install_test.go`). Optional one-line README + light `experiments/ab-simple/PROTOCOL.md` tip if needed for assert. **No** `internal/mcp` / domain / store changes required |
| Migration | **None** |
| Tip text | Single shared helper or const — **identical** stderr line for print and `--write` (keep current wording unless a typo fix; must still contain `reload` + `trace-mcp`) |
| DF-50 print | After successful stdout JSON encode: `fmt.Fprintf(os.Stderr, …tip…)`; exit 0. Failed encode: no tip (error path unchanged) |
| DF-50 write | Keep existing tip after successful upsert (may call same helper). Backup line stays separate |
| DF-50 stdout | Print mode must remain valid pretty JSON on stdout (tip **never** on stdout) |
| DF-22 help/README | Help: tip applies after **print or `--write`** (not write-only). README: same — tip after `install cursor` (print or write) / rebuild; keep prefer abs `--bin` + `trace_version` |
| DF-22 / DF-37 non-goals | **Do not** kill Cursor/MCP PIDs; **do not** add HTTP/daemon health; **do not** auto-reload Cursor; **do not** add new MCP tools; **do not** change mcp.json merge/backup semantics |
| DF-37 docs | Optional light PROTOCOL/README note: live MCP catalog may lag until reload/window refresh; confirm with `trace_version` / tool list. Product fix is tip visibility, not process control |
| Tests (required) | (1) **`TestInstallCursorPrintReloadTip`** (or extend `TestInstallCursorPrintSnippet`): print path stderr contains `reload` + `trace-mcp`; stdout still parses as `mcpServers.trace`. (2) Keep **`TestInstallCursorWriteMergeBackup`** / **`TestInstallCursorWriteCreateMissing`** tip asserts. (3) Optional help-string smoke if easy. (4) Carry-forward suites |
| Carry-forward | honesty A/B/C+G; Gate E/F; ablation; Gate H; compat; p0x; x0; Gate C `dry_run:false` untouched; P10 DF-21/22/32 binary behavior; P11-S01…S05 |
| Forbidden | New mig; daemon/HTTP/embeddings; full-rebuild indexer; forcing Cursor restart; new MCP tools; rewriting Phase 00–10 / P11-S01–S05 `done` history; S07+ product work |

## Effects on later scopes
- **S07** (seed-plan-review-polish): no install/reload coupling — serial after S06 review only. Light Depends note on S07 stubs.
- **S08 VERIFY:** include DF-50 print+write tip parity + help/README tip coverage + nine-tool/`trace_version` still present in evidence table.

## Exit
- [x] Thicken `01-mcp-install-reload.md` + `02-scope-review.md` + SCOPE-TODOS
- [x] Light upcoming Depends note (S07)
- [x] Board Notes; next **P11-S06-01**
- [x] Product Go — **not** this row
