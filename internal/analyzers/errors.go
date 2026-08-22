package analyzers

import "fmt"

// SkipError is returned when a path cannot be indexed (unsupported or binary).
// Callers that walk a tree may skip; IndexFile surfaces it as a clear error.
type SkipError struct {
	Path   string
	Reason string
}

func (e *SkipError) Error() string {
	return fmt.Sprintf("analyzers: skip %q: %s", e.Path, e.Reason)
}

func unsupportedExt(path string) error {
	return &SkipError{
		Path: path,
		Reason: "unsupported extension (Tier-2 deferred or Tier-3 path-only); " +
			"see docs/INDEX_LANG_POLICY.md",
	}
}

func binaryFile(path string) error {
	return &SkipError{Path: path, Reason: "binary content (NUL in first 8KiB)"}
}
