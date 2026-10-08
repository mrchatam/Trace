# S04 — Envelope + schema + descriptions — SCOPE-TODOS

Issues: [#157](https://github.com/mrchatam/Trace/issues/157), [#158](https://github.com/mrchatam/Trace/issues/158), [#159](https://github.com/mrchatam/Trace/issues/159), [#161](https://github.com/mrchatam/Trace/issues/161).

| ID | Item | Status |
|----|------|--------|
| S04-0 | P45-S04-00 thicken 01/02 | **done** 2026-10-04 |
| S04-1 | `listEnvelope`; `schemaFor` enums/ranges; generated descriptions; `changeShowRow` | **done** 2026-10-04 (P45-S04-01) |
| S04-2 | Verification record (self-review) | **done** 2026-10-04 **PASS** (P45-S04-02) |

## Acceptance

- [x] One envelope for changes/regressions/agents/tasks lists, with `truncated` that is always a real signal — #157
- [x] `trace_agents recommend` empty result carries the install hint — #157
- [x] Enums + ranges on all 17 tools via explicit `InputSchema`; `trace_explore` validates like `trace_context` — #158
- [x] Descriptions generated from the enum vocabularies; drift guard fails the build — #159 (caught 2 real drifts on day one)
- [x] `trace_changes show` fully snake_case; `RowID` off the wire — #161
- [x] Key-shape sweep over 21 real tool calls finds zero uppercase keys
- [x] Store opened once per changes/regressions call