package install

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrchatam/Trace/internal/agents"
	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

const subagentHookSlug = "hook:harness:subagent"

// InstallAgentDefaults upserts bundled harness agent profiles and capability stubs into .trace/.
func InstallAgentDefaults(opts InstallOpts) error {
	root := projectRoot(opts)
	st, err := store.Open(root)
	if err != nil {
		return fmt.Errorf("install: agents: %w", err)
	}
	defer st.Close()

	var catalog agents.DefaultCatalog
	if catalogPath := strings.TrimSpace(opts.CatalogPath); catalogPath != "" {
		catalog, err = agents.LoadDefaultCatalog(catalogPath)
	} else {
		catalog, err = agents.LoadEmbeddedDefaultCatalog()
	}
	if err != nil {
		return fmt.Errorf("install: agents: %w", err)
	}

	svc := domain.New(st)
	ctx := context.Background()
	errOut := opts.ErrOut
	if errOut == nil {
		errOut = io.Discard
	}

	for _, bundle := range catalog.Agents {
		if err := upsertBundledAgent(ctx, svc, catalog.RegistryVersion, bundle); err != nil {
			return fmt.Errorf("install: agents: %w", err)
		}
		fmt.Fprintf(errOut, "install: agents: upserted %s\n", bundle.Slug)
	}

	if err := upsertSubagentHook(ctx, svc, opts); err != nil {
		return fmt.Errorf("install: agents: %w", err)
	}
	fmt.Fprintf(errOut, "install: agents: declared %s\n", subagentHookSlug)
	PrintBootstrapHintIfNeeded(root, errOut)
	return nil
}

func upsertBundledAgent(ctx context.Context, svc *domain.Service, registryVersion string, bundle agents.HarnessAgentBundle) error {
	phases, err := jsonArray(bundle.DeliberationPhases)
	if err != nil {
		return err
	}
	keywords, err := jsonArray(bundle.TaskKeywords)
	if err != nil {
		return err
	}
	source := strings.TrimSpace(bundle.RegistrySource)
	if source == "" {
		source = "bundled"
	}

	id := ""
	if existing, err := svc.GetHarnessAgentBySlug(ctx, bundle.Slug); err == nil {
		id = existing.ID
	}

	if _, err := svc.UpsertHarnessAgent(ctx, domain.HarnessAgentInput{
		ID:                 id,
		Slug:               bundle.Slug,
		Title:              bundle.Title,
		Description:        bundle.Description,
		SubagentType:       bundle.SubagentType,
		DeliberationPhases: phases,
		TaskKeywords:       keywords,
		RecommendSubagent:  bundle.RecommendSubagent,
		RegistrySource:     source,
		RegistryVersion:    registryVersion,
		ExternalURL:        bundle.ExternalURL,
		Requirements:       bundle.Requirements,
	}); err != nil {
		return fmt.Errorf("upsert %s: %w", bundle.Slug, err)
	}

	for _, slug := range bundle.Requirements {
		if err := ensureCapabilityStub(ctx, svc, slug); err != nil {
			return err
		}
	}
	return nil
}

func ensureCapabilityStub(ctx context.Context, svc *domain.Service, slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil
	}
	if _, err := svc.GetCapabilityBySlug(ctx, slug); err == nil {
		return nil
	}
	kind, title := capabilityKindTitle(slug)
	_, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind:   kind,
		Slug:   slug,
		Title:  title,
		Status: domain.CapabilityStatusUnknown,
	})
	if err != nil {
		return fmt.Errorf("stub capability %s: %w", slug, err)
	}
	return nil
}

func capabilityKindTitle(slug string) (kind, title string) {
	switch {
	case strings.HasPrefix(slug, "skill:"):
		return domain.CapabilityKindSkill, slug
	case strings.HasPrefix(slug, "mcp:"):
		return domain.CapabilityKindMCP, slug
	case strings.HasPrefix(slug, "hook:"):
		return domain.CapabilityKindHook, slug
	default:
		return domain.CapabilityKindSkill, slug
	}
}

func upsertSubagentHook(ctx context.Context, svc *domain.Service, opts InstallOpts) error {
	status := domain.CapabilityStatusUnknown
	if detectSubagentHarness(opts) {
		status = domain.CapabilityStatusAvailable
	}
	_, err := svc.UpsertCapability(ctx, domain.CapabilityInput{
		Kind:   domain.CapabilityKindHook,
		Slug:   subagentHookSlug,
		Title:  "Harness subagent delegation",
		Status: status,
		Body:   "When AVAILABLE, loop recommendations may suggest a fresh subagent for independent review.",
	})
	return err
}

func detectSubagentHarness(opts InstallOpts) bool {
	if os.Getenv("TRACE_HARNESS_SUBAGENT") == "1" {
		return true
	}
	root := projectRoot(opts)
	if isDir(filepath.Join(root, ".cursor")) {
		return true
	}
	home := opts.HomeDir
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	if home != "" && isDir(filepath.Join(home, ".cursor", "rules")) {
		return true
	}
	return false
}

func isDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

func jsonArray(v []string) (string, error) {
	if v == nil {
		v = []string{}
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
