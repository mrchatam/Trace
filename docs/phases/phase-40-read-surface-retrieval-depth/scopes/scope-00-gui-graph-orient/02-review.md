# P40-S00-02 — Review (G5 GUI graph orient)

## Metadata
- id: P40-S00-02
- todo_ids: [P40-S00-02]
- role: reviewer
- skills: [code-review-and-quality, accessibility-reviewer, web-design-guidelines]
- mcps: [user-trace]
- verification: mixed

## Objective

Fresh independent review of S00-01 G5 implementation vs REMEDIATION-PLAN G5, GAP-REGISTRY G-008, Law 19, and M-001 moat charter.

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [01-implement.md](01-implement.md) — G5-A1–A7 acceptance map
- [00-PLANNER.md](00-PLANNER.md) — locks + rejects
- [REMEDIATION-PLAN G5](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-06-remediation-plan/REMEDIATION-PLAN.md)
- [GAP-REGISTRY G-008](../../../phase-38-retrieval-context-peer-gaps/scopes/scope-04-gap-registry/GAP-REGISTRY.md)

## Session start

Follow agent-loop-protocol Session start. **Fresh subagent** — do not share implementer session.

## Locked defaults

| Item | Value |
|------|-------|
| APPROVE bar | Medium+ confidence; zero open blocker/high |
| Spawn trigger | Blocker/high → spawn 02a/02b pair below this row |
| Law 19 | Graph retrieval in Go library; `web/` calls HTTP only |
| Cap defaults | UI_CAP=100, SEED_MAX_NODES=40, EXPAND_MAX_NODES=50 unchanged |
| MCP/compiler | **Untouched** in S00 — no `trace_explore` |

## Review checklist

### A — G5 gap closure

- [ ] First-visit orient panel on `/` (G5-A1)
- [ ] Moat-first copy — task loop before graph-only (G5-A2)
- [ ] Dismiss persistence (G5-A3)
- [ ] Budget/confidence labels honest (G5-A4)
- [ ] Install hook narrative in CONTRIBUTING or install path (G5-A6)

### B — M-001 moat

- [ ] Orient does not claim graph replaces task loop / gates / review
- [ ] Links or copy steer to Tasks + Loop
- [ ] No query-only or graph-only product drift

### C — Laws 6–7

- [ ] No cap constant increases in `overviewCompose.ts`
- [ ] Truncation banners still present when budget exceeded
- [ ] No full-graph dump behavior added

### D — Law 19

- [ ] No new retrieval/business logic in `web/` beyond presentation
- [ ] If new HTTP route added: library function exists first + handler thin

### E — Regression

- [ ] Seed compose overview still loads
- [ ] Manual center, expand, Inspector, empty states work
- [ ] `overviewCompose.test.ts` green if run

### F — Rejects

- [ ] No Graphify port / static graph.html replacement
- [ ] No parallel SQLite from browser
- [ ] No compiler/MCP/G2 changes smuggled in

### G — Live verification commands

```bash
# Web unit tests (node:test — no vitest in web/)
cd web && node --experimental-strip-types --test src/lib/overviewCompose.test.ts
cd web && npm run build

# Law 19 — no httpapi fork for orient-only
grep -rn 'orient' internal/httpapi/ web/src/ || true

# Caps unchanged
grep -E 'UI_CAP|SEED_MAX_NODES|EXPAND_MAX_NODES|SEED_CAP|DEPTH' web/src/lib/overviewCompose.ts

# Truncation/budget UX preserved
grep -n 'graphTruncated\|graph-budget-line\|Laws 6–7' web/src/screens/Graph.tsx

# Moat copy spot-check (Explore nav label + task loop)
grep -iE 'loop|gate|task|Explore' web/src/components/GraphOrientPanel.tsx web/src/screens/Graph.tsx web/src/layout/Nav.tsx

# Install hook narrative
grep -iE 'trace serve|graph-first|Explore' CONTRIBUTING.md
```

## Exit criteria

- [ ] APPROVE or spawn with evidence (blocker/high findings)
- [ ] Law 19 adapter boundary verified
- [ ] Board row → `done` with verdict + confidence in Notes

## Next

`P40-S01-00`
