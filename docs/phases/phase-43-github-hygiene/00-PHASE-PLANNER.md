# Phase 43 — GitHub hygiene (issues + Actions remediation)

**Status:** Active (2026-08-22)  
**Predecessor:** Phase 42 complete (G1–G9 remediation)  
**Successor default:** `no successor`

## Goal

Triage and remediate GitHub-facing hygiene for `mrchatam/Trace`:

1. Open issues — fix or close with evidence
2. GitHub Actions on `main` — required jobs green; extended evals honest
3. Dependabot PR inventory — document for human merge (no blind merge)

## Locks

| Lock | Value |
|------|-------|
| Required CI | `go-test`, `web`, `embed-gui`, `openapi` must pass on `main` |
| Extended evals | `evals-extended` advisory (`continue-on-error: true`) but should pass after fixes |
| Phase 29 | Compat `no_daemon` must reflect opt-in HTTP carve-out, not pre-P29 blanket ban |
| Dependabot | Inventory only — human reviews major bumps (typescript 7, actions v7, sqlite 1.57) |

## Scopes

| Scope | Deliverable |
|-------|-------------|
| S00 | Triage + implement compat/perf/CI policy fixes |
| S01 | VERIFY + DR-HANDOFF |

## Board

[`docs/TODO/phase-43.md`](../TODO/phase-43.md)
