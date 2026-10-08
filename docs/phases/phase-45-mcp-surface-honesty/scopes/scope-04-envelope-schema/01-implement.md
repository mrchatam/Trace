# Phase 45 / S04 / Implement — Envelope + schema + descriptions

## Row `P45-S04-01` — **done** 2026-10-04 (audit-hardened). Issues [#157](https://github.com/mrchatam/Trace/issues/157), [#158](https://github.com/mrchatam/Trace/issues/158), [#159](https://github.com/mrchatam/Trace/issues/159), [#161](https://github.com/mrchatam/Trace/issues/161).

## 1. Uniform list envelope — `internal/mcp/envelope.go` (new)

```go
type listEnvelope struct {
    OK         bool   `json:"ok"`
    Items      any    `json:"items"`          // always a slice, never null
    Count      int    `json:"count"`
    Truncated  bool   `json:"truncated"`
    Limit      int    `json:"limit,omitempty"`
    NextCursor string `json:"next_cursor,omitempty"`
    Hint       string `json:"hint,omitempty"`
}
func newListEnvelope(items any, count int, truncated bool) listEnvelope
func (e listEnvelope) withLimit(int) / withCursor(string) / withHint(string) listEnvelope
```

Applied: `trace_changes list` (gains `limit` via `changesListLimit`), `trace_regressions list` (`regressions` → `items`), `trace_agents list` + `recommend` (recommend gains `EmptyCatalogHint`), `trace_tasks` (gains the missing `ok`). After the audit round: `truncated` is the **page-fill signal** (`len(items) == effective limit`) in changes/regressions, with the effective `limit` echoed — never a hardcoded `false`.

Deliberately **not** converted: `trace_capability list` / `missing` / `trace_impact finding list` — their domain-named keys (`capabilities`, `missing`, `findings`) mirror the CLI output documents byte-for-byte. See [`../../AUDIT.md`](../../AUDIT.md) residuals.

## 2. Schema enums/ranges — `internal/mcp/schema.go` (new)

`schemaFor[In](constraints)` infers via `jsonschema.For[In]`, applies `propConstraint{enum, minimum, maximum}` per property, panics on a missing property or an unresolvable schema. Shared vocabularies (single source for enums **and** descriptions): `capabilityActions/Kinds/Status`, `reviewActions/Results`, `changeActions`, `regressionAction`, `agentActions`, `planActions`, `impactActions/Ops`, `loopActions` (incl. `gate`) + `loopGateFors`, `linkRels`, `contextFormats`, `matchers`, `workStates` (8), `addKinds` (8), `provenanceStatus`.

Wired in `server.go`: all 17 tools use `Description: desc*` + explicit `InputSchema: schemaFor[...](...)`. Ranges only where the code enforces them: `trace_context depth 1..2 / max_layer 1..3`, `trace_impact depth 1..8`, `trace_search/changes/regressions/explore limit 1..64`, `trace_explore max_nodes 1..500 / depth 1..8`, `trace_tasks/capability limit 1..500`.

#158's second half: `trace_explore` now performs the same validation family as `trace_context` (`resolveExploreBounds`) instead of deferring to compiler defaults silently.

## 3. Generated descriptions — `internal/mcp/descriptions.go` (new)

All 17 descriptions generated from the shared vocabularies. The drift guard `TestToolDescriptionsCoverSchemaEnums` fails when a schema enum value is absent from its description — it **immediately caught two real drifts** (`trace_add.status` enum and `trace_transition.to_state` missing from their descriptions; both fixed). Audit round: `descVersion` now names `project_root`; `descAgents` uses `join(agentActions)`.

## 4. snake_case show DTO (#161) — `internal/mcp/tools_changes.go`

`changeShowRow` + `changeShowRowFrom(store.Change)`: every field tagged snake_case, `ParentChangeID`/`Actor`/`Reason`/`SourceType`/`LastVerifiedAt` omitempty, **`RowID` dropped** (internal SQLite rowid, never on the wire). `trace_changes show` returns `{ok, change, paths}` with `paths` rows likewise tagged.

## 5. Store-handle hygiene (audit round)

`toolChanges`/`toolRegressions` previously opened the store, then opened it **again** inside each action handler. Both now open once and pass `*store.Store` to `changesList/Show/Compare` and `regressionsList`.

## 6. Acceptance — all PASS

`internal/mcp/schema_test.go`: `TestToolSchemasCarryEnumsAndRanges`, `TestToolDescriptionsCoverSchemaEnums`, `TestServerExposesEveryToolDescription`, `TestChangeShowRowIsSnakeCase`, `TestListEnvelopeIsUniform` (5 tools × {ok, items, count, truncated} + items-is-array + count==len(items)), `TestAgentsRecommendEmptyCarriesHint`, helpers `toolListSchemas`/`propSchema`/`seedEnvelopeFixture`.

Updated in `mcp_test.go` (old bare-array shapes): `TestMCPChangesList`, `TestMCPRegressionsRegistered`, `TestMCPAgentsRecommend`.

## Next

`P45-S04-02`