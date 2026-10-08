package store

import (
	"os"
	"path/filepath"
)

// ProjectInitialized reports whether projectRoot has an initialized store, by
// stat only. It never creates .trace/ or trace.db, never acquires the lock and
// never opens the database, so a caller can ask "is this project usable?" on a
// locked or uninitialized root without side effects.
func ProjectInitialized(projectRoot string) bool {
	abs, err := filepath.Abs(projectRoot)
	if err != nil {
		return false
	}
	fi, err := os.Stat(filepath.Join(abs, traceDirName, dbFileName))
	return err == nil && fi.Mode().IsRegular()
}
