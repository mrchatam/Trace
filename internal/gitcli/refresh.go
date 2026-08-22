package gitcli

import (
	"context"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

// Refresh indexes commits since the durable watermark (last indexed HEAD tip).
// A second call with no new commits is a noop (NewCommits == 0) and does not
// wipe or rewrite the existing index.
func (r *Repo) Refresh(ctx context.Context) (vcs.RefreshResult, error) {
	head, err := r.Head(ctx)
	if err != nil {
		return vcs.RefreshResult{}, err
	}

	watermark, err := r.st.GetMeta(store.MetaVCSWatermark)
	if err != nil {
		return vcs.RefreshResult{}, &vcs.Error{Op: "Refresh", Err: err}
	}

	var revRange []string
	if watermark == "" {
		revRange = []string{"rev-list", "--reverse", head}
	} else if watermark == head {
		return vcs.RefreshResult{NewCommits: 0}, nil
	} else {
		revRange = []string{"rev-list", "--reverse", watermark + ".." + head}
	}

	out, err := r.run.run(ctx, revRange...)
	if err != nil {
		return vcs.RefreshResult{}, &vcs.Error{Op: "Refresh", Err: err}
	}
	oids := splitNonEmptyLines(out)
	if len(oids) == 0 {
		// Tip may match without rev-list output (e.g. same commit); still advance watermark.
		if err := r.st.SetMeta(store.MetaVCSWatermark, head); err != nil {
			return vcs.RefreshResult{}, &vcs.Error{Op: "Refresh", Err: err}
		}
		return vcs.RefreshResult{NewCommits: 0}, nil
	}

	seq, err := r.st.NextCommitSeq()
	if err != nil {
		return vcs.RefreshResult{}, &vcs.Error{Op: "Refresh", Err: err}
	}

	newCount := 0
	for _, oid := range oids {
		meta, err := r.fetchCommitMetaGit(ctx, oid)
		if err != nil {
			return vcs.RefreshResult{}, err
		}
		paths, err := r.diffTreePaths(ctx, oid)
		if err != nil {
			return vcs.RefreshResult{}, err
		}
		indexedPaths := make([]store.IndexedPathChange, len(paths))
		for i, p := range paths {
			indexedPaths[i] = store.IndexedPathChange{Path: p.Path, Status: p.Status}
		}
		if err := r.st.UpsertIndexedCommit(store.IndexedCommit{
			OID:         meta.OID,
			ParentOIDs:  meta.ParentOIDs,
			CommittedAt: meta.CommittedAt,
			Subject:     meta.Subject,
			Seq:         seq,
		}, indexedPaths); err != nil {
			return vcs.RefreshResult{}, &vcs.Error{Op: "Refresh", Err: err}
		}
		seq++
		newCount++
	}

	if err := r.st.SetMeta(store.MetaVCSWatermark, head); err != nil {
		return vcs.RefreshResult{}, &vcs.Error{Op: "Refresh", Err: err}
	}
	return vcs.RefreshResult{NewCommits: newCount}, nil
}

func (r *Repo) fetchCommitMetaGit(ctx context.Context, oid string) (vcs.CommitMeta, error) {
	out, err := r.run.run(ctx, "show", "-s", "--format=%H%x00%P%x00%cI%x00%s", oid)
	if err != nil {
		return vcs.CommitMeta{}, &vcs.Error{Op: "Refresh", Err: err}
	}
	metas, err := parseCommitLog(out)
	if err != nil {
		return vcs.CommitMeta{}, &vcs.Error{Op: "Refresh", Err: err}
	}
	if len(metas) == 0 {
		return vcs.CommitMeta{}, &vcs.Error{Op: "Refresh", Err: fmt.Errorf("empty meta for %s", oid)}
	}
	m := metas[0]
	// Truncate subject to a short summary (no patch bodies anywhere).
	if len(m.Subject) > 200 {
		m.Subject = m.Subject[:200]
	}
	m.Subject = strings.TrimSpace(m.Subject)
	return m, nil
}
