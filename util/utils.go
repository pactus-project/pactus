package util

import (
	"math/rand"
	"time"
)

const MaxUint64 = ^uint64(0)
const MinUint64 = 0
const MaxInt64 = int64(MaxUint64 >> 1)
const MinInt64 = -MaxInt64 - 1

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func Max64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func RandInt(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func RandInt64(max int64) int64 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int63n(max)
}
