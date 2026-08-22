package mcp

import (
	"context"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// assertMCPToolAllowed fails closed via domain Assert with slug mcp:<toolName>.
func assertMCPToolAllowed(ctx context.Context, st *store.Store, toolName string) error {
	return domain.New(st).AssertToolAllowed(ctx, "mcp:"+toolName)
}
