# Security and Trust Boundaries

## 1. Core principle

The system is an authority layer for project knowledge and agent workflow. It must not pretend that prompts are security boundaries.

## 2. Local deployment

Initial daemon should:

- bind only to loopback;
- optionally require a local token;
- avoid exposing source/project APIs to the LAN by default.

## 3. Agent trust

Agents are untrusted with respect to:

- truth claims;
- project-plan correctness;
- declared file scope;
- verification assertions.

They are trusted only to the extent that the execution environment gives them access.

## 4. Verification

Do not accept:

- “tests passed” without executable evidence where tests are required;
- “human verified” from an agent;
- “no regressions” based only on an LLM assertion.

## 5. Path scope

Path scope is an advisory coordination mechanism unless the execution environment enforces it.

For stronger isolation, use:

- worktrees;
- containers;
- OS-level sandboxing;
- restricted credentials.

## 6. Secrets

The graph should avoid ingesting secrets.

Store references/metadata rather than secret values.

Environment captures should have secret-redaction policies.

## 7. External tools and MCPs

Capability records should identify:

- tool origin;
- permissions;
- credential scope;
- destructive operations;
- network access.

High-risk capabilities should be visible in task context.

## 8. Prompt injection

Project text, source code, comments, tests, issue text, and external tool output may contain hostile instructions.

The planner must distinguish:

- project policy;
- task instruction;
- source content;
- retrieved evidence.

Retrieved source text must not automatically become system-level authority.
