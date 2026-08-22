# P00 / S01 / 01 — Go module + CLI stub scaffold

## Metadata
- id: P00-S01-01
- todo_ids: [P00-S01-01]
- role: implementer
- skills: [incremental-implementation]
- mcps: [Shell, Read, Write, Glob]
- agents: []
- verification: automated

## Objective
Create Go module **`github.com/mrchatam/Trace`**, thin `cmd/trace` stub, locked `internal/` package layout (stubs only), `.gitignore` (must include `.trace/`), smoke test, and brief README build notes. **No** SQLite, git, analyzer, retrieval, or causal features.

## References
- [agent-loop-protocol.md](../../../../rules/agent-loop-protocol.md)
- [project-rules.md](../../../../rules/project-rules.md)
- [skills-map.md](../../../../rules/skills-map.md)
- [C_FIRST_SCOPE.md](../../../../init/C_FIRST_SCOPE.md) implementation boundaries
- [phase-00 README](../../README.md) package sketch
- [AGENTS.md](../../../../../AGENTS.md)

## Session start
Agent → clarify if needed → Plan → implement.

## Locked defaults

| Item | Value |
|------|-------|
| Module path | `github.com/mrchatam/Trace` (exact; case-sensitive) |
| Go version in `go.mod` | `go 1.22` (host may be newer; do not bump without need) |
| CLI entry | `cmd/trace` → package `main`; binary name `trace` |
| CLI behavior (stub) | `-h` / default help text + `trace version` prints `0.0.0-dev` (or equivalent single line); **no** subcommands beyond help/version |
| Library vs CLI | Library under `internal/` **must not** import `cmd/…`; CLI may import `internal/…` |
| Store dir (gitignore only) | `.trace/` — create **no** DB this scope |
| Out of scope | HTTP, MCP, daemon, domain logic, tree-sitter, migrations, fixtures/evals trees |
| Depends on | Phase locks in `phase-00-foundation/README.md` |
| Repo state at start | Docs-only (no `*.go` / `go.mod` yet) |

### Internal package taxonomy (LOCKED)

Do **not** invent alternate names. Directory = import path suffix; Go package name = last path element unless noted.

| Path | Package | Owns (later scopes) | S01 content |
|------|---------|---------------------|-------------|
| `internal/store` | `store` | `.trace/` SQLite (S02) | stub only (`doc.go` or empty `.go` with package clause) |
| `internal/vcs` | `vcs` | VCS adapter **interface** (S03) | stub |
| `internal/gitcli` | `gitcli` | `git` CLI impl of `vcs` (S03) | stub |
| `internal/analyzers` | `analyzers` | tree-sitter TS/JS + Python (S04) | stub |
| `internal/domain` | `domain` | work/causal API (S05) | stub |
| `internal/retrieval` | `retrieval` | exact / FTS / graph expand (S06) | stub |
| `internal/compiler` | `compiler` | context compiler Layer 0–1 (S06) | stub |

**Naming decision (P00-00 residual):** use **`internal/compiler`**, not `contextx`, not `context`.

- Product docs call this the **context compiler**; `C_FIRST_SCOPE` / phase README already say `compiler`.
- Reject `internal/context` — collides with stdlib `context`.
- Reject `internal/contextx` — prior alternate wording; superseded by this lock.
- Reject `internal/sqlite` — prefer `store` (matches S02 / `.trace/` product language). C_FIRST_SCOPE “sqlite” is historical alias for this path.

**Do not create in S01:** `fixtures/`, `evals/` (S08), `internal/mcp`, `cmd` siblings other than `trace`, public `pkg/` (not needed for P0-X).

### Target tree

```text
/
  go.mod                          # module github.com/mrchatam/Trace ; go 1.22
  go.sum                          # only if deps appear (stubs need none)
  cmd/trace/main.go               # help + version stub
  internal/store/doc.go
  internal/vcs/doc.go
  internal/gitcli/doc.go
  internal/analyzers/doc.go
  internal/domain/doc.go
  internal/retrieval/doc.go
  internal/compiler/doc.go
  internal/smoke_test.go          # or cmd/trace/*_test.go — see exit criteria
  .gitignore                      # must list .trace/ and bin/
  README.md                       # short Build section only (design SoT stays in docs/)
  docs/                           # existing — do not relocate
  LICENSE                         # existing Apache-2.0 — leave intact
```

### `.gitignore` (minimum)

Must include at least:

```gitignore
.trace/
bin/
*.exe
```

Optional but allowed: OS/IDE noise (`.DS_Store`, `.idea/`, etc.).

### README Build section (minimum)

Document:

```bash
go test ./...
go build -o bin/trace ./cmd/trace
./bin/trace version
```

Do not replace design docs; keep the existing product README intro.

## Board rights
Implementer: update **status + notes only** on `P00-S01-01`. Do not spawn rows or rewrite later prompts.

## Exit criteria
- [ ] `go.mod` module path is exactly `github.com/mrchatam/Trace` and declares `go 1.22`
- [ ] All seven `internal/{store,vcs,gitcli,analyzers,domain,retrieval,compiler}` packages exist and compile
- [ ] **No** package named `context`, `contextx`, or `sqlite` under `internal/`
- [ ] `go test ./...` passes (at least one smoke test that imports or references the module)
- [ ] `go build -o bin/trace ./cmd/trace` succeeds; `bin/` is gitignored
- [ ] `./bin/trace version` (or documented equivalent) prints a stub version; help works
- [ ] `.gitignore` contains `.trace/`
- [ ] No HTTP/MCP/daemon packages or listeners
- [ ] LICENSE / `docs/` design baseline left intact (no design SoT rewrite)
- [ ] TODO.md Notes for `P00-S01-01` updated; status `done`

## Minimal todos
- [ ] `go mod init github.com/mrchatam/Trace` + set `go 1.22`
- [ ] `cmd/trace` help + version stub
- [ ] Create seven locked `internal/*` stub packages (names above only)
- [ ] Smoke test + `go test ./...` / `go build`
- [ ] `.gitignore` (`.trace/`, `bin/`) + README Build section
- [ ] Board status + notes (paths created, commands run)
