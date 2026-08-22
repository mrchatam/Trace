# P34-S02-01 — Implement embed + static defaults

## Metadata
- id: P34-S02-01
- todo_ids: [P34-S02-01]
- role: implementer
- skills: [incremental-implementation, tdd]
- mcps: []
- verification: automated
- hooks: []

## Objective

Ship the **real Explore SPA** into `internal/httpapi/embeddist` via a scripted build/sync pipeline, and keep StaticDir resolution so a **consumer temp root without `web/`** serves that embedded SPA (not stub/placeholder). Follow [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md). **Do not** implement auto-port (S03). Law 19 — httpapi adapter only.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — **L1, L2** (do not reopen)
- [00-PLANNER.md](00-PLANNER.md)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — SoT for pipeline + StaticDir + T1–T3
- Live anchors (read before coding):
  - `internal/httpapi/static.go` — disk → embed → `placeholderHTML`
  - `internal/httpapi/embed.go` — `//go:embed all:embeddist`
  - `internal/httpapi/embeddist/index.html` — **stub today** (`Embedded GUI stub`)
  - `internal/httpapi/embeddist/README.md` — still teaches two-artifact everyday
  - `internal/httpapi/server.go` — default StaticDir `<root>/web/dist`; refuse == root
  - `internal/httpapi/httpapi_test.go` — `TestStaticDirRefuseProjectRoot`, `TestStaticCSPAndEmbedFallback` (today only asserts `"Trace"`)
  - `cmd/trace/local_http.go` — `--static-dir` + usage strings (minimal honesty OK)
  - `web/package.json` — `npm run build` → `tsc -b && vite build` → `web/dist/**`
  - `web/dist/index.html` — real SPA markers: `#root` + `/assets/` module script
  - `scripts/` — Trace-root scripts exist (e.g. `repro-stray-trace-db.sh`); **no** `embed-gui.sh` yet
  - **No** Trace-root `Makefile` / `.github/workflows` — do not invent phantom CI

## Session start

Follow agent-loop-protocol Session start. Unattended: execute PLAN locks below; clarify only if PLAN contradicts L1/L2 or live anchors.

## Locked defaults (FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Pipeline **(A)** | From Trace root: `cd web && npm ci && npm run build` → sync `web/dist/**` into `internal/httpapi/embeddist/` → `go build` (or test) so `//go:embed` sees real assets |
| Canonical entrypoint | **`scripts/embed-gui.sh`** (create). Script must `cd` to Trace repo root (derive from `$0`), fail on missing Node/`npm`, and be idempotent |
| `//go:generate` | On or near `embed.go`, invoke the script with a path that works when generate runs from `internal/httpapi/` (e.g. `//go:generate ../../scripts/embed-gui.sh` **or** script self-roots — pick one and document in script header) |
| Optional Makefile | Thin Trace-root `Makefile` target `embed-gui` calling the same script — **optional**; OK to add; do **not** invent CI workflows |
| Sync rules | Replace embeddist **SPA assets** with `web/dist/**`. **Keep/regenerate** short `embeddist/README.md` (pipeline steps + last-resort stub note — **not** “two-artifact everyday”). Do not leave orphaned stub `index.html` after a successful sync |
| StaticDir order | **Unchanged:** `diskIndexExists(StaticDir)` → `embeddedIndexExists()` → inline `placeholderHTML`. Default empty flag → `filepath.Join(absRoot, "web", "dist")`. **No** Trace-module / GOPATH probe. **No** copy SPA into consumer `.trace/` |
| Consumer without `web/` | Missing disk `index.html` → **embedded real SPA** (post-sync) |
| Trace-checkout DX | If disk `<root>/web/dist/index.html` exists → **disk wins** (contributor Vite path without re-embed) |
| Explicit `--static-dir` | Abs that path; still **refuse** StaticDir == project root; missing index → embed → placeholder |
| Stub / placeholder | Last-resort only (empty/broken embed). Phrase `Embedded GUI stub` must **not** appear in shipped embed `index.html` after successful pipeline. Update `placeholderHTML` + stub tone so they do **not** teach consumer `web/` as primary |
| SPA vs stub markers | **Real:** `<div id="root">` (or `id="root"`) **and** `script type="module"` with `/assets/` src. **Stub:** body contains `Embedded GUI stub` and/or missing `#root` / missing `/assets/` module script |
| Auto-port | **Out of scope** — leave P32-PORT fail-if-busy / no hop until S03 |
| Docs | Rewrite `embeddist/README.md` here. Minimal help/usage honesty in `local_http.go` / `help.go` OK if they still teach consumer `web/` as primary. **S04** owns full quickstart / help flip |
| Law 19 | No business-logic fork in `web/`; httpapi serves assets only |

### Answers locked for implementer (planner gate)

1. **Exact files/scripts:** create `scripts/embed-gui.sh`; add `//go:generate` near `embed.go`; optional root `Makefile` `embed-gui`; sync real `web/dist` into `embeddist/*` + rewrite `embeddist/README.md`; tone `placeholderHTML` in `static.go`; minimal usage strings if needed; tests in `httpapi_test.go` (and helpers as needed).
2. **Resolution order after change:** **unchanged** disk-if-`index.html` → embed → placeholder; default path string stays `<root>/web/dist`.
3. **Tests:** PLAN **T1** (required), **T2**, **T3** (keep/extend existing refuse-root). Seed **T10** honesty: after pipeline, embedded `index.html` is not stub-only (can be asserted in T1 or a focused embed-content check).
4. **Stub vs real detection in tests:** Assert **not** `Embedded GUI stub`; assert body contains `id="root"` (or `#root` equivalent in HTML) **and** `/assets/` with `type="module"` (match live `web/dist/index.html` shape).

## Role work

### Minimal todos (execute in order)

1. **Script** — Add `scripts/embed-gui.sh`:
   - Resolve Trace root; `cd web && npm ci && npm run build`.
   - Sync `web/dist` → `internal/httpapi/embeddist` without deleting the README contract (e.g. sync assets then write README, or backup/restore README).
   - Exit non-zero on failure; short usage comment at top.
2. **Generate hook** — Add `//go:generate` on/near `embed.go` pointing at that script.
3. **Optional** — Root `Makefile` with `embed-gui:` → same script (only if cheap; not required for exit).
4. **Populate embed** — Run the script once so `embeddist/` contains real SPA (`index.html` + `assets/**`, etc.) and rewritten README. Commit those artifacts with the change (tests compile against package embed).
5. **Tone** — Update `placeholderHTML` (and any remaining stub HTML only if you keep a last-resort file — prefer real SPA in embeddist; stub is “pipeline not run” failure mode documented in README, not the default checked-in index after S02).
6. **Minimal CLI honesty** (optional, small) — If `local_http.go` usage still says “Prefer building web/dist; embedded stub is a fallback,” adjust to: consumers get embedded SPA; disk `web/dist` is contributor/override. Do **not** rewrite `docs/gui-quickstart.md` (S04).
7. **Tests (TDD-friendly):**
   - **T1** — Temp root + store open (`.trace/`), **no** `web/`; `httpapi.New` default StaticDir; `GET /` → 200; body has real SPA markers; **not** `Embedded GUI stub`. (Update or replace `TestStaticCSPAndEmbedFallback` assertions — keep CSP checks.)
   - **T2** — Temp root with planted `web/dist/index.html` containing a unique marker string → `GET /` serves **disk** marker (not embed SPA).
   - **T3** — Keep `TestStaticDirRefuseProjectRoot` (or equivalent) green.
8. **Regression guard** — Do **not** change listen/bind/addr-in-use / auto-port behavior. Focused test run only.
9. **Board** — Mark own row done with Notes: script path, key files, test names + pass evidence.

### Out of this row

- Auto free-port / `flag.Changed` / serve post-bind print (S03)
- Full docs/quickstart flip (S04)
- VERIFY / DR-HANDOFF (S05)
- Inventing Trace-root CI workflows
- Explore UI redesign

## Exit criteria

- [ ] `scripts/embed-gui.sh` exists and successfully populates `embeddist` with real SPA (README kept/rewritten)
- [ ] `//go:generate` wired near `embed.go`
- [ ] Consumer-like root without `web/` serves real SPA in automated test (**T1**)
- [ ] Disk-wins test (**T2**) and StaticDir == root refused (**T3**)
- [ ] Stub phrase absent from embedded `index.html`; placeholder/README do not teach consumer `web/` as primary
- [ ] No auto-port behavior changes
- [ ] Board Notes cite tests + key files
- [ ] Next **P34-S02-02**

## Todo updates

Status + notes on **P34-S02-01** only.

## Next

`P34-S02-02`
