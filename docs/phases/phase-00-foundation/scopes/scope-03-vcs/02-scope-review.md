# P00 / S03 / 02 — Scope review (VCS adapter + git CLI)

## Metadata
- id: P00-S03-02
- todo_ids: [P00-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- agents: [code-reviewer]
- verification: automated

## Objective
Independent review of S03 (VCS adapter + git CLI). Severity-tag findings; small fixes or spawn `a` implement + `b` review pairs with **full** prompts. Forward-only. May thicken **upcoming** prompts if this scope’s surface drifted.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md) (Reviewer loop — write proper spawn prompts)
- Sibling `01-*.md` + board Notes

## Session start
Agent → clarify → Plan → review (fresh subagent).

## Review focus
- Exit criteria honestly met vs thickened `01-vcs.md`?
- DR-GIT: git CLI only — no go-git/libgit2 in `go.mod` / imports
- G1: no source blobs / permanent full diffs in SQLite; content == `git show`
- DR-INCREMENTAL: Refresh watermark — second refresh without new commits is noop; no full history rewrite
- Interface + ≥1 fake/mock; `gitcli` is the only production backend
- Silent failures on git non-zero exit; missing temp-repo tests
- Cross-scope: S04 expects path/OID via VCS; store `files.git_oid` optional only

## Exit criteria
- [ ] Findings recorded
- [ ] blocker/high fixed or spawned
- [ ] Confidence medium or high (residuals listed if medium)
- [ ] TODO.md updated

## Minimal todos
- [ ] Compare claims vs evidence
- [ ] Fix or spawn
- [ ] Re-verify
- [ ] Board update
