package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

// cmdImpact is a thin G19 adapter: finding / alternative / report / walk / predict / compare.
func cmdImpact(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace impact finding|alternative|report|walk|predict|compare …\n")
		return exitUsage
	}
	switch args[0] {
	case "finding":
		return cmdImpactFinding(root, args[1:])
	case "alternative":
		return cmdImpactAlternative(root, args[1:])
	case "report":
		return cmdImpactReport(root, args[1:])
	case "walk":
		return cmdImpactWalk(root, args[1:])
	case "predict":
		return cmdImpactPredict(root, args[1:])
	case "compare":
		return cmdImpactCompare(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "impact: unknown subcommand %q (want finding|alternative|report|walk|predict|compare)\n", args[0])
		return exitUsage
	}
}

func cmdImpactFinding(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace impact finding add|list …\n")
		return exitUsage
	}
	switch args[0] {
	case "add":
		return cmdImpactFindingAdd(root, args[1:])
	case "list":
		return cmdImpactFindingList(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "impact finding: unknown subcommand %q (want add|list)\n", args[0])
		return exitUsage
	}
}

func cmdImpactFindingAdd(root string, args []string) int {
	fs := flag.NewFlagSet("impact finding add", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	decisionID := fs.String("decision", "", "decision UUID")
	class := fs.String("class", "", "SAFE|CAUTION|HIGH|DESTRUCTIVE|REVERSAL")
	kind := fs.String("kind", "", "AFFECTED_WORK|INVALIDATED_ASSUMPTION|WORK_AT_RISK|NEW_WORK|UNRESOLVED")
	uncertainty := fs.String("uncertainty", "", "KNOWN|LIKELY|POSSIBLE|UNKNOWN (default UNKNOWN)")
	body := fs.String("body", "", "finding body")
	relatedType := fs.String("related-type", "", "optional related entity type")
	relatedID := fs.String("related-id", "", "optional related entity id")
	id := fs.String("id", "", "optional UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *decisionID == "" || strings.TrimSpace(*class) == "" || strings.TrimSpace(*kind) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace impact finding add --decision <id> --class SAFE|CAUTION|HIGH|DESTRUCTIVE|REVERSAL --kind AFFECTED_WORK|… [--uncertainty …] [--body …] [--related-type …] [--related-id …] [--id …]\n")
		return exitUsage
	}

	svc, st, code := openImpact(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	f, err := svc.AddImpactFinding(context.Background(), *decisionID, domain.ImpactFindingInput{
		ID: *id, ImpactClass: *class, Uncertainty: *uncertainty, Kind: *kind,
		Body: *body, RelatedType: *relatedType, RelatedID: *relatedID,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "id": f.ID, "decision": f.DecisionID,
		"impact_class": f.ImpactClass, "uncertainty": f.Uncertainty, "kind": f.Kind,
	})
	return exitOK
}

func cmdImpactFindingList(root string, args []string) int {
	fs := flag.NewFlagSet("impact finding list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	decisionID := fs.String("decision", "", "decision UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *decisionID == "" {
		fmt.Fprintf(os.Stderr, "usage: trace impact finding list --decision <id>\n")
		return exitUsage
	}

	svc, st, code := openImpact(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	list, err := svc.ListImpactFindings(context.Background(), *decisionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	if list == nil {
		list = []store.DecisionImpactFinding{}
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "findings": list, "count": len(list),
	})
	return exitOK
}

func cmdImpactAlternative(root string, args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: trace impact alternative add|list|recommend …\n")
		return exitUsage
	}
	switch args[0] {
	case "add":
		return cmdImpactAlternativeAdd(root, args[1:])
	case "list":
		return cmdImpactAlternativeList(root, args[1:])
	case "recommend":
		return cmdImpactAlternativeRecommend(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "impact alternative: unknown subcommand %q (want add|list|recommend)\n", args[0])
		return exitUsage
	}
}

func cmdImpactAlternativeAdd(root string, args []string) int {
	fs := flag.NewFlagSet("impact alternative add", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	decisionID := fs.String("decision", "", "decision UUID")
	title := fs.String("title", "", "alternative title")
	body := fs.String("body", "", "alternative body")
	recommended := fs.Bool("recommended", false, "mark as recommended (clears siblings)")
	id := fs.String("id", "", "optional UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *decisionID == "" || strings.TrimSpace(*title) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace impact alternative add --decision <id> --title <t> [--body …] [--recommended] [--id …]\n")
		return exitUsage
	}

	svc, st, code := openImpact(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	a, err := svc.AddDecisionAlternative(context.Background(), *decisionID, domain.AlternativeInput{
		ID: *id, Title: *title, Body: *body, Recommended: *recommended,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "id": a.ID, "decision": a.DecisionID,
		"title": a.Title, "is_recommended": a.IsRecommended,
	})
	return exitOK
}

func cmdImpactAlternativeList(root string, args []string) int {
	fs := flag.NewFlagSet("impact alternative list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	decisionID := fs.String("decision", "", "decision UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *decisionID == "" {
		fmt.Fprintf(os.Stderr, "usage: trace impact alternative list --decision <id>\n")
		return exitUsage
	}

	svc, st, code := openImpact(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	list, err := svc.ListDecisionAlternatives(context.Background(), *decisionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	if list == nil {
		list = []store.DecisionAlternative{}
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "alternatives": list, "count": len(list),
	})
	return exitOK
}

func cmdImpactAlternativeRecommend(root string, args []string) int {
	fs := flag.NewFlagSet("impact alternative recommend", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	decisionID := fs.String("decision", "", "decision UUID")
	id := fs.String("id", "", "alternative UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *decisionID == "" || *id == "" {
		fmt.Fprintf(os.Stderr, "usage: trace impact alternative recommend --decision <id> --id <alternative_id>\n")
		return exitUsage
	}

	svc, st, code := openImpact(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	if err := svc.SetRecommendedAlternative(context.Background(), *decisionID, *id); err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok": true, "decision": *decisionID, "id": *id, "is_recommended": true,
	})
	return exitOK
}

func cmdImpactReport(root string, args []string) int {
	fs := flag.NewFlagSet("impact report", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	decisionID := fs.String("decision", "", "decision UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *decisionID == "" {
		fmt.Fprintf(os.Stderr, "usage: trace impact report --decision <id>\n")
		return exitUsage
	}

	svc, st, code := openImpact(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	rep, err := svc.ImpactReport(context.Background(), *decisionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok":                  true,
		"decision_id":         rep.Decision.ID,
		"affected_task_ids":   rep.AffectedTaskIDs,
		"findings":            rep.Findings,
		"alternatives":        rep.Alternatives,
		"overall_class":       rep.OverallClass,
		"overall_uncertainty": rep.OverallUncertainty,
		"has_unknown":         rep.HasUnknown,
		"incomplete":          rep.Incomplete,
	})
	return exitOK
}

// seedFlagList collects repeated --seed file:<uuid>|symbol:<uuid> values.
type seedFlagList []string

func (s *seedFlagList) String() string { return strings.Join(*s, ",") }
func (s *seedFlagList) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func cmdImpactWalk(root string, args []string) int {
	fs := flag.NewFlagSet("impact walk", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var seeds seedFlagList
	fs.Var(&seeds, "seed", "seed as file:<uuid> or symbol:<uuid> (repeatable)")
	depth := fs.Int("depth", retrieval.DefaultImpactDepth(), "BFS depth 1..2 (default 2)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if len(seeds) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace impact walk --seed file:<uuid>|symbol:<uuid> [--seed …] [--depth 1|2]\n")
		return exitUsage
	}

	walkSeeds := make([]retrieval.ImpactSeed, 0, len(seeds))
	for _, raw := range seeds {
		typ, id, ok := strings.Cut(raw, ":")
		if !ok || typ == "" || id == "" {
			fmt.Fprintf(os.Stderr, "impact walk: bad --seed %q (want file:<uuid> or symbol:<uuid>)\n", raw)
			return exitUsage
		}
		typ = strings.ToLower(strings.TrimSpace(typ))
		if typ != "file" && typ != "symbol" {
			fmt.Fprintf(os.Stderr, "impact walk: seed type must be file|symbol, got %q\n", typ)
			return exitUsage
		}
		walkSeeds = append(walkSeeds, retrieval.ImpactSeed{EntityType: typ, EntityID: strings.TrimSpace(id)})
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	defer st.Close()
	if code := failCLIDenied(domain.New(st), "impact", "impact"); code != exitOK {
		return code
	}

	eng := retrieval.New(st)
	res, err := eng.ImpactWalk(context.Background(), walkSeeds, *depth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok":             true,
		"seeds":          res.Seeds,
		"blast":          res.Blast,
		"affected_tests": res.AffectedTests,
		"blast_total":    res.BlastTotal,
		"blast_kept":     res.BlastKept,
		"truncated":      res.Truncated,
		"depth":          res.Depth,
	})
	return exitOK
}

func cmdImpactPredict(root string, args []string) int {
	fs := flag.NewFlagSet("impact predict", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	changeID := fs.String("change", "", "change UUID")
	depth := fs.Int("depth", retrieval.DefaultImpactDepth(), "BFS depth 1..2 (default 2)")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if strings.TrimSpace(*changeID) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace impact predict --change <id> [--depth 1|2]\n")
		return exitUsage
	}

	svc, st, code := openImpact(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	row, err := svc.PredictImpactForChange(context.Background(), *changeID, *depth)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	var payload domain.PredictedImpactPayload
	_ = json.Unmarshal([]byte(row.PredictedJSON), &payload)
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok":                 true,
		"change_id":          row.ChangeID,
		"depth":              row.Depth,
		"seeds":              payload.Seeds,
		"blast_keys":         payload.BlastKeys,
		"affected_test_keys": payload.AffectedTestKeys,
		"blast_total":        payload.BlastTotal,
		"blast_kept":         payload.BlastKept,
		"truncated":          payload.Truncated,
		"created_at":         row.CreatedAt,
	})
	return exitOK
}

func cmdImpactCompare(root string, args []string) int {
	fs := flag.NewFlagSet("impact compare", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	changeID := fs.String("change", "", "change UUID")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if strings.TrimSpace(*changeID) == "" {
		fmt.Fprintf(os.Stderr, "usage: trace impact compare --change <id>\n")
		return exitUsage
	}

	svc, st, code := openImpact(root)
	if code != exitOK {
		return code
	}
	defer st.Close()

	res, err := svc.CompareActualImpact(context.Background(), *changeID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return exitFail
	}
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"ok":          true,
		"change_id":   res.ChangeID,
		"matched":     res.Delta.Matched,
		"unexpected":  res.Delta.Unexpected,
		"missed":      res.Delta.Missed,
		"compared_at": res.ComparedAt,
	})
	return exitOK
}

func openImpact(root string) (*domain.Service, *store.Store, int) {
	svc, st, code := openDomain(root)
	if code != exitOK {
		return nil, nil, code
	}
	if code := failCLIDenied(svc, "impact", "impact"); code != exitOK {
		st.Close()
		return nil, nil, code
	}
	retrieval.WireDomainImpactWalker(svc, retrieval.New(st))
	return svc, st, exitOK
}

func openDomain(root string) (*domain.Service, *store.Store, int) {
	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return nil, nil, exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "impact: %v\n", err)
		return nil, nil, exitFail
	}
	return domain.New(st), st, exitOK
}
