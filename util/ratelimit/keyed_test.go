package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyedRateLimit(t *testing.T) {
	limit := NewKeyedRateLimit[string](2, time.Hour, 10)

	assert.True(t, limit.AllowRequest("alice"))
	assert.True(t, limit.AllowRequest("alice"))
	assert.False(t, limit.AllowRequest("alice"))

	// The budget of one key is not shared with the others.
	assert.True(t, limit.AllowRequest("bob"))
}

func TestKeyedRateLimitUnlimited(t *testing.T) {
	limit := NewKeyedRateLimit[string](0, time.Hour, 10)

	for range 100 {
		assert.True(t, limit.AllowRequest("alice"))
	}
	assert.Empty(t, limit.limits)
}

func TestKeyedRateLimitIsBounded(t *testing.T) {
	limit := NewKeyedRateLimit[int](1, time.Hour, 5)

	for key := range 100 {
		assert.True(t, limit.AllowRequest(key), "a new key gets a fresh budget")
		assert.LessOrEqual(t, len(limit.limits), 5)
	}
}

func TestKeyedRateLimitDropsExpiredEntries(t *testing.T) {
	limit := NewKeyedRateLimit[int](1, 10*time.Millisecond, 3)

	for key := range 3 {
		assert.True(t, limit.AllowRequest(key))
	}
	time.Sleep(30 * time.Millisecond)

	assert.True(t, limit.AllowRequest(100))
	assert.Len(t, limit.limits, 1)
}
