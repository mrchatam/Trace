# P34-S04-01 — Docs + residual tests

## Metadata
- id: P34-S04-01
- todo_ids: [P34-S04-01]
- role: implementer
- skills: [incremental-implementation, writing-guidelines]
- mcps: []
- verification: automated
- hooks: []

## Objective

Flip **consumer-facing** docs (and any residual help/AGENTS lag) so the story matches shipped S02 + S03: **`.trace/` only**, SPA from Trace **binary embed**, **`trace gui` primary**, **default auto free-port** (`7432`–`7441`), explicit **`--addr` pin-strict**. Assert PLAN **T8**. **No** Explore craft / product behavior changes (docs + optional tiny help polish + residual tests only).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — **L1–L4** (do not reopen)
- [00-PLANNER.md](00-PLANNER.md)
- [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) — docs touch list + T8
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md) — consumer layout / docs audit
- Live anchors (read before editing — **post-S02/S03**):
  - `docs/gui-quickstart.md` — **primary debt** (still: build `web/dist`, “embedded stub”, “does **not** auto-pick a free port”, manual `--addr` multi-project)
  - `web/README.md` — still “no auto free-port”; production build as default StaticDir story
  - `AGENTS.md` — Phase 34 focus OK; orchestrator paste still says Next `P34-S01-01`; P33 line still “two-artifact / stub until P34”
  - Root `README.md` — only points at gui-quickstart (low debt; fix link blurb if it implies consumer `web/`)
  - `cmd/trace/help.go` + `local_http.go` usage — **already** auto-port + embed honesty (S03/S02); spot-check only
  - `internal/httpapi/embeddist/README.md` — **already** consumer `.trace/` + pipeline (S02); polish if gaps
  - `internal/httpapi/addr_in_use.go` `FormatAddrInUseMessage` — auto-hop hint already (S03); do not re-implement hop
  - `internal/httpapi/static.go` `placeholderHTML` — packaging-mistake tone OK (S02); touch only if still teaches consumer `web/` primary
  - CONTRIBUTING.md — no gui consumer teaching found (RESEARCH); skip unless spot-check finds drift

## Session start

Follow agent-loop-protocol. Unattended: execute locked touch list below. Do **not** reopen L1–L4 or re-implement auto-port/embed.

## Locked defaults (FINAL — do not re-debate)

| Item | Value |
|------|-------|
| Primary consumer story | From any Trace-initialized repo: `trace gui` → real Explore SPA from **binary**; footprint **`.trace/` only**; **no** project `web/` required |
| Multi-project | Second default `trace gui` / `serve` **auto-picks** next free loopback port in **`7432`–`7441`**; prints + (`gui`) opens **chosen** URL. One process = one root (L4) |
| `--addr` | Document as **pin**: set on cmdline → **fail if busy** (no hop). User-facing copy may say “`--addr` set on cmdline” / pin-strict — **do not** invent a `flag.Changed` API in docs (S03 used `fs.Visit("addr")`; stdlib has no `Flag.Changed`) |
| Remove / demote | Consumer two-artifact primary; “build `web/dist` in your app”; “embedded stub” as everyday; “no auto free-port” / “does not auto-pick” for **default** bind |
| Contributor DX | Trace **checkout** may `cd web && npm run build` / Vite; disk `<root>/web/dist` wins when present; `scripts/embed-gui.sh` / `make embed-gui` for embed refresh — label **contributor**, not consumer |
| PATH honesty | Keep: `go build`/`cp` PATH install ≠ `trace install` (agents/MCP/hooks only) |
| Help / usage | S03 owns port lines; S02 owns static honesty. S04: **spot-check** only; edit only if residual consumer-`web/` or no-auto wording remains |
| Residual tests | PLAN **T1–T7, T10, T11** already PASS in S02/S03. S04 residual = **T8** (docs/help/quickstart assertions). Optional: extend existing help tests if you touch `help.go`/`--help` strings. **No** new hop/embed product tests required unless a doc change breaks an existing test |
| Out | VERIFY / DR-HANDOFF (S05); Explore UI; changing L1–L4; re-implementing listen hop |

### Answers locked for implementer (planner gate)

1. **Touch priority (order):** (1) `docs/gui-quickstart.md` **rewrite** — primary; (2) `web/README.md` — contributor-labeled + auto-port; (3) `AGENTS.md` — next-runnable + drop stale “two-artifact until P34” once S02/S03 shipped; (4) root `README.md` if needed; (5) polish `embeddist/README.md` if needed; (6) help/usage/`placeholderHTML` only if residual lies remain.
2. **T8 pass criteria:** After edits, greppable consumer surfaces must **not** teach: (a) app-repo `web/` / two-artifact as **required** consumer path; (b) “no auto free-port” / “does not auto-pick” for **default** `gui`/`serve`. Must **positively** state: embed SPA for consumers; default busy → next free port `7432`–`7441`; `--addr` pins.
3. **Residual PLAN matrix:** Mark T8 done in board Notes. Explicitly note T1–T7/T10/T11 = already done (S02/S03) — N/A to re-run unless regression. T9 = S05.

## Live debt map (verified 2026-08-21, post-S03-02)

| Path | Status vs L1–L3 | S04 action |
|------|-----------------|------------|
| `docs/gui-quickstart.md` | **FAIL** — secondary section builds `web/dist`; “embedded stub”; multi-project “does **not** auto-pick”; manual `--addr` + old addr-in-use sample | **Rewrite** primary + multi-project + static table |
| `web/README.md` | **FAIL** — “no auto free-port”; production build as default StaticDir story | Relabel contributor; document default auto-port; Vite `--addr` pin when proxying |
| `AGENTS.md` | Stale next-row (`P34-S01-01`); P33 “two-artifact / stub until P34” outdated after S02 | Update orchestrator paste + Current focus next → **`P34-S04-02`** or post-01 **`P34-S05-00`** as appropriate when leaving Notes; drop “until P34 embed” once accurate |
| Root `README.md` | Low debt (link only) | Touch only if wording implies consumer `web/` |
| `cmd/trace/help.go` | **PASS** (auto-port + embed) | Spot-check; no rewrite expected |
| `cmd/trace/local_http.go` usage | **PASS** | Spot-check |
| `embeddist/README.md` | **PASS** (consumer `.trace/`, pipeline) | Optional polish |
| `FormatAddrInUseMessage` | **PASS** (auto-hop hint) | Do not re-implement; quickstart sample stderr should match live copy |
| `placeholderHTML` | **PASS** (packaging mistake; no consumer `web/` required) | Touch only if drift |
| CONTRIBUTING.md | Clean | Skip unless drift |

## Role work

### Minimal todos (execute in order)

1. **`docs/gui-quickstart.md` (primary)** — Rewrite so:
   - Launch: `trace gui` in a Trace-initialized **app** repo (no `web/` step).
   - Assets: SPA from Trace binary; consumer needs `.trace/` only; disk `web/dist` = contributor/override when present.
   - Multi-project: second default bind **auto-hops** `7432`→`7441`; show two-process example **without** requiring `--addr`; document `--addr` as optional **pin** (fail if busy).
   - Replace old addr-in-use sample with live `FormatAddrInUseMessage` shape (auto-hop + pin tip) **or** exhausted-range tip — do not paste P32 “pick a free port yourself” as the happy path.
   - Static dir row: default path string may still be `<root>/web/dist`, but explain **resolution**: disk if index present → **embed** → placeholder; consumers rarely need `--static-dir`.
   - Demote “build SPA in Trace checkout / embed-gui” to a short **contributor** subsection (or point at `embeddist/README.md` / `web/README.md`).
   - Keep PATH ≠ `trace install`, security defaults, Law 19 / loopback unchanged.
   - OPEN-PORT-MULTI link: keep as historical design notes if useful, but **do not** let it override L3 happy-path teaching in this doc.

2. **`web/README.md`** — Label top as **contributor / Trace checkout DX**. Fix multi-root paragraph: default auto free-port; use `--addr` to pin (and align Vite proxy). Production build = feed embed pipeline / disk override — not consumer requirement.

3. **`AGENTS.md`** — Update orchestrator “Next runnable” + Current focus next-row to match board after this implement (at least clear stale `P34-S01-01`). Soften/remove “still two-artifact / stub fallback until P34 embed” now that S02 shipped.

4. **Root `README.md`** — Only if quickstart blurb implies consumer must build `web/`.

5. **Polish** — `embeddist/README.md` if any leftover two-artifact everyday wording; help/usage/`placeholderHTML` only on residual lies.

6. **T8 evidence** — Grep (or equivalent) consumer surfaces for forbidden phrases; list hits fixed. Optional: if you change help strings, keep `TestHelp*` / gui/serve help tests green (`go test ./cmd/trace/ -count=1` focused).

7. **Board Notes** — Files touched; T8 evidence; residual matrix N/A list (T1–T7/T10/T11 already done); confirm no product hop/embed code changes.

### Suggested quickstart shape (not mandatory verbatim)

```markdown
## Launch (primary)
trace gui   # consumer repo; SPA from binary; default :7432

## Multi-project
# Project A
trace gui                    # → :7432 (or next free if busy)
# Project B (another cwd / -C)
trace gui -C /path/to/other  # → auto :7433…7441; opens that URL
# Pin (strict)
trace gui --addr 127.0.0.1:7433   # fail if busy
```

### Out of this row

- VERIFY / DR-HANDOFF (**S05**)
- Re-implementing `ListenAndServe` hop / embed pipeline
- Explore UI / craft redesign
- Changing DESIGN-LOCKS L1–L4

## Exit criteria

- [ ] `docs/gui-quickstart.md` teaches embed + auto-port multi-project; no consumer two-artifact / no-auto default
- [ ] `web/README.md` contributor-labeled; auto-port accurate
- [ ] `AGENTS.md` next-row / P33 stub language not stale vs shipped S02/S03
- [ ] T8 satisfied with evidence in Notes; residual product tests N/A or green
- [ ] No product behavior changes beyond docs/help string polish
- [ ] Next **P34-S04-02**

## Todo updates

Status + notes on **P34-S04-01** only.

## Next

`P34-S04-02`
