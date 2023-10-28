package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	c1 := time.Now()
	c2 := Now()
	c3 := RoundNow(5)

	assert.NotEqual(t, c1, c2)
	assert.Equal(t, c1.Second(), c2.Second())
	assert.Equal(t, c3.Nanosecond(), 0)
	assert.Equal(t, c3.Second()%5, 0)
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

	assert.Equal(t, c1.Nanosecond(), 0)
	assert.Equal(t, c2.Nanosecond(), 0)
	assert.Equal(t, c3.Nanosecond(), 0)
	assert.Equal(t, c4.Nanosecond(), 0)
	assert.Equal(t, c5.Nanosecond(), 0)

	assert.Equal(t, c1.Second(), 10)
	assert.Equal(t, c2.Second(), 20)
	assert.Equal(t, c3.Second(), 30)
	assert.Equal(t, c4.Second(), 40)
	assert.Equal(t, c5.Second(), 50)
}
