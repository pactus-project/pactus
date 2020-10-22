package util

import (
	"time"
)

func Now() time.Time {
	return Canonical(time.Now())
}

func Canonical(t time.Time) time.Time {
	return t.Round(0).UTC()
}
