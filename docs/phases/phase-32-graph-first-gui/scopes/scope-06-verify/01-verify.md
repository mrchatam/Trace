# P32-S06-01 — VERIFY

## Metadata
- id: P32-S06-01
- todo_ids: [P32-S06-01]
- role: verify
- skills: [planning-and-task-breakdown, grinding-until-pass]
- mcps: []
- verification: mixed
- hooks: []

## Objective

Run the locked Phase 32 VERIFY floor after S00–S05. Aggregate prior scope evidence + live re-checks into **`VERIFY-NOTES.md`** (+ evidence dir). Confirm explorer bar vs DESIGN-LOCKS; Laws 6–7/19; builds; e2e smoke; **P32-PORT** addressed (**#1 shipped**, **#2 deferred**). **Leave `DR-HANDOFF.md` OPEN** — close owned by **P32-S06-02**. **No product code.** Do **not** start S06-02 or a successor phase.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [00-PLANNER.md](00-PLANNER.md) — VERIFY floor locks (FINAL — S06-00)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- [DR-HANDOFF.md](../../DR-HANDOFF.md) — remains **OPEN** until S06-02
- Prior artifacts (cite in notes):
  - S00 [`RESEARCH.md`](../scope-00-research/RESEARCH.md)
  - S01 [`UX-IA.md`](../scope-01-ux-ia/UX-IA.md)
  - S02 [`NO-GAPS.md`](../scope-02-api-gaps/NO-GAPS.md) + board Notes P32-S02-01/02
  - S03–S05 board Notes + reviews (graph-home, craft, port docs)
- Live anchors: `web/src/App.tsx`, `web/src/screens/Graph.tsx`, `web/src/components/Inspector.tsx`, `web/src/api/ops.ts`, `internal/httpapi/` (`FormatAddrInUseMessage` / `IsAddrInUse` / `DefaultAddr`), `cmd/trace/serve.go`, `docs/gui-quickstart.md`, `web/README.md`

## Session start

Follow agent-loop-protocol Session start. Unattended: do not stop after planning. This row runs verification and records evidence; it does **not** close DR-HANDOFF, decide successor, or change product bodies.

## Locked defaults (FINAL — S06-00)

| Item | Value |
|------|-------|
| Precondition | P32-00 … P32-S05-02 all `done`; S05-02 PASS high |
| Product / CSS / TS / Go changes | **Forbidden** (evidence + notes only). Failures → spawn remediation from this row or leave FAIL for S06-02 to spawn |
| Evidence dir | `experiments/runs/YYYY-MM-DD-p32-s06-01-verify/evidence/` |
| Notes artifact | `scopes/scope-06-verify/VERIFY-NOTES.md` (**required**) |
| DR-HANDOFF | Stays **OPEN** — S06-02 closes |
| Successor | **Out of scope** — S06-02 only (lean default **no successor**) |
| Graph tech | 2D `@xyflow/react` only — **no Three.js** |
| Budgets | `DEFAULT_MAX=50`, `UI_CAP=100` in `Graph.tsx`; neighborhood requires center + `max_nodes` |
| P32-PORT | **#1 + docs required**; **#2 deferred** (do not fail solely for missing auto-port) |

### Aggregate evidence map (S00–S05 — cite in VERIFY-NOTES)

| Scope | Must cite / re-check |
|-------|----------------------|
| S00 | RESEARCH peer bar + P32-PORT prefer #1 |
| S01 | UX-IA hybrid C + inspector order + Laws 6–7 budgets |
| S02 | `NO-GAPS.md`; `getImpact` in `ops.ts`; P32-PORT **#1** (`IsAddrInUse` / `FormatAddrInUseMessage` + serve help); **#2 not shipped** |
| S03 | Graph home `/` + Overview `/overview` + `/graph`→`/`; `Inspector.tsx` depth map; select≠expand; e2e s03+s05 |
| S04 | Craft A/B/C on depth shell; no Three.js; depth/IA frozen |
| S05 | `gui-quickstart` **Multi-project / ports** + `web/README`; OPEN-PORT #3/#4 closed; #2 still deferred |

### Fail vs residual (locked)

**Fail VERIFY for:** `web` build FAIL; focused P32-PORT Go tests FAIL; e2e s03+s05 FAIL; graph not home (index still Overview-only without Explore graph); inspector missing required sections (summary→why→context→impact→reviews→links for applicable types); `getImpact` missing from `ops.ts`; unbounded full-graph dump CTA / budgets removed; Three.js / 3D default; Law 19 fork (business logic in `web/` beyond adapters); P32-PORT **#1** missing (no friendly in-use path / no `--addr` guidance); docs claim auto-port while #2 not shipped; loopback default softened to public bind; `/v1/path` invented in product.

**Do not fail VERIFY solely for residuals below** (record in VERIFY-NOTES):

| Residual | Disposition |
|----------|-------------|
| P32-PORT **#2** auto free-port / `:0` | Deferred with reason — OPEN-PORT-MULTI |
| Serve stderr prints “listening on” before bind fails | S02 low residual — non-blocking |
| Sticky chrome `box-shadow` transition unused (blur present) | S04/S05 deferred nit |
| Canvas keyboard select via list (`onSelect`) | Acceptable; same handler |
| No explorer screenshots / media pipeline | S05 deferred — no docs/assets pipeline |
| Optional denser craft polish beyond S04 | Out of phase bar if S04 craft A/B/C present |

## Locked verify command floor

Run from repo root unless noted. Tee outputs into evidence dir. Use `date +%Y-%m-%d` for the run folder name.

### Block 0 — Evidence dir + preflight

```bash
cd /home/ali/Desktop/Trace
RUN_DATE=$(date +%Y-%m-%d)
EVID="experiments/runs/${RUN_DATE}-p32-s06-01-verify/evidence"
mkdir -p "$EVID"
{
  echo "verify_id=P32-S06-01"
  echo "date=$RUN_DATE"
  echo "git_sha=$(git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "precondition=P32-S05-02 PASS high; S00–S05 done"
} > "$EVID/00-run-metadata.txt"
```

**Pass:** `$EVID` exists; metadata cites S05-02 PASS + S00–S05 complete.

### Block 1 — Web build (required)

```bash
(cd web && npm run build) 2>&1 | tee "$EVID/01-web-build.txt"
```

**Pass:** exit 0.

### Block 2 — P32-PORT Go tests (required)

```bash
go test ./internal/httpapi/ -run 'TestIsAddrInUse|TestFormatAddrInUse' -count=1 2>&1 | tee "$EVID/02-httpapi-port.txt"
go test ./cmd/trace/ -run 'TestServe' -count=1 2>&1 | tee "$EVID/03-serve-tests.txt"
```

**Pass:** exit 0; includes friendly addr-in-use coverage (`TestServeAddrInUseFriendlyMessage` or equivalent under `TestServe`).

### Block 3 — Explorer e2e smoke (required)

```bash
(cd web && npm run test:e2e -- e2e/s03-depth.spec.ts e2e/s05-gates.spec.ts) 2>&1 | tee "$EVID/04-e2e-explorer.txt"
```

**Pass:** exit 0; tests passed (expect **6** or document actual count if suite grew without dropping depth/gates coverage).

### Block 4 — DESIGN-LOCKS + Laws + explorer bar spot-check (required)

```bash
{
  echo "=== graph home / routes ==="
  grep -n 'path=\"/\"\|path=\"/overview\"\|path=\"/graph\"\|Explore\|Graph' web/src/App.tsx web/src/components/Nav.tsx 2>/dev/null | head -80
  echo "=== budgets ==="
  grep -n 'DEFAULT_MAX\|UI_CAP' web/src/screens/Graph.tsx
  echo "=== inspector + getImpact ==="
  grep -n 'getImpact\|getWhy\|getContext\|summary\|reviews' web/src/components/Inspector.tsx web/src/api/ops.ts | head -60
  echo "=== no Three.js / no /v1/path in web/src ==="
  (grep -RIn 'three\|Three\.js\|/v1/path' web/src 2>/dev/null || echo "PASS: no three|/v1/path hits") | head -20
  echo "=== P32-PORT helpers + DefaultAddr ==="
  grep -n 'FormatAddrInUseMessage\|IsAddrInUse\|DefaultAddr\|7432\|7433' internal/httpapi/*.go cmd/trace/serve.go cmd/trace/help.go 2>/dev/null | head -80
} 2>&1 | tee "$EVID/05-locks-spotcheck.txt"
```

**Pass:** index is Graph / Explore-first; `/overview` secondary; budgets 50/100; inspector depth sections present; `getImpact` in `ops.ts`; no Three.js default; no `/v1/path`; `DefaultAddr` still `127.0.0.1:7432`; `FormatAddrInUseMessage` / `IsAddrInUse` present; help/serve mention distinct `--addr`.

### Block 5 — P32-PORT docs + security (required)

```bash
{
  echo "=== gui-quickstart Multi-project / ports ==="
  grep -n 'Multi-project\|7432\|7433\|auto-port\|allow-remote\|fail' docs/gui-quickstart.md | head -60
  echo "=== web README multi-root ==="
  grep -n 'addr\|7432\|Explore\|multi' web/README.md | head -40
  echo "=== FormatAddrInUseMessage live string ==="
  grep -n 'FormatAddrInUseMessage\|address already\|--addr' internal/httpapi/*.go | head -40
  echo "=== OPEN-PORT-MULTI status ==="
  head -20 docs/phases/phase-32-graph-first-gui/OPEN-PORT-MULTI.md
} 2>&1 | tee "$EVID/06-port-docs-security.txt"
```

**Pass checklist (all required in VERIFY-NOTES):**

1. **Behavior #1:** `FormatAddrInUseMessage` + `IsAddrInUse`; serve fail-on-conflict + friendly stderr; help/usage distinct `--addr` e.g. `127.0.0.1:7433`; default `127.0.0.1:7432`. **#2 not shipped.**
2. **Docs #3/#4:** `docs/gui-quickstart.md` section **Multi-project / ports** (default bind, fail-on-conflict, `-C` + `--addr`, **no auto-port claim**) + `web/README.md` multi-root/`--addr`. Wording consistent with live `FormatAddrInUseMessage`.
3. **Security:** loopback default; remote still `--allow-remote` + token — no public-bind softening in docs or serve defaults.
4. Pointer: [`OPEN-PORT-MULTI.md`](../../OPEN-PORT-MULTI.md) status (#1 + S05 docs; #2 defer).

### Block 6 — Residuals list (required in notes)

Record non-blocking residuals from the locked table. Do **not** fail solely for them. Overall PASS only if blocks 0–5 green.

## Do not

- Change product Go/TS/CSS/docs in this row (spawn or FAIL instead)
- Close DR-HANDOFF (that is **P32-S06-02**)
- Claim/scaffold successor beyond citing lean **no successor** (S06-02 only)
- Treat missing #2 auto-port as VERIFY FAIL
- Fake PASS without evidence dir + notes

## VERIFY-NOTES.md template (required)

Write `docs/phases/phase-32-graph-first-gui/scopes/scope-06-verify/VERIFY-NOTES.md`:

```markdown
# VERIFY-NOTES — P32-S06-01

**Date:** YYYY-MM-DD
**Git SHA:** …
**Overall:** PASS | FAIL
**Evidence:** experiments/runs/YYYY-MM-DD-p32-s06-01-verify/evidence/
**Precondition:** P32-S05-02 PASS high; S00–S05 done

| Block | Result | Notes |
|------:|--------|-------|
| 0 Preflight | PASS/FAIL | |
| 1 web build | PASS/FAIL | |
| 2 P32-PORT Go tests | PASS/FAIL | httpapi + TestServe |
| 3 e2e explorer smoke | PASS/FAIL | s03-depth + s05-gates |
| 4 DESIGN-LOCKS / Laws / explorer | PASS/FAIL | hybrid C; 50/100; Law 19; no Three.js |
| 5 P32-PORT docs + security | PASS/FAIL | #1+#3+#4; #2 deferred |
| 6 Residuals | listed | |

## Aggregate (S00–S05)
- S00 RESEARCH: …
- S01 UX-IA: …
- S02 NO-GAPS + getImpact + PORT #1: …
- S03 graph-home + inspector: …
- S04 craft A/B/C: …
- S05 port docs: …

## P32-PORT tick
- [ ] #1 behavior evidenced
- [ ] #3/#4 docs evidenced
- [ ] #2 deferred (explicit)
- [ ] loopback / --allow-remote intact

## Residuals (non-blocking)
- …

## Failures (if any)
- …

## DR-HANDOFF
remains OPEN — close owner **P32-S06-02**

## Next
P32-S06-02
```

## Todo updates

Status + notes on **P32-S06-01** only. Do not mark S06-02 done.

## Exit criteria

- [ ] `VERIFY-NOTES.md` with overall PASS/FAIL and blocks 0–6
- [ ] Evidence dir under `experiments/runs/…-p32-s06-01-verify/evidence/`
- [ ] **P32-PORT** explicitly ticked (#1 + docs; #2 deferred)
- [ ] DESIGN-LOCKS / Laws 6–7/19 / explorer bar addressed in notes
- [ ] Board Notes summarize results + evidence path
- [ ] DR-HANDOFF still **OPEN**
- [ ] Next: **P32-S06-02**

## Next

`P32-S06-02`
