# Decision Impact and Hypothetical Planning

## 1. Purpose

Users must remain in control of the project while receiving early warning about destructive or high-impact choices.

The system should advise, not silently prevent.

## 2. Decisions are first-class entities

A decision records:

- question;
- chosen option;
- alternatives;
- rationale;
- evidence;
- author;
- approver;
- assumptions;
- affected requirements;
- affected work;
- consequences.

## 3. Impact classes

```text
SAFE
CAUTION
HIGH_IMPACT
DESTRUCTIVE
REVERSAL
```

## 4. Knowledge confidence

Each impact finding is classified:

```text
KNOWN
LIKELY
POSSIBLE
UNKNOWN
```

Unknowns must be visible.

## 5. Impact analysis

For a proposed decision:

```text
proposed decision
   ↓
affected decisions
   ↓
assumptions
   ↓
architecture
   ↓
files/symbols
   ↓
tasks/scopes/phases
   ↓
completed work at risk
   ↓
new work
```

## 6. User interaction

A high-impact result should communicate:

- what conflicts;
- why it conflicts;
- what work is affected;
- what becomes invalid;
- what remains safe;
- what is uncertain;
- recommended route;
- alternative routes.

The user may proceed anyway.

## 7. Alternative routes

For major changes, the planner may propose:

- direct replacement;
- compatibility layer;
- staged migration;
- architecture-preserving alternative.

The system should explain tradeoffs rather than pretending one answer is always correct.

## 8. Hypothetical plans

Users should eventually be able to simulate:

```text
current plan
   ├── option A
   ├── option B
   └── option C
```

A simulation creates a temporary graph state.

Simulation must not mutate the real project until explicitly adopted.

## 9. Reversal

A reversal is a new state transition that points back toward an earlier direction.

It must record:

- what is being reversed;
- why;
- what work becomes obsolete;
- what work can be preserved;
- user approval where required.

## 10. Destructive-change policy

The system should not silently stop a user.

Instead:

```text
low impact      → continue
moderate impact → inform
high impact     → warn + alternatives
destructive     → explicit confirmation
reversal        → explicit reversal plan
```
