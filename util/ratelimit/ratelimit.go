package ratelimit

import (
	"sync"
	"time"
)

type RateLimit struct {
	lk sync.RWMutex

	referenceTime time.Time
	threshold     int
	counter       int
	window        time.Duration
}

// NewRateLimit initializes a new RateLimit instance with the given threshold and window duration.
func NewRateLimit(threshold int, window time.Duration) *RateLimit {
	return &RateLimit{
		referenceTime: time.Now(),
		threshold:     threshold,
		counter:       0,
		window:        window,
	}
}

func (r *RateLimit) diff() time.Duration {
	return time.Since(r.referenceTime)
}

// expired reports whether the window of the limiter is over.
func (r *RateLimit) expired() bool {
	r.lk.RLock()
	defer r.lk.RUnlock()

	return r.diff() > r.window
}

func (r *RateLimit) reset() {
	r.counter = 0
	r.referenceTime = time.Now()
}

// Exceeded reports whether the threshold is reached in the current window,
// without counting a request.
func (r *RateLimit) Exceeded() bool {
	r.lk.RLock()
	defer r.lk.RUnlock()

	if r.threshold == 0 || r.diff() > r.window {
		return false
	}

	return r.counter >= r.threshold
}

// AllowRequest increments the counter and checks if the rate limit is exceeded.
// If the threshold is zero, it allows all requests.
func (r *RateLimit) AllowRequest() bool {
	r.lk.Lock()
	defer r.lk.Unlock()

	// If the threshold is zero, allow all requests
	if r.threshold == 0 {
		return true
	}

	// Check if the window has expired and reset if necessary
	if r.diff() > r.window {
		r.reset()
	}

	r.counter++

	// Check if the threshold is exceeded
	return r.counter <= r.threshold
}
