# Evaluation and Research Plan

## 1. Main hypotheses

### H1 — Project understanding

A persistent project graph improves an agent’s ability to answer project-understanding questions.

### H2 — Progressive planning

Progressive planning improves long-horizon implementation outcomes compared with exhaustive upfront planning.

### H3 — Discovery-driven replanning

Recording discoveries and propagating them to future work reduces downstream rework.

### H4 — Decision impact

Decision-aware impact analysis catches meaningful consequences before destructive plan changes.

### H5 — Evidence-driven review

Independent review plus deterministic evidence reduces false completion and regressions.

### H6 — Progressive context

Task-specific retrieval reduces token usage without reducing task success.

### H7 — Capability-aware planning

Selecting only the relevant skills/tools/MCPs improves efficiency and reduces environmental noise.

## 2. Benchmark levels

### Repository understanding

Questions:

- why does this file exist?
- what responsibility does it have?
- what created it?
- what decisions constrain it?
- what does modifying it affect?

### Planning

Tasks:

- decompose a goal;
- plan a scope;
- identify unknowns;
- revise future tasks after discovery.

### Implementation

Measure:

- completion;
- regressions;
- unnecessary edits;
- review failures;
- tokens;
- latency.

### Decision analysis

Prompt mid-project decisions that deliberately conflict with existing plans.

Measure:

- impacted work identified;
- invalidated assumptions identified;
- useful alternatives;
- false positives;
- missed consequences.

### Verification

Inject controlled bad implementations and measure whether review rejects them.

## 3. Baselines

Compare:

1. raw agent;
2. raw agent + normal repository tools;
3. agent + code graph;
4. agent + project graph;
5. agent + graph + progressive context;
6. agent + graph + progressive planner;
7. agent + graph + planner + capability selection.

## 4. Metrics

### Quality

- task success;
- regression count;
- missed dependency count;
- false completion rate;
- review escape rate;
- decision-impact precision/recall.

### Efficiency

- input tokens;
- output tokens;
- total model calls;
- latency;
- human intervention count.

### Planning

- plan churn;
- downstream rework;
- discovery propagation accuracy;
- number of stale tasks caught.

## 5. Controlled honesty tests

Create tasks where the implementation agent is likely to make false claims, such as:

- test command intentionally unavailable;
- partial implementation;
- hidden regression;
- missing edge case;
- incorrect fixture;
- UI requirement impossible to verify automatically.

The reviewer must be judged on whether it independently discovers the problem.

## 6. Reproducibility

Benchmark datasets, task definitions, and scoring procedures should be public where licensing and repository constraints allow.
