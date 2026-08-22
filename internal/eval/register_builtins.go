package eval

import (
	"context"
	"time"

	"github.com/mrchatam/Trace/internal/eval/mechs"
)

func init() {
	for _, b := range mechs.All() {
		Register(builtinAdapter{inner: b})
	}
}

type builtinAdapter struct {
	inner mechs.Builtin
}

func (a builtinAdapter) ID() string { return a.inner.ID() }

func (a builtinAdapter) Run(ctx context.Context, in EvalInput) (EvalResult, error) {
	if in.Rules != nil {
		ctx = WithRules(ctx, in.Rules)
	}
	passed, summary, details, err := a.inner.Run(ctx, in.TaskID, in.Service)
	if err != nil {
		return EvalResult{}, err
	}
	return EvalResult{
		MechanismID: a.inner.ID(),
		Passed:      passed,
		Summary:     summary,
		DetailsJSON: details,
		RecordedAt:  time.Now().UTC().Format(time.RFC3339),
	}, nil
}
