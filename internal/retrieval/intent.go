package retrieval

import (
	"regexp"
	"strings"
	"unicode"
)

const (
	maxIntentKeywords     = 32
	maxIntentEntityHints  = 16
	maxIntentFTSRunes     = 512
	intentSummaryMaxRunes = 256
)

var (
	uuidRE   = regexp.MustCompile(`(?i)\b[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\b`)
	pathRE   = regexp.MustCompile(`(?:[\w.-]+/)+[\w.-]+\.(?:go|ts|tsx|js|jsx|py|md|json|yaml|yml|toml)\b`)
	symbolRE = regexp.MustCompile(`\b[A-Z][a-zA-Z0-9_]{2,}\b`)
)

var intentStopWords = map[string]struct{}{
	"a": {}, "an": {}, "and": {}, "are": {}, "as": {}, "at": {}, "be": {}, "by": {},
	"for": {}, "from": {}, "in": {}, "is": {}, "it": {}, "of": {}, "on": {}, "or": {},
	"that": {}, "the": {}, "this": {}, "to": {}, "was": {}, "with": {},
}

// IntentInput carries deterministic sources for G9 intent extraction.
type IntentInput struct {
	TaskTitle string
	TaskBody  string
	Query     string
}

// EntityHint is a structured hint extracted from task/query text (UUID, path, symbol).
type EntityHint struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

// Intent is the bounded structured output of ExtractIntent.
type Intent struct {
	Keywords    []string     `json:"keywords"`
	EntityHints []EntityHint `json:"entity_hints,omitempty"`
	Scope       string       `json:"scope,omitempty"`
	Source      string       `json:"source,omitempty"`
}

// ExtractIntent performs bounded rule-based keyword and entity extraction.
//
// G9 intent precedes retrieval channels: structured keyword/entity extraction from task+query
// feeds FTS query building. G1 (compiler ContextOptions.Query) merges query FTS hits into the
// packet after Expand — complementary, not duplicate. Intent never replaces task_id moat.
func ExtractIntent(in IntentInput) Intent {
	var sources []string
	seenKw := map[string]struct{}{}
	var keywords []string

	addTokens := func(text, source string) {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}
		if source != "" {
			found := false
			for _, s := range sources {
				if s == source {
					found = true
					break
				}
			}
			if !found {
				sources = append(sources, source)
			}
		}
		for _, tok := range tokenizeIntent(text) {
			if len(keywords) >= maxIntentKeywords {
				return
			}
			if _, ok := seenKw[tok]; ok {
				continue
			}
			seenKw[tok] = struct{}{}
			keywords = append(keywords, tok)
		}
	}

	addTokens(in.TaskTitle, "task")
	addTokens(in.TaskBody, "task")
	addTokens(in.Query, "query")

	combined := strings.Join([]string{in.TaskTitle, in.TaskBody, in.Query}, "\n")
	hints := extractEntityHints(combined)

	scope := intentScope(in)
	source := strings.Join(sources, ",")

	return Intent{
		Keywords:    keywords,
		EntityHints: hints,
		Scope:       scope,
		Source:      source,
	}
}

func intentScope(in IntentInput) string {
	hasTask := strings.TrimSpace(in.TaskTitle) != "" || strings.TrimSpace(in.TaskBody) != ""
	hasQuery := strings.TrimSpace(in.Query) != ""
	switch {
	case hasTask && hasQuery:
		return "task+query"
	case hasTask:
		return "task"
	case hasQuery:
		return "query"
	default:
		return ""
	}
}

func tokenizeIntent(text string) []string {
	var out []string
	var b strings.Builder
	flush := func() {
		if b.Len() == 0 {
			return
		}
		tok := strings.ToLower(b.String())
		b.Reset()
		if len(tok) < 2 {
			return
		}
		if _, stop := intentStopWords[tok]; stop {
			return
		}
		out = append(out, tok)
	}
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		flush()
	}
	flush()
	return out
}

func extractEntityHints(text string) []EntityHint {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	seen := map[string]struct{}{}
	var hints []EntityHint
	add := func(kind, value string) {
		if value == "" || len(hints) >= maxIntentEntityHints {
			return
		}
		key := kind + "\x00" + strings.ToLower(value)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		hints = append(hints, EntityHint{Kind: kind, Value: value})
	}
	for _, m := range uuidRE.FindAllString(text, -1) {
		add("uuid", strings.ToLower(m))
	}
	for _, m := range pathRE.FindAllString(text, -1) {
		add("path", m)
	}
	for _, m := range symbolRE.FindAllString(text, -1) {
		add("symbol", m)
	}
	return hints
}

// FTSQuery returns a bounded lexical query string for SQLite FTS5 MATCH.
func (i Intent) FTSQuery() string {
	parts := make([]string, 0, len(i.Keywords)+len(i.EntityHints))
	for _, kw := range i.Keywords {
		parts = append(parts, kw)
	}
	for _, h := range i.EntityHints {
		switch h.Kind {
		case "uuid", "path":
			parts = append(parts, h.Value)
		case "symbol":
			parts = append(parts, h.Value)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	q := strings.Join(parts, " ")
	runes := []rune(q)
	if len(runes) > maxIntentFTSRunes {
		q = string(runes[:maxIntentFTSRunes])
	}
	return strings.TrimSpace(q)
}

// SummaryKeywords returns a compact capped keyword string for packet metadata.
func (i Intent) SummaryKeywords() string {
	if len(i.Keywords) == 0 {
		return ""
	}
	s := strings.Join(i.Keywords, " ")
	runes := []rune(s)
	if len(runes) > intentSummaryMaxRunes {
		s = string(runes[:intentSummaryMaxRunes])
	}
	return strings.TrimSpace(s)
}
