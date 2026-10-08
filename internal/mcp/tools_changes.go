package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/store"
)

// ChangesInput mirrors `trace changes list|show|compare`.
type ChangesInput struct {
	Project  string  `json:"project,omitempty" jsonschema:"optional project root override"`
	Action   string  `json:"action" jsonschema:"list|show|compare"`
	TaskID   string  `json:"task_id,omitempty" jsonschema:"list: optional task UUID filter"`
	Limit    float64 `json:"limit,omitempty" jsonschema:"list: optional max rows (default 32, cap 64)"`
	ChangeID string  `json:"change_id,omitempty" jsonschema:"show: change UUID"`
	From     string  `json:"from,omitempty" jsonschema:"compare: start git OID"`
	To       string  `json:"to,omitempty" jsonschema:"compare: end git OID"`
}

type changeListRow struct {
	ID         string `json:"id"`
	TaskID     string `json:"task_id"`
	GitCommit  string `json:"git_commit"`
	Status     string `json:"status"`
	Reason     string `json:"reason"`
	SourceType string `json:"source_type"`
	CreatedAt  string `json:"created_at"`
}

type changeShowPathRow struct {
	Path     string `json:"path"`
	Status   string `json:"status"`
	SymbolID string `json:"symbol_id"`
}

// changeShowRow is the tagged wire form of store.Change. Marshalling the store
// struct directly emitted PascalCase keys (Actor, GitCommit, RowID, ...) inside
// the same document where `list` and `paths` are snake_case (#161). RowID is
// an internal SQLite rowid and is deliberately not on the wire.
type changeShowRow struct {
	ID             string  `json:"id"`
	TaskID         string  `json:"task_id"`
	GitCommit      string  `json:"git_commit"`
	ParentChangeID string  `json:"parent_change_id,omitempty"`
	Actor          string  `json:"actor,omitempty"`
	Reason         string  `json:"reason,omitempty"`
	Status         string  `json:"status"`
	SourceType     string  `json:"source_type,omitempty"`
	Confidence     float64 `json:"confidence"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	LastVerifiedAt *string `json:"last_verified_at,omitempty"`
}

func changeShowRowFrom(c store.Change) changeShowRow {
	return changeShowRow{
		ID:             c.ID,
		TaskID:         c.TaskID,
		GitCommit:      c.GitCommit,
		ParentChangeID: c.ParentChangeID,
		Actor:          c.Actor,
		Reason:         c.Reason,
		Status:         c.Status,
		SourceType:     c.SourceType,
		Confidence:     c.Confidence,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
		LastVerifiedAt: c.LastVerifiedAt,
	}
}

type changeShowResponse struct {
	OK     bool                `json:"ok"`
	Change changeShowRow       `json:"change"`
	Paths  []changeShowPathRow `json:"paths"`
}

func (s *Server) toolChanges(ctx context.Context, _ *sdkmcp.CallToolRequest, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	_, st, err := s.openStore(in.Project)
	if err != nil {
		return nil, nil, err
	}
	defer st.Close()
	if err := assertMCPToolAllowed(ctx, st, "trace_changes"); err != nil {
		return nil, nil, err
	}

	switch strings.ToLower(strings.TrimSpace(in.Action)) {
	case "list":
		return s.changesList(ctx, st, in)
	case "show":
		return s.changesShow(ctx, st, in)
	case "compare":
		return s.changesCompare(ctx, st, in)
	default:
		return nil, nil, fmt.Errorf("trace_changes: action must be list|show|compare")
	}
}

// changesListLimit mirrors the store's default/cap so the response can report
// the effective limit rather than silently clamping.
func changesListLimit(requested int) (int, bool) {
	if requested > 0 {
		if requested > 64 {
			return 64, true
		}
		return requested, true
	}
	return 32, true
}

func (s *Server) changesList(ctx context.Context, st *store.Store, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	svc := domain.New(st)
	taskID := strings.TrimSpace(in.TaskID)
	if taskID != "" {
		// Fail closed on a task that does not exist (#154 class): an empty
		// change list is only authoritative when the task is real.
		if _, err := svc.GetTask(ctx, taskID); err != nil {
			return nil, nil, wrapToolErr("trace_changes list", err)
		}
	}
	changes, err := svc.ListChanges(ctx, int(in.Limit), in.TaskID)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_changes list: %w", err)
	}
	out := make([]changeListRow, 0, len(changes))
	for _, c := range changes {
		out = append(out, changeListRow{
			ID:         c.ID,
			TaskID:     c.TaskID,
			GitCommit:  c.GitCommit,
			Status:     c.Status,
			Reason:     c.Reason,
			SourceType: c.SourceType,
			CreatedAt:  c.CreatedAt,
		})
	}
	effLimit, _ := changesListLimit(int(in.Limit))
	// truncated is the page-fill signal: the page returned exactly the
	// effective limit, so there may be more rows. The store has no cursor
	// plumbing, so a definitive answer is not available here.
	env := newListEnvelope(out, len(out), len(out) == effLimit).withLimit(effLimit)
	b, err := json.Marshal(env)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) changesShow(ctx context.Context, st *store.Store, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	changeID := strings.TrimSpace(in.ChangeID)
	if changeID == "" {
		return nil, nil, fmt.Errorf("trace_changes show: change_id is required")
	}
	svc := domain.New(st)
	c, err := svc.GetChange(ctx, changeID)
	if err != nil {
		return nil, nil, wrapToolErr("trace_changes show", err)
	}
	paths, err := svc.ListChangePaths(ctx, changeID)
	if err != nil {
		return nil, nil, wrapToolErr("trace_changes show", err)
	}
	pathRows := make([]changeShowPathRow, 0, len(paths))
	for _, p := range paths {
		pathRows = append(pathRows, changeShowPathRow{
			Path:     p.Path,
			Status:   p.Status,
			SymbolID: p.SymbolID,
		})
	}
	resp := changeShowResponse{OK: true, Change: changeShowRowFrom(c), Paths: pathRows}
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}

func (s *Server) changesCompare(ctx context.Context, st *store.Store, in ChangesInput) (*sdkmcp.CallToolResult, any, error) {
	from := strings.TrimSpace(in.From)
	to := strings.TrimSpace(in.To)
	if from == "" || to == "" {
		return nil, nil, fmt.Errorf("trace_changes compare: from and to are required")
	}
	result, err := domain.New(st).CompareStates(ctx, from, to)
	if err != nil {
		return nil, nil, fmt.Errorf("trace_changes compare: %w", err)
	}
	b, err := json.Marshal(result)
	if err != nil {
		return nil, nil, err
	}
	return textResult(string(b)), nil, nil
}
