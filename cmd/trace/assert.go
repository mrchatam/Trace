package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mrchatam/Trace/internal/domain"
)

// assertCLICommand fails closed via domain Assert with slug cli:<command>.
// Alias fold (cli:reindex → cli:index) and builtin AUTO_ALLOW live in domain.
func assertCLICommand(ctx context.Context, svc *domain.Service, command string) error {
	return svc.AssertToolAllowed(ctx, "cli:"+command)
}

func failCLIDenied(svc *domain.Service, command, prefix string) int {
	if err := assertCLICommand(context.Background(), svc, command); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", prefix, err)
		return exitFail
	}
	return exitOK
}
