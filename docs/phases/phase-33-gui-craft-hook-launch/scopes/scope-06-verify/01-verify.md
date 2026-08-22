# P33-S06-01 — VERIFY

## Metadata
- id: P33-S06-01
- todo_ids: [P33-S06-01]
- role: verify
- skills: [planning-and-task-breakdown, grinding-until-pass]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Run the locked Phase 33 VERIFY floor after S00–S05. Aggregate prior scope evidence + live re-checks into **`VERIFY-NOTES.md`** (+ evidence dir). Confirm Themes **A–C** (craft + Explore hook + `trace gui` docs/launch); Laws **6–7/19**; S05 deferred **canvas** shot capture-or-waive. **Leave `DR-HANDOFF.md` OPEN** — close owned by **P33-S06-02**. **No product code.** Do **not** start S06-02 or invent a successor phase.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S06-00)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — Themes A–C
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S06-02
- Prior PASS artifacts (cite in notes):
  - S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md) + [`REVIEW`](../scope-00-research/02-review.md) Notes
  - S01 [`DESIGN.md`](../scope-01-design-ux/DESIGN.md) / [`UX-IA.md`](../scope-01-design-ux/UX-IA.md) + [`REVIEW.md`](../scope-01-design-ux/REVIEW.md)
  - S02 [`REVIEW.md`](../scope-02-gui-launch/REVIEW.md)
  - S03 [`REVIEW.md`](../scope-03-explore-graph/REVIEW.md)
  - S04 [`REVIEW.md`](../scope-04-color-craft/REVIEW.md) + `evidence/explore-{light,dark}.png`
  - S05 [`REVIEW.md`](../scope-05-polish/REVIEW.md) — docs primary; canvas **deferred here**
- Live anchors: `cmd/trace/` (`gui`, `local_http.go`, help), `internal/httpapi/addr_in_use.go`, `web/src/screens/Graph.tsx`, `web/src/lib/overviewCompose.ts`, `web/src/styles/tokens.css`, `docs/gui-quickstart.md`, `README.md`, `web/README.md`, `AGENTS.md`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or change product bodies.

## Locked defaults (FINAL — S06-00)

| Item | Value |
|------|-------|
| Precondition | P33-00 … P33-S05-02 all `done`; S05-02 PASS high |
| Product / CSS / TS / Go / docs changes | **Forbidden** (evidence + notes only). Failures → spawn remediation from this row or leave FAIL for S06-02 to spawn |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p33-s06-01-verify/evidence/` |
| Notes artifact | `scopes/scope-06-verify/VERIFY-NOTES.md` (**required**) |
| Optional canvas PNG | `scopes/scope-06-verify/evidence/explore-canvas.png` **or** explicit waive in VERIFY-NOTES |
| DR-HANDOFF | Stays **OPEN** — S06-02 closes |
| Successor | **Out of scope** — S06-02 only (lean default **no successor**) |
| Docs primary (S05 PASS) | `gui-quickstart` H1+first fence = **`trace gui`**; README / `web/README` / AGENTS aligned; PATH = `go install …/cmd/trace@…` ≠ `trace install`; `serve` secondary |
| Addr-in-use | Dual-word `gui\|serve` already landed — **re-spot-check** only |
| Budgets (Laws 6–7) | seeds target **6** / ≤**8**; `SEED_MAX_NODES=40` / `DEPTH=2`; `UI_CAP=100`; expand ≤**50** (`overviewCompose.ts` / Graph) |
| Graph tech | 2D `@xyflow/react` only — **no Three.js** |
| Explore route | `/` Graph overview-on-open — **≠** Nav `/overview` |

### Themes A–C acceptance map (must tick in VERIFY-NOTES)

| Theme | Acceptance (spot-check + prior REVIEW cite) |
|-------|-----------------------------------------------|
| **A — Color/craft** | Forest-moss tokens + `--kind-*`/`--state-*`; chroma strip 3px + 14% fill; color-not-only labels; not gray-only; S04 PNGs and/or live; skills cited on S01/S04 board Notes |
| **B — Explore hook** | Open `/` → overview-on-open (not empty center-gate); interactive graph; budgets above; select≠expand; e2e `s03-depth` green; skills on S03 Notes |
| **C — Launch/docs** | `trace gui --help` (+ `--no-open` / shared serve flags); PATH teach `go install` ≠ `trace install`; docs primary `trace gui`; fail-on-conflict / no auto-port; addr-in-use `gui\|serve` |

### Aggregate evidence map (S00–S05 — cite in VERIFY-NOTES)

| Scope | Must cite / re-check |
|-------|----------------------|
| S00 | RESEARCH leans: Explore≠`/overview`; PATH≠install; reject full dump / UA auto-port |
| S01 | DESIGN + UX-IA; Operate forest-moss; D+B+C; budgets 6≤8 / 40 / UI_CAP=100 / depth 2; skills gate |
| S02 | `trace gui` + `--no-open`; land `http://{addr}/`; Law 19; PATH teach; no auto-port |
| S03 | overviewCompose + Graph loadOverview; e2e overview-first; canvas keyboard residual accepted |
| S04 | tokens + chroma; contrast Notes; `evidence/explore-{light,dark}.png` (list-heavy OK) |
| S05 | docs primary flip; EmptyState CTA; addr-in-use dual-word; craft literacy; **canvas deferred → this row** |

### Fail vs residual (locked)

**Fail VERIFY for:** `web` build FAIL; focused Go gui/httpapi tests FAIL; e2e `s03-depth` FAIL; Explore still center-gate-only on open; budgets removed / full-dump CTA; Explore relocated off `/` to `/overview`; gray-only / no kind chroma (Theme A regression); `trace gui` missing or docs still lead with `serve` as primary; PATH conflated with `trace install`; auto-port claimed while not shipped; Law 19 fork (business logic in `web/` beyond adapters); loopback default softened to public bind; Three.js / 3D default.

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES):

| Residual | Disposition |
|----------|-------------|
| Explore **canvas** screenshot missing | Capture under this scope `evidence/` **or** **explicit waive** with reason (S04 list PNGs remain valid craft evidence) |
| Canvas keyboard arrow-roving | Accepted S03 residual — out of phase |
| S04 PNGs list-heavy (not canvas-forward) | Accept — Theme A craft already PASS |
| Optional denser craft polish beyond S04 | Out of phase bar if Theme A present |
| Hosted SaaS / brew/deb packages | Out of scope |

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir. Use `date +%Y-%m-%d` for the run folder name.

### Block 0 — Evidence dir + preflight

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p33-s06-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P33-S06-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=P33-S05-02 PASS high; S00–S05 done"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists; metadata cites S05-02 PASS + S00–S05 complete.

### Block 1 — Web build (required)

```bash
(cd web && npm run build) 2>&1 | tee "$EVID/01-web-build.txt"
```

**Pass:** exit 0.

### Block 2 — Go launch + addr-in-use (required)

```bash
go test ./cmd/trace/ -count=1 2>&1 | tee "$EVID/02-cmd-trace.txt"
go test ./internal/httpapi/ -run 'FormatAddrInUse|IsAddrInUse' -count=1 2>&1 | tee "$EVID/03-httpapi-addr.txt"
go run ./cmd/trace gui --help 2>&1 | tee "$EVID/04-gui-help.txt"
```

**Pass:** exit 0 on tests; help shows `gui`, `--no-open`, loopback/`7432` (or shared serve flags); FormatAddrInUse asserts / live string includes `gui|serve`.

### Block 3 — Explore unit + e2e smoke (required)

```bash
(cd web && node --experimental-strip-types --test src/lib/overviewCompose.test.ts) 2>&1 | tee "$EVID/05-overview-compose.txt"
(cd web && npm run test:e2e -- e2e/s03-depth.spec.ts) 2>&1 | tee "$EVID/06-e2e-explore.txt"
```

**Pass:** unit 7/7 (or document count if suite grew without dropping coverage); e2e exit 0 (overview-first, select≠expand).

### Block 4 — Themes A–C + Laws spot-check (required)

```bash
{
  echo "=== Theme C docs primary ==="
  head -40 docs/gui-quickstart.md
  grep -n 'trace gui\|trace install\|go install\|Secondary\|serve' docs/gui-quickstart.md README.md web/README.md AGENTS.md | head -80
  echo "=== Theme B route + budgets ==="
  grep -n 'path=\"/\"\|path=\"/overview\"\|Graph\|Explore' web/src/App.tsx | head -40
  grep -n 'SEED_TARGET\|SEED_CAP\|SEED_MAX_NODES\|UI_CAP\|EXPAND_MAX_NODES\|DEPTH' web/src/lib/overviewCompose.ts
  echo "=== Theme A tokens / chroma ==="
  grep -n -- '--kind-\|--state-\|chroma\|graph-node__kind' web/src/styles/tokens.css web/src/styles/app.css 2>/dev/null | head -60
  echo "=== Laws / rejects ==="
  (grep -RIn 'three\|Three\.js\|/v1/path' web/src 2>/dev/null || echo "PASS: no three|/v1/path hits") | head -20
  grep -n 'FormatAddrInUseMessage\|gui|serve\|DefaultAddr\|7432' internal/httpapi/*.go cmd/trace/*.go 2>/dev/null | head -60
} 2>&1 | tee "$EVID/07-themes-locks-spotcheck.txt"
```

**Pass checklist (all required in VERIFY-NOTES):**

1. **Theme C:** quickstart H1+first fence `trace gui`; PATH `go install …/cmd/trace@…` ≠ `trace install`; serve secondary; multi-project fail-on-conflict / no auto-port; addr-in-use `gui|serve`.
2. **Theme B:** `App` index → Graph; budgets 6/8/40/100/50/depth2; overview compose present.
3. **Theme A:** `--kind-*` / `--state-*` + chroma strip rules; not gray-only; cite S04 evidence PNGs path.
4. **Laws 6–7/19:** no full dump; adapters only; loopback default intact; no Three.js.

### Block 5 — Design skills cites (required)

From board Notes / REVIEW artifacts for **S01, S03, S04** implement+review: confirm metadata/Notes cite **impeccable** + **ui-ux-pro-max** + **frontend-design**. Record paths in VERIFY-NOTES.

**Pass:** all three scopes cite the skills gate (DESIGN-LOCKS).

### Block 6 — Canvas shot or waive (required)

| Option | Action |
|--------|--------|
| **A — Capture** | With GUI running (optional live), save one Explore **canvas**-forward PNG to `scopes/scope-06-verify/evidence/explore-canvas.png` (or under `$EVID/`) and link in VERIFY-NOTES |
| **B — Waive** | Explicit sentence in VERIFY-NOTES: waive canvas shot; rely on S04 `explore-{light,dark}.png` + live Theme A spot-check; reason (e.g. no media pipeline / list PNGs sufficient for craft bar) |

**Do not fail** solely for choosing B.

### Block 7 — Residuals list (required in notes)

Record non-blocking residuals from the locked table. Overall PASS only if blocks 0–5 green (+ block 6 A or B recorded).

## Do not

- Change product Go/TS/CSS/docs in this row (spawn or FAIL instead)
- Close DR-HANDOFF (that is **P33-S06-02**)
- Claim/scaffold successor beyond citing lean **no successor** (S06-02 only)
- Reopen S05 docs flip or S04 tokens as “fix”
- Treat waived canvas shot or canvas keyboard residual as VERIFY FAIL
- Fake PASS without evidence dir + notes

## VERIFY-NOTES.md template (required)

Write `docs/phases/phase-33-gui-craft-hook-launch/scopes/scope-06-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — P33-S06-01

**Date:** YYYY-MM-DD
**Git SHA:** …
**Overall:** PASS | FAIL
**Evidence:** experiments/runs/YYYY-MM-DD-p33-s06-01-verify/evidence/
**Precondition:** P33-S05-02 PASS high; S00–S05 done

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS/FAIL | |
| 1 web build | PASS/FAIL | |
| 2 Go gui + addr-in-use + help | PASS/FAIL | |
| 3 overviewCompose + e2e s03-depth | PASS/FAIL | |
| 4 Themes A–C + Laws spot-check | PASS/FAIL | |
| 5 Design skills cites S01/S03/S04 | PASS/FAIL | |
| 6 Canvas shot | CAPTURED path… \| WAIVED reason… | |
| 7 Residuals | listed | |

## Themes A–C
- [ ] A craft — …
- [ ] B Explore hook — …
- [ ] C launch/docs — …

## Aggregate (S00–S05)
- S00 RESEARCH: …
- S01 DESIGN/UX-IA: …
- S02 gui launch: …
- S03 Explore: …
- S04 craft: … (PNG paths)
- S05 polish/docs: …

## Residuals (non-blocking)
- …

## Failures (if any)
- …

## DR-HANDOFF
remains OPEN — close owner **P33-S06-02**

## Next
P33-S06-02
```

## Todo updates

Status + notes on **P33-S06-01** only. Do not mark S06-02 done.

## Exit criteria

- [ ] `VERIFY-NOTES.md` with overall PASS/FAIL and blocks 0–7
- [ ] Evidence dir under `experiments/runs/…-p33-s06-01-verify/evidence/`
- [ ] Themes **A–C** explicitly ticked
- [ ] Canvas shot **captured or waived**
- [ ] DESIGN-LOCKS / Laws 6–7/19 addressed in notes
- [ ] Board Notes summarize results + evidence path
- [ ] DR-HANDOFF still **OPEN**
- [ ] Next: **P33-S06-02**

## Next

`P33-S06-02`
