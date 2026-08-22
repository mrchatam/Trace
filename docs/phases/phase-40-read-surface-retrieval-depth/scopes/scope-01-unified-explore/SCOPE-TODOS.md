# SCOPE-TODOS — G2 unified explore

| ID | Item | Status |
|----|------|--------|
| G2-0 | Scope planner — live re-verify 16→17 touch-list; thicken 01/02 | done (P40-S01-00) |
| G2-1 | Law spike gate — waived; desk-check locked at P40-00 / P40-S01-00 | done (P40-S01-00) |
| G2-2 | Add-tool decision — `trace_explore` as 17th read-only MCP tool | done (P40-00) |
| G2-3 | Library Explore compose spec + caps (`internal/retrieval/explore.go`) | pending (S01-01) |
| G2-4 | MCP `tools_explore.go` + server registration + `RegisteredToolNames` 17 | pending (S01-01) |
| G2-5 | Instructions addendum + stale hygiene 9/16→9/17 | pending (S01-01) |
| G2-6 | Tests G2-T1–T7 + MCP mirrors + tool-count migration | pending (S01-01) |
| G2-7 | Review moat + 17-tool + Law 19 | pending (S01-02) |

**G1 prerequisite shipped** — `internal/compiler/compiler.go:158–165` (`ContextOptions.Query` merge).

**Live baseline (P40-S01-00):** 16 `AddTool` calls (`server.go:59–226`); no `trace_explore`; `TestToolNamesRegistered` + `TestRegisteredToolNames_IncludesTracePlan` assert 16.

**Law spike desk-check (waived — implement must still satisfy):**

| Check | Expected |
|-------|----------|
| Task required | Empty `task_id` → validation error |
| Query optional | Merges via G1 path when set |
| Caps honest | Response includes truncation/budget fields |
| Write surface visible | 16 existing tools + explore = 17; loop/review/transition unchanged |
| Not CG equivalent | Response shape differs; task packet section required |
| Compose-first preserved | Instructions still rank manual compose for fine control |
