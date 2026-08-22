# P06 / S01 / 02 — Scope review (capability surface)

## Metadata
- id: P06-S01-02
- todo_ids: [P06-S01-02]
- role: reviewer
- skills: [code-review-and-quality]
- mcps: [Shell, Read, Grep, Glob]
- verification: automated

## Objective
Independent review of P06-S01-01 against S01-00 / `01-capability-surface.md` locks. APPROVE with evidence or spawn forward remediations. Write [REVIEW-NOTES.md](REVIEW-NOTES.md).

## Session start
Agent → clarify if needed → Plan → execute (review).

## Review focus
- Claims vs S01-00 locks:
  - mig **`010_capability_surface.sql`** (`capabilities` + `task_capability_requirements`)
  - package **`internal/domain`** + store + **`internal/compiler`** (no `internal/capability` / megastore / planner fork)
  - kinds `SKILL`\|`RULE`\|`MCP`\|`TOOL`\|`HOOK`; status `AVAILABLE`\|`UNAVAILABLE`\|`UNKNOWN`
  - APIs Upsert/Get/List + Require/Unrequire/ListRequired + `MissingCapabilities`
  - `BuiltinMCPCapabilitySpecs` six tools; **no** auto-seed on Open/init
  - Packet `required_capabilities` + `missing_capabilities` only (no catalog dump); SchemaVersion `"0.1"`
  - CLI `trace capability` declare/list/require/unrequire/missing; **no** new MCP tools
  - **No** new entity_links rels
- S02 selection hooks usable (catalog + require + MissingCapabilities + packet fields for `evals/capability`)
- Carry-forward bars (honesty A/B/C / Gate G/E/F / p0x / x0 / Gate C `dry_run:false`)
- No ontology bloat; no daemon/HTTP/embeddings primary; G19 (no domain fork in MCP)
- GC-03/04 remain deferred unless Notes promote

## Exit criteria
- [ ] APPROVE high/medium with evidence **or** spawns inserted
- [ ] REVIEW-NOTES.md written
- [ ] Board Notes updated; light S02 Depends note if hooks differ from stubs
