# P33-S05-02 — Polish review

**Verdict:** **PASS**  
**Confidence:** **high**  
**Date:** 2026-08-21  
**Scope:** Theme C docs flip + S02–S04 polish residuals  
**Skills loaded:** code-review-and-quality · writing-guidelines (docs primary-story clarity)  
**Evidence:** live docs re-read; `go test ./internal/httpapi/ -run 'FormatAddrInUse|IsAddrInUse' -count=1` ok; Graph EmptyState + `.empty .btn--primary`; no S05 `tokens.css` rewrite

## Preflight

- [x] Fresh review vs P33-S05-01 implementer Notes; S01–S04 prompts untouched
- [x] P33-S05-01 `done`; this row was sole pending under S05
- [x] No product-feature scope creep (docs + residual polish only)

## Checklist evidence

### Docs primary story
- [x] `docs/gui-quickstart.md` H1 = ``Trace GUI quickstart (`trace gui`)``; first fence = `trace gui`; opens Explore `/`
- [x] PATH §: `CGO_ENABLED=1 go install …/cmd/trace@latest`; explicit **≠** `trace install` (agents/MCP/hooks only)
- [x] `serve` / `./bin/trace serve` under **Secondary** (scripting / no-browser) — demoted, not deleted
- [x] Multi-project: fail-on-conflict; **no** auto-port; pick `--addr` yourself (quickstart §Multi-project / ports)
- [x] Lead: “Not a daemon, not a hosted SaaS product”; Cloud path = future/separate
- [x] `README.md` gui-quickstart bullet → `trace gui` primary + `serve` headless/scripting
- [x] `web/README.md` lead = Operator SPA for `trace gui`; Dev examples primary `trace gui --no-open`
- [x] `AGENTS.md` Current focus: docs primary = `trace gui`; serve secondary (no “until docs flip” lag)
- [x] Related: `trace gui --help` + `trace serve --help` alongside

### Residuals (S02–S04)
- [x] EmptyState: inline primary `Link.btn.btn--primary` “Open Tasks” (`data-testid="explore-empty-open-tasks"`); `.empty .btn--primary { color: #fff }` contrast polish in `app.css`
- [x] Explore **canvas** screenshot: **deferred to VERIFY** (no S05 evidence dir; S04 PNGs remain list-focused) — track in S06
- [x] Addr-in-use dual-word: `FormatAddrInUseMessage` → `gui|serve:` + “trace gui or trace serve”; test asserts; quickstart sample matches
- [x] Craft literacy one-liner in quickstart “What you see”: chroma strip + text labels (not color alone)

### Non-regression
- [x] No palette/token rewrite undoing S04 (`tokens.css` kind/state + chroma strip left alone; polish CSS is EmptyState CTA only)
- [x] No compose / route / budget changes in S05 Notes / residual paths
- [x] Law 19 / loopback defaults unchanged in docs (127.0.0.1; adapters; no `/rpc` MCP)
- [x] `go test ./internal/httpapi/ -run 'FormatAddrInUse|IsAddrInUse' -count=1` **PASS**

## Findings

| Sev | Finding | Disposition |
|-----|---------|-------------|
| low | No Explore **canvas** screenshot in docs tree | Accept — deferred to VERIFY (S06) as implementer noted |
| nit | README still uses `./bin/trace` for version/build elsewhere | OK — not the GUI launch primary story |

**No blocker / high.** No spawn (`P33-S05-02a`/`02b`).

## Residual risks

- VERIFY should confirm live `trace gui --help` + smoke Explore graph and optionally capture one canvas shot if docs stay screenshot-light.
- Canvas keyboard arrow-roving remains accepted S03 residual (out of S05).

## Upcoming thickenings (reviewer rights)

- **S06-00/01** — VERIFY floor: docs primary = `trace gui`; PATH≠`trace install`; addr-in-use dual-word landed; **capture or explicitly waive** canvas screenshot; re-check Laws 6–7/19 + craft not gray-only; skills cites on S01/S03/S04.
- **S06-02** — DR-HANDOFF: default **no successor**; carry only thin residuals (canvas shot optional, canvas keyboard out-of-phase).

## Next runnable

**P33-S06-00**
