# P34-S03-01 — Implement auto free-port

## Metadata
- id: P34-S03-01
- todo_ids: [P34-S03-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: []
- verification: automated
- hooks: []

## Objective

Ship **L3 auto free-port** for default `trace gui` / `trace serve` bind: when the default listen address is busy, hop to the next free **loopback** port (UA-increment), print the **chosen** URL, and (for `gui`) open the browser to that URL. Explicit `--addr` stays **strict fail-if-busy**. Follow [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md). Do **not** change StaticDir/embed (S02 done). Law 19 / loopback / refuse-remote unchanged.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — **L3, L4** (do not reopen)
- [00-PLANNER.md](00-PLANNER.md)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — auto-port table + T4–T7, T11
- Peer algorithm (shape only): `similar projects/.../viewer.mjs` `listen(port, 10)` + `portExplicit`
- Live anchors (read before coding):
  - `internal/httpapi/bind.go` — `DefaultAddr = "127.0.0.1:7432"`
  - `internal/httpapi/addr_in_use.go` — `IsAddrInUse`, `FormatAddrInUseMessage` (manual `--addr` only today)
  - `internal/httpapi/server.go` — `Options`, `New`, `ListenAndServe` (**single** `net.Listen(s.addr)`; no hop; `OnListening(s.addr)` after bind)
  - `cmd/trace/local_http.go` — shared `cmdLocalHTTP` for `gui`+`serve`; `addr` flag default `DefaultAddr`; **serve prints URL pre-listen** (~L205–207) — S03 debt; `gui` prints+opens in `OnListening`
  - `cmd/trace/help.go` — gui line still says “no auto free-port”
  - `cmd/trace/gui_test.go`, `serve_test.go` — busy `--addr` friendly fail; `setGUIListenHook` / `notifyContext` patterns

## Session start

Follow agent-loop-protocol Session start. Unattended: execute locks below; do not reopen L3 or invent OS `:0` / public defaults.

## Locked defaults (FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Algorithm | **UA-increment:** start at port from configured/default addr; on `EADDRINUSE` (`IsAddrInUse`) and **not** explicit `--addr`, retry **port+1** same host; **max 10** attempts → ports **`7432`–`7441`** inclusive when starting from `DefaultAddr` |
| Host | Keep parsed host from the initial addr (default `127.0.0.1`). Do **not** flip to `0.0.0.0` / public. Do **not** weaken `RefuseRemote` |
| Explicit `--addr` | **Strict:** one listen attempt; on busy → fail. Detect via **`fs.Lookup("addr").Changed`** after `Parse` — **never** string-equal to `DefaultAddr`. Includes `--addr 127.0.0.1:7432` when Changed |
| Shared helper | Implement hop in **`internal/httpapi`** (extend `ListenAndServe` and/or small helper next to `bind.go` / `addr_in_use.go`). Both `gui` and `serve` call the same path — **no** divergent hop logic in `cmd/trace` |
| Wire explicit flag | Add `Options` field (e.g. `AddrExplicit bool` / `StrictAddr bool`). `cmdLocalHTTP` sets it from `flag.Changed` before `httpapi.New` |
| After successful bind | Update server’s bound addr (`Addr()` / `OnListening` arg) to the **chosen** `host:port` before `OnListening` |
| Print | **Always** print actual listen URL **post-bind**. Move `serve` stderr print from pre-listen into post-bind (`OnListening` or shared callback). `gui` already post-bind — ensure URL matches hop |
| Open | `gui` opens `http://<chosen>/` unless `--no-open`; hook (`callGUIListenHook`) must see **chosen** addr |
| Exhausted (default) | Fail after 10 busy attempts; stderr mentions auto range (`7432`–`7441` or “10 ports”) **and** `--addr` tip |
| Strict busy | Keep/adjust `FormatAddrInUseMessage` for explicit path; friendly text **may** note that default bind auto-hops |
| Help / usage | Remove “no auto free-port” / “does not auto-pick” for **default**. State: default busy → next free port; `--addr` pins and fails if busy. **S04** still owns full quickstart / AGENTS flip |
| Out | Embed/StaticDir; `docs/gui-quickstart.md`; OS `:0` as Trace default; always-on daemon; multi-root one process (L4) |

### Answers locked for implementer (planner gate)

1. **Exact algorithm + max attempts:** UA-increment from PLAN — seed port from initial listen addr (default **7432**), host **127.0.0.1** from default, **+1** on `IsAddrInUse` while attempts remain and `!AddrExplicit`, **max 10** tries (`7432`…`7441`). Non-`EADDRINUSE` errors fail immediately (no hop).
2. **Shared helper location:** **`internal/httpapi`** (not a `cmd/trace`-only loop). Prefer extending `Server.ListenAndServe` using `Options.AddrExplicit` + `IsAddrInUse` + port increment via `ParseListenAddr` / `net.JoinHostPort`. CLI only passes the flag and unifies post-bind print.
3. **Tests:** PLAN **T4–T7** + **T11** required. **T5:** two concurrent default binds → distinct ports; `gui` open/hook sees chosen addr (use existing `setGUIListenHook` / `notifyContext` patterns). Prefer `httpapi` unit tests for hop/exhaust; CLI tests for `flag.Changed` + print/open. Avoid `t.Parallel` on tests that bind `7432`–`7441`.

## Role work

### Minimal todos (execute in order)

1. **httpapi API** — Add `AddrExplicit` (name OK if clear) on `Options`; store on `Server`. Constant for max attempts (`10`) next to `DefaultAddr` or helper.
2. **Listen hop** — In `ListenAndServe` (or helper it calls): loop listen → on `IsAddrInUse` && `!AddrExplicit` && attempts left → increment port, update `s.addr`, retry; on success call `OnListening` with **final** addr; on exhaust return error distinguishable for CLI messaging (or let CLI format from attempt count / sentinel).
3. **Messages** — Split or parameterize:
   - explicit busy → `FormatAddrInUseMessage` (update copy: default auto-hops; pin with `--addr`)
   - default exhausted → dedicated message (range + `--addr`)
4. **CLI wire** — `cmdLocalHTTP`: after `fs.Parse`, `addrExplicit := fs.Lookup("addr").Changed`; pass into `Options`. Remove **pre-listen** serve print; set `OnListening` for **both** modes (serve: print only; gui: print + hook + optional open) so URL is always post-bind chosen addr.
5. **Help / usage** — Update `help.go` gui port line + `local_http.go` `usageGUI` / `usageServe` second-project wording (auto on default; `--addr` to pin).
6. **Tests (TDD-friendly):**
   - **T11** — Nothing on `:7432` → default bind stays `127.0.0.1:7432` (no unnecessary hop).
   - **T4** — Occupy `:7432` → default bind → `:7433` (or next free); printed/OnListening addr matches.
   - **T5** — Two concurrent default `gui`/`serve` (short-lived via `notifyContext`) → distinct ports; gui hook/open URL uses chosen addr.
   - **T6** — Explicit `--addr` on occupied addr (including `--addr` equal to `DefaultAddr` string) → **fail**, no hop. Keep/extend `TestGuiAddrInUseFriendlyMessage` / `TestServeAddrInUseFriendlyMessage`.
   - **T7** — Occupy all of `7432`–`7441` → default → fail mentioning range + `--addr`.
7. **Regression** — Refuse-remote / token / StaticDir / embed untouched. Focused: `go test ./internal/httpapi/ ./cmd/trace/` (or narrower).
8. **Board** — Notes: files touched, test names + pass evidence.

### Out of this row

- Full docs/quickstart / AGENTS / `web/README.md` flip (**S04**)
- VERIFY / DR-HANDOFF (**S05**)
- Embed pipeline / StaticDir policy (**S02** done)
- Explore UI redesign

## Exit criteria

- [ ] Default busy → auto free loopback port (T4); free path stays `:7432` (T11)
- [ ] Two concurrent defaults → distinct ports + correct gui open/hook URL (T5)
- [ ] Explicit `--addr` busy still fails, including DefaultAddr string when Changed (T6)
- [ ] Auto exhausted → fail with range + `--addr` (T7)
- [ ] `gui` and `serve` share httpapi hop; serve print is post-bind
- [ ] Help/usage no longer claim “no auto free-port” for default
- [ ] Board Notes cite tests + key files
- [ ] Next **P34-S03-02**

## Todo updates

Status + notes on **P34-S03-01** only.

## Next

`P34-S03-02`
