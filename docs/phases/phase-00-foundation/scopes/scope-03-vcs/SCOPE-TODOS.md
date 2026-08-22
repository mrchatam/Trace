# Scope S03 — VCS adapter + git CLI

- [x] P00-S03-00 planner — 2026-08-15: locks DR-GIT (CLI + `internal/vcs`/`gitcli`); no libgit2/go-git; thin commit/path index in `.trace/trace.db`; no blobs; thickened `01-vcs.md`
- [x] P00-S03-01 implement — 2026-08-15: `vcs.Repository` + Fake; `gitcli` CLI backend; mig `002_vcs_index.sql`; incremental Refresh
- [x] P00-S03-02 review — 2026-08-15: APPROVE high; History watermark==HEAD gate fixed inline; no spawns — see REVIEW-NOTES.md
