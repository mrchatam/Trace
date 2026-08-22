// Package traceagents holds embedded harness agent catalog assets.
package traceagents

import _ "embed"

//go:embed default.json
var DefaultCatalogJSON []byte
