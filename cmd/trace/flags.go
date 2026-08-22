package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// flagsFirst reorders argv so flag.FlagSet can parse flags that appear after
// positionals. valueFlags maps long flag names (no leading dashes) to whether
// they take a separate value argument (false = boolean flag).
func flagsFirst(args []string, valueFlags map[string]bool) []string {
	var flags, pos []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			pos = append(pos, args[i+1:]...)
			break
		}
		if !strings.HasPrefix(a, "-") {
			pos = append(pos, a)
			continue
		}
		flags = append(flags, a)
		name := strings.TrimLeft(a, "-")
		if eq := strings.IndexByte(name, '='); eq >= 0 {
			name = name[:eq]
		} else if takes, ok := valueFlags[name]; ok && takes {
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				flags = append(flags, args[i])
			}
		}
	}
	return append(flags, pos...)
}

// peelServeProjectFlags extracts global-style -C/--project from serve|gui argv so
// `trace serve -C DIR` works as well as `trace -C DIR serve`.
func peelServeProjectFlags(globalRoot string, args []string) (projectRoot string, rest []string, err error) {
	projectRoot = globalRoot
	rest = make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "-C" || a == "--project":
			if i+1 >= len(args) {
				return "", nil, fmt.Errorf("flag %s requires a directory argument", a)
			}
			if err := mergeServeProjectRoot(&projectRoot, globalRoot, args[i+1]); err != nil {
				return "", nil, err
			}
			i++
		case strings.HasPrefix(a, "-C="):
			if err := mergeServeProjectRoot(&projectRoot, globalRoot, strings.TrimPrefix(a, "-C=")); err != nil {
				return "", nil, err
			}
		case strings.HasPrefix(a, "--project="):
			if err := mergeServeProjectRoot(&projectRoot, globalRoot, strings.TrimPrefix(a, "--project=")); err != nil {
				return "", nil, err
			}
		default:
			rest = append(rest, a)
		}
	}
	return projectRoot, rest, nil
}

func mergeServeProjectRoot(projectRoot *string, globalRoot, override string) error {
	override = strings.TrimSpace(override)
	if override == "" {
		return fmt.Errorf("project directory must not be empty")
	}
	absOver, err := filepath.Abs(override)
	if err != nil {
		return err
	}
	if strings.TrimSpace(*projectRoot) == "" {
		*projectRoot = absOver
		return nil
	}
	absCur, err := filepath.Abs(*projectRoot)
	if err != nil {
		return err
	}
	if absCur == absOver {
		return nil
	}
	cwd, _ := os.Getwd()
	absCWD, _ := filepath.Abs(cwd)
	absGlobal, _ := filepath.Abs(globalRoot)
	if absGlobal == absCWD {
		*projectRoot = absOver
		return nil
	}
	return fmt.Errorf("conflicting project roots %q and %q (use one of -C, --project, or --root)", absCur, absOver)
}
