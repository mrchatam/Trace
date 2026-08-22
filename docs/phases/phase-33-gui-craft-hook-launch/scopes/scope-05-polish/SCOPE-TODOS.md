# Scope 05 — board map

**S05 polish + docs** — primary launch story flip to **`trace gui`**. Serial: **P33-S05-00 → P33-S05-01 → P33-S05-02**.

| Order | Board ID | Prompt | Role | Artifact / duty |
|------:|----------|--------|------|-----------------|
| 586 | P33-S05-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock doc list + residuals (**done** after planner) |
| 587 | P33-S05-01 | [01-implement.md](01-implement.md) | Implementer | Docs flip + residual fixes |
| 588 | P33-S05-02 | [02-review.md](02-review.md) | Reviewer | **done** — [`REVIEW.md`](REVIEW.md) PASS → S06 |

## Locked (P33-S05-00)

| Item | Value |
|------|-------|
| Primary | `trace gui` + PATH `go install …/cmd/trace@…` |
| Secondary | `trace serve` / `./bin/trace serve` (scripting / no-browser) |
| Required docs | `docs/gui-quickstart.md`, `README.md`, `web/README.md`, `AGENTS.md` |
| No retokenize | S04 forest-moss + chroma strip **PASS** — leave tokens alone |
| EmptyState | Live already has inline Open Tasks — verify/polish only |
| Optional | Canvas shot for docs; addr-in-use `gui`\|`serve`; craft literacy one-liner |

## Residuals carried in

| From | Item | Disposition |
|------|------|-------------|
| S04 | EmptyState CTA / list-heavy evidence | Verify inline CTA; optional canvas shot |
| S04 | Do not reinvent craft tokens | Hard lock |
| S02 | Addr-in-use still serve-worded | Optional dual-word |
| S02 | PATH tip already in quickstart | Keep; flip **lead** only |

## Out of this scope

- Phase VERIFY + DR-HANDOFF (S06)
- Compose / budgets / routes / canvas keyboard residual
- Hosted SaaS; auto-port; PATH via `trace install`
