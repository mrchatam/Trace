# P34-S00-01 — Research

## Metadata
- id: P34-S00-01
- todo_ids: [P34-S00-01]
- role: implementer
- skills: [research, planning-and-task-breakdown, diagnosing-bugs]
- mcps: []
- agents: []
- verification: automated
- hooks: []

## Objective

Author **only** `scopes/scope-00-research/RESEARCH.md`: embed vs sidecar, StaticDir policy, auto-port algorithm under L3, consumer layout audit, peer cites. **No product code.** Do not edit DESIGN-LOCKS, INTAKE, CLI, or `web/` (except reading). Do **not** start S01.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — L1–L4 + P34-00 clarifications
- [INTAKE.md](../../INTAKE.md)
- [Phase 34 README](../../README.md)
- [00-PLANNER.md](00-PLANNER.md) — locked defaults + must-answer set
- [SCOPE-TODOS.md](SCOPE-TODOS.md) — planner-locked live facts
- **L3 supersession context (read, do not reopen):**
  - [P33 S00 RESEARCH](../../../phase-33-gui-craft-hook-launch/scopes/scope-00-research/RESEARCH.md) — **rejected** UA auto-port (“conflicts with P32-PORT / multi-project `--addr` story”)
  - [P32 S00 RESEARCH](../../../phase-32-graph-first-gui/scopes/scope-00-research/RESEARCH.md) — deferred UA auto-increment as optional #2; shipped fail + `--addr`
- Live Trace (read-only):
  - `internal/httpapi/static.go` — disk then embed then placeholder
  - `internal/httpapi/embed.go` — `//go:embed all:embeddist`
  - `internal/httpapi/embeddist/index.html` + `README.md` — stub + optional `cp` recipe + “two-artifact everyday”
  - `internal/httpapi/server.go` — default StaticDir `<root>/web/dist`; refuse StaticDir == root
  - `internal/httpapi/addr_in_use.go` — `IsAddrInUse` + `FormatAddrInUseMessage` (manual `--addr` hint)
  - `cmd/trace/local_http.go` — shared `gui`/`serve`; usageGUI says “does not auto-pick a free port”
  - `cmd/trace/help.go` — “Port conflict: pick a distinct `--addr` — no auto free-port.”
  - `docs/gui-quickstart.md` — disk wins / optional embeddist copy / two-artifact
- Peer (present): `similar projects/Understand-Anything/understand-anything-plugin/packages/viewer/bin/viewer.mjs`
  - Default port `5173`; `listen(attemptPort+1)` on `EADDRINUSE` while `!portExplicit` and attempts left (**10**); bind `127.0.0.1`; `--port` explicit → no hop; `--port 0` = OS free port; best-effort open + `--no-open`

## Session start

Follow agent-loop-protocol Session start. Unattended: prefer DESIGN-LOCKS/INTAKE; do not reopen L1–L4. Proceed without waiting for plan confirmation.

## Locked defaults (planner FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Artifact | `docs/phases/phase-34-gui-packaging-multiproject/scopes/scope-00-research/RESEARCH.md` **only** |
| Product / CLI / web / OpenAPI edits | **Forbidden** |
| L1 | Consumer = `.trace/` only — no required `web/` |
| L2 | Real SPA from Trace binary; disk `web/dist` = Trace-checkout DX; stub last-resort; **reject** SPA copy under consumer `.trace/` as primary |
| L2 lean | Prefer **`go:embed`** of real SPA; install-sidecar only if RESEARCH proves embed unblockable — still must not require consumer `web/` |
| L3 | Default busy → auto free **loopback** port; print + open chosen URL; `--addr` **strict** fail-if-busy |
| Supersession | P33/P32 **rejected/deferred** UA auto-port; **L3 overturns** for **default bind only** — document explicitly in RESEARCH |
| Shared bind | Prefer same auto-port path for `gui` **and** `serve` unless RESEARCH justifies split |
| Sequence | This row → `P34-S00-02`; S01 waits for S00-02 PASS |

## Planner-verified live facts (2026-08-21 — confirm; note drift only)

| Area | Fact |
|------|------|
| Static order | `diskIndexExists(StaticDir)` → `embeddedIndexExists()` → inline `placeholderHTML` (`static.go`) |
| Default StaticDir | Empty flag → `<projectRoot>/web/dist` (`server.go` / `local_http.go` flag help) |
| Embed | Stub only — copy mentions “Embedded GUI stub” / build `web/dist` (`embeddist/index.html`) |
| Embed recipe | Manual `npm run build` + `rm`/`cp -a dist → embeddist` then `go build` (`embeddist/README.md`); **no** Makefile/CI embed target in Trace root today |
| Port conflict | Fail + stderr `FormatAddrInUseMessage` with example `--addr 127.0.0.1:7433` |
| Default bind | `127.0.0.1:7432` (`httpapi.DefaultAddr`) |
| Help debt | `help.go` + `usageGUI`/`usageServe` teach no auto free-port + prefer disk `web/dist` |

## Research rejects (must appear in RESEARCH)

1. Requiring consumer `web/` or documenting two-artifact as consumer primary.
2. Copying SPA into consumer `.trace/` as primary asset path.
3. Keeping “no auto-port” for **default** bind (conflicts with L3).
4. Auto-hopping when user set explicit `--addr`.
5. Public / `0.0.0.0` default; always-on daemon; hosted SaaS.

## Must answer (all required in RESEARCH.md)

1. **Embed options:** Rank at least: (A) make/CI/`go generate` pipeline → populate `embeddist` + `go:embed`; (B) install-sidecar beside binary; (C) keep manual `cp` only. Recommend **one** for S01/S02 (L2 lean = A unless blocked).
2. **StaticDir policy:** Consumer default = embed (or empty StaticDir that skips disk). When does disk win — Trace checkout only, any `-C` with `web/dist`, and/or `--static-dir` only? Propose a detection rule if Trace-checkout DX needs “disk wins when present.”
3. **Auto-port:** Compare and recommend one of:
   - **UA-style:** start at default port (`7432`), increment by 1 on `EADDRINUSE`, max N attempts, loopback host fixed
   - **Scan range:** try `[7432, 7432+R)` then fail
   - **`:0` then advertise:** OS-assigned free port (note: may skip preferred default when free — usually worse UX)
   - Peer nuance: UA skips hop when `--port` explicit; Trace should map that to **default bind auto** vs **`--addr` strict**
4. **`--addr`:** Confirm L3 split — explicit pin = fail-if-busy (friendly message may still *mention* default auto behavior). How to detect “user set `--addr`” vs default (`flag.Changed` vs string-equal DefaultAddr — recommend one).
5. **Docs/layout audit:** List concrete paths that violate L1/L2 teaching today (seed list below); assign fix owner S02 vs S04.

### Docs audit seed list (expand / verify)

| Path | Suspected issue |
|------|-----------------|
| `docs/gui-quickstart.md` | Disk wins; optional embeddist; two-artifact everyday |
| `internal/httpapi/embeddist/README.md` | “Supported everyday path remains two-artifact” |
| `cmd/trace/help.go` | Serves `web/dist`; “no auto free-port” |
| `cmd/trace/local_http.go` | usage strings prefer `web/dist`; no auto-pick |
| `internal/httpapi/addr_in_use.go` | Manual `--addr` only (expected until S03) |
| Phase 34 README “Live repo baseline” | Snapshot — not consumer doc; optional cite |
| Root `README.md` / `AGENTS.md` / `web/README.md` | Skim for consumer `web/` teaching |

## Preflight (before writing RESEARCH.md)

- [ ] Confirm static.go order: disk → embed → placeholder
- [ ] Confirm embeddist is stub (not full Vite SPA assets)
- [ ] Confirm no root Makefile/CI embed sync target (or document if found)
- [ ] Confirm gui/serve help says no auto free-port
- [ ] Confirm FormatAddrInUseMessage hints manual `--addr`
- [ ] Skim gui-quickstart disk/embed wording
- [ ] Read UA `viewer.mjs` `listen` (~L369–391): increment, attempts=10, `portExplicit`
- [ ] Quote P33 RESEARCH reject line for L3 supersession section

## Role work

1. Run preflight; note drift from planner-locked facts as facts only.
2. Write `RESEARCH.md` using the **required template** below.
3. Update board **P34-S00-01** status + Notes (artifact path + lean bullets: embed / StaticDir / auto-port).

### RESEARCH.md template (required headings)

```markdown
# Phase 34 S00 — RESEARCH

## Summary
(3–6 sentences: consumer stub cause; embed lean; auto-port lean; docs debt; L3 vs P33.)

## Live baseline (verified)
| Area | Fact | Evidence path |
|------|------|---------------|
| Static order | … | internal/httpapi/static.go |
| Default StaticDir | … | server.go / local_http.go |
| Embed contents | stub / full | embeddist/ |
| Embed pipeline today | … | embeddist/README.md (+ Makefile/CI if any) |
| Port on conflict | fail + hint | addr_in_use.go |
| gui/serve help | … | help.go / local_http.go |
| Docs story | … | docs/gui-quickstart.md |

## L3 supersession (P32/P33 → P34)
- What P32/P33 said (cite RESEARCH paths + short quote)
- What L3 overturns (default bind only)
- What stays (explicit `--addr` strict; loopback; Law 19)

## Embed options
| Option | Pros | Cons | Rank / Recommend? |
|--------|------|------|-------------------|
| go:embed + release/make/CI pipeline | … | … | … |
| Install-sidecar beside binary | … | … | … |
| Manual cp only (status quo) | … | … | … |
| Other (if any) | … | … | … |

Recommendation for S01/S02: …

## StaticDir policy recommendation
- Consumer default: …
- Trace-checkout DX: … (detection rule if any)
- `--static-dir` explicit: …
- Reject: requiring consumer web/; SPA under consumer `.trace/`

## Auto-port algorithm
- Candidates compared: UA-increment / scan range / `:0` …
- **Recommend:** …
- Default bind behavior (start port, max attempts, host constraint): …
- Explicit `--addr` detection + fail-if-busy: …
- Print + open URL (gui vs serve): …
- Shared `gui`+`serve` vs split: …
- Peer borrow/reject notes (`viewer.mjs`): …

## Consumer layout / docs audit
| Path or doc | Issue vs L1/L2/L3 | Fix owner (S02/S03/S04) |
|-------------|-------------------|-------------------------|
| … | … | … |

## Rejected alternatives (short)
Must include research rejects 1–5.

## Handoff to S01 / S02 / S03
- S01: PLAN locks to take from this RESEARCH (pipeline cmds, StaticDir table, auto-port params, test matrix seeds)
- S02: embed + StaticDir
- S03: auto-port
```

## Exit criteria

- [ ] `RESEARCH.md` exists with **all** template headings non-empty
- [ ] Embed recommendation + StaticDir policy + auto-port algorithm present
- [ ] L3 supersession of P33/P32 no-auto-port documented with cites
- [ ] Research rejects 1–5 covered
- [ ] Board Notes cite artifact + leans
- [ ] No product code / no DESIGN-LOCKS edits / no sibling prompt rewrites beyond Notes

## Minimal todos

- [ ] Preflight live static/embed/port/docs + UA listen + P33 cite
- [ ] Author `RESEARCH.md` from template
- [ ] Self-check must-answer 1–5 + rejects
- [ ] Update board row **P34-S00-01** Notes

## Todo updates

Status + notes on **P34-S00-01** only.

## Next

`P34-S00-02`
