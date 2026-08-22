package domain_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/mrchatam/Trace/internal/domain"
	"github.com/mrchatam/Trace/internal/retrieval"
	"github.com/mrchatam/Trace/internal/store"
)

func openDomainWithImpact(t *testing.T) (*domain.Service, *store.Store) {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	svc := domain.New(st)
	retrieval.WireDomainImpactWalker(svc, retrieval.New(st))
	return svc, st
}

func mustChangeWithPath(t *testing.T, svc *domain.Service, st *store.Store, path string, symbolID string) store.Change {
	t.Helper()
	task := mustChangeTask(t, svc)
	c, err := svc.CreateChange(context.Background(), domain.ChangeInput{
		TaskID: task.ID,
		Reason: "impact compare fixture",
		Paths:  []domain.ChangePathInput{{Path: path, SymbolID: symbolID}},
	})
	if err != nil {
		t.Fatalf("CreateChange: %v", err)
	}
	return c
}

func walkSnapshot(t *testing.T, st *store.Store, seeds []domain.ImpactSeedRef, depth int) *domain.ImpactWalkSnapshot {
	t.Helper()
	eng := retrieval.New(st)
	walker := retrieval.DomainImpactWalker{Engine: eng}
	snap, err := walker.ImpactWalk(context.Background(), seeds, depth)
	if err != nil {
		t.Fatalf("ImpactWalk: %v", err)
	}
	return snap
}

func TestRecordPredictedImpactThenCompareActual(t *testing.T) {
	svc, st := openDomainWithImpact(t)
	ctx := context.Background()

	lib, err := st.UpsertFile("pkg/lib.go", "hlib", nil)
	if err != nil {
		t.Fatal(err)
	}
	consumer, err := st.UpsertFile("pkg/use.go", "huse", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileImports("pkg/use.go", []store.Import{{ImportedPath: "pkg/lib.go"}}); err != nil {
		t.Fatal(err)
	}

	c := mustChangeWithPath(t, svc, st, "pkg/lib.go", "")

	res := walkSnapshot(t, st, []domain.ImpactSeedRef{
		{EntityType: "file", EntityID: lib.ID},
	}, 1)
	if _, err := svc.RecordPredictedImpact(ctx, c.ID, res); err != nil {
		t.Fatalf("RecordPredictedImpact: %v", err)
	}

	extra, err := st.UpsertFile("pkg/extra.go", "hex", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileImports("pkg/extra.go", []store.Import{{ImportedPath: "pkg/lib.go"}}); err != nil {
		t.Fatal(err)
	}
	_ = consumer
	_ = extra

	cmp, err := svc.CompareActualImpact(ctx, c.ID)
	if err != nil {
		t.Fatalf("CompareActualImpact: %v", err)
	}
	if cmp.ChangeID != c.ID || cmp.ComparedAt == "" {
		t.Fatalf("result: %+v", cmp)
	}
	if len(cmp.Delta.Matched) == 0 && len(cmp.Delta.Unexpected) == 0 && len(cmp.Delta.Missed) == 0 {
		t.Fatalf("empty delta: %+v", cmp.Delta)
	}

	row, err := st.GetImpactPrediction(c.ID)
	if err != nil {
		t.Fatal(err)
	}
	if row.CompareJSON == "" || row.ComparedAt == "" {
		t.Fatalf("persisted compare: json=%q at=%q", row.CompareJSON, row.ComparedAt)
	}
	var stored domain.ImpactCompareDelta
	if err := json.Unmarshal([]byte(row.CompareJSON), &stored); err != nil {
		t.Fatal(err)
	}
	if len(stored.Unexpected) < 1 {
		t.Fatalf("want unexpected key after graph growth: %+v", stored)
	}
}

func TestImpactCompareUnexpectedAndMissed(t *testing.T) {
	svc, st := openDomainWithImpact(t)
	ctx := context.Background()

	seed, err := st.UpsertFile("core/a.go", "ha", nil)
	if err != nil {
		t.Fatal(err)
	}
	predictedOnly, err := st.UpsertFile("core/old.go", "hold", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileImports("core/old.go", []store.Import{{ImportedPath: "core/a.go"}}); err != nil {
		t.Fatal(err)
	}

	c := mustChangeWithPath(t, svc, st, "core/a.go", "")

	resPredict := walkSnapshot(t, st, []domain.ImpactSeedRef{
		{EntityType: "file", EntityID: seed.ID},
	}, 1)
	predictedKey := domain.EntityKey("file", predictedOnly.ID)
	foundPredicted := false
	for _, h := range resPredict.Blast {
		if domain.EntityKey(h.EntityType, h.EntityID) == predictedKey {
			foundPredicted = true
		}
	}
	if !foundPredicted {
		t.Fatalf("predicted blast missing %s: %+v", predictedKey, resPredict.Blast)
	}
	if _, err := svc.RecordPredictedImpact(ctx, c.ID, resPredict); err != nil {
		t.Fatal(err)
	}

	if err := st.ReplaceFileImports("core/old.go", nil); err != nil {
		t.Fatal(err)
	}
	unexpected, err := st.UpsertFile("core/new.go", "hnew", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileImports("core/new.go", []store.Import{{ImportedPath: "core/a.go"}}); err != nil {
		t.Fatal(err)
	}

	cmp, err := svc.CompareActualImpact(ctx, c.ID)
	if err != nil {
		t.Fatalf("CompareActualImpact: %v", err)
	}
	unexpectedKey := domain.EntityKey("file", unexpected.ID)
	if !containsString(cmp.Delta.Missed, predictedKey) {
		t.Fatalf("missed want %q: %+v", predictedKey, cmp.Delta.Missed)
	}
	if !containsString(cmp.Delta.Unexpected, unexpectedKey) {
		t.Fatalf("unexpected want %q: %+v", unexpectedKey, cmp.Delta.Unexpected)
	}
}

func TestImpactCompareFailClosedWithoutPrediction(t *testing.T) {
	svc, st := openDomainWithImpact(t)
	ctx := context.Background()

	c := mustChangeWithPath(t, svc, st, "missing/index.go", "")
	if _, err := st.UpsertFile("missing/index.go", "hx", nil); err != nil {
		t.Fatal(err)
	}

	_, err := svc.CompareActualImpact(ctx, c.ID)
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("want ErrValidation, got %v", err)
	}
	if ve.Msg != "no impact prediction for change" {
		t.Fatalf("msg: %q", ve.Msg)
	}
}

func TestPredictImpactForChangeFailClosedUnindexedPath(t *testing.T) {
	svc, _ := openDomainWithImpact(t)
	ctx := context.Background()

	task := mustChangeTask(t, svc)
	c, err := svc.CreateChange(ctx, domain.ChangeInput{
		TaskID: task.ID,
		Reason: "unindexed",
		Paths:  []domain.ChangePathInput{{Path: "not/indexed.go"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = svc.PredictImpactForChange(ctx, c.ID, 1)
	var ve *domain.ErrValidation
	if !errors.As(err, &ve) {
		t.Fatalf("want ErrValidation, got %v", err)
	}
	if ve.Msg != `path "not/indexed.go" is not indexed` {
		t.Fatalf("msg: %q", ve.Msg)
	}
}

func TestPredictImpactForChangeUpsertReplaces(t *testing.T) {
	svc, st := openDomainWithImpact(t)
	ctx := context.Background()

	lib, err := st.UpsertFile("pkg/up.go", "hup", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := mustChangeWithPath(t, svc, st, "pkg/up.go", "")

	first, err := svc.PredictImpactForChange(ctx, c.ID, 1)
	if err != nil {
		t.Fatalf("first predict: %v", err)
	}
	firstAt := first.CreatedAt

	neighbor, err := st.UpsertFile("pkg/neighbor.go", "hn", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileImports("pkg/neighbor.go", []store.Import{{ImportedPath: "pkg/up.go"}}); err != nil {
		t.Fatal(err)
	}
	_ = lib
	_ = neighbor

	second, err := svc.PredictImpactForChange(ctx, c.ID, 1)
	if err != nil {
		t.Fatalf("second predict: %v", err)
	}
	if second.CreatedAt == "" || second.CreatedAt == firstAt {
		t.Fatalf("upsert should refresh created_at: first=%q second=%q", firstAt, second.CreatedAt)
	}
	if second.CompareJSON != "" || second.ComparedAt != "" {
		t.Fatalf("re-predict should clear compare: %+v", second)
	}

	var payload domain.PredictedImpactPayload
	if err := json.Unmarshal([]byte(second.PredictedJSON), &payload); err != nil {
		t.Fatal(err)
	}
	neighborKey := domain.EntityKey("file", neighbor.ID)
	if !containsString(payload.BlastKeys, neighborKey) {
		t.Fatalf("second predict blast_keys want %q: %+v", neighborKey, payload.BlastKeys)
	}
}

func TestCompareUsesStoredDepth(t *testing.T) {
	svc, st := openDomainWithImpact(t)
	ctx := context.Background()

	a, err := st.UpsertFile("chain/a.go", "ha", nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := st.UpsertFile("chain/b.go", "hb", nil)
	if err != nil {
		t.Fatal(err)
	}
	cFile, err := st.UpsertFile("chain/c.go", "hc", nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileImports("chain/b.go", []store.Import{{ImportedPath: "chain/a.go"}}); err != nil {
		t.Fatal(err)
	}
	if err := st.ReplaceFileImports("chain/c.go", []store.Import{{ImportedPath: "chain/b.go"}}); err != nil {
		t.Fatal(err)
	}

	change := mustChangeWithPath(t, svc, st, "chain/a.go", "")
	if _, err := svc.PredictImpactForChange(ctx, change.ID, 1); err != nil {
		t.Fatal(err)
	}

	cmp, err := svc.CompareActualImpact(ctx, change.ID)
	if err != nil {
		t.Fatal(err)
	}
	cKey := domain.EntityKey("file", cFile.ID)
	if containsString(cmp.Delta.Unexpected, cKey) || containsString(cmp.Delta.Matched, cKey) {
		t.Fatalf("depth-1 predict/compare must not include hop-2 node %q: %+v", cKey, cmp.Delta)
	}
	_ = a
	_ = b
}

func containsString(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}

func TestImpactPredictionStoreRoundtrip(t *testing.T) {
	st, _ := openTempStoreHelper(t)
	changeID := "change-1"
	payload := `{"seeds":[{"entity_type":"file","entity_id":"f1"}],"blast_keys":["file:f2"],"affected_test_keys":[],"depth":1,"blast_total":1,"blast_kept":1,"truncated":false}`
	row, err := st.UpsertImpactPrediction(store.ImpactPrediction{
		ChangeID:      changeID,
		PredictedJSON: payload,
		Depth:         1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if row.ChangeID != changeID {
		t.Fatalf("row: %+v", row)
	}
	got, err := st.GetImpactPrediction(changeID)
	if err != nil {
		t.Fatal(err)
	}
	if got.PredictedJSON != payload {
		t.Fatalf("predicted_json: %q", got.PredictedJSON)
	}
	delta := `{"matched":[],"unexpected":["file:x"],"missed":["file:y"]}`
	updated, err := st.UpdateImpactPredictionCompare(changeID, delta, "2026-08-18T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if updated.CompareJSON != delta || updated.ComparedAt != "2026-08-18T00:00:00Z" {
		t.Fatalf("compare persist: %+v", updated)
	}
	if _, err := st.GetImpactPrediction("missing"); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("missing: %v", err)
	}
}

func openTempStoreHelper(t *testing.T) (*store.Store, string) {
	t.Helper()
	root := t.TempDir()
	s, err := store.Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s, root
}
