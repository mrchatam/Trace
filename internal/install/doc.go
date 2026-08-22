// Package install is the marker-gated install registry (Detect / Install / Uninstall).
//
// Tiers: STABLE | CONDITIONAL | OPT_IN. Cursor is STABLE; claude is CONDITIONAL
// (requires .claude/ or CLAUDE.md under the project root). Logic lives here;
// cmd/trace is a thin adapter only (G19).
package install
