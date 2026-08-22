package store

import (
	"database/sql"
	"strings"
	"testing"
)

func TestChangeStoreRoundtrip(t *testing.T) {
	s, _ := openTempStore(t)

	c, err := s.UpsertChange(Change{
		TaskID:    "task-1",
		GitCommit: "abcdef0",
		Actor:     "agent",
		Reason:    "wire store",
		Status:    ChangeStatusRecorded,
	})
	if err != nil {
		t.Fatalf("UpsertChange: %v", err)
	}
	if c.ID == "" || c.CreatedAt == "" || c.UpdatedAt == "" {
		t.Fatalf("expected ids/timestamps: %+v", c)
	}
	if c.Status != ChangeStatusRecorded || c.GitCommit != "abcdef0" {
		t.Fatalf("fields: %+v", c)
	}

	p, err := s.InsertChangePath(ChangePath{
		ChangeID: c.ID,
		Path:     "internal/store/changes.go",
		Status:   "A",
	})
	if err != nil {
		t.Fatalf("InsertChangePath: %v", err)
	}
	if p.Path != "internal/store/changes.go" || p.Status != "A" {
		t.Fatalf("path: %+v", p)
	}

	eff, err := s.UpsertEffect(Effect{
		ChangeID:  c.ID,
		Dimension: "correctness",
		Expected:  "compiles",
	})
	if err != nil {
		t.Fatalf("UpsertEffect: %v", err)
	}
	if eff.Comparison != EffectComparisonNone || eff.Expected != "compiles" {
		t.Fatalf("effect: %+v", eff)
	}

	got, err := s.GetChange(c.ID)
	if err != nil || got.TaskID != "task-1" {
		t.Fatalf("GetChange: %+v err=%v", got, err)
	}
	listed, err := s.ListChangesByTaskID("task-1")
	if err != nil || len(listed) != 1 || listed[0].ID != c.ID {
		t.Fatalf("ListChangesByTaskID: %+v err=%v", listed, err)
	}
	paths, err := s.ListChangePaths(c.ID)
	if err != nil || len(paths) != 1 || paths[0].Path != p.Path {
		t.Fatalf("ListChangePaths: %+v err=%v", paths, err)
	}
	effects, err := s.ListEffectsByChangeID(c.ID)
	if err != nil || len(effects) != 1 || effects[0].ID != eff.ID {
		t.Fatalf("ListEffectsByChangeID: %+v err=%v", effects, err)
	}
	byDim, err := s.GetEffectByChangeDimension(c.ID, "correctness")
	if err != nil || byDim.ID != eff.ID {
		t.Fatalf("GetEffectByChangeDimension: %+v err=%v", byDim, err)
	}

	if _, err := s.UpsertChange(Change{}); err == nil {
		t.Fatal("expected task_id required")
	}
	if _, err := s.GetChange(""); err == nil {
		t.Fatal("expected id required")
	}
	if _, err := s.GetChange("missing"); err == nil {
		t.Fatal("expected missing change")
	}
	if _, err := s.InsertChangePath(ChangePath{ChangeID: c.ID}); err == nil {
		t.Fatal("expected path required")
	}
}

func TestChangeSchemaHasNoBlobOrPatchColumns(t *testing.T) {
	s, _ := openTempStore(t)

	for _, table := range []string{"changes", "change_paths", "effects"} {
		var name string
		if err := s.db.QueryRow(`SELECT name FROM sqlite_master WHERE name=?`, table).Scan(&name); err != nil {
			t.Fatalf("missing table %s: %v", table, err)
		}
		cols, err := tableColumns(s.conn, table)
		if err != nil {
			t.Fatal(err)
		}
		if len(cols) == 0 {
			t.Fatalf("%s has no columns", table)
		}
		for _, c := range cols {
			lower := strings.ToLower(c)
			for _, bad := range []string{
				"blob", "patch", "diff", "content", "file_body", "source_text",
				"paths_json", "files_json", "file_content", "raw_source",
			} {
				if strings.Contains(lower, bad) && lower != "confidence" {
					t.Errorf("forbidden column %s.%s", table, c)
				}
			}
		}
	}

	pathCols, err := tableColumns(s.conn, "change_paths")
	if err != nil {
		t.Fatal(err)
	}
	wantPath := map[string]bool{"change_id": true, "path": true, "status": true, "symbol_id": true}
	if len(pathCols) != 4 {
		t.Fatalf("change_paths columns=%v want exactly change_id,path,status,symbol_id", pathCols)
	}
	for _, c := range pathCols {
		if !wantPath[c] {
			t.Errorf("unexpected change_paths column %q", c)
		}
	}

	bad, where, err := s.HasBlobLikeColumns()
	if err != nil {
		t.Fatalf("HasBlobLikeColumns: %v", err)
	}
	if bad {
		t.Fatalf("HasBlobLikeColumns true at %s", where)
	}
}

func tableColumns(db *sql.DB, table string) ([]string, error) {
	rows, err := db.Query(`PRAGMA table_info(` + table + `)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cols []string
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return nil, err
		}
		cols = append(cols, name)
	}
	return cols, rows.Err()
}
