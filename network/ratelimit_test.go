package network

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimit(t *testing.T) {
	threshold := uint8(5)
	window := 100 * time.Millisecond
	r := newRateLimit(threshold, window)

	t.Run("InitialState", func(t *testing.T) {
		assert.Equal(t, uint8(0), r.counter)
	})

	t.Run("IncrementWithinThreshold", func(t *testing.T) {
		for i := uint8(0); i < threshold; i++ {
			assert.True(t, r.increment())
		}
		assert.Equal(t, threshold, r.counter)
	})

	t.Run("ExceedThreshold", func(t *testing.T) {
		assert.False(t, r.increment())
	})

	t.Run("ResetAfterWindow", func(t *testing.T) {
		time.Sleep(window + 10*time.Millisecond)
		assert.True(t, r.increment())
		assert.Equal(t, uint8(1), r.counter)
	})

	t.Run("ResetMethod", func(t *testing.T) {
		r.reset()
		assert.Equal(t, uint8(0), r.counter)
		assert.True(t, r.increment())
		assert.Equal(t, uint8(1), r.counter)
	})

	t.Run("DiffMethod", func(t *testing.T) {
		assert.LessOrEqual(t, r.diff(), window)
	})
}
