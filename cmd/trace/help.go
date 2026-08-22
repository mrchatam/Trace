package main

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer) {
	fmt.Fprint(w, `trace — local-first project knowledge graph for AI coding agents

Usage:
  trace [-C <root>] <command> [args]

Global:
  -C, --project <dir>   Project root (default: cwd). Resolved with filepath.Abs only
                        (no parent .trace walk-up). One open store per root
                        (.trace/trace.lock; single-writer intentional). Serialize
                        CLI↔MCP (or parallel Trace) on one root, or use separate
                        -C / worktree roots.

Commands:
  help | -h | --help    Show this help (exit 0)
  version | --version   Print 0.0.0-dev (exit 0)
  init [--with-agent-defaults]
                        Create/open .trace/trace.db (not project-root trace.db);
                        print DB path; optional bundled harness agent catalog install
  index [paths...]      Index supported source files (file-local; no full rebuild)
  index status          JSON: head, last_indexed_commit, stale, hook_installed,
                        supported_languages (see docs/INDEX_LANG_POLICY.md)
  index watch [--debounce 300ms] [paths...]
                        Foreground fsnotify watcher; debounced IndexFile per change;
                        SIGINT exits (not a daemon). Git-hook remains primary freshness.
  reindex [paths...]    Alias of index (same file-local IndexFile path)
  add <kind> …          Create entity via domain (goal|task|decision|assumption|
                        discovery|plan-change|claim|evidence)
                        Promotion path: after a BLOCKING discovery, run
                        add task --from-discovery <id> [--goal-id <id>]
  link <rel> …          Link via domain (goal-task|decision-task|
                        discovery-plan-change|discovery-mentions-task|claim-evidence)
  transition …          Transition task work_state via domain.TransitionTask
                        (DONE needs Review PASS + --as-operator, or --allow-done hatch;
                        linked FAIL blocks DONE even with sibling PASS; --evidence alone
                        ≠ DONE; --as-operator is a conscious claim, flag≠identity / not
                        verified operator identity; Actor string ≠ auth; missing caps
                        fail-closed unless --allow-missing-caps; --allow-done does not
                        bypass missing caps; --allow-done prints a loud WARNING on success).
                        Optional --enforce runs deliberation gate (GateForDone) before
                        DONE; exit 1 when blocked. Without --enforce, behavior unchanged.
                        --enforce does not bypass review PASS/--as-operator domain checks
                        when gate allows.
  review create|set|get|show|list|residual …
                        Create review (+optional --task/--scope), set result,
                        get|show by --id, list [--task], or residual add|list
  impact …              Decision impact (finding add|list; alternative add|list|recommend;
                        report; walk --seed file:<id>|symbol:<id> [--depth 1|2];
                        predict --change <id> [--depth 1|2]; compare --change <id>)
  capability …          Capability catalog (declare|list|require|unrequire|missing|
                        decide|decisions)
  plan …                Progressive coarse planner (create-coarse|set-current|deep|show)
  seed import <file>    Import seed JSON v1 through domain APIs. Relative <file> resolves
                        under -C project root; absolute paths unchanged. Stdout JSON may
                        include promotion_candidates (orphan BLOCKING discoveries) plus
                        promotion_hint; import does not auto-spawn tasks.
  seed export [-o <file>] [--strict] [--enforce] [--task <id>]
                        Export seed JSON v1 (causal entities, links, plan tree,
                        findings/alternatives when present) to stdout or -o.
                        Recommended commit path: trace/graph.json (.trace/ stays local).
                        Sets exported_at_commit (git SHA evidence, not identity)
                        when -C root is a git repo.
                        Default export omits reviews, transitions, and task work_state.
                        After seed import, clone tasks are PENDING until the clone
                        operator transitions them.
                        --strict validates export honesty (GateForExport per active
                        task, or --task only). --enforce requires --strict; exit 1 and
                        no write on violation. --strict alone warns on stderr, exit 0.
                        Without --strict, thin graphs (discoveries=0 decisions=0) still
                        emit a stderr warn (write-before-export nudge); exit stays 0.
  tasks [--goal <id>]   List tasks as JSON array (id, title, work_state, goal_id);
                        empty → []; optional --goal filters by goal_id
  tasks conflicts [--task <id>]
                        Advisory overlaps between active tasks (JSON ok + conflicts[])
  why <type> <id>       Causal explanation (JSON WhyResult on stdout)
  context <task-id>     Context packet (json|markdown|both)
  loop next --task <id> Emit one bounded JSON loop packet for a seed task;
                        derives goal context from task.goal_id, no stdin
  loop apply [--in <path>]
                        Apply trace.loop.apply.v1 writes from --in JSON file or stdin
                        Promotion path: include writes.spawned_tasks[].discovery_id
  loop status --task <id> [--goal <id>]
                        Report trace.loop.status.v1 from persisted loop-step evidence.
                        Includes violations[] (edit gate parity). Optional
                        .trace/config.json enforce mode (off|warn|strict, default off):
                        warn/strict print stderr hints when violations present; exit stays 0.
  .trace/config.json    Local enforce mode: { "enforce": "off"|"warn"|"strict" }.
                        Default off when missing. Does not auto-enable --enforce on
                        transition or seed export (use explicit flags).
  loop gate --task <id> [--for orient|edit|execute|done|export]
                        Check deliberation gate; emits trace.loop.gate.v1 JSON on stdout.
                        Exit 0 allowed, 1 blocked, 2 usage or internal error (default --for edit)
  agents list           List harness agent catalog profiles; JSON array on stdout
  agents describe <slug>
                        Full agent profile + requirements + registry metadata
  agents recommend (--task <id> | --phase <PHASE>) [--goal-keywords "kw ..."]
                        Ranked harness agent recommendations (max 4); no spawn
  migrate status        Report applied / embed migration versions
  backup -o <path>      Snapshot .trace/trace.db (optional --include-token)
  restore <path>        Install snapshot into .trace/trace.db ([--force])
  auth set|clear|status Local .trace/access.token gate (TRACE_ACCESS_TOKEN)
  install detect|uninstall <target>|agents|cursor|claude|cursor-hook|git-hook …
                        Marker-gated install registry: detect (JSON, no writes);
                        install agents (upsert bundled harness agent catalog into
                        .trace/); install cursor [--write] [--bin] [--mcp-json]
                        (MCP + project enforcement rules .cursor/rules/trace-enforcement.mdc
                        + AGENTS.md block; print: MCP JSON stdout, rules hint stderr);
                        install claude [--write] (CONDITIONAL: needs .claude/ or
                        CLAUDE.md; MCP + claude rules + AGENTS.md block);
                        install cursor-hook [--write] (CONDITIONAL: needs .cursor/;
                        preToolUse hook → trace loop gate --for edit; TRACE_TASK_ID);
                        install git-hook [--write] (CONDITIONAL: needs
                        git work tree; post-commit + pre-push fragments, never wraps
                        git commit); uninstall <target>. Environment: TRACE_TASK_ID,
                        TRACE_PROJECT_ROOT. Cursor print/--write keeps DF-22 reload tip;
                        confirm with trace_version.
  changes capture [--since <oid>] [--all]
                        Promote indexed VCS commits into changes (meaningful paths
                        only unless --all); JSON array on stdout
  changes compare --from <oid> --to <oid>
                        Compare project states; JSON on stdout
  changes list [--task <uuid>] [--limit N]
                        List recent changes newest-first; JSON array on stdout
  changes show <change-id>
                        Show one change and path refs (no file content); JSON on stdout
  changes similar --path <prefix> | --kind <kind> [--limit N]
                        Query historical changes similar by path prefix or change kind; JSON on stdout
  patterns refresh      Rebuild change_patterns from stored changes; JSON {ok,patterns_updated}
  patterns list [--limit N]
                        List stored change_patterns rows; JSON array on stdout
  knowledge list [--topic <t>] [--limit N]
                        List engineering_knowledge rows; JSON array on stdout
  knowledge synthesize  Synthesize knowledge from patterns and cognition; JSON {ok,created,updated}
  knowledge tendencies [--limit N]
                        List change-pattern tendencies (help/hurt); JSON array on stdout
  search <query> [--limit N]
                        Full-text search over indexed entities; JSON {ok,hits,count} on stdout
  test run --task <id> [--paths path,...]
                        Run relevant tests for a task; record kind=test outcomes; JSON on stdout
  tests verifying --symbol <uuid> | --file <path>
                        List tests that validate a symbol or file via validates edges; JSON on stdout
  verify run --task <id> [--force-eval]
                        Coordinate test → verification → evaluation cycle; JSON on stdout
  verify invariants --task <id>
                        Advisory architectural invariant check on latest change paths; JSON on stdout
  eval rules            Show project eval-rules.json (or defaults); JSON on stdout
  eval results --task <id>
                        List stored test/verification/evaluation outcomes; JSON on stdout
  outcomes compare --task <id> --kind test|evaluation
                        Compare last two stored outcomes of kind; JSON {previous,current,delta}
  outcomes improvements --change <id> | --task <id>
                        List recorded improvements; JSON on stdout
  outcomes failed [--task <id>] [--limit N]
                        List stored test failures/errors; JSON on stdout
  outcomes worked [--task <id>] [--limit N]
                        List improvements and passing tests (worked approaches); JSON on stdout
  regressions list [--task <id>] [--change <id>] [--limit N]
                        List regression history; JSON on stdout
  serve [--addr host:port] [--allow-remote] [--token TOKEN] [--token-file PATH]
                        [--root DIR] [--static-dir DIR] [--cors-origin URL]
                        Opt-in local HTTP API (default 127.0.0.1:7432). Non-loopback
                        bind requires --allow-remote and a bearer token (--token,
                        --token-file, or auto-generated once). Default busy → next
                        free loopback port (7432–7441); --addr pins (fail if busy).
                        CORS deny (no *); optional --cors-origin exact Origin
                        reflect only. Do not set --static-dir to the project root
                        (refused). Serves embedded Explore SPA (or disk web/dist
                        when present); JSON under /v1.
  gui [--addr host:port] [--allow-remote] [--token TOKEN] [--token-file PATH]
                        [--root DIR] [--static-dir DIR] [--cors-origin URL] [--no-open]
                        Consumer: cd your-project && trace gui (GUI embedded in binary;
                        project needs only .trace/). -C/--project before or after gui.
                        Same bind/token policy as serve; opens browser to Explore (/).

Handoff:
  Thin handoff SoT — predecessor task body (+ linked decisions via retrieval) is
  the source of truth (no first-class handoff entity). Successors must Trace-pull
  with trace context / trace why (and trace tasks as needed).

Build note:
  Install CLI on PATH (from a Trace checkout — module not published to proxy yet):
    CGO_ENABLED=1 go build -o bin/trace ./cmd/trace
    cp -f bin/trace ~/.local/bin/trace   # or any dir on PATH
  Then from a Trace project: trace gui
  (CGO_ENABLED=1 required for full analyzer-linked binary.)
  MCP stdio server:
    CGO_ENABLED=0 go build -o bin/trace-mcp ./cmd/trace-mcp
  Note: trace install … configures agents/MCP/hooks — it does not put the
  trace binary on PATH.
  Note: go install …/cmd/trace@latest fails until a module version is tagged/published.

Exit codes: 0 success | 1 usage | 2 operational failure
`)
}
