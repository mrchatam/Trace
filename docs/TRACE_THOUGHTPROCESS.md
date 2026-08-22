# Trace Cognitive Deliberation, Verification, Evaluation, and Progressive Learning — Implementation Planning Prompt

You are an expert software architect and agent-systems engineer working on **Trace**, an AI-agent orchestration and project-intelligence system.

Your task is to produce a **comprehensive implementation plan** for a new core subsystem in Trace that makes connected AI agents progressively reason about projects, improve plans, validate implementations, detect regressions, and accumulate structured knowledge about changes and their effects over time.

Do **not** implement anything yet.

Your job is to deeply analyze the concept, inspect the existing Trace repository and architecture, identify how the new subsystem should fit into the current implementation, and produce a concrete implementation roadmap detailed enough that another coding agent could execute it with minimal architectural ambiguity.

---

# 1. Core Vision

Trace should not be just a task manager, planner, or agent loop.

Trace should act as an **externalized project cognition and engineering feedback system** for AI agents.

The agent should continuously improve its understanding of a project through a loop like:

```text
Understand
    ↓
Identify unknowns and gaps
    ↓
Investigate
    ↓
Explore possible approaches
    ↓
Make/revise decisions
    ↓
Construct/improve plan
    ↓
Critique plan
    ↓
Implement
    ↓
Test
    ↓
Verify
    ↓
Measure
    ↓
Score / evaluate
    ↓
Detect regressions
    ↓
Reflect
    ↓
Update project knowledge
    ↓
Replan when necessary
    ↓
Continue
```

This should NOT be implemented as a simplistic fixed loop.

Instead, Trace should provide a **state-driven deliberation controller** that can determine what kind of thinking or validation is needed next based on the current state of the project.

---

# 2. Important Principle

Trace should NOT attempt to store the agent's private raw chain-of-thought.

Instead, Trace should store **structured, durable outputs of reasoning** that are useful for future decisions and agents.

Examples:

* goals
* requirements
* constraints
* success criteria
* questions
* unknowns
* assumptions
* findings
* evidence
* hypotheses
* candidate approaches
* decisions
* rejected alternatives
* risks
* plans
* tasks
* dependencies
* changes
* expected effects
* observed effects
* tests
* verification results
* scores
* regressions
* historical relationships
* confidence
* lessons learned

The goal is to preserve the useful externalized state of reasoning, not raw hidden reasoning traces.

---

# 3. Two Graph Concepts

Trace should conceptually contain two interconnected graph dimensions.

## A. Project / Engineering Graph

This represents the actual project.

Examples:

```text
Goal
  ↓
Requirement
  ↓
Feature
  ↓
Task
  ↓
Files
  ↓
Symbols
  ↓
Tests
  ↓
Evidence
```

Relationships may include:

* depends_on
* implements
* affects
* blocks
* tested_by
* derived_from
* related_to
* contradicts
* replaces

Example:

```text
Goal:
Support multi-provider routing

    ↓

Requirement:
Provider selection occurs per request

    ↓

Task:
Implement provider selection

    ├── affects → router.ts
    ├── affects → provider.ts
    └── tested_by → routing.test.ts
```

## B. Deliberation / Cognitive Graph

This represents what the agent needs to think about next.

Examples:

```text
Orient
  ↓
Gap Detection
  ↓
Investigate
  ↓
Explore
  ↓
Evaluate
  ↓
Plan
  ↓
Critique
  ↓
Execute
  ↓
Verify
  ↓
Reflect
```

But this must NOT be a rigid linear workflow.

The system must be able to transition dynamically.

For example:

```text
Plan
  ↓
Critique
  ↓
New architectural gap discovered
  ↓
Investigate
  ↓
New dependency discovered
  ↓
Explore
  ↓
Decision
  ↓
Replan
```

The implementation should therefore support **state-driven transitions** rather than only hardcoded sequential phases.

---

# 4. Cognitive Operations / Deliberation Phases

Design a phase model for Trace.

Initially consider:

```text
ORIENT
INVESTIGATE
EXPLORE
PLAN
CRITIQUE
EXECUTE
TEST
VERIFY
EVALUATE
REFLECT
REPLAN
```

Each phase should have:

* purpose
* inputs
* outputs
* allowed operations
* expected structured artifacts
* entry conditions
* exit conditions
* possible next phases
* failure conditions
* resource/token considerations

Examples:

## ORIENT

Goal:

Determine what is actually being requested and establish the current state.

Should identify:

* objective
* scope
* constraints
* success criteria
* relevant existing Trace context
* current task/context

## INVESTIGATE

Goal:

Resolve unknowns that prevent confident planning.

The agent should create explicit unresolved questions such as:

```text
Does the project already have an authentication abstraction?
Where are sessions stored?
Which components depend on this behavior?
```

Questions should become persistent Trace objects.

The phase should terminate when blocking unknowns are sufficiently resolved, or when the system determines that further investigation has diminishing returns.

## EXPLORE

Allow multiple candidate approaches.

For example:

```text
Option A
Option B
Option C
```

Each should contain:

* expected benefits
* disadvantages
* risks
* affected components
* dependencies
* complexity
* confidence
* evidence

## PLAN

Convert understanding and decisions into an implementation plan and task/dependency graph.

The plan must connect back to:

* requirements
* goals
* decisions
* files/components
* verification strategy

## CRITIQUE

Attempt to break the plan.

Ask:

* What is missing?
* Which assumptions are unsupported?
* Are there hidden dependencies?
* Could tasks conflict?
* Does every requirement map to implementation work?
* Does every task have a verification strategy?
* Are there architectural risks?
* Is the plan actually executable?

Critique should produce structured findings rather than silently rewriting the plan.

## EXECUTE

Agents perform actual implementation work.

Trace observes and records meaningful changes.

## TEST

Run relevant automated/manual/instrumented tests according to task/change type and historical risk.

## VERIFY

Determine whether implementation actually satisfies requirements and intended behavior.

Passing tests alone should NOT automatically imply verification.

## EVALUATE

Measure quality and compare against a baseline.

Potential dimensions:

* correctness
* requirement coverage
* test confidence
* performance
* memory
* complexity
* maintainability
* security
* reliability
* regression risk

## REFLECT

Determine what was learned from implementation and evaluation.

Examples:

* new architectural facts
* invalidated assumptions
* unexpected dependencies
* useful historical relationships
* new risks
* failed hypotheses
* successful strategies

## REPLAN

Reopen planning when new evidence invalidates part of the current plan or project understanding.

---

# 5. State-Driven Phase Selection

Do not simply implement:

```text
A → B → C → D → E
```

Instead design a controller that looks at project state and determines what cognitive operation is most appropriate.

Example state:

```text
blocking_questions = 3
critical_assumptions = 1
plan_exists = true
plan_confidence = 0.58
requirement_coverage = 0.74
verification_status = incomplete
```

Potential actions:

```text
INVESTIGATE = 0.94
CRITIQUE = 0.71
PLAN = 0.44
EXECUTE = 0.08
```

The controller should select the most appropriate next operation based on policies/conditions.

Do not over-engineer machine learning for the first implementation.

Start with deterministic policies that are explicit, testable, inspectable, and extensible.

Consider a future abstraction where Trace can learn which deliberation sequences work best for different models/projects.

---

# 6. Entry and Exit Conditions

Every cognitive phase should have explicit entry and exit conditions.

Examples:

### INVESTIGATE

Entry:

```text
blocking unresolved question exists
OR
high-impact uncertainty exists
```

Exit:

```text
no blocking unknowns remain
AND
required evidence exists
```

### PLAN

Entry:

```text
requirements sufficiently understood
AND
critical unknowns resolved
```

Exit:

```text
requirements mapped to tasks
AND
dependencies mapped
AND
verification strategy exists
```

### CRITIQUE

Entry:

```text
candidate plan exists
```

Exit:

```text
no critical flaws
AND
coverage acceptable
AND
risk acceptable
```

### EXECUTE

Entry:

```text
plan confidence meets threshold
AND
no blocking unknowns
```

Exit:

```text
planned implementation completed
```

### VERIFY

Exit:

```text
required evidence obtained
AND
requirements verified
```

or:

```text
new gap/regression found
→ route back into deliberation
```

Design these rules carefully.

---

# 7. Uncertainty as a First-Class Concept

Introduce an explicit representation of uncertainty.

Example:

```text
Uncertainty
 ├── question
 ├── severity
 ├── affected_nodes
 ├── confidence
 ├── evidence
 ├── resolution
 └── status
```

Example:

```text
Question:
Will modifying the auth middleware break websocket authentication?

Severity:
high

Affected:
auth middleware
websocket subsystem
session manager

Confidence:
0.43

Status:
unresolved
```

High-impact unresolved uncertainty should influence phase selection and may block implementation.

---

# 8. Assumptions

Agents constantly make assumptions.

Trace should capture them.

Example:

```text
Assumption:
All API requests pass through middleware.ts

Confidence:
0.65

Evidence:
Observed on 8 existing routes

Impact:
high
```

If later evidence contradicts the assumption:

```text
Assumption invalidated
    ↓
Reflect
    ↓
Replan
```

Design how assumptions are created, updated, invalidated, and linked to affected decisions/tasks.

---

# 9. Decisions and Alternatives

Create a persistent decision representation.

Example:

```text
Decision D21

Chosen:
Use existing authentication provider architecture.

Alternatives:
Custom implementation
External auth service

Why:
- existing infrastructure
- lower complexity
- no new dependency
- compatible with deployment

Rejected because:
Custom implementation introduces unnecessary complexity/security surface.
```

Decisions should reference:

* reasoning artifacts
* evidence
* alternatives
* affected components
* related tasks

Consider supporting **reconsideration triggers**.

Example:

```text
Reconsider when:
- Redis becomes unavailable
- consistency requirement changes
- persistence requirement appears
```

When new project evidence satisfies a reconsideration trigger, Trace should be able to invalidate/reopen the decision.

---

# 10. Changes as First-Class Objects

Every meaningful implementation change should be modeled.

Potential structure:

```text
Change
 ├── id
 ├── parent_change
 ├── actor
 ├── timestamp
 ├── reason
 ├── originating_task
 ├── files
 ├── symbols
 ├── expected_effects
 ├── actual_effects
 ├── risks
 ├── tests
 ├── verification_runs
 ├── baseline
 ├── score_before
 ├── score_after
 └── status
```

The system should capture enough history to understand:

```text
what changed
why it changed
what was expected
what actually happened
what tests were run
what improved
what regressed
what decisions were affected
```

Do not duplicate Git unnecessarily.

Investigate how Trace should reference Git commits/diffs while keeping Trace-specific semantic information separate.

---

# 11. Test vs Verification vs Evaluation

Treat these as separate concepts.

## Test

Answers:

> Does a specific condition hold?

Example:

```text
login_test = PASS
```

## Verification

Answers:

> Does the implementation satisfy the requirement/goal?

Example:

```text
Requirement:
Users remain authenticated after restart.

Evidence:
integration test
database validation
manual inspection

Result:
VERIFIED
```

## Evaluation

Answers:

> How good is this implementation relative to a baseline or alternative?

Example:

```text
Correctness       0.98
Performance       0.87
Complexity        0.72
Coverage          0.94
Maintainability   0.81
```

Design the data model and execution flow for these separately.

---

# 12. Verification and Evaluation Gates

An agent should not be able to simply claim that a task is done after editing files.

Consider:

```text
IMPLEMENTATION
      ↓
TEST GATE
      ↓
VERIFICATION GATE
      ↓
REGRESSION GATE
      ↓
QUALITY/EVALUATION GATE
      ↓
PROMOTE
```

Tasks may therefore have states such as:

```text
implemented
verified
evaluated
promotable
blocked
rejected
needs_rework
```

Some failures should be hard blockers.

For example, a critical correctness regression should not be offset by a performance improvement.

---

# 13. Baselines

Evaluation needs a reference point.

Design a baseline concept.

Example:

```text
Baseline B100

Commit:
abc123

Tests:
182/182

Performance:
p95 = 310ms

Memory:
350MB

Quality:
0.81
```

A change is evaluated against an appropriate baseline.

After acceptance:

```text
B100 → B101
```

The system should preserve historical baselines so agents can understand project evolution.

---

# 14. Actual Effects vs Expected Effects

This is one of the most important capabilities.

Before implementation:

```text
Expected effects:

Latency: improve
DB load: decrease
Correctness: unchanged
Memory: unchanged
```

After implementation:

```text
Actual effects:

Latency: +22% improvement
DB load: -34%
Memory: +18%
Error rate: +0.1%
```

Trace should compare expectation vs reality.

This comparison should generate structured observations:

```text
expected_effect
actual_effect
supported / partially_supported / contradicted
confidence
evidence
```

---

# 15. Regression Detection

The system should detect regressions across:

* tests
* requirements
* performance
* memory
* error rate
* compatibility
* security
* behavior
* other project-specific metrics

A regression should be able to automatically trigger:

```text
Regression
    ↓
Question
    ↓
Investigation
    ↓
Hypothesis
    ↓
Experiment
    ↓
Re-evaluation
```

Do NOT blindly claim causality.

For example, distinguish:

```text
correlated_with
```

from:

```text
caused_by
```

A change may correlate with a regression without being proven as the cause.

Support hypotheses such as:

```text
Hypothesis:
Cache invalidation may be responsible for the error increase.

Confidence:
0.63
```

Then update confidence as evidence accumulates.

---

# 16. Experiments

Consider making experiments first-class objects.

An experiment might compare:

```text
Baseline
Candidate A
Candidate B
```

and evaluate:

* correctness
* latency
* memory
* CPU
* errors
* tests
* complexity
* other relevant metrics

This allows Trace to make evidence-backed decisions rather than relying only on subjective model reasoning.

Design the smallest practical version that supports future experimentation.

---

# 17. Historical Cause/Effect Knowledge

Trace should eventually learn relationships such as:

```text
auth middleware
    → websocket behavior

cache changes
    → memory usage

ORM query changes
    → API latency

serialization changes
    → compatibility behavior
```

But the first implementation should distinguish:

```text
observed relationship
```

from:

```text
proven causal relationship
```

Potential representation:

```text
Relationship
 ├── source
 ├── target
 ├── relationship_type
 ├── observations
 ├── positive_effects
 ├── negative_effects
 ├── confidence
 └── evidence
```

Historical observations should influence future testing and planning.

For example:

```text
This component historically affects websocket behavior.

Therefore:
expand regression testing when modifying it.
```

---

# 18. Risk-Adaptive Verification

Testing should not always be identical.

Trace should eventually use project history and change metadata to determine verification scope.

Example:

A low-risk internal symbol rename might require:

```text
build
unit tests
static analysis
```

A database schema change might require:

```text
unit tests
integration tests
migration validation
rollback validation
API compatibility tests
```

A historically fragile component may automatically trigger broader regression testing.

Design this as a future extensible policy system.

---

# 19. Reflection and Learning

After implementation/evaluation, Trace should ask:

```text
What did we learn?

Which assumptions were wrong?

Which architecture details became clearer?

Which dependencies were discovered?

Which hypotheses were supported/rejected?

Which changes had unexpected side effects?

Which testing strategies were useful?

Should future tasks involving these components get broader verification?
```

The result should update the persistent project knowledge.

This is what makes Trace progressively smarter over time.

---

# 20. Verification Debt / Quality Debt

Consider tracking incomplete validation separately from implementation completion.

Example:

```text
Implementation:
complete

Verification:
partial

Missing:
performance benchmark
manual compatibility validation
2 integration cases

Status:
verification debt
```

This should prevent Trace from presenting uncertain or insufficiently verified work as fully complete.

---

# 21. The Entire Feedback Architecture

Aim conceptually toward:

```text
USER
  ↓
UNDERSTAND
  ↓
IDENTIFY GAPS
  ↓
INVESTIGATE
  ↓
EXPLORE
  ↓
DECIDE
  ↓
PLAN
  ↓
CRITIQUE
  ↓
IMPLEMENT
  ↓
TEST
  ↓
VERIFY
  ↓
MEASURE
  ↓
SCORE
  ↓
COMPARE TO BASELINE
  ↓
DETECT REGRESSIONS
  ↓
REFLECT
  ↓
UPDATE GRAPH / KNOWLEDGE
  ↓
        ┌─────────────────────┐
        │                     │
        │ no major problems   │
        │                     │
        ▼                     │
      CONTINUE                │
                              │
        ▲                     │
        │                     │
        │ problem/gap         │
        │                     │
        └── INVESTIGATE ◄─────┘
```

But implement this through dynamic state transitions rather than a rigid fixed workflow.

---

# 22. Agent Interaction Model

Design how an underlying agent interacts with Trace.

Trace should provide the agent with relevant structured context rather than dumping the entire graph.

The agent should be able to ask Trace for things such as:

```text
current objective
open questions
blocking uncertainties
relevant requirements
related tasks
affected files
historical changes
known risks
relevant decisions
previous evaluations
known regressions
historical effects of changing these components
```

The agent should return structured artifacts such as:

```text
question
finding
assumption
hypothesis
decision
plan update
task update
expected effect
change intent
verification result
evaluation result
reflection
```

Design the protocol/API for this interaction.

---

# 23. Context Selection

Do not assume the whole graph should be sent to the model every time.

Design a context-selection strategy that retrieves only relevant information.

Consider relevance based on:

* current task
* affected files
* dependency relationships
* requirements
* active questions
* recent changes
* historical effects
* risks
* decisions
* regressions
* current deliberation phase

The design should minimize token cost while retaining useful context.

---

# 24. Model-Agnostic Architecture

Trace should not assume one specific LLM.

The deliberation system must work with different agents/models/harnesses.

The model should receive a phase-specific objective and structured Trace context.

Do not hardcode reasoning behavior that only works for a specific provider.

The architecture should make it possible to compare:

```text
Model A
Model B
Model C
```

and potentially:

```text
Model + Harness + Trace policy
```

configurations.

---

# 25. Observability and Auditability

Every important transition should be inspectable.

We should be able to answer:

```text
Why did Trace decide to investigate?

Why did it create this task?

Why did it choose this implementation?

Why did it reject another option?

Why did it trigger additional tests?

Why did it decide that something was a regression?

Why did it replan?

Why was the implementation considered complete?

What evidence supported that conclusion?
```

Design the event/history model needed for this.

---

# 26. Avoid Premature Complexity

This is an architectural planning task, not an invitation to build an enormous autonomous research system immediately.

Explicitly separate:

## MVP

What is necessary to make the core concept work.

## Phase 2

Important improvements once the foundation works.

## Future / Experimental

Ideas that should not be implemented until evidence shows they are useful.

Pay special attention to:

* token/inference cost
* latency
* state size
* graph complexity
* determinism
* debuggability
* migration complexity
* agent reliability
* false-positive regressions
* unnecessary reasoning loops

---

# 27. Existing Trace Repository

First inspect the current repository and implementation.

You must determine:

* current architecture
* language/framework
* current Trace core abstractions
* current project graph implementation
* current task/todo system
* current agent loop
* current state model
* current Git integration
* existing APIs/interfaces
* persistence layer
* event system
* tests
* extension points
* current implementation status
* what has already been implemented versus merely planned

Do NOT design a parallel architecture without first understanding the existing code.

Preserve useful existing work where possible.

---

# 28. Important Existing Feature

The current implementation already contains an initial **agent loop / progressive planner loop**, although it has not yet been thoroughly tested.

Treat that loop as an existing prototype.

Determine:

* what it currently does
* what abstractions are reusable
* what should be retained
* what should be generalized
* what should be replaced
* how it can become the execution mechanism for the broader deliberation controller

Do not automatically throw it away.

---

# 29. Required Deliverable

Produce a **complete implementation plan**.

The plan must include:

## A. Architecture

Describe the proposed architecture and all major components.

## B. Data model

Define the required entities and relationships.

At minimum investigate:

```text
Goal
Requirement
Constraint
Task
Question
Uncertainty
Assumption
Finding
Hypothesis
Option
Decision
Risk
Change
Test
Verification
Evaluation
Baseline
Metric
Observation
Regression
Experiment
Reflection
Evidence
Deliberation
Phase
Relationship
```

Do not blindly implement every entity if some should be merged.

Explain your final choices.

## C. Graph model

Explain:

* nodes
* edges
* graph mutations
* historical data
* semantic relationships
* versioning
* invalidation
* provenance

## D. Deliberation engine

Explain:

* phase model
* state representation
* transition logic
* entry conditions
* exit conditions
* stopping conditions
* retries
* loops
* budget limits
* failure handling

## E. Agent/Trace protocol

Define the interaction between Trace and external agents.

Include expected inputs/outputs and structured artifacts.

## F. Testing architecture

Explain:

* test discovery
* test selection
* test execution
* test results
* test evidence
* historical test relevance
* risk-adaptive test selection

## G. Verification architecture

Explain how requirements and goals are verified independently of simple test pass/fail.

## H. Evaluation architecture

Explain:

* metrics
* scoring
* baselines
* comparisons
* quality gates
* confidence
* regressions

## I. Change/effect tracking

Explain exactly how Trace tracks:

```text
change
→ expected effect
→ actual effect
→ evidence
→ outcome
```

## J. Historical learning

Explain how Trace accumulates knowledge about:

* changes
* side effects
* component relationships
* regressions
* successful approaches
* failed approaches
* testing patterns
* risk patterns

## K. Storage / persistence

Explain how this should fit into the current persistence model.

Do not introduce a graph database automatically unless the current project actually needs one.

Evaluate whether the existing database can support the graph representation.

## L. Git integration

Explain how Trace should associate semantic changes with:

* commits
* diffs
* branches
* worktrees
* files
* agents
* tasks

Avoid duplicating functionality Git already provides.

## M. APIs / interfaces

Define the important interfaces/types/modules.

## N. State machine / policy engine

Provide a concrete proposed transition model.

## O. Failure modes

Consider:

* hallucinated findings
* false verification
* flaky tests
* misleading benchmarks
* false regression detection
* incorrect causal attribution
* stale graph state
* conflicting decisions
* infinite deliberation loops
* runaway token cost
* agent refusing to follow phase requirements
* incomplete evidence
* partial implementation

Explain safeguards.

## P. Security / trust

Explain where agent output must be treated as untrusted and what should require actual machine-verifiable evidence.

## Q. Testing the Trace subsystem itself

Design tests for:

* phase selection
* transitions
* stopping conditions
* graph mutations
* change tracking
* verification
* scoring
* regression detection
* historical relationships
* agent protocol
* failure handling

## R. Migration strategy

Explain how the existing Trace implementation can transition to the new architecture incrementally.

## S. Implementation roadmap

Break implementation into logical milestones.

For every milestone provide:

* objective
* exact components/modules affected
* new types/interfaces
* dependencies
* tests
* acceptance criteria
* risks
* expected outcome

## T. Prioritization

Classify work as:

```text
Must Have
Should Have
Could Have
Future
```

---

# 30. Very Important: Challenge the Concept

Do not simply agree with this proposal.

Critically evaluate it.

Identify:

* concepts that are likely unnecessary
* concepts that overlap
* places where the design risks becoming over-engineered
* places where an LLM is likely to exploit or ignore the mechanism
* places where deterministic software should replace LLM judgment
* places where machine-verifiable evidence is required
* potential token/cost explosions
* likely performance bottlenecks
* graph complexity risks
* ways the system could create useless deliberation loops
* ways scoring could become misleading
* ways historical "cause/effect" relationships could become unreliable

Where appropriate, simplify the design.

The objective is not maximum complexity.

The objective is the **smallest architecture that can create the intended progressive reasoning + implementation validation + historical learning behavior**.

---

# 31. Desired Outcome

At the end of your plan, Trace should conceptually enable an agent to do something like:

```text
User:
"Add feature X."

Trace:
Here is the goal, relevant requirements, existing architecture,
current tasks, previous decisions, open questions, risks,
and historical effects of related changes.

Agent:
I need to investigate 3 unknowns.

Trace:
Records questions.

Agent:
Investigates and reports findings.

Trace:
Updates graph.

Agent:
Generates 3 approaches.

Trace:
Records alternatives.

Agent:
Chooses approach B based on evidence.

Trace:
Records decision.

Agent:
Creates implementation plan.

Trace:
Critiques the plan and identifies one missing dependency.

Agent:
Updates plan.

Trace:
Allows implementation.

Agent:
Makes changes.

Trace:
Records semantic changes and expected effects.

Trace:
Runs relevant tests based on task type and project history.

Trace:
Verifies requirements.

Trace:
Measures effects against baseline.

Trace:
Detects that latency improved but memory regressed.

Trace:
Blocks promotion because the regression exceeds policy.

Agent:
Investigates the regression.

Trace:
Records hypothesis and evidence.

Agent:
Makes a corrective change.

Trace:
Re-runs evaluation.

Trace:
Confirms improvement.

Trace:
Promotes the new baseline.

Trace:
Updates historical knowledge that this component
affects memory and latency.

Future agent:
Receives that knowledge automatically.
```

That is the target behavior.

---

# 32. Final instruction

Before presenting the implementation roadmap:

1. Inspect the existing Trace repository thoroughly.
2. Understand what is already implemented.
3. Map the proposed architecture onto the existing codebase.
4. Identify conflicts and missing abstractions.
5. Prefer incremental evolution over unnecessary rewrites.
6. Explicitly identify the smallest viable first implementation.
7. Explain what should be tested before expanding the system.
8. Clearly separate architecture that is justified now from speculative future capabilities.

Your final response should be an **engineering-grade implementation plan**, not generic conceptual advice.

It should contain enough concrete detail that a senior coding agent can begin implementation directly from the plan.
