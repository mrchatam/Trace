# PROMPT-B0 — Equipment checkout desk (Baseline, No Trace)

> **Arm:** B0  
> **Experiment:** E02 — ab-p25-gap-pass-validation

## STOP — verify workspace first

Your Cursor workspace **must** be:

`/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/B0`

**Turn 1:** Run `pwd`. If the path is not exactly the line above, **stop** and tell the operator to use File → Open Folder on that directory.

## Requirement

Build a full **equipment checkout desk** application.

It must support:
- Asset catalog (laptops, monitors, keys, etc.)
- Checkout and return lifecycle
- User roles (requester, desk staff, admin)
- REST API for core operations
- Desk staff / admin UI for managing checkouts
- Public read-only **availability** page (what is free now)
- Audit log of checkout events

This is intentionally underspecified. Decide stack, schema, API, auth, layout, and testing during implementation.

## Workspace

**Project root:** `/home/ali/Desktop/Trace/experiments/ab-p25-gap-pass-validation/runs/B0`

- Edit **only** under this path
- Do **not** read or modify `runs/G1`
- Do **not** use Trace

## Multitask orchestrator

If you delegate to subagents:

1. Include this exact path in **every** subagent prompt
2. Add: **do not use any other directory**
3. Subagents run `pwd` before any work

## Constraints

- Working app with `go test ./...` passing
- Document how to run the server in README
