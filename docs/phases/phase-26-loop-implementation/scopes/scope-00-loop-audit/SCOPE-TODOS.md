# Scope 00 todos — loop audit

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P26-S00-00 | planner | done | 2026-08-20: Locked `01-loop-audit.md` — corrected paths (`internal/domain/` not `internal/seed/`), INT table, threshold-options-only rule, AUDIT.md template, preflight. Verified live: `internal/loop/`, `internal/store/deliberation.go`, `internal/domain/seed_{export,import}.go`, `cmd/trace/loop.go`, `internal/mcp/`, `internal/install/enforcement.go`. Confirmed P25-2: `ParentOrchestratorRule` L26 defined, unused in `cursorRulesMDCContent`/`claudeFallbackRulesContent`. Next: P26-S00-01. |
| P26-S00-01 | implementer | pending | Write `AUDIT.md` (no product code). Prompt: `01-loop-audit.md` |
