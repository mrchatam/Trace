# P32-S05-01 — Polish implement

## Metadata
- id: P32-S05-01
- todo_ids: [P32-S05-01]
- role: implementer
- skills: [documentation-and-adrs, incremental-implementation, writing-guidelines]
- mcps: []
- verification: automated
- hooks: []

## Objective

Ship **production polish + docs** for Phase 32 explorer — **not** new explorer features.

1. **Docs (required):** Update [`docs/gui-quickstart.md`](../../../../gui-quickstart.md) with the **multi-project / port** pattern from **P32-PORT #1** (shipped in S02). Touch [`web/README.md`](../../../../../web/README.md) if DX copy still assumes a single `:7432` forever.
2. **Residuals (small):** Close or explicitly defer S03/S04 polish nits listed below — no IA/API/depth changes.
3. **Screenshots (optional but preferred):** If adding UI images to quickstart/README, show **Explore `/`** graph-home (canvas-first + inspector), not Overview CRUD as the hero.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md) — #1 shipped; #2 deferred; **#3/#4 → this scope**
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [00-PLANNER.md](00-PLANNER.md) — locks below are **final**
- S04 review Notes (board Order 563) — craft/screenshot callouts
- Live: `docs/gui-quickstart.md`, `web/README.md`, `internal/httpapi/addr_in_use.go` (`FormatAddrInUseMessage`), `cmd/trace/serve.go` / `help.go`, explorer `web/src/screens/Graph.tsx`, `components/Inspector.tsx`, `styles/tokens.css`, `styles/app.css`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Unattended: proceed without waiting for plan confirmation.

## Locked defaults (do not re-debate)

| Item | Value |
|------|-------|
| Scope | Polish + docs only — **no** new inspector sections, routes, API clients, or graph features |
| Port story | Document **S02 #1 as shipped**: fail-on-conflict; friendly in-use stderr; distinct `--addr` per project. **Do not** claim auto free-port / `:0` (#2 deferred) |
| Default bind | Still `127.0.0.1:7432` (`httpapi.DefaultAddr`); loopback-trust; remote still needs `--allow-remote` + token |
| One serve = one root | One `trace serve` = one `--root`/`-C`. Second project → second process + free port |
| Docs primary | **`docs/gui-quickstart.md`** must gain a clear multi-project / port section (or subsection) |
| Docs secondary | `web/README.md` — at least a one-liner or short note if it only shows `:7432` (Vite proxy stays on that bind unless user changes both) |
| OPEN-PORT | Close narrative for #3/#4 via docs (no mandatory helper script/flag this scope). Optional tiny helper only if already trivial — prefer docs |
| Craft / screenshots | Hero UI = Explore graph-home — see screenshot surfaces below. **Not** Phase 29 ops dual-card |
| Depth / IA / API | Frozen from S03/S04 — do not reshuffle order, select≠expand, budgets, `getImpact`, or invent `/v1/path` |
| Law 19 | Docs must not imply browser SoT or always-on daemon |
| Regression | If any CSS/TS touched: `cd web && npm run build` + `npm run test:e2e -- e2e/s03-depth.spec.ts e2e/s05-gates.spec.ts` still **PASS**. Docs-only → no test gate required |

### Must-answer doc locks (encode in quickstart + Notes)

| # | Topic | Locked content |
|---|--------|----------------|
| 1 | Default | `./bin/trace serve` → `http://127.0.0.1:7432` |
| 2 | Conflict | Second bind on same addr **fails** with friendly stderr (`FormatAddrInUseMessage`) — not silent hang |
| 3 | Multi-project | Example: project A on default; project B: `trace serve -C /path/to/other --addr 127.0.0.1:7433` then open that URL |
| 4 | Not shipped | No auto-port search / `--addr …:0` as default story — say “pick a free port” |
| 5 | Explorer pointer | After open: home is **Explore** graph + inspector (why/context/impact/…); Overview is secondary `/overview` |

### Locked quickstart snippet intent (wording may vary; must convey)

```bash
# Project A (default)
./bin/trace serve
# → http://127.0.0.1:7432

# Project B (distinct port — required when 7432 is taken)
./bin/trace serve -C /path/to/other --addr 127.0.0.1:7433
# → http://127.0.0.1:7433
```

On conflict, CLI already prints e.g.:

```text
serve: address already in use (127.0.0.1:7432)
hint: another process (often trace serve) is bound there.
  For a second project, pick a free port, e.g.:
    trace serve -C /path/to/other --addr 127.0.0.1:7433
```

Docs should point humans to that pattern **before** they hit the error (proactive multi-project section), not only “see stderr.”

### Residual polish targets (from S04-02)

| Residual | Severity | Action this row |
|----------|----------|-----------------|
| Sticky chrome `box-shadow` transition unused (blur present) | low/nit | **Optional** CSS fix in `app.css` / tokens — or Notes “deferred intentional” |
| Canvas keyboard select via list | nit | **Defer OK** — list path is valid `onSelect`; do not invent large a11y redesign |
| Packaging nits | low | Only if broken discoverability (e.g. stale Phase-29-only wording in `web/README`) |

**Do not** reopen S04 craft A/B/C debates. **Do not** strip depth.

### Optional screenshot / doc UI surfaces (if adding images or “what you see”)

Prefer citing live craft (S04 shipped):

| Surface | What to show / say |
|---------|-------------------|
| Shell | `.graph-shell` canvas-first; inspector ~`minmax(18rem, 26rem)` |
| Canvas | Taller `--graph-canvas-height` (`min(48rem, 100dvh-14rem)` or live token) |
| Inspector | Structured `PacketView` `dl`; raw JSON in `<details>` |
| Nodes | Calm center vs selected chrome (no glow) |
| Motion | `inspector-settle` / node transitions / sticky chrome — all honor `prefers-reduced-motion` |
| Paths | `web/src/styles/tokens.css`, `app.css` (`.graph-shell` / `.graph-inspector` / `.graph-node*`), `Inspector.tsx` |

Screenshot files (if any): place under an existing docs/assets convention if present; otherwise skip images and use prose + paths — **do not** invent a new media pipeline.

## Preflight (confirm in Notes)

1. `gui-quickstart.md` today: two-artifact path + security — **no** multi-project `--addr` section yet.
2. S02 #1 live: `IsAddrInUse` + `FormatAddrInUseMessage`; help mentions `127.0.0.1:7433`; **#2 not shipped**.
3. Explorer craft live at `/` (Graph home + Inspector) — Overview at `/overview`.

## Role work

### A — Docs (required)

1. Add **Multi-project / ports** (name flexible) to `docs/gui-quickstart.md`: default bind, fail-on-conflict, distinct `--addr` example, link or paraphrase CLI hint; state #2 auto-port **not** shipped.
2. Brief explorer orientation: open URL → graph home + inspector (not ops CRUD hero).
3. Update `web/README.md` if still single-port-only for multi-root readers.
4. Optionally one cross-link to `OPEN-PORT-MULTI.md` or ADR — keep quickstart human-short.

### B — Residual polish (optional / small)

1. Address sticky chrome box-shadow nit **or** list deferral with reason in Notes.
2. No feature work beyond packaging/copy/CSS nits.

### C — Verify

1. Docs accurate vs `FormatAddrInUseMessage` / serve help (re-read live strings).
2. If code changed: build + e2e s03+s05 PASS cited in Notes.
3. Board Notes cite paths changed.

## Exit criteria

- [ ] `docs/gui-quickstart.md` documents multi-project distinct `--addr` accurately vs S02 #1 (fail-on-conflict; no false auto-port claim)
- [ ] `web/README.md` not misleading for multi-project (updated or confirmed OK)
- [ ] Known S04 residuals addressed **or** explicitly deferred with reason in Notes
- [ ] If screenshots/UI callouts: Explore graph-home craft, not Overview-as-hero
- [ ] No new explorer features; Law 19 / bind defaults unchanged
- [ ] Board Notes cite paths (+ test evidence if code touched)

## Minimal todos

- [ ] Update `docs/gui-quickstart.md` — multi-project ports + brief explorer pointer
- [ ] Touch `web/README.md` if needed for ports / explorer wording
- [ ] Optional: residual CSS nit (chrome box-shadow) or defer in Notes
- [ ] Optional: screenshots of graph-home shell
- [ ] If code changed: `npm run build` + e2e s03-depth + s05-gates
- [ ] Update board row **P32-S05-01** Notes

## Todo updates

Status + notes on **P32-S05-01** only.

## Next

`P32-S05-02`
