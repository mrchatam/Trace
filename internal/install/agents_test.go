package install_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/install"
	"github.com/mrchatam/Trace/internal/store"
)

func bundledCatalogPath(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		candidate := filepath.Join(dir, "trace", "agents", "default.json")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("trace/agents/default.json not found")
		}
		dir = parent
	}
}

func openProjectDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

func TestInstallAgentsSeedsDefaults(t *testing.T) {
	dir := openProjectDir(t)
	opts := install.InstallOpts{
		Write:       true,
		ProjectRoot: dir,
		ErrOut:      os.Stderr,
	}
	if err := install.InstallAgentDefaults(opts); err != nil {
		t.Fatalf("InstallAgentDefaults: %v", err)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	agents, err := st.ListHarnessAgents()
	if err != nil {
		t.Fatal(err)
	}
	if len(agents) != 6 {
		t.Fatalf("want 6 agents, got %d: %+v", len(agents), agents)
	}
}

func TestBundledProfilesIncludeRequirements(t *testing.T) {
	dir := openProjectDir(t)
	opts := install.InstallOpts{
		Write:       true,
		ProjectRoot: dir,
		CatalogPath: bundledCatalogPath(t),
		ErrOut:      os.Stderr,
	}
	if err := install.InstallAgentDefaults(opts); err != nil {
		t.Fatalf("InstallAgentDefaults: %v", err)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	wantReqs := map[string]int{
		"agent:code-reviewer":        1,
		"agent:performance-reviewer": 1,
		"agent:security-reviewer":    1,
		"agent:nested-reviewer":      1,
		"agent:explore":              1,
		"agent:generalPurpose":       0,
	}
	agents, err := st.ListHarnessAgents()
	if err != nil {
		t.Fatal(err)
	}
	for _, agent := range agents {
		want, ok := wantReqs[agent.Slug]
		if !ok {
			t.Fatalf("unexpected agent slug %q", agent.Slug)
		}
		reqs, err := st.ListHarnessAgentRequirements(agent.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(reqs) != want {
			t.Fatalf("%s requirements: want %d got %d (%+v)", agent.Slug, want, len(reqs), reqs)
		}
		delete(wantReqs, agent.Slug)
	}
	if len(wantReqs) != 0 {
		t.Fatalf("missing agents: %v", wantReqs)
	}
}

func TestInstallAgentsIdempotent(t *testing.T) {
	dir := openProjectDir(t)
	opts := install.InstallOpts{
		Write:       true,
		ProjectRoot: dir,
		CatalogPath: bundledCatalogPath(t),
		ErrOut:      os.Stderr,
	}
	if err := install.InstallAgentDefaults(opts); err != nil {
		t.Fatalf("first install: %v", err)
	}
	if err := install.InstallAgentDefaults(opts); err != nil {
		t.Fatalf("second install: %v", err)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	agents, err := st.ListHarnessAgents()
	if err != nil {
		t.Fatal(err)
	}
	if len(agents) != 6 {
		t.Fatalf("want 6 agents after idempotent install, got %d", len(agents))
	}
	reqs, err := st.ListAllHarnessAgentRequirements()
	if err != nil {
		t.Fatal(err)
	}
	if len(reqs) != 5 {
		t.Fatalf("want 5 requirement rows, got %d", len(reqs))
	}
}

func TestSubagentHookDeclaredOnInstall(t *testing.T) {
	dir := openProjectDir(t)
	opts := install.InstallOpts{
		Write:       true,
		ProjectRoot: dir,
		CatalogPath: bundledCatalogPath(t),
		HomeDir:     t.TempDir(),
		ErrOut:      os.Stderr,
	}
	t.Setenv("TRACE_HARNESS_SUBAGENT", "")
	if err := install.InstallAgentDefaults(opts); err != nil {
		t.Fatalf("InstallAgentDefaults: %v", err)
	}
	st, err := store.Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc := domain.New(st)
	ctx := context.Background()
	hook, err := svc.GetCapabilityBySlug(ctx, "hook:harness:subagent")
	if err != nil {
		t.Fatalf("GetCapabilityBySlug hook: %v", err)
	}
	if hook.Kind != domain.CapabilityKindHook {
		t.Fatalf("hook kind: want HOOK got %q", hook.Kind)
	}
	if hook.Status != domain.CapabilityStatusUnknown {
		t.Fatalf("without harness signals want UNKNOWN, got %q", hook.Status)
	}

	projectDir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(projectDir, ".cursor"), 0o755); err != nil {
		t.Fatal(err)
	}
	opts2 := install.InstallOpts{
		Write:       true,
		ProjectRoot: projectDir,
		CatalogPath: bundledCatalogPath(t),
		HomeDir:     t.TempDir(),
		ErrOut:      os.Stderr,
	}
	if err := install.InstallAgentDefaults(opts2); err != nil {
		t.Fatalf("InstallAgentDefaults with .cursor: %v", err)
	}
	st2, err := store.Open(projectDir)
	if err != nil {
		t.Fatal(err)
	}
	defer st2.Close()
	svc2 := domain.New(st2)
	hook2, err := svc2.GetCapabilityBySlug(ctx, "hook:harness:subagent")
	if err != nil {
		t.Fatal(err)
	}
	if hook2.Status != domain.CapabilityStatusAvailable {
		t.Fatalf("with .cursor want AVAILABLE, got %q", hook2.Status)
	}
}
