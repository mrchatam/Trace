# Project Documentation Index

## Execution (agents)

- [AGENTS.md](AGENTS.md)
- [docs/TODO.md](docs/TODO.md) — **run-order index** (phase boards in [docs/TODO/](docs/TODO/))
- [docs/rules/agent-loop-protocol.md](docs/rules/agent-loop-protocol.md)
- [docs/rules/project-rules.md](docs/rules/project-rules.md)
- [docs/rules/skills-map.md](docs/rules/skills-map.md)
- [docs/phases/phase-00-foundation/README.md](docs/phases/phase-00-foundation/README.md)

## Core design

- README.md
- docs/ROADMAP.md
- docs/ARCHITECTURE.md
- docs/PROJECT_MODEL.md
- docs/DESIGN_DECISIONS.md

## Algorithms

- docs/PLANNING.md
- docs/RETRIEVAL_AND_CONTEXT.md
- docs/REVIEW_AND_VERIFICATION.md
- docs/DECISION_IMPACT.md
- docs/AGENT_ENVIRONMENT.md

## Engineering

- docs/STORAGE_AND_PERFORMANCE.md
- docs/SECURITY.md
- docs/EVALUATION.md
- CONTRIBUTING.md
- LICENSE

## Initialization registers (design SoT, not run board)

See [docs/init/README.md](docs/init/README.md).

## Research

- [docs/research/SIMILAR-PROJECTS-REVIEW-PROMPT.md](docs/research/SIMILAR-PROJECTS-REVIEW-PROMPT.md) — paste-ready agent prompt to mine techniques from `similar projects/`
- [docs/research/SIMILAR-PROJECTS-REVIEW-OUTPUT-TEMPLATE.md](docs/research/SIMILAR-PROJECTS-REVIEW-OUTPUT-TEMPLATE.md) — findings template for consistent review runs
- [docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md](docs/research/SIMILAR-PROJECTS-REVIEW-2026-08-16.md) — 2026-08-16 peer review; thin **Phase 12** closed (S01+S02); ranks 4–6 thin-cut boarded as **Phase 14**; ranks 7+ / S05 still deferred
- [docs/research/TRACE-GOALS-PROGRESS-2026-08-17.md](docs/research/TRACE-GOALS-PROGRESS-2026-08-17.md) — 2026-08-17 goals vs achieved vs remaining; §4 #1 → **Phase 14** scaffold
- [docs/research/PORTABLE-GRAPH-GIT-2026-08-17.md](docs/research/PORTABLE-GRAPH-GIT-2026-08-17.md) — 2026-08-17 clone-readable semantic graph via seed JSON; **GO** → **Phase 17 complete** (portable-graph-git)
## Dogfood experiments

- [experiments/LADDER.md](experiments/LADDER.md) — capability isolation ladder (open-ended)
- [experiments/RESULTS.md](experiments/RESULTS.md) — scored portfolio
- [experiments/README.md](experiments/README.md) — experiment index
- [experiments/BATCH-D04b-D32.md](experiments/BATCH-D04b-D32.md) — scored (D04b…D32)
- [experiments/G1-NATURAL.md](experiments/G1-NATURAL.md) — natural G1 prompts (no CLI cookbook)
- [experiments/BATCH-D03-D05-D07.md](experiments/BATCH-D03-D05-D07.md) — scored
- [experiments/BATCH-D04-D06-D11.md](experiments/BATCH-D04-D06-D11.md) — scored
- [experiments/DOGFOOD-FINDINGS.md](experiments/DOGFOOD-FINDINGS.md) — DF-* log (canonical DF-17+ after 2026-08-16 collision map)
- [experiments/NATURAL-RERUN.md](experiments/NATURAL-RERUN.md) — natural G1 portfolio re-run
- [experiments/BATCH-D09b-D30b-D34.md](experiments/BATCH-D09b-D30b-D34.md) — D09b/D30b/D34 batch
- [docs/phases/phase-09-dogfood-hardening/README.md](docs/phases/phase-09-dogfood-hardening/README.md) — Phase 09 (closed)
- [docs/phases/phase-10-integrity-surfaces/README.md](docs/phases/phase-10-integrity-surfaces/README.md) — Phase 10 (closed; historical `no successor`)
- [docs/phases/phase-11-residual-surfaces/README.md](docs/phases/phase-11-residual-surfaces/README.md) — Phase 11 (closed; 18 DFs; historical `no successor`)
- [docs/phases/phase-11-residual-hardening/SUPERSEDED.md](docs/phases/phase-11-residual-hardening/SUPERSEDED.md) — superseded early P11 scaffold
- [docs/phases/phase-12-peer-honesty-surfaces/README.md](docs/phases/phase-12-peer-honesty-surfaces/README.md) — Phase 12 thin peer honesty (closed; historical `no successor`)
- [docs/phases/phase-13-import-resolve-honesty/README.md](docs/phases/phase-13-import-resolve-honesty/README.md) — Phase 13 post–P12 DF-60…67 (closed; historical `no successor`)
- [docs/phases/phase-14-peer-impact-install-gates/README.md](docs/phases/phase-14-peer-impact-install-gates/README.md) — Phase 14 thin peer impact + install gates (closed; historical `no successor`)
- [docs/phases/phase-15-p14-residual-plan/README.md](docs/phases/phase-15-p14-residual-plan/README.md) — Phase 15 P14 residual remediation (closed; historical `no successor`; R1 MCP Assert)
- [docs/phases/phase-16-assert-root-and-surfaces/README.md](docs/phases/phase-16-assert-root-and-surfaces/README.md) — Phase 16 post-P15 open DFs (closed; DR-HANDOFF `no successor`; DF-72 thin `trace_impact`)
- [docs/phases/phase-17-portable-graph-git/README.md](docs/phases/phase-17-portable-graph-git/README.md) — Phase 17 portable graph git (closed; DR-HANDOFF `no successor`; DF-80…85 fixed; DF-86 deferred)
- [docs/phases/phase-17-portable-graph-git/DF-84-FORWARD.md](docs/phases/phase-17-portable-graph-git/DF-84-FORWARD.md) — P17 upcoming locks (plan tree, `exported_at_commit`, deferred git-hook)
- [docs/phases/phase-18-fts-clone-honesty/README.md](docs/phases/phase-18-fts-clone-honesty/README.md) — Phase 18 D40 residuals (boarded; DF-87/88/89 landed S01–S03; S04 VERIFY APPROVE; P18-S05-00 FINAL rebuild locks; P18-S05-01 rebuilt; next **P18-S05-02**)
- [docs/phases/phase-18-fts-clone-honesty/DF-88-DECISION.md](docs/phases/phase-18-fts-clone-honesty/DF-88-DECISION.md) — keep P17 seed exclude; clone PENDING expected
- [docs/phases/phase-18-fts-clone-honesty/DR-HANDOFF.md](docs/phases/phase-18-fts-clone-honesty/DR-HANDOFF.md) — default `no successor`; close on **P18-S05-02**
- [docs/phases/phase-16-allowlist-seed-impact-surfaces/SUPERSEDED.md](docs/phases/phase-16-allowlist-seed-impact-surfaces/SUPERSEDED.md) — superseded parallel P16 draft
- [docs/phases/phase-16-allowlist-seed-impact-surfaces/SUPERSEDED.md](docs/phases/phase-16-allowlist-seed-impact-surfaces/SUPERSEDED.md) — superseded parallel P16 scaffold (DF-72-as-fix; do not run)
- [experiments/ab-import-resolve/](experiments/ab-import-resolve/) — optional DF-60 isolation dogfood (prepare + rubric)
- [experiments/POST-P15-DOGFOOD.md](experiments/POST-P15-DOGFOOD.md) — post-P15 regression dogfood (DF-68)
- [experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md](experiments/_bughunt/post-p15/POST-P15-BUGHUNT.md) — post-P15 hunt (DF-75…78)
- [experiments/BATCH-D21-D23.md](experiments/BATCH-D21-D23.md) — D21–D23 combo (DF-70…74)
