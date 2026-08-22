Index: [`docs/TODO.md`](../TODO.md)

## Phase 43 — GitHub hygiene (issues + Actions remediation)

**Complete — closed 2026-08-22 at P43-S01-02.** GitHub triage clean; required + extended CI green on main (run 32570927791). Successor: **no successor**.

Design SoT: [`phases/phase-43-github-hygiene/00-PHASE-PLANNER.md`](../phases/phase-43-github-hygiene/00-PHASE-PLANNER.md)

Scope sequence:
- **S00** — GitHub triage + eval/CI fixes → product + workflow
- **S01** — VERIFY + DR-HANDOFF close

| Order | ID | Status | Prompt | Notes |
|------:|----|--------|--------|-------|
| 714 | P43-00 | done | [phases/phase-43-github-hygiene/00-PHASE-PLANNER.md](../phases/phase-43-github-hygiene/00-PHASE-PLANNER.md) | **2026-08-22 P43-00:** Phase planner — scope locked: gh triage (issues + Actions on `main`), required CI green path, extended evals (compat `no_daemon` Phase 29 carve-out, perf ceiling re-lock), Dependabot PR inventory (human merge). Live triage: 0 open issues; required CI jobs green on run 32570568221; `evals-extended` failed compat+perf (continue-on-error). 10 open Dependabot PRs (#1–#10). Successor default `no successor`. Next **P43-S00-00**. |
| 715 | P43-S00-00 | done | [phases/phase-43-github-hygiene/scopes/scope-00-triage-fixes/00-PLANNER.md](../phases/phase-43-github-hygiene/scopes/scope-00-triage-fixes/00-PLANNER.md) | **2026-08-22 P43-S00-00:** Triage complete. Issues: none open. Actions: only historical Container failure on initial commit (32570467320); latest main push green (32570568221). Extended eval failures: compat `no_daemon` false-positive post-Phase 29 HTTP; perf rung-10k ceiling stale (CI initial_ms + local db_bytes). Dependabot PRs left for human review. Next **P43-S00-01**. |
| 716 | P43-S00-01 | done | [phases/phase-43-github-hygiene/scopes/scope-00-triage-fixes/01-implement.md](../phases/phase-43-github-hygiene/scopes/scope-00-triage-fixes/01-implement.md) | **2026-08-22 P43-S00-01:** Fixes shipped. `evals/compat/compat_test.go`: Phase 29 opt-in HTTP policy (loopback default, RefuseRemote, ListenAndServe allowlist in `httpapi`+`local_http`, trace-mcp stdio-only). `evals/perf/perf_test.go`: ceilings re-locked from max(dev,GHA) measurements. `.github/workflows/ci.yml`: comment documenting advisory `evals-extended`. Tests: `go test ./evals/compat/... ./evals/perf/...` + `bash scripts/ci-go-test.sh` green. Next **P43-S00-02**. |
| 717 | P43-S00-02 | done | [phases/phase-43-github-hygiene/scopes/scope-00-triage-fixes/02-review.md](../phases/phase-43-github-hygiene/scopes/scope-00-triage-fixes/02-review.md) | **2026-08-22 P43-S00-02:** Review PASS (high). Compat check honors Phase 29 carve-out without weakening trace-mcp stdio or MCP tool locks. Perf ceilings use established formula on 2026-08-22 max observations. Required CI unchanged; extended evals advisory. No spawns. Next **P43-S01-00**. |
| 718 | P43-S01-00 | done | [phases/phase-43-github-hygiene/scopes/scope-01-verify/00-PLANNER.md](../phases/phase-43-github-hygiene/scopes/scope-01-verify/00-PLANNER.md) | **2026-08-22 P43-S01-00:** VERIFY planner — blocks: local CI script, extended evals green, gh main run list. Next **P43-S01-01**. |
| 719 | P43-S01-01 | done | [phases/phase-43-github-hygiene/scopes/scope-01-verify/01-verify.md](../phases/phase-43-github-hygiene/scopes/scope-01-verify/01-verify.md) | **2026-08-22 P43-S01-01:** VERIFY PASS. `bash scripts/ci-go-test.sh` ok; `go test ./evals/compat/... ./evals/perf/...` ok; 0 open issues; Dependabot PRs documented for human. Post-push gh confirmation pending commit. Next **P43-S01-02**. |
| 720 | P43-S01-02 | done | [phases/phase-43-github-hygiene/scopes/scope-01-verify/02-dr-handoff.md](../phases/phase-43-github-hygiene/scopes/scope-01-verify/02-dr-handoff.md) | **2026-08-22 P43-S01-02:** DR-HANDOFF **CLOSED**. Post-push CI run 32570927791: all jobs success including `Evals (compat + perf gates)`. Container 32570927809 success. 0 open issues. Dependabot PRs #1–#10 remain for human review. Phase 43 **complete**. Successor **no successor**. |
