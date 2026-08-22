# Scope 01 — board map

**S01 Trace live audit** — serial: **P38-S01-00 → P38-S01-01 → P38-S01-02**. Artifact: `TRACE-AUDIT.md`. **Investigate only — no product code.**

| Order | Board ID | Prompt | Role | Artifact |
|------:|----------|--------|------|----------|
| 650 | P38-S01-00 | [00-PLANNER.md](00-PLANNER.md) | Scope planner | Lock defaults; thicken 01/02 — **done 2026-08-22** |
| 651 | P38-S01-01 | [01-investigate.md](01-investigate.md) | Investigate | Author `TRACE-AUDIT.md` |
| 652 | P38-S01-02 | [02-review.md](02-review.md) | Reviewer | APPROVE / REQUEST CHANGES / SPAWN |

## Planner locks (P38-S01-00)

| Lock | Value |
|------|-------|
| Hypotheses | H2, H3, H5, H6, H9, H10 + H1 partial, H8 partial |
| Method | Live CLI/MCP + file:line (not docs-only) |
| Evidence | `experiments/runs/YYYY-MM-DD-p38-s01-651/evidence/` |
| Verdicts | confirmed gap \| not a gap \| inconclusive |
| Dogfood | Trace repo root; read-only `.trace/` |
| Spawn | Unbounded slice → S01-01a row before S02 |
| Non-goals | No product code; no REMEDIATION-PLAN; no ranked build |

## Hypothesis → investigation todo map

| H | Todo(s) | Primary evidence files |
|---|---------|------------------------|
| H2 | T1 | `h2-compiler-fts.txt`, `h2-context-packet.json` |
| H3 | T2 | `h3-layers-packet-depth2.json`, `h3-layers-designed-vs-shipped.md` |
| H5 | T4 | `h5-index-status.json`, `h5-index-langs.txt` |
| H6 | T5 | `h6-mcp-tool-list.txt`, `h6-mcp-surface.md` |
| H9 | T3 | `h9-intent-grep.txt`, `h9-intent-pipeline.md` |
| H10 | T6 | `h10-install-detect.json`, `h10-install-moat.md` |
| H1 (partial) | T7 | `h1-trace-partial.md`, `h1-trace-context-depth2.json` |
| H8 (partial) | T8 | `h8-gui-partial.md` (+ optional screenshot) |
| All | T0, T9 | Preflight + TRACE-AUDIT.md synthesis |

## Optional tools (S01)

| H | Trace CLI | Trace MCP | Codegraph MCP | Notes |
|---|-----------|-----------|---------------|-------|
| H1 partial | ✓ | ✓ | — | Trace-only; CG/UA in S02/S03 |
| H2 | ✓ | ✓ | — | UA read optional for contrast cite |
| H3 | ✓ | ✓ | — | |
| H5 | ✓ | — | defer watch compare to S02 | |
| H6 | ✓ | ✓ | defer CG 1-tool to S02 | |
| H8 partial | ✓ (serve) | — | — | GUI optional |
| H9 | ✓ | ✓ | — | |
| H10 | ✓ | — | — | |

## Live command minimum (review gate)

Reviewer expects evidence for at least:

```bash
trace context <task-id> --format json
trace search "<query>"
trace index status
trace install detect
go test ./internal/mcp/ -run TestToolNamesRegistered -count=1
# MCP: trace_version, trace_context, trace_search (or CLI equivalents)
```

## Out of scope (S01)

- H4, H7, H11 (other scopes)
- Full H1/H8 peer matrix (S02/S03/S04)
- PEER-CG.md, PEER-UA-GF.md, GAP-REGISTRY, REMEDIATION-PLAN
- Any Go/TS product change

## Spawn rule

If T1–T8 reveal an unbounded slice (e.g. full HTTP API audit), reviewer inserts **P38-S01-01a/01b** below row 652 — do not block S02 indefinitely without S05 note.
