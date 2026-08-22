package analyzers

import (
	"github.com/mrchatam/Trace/internal/store"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

const pySymbolQuery = `
(module
  (function_definition
    name: (identifier) @name) @function)

(module
  (class_definition
    name: (identifier) @name) @class)

(class_definition
  body: (block
    (function_definition
      name: (identifier) @name) @method))
`

const pyImportQuery = `
(import_statement
  name: (dotted_name) @path) @import

(import_statement
  name: (aliased_import
    name: (dotted_name) @path)) @import

(import_from_statement
  module_name: [(dotted_name) (relative_import)] @path
  name: (dotted_name) @symbol) @from

(import_from_statement
  module_name: [(dotted_name) (relative_import)] @path
  name: (aliased_import
    name: (dotted_name) @symbol)) @from

(import_from_statement
  module_name: [(dotted_name) (relative_import)] @path
  (wildcard_import)) @from_wild
`

func extractPython(content []byte) ([]store.Symbol, []store.Import, error) {
	lang := tree_sitter.NewLanguage(tree_sitter_python.Language())
	tree, err := parseWith(lang, content)
	if err != nil {
		return nil, nil, err
	}
	defer tree.Close()

	var symbols []store.Symbol
	err = runQuery(lang, tree, content, pySymbolQuery, func(names []string, match *tree_sitter.QueryMatch) {
		nameNode := captureByName(names, match, "name")
		if nameNode == nil {
			return
		}
		var kind string
		var rangeNode *tree_sitter.Node
		for _, capName := range []string{"function", "class", "method"} {
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
	err = runQuery(lang, tree, content, pyImportQuery, func(names []string, match *tree_sitter.QueryMatch) {
		pathNode := captureByName(names, match, "path")
		if pathNode == nil {
			return
		}
		path := pathNode.Utf8Text(content)
		wild := captureByName(names, match, "from_wild") != nil
		prov := store.ImportProvenanceExtracted
		if wild || path == "" {
			prov = store.ImportProvenanceAmbiguous
		}
		symNode := captureByName(names, match, "symbol")
		if symNode != nil {
			imports = append(imports, store.Import{
				ImportedPath: path,
				Symbol:       ptrStr(symNode.Utf8Text(content)),
				Provenance:   prov,
			})
			return
		}
		imports = append(imports, store.Import{ImportedPath: path, Symbol: nil, Provenance: prov})
	})
	if err != nil {
		return nil, nil, err
	}

	return symbols, imports, nil
}
