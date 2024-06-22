package ntp

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/pactus-project/pactus/util/logger"
)

const (
	maxClockOffset = time.Duration(math.MinInt64)
)

// QueryError is returned when a query from all NTP pools encounters an error
// and we have no valid response from them.
type QueryError struct{}

func (QueryError) Error() string {
	return "failed to get NTP query from all pools"
}

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
	interval  time.Duration
	threshold time.Duration
	ticker    *time.Ticker
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
		c.ticker = time.NewTicker(interval)
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
func NewNtpChecker(opts ...CheckerOption) *Checker {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defaultInterval := time.Minute
	defaultThreshold := time.Second

	// Initialize the checker with default values.
	checker := &Checker{
		ctx:       ctxWithCancel,
		cancel:    cancel,
		interval:  defaultInterval,
		threshold: defaultThreshold,
		querier:   RemoteQuerier{},
		ticker:    time.NewTicker(defaultInterval),
	}

	// Apply provided options to override default values.
	for _, opt := range opts {
		opt(checker)
	}

	return checker
}

func (c *Checker) Start() {
	for {
		select {
		case <-c.ctx.Done():
			return

		case <-c.ticker.C:
			offset, _ := c.queryClockOffset()

			c.lk.Lock()
			c.offset = offset
			c.lk.Unlock()

			if c.offset != maxClockOffset && c.IsOutOfSync() {
				logger.Error(
					"the system time is out of sync with the network time by more than one second",
					"threshold", c.threshold, "offset", offset)
			}
		}
	}
}

func (c *Checker) Stop() {
	c.cancel()
	c.ticker.Stop()
}

func (c *Checker) IsOutOfSync() bool {
	c.lk.RLock()
	defer c.lk.RUnlock()

	return c.offset.Abs() > c.threshold
}

func (c *Checker) ClockOffset() (time.Duration, error) {
	c.lk.RLock()
	defer c.lk.RUnlock()

	if c.offset == maxClockOffset {
		return 0, QueryError{}
	}

	return c.offset, nil
}

func (c *Checker) queryClockOffset() (time.Duration, error) {
	for _, server := range _pools {
		response, err := c.querier.Query(server)
		if err != nil {
			logger.Warn("ntp query error", "server", server, "error", err)

			continue
		}

		if err := response.Validate(); err != nil {
			logger.Warn("ntp validate error", "server", server, "error", err)

			continue
		}

		logger.Debug("successful ntp query", "offset", response.ClockOffset, "RTT", response.RTT)

		return response.ClockOffset, nil
	}

	logger.Error("failed to get ntp query from all pool")

	return maxClockOffset, QueryError{}
}
