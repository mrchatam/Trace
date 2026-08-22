package analyzers

import (
	"strconv"

	"github.com/mrchatam/Trace/internal/store"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

// Call-expression query for Jest/Vitest-style tests. Declaration queries
// (jsSymbolQuery / tsSymbolQuery) cannot see test()/it()/describe().
const jsTestCallQuery = `
(call_expression
  function: (identifier) @callee) @call

(call_expression
  function: (member_expression
    object: (identifier) @obj
    property: (property_identifier) @prop)) @call
`

var jstsTestCallees = map[string]bool{
	"test": true, "it": true, "xtest": true, "fit": true, "describe": true,
}

var jstsTestMemberCallees = map[string]bool{
	"test.only": true, "test.skip": true,
	"it.only": true, "it.skip": true,
	"describe.only": true, "describe.skip": true,
}

func jstsLanguage(langID string) *tree_sitter.Language {
	switch langID {
	case LangJavaScript:
		return tree_sitter.NewLanguage(tree_sitter_javascript.Language())
	case LangTypeScript:
		return tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript())
	case LangTSX:
		return tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTSX())
	default:
		return nil
	}
}

// extractJSTSTestCalls is gated by the JS/TS test-file path; do not call from
// LanguageAdapter.Extract (API version stays 1).
func extractJSTSTestCalls(langID string, content []byte) ([]store.Symbol, error) {
	lang := jstsLanguage(langID)
	if lang == nil {
		return nil, nil
	}
	tree, err := parseWith(lang, content)
	if err != nil {
		return nil, err
	}
	defer tree.Close()

	var symbols []store.Symbol
	seen := map[string]bool{}
	err = runQuery(lang, tree, content, jsTestCallQuery, func(names []string, match *tree_sitter.QueryMatch) {
		call := captureByName(names, match, "call")
		if call == nil {
			return
		}
		callee := ""
		if n := captureByName(names, match, "callee"); n != nil {
			callee = n.Utf8Text(content)
			if !jstsTestCallees[callee] {
				return
			}
		} else {
			obj := captureByName(names, match, "obj")
			prop := captureByName(names, match, "prop")
			if obj == nil || prop == nil {
				return
			}
			callee = obj.Utf8Text(content) + "." + prop.Utf8Text(content)
			if !jstsTestMemberCallees[callee] {
				return
			}
		}
		name := firstCallStringArg(call, content)
		if name == "" {
			name = callee
		}
		start, end := nodeLines(call)
		key := name + "\n" + callee + "\n" + strconv.Itoa(start)
		if seen[key] {
			return
		}
		seen[key] = true
		symbols = append(symbols, store.Symbol{
			Name:      name,
			Kind:      "test",
			StartLine: start,
			EndLine:   end,
		})
	})
	if err != nil {
		return nil, err
	}
	return symbols, nil
}

func firstCallStringArg(call *tree_sitter.Node, content []byte) string {
	args := call.ChildByFieldName("arguments")
	if args == nil {
		return ""
	}
	for i := uint(0); i < args.NamedChildCount(); i++ {
		ch := args.NamedChild(i)
		if ch == nil {
			continue
		}
		switch ch.Kind() {
		case "string", "template_string", "string_fragment":
			if s := unquoteString(ch.Utf8Text(content)); s != "" {
				return s
			}
		}
	}
	return ""
}
