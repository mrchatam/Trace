# P16-S04-00 — Seed rels + impact import (stub — thicken vs live)

## Metadata
- id: P16-S04-00
- todo_ids: [P16-S04-00]
- role: planner
- skills: [planning-and-task-breakdown]
- mcps: [Shell, Read, Grep, Glob, Write]
- verification: automated

## Objective
Lock **FINAL** for **DF-70** (seed `discovery_mentions_task`) and **DF-73** (seed impact findings/alternatives). **No product Go.**

## Live gap (P16-00)
`cmd/trace/seed.go` switch omits `LinkDiscoveryMentionsTask`; seedDocument rejects unknown top-level keys (no findings/alternatives). CLI/MCP `discovery-mentions-task` already works (DF-42).

## Inherited locks
- Accept `discovery_mentions_task` and hyphen alias `discovery-mentions-task` (same as CLI)
- Seed JSON may include `impact_findings` / `decision_alternatives` (00 locks exact keys + required fields)
- Unknown **other** top-level keys still rejected
- Use existing domain AddImpactFinding / AddAlternative — no new entity_links rels
- Named: `TestSeedImportDiscoveryMentionsTask`; `TestSeedImportImpactFindings`

## Planner work
1. [ ] Confirm seed switch + JSON shape vs domain APIs
2. [ ] Thicken 01/02; **FINAL**; next **P16-S04-01**
