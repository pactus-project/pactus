package consensusv2

import (
	"fmt"
)

// InvalidJustificationError is returned when the justification for a change-proposer
// vote is invalid.
type InvalidJustificationError struct {
	Reason string
}

func (e InvalidJustificationError) Error() string {
	return fmt.Sprintf("invalid justification: %s", e.Reason)
}

// ConfigError is returned when the config is not valid with a descriptive Reason message.
type ConfigError struct {
	Reason string
}

func (e ConfigError) Error() string {
	return e.Reason
}
