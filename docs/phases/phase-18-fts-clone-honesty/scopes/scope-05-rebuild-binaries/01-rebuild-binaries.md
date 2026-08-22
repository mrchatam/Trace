# P18 / S05 / 01 — rebuild CLI + MCP binaries

## Metadata
- id: P18-S05-01
- todo_ids: [P18-S05-01]
- role: implementer
- skills: [systematic-debugging]
- mcps: [Shell, Read, Write, Grep]
- verification: automated

## Objective
Rebuild workspace binaries per sibling **00-PLANNER FINAL**. Confirm `trace-mcp -h` lists current tools. Board **status + Notes only**. **No product Go** except optional README/`help.go` CGO0 one-liners for `trace-mcp`.

**Stop if sibling `00-PLANNER.md` is still DRAFT.** It is **FINAL** — rebuild.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — required **FINAL** (commands + catalog + DF-87 skip-vs-fail SoT)
- [SCOPE-TODOS.md](SCOPE-TODOS.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md) — **do not close**
- Experiment CGO split: `experiments/runs/2026-08-17-multi-cap/README.md`
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Unattended: no Plan-mode switch. Depends-on: **P18-S04-02 APPROVE** + **P18-S05-00 FINAL**. Do **not** start from S04-01. Do **not** reverse DF-88. Do **not** implement S01 FTS. Do **not** re-run S04 VERIFY as this row’s fail bar. Still rebuild even if VERIFY notes mention stale binaries and even if live `-h` already lists 10 (binaries dated 2026-08-17 17:32 are stale vs S01–S04). **No board spawn.** **Do not close DR-HANDOFF.**

## Locked defaults (FINAL — copy)

| Item | Value |
|------|-------|
| Cwd | repo root |
| CLI | `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace` |
| MCP | `CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp` |
| Sandbox modules | **Preferred:** prefix both with `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off`. Bare build OK if it succeeds. 403 on `segmentio/encoding` (or any proxy block) → retry with prefix; **not** a product defect. Prefix retry fail → **FAIL** |
| Catalog | `./bin/trace-mcp -h` includes all **10** names, including `trace_impact`. Missing any → **FAIL**. Must not advertise `trace_install` / `trace_decide` / `trace_plan` / `trace_index` |
| DF-87 live | **Optional to skip; not optional to ignore red.** Skip = **non-fail** (Notes required). Run red (`syntax error near "/"` / no packet / non-zero) = **FAIL** |
| Docs | Optional: README Build + `cmd/trace/help.go` MCP lines `CGO_ENABLED=1` → `CGO_ENABLED=0`. Those two lines only |
| Forbidden | Feature code; hosted MCP; skipping rebuild because tests already passed or `-h` already lists 10; closing DR-HANDOFF; mutating repo `.trace/` for the DF-87 check |

## Role work
1. `mkdir -p bin`. Rebuild **both** binaries (GOMODCACHE prefix preferred — see locked commands).
2. Run `./bin/trace version` and `./bin/trace-mcp -h`. Record the 10 tool names + `ls -l bin/trace bin/trace-mcp` mtimes in board Notes.
3. Optional DF-87 live check **on the fresh CLI** (recipe in 00-PLANNER). If skipped, Notes: `DF-87 live context: skipped` + reason. If run red → **FAIL**.
4. Optional: align README Build + `cmd/trace/help.go` MCP build line to `CGO_ENABLED=0`.
5. Board Notes → **P18-S05-02**. Do **not** write DR-HANDOFF close.

## Locked verify

```bash
mkdir -p bin
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp
./bin/trace version
./bin/trace-mcp -h
```

`-h` must list: `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_tasks`, `trace_capability`, `trace_impact`, `trace_version`.

Optional DF-87 (temp dir outside repo; skip = non-fail):

```bash
TMP=$(mktemp -d)
./bin/trace -C "$TMP" init
./bin/trace -C "$TMP" add task --title 'GET /notes/search' --id 22222222-2222-2222-2222-222222222222
./bin/trace -C "$TMP" context 22222222-2222-2222-2222-222222222222 --format json
rm -rf "$TMP"
```

Fail the live check only if it was **run** and stderr contains `syntax error near "/"` (or non-zero / no packet).

## Todo updates
Board **status + Notes only**. Notes must include: rebuild cmds used (GOMODCACHE y/n), both mtimes, 10-name catalog, DF-87 live PASS|skipped (reason).

## Exit criteria
- [ ] Both binaries rebuilt **this row** (new mtime vs 2026-08-17 17:32)
- [ ] `trace-mcp -h` lists 10 tools including `trace_impact`
- [ ] DF-87 live: PASS **or** skipped-with-Notes (red if run = not done)
- [ ] Board Notes; next **P18-S05-02**
- [ ] DR-HANDOFF still OPEN

## Minimal todos
- [ ] Rebuild CLI + MCP
- [ ] Confirm `-h` catalog + version
- [ ] Optional DF-87 live (or skip with Notes)
- [ ] Board sync (no DR-HANDOFF close)
