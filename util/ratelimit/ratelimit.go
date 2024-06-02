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

func (r *RateLimit) reset() {
	r.counter = 0
	r.referenceTime = time.Now()
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
