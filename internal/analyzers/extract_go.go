package analyzers

import (
	"github.com/mrchatam/Trace/internal/store"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

const goSymbolQuery = `
(function_declaration
  name: (identifier) @name) @function

(method_declaration
  name: (field_identifier) @name) @method

(type_spec
  name: (type_identifier) @name) @type

(type_alias
  name: (type_identifier) @name) @type
`

const goPackageQuery = `
(package_clause
  (package_identifier) @name)
`

const goImportQuery = `
(import_spec) @import
`

func extractGo(content []byte) ([]store.Symbol, []store.Import, error) {
	lang := tree_sitter.NewLanguage(tree_sitter_go.Language())
	tree, err := parseWith(lang, content)
	if err != nil {
		return nil, nil, err
	}
	defer tree.Close()

	var symbols []store.Symbol
	err = runQuery(lang, tree, content, goSymbolQuery, func(names []string, match *tree_sitter.QueryMatch) {
		nameNode := captureByName(names, match, "name")
		if nameNode == nil {
			return
		}
		var kind string
		var rangeNode *tree_sitter.Node
		for _, capName := range []string{"function", "method", "type"} {
			if n := captureByName(names, match, capName); n != nil {
				kind = capName
				rangeNode = n
				break
			}
		}
		if kind == "" || rangeNode == nil {
			return
		}
		symbols = append(symbols, symbolFromNamedNode(kind, nameNode, rangeNode, content))
	})
	if err != nil {
		return nil, nil, err
	}

	var imports []store.Import
	err = runQuery(lang, tree, content, goImportQuery, func(names []string, match *tree_sitter.QueryMatch) {
		spec := captureByName(names, match, "import")
		if spec == nil {
			return
		}
		pathNode := spec.ChildByFieldName("path")
		if pathNode == nil {
			return
		}
		path := unquoteString(pathNode.Utf8Text(content))
		if path == "" {
			return
		}
		var sym *string
		if nameNode := spec.ChildByFieldName("name"); nameNode != nil && nameNode.Kind() == "package_identifier" {
			sym = ptrStr(nameNode.Utf8Text(content))
		}
		imports = append(imports, store.Import{ImportedPath: path, Symbol: sym, Provenance: store.ImportProvenanceExtracted})
	})
	if err != nil {
		return nil, nil, err
	}

	return symbols, imports, nil
}

func extractGoPackageName(content []byte) string {
	lang := tree_sitter.NewLanguage(tree_sitter_go.Language())
	tree, err := parseWith(lang, content)
	if err != nil {
		return ""
	}
	defer tree.Close()

	var name string
	err = runQuery(lang, tree, content, goPackageQuery, func(names []string, match *tree_sitter.QueryMatch) {
		n := captureByName(names, match, "name")
		if n == nil || name != "" {
			return
		}
		name = n.Utf8Text(content)
	})
	if err != nil {
		return ""
	}
	return name
}
