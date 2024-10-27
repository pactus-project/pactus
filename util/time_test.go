package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoundNow(t *testing.T) {
	time1 := time.Now()
	time2 := RoundNow(1)
	time3 := RoundNow(5)

	assert.NotEqual(t, time1, time2)
	assert.Equal(t, time1.Second(), time2.Second())
	assert.Equal(t, int64(0), time2.UnixMicro()%1000000)
	assert.Equal(t, int64(0), time2.UnixMilli()%1000)
	assert.Equal(t, 0, time3.Nanosecond())
	assert.Equal(t, 0, time3.Second()%5)
}

func TestRoundingTime(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"2006-01-02T15:04:11.111111111Z", 10},
		{"2006-01-02T15:04:24.333333333Z", 20},
		{"2006-01-02T15:04:35.555555555Z", 30},
		{"2006-01-02T15:04:48.777777777Z", 40},
		{"2006-01-02T15:04:59.999999999Z", 50},
	}

	for _, tt := range tests {
		parsedTime, err := time.Parse(time.RFC3339Nano, tt.input)
		assert.NoError(t, err, "Failed to parse time")

		roundedTime := roundDownTime(parsedTime, 10)

		assert.Equal(t, 0, roundedTime.Nanosecond(), "Nanoseconds should be rounded to 0")
		assert.Equal(t, tt.expected, roundedTime.Second(), "Seconds should match the expected rounded value")
	}
}
