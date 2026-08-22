package domain

import "fmt"

// SeedExportHonestyViolation is a document-level graph honesty failure for seed export.
type SeedExportHonestyViolation struct {
	Message string
}

func seedLinkFromID(l SeedLink) string {
	if l.From != "" {
		return l.From
	}
	return l.FromID
}

// CollectSeedDocumentHonestyViolations checks export document richness and orphan causal links.
func CollectSeedDocumentHonestyViolations(doc SeedDocument) []SeedExportHonestyViolation {
	var out []SeedExportHonestyViolation

	if len(doc.Discoveries) == 0 && len(doc.Decisions) == 0 {
		out = append(out, SeedExportHonestyViolation{
			Message: "graph honesty: thin graph (discoveries=0 decisions=0); require discoveries≥1 OR decisions≥1",
		})
	}

	mentionsTask := map[string]bool{}
	affectsTask := map[string]bool{}
	for _, l := range doc.Links {
		from := seedLinkFromID(l)
		switch l.Rel {
		case RelDiscoveryMentionsTask, "discovery-mentions-task":
			if from != "" {
				mentionsTask[from] = true
			}
		case RelDecisionAffectsTask, "decision-task":
			if from != "" {
				affectsTask[from] = true
			}
		}
	}

	for _, d := range doc.Discoveries {
		if !mentionsTask[d.ID] {
			out = append(out, SeedExportHonestyViolation{
				Message: fmt.Sprintf("graph honesty: discovery %s missing discovery_mentions_task link", d.ID),
			})
		}
	}
	for _, d := range doc.Decisions {
		if !affectsTask[d.ID] {
			out = append(out, SeedExportHonestyViolation{
				Message: fmt.Sprintf("graph honesty: decision %s missing decision_affects_task link", d.ID),
			})
		}
	}

	return out
}
