# PROMPT-B0 — Incident tracker (Baseline, No Trace)

> **Arm:** B0  
> **Experiment:** E01 — ab-incident-tracker

## STOP — verify workspace first

Your Cursor workspace **must** be:

`/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/B0`

**Turn 1:** Run `pwd`. If the path is not exactly the line above, **stop** and tell the operator to use File → Open Folder on that directory. Do not proceed from the Trace monorepo root or any other folder.

## Requirement

Build a full **on-call incident tracker**.

It must support:
- Multiple severity levels for incidents
- User roles and permissions (reporter, responder, admin)
- Status lifecycle (report → acknowledge → resolve, or equivalent)
- Assignment and comment/timeline on incidents
- Search and filter across incidents
- REST API for all core operations
- Responder/admin dashboard for incident management
- Public read-only status page for leadership
- Activity/audit log

This is intentionally underspecified. Decide stack, schema, API, auth, layout, and testing during implementation.

## Workspace

**Project root:** `/home/ali/Desktop/Trace/experiments/ab-incident-tracker/runs/B0`

- Edit **only** under this path
- Do **not** read or modify `runs/G1`
- Do **not** use Trace (`trace` CLI, `.trace/`, Trace MCP)

## Multitask orchestrator

If you delegate to subagents:

1. Include this exact path in **every** subagent prompt
2. Add: **do not use any other directory**
3. Subagents run `pwd` before any work
4. You remain responsible for workspace correctness — wrong path invalidates the experiment

## Constraints

- Deliver a working app with tests passing (`go test ./...`)
- Document how to run the server and default credentials in README
