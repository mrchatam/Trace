package x0

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// QueryBank is the Gate C understanding question set + GT keys.
type QueryBank struct {
	SchemaVersion int     `json:"schema_version"`
	TaskFamily    string  `json:"task_family"`
	Fixture       string  `json:"fixture"`
	Queries       []Query `json:"queries"`
}

// Query is one agent-facing understanding question with grading keys.
type Query struct {
	ID                string          `json:"id"`
	Theme             string          `json:"theme"`
	Question          string          `json:"question"`
	GT                json.RawMessage `json:"gt"`
	MustIncludeAny    []string        `json:"must_include_any"`
	MustIncludeAll    []string        `json:"must_include_all"`
	MustIncludeAllAlt []string        `json:"must_include_all_alt"`
	CriticalMissIfAny []string        `json:"critical_miss_if_any"`
}

// AnswerPack is one recorded (or live) agent run's answers.
type AnswerPack struct {
	Condition string       `json:"condition"`
	RunID     string       `json:"run_id"`
	Model     string       `json:"model"`
	ToolsUsed []string     `json:"tools_used"`
	Answers   []PackAnswer `json:"answers"`
	Notes     string       `json:"notes,omitempty"`
	LatencyMS *float64     `json:"latency_ms,omitempty"`
	Tokens    *float64     `json:"tokens,omitempty"`
}

// PackAnswer is free-text plus structured asserts used for fair grading.
// Grading prefers Assert tokens (entity ids, rels, facts) over raw Text substring
// so negations like "no IN_PROGRESS" do not falsely score correct.
type PackAnswer struct {
	ID     string   `json:"id"`
	Text   string   `json:"text"`
	Assert []string `json:"assert"`
}

// Grade is per-query scoring.
type Grade string

const (
	GradeCorrect      Grade = "correct"
	GradeIncorrect    Grade = "incorrect"
	GradeCriticalMiss Grade = "critical_miss"
)

// PerQueryGrade is one graded query row for the quality object.
type PerQueryGrade struct {
	ID    string `json:"id"`
	Grade Grade  `json:"grade"`
}

// UnderstandingQuality is the schema v1 quality object for understanding runs.
type UnderstandingQuality struct {
	UnderstandingAccuracy float64         `json:"understanding_accuracy"`
	Correct               int             `json:"correct"`
	Total                 int             `json:"total"`
	CriticalMisses        int             `json:"critical_misses"`
	PerQuery              []PerQueryGrade `json:"per_query"`
}

// LoadQueryBank reads queries.json from path.
func LoadQueryBank(path string) (*QueryBank, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var bank QueryBank
	if err := json.Unmarshal(b, &bank); err != nil {
		return nil, err
	}
	if len(bank.Queries) < 5 {
		return nil, fmt.Errorf("query bank needs ≥5 queries, got %d", len(bank.Queries))
	}
	return &bank, nil
}

// LoadAnswerPack reads a recorded answer pack JSON.
func LoadAnswerPack(path string) (*AnswerPack, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pack AnswerPack
	if err := json.Unmarshal(b, &pack); err != nil {
		return nil, err
	}
	if pack.RunID == "" {
		return nil, fmt.Errorf("answer pack missing run_id: %s", path)
	}
	return &pack, nil
}

func findAnswer(pack *AnswerPack, queryID string) PackAnswer {
	for _, a := range pack.Answers {
		if a.ID == queryID {
			return a
		}
	}
	return PackAnswer{}
}

func norm(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func tokenMatch(haystack []string, needle string) bool {
	n := norm(needle)
	if n == "" {
		return true
	}
	for _, h := range haystack {
		hn := norm(h)
		// Assert token equals needle, or assert token contains needle.
		// Do NOT treat "needle contains assert token" as a hit — that false-positives
		// when a UUID assert matches a longer critical_miss phrase that mentions it.
		if hn == n || strings.Contains(hn, n) {
			return true
		}
	}
	return false
}

func anyToken(haystack []string, needles []string) bool {
	for _, n := range needles {
		if tokenMatch(haystack, n) {
			return true
		}
	}
	return false
}

func allTokens(haystack []string, needles []string) bool {
	if len(needles) == 0 {
		return true
	}
	for _, n := range needles {
		if !tokenMatch(haystack, n) {
			return false
		}
	}
	return true
}

// GradeAnswer grades structured asserts (preferred) against query GT keys.
// Empty assert ⇒ incorrect (unknown / refused), unless critical tokens appear in assert.
func GradeAnswer(q Query, ans PackAnswer) Grade {
	assert := ans.Assert
	if len(assert) == 0 {
		return GradeIncorrect
	}
	if anyToken(assert, q.CriticalMissIfAny) {
		return GradeCriticalMiss
	}
	primaryOK := anyToken(assert, q.MustIncludeAny) && allTokens(assert, q.MustIncludeAll)
	altOK := false
	if len(q.MustIncludeAllAlt) > 0 {
		altOK = anyToken(assert, q.MustIncludeAny) && allTokens(assert, q.MustIncludeAllAlt)
	}
	if primaryOK || altOK {
		return GradeCorrect
	}
	return GradeIncorrect
}

// GradePack grades all queries in bank order against the pack.
func GradePack(bank *QueryBank, pack *AnswerPack) UnderstandingQuality {
	out := UnderstandingQuality{
		Total:    len(bank.Queries),
		PerQuery: make([]PerQueryGrade, 0, len(bank.Queries)),
	}
	for _, q := range bank.Queries {
		g := GradeAnswer(q, findAnswer(pack, q.ID))
		out.PerQuery = append(out.PerQuery, PerQueryGrade{ID: q.ID, Grade: g})
		if g == GradeCorrect {
			out.Correct++
		}
		if g == GradeCriticalMiss {
			out.CriticalMisses++
		}
	}
	if out.Total > 0 {
		out.UnderstandingAccuracy = float64(out.Correct) / float64(out.Total)
	}
	return out
}

// MeanAccuracy returns mean understanding_accuracy across qualities.
func MeanAccuracy(qs []UnderstandingQuality) float64 {
	if len(qs) == 0 {
		return 0
	}
	var sum float64
	for _, q := range qs {
		sum += q.UnderstandingAccuracy
	}
	return sum / float64(len(qs))
}

// MeanCriticalMisses returns mean critical_misses across qualities.
func MeanCriticalMisses(qs []UnderstandingQuality) float64 {
	if len(qs) == 0 {
		return 0
	}
	var sum float64
	for _, q := range qs {
		sum += float64(q.CriticalMisses)
	}
	return sum / float64(len(qs))
}

// QualityJSON marshals UnderstandingQuality for MetricsRun.Quality.
func QualityJSON(q UnderstandingQuality) (json.RawMessage, error) {
	return json.Marshal(q)
}
