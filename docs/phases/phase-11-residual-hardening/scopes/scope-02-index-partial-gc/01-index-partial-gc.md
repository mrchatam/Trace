# P11-S02-01 — Index partial-path GC (STUB)

## Metadata
- id: P11-S02-01
- todo_ids: [P11-S02-01]
- role: implementer
- verification: automated

## Objective
Implement DF-40 per **P11-S02-00** FINAL locks.

## Depends-on
- P11-S02-00 done; Phase 10 DF-20 tests stay green

## Exit criteria (outline)
- [ ] Named regression: rename → `index <new-path>` → old path/symbols gone (or locked equivalent)
- [ ] `TestIndexGCAfterPathRename` / argv isolation / incremental isolation still PASS
- [ ] No full-rebuild architecture; no MCP index
- [ ] CGO1 cmd/trace + p0x/x0 + product `./...` PASS
