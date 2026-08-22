package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/gitcli"
	"github.com/mrchatam/Trace/internal/store"
	"github.com/mrchatam/Trace/internal/vcs"
)

// StateCompareResult is the path delta between two git states plus linked changes.
type StateCompareResult struct {
	From      string   `json:"from"`
	To        string   `json:"to"`
	Added     []string `json:"added"`
	Removed   []string `json:"removed"`
	Modified  []string `json:"modified"`
	ChangeIDs []string `json:"change_ids"`
}

// CompareStates diffs two git OIDs (from..to) and links promoted change ids when present.
// Path lists are repo-relative; no blob content is read or stored.
func (s *Service) CompareStates(ctx context.Context, fromOID, toOID string) (StateCompareResult, error) {
	from, err := normalizeGitCommit(fromOID)
	if err != nil {
		return StateCompareResult{}, err
	}
	to, err := normalizeGitCommit(toOID)
	if err != nil {
		return StateCompareResult{}, err
	}
	if to == "" {
		return StateCompareResult{}, &ErrValidation{Msg: "to OID is required"}
	}

	repo, err := gitcli.OpenWithStore(s.store.ProjectRoot(), s.store)
	if err != nil {
		return StateCompareResult{}, fmt.Errorf("domain: compare states: %w", err)
	}
	defer repo.Close()

	changes, err := repo.DiffNameStatus(ctx, from, to)
	if err != nil {
		return StateCompareResult{}, fmt.Errorf("domain: compare states: %w", err)
	}

	result := StateCompareResult{
		From:     from,
		To:       to,
		Added:    []string{},
		Removed:  []string{},
		Modified: []string{},
	}
	for _, ch := range changes {
		path := store.NormalizePath(ch.Path)
		if path == "" {
			continue
		}
		switch strings.ToUpper(strings.TrimSpace(ch.Status)) {
		case "A":
			result.Added = append(result.Added, path)
		case "D":
			result.Removed = append(result.Removed, path)
		default:
			result.Modified = append(result.Modified, path)
		}
	}

	result.ChangeIDs, err = s.linkedChangeIDs(ctx, repo, from, to)
	if err != nil {
		return StateCompareResult{}, err
	}
	return result, nil
}

func (s *Service) linkedChangeIDs(ctx context.Context, repo vcs.Repository, fromOID, toOID string) ([]string, error) {
	seen := map[string]struct{}{}
	var oids []string
	addOID := func(oid string) {
		oid = strings.TrimSpace(oid)
		if oid == "" {
			return
		}
		if _, ok := seen[oid]; ok {
			return
		}
		seen[oid] = struct{}{}
		oids = append(oids, oid)
	}

	addOID(fromOID)
	addOID(toOID)
	commits, err := repo.CommitsBetween(ctx, fromOID, toOID)
	if err != nil {
		return nil, fmt.Errorf("domain: compare states: %w", err)
	}
	for _, c := range commits {
		addOID(c.OID)
	}

	var ids []string
	idSeen := map[string]struct{}{}
	for _, oid := range oids {
		ch, err := s.store.GetChangeByGitCommit(oid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return nil, err
		}
		if ch.ID == "" {
			continue
		}
		if _, ok := idSeen[ch.ID]; ok {
			continue
		}
		idSeen[ch.ID] = struct{}{}
		ids = append(ids, ch.ID)
	}
	if ids == nil {
		ids = []string{}
	}
	return ids, nil
}
