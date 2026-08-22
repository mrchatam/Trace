# scope-01-architecture todos

Board: [`docs/TODO/phase-29.md`](../../../../TODO/phase-29.md) orders **506–508**.

| Order | ID | Prompt | Status |
|------:|----|--------|--------|
| 506 | P29-S01-00 | [00-PLANNER.md](00-PLANNER.md) | done (planner) |
| 507 | P29-S01-01 | [01-adr-and-openapi.md](01-adr-and-openapi.md) | done |
| 508 | P29-S01-02 | [02-review.md](02-review.md) | done (review PASS; inline gate/reset + seed path; no 02a/02b) |

## Artifact targets (locked S01-00)

| Artifact | Path |
|----------|------|
| ADR | `docs/adr/ADR-HTTP-API-GUI.md` |
| OpenAPI | `api/openapi.yaml` |

## Notes

- Planner locks RESEARCH §1–7 (auth loopback-trust, disk static first, no CORS `*`, no MCP `/rpc`, budgeted graph, paths above).
- Implementer authors ADR + OpenAPI only — no `internal/httpapi` / `web/` / `serve`.
- Reviewer (P29-S01-02): added loop gate/reset paths; seed path confinement; thickened S02 prompts. Next board row **P29-S02-00**.
