package analyzers

import (
	"strconv"

	"github.com/mrchatam/Trace/internal/store"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

// File is named extract_javascript.go (not extract_js.go): Go treats *_js.go as
// GOOS=js and would exclude the file on linux/darwin.

// JavaScript grammar: class names are identifier (not type_identifier).
const jsSymbolQuery = `
(function_declaration
  name: (identifier) @name) @function

(generator_function_declaration
  name: (identifier) @name) @function

(class_declaration
  name: (identifier) @name) @class

(method_definition
  name: (property_identifier) @name) @method
`

// TypeScript/TSX grammar: class names are type_identifier.
const tsSymbolQuery = `
(function_declaration
  name: (identifier) @name) @function

(generator_function_declaration
  name: (identifier) @name) @function

(class_declaration
  name: (type_identifier) @name) @class

(method_definition
  name: (property_identifier) @name) @method
`

const jsImportStmtQuery = `(import_statement) @stmt`

func extractJS(content []byte) ([]store.Symbol, []store.Import, error) {
	lang := tree_sitter.NewLanguage(tree_sitter_javascript.Language())
	return extractJSTS(lang, content, jsSymbolQuery)
}

func extractJSTS(lang *tree_sitter.Language, content []byte, symbolQuery string) ([]store.Symbol, []store.Import, error) {
	tree, err := parseWith(lang, content)
	if err != nil {
		return nil, nil, err
	}
	defer tree.Close()

	var symbols []store.Symbol
	err = runQuery(lang, tree, content, symbolQuery, func(names []string, match *tree_sitter.QueryMatch) {
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
		// Skip constructors — not required as methods for minimal graph.
		if kind == "method" && nameNode.Utf8Text(content) == "constructor" {
			return
		}
		symbols = append(symbols, symbolFromNamedNode(kind, nameNode, rangeNode, content))
	})
	if err != nil {
		return nil, nil, err
	}

	var imports []store.Import
	err = runQuery(lang, tree, content, jsImportStmtQuery, func(names []string, match *tree_sitter.QueryMatch) {
		stmt := captureByName(names, match, "stmt")
		if stmt == nil {
			return
		}
		imports = append(imports, extractJSImportStatement(stmt, content)...)
	})
	if err != nil {
		return nil, nil, err
	}

	return symbols, imports, nil
}

func extractJSImportStatement(stmt *tree_sitter.Node, content []byte) []store.Import {
	src := stmt.ChildByFieldName("source")
	if src == nil {
		return nil
	}
	path := unquoteString(src.Utf8Text(content))
	clause := stmt.ChildByFieldName("import") // may be nil for side-effect import
	if clause == nil {
		// Some grammars use anonymous child; fall back to scanning.
		clause = findChildKind(stmt, "import_clause")
	}

	var named []string
	if clause != nil {
		collectNamedImportSymbols(clause, content, &named)
	}
	if len(named) > 0 {
		out := make([]store.Import, 0, len(named))
		for _, sym := range named {
			out = append(out, store.Import{ImportedPath: path, Symbol: ptrStr(sym), Provenance: store.ImportProvenanceExtracted})
		}
		return out
	}
	return []store.Import{{ImportedPath: path, Symbol: nil, Provenance: store.ImportProvenanceExtracted}}
}

func findChildKind(n *tree_sitter.Node, kind string) *tree_sitter.Node {
	for i := uint(0); i < n.ChildCount(); i++ {
		ch := n.Child(i)
		if ch != nil && ch.Kind() == kind {
			return ch
		}
	}
	return nil
}

func collectNamedImportSymbols(n *tree_sitter.Node, content []byte, out *[]string) {
	if n == nil {
		return
	}
	if n.Kind() == "import_specifier" {
		name := n.ChildByFieldName("name")
		if name != nil {
			*out = append(*out, name.Utf8Text(content))
		}
		return
	}
	for i := uint(0); i < n.NamedChildCount(); i++ {
		ch := n.NamedChild(i)
		collectNamedImportSymbols(ch, content, out)
	}
}

const jsExportStmtQuery = `(export_statement) @export`

// jstsExportedKeys is IndexFile-only (API version stays 1). Keys are
// kind+"\n"+name+"\n"+startLine for declarations under export_statement.
func jstsExportedKeys(langID string, content []byte) (map[string]bool, error) {
	lang := jstsLanguage(langID)
	if lang == nil {
		return nil, nil
	}
	tree, err := parseWith(lang, content)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	out := map[string]bool{}
	err = runQuery(lang, tree, content, jsExportStmtQuery, func(names []string, match *tree_sitter.QueryMatch) {
		stmt := captureByName(names, match, "export")
		if stmt == nil {
			return
		}
		collectJSTSExportDecls(stmt, content, out)
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func collectJSTSExportDecls(exportStmt *tree_sitter.Node, content []byte, out map[string]bool) {
	if exportStmt == nil {
		return
	}
	for i := uint(0); i < exportStmt.NamedChildCount(); i++ {
		ch := exportStmt.NamedChild(i)
		if ch == nil {
			continue
		}
		switch ch.Kind() {
		case "function_declaration", "generator_function_declaration":
			recordJSTSExport(ch, "function", content, out)
		case "class_declaration":
			recordJSTSExport(ch, "class", content, out)
			collectJSTSExportedMethods(ch, content, out)
		}
	}
}

func collectJSTSExportedMethods(classDecl *tree_sitter.Node, content []byte, out map[string]bool) {
	body := classDecl.ChildByFieldName("body")
	if body == nil {
		return
	}
	for i := uint(0); i < body.NamedChildCount(); i++ {
		ch := body.NamedChild(i)
		if ch == nil || ch.Kind() != "method_definition" {
			continue
		}
		nameNode := ch.ChildByFieldName("name")
		if nameNode == nil || nameNode.Utf8Text(content) == "constructor" {
			continue
		}
		recordNamedExport(nameNode, ch, "method", content, out)
	}
}

func recordJSTSExport(decl *tree_sitter.Node, kind string, content []byte, out map[string]bool) {
	nameNode := decl.ChildByFieldName("name")
	if nameNode == nil {
		return
	}
	recordNamedExport(nameNode, decl, kind, content, out)
}

func recordNamedExport(nameNode, rangeNode *tree_sitter.Node, kind string, content []byte, out map[string]bool) {
	name := nameNode.Utf8Text(content)
	if name == "" {
		return
	}
	start, _ := nodeLines(rangeNode)
	out[exportKey(kind, name, start)] = true
}

func exportKey(kind, name string, startLine int) string {
	return kind + "\n" + name + "\n" + strconv.Itoa(startLine)
}
