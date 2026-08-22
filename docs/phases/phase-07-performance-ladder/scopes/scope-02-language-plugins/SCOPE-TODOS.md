# Scope S02 — Language plugins / adapters

**Depends-on:** S01 APPROVE (`P07-S01-02` done 2026-08-16). **Live T0:** `cmd/trace` `isT0SkipDir`/`isT0SkipPath`/`walkIndexable` — S02 must not regress when extending `DetectLanguage`.

| ID | Role | Status | Notes |
|----|------|--------|-------|
| P07-S02-00 | planner | done | 2026-08-16: locked **Go** + `tree-sitter-go` **v0.25.0**; DetectLanguage/`extract` switches; golden required; no Gate H invent |
| P07-S02-01 | implement | done | 2026-08-16: Go adapter + golden; see board Notes |
| P07-S02-02 | review | done | 2026-08-16: APPROVE high — [REVIEW-NOTES.md](REVIEW-NOTES.md); next P07-S03-00 |

## Checklist

- [x] P07-S02-00 planner
- [x] P07-S02-01 implement
- [x] P07-S02-02 review

## Phase locks (FINAL — P07-S02-00)

| Item | Value |
|------|-------|
| Language | **Go** (exactly one new adapter) |
| Grammar | `github.com/tree-sitter/tree-sitter-go` **v0.25.0** (`bindings/go`) |
| Package | `internal/analyzers` — switch extension, not plugin registry |
| CGO | Analyzers-only |
| S01 | Consume T0 walk; do not regress |
| Gate H | No threshold invent in S02 |
| Carry-forward | Honesty; Gates E/F/G; capability ablation; p0x; x0; Gate C intact |
