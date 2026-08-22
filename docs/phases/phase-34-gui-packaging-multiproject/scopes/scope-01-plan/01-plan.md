# P34-S01-01 — Plan

## Metadata
- id: P34-S01-01
- todo_ids: [P34-S01-01]
- role: implementer
- skills: [planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Author **only** `scopes/scope-01-plan/PLAN.md` from S00 RESEARCH (PASS) + DESIGN-LOCKS L1–L4 + live repo. Lock embed/release pipeline, StaticDir policy, auto-port + flags, test matrix, docs touch list, and S02→S05 handoffs. **No product code.** Do **not** implement embed or auto-port here.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [00-PLANNER.md](00-PLANNER.md)
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md) — **PASS**; do not re-debate leans
- Live anchors (read-only): `internal/httpapi/static.go`, `embed.go`, `embeddist/`, `addr_in_use.go`, `bind.go`, `server.go`, `cmd/trace/local_http.go`, `docs/gui-quickstart.md`, `web/package.json`, `web/dist/index.html` (SPA marker reference)

## Session start

Follow agent-loop-protocol Session start. Unattended: bind to RESEARCH + L1–L4; do not reopen locks. Prefer the locked defaults below over inventing alternatives.

## Locked defaults (FINAL — copy into PLAN; do not reopen)

| Item | Value |
|------|-------|
| Artifact | `docs/phases/phase-34-gui-packaging-multiproject/scopes/scope-01-plan/PLAN.md` **only** |
| Product edits | **Forbidden** (no Go/`web/`/Makefile/CI files in this row) |
| Board order | **S02 embed → S03 auto-port → S04 docs → S05 VERIFY** — PLAN must not reverse |
| L1 | Consumer footprint = `.trace/` only; no required consumer `web/` |
| L2 | Real Explore SPA in release `go:embed`; disk `<root>/web/dist` = contributor DX only |
| L3 | Default-bind auto free-port; explicit `--addr` strict fail-if-busy |
| L4 | One process = one project root; multi-project = N processes × N ports |
| Reject | SPA under consumer `.trace/`; public/`0.0.0.0` default; always-on daemon; install-sidecar **(B)** unless PLAN proves Node unblockable (not evidenced); manual-`cp`-only **(C)** as primary; Trace-module path probe for StaticDir |

### Embed pipeline (RESEARCH lean **A** — lock in PLAN)

| Step | Command / rule |
|------|----------------|
| Build SPA | From Trace repo root: `cd web && npm ci && npm run build` (`web/package.json` → `tsc -b && vite build`) |
| Sync into embed | Replace contents of `internal/httpapi/embeddist/` with `web/dist/**` **except** keep/regenerate a short `embeddist/README.md` (pipeline + last-resort stub note — **not** “two-artifact everyday”) |
| Go build | `go build -o bin/trace ./cmd/trace` (or release equivalent) so `//go:embed all:embeddist` picks up real assets |
| Automation name | PLAN must name **one primary** Trace-root entrypoint: prefer **`scripts/embed-gui.sh`** + `//go:generate` on/near `internal/httpapi/embed.go` invoking it; optional thin **`Makefile` target `embed-gui`** wrapping the same. Document both contributor DX and future CI/release: run embed sync **before** `go build` when shipping full SPA |
| Today’s gap | No Trace-root `Makefile`; no `.github/workflows` — PLAN may **schedule creating** script/`go:generate`/(optional) Makefile in **S02**; do not invent a phantom existing CI job |
| VERIFY stub-fail | S05 **fails** if shipped/embedded UI is still stub while release intended full SPA. Stub markers: body text `Embedded GUI stub` and/or missing `#root` / `/assets/` module script. Real SPA markers (from live `web/dist/index.html`): `<div id="root">` + `script type="module"` under `/assets/` |
| Last-resort stub | Keep a stub path only when embed empty/dev mistake; placeholderHTML remains after embed miss — tone update in S02 |

### StaticDir resolution (lock table into PLAN)

| Context | Behavior |
|---------|----------|
| Empty `--static-dir` | Resolve candidate to `filepath.Join(absRoot, "web", "dist")` (keep today’s default path string) |
| Consumer project (no `web/`) | Candidate has no `index.html` → **embedded real SPA** (post-S02) → else placeholder |
| Trace checkout with `web/dist` | Opportunistic: if `diskIndexExists(StaticDir)` → **disk wins** (contributor DX). **No** Trace-module / GOPATH probe |
| Explicit `--static-dir DIR` | Abs that path as StaticDir; still **refuse** StaticDir == project root; missing `index.html` → embed → placeholder |
| Order (unchanged semantics) | disk-if-`index.html` → `embeddedIndexExists()` → inline `placeholderHTML` |
| Consumer story | Must **not** require `web/`; do **not** copy SPA into `.trace/` |

### Auto-port / bind (lock table into PLAN)

| Case | Behavior |
|------|----------|
| Commands | **Both** `trace gui` and `trace serve` — **shared** bind helper (prefer `httpapi` listen helper or shared path in `local_http.go`; PLAN picks one location, S03 implements) |
| Default host/port | Host `127.0.0.1` from `DefaultAddr`; start port **7432** |
| Default free | Bind `:7432`; print `http://127.0.0.1:7432/`; `gui` opens that URL unless `--no-open` |
| Default busy (`EADDRINUSE`) | UA-increment: try `7433` … up to **max 10** attempts (**ports `7432`–`7441`** inclusive); on success print **chosen** URL; `gui` opens **chosen** URL |
| Exhausted | Fail with updated message: auto range exhausted + suggest `--addr` |
| Explicit `--addr` | Detect via **`FlagSet` / `flag.Changed` on `addr`** — **not** string-equal to `DefaultAddr`. Busy → **strict fail** (friendly text may mention default auto). Includes `--addr 127.0.0.1:7432` |
| `serve` print timing | Move listen URL print to **post-bind** success (today pre-listen at `local_http.go` ~205 is debt for S03) |
| `gui` open | Keep `OnListening` / open after successful listen; URL must match **actual** bound addr after hop |
| Out of S03 | No public bind default; no OS `:0` as Trace default |

### Docs owners (RESEARCH audit — do not thin)

| Path | Primary owner |
|------|----------------|
| `docs/gui-quickstart.md` | **S04** |
| `cmd/trace/help.go` | **S03** (port) + **S04**/minimal **S02** (static) |
| `cmd/trace/local_http.go` usage strings | **S03** + **S02**/S04 |
| `internal/httpapi/addr_in_use.go` | **S03** |
| `internal/httpapi/embeddist/README.md` | **S02** rewrite + **S04** polish |
| `web/README.md` | **S04** (contributor-labeled) |
| Root `README.md` | **S04** if quickstart link implies consumer `web/` |
| `AGENTS.md` | **S04** / phase close optional |
| `placeholderHTML` / stub copy | **S02** tone; **S04** if user-visible |

## Must lock in PLAN.md (checklist for author)

1. Exact release/dev commands to populate `embeddist` (script/`go generate`/optional `make embed-gui` + `web` build + sync rules + VERIFY stub-fail criteria).
2. StaticDir resolution table (consumer / Trace-checkout / `--static-dir`) matching locked table above.
3. Auto-port: both commands; algorithm; max 10; `flag.Changed`; post-bind serve print; gui open chosen URL.
4. Test matrix rows covering S02/S03/S04/S05 (seed IDs below — expand wording OK, do not drop assertions).
5. Docs/help touch list for S04 (full audit table owners).
6. Handoff bullets S02→S05; out-of-scope list (sidecar B, `.trace/` SPA, public bind, Explore UI redesign).

## Seed test matrix (required rows in PLAN — may add more)

| ID | Scope | Assertion |
|----|-------|-----------|
| T1 | S02 | Temp consumer root with `.trace/` init, **no** `web/` → `GET /` (or static index) is **real SPA** (`#root` and/or `/assets/` module script); **not** stub phrase `Embedded GUI stub` |
| T2 | S02 | Root with planted `web/dist/index.html` marker → **disk wins** over embed |
| T3 | S02 | `StaticDir == project root` still refused |
| T4 | S03 | Default bind busy → next port (e.g. `:7433`) + printed URL matches bound addr |
| T5 | S03 | Two concurrent default `gui`/`serve` → distinct ports; `gui` open/hook sees chosen addr |
| T6 | S03 | Explicit `--addr` busy → **fail** (including `--addr` equal to DefaultAddr string when `Changed`) |
| T7 | S03 | Auto exhausted (10 busy ports) → fail mentioning range + `--addr` |
| T8 | S04 | Help/usage/quickstart no longer teach consumer two-artifact primary or “no auto free-port” for default |
| T9 | S05 | VERIFY floor: T1 + concurrent free-port (T4/T5) + docs consumer story (T8); evidence under `experiments/runs/` |

## PLAN.md template (required headings — all non-empty)

```markdown
# Phase 34 S01 — PLAN

## Summary
(2–4 sentences: L2 real embed + L3 default auto-port; consumer `.trace/` only.)

## Embed / release pipeline
- Primary command(s): …
- Steps: web build → sync embeddist → go build
- README / stub / VERIFY stub-fail rule
- Files S02 may create (script, go:generate, optional Makefile) — plan only

## StaticDir resolution policy
| Context | Behavior |
|---------|----------|
| Consumer project (no web/) | … |
| Trace checkout with web/dist | … |
| Explicit --static-dir | … |
| Refuse root | … |

## Auto-port / bind policy
| Case | Behavior |
|------|----------|
| Default addr free | … |
| Default addr busy | … |
| Explicit --addr busy | … |
| Exhausted | … |
| Commands covered | gui + serve (shared) |
| Detection of explicit addr | flag.Changed |
| serve print / gui open | post-bind; open chosen URL |

## Test matrix
| ID | Scope | Assertion |
|----|-------|-----------|
| T1 | … | … |
(include at least T1–T9 seeds)

## Docs / help touch list (S04)
(table or bullets from RESEARCH audit — owners S02/S03/S04)

## Out of scope / deferred
- install-sidecar (B) unless Node unblockable
- SPA copy under consumer `.trace/`
- public bind default; OS :0 as default
- Explore UI / craft redesign
- Reversing board order (auto-port before embed)

## Handoff
- S02: …
- S03: …
- S04: …
- S05 VERIFY floor: …
```

## Minimal todos

1. Read RESEARCH.md + DESIGN-LOCKS + live anchors listed above (confirm stub-only embeddist, no root Makefile/CI, DefaultAddr `127.0.0.1:7432`, serve pre-listen print).
2. Write `PLAN.md` filling **every** template heading with the locked defaults (no product code).
3. Ensure test matrix includes T1–T9 seeds; docs list matches RESEARCH audit owners.
4. Board Notes on **P34-S01-01** cite `…/PLAN.md` + 1–2 lock bullets.
5. Self-check exit criteria; do **not** start S02.

## Exit criteria

- [ ] `PLAN.md` exists at `scopes/scope-01-plan/PLAN.md` with all template headings non-empty
- [ ] Embed pipeline commands + StaticDir table + auto-port table + test matrix (T1–T9+) + docs touch list present
- [ ] Board order S02→S03 preserved; rejects (`.trace/` SPA, public default, sidecar-as-primary) explicit
- [ ] Board Notes cite `PLAN.md` + lock bullets
- [ ] No product code / no `Makefile`/`scripts/`/`embeddist` edits in this row

## Todo updates

Status + notes on **P34-S01-01** only.

## Next

`P34-S02-00` (after this row `done`)
