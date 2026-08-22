# H — Verification Strategy

## Principles

- Implementers do not self-certify success into verified state.
- Deterministic checks first; LLM review second; human when required.
- Different review layers have different questions.

## Layers

### 1. Todo (task) review

**Question:** Does the repository state satisfy this task’s exit criteria?

**Required:**

- Reconstruct exit criteria from the task packet (not implementer chat).
- Inspect diff / file state via Git.
- Run every automatable verification item; store Evidence.
- Attempt falsification (missing edge cases, scope violations).
- Result: PASS | FAIL | UNCERTAIN with residuals.

**Forbidden:** Promoting DONE based only on implementer narrative.

### 2. Scope review

**Question:** Do completed tasks form a coherent S0 (or later scope)?

**Required:**

- Cross-task contracts (CLI ↔ core ↔ schema).
- Laws compliance spot-check.
- Downstream enablement (can X0 run?).
- Discoveries filed for gaps.
- May open PlanChange for future work; must not silently rewrite history.

### 3. Phase review

**Question:** Did the phase move toward the goal? Are assumptions still valid?

Used from Phase 2 onward (Gate C report is the first phase-style review).

### 4. Project-level review

**Question:** Do Gates C–H and product thesis still justify continued investment?

Triggered after major experiments or yearly/major milestones.

## Fix loops

```text
implement → todo review → fix → todo review
```

After **3** consecutive FAIL cycles on the same task: stop patching; run cause investigation (hypotheses, reproduction, possibly split task or replan). Record as Discovery if systemic.

## Authority matrix (v0)

| Action | Implementer | Reviewer | Human |
|--------|-------------|----------|-------|
| Edit code in task scope | yes | minor fixes OK | yes |
| Add Claim | yes | yes | yes |
| Add Evidence | yes (capture) | yes | yes |
| PASS/FAIL Review | no | yes | yes |
| Promote VerifiedFact | no | yes if policy met | yes |
| Mark Task DONE | no | yes if PASS | yes |
| Adopt PlanChange affecting future tasks | propose only | propose | ack required |
| Override impact/warning | no | no | yes |

## Evidence strength (reminder)

Assertion < artifact claim < independent observation < deterministic evidence < reproduced behavior < human verification.

## S0-specific verification

Every T00x lists exit criteria and verification. Scope S0 closes only via T014 checklist in `C_FIRST_SCOPE.md`.

## Anti-patterns

- Endless “one more fix” without cause investigation.
- Reviewer copying implementer checklist blindly.
- Treating token confidence as proof.
- Skipping honesty scenarios because “tests pass.”
