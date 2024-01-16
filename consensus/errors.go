package consensus

import (
	"fmt"

	"github.com/pactus-project/pactus/types/vote"
)

// invalidJustificationError is returned when the justification for a change-proposer
// vote is invalid.
type invalidJustificationError struct {
	JustType vote.JustType
	Reason   string
}

func (e invalidJustificationError) Error() string {
	return fmt.Sprintf("invalid justification: %s, reason: %s",
		e.JustType.String(), e.Reason)
}

// ConfigError is returned when the config is not valid with a descriptive Reason message.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}
