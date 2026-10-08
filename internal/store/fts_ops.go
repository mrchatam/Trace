package store

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// opsTokKind classifies one token of a boolean FTS query.
type opsTokKind int

const (
	opsLiteral opsTokKind = iota
	opsPhrase
	opsAnd
	opsOr
	opsNot
	opsNear
)

// opsTok is one token of a boolean FTS query.
type opsTok struct {
	kind   opsTokKind
	text   string // literal or phrase content, already sanitized
	prefix bool   // trailing '*' (prefix match)
	nearN  int    // NEAR(n) distance; 0 means bare NEAR
}

// nearOpsRe matches NEAR(n). NEAR is recognized only in the exact spellings
// below so a query for the word "near" stays a literal.
var nearOpsRe = regexp.MustCompile(`^(?:NEAR|Near|near)\((\d{1,3})\)$`)

// tokenizeOpsQuery splits q into an operator-aware token stream.
//
// Operator keywords are recognized only in uppercase (AND/OR/NOT/NEAR), so
// "and" in prose stays a search term. Double quotes delimit a phrase and a
// trailing '*' marks a prefix. Everything else is sanitized with the same
// letter/number charset the AND path uses, so punctuation can never reach the
// MATCH expression as syntax.
func tokenizeOpsQuery(q string) []opsTok {
	runes := []rune(q)
	var toks []opsTok
	var raw strings.Builder

	flush := func() {
		r := raw.String()
		raw.Reset()
		if r == "" {
			return
		}
		switch r {
		case "AND":
			toks = append(toks, opsTok{kind: opsAnd})
			return
		case "OR":
			toks = append(toks, opsTok{kind: opsOr})
			return
		case "NOT":
			toks = append(toks, opsTok{kind: opsNot})
			return
		case "NEAR":
			toks = append(toks, opsTok{kind: opsNear})
			return
		}
		if m := nearOpsRe.FindStringSubmatch(r); m != nil {
			n, _ := strconv.Atoi(m[1])
			toks = append(toks, opsTok{kind: opsNear, nearN: n})
			return
		}
		text := strings.Join(ftsTokens(r), " ")
		if text == "" {
			return
		}
		toks = append(toks, opsTok{kind: opsLiteral, text: text, prefix: strings.HasSuffix(r, "*")})
	}

	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case r == '"':
			flush()
			j := i + 1
			var phrase strings.Builder
			for j < len(runes) && runes[j] != '"' {
				if runes[j] == '\\' && j+1 < len(runes) {
					j++
				}
				phrase.WriteRune(runes[j])
				j++
			}
			if j < len(runes) {
				i = j // consume the closing quote
			} else {
				i = len(runes)
			}
			if text := strings.Join(ftsTokens(phrase.String()), " "); text != "" {
				toks = append(toks, opsTok{kind: opsPhrase, text: text})
			}
		case unicode.IsSpace(r):
			flush()
		default:
			raw.WriteRune(r)
		}
	}
	flush()
	return toks
}

// opsOperands returns the literal/phrase terms of a token stream, in order.
func opsOperands(toks []opsTok) []string {
	var out []string
	for _, t := range toks {
		if t.kind == opsLiteral || t.kind == opsPhrase {
			out = append(out, t.text)
		}
	}
	return out
}

// buildOpsMatch folds a token stream into a fully parenthesized FTS5
// expression. Every binary merge is parenthesized, so the result does not
// depend on FTS5 operator precedence.
//
// ok is false for a malformed stream (leading or doubled binary operator, a
// binary operator with no right operand, a doubled NOT, an out-of-range NEAR
// distance, or no operands at all) so the caller can fall back to MatchAnd
// instead of emitting something SQLite would reject at query time.
func buildOpsMatch(toks []opsTok) (string, bool) {
	var expr, infix string
	negate := false

	for _, t := range toks {
		switch t.kind {
		case opsAnd, opsOr:
			if expr == "" || infix != "" {
				return "", false
			}
			if t.kind == opsOr {
				infix = "OR"
			} else {
				infix = "AND"
			}
			continue
		case opsNear:
			if expr == "" || infix != "" {
				return "", false
			}
			if t.nearN == 0 {
				infix = "NEAR"
			} else if t.nearN <= maxNearDistance {
				infix = "NEAR(" + strconv.Itoa(t.nearN) + ")"
			} else {
				return "", false
			}
			continue
		case opsNot:
			// A NOT must introduce the operand that follows it; after a
			// complete expression it would be a syntax error in FTS5.
			if expr != "" || negate {
				return "", false
			}
			negate = true
			continue
		}

		operand := quoteFTS(t.text)
		if t.prefix {
			operand += "*"
		}
		if negate {
			operand = "NOT (" + operand + ")"
			negate = false
		}
		if expr == "" {
			expr = operand
			continue
		}
		if infix == "" {
			return "", false
		}
		expr = "(" + expr + " " + infix + " " + operand + ")"
		infix = ""
	}

	if negate || infix != "" || expr == "" {
		return "", false
	}
	return expr, true
}
