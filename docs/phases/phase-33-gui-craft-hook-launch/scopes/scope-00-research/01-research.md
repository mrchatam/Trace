# P33-S00-01 — Research

## Metadata
- id: P33-S00-01
- todo_ids: [P33-S00-01]
- role: implementer
- skills: [research, planning-and-task-breakdown]
- mcps: []
- agents: []
- verification: automated
- hooks: []

## Objective

Author **only** `scopes/scope-00-research/RESEARCH.md`: peer launch matrix (Graphify `graph.html`, Understand-Anything open-browser/`npx` viewer), **Laws 6–7–safe** project-overview-graph options for Explore, PATH/`trace gui` launch recommendations for S01–S02. **No product code.** Do not edit DESIGN-LOCKS, INTAKE, OpenAPI, CLI, or `web/`.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [docs/init/G_PROJECT_LAWS.md](../../../../init/G_PROJECT_LAWS.md) — Law **6** (no full-project graph dumps by default), Law **7** (progressive expand), Law **19** (adapters)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — Themes A–C (do not reopen)
- [INTAKE.md](../../INTAKE.md) — complaints + seed peer table
- [Phase 33 README](../../README.md)
- [00-PLANNER.md](00-PLANNER.md) — locked defaults + must-answer set
- Live Trace:
  - `web/src/App.tsx` — index → `Graph`; `/graph` → `/`; **also** `/overview` → ops `Overview` screen (not Theme B hook)
  - `web/src/layout/Nav.tsx` — label **Explore** (`/`); separate nav **Overview** (`/overview`)
  - `web/src/screens/Graph.tsx` — `DEFAULT_MAX=50`, `UI_CAP=100`; empty until center
  - `web/src/api/ops.ts` — `search`, `getGraph` (center + `max_nodes` required), `listTasks`, `getProject`, …
  - `api/openapi.yaml` — `/v1/graph` budgeted; `/v1/search`; seed export not full-graph body
  - `cmd/trace/root.go` — `serve` / `install` cases; **no** `gui`
  - `cmd/trace/serve.go` — default `127.0.0.1:7432`; P32-PORT friendly in-use (explicit `--addr`; no auto-pick port)
  - `cmd/trace/install.go` + help — agents/MCP/hooks only (`detect|uninstall|agents|cursor|claude|cursor-hook|git-hook`)
  - `docs/gui-quickstart.md` — primary story still `./bin/trace serve`
  - `go.mod` — module `github.com/mrchatam/Trace` (for `go install` ranking)
- Peers (verified present — cite these paths):
  - Graphify open-file hook: `similar projects/graphify/worked/rsl-siege-manager/README.md` (`Start-Process` / open `graph.html`); artifact `…/graph.html`
  - Graphify community/cluster energy (inspiration only): `similar projects/graphify/worked/mixed-corpus/raw/cluster.py`; review notes in `…/rsl-siege-manager/review.md`
  - UA viewer: `similar projects/Understand-Anything/understand-anything-plugin/packages/viewer/README.md` (`npx …viewer.tgz`, prints URL, opens browser)
  - UA open-browser impl: `…/packages/viewer/bin/viewer.mjs` (`--no-open`; platform `open`/`start`/`xdg-open`; **auto-increments port if taken** — reject for Trace)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: prefer DESIGN-LOCKS/INTAKE; **do not** reopen Themes A–C. Proceed without waiting for plan confirmation.

## Locked defaults (planner FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Artifact | `docs/phases/phase-33-gui-craft-hook-launch/scopes/scope-00-research/RESEARCH.md` **only** |
| Product / CLI / web / OpenAPI edits | **Forbidden** |
| Explore today | Index already **is** Graph (`App.tsx`); gap = **center-first empty gate**, not “missing Graph route” |
| Explore ≠ | Nav **Overview** (`/overview`) — ops screen; do **not** describe that as the Theme B hook |
| Hook upgrade | Research **project overview on open** under Laws 6–7 — not a second SPA |
| Laws | 6–7 budgeted/progressive; reject unbounded full dump as default; Law 19 if API discussed |
| Launch lean | Recommend **`trace gui`** primary; flag `-gui` secondary only; reuse `serve` + P32-PORT |
| Open browser | Best-effort; listen success ≠ browser success (S02 will implement) |
| PATH ≠ | `trace install …` (agents/MCP/hooks) — never conflate |
| Sequence | This row → `P33-S00-02`; S01 waits for S00-02 PASS |
| Peers missing? | Both folders **present** — cite paths above; if a file missing, fall back to INTAKE seed + peer README |

## Research rejects (must appear in RESEARCH “Rejected alternatives” or matrix Borrow/reject)

Document explicitly — do not soft-pedal:

1. **Graphify-style unbounded whole-graph HTML as Trace default** — energy to borrow (interactive canvas hook); dump-as-default to reject (Laws 6–7).
2. **UA auto-increment listen port** — conflicts with P32-PORT / Trace multi-project story (friendly in-use + user picks `--addr`). Borrow open-browser + one-command-from-cwd; reject silent port hopping.
3. **PATH via `trace install …`** — agents/MCP/hooks only; binary PATH is separate (`go install` / symlink / package).
4. **Treating `/overview` as Explore** — Theme B target is `/` Graph hook upgrade, not the Overview ops route.
5. **Three.js-first / always-on daemon / hosted SaaS** — out of DESIGN-LOCKS; short reject bullets only.

## Must answer (handoff — all required in RESEARCH.md)

1. **Peer launch:** What to **borrow** vs **reject** from Graphify (static whole-graph HTML), UA (`npx` viewer + open browser), Trace today (`trace serve`, docs `./bin/…`)?
2. **Overview under Laws 6–7:** Options (clusters / caps / progressive expand / seeded multi-center / client composition of search+graph+tasks, …) + **one** recommendation for S01.
3. **API:** Prefer reuse of existing `/v1` (`search` + `getGraph` + client composition)? Or thin OpenAPI gap? **Never** recommend unbounded full-graph dump default. Cap any proposed `max_nodes` / seed count.
4. **PATH:** Rank `go install` / make-or-symlink / package; how docs should teach; keep distinct from `trace install`.
5. **CLI for S02:** Confirm **`trace gui`** primary (cwd/`-C`); relationship to `serve`; open-browser failure mode lean.

## Preflight (before writing RESEARCH.md)

- [ ] Confirm `App.tsx` index → Graph; Nav “Explore”; Graph empty-until-center + budgets; note `/overview` is separate
- [ ] Confirm `root.go` has `serve`/`install`, no `gui`
- [ ] Confirm `gui-quickstart.md` leads with `./bin/trace serve`
- [ ] Skim Graphify: `worked/rsl-siege-manager/README.md` + `graph.html` (+ optional cluster.py / review.md)
- [ ] Skim UA: `packages/viewer/README.md` + `bin/viewer.mjs` (open-browser + port behavior)
- [ ] Skim OpenAPI `/v1/graph` + `/v1/search` descriptions (budget language)

## Role work

1. Run preflight; note any drift from locked table in RESEARCH (facts only — do not change locks).
2. Write `RESEARCH.md` using the **required template** below. Every heading non-empty; claims cite path or peer README (prefer the verified peer paths in References).
3. Keep Summary + Handoff short enough that S01/S02 planners can paste leans without re-reading peers.
4. Ensure **Research rejects** list above appear in matrix and/or Rejected alternatives.
5. Update board **P33-S00-01** status + Notes (artifact path + 1–2 lean bullets).

### RESEARCH.md template (required headings)

```markdown
# Phase 33 S00 — RESEARCH

## Summary
(3–6 sentences: launch friction; Explore hook gap; overview-budget lean; PATH lean.)

## Live baseline (verified)
| Area | Fact | Evidence path |
|------|------|---------------|
| Explore home | … | web/src/App.tsx |
| Explore ≠ Overview | … | web/src/layout/Nav.tsx (`/` vs `/overview`) |
| Empty gate | … | web/src/screens/Graph.tsx |
| Graph API | … | api/openapi.yaml / ops.ts |
| CLI launch | … | cmd/trace/root.go, serve.go |
| Docs story | … | docs/gui-quickstart.md |
| Agents install | … | cmd/trace/install.go / help |

## Peer launch matrix
| Dimension | Graphify | Understand-Anything | Trace today | Borrow / reject |
|-----------|----------|---------------------|-------------|-----------------|
| How GUI opens | … | … | `trace serve`; docs `./bin/…` | … |
| Auto-open browser | … | … | No | … |
| Port on conflict | n/a (file) | auto-increment (reject for Trace) | P32 friendly in-use + pick `--addr` | … |
| PATH / install | … | … | build to `bin/`; `install` = agents/MCP | … |
| Viz hook | graph.html (file / community aggregate) | interactive dashboard | Explore = center-first Graph | … |
| Daemon | none for viz | local viewer/Vite | opt-in HTTP serve | … |

Cite peer sources with paths from References (Graphify README/graph.html; UA viewer README + viewer.mjs).

## Explore overview under Laws 6–7
### Options considered
At least three of: (A) cluster/community nodes with caps; (B) seeded multi-center neighborhoods merged client-side; (C) progressive expand from overview seeds; (D) search/task/plan seeds → then getGraph; (E) other — with Law 6–7 risk notes.

### Recommendation for S01
One primary option + why (hook feel vs dump risk). Name budgets (seed count, max_nodes, depth) as leans — S01 may refine numbers.

### API implication
`reuse` | `thin_gap` | `reuse_then_gap_later` — justify. Explicit reject of unbounded dump / seed-export-as-graph-body.

## PATH / install options
| Option | Pros | Cons | Rank / Recommend? |
|--------|------|------|-------------------|
| `go install github.com/mrchatam/Trace/cmd/trace@…` | … | … | … |
| make / symlink from `bin/trace` | … | … | … |
| package (deb/brew/…) | … | … | … |

Document teaching: primary user story vs contributor build. **Must** state PATH ≠ `trace install` agents/MCP.

## CLI shape for S02
- Primary: `trace gui` (cwd / `-C`)
- Inherit from serve: `--addr`, loopback default, P32-PORT messaging, token/allow-remote posture
- Open-browser: best-effort; failure → stderr tip; exit 0 if listening
- Keep `serve` for scripting; demote in docs (S05 owns full flip)
- `-gui` flag: secondary only if argued; prefer subcommand

## Rejected alternatives (short)
Bullets must include: Graphify-style static full dump as Trace default; UA auto-port increment; overloading `trace install` for PATH; `/overview` as Explore hook; Three.js; always-on daemon.

## Handoff to S01 / S02 / S03
- S01: overview IA lean + budget numbers + color note (“S04 owns full shell; S01 tokens only”)
- S02: CLI + PATH ranked choice
- S03: data path lean (which ops / composition) — no implement detail
```

## Exit criteria

- [ ] `RESEARCH.md` exists at scope path with **all** template headings non-empty
- [ ] Peer matrix filled with borrow/reject column (includes port-on-conflict row)
- [ ] Research rejects 1–5 covered in matrix and/or Rejected alternatives
- [ ] One overview recommendation + explicit API lean (no unbounded dump)
- [ ] PATH ranked + PATH ≠ `trace install` stated
- [ ] Explicit **`trace gui`** primary recommendation for S02
- [ ] Board Notes cite `…/scope-00-research/RESEARCH.md` + lean one-liners
- [ ] No product code / no DESIGN-LOCKS edits / no sibling prompt rewrites

## Minimal todos

- [ ] Preflight live routes + serve/install + peers (paths in References)
- [ ] Author `RESEARCH.md` from template
- [ ] Self-check must-answer 1–5 + research rejects
- [ ] Update board row **P33-S00-01** Notes

## Todo updates

Status + notes on **P33-S00-01** only.

## Next

`P33-S00-02`
