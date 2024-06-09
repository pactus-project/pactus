package ntp

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/beevik/ntp"
	"github.com/pactus-project/pactus/util/errors"
	"github.com/pactus-project/pactus/util/logger"
)

const (
	maxClockOffset = time.Duration(math.MinInt64)
)

var _pools = []string{
	"pool.ntp.org",
	"time.google.com",
	"time.cloudflare.com",
	"time.apple.com",
	"time.windows.com",
	"ntp.ubuntu.com",
}

type Checker struct {
	lock   sync.RWMutex
	ctx    context.Context
	cancel func()

	ticker    *time.Ticker
	threshold time.Duration
	offset    time.Duration
	interval  time.Duration
}

func NewNtpChecker(interval, threshold time.Duration) *Checker {
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	server := &Checker{
		ctx:       ctxWithCancel,
		cancel:    cancel,
		interval:  interval,
		threshold: threshold,
		ticker:    time.NewTicker(interval),
	}

	return server
}

func (c *Checker) Start() {
	for {
		select {
		case <-c.ctx.Done():
			return

		case <-c.ticker.C:
			offset := c.clockOffset()
			c.lock.Lock()
			c.offset = offset
			c.lock.Unlock()

			if c.offset == maxClockOffset {
				logger.Error("error on getting clock offset")
			} else if c.OutOfSync(offset) {
				logger.Error(
					"The node is out of sync with the network time",
					"threshold", c.threshold,
					"offset", offset,
					"threshold(secs)", c.threshold.Seconds(),
					"offset(secs)", offset.Seconds(),
				)
			}
		}
	}
}

func (c *Checker) Stop() {
	c.cancel()
	c.ticker.Stop()
}

func (c *Checker) OutOfSync(offset time.Duration) bool {
	return math.Abs(float64(offset)) > float64(c.threshold)
}

func (c *Checker) GetClockOffset() (time.Duration, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	offset := c.offset

	if offset == maxClockOffset {
		return 0, errors.Errorf(errors.ErrNtpError, "unable to get clock offset")
	}

	return offset, nil
}

func (*Checker) clockOffset() time.Duration {
	for _, server := range _pools {
		response, err := ntp.Query(server)
		if err != nil {
			logger.Warn("ntp error", "server", server, "error", err)

			continue
		}

		if err := response.Validate(); err != nil {
			logger.Warn("ntp validate error", "server", server, "error", err)

			continue
		}

		return response.ClockOffset
	}

	logger.Error("failed to get ntp query from all pool, set default max clock offset")

	return maxClockOffset
}
