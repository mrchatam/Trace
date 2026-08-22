package compiler

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"

	"github.com/mrchatam/Trace/internal/store"
)

const indexHonestyNotice = "indexed view may lag disk — reindex or Read live"

// buildIndexHonesty compares file items' indexed content_hash to sha256(disk).
// Prefer false-fresh: missing row, missing file, or I/O/hash errors are omitted.
// universe should be the pre-trim pipeline (post Layer-1 / MaxCandidateHits).
func buildIndexHonesty(st *store.Store, universe []Item) *IndexHonesty {
	if st == nil {
		return nil
	}
	root := st.ProjectRoot()
	seen := make(map[string]struct{})
	var stale []string
	for _, it := range universe {
		if it.EntityType != "file" || it.EntityID == "" {
			continue
		}
		rec, err := st.GetFileByID(it.EntityID)
		if err != nil || rec.Path == "" || rec.ContentHash == "" {
			continue // false-fresh
		}
		if _, ok := seen[rec.Path]; ok {
			continue
		}
		diskPath := filepath.Join(root, filepath.FromSlash(rec.Path))
		sum, err := sha256FileHex(diskPath)
		if err != nil {
			continue // false-fresh
		}
		if sum == rec.ContentHash {
			continue
		}
		seen[rec.Path] = struct{}{}
		stale = append(stale, rec.Path)
	}
	if len(stale) == 0 {
		return nil
	}
	// Unique → sort → first 8 (deterministic lex banner); expose total when truncated.
	sort.Strings(stale)
	staleTotal := len(stale)
	truncated := staleTotal > maxIndexHonestyStalePaths
	if truncated {
		stale = stale[:maxIndexHonestyStalePaths]
	}
	return &IndexHonesty{
		StalePaths:     stale,
		StaleTotal:     staleTotal,
		StaleTruncated: truncated,
		Notice:         indexHonestyNotice,
	}
}

func sha256FileHex(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), nil
}
