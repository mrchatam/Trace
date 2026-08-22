# Agent Environment and Capabilities

## 1. Why environment belongs in the project model

A task is constrained not only by the codebase but also by what the agent can actually do.

Relevant environment objects include:

- skills;
- rules;
- tools;
- MCP servers;
- hooks;
- commands;
- agents;
- models;
- permissions.

## 2. Capability model

```text
Task
 ├── requires → Skill
 ├── requires → Tool
 ├── requires → MCP
 ├── constrained_by → Rule
 └── verified_by → Capability
```

## 3. Rules

Rules are relatively persistent project constraints.

Examples:

- coding conventions;
- security requirements;
- secret handling;
- architecture invariants.

Do not duplicate task-specific procedural instructions into always-on rules.

## 4. Skills

Skills are dynamically relevant procedures or domain knowledge.

A skill can define:

- name;
- description;
- trigger conditions;
- prerequisites;
- required tools;
- token/cost hints;
- project scope.

Only relevant skills should be activated for a task.

## 5. MCPs and tools

Represent:

- capability offered;
- version;
- availability;
- permission;
- environment;
- task compatibility.

The planner should check capability availability before assigning work that depends on it.

## 6. Hooks

Hooks may provide:

- pre-edit checks;
- post-edit checks;
- post-test automation;
- verification;
- integration steps.

Hooks should be represented as capabilities/policies rather than opaque magic.

## 7. Capability-aware planning

Example:

```text
UI verification task
  ↓
requires:
  browser
  visual-testing skill
  accessibility skill
  browser MCP
  human verification
```

If the environment lacks the needed capability:

- warn;
- propose an alternative;
- request configuration;
- or change the task plan.

## 8. Environment history

Capability changes can affect planning.

Track:

- added/removed capabilities;
- versions;
- compatibility;
- last verified time.

A future task that was planned against an unavailable capability should become re-evaluable.

## 9. Agent adapters

An agent adapter maps the project model to an execution environment.

Examples:

- Claude Code;
- Cursor CLI;
- Codex;
- OpenHands;
- future agents.

The planner should not depend on a single harness.

## 10. Minimal context principle

Environment selection follows the same rule as project retrieval:

> provide only the capabilities needed for the task.
