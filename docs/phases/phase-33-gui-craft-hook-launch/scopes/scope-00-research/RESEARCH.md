# Phase 33 S00 — RESEARCH

## Summary

Trace’s launch story still centers on `./bin/trace serve` with no browser open and no `gui` subcommand — friction versus peers that one-command from a project path. Explore already *is* the Graph screen (`/`), but the UX is still an empty “pick center” gate (`DEFAULT_MAX=50` / `UI_CAP=100`), not a Graphify-like project overview hook; Nav **Overview** (`/overview`) is a separate ops screen and must not be conflated. Under Laws 6–7, the overview lean is **client composition of budgeted seeds** (search + tasks + project → parallel `getGraph` neighborhoods, progressive expand) — never an unbounded full-graph dump. PATH lean: document `go install github.com/mrchatam/Trace/cmd/trace@…` as primary; keep binary PATH distinct from `trace install` (agents/MCP/hooks). For S02: ship **`trace gui`** as primary (reuse `serve` + P32-PORT + best-effort open browser).

## Live baseline (verified)

| Area | Fact | Evidence path |
|------|------|---------------|
| Explore home | Index route renders `Graph`; `/graph` redirects to `/` | `web/src/App.tsx` |
| Explore ≠ Overview | Nav label **Explore** → `/`; separate **Overview** → `/overview` (ops screen) | `web/src/layout/Nav.tsx` |
| Empty gate | `DEFAULT_MAX=50`, `UI_CAP=100`; empty until center selected (“Pick center” / EmptyState) | `web/src/screens/Graph.tsx` |
| Graph API | `/v1/graph` requires `center` + `max_nodes`; summary “never full dump”; `/v1/search` capped; client `getGraph` / `search` / `listTasks` / `getProject` | `api/openapi.yaml`, `web/src/api/ops.ts` |
| CLI launch | `root.go` has `serve` + `install`; **no** `gui`; default `127.0.0.1:7432`; friendly in-use + pick `--addr` (no auto-port) | `cmd/trace/root.go`, `cmd/trace/serve.go` |
| Docs story | Quickstart leads with `go build -o bin/trace` then `./bin/trace serve` | `docs/gui-quickstart.md` |
| Agents install | `trace install detect\|uninstall\|agents\|cursor\|claude\|cursor-hook\|git-hook` — not PATH for binary | `cmd/trace/install.go` |

No drift from planner-locked table (SCOPE-TODOS / 01-research locked defaults).

## Peer launch matrix

| Dimension | Graphify | Understand-Anything | Trace today | Borrow / reject |
|-----------|----------|---------------------|-------------|-----------------|
| How GUI opens | After extract: open static `graph.html` (`Start-Process .\graphify-out\graph.html`) | One-shot `npx …viewer.tgz <project>` prints tokenized URL | `trace serve`; docs push `./bin/trace serve` | **Borrow** one-command-from-cwd UX (UA) + immediate interactive canvas feel (Graphify). **Reject** docs-as-`./bin/…` primary; prefer `trace gui` on PATH |
| Auto-open browser | Manual file open (OS shell) | Best-effort `open` / `start` / `xdg-open`; `--no-open` opt-out (`viewer.mjs`) | No | **Borrow** UA best-effort open + `--no-open`-style opt-out for S02. Listen success ≠ browser success |
| Port on conflict | n/a (file viz) | Auto-increments if taken (`listen(attemptPort+1)`, up to 10) unless `--port` explicit | P32 friendly in-use stderr + user picks `--addr` | **Reject** UA auto-port for Trace. **Borrow** Trace/P32 explicit `--addr` multi-project story |
| PATH / install | `uv tool install graphifyy` → `graphify` on PATH | `npx` release tarball (no global binary required) | Build to `bin/`; `install` = agents/MCP/hooks only | **Borrow** “binary/tool on PATH” teaching. **Reject** conflating PATH with `trace install …` |
| Viz hook | Interactive force `graph.html` (whole-graph artifact; communities in report / `cluster.py`) | Interactive dashboard over existing KG | Explore = center-first Graph; empty until center | **Borrow** Graphify pan/zoom/click energy as Explore hook. **Reject** Graphify-style unbounded whole-graph HTML/API as Trace default (Laws 6–7) |
| Daemon | None for viz (static file) | Local viewer HTTP on loopback + token | Opt-in HTTP `serve` (Law 19 adapter) | Keep opt-in local serve; **reject** always-on daemon / hosted SaaS |

**Peer cites:** Graphify open-file — `similar projects/graphify/worked/rsl-siege-manager/README.md` (+ artifact `…/graph.html`); community energy — `…/mixed-corpus/raw/cluster.py`, review notes `…/rsl-siege-manager/review.md`. UA — `similar projects/Understand-Anything/understand-anything-plugin/packages/viewer/README.md`; open-browser + auto-port — `…/packages/viewer/bin/viewer.mjs`.

## Explore overview under Laws 6–7

### Options considered

| Option | Idea | Law 6–7 risk |
|--------|------|--------------|
| **(A) Cluster / community nodes with caps** | Graphify-like community hubs (`cluster.py` Leiden energy; review notes communities as headline structure) as overview nodes, expand into members | Hook feel high; needs cluster summary data. Risk: cluster payload that reintroduces near-full dump if uncapped — must cap communities and members |
| **(B) Seeded multi-center neighborhoods merged client-side** | N centers → N `getGraph` calls → merge/dedupe in UI under a total node budget | Reuses budgeted API; risk is N×max_nodes blowup — need seed count × per-seed `max_nodes` ≤ UI cap |
| **(C) Progressive expand from overview seeds** | Show seeds first; expand neighborhood on click / double-click / “expand” (existing re-center pattern) | Strong Law 7 fit; risk if “expand all” is default — keep expand explicit |
| **(D) Search / task / plan seeds → then getGraph** | Derive centers from `getProject` + `listTasks` (active) + `/v1/search` (capabilities / entities) then budgeted graph | Reuse-only; risk if search limits are huge or seeds unbounded — cap seed list |
| **(E) Seed-export / status as graph body** | Treat `/v1/seed/export` or status as the Explore canvas | **Reject** — OpenAPI already marks seed status as *not* a full graph body; export is portable graph for clone, not progressive Explore |

### Recommendation for S01

**Primary: (D) + (B) + (C)** — on Explore open, derive a small seed set from existing ops (`getProject`, `listTasks`, `search`), fetch **budgeted** neighborhoods via `getGraph` for those seeds, merge client-side, then **progressive expand** on user action (no empty “pick center” gate as the first paint).

**Budget leans (S01 may refine):**

| Lean | Value |
|------|-------|
| Seed count | **4–8** centers (hard cap ≤ 8) |
| Per-seed `max_nodes` | **30–50** (≤ current `DEFAULT_MAX`) |
| Merged UI cap | **≤ 100–120** visible nodes (honor today’s `UI_CAP=100` unless S01 argues a small bump ≤ 150) |
| Depth | default **2** (OpenAPI default) |
| Expand | user-driven re-center / expand only; no “load all” |

**(A) clusters:** inspiration only for S01 IA (group by kind/state visually); do **not** require a Leiden/community API in Phase 33. Defer real cluster endpoints to a later thin gap if needed.

Color note for S01: **tokens + Explore overview IA only**; **S04 owns full shell colorize/craft**.

### API implication

**`reuse`** (preferred) — compose `search` + `listTasks` + `getProject` + `getGraph` in the Explore client. OpenAPI already forbids unbounded graph (`center` + `max_nodes` required; “never full dump”).

**`reuse_then_gap_later`** only if S01/S03 prove seed quality is inadequate without a capped “overview seeds” helper — any gap must stay budgeted (seed count + `max_nodes`), Law 19 adapter-only, **never** unbounded dump or seed-export-as-graph-body.

**Explicit reject:** unbounded full-graph dump as default; Graphify whole-`graph.html` dump semantics as Trace API/GUI default; treating `/v1/seed/export` as Explore payload.

## PATH / install options

| Option | Pros | Cons | Rank / Recommend? |
|--------|------|------|-------------------|
| `go install github.com/mrchatam/Trace/cmd/trace@…` | Puts `trace` on user PATH; matches Go module (`go.mod`); one documented line for users with Go | Requires Go toolchain; version pin discipline (`@latest` vs tag) | **#1 — primary user story** |
| make / symlink from `bin/trace` | Fast for contributors already building; explicit local binary | Not “installed product”; easy to regress to `./bin/…` docs | **#2 — contributor / DIY** |
| package (deb/brew/…) | Best non-Go UX | Out of Phase 33 core; packaging + release ops | **#3 — later / optional** |

**Teaching:** Primary docs (S05) teach `go install …` then `trace gui` from a Trace-initialized repo. Contributor path: `go build -o bin/trace` + optional symlink — not the headline. **PATH ≠ `trace install`:** `trace install …` remains agents/MCP/hooks only (`install.go`); never document it as “install Trace on PATH.”

## CLI shape for S02

- **Primary:** `trace gui` (project from cwd / `-C` / `--root` parity with serve)
- **Inherit from serve:** `--addr`, loopback default `127.0.0.1:7432`, P32-PORT friendly in-use messaging (no silent port hop), token / `--allow-remote` posture unchanged
- **Open-browser:** best-effort after listen success; on failure print stderr tip with URL; **exit 0 if listening** (browser failure ≠ serve failure); support opt-out flag akin to UA `--no-open`
- **Keep `serve`:** scripting / CI / headless; demote in user-facing docs (S05 owns full flip from `gui-quickstart.md`)
- **`-gui` flag:** secondary only; prefer subcommand `trace gui` (DESIGN-LOCKS Theme C)

## Rejected alternatives (short)

- **Graphify-style static / unbounded whole-graph HTML as Trace default** — borrow interactive canvas energy only; reject dump-as-default (Laws 6–7).
- **UA auto-increment listen port** — conflicts with P32-PORT / multi-project `--addr` story; reject silent port hopping.
- **PATH via `trace install …`** — agents/MCP/hooks only; binary PATH is `go install` / symlink / package.
- **Treating `/overview` as Explore** — Theme B target is `/` Graph hook upgrade, not the Overview ops route.
- **Three.js-first** — out of DESIGN-LOCKS; keep React + xyflow.
- **Always-on daemon** — remain opt-in local HTTP only.
- **Hosted SaaS** — out of Phase 33; cloud is future OpenAPI consumer only.

## Handoff to S01 / S02 / S03

- **S01:** Overview IA = seed-composed Explore on open (D+B+C); budgets above; Explore ≠ `/overview`; color = **tokens only** (S04 owns full shell). Skills: impeccable + ui-ux-pro-max + frontend-design.
- **S02:** Implement **`trace gui`** primary; reuse serve + P32-PORT; best-effort browser open; PATH docs lean **`go install github.com/mrchatam/Trace/cmd/trace@…`**; never overload `trace install` for PATH.
- **S03:** Data path = client composition of `getProject` / `listTasks` / `search` → `getGraph` (ops.ts); progressive expand via existing center+budget calls; no unbounded dump; no new SPA.
