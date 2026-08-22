# DESIGN-LOCKS — Phase 33

**Human-promoted 2026-08-21** from post–P32 dissatisfaction. Light clarifications from **P33-00** (2026-08-21) against live repo — do not reopen Themes A–C.

| Lock | Value |
|------|-------|
| Predecessor | Phase 32 **closed** — do not rewrite P32 history; spawn forward |
| Theme A — Color / craft | Replace bland UI with a **deliberate color system** + craft pass. **Required skills on implement/review:** `impeccable` (esp. colorize/bolder/polish + craft-floor), `ui-ux-pro-max`, `frontend-design`; also use `design-taste-frontend` / `high-end-visual-design` / `emil-design-eng` where helpful. Avoid AI-default purple/cream/broadsheet clusters unless brief demands. |
| Theme B — Explore hook | **Explore (`/`)** is an **interactive project graph** (Graphify energy: force/pan/zoom/click) as the user hook. Inspector/panels remain. Laws **6–7**: no unbounded full dump — S00/S01 choose budgeted overview (clusters, caps, progressive expand). **Live gap:** index already routes to `Graph.tsx`, but UX is still **center-first neighborhood** (empty until pick center) — Phase 33 upgrades the **open experience** to a project overview hook, not a second SPA. |
| Theme C — Launch | **`trace gui`** from a repo (cwd/`-C`) starts local GUI and **opens the browser**. Install path: `trace` on PATH via documented `go install` / package / symlink — not `./bin/trace serve` as the primary story. Reuse `serve` + P32-PORT behavior. Keep `serve` for scripting; demote in user-facing docs (S05). |
| PATH vs `trace install` | Existing `trace install …` targets agents/MCP/hooks — **not** PATH placement of the binary. S00/S02 must not conflate the two. |
| Stack | Keep `web/` React+Vite+xyflow; Law 19 adapters only |
| Skills gate | S01, S03, S04 implement **and** review prompts list impeccable + ui-ux-pro-max + frontend-design in metadata |
| Out of scope | Hosted SaaS; always-on daemon; Three.js-first; silent delete of root `trace.db`; unbounded full-graph dump API |

## Success sketch

1. Open GUI → Explore immediately feels like a **colored, interactive graph** of the project (hook) — not an empty “pick center” gate.
2. From any Trace-initialized repo: `trace gui` → browser opens to that project’s GUI (PATH install assumed).
3. UI no longer reads as gray/bland — color encodes kind/state and brand atmosphere.
