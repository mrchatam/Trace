# E — Assumption Register

| ID | Assumption | Why we assume it | Confidence | Validation | Status |
|----|------------|------------------|------------|------------|--------|
| A1 | A persistent causal/work graph can beat repo+Git tools on understanding tasks when fairly seeded | Core thesis H1 | Medium | Experiment X0 / Gate C | EXPERIMENT_REQUIRED |
| A2 | Progressive planning beats exhaustive upfront planning on multi-hour agent work | Thesis H2 | Medium-low | Gate D after planner exists | EXPERIMENT_REQUIRED |
| A3 | Independent review + deterministic evidence reduces false completions enough to justify cost | Thesis H5 | Medium | Honesty suite | EXPERIMENT_REQUIRED |
| A4 | Users/agents will invest in recording goals/decisions/discoveries | Without this, graph stays empty | Medium-low | Dogfood + agent seeding measurement (after human GT) | EXPERIMENT_REQUIRED |
| A5 | SQLite + indexes suffice until ≥200k LOC structural graphs | Avoid premature infra | Medium | Gate H ladders | ACCEPTED_RISK until measured |
| A6 | Exact + lexical + graph retrieval suffice for P0-X/X0 without embeddings | Reduce moving parts | Medium | P0-X (P00-S09-01) then X0 | PROVISIONAL (P0-X held; X0 still open) |
| A7 | Thin event log is enough for provenance/rebuild | Simplicity | Medium | Recovery drill later | PROVISIONAL |
| A8 | TS/JS+Python analyzer quality is “good enough” for fixture repos | Ship slice | Medium | Structural query tests | PROVISIONAL |
| A9 | Library + CLI (no daemon) is sufficient to validate the core model | User Round 1 | High | P0-X | ANSWERED / accepted |
| A10 | Single-user local-first is the only deployment that matters for proof | Focus | High | Product non-goals | ACCEPTED_RISK |
| A11 | Planted ground-truth in synthetic fixtures can stand in for “real” causality for foundation tests | Needed for eval | Medium | External validity risk; OSS later | ACCEPTED_RISK for P0-X/X0 |
| A12 | MCP should follow a validated CLI context API (adapter, not architecture driver) | User Round 1 | High | Post–P0-X MCP task | ANSWERED / accepted |
| A13 | Plan churn controls will prevent bureaucracy failure mode | Grill G10 | Medium | Plan churn metric | EXPERIMENT_REQUIRED |
| A14 | Prompt-injection via retrieved code is mitigated by labeling + non-elevation | Security model | Medium-low | Red-team later | ACCEPTED_RISK early |
| A15 | The settled 7-point P0-X bar (incl. structural + incremental) is enough to close P0 before agent Gate C | User Round 2 | High | P00-S09-01 | VALIDATED |
| A16 | Go is a good long-term fit for CLI + eventual daemon concurrency | User Q1=c | Medium-high | Contributor experience; Gate H | ACCEPTED_RISK |
| A17 | tree-sitter “minimal symbols/imports” is sufficient structural foundation for P0-X | User Q4=a | Medium | P00-S09-01 (criterion-2 + #7 sibling isolation) | VALIDATED |
| A18 | Abstracting Git behind an interface while using CLI first avoids later rewrite cost | User Q5=a | High | T003 design review | ACCEPTED_RISK |

Confidence is not evidence. Status moves only when validation runs or a Decision supersedes the assumption.
