# P33-S02-02 — Launch review

## Metadata
- id: P33-S02-02
- todo_ids: [P33-S02-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of **`trace gui` + PATH teach** against DESIGN-LOCKS Theme C, Law 19, RESEARCH rejects (UA auto-port; PATH ≠ agents-install), and S01 land-on-Explore-`/`. Spawn `P33-S02-02a`/`02b` if blocker/high. Do **not** flip full `gui-quickstart` primary story (S05) unless trivial help-only. No Explore UI scope creep into S03.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md) — Theme C; PATH vs `trace install`
- [01-implement.md](01-implement.md) — locked defaults + test table
- [`../scope-00-research/RESEARCH.md`](../scope-00-research/RESEARCH.md)
- [`../scope-01-design-ux/UX-IA.md`](../scope-01-design-ux/UX-IA.md) — route clarity
- Live evidence: `cmd/trace/` (`root.go`, `gui`/`serve`, help, tests), `internal/httpapi/`, spot `docs/gui-quickstart.md`

## Session start

**Fresh subagent.** Follow agent-loop-protocol Session start.

## Checklist

### Theme C / CLI

- [ ] `gui` subcommand present in `root.go`; help accurate (`--no-open`, inherit serve flags)
- [ ] **No** requirement for global `-gui` this scope (subcommand primary)
- [ ] `serve` still works for scripting; gui does not delete/break serve
- [ ] Reuses `httpapi` listen path (Law 19); no parallel business logic in CLI

### Bind / port / security

- [ ] Default loopback `127.0.0.1:7432`; no public bind default
- [ ] Non-loopback refused without `--allow-remote` (+ token posture matches serve)
- [ ] Port conflict: friendly fail + user `--addr` — **no** auto-increment / silent hop (reject UA auto-port)
- [ ] Not an always-on daemon

### Browser open

- [ ] Opens best-effort after successful listen
- [ ] URL printed always
- [ ] Open failure → tip; listen success still OK (exit semantics)
- [ ] `--no-open` skips open (CI/headless)
- [ ] Opened URL targets Explore **`/`** — **not** `/overview` (S01 UX-IA / Theme B)

### PATH ≠ agents install (required)

- [ ] PATH teach is **`go install github.com/mrchatam/Trace/cmd/trace@…`** (help and/or minimal docs)
- [ ] Contributor build/symlink is secondary only
- [ ] **`trace install …` not documented or coded as PATH placement** (still agents/MCP/hooks)
- [ ] No overload of `cmdInstall` for binary PATH

### Docs / scope boundaries

- [ ] Full quickstart primary flip **not** required here (S05); trivial help/`gui` tip OK
- [ ] No Explore seed/graph product work (S03); no shell colorize (S04)

### Tests / evidence

- [ ] Tests cover: help; remote refuse; addr-in-use; `--no-open` / mocked open URL `/`; no auto-port
- [ ] Board Notes on P33-S02-01 have runnable evidence
- [ ] Confidence **medium** or **high** with residual risks listed if medium

### Findings disposition

- [ ] blocker/high: small inline fix **or** spawn `P33-S02-02a` (implement) + `P33-S02-02b` (review) immediately below this row
- [ ] Do not rewrite S00/S01 `done` history; thicken **upcoming** S03+ only if launch gaps affect them

## Exit criteria

- [ ] No open blocker/high without pending follow-up
- [ ] Confidence medium/high with evidence
- [ ] Next runnable **P33-S03-00** (or spawned S02 fix pair if any)

## Todo updates

Status + notes on **P33-S02-02** only (plus spawn rows if created).
