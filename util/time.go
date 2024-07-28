package util

import (
	"time"
)

// RoundNow returns the result of rounding sec to the current time in UTC.
// The rounding behavior is rounding down.
func RoundNow(sec int) time.Time {
	return roundDownTime(time.Now(), sec)
}

func roundDownTime(t time.Time, sec int) time.Time {
	t = t.Truncate(time.Duration(sec) * time.Second)
	t = t.UTC()

	return t
}
