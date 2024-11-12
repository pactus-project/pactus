package sync

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigCheck(t *testing.T) {
	c := DefaultConfig()
	assert.NoError(t, c.BasicCheck())
	assert.Equal(t, c.SessionTimeout(), 10*time.Second)
}
