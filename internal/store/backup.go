package store

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ErrRestoreExists is returned when Restore would overwrite an existing
// trace.db without force=true.
var ErrRestoreExists = errors.New("store: restore target .trace/trace.db already exists (pass force to overwrite)")

// BackupOptions controls optional companions written next to the DB snapshot.
type BackupOptions struct {
	// IncludeToken also copies .trace/access.token beside the snapshot
	// (destPath + ".token"). Default false — token is excluded.
	IncludeToken bool
}

// BackupTo writes a consistent snapshot of this store's trace.db to destPath
// using VACUUM INTO while the exclusive lock is held.
func (s *Store) BackupTo(destPath string) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("store: backup: nil store")
	}
	absDest, err := filepath.Abs(destPath)
	if err != nil {
		return fmt.Errorf("store: backup dest: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(absDest), 0o755); err != nil {
		return fmt.Errorf("store: mkdir backup dest: %w", err)
	}
	// Remove stale dest so VACUUM INTO can create a fresh file.
	_ = os.Remove(absDest)
	if _, err := s.db.Exec(`VACUUM INTO ?`, absDest); err != nil {
		return fmt.Errorf("store: vacuum into backup: %w", err)
	}
	return nil
}

// Backup opens the project store (lock + auth), snapshots trace.db to destPath,
// then closes. Optional IncludeToken copies access.token to destPath+".token".
func Backup(projectRoot, destPath string, opts BackupOptions) error {
	s, err := Open(projectRoot)
	if err != nil {
		return err
	}
	defer s.Close()
	if err := s.BackupTo(destPath); err != nil {
		return err
	}
	if opts.IncludeToken {
		src := filepath.Join(s.projectRoot, traceDirName, accessTokenFileName)
		tok, err := readAccessTokenFile(src)
		if err != nil {
			return err
		}
		if tok != "" {
			side := destPath + ".token"
			if err := os.WriteFile(side, []byte(tok+"\n"), 0o600); err != nil {
				return fmt.Errorf("store: backup token: %w", err)
			}
		}
	}
	return nil
}

// Restore installs a trace.db snapshot into <projectRoot>/.trace/trace.db.
// Fails with ErrRestoreExists if the target exists and force is false.
// Fails with ErrLocked if another Open holds the lock.
// After install, opens the store (migrate) and rebinds projects.root_path to
// the current Abs root. Does not copy access.token.
func Restore(projectRoot, backupPath string, force bool) error {
	absRoot, traceDir, err := resolveTraceDir(projectRoot)
	if err != nil {
		return err
	}
	absBak, err := filepath.Abs(backupPath)
	if err != nil {
		return fmt.Errorf("store: restore source: %w", err)
	}
	if _, err := os.Stat(absBak); err != nil {
		return fmt.Errorf("store: restore source: %w", err)
	}

	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		return fmt.Errorf("store: mkdir .trace: %w", err)
	}

	lock, err := acquireTraceLock(traceDir)
	if err != nil {
		return err
	}
	defer func() { _ = lock.release() }()

	dbPath := filepath.Join(traceDir, dbFileName)
	if _, err := os.Stat(dbPath); err == nil {
		if !force {
			return ErrRestoreExists
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("store: stat target db: %w", err)
	}

	if err := copyFile(absBak, dbPath); err != nil {
		return fmt.Errorf("store: install backup: %w", err)
	}
	_ = lock.release()
	lock = nil

	s, err := openStore(absRoot, true)
	if err != nil {
		return err
	}
	return s.Close()
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp := dst + ".tmp"
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, dst); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}
