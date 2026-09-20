// Command trace is the Trace CLI entrypoint (thin adapter over library APIs).
package main

import (
	"os"

	"github.com/mrchatam/Trace/internal/buildmeta"
)

// version mirrors buildmeta for CLI `trace version` (ldflags-compatible via buildmeta.Version).
var version = buildmeta.String()

const (
	exitOK    = 0
	exitUsage = 1
	exitFail  = 2
)

func main() {
	os.Exit(run(os.Args[1:]))
}
