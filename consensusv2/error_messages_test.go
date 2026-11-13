package consensusv2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMessages(t *testing.T) {
	assert.Equal(t, "invalid justification: bad proof", InvalidJustificationError{Reason: "bad proof"}.Error())
	assert.Equal(t, "bad config", ConfigError{Reason: "bad config"}.Error())
}
