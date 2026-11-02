package ntp

import (
	"context"
	"errors"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/pactus-project/gopkg/logger"
	"github.com/pactus-project/gopkg/scheduler"
)

const (
	maxClockOffset = time.Duration(math.MinInt64)

	// maxSamples is the number of valid responses used to compute the clock offset.
	maxSamples = 3
)

// QueryError is returned when a query from all NTP pools encounters an error
// and we have no valid response from them.
type QueryError struct{}

func (QueryError) Error() string {
	return "failed to get NTP query from all pools"
}

// ErrNotMeasured is returned by ClockOffset until the first measurement is done.
var ErrNotMeasured = errors.New("the clock offset is not measured yet")

var _pools = []string{
	"pool.ntp.org",
	"time.google.com",
	"time.cloudflare.com",
	"time.apple.com",
	"time.windows.com",
	"ntp.ubuntu.com",
}

// Checker represents a NTP checker that periodically checks the system time against the network time.
type Checker struct {
	lk sync.RWMutex

	ctx       context.Context
	cancel    func()
	querier   Querier
	offset    time.Duration
	measured  bool
	interval  time.Duration
	threshold time.Duration
}

// CheckerOption defines the type for functions that configure a Checker.
type CheckerOption func(*Checker)

// WithQuerier sets the Querier for the Checker.
func WithQuerier(querier Querier) CheckerOption {
	return func(c *Checker) {
		c.querier = querier
	}
}

// WithInterval sets the interval at which the checker will run.
// Default interval is 1 minute.
func WithInterval(interval time.Duration) CheckerOption {
	return func(c *Checker) {
		c.interval = interval
	}
}

// WithThreshold sets the threshold for determining if the system time is out of sync.
// Default threshold is 1 second.
func WithThreshold(threshold time.Duration) CheckerOption {
	return func(c *Checker) {
		c.threshold = threshold
	}
}

// NewNtpChecker creates a new Checker with the provided options.
// If no options are provided, it uses default values for interval and threshold.
func NewNtpChecker(ctx context.Context, opts ...CheckerOption) *Checker {
	ctx, cancel := context.WithCancel(ctx)
	defaultInterval := time.Minute
	defaultThreshold := time.Second

	// Initialize the checker with default values.
	checker := &Checker{
		ctx:       ctx,
		cancel:    cancel,
		interval:  defaultInterval,
		threshold: defaultThreshold,
		querier:   RemoteQuerier{},
	}

	// Apply provided options to override default values.
	for _, opt := range opts {
		opt(checker)
	}

	return checker
}

// Start measures the clock offset once per interval.
func (c *Checker) Start() {
	scheduler.Every(c.interval).Do(c.ctx, func(_ context.Context) {
		c.measure()
	})
}

func (c *Checker) measure() {
	offset, _ := c.queryClockOffset()

	c.lk.Lock()
	c.offset = offset
	c.measured = true
	c.lk.Unlock()

	if offset != maxClockOffset && offset.Abs() > c.threshold {
		logger.Error(
			"the system time is out of sync with the network time",
			"threshold", c.threshold, "offset", offset,
		)
	}
}

func (c *Checker) Stop() {
	c.cancel()
}

// IsOutOfSync reports whether the clock is out of sync with the network time.
// It is false until the first measurement is done.
func (c *Checker) IsOutOfSync() bool {
	c.lk.RLock()
	defer c.lk.RUnlock()

	return c.measured && c.offset.Abs() > c.threshold
}

// ClockOffset returns the last measured clock offset. It returns an error
// until the first measurement is done, and when the last one has failed.
func (c *Checker) ClockOffset() (time.Duration, error) {
	c.lk.RLock()
	defer c.lk.RUnlock()

	if !c.measured {
		return 0, ErrNotMeasured
	}

	if c.offset == maxClockOffset {
		return 0, QueryError{}
	}

	return c.offset, nil
}

// queryClockOffset queries the pools in order until it has maxSamples valid
// responses, and returns the median offset, so that a single wrong server
// can't skew the result.
func (c *Checker) queryClockOffset() (time.Duration, error) {
	offsets := make([]time.Duration, 0, maxSamples)
	for _, server := range _pools {
		response, err := c.querier.Query(server)
		if err != nil {
			logger.Debug("ntp query error", "server", server, "error", err)

			continue
		}

		if err := response.Validate(); err != nil {
			logger.Debug("ntp validate error", "server", server, "error", err)

			continue
		}

		logger.Debug("successful ntp query", "offset", response.ClockOffset, "RTT", response.RTT)

		offsets = append(offsets, response.ClockOffset)
		if len(offsets) == maxSamples {
			break
		}
	}

	if len(offsets) == 0 {
		logger.Error("failed to get ntp query from all pool")

		return maxClockOffset, QueryError{}
	}

	return medianOffset(offsets), nil
}

// medianOffset returns the median of offsets, which must not be empty.
func medianOffset(offsets []time.Duration) time.Duration {
	slices.Sort(offsets)

	mid := len(offsets) / 2
	if len(offsets)%2 == 1 {
		return offsets[mid]
	}

	return offsets[mid-1] + (offsets[mid]-offsets[mid-1])/2
}
