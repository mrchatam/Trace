# P09-S03-02 — Scope review notes (2026-08-16)

Independent review of DF-03/DF-05 install-wire vs `00-PLANNER.md` / `01-install-wire.md` locks + `P09-S03-01` Notes. Fresh session; claims re-verified in-repo (no implementer session shared).

## Verdict

**APPROVE** — no blocker / high findings. Confidence: **high**.

## Evidence checklist

| Criterion | Result |
|-----------|--------|
| Command is `trace install cursor` (not competing `mcp snippet`) | Pass — `root.go` `case "install"`; `cmdInstall` requires `cursor`; unknown target → usage |
| Default print: merge-ready `{"mcpServers":{"trace":…}}` with `args: ["-C","${workspaceFolder}"]` | Pass — `install.go` + `TestInstallCursorPrintSnippet` |
| `--write` upserts only `trace`; other `mcpServers` preserved | Pass — `TestInstallCursorWriteMergeBackup` keeps `other` |
| Backup before overwrite; UTC `.bak.<stamp>`; path on stderr | Pass — `*.bak.YYYYMMDDTHHMMSSZ`; stderr asserted |
| Invalid JSON fail-closed (exit 2 / `exitFail`); no write / no backup | Pass — `TestInstallCursorWriteInvalidJSON` |
| Default `command` = `trace-mcp`; `--bin` / `--mcp-json` honored | Pass — print + write tests; create-missing dirs OK |
| Help + README + `experiments/ab-simple/PROTOCOL.md` document install + DF-05 | Pass — help line; README Install section; PROTOCOL DF-05 run-folder footgun |
| Thin `cmd/trace` only (G19); **no** new MCP tools / daemon / HTTP / mig | Pass — six tools still (`trace_why`/`context`/`add`/`link`/`transition`/`review`); schema through `010_*` only |
| S02 `trace tasks` still green; S01 Why/context-with-review still green | Pass — `TestTasksListAfterSeed` / `TestWhyAndContextWithLinkedReview` in `./...` |
| Carry-forward: honesty A/B/C + Gate G, p0x, x0, `./...` | Pass (independent re-run) |

## Verify (independent re-run)

```text
CGO_ENABLED=0 go test ./evals/honesty/... -count=1                                                    PASS
CGO_ENABLED=1 go test ./cmd/trace/... ./evals/p0x/... ./evals/x0/... ./evals/honesty/... -count=1   PASS
CGO_ENABLED=1 go test ./... -count=1                                                                  PASS
```

## Findings

None at blocker / high / medium.

### Low (no spawn)

- If existing `mcpServers` is valid JSON but **not** an object (e.g. array), the type assert fails and a new map replaces it — sibling servers cannot be preserved in that degenerate shape. Normal Cursor `mcp.json` objects are covered by merge tests.

### Nit

- Empty existing file (size 0) skips backup and writes fresh snippet — acceptable for first-time create; matches “backup before overwrite of existing **non-empty**” behavior in code.

## Spawns

None.

## Residuals

- MCP list-tasks / `trace_tasks` remains out of scope (discoverability stays CLI `trace tasks`).
- S01 residual unchanged: `plan_scope` ExactLookup still out; scope-only review expand path untested.
- Degenerate non-object `mcpServers` (low above).

## S04 compatibility

S04 VERIFY stubs remain compatible: spot-check `trace install cursor` snippet shape + DF-05 docs; do **not** require a new MCP list-tasks tool. Light note added on S04 `01-verify.md`.

## Next

**P09-S04-00** (phase VERIFY planner).
