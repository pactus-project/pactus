package network

import (
	"sync"
	"time"
)

type rateLimit struct {
	referenceTime time.Time
	threshold     uint8
	counter       uint8
	mu            sync.Mutex
	window        time.Duration
}

func newRateLimit(threshold uint8, window time.Duration) *rateLimit {
	return &rateLimit{
		referenceTime: time.Now(),
		threshold:     threshold,
		counter:       0,
		mu:            sync.Mutex{},
		window:        window,
	}
}

func (r *rateLimit) diff() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	return time.Since(r.referenceTime)
}

func (r *rateLimit) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter = 0
	r.referenceTime = time.Now()
}

func (r *rateLimit) increment() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the window has expired and reset if necessary
	if time.Since(r.referenceTime) > r.window {
		r.counter = 0
		r.referenceTime = time.Now()
	}

	r.counter++

	// Check if the threshold is exceeded
	if r.counter > r.threshold {
		return false
	}

	return true
}
