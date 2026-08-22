// Command trace is the Trace CLI entrypoint (thin adapter over library APIs).
package main

import (
	"os"
)

const version = "0.0.0-dev"

const (
	exitOK    = 0
	exitUsage = 1
	exitFail  = 2
)

func main() {
	os.Exit(run(os.Args[1:]))
}
