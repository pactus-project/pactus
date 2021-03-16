package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusFlags(t *testing.T) {
	s := NewStatus()

	assert.False(t, s.IsProposed())
	assert.False(t, s.IsPrepared())
	assert.False(t, s.IsPreCommitted())
	assert.False(t, s.IsCommitted())

	s.SetProposed(true)
	s.SetPrepared(true)
	s.SetPreCommitted(true)
	s.SetCommitted(true)

	assert.True(t, s.IsProposed())
	assert.True(t, s.IsPrepared())
	assert.True(t, s.IsPreCommitted())
	assert.True(t, s.IsCommitted())
}
