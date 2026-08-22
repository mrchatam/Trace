# Scope 04 todos — installer fix

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P26-S04-00 | planner | done | Confirmed P25-2 gap via `rg` in `internal/install/enforcement.go`: `ParentOrchestratorRule` defined; `cursorRulesMDCContent` and `claudeFallbackRulesContent` present with no usage hit. Locked `01-implement.md` and `02-review.md` defaults/checklist for S04 wiring + generated cursor rules assertion. |
| P26-S04-01 | implementer | pending | Wire ParentOrchestratorRule + test (do not start from planner row). |
| P26-S04-02 | reviewer | pending | Independent review. |
