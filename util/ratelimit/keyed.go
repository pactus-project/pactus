package ratelimit

import (
	"sync"
	"time"
)

// KeyedRateLimit applies an independent RateLimit to each key, for example to each peer.
// The table is bounded: when it is full, expired entries are dropped first and,
// if there is still no room, an arbitrary entry is evicted.
type KeyedRateLimit[K comparable] struct {
	lk sync.Mutex

	threshold int
	window    time.Duration
	maxKeys   int
	limits    map[K]*RateLimit
}

// NewKeyedRateLimit creates a KeyedRateLimit that allows threshold requests
// per window for each key, and tracks at most maxKeys keys.
// If the threshold is zero, all requests are allowed.
func NewKeyedRateLimit[K comparable](threshold int, window time.Duration, maxKeys int) *KeyedRateLimit[K] {
	return &KeyedRateLimit[K]{
		threshold: threshold,
		window:    window,
		maxKeys:   max(maxKeys, 1),
		limits:    make(map[K]*RateLimit),
	}
}

// AllowRequest increments the counter of the key and checks if its rate limit is exceeded.
func (k *KeyedRateLimit[K]) AllowRequest(key K) bool {
	if k.threshold == 0 {
		return true
	}

	k.lk.Lock()
	limit, ok := k.limits[key]
	if !ok {
		k.makeRoom()
		limit = NewRateLimit(k.threshold, k.window)
		k.limits[key] = limit
	}
	k.lk.Unlock()

	return limit.AllowRequest()
}

// makeRoom must be called with the lock held.
func (k *KeyedRateLimit[K]) makeRoom() {
	if len(k.limits) < k.maxKeys {
		return
	}

	for key, limit := range k.limits {
		if limit.expired() {
			delete(k.limits, key)
		}
	}

	for key := range k.limits {
		if len(k.limits) < k.maxKeys {
			return
		}
		delete(k.limits, key)
	}
}
