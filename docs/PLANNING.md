# Progressive Planning

## 1. Philosophy

The planner must not attempt to predict the complete future of a project in exhaustive detail.

Instead:

> Plan broadly early, deeply near execution, and continuously revise future work as the project reveals new information.

## 2. Planning hierarchy

```text
Goal
 ↓
Project plan
 ↓
Phase
 ↓
Scope
 ↓
Task
 ↓
Implementation detail
```

Planning detail increases as uncertainty decreases.

## 3. Initial project planning

The initial planner should produce:

- goal interpretation;
- constraints;
- quality bar;
- coarse architecture;
- phases;
- known dependencies;
- known decisions;
- assumptions;
- unresolved questions;
- minimal early tasks;
- areas explicitly expected to require later discovery.

It should not fabricate complete downstream todo lists.

## 4. Scope planning

At scope start, the planner performs a deep refresh using:

- current code;
- project graph;
- current plan;
- decisions;
- assumptions;
- discoveries;
- previous review findings;
- future tasks likely to be affected;
- current environment capabilities.

It produces:

- current objective;
- current architecture;
- locked decisions;
- exact task boundaries;
- exit criteria;
- likely downstream effects;
- unresolved questions;
- future work changes.

## 5. Implementation

An implementer receives:

- task packet;
- direct code context;
- relevant decisions;
- relevant assumptions;
- selected environment capabilities;
- verification plan.

The implementer may:

- modify the assigned code;
- report discoveries;
- propose plan changes.

It should not silently rewrite future planning.

## 6. Review

A fresh reviewer reconstructs the required state independently.

The reviewer may:

- reject;
- request fixes;
- make minor in-scope fixes;
- create/modify future tasks;
- report discoveries;
- trigger cause investigation.

## 7. Scope replan

After meaningful review findings or discoveries:

```text
current implementation knowledge
        +
new evidence
        ↓
scope replan
        ↓
future task graph updates
```

Only affected future work should be revised.

## 8. Phase review

At the end of a phase, evaluate:

- goal alignment;
- requirements satisfied;
- architecture direction;
- assumptions still valid;
- discoveries;
- future phase impact;
- technical debt introduced;
- unresolved risks.

The phase reviewer can recommend a phase replan.

## 9. Unknowns

The planner explicitly tracks:

```text
KNOWN
ASSUMED
UNRESOLVED
EXPECTED_DISCOVERY
UNKNOWN
```

Unresolved or uncertain information should not be silently upgraded into established facts.

## 10. Minimal scope

Tasks should be kept small enough that:

- intent is unambiguous;
- verification is tractable;
- review is cheap;
- failures remain localized.

The planner may split tasks when review or evidence becomes too broad.

## 11. Forward movement

Normal work should advance the current project state.

When a solution must be abandoned:

- record why;
- mark affected work;
- create a replacement route;
- optionally create an explicit reversal state.

Never silently rewrite history.

## 12. Replanning triggers

A replan should be considered after:

- architecture-changing discovery;
- repeated failure;
- invalidated assumption;
- major user decision;
- unexpected dependency;
- review finding affecting future work;
- capability/environment change;
- scope boundary violation;
- significant external constraint change.
