# Review and Verification

## 1. Core principle

An agent's claim is not evidence.

```text
agent claim
   ↓
verification
   ↓
evidence
   ↓
verified state
```

The project should make incorrect claims hard to promote into authoritative state.

## 2. Evidence hierarchy

Suggested strength levels:

1. agent assertion;
2. claimed artifact/change;
3. independent observation;
4. deterministic evidence;
5. independently reproduced behavior;
6. human verification.

The exact acceptance policy depends on task risk.

## 3. Verification plan

Every meaningful task should have explicit verification items.

Example:

```text
V1 expected file exists
V2 API compiles
V3 success case works
V4 failure case works
V5 regression test passes
V6 no unrelated public behavior changed
```

Each item gets:

- status;
- evidence;
- verifier;
- confidence;
- residuals.

## 4. Todo review

Question:

> Does the actual repository state satisfy this task?

Process:

1. reconstruct requirements independently;
2. inspect actual repository;
3. inspect diff/history;
4. verify every exit criterion;
5. reproduce critical claims;
6. attack likely failure modes;
7. record evidence.

Review prompts should use:

- falsification;
- counterexample search;
- boundary analysis;
- assumption challenge;
- regression search.

## 5. Scope review

Question:

> Do the completed tasks form a coherent implementation?

Review:

- cross-task contracts;
- architectural boundaries;
- shared interfaces;
- dependencies;
- future task assumptions;
- discoveries;
- scope creep.

## 6. Phase review

Question:

> Did the phase move the project toward the intended goal?

Review:

- goal alignment;
- requirements;
- architectural trajectory;
- assumptions;
- technical debt;
- discovered constraints;
- future-phase validity;
- user decisions.

## 7. Independent reviewers

For meaningful review:

- fresh agent context;
- different agent identity;
- no dependence on implementer narrative;
- optionally different model for high-risk work.

Do not treat multiple reviewers as a voting system. Disagreement should trigger deeper investigation.

## 8. Adversarial review

High-risk tasks may include an explicit adversarial reviewer:

> Try to prove the implementation should not be accepted.

It should search for:

- counterexamples;
- missed edge cases;
- incorrect assumptions;
- regressions;
- incomplete integration;
- false verification claims.

## 9. Fix loops

Recommended flow:

```text
implement
 ↓
review
 ↓
fix
 ↓
review
```

After repeated failures:

```text
failure
 ↓
cause investigation
 ↓
hypothesis matrix
 ↓
reproduction
 ↓
new fix plan
```

Do not permit endless patch loops.

## 10. Verification policy

Task classes may require different evidence.

Examples:

### Refactor

- build/typecheck;
- relevant tests;
- diff review.

### API change

- unit tests;
- integration tests;
- schema checks;
- compatibility checks.

### UI

- automated tests;
- browser verification;
- visual check;
- human verification when appropriate.

### Security-sensitive

- targeted automated checks;
- independent review;
- human gate where necessary.

## 11. Human verification

Machine verification must not pretend to prove subjective or physically interactive outcomes.

Human evidence includes:

- verifier;
- date;
- environment;
- method;
- result;
- artifacts;
- notes.

## 12. Truth gradient

A fact may move through:

```text
UNVERIFIED
   ↓
CLAIMED
   ↓
OBSERVED
   ↓
REPRODUCED
   ↓
VERIFIED
   ↓
HUMAN_VERIFIED
```

Confidence refers to evidence quality, not model self-confidence.
