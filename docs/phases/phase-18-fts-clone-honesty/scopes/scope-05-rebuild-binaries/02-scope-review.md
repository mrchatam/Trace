# P18 / S05 / 02 — rebuild review / handoff close

## Metadata
- id: P18-S05-02
- todo_ids: [P18-S05-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Independent review of S05 rebuild vs FINAL. **Close DR-HANDOFF** (default **`no successor`**). S04-00 FINAL: S04 **starts** Notes only — this row owns close. Mark Phase 18 complete when APPROVE. Do **not** rewrite P17 or P18-00 history. Next runnable **none** unless this row records a successor.

**Stop if sibling `00-PLANNER.md` is still DRAFT.** It is **FINAL**.

## References
- Sibling [00-PLANNER.md](00-PLANNER.md) — **FINAL**
- Sibling [01-rebuild-binaries.md](01-rebuild-binaries.md)
- [../../DR-HANDOFF.md](../../DR-HANDOFF.md)
- S04 [VERIFY-NOTES.md](../scope-04-phase-verify/VERIFY-NOTES.md)
- P17 DR-HANDOFF (historical `no successor` — do not rewrite)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Fresh subagent ≠ S05-01. Unattended: no Plan-mode switch. Re-run rebuild cmds (preferred) **or** confirm mtime vs S05-01 Notes. **Must** re-run `./bin/trace version` and `./bin/trace-mcp -h`. Own DR-HANDOFF completion. Do **not** auto-scaffold Phase 19. Do **not** re-run S04 named DF suite as this row’s fail bar.

## Review focus
Reject APPROVE if binaries were not rebuilt this scope (mtime still 2026-08-17 17:32 / Notes claim rebuild but files unchanged). Reject treating a skipped DF-87 live check as a fail. Reject treating a **run-red** DF-87 live check as non-fail. Reject closing with a successor unless this row explicitly records one. Reject reversing DF-88. Reject claiming this is research S05.

## Checklist

| # | Check | How |
|---|--------|-----|
| 1 | `bin/trace` rebuilt with CGO1 | `go build` this session **or** mtime vs S05-01 Notes (must be newer than 2026-08-17 17:32) |
| 2 | `bin/trace-mcp` rebuilt with CGO0 | same |
| 3 | `-h` lists 10 tools including `trace_impact` | **must** re-run `./bin/trace-mcp -h` |
| 4 | `./bin/trace version` exits 0 | **must** re-run |
| 5 | No product feature Go (docs CGO line ok) | Diff — README/`help.go` MCP CGO0 one-liners allowed |
| 6 | DF-88 exclude not reversed | Grep seed export / omit test not deleted |
| 7 | P17 historical `no successor` intact | Read P17 DR-HANDOFF — do not edit |
| 8 | DF-87 live skip-vs-fail honored | If S05-01 skipped → non-fail. If S05-01 ran red → **reject**. Optional re-run |
| 9 | This phase DR-HANDOFF closed | Write close on [../../DR-HANDOFF.md](../../DR-HANDOFF.md) = **`no successor`** unless this row records otherwise |

## Locked verify (re-run)

**Must:**

```bash
./bin/trace version
./bin/trace-mcp -h
```

**Preferred** (independent rebuild):

```bash
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
GOMODCACHE=$HOME/go/pkg/mod GOPROXY=off CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp
```

Same 403 class as 00-PLANNER: prefix retry is **not** a product defect.

Optional DF-87: same recipe as 00-PLANNER. If S05-01 skipped, skip here too (**non-fail**). If S05-01 ran PASS, re-run is extra evidence.

## Exit criteria
- [x] APPROVE high (or medium with residuals listed)
- [x] REVIEW-NOTES.md
- [x] DR-HANDOFF closed = **`no successor`** unless this row records otherwise
- [x] Board Notes; **Phase 18 complete**; next runnable **none** (unless successor recorded)

## Minimal todos
- [x] Confirm binaries + catalog (`version` + `-h`)
- [x] REVIEW-NOTES.md + DR-HANDOFF close
- [x] Board sync
