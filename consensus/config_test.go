package consensus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.SanityCheck())

	c.DeltaDuration = -1 * time.Second
	assert.Error(t, c.SanityCheck())
	c.TimeoutPrecommit = -1 * time.Second
	assert.Error(t, c.SanityCheck())
	c.TimeoutPrepare = -1 * time.Second
	assert.Error(t, c.SanityCheck())
}

func TestTimeoutWithDelta(t *testing.T) {
	c := DefaultConfig()

	assert.Equal(t, c.PrepareTimeout(0), c.TimeoutPrepare)
	assert.Equal(t, c.PrepareTimeout(1), c.TimeoutPrepare+c.DeltaDuration)

	assert.Equal(t, c.PrecommitTimeout(0), c.TimeoutPrecommit)
	assert.Equal(t, c.PrecommitTimeout(1), c.TimeoutPrecommit+c.DeltaDuration)
}
