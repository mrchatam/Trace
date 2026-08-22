# DESIGN-LOCKS — Phase 38

**Investigation-first peer gap analysis.** No implement.

| Lock | Value |
|------|-------|
| Theme | Trace vs peers — **retrieval, context, code graph, MCP, GUI explore, harness orient** |
| Must produce | `GAP-REGISTRY.md`, `SATURATION-NOTES.md`, `REMEDIATION-PLAN.md` |
| Must not | Product Go/TS changes; “fix while investigating”; premature PLAN before saturation gate |
| Investigation tools | See [`PEER-FIXTURES.md`](PEER-FIXTURES.md) |
| Saturation gate | **S05** — reviewer must sign “confident exit” or spawn more S01–S04 rows |
| Planning gate | **S06** only after S05 PASS (saturation APPROVE) |
| Successor | **Remediation implementation phase** (Phase 39+), human-promoted — not auto-built in P38 |
| Laws | G6 no dump; G19 adapters; local-first; DR-NOSSEM noted not silently violated |

## Investigation row rules

- Each investigate row: **read + compare + cite** (file:line or peer path + mechanism)  
- Spawn new board rows when a gap needs **dedicated** peer or Trace slice — do not cram unbounded scope into one row  
- “Potential improvement” ≠ “accept for build” until REMEDIATION-PLAN ranks it  

## Saturation exit criteria (S05)

Investigation may close only if reviewer documents:

- [ ] Every H1–H11 in INTAKE has **verified / rejected / deferred** with evidence pointer  
- [ ] At least one **live** Trace command trace per major gap claim (CLI or MCP)  
- [ ] Each primary peer (CG, UA, GF) has **mechanism cite** (not README-only)  
- [ ] Cross-matrix lists **Trace strengths** peers lack (moat row)  
- [ ] Spawn list empty or explicitly deferred with trigger  
- [ ] Reviewer confidence **high** that new rows would duplicate, not discover  

## REMEDIATION-PLAN shape (S06)

- Ranked **G1–Gn** improvement themes  
- Each: problem, evidence (GAP id), peer pattern, Trace-law fit, **proposed future phase/sketch**, risk, **not in P38**  
- Explicit **reject** list (anti-patterns from P24 §4)
