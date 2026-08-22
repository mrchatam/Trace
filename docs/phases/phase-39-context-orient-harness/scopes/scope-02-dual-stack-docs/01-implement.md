# P39-S02-01 — Implement G4 dual-stack docs

## Metadata
- id: P39-S02-01
- todo_ids: [P39-S02-01]
- role: implementer
- skills: [documentation-and-adrs, writing-for-agents, writing-guidelines]
- verification: automated

## Objective

Author **G4** dual-stack documentation: Trace + Codegraph as **complementary optional stacks** in `CONTRIBUTING.md` and `AGENTS.md` (G-011). **Doc-only — no product code.**

Full recipe lives here; S01 already shipped moat-first MCP orient + pointer-only Codegraph stub — **do not duplicate** (see coordination table).

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](../../../../rules/project-rules.md) — Law 19
- [00-PLANNER.md](00-PLANNER.md) — checklist G4-D1–D8
- [REMEDIATION-PLAN G4](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [PEER-CG §5](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md)
- [PEER-FIXTURES.md](../../../phase-38-retrieval-context-peer-gaps/PEER-FIXTURES.md)

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute).

## Locked defaults

| Item | Value |
|------|-------|
| GAP id | G-011 |
| Product code | **Forbidden** |
| Primary edits | `CONTRIBUTING.md`, `AGENTS.md` |
| Optional | `README.md` one-paragraph pointer (only if CONTRIBUTING section exceeds ~80 lines) |
| Tone | Complementary optional stack — Trace moat primary for task/evidence work |
| Must cover | G4-D1–D8 (checklist below) |
| Must not | Mandatory dual-index; bundled MCP; Trace core reads/writes `.codegraph/`; re-litigate moat-first/reload prose |

## S01 coordination — do not duplicate

S01 shipped these; S02 **extends**, does not rewrite:

| Location | S01 owns (keep as-is) | S02 adds |
|----------|----------------------|----------|
| `CONTRIBUTING.md` **Agent workflow (local)** | Moat-first orient (`:68–69`); harness vs enforcement (`:70`); MCP reload + 9/16 + offline rebuild (`:72`) | New **`## Trace + Codegraph (optional dual-stack)`** section after Agent workflow block |
| `CONTRIBUTING.md` line 72 pointer | "see the **dual-stack** section (Phase 39 S02…)" | Replace placeholder with markdown link to new section anchor (e.g. `#trace--codegraph-optional-dual-stack`) |
| `internal/mcp/instructions.go` | Codegraph complement pointer-only (`:23–26`) | **Do not edit** — already points to CONTRIBUTING dual-stack |
| `AGENTS.md` | Agent workflow list (`:13–20`); orchestrator snippet | New **`### Optional Codegraph (complement)`** under Agent workflow — 3–5 bullets + link to CONTRIBUTING |

**Explicit non-touch:** `internal/`, `cmd/`, `web/`, `trace/` product paths; `instructions.go` (S01 complete).

## G4-D1–D8 acceptance checklist (implementer self-check)

| ID | Requirement | Implementation hint |
|----|-------------|---------------------|
| G4-D1 | Section title identifies Trace + Codegraph as **complementary**, not merged product | Heading like `Trace + Codegraph (optional dual-stack)`; opening sentence: two local-first tools, different jobs |
| G4-D2 | Storage: `.trace/` (task/plan/evidence) vs `.codegraph/` (symbol graph) — separate indexes | Table or bullets: `.trace/trace.db` vs `.codegraph/`; both gitignored; no shared SQLite |
| G4-D3 | When to use Trace: task loop, gates, plan tree, evidence, progressive task packet | Bullet list: `trace_tasks`, `trace_context`, `trace_loop`, gate, `trace_review`, `trace_plan`, portable graph — moat-first |
| G4-D4 | When to use Codegraph: symbol exploration, call paths, blast radius — **optional** per repo | Bullet list: `codegraph init`, `codegraph_explore`, per-project index; "optional — only when symbol graph helps" |
| G4-D5 | Law 19: each stack is adapter/MCP over its own store; Trace core does not index into `.codegraph/` | Explicit boundary: Trace CLI/MCP never opens `.codegraph/`; Codegraph MCP is separate harness config |
| G4-D6 | Setup recipe: `trace index` / install path + optional `codegraph init` — neither required for the other | Numbered setup: (1) Trace path, (2) optional CG path; order independent |
| G4-D7 | Reject language explicit: no default dual-index, no bundled dual MCP in Trace product | Short **Not shipping** subsection: reject mandatory dual-index, bundled trace+CG MCP, Trace indexing CG data |
| G4-D8 | Link to Phase 38 PEER-CG complement note / PEER-FIXTURES | Relative links to `docs/phases/phase-38-retrieval-context-peer-gaps/scopes/scope-02-codegraph-peer/PEER-CG.md` §5 and `PEER-FIXTURES.md` |

## Suggested section structure (CONTRIBUTING)

Insert **after** the Agent workflow block (after `:72`, before `## Review`):

```markdown
## Trace + Codegraph (optional dual-stack)

<opening: complementary, not merged — satisfies G4-D1>

### Storage (separate indexes)
<G4-D2>

### When to use Trace
<G4-D3 — cross-link moat-first section above, don't repeat full playbook>

### When to use Codegraph (optional)
<G4-D4>

### Setup
<G4-D6>

### Law 19 — adapter boundaries
<G4-D5>

### Not shipping (product rejects)
<G4-D7>

### Investigation context
<G4-D8 links>
```

**AGENTS.md** — after item 5 under Agent workflow (`:20`), add:

```markdown
### Optional Codegraph (complement)

Symbol-level code exploration is **optional** via separate Codegraph MCP — Trace owns task loop + evidence. See [CONTRIBUTING — Trace + Codegraph](CONTRIBUTING.md#trace--codegraph-optional-dual-stack). Orchestrator/harness may register both MCP servers; user's choice per repo.
```

## Touch-list

| Step | File | Action |
|------|------|--------|
| 1 | `CONTRIBUTING.md` | Add `## Trace + Codegraph (optional dual-stack)` (G4-D1–D8); update `:72` pointer to section anchor |
| 2 | `AGENTS.md` | Add `### Optional Codegraph (complement)` under Agent workflow |
| 3 | `README.md` | Optional one-line cross-link under agent/contributing area — **only** if dual-stack section is long |

## Implementation order

```text
1. Draft CONTRIBUTING dual-stack section (G4-D1–D8)
2. Wire line-72 pointer to section anchor (replace Phase 39 S02 placeholder)
3. Add AGENTS.md complement subsection (pointer + moat reminder)
4. Optional README one-liner
5. Self-check G4-D1–D8 table + confirm no product files touched
```

## Live verification commands (doc-only)

```bash
# Confirm no product diff
git diff --name-only | grep -E '^(internal/|cmd/|web/)' && echo FAIL || echo OK-no-product

# G4-D1 title/complement language
grep -n 'Trace + Codegraph\|complement\|dual-stack' CONTRIBUTING.md AGENTS.md

# G4-D2 storage separation
grep -n '\.trace/\|\.codegraph/' CONTRIBUTING.md

# G4-D5 Law 19 / no cross-index
grep -n 'Law 19\|does not index\|adapter' CONTRIBUTING.md

# G4-D7 reject language
grep -n 'Not shipping\|bundled\|mandatory dual' CONTRIBUTING.md

# G4-D8 investigation links
grep -n 'PEER-CG\|PEER-FIXTURES' CONTRIBUTING.md

# S01 pointer wired (no stale placeholder)
grep -n 'Phase 39 S02' CONTRIBUTING.md   # expect 0 after implement (anchor link replaces)
grep -n 'dual-stack' CONTRIBUTING.md internal/mcp/instructions.go
```

## Minimal todos

- [ ] Add CONTRIBUTING dual-stack section satisfying G4-D1–D8
- [ ] Update Agent workflow pointer (`:72`) to section anchor
- [ ] Add AGENTS.md optional complement subsection
- [ ] Self-check: no product files in diff
- [ ] Board row → `done` with G4-D1–D8 pass/fail in Notes

## Exit criteria

- [ ] G4-D1–D8 checklist satisfied (self-check in Notes)
- [ ] No product code files changed
- [ ] S01 moat-first/reload prose not duplicated (coordination table respected)
- [ ] Board row → `done` with Notes listing doc paths + checklist summary

## Next

`P39-S02-02`
