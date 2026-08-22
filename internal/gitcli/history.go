package gitcli

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

const defaultHistoryLimit = 50

// History implements vcs.Repository. Prefers the thin index after Refresh;
// falls back to git log.
func (r *Repo) History(ctx context.Context, path string, limit int) ([]vcs.CommitMeta, error) {
	path = store.NormalizePath(path)
	if path == "" {
		return nil, &vcs.Error{Op: "History", Err: vcs.ErrInvalidPath}
	}
	if limit <= 0 {
		limit = defaultHistoryLimit
	}

	// Prefer the thin index only when Refresh has caught up to HEAD.
	// Otherwise a non-empty but stale index would hide newer commits
	// (and mid-Refresh partial indexes would truncate History).
	if wm, err := r.st.GetMeta(store.MetaVCSWatermark); err == nil && wm != "" {
		if head, err := r.Head(ctx); err == nil && wm == head {
			if indexed, err := r.st.ListIndexedHistory(path, limit); err == nil && len(indexed) > 0 {
				return toCommitMetas(indexed), nil
			}
		}
	}

	args := []string{
		"log",
		"--format=%H%x00%P%x00%cI%x00%s",
		"-n", strconv.Itoa(limit),
		"--", path,
	}
	out, err := r.run.run(ctx, args...)
	if err != nil {
		return nil, &vcs.Error{Op: "History", Err: err}
	}
	return parseCommitLog(out)
}

// LastChanged implements vcs.Repository.
func (r *Repo) LastChanged(ctx context.Context, path string) (vcs.CommitMeta, error) {
	hist, err := r.History(ctx, path, 1)
	if err != nil {
		return vcs.CommitMeta{}, err
	}
	if len(hist) == 0 {
		return vcs.CommitMeta{}, &vcs.Error{Op: "LastChanged", Err: vcs.ErrNotFound}
	}
	return hist[0], nil
}

// CommitsBetween implements vcs.Repository — oldest first, fromRev..toRev.
func (r *Repo) CommitsBetween(ctx context.Context, fromRev, toRev string) ([]vcs.CommitMeta, error) {
	if toRev == "" {
		return nil, &vcs.Error{Op: "CommitsBetween", Err: fmt.Errorf("toRev required")}
	}
	rangeSpec := toRev
	if fromRev != "" {
		rangeSpec = fromRev + ".." + toRev
	}

	out, err := r.run.run(ctx, "rev-list", "--reverse", rangeSpec)
	if err != nil {
		return nil, &vcs.Error{Op: "CommitsBetween", Err: err}
	}
	oids := splitNonEmptyLines(out)
	if len(oids) == 0 {
		return nil, nil
	}

	// Prefer index when all OIDs are present; otherwise fill from git.
	indexed, _ := r.st.ListIndexedCommitsByOIDs(oids)
	if len(indexed) == len(oids) {
		return toCommitMetas(indexed), nil
	}

	metas := make([]vcs.CommitMeta, 0, len(oids))
	for _, oid := range oids {
		m, err := r.commitMeta(ctx, oid)
		if err != nil {
			return nil, err
		}
		metas = append(metas, m)
	}
	return metas, nil
}

// Changes implements vcs.Repository.
func (r *Repo) Changes(ctx context.Context, commitOID string) ([]vcs.PathChange, error) {
	if commitOID == "" {
		return nil, &vcs.Error{Op: "Changes", Err: fmt.Errorf("commit OID required")}
	}

	if paths, err := r.st.ListIndexedCommitPaths(commitOID); err == nil && len(paths) > 0 {
		out := make([]vcs.PathChange, len(paths))
		for i, p := range paths {
			out[i] = vcs.PathChange{Path: p.Path, Status: p.Status}
		}
		return out, nil
	}

	return r.diffTreePaths(ctx, commitOID)
}

func (r *Repo) diffTreePaths(ctx context.Context, commitOID string) ([]vcs.PathChange, error) {
	out, err := r.run.run(ctx, "diff-tree", "--no-commit-id", "--name-status", "-r", "--root", commitOID)
	if err != nil {
		return nil, &vcs.Error{Op: "Changes", Err: err}
	}
	return parseNameStatus(out), nil
}

// DiffNameStatus implements vcs.Repository — git diff --name-status from..to.
func (r *Repo) DiffNameStatus(ctx context.Context, fromRev, toRev string) ([]vcs.PathChange, error) {
	toRev = strings.TrimSpace(toRev)
	if toRev == "" {
		return nil, &vcs.Error{Op: "DiffNameStatus", Err: fmt.Errorf("toRev required")}
	}
	if err := r.verifyCommit(ctx, fromRev); err != nil {
		return nil, err
	}
	if err := r.verifyCommit(ctx, toRev); err != nil {
		return nil, err
	}

	rangeSpec := toRev
	if strings.TrimSpace(fromRev) != "" {
		rangeSpec = strings.TrimSpace(fromRev) + ".." + toRev
	}
	out, err := r.run.run(ctx, "diff", "--name-status", rangeSpec)
	if err != nil {
		return nil, &vcs.Error{Op: "DiffNameStatus", Err: err}
	}
	return parseNameStatus(out), nil
}

func (r *Repo) verifyCommit(ctx context.Context, oid string) error {
	oid = strings.TrimSpace(oid)
	if oid == "" {
		return nil
	}
	if _, err := r.run.run(ctx, "rev-parse", "--verify", oid+"^{commit}"); err != nil {
		return &vcs.Error{Op: "verifyCommit", Err: err}
	}
	return nil
}

func (r *Repo) commitMeta(ctx context.Context, oid string) (vcs.CommitMeta, error) {
	if c, err := r.st.GetIndexedCommit(oid); err == nil {
		return toCommitMeta(c), nil
	}
	out, err := r.run.run(ctx, "show", "-s", "--format=%H%x00%P%x00%cI%x00%s", oid)
	if err != nil {
		return vcs.CommitMeta{}, &vcs.Error{Op: "commitMeta", Err: err}
	}
	metas, err := parseCommitLog(out)
	if err != nil || len(metas) == 0 {
		if err == nil {
			err = vcs.ErrNotFound
		}
		return vcs.CommitMeta{}, &vcs.Error{Op: "commitMeta", Err: err}
	}
	return metas[0], nil
}

func parseCommitLog(out string) ([]vcs.CommitMeta, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return nil, nil
	}
	var metas []vcs.CommitMeta
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\x00")
		if len(parts) < 4 {
			return nil, fmt.Errorf("gitcli: bad commit log line %q", line)
		}
		parents := strings.Fields(strings.TrimSpace(parts[1]))
		metas = append(metas, vcs.CommitMeta{
			OID:         parts[0],
			ParentOIDs:  parents,
			CommittedAt: normalizeTime(parts[2]),
			Subject:     parts[3],
		})
	}
	return metas, nil
}

func parseNameStatus(out string) []vcs.PathChange {
	var changes []vcs.PathChange
	for _, line := range splitNonEmptyLines(out) {
		fields := strings.Split(line, "\t")
		if len(fields) < 2 {
			continue
		}
		status := fields[0]
		path := fields[len(fields)-1]
		// Renames: R100\told\tnew — keep status letter + new path.
		if len(status) > 0 {
			status = string(status[0])
		}
		changes = append(changes, vcs.PathChange{
			Path:   store.NormalizePath(path),
			Status: status,
		})
	}
	return changes
}

func splitNonEmptyLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}

func normalizeTime(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	// git %cI is already ISO-8601; accept as-is if parseable.
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC().Format(time.RFC3339)
	}
	return s
}

func toCommitMeta(c store.IndexedCommit) vcs.CommitMeta {
	return vcs.CommitMeta{
		OID:         c.OID,
		ParentOIDs:  append([]string(nil), c.ParentOIDs...),
		CommittedAt: c.CommittedAt,
		Subject:     c.Subject,
	}
}

func toCommitMetas(cs []store.IndexedCommit) []vcs.CommitMeta {
	out := make([]vcs.CommitMeta, len(cs))
	for i, c := range cs {
		out[i] = toCommitMeta(c)
	}
	return out
}
