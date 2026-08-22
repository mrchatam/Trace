package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/install"
	"github.com/mrchatam/Trace/internal/store"
)

// installCursorReloadTip aliases the library tip so existing TestInstallCursor* keep matching.
const installCursorReloadTip = install.CursorReloadTip

func cmdInstall(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace install detect|uninstall <target>|agents|cursor|claude|cursor-hook|git-hook …\n")
		return exitUsage
	}
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	switch args[0] {
	case "detect":
		return cmdInstallDetect(abs, args[1:])
	case "uninstall":
		return cmdInstallUninstall(abs, args[1:])
	case install.TargetCursor:
		return cmdInstallCursor(abs, args[1:])
	case install.TargetClaude:
		return cmdInstallClaude(abs, args[1:])
	case install.TargetGitHook:
		return cmdInstallGitHook(abs, args[1:])
	case install.TargetCursorHook:
		return cmdInstallCursorHook(abs, args[1:])
	case "agents":
		return cmdInstallAgents(abs, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown install target: %s\n", args[0])
		fmt.Fprintf(os.Stderr, "usage: trace install detect|uninstall <target>|agents|cursor|claude|cursor-hook|git-hook …\n")
		return exitUsage
	}
}

func cmdInstallAgents(projectRoot string, args []string) int {
	fs := flag.NewFlagSet("install agents", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	args = flagsFirst(args, nil)
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace install agents\n")
		return exitUsage
	}
	opts := install.InstallOpts{
		Write:       true,
		ProjectRoot: projectRoot,
		ErrOut:      os.Stderr,
	}
	if err := install.InstallAgentDefaults(opts); err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdInstallDetect(projectRoot string, args []string) int {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace install detect\n")
		return exitUsage
	}
	infos := install.ListTargets(install.InstallOpts{ProjectRoot: projectRoot})
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(infos); err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdInstallUninstall(projectRoot string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace install uninstall <target> [--mcp-json path]\n")
		return exitUsage
	}
	targetID := args[0]
	rest := args[1:]
	tgt, err := install.Lookup(targetID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitUsage
	}

	fs := flag.NewFlagSet("install uninstall", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	mcpJSON := fs.String("mcp-json", "", "path to mcp.json (cursor; default $HOME/.cursor/mcp.json)")
	rest = flagsFirst(rest, map[string]bool{"mcp-json": true})
	if err := fs.Parse(rest); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace install uninstall <target> [--mcp-json path]\n")
		return exitUsage
	}

	opts := install.InstallOpts{
		MCPJSON:     *mcpJSON,
		ProjectRoot: projectRoot,
		ErrOut:      os.Stderr,
	}
	if err := tgt.Uninstall(opts); err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	if targetID == install.TargetGitHook {
		if err := setHookInstalledFlagForProject(projectRoot, false); err != nil {
			fmt.Fprintf(os.Stderr, "install: %v\n", err)
			return exitFail
		}
	}
	return exitOK
}

func cmdInstallCursor(projectRoot string, args []string) int {
	fs := flag.NewFlagSet("install cursor", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	write := fs.Bool("write", false, "upsert mcpServers.trace into Cursor mcp.json")
	bin := fs.String("bin", "trace-mcp", "trace-mcp binary (PATH name or absolute path)")
	mcpJSON := fs.String("mcp-json", "", "path to mcp.json (default: $HOME/.cursor/mcp.json)")
	args = flagsFirst(args, map[string]bool{"write": false, "bin": true, "mcp-json": true})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace install cursor [--write] [--bin path] [--mcp-json path]\n")
		return exitUsage
	}

	tgt, err := install.Lookup(install.TargetCursor)
	if err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	opts := install.InstallOpts{
		Write:       *write,
		Bin:         *bin,
		MCPJSON:     *mcpJSON,
		ProjectRoot: projectRoot,
		Out:         os.Stdout,
		ErrOut:      os.Stderr,
	}
	if err := tgt.Install(opts); err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	install.PrintBootstrapHintIfNeeded(projectRoot, os.Stderr)
	return exitOK
}

func cmdInstallClaude(projectRoot string, args []string) int {
	fs := flag.NewFlagSet("install claude", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	write := fs.Bool("write", false, "write .claude/trace-mcp.json when marker present")
	bin := fs.String("bin", "trace-mcp", "trace-mcp binary (PATH name or absolute path)")
	args = flagsFirst(args, map[string]bool{"write": false, "bin": true})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace install claude [--write] [--bin path]\n")
		return exitUsage
	}

	tgt, err := install.Lookup(install.TargetClaude)
	if err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	opts := install.InstallOpts{
		Write:       *write,
		Bin:         *bin,
		ProjectRoot: projectRoot,
		Out:         os.Stdout,
		ErrOut:      os.Stderr,
	}
	if err := tgt.Install(opts); err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	install.PrintBootstrapHintIfNeeded(projectRoot, os.Stderr)
	return exitOK
}

func cmdInstallCursorHook(projectRoot string, args []string) int {
	fs := flag.NewFlagSet("install cursor-hook", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	write := fs.Bool("write", false, "install preToolUse hook calling trace loop gate")
	args = flagsFirst(args, map[string]bool{"write": false})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace install cursor-hook [--write]\n")
		return exitUsage
	}

	tgt, err := install.Lookup(install.TargetCursorHook)
	if err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	opts := install.InstallOpts{
		Write:       *write,
		ProjectRoot: projectRoot,
		ErrOut:      os.Stderr,
	}
	if err := tgt.Install(opts); err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdInstallGitHook(projectRoot string, args []string) int {
	fs := flag.NewFlagSet("install git-hook", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	write := fs.Bool("write", false, "write post-commit and pre-push hook fragments")
	args = flagsFirst(args, map[string]bool{"write": false})
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace install git-hook [--write]\n")
		return exitUsage
	}

	tgt, err := install.Lookup(install.TargetGitHook)
	if err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	opts := install.InstallOpts{
		Write:       *write,
		ProjectRoot: projectRoot,
		Out:         os.Stdout,
		ErrOut:      os.Stderr,
	}
	if err := tgt.Install(opts); err != nil {
		fmt.Fprintf(os.Stderr, "install: %v\n", err)
		return exitFail
	}
	if *write {
		if err := setHookInstalledFlagForProject(projectRoot, true); err != nil {
			fmt.Fprintf(os.Stderr, "install: %v\n", err)
			return exitFail
		}
	}
	return exitOK
}

func setHookInstalledFlagForProject(projectRoot string, installed bool) error {
	st, err := store.Open(projectRoot)
	if err != nil {
		return err
	}
	defer st.Close()
	return setHookInstalledFlag(st, installed)
}
