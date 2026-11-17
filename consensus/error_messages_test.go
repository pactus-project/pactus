package consensus

import (
	"testing"

	"github.com/pactus-project/pactus/types/vote"
	"github.com/stretchr/testify/assert"
)

func TestErrorMessages(t *testing.T) {
	assert.Equal(t, "invalid justification: JustInitYes, reason: bad proof",
		InvalidJustificationError{JustType: vote.JustTypeInitYes, Reason: "bad proof"}.Error())
	assert.Equal(t, "bad config", ConfigError{Reason: "bad config"}.Error())
}
