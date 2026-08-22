package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func cmdWhy(root string, args []string) int {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: trace why <type> <id>\n")
		return exitUsage
	}
	entityType, entityID := args[0], args[1]

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "why: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "why: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "why", "why"); code != exitOK {
		return code
	}

	eng := retrieval.New(st)
	if repo, rerr := tryOpenGit(abs, st); rerr == nil {
		defer repo.Close()
		eng = eng.WithVCS(repo)
	}

	res, err := eng.Why(context.Background(), entityType, entityID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "why: %v\n", err)
		return exitFail
	}
	payload, err := encodeWhyWithImpact(context.Background(), domain.New(st), entityType, entityID, res)
	if err != nil {
		fmt.Fprintf(os.Stderr, "why: %v\n", err)
		return exitFail
	}
	fmt.Println(string(payload))
	return exitOK
}

func encodeWhyWithImpact(ctx context.Context, svc *domain.Service, entityType, entityID string, res retrieval.WhyResult) ([]byte, error) {
	impact, err := svc.ImpactSummariesForWhySeed(ctx, retrieval.NormalizeEntityType(entityType), entityID)
	if err != nil {
		return nil, err
	}
	out := struct {
		retrieval.WhyResult
		Impact []domain.DecisionImpact `json:"impact,omitempty"`
	}{WhyResult: res, Impact: impact}
	return json.MarshalIndent(out, "", "  ")
}
