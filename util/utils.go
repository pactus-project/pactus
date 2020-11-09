package util

import (
	"math/rand"
	"time"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func RandInt(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func RandInt64(max int64) int64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int63n(max)
}
