# Project Model

## 1. Entity classes

### Intent entities

- Goal
- Requirement
- Constraint
- Decision
- Assumption

### Work entities

- Project
- Phase
- Scope
- Task
- Review
- Evidence
- Discovery
- PlanVersion
- PlanChange

### Code entities

- Repository
- Directory
- File
- Symbol
- Module

### History entities

- Commit
- Change
- ProjectState

### Environment entities

- Skill
- Rule
- Tool
- MCP
- Hook
- Agent
- Model
- Capability

## 2. Relationship examples

```text
Goal -> Requirement
Requirement -> Task

Decision -> Task
Decision -> Assumption
Decision -> File

Phase -> Scope
Scope -> Task
Task -> Task
Task -> File
Task -> Commit
Task -> Discovery
Task -> Decision

File -> Symbol
Symbol -> Symbol
File -> File

Discovery -> Task
Discovery -> Assumption
Discovery -> PlanChange

Task -> Evidence
Review -> Evidence
```

## 3. Provenance

Every important semantic fact must carry provenance.

Required provenance fields:

- source type;
- source entity;
- evidence references;
- created_at;
- last_verified_at;
- confidence;
- status.

Suggested source types:

```text
DETERMINISTIC
USER_ASSERTED
AGENT_INFERRED
AGENT_PROPOSED
HUMAN_VERIFIED
```

Facts can be marked `STALE`.

## 4. Evidence

A completion claim is not a fact.

Model:

```text
Claim
  ↓
Evidence
  ↓
Verification
  ↓
Verified Fact
```

Evidence can include:

- command + exit code;
- test report;
- diff;
- static analysis result;
- browser observation;
- screenshot/video;
- another agent's independent verification;
- human observation.

## 5. Task state

Suggested task states:

```text
PENDING
IN_PROGRESS
AWAITING_REVIEW
BLOCKED
FAILED
DONE
STALE
SKIPPED
```

A state transition should record:

- actor;
- previous state;
- new state;
- reason;
- evidence references.

## 6. Plan state

A plan is not immutable.

A plan has:

- version;
- scope;
- assumptions;
- confidence;
- known unknowns;
- decisions;
- tasks;
- dependencies.

Plan changes are first-class objects.

## 7. Decision model

A Decision contains:

- question;
- chosen option;
- alternatives;
- rationale;
- evidence;
- affected requirements;
- assumptions;
- consequences;
- author;
- approver;
- status.

Decision statuses:

```text
PROPOSED
ACCEPTED
REJECTED
SUPERSEDED
```

## 8. Discovery model

A Discovery captures unexpected information learned during work.

Fields:

- observation;
- evidence;
- discovered_by;
- discovery_type;
- affected entities;
- invalidated assumptions;
- recommended action;
- resolution;
- confidence.

## 9. Forward-state model

Never erase meaningful project history.

A reversal is a new state transition with explicit rationale.

Useful categories:

```text
PROGRESS
CORRECTION
REPLAN
REVERSAL
```

## 10. Task capability requirements

Tasks may declare:

- required skills;
- required tools;
- required MCPs;
- required rules;
- required verification capabilities;
- permission requirements.

The planner resolves these against the project's current environment.
