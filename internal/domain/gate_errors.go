package domain

import "fmt"

// PrematureImplementation is returned when a harness gate blocks action because
// deliberation policy has not reached the required phase.
type PrematureImplementation struct {
	TaskID           string
	For              string
	RecommendedPhase string
	ReasonCode       string
	Message          string
}

func (e *PrematureImplementation) Error() string {
	if e.Message != "" {
		return e.Message
	}
	msg := fmt.Sprintf("premature implementation: task %s blocked for %s", e.TaskID, e.For)
	if e.RecommendedPhase != "" {
		msg += fmt.Sprintf(" (recommended phase %s", e.RecommendedPhase)
		if e.ReasonCode != "" {
			msg += fmt.Sprintf(", %s", e.ReasonCode)
		}
		msg += ")"
	} else if e.ReasonCode != "" {
		msg += fmt.Sprintf(" (%s)", e.ReasonCode)
	}
	return msg
}

// Code returns the stable machine-readable violation code.
func (e *PrematureImplementation) Code() string {
	return "premature_implementation"
}
