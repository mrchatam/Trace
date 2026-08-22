package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/config"
	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/install"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

func cmdInit(root string, args []string) int {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	withAgentDefaults := fs.Bool("with-agent-defaults", false, "install bundled harness agent catalog after init")
	args = flagsFirst(args, map[string]bool{"with-agent-defaults": false})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace init [--with-agent-defaults]\n")
		return exitUsage
	}
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "init: %v\n", err)
		return exitFail
	}
	defer st.Close()

	if repo, err := tryOpenGit(abs, st); err == nil {
		defer repo.Close()
		if _, rerr := repo.Refresh(context.Background()); rerr != nil {
			fmt.Fprintf(os.Stderr, "init: git refresh: %v\n", rerr)
			return exitFail
		}
	}

	fmt.Println(st.DBPath())
	config.WarnIfTraceDirWithoutConfig(abs, os.Stderr)

	if *withAgentDefaults {
		if err := install.InstallAgentDefaults(install.InstallOpts{
			Write:       true,
			ProjectRoot: abs,
			ErrOut:      os.Stderr,
		}); err != nil {
			fmt.Fprintf(os.Stderr, "init: %v\n", err)
			return exitFail
		}
	}
	return exitOK
}

// tryOpenGit returns a gitcli repo bound to an already-open store, or an error
// when root is not a git work tree.
func tryOpenGit(absRoot string, st *store.Store) (vcs.Repository, error) {
	repo, err := gitcli.OpenWithStore(absRoot, st)
	if err != nil {
		if isNotRepo(err) {
			return nil, err
		}
		return nil, err
	}
	return repo, nil
}

func isNotRepo(err error) bool {
	var ve *vcs.Error
	if errors.As(err, &ve) {
		return errors.Is(ve.Err, vcs.ErrNotRepo)
	}
	return errors.Is(err, vcs.ErrNotRepo)
}
