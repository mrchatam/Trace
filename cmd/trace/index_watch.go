package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mrchatam/Trace/internal/analyzers"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

// indexWatchNotifyContext is injectable for short-lived watch tests (G7-F5).
var indexWatchNotifyContext = func() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}

func cmdIndexWatch(root string, args []string) int {
	ctx, cancel := indexWatchNotifyContext()
	defer cancel()
	return runIndexWatch(ctx, root, args)
}

func runIndexWatch(ctx context.Context, root string, args []string) int {
	fs := flag.NewFlagSet("index watch", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	debounce := fs.Duration("debounce", 300*time.Millisecond, "debounce window per changed path")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	watchPaths := fs.Args()

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "index watch: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "index watch: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "index", "index watch"); code != exitOK {
		return code
	}

	dirs, err := collectWatchDirs(abs, watchPaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "index watch: %v\n", err)
		return exitFail
	}
	if len(dirs) == 0 {
		fmt.Fprintf(os.Stderr, "index watch: no directories to watch\n")
		return exitFail
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "index watch: %v\n", err)
		return exitFail
	}
	defer watcher.Close()

	for _, d := range dirs {
		if err := watcher.Add(d); err != nil {
			fmt.Fprintf(os.Stderr, "index watch: add %s: %v\n", d, err)
			return exitFail
		}
	}

	var repo vcs.Repository
	if r, rerr := tryOpenGit(abs, st); rerr == nil {
		repo = r
		defer repo.Close()
	}

	fmt.Fprintf(os.Stderr, "watching %d director%s under %s (debounce %s)\n",
		len(dirs), pluralSuffix(len(dirs), "y", "ies"), abs, debounce)

	indexCtx := context.Background()
	var pending sync.Map // path -> *time.Timer

	scheduleIndex := func(absPath string) {
		rel, normAbs, err := normalizeProjectPath(abs, absPath)
		if err != nil {
			return
		}
		if _, err := os.Stat(normAbs); err != nil {
			return
		}
		if isT0SkipPath(rel) {
			return
		}
		if _, ok := analyzers.DetectLanguage(rel); !ok {
			return
		}

		if prev, loaded := pending.Load(absPath); loaded {
			prev.(*time.Timer).Stop()
		}
		t := time.AfterFunc(*debounce, func() {
			pending.Delete(absPath)
			if err := indexOne(indexCtx, st, repo, abs, rel, normAbs); err != nil {
				var skip *analyzers.SkipError
				if errors.As(err, &skip) {
					return
				}
				fmt.Fprintf(os.Stderr, "index watch: %s: %v\n", rel, err)
				return
			}
			fmt.Fprintf(os.Stderr, "indexed %s\n", rel)
		})
		pending.Store(absPath, t)
	}

	for {
		select {
		case <-ctx.Done():
			pending.Range(func(key, value any) bool {
				value.(*time.Timer).Stop()
				return true
			})
			return exitOK
		case ev, ok := <-watcher.Events:
			if !ok {
				return exitOK
			}
			if ev.Has(fsnotify.Chmod) && !ev.Has(fsnotify.Write|fsnotify.Create|fsnotify.Rename) {
				continue
			}
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
				continue
			}
			path := ev.Name
			info, err := os.Stat(path)
			if err != nil {
				continue
			}
			if info.IsDir() {
				if !isT0SkipDir(info.Name()) {
					_ = watcher.Add(path)
				}
				continue
			}
			scheduleIndex(path)
		case err, ok := <-watcher.Errors:
			if !ok {
				return exitOK
			}
			fmt.Fprintf(os.Stderr, "index watch: %v\n", err)
		}
	}
}

func pluralSuffix(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

func collectWatchDirs(root string, paths []string) ([]string, error) {
	if len(paths) == 0 {
		paths = []string{"."}
	}
	seen := make(map[string]struct{})
	var dirs []string
	addDir := func(d string) {
		d = filepath.Clean(d)
		if _, ok := seen[d]; ok {
			return
		}
		seen[d] = struct{}{}
		dirs = append(dirs, d)
	}

	for _, p := range paths {
		_, absPath, err := normalizeProjectPath(root, p)
		if err != nil {
			return nil, err
		}
		info, err := os.Stat(absPath)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			addDir(filepath.Dir(absPath))
			continue
		}
		err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !d.IsDir() {
				return nil
			}
			if isT0SkipDir(d.Name()) {
				return filepath.SkipDir
			}
			addDir(path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return dirs, nil
}
