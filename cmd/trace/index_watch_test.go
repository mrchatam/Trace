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
