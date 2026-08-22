# Trace — thin convenience targets (optional). Prefer scripts/ for CI-less workflows.

.PHONY: embed-gui

embed-gui:
	./scripts/embed-gui.sh
