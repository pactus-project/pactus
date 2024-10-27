package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	threshold := 5
	window := 100 * time.Millisecond
	rateLimit := NewRateLimit(threshold, window)

	t.Run("InitialState", func(t *testing.T) {
		assert.Equal(t, 0, rateLimit.counter)
	})

	t.Run("AllowRequestWithinThreshold", func(t *testing.T) {
		for i := 0; i < threshold; i++ {
			assert.True(t, rateLimit.AllowRequest())
		}
		assert.Equal(t, threshold, rateLimit.counter)
	})

	t.Run("ExceedThreshold", func(t *testing.T) {
		assert.False(t, rateLimit.AllowRequest())
	})

	t.Run("ResetAfterWindow", func(t *testing.T) {
		time.Sleep(window + 10*time.Millisecond)
		assert.True(t, rateLimit.AllowRequest())
		assert.Equal(t, 1, rateLimit.counter)
	})

	t.Run("ResetMethod", func(t *testing.T) {
		rateLimit.reset()
		assert.Equal(t, 0, rateLimit.counter)
		assert.True(t, rateLimit.AllowRequest())
		assert.Equal(t, 1, rateLimit.counter)
	})

	t.Run("DiffMethod", func(t *testing.T) {
		assert.LessOrEqual(t, rateLimit.diff(), window)
	})
}

func TestRateLimitZeroThreshold(t *testing.T) {
	window := 100 * time.Millisecond
	r := NewRateLimit(0, window)

	assert.True(t, r.AllowRequest())
	assert.Zero(t, r.counter)
}
