// Package compiler builds progressive Layer 0–3 task context packets (DR-PACK).
//
// JSON is the canonical packet form; Markdown is a labeled render of the same
// structure. Retrieved project text is marked trust=untrusted_data (G4/A14);
// structured system fields (ids, enums, reason codes, budgets) may be
// trust=system. Token/item budgets are enforced; there is no dump API (G6).
// Layers 2–3 are opt-in via ContextOptions.MaxLayer (default 1 = L0–L1 only).
package compiler
