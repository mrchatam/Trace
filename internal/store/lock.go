package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

// ErrLocked is returned when another process or open Store already holds
// <projectRoot>/.trace/trace.lock after the short Open wait budget.
// Single-writer is intentional: serialize CLI↔MCP (or parallel Trace on one
// root), or use separate -C / worktree roots.
var ErrLocked = errors.New("store: project root locked (.trace/trace.lock held by another process); serialize CLI↔MCP (or parallel Trace) on one root, or use separate -C / worktree roots")

// ErrNotInitialized is returned by OpenExisting when <projectRoot>/.trace/trace.db
// is missing or is not a regular file. OpenExisting must not create .trace/ or
// the database.
var ErrNotInitialized = errors.New("store: not initialized (missing .trace/trace.db)")

const (
	lockFileName = "trace.lock"
	// defaultLockWait is the Open acquire budget for brief soft races (DF-47).
	defaultLockWait = 350 * time.Millisecond
	lockRetryStep   = 25 * time.Millisecond
	lockWaitEnv     = "TRACE_LOCK_WAIT_MS"
)

type flockHandle struct {
	f *os.File
}

func lockWaitBudget() time.Duration {
	raw := os.Getenv(lockWaitEnv)
	if raw == "" {
		return defaultLockWait
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms < 0 {
		return defaultLockWait
	}
	return time.Duration(ms) * time.Millisecond
}

func acquireTraceLock(traceDir string) (*flockHandle, error) {
	path := filepath.Join(traceDir, lockFileName)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, fmt.Errorf("store: open lock: %w", err)
	}
	deadline := time.Now().Add(lockWaitBudget())
	for {
		err := unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB)
		if err == nil {
			return &flockHandle{f: f}, nil
		}
		if !errors.Is(err, unix.EWOULDBLOCK) && !errors.Is(err, unix.EAGAIN) {
			_ = f.Close()
			return nil, fmt.Errorf("store: flock: %w", err)
		}
		if !time.Now().Before(deadline) {
			_ = f.Close()
			return nil, ErrLocked
		}
		sleep := lockRetryStep
		if rem := time.Until(deadline); rem < sleep {
			sleep = rem
		}
		if sleep > 0 {
			time.Sleep(sleep)
		}
	}
}

func (h *flockHandle) release() error {
	if h == nil || h.f == nil {
		return nil
	}
	_ = unix.Flock(int(h.f.Fd()), unix.LOCK_UN)
	err := h.f.Close()
	h.f = nil
	return err
}
