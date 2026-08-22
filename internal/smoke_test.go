package internal_test

import (
	"testing"

	_ "github.com/mrchatam/Trace/internal/analyzers"
	_ "github.com/mrchatam/Trace/internal/compiler"
	_ "github.com/mrchatam/Trace/internal/domain"
	_ "github.com/mrchatam/Trace/internal/gitcli"
	_ "github.com/mrchatam/Trace/internal/retrieval"
	_ "github.com/mrchatam/Trace/internal/store"
	_ "github.com/mrchatam/Trace/internal/vcs"
)

// TestScaffoldPackagesCompile is a smoke check that the locked S01 package
// layout resolves and compiles under module github.com/mrchatam/Trace.
func TestScaffoldPackagesCompile(t *testing.T) {
	t.Log("module github.com/mrchatam/Trace scaffold packages present")
}
