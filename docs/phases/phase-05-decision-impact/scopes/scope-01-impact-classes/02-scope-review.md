# P05 / S01 / 02 — Scope review (impact classes)

## Metadata
- id: P05-S01-02
- todo_ids: [P05-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of P05-S01-01 against S01-00 / `01-impact-classes.md` locks. APPROVE with evidence or spawn forward remediations. Write [REVIEW-NOTES.md](REVIEW-NOTES.md).

## Session start
Agent → clarify if needed → Plan → execute (review).

## Review focus (required)

### Claims vs locks
- Mig **`009_decision_impact.sql`** only (additive; tables `decision_impact_findings` + `decision_alternatives`; no rewrite of 001–008; no ALTER on `decisions`)
- Package is **`internal/domain`** + store helpers — **no** `internal/impact`, **no** planner fork
- Enums exact: impact `SAFE|CAUTION|HIGH|DESTRUCTIVE|REVERSAL`; uncertainty `KNOWN|LIKELY|POSSIBLE|UNKNOWN` (empty→UNKNOWN); kinds `AFFECTED_WORK|INVALIDATED_ASSUMPTION|WORK_AT_RISK|NEW_WORK|UNRESOLVED`
- APIs: Add/List findings; Add/List/SetRecommended alternatives (single recommended); `ImpactReport` with fail-closed `HasUnknown` / `Incomplete` / OverallClass rollup order
- **No new** `entity_links` rels; `CreateDecision` / `GetDecision` / `LinkDecisionTask` (`decision_affects_task`) behavior unchanged
- Thin CLI `trace impact` (finding/alternative/report) G19 — no business logic in `cmd/trace`; MCP impact tools absent
- DR-NOIMP: no commercial auto-engine / embeddings / silent plan mutation; `plan simulate` absent

### Gate F readiness (S02)
- Confirm S02 stubs / Depends list usable plant hooks: `AddImpactFinding`, `LinkDecisionTask`, `ImpactReport` fields (`HasUnknown`, `Incomplete`, `OverallClass`, findings slice includes UNKNOWN)
- Light-thicken upcoming S02 prompts if hooks differ from what S01 shipped (forward-only)

### Carry-forward bars
- Honesty Paths A/B/C + Gate G `TestHonestyEscapeRateGateGPrelim`
- Gate E `TestPlantedDiscoveryReplan`
- p0x 7/7; x0; `CGO_ENABLED=0` domain+store; `CGO_ENABLED=1` `./...`
- Gate C `docs/verification/gate-c-x0/` untouched (`dry_run:false` intact)

### Laws / hygiene
- G19 reverse-import: domain must not import `cmd/trace`
- No daemon/HTTP/embeddings primary
- VerifiedFact still out

## Suggested evidence commands

```text
CGO_ENABLED=0 go test ./internal/domain/... ./internal/store/... ./evals/honesty/... ./evals/replan/...
CGO_ENABLED=1 go test ./evals/p0x/... ./evals/x0/... ./...
```

## Exit criteria
- [ ] Verdict + confidence + REVIEW-NOTES.md
- [ ] No open blocker/high without spawn (`P05-S01-02a`/`02b` forward)
- [ ] Board status + Notes; S02 Depends note accurate if hooks changed
