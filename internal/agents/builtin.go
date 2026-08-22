package agents

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	traceagents "github.com/mrchatam/Trace/trace/agents"
)

const defaultCatalogSchemaVersion = 1

// HarnessAgentBundle is one profile from the bundled default catalog JSON.
type HarnessAgentBundle struct {
	Slug               string   `json:"slug"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	SubagentType       string   `json:"subagent_type"`
	DeliberationPhases []string `json:"deliberation_phases"`
	TaskKeywords       []string `json:"task_keywords"`
	RecommendSubagent  bool     `json:"recommend_subagent"`
	RegistrySource     string   `json:"registry_source"`
	ExternalURL        string   `json:"external_url,omitempty"`
	Requirements       []string `json:"requirements"`
}

// DefaultCatalog is the bundled harness agent catalog file (schema v1).
type DefaultCatalog struct {
	SchemaVersion   int                  `json:"schema_version"`
	RegistryVersion string               `json:"registry_version"`
	Agents          []HarnessAgentBundle `json:"agents"`
}

// LoadEmbeddedDefaultCatalog parses and validates the committed bundled catalog.
func LoadEmbeddedDefaultCatalog() (DefaultCatalog, error) {
	return parseDefaultCatalog(traceagents.DefaultCatalogJSON, "embedded default catalog")
}

// LoadDefaultCatalog parses and validates bundled catalog JSON from path.
func LoadDefaultCatalog(path string) (DefaultCatalog, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return DefaultCatalog{}, fmt.Errorf("agents: read default catalog: %w", err)
	}
	return parseDefaultCatalog(raw, path)
}

func parseDefaultCatalog(raw []byte, label string) (DefaultCatalog, error) {
	var doc DefaultCatalog
	if err := json.Unmarshal(raw, &doc); err != nil {
		return DefaultCatalog{}, fmt.Errorf("agents: parse %s: %w", label, err)
	}
	if err := validateDefaultCatalog(doc); err != nil {
		return DefaultCatalog{}, err
	}
	return doc, nil
}

func validateDefaultCatalog(doc DefaultCatalog) error {
	if doc.SchemaVersion != defaultCatalogSchemaVersion {
		return fmt.Errorf("agents: default catalog schema_version want %d got %d", defaultCatalogSchemaVersion, doc.SchemaVersion)
	}
	if strings.TrimSpace(doc.RegistryVersion) == "" {
		return fmt.Errorf("agents: default catalog registry_version is required")
	}
	if len(doc.Agents) == 0 {
		return fmt.Errorf("agents: default catalog must include at least one agent")
	}
	seen := map[string]bool{}
	for i, a := range doc.Agents {
		slug := strings.TrimSpace(a.Slug)
		if slug == "" {
			return fmt.Errorf("agents: default catalog agents[%d]: slug is required", i)
		}
		if !strings.HasPrefix(slug, "agent:") {
			return fmt.Errorf("agents: default catalog agents[%d]: slug must start with agent:", i)
		}
		if seen[slug] {
			return fmt.Errorf("agents: default catalog duplicate slug %q", slug)
		}
		seen[slug] = true
		if strings.TrimSpace(a.Title) == "" {
			return fmt.Errorf("agents: default catalog %s: title is required", slug)
		}
		if strings.TrimSpace(a.SubagentType) == "" {
			return fmt.Errorf("agents: default catalog %s: subagent_type is required", slug)
		}
		if a.DeliberationPhases == nil {
			return fmt.Errorf("agents: default catalog %s: deliberation_phases must be a JSON array", slug)
		}
		if a.TaskKeywords == nil {
			return fmt.Errorf("agents: default catalog %s: task_keywords must be a JSON array", slug)
		}
	}
	return nil
}
