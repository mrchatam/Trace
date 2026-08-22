# P18-S05-00 — rebuild CLI + MCP binaries (FINAL)

## Metadata
- id: P18-S05-00
- todo_ids: [P18-S05-00]
- role: planner
- skills: [planning-and-task-breakdown, writing-for-agents]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** for an in-phase **rebuild** after VERIFY: fresh `bin/trace` + `bin/trace-mcp` so dogfood/MCP sessions cannot SKIP on a stale binary. **No product Go.** Lesson: stale `bin/trace-mcp` caused experiment SKIP. **This row does not rebuild** (planner-only). **S05-02** closes DR-HANDOFF. Stop if [../../00-PHASE-PLANNER.md](../../00-PHASE-PLANNER.md) is not FINAL (it is).

## References
- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [phase README](../../README.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- Live: `cmd/trace`, `cmd/trace-mcp`, `internal/mcp.RegisteredToolNames` (10 tools)
- Experiment convention: `experiments/runs/2026-08-17-multi-cap/README.md` (`CGO_ENABLED=1` trace, `CGO_ENABLED=0` trace-mcp)
- GOMODCACHE class: S01-02 REVIEW-NOTES (sandbox 403 `segmentio/encoding`); S04 VERIFY product bar
- [docs/TODO.md](../../../../TODO.md)

## Session start
Follow agent-loop-protocol. Unattended: no Plan-mode switch; **no rebuild this row**; **no product Go**. Depends-on: **P18-S04-02 APPROVE**. This is **not** a successor phase. **S02-00 FINAL:** DF-88 is docs/help/comments only — this scope still rebuilds binaries after VERIFY; do **not** reverse exclude. S02 must not retarget MCP CGO build-note lines (S05-01 may still correct them to CGO0).

## Live confirm (2026-08-18 — this planner; inspection only, no rebuild)

| Check | Result |
|-------|--------|
| `RegisteredToolNames()` | **10**, registration order: `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_tasks`, `trace_capability`, `trace_impact`, `trace_version` (`internal/mcp/server.go`) |
| Live `./bin/trace-mcp -h` | Same **10** names including **`trace_impact`** (help wraps across 3 lines; matches `cmd/trace-mcp/main.go`) |
| Live `./bin/trace version` | `0.0.0-dev` |
| Existing binaries | `bin/trace` + `bin/trace-mcp` mtime **2026-08-17 17:32** — catalog already 10; binaries still **stale** vs S01–S04 product. Rebuild is **required**, not optional because `-h` already lists 10 |
| README Build MCP line | still `CGO_ENABLED=1 go build -o bin/trace-mcp ./cmd/trace-mcp` |
| `cmd/trace/help.go` MCP line | still `CGO_ENABLED=1 go build -o bin/trace-mcp ./cmd/trace-mcp` |
| Help tests | **no** assertion of MCP CGO=1 — optional CGO0 one-liners will not trip a named help test |

## Depends-on S04 FINAL (VERIFY planner 2026-08-18 — do not re-lock product tests)

Sibling [../scope-04-phase-verify/00-PLANNER.md](../scope-04-phase-verify/00-PLANNER.md) **FINAL** delivers:

- Named DF-87/88/89 `-run` filters + carry-forward (incl. P17 seed-export keepers; two-clone **not** required)
- Evidence path `VERIFY-NOTES.md`; S04-01 **starts** DR-HANDOFF Notes (`no successor`); S04-02 re-verifies product only (**APPROVE** high)
- Stale `bin/trace` / `bin/trace-mcp` is **this scope’s job**, not a VERIFY fail
- Do **not** re-run the S04 named DF suite as this scope’s fail bar

S05-02 still **owns DR-HANDOFF close**. Do not reverse DF-88.

## Locked defaults (FINAL)

| Item | Value |
|------|-------|
| Cwd | repo root (`github.com/mrchatam/Trace`) |
| `bin/trace` | `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace` (tree-sitter) |
| `bin/trace-mcp` | `CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp` (no analyzers; CGO1 also links but is **not** the convention) |
| Prefetch / sandbox | **Preferred prefix on both builds:** `GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off` (S04 product-bar class). Bare `CGO_ENABLED=… go build` is OK **if it succeeds**. If module fetch **403**s (`segmentio/encoding` or any proxy block), retry with that prefix — **not a product defect**, not a FAIL. Prefix retry fail → **FAIL** (cannot produce binaries) |
| Help catalog | `./bin/trace-mcp -h` lists **all 10** names including **`trace_impact`**. SoT = `RegisteredToolNames()` order. Missing any name → **FAIL**. Must **not** advertise `trace_install`, `trace_decide`, `trace_plan`, `trace_index` |
| CLI identity | `./bin/trace version` exits 0 (prints `0.0.0-dev` today) |
| Optional DF-87 | See **DF-87 live context — skip vs fail** below. **Skip = non-fail.** **Run red = FAIL.** |
| Docs honesty | If README / `cmd/trace/help.go` still say `CGO_ENABLED=1` for `trace-mcp`, S05-01 **may** correct **those two lines only** to CGO0. No other docs drive-by |
| DR-HANDOFF | **S05-02 owns close** = default **`no successor`**. This planner + S05-01 **must not** close |
| Forbidden | Product feature Go; hosted MCP; reversing DF-88 exclude; rewriting P17/`P18-00` history; treating this as research S05; re-running S04 named suite as fail bar; closing DR-HANDOFF on 00 or 01 |

### Locked rebuild commands (FINAL)

`mkdir -p bin` if needed. Run from repo root.

```bash
# Preferred (sandbox 403 class — S01-02 / S04 VERIFY)
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp

# Bare OK only if it succeeds (same CGO split)
# CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
# CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp
```

### Locked help-tool assertion (FINAL — fail bar)

```bash
./bin/trace version
./bin/trace-mcp -h
```

`-h` stdout **must** contain each of: `trace_why`, `trace_context`, `trace_add`, `trace_link`, `trace_transition`, `trace_review`, `trace_tasks`, `trace_capability`, `trace_impact`, `trace_version`.

Live help wrapping (must stay these names):

```text
usage: trace-mcp [-C|--project <dir>]
  Thin MCP stdio server (official go-sdk). Tools: trace_why, trace_context,
  trace_add, trace_link, trace_transition, trace_review,
  trace_tasks, trace_capability, trace_impact, trace_version.
```

Record the 10 names + both binary mtimes in S05-01 board Notes.

### DF-87 live context — skip vs fail (FINAL)

S01/S04 already proved DF-87 with named tests. This is a **thin live CLI check** on the **fresh** `bin/trace`, not a re-prove of the S04 suite.

| Path | Fail bar |
|------|----------|
| **Skipped** | **Non-fail.** Notes **must** say `DF-87 live context: skipped` + one-line reason |
| **Run, green** | Extra evidence. Notes: `DF-87 live context: PASS` |
| **Run, red** | **FAIL this row.** Optional does **not** mean “ignore a red result” |

**Red (FAIL if run):** non-zero `context` exit; **or** stderr contains `syntax error near "/"`; **or** no packet on stdout.

**Green:** exit 0, stdout is a context packet (JSON default), stderr does **not** contain `syntax error near "/"`. Expand-only / empty FTS hits is **OK** — the defect was abort-on-slash MATCH.

**Recipe** (temp dir **outside** the repo; do **not** write repo `.trace/`):

```bash
TMP=$(mktemp -d)
./bin/trace -C "$TMP" init
./bin/trace -C "$TMP" add task --title 'GET /notes/search' --id 22222222-2222-2222-2222-222222222222
./bin/trace -C "$TMP" context 22222222-2222-2222-2222-222222222222 --format json
# then rm -rf "$TMP"
```

If `init`/`add` cannot run (no sandbox tmp, etc.), **skip** (non-fail) — do not invent a repo-`.trace/` mutation. Do **not** re-run S04 named DF tests as a substitute fail bar.

S05-02: if S05-01 skipped, reviewer skip is also **non-fail**. If S05-01 ran PASS, reviewer may re-run or accept Notes. If S05-01 ran red and marked done anyway → **reject APPROVE**.

## Planner work
1. [x] Confirm live `-h` vs `RegisteredToolNames` (10 including `trace_impact`)
2. [x] Lock rebuild cmds + GOMODCACHE 403 class + help assertion + DF-87 skip-vs-fail as **FINAL**
3. [x] Thicken 01/02/SCOPE-TODOS; stamp DR-HANDOFF close owner (still OPEN)
4. [x] Mark **FINAL**; next **P18-S05-01**

## Exit criteria
- [x] This prompt **FINAL** with rebuild cmds + 10-tool list + DF-87 skip-vs-fail
- [x] 01/02 thickened
- [x] Board Notes; next **P18-S05-01**
- [x] No product Go this row; **no rebuild** this row; DR-HANDOFF **not** closed

## Out of scope
- Rebuilding `bin/trace` / `bin/trace-mcp` (S05-01)
- Closing DR-HANDOFF / marking Phase 18 complete (S05-02)
- Product feature Go / new MCP tools / daemon / hosted MCP
- Re-running S04 named DF-87/88/89 suite as this scope’s fail bar
- Reversing DF-88 exclude; rewriting P17/`P18-00` history
- Auto-boarding Phase 19 / research S05 / `plan simulate` / D21+

## Next
**P18-S05-01** (implementer rebuild + `-h` catalog; optional DF-87 live check).
