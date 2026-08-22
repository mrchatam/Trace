# P39-S02-02 — Review G4 dual-stack docs

## Metadata
- id: P39-S02-02
- todo_ids: [P39-S02-02]
- role: reviewer
- skills: [documentation-and-adrs, writing-for-agents, writing-guidelines]
- verification: mixed

## Objective

Fresh independent review of S02-01 docs vs G-011 gap closure, H11 doc-only lock, Law 19 accuracy, and S01 coordination (no moat-first duplication).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — G4-D1–D8 + coordination table
- [00-PLANNER.md](00-PLANNER.md)
- [REMEDIATION-PLAN G4](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [PEER-CG §5](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md)

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high without trivial fix → spawn 02a/02b pair below this row |
| Product code | **Must be absent** from diff |
| S01 boundary | Moat-first/reload/9/16 stays in Agent workflow block; dual-stack is additive section |
| Rejects | Mandatory dual-index; bundled MCP; Trace core indexes `.codegraph/` |

## Review checklist

### A — G4-D1–D8 (per-item)

- [ ] **G4-D1** — CONTRIBUTING section title/lead states Trace + Codegraph are **complementary**, not one merged product
- [ ] **G4-D2** — `.trace/` (task/plan/evidence) vs `.codegraph/` (symbol graph) documented as **separate** indexes; both gitignored
- [ ] **G4-D3** — Trace use-cases: task loop, gates, plan tree, evidence, progressive task packet (moat primary)
- [ ] **G4-D4** — Codegraph use-cases: symbol exploration, call paths, blast radius; framed **optional** per repo
- [ ] **G4-D5** — Law 19 accurate: adapter/MCP per store; Trace core does **not** index into `.codegraph/`
- [ ] **G4-D6** — Setup: `trace index`/install + optional `codegraph init`; neither required for the other; order independent
- [ ] **G4-D7** — Explicit reject: no default dual-index, no bundled dual MCP in Trace product
- [ ] **G4-D8** — Valid relative links to PEER-CG §5 and PEER-FIXTURES (paths exist)

### B — H11 lock

- [ ] Doc-only — diff limited to `CONTRIBUTING.md`, `AGENTS.md`, optional `README.md`
- [ ] No mandatory dual-index language ("must run both", "required dual-index")
- [ ] No bundled MCP claim ("Trace ships Codegraph", "single MCP for both")

### C — M-001 moat

- [ ] Trace task moat presented as **primary** workflow for directed work + evidence
- [ ] Codegraph framed as optional complement for code-graph reads — not replacement for task loop

### D — S01 coordination (non-duplication)

- [ ] Agent workflow block (`CONTRIBUTING.md:64–72`) moat-first + reload prose **unchanged in substance** (pointer may update to anchor)
- [ ] No second full moat-first playbook in dual-stack section (cross-link OK)
- [ ] `internal/mcp/instructions.go` **not edited** in S02 diff (S01 pointer-only stub sufficient)
- [ ] Line `:72` placeholder "Phase 39 S02" replaced with working section link

### E — Law 19

- [ ] Adapter/storage boundaries accurate vs `docs/rules/project-rules.md`
- [ ] No false claim Trace core indexes or merges Codegraph data
- [ ] No `.codegraph/` writes from Trace CLI/MCP implied

### F — AGENTS.md

- [ ] Optional complement subsection present under Agent workflow
- [ ] Points to CONTRIBUTING dual-stack section (not standalone duplicate recipe)

### G — Links

- [ ] PEER-CG path resolves: `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md`
- [ ] PEER-FIXTURES path resolves: `docs/phases/phase-38-retrieval-context-peer-gaps/PEER-FIXTURES.md`

## Live verification commands

```bash
# Doc-only diff
git diff --name-only
git diff --name-only | grep -E '^(internal/|cmd/|web/)' && echo FAIL-product-touch || echo OK-doc-only

# Per-checklist greps
grep -n 'Trace + Codegraph\|complement\|dual-stack' CONTRIBUTING.md AGENTS.md
grep -n '\.trace/\|\.codegraph/' CONTRIBUTING.md
grep -n 'Law 19\|does not index\|adapter' CONTRIBUTING.md
grep -n 'Not shipping\|bundled\|mandatory dual\|must run both' CONTRIBUTING.md
grep -n 'PEER-CG\|PEER-FIXTURES' CONTRIBUTING.md

# S01 coordination
grep -n 'Moat-first orient\|9/16\|trace_version' CONTRIBUTING.md
grep -n 'Phase 39 S02' CONTRIBUTING.md   # expect 0 post-S02-01
git diff internal/mcp/instructions.go 2>/dev/null | wc -l   # expect 0

# Link targets exist
test -f docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md && echo PEER-CG-ok
test -f docs/phases/phase-38-retrieval-context-peer-gaps/PEER-FIXTURES.md && echo PEER-FIXTURES-ok
```

## Exit criteria

- [ ] APPROVE (medium+ confidence) or spawn repair pair
- [ ] All G4-D1–D8 verified with evidence in Notes
- [ ] Board row → `done` with verdict + findings counts in Notes

## Next

`P39-S03-00`
