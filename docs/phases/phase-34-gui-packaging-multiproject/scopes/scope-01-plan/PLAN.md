# Phase 34 S01 — PLAN

## Summary

Phase 34 ships the **real Explore SPA** inside the Trace binary (`go:embed` of built `web/dist` into `internal/httpapi/embeddist`) so consumer projects need **only** `.trace/` (L1/L2) — no required consumer `web/`. Default `trace gui` / `trace serve` bind auto-hops on `EADDRINUSE` from `127.0.0.1:7432` through at most 10 ports (`7432`–`7441`), prints the chosen URL, and `gui` opens that URL (L3); explicit `--addr` stays strict fail-if-busy via `flag.Changed`. Board order is fixed: **S02 embed → S03 auto-port → S04 docs → S05 VERIFY**. One process remains one project root (L4).

## Embed / release pipeline

- **Primary command(s):**
  - **Canonical Trace-root entrypoint:** `scripts/embed-gui.sh` (create in **S02**).
  - **`//go:generate`** on or near `internal/httpapi/embed.go` invoking that script (S02).
  - **Optional thin wrapper:** Trace-root `Makefile` target `embed-gui` calling the same script (S02 may add; none exists today).
- **Steps (lean A — lock):**
  1. From Trace repo root: `cd web && npm ci && npm run build` (`web/package.json` → `tsc -b && vite build` → `web/dist/**`).
  2. Replace contents of `internal/httpapi/embeddist/` with `web/dist/**`, **except** keep/regenerate a short `embeddist/README.md` (pipeline steps + last-resort stub note — **not** “two-artifact everyday”).
  3. `go build -o bin/trace ./cmd/trace` (or release equivalent) so `//go:embed all:embeddist` picks up real assets.
- **Contributor DX:** Run embed sync before `go build` when iterating on embedded UI; for Vite hot path, use disk `web/dist` (StaticDir opportunistic — see below) without re-embedding.
- **Future CI/release:** When CI exists, run embed sync **before** release `go build`. Today: **no** Trace-root `Makefile`, **no** Trace-root `.github/workflows` — do **not** invent a phantom job; S02 may schedule creating script/`go:generate`/(optional) Makefile only.
- **README / stub:**
  - S02 rewrites `embeddist/README.md` away from two-artifact-as-everyday.
  - Keep a stub path only when embed is empty/dev mistake; `placeholderHTML` remains after embed miss — tone update in S02 (consumer-facing copy must not teach requiring `web/` as primary).
- **VERIFY stub-fail (S05):** Fail if shipped/embedded UI is still stub while release intended full SPA.
  - **Stub markers:** body text `Embedded GUI stub` and/or missing `#root` / missing `/assets/` module script.
  - **Real SPA markers** (from live `web/dist/index.html`): `<div id="root">` + `script type="module"` under `/assets/`.
- **Files S02 may create (plan only — not this row):** `scripts/embed-gui.sh`; `//go:generate` directive on/near `embed.go`; optional root `Makefile` `embed-gui`; sync of real SPA into `embeddist/` (plus README rewrite); minimal help/usage honesty for static story if needed; `placeholderHTML` / stub tone.

## StaticDir resolution policy

Order (unchanged semantics): **disk-if-`index.html` → `embeddedIndexExists()` → inline `placeholderHTML`**. No Trace-module / GOPATH probe. Do **not** copy SPA into consumer `.trace/`.

| Context | Behavior |
|---------|----------|
| Empty `--static-dir` | Resolve candidate to `filepath.Join(absRoot, "web", "dist")` (keep today’s default path string). |
| Consumer project (no web/) | Candidate has no `index.html` → **embedded real SPA** (post-S02) → else placeholder. Consumers must **not** require `web/`. |
| Trace checkout with web/dist | Opportunistic: if `diskIndexExists(StaticDir)` → **disk wins** (contributor DX). |
| Explicit `--static-dir DIR` | Abs that path as StaticDir; still **refuse** StaticDir == project root; missing `index.html` → embed → placeholder. |
| Refuse root | `StaticDir == project root` refused (would expose `.trace/` / source) — keep current refusal. |

## Auto-port / bind policy

**Shared helper location (lock):** Implement listen/retry in **`internal/httpapi`** (prefer extending bind/listen near `bind.go` / `ListenAndServe`, or a small dedicated helper in that package). Both `trace gui` and `trace serve` in `cmd/trace/local_http.go` call the same helper — no divergent hop logic. Update `FormatAddrInUseMessage` for auto-exhausted vs strict `--addr` paths in S03.

| Case | Behavior |
|------|----------|
| Default addr free | Bind `127.0.0.1:7432` (`DefaultAddr`); print `http://127.0.0.1:7432/`; `gui` opens that URL unless `--no-open`. |
| Default addr busy | UA-increment: try `7433` … up to **max 10** attempts (**ports `7432`–`7441`** inclusive); on success print **chosen** URL; `gui` opens **chosen** URL. Host stays loopback from default. |
| Explicit `--addr` busy | **Strict fail** (friendly text may mention default auto). Includes `--addr 127.0.0.1:7432` when flag was set. |
| Exhausted | Fail: auto range exhausted + suggest `--addr`. |
| Commands covered | **Both** `gui` + `serve` (shared). |
| Detection of explicit addr | **`FlagSet` / `flag.Changed` on `addr`** — **not** string-equal to `DefaultAddr`. |
| serve print / gui open | Move `serve` listen URL print to **post-bind** success (today’s pre-listen print at `local_http.go` ~205 is S03 debt). Keep `gui` `OnListening` / open after successful listen; URL must match **actual** bound addr after hop. |

**Out of S03:** No public/`0.0.0.0` default; no OS `:0` as Trace default; no always-on daemon.

## Test matrix

| ID | Scope | Assertion |
|----|-------|-----------|
| T1 | S02 | Temp consumer root with `.trace/` init, **no** `web/` → `GET /` (or static index) is **real SPA** (`#root` and/or `/assets/` module script); **not** stub phrase `Embedded GUI stub`. |
| T2 | S02 | Root with planted `web/dist/index.html` marker → **disk wins** over embed. |
| T3 | S02 | `StaticDir == project root` still refused. |
| T4 | S03 | Default bind busy → next port (e.g. `:7433`) + printed URL matches bound addr. |
| T5 | S03 | Two concurrent default `gui`/`serve` → distinct ports; `gui` open/hook sees chosen addr. |
| T6 | S03 | Explicit `--addr` busy → **fail** (including `--addr` equal to DefaultAddr string when `Changed`). |
| T7 | S03 | Auto exhausted (10 busy ports) → fail mentioning range + `--addr`. |
| T8 | S04 | Help/usage/quickstart no longer teach consumer two-artifact primary or “no auto free-port” for default. |
| T9 | S05 | VERIFY floor: T1 + concurrent free-port (T4/T5) + docs consumer story (T8); evidence under `experiments/runs/`. |
| T10 | S02/S05 | After intended release embed pipeline, embedded `index.html` must not be stub-only (VERIFY stub-fail criteria). |
| T11 | S03 | Default free path still binds `:7432` first (no unnecessary hop when free). |

## Docs / help touch list (S04)

| Path | Primary owner | Notes |
|------|---------------|-------|
| `docs/gui-quickstart.md` | **S04** | Primary consumer story: `.trace/` only + embedded SPA; drop two-artifact primary teaching. |
| `cmd/trace/help.go` | **S03** (port) + **S04**/minimal **S02** (static) | Remove “no auto free-port”; static story honesty. |
| `cmd/trace/local_http.go` usage strings | **S03** + **S02**/S04 | Prefer embed for consumers; auto-port on default; `--addr` strict. |
| `internal/httpapi/addr_in_use.go` | **S03** | Distinct messaging for auto-exhausted vs explicit busy. |
| `internal/httpapi/embeddist/README.md` | **S02** rewrite + **S04** polish | Pipeline + last-resort stub — not everyday two-artifact. |
| `web/README.md` | **S04** (contributor-labeled) | Contributor DX; label auto-port / multi-process correctly. |
| Root `README.md` | **S04** if needed | Only if quickstart link implies consumer `web/`. |
| `AGENTS.md` | **S04** / phase close optional | Align Phase 34 focus notes after ship. |
| `placeholderHTML` / stub copy | **S02** tone; **S04** if user-visible | Last-resort only; do not teach consumer `web/` as primary. |

## Out of scope / deferred

- install-sidecar **(B)** unless Node unblockable in release path (not evidenced — reject as primary).
- SPA copy under consumer `.trace/` (L1/L2 reject).
- public bind default; OS `:0` as Trace default; always-on daemon; hosted SaaS.
- Explore UI / craft redesign (Phase 33 closed — spawn forward only if needed later).
- Reversing board order (auto-port before embed) — **forbidden**.
- Manual-`cp`-only **(C)** as primary release path.
- Trace-module path probe for StaticDir.
- Creating CI workflows in this phase unless a later row explicitly schedules them (S02 may add script/`go:generate`/optional Makefile only).

## Handoff

- **S02:** Real SPA into `embeddist` via `scripts/embed-gui.sh` + `go:generate` (+ optional `make embed-gui`); keep StaticDir opportunistic disk→embed→placeholder; rewrite embeddist README; stub/placeholder tone; tests T1–T3 (+ T10 seed); **no** auto-port.
- **S03:** Shared httpapi listen helper for `gui`+`serve`; UA-increment `7432`–`7441`; `flag.Changed` strict `--addr`; post-bind `serve` print; `gui` open chosen URL; update addr-in-use / usage / help port lines; tests T4–T7 (+ T11).
- **S04:** Docs/help touch list above; T8; contributor vs consumer labeling; polish embeddist README if needed.
- **S05 VERIFY floor:** T1 (consumer no `web/` → real SPA) + concurrent free-port (T4/T5) + docs consumer story (T8) + stub-fail if embed still stub when full SPA intended (T9/T10); evidence under `experiments/runs/`; DR-HANDOFF close when gate passes.
