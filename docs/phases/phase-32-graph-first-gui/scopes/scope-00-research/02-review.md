# P32-S00-02 — Research review

## Metadata
- id: P32-S00-02
- todo_ids: [P32-S00-02]
- role: reviewer
- skills: [code-review-and-quality, planning-and-task-breakdown]
- mcps: []
- verification: automated
- hooks: []

## Objective

Independent review of `RESEARCH.md` against DESIGN-LOCKS, OPEN-PORT-MULTI, and live repo facts. Fix trivial RESEARCH doc gaps inline; spawn or thicken **upcoming** prompts (S01+) if structural. **No product code.**

## References

- [docs/rules/agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [DESIGN-LOCKS.md](../../DESIGN-LOCKS.md)
- [OPEN-PORT-MULTI.md](../../OPEN-PORT-MULTI.md)
- [01-research.md](01-research.md)
- Artifact: [RESEARCH.md](RESEARCH.md) (must exist after S00-01)
- Spot-check live: `web/src/App.tsx`, `web/src/screens/Graph.tsx`, `web/src/api/ops.ts`, `api/openapi.yaml`, `internal/httpapi/bind.go`

## Session start

Follow agent-loop-protocol Session start (Agent → clarify → Plan → execute). Fresh context — do not share the S00-01 implementer session.

## Checklist

### Locks & laws
- [ ] Peer bar matches locks (between Graphify and UA; Trace plan/task/decision semantics kept)
- [ ] Depth-before-visual ordering preserved (S03 depth list before S04 craft; no “visual first” sneak)
- [ ] Laws 6–7 / 19 called out (budgeted neighborhood; adapters only; no full-graph dump default)
- [ ] No silent reopen of 3D / second SPA / hosted SaaS / public bind defaults / MCP `/rpc` in browser

### Artifact quality
- [ ] All RESEARCH template headings present and non-empty
- [ ] `web/` inventory matches live routes (Overview index; Graph is `/graph`, not home)
- [ ] Graph inspector gap called out (Graph has no rich inspector; why/context only on TaskDetail today)
- [ ] API reuse map cites real `/v1` ops; `getImpact` OpenAPI vs `ops.ts` wrapper gap accurate if still true
- [ ] Borrow vs reject peer rows are actionable for S01 (not marketing fluff)

### P32-PORT (required)
- [ ] Section present; light review confirmed against OPEN-PORT-MULTI + live Listen
- [ ] Actionable S02 recommendation — **minimum prefer #1** (friendly `EADDRINUSE` / in-use message + `--addr` guidance)
- [ ] States **S02 always owns ship** even if API is `NO-GAPS.md`
- [ ] Peer port pattern (e.g. UA auto-increment) cited without forcing #2 as sole story
- [ ] Does not treat docs-only (#3) as complete S02 closure without #1 (or explicit discouraged deferral)

### Forward board rights
- [ ] If RESEARCH is PASS: next is **P32-S01-00**; optionally thicken S01 planner stubs only if gaps are IA-shaped
- [ ] If blocker/high: inline RESEARCH fix **or** spawn `P32-S00-02a` / `02b` immediately below this row (full prompts)

## Findings format

List by severity: `blocker` | `high` | `medium` | `low` | `nit`. Every blocker/high needs inline fix or spawn before PASS.

## Exit criteria

- [ ] No open blocker/high without spawn or inline fix
- [ ] Confidence **medium** or **high** with evidence in board Notes (cite RESEARCH path + checklist ticks)
- [ ] Next: **P32-S01-00** (unless spawn remediation rows inserted first)

## Todo updates

Status + notes on **P32-S00-02**; may thicken upcoming **S01** prompts only (forward-only).

## Next

`P32-S01-00` (or spawned `P32-S00-02a` if remediation required)
