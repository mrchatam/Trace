# P33-S02-01 — Implement `trace gui` + PATH

## Metadata
- id: P33-S02-01
- todo_ids: [P33-S02-01]
- role: implementer
- skills: [incremental-implementation]
- mcps: []
- verification: automated
- hooks: []

## Objective

Ship **`trace gui`**: bind project (cwd/`-C`/`--root`), start local GUI HTTP (reuse serve + `httpapi`), print URL, **best-effort open default browser** to Explore **`/`**. Deliver PATH teach story (`go install` primary ≠ `trace install`). Loopback + P32-PORT unchanged. **Law 19** — CLI adapter only; no business-logic fork. **No Explore UI work** (S03). **No full quickstart primary flip** (S05).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — Theme C + PATH ≠ `trace install`
- [00-PLANNER.md](00-PLANNER.md) — locked defaults
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md) — Theme C CLI shape; reject UA auto-port
- [`../scope-01-design-ux/UX-IA.md`](../scope-01-design-ux/UX-IA.md) — land on Explore `/`, not `/overview`
- Live: `cmd/trace/serve.go`, `cmd/trace/root.go`, `cmd/trace/help.go`, `cmd/trace/serve_test.go`, `cmd/trace/install.go`, `internal/httpapi/` (`DefaultAddr`, `RefuseRemote`, `FormatAddrInUseMessage`, `ListenAndServe`)
- Peer pattern (borrow open only): `similar projects/Understand-Anything/…/packages/viewer/bin/viewer.mjs` (`--no-open`; `open`/`start`/`xdg-open`) — **reject** its auto-port

## Session start

Follow agent-loop-protocol Session start (Agent → clarify if blocked → Plan → execute). Prefer RESEARCH Theme C; do not reopen Themes A–B or S01 craft.

## Locked defaults (planner — do not re-debate)

| Item | Value |
|------|-------|
| Primary CLI | Subcommand **`trace gui`** only for S02. **Do not** add global `-gui` / `--gui` (secondary deferred; Theme C prefers subcommand) |
| Keep | `trace serve` unchanged for scripting/CI/headless |
| Flag surface | Inherit serve: `--addr`, `--allow-remote`, `--token`, `--token-file`, `--root`, `--static-dir`, `--cors-origin` + **`--no-open`**. Global `-C` / `--project` parity with serve |
| Defaults | Loopback **`127.0.0.1:7432`** (`httpapi.DefaultAddr`); refuse non-loopback without `--allow-remote`; token posture identical to serve |
| Port conflict | Reuse **P32-PORT** (`IsAddrInUse` + `FormatAddrInUseMessage`) — **fail**; user picks distinct `--addr`. **Reject** UA auto-increment / silent port hop / treating `:0` as multi-project |
| Listen reuse | Share serve’s `httpapi.New` + listen path (extract small shared helper in `cmd/trace` if needed). Prefer one code path for bind/policy; do not fork httpapi business behavior |
| Browser | After **successful TCP listen**, best-effort open; **always** print URL to stderr (or stdout consistent with serve’s listen line). Open failure → stderr tip with URL; **exit 0 if listening**. `--no-open` skips open entirely |
| Open impl | Small helper in `cmd/trace` (e.g. `openBrowser(url string) error`): darwin `open`, windows `cmd /c start`, else `xdg-open` via `os/exec` — **no new module deps**. Injectable opener for tests |
| Landing URL | **`http://{addr}/`** (Explore Graph home per UX-IA / App routes). **Never** force `/overview`. Optional trailing path only if it still routes to Graph (`/` or redirect-equivalent); default **`/`** |
| PATH #1 | Document **`go install github.com/mrchatam/Trace/cmd/trace@latest`** (or pinned tag) in help Build note + minimal install tip |
| PATH #2 | Contributor: existing `CGO_ENABLED=1 go build -o bin/trace ./cmd/trace` + optional symlink tip — **not** headline. No new Makefile required (repo has none at root) |
| PATH ≠ install | **Never** teach `trace install …` as putting `trace` on PATH (`install.go` = agents/MCP/hooks only) |
| Docs scope | Help accurate for `gui`; **minimal** PATH/`gui` tip OK in `docs/gui-quickstart.md` (short section). **S05** owns rewriting quickstart so `trace gui` is primary story |
| Law 19 | CLI → `httpapi` / libraries only |
| Out | Auto free-port; public bind default; always-on daemon; Explore seed/graph UX; shell colorize; packaging (brew/deb); conflating PATH with `trace install` |

### Must-answer locks (resolved by S02-00)

1. **Flags** — serve set + `--no-open`; no `-gui` this scope.
2. **Open** — platform `open`/`start`/`xdg-open`; mockable; CI uses `--no-open`.
3. **PATH** — `go install …/cmd/trace@…` primary in help (+ short docs tip); contributor build/symlink secondary.
4. **Tests** — help; remote refuse; addr-in-use; `--no-open` / mocked open URL = `/`; no auto-port.
5. **Land** — open `http://127.0.0.1:7432/` (or chosen `--addr`) → Explore `/`.

## Role work

### 1. Wire `gui` command

- Add `case "gui":` in `cmd/trace/root.go` → `cmdGui`.
- Implement `cmdGui` mirroring `cmdServe` flag parse + root/`--root` conflict rules + `httpapi.New` options.
- Add `--no-open` bool (default false = open).
- After successful listen: print listen line (prefer `trace gui: listening on http://…`); open `http://{addr}/` unless `--no-open`.
- On open error: stderr tip (URL + how to open manually); continue serving; exit code still success on clean shutdown.
- On addr-in-use: same friendly message path as serve (update hint text to mention `trace gui` **or** keep serve wording — either OK if still tells user to pick `--addr`).

**Open-after-listen:** Today `ListenAndServe` binds then blocks. Prefer a minimal extension (e.g. optional `OnListening` callback after `net.Listen` succeeds, or shared CLI helper that listens then invokes callback then serves) so the browser opens only when bind succeeded — **not** a fixed sleep race. Do not auto-retry other ports on failure.

### 2. Help

- Document `gui` in `cmd/trace/help.go` (flags incl. `--no-open`; note opens browser to GUI; keep `serve` entry).
- Build note: add PATH line for `go install github.com/mrchatam/Trace/cmd/trace@latest` (CGO note may still apply for full analyzers — be honest if `go install` needs `CGO_ENABLED=1` for this module).

### 3. PATH deliverable (docs/help only — no agents-install changes)

- Help Build note = primary teach.
- Optional ≤10-line tip in `docs/gui-quickstart.md` (e.g. “Install CLI on PATH”) with `go install` then `trace gui` — **do not** replace the whole `./bin/trace serve` walkthrough (S05).
- Do **not** change `cmdInstall` / install targets to place binary on PATH.

### 4. Tests (`cmd/trace/gui_test.go` and/or extend `serve_test.go`)

| Case | Expect |
|------|--------|
| Help / `gui --help` | Mentions `gui`, `--no-open`, listen defaults |
| Remote refuse | `gui --addr 0.0.0.0:…` without `--allow-remote` fails (mirror serve) |
| Addr in use | Occupied `--addr` fails with in-use messaging; **no** second port attempt |
| `--no-open` | Opener not called; process still binds (short-lived / cancel ctx pattern as serve tests allow) |
| Open URL | When opener invoked, URL is `http://{addr}/` — **not** `…/overview` |
| Open fail | Injected opener error → non-zero **not** required if listen OK (assert tip / still exitOK on cancelled listen if testable) |

Prefer injectable `openBrowser` var/func for unit tests; do not require a real GUI browser in CI.

### 5. Self-check

- `go test ./cmd/trace/ -count=1` (or targeted gui/serve tests) green.
- Manual smoke optional: `trace gui --no-open` prints URL ending in `/`.

## Exit criteria

- [ ] `trace gui` dispatched; help documents it + `--no-open`
- [ ] Reuses serve/httpapi listen + loopback/remote/token/P32-PORT behavior; **no** auto-port
- [ ] Best-effort browser open to Explore **`/`**; open fail ≠ listen fail; `--no-open` works
- [ ] PATH story: `go install github.com/mrchatam/Trace/cmd/trace@…` in help (+ minimal docs tip OK); **not** via `trace install`
- [ ] Tests cover help / remote refuse / in-use / no-open / landing path
- [ ] Board Notes with evidence commands (`go test …`, sample help snippet)
- [ ] No Explore graph UX changes; no S05 full docs rewrite

## Minimal todos

- [ ] Shared listen/gui wiring + `cmdGui` + root/help
- [ ] `openBrowser` helper + `--no-open` + land on `/`
- [ ] PATH teach in help (+ optional short quickstart tip)
- [ ] Tests per table above
- [ ] Board Notes with evidence

## Todo updates

Status + notes on **P33-S02-01** only.

## Next

`P33-S02-02`
