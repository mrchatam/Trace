package domain

import (
	"context"
	"sort"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

const (
	evidenceQueryDefaultLimit = 32
	evidenceQueryMaxLimit     = 64
)

// EvidenceQueryOpts bounds list queries (failed/worked/regressions).
type EvidenceQueryOpts struct {
	TaskID   string
	ChangeID string
	Limit    int
}

// ValidatingTest is a test artifact that validates a symbol or file.
type ValidatingTest struct {
	EdgeProvenance string `json:"edge_provenance"`
	TestFileID     string `json:"test_file_id"`
	TestFilePath   string `json:"test_file_path"`
	TestSymbolID   string `json:"test_symbol_id,omitempty"`
	TestSymbolName string `json:"test_symbol_name,omitempty"`
}

// FailedOutcomeRow is a stored kind=test row with fail or error status.
type FailedOutcomeRow struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	TestName   string `json:"test_name"`
	TestStatus string `json:"test_status"`
	Summary    string `json:"summary,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// WorkedApproachRow is an improvement or passing test outcome.
type WorkedApproachRow struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Summary   string `json:"summary,omitempty"`
	TestName  string `json:"test_name,omitempty"`
	TaskID    string `json:"task_id,omitempty"`
	CreatedAt string `json:"created_at"`
}

// RegressionQueryRow is a regression with optional associated change ids.
type RegressionQueryRow struct {
	ID          string   `json:"id"`
	TaskID      string   `json:"task_id"`
	Dimension   string   `json:"dimension"`
	Attribution string   `json:"attribution"`
	Summary     string   `json:"summary"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"created_at"`
	ChangeIDs   []string `json:"change_ids,omitempty"`
}

func clampEvidenceLimit(limit int) int {
	if limit <= 0 {
		return evidenceQueryDefaultLimit
	}
	if limit > evidenceQueryMaxLimit {
		return evidenceQueryMaxLimit
	}
	return limit
}

func (s *Service) validatingTestFromEdge(edge store.CodeEdge) (ValidatingTest, error) {
	fromFile, err := s.store.GetFileByID(edge.FromFileID)
	if err != nil {
		return ValidatingTest{}, err
	}
	row := ValidatingTest{
		EdgeProvenance: edge.Provenance,
		TestFileID:     fromFile.ID,
		TestFilePath:   fromFile.Path,
	}
	if edge.FromSymbolID != nil && strings.TrimSpace(*edge.FromSymbolID) != "" {
		sym, _, err := s.store.GetSymbolByID(*edge.FromSymbolID)
		if err != nil {
			return ValidatingTest{}, err
		}
		row.TestSymbolID = sym.ID
		row.TestSymbolName = sym.Name
	}
	return row, nil
}

// ListTestsValidatingSymbol returns tests with reverse validates edges to symbolID.
func (s *Service) ListTestsValidatingSymbol(ctx context.Context, symbolID string) ([]ValidatingTest, error) {
	_ = ctx
	symbolID = strings.TrimSpace(symbolID)
	if symbolID == "" {
		return nil, &ErrValidation{Msg: "symbol id is required"}
	}
	if _, _, err := s.store.GetSymbolByID(symbolID); err != nil {
		return nil, err
	}
	edges, err := s.store.ListValidatesForSymbol(symbolID)
	if err != nil {
		return nil, err
	}
	out := make([]ValidatingTest, 0, len(edges))
	for _, e := range edges {
		row, err := s.validatingTestFromEdge(e)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, nil
}

// ListTestsValidatingFile resolves path to a file and returns validating tests.
func (s *Service) ListTestsValidatingFile(ctx context.Context, path string) ([]ValidatingTest, error) {
	_ = ctx
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, &ErrValidation{Msg: "file path is required"}
	}
	f, err := s.store.GetFileByPath(path)
	if err != nil {
		return nil, err
	}
	edges, err := s.store.ListValidatesForFile(f.ID)
	if err != nil {
		return nil, err
	}
	out := make([]ValidatingTest, 0, len(edges))
	for _, e := range edges {
		row, err := s.validatingTestFromEdge(e)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, nil
}

// ListFailedOutcomes returns kind=test rows with test_status fail or error, newest first.
func (s *Service) ListFailedOutcomes(ctx context.Context, opts EvidenceQueryOpts) ([]FailedOutcomeRow, error) {
	_ = ctx
	taskID := strings.TrimSpace(opts.TaskID)
	if taskID != "" {
		if _, err := s.store.GetTask(taskID); err != nil {
			return nil, err
		}
	}
	rows, err := s.store.ListFailedTestOutcomes(clampEvidenceLimit(opts.Limit), taskID)
	if err != nil {
		return nil, err
	}
	out := make([]FailedOutcomeRow, 0, len(rows))
	for _, o := range rows {
		out = append(out, FailedOutcomeRow{
			ID: o.ID, TaskID: o.TaskID, TestName: o.TestName,
			TestStatus: o.TestStatus, Summary: o.Summary, CreatedAt: o.CreatedAt,
		})
	}
	return out, nil
}

type workedCandidate struct {
	id, kind, summary, testName, taskID, createdAt string
}

func (s *Service) listWorkedCandidates(taskID string) ([]workedCandidate, error) {
	var candidates []workedCandidate
	if taskID != "" {
		imps, err := s.store.ListImprovementsByTaskID(taskID)
		if err != nil {
			return nil, err
		}
		for _, imp := range imps {
			candidates = append(candidates, workedCandidate{
				id: imp.ID, kind: "improvement", summary: imp.Summary,
				taskID: imp.TaskID, createdAt: imp.CreatedAt,
			})
		}
		passes, err := s.store.ListPassingTestOutcomes(0, taskID)
		if err != nil {
			return nil, err
		}
		for _, o := range passes {
			candidates = append(candidates, workedCandidate{
				id: o.ID, kind: "test_pass", testName: o.TestName,
				taskID: o.TaskID, createdAt: o.CreatedAt,
			})
		}
		return candidates, nil
	}
	imps, err := s.store.ListAllImprovements()
	if err != nil {
		return nil, err
	}
	for _, imp := range imps {
		candidates = append(candidates, workedCandidate{
			id: imp.ID, kind: "improvement", summary: imp.Summary,
			taskID: imp.TaskID, createdAt: imp.CreatedAt,
		})
	}
	passes, err := s.store.ListPassingTestOutcomes(0, "")
	if err != nil {
		return nil, err
	}
	for _, o := range passes {
		candidates = append(candidates, workedCandidate{
			id: o.ID, kind: "test_pass", testName: o.TestName,
			taskID: o.TaskID, createdAt: o.CreatedAt,
		})
	}
	return candidates, nil
}

// ListWorkedApproaches unions improvements and passing test outcomes; deduped by id, newest first.
func (s *Service) ListWorkedApproaches(ctx context.Context, opts EvidenceQueryOpts) ([]WorkedApproachRow, error) {
	_ = ctx
	taskID := strings.TrimSpace(opts.TaskID)
	if taskID != "" {
		if _, err := s.store.GetTask(taskID); err != nil {
			return nil, err
		}
	}
	candidates, err := s.listWorkedCandidates(taskID)
	if err != nil {
		return nil, err
	}
	sort.Slice(candidates, func(i, j int) bool {
		a, b := candidates[i], candidates[j]
		if a.createdAt != b.createdAt {
			return a.createdAt > b.createdAt
		}
		return a.id > b.id
	})
	seen := map[string]struct{}{}
	limit := clampEvidenceLimit(opts.Limit)
	out := make([]WorkedApproachRow, 0, limit)
	for _, c := range candidates {
		if _, dup := seen[c.id]; dup {
			continue
		}
		seen[c.id] = struct{}{}
		out = append(out, WorkedApproachRow{
			ID: c.id, Kind: c.kind, Summary: c.summary, TestName: c.testName,
			TaskID: c.taskID, CreatedAt: c.createdAt,
		})
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (s *Service) regressionChangeIDs(regressionID string) ([]string, error) {
	links, err := s.store.ListLinksFrom(EntityRegression, regressionID)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, l := range links {
		if l.Rel == RelRegressionAssociatedChange && l.ToType == EntityChange {
			ids = append(ids, l.ToID)
		}
	}
	return ids, nil
}

func (s *Service) mapRegressionRows(rows []store.Regression) ([]RegressionQueryRow, error) {
	out := make([]RegressionQueryRow, 0, len(rows))
	for _, r := range rows {
		changeIDs, err := s.regressionChangeIDs(r.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, RegressionQueryRow{
			ID: r.ID, TaskID: r.TaskID, Dimension: r.Dimension,
			Attribution: r.Attribution, Summary: r.Summary, Status: r.Status,
			CreatedAt: r.CreatedAt, ChangeIDs: changeIDs,
		})
	}
	return out, nil
}

// ListRegressions returns regressions filtered by task or change, else project-wide newest-first.
func (s *Service) ListRegressions(ctx context.Context, opts EvidenceQueryOpts) ([]RegressionQueryRow, error) {
	_ = ctx
	changeID := strings.TrimSpace(opts.ChangeID)
	taskID := strings.TrimSpace(opts.TaskID)
	limit := clampEvidenceLimit(opts.Limit)

	if changeID != "" {
		if _, err := s.store.GetChange(changeID); err != nil {
			return nil, err
		}
		rows, err := s.store.ListRegressionsByChangeID(changeID)
		if err != nil {
			return nil, err
		}
		sort.Slice(rows, func(i, j int) bool {
			if rows[i].CreatedAt != rows[j].CreatedAt {
				return rows[i].CreatedAt > rows[j].CreatedAt
			}
			return rows[i].ID > rows[j].ID
		})
		if len(rows) > limit {
			rows = rows[:limit]
		}
		return s.mapRegressionRows(rows)
	}
	if taskID != "" {
		if _, err := s.store.GetTask(taskID); err != nil {
			return nil, err
		}
	}
	rows, err := s.store.ListRegressionsRecent(limit, taskID)
	if err != nil {
		return nil, err
	}
	return s.mapRegressionRows(rows)
}
