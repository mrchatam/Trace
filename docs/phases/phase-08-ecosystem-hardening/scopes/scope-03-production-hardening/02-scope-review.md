# P08 / S03 / 02 — Scope review (production hardening)

## Metadata
- id: P08-S03-02
- todo_ids: [P08-S03-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of S03 production-hardening deliverables vs planner locks. Spawn remediations on blocker/high. Forward-only.

## References
- [../../../../rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- Sibling [00-PLANNER.md](00-PLANNER.md) + [01-production-hardening.md](01-production-hardening.md)
- [docs/TODO.md](../../../../TODO.md)

## Session start
Agent → clarify → Plan → execute (review).

## Focus checklist

### Migrate hygiene
- [ ] Embed `schema/NNN_*.sql` + `schema_migrations` still the apply path on Open
- [ ] Forward-only (no downgrade CLI / no rewrite of applied migs)
- [ ] `trace migrate status` (or equivalent) reports applied/max consistently
- [ ] Prefer **no** `011_*`; if present, additive + justified in Notes (not analyzer API version)

### Backup / restore
- [ ] Snapshot is **`trace.db`** (not whole repo / not source trees)
- [ ] Taken while respecting exclusive `trace.lock`
- [ ] Restore into path-local `<absRoot>/.trace/`; rebinds `projects.root_path` to current Abs
- [ ] Round-trip test green; `HasBlobLikeColumns` remains false
- [ ] `access.token` excluded from backup by default

### Local auth
- [ ] Optional `.trace/access.token` + `TRACE_ACCESS_TOKEN` gate on Open
- [ ] Fail-closed when token present and env missing/wrong
- [ ] `trace auth set|clear|status` thin; status does not leak secret
- [ ] **No** cloud OAuth / hosted IdP / daemon auth theater
- [ ] MCP: no new tools; Open inherits gate (G19)

### S02 / bars
- [ ] Path-local bind unchanged; no shared-parent `.trace`
- [ ] Lock still fail-closed; backup/restore do not bypass lock
- [ ] Carry-forward: Gate H + honesty A/B/C + Gate G/E/F + ablation + p0x + x0 + Gate C `dry_run:false`
- [ ] Upcoming S04 VERIFY stubs still compatible (checklist must cover migrate/backup/auth)

## Verify commands (re-run)

```bash
CGO_ENABLED=0 go test ./internal/store/... ./internal/vcs/... ./internal/gitcli/... ./internal/domain/... ./internal/mcp/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... -count=1
CGO_ENABLED=1 go test ./cmd/trace/... ./internal/analyzers/... -count=1
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./evals/honesty/... ./evals/replan/... ./evals/impact/... ./evals/capability/... ./evals/perf/... -count=1
CGO_ENABLED=1 go test ./... -count=1
```

## Role work
1. Diff claims (01 Notes + locks) vs repo evidence.
2. Findings by severity; blocker/high → inline fix or spawn `02a`/`02b`.
3. Light-thicken **upcoming** S04 Depends if checklist items need naming (do not rewrite `done` history).
4. Write [REVIEW-NOTES.md](REVIEW-NOTES.md); board status + Notes.

## Exit criteria
- [ ] APPROVE (medium/high with residuals listed) or spawn with evidence
- [ ] No open blocker/high without pending follow-up
- [ ] Board status + Notes; next **P08-S04-00** when APPROVE

## Out of scope
- Owning `evals/compat` harness creation (S04)
- Rewriting S01/S02 done prompts
- Cloud auth / daemon promotion

## Minimal todos
- [ ] Diff claims vs repo; re-run locked tests; write REVIEW-NOTES
- [ ] APPROVE or spawn; update board Notes
