package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	threshold := 5
	window := 100 * time.Millisecond
	r := NewRateLimit(threshold, window)

	t.Run("InitialState", func(t *testing.T) {
		assert.Equal(t, 0, r.counter)
	})

	t.Run("AllowRequestWithinThreshold", func(t *testing.T) {
		for i := 0; i < threshold; i++ {
			assert.True(t, r.AllowRequest())
		}
		assert.Equal(t, threshold, r.counter)
	})

	t.Run("ExceedThreshold", func(t *testing.T) {
		assert.False(t, r.AllowRequest())
	})

	t.Run("ResetAfterWindow", func(t *testing.T) {
		time.Sleep(window + 10*time.Millisecond)
		assert.True(t, r.AllowRequest())
		assert.Equal(t, 1, r.counter)
	})

	t.Run("ResetMethod", func(t *testing.T) {
		r.reset()
		assert.Equal(t, 0, r.counter)
		assert.True(t, r.AllowRequest())
		assert.Equal(t, 1, r.counter)
	})

	t.Run("DiffMethod", func(t *testing.T) {
		assert.LessOrEqual(t, r.diff(), window)
	})
}

func TestRateLimitZeroThreshold(t *testing.T) {
	window := 100 * time.Millisecond
	r := NewRateLimit(0, window)

	assert.True(t, r.AllowRequest())
	assert.Zero(t, r.counter)
}
