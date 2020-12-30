package util

import (
	"time"
)

func Now() time.Time {
	return RoundTime(time.Now(), 0)
}

func RoundNow(sec int) time.Time {
	return RoundTime(time.Now(), sec)
}

func RoundTime(t time.Time, sec int) time.Time {
	return t.Round(time.Duration(sec) * time.Second).UTC()
}
