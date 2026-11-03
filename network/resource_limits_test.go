package network

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceLimitsAreBounded(t *testing.T) {
	conf := DefaultConfig()
	conf.MaxConns = 1000

	limits := resourceLimits(conf)
	system := limits.ToPartialLimitConfig().System

	assert.Less(t, int(system.ConnsInbound), math.MaxInt, "inbound connections should be capped")
	assert.GreaterOrEqual(t, int(system.ConnsInbound), conf.MaxConns+conf.MaxConns/4,
		"the cap should not be below the connection manager high watermark")
	assert.Less(t, int(system.Conns), math.MaxInt)
}
