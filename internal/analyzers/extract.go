package analyzers

import (
	"fmt"
	"strings"

	"github.com/mrchatam/Trace/internal/store"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func nodeLines(n *tree_sitter.Node) (start, end int) {
	start = int(n.StartPosition().Row) + 1
	endRow := n.EndPosition().Row
	endCol := n.EndPosition().Column
	// EndPosition is exclusive; a column-0 end means the node ended at EOL of previous row.
	if endCol == 0 && endRow > n.StartPosition().Row {
		endRow--
	}
	end = int(endRow) + 1
	if end < start {
		end = start
	}
	return start, end
}

func unquoteString(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		switch s[0] {
		case '"', '\'', '`':
			if s[len(s)-1] == s[0] {
				return s[1 : len(s)-1]
			}
		}
	}
	return s
}

func parseWith(lang *tree_sitter.Language, content []byte) (*tree_sitter.Tree, error) {
	parser := tree_sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(lang); err != nil {
		return nil, fmt.Errorf("set language: %w", err)
	}
	tree := parser.Parse(content, nil)
	if tree == nil {
		return nil, fmt.Errorf("parse returned nil tree")
	}
	return tree, nil
}

// runQuery executes a tree-sitter query and invokes fn for each match.
// Capture names are available via query.CaptureNames()[capture.Index].
func runQuery(lang *tree_sitter.Language, tree *tree_sitter.Tree, content []byte, querySrc string, fn func(captureNames []string, match *tree_sitter.QueryMatch)) error {
	query, qerr := tree_sitter.NewQuery(lang, querySrc)
	if qerr != nil {
		return fmt.Errorf("query: %v", qerr)
	}
	defer query.Close()

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	names := query.CaptureNames()
	matches := cursor.Matches(query, tree.RootNode(), content)
	for match := matches.Next(); match != nil; match = matches.Next() {
		fn(names, match)
	}
	return nil
}

func captureByName(names []string, match *tree_sitter.QueryMatch, want string) *tree_sitter.Node {
	for _, cap := range match.Captures {
		if int(cap.Index) < len(names) && names[cap.Index] == want {
			n := cap.Node
			return &n
		}
	}
	return nil
}

func symbolFromNamedNode(kind string, nameNode, rangeNode *tree_sitter.Node, content []byte) store.Symbol {
	name := nameNode.Utf8Text(content)
	start, end := nodeLines(rangeNode)
	return store.Symbol{
		Name:      name,
		Kind:      kind,
		StartLine: start,
		EndLine:   end,
	}
}

func ptrStr(s string) *string {
	if s == "" {
		return nil
	}
	v := s
	return &v
}
