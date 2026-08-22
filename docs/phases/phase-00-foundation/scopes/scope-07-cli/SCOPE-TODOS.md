# Scope S07 — trace CLI

- [x] P00-S07-00 planner — 2026-08-15: locked stdlib argv (no cobra); `-C` root; commands init/index/reindex/add/link/transition/seed/why/context; seed JSON v1; wiring matrix to store/gitcli/analyzers/domain/retrieval/compiler; CGO for binary; exit 0/1/2; thickened 01-cli.md + light S08 Depends; no product Go; **P00-S07-01 ready**
- [x] P00-S07-01 implement — 2026-08-15: thin CLI adapter + seed v1 + tests; `CGO_ENABLED=1 go test ./...` + build `bin/trace` ok
- [x] P00-S07-02 review — 2026-08-15: APPROVE high; inline AtRev fallthrough fix; no spawns; next P00-S08-00
