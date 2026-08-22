# Scope S03 — Agent X0 harness (CLI)

**Depends-on:** Board order after S02 (done). Hard need: Phase 00 CLI `why`/`context` + `fixtures/x0` (live).

**Locks (P01-S03-00):**
| Item | Value |
|------|-------|
| B0 | Agent + ordinary repo tools; **no** `trace why`/`context` |
| G1 | Agent + `trace` CLI `why`/`context` (scripted; may still read repo) |
| Corpus | `fixtures/x0` + abs `seed/gt.json`; Task `22222222-…222` |
| Package | `evals/x0` + committed `schema.json` v1 |
| Metrics | Temp `metrics-b0.json` + `metrics-g1.json`; `dry_run: true`; both conditions |
| Dry-run ≠ Gate C | Emit + validate schema only; no superiority scoring |
| Keep separate | `evals/p0x` (P0-X 7/7); `evals/honesty` (H5) — do not merge |
| MCP | Not required for X0 (DR-AGENT) |
| DONE | Prefer no DONE in dry-run; if needed → Review PASS or `allow_done` |

**S01 note:** EvidenceIDs alone ≠ DONE.

**S02 note:** Honesty stays in `evals/honesty`; X0 is B0/G1 harness only.

- [x] P01-S03-00 planner
- [x] P01-S03-01 implement
- [x] P01-S03-02 review (+ spawns as needed)
