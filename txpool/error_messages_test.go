package txpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorMessages(t *testing.T) {
	assert.Equal(t, "bad config", ConfigError{Reason: "bad config"}.Error())
	assert.Equal(t, "transaction fee is below the minimum of 0.01 PAC", InvalidFeeError{MinimumFee: 1e7}.Error())
}
