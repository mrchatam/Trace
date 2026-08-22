package analyzers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
)

var conventionalLayers = []string{"cmd", "internal", "pkg"}

type architectureOverlay struct {
	Layers []architectureOverlayLayer `json:"layers"`
}

type architectureOverlayLayer struct {
	Name     string   `json:"name"`
	Prefixes []string `json:"prefixes"`
}

func computeArchitecturalBoundaryEdges(st *store.Store, p string) ([]store.CodeEdge, error) {
	p = store.NormalizePath(p)
	f, err := st.GetFileByPath(p)
	if err != nil {
		return nil, err
	}
	layer, prov, ok, err := resolveArchitecturalLayer(st, p)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	toPath := store.ArchitectureLayerPath(layer)
	if toPath == "" {
		return nil, nil
	}
	hash := sha256Hex([]byte("trace-layer:" + layer))
	toFile, err := st.UpsertFile(toPath, hash, nil)
	if err != nil {
		return nil, err
	}
	return []store.CodeEdge{{
		FromFileID: f.ID,
		ToFileID:   toFile.ID,
		Rel:        store.RelArchitecturalBoundary,
		Provenance: prov,
	}}, nil
}

func resolveArchitecturalLayer(st *store.Store, p string) (layer, provenance string, ok bool, err error) {
	root := st.ProjectRoot()
	overlay, err := loadArchitectureOverlay(root)
	if err != nil {
		return "", "", false, err
	}
	if overlay != nil {
		if name, matched := overlayLayerForPath(overlay, p); matched {
			return name, store.ImportProvenanceExtracted, true, nil
		}
	}
	layer, ok = conventionalLayerForPath(p)
	if !ok {
		return "", "", false, nil
	}
	if goModPresent(root) {
		return layer, store.ImportProvenanceExtracted, true, nil
	}
	return layer, store.ImportProvenanceInferred, true, nil
}

func conventionalLayerForPath(p string) (string, bool) {
	p = store.NormalizePath(p)
	for _, seg := range strings.Split(p, "/") {
		for _, layer := range conventionalLayers {
			if seg == layer {
				return layer, true
			}
		}
	}
	return "", false
}

func goModPresent(root string) bool {
	if root == "" {
		return false
	}
	b, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "module ") {
			return true
		}
	}
	return false
}

func loadArchitectureOverlay(root string) (*architectureOverlay, error) {
	if root == "" {
		return nil, nil
	}
	b, err := os.ReadFile(filepath.Join(root, "trace", "architecture.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var overlay architectureOverlay
	if err := json.Unmarshal(b, &overlay); err != nil {
		return nil, err
	}
	return &overlay, nil
}

func overlayLayerForPath(overlay *architectureOverlay, p string) (string, bool) {
	if overlay == nil {
		return "", false
	}
	p = store.NormalizePath(p)
	bestLen := -1
	best := ""
	for _, layer := range overlay.Layers {
		if layer.Name == "" {
			continue
		}
		for _, prefix := range layer.Prefixes {
			prefix = strings.Trim(store.NormalizePath(prefix), "/")
			if prefix == "" || !pathHasDirPrefix(p, prefix) {
				continue
			}
			if len(prefix) > bestLen {
				bestLen = len(prefix)
				best = layer.Name
			}
		}
	}
	if best == "" {
		return "", false
	}
	return best, true
}

func pathHasDirPrefix(p, prefix string) bool {
	if p == prefix {
		return true
	}
	return strings.HasPrefix(p, prefix+"/")
}
