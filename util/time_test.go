package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoundNow(t *testing.T) {
	c1 := time.Now()
	c2 := RoundNow(1)
	c3 := RoundNow(5)

	assert.NotEqual(t, c1, c2)
	assert.Equal(t, c1.Second(), c2.Second())
	assert.Equal(t, int64(0), c2.UnixMicro()%1000000)
	assert.Equal(t, int64(0), c2.UnixMilli()%1000)
	assert.Equal(t, 0, c3.Nanosecond())
	assert.Equal(t, 0, c3.Second()%5)
}

func TestRoundingTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:11.111111111Z")
	t2, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:24.333333333Z")
	t3, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:35.555555555Z")
	t4, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:48.777777777Z")
	t5, _ := time.Parse(time.RFC3339Nano, "2006-01-02T15:04:59.999999999Z")
	c1 := roundDownTime(t1, 10)
	c2 := roundDownTime(t2, 10)
	c3 := roundDownTime(t3, 10)
	c4 := roundDownTime(t4, 10)
	c5 := roundDownTime(t5, 10)

	assert.Equal(t, 0, c1.Nanosecond())
	assert.Equal(t, 0, c2.Nanosecond())
	assert.Equal(t, 0, c3.Nanosecond())
	assert.Equal(t, 0, c4.Nanosecond())
	assert.Equal(t, 0, c5.Nanosecond())

	assert.Equal(t, 10, c1.Second())
	assert.Equal(t, 20, c2.Second())
	assert.Equal(t, 30, c3.Second())
	assert.Equal(t, 40, c4.Second())
	assert.Equal(t, 50, c5.Second())
}
