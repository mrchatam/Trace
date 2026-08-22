# Scope 04 — board map

**S04 docs + residual tests (T8).** Serial: **P34-S04-00 → P34-S04-01 → P34-S04-02**. Start only after **P34-S03-02** PASS (confirmed).

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 604 | P34-S04-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock docs touch list + thicken 01/02 (this row) |
| 605 | P34-S04-01 | [01-implement.md](01-implement.md) | Implementer | Quickstart/`web`/AGENTS flip + T8; residual product tests N/A |
| 606 | P34-S04-02 | [02-review.md](02-review.md) | Reviewer | L1–L3 docs checklist + T8; small fixes OK |

## Inputs (verified baseline 2026-08-21, post-S03-02)

| Source | Fact |
|--------|------|
| [`../scope-01-plan/PLAN.md`](../scope-01-plan/PLAN.md) | Docs touch list; **T8** = help/usage/quickstart no consumer two-artifact / no “no auto free-port” for default |
| L1 / L2 / L3 | Consumer `.trace/` only; SPA from binary; default busy → auto free port; `--addr` pin strict |
| S02 | Real SPA in `embeddist`; `placeholderHTML` packaging tone; `embeddist/README` consumer `.trace/` |
| S03 | UA-incr `7432`–`7441` in `httpapi.ListenAndServe`; CLI `fs.Visit("addr")` ≡ pin; help/usage auto-port; `FormatAddrInUseMessage` auto-hop hint |
| Live `docs/gui-quickstart.md` | **Flipped** (S04-01) — consumer `trace gui` + binary embed + `.trace/` only; auto-hop `7432`–`7441`; `--addr` pin |
| Live `web/README.md` | **Flipped** (S04-01) — contributor DX; auto-port + pin accurate |
| Live `AGENTS.md` / `docs/TODO.md` | Next-runnable → **P34-S05-00** (S04-02 review); P33 “two-artifact until P34” removed |
| Live help/usage | **Already flipped** (S03/S02) — spot-check PASS in S04-02 |
| Live root `README.md` | Points at gui-quickstart only — no consumer `web/` implication |
| Residual tests | T1–T7, T10, T11 done in S02/S03; S04 owns **T8** only; T9 = S05 |

## Locked answers (P34-S04-00)

1. **Primary rewrite:** `docs/gui-quickstart.md` — embed + auto-port multi-project; demote Trace-checkout `web/` to contributor.
2. **Contributor:** `web/README.md` labeled DX; auto-port accurate; Vite + `--addr` pin when proxying.
3. **T8:** Grep-backed; no product hop/embed re-implementation; help already mostly done.
4. **`--addr` copy:** “set on cmdline” / pin-strict — do not document a nonexistent `flag.Changed` field API.

## Out of this scope

- VERIFY / DR-HANDOFF close (**S05**).
- Re-implementing embed pipeline or listen hop.
- Explore UI / craft redesign.
- Changing DESIGN-LOCKS L1–L4.
