# INTAKE — Phase 33 GUI satisfaction gaps

**Human 2026-08-21** after Phase 32 close. Do not reopen Phase 32 rows.

## Complaints

1. **Bland / no coloring** — visual craft still too muted; needs a strong palette and design pass.
2. **Explore is not a hook** — landing (“Explore”) must show an **interactive project graph** like Graphify (pan/zoom/click), not a thin ops overview.
3. **Launch UX** — telling users to run `./bin/trace serve` is impractical. Want **installed CLI** so from a repo terminal: something like `trace gui` opens the GUI for **that** cwd/project.

## Light peer launch notes (seed for S00 — not full RCA)

| Peer | How GUI opens | Takeaway |
|------|---------------|----------|
| **Understand-Anything** | In-agent `/understand-dashboard`; standalone `npx …viewer.tgz <project>` prints tokenized URL and **opens browser**; Vite `dev:dashboard` from clone | One command from project path; auto-open browser; PATH/`npx` not `./bin/…` |
| **Graphify** | Writes `graph.html` (+ json); user opens file (`Start-Process` / browser) | Interactive **whole-graph** HTML is the hook; zero long-running daemon for viz |
| **Trace today** | `trace serve` (default `:7432`); docs push `./bin/trace serve`; no `gui` subcommand; no auto-open browser | Friction; multi-project needs `--addr` (P32-PORT may have improved errors) |

## Recommended product direction (locks in DESIGN-LOCKS)

- Design: **impeccable** (`colorize` / `bolder` / craft floor) + **ui-ux-pro-max** + **frontend-design** (+ design-taste / high-end as needed) — opinionated color system, not gray monochrome.
- Explore: Graphify-class **interactive canvas as home**; still honor Laws 6–7 (budgeted / progressive “project overview” — not unbounded dump).
- Launch: ship **`trace gui`** (subcommand; cwd = project) = serve + open browser; ensure **install puts `trace` on PATH** (go install / symlink docs). Flag form `trace -gui` is secondary if needed.
