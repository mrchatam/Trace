# Skills / MCP / hooks map (Trace)

Route concerns to skills and tools. Prompts list a **subset**; do not load everything.

## Always consider

| Concern | Skills / tools |
|---------|----------------|
| Unclear requirements | `ask-questions-if-underspecified`, `grilling`, `brainstorming` |
| Session planning | Switch to Plan mode after clarification (see agent-loop-protocol) |
| Task breakdown | `planning-and-task-breakdown` |
| Code changes | language skills; `incremental-implementation` |
| Tests failing / bugs | `systematic-debugging`, `diagnosing-bugs`, `grinding-until-pass` |
| Review | `code-review-and-quality`, `parallel-code-review`, `review-bugbot`, `review-security` |
| Security-sensitive | `security-and-hardening`, `semgrep` (if available) |
| Docs / ADRs | `documentation-and-adrs` |
| Git commit (only if user asks) | git safety in user rules |
| Library docs | MCP `user-context7` |
| Codegraph (if indexed) | MCP `user-codegraph` |

## By role

| Role | Prefer |
|------|--------|
| Phase / scope planner | planning-and-task-breakdown, grilling (if locks unclear), brainstorming (alternatives) |
| Implementer | incremental-implementation, tdd when logic-heavy, context7 for Go/sqlite/tree-sitter |
| Reviewer | code-review-and-quality, silent-failure-hunter mindset, security-reviewer when trust boundaries touched |
| P0-X / eval | Follow `docs/init/I_BENCHMARK_PLAN.md`; deterministic scripts over LLM judges |

## MCPs (typical)

| MCP | Use |
|-----|-----|
| Shell / Read / Grep / Glob / Write | Default coding |
| user-context7 | Go, tree-sitter, sqlite docs |
| user-memory | Persist durable project decisions when asked |
| cursor-ide-browser | Only if a future UI task needs it (not P0) |

## Hooks

None required for P0-X. If repo later adds Cursor hooks (format, secret scan), list them on the prompt’s `hooks:` metadata.

## Human verification triggers

Use `verification: human` when:

- Subjective UX judgment
- Security accept/reject of a destructive policy
- Listening / sensory outcomes (N/A for Trace P0)
- Explicit user gate in the prompt
