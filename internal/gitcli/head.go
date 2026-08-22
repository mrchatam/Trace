package gitcli

import (
	"context"
	"path/filepath"
	"strings"
)

// HeadAtRoot returns git rev-parse HEAD for projectRoot when it is a git work tree.
// Returns ("", nil) when root is not a repository, git is missing, or HEAD is unavailable.
func HeadAtRoot(ctx context.Context, projectRoot string) (string, error) {
	abs, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", err
	}
	r := newRunner(abs)
	out, err := r.run(ctx, "rev-parse", "--is-inside-work-tree")
	if err != nil {
		return "", nil
	}
	if strings.TrimSpace(out) != "true" {
		return "", nil
	}
	head, err := r.run(ctx, "rev-parse", "HEAD")
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(head), nil
}
