package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/mrchatam/Trace/internal/store"
)

func TestIndexWatchDebounced(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() {
		done <- runIndexWatch(ctx, dir, []string{"--debounce", "50ms", "."})
	}()

	// Allow watcher setup.
	time.Sleep(100 * time.Millisecond)

	goPath := filepath.Join(dir, "watchme.go")
	body := []byte("package watchme\nfunc Watched() {}\n")
	if err := os.WriteFile(goPath, body, 0o644); err != nil {
		t.Fatal(err)
	}

	// Wait for debounced index (watch holds store lock until exit).
	time.Sleep(200 * time.Millisecond)
	cancel()
	if code := <-done; code != exitOK {
		t.Fatalf("watch exit: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	if _, ferr := st.GetFileByPath("watchme.go"); ferr != nil {
		t.Fatalf("watchme.go not indexed: %v", ferr)
	}
	syms, serr := st.ListSymbolsByPath("watchme.go")
	if serr != nil || len(syms) == 0 {
		t.Fatalf("symbols: err=%v len=%d", serr, len(syms))
	}
}

// TestIndexWatchHealsDrift pins watch-mode drift healing: files the event
// stream cannot see (a directory created with contents — only the dir event
// fires — and fsnotify blind spots generally) are indexed by the drift sweep
// after the debounce window, so the store converges to the tree without a
// manual `trace index`.
func TestIndexWatchHealsDrift(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() {
		done <- runIndexWatch(ctx, dir, []string{"--debounce", "50ms", "."})
	}()

	// Allow watcher setup.
	time.Sleep(100 * time.Millisecond)

	// A new directory with a file already inside: the watcher receives only
	// the dir Create event, so per-path timers never fire for inner.go.
	sub := filepath.Join(dir, "late")
	if err := os.Mkdir(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "inner.go"), []byte("package late\nfunc Inner() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Root-level new file: covered by the event path, healed identically if
	// the event is dropped.
	if err := os.WriteFile(filepath.Join(dir, "late_root.go"), []byte("package late\nfunc Root() {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Wait past the debounce window plus the sweep (watch holds the store
	// lock until exit, so assertions run after cancel).
	time.Sleep(600 * time.Millisecond)
	cancel()
	if code := <-done; code != exitOK {
		t.Fatalf("watch exit: %d", code)
	}

	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	for _, rel := range []string{"late/inner.go", "late_root.go"} {
		if _, ferr := st.GetFileByPath(rel); ferr != nil {
			t.Fatalf("drift heal missed %s: %v", rel, ferr)
		}
		syms, serr := st.ListSymbolsByPath(rel)
		if serr != nil || len(syms) == 0 {
			t.Fatalf("%s symbols: err=%v len=%d", rel, serr, len(syms))
		}
	}
}

func TestIndexWatchForegroundExit(t *testing.T) {
	dir := t.TempDir()
	if code := run([]string{"-C", dir, "init"}); code != exitOK {
		t.Fatalf("init: %d", code)
	}

	prev := indexWatchNotifyContext
	indexWatchNotifyContext = func() (context.Context, context.CancelFunc) {
		return context.WithCancel(context.Background())
	}
	t.Cleanup(func() { indexWatchNotifyContext = prev })

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() {
		done <- runIndexWatch(ctx, dir, []string{"--debounce", "50ms"})
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case code := <-done:
		if code != exitOK {
			t.Fatalf("watch exit: %d want %d", code, exitOK)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watch did not exit after cancel")
	}
}
