# S02 — Cognitive artifacts — scope todos

| Order | ID | Role | Status cue |
|------:|----|------|------------|
| 1 | P20-S02-00 | scope planner | done — 016 schema + invalidate/reconsider APIs FINAL |
| 2 | P20-S02-01 | implementer | done — 016 + domain APIs + 14 named tests; ceiling 16 |
| 3 | P20-S02-02 | reviewer | pending — merge discipline + blocking count + Law 11 |

**Depends on:** S01 `PolicyInputs.blocking_uncertainty_count` + `ApplyDeliberationTransition`. **Feeds:** S06 fills PolicyInputs from `CountBlockingUncertainties`; S03 may link Hypothesis; S07 seed residual.

**Locks (do not re-debate in 01):** migration **016**; tables `uncertainties` / `hypotheses` / `decision_reconsiderations`; no Finding fork; Assumption invalidate = STALE/SUPERSEDED; blocking query task-scoped.
