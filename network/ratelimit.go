package network

import (
	"time"
)

type rateLimit struct {
	referenceTime time.Time
	threshold     int
	counter       int
	window        time.Duration
}

func newRateLimit(threshold int, window time.Duration) *rateLimit {
	return &rateLimit{
		referenceTime: time.Now(),
		threshold:     threshold,
		counter:       0,
		window:        window,
	}
}

func (r *rateLimit) diff() time.Duration {
	return time.Since(r.referenceTime)
}

func (r *rateLimit) reset() {
	r.counter = 0
	r.referenceTime = time.Now()
}

func (r *rateLimit) increment() bool {

	// Check if the window has expired and reset if necessary
	if r.diff() > r.window {
		r.reset()
	}

	r.counter++

	// Check if the threshold is exceeded
	return r.counter <= r.threshold
}
