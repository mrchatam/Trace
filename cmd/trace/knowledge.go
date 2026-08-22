package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

type knowledgeSynthesizeResponse struct {
	OK      bool `json:"ok"`
	Created int  `json:"created"`
	Updated int  `json:"updated"`
}

func cmdKnowledge(root string, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: trace knowledge list|synthesize|tendencies\n")
		return exitUsage
	}
	switch args[0] {
	case "list":
		return cmdKnowledgeList(root, args[1:])
	case "synthesize":
		return cmdKnowledgeSynthesize(root, args[1:])
	case "tendencies":
		return cmdKnowledgeTendencies(root, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown knowledge subcommand: %s\n", args[0])
		return exitUsage
	}
}

func cmdKnowledgeList(root string, args []string) int {
	fs := flag.NewFlagSet("knowledge list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	topic := fs.String("topic", "", "filter by topic")
	limit := fs.Int("limit", 0, "max rows (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"topic": true, "limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace knowledge list [--topic <t>] [--limit N]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge list: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge list: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "knowledge", "knowledge list"); code != exitOK {
		return code
	}

	rows, err := svc.ListEngineeringKnowledge(context.Background(), domain.ListEngineeringKnowledgeOpts{
		Topic: *topic,
		Limit: *limit,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge list: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(rows); err != nil {
		fmt.Fprintf(os.Stderr, "knowledge list: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdKnowledgeSynthesize(root string, args []string) int {
	fs := flag.NewFlagSet("knowledge synthesize", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace knowledge synthesize\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge synthesize: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge synthesize: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "knowledge", "knowledge synthesize"); code != exitOK {
		return code
	}

	result, err := svc.SynthesizeKnowledge(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge synthesize: %v\n", err)
		return exitFail
	}
	resp := knowledgeSynthesizeResponse{OK: true, Created: result.Created, Updated: result.Updated}
	if err := json.NewEncoder(os.Stdout).Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "knowledge synthesize: %v\n", err)
		return exitFail
	}
	return exitOK
}

func cmdKnowledgeTendencies(root string, args []string) int {
	fs := flag.NewFlagSet("knowledge tendencies", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	limit := fs.Int("limit", 0, "max rows (default 32, cap 64)")
	if err := fs.Parse(flagsFirst(args, map[string]bool{"limit": true})); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "usage: trace knowledge tendencies [--limit N]\n")
		return exitUsage
	}

	abs, err := resolveRoot(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge tendencies: %v\n", err)
		return exitFail
	}
	st, err := store.Open(abs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge tendencies: %v\n", err)
		return exitFail
	}
	defer st.Close()
	svc := domain.New(st)
	if code := failCLIDenied(svc, "knowledge", "knowledge tendencies"); code != exitOK {
		return code
	}

	rows, err := svc.ListTendencies(context.Background(), *limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "knowledge tendencies: %v\n", err)
		return exitFail
	}
	if err := json.NewEncoder(os.Stdout).Encode(rows); err != nil {
		fmt.Fprintf(os.Stderr, "knowledge tendencies: %v\n", err)
		return exitFail
	}
	return exitOK
}
