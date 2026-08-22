# P33-S02-02 — Launch review

**Verdict:** **PASS**  
**Confidence:** **high**  
**Date:** 2026-08-21  
**Scope:** `trace gui` + PATH teach (Theme C); Law 19; land Explore `/`  
**Skills loaded:** code-review-and-quality  
**Evidence:** `go test ./cmd/trace/ -count=1` ok; targeted Gui/Serve + httpapi AddrInUse/RefuseRemote ok

## Checklist evidence

### Theme C / CLI
- [x] `case "gui":` in `root.go` → `cmdGui` → shared `cmdLocalHTTP(..., localHTTPGUI)`
- [x] Help documents `gui` + `--no-open` + default `127.0.0.1:7432`; no global `-gui`/`--gui`
- [x] `serve` retained (`cmdServe` / `localHTTPServe`); gui does not remove serve
- [x] Law 19: CLI adapter only — `httpapi.New` + `ListenAndServe`; `OnListening` for open; no business fork

### Bind / port / security
- [x] Default `httpapi.DefaultAddr` (`127.0.0.1:7432`); `RefuseRemote` without `--allow-remote`
- [x] Token posture shared with serve (generated bearer off-loopback when allowed)
- [x] Addr-in-use → `IsAddrInUse` + `FormatAddrInUseMessage`; **no** auto-port / silent hop (tests + usage copy)
- [x] Opt-in foreground listen (signal cancel); not an always-on daemon

### Browser open
- [x] `OnListening` after successful `net.Listen` (not sleep race)
- [x] Always prints `trace gui: listening on http://…` on successful bind
- [x] Open fail → tip + continue; exit OK on cancel (`TestGuiOpenFailStillListens`)
- [x] `--no-open` skips opener (`TestGuiNoOpenDoesNotCallOpener`)
- [x] Open URL = `http://{addr}/` — not `/overview` (`TestGuiOpenURLLandsOnExploreRoot`)

### PATH ≠ agents install
- [x] Help Build note: `CGO_ENABLED=1 go install github.com/mrchatam/Trace/cmd/trace@latest` then `trace gui`
- [x] Contributor `go build -o bin/trace` secondary
- [x] Explicit: `trace install …` ≠ PATH; `cmdInstall` unchanged (agents/MCP/hooks)
- [x] Minimal tip in `docs/gui-quickstart.md` §Install CLI on PATH (S05 still owns full primary flip)

### Docs / scope boundaries
- [x] No full quickstart rewrite; no Explore seed/graph (S03); no shell colorize (S04)

### Tests / evidence
- [x] Help / remote refuse / addr-in-use / no-open / land `/` / open-fail covered in `gui_test.go`
- [x] P33-S02-01 Notes cite `go test ./cmd/trace/ -count=1`

## Findings

| Sev | Finding | Disposition |
|-----|---------|-------------|
| low | `FormatAddrInUseMessage` still says `serve:` / `trace serve` when gui hits conflict (implement prompt allowed either wording) | Optional dual-word in **S05** residual; does not block — still tells user to pick `--addr` |
| nit | `OnListening` reports configured `s.addr`, not `ln.Addr()` | Irrelevant under reject-`:0` / auto-port; leave as-is |

**No blocker / high.** No spawn rows (`P33-S02-02a`/`02b`).

## Upcoming thickenings (reviewer rights)

- **S03** — CLI land contract: `trace gui` opens `http://{addr}/`; Explore must stay at `/` (do not relocate graph to `/overview` without CLI change — out of S03).
- **S05** — PATH + short `trace gui` tip already shipped; rewrite lead from `./bin/trace serve` → `trace gui`; keep serve secondary; optional addr-in-use copy mentions gui|serve.

## Next runnable

**P33-S03-00**
