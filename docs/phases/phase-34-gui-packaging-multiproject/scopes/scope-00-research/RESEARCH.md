# Phase 34 S00 — RESEARCH

## Summary

Consumer dogfood shows a stub UI because default `StaticDir` is `<root>/web/dist`, app repos have no `web/`, and `embeddist` is still a stub HTML page — not the Phase 33 Explore SPA. L2 lean is a release/`go generate`/make/CI pipeline that populates `embeddist` with real Vite output and ships it via existing `//go:embed all:embeddist`; install-sidecar is unnecessary unless that pipeline proves unblockable. L3 requires default-bind auto free-port (UA-style increment from `7432`, loopback host fixed, max ~10 attempts); explicit `--addr` stays fail-if-busy (`flag.Changed`). Docs/help still teach two-artifact + “no auto free-port”; those lines are debt for S02/S03/S04. P33/P32 rejected or deferred UA auto-port; **L3 overturns that only for the default bind happy path**.

## Live baseline (verified)

| Area | Fact | Evidence path |
|------|------|---------------|
| Static order | Prefer disk when `StaticDir/index.html` exists → else `embeddedIndexExists()` → else inline `placeholderHTML` | `internal/httpapi/static.go` |
| Default StaticDir | Empty flag → `filepath.Join(abs, "web", "dist")`; refuse StaticDir == project root | `internal/httpapi/server.go`; flag help `cmd/trace/local_http.go` |
| Embed contents | **Stub only** — `index.html` + `README.md` (no Vite assets) | `internal/httpapi/embeddist/` |
| Embed pipeline today | Manual `npm ci && npm run build` + `rm`/`cp -a dist → embeddist` then `go build`; README says everyday path remains two-artifact | `internal/httpapi/embeddist/README.md` |
| Makefile/CI embed sync | **None** — no Trace-root `Makefile`; no Trace-root `.github/workflows` | repo root (2026-08-21) |
| Port on conflict | Fail; stderr `FormatAddrInUseMessage` with manual `--addr 127.0.0.1:7433` example | `internal/httpapi/addr_in_use.go`; `cmd/trace/local_http.go` |
| Default bind | `127.0.0.1:7432` | `internal/httpapi/bind.go` (`DefaultAddr`) |
| gui/serve help | Prefer `web/dist`; “does not auto-pick a free port” / “no auto free-port” | `cmd/trace/local_http.go` usage; `cmd/trace/help.go` |
| Docs story | Disk wins when present; optional embeddist copy; everyday use without rewriting embeddist | `docs/gui-quickstart.md` |
| gui open | After listen: print URL; best-effort browser open unless `--no-open` | `cmd/trace/local_http.go` (`OnListening`) |
| serve open | Prints listen URL; no browser open | same file (`localHTTPServe`) |

**Drift vs planner-locked facts:** None material. Confirmed stub-only embeddist, disk→embed→placeholder order, no root Makefile/CI embed target, fail+manual `--addr` messaging.

## L3 supersession (P32/P33 → P34)

- **What P32 said:** Port research preferred friendly `EADDRINUSE` + `--addr` examples as Phase 32 minimum; treated UA auto-increment as **optional / defer #2**, not the default story ([`phase-32…/scopes/scope-00-research/RESEARCH.md`](../../../phase-32-graph-first-gui/scopes/scope-00-research/RESEARCH.md) — “**defer/#2** auto-increment (optional)”).
- **What P33 said (reject):** Peer matrix: “**Reject** UA auto-port for Trace. **Borrow** Trace/P32 explicit `--addr` multi-project story.” Rejected alternatives: “**UA auto-increment listen port** — conflicts with P32-PORT / multi-project `--addr` story; reject silent port hopping.” ([`phase-33…/scopes/scope-00-research/RESEARCH.md`](../../../phase-33-gui-craft-hook-launch/scopes/scope-00-research/RESEARCH.md)).
- **What L3 overturns:** For **default bind only** (`127.0.0.1:7432` when the user did not set `--addr`), busy → auto try next free **loopback** port; print chosen URL; `gui` opens that URL. Documented in [`DESIGN-LOCKS.md`](../../DESIGN-LOCKS.md) L3 + clarifications.
- **What stays:** Explicit `--addr` → **strict** fail-if-busy (friendly message may mention default auto behavior); loopback default; Law 19 (HTTP/UI adapters only); L4 one process = one project root; no public/`0.0.0.0` default; no always-on daemon; no hosted SaaS.

## Embed options

| Option | Pros | Cons | Rank / Recommend? |
|--------|------|------|-------------------|
| **(A) go:embed + release/make/CI/`go generate` pipeline** populate `embeddist` from `web/dist` | Matches existing `//go:embed all:embeddist`; single binary for consumers; L2 lean; stub stays last-resort only when pipeline skipped | Needs Node in release/CI; binary size grows with SPA; must keep `embeddist` sync honest in VERIFY | **#1 Recommend** |
| **(B) Install-sidecar beside binary** | Smaller Go binary; swap UI without rebuild | Second artifact for installers/PATH layouts; easy to desync; still must not require consumer `web/` | #2 only if A unblockable |
| **(C) Manual `cp` only (status quo)** | Already documented | Everyday path stays two-artifact; release often ships stub; dogfood failure mode | **Reject** as primary |
| Other: copy SPA under consumer `.trace/` | Would “fix” missing disk without Trace rebuild | Violates L1/L2 reject (`.trace/` = store/lock/token) | **Reject** |

**Recommendation for S01/S02:** **(A)** — add a Trace-root sync target (e.g. `go generate` and/or `make embed-gui` / CI release step): `web` build → replace `internal/httpapi/embeddist` contents (keep a short README or generate note) → `go build`. S05 VERIFY must fail if shipped embed is still stub when a full SPA was supposed to be built in. Sidecar only if generate/CI cannot run Node in the release path (not evidenced today).

## StaticDir policy recommendation

- **Consumer default:** Empty `--static-dir` resolves to `<root>/web/dist` as a **candidate only**. If `index.html` is absent (typical app repo), fall through to **embedded real SPA** (post-S02). Consumers must **not** need `web/`.
- **Trace-checkout DX:** Same opportunistic rule — when contributor builds `web/dist`, **disk wins** so Vite iteration is immediate without re-embedding. Detection rule: **no special Trace-module probe required** — `diskIndexExists(StaticDir)` is enough and matches current `static.go`. Optional later hardening: only treat default disk path as “contributor DX” in docs; code can stay opportunistic for any `-C` that happens to have `web/dist`.
- **`--static-dir` explicit:** Always use that absolute path as StaticDir (still refuse == project root). If its `index.html` is missing, fall through to embed (current behavior), then placeholder.
- **Reject:** Requiring consumer `web/`; documenting two-artifact as consumer primary; copying SPA into consumer `.trace/` as primary asset path.
- **Note:** Today’s order is already “disk if present → embed”; the product bug is **stub embed**, not “wrong order.” S02 ships real embed; StaticDir default path can stay unless PLAN wants empty StaticDir meaning “skip disk probe” (unnecessary if embed is real).

## Auto-port algorithm

- **Candidates compared:**
  - **UA-increment:** start at configured/default port, `+1` on `EADDRINUSE` while attempts remain and port not explicit; host fixed `127.0.0.1` (`viewer.mjs` `listen`, attempts=10, `portExplicit`).
  - **Scan range:** try `[7432, 7432+R)` then fail — same idea as UA with explicit ceiling.
  - **`:0` then advertise:** OS free port; skips preferred `7432` even when free — worse discoverability/docs.
- **Recommend:** **UA-style increment** mapped to Trace: start at port from default addr (`7432`), host stay **loopback from default** (`127.0.0.1`), max **10** attempts (ports `7432`…`7441`), then fail with an updated message that mentions exhausted auto range + `--addr`.
- **Default bind behavior:** Auto only when user did **not** set `--addr`. On each successful bind, print `http://<chosen>/`; for `gui`, open that URL unless `--no-open`.
- **Explicit `--addr` detection:** Prefer **`FlagSet.Lookup("addr").Changed` / visit** (or equivalent after parse) — **not** string-equal to `DefaultAddr`. Reason: `--addr 127.0.0.1:7432` must remain strict fail-if-busy.
- **Print + open:** Share bind retry for **`gui` and `serve`**. Keep product difference: `gui` opens browser; `serve` prints only (align print timing with **post-bind** success — today’s serve pre-listen print should move with S03).
- **Shared vs split:** **Shared** auto-port path (same local HTTP entry in `local_http.go` / httpapi listen helper) unless IMPLEMENT discovers a hard split; L3 prefers shared.
- **Peer borrow/reject (`viewer.mjs`):**
  - **Borrow:** increment-on-`EADDRINUSE`, attempt budget (~10), loopback bind, skip hop when port/addr explicit, print URL, best-effort open + `--no-open`.
  - **Reject / map:** UA `--port 0` = OS free port — do **not** make Trace default; optional later only. UA tokens in URL — Trace keeps existing token policy (unchanged by this research).

## Consumer layout / docs audit

| Path or doc | Issue vs L1/L2/L3 | Fix owner (S02/S03/S04) |
|-------------|-------------------|-------------------------|
| `docs/gui-quickstart.md` | Disk wins; optional embeddist; everyday without rewriting embeddist — consumer two-artifact teaching | **S04** (primary); S02 may touch if needed for honesty |
| `internal/httpapi/embeddist/README.md` | “Supported everyday path remains two-artifact” | **S02** (pipeline README rewrite) + S04 polish |
| `cmd/trace/help.go` | “Serves web/dist”; “Port conflict… no auto free-port” | **S03** (port lines) + **S04**/minimal S02 (static story) |
| `cmd/trace/local_http.go` | usage: prefer `web/dist`; “does not auto-pick a free port” | **S03** + **S02**/S04 |
| `internal/httpapi/addr_in_use.go` | Manual `--addr` only (expected until auto-port) | **S03** (message for auto vs strict) |
| `web/README.md` | Second root → distinct `--addr`; “no auto free-port”; production build as default StaticDir story | **S04** (contributor-labeled); S03 for port |
| Root `README.md` | Points at gui-quickstart only — low debt | **S04** if quickstart link text implies consumer `web/` |
| `AGENTS.md` | Phase 34 focus already states consumer `.trace/` + embed + auto-port; still notes P33 two-artifact until P34 | **S04**/phase close — optional; not consumer-facing primary |
| Phase 34 README “Live repo baseline” | Accurate snapshot of pre-fix debt — not consumer doc | Cite only; optional refresh at VERIFY |
| `internal/httpapi/static.go` `placeholderHTML` | Still teaches `cd web && npm run build` / `web/dist` | **S02** tone for last-resort; **S04** if user-visible |
| `internal/httpapi/embeddist/index.html` stub copy | Teaches build disk / copy embeddist | Replaced by real SPA in **S02**; stub last-resort copy **S02**/S04 |
| CONTRIBUTING.md | No gui/`web/dist` consumer teaching found | None |

## Rejected alternatives (short)

1. **Requiring consumer `web/`** or documenting two-artifact as consumer primary — violates L1/L2; root cause of stub dogfood.
2. **Copying SPA into consumer `.trace/`** as primary asset path — L2 reject; `.trace/` is store/lock/token only.
3. **Keeping “no auto-port” for default bind** — conflicts with L3; overturns P33 reject for happy path only.
4. **Auto-hopping when user set explicit `--addr`** — L3 keeps pin strict; use `flag.Changed`, not DefaultAddr string compare.
5. **Public / `0.0.0.0` default; always-on daemon; hosted SaaS** — Law 19 / Phase 29 carve-out / out of scope.

## Handoff to S01 / S02 / S03

- **S01:** PLAN locks — (1) embed pipeline commands (`go generate`/make/CI) + VERIFY stub-fail; (2) StaticDir table (opportunistic disk-if-present → real embed → placeholder; `--static-dir` explicit; refuse root); (3) auto-port params (start `7432`, host `127.0.0.1`, max 10, `flag.Changed` for strict); (4) test matrix seeds: consumer temp root no `web/` → real SPA marker; default busy → `:7433` + printed/opened URL; `--addr` busy → fail; Trace checkout with `web/dist` → disk wins.
- **S02:** Real SPA into `embeddist` via pipeline; keep/confirm StaticDir resolution; rewrite embeddist README away from two-artifact everyday; minimal help honesty OK; **no** auto-port.
- **S03:** Shared `gui`+`serve` auto-port; update `FormatAddrInUseMessage` / usage / help; post-bind print; `gui` open chosen URL; explicit `--addr` unchanged fail path.
