package gitcli

import (
	"context"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

// ShowFile implements vcs.Repository — bytes from git show, never SQLite.
func (r *Repo) ShowFile(ctx context.Context, rev, path string) ([]byte, error) {
	path = store.NormalizePath(path)
	if path == "" || strings.Contains(path, "..") {
		return nil, &vcs.Error{Op: "ShowFile", Err: vcs.ErrInvalidPath}
	}
	if rev == "" {
		return nil, &vcs.Error{Op: "ShowFile", Err: fmt.Errorf("rev required")}
	}
	spec := rev + ":" + path
	out, err := r.run.runBytes(ctx, "show", spec)
	if err != nil {
		low := strings.ToLower(err.Error())
		if strings.Contains(low, "does not exist") ||
			strings.Contains(low, "exists on disk, but not in") ||
			strings.Contains(low, "bad revision") ||
			strings.Contains(low, "path not in") ||
			strings.Contains(low, "fatal: path") {
			return nil, &vcs.Error{Op: "ShowFile", Err: vcs.ErrNotFound}
		}
		return nil, &vcs.Error{Op: "ShowFile", Err: err}
	}
	return out, nil
}
